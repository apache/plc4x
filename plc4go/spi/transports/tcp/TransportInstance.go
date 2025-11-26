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

	stdErrors "errors"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance

	RemoteAddress  *net.TCPAddr
	LocalAddress   *net.TCPAddr
	ConnectTimeout uint32

	transport *Transport

	tcpConn net.Conn
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

	m.reader = bufio.NewReaderSize(m.tcpConn, 100000)

	m.connected.Store(true)
	return nil
}

func (m *TransportInstance) Reset() {
	if m.tcpConn == nil {
		m.log.Trace().Msg("No connection to reset")
		return
	}
	_ = m.tcpConn.SetReadDeadline(time.Now().Add(1))
	_, _ = m.tcpConn.Read(make([]byte, 4096))
	m.reader = bufio.NewReader(m.tcpConn)
	m.log.Trace().Msg("Connection reset")
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
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.tcpConn.SetReadDeadline(deadline); err != nil {
			return errors.Wrap(err, "error setting read deadline")
		}
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
	return m.tcpConn.SetDeadline(deadline)
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
	if stdErrors.As(err, &opErr) {
		if opErr.Timeout() {
			return transports.TransportErrorRetryable
		}
		if transports.IsTransientSyscallError(opErr.Err) {
			return transports.TransportErrorTransient
		}
		var syscallErr syscall.Errno
		if stdErrors.As(opErr.Err, &syscallErr) {
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
