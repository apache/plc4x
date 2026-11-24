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

package virtualcan

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/can"
)

// maxBufferedFrames bounds what an instance holds for a reader which doesn't keep up. Past it the
// oldest frames are dropped, the way a CAN controller overruns, so a stalled reader can't make the
// whole process grow.
const maxBufferedFrames = 4096

// errNotConnected is what every operation on a closed instance returns.
var errNotConnected = errors.New("working on a unconnected connection")

// TransportInstance is one node on a virtual bus.
type TransportInstance struct {
	BusName string

	mutex     sync.Mutex
	bus       *bus
	buffer    []byte
	notify    chan struct{}
	dropped   uint64
	connected bool

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTransportInstance(busName string, _options ...options.WithOption) *TransportInstance {
	return &TransportInstance{
		BusName: busName,
		notify:  make(chan struct{}, 1),
		log:     options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

func (m *TransportInstance) Connect(_ context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.connected {
		return errors.New("already connected")
	}
	m.buffer = nil
	m.connected = true
	m.bus = join(m.BusName, m)
	return nil
}

func (m *TransportInstance) Close() error {
	m.mutex.Lock()
	if !m.connected {
		m.mutex.Unlock()
		return nil
	}
	m.connected = false
	m.bus = nil
	m.buffer = nil
	m.mutex.Unlock()
	leave(m.BusName, m)
	m.wake()
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.connected
}

// Reset does nothing; see the SocketCAN transport for why a frame stream has nothing to reset.
func (m *TransportInstance) Reset() {}

// receive appends records written by another member.
func (m *TransportInstance) receive(records []byte) {
	m.mutex.Lock()
	if !m.connected {
		m.mutex.Unlock()
		return
	}
	m.buffer = append(m.buffer, records...)
	if overflow := len(m.buffer) - maxBufferedFrames*can.FrameSize; overflow > 0 {
		// Drop whole frames from the front, never part of one.
		overflow = (overflow + can.FrameSize - 1) / can.FrameSize * can.FrameSize
		m.buffer = append([]byte(nil), m.buffer[overflow:]...)
		m.dropped += uint64(overflow / can.FrameSize)
	}
	m.mutex.Unlock()
	m.wake()
}

func (m *TransportInstance) wake() {
	select {
	case m.notify <- struct{}{}:
	default:
	}
}

// waitFor blocks until numBytes are buffered, the context ends or the instance is closed.
func (m *TransportInstance) waitFor(ctx context.Context, numBytes uint32) error {
	for {
		m.mutex.Lock()
		connected, available := m.connected, uint32(len(m.buffer))
		m.mutex.Unlock()
		if !connected {
			return errNotConnected
		}
		if available >= numBytes {
			return nil
		}
		select {
		case <-m.notify:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if !m.connected {
		return 0, errNotConnected
	}
	return uint32(len(m.buffer)), nil
}

func (m *TransportInstance) FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool)) error {
	for position := uint32(0); ; position++ {
		peeked, err := m.PeekReadableBytes(ctx, position+1)
		if err != nil {
			return err
		}
		if !until(uint(position), peeked[position], nil) {
			return nil
		}
	}
}

func (m *TransportInstance) PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error) {
	if err := m.waitFor(ctx, numBytes); err != nil {
		return nil, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if uint32(len(m.buffer)) < numBytes {
		return nil, errNotConnected
	}
	return append([]byte(nil), m.buffer[:numBytes]...), nil
}

func (m *TransportInstance) Read(ctx context.Context, numBytes uint32) ([]byte, error) {
	if err := m.waitFor(ctx, numBytes); err != nil {
		return nil, err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if uint32(len(m.buffer)) < numBytes {
		// The front frames were dropped by an overrun, or the instance closed, between the wait and
		// the lock.
		return nil, errors.New("buffered frames were dropped while reading")
	}
	data := append([]byte(nil), m.buffer[:numBytes]...)
	m.buffer = m.buffer[numBytes:]
	return data, nil
}

// Write puts whole frame records on the bus.
func (m *TransportInstance) Write(_ context.Context, data []byte) error {
	if len(data)%can.FrameSize != 0 {
		return errors.Errorf("virtualcan writes whole %d byte frame records, got %d bytes", can.FrameSize, len(data))
	}
	for offset := 0; offset < len(data); offset += can.FrameSize {
		if _, err := can.Decode(data[offset : offset+can.FrameSize]); err != nil {
			return errors.Wrap(err, "refusing to put a malformed frame on the bus")
		}
	}
	m.mutex.Lock()
	b := m.bus
	m.mutex.Unlock()
	if b == nil {
		return errNotConnected
	}
	b.deliver(m, append([]byte(nil), data...))
	return nil
}

// DroppedFrames is how many frames were discarded because the reader fell behind.
func (m *TransportInstance) DroppedFrames() uint64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.dropped
}

func (m *TransportInstance) ClassifyError(err error) transports.TransportErrorKind {
	switch {
	case err == nil:
		return transports.TransportErrorUnknown
	case errors.Is(err, context.DeadlineExceeded):
		return transports.TransportErrorRetryable
	}
	return transports.TransportErrorFatal
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("virtualcan:%s", m.BusName)
}
