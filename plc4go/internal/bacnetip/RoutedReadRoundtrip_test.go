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
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiTransports "github.com/apache/plc4x/plc4go/pkg/api/transports"
	model "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// fakeBacnetRouter plays a BACnet network-layer router: it requires every
// confirmed request to carry the expected destination specifier (DNET/DADR)
// and answers with a ComplexAck whose NPDU carries the routed device's
// SOURCE specifier — exactly what a real router delivers back (ASHRAE 135
// clause 6.2.4).
type fakeBacnetRouter struct {
	conn         *net.UDPConn
	wg           sync.WaitGroup
	log          zerolog.Logger
	t            *testing.T
	wantDNET     uint16
	wantDADR     []uint8
	badRequests  int
	goodRequests int
	mu           sync.Mutex
}

func startFakeBacnetRouter(t *testing.T, log zerolog.Logger, dnet uint16, dadr []uint8) *fakeBacnetRouter {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	require.NoError(t, err)
	r := &fakeBacnetRouter{conn: conn, log: log, t: t, wantDNET: dnet, wantDADR: dadr}
	r.wg.Add(1)
	go r.serve()
	return r
}

func (r *fakeBacnetRouter) port() int { return r.conn.LocalAddr().(*net.UDPAddr).Port }

func (r *fakeBacnetRouter) counts() (good, bad int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.goodRequests, r.badRequests
}

func (r *fakeBacnetRouter) serve() {
	defer r.wg.Done()
	buf := make([]byte, 4096)
	for {
		n, src, err := r.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		data := make([]byte, n)
		copy(data, buf[:n])

		bvlc, err := model.BVLCParse[model.BVLC](r.t.Context(), data)
		if err != nil {
			continue
		}
		npdu := bvlc.(interface{ GetNpdu() model.NPDU }).GetNpdu()
		cr, ok := npdu.GetApdu().(model.APDUConfirmedRequest)
		if !ok {
			continue
		}

		// A router can only forward specifier-carrying NPDUs.
		dnet := npdu.GetDestinationNetworkAddress()
		if !npdu.GetControl().GetDestinationSpecified() || dnet == nil || *dnet != r.wantDNET {
			r.mu.Lock()
			r.badRequests++
			r.mu.Unlock()
			r.log.Warn().Msg("fake router: request without expected destination specifier — dropping (as a real router would)")
			continue
		}
		r.mu.Lock()
		r.goodRequests++
		r.mu.Unlock()

		// Reply as the routed device: ComplexAck with SOURCE specifier.
		serviceAck := model.NewBACnetServiceAckReadProperty(
			0,
			model.CreateBACnetContextTagObjectIdentifier(0, uint16(model.BACnetObjectType_ANALOG_INPUT), 1),
			model.CreateBACnetPropertyIdentifierTagged(1, uint32(model.BACnetPropertyIdentifier_PRESENT_VALUE)),
			nil,
			constructedDataFromTag(model.CreateBACnetApplicationTagReal(42.25)),
		)
		apdu := model.NewAPDUComplexAck(false, false, cr.GetInvokeId(), nil, nil, serviceAck, nil, nil)
		control := model.NewNPDUControl(false, false, true, false, model.NPDUNetworkPriority_NORMAL_MESSAGE)
		snet := r.wantDNET
		slen := uint8(len(r.wantDADR))
		replyNpdu := model.NewNPDU(1, control, nil, nil, nil, &snet, &slen, r.wantDADR, nil, nil, apdu)
		reply, err := model.NewBVLCOriginalUnicastNPDU(replyNpdu).Serialize()
		if err != nil {
			r.log.Error().Err(err).Msg("fake router: serialize reply")
			continue
		}
		if _, err := r.conn.WriteToUDP(reply, src); err != nil {
			return
		}
	}
}

func (r *fakeBacnetRouter) stop() {
	_ = r.conn.Close()
	r.wg.Wait()
}

// TestNativeBacnetRead_RoutedViaRouter pins the routed round trip end to end
// at driver level: a connection with RemoteNetwork/RemoteAddress emits
// destination-specified requests, and the source-specified reply a router
// delivers back must satisfy the read.
func TestNativeBacnetRead_RoutedViaRouter(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	dadr := []uint8{192, 168, 102, 2, 0xBA, 0xC0}
	router := startFakeBacnetRouter(t, log, 3001, dadr)
	t.Cleanup(router.stop)

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000&RemoteNetwork=3001&RemoteAddress=192.168.102.2:47808", router.port())
	connCtx, connCancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(connCancel)
	conn, err := dm.GetConnection(connCtx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, conn.Close())
	})

	code, val, ok := readPresentValue(t, conn)
	good, bad := router.counts()
	t.Logf("router saw %d routed request(s), %d unroutable", good, bad)
	assert.Zero(t, bad, "driver must never send specifier-less requests on a routed connection")
	require.True(t, ok, "routed read did not complete (reply rejected or timed out)")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.InDelta(t, 42.25, val, 0.001)
}
