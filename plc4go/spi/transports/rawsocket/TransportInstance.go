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

package rawsocket

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
)

// TransportInstance is one AF_PACKET socket bound to one interface.
type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance

	InterfaceName string
	Filter        Filter

	transport *Transport

	stateChangeMutex sync.Mutex
	socket           packetSocket
	reader           *bufio.Reader
	connected        atomic.Bool

	writeMutex sync.Mutex

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

// packetSocket is the platform part: one frame per read and per write.
type packetSocket interface {
	readFrame(buffer []byte) (int, error)
	writeFrame(frame []byte) error
	setReadDeadline(deadline time.Time) error
	setWriteDeadline(deadline time.Time) error
	hardwareAddress() net.HardwareAddr
	close() error
}

func NewTransportInstance(interfaceName string, filter Filter, transport *Transport, _options ...options.WithOption) *TransportInstance {
	instance := &TransportInstance{
		InterfaceName: interfaceName,
		Filter:        filter,
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
	program, err := m.Filter.Program()
	if err != nil {
		return errors.Wrap(err, "error compiling the packet filter")
	}
	socket, err := openSocket(m.InterfaceName, program)
	if err != nil {
		return errors.Wrapf(err, "error opening a raw socket on %s", m.InterfaceName)
	}
	m.socket = socket
	m.reader = bufio.NewReaderSize(&recordReader{socket: socket}, 8*(MaxFrameSize+RecordHeaderSize))
	m.connected.Store(true)
	m.log.Debug().Str("interface", m.InterfaceName).Msg("raw socket open")
	return nil
}

func (m *TransportInstance) Close() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if !m.connected.Swap(false) {
		return nil
	}
	return m.socket.close()
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

// Reset does nothing: the stream is whole records, there is no partial state to clear.
func (m *TransportInstance) Reset() {}

func (m *TransportInstance) GetReader() transports.ExtendedReader {
	return m.reader
}

func (m *TransportInstance) SetReadDeadline(deadline time.Time) error {
	return m.socket.setReadDeadline(deadline)
}

// HardwareAddress is the interface's MAC address, the source address of what is sent.
func (m *TransportInstance) HardwareAddress() net.HardwareAddr {
	if !m.IsConnected() {
		return nil
	}
	return m.socket.hardwareAddress()
}

// LocalAddress is the interface's first IPv4 address, if it has one.
func (m *TransportInstance) LocalAddress() net.Addr {
	networkInterface, err := net.InterfaceByName(m.InterfaceName)
	if err != nil {
		return nil
	}
	addresses, err := networkInterface.Addrs()
	if err != nil {
		return nil
	}
	for _, address := range addresses {
		if ipNet, ok := address.(*net.IPNet); ok && ipNet.IP.To4() != nil {
			return ipNet
		}
	}
	return nil
}

// Write sends one or more frame records, each as its own frame.
func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()
	deadline, _ := ctx.Deadline()
	if err := m.socket.setWriteDeadline(deadline); err != nil {
		return errors.Wrap(err, "error setting the write deadline")
	}
	for len(data) > 0 {
		length, err := RecordFrameLength(data)
		if err != nil {
			return err
		}
		if len(data) < RecordHeaderSize+length {
			return errors.Errorf("record announces %d bytes, %d follow", length, len(data)-RecordHeaderSize)
		}
		if err := m.socket.writeFrame(data[RecordHeaderSize : RecordHeaderSize+length]); err != nil {
			return errors.Wrapf(err, "error sending on %s", m.InterfaceName)
		}
		data = data[RecordHeaderSize+length:]
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
		return transports.TransportErrorTransient
	case transports.IsTransientSyscallError(err):
		return transports.TransportErrorTransient
	}
	return transports.TransportErrorFatal
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("raw-socket:%s", m.InterfaceName)
}

// recordReader turns one-frame-per-read into the record stream.
type recordReader struct {
	socket  packetSocket
	pending []byte
	buffer  [RecordHeaderSize + MaxFrameSize + 64]byte
}

func (r *recordReader) Read(p []byte) (int, error) {
	for len(r.pending) == 0 {
		n, err := r.socket.readFrame(r.buffer[RecordHeaderSize:])
		if err != nil {
			return 0, err
		}
		if n < MinFrameSize || n > MaxFrameSize {
			// Runts and oversized frames can't be a record; skip them rather than break the stream.
			continue
		}
		binary.BigEndian.PutUint16(r.buffer[:RecordHeaderSize], uint16(n))
		r.pending = r.buffer[:RecordHeaderSize+n]
	}
	n := copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}
