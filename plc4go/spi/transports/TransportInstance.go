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
	"github.com/pkg/errors"
)

type TransportInstance interface {
	Connect() error
	Close() error

	IsConnected() bool

	// FillBuffer fills the buffer `until` false (Useful in conjunction if you want GetNumBytesAvailableInBuffer)
	FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error
	// GetNumBytesAvailableInBuffer returns the bytes currently available in buffer (!!!Careful: if you looking for a termination you have to use FillBuffer)
	GetNumBytesAvailableInBuffer() (uint32, error)
	PeekReadableBytes(numBytes uint32) ([]uint8, error)
	Read(numBytes uint32) ([]uint8, error)

	Write(data []uint8) error
}

type TestTransportInstance interface {
	TransportInstance
	FillReadBuffer(data []uint8) error
	GetNumDrainableBytes() uint32
	DrainWriteBuffer(numBytes uint32) ([]uint8, error)
}

type DefaultBufferedTransportInstance struct {
	*bufio.Reader
}

func (m *DefaultBufferedTransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	if m.Reader == nil {
		return 0, nil
	}
	_, _ = m.Peek(1)
	return uint32(m.Buffered()), nil
}

func (m *DefaultBufferedTransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error {
	if m.Reader == nil {
		return nil
	}
	nBytes := uint32(1)
	for {
		bytes, err := m.PeekReadableBytes(nBytes)
		if err != nil {
			return errors.Wrap(err, "Error while peeking")
		}
		if keepGoing := until(uint(nBytes-1), bytes[len(bytes)-1], m.Reader); !keepGoing {
			return nil
		}
		nBytes++
	}
}

func (m *DefaultBufferedTransportInstance) PeekReadableBytes(numBytes uint32) ([]uint8, error) {
	if m.Reader == nil {
		return nil, errors.New("error peeking from transport. No reader available")
	}
	return m.Peek(int(numBytes))
}

func (m *DefaultBufferedTransportInstance) Read(numBytes uint32) ([]uint8, error) {
	if m.Reader == nil {
		return nil, errors.New("error reading from transport. No reader available")
	}
	data := make([]uint8, numBytes)
	for i := uint32(0); i < numBytes; i++ {
		val, err := m.ReadByte()
		if err != nil {
			return nil, errors.Wrap(err, "error reading")
		}
		data[i] = val
	}
	return data, nil
}
