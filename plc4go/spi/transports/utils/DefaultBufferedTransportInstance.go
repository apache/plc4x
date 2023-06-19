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
	"runtime/debug"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type DefaultBufferedTransportInstanceRequirements interface {
	GetReader() transports.ExtendedReader
	Connect() error
}

type DefaultBufferedTransportInstance interface {
	ConnectWithContext(ctx context.Context) error
	GetNumBytesAvailableInBuffer() (uint32, error)
	FillBuffer(until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool) error
	PeekReadableBytes(numBytes uint32) ([]byte, error)
	Read(numBytes uint32) ([]byte, error)
}

func NewDefaultBufferedTransportInstance(defaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements, _options ...options.WithOption) DefaultBufferedTransportInstance {
	return &defaultBufferedTransportInstance{
		DefaultBufferedTransportInstanceRequirements: defaultBufferedTransportInstanceRequirements,
		log: options.ExtractCustomLogger(_options...),
	}
}

type defaultBufferedTransportInstance struct {
	DefaultBufferedTransportInstanceRequirements

	log zerolog.Logger
}

// ConnectWithContext is a compatibility implementation for those transports not implementing this function
func (m *defaultBufferedTransportInstance) ConnectWithContext(ctx context.Context) error {
	ch := make(chan error, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			}
		}()
		ch <- m.Connect()
		close(ch)
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *defaultBufferedTransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if m.GetReader() == nil {
		return 0, nil
	}
	_, _ = m.GetReader().Peek(1)
	return uint32(m.GetReader().Buffered()), nil
}

func (m *defaultBufferedTransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader transports.ExtendedReader) bool) error {
	if m.GetReader() == nil {
		return nil
	}
	nBytes := uint32(1)
	for {
		bytes, err := m.PeekReadableBytes(nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		if keepGoing := until(uint(nBytes-1), bytes[len(bytes)-1], m.GetReader()); !keepGoing {
			return nil
		}
		nBytes++
	}
}

func (m *defaultBufferedTransportInstance) PeekReadableBytes(numBytes uint32) ([]byte, error) {
	if m.GetReader() == nil {
		return nil, errors.New("error peeking from transport. No reader available")
	}
	return m.GetReader().Peek(int(numBytes))
}

func (m *defaultBufferedTransportInstance) Read(numBytes uint32) ([]byte, error) {
	if m.GetReader() == nil {
		return nil, errors.New("error reading from transport. No reader available")
	}
	data := make([]byte, numBytes)
	for i := uint32(0); i < numBytes; i++ {
		val, err := m.GetReader().ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}
