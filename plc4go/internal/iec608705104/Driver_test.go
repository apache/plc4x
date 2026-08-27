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
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

func newTestDriver(t *testing.T) *Driver {
	t.Helper()
	driver, ok := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...).(*Driver)
	require.True(t, ok)
	return driver
}

func TestDriver_Identity(t *testing.T) {
	driver := newTestDriver(t)
	assert.Equal(t, "iec-60870-5-104", driver.GetProtocolCode())
	assert.Equal(t, "IEC 60870-5-104", driver.GetProtocolName())
	assert.Equal(t, "tcp", driver.GetDefaultTransport())
	assert.Equal(t, []string{"tcp", "test"}, driver.GetSupportedTransportCodes())
	assert.False(t, driver.SupportsDiscovery(), "there is no IEC 60870-5-104 discovery in plc4j either")
	assert.Equal(t, uint16(2404), readWriteModel.Constant_DEFAULTPORT, "the port the standard assigns")

	tag, err := driver.GetPlcTagHandler().ParseTag("10/13")
	require.NoError(t, err)
	assert.Equal(t, "10/13", tag.GetAddressString())
}

// A caller must not be able to reach into the driver's own slice and rewrite what it claims to
// support.
func TestDriver_GetSupportedTransportCodesIsACopy(t *testing.T) {
	driver := newTestDriver(t)
	handedOut := driver.GetSupportedTransportCodes()
	handedOut[0] = "udp"
	assert.Equal(t, []string{"tcp", "test"}, driver.GetSupportedTransportCodes())
}

// The APCI is a length-prefixed stream format, so a datagram transport is refused up front rather
// than at the first read.
func TestDriver_GetConnectionRefusesAnUnsupportedTransport(t *testing.T) {
	driver := newTestDriver(t)

	connection, err := driver.GetConnection(
		testutils.TestContext(t),
		url.URL{Scheme: "udp", Host: "127.0.0.1"},
		map[string]transports.Transport{},
		map[string][]string{})

	assert.Error(t, err)
	assert.Nil(t, connection)
}

// A supported transport which isn't registered is an error too, rather than a nil transport nobody
// notices until the first frame.
func TestDriver_GetConnectionRefusesAMissingTransport(t *testing.T) {
	driver := newTestDriver(t)

	connection, err := driver.GetConnection(
		testutils.TestContext(t),
		url.URL{Scheme: "tcp", Host: "127.0.0.1"},
		map[string]transports.Transport{},
		map[string][]string{})

	assert.Error(t, err)
	assert.Nil(t, connection)
}

// A connection string with an unusable option is refused before a socket is opened.
func TestDriver_GetConnectionRefusesABadOption(t *testing.T) {
	driver := newTestDriver(t)

	connection, err := driver.GetConnection(
		testutils.TestContext(t),
		url.URL{Scheme: "tcp", Host: "127.0.0.1"},
		map[string]transports.Transport{},
		map[string][]string{"request-timeout-ms": {"not a number"}})

	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_CheckTagAddress(t *testing.T) {
	driver := newTestDriver(t)
	assert.NoError(t, driver.CheckTagAddress("*/*"))
	assert.Error(t, driver.CheckTagAddress("nonsense"))
}
