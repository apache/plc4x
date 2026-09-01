/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package modbus

import (
	"bytes"
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newTestAsciiCodec spins a MessageCodecAscii up on a connected (but not running) test-transport.
// The transport instance is what the tests feed frames into and read the codec's output back out
// of.
func newTestAsciiCodec(t *testing.T) (*MessageCodecAscii, *test.TransportInstance) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)

	transport := test.NewTransport(_options...)
	transportInstance, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
	require.NoError(t, err)
	testTransportInstance, ok := transportInstance.(*test.TransportInstance)
	require.True(t, ok)
	require.NoError(t, testTransportInstance.Connect(t.Context()))
	t.Cleanup(func() {
		assert.NoError(t, testTransportInstance.Close())
	})

	return NewMessageCodecAscii(testTransportInstance, _options...), testTransportInstance
}

// asciiFrame builds the frame an ADU goes on the wire as: ':' + hex(address, PDU, LRC) + CR + LF.
func asciiFrame(t *testing.T, address uint8, pdu readWriteModel.ModbusPDU) []byte {
	t.Helper()
	theBytes, err := readWriteModel.NewModbusAsciiADU(address, pdu).Serialize()
	require.NoError(t, err)
	return encodeAsciiFrame(theBytes)
}

// The hex layer is the whole difference to RTU, so it is worth pinning down on its own. The
// example is the one the ASCII parser-serializer testsuite uses: address 1, read holding registers
// from 0, quantity 10, whose bytes sum to 0x0E and therefore carry an LRC of 0xF2.
func TestEncodeAsciiFrame(t *testing.T) {
	tests := []struct {
		name        string
		binaryBytes []byte
		want        string
	}{
		{name: "an empty payload is still a frame", binaryBytes: nil, want: ":\r\n"},
		{name: "nibbles are spelled out in upper case hex", binaryBytes: []byte{0xAB, 0xCD, 0xEF}, want: ":ABCDEF\r\n"},
		{name: "a leading zero nibble is not dropped", binaryBytes: []byte{0x00, 0x0A}, want: ":000A\r\n"},
		{
			name:        "a read holding registers request",
			binaryBytes: []byte{0x01, 0x03, 0x00, 0x00, 0x00, 0x0A, 0xF2},
			want:        ":01030000000AF2\r\n",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, string(encodeAsciiFrame(testCase.binaryBytes)))
		})
	}
}

func TestAsciiFrameEnd(t *testing.T) {
	tests := []struct {
		name string
		data string
		want int
	}{
		{name: "a terminated frame ends at its CR", data: ":010304\r\n", want: 7},
		{name: "the terminator of the first frame wins", data: ":010304\r\n:020305\r\n", want: 7},
		{name: "an unterminated frame has no end", data: ":010304", want: -1},
		{name: "a lone CR is not a terminator", data: ":010304\r", want: -1},
		{name: "a lone LF is not a terminator", data: ":010304\n", want: -1},
		{name: "nothing at all has no end", data: "", want: -1},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, asciiFrameEnd([]byte(testCase.data)))
		})
	}
}

func TestDecodeAsciiPayload(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    []byte
		wantErr string
	}{
		{name: "upper case hex", payload: "0103AB", want: []byte{0x01, 0x03, 0xAB}},
		{name: "lower case hex is accepted too", payload: "0103ab", want: []byte{0x01, 0x03, 0xAB}},
		{name: "too short to be an ADU", payload: "0103", wantErr: "at least 6"},
		{name: "an odd number of characters", payload: "0103AB0", wantErr: "odd number"},
		{name: "a stray colon from a truncated frame", payload: "0103:1", wantErr: "not hex encoded"},
		{name: "something that isn't hex at all", payload: "zzzzzz", wantErr: "not hex encoded"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := decodeAsciiPayload([]byte(testCase.payload))
			if testCase.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, got)
		})
	}
}

// Sending something that isn't an ASCII ADU is a programming error somewhere up the stack, but it
// must surface as an error rather than take the process down with a failed type assertion. In
// particular a Modbus TCP ADU must not slip through - that is exactly the bug the ASCII driver had
// while it was still wired to the TCP codec.
func TestMessageCodecAscii_SendRejectsAForeignMessage(t *testing.T) {
	tests := []struct {
		name    string
		message spi.Message
	}{
		{name: "not an ADU at all", message: notAnAdu{}},
		{name: "a Modbus TCP ADU", message: readWriteModel.NewModbusTcpADU(1, 1,
			readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 1))},
		{name: "a Modbus RTU ADU", message: readWriteModel.NewModbusRtuADU(1,
			readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 1))},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			codec, _ := newTestAsciiCodec(t)

			err := codec.Send(testutils.TestContext(t), "test", testCase.message)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ModbusAsciiADU")
		})
	}
}

// What goes out is printable: a colon, hex characters, and a CR/LF. This is the exact frame the
// ASCII parser-serializer testsuite spells out in binary.
func TestMessageCodecAscii_SendWritesAFramedRequest(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	request := readWriteModel.NewModbusAsciiADU(1, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 10))

	require.NoError(t, codec.Send(testutils.TestContext(t), "test", request))

	written := transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
	assert.Equal(t, ":01030000000AF2\r\n", string(written))
}

func TestMessageCodecAscii_ReceiveFramesAResponse(t *testing.T) {
	tests := []struct {
		name    string
		address uint8
		pdu     readWriteModel.ModbusPDU
	}{
		{
			name:    "a read response",
			address: 1,
			pdu:     readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A, 0x00, 0x2B}),
		},
		{
			name:    "a single-byte read response is the smallest read frame",
			address: 42,
			pdu:     readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01}),
		},
		{
			name:    "a write echo",
			address: 1,
			pdu:     readWriteModel.NewModbusPDUWriteSingleRegisterResponse(7, 0x1234),
		},
		{
			name:    "an exception",
			address: 3,
			pdu:     readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			codec, transportInstance := newTestAsciiCodec(t)
			transportInstance.FillReadBuffer(asciiFrame(t, testCase.address, testCase.pdu))

			message, err := codec.Receive(testutils.TestContext(t))

			require.NoError(t, err)
			require.NotNil(t, message)
			adu, ok := message.(readWriteModel.ModbusAsciiADU)
			require.True(t, ok, "got %T", message)
			assert.Equal(t, testCase.address, adu.GetAddress())
			assert.Equal(t, testCase.pdu, adu.GetPdu())
			// The frame and nothing but the frame was consumed.
			available, err := transportInstance.GetNumBytesAvailableInBuffer()
			require.NoError(t, err)
			assert.Zero(t, available)
		})
	}
}

// Devices are specified to send upper case hex, but a reader has to cope with either case.
func TestMessageCodecAscii_ReceiveAcceptsLowerCaseHex(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0xAB})
	frame := asciiFrame(t, 1, pdu)
	transportInstance.FillReadBuffer(bytes.ToLower(frame))

	message, err := codec.Receive(testutils.TestContext(t))

	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, pdu, message.(readWriteModel.ModbusAsciiADU).GetPdu())
}

// Two frames back to back are two messages, not one - the second one must still be there after the
// first was consumed.
func TestMessageCodecAscii_ReceiveFramesBackToBackResponses(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	first := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})
	second := readWriteModel.NewModbusPDUWriteSingleRegisterResponse(7, 0x1234)
	transportInstance.FillReadBuffer(append(asciiFrame(t, 1, first), asciiFrame(t, 2, second)...))

	firstMessage, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, firstMessage)
	assert.Equal(t, first, firstMessage.(readWriteModel.ModbusAsciiADU).GetPdu())

	secondMessage, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, secondMessage)
	assert.Equal(t, uint8(2), secondMessage.(readWriteModel.ModbusAsciiADU).GetAddress())
	assert.Equal(t, second, secondMessage.(readWriteModel.ModbusAsciiADU).GetPdu())
}

// A frame whose terminator hasn't arrived yet must be waited for, not consumed - a serial line
// hands the bytes over as they trickle in.
func TestMessageCodecAscii_ReceiveWaitsForTheRestOfAFrame(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A, 0x00, 0x2B})
	frame := asciiFrame(t, 1, pdu)
	require.Len(t, frame, 19)
	transportInstance.FillReadBuffer(frame[:12])

	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	assert.Nil(t, message)
	available, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.EqualValues(t, 12, available, "a partial frame must not be consumed")

	transportInstance.FillReadBuffer(frame[12:])
	// The test transport only moves bytes from its channel into its buffer once somebody asks for
	// more than it holds; the real ones fill themselves in GetNumBytesAvailableInBuffer.
	_, err = transportInstance.PeekReadableBytes(testutils.TestContext(t), uint32(len(frame)))
	require.NoError(t, err)

	message, err = codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, pdu, message.(readWriteModel.ModbusAsciiADU).GetPdu())
}

// Less than the smallest frame can't be anything yet either.
func TestMessageCodecAscii_ReceiveWaitsForAMinimalFrame(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	transportInstance.FillReadBuffer([]byte(":0103"))

	// The buffer never fills up to the minimum, so the fill has to be bounded by the context.
	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 50*time.Millisecond)
	defer cancel()
	message, err := codec.Receive(ctx)

	require.NoError(t, err)
	assert.Nil(t, message)
	available, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.EqualValues(t, 5, available)
}

// A frame whose LRC doesn't add up never happened as far as the codec is concerned: the generated
// parser rejects it and the codec resynchronizes rather than handing a corrupted response up.
func TestMessageCodecAscii_ReceiveRejectsABrokenLrc(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	frame := asciiFrame(t, 1, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))
	// The last two hex characters before the CR/LF are the LRC.
	frame[len(frame)-3] = '0'
	frame[len(frame)-4] = '0'
	transportInstance.FillReadBuffer(frame)

	message, err := codec.Receive(testutils.TestContext(t))

	require.NoError(t, err)
	assert.Nil(t, message, "a frame with a bad LRC must not be handed up")
	assert.NotZero(t, codec.resyncSkippedBytes, "the codec should be resynchronizing")
}

// The same, but with a good frame behind the broken one: the response that did survive the line
// has to come through once the codec has resynchronized onto it.
func TestMessageCodecAscii_ReceiveResynchronizesOntoTheNextGoodFrame(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	broken := asciiFrame(t, 1, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))
	broken[len(broken)-3] = '0'
	broken[len(broken)-4] = '0'
	good := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2B})
	transportInstance.FillReadBuffer(append(broken, asciiFrame(t, 1, good)...))

	message := receiveUntilAnAsciiMessageArrives(t, codec)

	require.NotNil(t, message)
	assert.Equal(t, good, message.(readWriteModel.ModbusAsciiADU).GetPdu())
	assert.Zero(t, codec.resyncSkippedBytes, "the resync counter resets once the stream is back in sync")
}

// Garbage in front of a frame - the tail of a frame that started before we started listening, a
// frame that got cut off mid-flight, or line noise - is skipped a byte at a time until a frame
// parses.
func TestMessageCodecAscii_ReceiveSkipsLeadingGarbage(t *testing.T) {
	tests := []struct {
		name    string
		garbage []byte
	}{
		{name: "the tail of a frame we joined late", garbage: []byte("0304000000\r\n")},
		{name: "a frame that was cut off before its terminator", garbage: []byte(":010304AB")},
		{name: "characters that aren't hex", garbage: []byte(":hello world\r\n")},
		{name: "an odd number of hex characters", garbage: []byte(":010304A\r\n")},
		{name: "nothing but colons", garbage: []byte("::::::")},
		{name: "zeroes", garbage: []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			codec, transportInstance := newTestAsciiCodec(t)
			pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})
			transportInstance.FillReadBuffer(append(testCase.garbage, asciiFrame(t, 9, pdu)...))

			message := receiveUntilAnAsciiMessageArrives(t, codec)

			require.NotNil(t, message)
			adu, ok := message.(readWriteModel.ModbusAsciiADU)
			require.True(t, ok, "got %T", message)
			assert.EqualValues(t, 9, adu.GetAddress())
			assert.Equal(t, pdu, adu.GetPdu())
		})
	}
}

// A colon that never gets terminated would stall the codec forever if it were waited for
// unconditionally, so once more bytes than the longest possible frame have piled up behind it the
// colon is thrown away and the codec looks for the next one.
func TestMessageCodecAscii_ReceiveGivesUpOnAnOverlongUnterminatedFrame(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	// A colon followed by nothing but hex characters: a frame this long cannot exist.
	transportInstance.FillReadBuffer(append([]byte{asciiFrameStart}, bytes.Repeat([]byte("0"), int(modbusAsciiMaxSize))...))

	message, err := codec.Receive(testutils.TestContext(t))

	require.NoError(t, err)
	assert.Nil(t, message)
	assert.NotZero(t, codec.resyncSkippedBytes, "the colon should have been thrown away")
}

// The mirror image: an unterminated frame that could still turn into a real one is waited for
// rather than skipped, which is what makes a fragmented frame recoverable at all.
func TestMessageCodecAscii_ReceiveWaitsOutAnUnterminatedFrameThatCouldStillArrive(t *testing.T) {
	codec, transportInstance := newTestAsciiCodec(t)
	transportInstance.FillReadBuffer([]byte(":0103140000"))

	message, err := codec.Receive(testutils.TestContext(t))

	require.NoError(t, err)
	assert.Nil(t, message)
	assert.Zero(t, codec.resyncSkippedBytes, "nothing was skipped - the codec is waiting for the terminator")
	available, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.EqualValues(t, 11, available)
}

// receiveUntilAnAsciiMessageArrives drives the codec until it produces a message. Resynchronization
// can need more than one call whenever a skipped byte leaves the buffer too short to decide
// anything.
func receiveUntilAnAsciiMessageArrives(t *testing.T, codec *MessageCodecAscii) spi.Message {
	t.Helper()
	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 2*time.Second)
	defer cancel()
	for range 64 {
		message, err := codec.Receive(ctx)
		require.NoError(t, err)
		if message != nil {
			return message
		}
	}
	t.Fatal("codec never produced a message")
	return nil
}
