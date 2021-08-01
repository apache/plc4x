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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Transport struct {
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
	log.Trace().Msg("create transport instance")
	transportInstance := NewTransportInstance(&m)

	castFunc := func(typ interface{}) (transports.TransportInstance, error) {
		if transportInstance, ok := typ.(transports.TransportInstance); ok {
			return transportInstance, nil
		}
		return nil, errors.Errorf("couldn't cast to TransportInstance. Actual instance: %T", typ)
	}
	return castFunc(transportInstance)
}

type TransportInstance struct {
	readBuffer  []byte
	writeBuffer []byte
	connected	bool
	transport   *Transport
}

func NewTransportInstance(transport *Transport) *TransportInstance {
	return &TransportInstance{
		readBuffer:  []byte{},
		writeBuffer: []byte{},
		connected:   false,
		transport:   transport,
	}
}

func (m *TransportInstance) Connect() error {
	log.Trace().Msg("Connect")
	m.connected = true
	return nil
}

func (m *TransportInstance) Close() error {
	log.Trace().Msg("Close")
	m.connected = false
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected
}
func (m *TransportInstance) GetNumReadableBytes() (uint32, error) {
	readableBytes := len(m.readBuffer)
	log.Trace().Msgf("return number of readable bytes %d", readableBytes)
	return uint32(readableBytes), nil
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
	log.Trace().Msgf("Peek %d readable bytes", numBytes)
	return m.readBuffer[0:numBytes], nil
}

func (m *TransportInstance) Read(numBytes uint32) ([]uint8, error) {
	log.Trace().Msgf("Read num bytes %d", numBytes)
	data := m.readBuffer[0:int(numBytes)]
	m.readBuffer = m.readBuffer[int(numBytes):]
	return data, nil
}

func (m *TransportInstance) Write(data []uint8) error {
	log.Trace().Msgf("Write data 0x%x", data)
	m.writeBuffer = append(m.writeBuffer, data...)
	return nil
}

func (m *TransportInstance) FillReadBuffer(data []uint8) error {
	log.Trace().Msgf("FillReadBuffer with 0x%x", data)
	m.readBuffer = append(m.readBuffer, data...)
	return nil
}

func (m *TransportInstance) GetNumDrainableBytes() uint32 {
	log.Trace().Msg("get number of drainable bytes")
	return uint32(len(m.writeBuffer))
}

func (m *TransportInstance) DrainWriteBuffer(numBytes uint32) ([]uint8, error) {
	log.Trace().Msgf("Drain write buffer with number of bytes %d", numBytes)
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	return data, nil
}
