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
	"math"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcBYTE struct {
	PlcSimpleValueAdapter
	value uint8
}

func NewPlcBYTE(value uint8) PlcBYTE {
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

func (m PlcBYTE) GetUint8() uint8 {
	return m.value
}

func (m PlcBYTE) GetUint16() uint16 {
	return uint16(m.GetUint8())
}

func (m PlcBYTE) GetUint32() uint32 {
	return uint32(m.GetUint8())
}

func (m PlcBYTE) GetUint64() uint64 {
	return uint64(m.GetUint8())
}

func (m PlcBYTE) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcBYTE) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint8())
	}
	return 0
}

func (m PlcBYTE) GetInt16() int16 {
	return int16(m.GetUint8())
}

func (m PlcBYTE) GetInt32() int32 {
	return int32(m.GetUint8())
}

func (m PlcBYTE) GetInt64() int64 {
	return int64(m.GetUint8())
}

func (m PlcBYTE) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint8())
}

func (m PlcBYTE) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint8())
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
