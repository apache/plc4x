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
	cfg            serialConfig

	connected        atomic.Bool
	stateChangeMutex sync.Mutex

	transport  *Transport
	serialPort serialport.Port
	armer      *deadlineReader
	pacer      *pacer
	reader     *bufio.Reader

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTransportInstance(serialPortName string, baudRate uint, connectTimeout uint32, transport *Transport, _options ...options.WithOption) *TransportInstance {
	cfg := defaultSerialConfig()
	cfg.port.BaudRate = baudRate
	cfg.connectTimeout = connectTimeout
	return NewTransportInstanceWithConfig(serialPortName, cfg, transport, _options...)
}

func NewTransportInstanceWithConfig(serialPortName string, cfg serialConfig, transport *Transport, _options ...options.WithOption) *TransportInstance {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	transportInstance := &TransportInstance{
		SerialPortName: serialPortName,
		BaudRate:       cfg.port.BaudRate,
		ConnectTimeout: cfg.connectTimeout,
		cfg:            cfg,
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

	var serialPort serialport.Port
	var err error
	if m.cfg.reusePort {
		if m.transport == nil {
			return errors.New("reuse-port requires a transport-managed instance")
		}
		serialPort, err = m.transport.registry.acquire(m.SerialPortName, sharedPortConfig{
			port:            m.cfg.port,
			dtr:             m.cfg.dtr,
			rts:             m.cfg.rts,
			interframeDelay: m.cfg.interframeDelay,
		})
		if err != nil {
			return errors.Wrap(err, "error acquiring shared serial port")
		}
		m.serialPort = serialPort
		// dtr/rts are applied by the registry at open; pacing lives inside
		// the shared port's write path.
		m.pacer = nil
	} else {
		serialPort, err = serialport.Open(m.SerialPortName, m.cfg.port)
		if err != nil {
			return errors.Wrap(err, "error connecting to serial port")
		}
		m.serialPort = serialPort
		m.applyModemLines(serialPort)
		m.pacer = newPacer(m.cfg.interframeDelay)
	}
	m.armer = newDeadlineReader(m.serialPort, m.cfg.readTimeout, m.pacer)
	m.reader = bufio.NewReader(m.armer)
	m.connected.Store(true)
	return nil
}

// applyModemLines asserts DTR/RTS when configured. Failures are warnings,
// not connection errors: ptys and adapters without modem lines must stay
// usable, and a missing line rarely prevents communication.
func (m *TransportInstance) applyModemLines(port serialport.Port) {
	controlPort, ok := port.(serialport.ControlPort)
	if !ok {
		return
	}
	if m.cfg.dtr {
		if err := controlPort.SetDTR(true); err != nil {
			m.log.Warn().Err(err).Msg("could not assert DTR")
		}
	}
	if m.cfg.rts {
		if err := controlPort.SetRTS(true); err != nil {
			m.log.Warn().Err(err).Msg("could not assert RTS")
		}
	}
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
	m.serialPort = nil
	m.armer = nil
	m.pacer = nil
	m.connected.Store(false)
	if err != nil {
		return errors.Wrap(err, "error closing serial port")
	}
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
	if !m.connected.Load() {
		return errors.New("error writing to transport. Not connected")
	}
	if m.serialPort == nil {
		return errors.New("error writing to transport. No writer available")
	}
	var deadline time.Time
	if ctxDeadline, ok := ctx.Deadline(); ok {
		deadline = ctxDeadline
	} else if m.cfg.writeTimeout > 0 {
		deadline = time.Now().Add(m.cfg.writeTimeout)
	}
	if err := m.pacer.waitTurn(deadline); err != nil {
		return errors.Wrap(err, "error pacing write")
	}
	if err := m.serialPort.SetWriteDeadline(deadline); err != nil {
		return errors.Wrap(err, "error setting write deadline")
	}
	num, err := m.serialPort.Write(data)
	m.pacer.noteActivity()
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
	armer := m.armer
	if armer == nil {
		return errors.New("error setting read deadline. Not connected")
	}
	armer.setExplicitDeadline(deadline)
	return nil
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
