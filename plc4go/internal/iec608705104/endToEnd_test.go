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

// The other tests in this package all cut the path short: the connection and subscriber tests drive a
// stub codec and the codec tests call Receive directly, so nothing there wires the real codec's
// receive worker to the connection. That gap is exactly where a push-only protocol goes wrong - a
// codec without a custom message handler stops reading the transport as soon as its expectations
// drain to zero, which here means right after the connect handshake, and then the station's
// telemetry sits in the transport buffer forever. The tests below run the whole path: octets into a
// test transport, out of the connection as a subscription event.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// sentFrames drains everything the connection has put on the wire.
func sentFrames(t *testing.T, transportInstance *test.TransportInstance) []byte {
	t.Helper()
	return transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
}

// waitForSentFrames waits until at least numBytes octets have been written and hands them back.
func waitForSentFrames(t *testing.T, transportInstance *test.TransportInstance, numBytes uint32) []byte {
	t.Helper()
	require.Eventually(t, func() bool {
		return transportInstance.GetNumDrainableBytes() >= numBytes
	}, 20*time.Second, time.Millisecond, "only %d of %d octets reached the transport",
		transportInstance.GetNumDrainableBytes(), numBytes)
	return sentFrames(t, transportInstance)
}

// newRunningConnection connects a connection whose codec really runs, answering the handshake the way
// a controlled station would: each activation is answered once it is on the wire, which keeps this
// deterministic (the expectation is registered before the frame is sent).
func newRunningConnection(t *testing.T) (*Connection, *test.TransportInstance) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec, transportInstance := newTestCodec(t)
	connection := NewConnection(
		DefaultConfiguration(),
		codec,
		map[string][]string{},
		NewTagHandler(_options...),
		_options...,
	)

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	assert.Equal(t, mustDecodeHex(t, "680443000000"), waitForSentFrames(t, transportInstance, 6),
		"the handshake opens with a test frame activation")
	transportInstance.FillReadBuffer(mustDecodeHex(t, "680483000000"))

	assert.Equal(t, mustDecodeHex(t, "680407000000"), waitForSentFrames(t, transportInstance, 6),
		"and goes on with a start-data-transfer activation")
	transportInstance.FillReadBuffer(mustDecodeHex(t, "68040b000000"))

	select {
	case err := <-connectErrors:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	require.True(t, connection.IsConnected())
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	return connection, transportInstance
}

// A codec with no outstanding expectation still has to read the transport. This driver registers
// exactly two expectations in its whole life (the handshake) and everything afterwards is
// unsolicited, so a receive worker which parks once the expectations drain would deliver nothing at
// all after connecting.
func TestEndToEnd_CodecReceivesWithoutAnyExpectation(t *testing.T) {
	codec, transportInstance := newTestCodec(t)
	require.NoError(t, codec.Connect(testutils.TestContext(t)))
	t.Cleanup(func() { assert.NoError(t, codec.Disconnect()) })

	transportInstance.FillReadBuffer(mustDecodeHex(t, "681a04000200010414000a0001000000020000000300000004000000"))

	select {
	case message := <-codec.GetDefaultIncomingMessageChannel():
		require.NotNil(t, message)
		dataFrame, ok := message.(readWriteModel.APDUIFormat)
		require.True(t, ok, "%T is not an I-format frame", message)
		assert.Equal(t, uint16(10), dataFrame.GetAsdu().GetAsduAddressField())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the pushed frame was never read off the transport")
	}
}

// The whole subscription path: an ASDU the station pushes travels through the running codec, the
// connection's incoming-message worker and the subscriber into a consumer - carrying its quality.
func TestEndToEnd_PushedAsduReachesASubscriber(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)

	handle, responseCode := subscribe(t, connection, changeOfState, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.Empty(t, sentFrames(t, transportInstance), "subscribing puts nothing on the wire")

	events := make(chan apiModel.PlcSubscriptionEvent, 8)
	registration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		events <- event
	})
	require.NotNil(t, registration)
	t.Cleanup(registration.Unregister)

	// A single point information at ASDU 10, information object 13, on and flagged invalid.
	transportInstance.FillReadBuffer(iFormatFrame(0, 0, singlePointAsdu(13, 0x81)))

	select {
	case event := <-events:
		assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("breaker"))
		assert.Equal(t, "10/13", event.GetAddress("breaker"))
		value := event.GetValue("breaker")
		assert.True(t, nested(t, value, fieldValue).GetBool())
		assert.True(t, nested(t, value, fieldQuality, "invalid").GetBool(),
			"the quality travels with the value rather than being dropped")
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the pushed ASDU never reached the subscriber")
	}
}

// A station probing the link over the real codec has to get its confirmation back on the wire.
func TestEndToEnd_TestFrameIsConfirmedOnTheWire(t *testing.T) {
	_, transportInstance := newRunningConnection(t)

	transportInstance.FillReadBuffer(mustDecodeHex(t, "680443000000"))

	assert.Equal(t, mustDecodeHex(t, "680483000000"), waitForSentFrames(t, transportInstance, 6))
}

// The acknowledgement window over the real codec: eight I-format frames in one chunk, one S-format
// acknowledgement out, carrying the sequence number of the frame expected next.
func TestEndToEnd_AcknowledgementReachesTheWire(t *testing.T) {
	_, transportInstance := newRunningConnection(t)

	var stream []byte
	for sendSequenceNo := range uint16(defaultAckThreshold) {
		stream = append(stream, iFormatFrame(sendSequenceNo, 0, singlePointAsdu(13, 0x01))...)
	}
	transportInstance.FillReadBuffer(stream)

	// 680401001000: an S-format frame whose receive sequence number field is 8 << 1.
	assert.Equal(t, mustDecodeHex(t, "680401001000"), waitForSentFrames(t, transportInstance, 6))
}
