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
)

type transportInstanceDrivenExtendedReader struct {
	*TransportInstance

	ctx context.Context
}

func (t *transportInstanceDrivenExtendedReader) Read(p []byte) (n int, err error) {
	numBytes := uint32(len(p))
	availableBytes := t.availableBytes()
	if availableBytes < numBytes {
		t.log.Trace().Uint32("numBytes", numBytes).Uint32("availableBytes", availableBytes).Msg("Trying transfer now")
		availableBytes = t.transferFromChannel(t.ctx)
	} else {
		t.log.Trace().Msg("enough bytes available")
	}
	return bufio.NewReader(bytes.NewReader(t.readBuffer)).Read(p)
}

func (t *transportInstanceDrivenExtendedReader) ReadByte() (byte, error) {
	numBytes := uint32(1)
	availableBytes := t.availableBytes()
	if availableBytes < numBytes {
		t.log.Trace().Uint32("numBytes", numBytes).Uint32("availableBytes", availableBytes).Msg("Trying transfer now")
		availableBytes = t.transferFromChannel(t.ctx)
	} else {
		t.log.Trace().Msg("enough bytes available")
	}
	return bufio.NewReader(bytes.NewReader(t.readBuffer)).ReadByte()
}

func (t *transportInstanceDrivenExtendedReader) Peek(n int) ([]byte, error) {
	numBytes := uint32(n)
	availableBytes := t.availableBytes()
	if availableBytes < numBytes {
		t.log.Trace().Uint32("numBytes", numBytes).Uint32("availableBytes", availableBytes).Msg("Trying transfer now")
		availableBytes = t.transferFromChannel(t.ctx)
	} else {
		t.log.Trace().Msg("enough bytes available")
	}
	return bufio.NewReader(bytes.NewReader(t.readBuffer)).Peek(n)
}

func (t *transportInstanceDrivenExtendedReader) Buffered() int {
	return len(t.readBuffer)
}
