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
	"context"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

// stubCodec is a minimal spi.MessageCodec test double. Only the methods which
// the connection-lifecycle actually needs are implemented; everything else is
// inherited from the (nil) embedded interface and would panic if it ever got
// called, which is exactly what we want from a test double.
type stubCodec struct {
	spi.MessageCodec

	transportInstance transports.TransportInstance
	sentRequests      []spi.Message
	accepts           []spi.AcceptsMessage
	// respond is invoked for every SendRequest. Returning an error makes the
	// codec report a send-failure, otherwise the (optional) reply is handed to
	// the handleMessage / handleError callbacks.
	respond func(message spi.Message, handleMessage spi.HandleMessage, handleError spi.HandleError) error
}

func (s *stubCodec) GetTransportInstance() transports.TransportInstance {
	return s.transportInstance
}

func (s *stubCodec) SendRequest(_ context.Context, _ string, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	s.sentRequests = append(s.sentRequests, message)
	s.accepts = append(s.accepts, acceptsMessage)
	if s.respond == nil {
		return nil
	}
	return s.respond(message, handleMessage, handleError)
}

// newStubbedConnection builds a Connection which talks to a stubCodec. The udp
// transport-instance is never connected, it only supplies the local address the
// request-builders need.
func newStubbedConnection(t *testing.T, connectionOptions map[string][]string, respond func(message spi.Message, handleMessage spi.HandleMessage, handleError spi.HandleError) error) (*Connection, *stubCodec) {
	t.Helper()
	transportInstance := udp.NewTransportInstance(
		&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3671}, nil, false, nil)
	codec := &stubCodec{transportInstance: transportInstance, respond: respond}
	connection := &Connection{
		messageCodec: codec,
		options:      connectionOptions,
		metadata:     &ConnectionMetadata{},
		log:          testutils.ProduceTestingLogger(t),
	}
	return connection, codec
}

func respondWithConnectionState(status driverModel.Status) func(spi.Message, spi.HandleMessage, spi.HandleError) error {
	return func(_ spi.Message, handleMessage spi.HandleMessage, _ spi.HandleError) error {
		return handleMessage(driverModel.NewConnectionStateResponse(1, status))
	}
}

func Test_Connection_IsConnected(t *testing.T) {
	t.Run("no codec at all", func(t *testing.T) {
		connection := &Connection{log: testutils.ProduceTestingLogger(t)}
		assert.False(t, connection.IsConnected())
	})
	t.Run("gateway answers", func(t *testing.T) {
		connection, codec := newStubbedConnection(t, nil, respondWithConnectionState(driverModel.Status_NO_ERROR))
		assert.True(t, connection.IsConnected(), "a successful ping means we are connected")
		require.Len(t, codec.sentRequests, 1)
		assert.IsType(t, driverModel.NewConnectionStateRequest(0, newTestHPAIControlEndpoint()), codec.sentRequests[0])
	})
	t.Run("gateway reports an error", func(t *testing.T) {
		connection, _ := newStubbedConnection(t, nil,
			func(_ spi.Message, _ spi.HandleMessage, handleError spi.HandleError) error {
				return handleError(errors.New("nope"))
			})
		assert.False(t, connection.IsConnected())
	})
	t.Run("sending fails", func(t *testing.T) {
		connection, _ := newStubbedConnection(t, nil,
			func(spi.Message, spi.HandleMessage, spi.HandleError) error {
				return errors.New("transport is gone")
			})
		assert.False(t, connection.IsConnected())
	})
	t.Run("invalidated connection", func(t *testing.T) {
		connection, codec := newStubbedConnection(t, nil, respondWithConnectionState(driverModel.Status_NO_ERROR))
		connection.invalidated.Store(true)
		assert.False(t, connection.IsConnected())
		assert.Empty(t, codec.sentRequests, "an invalidated connection must not even try to ping")
	})
}

func Test_Connection_connectionStateTimer(t *testing.T) {
	t.Run("start creates the ticker and stop tears it down", func(t *testing.T) {
		connection, _ := newStubbedConnection(t, nil, respondWithConnectionState(driverModel.Status_NO_ERROR))

		connection.startConnectionStateTimer()
		require.NotNil(t, connection.connectionStateTimer, "a keepalive ticker must be created on connect")
		require.NotNil(t, connection.quitConnectionStateTimer)
		quit := connection.quitConnectionStateTimer

		// Starting twice must not replace the running ticker.
		connection.startConnectionStateTimer()
		assert.True(t, quit == connection.quitConnectionStateTimer, "starting twice must not swap the quit channel")

		connection.stopConnectionStateTimer()
		assert.Nil(t, connection.connectionStateTimer)
		assert.Nil(t, connection.quitConnectionStateTimer)
		select {
		case <-quit:
		default:
			t.Error("the keepalive goroutine was never told to quit")
		}
		// Stopping twice must be safe.
		assert.NotPanics(t, connection.stopConnectionStateTimer)
		connection.wg.Wait()
	})
	// The keepalive used to be reset by interceptIncomingMessage, so on any installation
	// with group telegrams more often than every 60s the ConnectionStateRequest was never
	// sent and the gateway dropped the tunnel. The heartbeat has to run at a fixed rate.
	t.Run("incoming traffic must not postpone the keepalive", func(t *testing.T) {
		var keepalives atomic.Int32
		connection, _ := newStubbedConnection(t, nil,
			func(message spi.Message, handleMessage spi.HandleMessage, _ spi.HandleError) error {
				if _, ok := message.(driverModel.ConnectionStateRequest); ok {
					keepalives.Add(1)
				}
				return handleMessage(driverModel.NewConnectionStateResponse(1, driverModel.Status_NO_ERROR))
			})
		// Without a ticker the interceptor must be a no-op instead of a nil-dereference.
		assert.NotPanics(t, func() { connection.interceptIncomingMessage(nil) })

		connection.keepaliveInterval = 20 * time.Millisecond
		connection.startConnectionStateTimer()
		t.Cleanup(func() {
			connection.stopConnectionStateTimer()
			connection.wg.Wait()
		})

		// Feed the interceptor far more often than the keepalive interval.
		const wantKeepalives = 3
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) && keepalives.Load() < wantKeepalives {
			assert.NotPanics(t, func() { connection.interceptIncomingMessage(nil) })
			time.Sleep(time.Millisecond)
		}
		assert.GreaterOrEqual(t, keepalives.Load(), int32(wantKeepalives),
			"the heartbeat has to run at a fixed interval, independent of the incoming traffic")
	})
	t.Run("keepalive sends a connection state request", func(t *testing.T) {
		var gotStatusRequest atomic.Bool
		connection, codec := newStubbedConnection(t, nil,
			func(message spi.Message, handleMessage spi.HandleMessage, _ spi.HandleError) error {
				if _, ok := message.(driverModel.ConnectionStateRequest); ok {
					gotStatusRequest.Store(true)
				}
				return handleMessage(driverModel.NewConnectionStateResponse(1, driverModel.Status_NO_ERROR))
			})
		connection.CommunicationChannelId = 42

		connection.sendKeepalive()

		assert.True(t, gotStatusRequest.Load())
		require.Len(t, codec.sentRequests, 1)
		connectionStateRequest, ok := codec.sentRequests[0].(driverModel.ConnectionStateRequest)
		require.True(t, ok)
		assert.Equal(t, uint8(42), connectionStateRequest.GetCommunicationChannelId())
	})
	t.Run("keepalive survives a failing gateway", func(t *testing.T) {
		connection, _ := newStubbedConnection(t, nil,
			func(spi.Message, spi.HandleMessage, spi.HandleError) error {
				return errors.New("gateway is gone")
			})
		assert.NotPanics(t, connection.sendKeepalive)
	})
}

func Test_Connection_getTunnelConnectionType(t *testing.T) {
	tests := []struct {
		name    string
		options map[string][]string
		want    driverModel.KnxLayer
	}{
		{"no option", nil, driverModel.KnxLayer_TUNNEL_LINK_LAYER},
		{"empty option", map[string][]string{"connection-type": {}}, driverModel.KnxLayer_TUNNEL_LINK_LAYER},
		{"link layer", map[string][]string{"connection-type": {"LINK_LAYER"}}, driverModel.KnxLayer_TUNNEL_LINK_LAYER},
		{"lowercase", map[string][]string{"connection-type": {"link_layer"}}, driverModel.KnxLayer_TUNNEL_LINK_LAYER},
		{"raw", map[string][]string{"connection-type": {"RAW"}}, driverModel.KnxLayer_TUNNEL_RAW},
		{"busmonitor", map[string][]string{"connection-type": {"BUSMONITOR"}}, driverModel.KnxLayer_TUNNEL_BUSMONITOR},
		{"prefixed", map[string][]string{"connection-type": {"TUNNEL_BUSMONITOR"}}, driverModel.KnxLayer_TUNNEL_BUSMONITOR},
		{"garbage falls back", map[string][]string{"connection-type": {"nonsense"}}, driverModel.KnxLayer_TUNNEL_LINK_LAYER},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, _ := newStubbedConnection(t, tt.options, nil)
			assert.Equal(t, tt.want, connection.getTunnelConnectionType())
		})
	}
}

func Test_Connection_getRequestTimeout(t *testing.T) {
	tests := []struct {
		name    string
		options map[string][]string
		want    time.Duration
	}{
		{"no option", nil, defaultRequestTimeout},
		{"empty option", map[string][]string{"request-timeout-ms": {}}, defaultRequestTimeout},
		{"explicit value", map[string][]string{"request-timeout-ms": {"1500"}}, 1500 * time.Millisecond},
		{"zero falls back", map[string][]string{"request-timeout-ms": {"0"}}, defaultRequestTimeout},
		{"garbage falls back", map[string][]string{"request-timeout-ms": {"soon"}}, defaultRequestTimeout},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection, _ := newStubbedConnection(t, tt.options, nil)
			assert.Equal(t, tt.want, connection.getRequestTimeout())
		})
	}
}

// Test_ConnectionMetadata_GetConnectionAttributes guards against the copy-paste
// bug which rendered the multicast- and mac-address from the serial number.
func Test_ConnectionMetadata_GetConnectionAttributes(t *testing.T) {
	metadata := &ConnectionMetadata{
		KnxMedium:              driverModel.KnxMedium_MEDIUM_TP1,
		GatewayName:            "gateway",
		GatewayKnxAddress:      "1.1.0",
		ClientKnxAddress:       "1.1.255",
		ProjectNumber:          1,
		InstallationNumber:     2,
		DeviceSerialNumber:     []byte{1, 2, 3, 4, 5, 6},
		DeviceMulticastAddress: []byte{224, 0, 23, 12},
		DeviceMacAddress:       []byte{10, 20, 30, 40, 50, 60},
		SupportedServices:      []string{"core", "tunneling"},
	}

	attributes := metadata.GetConnectionAttributes()

	assert.Equal(t, "1 2 3 4 5 6", attributes["DeviceSerialNumber"])
	assert.Equal(t, "224.0.23.12", attributes["DeviceMulticastAddress"])
	assert.Equal(t, "10:20:30:40:50:60", attributes["DeviceMacAddress"])
	assert.Equal(t, "core, tunneling", attributes["SupportedServices"])
}

// A T_ACK is sent for a T_Data_Connected, and for nothing else (KNX 3/3/4, transport layer):
// not for an unnumbered T_Data_Individual, and not for a control PDU like a T_ACK, T_NAK or
// T_Connect, which would answer the device's acknowledgement with one of ours.
func Test_Connection_handleIncomingTunnelingRequest_acknowledgesOnlyConnectedData(t *testing.T) {
	device := driverModel.NewKnxAddress(1, 1, 5)
	client := driverModel.NewKnxAddress(1, 1, 250)
	frameTo := func(apdu driverModel.Apdu) driverModel.TunnelingRequest {
		return driverModel.NewTunnelingRequest(
			driverModel.NewTunnelingRequestDataBlock(1, 0),
			driverModel.NewLDataInd(0, nil, driverModel.NewLDataExtended(true, true, driverModel.CEMIPriority_SYSTEM,
				false, false, false, 6, 0, device, KnxAddressToByteArray(client), apdu)))
	}
	descriptor := driverModel.NewApduDataDeviceDescriptorResponse(0, []byte{0x07, 0xB0})
	for _, test := range []struct {
		name    string
		apdu    driverModel.Apdu
		wantAck bool
	}{
		{"T_Data_Connected", driverModel.NewApduDataContainer(true, 3, descriptor), true},
		{"T_Data_Individual", driverModel.NewApduDataContainer(false, 0, descriptor), false},
		{"T_ACK", driverModel.NewApduControlContainer(true, 3, driverModel.NewApduControlAck()), false},
		{"T_NAK", driverModel.NewApduControlContainer(true, 3, driverModel.NewApduControlNack()), false},
		{"T_Connect", driverModel.NewApduControlContainer(false, 0, driverModel.NewApduControlConnect()), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			connection, codec := newStubbedConnection(t, nil, nil)
			connection.ClientKnxAddress = client
			connection.CommunicationChannelId = 1

			connection.handleIncomingTunnelingRequest(t.Context(), frameTo(test.apdu))
			connection.wg.Wait()

			if !test.wantAck {
				assert.Empty(t, codec.sentRequests)
				return
			}
			require.Len(t, codec.sentRequests, 1)
			frame := codec.sentRequests[0].(driverModel.TunnelingRequest).GetCemi().(driverModel.LDataReq).GetDataFrame().(driverModel.LDataExtended)
			ack, ok := frame.GetApdu().(driverModel.ApduControlContainer)
			require.True(t, ok)
			assert.IsType(t, driverModel.NewApduControlAck(), ack.GetControlApdu())
			assert.Equal(t, uint8(3), ack.GetCounter())
			assert.Equal(t, KnxAddressToByteArray(device), frame.GetDestinationAddress())
		})
	}
}

// The expectations of a device disconnect and of an authorization type-asserted the frames they
// were offered without checking: the confirmation of any other control PDU to the same device,
// or any L_Data.con at all, panicked the codec worker, which lost the frame.
func Test_Connection_deviceExpectationsDontPanicOnOtherFrames(t *testing.T) {
	device := driverModel.NewKnxAddress(1, 1, 5)
	client := driverModel.NewKnxAddress(1, 1, 250)
	confirmation := func(apdu driverModel.Apdu) driverModel.TunnelingRequest {
		return driverModel.NewTunnelingRequest(
			driverModel.NewTunnelingRequestDataBlock(1, 0),
			driverModel.NewLDataCon(0, nil, driverModel.NewLDataExtended(true, true, driverModel.CEMIPriority_SYSTEM,
				false, false, false, 6, 0, client, KnxAddressToByteArray(device), apdu)))
	}
	ackConfirmation := confirmation(driverModel.NewApduControlContainer(true, 0, driverModel.NewApduControlAck()))

	// Nothing answers: each request fails right after it was sent.
	connection, codec := newStubbedConnection(t, nil, func(_ spi.Message, _ spi.HandleMessage, handleError spi.HandleError) error {
		return handleError(errors.New("no answer"))
	})
	connection.ClientKnxAddress = client
	connection.CommunicationChannelId = 1
	connection.DeviceConnections = map[uint16]*KnxDeviceConnection{}
	connection.addDeviceConnection(&KnxDeviceConnection{address: device})

	_, _ = connection.sendDeviceDisconnectionRequest(t.Context(), device)
	_, _ = connection.sendDeviceAuthentication(t.Context(), device, 0, []byte{1, 2, 3, 4})
	require.Len(t, codec.accepts, 2)
	for _, accepts := range codec.accepts {
		assert.NotPanics(t, func() {
			assert.False(t, accepts(ackConfirmation))
		})
	}
	assert.True(t, codec.accepts[0](confirmation(driverModel.NewApduControlContainer(false, 0, driverModel.NewApduControlDisconnect()))))
}
