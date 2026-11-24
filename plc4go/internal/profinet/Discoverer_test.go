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

// The simulated devices answer from the DCP frame layout of IEC 61158-6-10, written out byte by
// byte, so passing means the generated model reads real Identify responses, not just its own.

import (
	"context"
	"encoding/binary"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
	"github.com/apache/plc4x/plc4go/spi/transports/virtualeth"
)

type simulatedDevice struct {
	mac            net.HardwareAddr
	name, typeName string
	vendor, device uint16
	ip, mask       net.IP
	// tagged answers with an 802.1Q tag, as some devices do.
	tagged bool
	// wrongXid answers with a transaction id that isn't the request's.
	wrongXid bool
}

// dcpBlock is option, suboption, length, block info, data, padded to an even length.
func dcpBlock(option, suboption byte, blockInfo uint16, data []byte) []byte {
	block := []byte{option, suboption}
	block = binary.BigEndian.AppendUint16(block, uint16(2+len(data)))
	block = binary.BigEndian.AppendUint16(block, blockInfo)
	block = append(block, data...)
	if len(block)%2 != 0 {
		block = append(block, 0)
	}
	return block
}

func (s simulatedDevice) response(destination net.HardwareAddr, xid uint32) []byte {
	var blocks []byte
	blocks = append(blocks, dcpBlock(0x02, 0x01, 0, []byte(s.typeName))...)
	blocks = append(blocks, dcpBlock(0x02, 0x02, 0, []byte(s.name))...)
	blocks = append(blocks, dcpBlock(0x02, 0x03, 0, binary.BigEndian.AppendUint16(binary.BigEndian.AppendUint16(nil, s.vendor), s.device))...)
	blocks = append(blocks, dcpBlock(0x02, 0x04, 0, []byte{0x01, 0x00})...) // IO device
	ipData := append(append(append([]byte{}, s.ip.To4()...), s.mask.To4()...), 0, 0, 0, 0)
	blocks = append(blocks, dcpBlock(0x01, 0x02, 0x0001, ipData)...)
	if s.wrongXid {
		xid++
	}
	pdu := []byte{0xFE, 0xFF, 0x05, 0x01}
	pdu = binary.BigEndian.AppendUint32(pdu, xid)
	pdu = binary.BigEndian.AppendUint16(pdu, 0)
	pdu = binary.BigEndian.AppendUint16(pdu, uint16(len(blocks)))
	pdu = append(pdu, blocks...)
	frame := append(append([]byte{}, destination...), s.mac...)
	if s.tagged {
		frame = append(frame, 0x81, 0x00, 0x00, 0x00)
	}
	frame = append(frame, 0x88, 0x92)
	frame = append(frame, pdu...)
	for len(frame) < 60 {
		frame = append(frame, 0)
	}
	return frame
}

// dcpOffset is where the DCP PDU starts: after the EtherType, which follows an 802.1Q tag if the
// frame still carries one (a Linux receiver with VLAN offload gets it stripped).
func dcpOffset(frame []byte) int {
	if len(frame) >= 18 && binary.BigEndian.Uint16(frame[12:]) == 0x8100 {
		return 18
	}
	return 14
}

// recorded is the requests the simulated devices saw.
type recorded struct {
	mutex    sync.Mutex
	requests [][]byte
}

func (r *recorded) all() [][]byte {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return append([][]byte(nil), r.requests...)
}

// runDevices answers Identify requests on a segment and records the requests.
func runDevices(t *testing.T, segment string, devices ...simulatedDevice) *recorded {
	t.Helper()
	requests := &recorded{}
	for _, device := range devices {
		instance := virtualeth.NewTransportInstance(segment, rawsocket.Filter{EtherTypes: []uint16{etherTypeProfinet}}, device.mac)
		require.NoError(t, instance.Connect(t.Context()))
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		go func() {
			defer close(done)
			for ctx.Err() == nil {
				frame, err := readFrame(ctx, instance)
				if err != nil {
					continue
				}
				requests.mutex.Lock()
				requests.requests = append(requests.requests, frame)
				requests.mutex.Unlock()
				// Request: FE FE 05 00 xid(4) responseDelay(2) length(2).
				pdu := frame[dcpOffset(frame):]
				if len(pdu) < 8 || binary.BigEndian.Uint16(pdu) != 0xFEFE {
					continue
				}
				record, err := rawsocket.EncodeRecord(device.response(frame[6:12], binary.BigEndian.Uint32(pdu[4:])))
				require.NoError(t, err)
				require.NoError(t, instance.Write(ctx, record))
			}
		}()
		t.Cleanup(func() {
			cancel()
			<-done
			_ = instance.Close()
		})
	}
	return requests
}

func testDiscoverer(t *testing.T, segment string) *Discoverer {
	d := NewDiscoverer(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	d.openStation = func(interfaceName string) (station, error) {
		return virtualeth.NewTransportInstance(interfaceName, rawsocket.Filter{EtherTypes: []uint16{etherTypeProfinet}}, nil), nil
	}
	d.interfaces = func() ([]string, error) { return []string{segment}, nil }
	return d
}

func TestDiscoverIdentifiesDevices(t *testing.T) {
	segment := t.Name()
	requests := runDevices(t, segment,
		simulatedDevice{mac: net.HardwareAddr{0x00, 0x0E, 0x8C, 0x01, 0x02, 0x03}, name: "et200sp-1", typeName: "ET 200SP", vendor: 0x002A, device: 0x0313, ip: net.IPv4(192, 168, 0, 10), mask: net.IPv4(255, 255, 255, 0)},
		simulatedDevice{mac: net.HardwareAddr{0x00, 0xD0, 0xC9, 0x0A, 0x0B, 0x0C}, name: "adam-6150pn", typeName: "ADAM6150", vendor: 0x01A4, device: 0x0001, ip: net.IPv4zero, mask: net.IPv4zero, tagged: true},
		simulatedDevice{mac: net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x99}, name: "stray", typeName: "X", ip: net.IPv4zero, mask: net.IPv4zero, wrongXid: true},
	)
	d := testDiscoverer(t, segment)

	var mutex sync.Mutex
	items := map[string]apiModel.PlcDiscoveryItem{}
	require.NoError(t, d.scan(t.Context(), segment, 300*time.Millisecond, func(item apiModel.PlcDiscoveryItem) {
		mutex.Lock()
		defer mutex.Unlock()
		items[item.GetAttributes()["deviceName"].GetString()] = item
	}))

	require.Len(t, items, 2, "the answer with the wrong transaction id is ignored")
	et := items["et200sp-1"]
	require.NotNil(t, et)
	assert.Equal(t, "profinet", et.GetProtocolCode())
	assert.Equal(t, "raw-socket", et.GetTransportCode())
	transportUrl := et.GetTransportUrl()
	assert.Equal(t, "raw-socket://"+segment, transportUrl.String())
	assert.Equal(t, "ET 200SP - et200sp-1", et.GetName())
	assert.Equal(t, map[string][]string{"device-name": {"et200sp-1"}, "device-mac": {"00:0e:8c:01:02:03"}, "device-ip": {"192.168.0.10"}}, et.GetOptions())
	attributes := et.GetAttributes()
	assert.Equal(t, "192.168.0.10", attributes["ipAddress"].GetString())
	assert.Equal(t, "255.255.255.0", attributes["subnetMask"].GetString())
	assert.Equal(t, "002A", attributes["vendorId"].GetString())
	assert.Equal(t, "0313", attributes["deviceId"].GetString())
	assert.Equal(t, "DEVICE", attributes["role"].GetString())
	assert.Equal(t, "dcp", attributes["packetType"].GetString())

	adam := items["adam-6150pn"]
	require.NotNil(t, adam, "a VLAN tagged answer is understood")
	assert.NotContains(t, adam.GetOptions(), "device-ip", "a device without an IP address has none to connect with yet")

	// The request on the wire: multicast, VLAN priority 0, Identify All.
	all := requests.all()
	require.NotEmpty(t, all)
	request := all[0]
	assert.Equal(t, []byte{0x01, 0x0E, 0xCF, 0x00, 0x00, 0x00}, request[0:6])
	assert.Equal(t, []byte{0x81, 0x00, 0x00, 0x00, 0x88, 0x92}, request[12:18])
	pdu := request[18:]
	assert.Equal(t, []byte{0xFE, 0xFE, 0x05, 0x00}, pdu[0:4], "frame id, Identify, request")
	assert.Equal(t, uint16(identifyResponseDelay), binary.BigEndian.Uint16(pdu[8:]))
	assert.Equal(t, uint16(4), binary.BigEndian.Uint16(pdu[10:]), "one 4 byte block")
	assert.Equal(t, []byte{0xFF, 0xFF, 0x00, 0x00}, pdu[12:16], "the All selector")
	assert.GreaterOrEqual(t, len(request), 60)
}

func TestDiscoverOptions(t *testing.T) {
	d := testDiscoverer(t, "default")
	names, window, err := d.plan([]options.WithDiscoveryOption{
		options.WithDiscoveryOptionDeviceName("eth1"),
		options.WithDiscoveryOptionDeviceName("eth2"),
		options.WithDiscoveryOptionProtocolSpecific(DiscoveryWindowOption, 1500),
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"eth1", "eth2"}, names)
	assert.Equal(t, 1500*time.Millisecond, window)
	names, window, err = d.plan(nil)
	require.NoError(t, err)
	assert.Equal(t, []string{"default"}, names)
	assert.Equal(t, defaultDiscoveryWindow, window)
}

func TestDiscoverRunsInBackground(t *testing.T) {
	segment := t.Name()
	runDevices(t, segment, simulatedDevice{mac: net.HardwareAddr{0x00, 0x0E, 0x8C, 0x00, 0x00, 0x01}, name: "one", typeName: "T", ip: net.IPv4zero, mask: net.IPv4zero})
	d := testDiscoverer(t, segment)
	found := make(chan string, 1)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	start := time.Now()
	require.NoError(t, d.Discover(ctx, func(item apiModel.PlcDiscoveryItem) {
		found <- item.GetAttributes()["deviceName"].GetString()
	}, options.WithDiscoveryOptionProtocolSpecific(DiscoveryWindowOption, 500)))
	assert.Less(t, time.Since(start), 100*time.Millisecond, "Discover returns at once")
	select {
	case name := <-found:
		assert.Equal(t, "one", name)
	case <-ctx.Done():
		t.Fatal("no device found")
	}
	require.NoError(t, d.Close())
}
