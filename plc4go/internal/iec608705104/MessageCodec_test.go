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

package iec608705104

import (
	"context"
	"encoding/hex"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newTestCodec spins a MessageCodec up on a connected (but not running) test transport. The
// transport instance is what the tests feed frames into and read the codec's output back out of.
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

// receive frames one APDU off the transport, insisting that one really comes out.
func receive(t *testing.T, codec *MessageCodec) spi.Message {
	t.Helper()
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message, "expected an APDU")
	return message
}

// shortLivedContext bounds a receive which is expected to find nothing. Waiting for octets which
// aren't coming is exactly what a transport is supposed to do, so a receive that must come back
// empty needs a deadline of its own.
func shortLivedContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 100*time.Millisecond)
	t.Cleanup(cancel)
	return ctx
}

// mustDecodeHex is one of the frames the parser-serializer testsuite captured off a real station.
func mustDecodeHex(t *testing.T, frame string) []byte {
	t.Helper()
	theBytes, err := hex.DecodeString(frame)
	require.NoError(t, err)
	return theBytes
}

func TestApduFrameSize(t *testing.T) {
	tests := []struct {
		name    string
		buffer  []byte
		want    uint32
		wantErr error
	}{
		{name: "the shortest APDU there is", buffer: []byte{0x68, 0x04}, want: 6},
		{name: "a single point information", buffer: []byte{0x68, 0x0E}, want: 16},
		{name: "the longest APDU there is", buffer: []byte{0x68, 0xFD}, want: 255},
		{name: "an empty buffer sizes nothing", buffer: nil, wantErr: errNotEnoughData},
		{name: "the start byte alone sizes nothing", buffer: []byte{0x68}, wantErr: errNotEnoughData},
		{name: "a wrong start byte cannot start an APDU", buffer: []byte{0x69, 0x04}, wantErr: errOutOfSync},
		{name: "a wrong start byte is refused before the length arrives", buffer: []byte{0x00}, wantErr: errOutOfSync},
		{name: "a length below the control field cannot be", buffer: []byte{0x68, 0x03}, wantErr: errOutOfSync},
		{name: "a zero length cannot be", buffer: []byte{0x68, 0x00}, wantErr: errOutOfSync},
		{name: "a length beyond what one octet holds cannot be", buffer: []byte{0x68, 0xFE}, wantErr: errOutOfSync},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			size, err := apduFrameSize(testCase.buffer)
			if testCase.wantErr != nil {
				assert.ErrorIs(t, err, testCase.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, size)
		})
	}
}

// The whole handshake plus a data frame, back to back in one chunk - which is exactly how a station
// sends it (the testsuite's own frames arrive concatenated for the same reason).
func TestMessageCodec_ReceiveFramesTheStream(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	var stream []byte
	stream = append(stream, mustDecodeHex(t, "680483000000")...) // TESTFR con
	stream = append(stream, mustDecodeHex(t, "68040b000000")...) // STARTDT con
	stream = append(stream, mustDecodeHex(t, "680401001a00")...) // S-format, N(R) = 13
	stream = append(stream, mustDecodeHex(t, "681a04000200010414000a0001000000020000000300000004000000")...)
	transportInstance.FillReadBuffer(stream)

	_, ok := receive(t, codec).(readWriteModel.APDUUFormatTestFrameConfirmation)
	assert.True(t, ok, "the test frame confirmation")
	_, ok = receive(t, codec).(readWriteModel.APDUUFormatStartDataTransferConfirmation)
	assert.True(t, ok, "the start-data-transfer confirmation")

	supervisory, ok := receive(t, codec).(readWriteModel.APDUSFormat)
	require.True(t, ok)
	assert.Equal(t, uint16(13), supervisory.GetReceiveSequenceNo()>>1)

	dataFrame, ok := receive(t, codec).(readWriteModel.APDUIFormat)
	require.True(t, ok)
	assert.Equal(t, uint16(2), dataFrame.GetCommand()>>1, "the send sequence number")
	assert.Equal(t, uint16(1), dataFrame.GetReceiveSequenceNo()>>1)
	assert.Equal(t, uint16(10), dataFrame.GetAsdu().GetAsduAddressField())
	assert.Len(t, dataFrame.GetAsdu().GetInformationObjects(), 4)

	// Every one of the four frames was consumed exactly as far as it reached.
	remaining, err := transportInstance.GetNumBytesAvailableInBuffer()
	require.NoError(t, err)
	assert.Zero(t, remaining)
}

// A frame which arrives in pieces must be waited for rather than half consumed.
func TestMessageCodec_ReceiveWaitsForAFragmentedFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)
	frame := mustDecodeHex(t, "681a04000200010414000a0001000000020000000300000004000000")

	// The start octet on its own says nothing about how long the frame is.
	transportInstance.FillReadBuffer(frame[:1])
	message, err := codec.Receive(shortLivedContext(t))
	require.NoError(t, err)
	assert.Nil(t, message, "a fragment must not be consumed")

	// The length octet arrives, but not the body behind it.
	transportInstance.FillReadBuffer(frame[1:10])
	message, err = codec.Receive(shortLivedContext(t))
	require.NoError(t, err)
	assert.Nil(t, message, "a fragment must not be consumed")

	transportInstance.FillReadBuffer(frame[10:])
	dataFrame, ok := receive(t, codec).(readWriteModel.APDUIFormat)
	require.True(t, ok)
	assert.Len(t, dataFrame.GetAsdu().GetInformationObjects(), 4)
}

// Joining the stream halfway through a frame leaves octets at the head which cannot start one. They
// are thrown away until the stream lines up again - plc4j's codec throws instead, which drops the
// TCP connection and every subscription on it over a single corrupted octet.
func TestMessageCodec_ReceiveResynchronizes(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	var stream []byte
	stream = append(stream, 0x00, 0x0A, 0x68, 0x03) // tail of a frame plus an impossible length
	stream = append(stream, mustDecodeHex(t, "680483000000")...)
	transportInstance.FillReadBuffer(stream)

	_, ok := receive(t, codec).(readWriteModel.APDUUFormatTestFrameConfirmation)
	assert.True(t, ok, "the frame behind the garbage")
	assert.Zero(t, codec.resyncSkippedBytes, "the resync counter is cleared once a frame parses")
}

// A frame whose length octet is right but whose body doesn't parse is dropped whole: the length said
// where the next frame starts, so the frames behind it need not pay for it.
func TestMessageCodec_ReceiveDropsAnUnparseableFrameWhole(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// An I-format frame claiming a type identification the model has no case for at all.
	broken := []byte{0x68, 0x08, 0x02, 0x00, 0x00, 0x00, 0x2A, 0x01, 0x03, 0x00}
	var stream []byte
	stream = append(stream, broken...)
	stream = append(stream, mustDecodeHex(t, "680483000000")...)
	transportInstance.FillReadBuffer(stream)

	message, err := codec.Receive(shortLivedContext(t))
	require.NoError(t, err)
	assert.Nil(t, message, "the broken frame yields nothing")

	_, ok := receive(t, codec).(readWriteModel.APDUUFormatTestFrameConfirmation)
	assert.True(t, ok, "the frame behind the broken one")
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

func TestMessageCodec_Send(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	require.NoError(t, codec.Send(testutils.TestContext(t), "test",
		readWriteModel.NewAPDUUFormatStartDataTransferActivation(commandStartDataTransferActivation)))

	written := transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
	assert.Equal(t, mustDecodeHex(t, "680407000000"), written)
}

// An S-format acknowledgement carries the sequence number shifted left by one, because the low bit
// of the control field octet pair is the format discriminator.
func TestMessageCodec_SendAcknowledgement(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	require.NoError(t, codec.Send(testutils.TestContext(t), "test",
		readWriteModel.NewAPDUSFormat(commandSupervisoryFormat, 13<<1)))

	written := transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
	assert.Equal(t, mustDecodeHex(t, "680401001a00"), written)
}

func TestMessageCodec_SendRefusesSomethingElse(t *testing.T) {
	codec, _ := newTestCodec(t)
	assert.Error(t, codec.Send(testutils.TestContext(t), "test", nil))
}
