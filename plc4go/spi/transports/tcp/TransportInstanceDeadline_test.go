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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Write used to arm a READ deadline from the request context (copy-paste from
// the read path). That deadline was sticky: once it expired, every socket read
// failed instantly with i/o timeout, so a reply arriving after the request
// deadline sat unread in the kernel buffer until the NEXT Write re-armed a
// fresh deadline. A late reply must instead reach the reader - the codec's
// response matching decides what to do with it, not the transport.
func TestTransportInstance_WriteMustNotArmReadDeadline(t *testing.T) {
	ti, serverConn := connectedTestInstance(t)

	writeCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	require.NoError(t, ti.Write(writeCtx, []byte{0x00, 0x01}))

	// Let the request deadline expire, then deliver the (late) reply.
	time.Sleep(150 * time.Millisecond)
	frame := []byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x04, 0x01, 0x03, 0x02, 0x2a}
	_, err := serverConn.Write(frame)
	require.NoError(t, err)

	// The reply must still be picked up by the buffered read path.
	waitBuffered(t, ti, uint32(len(frame)))

	readCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	got, err := ti.Read(readCtx, uint32(len(frame)))
	require.NoError(t, err, "a reply arriving after the write context deadline must still be readable")
	assert.Equal(t, frame, got)
}

// Write never armed a WRITE deadline at all, so a wedged peer (nothing
// draining the connection, kernel buffers full) blocked Write() indefinitely
// - past any caller context deadline. This is the unbounded half of the
// copy-paste bug.
func TestTransportInstance_WriteHonorsContextDeadline(t *testing.T) {
	ti, serverConn := connectedTestInstance(t)

	// Shrink both kernel buffers as far as the OS allows and never read on the
	// server side, so a large enough write must block in the kernel.
	require.NoError(t, ti.tcpConn.(*net.TCPConn).SetWriteBuffer(1))
	require.NoError(t, serverConn.(*net.TCPConn).SetReadBuffer(1))

	payload := make([]byte, 64*1024*1024)
	writeCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	start := time.Now()
	go func() {
		errCh <- ti.Write(writeCtx, payload)
	}()

	select {
	case err := <-errCh:
		require.Error(t, err, "a write the peer never drains must fail once the context deadline passes")
		var netErr net.Error
		require.ErrorAs(t, err, &netErr)
		assert.True(t, netErr.Timeout(), "expected a timeout error, got: %v", err)
		assert.Less(t, time.Since(start), 3*time.Second, "write must return promptly after the deadline")
	case <-time.After(5 * time.Second):
		t.Fatal("Write blocked past the context deadline (no write deadline armed)")
	}
}

// The SetReadDeadline METHOD used to call tcpConn.SetDeadline, arming BOTH
// directions: a read deadline set by the buffered read path (or a driver)
// silently bounded - or, once expired, instantly failed - subsequent writes.
func TestTransportInstance_SetReadDeadlineMustNotBoundWrites(t *testing.T) {
	ti, serverConn := connectedTestInstance(t)

	// Keep the connection drained so the write below cannot block.
	go func() {
		buf := make([]byte, 1024)
		for {
			if _, err := serverConn.Read(buf); err != nil {
				return
			}
		}
	}()

	require.NoError(t, ti.SetReadDeadline(time.Now().Add(-time.Second)))

	err := ti.Write(context.Background(), []byte{0x00, 0x01})
	require.NoError(t, err, "an (expired) READ deadline must not fail writes")
}

// A write deadline armed by one bounded Write must not stick to the
// connection: a later Write whose context carries NO deadline would otherwise
// fail instantly once the old deadline expired.
func TestTransportInstance_ExpiredWriteDeadlineDoesNotStick(t *testing.T) {
	ti, serverConn := connectedTestInstance(t)

	// Keep the connection drained so writes cannot block.
	go func() {
		buf := make([]byte, 1024)
		for {
			if _, err := serverConn.Read(buf); err != nil {
				return
			}
		}
	}()

	// A Write whose context deadline already passed fails - that is the
	// caller's bounded contract - and leaves an expired deadline armed.
	expiredCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()
	require.Error(t, ti.Write(expiredCtx, []byte{0x00, 0x01}))

	// A deadline-less Write afterwards must succeed.
	err := ti.Write(context.Background(), []byte{0x00, 0x01})
	require.NoError(t, err, "an expired write deadline from a previous bounded Write must not fail later unbounded writes")
}

// With the sticky read deadline gone, something else has to keep the receive
// worker's bare polls bounded: DefaultCodec drives Receive with a deadline-less
// long-lived context, and bufio's Peek(1) on a silent connection would
// otherwise block forever - piling up one abandoned Receive goroutine per
// receiveTimeout cycle, all racing on the same bufio.Reader once data finally
// arrives. The fallback read deadline (mirroring the serial transport's
// deadlineReader) bounds every deadline-less read.
func TestTransportInstance_BarePollIsBoundedByFallback(t *testing.T) {
	ti, _ := connectedTestInstance(t)

	done := make(chan error, 1)
	go func() {
		_, err := ti.GetReader().Peek(1)
		done <- err
	}()

	select {
	case err := <-done:
		require.Error(t, err, "Peek on a silent connection must fail with a timeout, not return data")
		var netErr net.Error
		require.ErrorAs(t, err, &netErr)
		assert.True(t, netErr.Timeout(), "expected a timeout error, got: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("bare Peek(1) on a silent connection blocked forever (no fallback read deadline)")
	}
}
