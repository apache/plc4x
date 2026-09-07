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

package serial

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
)

func acquireTwo(t *testing.T, registry *sharedPortRegistry, slavePath string, cfg sharedPortConfig) (serialport.Port, serialport.Port) {
	t.Helper()
	first, err := registry.acquire(slavePath, cfg)
	require.NoError(t, err)
	second, err := registry.acquire(slavePath, cfg)
	require.NoError(t, err)
	return first, second
}

func TestSharedPort_BroadcastToAllSubscribers(t *testing.T) {
	master, slavePath := openPTY(t)
	registry := newSharedPortRegistry(zerolog.Nop())
	cfg := sharedPortConfig{port: serialport.Config{BaudRate: 9600}}
	first, second := acquireTwo(t, registry, slavePath, cfg)
	defer first.Close()
	defer second.Close()

	_, err := master.Write([]byte("hello"))
	require.NoError(t, err)

	for _, sub := range []serialport.Port{first, second} {
		require.NoError(t, sub.SetReadDeadline(time.Now().Add(5*time.Second)))
		buf := make([]byte, 16)
		n, err := sub.Read(buf)
		require.NoError(t, err)
		assert.Equal(t, "hello", string(buf[:n]), "every subscriber sees the full traffic")
	}
}

func TestSharedPort_WritesSerializedAndPaced(t *testing.T) {
	master, slavePath := openPTY(t)
	registry := newSharedPortRegistry(zerolog.Nop())
	cfg := sharedPortConfig{port: serialport.Config{BaudRate: 9600}, interframeDelay: 60 * time.Millisecond}
	first, second := acquireTwo(t, registry, slavePath, cfg)
	defer first.Close()
	defer second.Close()

	start := time.Now()
	var wg sync.WaitGroup
	for _, sub := range []serialport.Port{first, second} {
		wg.Add(1)
		go func(p serialport.Port) {
			defer wg.Done()
			_, err := p.Write([]byte("XXXX"))
			assert.NoError(t, err)
		}(sub)
	}
	wg.Wait()
	elapsed := time.Since(start)

	buf := make([]byte, 8)
	total := 0
	for total < 8 {
		n, err := master.Read(buf[total:])
		require.NoError(t, err)
		total += n
	}
	assert.Equal(t, "XXXXXXXX", string(buf), "both payloads intact")
	assert.GreaterOrEqual(t, elapsed, 60*time.Millisecond, "second write must wait out the inter-frame gap")
}

func TestSharedPort_ConfigMismatchFails(t *testing.T) {
	_, slavePath := openPTY(t)
	registry := newSharedPortRegistry(zerolog.Nop())
	first, err := registry.acquire(slavePath, sharedPortConfig{port: serialport.Config{BaudRate: 9600}})
	require.NoError(t, err)
	defer first.Close()

	_, err = registry.acquire(slavePath, sharedPortConfig{port: serialport.Config{BaudRate: 19200}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), slavePath)
}

func TestSharedPort_RefcountLifecycle(t *testing.T) {
	master, slavePath := openPTY(t)
	registry := newSharedPortRegistry(zerolog.Nop())
	cfg := sharedPortConfig{port: serialport.Config{BaudRate: 9600}}
	first, second := acquireTwo(t, registry, slavePath, cfg)

	// First close: the second subscriber keeps working.
	require.NoError(t, first.Close())
	_, err := master.Write([]byte{0x01})
	require.NoError(t, err)
	require.NoError(t, second.SetReadDeadline(time.Now().Add(5*time.Second)))
	buf := make([]byte, 1)
	_, err = second.Read(buf)
	require.NoError(t, err)

	// Reads on the closed subscription fail with os.ErrClosed.
	_, err = first.Read(buf)
	require.ErrorIs(t, err, os.ErrClosed)

	// Last close releases the port; a fresh acquire re-opens it.
	require.NoError(t, second.Close())
	third, err := registry.acquire(slavePath, cfg)
	require.NoError(t, err)
	require.NoError(t, third.Close())
}

func TestSharedPort_OverflowDropsOldest(t *testing.T) {
	master, slavePath := openPTY(t)
	registry := newSharedPortRegistry(zerolog.Nop())
	registry.bufferCapacity = 8 // tiny buffer to force overflow
	cfg := sharedPortConfig{port: serialport.Config{BaudRate: 9600}}
	sub, err := registry.acquire(slavePath, cfg)
	require.NoError(t, err)
	defer sub.Close()

	_, err = master.Write([]byte("0123456789ABCDEF"))
	require.NoError(t, err)

	// Wait for delivery, then confirm reads still work and end with the
	// freshest bytes (drop-oldest, never drop-newest).
	require.NoError(t, sub.SetReadDeadline(time.Now().Add(5*time.Second)))
	collected := make([]byte, 0, 16)
	buf := make([]byte, 16)
	for {
		require.NoError(t, sub.SetReadDeadline(time.Now().Add(200*time.Millisecond)))
		n, err := sub.Read(buf)
		if err != nil {
			break // deadline: stream drained
		}
		collected = append(collected, buf[:n]...)
	}
	require.NotEmpty(t, collected)
	assert.Equal(t, byte('F'), collected[len(collected)-1], "newest byte must survive")
	assert.LessOrEqual(t, len(collected), 16)
	assert.Positive(t, sub.(*subscription).ring.droppedTotal(), "drop counter must have risen")
}

func TestSharedPort_SubscriptionIsNotControlPort(t *testing.T) {
	_, slavePath := openPTY(t)
	registry := newSharedPortRegistry(zerolog.Nop())
	sub, err := registry.acquire(slavePath, sharedPortConfig{port: serialport.Config{BaudRate: 9600}})
	require.NoError(t, err)
	defer sub.Close()
	_, isControl := sub.(serialport.ControlPort)
	assert.False(t, isControl, "modem control must not leak through shared subscriptions")
}
