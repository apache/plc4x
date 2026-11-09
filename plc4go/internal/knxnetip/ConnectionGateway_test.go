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
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
)

// Tests of the connection against fakeKnxGateway, over the real codec and udp transport.

// Connecting used to panic on every gateway: each service family of the SEARCH_RESPONSE was
// type-asserted to an interface with GetTypeName, which no generated model type has.
func Test_Connection_Connect_listsTheServiceFamilies(t *testing.T) {
	gateway := newFakeKnxGateway(t)
	connection := connectToFakeGateway(t, t.Context(), gateway)

	assert.Equal(t, []string{"KnxNetIpCore v1", "KnxNetIpTunneling v1"}, connection.metadata.SupportedServices)
	assert.Empty(t, gateway.Violations())
}

// The connection created its codec without the options, so the codec logged to the global
// logger instead of the connection's.
func Test_Connection_codecUsesTheConnectionLogger(t *testing.T) {
	var global bytes.Buffer
	previous := log.Logger
	log.Logger = zerolog.New(&global).Level(zerolog.TraceLevel)
	t.Cleanup(func() { log.Logger = previous })

	gateway := newFakeKnxGateway(t)
	connectToFakeGateway(t, t.Context(), gateway)

	assert.Empty(t, global.String())
}

// Individual addresses were compared with == and != on the generated model types, which are
// pointers, so a frame never matched the address it was meant for: a device that does not exist
// was only noticed when the context ran out, instead of from the gateway's negative L_Data.con.
func Test_Connection_DeviceConnect_deviceThatDoesNotExist(t *testing.T) {
	gateway := newFakeKnxGateway(t)
	connection := connectToFakeGateway(t, t.Context(), gateway)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	select {
	case result := <-connection.DeviceConnect(ctx, driverModel.NewKnxAddress(1, 1, 9)):
		require.Error(t, result.err)
		assert.ErrorContains(t, result.err, "error connecting to device at: 1.1.9")
		assert.NotErrorIs(t, result.err, context.DeadlineExceeded)
	case <-ctx.Done():
		t.Fatal("DeviceConnect did not complete")
	}
	_, connected := connection.getDeviceConnection(driverModel.NewKnxAddress(1, 1, 9))
	assert.False(t, connected)
	assert.Empty(t, gateway.Violations())
}

// The device connections were keyed by the address pointer, so a lookup with another instance of
// the same address never found the connection.
func Test_Connection_deviceConnectionsAreKeyedByAddressValue(t *testing.T) {
	connection := &Connection{DeviceConnections: map[uint16]*KnxDeviceConnection{}}
	deviceConnection := &KnxDeviceConnection{address: driverModel.NewKnxAddress(1, 2, 3)}
	connection.addDeviceConnection(deviceConnection)

	found, ok := connection.getDeviceConnection(driverModel.NewKnxAddress(1, 2, 3))
	assert.True(t, ok)
	assert.Same(t, deviceConnection, found)
	_, ok = connection.getDeviceConnection(driverModel.NewKnxAddress(1, 2, 4))
	assert.False(t, ok)
	addresses := connection.deviceConnectionAddresses()
	require.Len(t, addresses, 1)
	assert.Equal(t, "1.2.3", KnxAddressToString(addresses[0]))

	connection.removeDeviceConnection(driverModel.NewKnxAddress(1, 2, 3))
	_, ok = connection.getDeviceConnection(driverModel.NewKnxAddress(1, 2, 3))
	assert.False(t, ok)
}
