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

// newTestRtuCodec spins a MessageCodecRtu up on a connected (but not running) test-transport. The
// transport instance is what the tests feed frames into and read the codec's output back out of.
func newTestRtuCodec(t *testing.T) (*MessageCodecRtu, *test.TransportInstance) {
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

	return NewMessageCodecRtu(testTransportInstance, _options...), testTransportInstance
}

// rtuFrame serializes an RTU ADU the way it goes onto the wire, CRC included.
func rtuFrame(t *testing.T, address uint8, pdu readWriteModel.ModbusPDU) []byte {
	t.Helper()
	theBytes, err := readWriteModel.NewModbusRtuADU(address, pdu).Serialize()
	require.NoError(t, err)
	return theBytes
}

func TestRtuResponseSize(t *testing.T) {
	tests := []struct {
		name   string
		header []byte
		want   int
	}{
		{name: "read coils carries its byte count", header: []byte{0x01, 0x01, 0x02}, want: 7},
		{name: "read discrete inputs carries its byte count", header: []byte{0x01, 0x02, 0x01}, want: 6},
		{name: "read holding registers carries its byte count", header: []byte{0x01, 0x03, 0x14}, want: 25},
		{name: "read input registers carries its byte count", header: []byte{0x01, 0x04, 0x02}, want: 7},
		{name: "read file record carries its byte count", header: []byte{0x01, 0x14, 0x06}, want: 11},
		{name: "write file record carries its byte count", header: []byte{0x01, 0x15, 0x09}, want: 14},
		{name: "read write multiple registers carries its byte count", header: []byte{0x01, 0x17, 0x04}, want: 9},
		{name: "write single coil echoes a fixed size", header: []byte{0x01, 0x05, 0xFF}, want: 8},
		{name: "write single register echoes a fixed size", header: []byte{0x01, 0x06, 0x00}, want: 8},
		{name: "write multiple coils echoes a fixed size", header: []byte{0x01, 0x0F, 0x00}, want: 8},
		{name: "write multiple registers echoes a fixed size", header: []byte{0x01, 0x10, 0x00}, want: 8},
		{name: "an exception is always five bytes", header: []byte{0x01, 0x83, 0x02}, want: 5},
		{name: "the exception of an unknown function is still five bytes", header: []byte{0x01, 0xFF, 0x00}, want: 5},
		{name: "an unknown function code has no size", header: []byte{0x01, 0x2B, 0x00}, want: -1},
		{name: "a function code of zero has no size", header: []byte{0x01, 0x00, 0x00}, want: -1},
		{name: "a short header has no size", header: []byte{0x01, 0x03}, want: -1},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, rtuResponseSize(testCase.header))
		})
	}
}

// Sending something that isn't an RTU ADU is a programming error somewhere up the stack, but it
// must surface as an error rather than take the process down with a failed type assertion. In
// particular a Modbus TCP ADU must not slip through - that is exactly the bug the RTU driver had
// while it was still wired to the TCP codec.
func TestMessageCodecRtu_SendRejectsAForeignMessage(t *testing.T) {
	tests := []struct {
		name    string
		message spi.Message
	}{
		{name: "not an ADU at all", message: notAnAdu{}},
		{name: "a Modbus TCP ADU", message: readWriteModel.NewModbusTcpADU(1, 1,
			readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 1))},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			codec, _ := newTestRtuCodec(t)

			err := codec.Send(testutils.TestContext(t), "test", testCase.message)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ModbusRtuADU")
		})
	}
}

// What goes out on a serial line is address, PDU and CRC - no MBAP header.
func TestMessageCodecRtu_SendWritesAFramedRequest(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	request := readWriteModel.NewModbusRtuADU(17, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 10))

	require.NoError(t, codec.Send(testutils.TestContext(t), "test", request))

	written := transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
	assert.Equal(t, rtuFrame(t, 17, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 10)), written)
	// The CRC is the last two bytes and is what tells the frame apart from an MBAP-framed one.
	require.Len(t, written, 8)
	assert.Equal(t, uint8(17), written[0])
}

func TestMessageCodecRtu_ReceiveFramesAResponse(t *testing.T) {
	tests := []struct {
		name    string
		address uint8
		pdu     readWriteModel.ModbusPDU
	}{
		{
			name:    "a read response is sized by its byte count",
			address: 1,
			pdu:     readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A, 0x00, 0x2B}),
		},
		{
			name:    "a single-byte read response is the smallest read frame",
			address: 42,
			pdu:     readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01}),
		},
		{
			name:    "a write echo has a fixed size",
			address: 1,
			pdu:     readWriteModel.NewModbusPDUWriteSingleRegisterResponse(7, 0x1234),
		},
		{
			name:    "an exception is five bytes",
			address: 3,
			pdu:     readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			codec, transportInstance := newTestRtuCodec(t)
			frame := rtuFrame(t, testCase.address, testCase.pdu)
			transportInstance.FillReadBuffer(frame)

			message, err := codec.Receive(testutils.TestContext(t))

			require.NoError(t, err)
			require.NotNil(t, message)
			adu, ok := message.(readWriteModel.ModbusRtuADU)
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

// Two frames back to back are two messages, not one - the second one must still be there after the
// first was consumed.
func TestMessageCodecRtu_ReceiveFramesBackToBackResponses(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	first := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})
	second := readWriteModel.NewModbusPDUWriteSingleRegisterResponse(7, 0x1234)
	transportInstance.FillReadBuffer(append(rtuFrame(t, 1, first), rtuFrame(t, 2, second)...))

	firstMessage, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, firstMessage)
	assert.Equal(t, first, firstMessage.(readWriteModel.ModbusRtuADU).GetPdu())

	secondMessage, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, secondMessage)
	assert.Equal(t, uint8(2), secondMessage.(readWriteModel.ModbusRtuADU).GetAddress())
	assert.Equal(t, second, secondMessage.(readWriteModel.ModbusRtuADU).GetPdu())
}

// A frame that hasn't fully arrived yet must be waited for, not consumed - a serial line hands the
// bytes over as they trickle in.
func TestMessageCodecRtu_ReceiveWaitsForTheRestOfAFrame(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A, 0x00, 0x2B})
	frame := rtuFrame(t, 1, pdu)
	require.Len(t, frame, 9)
	transportInstance.FillReadBuffer(frame[:6])

	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	assert.Nil(t, message)
	available, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.EqualValues(t, 6, available, "a partial frame must not be consumed")

	transportInstance.FillReadBuffer(frame[6:])
	// The test transport only moves bytes from its channel into its buffer once somebody asks for
	// more than it holds; the real ones fill themselves in GetNumBytesAvailableInBuffer.
	_, err = transportInstance.PeekReadableBytes(testutils.TestContext(t), uint32(len(frame)))
	require.NoError(t, err)

	message, err = codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, pdu, message.(readWriteModel.ModbusRtuADU).GetPdu())
}

// Less than the smallest frame can't be anything yet either.
func TestMessageCodecRtu_ReceiveWaitsForAMinimalFrame(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	transportInstance.FillReadBuffer([]byte{0x01, 0x03, 0x02})

	// The buffer never fills up to the minimum, so the fill has to be bounded by the context.
	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 50*time.Millisecond)
	defer cancel()
	message, err := codec.Receive(ctx)

	require.NoError(t, err)
	assert.Nil(t, message)
	available, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.EqualValues(t, 3, available)
}

// A frame whose CRC doesn't add up never happened as far as the codec is concerned: the generated
// parser rejects it and the codec resynchronizes rather than handing a corrupted response up.
func TestMessageCodecRtu_ReceiveRejectsABrokenCrc(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	frame := rtuFrame(t, 1, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))
	frame[len(frame)-1] ^= 0xFF // one flipped bit in the CRC is all it takes
	transportInstance.FillReadBuffer(frame)

	message, err := codec.Receive(testutils.TestContext(t))

	require.NoError(t, err)
	assert.Nil(t, message, "a frame with a bad CRC must not be handed up")
	assert.NotZero(t, codec.resyncSkippedBytes, "the codec should be resynchronizing")
}

// The same, but with a good frame behind the broken one: the response that did survive the line
// has to come through once the codec has resynchronized onto it.
func TestMessageCodecRtu_ReceiveResynchronizesOntoTheNextGoodFrame(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	broken := rtuFrame(t, 1, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))
	broken[len(broken)-1] ^= 0xFF
	good := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2B})
	transportInstance.FillReadBuffer(append(broken, rtuFrame(t, 1, good)...))

	message := receiveUntilAMessageArrives(t, codec)

	require.NotNil(t, message)
	assert.Equal(t, good, message.(readWriteModel.ModbusRtuADU).GetPdu())
	assert.Zero(t, codec.resyncSkippedBytes, "the resync counter resets once the stream is back in sync")
}

// Garbage in front of a frame - the tail of a frame that started before we started listening, or
// line noise - is skipped a byte at a time until a frame parses.
func TestMessageCodecRtu_ReceiveSkipsLeadingGarbage(t *testing.T) {
	tests := []struct {
		name    string
		garbage []byte
	}{
		{name: "unknown function codes", garbage: []byte{0x11, 0x2B, 0x63, 0x00}},
		{name: "a write echo that never made it intact", garbage: []byte{0x01, 0x05, 0x00, 0x01, 0xFF, 0x00, 0x11, 0x22}},
		{name: "zeroes", garbage: []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			codec, transportInstance := newTestRtuCodec(t)
			pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})
			transportInstance.FillReadBuffer(append(testCase.garbage, rtuFrame(t, 9, pdu)...))

			message := receiveUntilAMessageArrives(t, codec)

			require.NotNil(t, message)
			adu, ok := message.(readWriteModel.ModbusRtuADU)
			require.True(t, ok, "got %T", message)
			assert.Equal(t, uint8(9), adu.GetAddress())
			assert.Equal(t, pdu, adu.GetPdu())
		})
	}
}

// A candidate at a clean stream position that claims more bytes than have arrived is waited for
// rather than skipped - that is how a genuinely fragmented frame is recognized. The price, which
// plc4j's ModbusRtuMessageCodec pays too, is that garbage which happens to look like a large read
// response stalls the codec until enough bytes to reject it have come in.
func TestMessageCodecRtu_ReceiveWaitsOutAnOversizedCandidateAtACleanPosition(t *testing.T) {
	codec, transportInstance := newTestRtuCodec(t)
	// address 1, function code 3 (read holding registers), byte count 250 - a 255 byte frame.
	transportInstance.FillReadBuffer([]byte{0x01, 0x03, 0xFA, 0x00})

	message, err := codec.Receive(testutils.TestContext(t))

	require.NoError(t, err)
	assert.Nil(t, message)
	assert.Zero(t, codec.resyncSkippedBytes, "nothing was skipped - the codec is waiting for the rest")
	available, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.EqualValues(t, 4, available)
}

// receiveUntilAMessageArrives drives the codec until it produces a message. Resynchronization can
// need more than one call whenever a skipped byte leaves the buffer too short to decide anything.
func receiveUntilAMessageArrives(t *testing.T, codec *MessageCodecRtu) spi.Message {
	t.Helper()
	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 2*time.Second)
	defer cancel()
	for range 32 {
		message, err := codec.Receive(ctx)
		require.NoError(t, err)
		if message != nil {
			return message
		}
	}
	t.Fatal("codec never produced a message")
	return nil
}
