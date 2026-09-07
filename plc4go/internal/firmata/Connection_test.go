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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newTestConnection builds a connection on a test transport whose codec isn't running. Writing and
// subscribing only ever send, so they work as they do in the field; incoming messages are handed to
// the connection directly, which is what the codec's worker would otherwise do.
func newTestConnection(t *testing.T) (*Connection, *test.TransportInstance) {
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
	return connection, transportInstance
}

// reportFirmwareFrame is what a board answers a system reset with: its version and the name of the
// sketch, one 7 bit character per two bytes.
func reportFirmwareFrame(major byte, minor byte, name string) []byte {
	frame := []byte{0xF0, 0x79, major, minor}
	for _, character := range []byte(name) {
		frame = append(frame, character, 0x00)
	}
	return append(frame, 0xF7)
}

// parseMessage parses a frame the way the codec would, so a test can hand the connection exactly
// what a board would have put on the wire.
func parseMessage(t *testing.T, frame []byte) readWriteModel.FirmataMessage {
	t.Helper()
	message, err := readWriteModel.FirmataMessageParse[readWriteModel.FirmataMessage](t.Context(), frame, true)
	require.NoError(t, err)
	return message
}

// The connect handshake is plc4j's: send a system reset, then wait for the board to report its
// firmware. Only once that arrives is there really a board on the other end.
func TestConnection_ConnectHandshake(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	// The system reset goes out first ...
	var request stubRequest
	select {
	case request = <-codec.requests:
	case <-time.After(20 * time.Second):
		require.FailNow(t, "no system reset was sent")
	}
	systemReset, ok := request.message.(readWriteModel.FirmataMessageCommand)
	require.True(t, ok)
	assert.IsType(t, readWriteModel.NewFirmataCommandSystemReset(), systemReset.GetCommand())

	// ... and the board answers with its protocol version followed by the firmware report. The
	// version is of no interest to the driver and must not be mistaken for the report.
	protocolVersion := readWriteModel.NewFirmataMessageCommand(
		readWriteModel.NewFirmataCommandProtocolVersion(0x02, 0x05))
	assert.False(t, request.acceptsMessage(protocolVersion), "the protocol version is not the firmware report")

	report := parseMessage(t, reportFirmwareFrame(0x02, 0x05, "StandardFirmata.ino"))
	require.True(t, request.acceptsMessage(report))
	require.NoError(t, request.handleMessage(report))

	select {
	case err := <-connectErrors:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	assert.True(t, connection.IsConnected())

	require.NoError(t, connection.Close())
	assert.False(t, connection.IsConnected())
}

// A board which doesn't report its firmware isn't a board we can work with. plc4j closes the
// connection and fails the connect the same way.
func TestConnection_ConnectFailsWithoutAFirmwareReport(t *testing.T) {
	codec := newStubCodec()
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 250 * time.Millisecond
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())

	err := connection.Connect(testutils.TestContext(t))

	assert.Error(t, err)
	assert.False(t, connection.IsConnected())
	assert.Len(t, codec.getSent(), 2, "the system reset and the report-firmware request went out even though nothing came back")
}

// A board which reports an error instead of its firmware fails the connect just the same.
func TestConnection_ConnectFailsOnAHandshakeError(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	select {
	case request := <-codec.requests:
		require.NoError(t, request.handleError(errors.New("the board is gone")))
	case <-time.After(20 * time.Second):
		require.FailNow(t, "no system reset was sent")
	}

	select {
	case err := <-connectErrors:
		assert.Error(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	assert.False(t, connection.IsConnected())
}

// The firmware report is asked for rather than waited on. StandardFirmata's systemResetCallback
// only clears the pin state - printFirmwareVersion runs from setup() - so a board volunteers the
// report only when it really restarts. Over serial that is masked, as opening the port toggles DTR
// and resets the AVR, but nothing resets a WiFi or Ethernet board when a TCP connection is opened,
// which is the transport this driver newly speaks. plc4j's onConnect only ever waits.
func TestConnection_HandshakeAsksForTheFirmwareReport(t *testing.T) {
	codec := newStubCodec()
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 250 * time.Millisecond
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())

	// Nothing is fed back, so the handshake runs into its timeout - what matters is what went out.
	assert.Error(t, connection.Connect(testutils.TestContext(t)))

	assert.Equal(t, []byte{
		0xFF,             // system reset
		0xF0, 0x79, 0xF7, // report firmware: asked for, not waited for
	}, codec.sentBytesOf(t))
}

// A board which only answers the request - the TCP case, where the reset never restarts it - still
// completes the handshake.
func TestConnection_HandshakeCompletesOnTheAnsweredRequest(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()

	var request stubRequest
	select {
	case request = <-codec.requests:
	case <-time.After(20 * time.Second):
		require.FailNow(t, "no system reset was sent")
	}
	// Wait for the request to be out before answering it: it is the request the board reacts to.
	require.Eventually(t, func() bool { return len(codec.getSent()) == 2 }, 20*time.Second, 10*time.Millisecond)

	report := parseMessage(t, reportFirmwareFrame(0x02, 0x05, "StandardFirmataWiFi.ino"))
	require.True(t, request.acceptsMessage(report))
	require.NoError(t, request.handleMessage(report))

	select {
	case err := <-connectErrors:
		require.NoError(t, err)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the connect never returned")
	}
	require.NoError(t, connection.Close())
}

// A transport which carries the reset but not the report-firmware request fails the connect right
// away rather than waiting out the timeout for a report which was never asked for. Both halves of
// that need pinning: the error has to name the send which failed (a handshake which swallowed it and
// then timed out would report an error too), and the connect has to be back long before the request
// timeout - deliberately set far above the bound below - could have expired.
func TestConnection_ConnectFailsWhenTheFirmwareRequestCantBeSent(t *testing.T) {
	codec := newStubCodec()
	// The system reset gets out, the report-firmware request behind it doesn't.
	codec.failSendsAfter(1)
	configuration := DefaultConfiguration()
	configuration.requestTimeout = time.Minute
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())

	started := time.Now()
	err := connection.Connect(testutils.TestContext(t))
	elapsed := time.Since(started)

	require.Error(t, err)
	assert.ErrorContains(t, err, "error requesting the firmware report",
		"the connect has to fail on the send which failed")
	assert.ErrorContains(t, err, "the transport is gone", "and it has to carry the transport's error")
	assert.Less(t, elapsed, 10*time.Second,
		"the connect has to fail on the failed send instead of waiting out the %s request timeout",
		configuration.requestTimeout)
	assert.False(t, connection.IsConnected())
	assert.Len(t, codec.getSent(), 1)
}

// A transport which is gone can't even carry the system reset.
func TestConnection_ConnectFailsWhenTheResetCantBeSent(t *testing.T) {
	codec := newStubCodec()
	codec.failSends()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())

	assert.Error(t, connection.Connect(testutils.TestContext(t)))
	assert.False(t, connection.IsConnected())
}

// Everything the board pushes after the handshake arrives through the codec's default channel, which
// the connection pumps into its value caches and on to the subscribers.
func TestConnection_IncomingWorker(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })

	connectErrors := make(chan error, 1)
	go func() {
		connectErrors <- connection.Connect(testutils.TestContext(t))
	}()
	select {
	case request := <-codec.requests:
		require.NoError(t, request.handleMessage(parseMessage(t, reportFirmwareFrame(0x02, 0x05, "x"))))
	case <-time.After(20 * time.Second):
		require.FailNow(t, "no system reset was sent")
	}
	require.NoError(t, <-connectErrors)

	codec.incoming <- readWriteModel.NewFirmataMessageAnalogIO(3, []int8{0x05, 0x01})
	assert.Eventually(t, func() bool {
		return connection.analogValue(3) == 0x85
	}, 20*time.Second, time.Millisecond, "the analog sample never made it into the cache")
}

// Close has to work whether or not Connect ever ran, as a failed connect closes on its way out.
func TestConnection_CloseWithoutConnect(t *testing.T) {
	connection, _ := newTestConnection(t)
	assert.NoError(t, connection.Close())
	assert.NoError(t, connection.Close())
}

func TestConnection_Ping(t *testing.T) {
	connection, _ := newTestConnection(t)
	// Firmata has nothing to ping with - it acknowledges nothing - so all a ping can report is
	// whether the connection is still up.
	assert.Error(t, connection.Ping(testutils.TestContext(t)), "an unconnected connection can't be pinged")
	connection.SetConnected(true)
	assert.NoError(t, connection.Ping(testutils.TestContext(t)))
}

func TestDecodeAnalogValue(t *testing.T) {
	tests := []struct {
		name    string
		data    []int8
		want    int16
		wantErr bool
	}{
		{name: "zero", data: []int8{0x00, 0x00}, want: 0},
		{name: "the least significant part comes first", data: []int8{0x05, 0x00}, want: 5},
		{name: "the most significant part is shifted by seven", data: []int8{0x00, 0x01}, want: 128},
		{name: "both parts", data: []int8{0x7F, 0x7F}, want: 16383},
		{name: "the eighth bit of a data byte is not part of the value", data: []int8{-1, -1}, want: 16383},
		{name: "a truncated message", data: []int8{0x00}, wantErr: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := decodeAnalogValue(testCase.data)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, value)
		})
	}
}

func TestDecodeDigitalPortValue(t *testing.T) {
	tests := []struct {
		name    string
		data    []int8
		want    uint8
		wantErr bool
	}{
		{name: "zero", data: []int8{0x00, 0x00}, want: 0x00},
		{name: "the first seven pins live in the first byte", data: []int8{0x7F, 0x00}, want: 0x7F},
		{name: "the eighth pin lives in the second byte", data: []int8{0x00, 0x01}, want: 0x80},
		{name: "all eight pins", data: []int8{0x7F, 0x01}, want: 0xFF},
		{name: "a truncated message", data: []int8{0x00}, wantErr: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := decodeDigitalPortValue(testCase.data)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, value)
		})
	}
}

// The value caches only report a change once, however often the board resends the same value.
func TestConnection_ValueCaches(t *testing.T) {
	connection, _ := newTestConnection(t)

	assert.True(t, connection.updateAnalogValue(3, 512))
	assert.False(t, connection.updateAnalogValue(3, 512))
	assert.True(t, connection.updateAnalogValue(3, 513))
	assert.Equal(t, int16(513), connection.analogValue(3))
	assert.Equal(t, unknownAnalogValue, connection.analogValue(4), "a pin the board never reported is unknown")

	// Port 1 covers pins 8 to 15; only the pins which really changed come back, and a pin nobody
	// has heard from yet counts as low rather than as unknown.
	changed, err := connection.updateDigitalValues(1, []int8{0x05, 0x00})
	require.NoError(t, err)
	assert.Equal(t, []uint8{8, 10}, changed)
	assert.True(t, connection.digitalValue(8))
	assert.False(t, connection.digitalValue(9))
	assert.True(t, connection.digitalValue(10))

	changed, err = connection.updateDigitalValues(1, []int8{0x07, 0x00})
	require.NoError(t, err)
	assert.Equal(t, []uint8{9}, changed)
	assert.True(t, connection.digitalValue(9))

	changed, err = connection.updateDigitalValues(1, []int8{0x07, 0x00})
	require.NoError(t, err)
	assert.Empty(t, changed, "the same report again is no change")
}

// A malformed message must not take the incoming worker down.
func TestConnection_HandleIncomingMessageSurvivesGarbage(t *testing.T) {
	connection, _ := newTestConnection(t)

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageAnalogIO(3, []int8{0x01}))
	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(0, []int8{0x01}))
	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageCommand(readWriteModel.NewFirmataCommandSystemReset()))

	assert.Equal(t, unknownAnalogValue, connection.analogValue(3), "a truncated message must not be cached")
}

func TestConnection_String(t *testing.T) {
	connection, _ := newTestConnection(t)
	assert.Contains(t, connection.String(), "firmata.Connection")
}

// Driver.GetConnection hands the connection its options as append(d._options, ...), which leaves
// spare capacity in the backing array. Both request builders then append their own logger option to
// that very slice, and two appends onto a slice with spare capacity write into the same slot - which
// is a data race whose failure mode is a torn interface value read by options.ExtractCustomLogger.
// Only -race sees it, which is why this test exists at all.
func TestConnection_RequestBuildersDoNotRaceOnTheOptions(t *testing.T) {
	base := testutils.EnrichOptionsWithOptionsForTesting(t)
	// Exactly the shape the driver hands over: a list with room left in its backing array.
	_options := make([]options.WithOption, 0, len(base)+4)
	_options = append(_options, base...)
	require.Greater(t, cap(_options), len(_options), "this test needs spare capacity to race on")

	connection := NewConnection(DefaultConfiguration(), newStubCodec(), map[string][]string{}, NewTagHandler(), _options...)
	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() { assert.NotNil(t, connection.WriteRequestBuilder()) })
		wg.Go(func() { assert.NotNil(t, connection.SubscriptionRequestBuilder()) })
	}
	wg.Wait()
}
