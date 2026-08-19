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

package abeth

import (
	"context"
	"encoding/hex"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newTestTransport builds the in-process test transport the codec and driver tests run on.
func newTestTransport(t *testing.T, _options ...options.WithOption) *test.Transport {
	t.Helper()
	return test.NewTransport(_options...)
}

// newTestCodec spins a MessageCodec up on a connected (but not running) test transport. The
// transport instance is what the tests feed packets into and read the codec's output back out of.
func newTestCodec(t *testing.T) (*MessageCodec, *test.TransportInstance) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)

	transport := newTestTransport(t, _options...)
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

func mustDecodeHex(t *testing.T, hexString string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(hexString)
	require.NoError(t, err)
	return decoded
}

func TestPacketSizeFromHeader(t *testing.T) {
	tests := []struct {
		name   string
		header []byte
		want   uint32
	}{
		{name: "a header-only packet is 28 bytes", header: []byte{0x01, 0x01, 0x00, 0x00}, want: 28},
		{name: "the length field counts what follows the header", header: []byte{0x01, 0x07, 0x00, 0x0E}, want: 42},
		{name: "the length field is big-endian", header: []byte{0x02, 0x07, 0x01, 0x00}, want: 284},
		// A uint16 total would wrap 0xFFFF + 28 back to 27, making Read consume less than a frame
		// and leaving the receive worker unable to ever resynchronize.
		{name: "the biggest length doesn't wrap the total", header: []byte{0x02, 0x07, 0xFF, 0xFF}, want: 65563},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, packetSizeFromHeader(testCase.header))
		})
	}
}

func TestMessageCodec_SendSerializesBigEndian(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	request := readWriteModel.NewCIPEncapsulationConnectionRequest(0, 0, connectionRequestSenderContext, 0)
	require.NoError(t, codec.Send(testutils.TestContext(t), "connect", request))

	// Byte for byte the "Connection Request" case of the parser-serializer testsuite.
	assert.Equal(t,
		mustDecodeHex(t, "01010000000000000000000000040005000000000000000000000000"),
		transportInstance.DrainWriteBuffer(28))
}

func TestMessageCodec_SendRejectsForeignMessages(t *testing.T) {
	codec, _ := newTestCodec(t)
	// A DF1 message on its own is not a CIP encapsulation packet, so it can't be framed.
	err := codec.Send(testutils.TestContext(t), "nonsense",
		readWriteModel.NewDF1CommandRequestMessage(0, 5, 0, 1,
			readWriteModel.NewDF1RequestProtectedTypedLogicalRead(2, 7, 0x89, 3, 0)))
	assert.Error(t, err)
}

func TestMessageCodec_ReceiveFramesTheStream(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// A connection response (header only) directly followed by a read response with a 24 byte
	// payload - both cases of the parser-serializer testsuite, back to back in one chunk.
	transportInstance.FillReadBuffer(mustDecodeHex(t,
		"02010000000003320000000000040005000000000000000000000000"+
			"02070020000003320000000040000000000000000000000000000000000508004f000401910101000900040405001f02010003000404050000024000"))

	connectionResponse, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	typedConnectionResponse, ok := connectionResponse.(readWriteModel.CIPEncapsulationConnectionResponse)
	require.True(t, ok, "%T is not a connection response", connectionResponse)
	assert.Equal(t, uint32(818), typedConnectionResponse.GetSessionHandle())

	readResponse, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	typedReadResponse, ok := readResponse.(readWriteModel.CIPEncapsulationReadResponse)
	require.True(t, ok, "%T is not a read response", readResponse)
	assert.Equal(t, uint32(818), typedReadResponse.GetSessionHandle())
	assert.Equal(t, uint16(1025), typedReadResponse.GetResponse().GetTransactionCounter())
	logicalRead, ok := typedReadResponse.GetResponse().(readWriteModel.DF1CommandResponseMessageProtectedTypedLogicalRead)
	require.True(t, ok)
	assert.Len(t, logicalRead.GetData(), 24, "the payload length follows from packetLen - 8")
}

func TestMessageCodec_ReceiveWaitsForTheWholeFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// A partial frame leaves the codec waiting for the rest of it, so the two incomplete steps
	// below get a context of their own that runs out - which is what a real transport going quiet
	// looks like from in here.
	receivePartial := func() any {
		t.Helper()
		ctx, cancel := context.WithTimeout(testutils.TestContext(t), 250*time.Millisecond)
		defer cancel()
		message, err := codec.Receive(ctx)
		require.NoError(t, err)
		return message
	}

	// Only the first two bytes of a 28 byte connection response: not even the length field is
	// complete, so there is nothing to decide a frame size from.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "0201"))
	assert.Nil(t, receivePartial(), "a packet whose length isn't known yet must not be framed")

	// The length is decidable now, but the body is still missing.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "00000000033200000000"))
	assert.Nil(t, receivePartial(), "an incomplete packet must not be framed")

	// And now the rest of it.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "00040005000000000000000000000000"))
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	_, ok := message.(readWriteModel.CIPEncapsulationConnectionResponse)
	assert.True(t, ok, "%T is not a connection response", message)
}

func TestMessageCodec_ReceiveDropsAMalformedFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// 0x9999 is not a command type the type switch knows, so the packet can't be parsed. The frame
	// is consumed all the same, which is what keeps the stream in sync for the packet after it.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "99990000000003320000000000040005000000000000000000000000"))
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	assert.Nil(t, message)

	transportInstance.FillReadBuffer(mustDecodeHex(t, "02010000000003320000000000040005000000000000000000000000"))
	message, err = codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message, "the codec has to resynchronize on the packet after a bad one")
	_, ok := message.(readWriteModel.CIPEncapsulationConnectionResponse)
	assert.True(t, ok, "%T is not a connection response", message)
}
