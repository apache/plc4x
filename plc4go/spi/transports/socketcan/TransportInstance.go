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

package socketcan

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/can"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
)

// TransportInstance is one raw CAN socket bound to one interface.
type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance

	InterfaceName string
	// ErrorMask is the error classes the kernel delivers as error frames (0: none). See the package
	// comment.
	ErrorMask uint32

	transport *Transport

	stateChangeMutex sync.Mutex
	file             *os.File
	reader           *bufio.Reader
	connected        atomic.Bool

	writeMutex sync.Mutex

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTransportInstance(interfaceName string, transport *Transport, _options ...options.WithOption) *TransportInstance {
	instance := &TransportInstance{
		InterfaceName: interfaceName,
		transport:     transport,
		log:           options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
	instance.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(instance, _options...)
	return instance
}

func (m *TransportInstance) Connect(_ context.Context) error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if m.connected.Load() {
		return errors.New("already connected")
	}
	file, err := openSocket(m.InterfaceName, m.ErrorMask)
	if err != nil {
		return errors.Wrapf(err, "error opening a SocketCAN socket on %s", m.InterfaceName)
	}
	m.file = file
	m.reader = bufio.NewReaderSize(&frameReader{file: file}, 64*can.FrameSize)
	m.connected.Store(true)
	m.log.Debug().Str("interface", m.InterfaceName).Msg("SocketCAN socket open")
	return nil
}

func (m *TransportInstance) Close() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if !m.connected.Swap(false) {
		return nil
	}
	if err := m.file.Close(); err != nil {
		return errors.Wrap(err, "error closing the SocketCAN socket")
	}
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

// Reset does nothing: every record on the stream is a whole frame, so there is no partial state a
// reset could clear without throwing away frames which have already arrived.
func (m *TransportInstance) Reset() {}

// GetReader is for DefaultBufferedTransportInstance.
func (m *TransportInstance) GetReader() transports.ExtendedReader {
	return m.reader
}

// SetReadDeadline is for DefaultBufferedTransportInstance.
func (m *TransportInstance) SetReadDeadline(deadline time.Time) error {
	return m.file.SetReadDeadline(deadline)
}

// Write sends one or more whole frame records. Each goes out as its own write, which is what a raw
// CAN socket requires.
func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	if len(data)%can.FrameSize != 0 {
		return errors.Errorf("SocketCAN writes whole %d byte frame records, got %d bytes", can.FrameSize, len(data))
	}
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()
	deadline, _ := ctx.Deadline()
	if err := m.file.SetWriteDeadline(deadline); err != nil {
		return errors.Wrap(err, "error setting the write deadline")
	}
	record := make([]byte, can.FrameSize)
	for offset := 0; offset < len(data); offset += can.FrameSize {
		copy(record, data[offset:offset+can.FrameSize])
		// The kernel wants the identifier in host byte order.
		binary.NativeEndian.PutUint32(record[0:4], binary.BigEndian.Uint32(record[0:4]))
		n, err := m.file.Write(record)
		if err != nil {
			return errors.Wrapf(err, "error writing to %s", m.InterfaceName)
		}
		if n != can.FrameSize {
			return errors.Errorf("short CAN write on %s: %d of %d bytes", m.InterfaceName, n, can.FrameSize)
		}
	}
	return nil
}

func (m *TransportInstance) ClassifyError(err error) transports.TransportErrorKind {
	switch {
	case err == nil:
		return transports.TransportErrorUnknown
	case transports.ErrorIs(err, os.ErrDeadlineExceeded):
		return transports.TransportErrorRetryable
	case transports.ErrorIs(err, syscall.ENOBUFS):
		// The interface's transmit queue is full: the bus is saturated or nobody acknowledges.
		return transports.TransportErrorTransient
	case transports.IsTransientSyscallError(err):
		return transports.TransportErrorTransient
	}
	return transports.TransportErrorFatal
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("socketcan:%s", m.InterfaceName)
}

// frameReader turns the socket's one-frame-per-read datagrams into the record stream, converting
// the identifier from host to network byte order.
type frameReader struct {
	file    *os.File
	pending []byte
	frame   [can.FrameSize]byte
}

var _ io.Reader = (*frameReader)(nil)

func (r *frameReader) Read(p []byte) (int, error) {
	if len(r.pending) == 0 {
		n, err := r.file.Read(r.frame[:])
		if err != nil {
			return 0, err
		}
		if n != can.FrameSize {
			return 0, errors.Errorf("unexpected %d byte read from a raw CAN socket", n)
		}
		binary.BigEndian.PutUint32(r.frame[0:4], binary.NativeEndian.Uint32(r.frame[0:4]))
		r.pending = r.frame[:]
	}
	n := copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}
