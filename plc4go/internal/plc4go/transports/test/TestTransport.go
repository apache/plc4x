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
	"net/url"
	"github.com/apache/plc4x/plc4go/internal/plc4go/transports"
)

type TestTransport struct {
	transports.Transport
}

func NewTestTransport() *TestTransport {
	return &TestTransport{}
}

func (m TestTransport) GetTransportCode() string {
	return "test"
}

func (m TestTransport) GetTransportName() string {
	return "Test Transport"
}

func (m TestTransport) CreateTransportInstance(transportUrl url.URL, options map[string][]string) (transports.TransportInstance, error) {
	transportInstance := NewTestTransportInstance(&m)

	castFunc := func(typ interface{}) (transports.TransportInstance, error) {
		if transportInstance, ok := typ.(transports.TransportInstance); ok {
			return transportInstance, nil
		}
		return nil, errors.New("couldn't cast to TransportInstance")
	}
	return castFunc(transportInstance)
}

type TestTransportInstance struct {
	readBuffer  []byte
	writeBuffer []byte
	transport   *TestTransport
}

func NewTestTransportInstance(transport *TestTransport) *TestTransportInstance {
	return &TestTransportInstance{
		readBuffer:  []byte{},
		writeBuffer: []byte{},
		transport:   transport,
	}
}

func (m *TestTransportInstance) Connect() error {
	return nil
}

func (m *TestTransportInstance) Close() error {
	return nil
}

func (m *TestTransportInstance) GetNumReadableBytes() (uint32, error) {
	return uint32(len(m.readBuffer)), nil
}

func (m *TestTransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
	return m.readBuffer[0:numBytes], nil
}

func (m *TestTransportInstance) Read(numBytes uint32) ([]uint8, error) {
	data := m.readBuffer[0:int(numBytes)]
	m.readBuffer = m.readBuffer[int(numBytes):]
	return data, nil
}

func (m *TestTransportInstance) Write(data []uint8) error {
	m.writeBuffer = append(m.writeBuffer, data...)
	return nil
}

func (m *TestTransportInstance) FillReadBuffer(data []uint8) error {
	m.readBuffer = append(m.readBuffer, data...)
	return nil
}

func (m *TestTransportInstance) GetNumDrainableBytes() uint32 {
	return uint32(len(m.writeBuffer))
}

func (m *TestTransportInstance) DrainWriteBuffer(numBytes uint32) ([]uint8, error) {
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	return data, nil
}
