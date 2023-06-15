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
	"encoding/hex"
	"math"
	"sync"

	"github.com/apache/plc4x/plc4go/spi/options"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type TransportInstance struct {
	readBuffer       []byte
	writeBuffer      []byte
	connected        bool
	transport        *Transport
	writeInterceptor func(transportInstance *TransportInstance, data []byte)

	dataMutex        sync.RWMutex
	stateChangeMutex sync.RWMutex

	log zerolog.Logger
}

func NewTransportInstance(transport *Transport, _options ...options.WithOption) *TransportInstance {
	return &TransportInstance{
		readBuffer:  []byte{},
		writeBuffer: []byte{},
		connected:   false,
		transport:   transport,

		log: options.ExtractCustomLogger(_options...),
	}
}

func (m *TransportInstance) Connect() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if m.connected {
		m.log.Warn().Msg("already connected")
		return nil
	}
	m.log.Trace().Msg("Connect")
	m.connected = true
	return nil
}

func (m *TransportInstance) ConnectWithContext(_ context.Context) error {
	return m.Connect()
}

func (m *TransportInstance) Close() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	m.log.Trace().Msg("Close")
	m.connected = false
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	m.stateChangeMutex.RLock()
	defer m.stateChangeMutex.RUnlock()
	return m.connected
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	readableBytes := len(m.readBuffer)
	m.log.Trace().Msgf("return number of readable bytes %d", readableBytes)
	return uint32(readableBytes), nil
}

func (m *TransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error {
	m.log.Trace().Msg("Fill the buffer")
	nBytes := uint32(1)
	for {
		m.log.Trace().Msgf("Peeking %d bytes", nBytes)
		_bytes, err := m.PeekReadableBytes(nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		m.dataMutex.RLock()
		reader := bufio.NewReader(bytes.NewReader(m.readBuffer))
		if keepGoing := until(uint(nBytes-1), _bytes[len(_bytes)-1], reader); !keepGoing {
			m.log.Trace().Msgf("Stopped after %d bytes", nBytes)
			m.dataMutex.RUnlock()
			return nil
		}
		m.dataMutex.RUnlock()
		nBytes++
	}
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]byte, error) {
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	availableBytes := uint32(math.Min(float64(numBytes), float64(len(m.readBuffer))))
	m.log.Trace().Msgf("Peek %d readable bytes (%d available bytes)", numBytes, availableBytes)
	var err error
	if availableBytes != numBytes {
		err = errors.New("not enough bytes available")
	}
	if availableBytes == 0 {
		m.log.Trace().Msg("No bytes available")
		return nil, err
	}
	return m.readBuffer[0:availableBytes], err
}

func (m *TransportInstance) Read(numBytes uint32) ([]byte, error) {
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	nBytes := len(m.readBuffer)
	m.log.Trace().Msgf("Read num bytes %d (of %d available)", numBytes, nBytes)
	if nBytes < 1 {
		return nil, errors.Errorf("Only %d bytes available. Requested %d", nBytes, numBytes)
	}
	data := m.readBuffer[0:int(numBytes)]
	m.readBuffer = m.readBuffer[int(numBytes):]
	m.log.Trace().Msgf("New buffer size %d", nBytes)
	return data, nil
}

func (m *TransportInstance) SetWriteInterceptor(writeInterceptor func(transportInstance *TransportInstance, data []byte)) {
	m.writeInterceptor = writeInterceptor
}

func (m *TransportInstance) Write(data []byte) error {
	if m.writeInterceptor != nil {
		m.log.Trace().Msgf("Passing data to write interceptor\n%s", hex.Dump(data))
		m.writeInterceptor(m, data)
	}
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.log.Trace().Msgf("Write data to write buffer\n%s", hex.Dump(data))
	m.writeBuffer = append(m.writeBuffer, data...)
	return nil
}

func (m *TransportInstance) FillReadBuffer(data []byte) {
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.log.Trace().Msgf("fill read buffer with \n%s (%d bytes). (Adding to %d bytes existing)", hex.Dump(data), len(data), len(m.readBuffer))
	m.readBuffer = append(m.readBuffer, data...)
}

func (m *TransportInstance) GetNumDrainableBytes() uint32 {
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	m.log.Trace().Msg("get number of drainable bytes")
	return uint32(len(m.writeBuffer))
}

func (m *TransportInstance) DrainWriteBuffer(numBytes uint32) []byte {
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.log.Trace().Msgf("Drain write buffer with number of bytes %d", numBytes)
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	return data
}

func (m *TransportInstance) String() string {
	return "test" //TODO: maybe use plc4xgen
}
