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

type PlcBYTE struct {
	PlcSimpleValueAdapter
	value byte
}

func NewPlcBYTE(value byte) PlcBYTE {
	return PlcBYTE{
		value: value,
	}
}

func (m PlcBYTE) GetRaw() []byte {
	return []byte{m.value}
}

func (m PlcBYTE) IsBool() bool {
	return true
}

func (m PlcBYTE) GetBoolLength() uint32 {
	return 8
}

func (m PlcBYTE) GetBool() bool {
	return m.value&1 == 1
}

func (m PlcBYTE) GetBoolAt(index uint32) bool {
	if index > 7 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcBYTE) GetBoolArray() []bool {
	return []bool{m.value&1 == 1, m.value>>1&1 == 1,
		m.value>>2&1 == 1, m.value>>3&1 == 1,
		m.value>>4&1 == 1, m.value>>5&1 == 1,
		m.value>>6&1 == 1, m.value>>7&1 == 1}
}

func (m PlcBYTE) IsByte() bool {
	return true
}

func (m PlcBYTE) GetByte() byte {
	return m.value
}

func (m PlcBYTE) IsUint8() bool {
	return true
}
func (m PlcBYTE) GetUint8() uint8 {
	return m.value
}

func (m PlcBYTE) IsString() bool {
	return true
}

func (m PlcBYTE) GetString() string {
	var strVal string
	for i, val := range m.GetBoolArray() {
		if i > 0 {
			strVal = strVal + ", "
		}
		if val {
			strVal = strVal + "true"
		} else {
			strVal = strVal + "false"
		}
	}
	return strVal
}

func (m PlcBYTE) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.BYTE
}

func (m PlcBYTE) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcBYTE) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteByte("PlcBYTE", m.value)
}

func (m PlcBYTE) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 8, m.value)
}
