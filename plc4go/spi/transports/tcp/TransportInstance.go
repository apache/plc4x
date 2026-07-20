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
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// defaultReadPollTimeout bounds reads that carry no explicit deadline (the
// codec's receive worker polls with a deadline-less long-lived context): a
// silent connection must yield a retryable timeout instead of blocking a
// worker goroutine forever. Mirrors the serial transport's readTimeout
// default.
const defaultReadPollTimeout = time.Second

type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance

	RemoteAddress  *net.TCPAddr
	LocalAddress   *net.TCPAddr
	ConnectTimeout uint32

	transport *Transport

	tcpConn net.Conn
	armer   *deadlineReader
	reader  *bufio.Reader

	connected        atomic.Bool
	stateChangeMutex sync.RWMutex

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTcpTransportInstance(remoteAddress *net.TCPAddr, connectTimeout uint32, transport *Transport, _options ...options.WithOption) *TransportInstance {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	transportInstance := &TransportInstance{
		RemoteAddress:  remoteAddress,
		ConnectTimeout: connectTimeout,
		transport:      transport,

		log: customLogger,
	}
	transportInstance.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(transportInstance, _options...)
	return transportInstance
}

func (m *TransportInstance) Connect(ctx context.Context) error {
	if m.connected.Load() {
		return errors.New("already connected")
	}
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if m.RemoteAddress == nil {
		return errors.New("Required remote address missing")
	}
	var err error
	var d net.Dialer
	m.tcpConn, err = d.DialContext(ctx, "tcp", m.RemoteAddress.String())
	if err != nil {
		return errors.Wrap(err, "error connecting to remote address")
	}

	m.LocalAddress = m.tcpConn.LocalAddr().(*net.TCPAddr)

	m.armer = newDeadlineReader(m.tcpConn, defaultReadPollTimeout)
	m.reader = bufio.NewReaderSize(m.armer, 100000)

	m.connected.Store(true)
	return nil
}

// Reset is deliberately a no-op for TCP, matching the UDP and serial
// transports. It used to poke the read deadline, drain one socket read into a
// discard buffer, and swap m.reader — but it is called from OUTSIDE the
// receive worker (connection-cache lease grants, DefaultCodec retryable-error
// handling) while the worker concurrently drives
// GetNumBytesAvailableInBuffer/Peek/Read on the same conn and reader. Every
// one of those mutations races the worker, and on a STREAM they are doubly
// destructive: the drain can consume a fresh reply, and dropping bytes the
// reader had already buffered desyncs the protocol framing for everything
// that follows (the swap also silently shrank the reader from the 100k
// buffer Connect creates to bufio's 4k default). Stale complete frames from
// an abandoned exchange are already handled by the codec's response
// matching; reconnection semantics live in Close/Connect.
// See TestTransportInstance_ResetKeepsBufferedBytes /
// TestTransportInstance_ResetDoesNotRaceReceiveWorker.
func (m *TransportInstance) Reset() {
	m.log.Trace().Msg("Reset is a no-op for TCP (see doc comment)")
}

func (m *TransportInstance) Close() error {
	defer utils.StopWarn(m.log)()
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if !m.connected.Load() {
		return nil
	}
	if err := m.tcpConn.Close(); err != nil {
		return errors.Wrap(err, "error closing connection")
	}
	m.connected.Store(false)
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
	if !m.connected.Load() {
		return errors.New("error writing to transport. Not connected")
	}
	// Bound the write itself by the caller's context deadline. This used to arm
	// a READ deadline (copy-paste from the read path): the write stayed
	// unbounded, and the sticky read deadline broke reads of replies arriving
	// after it had expired.
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.tcpConn.SetWriteDeadline(deadline); err != nil {
			return errors.Wrap(err, "error setting write deadline")
		}
	} else if err := m.tcpConn.SetWriteDeadline(time.Time{}); err != nil {
		// Clear whatever deadline a previous bounded Write left armed - an
		// expired one would fail this (unbounded) write instantly.
		return errors.Wrap(err, "error clearing write deadline")
	}
	num, err := m.tcpConn.Write(data)
	if err != nil {
		return errors.Wrap(err, "error writing")
	}
	if num != len(data) {
		return errors.New("error writing: not all bytes written")
	}
	return nil
}

func (m *TransportInstance) GetReader() transports.ExtendedReader {
	return m.reader
}

func (m *TransportInstance) SetReadDeadline(deadline time.Time) error {
	// Routed through the deadlineReader (NOT tcpConn.SetDeadline, which arms
	// both directions and made read deadlines fail writes): the explicit
	// deadline wins over the poll fallback for the next read(s) and, once
	// observed expired, auto-clears instead of sticking to the connection.
	armer := m.armer
	if armer == nil {
		return errors.New("error setting read deadline. Not connected")
	}
	armer.setExplicitDeadline(deadline)
	return nil
}

func (m *TransportInstance) String() string {
	localAddress := ""
	if m.LocalAddress != nil {
		localAddress = m.LocalAddress.String() + "->"
	}
	return fmt.Sprintf("tcp:%s%s", localAddress, m.RemoteAddress)
}

// ClassifyError attempts to map common network errors to transport severity categories.
func (m *TransportInstance) ClassifyError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transports.ErrorIs(err, io.EOF) || transports.ErrorIs(err, net.ErrClosed) {
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
		var syscallErr syscall.Errno
		if transports.ErrorAs(opErr.Err, &syscallErr) {
			switch syscallErr {
			case syscall.ECONNREFUSED, syscall.ECONNRESET, syscall.EPIPE, syscall.ENETDOWN, syscall.ENETUNREACH:
				return transports.TransportErrorFatal
			case syscall.ETIMEDOUT:
				return transports.TransportErrorRetryable
			}
		}
		return transports.TransportErrorFatal
	}
	return transports.TransportErrorFatal
}
