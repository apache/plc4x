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

package virtualeth

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

// maxBufferedBytes bounds what an instance holds for a reader that doesn't keep up; past it the
// oldest frames are dropped.
const maxBufferedBytes = 1 << 20

var errNotConnected = errors.New("working on a unconnected connection")

// TransportInstance is one station on a virtual segment.
type TransportInstance struct {
	Segment string
	Filter  rawsocket.Filter

	mac net.HardwareAddr

	mutex     sync.Mutex
	segment   *segment
	buffer    []byte
	notify    chan struct{}
	dropped   uint64
	connected bool

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

// NewTransportInstance makes a station; a nil mac gets a random locally administered one.
func NewTransportInstance(segmentName string, filter rawsocket.Filter, mac net.HardwareAddr, _options ...options.WithOption) *TransportInstance {
	if mac == nil {
		mac = RandomMAC()
	}
	return &TransportInstance{
		Segment: segmentName,
		Filter:  filter,
		mac:     mac,
		notify:  make(chan struct{}, 1),
		log:     options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

// HardwareAddress is the station's MAC address.
func (m *TransportInstance) HardwareAddress() net.HardwareAddr {
	return m.mac
}

func (m *TransportInstance) Connect(_ context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.connected {
		return errors.New("already connected")
	}
	m.buffer = nil
	m.connected = true
	m.segment = join(m.Segment, m)
	return nil
}

func (m *TransportInstance) Close() error {
	m.mutex.Lock()
	if !m.connected {
		m.mutex.Unlock()
		return nil
	}
	m.connected = false
	m.segment = nil
	m.buffer = nil
	m.mutex.Unlock()
	leave(m.Segment, m)
	m.wake()
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.connected
}

func (m *TransportInstance) Reset() {}

func (m *TransportInstance) receive(frame []byte) {
	if !m.Filter.Matches(frame) {
		return
	}
	record, err := rawsocket.EncodeRecord(frame)
	if err != nil {
		return
	}
	m.mutex.Lock()
	if !m.connected {
		m.mutex.Unlock()
		return
	}
	m.buffer = append(m.buffer, record...)
	for len(m.buffer) > maxBufferedBytes {
		length, err := rawsocket.RecordFrameLength(m.buffer)
		if err != nil {
			m.buffer = nil
			break
		}
		m.buffer = m.buffer[rawsocket.RecordHeaderSize+length:]
		m.dropped++
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
		return nil, errors.New("buffered frames were dropped while reading")
	}
	data := append([]byte(nil), m.buffer[:numBytes]...)
	m.buffer = m.buffer[numBytes:]
	return data, nil
}

// Write puts frame records on the segment.
func (m *TransportInstance) Write(_ context.Context, data []byte) error {
	m.mutex.Lock()
	s := m.segment
	m.mutex.Unlock()
	if s == nil {
		return errNotConnected
	}
	for len(data) > 0 {
		length, err := rawsocket.RecordFrameLength(data)
		if err != nil {
			return err
		}
		if len(data) < rawsocket.RecordHeaderSize+length {
			return errors.Errorf("record announces %d bytes, %d follow", length, len(data)-rawsocket.RecordHeaderSize)
		}
		s.deliver(m, append([]byte(nil), data[rawsocket.RecordHeaderSize:rawsocket.RecordHeaderSize+length]...))
		data = data[rawsocket.RecordHeaderSize+length:]
	}
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
	return fmt.Sprintf("virtual-eth:%s(%s)", m.Segment, m.mac)
}
