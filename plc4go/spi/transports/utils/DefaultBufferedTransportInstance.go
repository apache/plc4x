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

package utils

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

type DefaultBufferedTransportInstanceRequirements interface {
	GetReader() transports.ExtendedReader
	IsConnected() bool
	SetReadDeadline(deadline time.Time) error
}

type DefaultBufferedTransportInstance interface {
	GetNumBytesAvailableInBuffer() (uint32, error)
	FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool)) error
	PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error)
	Read(ctx context.Context, numBytes uint32) ([]byte, error)
}

func NewDefaultBufferedTransportInstance(defaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements, _options ...options.WithOption) DefaultBufferedTransportInstance {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &defaultBufferedTransportInstance{
		DefaultBufferedTransportInstanceRequirements: defaultBufferedTransportInstanceRequirements,
		log: customLogger,
	}
}

type defaultBufferedTransportInstance struct {
	DefaultBufferedTransportInstanceRequirements

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func (m *defaultBufferedTransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if !m.IsConnected() {
		return 0, errors.New("working on a unconnected connection")
	}
	if m.GetReader() == nil {
		return 0, nil
	}
	_, _ = m.GetReader().Peek(1)
	return uint32(m.GetReader().Buffered()), nil
}

func (m *defaultBufferedTransportInstance) FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool) error {
	if !m.IsConnected() {
		return errors.New("working on a unconnected connection")
	}
	if m.GetReader() == nil {
		return nil
	}
	nBytes := uint32(1)
	for ctx.Err() == nil {
		bytes, err := m.PeekReadableBytes(ctx, nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		if keepGoing := until(uint(nBytes-1), bytes[len(bytes)-1], m.GetReader()); !keepGoing {
			return nil
		}
		nBytes++
	}
	return errors.Wrap(ctx.Err(), "Timeout while filling buffer")
}

func (m *defaultBufferedTransportInstance) PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	if m.GetReader() == nil {
		return nil, errors.New("error peeking from transport. No reader available")
	}
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.SetReadDeadline(deadline); err != nil {
			return nil, errors.Wrap(err, "error setting read deadline")
		}
	}
	return m.GetReader().Peek(int(numBytes))
}

func (m *defaultBufferedTransportInstance) Read(ctx context.Context, numBytes uint32) ([]byte, error) {
	if !m.IsConnected() {
		return nil, errors.New("working on a unconnected connection")
	}
	if m.GetReader() == nil {
		return nil, errors.New("error reading from transport. No reader available")
	}
	if deadline, ok := ctx.Deadline(); ok {
		m.log.Trace().Time("deadline", deadline).Msg("deadline set")
		if err := m.SetReadDeadline(deadline); err != nil {
			return nil, errors.Wrap(err, "error setting read deadline")
		}
	}
	data := make([]byte, numBytes)
	for i := range numBytes {
		val, err := m.GetReader().ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}
