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

package test

import (
	"bufio"
	"bytes"
	"context"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"net/url"
)

type Transport struct {
	preregisteredInstances map[url.URL]transports.TransportInstance
}

func NewTransport() *Transport {
	return &Transport{preregisteredInstances: map[url.URL]transports.TransportInstance{}}
}

func (m *Transport) GetTransportCode() string {
	return "test"
}

func (m *Transport) GetTransportName() string {
	return "Test Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, _ map[string][]string) (transports.TransportInstance, error) {
	if preregisteredInstance, ok := m.preregisteredInstances[transportUrl]; ok {
		log.Trace().Msgf("Returning pre registered instance for %v", transportUrl)
		return preregisteredInstance, nil
	}
	log.Trace().Msg("create transport instance")
	return NewTransportInstance(m), nil
}

func (m *Transport) AddPreregisteredInstances(transportUrl url.URL, preregisteredInstance transports.TransportInstance) error {
	if _, ok := m.preregisteredInstances[transportUrl]; ok {
		return errors.Errorf("registered instance for %v already registered", transportUrl)
	}
	m.preregisteredInstances[transportUrl] = preregisteredInstance
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}

type TransportInstance struct {
	readBuffer       []byte
	writeBuffer      []byte
	connected        bool
	transport        *Transport
	writeInterceptor func(transportInstance *TransportInstance, data []byte)
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

func (m *TransportInstance) ConnectWithContext(_ context.Context) error {
	return m.Connect()
}

func (m *TransportInstance) Close() error {
	log.Trace().Msg("Close")
	m.connected = false
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	readableBytes := len(m.readBuffer)
	log.Trace().Msgf("return number of readable bytes %d", readableBytes)
	return uint32(readableBytes), nil
}

func (m *TransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error {
	nBytes := uint32(1)
	for {
		_bytes, err := m.PeekReadableBytes(nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		if keepGoing := until(uint(nBytes-1), _bytes[len(_bytes)-1], bufio.NewReader(bytes.NewReader(m.readBuffer))); !keepGoing {
			return nil
		}
		nBytes++
	}
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]byte, error) {
	log.Trace().Msgf("Peek %d readable bytes", numBytes)
	availableBytes := uint32(math.Min(float64(numBytes), float64(len(m.readBuffer))))
	var err error
	if availableBytes != numBytes {
		err = errors.New("not enough bytes available")
	}
	if availableBytes == 0 {
		return nil, err
	}
	return m.readBuffer[0:availableBytes], nil
}

func (m *TransportInstance) Read(numBytes uint32) ([]byte, error) {
	log.Trace().Msgf("Read num bytes %d", numBytes)
	data := m.readBuffer[0:int(numBytes)]
	m.readBuffer = m.readBuffer[int(numBytes):]
	return data, nil
}

func (m *TransportInstance) SetWriteInterceptor(writeInterceptor func(transportInstance *TransportInstance, data []byte)) {
	m.writeInterceptor = writeInterceptor
}

func (m *TransportInstance) Write(data []byte) error {
	if m.writeInterceptor != nil {
		m.writeInterceptor(m, data)
	}
	log.Trace().Msgf("Write data %#x", data)
	m.writeBuffer = append(m.writeBuffer, data...)
	return nil
}

func (m *TransportInstance) FillReadBuffer(data []byte) {
	log.Trace().Msgf("FillReadBuffer with %#x", data)
	m.readBuffer = append(m.readBuffer, data...)
}

func (m *TransportInstance) GetNumDrainableBytes() uint32 {
	log.Trace().Msg("get number of drainable bytes")
	return uint32(len(m.writeBuffer))
}

func (m *TransportInstance) DrainWriteBuffer(numBytes uint32) []byte {
	log.Trace().Msgf("Drain write buffer with number of bytes %d", numBytes)
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	return data
}

func (m *TransportInstance) String() string {
	return "test"
}
