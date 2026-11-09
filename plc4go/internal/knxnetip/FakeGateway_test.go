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
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

// fakeGatewayAddress, fakeTunnelAddress: the individual addresses of the fake interface and of
// the one tunnel it offers.
var (
	fakeGatewayAddress = driverModel.NewKnxAddress(1, 1, 0)
	fakeTunnelAddress  = driverModel.NewKnxAddress(1, 1, 250)
)

// fakeKnxDevice is a System B device on the line behind the fake gateway. It implements the
// connection oriented transport layer (T_Connect, T_Data_Connected with sequence numbers and
// T_ACK, T_Disconnect) and answers DeviceDescriptor_Read, PropertyValue_Read and Memory_Read.
type fakeKnxDevice struct {
	address    driverModel.KnxAddress
	properties map[[2]uint8][]byte // (object index, property id) -> value of element 1
	memory     map[uint16][]byte   // start address -> content

	connected bool
	seqRecv   uint8
	seqSend   uint8
	pending   bool
}

// fakeKnxGateway is a KNXnet/IP tunnelling interface on loopback with a TP1 line behind it. It
// is strict where a real interface or device is: every frame the client should not have sent is
// recorded as a violation, and a TUNNELING_REQUEST the client does not acknowledge within a
// second is one too.
type fakeKnxGateway struct {
	t          *testing.T
	conn       *net.UDPConn
	serviceIds []driverModel.ServiceId

	mu          sync.Mutex
	client      *net.UDPAddr
	channel     uint8
	recvSeq     uint8
	devices     map[uint16]*fakeKnxDevice
	groupValues map[uint16][]byte // group address -> GroupValue_Response payload (first byte, data)
	violations  []string
	heartbeats  int
	requests    []driverModel.LDataExtended // every L_Data.req of the client

	outbox chan driverModel.CEMI
	acks   chan uint8
	done   chan struct{}
	wg     sync.WaitGroup
}

func newFakeKnxGateway(t *testing.T, devices ...*fakeKnxDevice) *fakeKnxGateway {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	g := &fakeKnxGateway{
		t:           t,
		conn:        conn,
		serviceIds:  []driverModel.ServiceId{driverModel.NewKnxNetIpCore(1), driverModel.NewKnxNetIpTunneling(1)},
		channel:     7,
		devices:     map[uint16]*fakeKnxDevice{},
		groupValues: map[uint16][]byte{},
		outbox:      make(chan driverModel.CEMI, 100),
		acks:        make(chan uint8, 100),
		done:        make(chan struct{}),
	}
	for _, d := range devices {
		raw := KnxAddressToByteArray(d.address)
		g.devices[uint16(raw[0])<<8|uint16(raw[1])] = d
	}
	g.wg.Go(g.receive)
	g.wg.Go(g.sendLoop)
	t.Cleanup(func() {
		close(g.done)
		_ = g.conn.Close()
		g.wg.Wait()
	})
	return g
}

// newFakeSystemBDevice returns a device with the max APDU length property and the given memory.
func newFakeSystemBDevice(address driverModel.KnxAddress, memory map[uint16][]byte) *fakeKnxDevice {
	return &fakeKnxDevice{
		address: address,
		properties: map[[2]uint8][]byte{
			{0, 56}: {0, 15},               // PID_MAX_APDULENGTH
			{3, 13}: {0, 131, 10, 1, 0x12}, // PID_PROGRAM_VERSION
		},
		memory: memory,
	}
}

func (g *fakeKnxGateway) url() url.URL {
	return url.URL{Scheme: "udp", Host: g.conn.LocalAddr().String()}
}

func (g *fakeKnxGateway) violation(format string, args ...any) {
	g.violations = append(g.violations, fmt.Sprintf(format, args...))
}

// Violations returns what the client did wrong so far.
func (g *fakeKnxGateway) Violations() []string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]string(nil), g.violations...)
}

// Heartbeats returns the number of CONNECTIONSTATE_REQUESTs the client sent.
func (g *fakeKnxGateway) Heartbeats() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.heartbeats
}

// Requests returns every L_Data.req the client sent.
func (g *fakeKnxGateway) Requests() []driverModel.LDataExtended {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]driverModel.LDataExtended(nil), g.requests...)
}

// SetGroupValue sets what a GroupValue_Read of the group address is answered with.
func (g *fakeKnxGateway) SetGroupValue(groupAddress []byte, payload []byte) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.groupValues[uint16(groupAddress[0])<<8|uint16(groupAddress[1])] = payload
}

// SendGroupValueWrite puts a GroupValue_Write of a device on the line, which the gateway passes
// on to the tunnel.
func (g *fakeKnxGateway) SendGroupValueWrite(source driverModel.KnxAddress, groupAddress []byte, dataFirstByte int8, data []byte) {
	g.outbox <- driverModel.NewLDataInd(0, nil, driverModel.NewLDataExtended(true, true, driverModel.CEMIPriority_LOW,
		false, false, true, 6, 0, source, groupAddress,
		driverModel.NewApduDataContainer(false, 0, driverModel.NewApduDataGroupValueWrite(dataFirstByte, data))))
}

func (g *fakeKnxGateway) send(to *net.UDPAddr, message driverModel.KnxNetIpMessage) {
	data, err := message.Serialize()
	if err != nil {
		g.t.Errorf("fake gateway: serializing %T: %v", message, err)
		return
	}
	_, _ = g.conn.WriteToUDP(data, to)
}

func (g *fakeKnxGateway) receive() {
	buf := make([]byte, 1500)
	for {
		n, from, err := g.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		message, err := driverModel.KnxNetIpMessageParse[driverModel.KnxNetIpMessage](context.Background(), buf[:n])
		if err != nil {
			g.mu.Lock()
			g.violation("malformed frame from %s: %v", from, err)
			g.mu.Unlock()
			continue
		}
		g.handle(message, from)
	}
}

func (g *fakeKnxGateway) handle(message driverModel.KnxNetIpMessage, from *net.UDPAddr) {
	g.mu.Lock()
	defer g.mu.Unlock()
	switch m := message.(type) {
	case driverModel.SearchRequest:
		g.send(from, driverModel.NewSearchResponse(
			driverModel.NewHPAIControlEndpoint(driverModel.HostProtocolCode_IPV4_UDP,
				driverModel.NewIPAddress([]byte{127, 0, 0, 1}), uint16(g.conn.LocalAddr().(*net.UDPAddr).Port)),
			driverModel.NewDIBDeviceInfo(54, driverModel.KnxMedium_MEDIUM_TP1, driverModel.NewDeviceStatus(false),
				fakeGatewayAddress, driverModel.NewProjectInstallationIdentifier(0, 0), []byte{0, 0xfa, 0, 0, 0x11, 0},
				driverModel.NewIPAddress([]byte{224, 0, 23, 12}), driverModel.NewMACAddress([]byte{2, 0, 0, 0, 0x11, 0}),
				append([]byte("Fake Gateway"), make([]byte, 18)...)),
			driverModel.NewDIBSuppSvcFamilies(2, g.serviceIds),
		))
	case driverModel.ConnectionRequest:
		g.client, g.recvSeq = from, 0
		g.send(from, driverModel.NewConnectionResponse(g.channel, driverModel.Status_NO_ERROR,
			driverModel.NewHPAIDataEndpoint(driverModel.HostProtocolCode_IPV4_UDP, driverModel.NewIPAddress([]byte{127, 0, 0, 1}),
				uint16(g.conn.LocalAddr().(*net.UDPAddr).Port)),
			driverModel.NewConnectionResponseDataBlockTunnelConnection(fakeTunnelAddress)))
	case driverModel.ConnectionStateRequest:
		g.heartbeats++
		g.send(from, driverModel.NewConnectionStateResponse(m.GetCommunicationChannelId(), driverModel.Status_NO_ERROR))
	case driverModel.DisconnectRequest:
		g.send(from, driverModel.NewDisconnectResponse(m.GetCommunicationChannelId(), driverModel.Status_NO_ERROR))
	case driverModel.TunnelingResponse:
		select {
		case g.acks <- m.GetTunnelingResponseDataBlock().GetSequenceCounter():
		default:
		}
	case driverModel.TunnelingRequest:
		block := m.GetTunnelingRequestDataBlock()
		if block.GetCommunicationChannelId() != g.channel {
			g.violation("TUNNELING_REQUEST for unknown channel %d", block.GetCommunicationChannelId())
			return
		}
		switch block.GetSequenceCounter() {
		case g.recvSeq:
		case g.recvSeq - 1: // a repetition: acknowledge again, don't execute again
			g.send(from, driverModel.NewTunnelingResponse(driverModel.NewTunnelingResponseDataBlock(g.channel, block.GetSequenceCounter(), driverModel.Status_NO_ERROR)))
			return
		default:
			g.violation("TUNNELING_REQUEST sequence %d, expected %d", block.GetSequenceCounter(), g.recvSeq)
			return
		}
		g.recvSeq++
		g.send(from, driverModel.NewTunnelingResponse(driverModel.NewTunnelingResponseDataBlock(g.channel, block.GetSequenceCounter(), driverModel.Status_NO_ERROR)))
		g.handleCEMI(m.GetCemi())
	default:
		g.violation("unexpected %T", message)
	}
}

// handleCEMI puts an L_Data.req on the line: the client gets the L_Data.con, the devices the frame.
func (g *fakeKnxGateway) handleCEMI(cemi driverModel.CEMI) {
	req, ok := cemi.(driverModel.LDataReq)
	if !ok {
		g.violation("unexpected cEMI %T", cemi)
		return
	}
	frame, ok := req.GetDataFrame().(driverModel.LDataExtended)
	if !ok {
		g.violation("unexpected data frame %T", req.GetDataFrame())
		return
	}
	g.requests = append(g.requests, frame)
	confirm := func(ok bool) {
		g.outbox <- driverModel.NewLDataCon(0, nil, driverModel.NewLDataExtended(frame.GetFrameType(), frame.GetNotRepeated(),
			frame.GetPriority(), frame.GetAcknowledgeRequested(), !ok, frame.GetGroupAddress(), frame.GetHopCount(),
			frame.GetExtendedFrameFormat(), fakeTunnelAddress, frame.GetDestinationAddress(), frame.GetApdu()))
	}
	destination := uint16(frame.GetDestinationAddress()[0])<<8 | uint16(frame.GetDestinationAddress()[1])
	if frame.GetGroupAddress() {
		confirm(true)
		container, ok := frame.GetApdu().(driverModel.ApduDataContainer)
		if !ok {
			return
		}
		if _, ok := container.GetDataApdu().(driverModel.ApduDataGroupValueRead); ok {
			if payload, ok := g.groupValues[destination]; ok {
				g.outbox <- driverModel.NewLDataInd(0, nil, driverModel.NewLDataExtended(true, true, driverModel.CEMIPriority_LOW,
					false, false, true, 6, 0, driverModel.NewKnxAddress(1, 1, 5), frame.GetDestinationAddress(),
					driverModel.NewApduDataContainer(false, 0, driverModel.NewApduDataGroupValueResponse(int8(payload[0]), payload[1:]))))
			}
		}
		return
	}
	device, ok := g.devices[destination]
	if !ok {
		confirm(false) // nobody acknowledges the frame on the line
		return
	}
	confirm(true)
	g.handleDeviceFrame(device, frame)
}

func (g *fakeKnxGateway) deviceSends(device *fakeKnxDevice, apdu driverModel.Apdu) {
	g.outbox <- driverModel.NewLDataInd(0, nil, driverModel.NewLDataExtended(true, true, driverModel.CEMIPriority_LOW,
		false, false, false, 6, 0, device.address, KnxAddressToByteArray(fakeTunnelAddress), apdu))
}

func (g *fakeKnxGateway) handleDeviceFrame(device *fakeKnxDevice, frame driverModel.LDataExtended) {
	switch apdu := frame.GetApdu().(type) {
	case driverModel.ApduControlContainer:
		switch apdu.GetControlApdu().(type) {
		case driverModel.ApduControlConnect:
			device.connected, device.seqRecv, device.seqSend, device.pending = true, 0, 0, false
		case driverModel.ApduControlDisconnect:
			device.connected = false
		case driverModel.ApduControlAck:
			if !device.connected || !device.pending || apdu.GetCounter() != device.seqSend {
				g.violation("T_ACK %d to %s, which waits for no acknowledgement of it", apdu.GetCounter(), KnxAddressToString(device.address))
				return
			}
			device.pending = false
			device.seqSend = (device.seqSend + 1) & 0x0F
		default:
			g.violation("unexpected control apdu %T", apdu.GetControlApdu())
		}
	case driverModel.ApduDataContainer:
		if !apdu.GetNumbered() || !device.connected {
			g.violation("%T to %s without a connection", apdu.GetDataApdu(), KnxAddressToString(device.address))
			return
		}
		if apdu.GetCounter() != device.seqRecv {
			g.violation("T_Data_Connected sequence %d to %s, expected %d", apdu.GetCounter(), KnxAddressToString(device.address), device.seqRecv)
			return
		}
		if device.pending {
			g.violation("request to %s before it acknowledged the response %d", KnxAddressToString(device.address), device.seqSend)
			return
		}
		g.deviceSends(device, driverModel.NewApduControlContainer(true, apdu.GetCounter(), driverModel.NewApduControlAck()))
		device.seqRecv = (device.seqRecv + 1) & 0x0F
		response := g.execute(device, apdu.GetDataApdu())
		if response == nil {
			return
		}
		device.pending = true
		g.deviceSends(device, driverModel.NewApduDataContainer(true, device.seqSend, response))
	}
}

func (g *fakeKnxGateway) execute(device *fakeKnxDevice, request driverModel.ApduData) driverModel.ApduData {
	switch request := request.(type) {
	case driverModel.ApduDataDeviceDescriptorRead:
		return driverModel.NewApduDataDeviceDescriptorResponse(0, []byte{0x07, 0xB0})
	case driverModel.ApduDataMemoryRead:
		for start, content := range device.memory {
			if request.GetAddress() >= start && int(request.GetAddress()-start)+int(request.GetNumBytes()) <= len(content) {
				offset := int(request.GetAddress() - start)
				return driverModel.NewApduDataMemoryResponse(request.GetAddress(), content[offset:offset+int(request.GetNumBytes())])
			}
		}
		return driverModel.NewApduDataMemoryResponse(request.GetAddress(), nil)
	case driverModel.ApduDataOther:
		read, ok := request.GetExtendedApdu().(driverModel.ApduDataExtPropertyValueRead)
		if !ok {
			g.violation("unexpected extended apdu %T", request.GetExtendedApdu())
			return nil
		}
		value, ok := device.properties[[2]uint8{read.GetObjectIndex(), read.GetPropertyId()}]
		if !ok || read.GetIndex() != 1 || read.GetCount() != 1 {
			return driverModel.NewApduDataOther(driverModel.NewApduDataExtPropertyValueResponse(read.GetObjectIndex(), read.GetPropertyId(), 0, read.GetIndex(), nil))
		}
		return driverModel.NewApduDataOther(driverModel.NewApduDataExtPropertyValueResponse(read.GetObjectIndex(), read.GetPropertyId(), 1, 1, value))
	default:
		g.violation("unexpected apdu %T", request)
		return nil
	}
}

// sendLoop passes the frames of the line to the tunnel one at a time, each acknowledged before
// the next, as the tunnelling protocol demands.
func (g *fakeKnxGateway) sendLoop() {
	var seq uint8
	for {
		var cemi driverModel.CEMI
		select {
		case <-g.done:
			return
		case cemi = <-g.outbox:
		}
		g.mu.Lock()
		client := g.client
		g.mu.Unlock()
		if client == nil {
			continue
		}
		g.send(client, driverModel.NewTunnelingRequest(driverModel.NewTunnelingRequestDataBlock(g.channel, seq), cemi))
		timeout := time.NewTimer(time.Second)
	wait:
		for {
			select {
			case <-g.done:
				timeout.Stop()
				return
			case ack := <-g.acks:
				if ack == seq {
					break wait
				}
			case <-timeout.C:
				g.mu.Lock()
				g.violation("TUNNELING_REQUEST %d (%T) was never acknowledged", seq, cemi)
				g.mu.Unlock()
				break wait
			}
		}
		timeout.Stop()
		seq++
	}
}

// connectToFakeGateway opens a plc4go connection to the gateway over the real udp transport and
// closes it when the test ends.
func connectToFakeGateway(t *testing.T, ctx context.Context, g *fakeKnxGateway) *Connection {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	transportInstance, err := udp.NewTransport(_options...).CreateTransportInstance(g.url(), map[string][]string{}, _options...)
	require.NoError(t, err)
	connection := NewConnection(transportInstance, map[string][]string{}, NewTagHandler(), _options...)
	require.NoError(t, connection.Connect(ctx))
	t.Cleanup(func() {
		if err := connection.Close(); err != nil && !errors.Is(err, context.Canceled) {
			t.Logf("closing the connection: %v", err)
		}
		_ = connection.messageCodec.Disconnect()
	})
	return connection
}
