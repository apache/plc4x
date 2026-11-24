//go:build linux

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

package profinet

// Connects to a real Profinet device when PROFINET_TEST_INTERFACE, PROFINET_TEST_DEVICE (its name
// of station) and PROFINET_TEST_GSD (a directory with its GSDML) are set. tools/pnet-peer runs the
// open-source p-net sample device for this in Docker; a Pi with p-net or any device works too.

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

func TestConnectToDevice(t *testing.T) {
	iface, name, gsd := os.Getenv("PROFINET_TEST_INTERFACE"), os.Getenv("PROFINET_TEST_DEVICE"), os.Getenv("PROFINET_TEST_GSD")
	if iface == "" || name == "" || gsd == "" {
		t.Skip("set PROFINET_TEST_INTERFACE, PROFINET_TEST_DEVICE and PROFINET_TEST_GSD to run against a device")
	}
	options := url.Values{"device-name": {name}, "gsd-directory": {gsd}, "reduction-ratio": {"32"}}
	if ip := os.Getenv("PROFINET_TEST_LOCAL_IP"); ip != "" {
		options.Set("local-ip", ip)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	connection, err := driver.GetConnection(ctx, url.URL{Scheme: rawsocket.TransportCode, Host: iface},
		map[string]transports.Transport{rawsocket.TransportCode: rawsocket.NewTransport()}, options)
	require.NoError(t, err)
	defer func() { _ = connection.Close() }()

	c := connection.(*Connection)
	for _, s := range c.ar.plan.Submodules {
		t.Logf("slot %d subslot 0x%04X module 0x%08X submodule 0x%08X in %d out %d (%s)", s.Slot, s.Subslot, s.ModuleIdent, s.SubmoduleIdent, s.Layout.InputLength, s.Layout.OutputLength, s.Resolved.ModuleID())
	}
	// Several watchdog periods: the relation must stay up on our cyclic frames.
	time.Sleep(5 * time.Second)
	assert.True(t, connection.IsConnected(), "the relation stays up")
	_, status, last := c.ar.inputImage()
	assert.WithinDuration(t, time.Now(), last, time.Second, "input frames keep arriving")
	t.Logf("input data status 0x%02X", status)
}
