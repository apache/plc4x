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
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TransportInstance struct {
	LocalAddress  *net.UDPAddr
	RemoteAddress *net.UDPAddr
	// SoReUse historically toggled SO_REUSEPORT via github.com/libp2p/go-reuseport.
	// In practice that does not deliver broadcast traffic to multiple host-local
	// BACnet stacks (the kernel hashes each broadcast to one socket on Linux,
	// and behaviour is platform-dependent elsewhere), so this field is kept for
	// API compatibility but no longer changes how the socket is bound.
	//
	// Deprecated: no-op as of plc4go 1.0; will be removed once the API contract
	// can absorb the signature change.
	SoReUse bool

	transport *Transport
	udpConn   *net.UDPConn
	reader    *bufio.Reader

	connected        atomic.Bool
	stateChangeMutex sync.RWMutex

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTransportInstance(localAddress *net.UDPAddr, remoteAddress *net.UDPAddr, soReUse bool, transport *Transport, _options ...options.WithOption) *TransportInstance {
	logger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &TransportInstance{
		LocalAddress:  localAddress,
		RemoteAddress: remoteAddress,
		SoReUse:       soReUse,
		transport:     transport,

		log: logger,
	}
}

func (m *TransportInstance) Connect(ctx context.Context) error {
	if m.connected.Load() {
		return errors.New("already connected")
	}
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	// If we haven't provided a local address, have the system figure it out by dialing
	// the remote address and then using that connections local address as local address.
	if m.LocalAddress == nil && m.RemoteAddress != nil {
		var d net.Dialer
		udpTest, err := d.DialContext(ctx, "udp", m.RemoteAddress.String())
		if err != nil {
			return errors.Wrapf(err, "error connecting to remote address '%s'", m.RemoteAddress)
		}
		m.LocalAddress = udpTest.LocalAddr().(*net.UDPAddr)
		err = udpTest.Close()
		if err != nil {
			return errors.Wrap(err, "error closing test-connection")
		}
	}

	// "connect" to the remote
	var err error
	if m.RemoteAddress != nil {
		if m.udpConn, err = net.DialUDP("udp", m.LocalAddress, m.RemoteAddress); err != nil {
			return errors.Wrapf(err, "error connecting to remote address '%s'", m.RemoteAddress)
		}
	} else {
		// Listen-only mode (no RemoteAddress). The SoReUse field used to switch
		// to a SO_REUSEPORT-enabled bind via libp2p/go-reuseport; that did not
		// actually solve the "multiple BACnet stacks on one host" problem (see
		// the SoReUse doc comment), so both branches now use stdlib net.ListenUDP.
		if m.udpConn, err = net.ListenUDP("udp", m.LocalAddress); err != nil {
			return errors.Wrapf(err, "error connecting to local address '%s'", m.LocalAddress)
		}
	}

	// Passive bufio.Reader over the UDP socket — same pattern the TCP
	// transport uses. The codec's Receive worker drives reads through
	// PeekReadableBytes/Read/FillBuffer with a deadline set from the request
	// context, so we don't need a separate pump goroutine.
	m.reader = bufio.NewReader(m.udpConn)

	m.connected.Store(true)

	return nil
}

func (m *TransportInstance) Reset() {
	if m.udpConn == nil {
		m.log.Trace().Msg("No connection to reset")
		return
	}
	_ = m.udpConn.SetReadDeadline(time.Now().Add(1))
	_, _, _ = m.udpConn.ReadFromUDP(make([]byte, 4096))
	m.reader = bufio.NewReader(m.udpConn)
	m.log.Trace().Msg("Connection reset")
}

func (m *TransportInstance) Close() error {
	defer utils.StopWarn(m.log)()
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if !m.connected.Load() {
		return nil
	}
	err := m.udpConn.Close()
	if err != nil {
		return errors.Wrap(err, "error closing connection")
	}
	m.connected.Store(false)
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if !m.IsConnected() {
		return 0, errors.New("working on a unconnected connection")
	}
	if m.reader == nil {
		return 0, nil
	}
	// Use a fresh, short read deadline for this poll. Read/PeekReadableBytes set
	// a sticky SetReadDeadline from the request context; once that deadline has
	// passed, a deadline-less Peek here would keep failing with i/o timeout and
	// the codec would never observe further inbound datagrams.
	if m.udpConn != nil {
		_ = m.udpConn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	}
	_, _ = m.reader.Peek(1)
	return uint32(m.reader.Buffered()), nil
}

func (m *TransportInstance) FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool)) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	nBytes := uint32(1)
	for ctx.Err() == nil {
		_bytes, err := m.PeekReadableBytes(ctx, nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		if keepGoing := until(uint(nBytes-1), _bytes[len(_bytes)-1], m.reader); !keepGoing {
			return nil
		}
		nBytes++
	}
	return errors.Wrap(ctx.Err(), "Timeout while filling buffer")
}

func (m *TransportInstance) PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.udpConn.SetReadDeadline(deadline); err != nil {
			return nil, errors.Wrap(err, "error setting read deadline")
		}
	}
	return m.reader.Peek(int(numBytes))
}

func (m *TransportInstance) Read(ctx context.Context, numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	data := make([]byte, numBytes)
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.udpConn.SetReadDeadline(deadline); err != nil {
			return nil, errors.Wrap(err, "error setting read deadline")
		}
	}
	for i := range numBytes {
		val, err := m.reader.ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}

func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.udpConn.SetWriteDeadline(deadline); err != nil {
			return errors.Wrap(err, "error setting read deadline")
		}
	}
	var num int
	var err error
	// A connected UDP socket (obtained via net.DialUDP) rejects WriteToUDP with
	// "use of WriteTo with pre-connected connection" — we have to use the plain
	// Write() path instead. udpConn.RemoteAddr() is nil for ListenUDP sockets and
	// the connected remote for DialUDP sockets, so that's the right discriminator.
	if m.udpConn.RemoteAddr() != nil {
		num, err = m.udpConn.Write(data)
	} else if m.RemoteAddress != nil {
		num, err = m.udpConn.WriteToUDP(data, m.RemoteAddress)
	} else {
		num, err = m.udpConn.Write(data)
	}
	if err != nil {
		return errors.Wrapf(err, "error writing (remote address: %s)", m.RemoteAddress)
	}
	if num != len(data) {
		return errors.Errorf("error writing: not all bytes written (Expected %d, Actual %d)", len(data), num)
	}
	return nil
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("udp:%s->%s", m.LocalAddress, m.RemoteAddress)
}

func (m *TransportInstance) ClassifyError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transports.ErrorIs(err, io.EOF) || transports.ErrorIs(err, net.ErrClosed) || transports.ErrorIs(err, syscall.EPIPE) {
		return transports.TransportErrorFatal
	}
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return transports.TransportErrorRetryable
		}
	}
	if transports.IsTransientSyscallError(err) {
		return transports.TransportErrorTransient
	}
	var opErr *net.OpError
	if transports.ErrorAs(err, &opErr) && opErr != nil {
		if opErr.Timeout() {
			return transports.TransportErrorRetryable
		}
		if transports.IsTransientSyscallError(opErr.Err) {
			return transports.TransportErrorTransient
		}
		if syscallErr, ok := opErr.Err.(syscall.Errno); ok {
			switch syscallErr {
			case syscall.ECONNRESET, syscall.ECONNREFUSED, syscall.ENETDOWN, syscall.ENETUNREACH:
				return transports.TransportErrorFatal
			case syscall.ETIMEDOUT:
				return transports.TransportErrorRetryable
			}
		}
	}
	return transports.TransportErrorFatal
}
