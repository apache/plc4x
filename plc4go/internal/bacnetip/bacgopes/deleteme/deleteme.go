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

package deleteme

import (
	"context"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type MessageBridge interface {
	spi.Message
	PDUData
	IsMessageBridge()
}

type messageBridge struct {
	PDUData
}

func NewMessageBridge(bytes ...byte) MessageBridge {
	return &messageBridge{NewPDUData(NewArgs(bytes), NoKWArgs)}
}

var _ MessageBridge = (*messageBridge)(nil)

func (m *messageBridge) Serialize() ([]byte, error) {
	return m.GetPduData(), nil
}

func (m *messageBridge) SerializeWithWriteBuffer(_ context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteByteArray("Bytes", m.GetPduData())
}

func (m *messageBridge) GetLengthInBytes(_ context.Context) uint16 {
	return uint16(len(m.GetPduData()))
}

func (m *messageBridge) GetLengthInBits(ctx context.Context) uint16 {
	return m.GetLengthInBytes(ctx) * 8
}

func (*messageBridge) IsMessageBridge() {
}
