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
	"context"
	"fmt"
	"io"
)

type TransportInstance interface {
	fmt.Stringer
	io.Closer

	// Connect connects the transport instance
	Connect(ctx context.Context) error

	// IsConnected returns true if the transport is connected
	IsConnected() bool

	// FillBuffer fills the buffer `until` false (Useful in conjunction if you want GetNumBytesAvailableInBuffer)
	FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader ExtendedReader) (keepGoing bool)) error
	// GetNumBytesAvailableInBuffer returns the bytes currently available in buffer (!!!Careful: if you looking for a termination you have to use FillBuffer)
	GetNumBytesAvailableInBuffer() (uint32, error)
	// PeekReadableBytes returns the bytes currently available in buffer without consuming them
	PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error)
	// Read consumes the bytes from the transport
	Read(ctx context.Context, numBytes uint32) ([]byte, error)
	// Write writes data to the transport
	Write(ctx context.Context, data []byte) error
	// Reset resets the transport instance
	Reset()
}
