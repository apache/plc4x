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
	"strings"
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

// The responses of a device were acknowledged from inside the expectation handler, which runs on
// the codec worker while it holds the expectation lock; registering the T_ACK's expectation then
// waited for that lock forever. The codec stopped, and with it the acknowledgement of the
// gateway's tunnelling requests.
func Test_Connection_DeviceConnect_existingDevice(t *testing.T) {
	gateway := newFakeKnxGateway(t, newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), nil))
	connection := connectToFakeGateway(t, t.Context(), gateway)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	select {
	case result := <-connection.DeviceConnect(ctx, driverModel.NewKnxAddress(1, 1, 5)):
		require.NoError(t, result.err)
		assert.Equal(t, uint16(0x07B0), result.connection.deviceDescriptor)
		assert.Equal(t, uint16(15), result.connection.maxApdu)
	case <-ctx.Done():
		t.Fatal("DeviceConnect did not complete")
	}
	_, connected := connection.getDeviceConnection(driverModel.NewKnxAddress(1, 1, 5))
	assert.True(t, connected)

	// The gateway still gets its tunnelling requests acknowledged.
	gateway.SendGroupValueWrite(driverModel.NewKnxAddress(1, 1, 5), []byte{0x09, 0x01}, 1, nil)
	time.Sleep(1500 * time.Millisecond)
	assert.NotContains(t, strings.Join(gateway.Violations(), "\n"), "never acknowledged")
}

// DeviceReadProperty, DeviceReadPropertyDescriptor, DeviceReadMemory and DeviceAuthenticate
// connect to the device first if they are not connected yet, but ignored the connection
// DeviceConnect returned and reported "unable to connect to device" every time.
func Test_Connection_DeviceReadProperty_connectsFirst(t *testing.T) {
	gateway := newFakeKnxGateway(t, newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), nil))
	connection := connectToFakeGateway(t, t.Context(), gateway)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	result := <-connection.DeviceReadProperty(ctx, driverModel.NewKnxAddress(1, 1, 5), 0, 56, 1, 1)
	require.NoError(t, result.err)
	assert.Equal(t, uint16(15), result.value.GetUint16())
}

// DeviceReadMemory decoded the memory with the datapoint parser, which expects the layout of a
// group value payload: a reserved byte in front of every value of eight bits or more. Memory
// holds the bare values, so a UINT read ran out of data and a USINT read skipped every other
// byte.
func Test_Connection_DeviceReadMemory_decodesBareValues(t *testing.T) {
	memory := map[uint16][]byte{0x4000: {0x00, 0x02, 0x09, 0x01, 0x09, 0x02}}
	for _, test := range []struct {
		datapointType driverModel.KnxDatapointType
		numElements   uint8
		want          []uint64
	}{
		{driverModel.KnxDatapointType_UINT, 3, []uint64{2, 0x0901, 0x0902}},
		{driverModel.KnxDatapointType_USINT, 4, []uint64{0, 2, 9, 1}},
		{driverModel.KnxDatapointType_UDINT, 1, []uint64{0x00020901}},
		{driverModel.KnxDatapointType_UINT, 1, []uint64{2}},
	} {
		t.Run(test.datapointType.String(), func(t *testing.T) {
			gateway := newFakeKnxGateway(t, newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), memory))
			connection := connectToFakeGateway(t, t.Context(), gateway)

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
			defer cancel()
			result := <-connection.DeviceReadMemory(ctx, driverModel.NewKnxAddress(1, 1, 5), 0x4000, test.numElements, &test.datapointType)
			require.NoError(t, result.err)
			var got []uint64
			if result.value.IsList() {
				for _, value := range result.value.GetList() {
					got = append(got, value.GetUint64())
				}
			} else {
				got = append(got, result.value.GetUint64())
			}
			assert.Equal(t, test.want, got)
			assert.Empty(t, gateway.Violations())
		})
	}
}

// A device answers a Memory_Read of memory it does not have with no data. That shrank the
// assumed APDU length until no byte fit into a request any more, and the read never ended.
func Test_Connection_DeviceReadMemory_noData(t *testing.T) {
	gateway := newFakeKnxGateway(t, newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), nil))
	connection := connectToFakeGateway(t, t.Context(), gateway)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	uint16Type := driverModel.KnxDatapointType_UINT
	select {
	case result := <-connection.DeviceReadMemory(ctx, driverModel.NewKnxAddress(1, 1, 5), 0x5000, 1, &uint16Type):
		assert.ErrorContains(t, result.err, "no data")
	case <-ctx.Done():
		t.Fatal("DeviceReadMemory did not complete")
	}
}

// A connection outlives the context it was connected with (the connection cache cancels it as
// soon as the connection is handed out), but the tunnelling ACKs, the heartbeat and everything
// else written to the socket later failed with "i/o timeout" once that context's deadline had
// passed, and the tunnelling handler sent its T_ACKs on the cancelled context.
func Test_Connection_outlivesTheConnectContext(t *testing.T) {
	gateway := newFakeKnxGateway(t, newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), nil))
	connectCtx, cancelConnect := context.WithTimeout(t.Context(), time.Second)
	connection := connectToFakeGateway(t, connectCtx, gateway)
	cancelConnect()
	time.Sleep(1500 * time.Millisecond) // past the deadline of the connect context

	gateway.SendGroupValueWrite(driverModel.NewKnxAddress(1, 1, 5), []byte{0x09, 0x01}, 1, nil)
	time.Sleep(1200 * time.Millisecond) // the gateway waits a second for the ACK
	assert.Empty(t, gateway.Violations())

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	result := <-connection.DeviceReadProperty(ctx, driverModel.NewKnxAddress(1, 1, 5), 0, 56, 1, 1)
	require.NoError(t, result.err)
	assert.Equal(t, uint16(15), result.value.GetUint16())
	require.NoError(t, connection.Ping(ctx))
	time.Sleep(1200 * time.Millisecond)
	assert.Empty(t, gateway.Violations())
}

// The tunnelling handler acknowledged unsolicited T_Data_Connected frames of a device on the
// context the connection was opened with; once that was cancelled, nothing was acknowledged.
func Test_Connection_acknowledgesUnsolicitedDataAfterTheConnectContextEnded(t *testing.T) {
	device := newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), nil)
	gateway := newFakeKnxGateway(t, device)
	connectCtx, cancelConnect := context.WithCancel(t.Context())
	connectToFakeGateway(t, connectCtx, gateway)
	cancelConnect()

	// A device repeats a response whose T_ACK it did not get.
	gateway.mu.Lock()
	device.connected, device.pending, device.seqSend = true, true, 0
	gateway.deviceSends(device, driverModel.NewApduDataContainer(true, 0, driverModel.NewApduDataDeviceDescriptorResponse(0, []byte{0x07, 0xB0})))
	gateway.mu.Unlock()

	require.Eventually(t, func() bool {
		gateway.mu.Lock()
		defer gateway.mu.Unlock()
		return !device.pending
	}, 3*time.Second, 10*time.Millisecond, "the response was never acknowledged")
	assert.Empty(t, gateway.Violations())
}

// The GroupValue_Response was matched by comparing its source, an individual address, with the
// group address that was read, via reflect.DeepEqual: never equal, so a read of a group address
// without a cached value always timed out.
func Test_Connection_ReadGroupAddress(t *testing.T) {
	gateway := newFakeKnxGateway(t)
	gateway.SetGroupValue([]byte{0x09, 0x03}, []byte{0x00, 45})
	connection := connectToFakeGateway(t, t.Context(), gateway)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	usint := driverModel.KnxDatapointType_USINT
	select {
	case result := <-connection.ReadGroupAddress(ctx, []byte{0x09, 0x03}, &usint):
		require.NoError(t, result.err)
		assert.Equal(t, uint8(45), result.value.GetUint8())
	case <-ctx.Done():
		t.Fatal("ReadGroupAddress did not complete")
	}
	assert.Empty(t, gateway.Violations())
}

// The reader and the browser were created without the connection's options, so they (and the
// default browser under the browser) logged to the global logger.
func Test_Connection_readerAndBrowserUseTheConnectionLogger(t *testing.T) {
	var global bytes.Buffer
	previous := log.Logger
	log.Logger = zerolog.New(&global).Level(zerolog.TraceLevel)
	t.Cleanup(func() { log.Logger = previous })

	gateway := newFakeKnxGateway(t, newFakeSystemBDevice(driverModel.NewKnxAddress(1, 1, 5), nil))
	connection := connectToFakeGateway(t, t.Context(), gateway)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	readRequest, err := connection.ReadRequestBuilder().AddTagAddress("maxApdu", "1.1.5#0/56").Build()
	require.NoError(t, err)
	readResult := <-readRequest.Execute(ctx)
	require.NoError(t, readResult.GetErr())
	// The group objects of a device that does not exist: the browser logs the failure.
	browseRequest, err := connection.BrowseRequestBuilder().AddQuery("objects", "1.1.9#com-obj").Build()
	require.NoError(t, err)
	browseResult := <-browseRequest.Execute(ctx)
	require.NoError(t, browseResult.GetErr())

	assert.Empty(t, global.String())
}
