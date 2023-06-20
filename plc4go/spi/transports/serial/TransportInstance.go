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
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	transportUtils "github.com/apache/plc4x/plc4go/spi/transports/utils"

	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type TransportInstance struct {
	transportUtils.DefaultBufferedTransportInstance

	SerialPortName string
	BaudRate       uint
	ConnectTimeout uint32

	connected        atomic.Bool
	stateChangeMutex sync.Mutex

	transport  *Transport
	serialPort io.ReadWriteCloser
	reader     *bufio.Reader

	log zerolog.Logger
}

func NewTransportInstance(serialPortName string, baudRate uint, connectTimeout uint32, transport *Transport, _options ...options.WithOption) *TransportInstance {
	transportInstance := &TransportInstance{
		SerialPortName: serialPortName,
		BaudRate:       baudRate,
		ConnectTimeout: connectTimeout,
		transport:      transport,

		log: options.ExtractCustomLogger(_options...),
	}
	transportInstance.DefaultBufferedTransportInstance = transportUtils.NewDefaultBufferedTransportInstance(transportInstance, _options...)
	return transportInstance
}

func (m *TransportInstance) Connect() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if m.connected.Load() {
		return errors.New("Already connected")
	}

	var err error
	config := serial.OpenOptions{PortName: m.SerialPortName, BaudRate: m.BaudRate, DataBits: 8, StopBits: 1, MinimumReadSize: 0, InterCharacterTimeout: 100 /*, RTSCTSFlowControl: true*/}
	m.serialPort, err = serial.Open(config)
	if err != nil {
		return errors.Wrap(err, "error connecting to serial port")
	}
	// Add a logging layer ...
	/*logFile, err := ioutil.TempFile(os.TempDir(), "transport-logger")
	if err != nil {
		m.log.Error().Msg("Error creating file for logging transport requests")
	} else {
		fileLogger := zerolog.New(logFile).With().Logger()
		m.serialPort = utils.NewTransportLogger(m.serialPort, utils.WithLogger(fileLogger))
		m.log.Trace().Msgf("Logging Transport to file %s", logFile.Name())
	}*/
	m.reader = bufio.NewReader(m.serialPort)

	return nil
}

func (m *TransportInstance) Close() error {
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

func (m *TransportInstance) Write(data []byte) error {
	if !m.connected.Load() {
		return errors.New("error writing to transport. Not connected")
	}
	if m.serialPort == nil {
		return errors.New("error writing to transport. No writer available")
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

func (m *TransportInstance) String() string {
	return fmt.Sprintf("serial:%s:%d", m.SerialPortName, m.BaudRate)
}
