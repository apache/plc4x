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

package udp

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ackBytes is the exact 9-byte BACnet/IP SimpleAck observed in the field
// (81 0a 00 09 ...). Bytes [2],[3] are the BVLC length = 0x0009.
var ackBytes = []byte{0x81, 0x0a, 0x00, 0x09, 0x01, 0x00, 0x20, 0x06, 0x0f}

// pollOneLikeCodec drives the transport exactly like bacnetip.MessageCodec.Receive:
// GetNumBytesAvailableInBuffer (which arms a fresh short read deadline), then
// PeekReadableBytes/Read with a *deadline-less* background context (the codec
// polls with m.ctx, which is context.Background()-derived). It loops on the
// same ~10ms cadence as the default codec ReceiveWork loop.
func pollOneLikeCodec(t *testing.T, ti *TransportInstance, overall time.Duration) ([]byte, bool) {
	t.Helper()
	deadline := time.Now().Add(overall)
	for time.Now().Before(deadline) {
		num, err := ti.GetNumBytesAvailableInBuffer()
		if err == nil && num >= 4 {
			hdr, err := ti.PeekReadableBytes(context.Background(), 4)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			size := uint32(uint16(hdr[2])<<8 | uint16(hdr[3]))
			if num < size {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			data, err := ti.Read(context.Background(), size)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			return data, true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, false
}

// TestTransportInstance_DelayedReplyAfterIdle reproduces the field scenario:
// a request/reply works, the connection then sits idle well past the per-poll
// read deadline, and a *delayed* reply arrives. It must still reach the codec.
func TestTransportInstance_DelayedReplyAfterIdle(t *testing.T) {
	for _, idle := range []time.Duration{0, 500 * time.Millisecond, 2 * time.Second} {
		t.Run(idle.String(), func(t *testing.T) {
			peer, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
			require.NoError(t, err)
			t.Cleanup(func() { _ = peer.Close() })

			ti := NewTransportInstance(nil, peer.LocalAddr().(*net.UDPAddr), false, nil)
			require.NoError(t, ti.Connect(t.Context()))
			t.Cleanup(func() { _ = ti.Close() })

			// First exchange: client writes so the peer learns the client addr.
			require.NoError(t, ti.Write(context.Background(), ackBytes))
			buf := make([]byte, 4096)
			require.NoError(t, peer.SetReadDeadline(time.Now().Add(2*time.Second)))
			_, clientAddr, err := peer.ReadFromUDP(buf)
			require.NoError(t, err)

			// Immediate reply must arrive.
			_, err = peer.WriteToUDP(ackBytes, clientAddr)
			require.NoError(t, err)
			_, ok := pollOneLikeCodec(t, ti, 2*time.Second)
			require.True(t, ok, "immediate reply lost")

			// Idle gap: the receive loop keeps polling (arming 10ms deadlines that
			// keep expiring) with no traffic.
			if idle > 0 {
				stop := time.Now().Add(idle)
				for time.Now().Before(stop) {
					_, _ = ti.GetNumBytesAvailableInBuffer()
					time.Sleep(10 * time.Millisecond)
				}
			}

			// Delayed reply after the idle gap.
			_, err = peer.WriteToUDP(ackBytes, clientAddr)
			require.NoError(t, err)
			_, ok = pollOneLikeCodec(t, ti, 2*time.Second)
			require.True(t, ok, "delayed reply after %s idle gap was lost", idle)
		})
	}
}
