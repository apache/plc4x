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
	"context"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TransportInstance struct {
	readChannel chan []byte
	readBuffer  []byte
	writeBuffer []byte
	dataMutex   sync.RWMutex

	transport *Transport

	writeInterceptor func(transportInstance *TransportInstance, data []byte)
	simulatedLatency time.Duration

	connected        atomic.Bool
	stateChangeMutex sync.RWMutex

	log zerolog.Logger
}

var _ transports.TransportInstance = (*TransportInstance)(nil)

func NewTransportInstance(transport *Transport, _options ...options.WithOption) *TransportInstance {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	shouldTrace, found := ExtractTraceTransportInstance(_options...)
	if found && !shouldTrace {
		if customLogger.GetLevel() < zerolog.InfoLevel {
			customLogger = customLogger.Level(zerolog.InfoLevel)
		}
	}
	simulatedLatency, found := ExtractSimulatedLatency(_options...)
	if !found {
		simulatedLatency = 10 * time.Millisecond
	}
	return &TransportInstance{
		readChannel: make(chan []byte, 10),
		readBuffer:  []byte{},
		writeBuffer: []byte{},
		transport:   transport,

		simulatedLatency: simulatedLatency,

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

// WithSimulatedLatency adds simulated latency to the transport instance
func WithSimulatedLatency(latency time.Duration) options.WithOption {
	return withSimulatedLatency{latency: latency}
}

// ExtractSimulatedLatency to extract the simulated latency of the transport instance
func ExtractSimulatedLatency(options ...options.WithOption) (latency time.Duration, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withSimulatedLatency:
			latency, found = option.latency, true
		}
	}
	return
}

type withSimulatedLatency struct {
	options.Option
	latency time.Duration
}

func (m *TransportInstance) Connect(_ context.Context) error {
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

func (m *TransportInstance) Close() error {
	defer utils.StopWarn(m.log)()
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
	return m.availableBytes(), nil
}

func (m *TransportInstance) FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool)) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	m.log.Trace().Msg("Fill the buffer")
	nBytes := uint32(1)
	for ctx.Err() == nil {
		m.log.Trace().Dur("simulatedLatency", m.simulatedLatency).Msg("Sleeping simulatedLatency")
		timer := time.NewTimer(m.simulatedLatency)
		select {
		case <-ctx.Done():
		case <-timer.C:
		}
		m.log.Trace().Uint32("nBytes", nBytes).Msg("Peeking bytes")
		_bytes, err := m.PeekReadableBytes(ctx, nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		m.log.Trace().Msg("calling until callback")
		if keepGoing := until(uint(nBytes-1), _bytes[len(_bytes)-1], &transportInstanceDrivenExtendedReader{m, ctx}); !keepGoing {
			m.log.Trace().Uint32("nBytes", nBytes).Msg("Stopped after nBytes")
			return nil
		}
		m.log.Trace().Uint32("nBytes", nBytes).Msg("Keep going")
		nBytes++
	}
	return errors.Wrap(ctx.Err(), "Timeout while filling buffer")
}
func (m *TransportInstance) PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	m.log.Trace().Dur("simulatedLatency", m.simulatedLatency).Msg("Sleeping simulatedLatency")
	timer := time.NewTimer(m.simulatedLatency)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
	}
	availableBytes := m.availableBytes()
	if availableBytes < numBytes {
		m.log.Trace().Uint32("numBytes", numBytes).Uint32("availableBytes", availableBytes).Msg("Trying transfer now")
		availableBytes = m.transferFromChannel(ctx)
	} else {
		m.log.Trace().Msg("enough bytes available")
	}
	m.log.Trace().
		Uint32("numBytes", numBytes).
		Uint32("availableBytes", availableBytes).
		Msg("Peek numBytes readable bytes (of availableBytes available)")
	var err error
	if availableBytes < numBytes {
		m.log.Trace().Msg("not enough bytes available")
		return m.readBuffer[:], bufio.ErrBufferFull
	}
	m.log.Trace().Msg("enough bytes available")
	peekAble := m.peek()
	m.log.Trace().Int("peekAbleLen", len(peekAble)).Msg("New buffer size peekAbleLen")
	return peekAble[:numBytes], err
}

func (m *TransportInstance) Read(ctx context.Context, numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	availableBytes := m.availableBytes()
	m.log.Trace().
		Uint32("numBytes", numBytes).
		Uint32("availableBytes", availableBytes).
		Msg("Read num bytes numBytes (of availableBytes available)")
	if availableBytes < 1 {
		return nil, errors.Errorf("Only %d bytes available. Requested %d", availableBytes, numBytes)
	}
	if availableBytes < numBytes {
		m.log.Trace().Uint32("numBytes", numBytes).Uint32("availableBytes", availableBytes).Msg("Trying transfer now")
		availableBytes = m.transferFromChannel(ctx)
	}
	m.log.Trace().Dur("simulatedLatency", m.simulatedLatency).Msg("Sleeping simulatedLatency")
	timer := time.NewTimer(m.simulatedLatency)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
	}
	return m.read(int(numBytes)), nil
}

func (m *TransportInstance) SetWriteInterceptor(writeInterceptor func(transportInstance *TransportInstance, data []byte)) {
	m.log.Trace().Msg("Set write interceptor")
	m.writeInterceptor = writeInterceptor
}

func (m *TransportInstance) Write(ctx context.Context, data []byte) error {
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
	m.log.Trace().Dur("simulatedLatency", m.simulatedLatency).Msg("Sleeping simulatedLatency")
	timer := time.NewTimer(m.simulatedLatency)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	m.log.Trace().
		Hex("data", data).
		Str("hexDump", hex.Dump(data)).
		Msg("Write data to write buffer")
	m.dataMutex.Lock()
	m.writeBuffer = append(m.writeBuffer, data...)
	m.dataMutex.Unlock()
	return nil
}

func (m *TransportInstance) FillReadBuffer(data []byte) {
	if !m.IsConnected() {
		m.log.Error().Msg("working on a unconnected connection")
		return
	}
	m.log.Trace().
		Int("nBytes", len(data)).
		Int("existingBytes", len(m.readBuffer)).
		Int("readChannelSize", len(m.readChannel)).
		Str("hexDump", hex.Dump(data)).
		Msg("fill read buffer with hexDump (nBytes bytes). (Adding to existingBytes bytes existing)")
	m.readChannel <- data
}

func (m *TransportInstance) GetNumDrainableBytes() uint32 {
	if !m.IsConnected() {
		m.log.Error().Msg("working on a unconnected connection")
		return 0
	}
	m.dataMutex.RLock()
	m.log.Trace().Msg("get number of drainable bytes")
	writeBufLen := uint32(len(m.writeBuffer))
	m.dataMutex.RUnlock()
	return writeBufLen
}

func (m *TransportInstance) DrainWriteBuffer(numBytes uint32) []byte {
	if !m.IsConnected() {
		m.log.Error().Msg("working on a unconnected connection")
		return nil
	}
	m.dataMutex.Lock()
	m.log.Trace().
		Uint32("numBytes", numBytes).
		Msg("Drain write buffer with number of bytes")
	data := m.writeBuffer[0:int(numBytes)]
	m.writeBuffer = m.writeBuffer[int(numBytes):]
	m.dataMutex.Unlock()
	return data
}

func (m *TransportInstance) String() string {
	return "test"
}

func (m *TransportInstance) availableBytes() uint32 {
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	return uint32(len(m.readBuffer))
}

func (m *TransportInstance) peek() []byte {
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	return m.readBuffer[0:len(m.readBuffer)]
}

func (m *TransportInstance) read(numBytes int) []byte {
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	data := m.readBuffer[0:numBytes]
	m.readBuffer = m.readBuffer[numBytes:]
	return data
}

func (m *TransportInstance) appendRead(newBytes ...byte) (totalAvailableBytes uint32) {
	m.dataMutex.Lock()
	defer m.dataMutex.Unlock()
	m.readBuffer = append(m.readBuffer, newBytes...)
	return uint32(len(m.readBuffer))
}

func (m *TransportInstance) transferFromChannel(ctx context.Context) (totalAvailableBytes uint32) {
	totalAvailableBytes = m.availableBytes()
	m.log.Trace().Msg("waiting for transfer")
	start := time.Now()
	select {
	case <-ctx.Done():
		m.log.Trace().Msg("Context done")
	case newBytes := <-m.readChannel:
		m.log.Trace().Dur("time", time.Since(start)).Int("nBytes", len(newBytes)).Msg("Got new bytes")
		totalAvailableBytes = m.appendRead(newBytes...)
	}
	return totalAvailableBytes
}
