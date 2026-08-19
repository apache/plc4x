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

package firmata

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// newTestCodec spins a MessageCodec up on a connected (but not running) test transport. The
// transport instance is what the tests feed messages into and read the codec's output back out of.
func newTestCodec(t *testing.T) (*MessageCodec, *test.TransportInstance) {
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

	return NewMessageCodec(testTransportInstance, _options...), testTransportInstance
}

func TestFirmataFrameSize(t *testing.T) {
	tests := []struct {
		name    string
		buffer  []byte
		want    int
		wantErr error
	}{
		{name: "an analog-IO message is three bytes", buffer: []byte{0xE1, 0x05, 0x02}, want: 3},
		{name: "a digital-IO message is three bytes", buffer: []byte{0x92, 0x05, 0x00}, want: 3},
		{name: "a subscribe-analog message is two bytes", buffer: []byte{0xC1, 0x01}, want: 2},
		{name: "a subscribe-digital message is two bytes", buffer: []byte{0xD1, 0x01}, want: 2},
		{name: "a set-pin-mode command is three bytes", buffer: []byte{0xF4, 0x04, 0x01}, want: 3},
		{name: "a set-digital-pin-value command is three bytes", buffer: []byte{0xF5, 0x04, 0x01}, want: 3},
		{name: "a protocol-version command is three bytes", buffer: []byte{0xF9, 0x02, 0x05}, want: 3},
		{name: "a system reset is a single byte", buffer: []byte{0xFF}, want: 1},
		{name: "a sysex ends at its terminator", buffer: []byte{0xF0, 0x79, 0x02, 0x05, 0xF7}, want: 5},
		{name: "a sysex ends at the first terminator", buffer: []byte{0xF0, 0x6D, 0x02, 0xF7, 0xFF}, want: 4},
		{name: "an empty buffer sizes nothing", buffer: nil, wantErr: errNotEnoughData},
		{name: "a sysex without its terminator isn't complete", buffer: []byte{0xF0, 0x79, 0x02}, wantErr: errNotEnoughData},
		{name: "a payload byte cannot start a message", buffer: []byte{0x05}, wantErr: errOutOfSync},
		{name: "an undefined message type cannot start a message", buffer: []byte{0x80, 0x00}, wantErr: errOutOfSync},
		{name: "an undefined system sub-command cannot start a message", buffer: []byte{0xF1, 0x00}, wantErr: errOutOfSync},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			size, err := firmataFrameSize(testCase.buffer)
			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, size)
		})
	}
}

// receive frames one message off the transport, insisting that one really comes out.
func receive(t *testing.T, codec *MessageCodec) spi.Message {
	t.Helper()
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message, "expected a message")
	return message
}

func TestMessageCodec_ReceiveFramesTheStream(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// Three messages of three different lengths, back to back in a single chunk, followed by the
	// sysex the board answers a system reset with.
	transportInstance.FillReadBuffer([]byte{
		0xE1, 0x05, 0x02, // analog IO, pin 1
		0x92, 0x03, 0x01, // digital IO, port 2
		0xFF,                                     // system reset
		0xF0, 0x79, 0x02, 0x05, 0x41, 0x00, 0xF7, // report firmware "A" 2.5
	})

	analogIO, ok := receive(t, codec).(readWriteModel.FirmataMessageAnalogIO)
	require.True(t, ok)
	assert.Equal(t, uint8(1), analogIO.GetPin())
	value, err := decodeAnalogValue(analogIO.GetData())
	require.NoError(t, err)
	assert.Equal(t, int16(0x105), value)

	digitalIO, ok := receive(t, codec).(readWriteModel.FirmataMessageDigitalIO)
	require.True(t, ok)
	assert.Equal(t, uint8(2), digitalIO.GetPinBlock())

	systemReset, ok := receive(t, codec).(readWriteModel.FirmataMessageCommand)
	require.True(t, ok)
	assert.IsType(t, readWriteModel.NewFirmataCommandSystemReset(), systemReset.GetCommand())

	report, ok := firmwareReport(receive(t, codec))
	require.True(t, ok)
	assert.Equal(t, uint8(2), report.GetMajorVersion())
	assert.Equal(t, uint8(5), report.GetMinorVersion())
	assert.Equal(t, "A", string(report.GetFileName()))

	// Every one of the four messages was consumed exactly as far as it reached.
	remaining, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.Zero(t, remaining)
}

// A message which arrives in pieces must be waited for rather than half consumed.
func TestMessageCodec_ReceiveWaitsForAFragmentedMessage(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	transportInstance.FillReadBuffer([]byte{0xE1, 0x05})
	message, err := codec.Receive(shortLivedContext(t))
	require.NoError(t, err)
	assert.Nil(t, message, "a fragment must not be consumed")

	transportInstance.FillReadBuffer([]byte{0x02})
	analogIO, ok := receive(t, codec).(readWriteModel.FirmataMessageAnalogIO)
	require.True(t, ok)
	assert.Equal(t, uint8(1), analogIO.GetPin())
}

// A sysex has to be waited for until its terminator has arrived, however long that takes.
func TestMessageCodec_ReceiveWaitsForASysexTerminator(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	transportInstance.FillReadBuffer([]byte{0xF0, 0x79, 0x02, 0x05, 0x41, 0x00})
	message, err := codec.Receive(shortLivedContext(t))
	require.NoError(t, err)
	assert.Nil(t, message, "an unterminated sysex must not be consumed")

	transportInstance.FillReadBuffer([]byte{0xF7})
	report, ok := firmwareReport(receive(t, codec))
	require.True(t, ok)
	assert.Equal(t, "A", string(report.GetFileName()))
}

// Joining a serial line in the middle of a message leaves payload bytes at the head of the stream.
// They can't start a message - payload bytes always have their most significant bit clear - so they
// are thrown away until the stream lines up again. plc4j throws here, which on a shared UART means
// the connection never recovers.
func TestMessageCodec_ReceiveResynchronizes(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	transportInstance.FillReadBuffer([]byte{
		0x05, 0x02, // the tail of a message we joined halfway through
		0x80, // an undefined message type
		0xF1, // an undefined system sub-command
		0xE1, 0x05, 0x02,
	})

	analogIO, ok := receive(t, codec).(readWriteModel.FirmataMessageAnalogIO)
	require.True(t, ok)
	assert.Equal(t, uint8(1), analogIO.GetPin())
	assert.Zero(t, codec.resyncSkippedBytes, "the resync counter is cleared once a message parses")
}

func TestMessageCodec_ReceiveRefusesAnUnconnectedTransport(t *testing.T) {
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	transport := test.NewTransport(_options...)
	transportInstance, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
	require.NoError(t, err)
	codec := NewMessageCodec(transportInstance, _options...)

	message, err := codec.Receive(testutils.TestContext(t))
	assert.Error(t, err)
	assert.Nil(t, message)
}

// shortLivedContext bounds a receive which is expected to find nothing. Waiting for bytes which
// aren't coming is exactly what a transport is supposed to do, so a receive that must come back
// empty needs a deadline of its own.
func shortLivedContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 100*time.Millisecond)
	t.Cleanup(cancel)
	return ctx
}

func TestMessageCodec_Send(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	require.NoError(t, codec.Send(testutils.TestContext(t), "test",
		readWriteModel.NewFirmataMessageSubscribeDigitalPinValue(2, true)))
	require.NoError(t, codec.Send(testutils.TestContext(t), "test",
		readWriteModel.NewFirmataMessageCommand(readWriteModel.NewFirmataCommandSystemReset())))

	assert.Equal(t, []byte{0xD2, 0x01, 0xFF}, transportInstance.DrainWriteBuffer(3))
}

// Sending something that isn't a firmata message is a programming error somewhere up the stack, but
// it must surface as an error rather than take the process down with a failed type assertion.
func TestMessageCodec_SendRejectsAForeignMessage(t *testing.T) {
	codec, _ := newTestCodec(t)

	err := codec.Send(testutils.TestContext(t), "test", notAFirmataMessage{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "FirmataMessage")
}

type notAFirmataMessage struct{}

func (notAFirmataMessage) String() string { return "notAFirmataMessage" }
func (notAFirmataMessage) Serialize() ([]byte, error) {
	return nil, nil
}
func (notAFirmataMessage) SerializeWithWriteBuffer(context.Context, utils.WriteBuffer) error {
	return nil
}
func (notAFirmataMessage) GetLengthInBytes(context.Context) uint64 { return 0 }
func (notAFirmataMessage) GetLengthInBits(context.Context) uint64  { return 0 }
