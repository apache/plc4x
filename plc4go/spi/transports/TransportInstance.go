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

package transports

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type TransportInstance interface {
	fmt.Stringer
	Connect() error
	ConnectWithContext(ctx context.Context) error
	Close() error

	IsConnected() bool

	// FillBuffer fills the buffer `until` false (Useful in conjunction if you want GetNumBytesAvailableInBuffer)
	FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error
	// GetNumBytesAvailableInBuffer returns the bytes currently available in buffer (!!!Careful: if you looking for a termination you have to use FillBuffer)
	GetNumBytesAvailableInBuffer() (uint32, error)
	PeekReadableBytes(numBytes uint32) ([]byte, error)
	Read(numBytes uint32) ([]byte, error)

	Write(data []byte) error
}

type DefaultBufferedTransportInstanceRequirements interface {
	GetReader() *bufio.Reader
	Connect() error
}

type DefaultBufferedTransportInstance interface {
	ConnectWithContext(ctx context.Context) error
	GetNumBytesAvailableInBuffer() (uint32, error)
	FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error
	PeekReadableBytes(numBytes uint32) ([]byte, error)
	Read(numBytes uint32) ([]byte, error)
}

func NewDefaultBufferedTransportInstance(defaultBufferedTransportInstanceRequirements DefaultBufferedTransportInstanceRequirements) DefaultBufferedTransportInstance {
	return &defaultBufferedTransportInstance{defaultBufferedTransportInstanceRequirements}
}

type defaultBufferedTransportInstance struct {
	DefaultBufferedTransportInstanceRequirements
}

// ConnectWithContext is a compatibility implementation for those transports not implementing this function
func (m *defaultBufferedTransportInstance) ConnectWithContext(ctx context.Context) error {
	ch := make(chan error, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal().Interface("err", err).Msg("connect panic-ed")
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

func (m *defaultBufferedTransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error {
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
