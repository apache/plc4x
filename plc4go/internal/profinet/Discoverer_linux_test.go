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

// DCP discovery over real interfaces, when PROFINET_TEST_PAIR names a connected pair (a veth pair):
// the discoverer scans the first, a simulated device answers on the second.

import (
	"context"
	"encoding/binary"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

func TestDiscoverOverRawSockets(t *testing.T) {
	pair := strings.Split(os.Getenv("PROFINET_TEST_PAIR"), ",")
	if len(pair) != 2 {
		t.Skip("set PROFINET_TEST_PAIR to two connected interfaces (a veth pair) to run")
	}
	deviceSide := rawsocket.NewTransportInstance(pair[1], rawsocket.Filter{EtherTypes: []uint16{etherTypeProfinet}}, rawsocket.NewTransport())
	require.NoError(t, deviceSide.Connect(t.Context()))
	defer func() { _ = deviceSide.Close() }()
	device := simulatedDevice{mac: deviceSide.HardwareAddress(), name: "veth-device", typeName: "Simulated", vendor: 0x1234, device: 0x5678, ip: net.IPv4(10, 0, 0, 2), mask: net.IPv4(255, 0, 0, 0)}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go func() {
		for ctx.Err() == nil {
			frame, err := readFrame(ctx, deviceSide)
			if err != nil {
				continue
			}
			pdu := frame[dcpOffset(frame):]
			if len(pdu) < 8 || binary.BigEndian.Uint16(pdu) != 0xFEFE {
				continue
			}
			record, _ := rawsocket.EncodeRecord(device.response(frame[6:12], binary.BigEndian.Uint32(pdu[4:])))
			_ = deviceSide.Write(ctx, record)
		}
	}()

	d := NewDiscoverer(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	var items []apiModel.PlcDiscoveryItem
	require.NoError(t, d.scan(t.Context(), pair[0], time.Second, func(item apiModel.PlcDiscoveryItem) {
		items = append(items, item)
	}))
	require.Len(t, items, 1)
	assert.Equal(t, "Simulated - veth-device", items[0].GetName())
	assert.Equal(t, device.mac.String(), items[0].GetAttributes()["macAddress"].GetString())
	assert.Equal(t, []string{"10.0.0.2"}, items[0].GetOptions()["device-ip"])
}
