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

package umas

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

func newTestDriver(t *testing.T) *Driver {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	driver, ok := NewDriver(_options...).(*Driver)
	require.True(t, ok)
	t.Cleanup(func() { assert.NoError(t, driver.Close()) })
	return driver
}

func TestDriver_Identity(t *testing.T) {
	driver := newTestDriver(t)
	assert.Equal(t, "umas", driver.GetProtocolCode())
	assert.Equal(t, "UMAS (Schneider Electric)", driver.GetProtocolName())
	assert.Equal(t, defaultTransportCode, driver.GetDefaultTransport())
	assert.Equal(t, []string{"tcp", "test"}, driver.GetSupportedTransportCodes())
}

// UMAS is tunneled inside Modbus/TCP, so it listens on the Modbus port. The number comes from the
// mspec constant rather than from a literal here.
func TestDriver_DefaultPort(t *testing.T) {
	assert.Equal(t, 502, UmasPort)
	assert.Equal(t, int(readWriteModel.Constant_UMASTCPDEFAULTPORT), UmasPort)
}

func TestDriver_TagsAreSymbolNames(t *testing.T) {
	driver := newTestDriver(t)
	tag, err := driver.GetPlcTagHandler().ParseTag("g_r32")
	require.NoError(t, err)
	assert.Equal(t, "g_r32", tag.GetAddressString())

	_, err = driver.GetPlcTagHandler().ParseTag("9nope")
	assert.Error(t, err)
}

// Modbus/TCP is a stream protocol, so a datagram transport can't carry it. Refusing here rather than
// at the first read keeps a mistyped connection string from looking like a PLC which never answers.
func TestDriver_RefusesTransportsItCannotSpeak(t *testing.T) {
	driver := newTestDriver(t)
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)

	availableTransports := map[string]transports.Transport{
		"test": test.NewTransport(_options...),
	}
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "udp", Host: "localhost"}, availableTransports, map[string][]string{})
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_ReportsAMissingTransport(t *testing.T) {
	driver := newTestDriver(t)
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "tcp", Host: "localhost"}, map[string]transports.Transport{}, map[string][]string{})
	assert.Error(t, err)
	assert.Nil(t, connection)
}

func TestDriver_ReportsBadDriverOptions(t *testing.T) {
	driver := newTestDriver(t)
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	availableTransports := map[string]transports.Transport{
		"test": test.NewTransport(_options...),
	}
	connection, err := driver.GetConnection(testutils.TestContext(t),
		url.URL{Scheme: "test", Host: "localhost"}, availableTransports,
		map[string][]string{"request-timeout-ms": {"never"}})
	assert.Error(t, err)
	assert.Nil(t, connection)
}
