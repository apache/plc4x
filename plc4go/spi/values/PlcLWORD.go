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

type PlcLWORD struct {
	PlcSimpleValueAdapter
	value uint64
}

func NewPlcLWORD(value uint64) PlcLWORD {
	return PlcLWORD{
		value: value,
	}
}

func (m PlcLWORD) IsRaw() bool {
	return true
}

func (m PlcLWORD) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcLWORD) IsBool() bool {
	return true
}

func (m PlcLWORD) GetBoolLength() uint32 {
	return 64
}

func (m PlcLWORD) GetBool() bool {
	return m.value&1 == 1
}

func (m PlcLWORD) GetBoolAt(index uint32) bool {
	if index > 63 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcLWORD) GetBoolArray() []bool {
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
		m.value>>30&1 == 1, m.value>>31&1 == 1,
		m.value>>32&1 == 1, m.value>>33&1 == 1,
		m.value>>34&1 == 1, m.value>>35&1 == 1,
		m.value>>36&1 == 1, m.value>>37&1 == 1,
		m.value>>38&1 == 1, m.value>>39&1 == 1,
		m.value>>40&1 == 1, m.value>>41&1 == 1,
		m.value>>42&1 == 1, m.value>>43&1 == 1,
		m.value>>44&1 == 1, m.value>>45&1 == 1,
		m.value>>46&1 == 1, m.value>>47&1 == 1,
		m.value>>48&1 == 1, m.value>>49&1 == 1,
		m.value>>50&1 == 1, m.value>>51&1 == 1,
		m.value>>52&1 == 1, m.value>>53&1 == 1,
		m.value>>54&1 == 1, m.value>>55&1 == 1,
		m.value>>56&1 == 1, m.value>>57&1 == 1,
		m.value>>58&1 == 1, m.value>>59&1 == 1,
		m.value>>60&1 == 1, m.value>>61&1 == 1,
		m.value>>62&1 == 1, m.value>>63&1 == 1}
}

func (m PlcLWORD) IsByte() bool {
	return m.IsUint8()
}

func (m PlcLWORD) GetByte() byte {
	return m.GetUint8()
}

func (m PlcLWORD) IsUint8() bool {
	return m.value <= math.MaxUint8
}

func (m PlcLWORD) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetUint16())
	}
	return 0
}

func (m PlcLWORD) IsUint16() bool {
	return m.value <= math.MaxUint16
}

func (m PlcLWORD) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetUint32())
	}
	return 0
}

func (m PlcLWORD) IsUint32() bool {
	return m.value <= math.MaxUint32
}

func (m PlcLWORD) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetUint64())
	}
	return 0
}

func (m PlcLWORD) GetUint64() uint64 {
	return m.value
}

func (m PlcLWORD) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcLWORD) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint64())
	}
	return 0
}

func (m PlcLWORD) IsInt16() bool {
	return m.value < math.MaxInt16
}

func (m PlcLWORD) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetUint64())
	}
	return 0
}

func (m PlcLWORD) IsInt32() bool {
	return m.value < math.MaxInt32
}

func (m PlcLWORD) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetUint64())
	}
	return 0
}

func (m PlcLWORD) IsInt64() bool {
	return m.value < math.MaxInt64
}

func (m PlcLWORD) GetInt64() int64 {
	if m.IsInt64() {
		return int64(m.GetUint64())
	}
	return 0
}

func (m PlcLWORD) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint32())
}

func (m PlcLWORD) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint32())
}

func (m PlcLWORD) IsString() bool {
	return true
}

func (m PlcLWORD) GetString() string {
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

func (m PlcLWORD) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.LWORD
}

func (m PlcLWORD) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcLWORD) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint64("PlcLWORD", 64, m.value)
}

func (m PlcLWORD) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 64, m.value)
}
