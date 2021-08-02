/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"math"
)

type PlcLREAL struct {
	value float64
	PlcSimpleNumericValueAdapter
}

func NewPlcLREAL(value float64) PlcLREAL {
	return PlcLREAL{
		value: value,
	}
}

func (m PlcLREAL) GetBoolean() bool {
	if m.value == 0.0 {
		return false
	}
	return true
}

func (m PlcLREAL) IsUint8() bool {
	return m.value >= 0 && m.value <= math.MaxUint8
}

func (m PlcLREAL) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsUint16() bool {
	return m.value >= 0 && m.value <= math.MaxUint16
}

func (m PlcLREAL) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsUint32() bool {
	return m.value >= 0 && m.value <= math.MaxUint32
}

func (m PlcLREAL) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsUint64() bool {
	return m.value >= 0 && m.value <= math.MaxUint64
}

func (m PlcLREAL) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsInt8() bool {
	return m.value >= math.MinInt8 && m.value <= math.MaxInt8
}

func (m PlcLREAL) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsInt16() bool {
	return m.value >= math.MinInt16 && m.value <= math.MaxInt16
}

func (m PlcLREAL) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsInt32() bool {
	return m.value >= math.MinInt32 && m.value <= math.MaxInt32
}

func (m PlcLREAL) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsInt64() bool {
	return m.value >= math.MinInt64 && m.value <= math.MaxInt64
}

func (m PlcLREAL) GetInt64() int64 {
	if m.IsInt64() {
		return int64(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) IsFloat32() bool {
	return m.value >= -math.MaxFloat32 && m.value <= math.MaxFloat32
}

func (m PlcLREAL) GetFloat32() float32 {
	if m.IsInt64() {
		return float32(m.GetFloat64())
	}
	return 0
}

func (m PlcLREAL) GetFloat64() float64 {
	return m.value
}

func (m PlcLREAL) IsString() bool {
	return true
}

func (m PlcLREAL) GetString() string {
	return fmt.Sprintf("%g", m.GetFloat64())
}

func (m PlcLREAL) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteFloat64("PlcLREAL", 64, m.value)
}
