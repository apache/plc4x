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

package tcp

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// connectedTestInstance dials a local listener and returns the transport
// instance plus the server-side conn for injecting bytes.
func connectedTestInstance(t *testing.T) (*TransportInstance, net.Conn) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = listener.Close() })

	serverConnCh := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		serverConnCh <- conn
	}()

	addr := listener.Addr().(*net.TCPAddr)
	ti := NewTcpTransportInstance(addr, 1, nil, testutils.EnrichOptionsWithOptionsForTesting(t)...)
	require.NoError(t, ti.Connect(context.Background()))
	t.Cleanup(func() { _ = ti.Close() })

	serverConn := <-serverConnCh
	t.Cleanup(func() { _ = serverConn.Close() })
	return ti, serverConn
}

// waitBuffered polls until the instance's reader has buffered n bytes.
func waitBuffered(t *testing.T, ti *TransportInstance, n uint32) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		avail, err := ti.GetNumBytesAvailableInBuffer()
		require.NoError(t, err)
		if avail >= n {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %d buffered bytes", n)
}

// A frame the peer already delivered - and that the receive path already
// buffered - must survive a lease-time Reset. The old Reset swapped m.reader,
// silently discarding buffered bytes: on a stream that both loses the reply
// and desyncs the framing for everything after it (the TCP analogue of the
// BACnet SimpleAck loss fixed in the UDP transport).
func TestTransportInstance_ResetKeepsBufferedBytes(t *testing.T) {
	ti, serverConn := connectedTestInstance(t)

	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x01, 0x03, 0x02, 0x2a}
	_, err := serverConn.Write(frame)
	require.NoError(t, err)
	waitBuffered(t, ti, uint32(len(frame)))

	ti.Reset() // cache-lease analogue

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	got, err := ti.Read(ctx, uint32(len(frame)))
	require.NoError(t, err, "buffered frame must survive Reset")
	assert.Equal(t, frame, got, "frame bytes must be intact after Reset")
}

// Reset must not race the receive worker. The worker side is modelled the way
// DefaultCodec drives it: a tight GetNumBytesAvailableInBuffer/Read poll loop,
// while another goroutine (the connection cache granting a lease) calls
// Reset(). Run with -race.
func TestTransportInstance_ResetDoesNotRaceReceiveWorker(t *testing.T) {
	ti, serverConn := connectedTestInstance(t)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Peer keeps trickling bytes so the worker's reads have real traffic.
	wg.Go(func() {
		payload := []byte{0xAA, 0xBB, 0xCC, 0xDD}
		for {
			select {
			case <-stop:
				return
			default:
				_, _ = serverConn.Write(payload)
				time.Sleep(time.Millisecond)
			}
		}
	})

	// Receive worker analogue.
	wg.Go(func() {
		for {
			select {
			case <-stop:
				return
			default:
				if avail, err := ti.GetNumBytesAvailableInBuffer(); err == nil && avail > 0 {
					ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
					_, _ = ti.Read(ctx, avail)
					cancel()
				}
			}
		}
	})

	// Cache-lease analogue: Reset storm from a foreign goroutine.
	for range 200 {
		ti.Reset()
		time.Sleep(500 * time.Microsecond)
	}

	close(stop)
	wg.Wait()
}
