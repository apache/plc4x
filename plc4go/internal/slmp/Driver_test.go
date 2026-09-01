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

package slmp

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

func newTestDriver(t *testing.T) *Driver {
	t.Helper()
	driver, ok := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...).(*Driver)
	require.True(t, ok)
	t.Cleanup(func() {
		assert.NoError(t, driver.Close())
	})
	return driver
}

func newTestTransport(t *testing.T, _options ...options.WithOption) *test.Transport {
	t.Helper()
	return test.NewTransport(_options...)
}

func TestDriver_Identity(t *testing.T) {
	driver := newTestDriver(t)
	assert.Equal(t, "slmp", driver.GetProtocolCode())
	assert.Equal(t, "SLMP (MELSEC) 3E", driver.GetProtocolName())
	assert.Equal(t, defaultTransportCode, driver.GetDefaultTransport())
	assert.Equal(t, []string{"tcp", "test"}, driver.GetSupportedTransportCodes())
	assert.False(t, driver.SupportsDiscovery(), "there is no slmp discovery in plc4j either")

	tag, err := driver.GetPlcTagHandler().ParseTag("D350")
	require.NoError(t, err)
	assert.Equal(t, "D350:WORD", tag.GetAddressString())
}

// TestDriver_DefaultPortComesFromTheMspec keeps the default port and the wire spec from drifting
// apart; plc4j's SlmpTcpTransportConfiguration reads the same constant.
func TestDriver_DefaultPortComesFromTheMspec(t *testing.T) {
	assert.Equal(t, 5007, SlmpPort)
}

// TestDriver_GetSupportedTransportCodesIsACopy keeps a caller from reaching into the driver's own
// slice and rewriting what the driver claims to support.
func TestDriver_GetSupportedTransportCodesIsACopy(t *testing.T) {
	driver := newTestDriver(t)
	handedOut := driver.GetSupportedTransportCodes()
	handedOut[0] = "udp"
	assert.Equal(t, []string{"tcp", "test"}, driver.GetSupportedTransportCodes())
}

func TestDriver_GetConnectionRejectsUnsupportedTransports(t *testing.T) {
	driver := newTestDriver(t)
	// A 3E frame is length-delimited over a byte stream, so a datagram transport can't carry it -
	// even when one is registered.
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "udp", Host: "192.168.0.1"},
		map[string]transports.Transport{"udp": nil},
		map[string][]string{})
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_GetConnectionRejectsAMissingTransport(t *testing.T) {
	driver := newTestDriver(t)
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "tcp", Host: "192.168.0.1"},
		map[string]transports.Transport{},
		map[string][]string{})
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_GetConnectionRejectsBadOptions(t *testing.T) {
	driver := newTestDriver(t)
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	transport := newTestTransport(t, _options...)
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "test", Host: "hurz"},
		map[string]transports.Transport{"test": transport},
		// Beyond the unsigned 16-bit monitoring-timer field in the 3E frame.
		map[string][]string{"monitoring-timer": {"65536"}})
	assert.Error(t, err)
	assert.Nil(t, connection)
}

// TestDriver_GetConnectionUsesTheSlmpPort makes sure the driver hands the transport the port an SLMP
// device listens on by default, so a connection string without one works.
func TestDriver_GetConnectionUsesTheSlmpPort(t *testing.T) {
	driver := newTestDriver(t)
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	transport := newTestTransport(t, _options...)
	driverOptions := map[string][]string{}
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "test", Host: "hurz"},
		map[string]transports.Transport{"test": transport},
		driverOptions)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})
	assert.Equal(t, []string{"5007"}, driverOptions["defaultTcpPort"])
	// SLMP has no handshake, so a connected transport is a connected device.
	assert.True(t, connection.IsConnected())
	metadata := connection.GetMetadata()
	assert.True(t, metadata.CanRead())
	assert.True(t, metadata.CanWrite())
	assert.True(t, metadata.CanSubscribe())
}
