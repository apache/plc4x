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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/cache"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiTransports "github.com/apache/plc4x/plc4go/pkg/api/transports"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// TestNativeBacnetWrite_LeaseResetOnLiveSocket covers the one cache path the
// earlier roundtrip tests missed: a cache WITHOUT maxIdleTime (the production
// default) re-leases the SAME live connection, and every successful lease runs
// TransportInstance.Reset() — deadline poke, drain read, bufio.Reader swap —
// concurrently with the codec's 10ms receive-poll worker on the same socket.
// The write's SimpleAck (arriving ~1ms after TX, like the field pcap) must
// still resolve Execute on every iteration.
//
// Field shape (CI run 29414550340): read lease, return, idle gap, write lease
// (Reset on live socket), TX +300ms, ack +0.8ms — ack visible on the host
// interface with a good checksum but never surfaced by the codec.
func TestNativeBacnetWrite_LeaseResetOnLiveSocket(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeWriteDevice(t, log)
	t.Cleanup(device.stop)

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)
	// No WithMaxIdleTime: the connection is never recycled, so every lease
	// after the first resets the LIVE transport instead of reconnecting.
	connCache := cache.NewPlcConnectionCache(dm, cache.WithCustomLogger(log))
	t.Cleanup(func() { assert.NoError(t, connCache.Close()) })

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=3000", device.port())

	// Lease #1: priming read (creates the connection; scan analogue).
	{
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		conn, err := connCache.GetConnection(ctx, connStr)
		cancel()
		require.NoError(t, err)
		code, val, ok := readPresentValue(t, conn)
		require.True(t, ok, "priming read via cache hung")
		assert.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.InDelta(t, 23.5, val, 0.001)
		require.NoError(t, conn.Close())
	}

	// Short idle, then many write leases. Each lease runs Reset() against the
	// live socket while the receive worker polls — every iteration is one
	// spin of the race wheel.
	time.Sleep(500 * time.Millisecond)
	for i := 0; i < 50; i++ {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		conn, err := connCache.GetConnection(ctx, connStr)
		cancel()
		require.NoError(t, err, "lease %d failed", i)
		code, ok := writePresentValue(t, conn)
		require.NoError(t, conn.Close())
		require.True(t, ok, "iteration %d: write did not complete (Execute never resolved) — SimpleAck lost around lease-time Reset", i)
		assert.Equal(t, apiModel.PlcResponseCode_OK, code, "iteration %d", i)
	}
}

// TestNativeBacnetWrite_LeaseResetWithIdleGap is the exact field ordering with
// a real idle gap per iteration: read, return, idle >1s, write lease (Reset on
// live socket), delayed device reply. Slower per iteration but closest to the
// failing CI run's shape.
func TestNativeBacnetWrite_LeaseResetWithIdleGap(t *testing.T) {
	if testing.Short() {
		t.Skip("multi-second idle gaps")
	}
	log := testutils.ProduceTestingLogger(t)
	device := startDelayedWriteDevice(t, log, 300*time.Millisecond)
	t.Cleanup(device.stop)

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)
	connCache := cache.NewPlcConnectionCache(dm, cache.WithCustomLogger(log))
	t.Cleanup(func() { assert.NoError(t, connCache.Close()) })

	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=5000", device.port())

	{
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		conn, err := connCache.GetConnection(ctx, connStr)
		cancel()
		require.NoError(t, err)
		code, val, ok := readPresentValue(t, conn)
		require.True(t, ok, "priming read via cache hung")
		assert.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.InDelta(t, 23.5, val, 0.001)
		require.NoError(t, conn.Close())
	}

	for i := 0; i < 3; i++ {
		time.Sleep(1500 * time.Millisecond) // idle gap on the returned lease
		ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
		conn, err := connCache.GetConnection(ctx, connStr)
		cancel()
		require.NoError(t, err, "lease %d failed", i)
		code, ok := writePresentValue(t, conn)
		require.NoError(t, conn.Close())
		require.True(t, ok, "iteration %d: write after idle gap did not complete", i)
		assert.Equal(t, apiModel.PlcResponseCode_OK, code, "iteration %d", i)
	}
}
