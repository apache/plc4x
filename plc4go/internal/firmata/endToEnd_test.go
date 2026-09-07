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

// The other tests in this package all cut the path short: the connection tests drive a stub codec
// and the codec tests call Receive directly, so nothing ever wires the real codec's receive worker
// to the connection. That gap is what let a parked receive worker through - a codec without a
// custom message handler stops reading the transport as soon as its expectations drain to zero,
// which for a push-only protocol means after the connect handshake. The tests here run the whole
// path: bytes into a test transport, out of the connection as a subscription event.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// reportFirmwareRequestFrame is what the connect handshake asks the board for after resetting it.
var reportFirmwareRequestFrame = []byte{0xF0, 0x79, 0xF7}

// newRunningConnection connects a connection whose codec really runs, answering the handshake the
// way a board would: the system reset and the report-firmware request go out on the transport, and
// the firmware report is pushed back once they have.
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

	// Answering only once the request is on the wire is what keeps this deterministic: the
	// expectation matching the report is registered before the request is sent.
	require.Eventually(t, func() bool {
		return transportInstance.GetNumDrainableBytes() >= 4
	}, 20*time.Second, time.Millisecond, "the handshake never reached the transport")
	assert.Equal(t,
		append([]byte{0xFF}, reportFirmwareRequestFrame...),
		sentBytes(t, transportInstance),
		"the handshake is a system reset followed by an explicit report-firmware request")

	transportInstance.FillReadBuffer(reportFirmwareFrame(0x02, 0x05, "StandardFirmata.ino"))
	select {
	case err := <-connectErrors:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	return connection, transportInstance
}

// A codec with no outstanding expectation still has to read the transport: firmata registers exactly
// one expectation in its whole life (the handshake), and everything the board says afterwards is
// unsolicited. A receive worker which parks once the expectations drain to zero would leave those
// messages sitting in the transport buffer forever.
func TestEndToEnd_CodecReceivesWithoutAnyExpectation(t *testing.T) {
	codec, transportInstance := newTestCodec(t)
	require.NoError(t, codec.Connect(testutils.TestContext(t)))
	t.Cleanup(func() { assert.NoError(t, codec.Disconnect()) })

	transportInstance.FillReadBuffer([]byte{0xE3, 0x05, 0x01})

	select {
	case message := <-codec.GetDefaultIncomingMessageChannel():
		require.NotNil(t, message)
		analogIO, ok := message.(interface{ GetPin() uint8 })
		require.True(t, ok, "%T is not an analog-IO message", message)
		assert.Equal(t, uint8(3), analogIO.GetPin())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the pushed message was never read off the transport")
	}
}

// The whole subscription path: a pin update the board pushes has to travel through the running
// codec, the connection's incoming-message worker and the subscriber into a consumer.
func TestEndToEnd_PushedAnalogUpdateReachesASubscriber(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)

	handle, responseCode := subscribe(t, connection, "dial", "analog:2")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.Equal(t, []byte{0xC2, 0x01}, sentBytes(t, transportInstance), "reporting was switched on for pin 2")

	events := make(chan apiModel.PlcSubscriptionEvent, 4)
	registration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		events <- event
	})
	require.NotNil(t, registration)
	t.Cleanup(registration.Unregister)

	// Analog IO on pin 2: 0x05 as the low 7 bits and 0x01 as the next 7, which is 133.
	transportInstance.FillReadBuffer([]byte{0xE2, 0x05, 0x01})

	select {
	case event := <-events:
		assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("dial"))
		assert.Equal(t, int16(133), event.GetValue("dial").GetInt16())
		assert.Equal(t, "analog:2", event.GetAddress("dial"))
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the pushed analog sample never reached the subscriber")
	}
}

// The same for a digital port, which is the other half of what a board pushes.
func TestEndToEnd_PushedDigitalUpdateReachesASubscriber(t *testing.T) {
	connection, transportInstance := newRunningConnection(t)

	handle, responseCode := subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.Equal(t, []byte{0xF4, 0x09, 0x00, 0xD1, 0x01}, sentBytes(t, transportInstance))

	events := make(chan apiModel.PlcSubscriptionEvent, 4)
	registration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		events <- event
	})
	require.NotNil(t, registration)
	t.Cleanup(registration.Unregister)

	// Digital IO for port 1 (pins 8 to 15) with bit 1 - pin 9 - set.
	transportInstance.FillReadBuffer([]byte{0x91, 0x02, 0x00})

	select {
	case event := <-events:
		assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("button"))
		assert.True(t, event.GetValue("button").GetBool())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the pushed digital port never reached the subscriber")
	}
}
