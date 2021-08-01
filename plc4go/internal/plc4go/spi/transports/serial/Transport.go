//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package serial

import (
	"bufio"
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/jacobsa/go-serial/serial"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"
)

type Transport struct {
}

func NewTransport() *Transport {
	return &Transport{}
}

func (m Transport) GetTransportCode() string {
	return "serial"
}

func (m Transport) GetTransportName() string {
	return "Serial Transport"
}

func (m Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
	return m.CreateTransportInstanceForLocalAddress(transportUrl, options, nil)
}

func (m Transport) CreateTransportInstanceForLocalAddress(transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr) (transports.TransportInstance, error) {
	var serialPortName = transportUrl.Path

	var baudRate = uint(115200)
	if val, ok := options["baud-rate"]; ok {
		parsedBaudRate, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "error setting connect-timeout")
		} else {
			baudRate = uint(parsedBaudRate)
		}
	}

	var connectTimeout uint32 = 1000
	if val, ok := options["connect-timeout"]; ok {
		parsedConnectTimeout, err := strconv.ParseUint(val[0], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "error setting connect-timeout")
		} else {
			connectTimeout = uint32(parsedConnectTimeout)
		}
	}

	transportInstance := NewTransportInstance(serialPortName, baudRate, connectTimeout, &m)

	castFunc := func(typ interface{}) (transports.TransportInstance, error) {
		if transportInstance, ok := typ.(transports.TransportInstance); ok {
			return transportInstance, nil
		}
		return nil, errors.New("couldn't cast to TransportInstance")
	}
	return castFunc(transportInstance)
}

type TransportInstance struct {
	SerialPortName string
	BaudRate       uint
	ConnectTimeout uint32
	transport      *Transport
	serialPort     io.ReadWriteCloser
	reader         *bufio.Reader
}

func NewTransportInstance(serialPortName string, baudRate uint, connectTimeout uint32, transport *Transport) *TransportInstance {
	return &TransportInstance{
		SerialPortName: serialPortName,
		BaudRate:       baudRate,
		ConnectTimeout: connectTimeout,
		transport:      transport,
	}
}

func (m *TransportInstance) Connect() error {
	var err error
	config := serial.OpenOptions{PortName: m.SerialPortName, BaudRate: m.BaudRate, DataBits: 8, StopBits: 1, MinimumReadSize: 1}
	fmt.Printf("%v", config)
	m.serialPort, err = serial.Open(config)
	if err != nil {
		return errors.Wrap(err, "error connecting to serial port")
	}

	m.reader = bufio.NewReader(m.serialPort)

	return nil
}

func (m *TransportInstance) Close() error {
	if m.serialPort == nil {
		return nil
	}
	err := m.serialPort.Close()
	if err != nil {
		return errors.Wrap(err, "error closing serial port")
	}
	m.serialPort = nil
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.serialPort != nil
}

func (m *TransportInstance) GetNumReadableBytes() (uint32, error) {
	if m.reader == nil {
		return 0, nil
	}
	peekChan := make(chan bool)
	go func() {
		_, _ = m.reader.Peek(1)
		peekChan <- true
	}()
	select {
	case <-peekChan:
		return uint32(m.reader.Buffered()), nil
	case <-time.After(10 * time.Millisecond):
		return 0, nil
	}
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
	if m.reader == nil {
		return nil, errors.New("error peeking from transport. No reader available")
	}
	return m.reader.Peek(int(numBytes))
}

func (m *TransportInstance) Read(numBytes uint32) ([]uint8, error) {
	if m.reader == nil {
		return nil, errors.New("error reading from transport. No reader available")
	}
	data := make([]uint8, numBytes)
	for i := uint32(0); i < numBytes; i++ {
		val, err := m.reader.ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}

func (m *TransportInstance) Write(data []uint8) error {
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
