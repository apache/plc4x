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

package bacnetip

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

func TestNewDiscoverer_DefaultTimeout(t *testing.T) {
	d := NewDiscoverer()
	assert.Equal(t, 5*time.Second, d.discoveryTimeout)
}

func TestSetDiscoveryTimeout_Override(t *testing.T) {
	d := NewDiscoverer()
	d.SetDiscoveryTimeout(2 * time.Second)
	assert.Equal(t, 2*time.Second, d.discoveryTimeout)
}

func TestSetDiscoveryTimeout_ZeroFallsBackToDefault(t *testing.T) {
	d := NewDiscoverer()
	d.SetDiscoveryTimeout(0)
	assert.Equal(t, 5*time.Second, d.discoveryTimeout)
}

func TestSetDiscoveryTimeout_NegativeFallsBackToDefault(t *testing.T) {
	d := NewDiscoverer()
	d.SetDiscoveryTimeout(-1 * time.Second)
	assert.Equal(t, 5*time.Second, d.discoveryTimeout)
}

func TestHandleIncomingBVLCs_DispatchesIAm(t *testing.T) {
	// Real IAm captured from a device advertising instance 3001 (vendor 999):
	// BVLC(Original-Broadcast-NPDU) / NPDU / APDU(unconfirmed, IAm).
	iamBytes := []byte{
		0x81, 0x0b, 0x00, 0x15, // BVLC
		0x01, 0x00, // NPDU
		0x10, 0x00, // APDU unconfirmed, service IAm
		0xc4, 0x02, 0x00, 0x0b, 0xb9, // object-id: device 3001
		0x22, 0x05, 0xc4, // max-apdu 1476
		0x91, 0x00, // segmentation: both
		0x22, 0x03, 0xe7, // vendor 999
	}
	bvlc, err := driverModel.BVLCParse[driverModel.BVLC](context.Background(), iamBytes)
	require.NoError(t, err)

	d := NewDiscoverer()
	ch := make(chan receivedBvlcMessage, 1)
	ch <- receivedBvlcMessage{bvlc: bvlc, addr: &net.UDPAddr{IP: net.IPv4(192, 168, 100, 2), Port: 47808}}

	got := make(chan apiModel.PlcDiscoveryItem, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan struct{})
	go func() {
		d.handleIncomingBVLCs(ctx, func(item apiModel.PlcDiscoveryItem) { got <- item }, ch)
		close(done)
	}()

	select {
	case item := <-got:
		assert.Contains(t, item.GetName(), "3001")
		transportURL := item.GetTransportUrl()
		assert.Equal(t, "192.168.100.2", transportURL.Hostname())
	case <-time.After(2 * time.Second):
		t.Fatal("callback was not invoked for IAm")
	}

	// Cancelling the context must make the handler return (no hang).
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handleIncomingBVLCs did not return after context cancellation")
	}
}

func TestResolveBacnetUDPAddr(t *testing.T) {
	// Host only — default port applies.
	addr, err := resolveBacnetUDPAddr("192.168.1.50", 47808)
	require.NoError(t, err)
	assert.Equal(t, "192.168.1.50", addr.IP.String())
	assert.Equal(t, 47808, addr.Port)

	// Host:port — explicit port wins.
	addr, err = resolveBacnetUDPAddr("10.0.0.5:47809", 47808)
	require.NoError(t, err)
	assert.Equal(t, "10.0.0.5", addr.IP.String())
	assert.Equal(t, 47809, addr.Port)

	// Invalid host.
	_, err = resolveBacnetUDPAddr("not-an-ip", 47808)
	require.Error(t, err)
}

func TestExtractProtocolSpecificOptions_RemoteAddress(t *testing.T) {
	opts := []options.WithDiscoveryOption{
		options.WithDiscoveryOptionProtocolSpecific("remote-address", "192.168.1.50"),
		options.WithDiscoveryOptionProtocolSpecific("bacnet-port", 47808),
	}
	specific, err := extractProtocolSpecificOptions(opts)
	require.NoError(t, err)
	assert.Equal(t, "192.168.1.50", specific.remoteAddress)
}

// iAmFrame is a real IAm of device 3001 (vendor 999) as an Original-Broadcast-NPDU.
var iAmFrame = []byte{
	0x81, 0x0b, 0x00, 0x15, // BVLC
	0x01, 0x00, // NPDU
	0x10, 0x00, // APDU unconfirmed, service IAm
	0xc4, 0x02, 0x00, 0x0b, 0xb9, // object-id: device 3001
	0x22, 0x05, 0xc4, // max-apdu 1476
	0x91, 0x00, // segmentation: both
	0x22, 0x03, 0xe7, // vendor 999
}

// The discovery socket shares the BACnet port, so it gets datagrams meant for other sockets. A BVLC-Result (no
// NPDU), a network layer message and a confirmed request are skipped without a panic, and an IAm behind them is
// still reported.
func TestHandleIncomingBVLCs_SkipsWhatIsNoUnconfirmedRequest(t *testing.T) {
	src := &net.UDPAddr{IP: net.IPv4(192, 168, 100, 2), Port: 47808}
	nlmControl := driverModel.NewNPDUControl(true, false, false, false, driverModel.NPDUNetworkPriority_NORMAL_MESSAGE)
	// ReadProperty(device 1, object-name), a confirmed request
	confirmed, err := driverModel.BVLCParse[driverModel.BVLC](t.Context(), []byte{
		0x81, 0x0a, 0x00, 0x11, // BVLC Original-Unicast-NPDU
		0x01, 0x04, // NPDU, expecting reply
		0x00, 0x05, 0x01, 0x0c, // APDU confirmed, invoke id 1, ReadProperty
		0x0c, 0x02, 0x00, 0x00, 0x01, // object-id: device 1
		0x19, 0x4d, // property: object-name
	})
	require.NoError(t, err)
	iAm, err := driverModel.BVLCParse[driverModel.BVLC](t.Context(), iAmFrame)
	require.NoError(t, err)

	ch := make(chan receivedBvlcMessage, 8)
	ch <- receivedBvlcMessage{bvlc: driverModel.NewBVLCResult(driverModel.BVLCResultCode_SUCCESSFUL_COMPLETION), addr: src}
	ch <- receivedBvlcMessage{bvlc: driverModel.NewBVLCOriginalUnicastNPDU(driverModel.NewNPDU(1, nlmControl, nil, nil, nil, nil, nil, nil, nil, driverModel.NewNLMWhoIsRouterToNetwork(nil), nil)), addr: src}
	ch <- receivedBvlcMessage{bvlc: confirmed, addr: src}
	ch <- receivedBvlcMessage{bvlc: iAm, addr: src}

	got := make(chan apiModel.PlcDiscoveryItem, 4)
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan struct{})
	go func() {
		defer close(done)
		NewDiscoverer().handleIncomingBVLCs(ctx, func(item apiModel.PlcDiscoveryItem) { got <- item }, ch)
	}()
	select {
	case item := <-got:
		assert.Contains(t, item.GetName(), "3001")
	case <-time.After(2 * time.Second):
		t.Fatal("the IAm behind the skipped datagrams was not reported")
	}
	cancel()
	<-done
	assert.Empty(t, got)
}

// A discovery socket's reader skips a datagram that does not parse, and ends when the socket is closed, also
// when nobody takes what it read any more: Close waits for it.
func TestReceiveBVLCs_EndsWithTheSocket(t *testing.T) {
	conn, err := net.ListenPacket("udp4", "127.0.0.1:0")
	require.NoError(t, err)
	sender, err := net.Dial("udp4", conn.LocalAddr().String())
	require.NoError(t, err)
	defer func() { _ = sender.Close() }()

	d := NewDiscoverer()
	ctx, cancel := context.WithCancel(t.Context())
	incoming := make(chan receivedBvlcMessage) // unbuffered: nobody reads after the first message
	d.wg.Go(func() { d.receiveBVLCs(ctx, conn, incoming) })

	_, err = sender.Write([]byte{0x81, 0x0b, 0x00}) // shorter than a BVLC header
	require.NoError(t, err)
	_, err = sender.Write(iAmFrame)
	require.NoError(t, err)
	select {
	case msg := <-incoming:
		_, ok := msg.bvlc.(driverModel.BVLCOriginalBroadcastNPDU)
		assert.True(t, ok, "got %T", msg.bvlc)
	case <-time.After(2 * time.Second):
		t.Fatal("the IAm after the malformed datagram was not handed on")
	}

	// The window ends: nobody reads incoming any more, a late datagram arrives, the socket closes
	cancel()
	_, err = sender.Write(iAmFrame)
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond)
	require.NoError(t, conn.Close())

	closed := make(chan struct{})
	go func() {
		_ = d.Close()
		close(closed)
	}()
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not return: a discovery reader outlived its socket")
	}
}
