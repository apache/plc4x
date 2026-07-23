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
	"github.com/apache/plc4x/plc4go/pkg/api/cache"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiTransports "github.com/apache/plc4x/plc4go/pkg/api/transports"
	model "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// fakeBacnetDevice is a minimal UDP BACnet/IP device used to exercise the real
// driver's read round-trip. On every confirmed request it replies with a
// ReadProperty ComplexAck carrying ANALOG_INPUT,1/PRESENT_VALUE = 23.5, echoing
// the request's invoke id so the Reader's expectation matches.
type fakeBacnetDevice struct {
	conn    *net.UDPConn
	wg      sync.WaitGroup
	log     zerolog.Logger
	gotReqs int
	ctx     context.Context
	mu      sync.Mutex
}

func startFakeBacnetDevice(t *testing.T, log zerolog.Logger) *fakeBacnetDevice {
	t.Helper()
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0}
	conn, err := net.ListenUDP("udp4", addr)
	require.NoError(t, err)
	d := &fakeBacnetDevice{conn: conn, log: log, ctx: t.Context()}
	d.wg.Add(1)
	go d.serve()
	return d
}

func (d *fakeBacnetDevice) port() int { return d.conn.LocalAddr().(*net.UDPAddr).Port }

func (d *fakeBacnetDevice) requestCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.gotReqs
}

func (d *fakeBacnetDevice) serve() {
	defer d.wg.Done()
	buf := make([]byte, 4096)
	for {
		n, src, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			return // socket closed
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		d.log.Info().Int("bytes", n).Stringer("src", src).Msg("fake device received packet")

		invokeId, ok := d.extractInvokeId(data)
		if !ok {
			d.log.Warn().Msg("fake device: could not extract invoke id; ignoring")
			continue
		}
		d.mu.Lock()
		d.gotReqs++
		d.mu.Unlock()

		resp := d.buildReadPropertyAck(invokeId)
		theBytes, err := resp.Serialize()
		if err != nil {
			d.log.Error().Err(err).Msg("fake device: serialize response")
			continue
		}
		if _, err := d.conn.WriteToUDP(theBytes, src); err != nil {
			d.log.Error().Err(err).Msg("fake device: write response")
			continue
		}
		d.log.Info().Uint8("invokeId", invokeId).Int("bytes", len(theBytes)).Stringer("dst", src).Msg("fake device sent ComplexAck")
	}
}

func (d *fakeBacnetDevice) extractInvokeId(data []byte) (uint8, bool) {
	bvlc, err := model.BVLCParse[model.BVLC](d.ctx, data)
	if err != nil {
		d.log.Error().Err(err).Msg("fake device: parse BVLC")
		return 0, false
	}
	npduRetriever, ok := bvlc.(interface{ GetNpdu() model.NPDU })
	if !ok {
		return 0, false
	}
	apdu := npduRetriever.GetNpdu().GetApdu()
	cr, ok := apdu.(model.APDUConfirmedRequest)
	if !ok {
		d.log.Warn().Msgf("fake device: not a confirmed request: %T", apdu)
		return 0, false
	}
	return cr.GetInvokeId(), true
}

func (d *fakeBacnetDevice) buildReadPropertyAck(invokeId uint8) model.BVLC {
	serviceAck := model.NewBACnetServiceAckReadProperty(
		0,
		model.CreateBACnetContextTagObjectIdentifier(0, uint16(model.BACnetObjectType_ANALOG_INPUT), 1),
		model.CreateBACnetPropertyIdentifierTagged(1, uint32(model.BACnetPropertyIdentifier_PRESENT_VALUE)),
		nil,
		constructedDataFromTag(model.CreateBACnetApplicationTagReal(23.5)),
	)
	apdu := model.NewAPDUComplexAck(false, false, invokeId, nil, nil, serviceAck, nil, nil)
	return wrapAPDU(apdu, false, nil)
}

func (d *fakeBacnetDevice) stop() {
	_ = d.conn.Close()
	d.wg.Wait()
}

// readPresentValue issues a single ReadRequest for ANALOG_INPUT,1/PRESENT_VALUE
// against the given connection and returns the response code + float value.
func readPresentValue(t *testing.T, conn plc4go.PlcConnection) (apiModel.PlcResponseCode, float32, bool) {
	t.Helper()
	rr, err := conn.ReadRequestBuilder().AddTagAddress("pv", "ANALOG_INPUT,1/PRESENT_VALUE").Build()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 8*time.Second)
	t.Cleanup(cancel)
	select {
	case <-ctx.Done():
		return 0, 0, false
	case res := <-rr.Execute(ctx):
		if res.GetErr() != nil {
			t.Logf("read error: %v", res.GetErr())
			return 0, 0, false
		}
		resp := res.GetResponse()
		code := resp.GetResponseCode("pv")
		v := resp.GetValue("pv")
		if v != nil && v.IsFloat32() {
			return code, v.GetFloat32(), true
		}
		return code, 0, true
	}
}

func TestNativeBacnetRead_DirectConnection(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeBacnetDevice(t, log)
	t.Cleanup(device.stop)

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000", device.port())
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	conn, err := dm.GetConnection(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, conn.Close())
	})

	code, val, ok := readPresentValue(t, conn)
	t.Logf("device received %d request(s)", device.requestCount())
	require.True(t, ok, "read did not complete (timed out waiting for response)")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.InDelta(t, 23.5, val, 0.001)
}

func TestNativeBacnetRead_ViaCache(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeBacnetDevice(t, log)
	t.Cleanup(device.stop)

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)
	connCache := cache.NewPlcConnectionCache(dm, cache.WithCustomLogger(log))
	t.Cleanup(func() {
		assert.NoError(t, connCache.Close())
	})

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000", device.port())
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	conn, err := connCache.GetConnection(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, conn.Close())
	})

	code, val, ok := readPresentValue(t, conn)
	t.Logf("device received %d request(s)", device.requestCount())
	require.True(t, ok, "read via cache did not complete (timed out waiting for response)")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.InDelta(t, 23.5, val, 0.001)
}
