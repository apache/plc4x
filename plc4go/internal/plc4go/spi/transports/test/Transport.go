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
package test

import (
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"net/url"
)

type Transport struct {
	transports.Transport
}

func NewTransport() *Transport {
	return &Transport{}
}

func (m Transport) GetTransportCode() string {
	return "test"
}

func (m Transport) GetTransportName() string {
	return "Test Transport"
}

func (m Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
	transportInstance := NewTransportInstance(&m)

	castFunc := func(typ interface{}) (transports.TransportInstance, error) {
		if transportInstance, ok := typ.(transports.TransportInstance); ok {
			return transportInstance, nil
		}
		return nil, errors.New("couldn't cast to TransportInstance")
	}
	return castFunc(transportInstance)
}

type TransportInstance struct {
	readBuffer  []byte
	writeBuffer []byte
	transport   *Transport
}

func NewTransportInstance(transport *Transport) *TransportInstance {
	return &TransportInstance{
		readBuffer:  []byte{},
		writeBuffer: []byte{},
		transport:   transport,
	}
}

func (m *TransportInstance) Connect() error {
	return nil
}

func (m *TransportInstance) Close() error {
	return nil
}

func (m *TransportInstance) GetNumReadableBytes() (uint32, error) {
	return uint32(len(m.readBuffer)), nil
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
	return m.readBuffer[0:numBytes], nil
}

func (m *TransportInstance) Read(numBytes uint32) ([]uint8, error) {
	data := m.readBuffer[0:int(numBytes)]
	m.readBuffer = m.readBuffer[int(numBytes):]
	return data, nil
}

func (m *TransportInstance) Write(data []uint8) error {
	m.writeBuffer = append(m.writeBuffer, data...)
	return nil
}

func (m *TransportInstance) FillReadBuffer(data []uint8) error {
	m.readBuffer = append(m.readBuffer, data...)
	return nil
}

func (m *TransportInstance) GetNumDrainableBytes() uint32 {
	return uint32(len(m.writeBuffer))
}

func (m *TransportInstance) DrainWriteBuffer(numBytes uint32) ([]uint8, error) {
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	return data, nil
}
