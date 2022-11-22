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

package values

import (
	"encoding/binary"
	"fmt"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcBOOL struct {
	PlcSimpleValueAdapter
	value bool
}

func NewPlcBOOL(value bool) PlcBOOL {
	return PlcBOOL{
		value: value,
	}
}

func (m PlcBOOL) GetRaw() []byte {
	if m.value {
		return []byte{0x01}
	}
	return []byte{0x00}
}

func (m PlcBOOL) IsBool() bool {
	return true
}

func (m PlcBOOL) GetBoolLength() uint32 {
	return 1
}

func (m PlcBOOL) GetBool() bool {
	return m.value
}

func (m PlcBOOL) GetBoolAt(index uint32) bool {
	if index == 0 {
		return m.value
	}
	return false
}

func (m PlcBOOL) GetBoolArray() []bool {
	return []bool{m.value}
}

func (m PlcBOOL) IsString() bool {
	return true
}

func (m PlcBOOL) GetString() string {
	if m.GetBool() {
		return "true"
	} else {
		return "false"
	}
}

func (m PlcBOOL) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.BOOL
}

func (m PlcBOOL) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcBOOL) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteBit("PlcBOOL", m.value)
}

func (m PlcBOOL) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 1, m.value)
}
