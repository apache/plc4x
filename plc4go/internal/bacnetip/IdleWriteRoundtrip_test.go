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

// delayedWriteDevice answers WriteProperty with a SimpleAck after a fixed delay,
// mimicking the ~300ms round-trip observed in the field (instant replies in the
// other fakes cannot exercise the delayed-arrival path).
type delayedWriteDevice struct {
	conn  *net.UDPConn
	wg    sync.WaitGroup
	log   zerolog.Logger
	delay time.Duration
	ctx   context.Context
}

func startDelayedWriteDevice(t *testing.T, log zerolog.Logger, delay time.Duration) *delayedWriteDevice {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	require.NoError(t, err)
	d := &delayedWriteDevice{conn: conn, log: log, delay: delay, ctx: t.Context()}
	d.wg.Add(1)
	go d.serve()
	return d
}

func (d *delayedWriteDevice) port() int { return d.conn.LocalAddr().(*net.UDPAddr).Port }
func (d *delayedWriteDevice) stop()     { _ = d.conn.Close(); d.wg.Wait() }

func (d *delayedWriteDevice) serve() {
	defer d.wg.Done()
	buf := make([]byte, 4096)
	for {
		n, src, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		data := make([]byte, n)
		copy(data, buf[:n])
		bvlc, err := model.BVLCParse[model.BVLC](d.ctx, data)
		if err != nil {
			continue
		}
		npduRetriever, ok := bvlc.(interface{ GetNpdu() model.NPDU })
		if !ok {
			continue
		}
		cr, ok := npduRetriever.GetNpdu().GetApdu().(model.APDUConfirmedRequest)
		if !ok {
			continue
		}
		invokeId := cr.GetInvokeId()
		choice := cr.GetServiceRequest().GetServiceChoice()
		go func() {
			time.Sleep(d.delay)
			var resp model.BVLC
			if choice == model.BACnetConfirmedServiceChoice_READ_PROPERTY {
				resp = buildReadAckFor(invokeId)
			} else {
				resp = wrapAPDU(model.NewAPDUSimpleAck(invokeId, choice), false)
			}
			theBytes, err := resp.Serialize()
			if err != nil {
				return
			}
			_, _ = d.conn.WriteToUDP(theBytes, src)
		}()
	}
}

// TestNativeBacnetWrite_DelayedReplyAfterIdle combines both field conditions: a
// multi-second idle gap AND a ~400ms device reply delay on the write.
func TestNativeBacnetWrite_DelayedReplyAfterIdle(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startDelayedWriteDevice(t, log, 400*time.Millisecond)
	t.Cleanup(device.stop)

	conn := newWriteTestConnection(t, log, device.port())

	code, val, ok := readPresentValue(t, conn)
	require.True(t, ok, "priming read hung")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.InDelta(t, 23.5, val, 0.001)

	time.Sleep(3 * time.Second)

	code, ok = writePresentValue(t, conn)
	require.True(t, ok, "delayed write after idle gap did not resolve")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
}

// TestNativeBacnetWrite_AfterIdleGap reproduces the real field ordering: a read
// succeeds, then the connection sits idle for several seconds (the receive loop
// keeps polling with expiring read deadlines), then a WriteProperty is issued.
// The SimpleAck must still resolve the Execute channel.
func TestNativeBacnetWrite_AfterIdleGap(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	t.Cleanup(device.stop)

	conn := newWriteTestConnection(t, log, device.port())

	// One read to prime the connection (mirrors discovery/poll).
	code, val, ok := readPresentValue(t, conn)
	require.True(t, ok, "priming read hung")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.InDelta(t, 23.5, val, 0.001)

	// Idle gap well beyond the 10ms per-poll read deadline.
	time.Sleep(3 * time.Second)

	code, ok = writePresentValue(t, conn)
	t.Logf("device received %d request(s)", device.requestCount())
	require.True(t, ok, "write after idle gap did not complete (Execute never resolved)")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
}

// TestNativeBacnetWrite_AfterIdleViaCache mirrors the real ugagent path: the
// connection is leased from a PlcConnectionCache with a short max-idle-time, a
// read runs, the lease is returned and sits idle past maxIdleTime (so the cache
// proactively reconnects on the next lease, opening a NEW ephemeral socket),
// then a write is issued. The SimpleAck must still resolve.
func TestNativeBacnetWrite_AfterIdleViaCache(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	t.Cleanup(device.stop)

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)
	connCache := cache.NewPlcConnectionCache(dm,
		cache.WithCustomLogger(log),
		cache.WithMaxIdleTime(500*time.Millisecond),
	)
	t.Cleanup(func() { assert.NoError(t, connCache.Close()) })

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000", device.port())

	// Lease #1: a read.
	{
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		conn, err := connCache.GetConnection(ctx, connStr)
		require.NoError(t, err)
		code, val, ok := readPresentValue(t, conn)
		require.True(t, ok, "priming read via cache hung")
		assert.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.InDelta(t, 23.5, val, 0.001)
		require.NoError(t, conn.Close()) // return lease to cache -> becomes idle
	}

	// Idle past maxIdleTime so the next lease triggers a reconnect (new socket).
	time.Sleep(2 * time.Second)

	// Lease #2: a write on the (reconnected) connection.
	{
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		conn, err := connCache.GetConnection(ctx, connStr)
		require.NoError(t, err)
		code, ok := writePresentValue(t, conn)
		t.Logf("device received %d request(s)", device.requestCount())
		require.NoError(t, conn.Close())
		require.True(t, ok, "write after idle-reconnect via cache did not complete (Execute never resolved)")
		assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	}
}
