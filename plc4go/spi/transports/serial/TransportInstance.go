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

package serial

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance

	SerialPortName string
	BaudRate       uint
	ConnectTimeout uint32

	connected        atomic.Bool
	stateChangeMutex sync.Mutex

	transport  *Transport
	serialPort serialport.Port
	reader     *bufio.Reader

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTransportInstance(serialPortName string, baudRate uint, connectTimeout uint32, transport *Transport, _options ...options.WithOption) *TransportInstance {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	transportInstance := &TransportInstance{
		SerialPortName: serialPortName,
		BaudRate:       baudRate,
		ConnectTimeout: connectTimeout,
		transport:      transport,

		log: customLogger,
	}
	transportInstance.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(transportInstance, _options...)
	return transportInstance
}

func (m *TransportInstance) Connect(ctx context.Context) error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if m.connected.Load() {
		return errors.New("Already connected")
	}

	serialPort, err := serialport.Open(m.SerialPortName, serialport.Config{BaudRate: m.BaudRate})
	if err != nil {
		return errors.Wrap(err, "error connecting to serial port")
	}
	m.serialPort = serialPort
	m.reader = bufio.NewReader(m.serialPort)
	m.connected.Store(true)

	return nil
}

func (m *TransportInstance) Reset() {
	// No-Op
}

func (m *TransportInstance) Close() error {
	defer utils.StopWarn(m.log)()
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()

	if m.serialPort == nil {
		return nil
	}
	err := m.serialPort.Close()
	if err != nil {
		return errors.Wrap(err, "error closing serial port")
	}
	m.serialPort = nil

	m.connected.Store(false)
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.serialPort != nil
}

func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
	if !m.connected.Load() {
		return errors.New("error writing to transport. Not connected")
	}
	if m.serialPort == nil {
		return errors.New("error writing to transport. No writer available")
	}
	if deadline, ok := ctx.Deadline(); ok {
		if err := m.serialPort.SetWriteDeadline(deadline); err != nil {
			return errors.Wrap(err, "error setting write deadline")
		}
	}
	num, err := m.serialPort.Write(data)
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
	serialPort := m.serialPort
	if serialPort == nil {
		return errors.New("error setting read deadline. No serial port available")
	}
	return serialPort.SetReadDeadline(deadline)
}

func (m *TransportInstance) String() string {
	return fmt.Sprintf("serial:%s:%d", m.SerialPortName, m.BaudRate)
}

func (m *TransportInstance) ClassifyError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transports.IsTransientSyscallError(err) {
		return transports.TransportErrorTransient
	}
	if transports.ErrorIs(err, io.EOF) {
		return transports.TransportErrorFatal
	}
	lower := strings.ToLower(err.Error())
	switch {
	case strings.Contains(lower, "timeout"):
		return transports.TransportErrorRetryable
	case strings.Contains(lower, "temporarily unavailable"):
		return transports.TransportErrorTransient
	case strings.Contains(lower, "closed"):
		return transports.TransportErrorFatal
	}
	return transports.TransportErrorFatal
}
