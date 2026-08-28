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
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

func TestDriver_Metadata(t *testing.T) {
	driver := NewDriver()

	assert.Equal(t, "firmata", driver.GetProtocolCode())
	assert.Equal(t, "Firmata", driver.GetProtocolName())
	// plc4j FirmataDriver.getDefaultTransportCode: firmata's home is a UART.
	assert.Equal(t, "serial", driver.GetDefaultTransport())
	// plc4j FirmataDriver.getSupportedTransportCodes.
	assert.Equal(t, []string{"serial", "tcp", "test"}, driver.GetSupportedTransportCodes())
	// A firmata board can be asked for its capabilities, but this driver doesn't discover.
	assert.False(t, driver.SupportsDiscovery())
	assert.Error(t, driver.Discover(t.Context(), nil))
}

func TestDriver_CheckTagAddress(t *testing.T) {
	driver := NewDriver()

	assert.NoError(t, driver.CheckTagAddress("digital:4[0..1]:PULLUP"))
	assert.NoError(t, driver.CheckTagAddress("analog:4"))
	assert.Error(t, driver.CheckTagAddress("holding-register:4"))
	// Browsing isn't supported, so no query is valid.
	assert.Error(t, driver.CheckQuery("digital:*"))
}

// Firmata speaks a byte stream, which rules out the datagram transports. Refusing them up front
// keeps a mistyped connection string from looking like a board which never answers.
func TestDriver_GetConnectionRefusesUnsupportedTransports(t *testing.T) {
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)

	connection, err := driver.GetConnection(
		t.Context(),
		url.URL{Scheme: "udp", Host: "1.2.3.4"},
		map[string]transports.Transport{"udp": test.NewTransport()},
		map[string][]string{},
	)
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_GetConnectionNeedsARegisteredTransport(t *testing.T) {
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)

	connection, err := driver.GetConnection(
		t.Context(),
		url.URL{Scheme: "tcp", Host: "1.2.3.4"},
		map[string]transports.Transport{},
		map[string][]string{},
	)
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_GetConnectionRejectsABadConfiguration(t *testing.T) {
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)

	connection, err := driver.GetConnection(
		t.Context(),
		url.URL{Scheme: "test"},
		map[string]transports.Transport{"test": test.NewTransport()},
		map[string][]string{"request-timeout-ms": {"not a number"}},
	)
	assert.Error(t, err)
	assert.Nil(t, connection)
}

// The whole way in: the driver builds the codec and the connection, and the connection runs its
// handshake with the board on the other end of the transport.
func TestDriver_GetConnection(t *testing.T) {
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	driver := NewDriver(_options...)

	transportUrl := url.URL{Scheme: "test"}
	transport := test.NewTransport(_options...)
	transportInstance := test.NewTransportInstance(transport, _options...)
	require.NoError(t, transport.AddPreregisteredInstances(transportUrl, transportInstance))

	// Answer the system reset with a firmware report as soon as it shows up on the wire. This is
	// the only thing ever pushed into the read buffer, so it can't race with the codec's worker
	// pulling bytes out of it.
	go func() {
		deadline := time.Now().Add(20 * time.Second)
		for time.Now().Before(deadline) {
			if transportInstance.GetNumDrainableBytes() >= 1 {
				transportInstance.DrainWriteBuffer(1)
				transportInstance.FillReadBuffer(reportFirmwareFrame(0x02, 0x05, "StandardFirmata.ino"))
				return
			}
			time.Sleep(time.Millisecond)
		}
	}()

	connection, err := driver.GetConnection(
		t.Context(),
		transportUrl,
		map[string]transports.Transport{"test": transport},
		map[string][]string{},
	)
	require.NoError(t, err)
	require.NotNil(t, connection)
	assert.True(t, connection.IsConnected())
	// The default port for a firmata board reachable over the network, mirroring plc4j's
	// FirmataTcpTransportConfiguration.
	assert.Equal(t, "3030", defaultTcpPort)
	assert.NoError(t, connection.Close())
}
