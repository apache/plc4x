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
	"fmt"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"math"
	"strconv"
)

type PlcLINT struct {
	value int64
	PlcSimpleNumericValueAdapter
}

func NewPlcLINT(value int64) PlcLINT {
	return PlcLINT{
		value: value,
	}
}

func (m PlcLINT) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	_ = m.Serialize(buf)
	return buf.GetBytes()
}

func (m PlcLINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcLINT) IsByte() bool {
	return m.IsUint8()
}

func (m PlcLINT) GetByte() byte {
	return m.GetUint8()
}

func (m PlcLINT) IsUint8() bool {
	return m.value >= 0 && m.value <= math.MaxUint8
}

func (m PlcLINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) IsUint16() bool {
	return m.value >= 0 && m.value <= math.MaxUint16
}

func (m PlcLINT) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) IsUint32() bool {
	return m.value >= 0 && m.value <= math.MaxUint32
}

func (m PlcLINT) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) IsUint64() bool {
	return m.value >= 0
}

func (m PlcLINT) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) IsInt8() bool {
	return m.value >= math.MinInt8 && m.value <= math.MaxInt8
}

func (m PlcLINT) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) IsInt16() bool {
	return m.value >= math.MinInt16 && m.value <= math.MaxInt16
}

func (m PlcLINT) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) IsInt32() bool {
	return m.value >= math.MinInt32 && m.value <= math.MaxInt32
}

func (m PlcLINT) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetInt64())
	}
	return 0
}

func (m PlcLINT) GetInt64() int64 {
	return m.value
}

func (m PlcLINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetInt64())
}

func (m PlcLINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetInt64())
}

func (m PlcLINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcLINT) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.LINT
}

func (m PlcLINT) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteInt64("PlcLINT", 64, m.value)
}

func (m PlcLINT) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 64, m.value)
}
