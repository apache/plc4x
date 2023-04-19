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

package _default

import (
	"bufio"
	"context"
)

type testTransportInstance struct {
}

func (t testTransportInstance) String() string {
	// NO-OP
	return ""
}

func (t testTransportInstance) Connect() error {
	// NO-OP
	return nil
}

func (t testTransportInstance) ConnectWithContext(ctx context.Context) error {
	// NO-OP
	return nil
}

func (t testTransportInstance) Close() error {
	// NO-OP
	return nil
}

func (t testTransportInstance) IsConnected() bool {
	// NO-OP
	return false
}

func (t testTransportInstance) FillBuffer(until func(pos uint, currentByte byte, reader *bufio.Reader) bool) error {
	// NO-OP
	return nil
}

func (t testTransportInstance) GetNumBytesAvailableInBuffer() (uint32, error) {
	// NO-OP
	return 0, nil
}

func (t testTransportInstance) PeekReadableBytes(numBytes uint32) ([]byte, error) {
	// NO-OP
	return nil, nil
}

func (t testTransportInstance) Read(numBytes uint32) ([]byte, error) {
	// NO-OP
	return nil, nil
}

func (t testTransportInstance) Write(data []byte) error {
	// NO-OP
	return nil
}
