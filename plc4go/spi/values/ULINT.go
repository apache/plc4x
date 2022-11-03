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
	"math"
	"strconv"
)

type PlcULINT struct {
	value uint64
	PlcSimpleNumericValueAdapter
}

func NewPlcULINT(value uint64) PlcULINT {
	return PlcULINT{
		value: value,
	}
}

func (m PlcULINT) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcULINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcULINT) IsByte() bool {
	return m.IsUint8()
}

func (m PlcULINT) GetByte() byte {
	return m.GetUint8()
}

func (m PlcULINT) IsUint8() bool {
	return m.value <= math.MaxUint8
}

func (m PlcULINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetUint16())
	}
	return 0
}

func (m PlcULINT) IsUint16() bool {
	return m.value <= math.MaxUint16
}

func (m PlcULINT) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetUint32())
	}
	return 0
}

func (m PlcULINT) IsUint32() bool {
	return m.value <= math.MaxUint32
}

func (m PlcULINT) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetUint64())
	}
	return 0
}

func (m PlcULINT) GetUint64() uint64 {
	return m.value
}

func (m PlcULINT) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcULINT) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint64())
	}
	return 0
}

func (m PlcULINT) IsInt16() bool {
	return m.value < math.MaxInt16
}

func (m PlcULINT) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetUint64())
	}
	return 0
}

func (m PlcULINT) IsInt32() bool {
	return m.value < math.MaxInt32
}

func (m PlcULINT) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetUint64())
	}
	return 0
}

func (m PlcULINT) IsInt64() bool {
	return m.value < math.MaxInt64
}

func (m PlcULINT) GetInt64() int64 {
	if m.IsInt64() {
		return int64(m.GetUint64())
	}
	return 0
}

func (m PlcULINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint32())
}

func (m PlcULINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint32())
}

func (m PlcULINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcULINT) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.ULINT
}

func (m PlcULINT) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcULINT) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint64("PlcUINT", 64, m.value)
}

func (m PlcULINT) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 64, m.value)
}
