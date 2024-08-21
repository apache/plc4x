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

package tests

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

var StartTime = time.Time{}.Add(1 * time.Hour)

type DummyMessage struct {
	Data []byte
}

func NewDummyMessage(data ...byte) *DummyMessage {
	return &DummyMessage{Data: data}
}

func (d DummyMessage) String() string {
	return hex.EncodeToString(d.Data)
}

func (d DummyMessage) Serialize() ([]byte, error) {
	return d.Data, nil
}

func (d DummyMessage) SerializeWithWriteBuffer(_ context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteByteArray("data", d.Data)
}

func (d DummyMessage) GetLengthInBytes(_ context.Context) uint16 {
	return uint16(len(d.Data))
}

func (d DummyMessage) GetLengthInBits(_ context.Context) uint16 {
	return uint16(len(d.Data))
}

type AssertionError struct {
	Message string
}

func (a AssertionError) Error() string {
	return fmt.Sprintf("AssertionError: %s", a.Message)
}
