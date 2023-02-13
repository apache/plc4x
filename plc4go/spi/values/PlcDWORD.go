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
	"context"
	"encoding/binary"
	"fmt"
	"math"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcDWORD struct {
	PlcSimpleValueAdapter
	value uint32
}

func NewPlcDWORD(value uint32) PlcDWORD {
	return PlcDWORD{
		value: value,
	}
}

func (m PlcDWORD) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcDWORD) IsBool() bool {
	return true
}

func (m PlcDWORD) GetBoolLength() uint32 {
	return 32
}

func (m PlcDWORD) GetBool() bool {
	return m.value&1 == 1
}

func (m PlcDWORD) GetBoolAt(index uint32) bool {
	if index > 31 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcDWORD) GetBoolArray() []bool {
	return []bool{m.value&1 == 1, m.value>>1&1 == 1,
		m.value>>2&1 == 1, m.value>>3&1 == 1,
		m.value>>4&1 == 1, m.value>>5&1 == 1,
		m.value>>6&1 == 1, m.value>>7&1 == 1,
		m.value>>8&1 == 1, m.value>>9&1 == 1,
		m.value>>10&1 == 1, m.value>>11&1 == 1,
		m.value>>12&1 == 1, m.value>>13&1 == 1,
		m.value>>14&1 == 1, m.value>>15&1 == 1,
		m.value>>16&1 == 1, m.value>>17&1 == 1,
		m.value>>18&1 == 1, m.value>>19&1 == 1,
		m.value>>20&1 == 1, m.value>>21&1 == 1,
		m.value>>22&1 == 1, m.value>>23&1 == 1,
		m.value>>24&1 == 1, m.value>>25&1 == 1,
		m.value>>26&1 == 1, m.value>>27&1 == 1,
		m.value>>28&1 == 1, m.value>>29&1 == 1,
		m.value>>30&1 == 1, m.value>>31&1 == 1}
}

func (m PlcDWORD) IsUint8() bool {
	return m.value <= math.MaxUint8
}

func (m PlcDWORD) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetUint16())
	}
	return 0
}

func (m PlcDWORD) IsUint16() bool {
	return m.value <= math.MaxUint16
}

func (m PlcDWORD) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetUint32())
	}
	return 0
}

func (m PlcDWORD) GetUint32() uint32 {
	return m.value
}

func (m PlcDWORD) GetUint64() uint64 {
	return uint64(m.GetUint32())
}

func (m PlcDWORD) IsByte() bool {
	return m.IsUint8()
}

func (m PlcDWORD) GetByte() byte {
	return m.GetUint8()
}

func (m PlcDWORD) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcDWORD) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint32())
	}
	return 0
}

func (m PlcDWORD) IsInt16() bool {
	return m.value < math.MaxInt16
}

func (m PlcDWORD) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetUint32())
	}
	return 0
}

func (m PlcDWORD) IsInt32() bool {
	return m.value < math.MaxInt32
}

func (m PlcDWORD) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetUint32())
	}
	return 0
}

func (m PlcDWORD) GetInt64() int64 {
	return int64(m.GetUint32())
}

func (m PlcDWORD) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint32())
}

func (m PlcDWORD) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint32())
}

func (m PlcDWORD) IsString() bool {
	return true
}

func (m PlcDWORD) GetString() string {
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

func (m PlcDWORD) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DWORD
}

func (m PlcDWORD) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcDWORD) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint32("PlcDINT", 32, m.value)
}

func (m PlcDWORD) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 32, m.value)
}
