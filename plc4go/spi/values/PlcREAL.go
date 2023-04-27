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

type PlcREAL struct {
	value float32
	PlcSimpleNumericValueAdapter
}

func NewPlcREAL(value float32) PlcREAL {
	return PlcREAL{
		value: value,
	}
}

func (m PlcREAL) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcREAL) GetBool() bool {
	if m.value == 0.0 {
		return false
	}
	return true
}

func (m PlcREAL) IsByte() bool {
	return m.IsUint8()
}

func (m PlcREAL) GetByte() byte {
	return m.GetUint8()
}

func (m PlcREAL) IsUint8() bool {
	return m.value >= 0 && m.value <= math.MaxUint8
}

func (m PlcREAL) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsUint16() bool {
	return m.value >= 0 && m.value <= math.MaxUint16
}

func (m PlcREAL) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsUint32() bool {
	return m.value >= 0 && m.value <= math.MaxUint32
}

func (m PlcREAL) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsUint64() bool {
	return m.value >= 0 && m.value <= math.MaxUint64
}

func (m PlcREAL) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsInt8() bool {
	return m.value >= math.MinInt8 && m.value <= math.MaxInt8
}

func (m PlcREAL) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsInt16() bool {
	return m.value >= math.MinInt16 && m.value <= math.MaxInt16
}

func (m PlcREAL) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsInt32() bool {
	return m.value >= math.MinInt32 && m.value <= math.MaxInt32
}

func (m PlcREAL) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) IsInt64() bool {
	return m.value >= math.MinInt64 && m.value <= math.MaxInt64
}

func (m PlcREAL) GetInt64() int64 {
	if m.IsInt64() {
		return int64(m.GetFloat32())
	}
	return 0
}

func (m PlcREAL) GetFloat32() float32 {
	return m.value
}

func (m PlcREAL) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetFloat32())
}

func (m PlcREAL) IsString() bool {
	return true
}

func (m PlcREAL) GetString() string {
	return fmt.Sprintf("%g", m.GetFloat32())
}

func (m PlcREAL) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.REAL
}

func (m PlcREAL) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcREAL) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteFloat32("PlcREAL", 32, m.value)
}

func (m PlcREAL) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 32, m.value)
}
