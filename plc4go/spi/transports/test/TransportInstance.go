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
	"sync/atomic"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type TransportInstance struct {
	readBuffer       []byte
	writeBuffer      []byte
	transport        *Transport
	writeInterceptor func(transportInstance *TransportInstance, data []byte)

	dataMutex        sync.RWMutex
	connected        atomic.Bool
	stateChangeMutex sync.RWMutex

	log zerolog.Logger
}

func NewTransportInstance(transport *Transport, _options ...options.WithOption) *TransportInstance {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	shouldTrace, found := ExtractTraceTransportInstance(_options...)
	if found && !shouldTrace {
		if customLogger.GetLevel() < zerolog.InfoLevel {
			customLogger = customLogger.Level(zerolog.InfoLevel)
		}
	}
	return &TransportInstance{
		readBuffer:  []byte{},
		writeBuffer: []byte{},
		transport:   transport,

		log: customLogger,
	}
}

// WithTraceTransportInstance enables tracing of the test transport instance
func WithTraceTransportInstance(trace bool) options.WithOption {
	return withTraceTransportInstance{trace: trace}
}

// ExtractTraceTransportInstance to extract the flag indicating that transport instance should be traced
func ExtractTraceTransportInstance(options ...options.WithOption) (trace bool, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withTraceTransportInstance:
			trace, found = option.trace, true
		}
	}
	return
}

type withTraceTransportInstance struct {
	options.Option
	trace bool
}

func (m *TransportInstance) Connect() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if m.connected.Load() {
		m.log.Warn().Msg("already connected")
		return nil
	}
	m.log.Trace().Msg("Connect")
	m.connected.Store(true)
	return nil
}

func (m *TransportInstance) ConnectWithContext(_ context.Context) error {
	return m.Connect()
}

func (m *TransportInstance) Close() error {
	m.stateChangeMutex.Lock()
	defer m.stateChangeMutex.Unlock()
	if !m.connected.Load() {
		return nil
	}
	m.log.Trace().Msg("Close")
	m.connected.Store(true)
	return nil
}

func (m *TransportInstance) IsConnected() bool {
	return m.connected.Load()
}

func (m *TransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if !m.IsConnected() {
		panic(errors.New("working on a unconnected connection"))
	}
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	readableBytes := len(m.readBuffer)
	m.log.Trace().Int("readableBytes", readableBytes).Msg("return number of readable bytes")
	return uint32(readableBytes), nil
}

func (m *TransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	m.log.Trace().Msg("Fill the buffer")
	nBytes := uint32(1)
	for {
		m.log.Trace().Uint32("nBytes", nBytes).Msg("Peeking bytes")
		_bytes, err := m.PeekReadableBytes(nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		m.dataMutex.RLock()
		reader := bufio.NewReader(bytes.NewReader(m.readBuffer))
		if keepGoing := until(uint(nBytes-1), _bytes[len(_bytes)-1], reader); !keepGoing {
			m.log.Trace().Uint32("nBytes", nBytes).Msg("Stopped after nBytes")
			m.dataMutex.RUnlock()
			return nil
		}
		m.dataMutex.RUnlock()
		nBytes++
	}
}

func (m *TransportInstance) PeekReadableBytes(numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	availableBytes := uint32(math.Min(float64(numBytes), float64(len(m.readBuffer))))
	m.log.Trace().
		Uint32("numBytes", numBytes).
		Uint32("availableBytes", availableBytes).
		Msg("Peek numBytes readable bytes (of availableBytes available)")
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
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	availableBytes := uint32(math.Min(float64(numBytes), float64(len(m.readBuffer))))
	m.log.Trace().
		Uint32("numBytes", numBytes).
		Uint32("availableBytes", availableBytes).
		Msg("Read num bytes numBytes (of availableBytes available)")
	if availableBytes < 1 {
		return nil, errors.Errorf("Only %d bytes available. Requested %d", availableBytes, numBytes)
	}
	data := m.readBuffer[0:int(numBytes)]
	m.readBuffer = m.readBuffer[int(numBytes):]
	m.log.Trace().Uint32("availableBytes", availableBytes).Msg("New buffer size availableBytes")
	return data, nil
}

func (m *TransportInstance) SetWriteInterceptor(writeInterceptor func(transportInstance *TransportInstance, data []byte)) {
	m.log.Trace().Msg("Set write interceptor")
	m.writeInterceptor = writeInterceptor
}

func (m *TransportInstance) Write(data []byte) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	if m.writeInterceptor != nil {
		m.log.Trace().
			Hex("data", data).
			Str("hexDump", hex.Dump(data)).
			Msg("Passing data to write interceptor")
		m.writeInterceptor(m, data)
	}
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.log.Trace().
		Hex("data", data).
		Str("hexDump", hex.Dump(data)).
		Msg("Write data to write buffer")
	m.writeBuffer = append(m.writeBuffer, data...)
	return nil
}

func (m *TransportInstance) FillReadBuffer(data []byte) {
	if !m.IsConnected() {
		m.log.Error().Msg("working on a unconnected connection")
		return
	}
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.log.Trace().
		Int("nBytes", len(data)).
		Int("existingBytes", len(m.readBuffer)).
		Str("hexDump", hex.Dump(data)).
		Msg("fill read buffer with hexDump (nBytes bytes). (Adding to existingBytes bytes existing)")
	m.readBuffer = append(m.readBuffer, data...)
}

func (m *TransportInstance) GetNumDrainableBytes() uint32 {
	if !m.IsConnected() {
		m.log.Error().Msg("working on a unconnected connection")
		return 0
	}
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	m.log.Trace().Msg("get number of drainable bytes")
	return uint32(len(m.writeBuffer))
}

func (m *TransportInstance) DrainWriteBuffer(numBytes uint32) []byte {
	if !m.IsConnected() {
		m.log.Error().Msg("working on a unconnected connection")
		return nil
	}
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.log.Trace().
		Uint32("numBytes", numBytes).
		Msg("Drain write buffer with number of bytes")
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	return data
}

func (m *TransportInstance) String() string {
	return "test"
}
