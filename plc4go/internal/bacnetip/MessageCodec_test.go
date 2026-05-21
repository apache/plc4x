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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// fakeTransportInstance is an in-memory transport tailored for codec-level
// tests: bytes pushed via PushReceived land directly in the read buffer (no
// channel pump), and bytes written by the codec land in writeBuffer. Lets us
// exercise Send/Receive without standing up a real socket.
type fakeTransportInstance struct {
	mu          sync.Mutex
	readBuffer  bytes.Buffer
	writeBuffer bytes.Buffer
	connected   bool
}

func newFakeTransport() *fakeTransportInstance {
	return &fakeTransportInstance{connected: true}
}

func (f *fakeTransportInstance) PushReceived(data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.readBuffer.Write(data)
}

func (f *fakeTransportInstance) WrittenBytes() []byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]byte(nil), f.writeBuffer.Bytes()...)
}

func (f *fakeTransportInstance) Connect(_ context.Context) error           { f.connected = true; return nil }
func (f *fakeTransportInstance) ConnectWithContext(_ context.Context) error { return f.Connect(nil) }
func (f *fakeTransportInstance) Close() error                              { f.connected = false; return nil }
func (f *fakeTransportInstance) IsConnected() bool                         { return f.connected }
func (f *fakeTransportInstance) String() string                            { return "fake" }
func (f *fakeTransportInstance) Reset()                                    {}

func (f *fakeTransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return uint32(f.readBuffer.Len()), nil
}

func (f *fakeTransportInstance) PeekReadableBytes(_ context.Context, n uint32) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if uint32(f.readBuffer.Len()) < n {
		return nil, fmt.Errorf("not enough bytes: want %d, have %d", n, f.readBuffer.Len())
	}
	buf := f.readBuffer.Bytes()
	return append([]byte(nil), buf[:n]...), nil
}

func (f *fakeTransportInstance) Read(_ context.Context, n uint32) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if uint32(f.readBuffer.Len()) < n {
		return nil, fmt.Errorf("not enough bytes")
	}
	out := make([]byte, n)
	_, _ = f.readBuffer.Read(out)
	return out, nil
}

func (f *fakeTransportInstance) Write(_ context.Context, data []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.writeBuffer.Write(data)
	return nil
}

func (f *fakeTransportInstance) FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool)) error {
	// Implemented for completeness; the codec doesn't drive this path in
	// these tests. Mirrors the contract of DefaultBufferedTransportInstance.
	nBytes := uint32(1)
	for ctx.Err() == nil {
		b, err := f.PeekReadableBytes(ctx, nBytes)
		if err != nil {
			return err
		}
		if !until(uint(nBytes-1), b[len(b)-1], bufio.NewReader(bytes.NewReader(b))) {
			return nil
		}
		nBytes++
	}
	return ctx.Err()
}

var _ transports.TransportInstance = (*fakeTransportInstance)(nil)

func newTestCodec(t *testing.T) (*MessageCodec, *fakeTransportInstance) {
	t.Helper()
	ti := newFakeTransport()
	codec := NewMessageCodec(ti)
	t.Cleanup(func() { _ = codec.Disconnect() })
	return codec, ti
}

// makeWhoIsBVLC builds a minimal valid BVLC frame (BVLCOriginalUnicastNPDU
// wrapping an unconfirmed WhoIs APDU). Used as a payload for both Send and
// Receive round-trip tests because it serializes to a small, deterministic
// byte sequence.
func makeWhoIsBVLC(t *testing.T) readWriteModel.BVLC {
	t.Helper()
	whoIs := readWriteModel.NewBACnetUnconfirmedServiceRequestWhoIs(nil, nil)
	apdu := readWriteModel.NewAPDUUnconfirmedRequest(whoIs)
	return wrapAPDU(apdu, false)
}

func TestMessageCodec_Send_SerializesBVLCToTransport(t *testing.T) {
	codec, ti := newTestCodec(t)
	bvlc := makeWhoIsBVLC(t)
	expectedBytes, err := bvlc.Serialize()
	require.NoError(t, err)

	require.NoError(t, codec.Send(context.Background(), "test", bvlc))

	assert.Equal(t, expectedBytes, ti.WrittenBytes())
}

func TestMessageCodec_Send_RejectsNonBVLC(t *testing.T) {
	codec, _ := newTestCodec(t)
	apdu := readWriteModel.NewAPDUUnconfirmedRequest(
		readWriteModel.NewBACnetUnconfirmedServiceRequestWhoIs(nil, nil),
	)
	assert.Panics(t, func() {
		_ = codec.Send(context.Background(), "test", apdu)
	}, "Send must panic on non-BVLC payload — callers must wrap with wrapAPDU first")
}

func TestMessageCodec_Receive_ParsesBVLCFromTransport(t *testing.T) {
	codec, ti := newTestCodec(t)
	bvlc := makeWhoIsBVLC(t)
	raw, err := bvlc.Serialize()
	require.NoError(t, err)
	ti.PushReceived(raw)

	msg, err := codec.Receive(context.Background())
	require.NoError(t, err)
	require.NotNil(t, msg, "Receive should return the parsed BVLC")

	parsed, ok := msg.(readWriteModel.BVLC)
	require.True(t, ok, "Receive should return a BVLC, got %T", msg)
	assert.Equal(t, bvlc.GetBvlcFunction(), parsed.GetBvlcFunction())
}

func TestMessageCodec_Receive_NotEnoughBytesReturnsNil(t *testing.T) {
	codec, ti := newTestCodec(t)
	// 2 bytes is below the 4-byte minimum BVLC header — Receive must
	// short-circuit with (nil, nil), not block or err.
	ti.PushReceived([]byte{0x81, 0x0a})
	msg, err := codec.Receive(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, msg, "<4 buffered bytes should not yield a message")
}

func TestMessageCodec_Receive_PartialPacketReturnsNil(t *testing.T) {
	codec, ti := newTestCodec(t)
	// BVLC header claims packet-size 100 (0x0064), but we only feed 4 bytes.
	ti.PushReceived([]byte{0x81, 0x0a, 0x00, 0x64})
	msg, err := codec.Receive(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, msg, "buffer < declared packet-size should not yield a message")
}

func TestMessageCodec_ExpectationsMatchUnsolicitedMessages(t *testing.T) {
	// Regression for the always-true handleCustomMessage bug:
	// expectations must still match incoming messages even when a
	// customMessageHandling (keepReceiveLoopActive) is registered, because
	// our handler returns false and lets HandleMessages run.
	codec, ti := newTestCodec(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	got := make(chan spi.Message, 1)
	codec.Expect(ctx, "test-expect",
		func(_ spi.Message) bool { return true },
		func(msg spi.Message) error { got <- msg; return nil },
		func(err error) error { return err },
	)

	require.NoError(t, codec.Connect(ctx))
	raw, err := makeWhoIsBVLC(t).Serialize()
	require.NoError(t, err)
	ti.PushReceived(raw)

	select {
	case msg := <-got:
		_, ok := msg.(readWriteModel.BVLC)
		assert.True(t, ok, "expectation should receive the parsed BVLC, got %T", msg)
	case <-time.After(3 * time.Second):
		t.Fatal("expectation handler never fired — keepReceiveLoopActive may be swallowing messages")
	}
}

func TestKeepReceiveLoopActive_AlwaysReturnsFalse(t *testing.T) {
	// keepReceiveLoopActive's only job is to keep the receive worker awake
	// when there are zero expectations. Its return value of false is what
	// lets HandleMessages → defaultIncomingMessageChannel still run. If
	// somebody flips this to true, every Read/Write/Subscribe times out
	// (the original bug we hit during integration testing).
	bvlc := makeWhoIsBVLC(t)
	assert.False(t,
		keepReceiveLoopActive(context.Background(), (_default.DefaultCodecRequirements)(nil), bvlc),
		"keepReceiveLoopActive must return false or expectation matching breaks")
}
