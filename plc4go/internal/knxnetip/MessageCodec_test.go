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

package knxnetip

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

func newTestHPAIControlEndpoint() model.HPAIControlEndpoint {
	return model.NewHPAIControlEndpoint(
		model.HostProtocolCode_IPV4_UDP, model.NewIPAddress([]byte{192, 168, 42, 11}), 3671)
}

func newTestSearchResponse() model.SearchResponse {
	return model.NewSearchResponse(
		newTestHPAIControlEndpoint(),
		model.NewDIBDeviceInfo(
			54,
			model.KnxMedium_MEDIUM_TP1,
			model.NewDeviceStatus(false),
			model.NewKnxAddress(1, 1, 0),
			model.NewProjectInstallationIdentifier(1, 2),
			[]byte{1, 2, 3, 4, 5, 6},
			model.NewIPAddress([]byte{224, 0, 23, 12}),
			model.NewMACAddress([]byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}),
			[]byte("test-gateway"),
		),
		model.NewDIBSuppSvcFamilies(2, nil),
	)
}

func newTestTunnelingRequest(communicationChannelId, sequenceCounter uint8) model.TunnelingRequest {
	return model.NewTunnelingRequest(
		model.NewTunnelingRequestDataBlock(communicationChannelId, sequenceCounter),
		model.NewLDataInd(0, nil,
			model.NewLDataExtended(
				true, true, model.CEMIPriority_LOW, false, false, true, 6, 0,
				model.NewKnxAddress(1, 1, 1),
				[]byte{0x08, 0x01},
				model.NewApduDataContainer(false, 0,
					model.NewApduDataGroupValueWrite(1, nil)),
			),
		),
	)
}

// newTestMessageCodec spins up a MessageCodec on top of a (connected, but not
// running) test-transport. The returned transport instance can be used to check
// what the codec wrote and the returned slice-pointer collects every message the
// message-interceptor got to see.
func newTestMessageCodec(t *testing.T) (*MessageCodec, *test.TransportInstance, *[]spi.Message) {
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

	var intercepted []spi.Message
	codec := NewMessageCodec(testTransportInstance, func(message spi.Message) {
		intercepted = append(intercepted, message)
	}, _options...)
	return codec, testTransportInstance, &intercepted
}

// Test_CustomMessageHandling_passesNonTunnelingFramesOn is the regression test for
// the single-value type-assertions which used to panic (and thus silently drop)
// every frame which was neither a TunnelingRequest nor a TunnelingResponse. Those
// are exactly the frames the connect-handshake and the keepalive wait for.
func Test_CustomMessageHandling_passesNonTunnelingFramesOn(t *testing.T) {
	tests := []struct {
		name    string
		message spi.Message
	}{
		{
			name:    "SearchResponse",
			message: newTestSearchResponse(),
		},
		{
			name:    "ConnectionResponse",
			message: model.NewConnectionResponse(1, model.Status_NO_ERROR, nil, nil),
		},
		{
			name:    "ConnectionStateResponse",
			message: model.NewConnectionStateResponse(1, model.Status_NO_ERROR),
		},
		{
			name:    "DisconnectResponse",
			message: model.NewDisconnectResponse(1, model.Status_NO_ERROR),
		},
		{
			name:    "DisconnectRequest",
			message: model.NewDisconnectRequest(1, newTestHPAIControlEndpoint()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec, transportInstance, intercepted := newTestMessageCodec(t)

			handler := CustomMessageHandling(testutils.ProduceTestingLogger(t))

			assert.NotPanics(t, func() {
				assert.False(t, handler(t.Context(), codec, tt.message),
					"non-tunneling frames must not be marked as handled, otherwise no expectation ever sees them")
			})

			assert.Equal(t, []spi.Message{tt.message}, *intercepted,
				"the message interceptor must see the frame")
			assert.Zero(t, transportInstance.GetNumDrainableBytes(),
				"no response should have been sent for a non-tunneling frame")
		})
	}
}

// Test_CustomMessageHandling_tunnelingResponseIsPassedOn is the regression test for the
// tunneling ACK which used to be swallowed here: DefaultCodec.ReceiveWork skips the
// expectations completely for a message the custom handler reports as handled, so a
// swallowed ACK made every correlated tunneling-request (e.g. a group-address write) run
// into its request-timeout.
func Test_CustomMessageHandling_tunnelingResponseIsPassedOn(t *testing.T) {
	codec, transportInstance, intercepted := newTestMessageCodec(t)

	handler := CustomMessageHandling(testutils.ProduceTestingLogger(t))

	tunnelingResponse := model.NewTunnelingResponse(
		model.NewTunnelingResponseDataBlock(1, 2, model.Status_NO_ERROR))

	assert.False(t, handler(t.Context(), codec, tunnelingResponse),
		"an ACK must reach the expectations, otherwise no tunneling request can ever be correlated")
	assert.Equal(t, []spi.Message{spi.Message(tunnelingResponse)}, *intercepted)
	assert.Zero(t, transportInstance.GetNumDrainableBytes(), "an ACK is not answered")
}

// Test_MessageCodec_tunnelingAckSatisfiesAnExpectation drives the very same case through
// the running codec, which is what the write-path relies on.
func Test_MessageCodec_tunnelingAckSatisfiesAnExpectation(t *testing.T) {
	codec, transportInstance, _ := newTestMessageCodec(t)
	require.NoError(t, codec.Connect(t.Context()))
	t.Cleanup(func() { _ = codec.Disconnect() })

	tunnelingResponse := model.NewTunnelingResponse(
		model.NewTunnelingResponseDataBlock(0x2a, 0x0b, model.Status_NO_ERROR))

	handled := make(chan model.TunnelingResponse, 1)
	codec.Expect(t.Context(), "tunneling_ack",
		func(message spi.Message) bool {
			_, ok := message.(model.TunnelingResponse)
			return ok
		},
		func(message spi.Message) error {
			handled <- message.(model.TunnelingResponse)
			return nil
		},
		func(err error) error { return nil },
	)

	theBytes, err := tunnelingResponse.Serialize()
	require.NoError(t, err)
	transportInstance.FillReadBuffer(theBytes)
	// The test transport only moves the bytes from its channel into the read buffer on an
	// actual read, while the codec only reads once it sees them as available.
	_, err = transportInstance.PeekReadableBytes(t.Context(), uint32(len(theBytes)))
	require.NoError(t, err)

	select {
	case got := <-handled:
		assert.Equal(t, uint8(0x2a), got.GetTunnelingResponseDataBlock().GetCommunicationChannelId())
		assert.Equal(t, uint8(0x0b), got.GetTunnelingResponseDataBlock().GetSequenceCounter())
	case <-time.After(10 * time.Second):
		t.Fatal("the expectation never saw the tunneling ack")
	}
}

func Test_CustomMessageHandling_tunnelingRequestIsAcknowledged(t *testing.T) {
	codec, transportInstance, intercepted := newTestMessageCodec(t)

	handler := CustomMessageHandling(testutils.ProduceTestingLogger(t))

	tunnelingRequest := newTestTunnelingRequest(0x2a, 0x0b)

	assert.False(t, handler(t.Context(), codec, tunnelingRequest),
		"a tunneling request still has to be passed on for further handling")
	assert.Equal(t, []spi.Message{spi.Message(tunnelingRequest)}, *intercepted)

	// The gateway must have gotten an ACK for the request.
	numDrainable := transportInstance.GetNumDrainableBytes()
	require.NotZero(t, numDrainable)
	rawAck := transportInstance.DrainWriteBuffer(numDrainable)
	parsedAck, err := model.KnxNetIpMessageParse[model.KnxNetIpMessage](context.Background(), rawAck)
	require.NoError(t, err)
	ack, ok := parsedAck.(model.TunnelingResponse)
	require.True(t, ok, "expected a TunnelingResponse, got %T", parsedAck)
	assert.Equal(t, uint8(0x2a), ack.GetTunnelingResponseDataBlock().GetCommunicationChannelId())
	assert.Equal(t, uint8(0x0b), ack.GetTunnelingResponseDataBlock().GetSequenceCounter())
	assert.Equal(t, model.Status_NO_ERROR, ack.GetTunnelingResponseDataBlock().GetStatus())
}

// foreignCodec is a _default.DefaultCodecRequirements implementation which is NOT our
// *MessageCodec, which is what the comma-ok type assertion in CustomMessageHandling has to
// cope with.
type foreignCodec struct {
	sent []spi.Message
}

func (f *foreignCodec) GetCodec() spi.MessageCodec { return nil }

func (f *foreignCodec) Send(_ context.Context, _ string, message spi.Message) error {
	f.sent = append(f.sent, message)
	return nil
}

func (f *foreignCodec) Receive(_ context.Context) (spi.Message, error) { return nil, nil }

// Test_CustomMessageHandling_foreignCodec makes sure we don't blow up if the
// handler is ever attached to something which isn't our own codec.
func Test_CustomMessageHandling_foreignCodec(t *testing.T) {
	handler := CustomMessageHandling(testutils.ProduceTestingLogger(t))

	t.Run("a frame without an interceptor to call", func(t *testing.T) {
		codec := &foreignCodec{}
		assert.NotPanics(t, func() {
			assert.False(t, handler(t.Context(), codec, newTestSearchResponse()))
		})
		assert.Empty(t, codec.sent)
	})
	t.Run("a tunneling request is still acknowledged", func(t *testing.T) {
		codec := &foreignCodec{}
		tunnelingRequest := newTestTunnelingRequest(0x2a, 0x0b)
		assert.NotPanics(t, func() {
			assert.False(t, handler(t.Context(), codec, tunnelingRequest))
		})
		require.Len(t, codec.sent, 1)
		ack, ok := codec.sent[0].(model.TunnelingResponse)
		require.True(t, ok, "expected a TunnelingResponse, got %T", codec.sent[0])
		assert.Equal(t, uint8(0x2a), ack.GetTunnelingResponseDataBlock().GetCommunicationChannelId())
		assert.Equal(t, uint8(0x0b), ack.GetTunnelingResponseDataBlock().GetSequenceCounter())
	})
	t.Run("our own codec without an interceptor", func(t *testing.T) {
		codec, _, _ := newTestMessageCodec(t)
		codec.messageInterceptor = nil
		assert.NotPanics(t, func() {
			assert.False(t, handler(t.Context(), codec, newTestSearchResponse()))
		})
	})
}
