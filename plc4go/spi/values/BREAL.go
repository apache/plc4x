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
	"math/big"
)

type PlcBREAL struct {
	value *big.Float
	PlcSimpleNumericValueAdapter
}

func NewPlcBREAL(value *big.Float) PlcBREAL {
	return PlcBREAL{
		value: value,
	}
}

func (m PlcBREAL) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	m.Serialize(buf)
	return buf.GetBytes()
}

func (m PlcBREAL) GetBoolean() bool {
	if m.isZero() {
		return false
	}
	return true
}

func (m PlcBREAL) IsUint8() bool {
	return m.isGreaterOrEqual(0.0) && m.isLowerOrEqual(math.MaxUint8)
}

func (m PlcBREAL) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsUint16() bool {
	return m.isGreaterOrEqual(0.0) && m.isLowerOrEqual(math.MaxUint16)
}

func (m PlcBREAL) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsUint32() bool {
	return m.isGreaterOrEqual(0.0) && m.isLowerOrEqual(math.MaxUint32)
}

func (m PlcBREAL) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsUint64() bool {
	return m.isGreaterOrEqual(0.0) && m.isLowerOrEqual(math.MaxUint64)
}

func (m PlcBREAL) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsByte() bool {
	return m.IsUint8()
}

func (m PlcBREAL) GetByte() byte {
	return m.GetUint8()
}

func (m PlcBREAL) IsInt8() bool {
	return m.isGreaterOrEqual(math.MinInt8) && m.isLowerOrEqual(math.MaxInt8)
}

func (m PlcBREAL) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsInt16() bool {
	return m.isGreaterOrEqual(math.MinInt16) && m.isLowerOrEqual(math.MaxInt16)
}

func (m PlcBREAL) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsInt32() bool {
	return m.isGreaterOrEqual(math.MinInt32) && m.isLowerOrEqual(math.MaxInt32)
}

func (m PlcBREAL) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) IsInt64() bool {
	return m.isGreaterOrEqual(math.MinInt64) && m.isLowerOrEqual(math.MaxInt64)
}

func (m PlcBREAL) GetInt64() int64 {
	if m.IsInt64() {
		return int64(m.GetFloat32())
	}
	return 0
}

func (m PlcBREAL) GetFloat32() float32 {
	f, _ := m.value.Float32()
	return f
}

func (m PlcBREAL) GetFloat64() float64 {
	f, _ := m.value.Float64()
	return f
}

func (m PlcBREAL) IsString() bool {
	return true
}

func (m PlcBREAL) GetString() string {
	return fmt.Sprintf("%g", m.GetFloat64())
}

func (m PlcBREAL) GetPLCValueType() apiValues.PLCValueType {
	return apiValues.BREAL
}

func (m PlcBREAL) isZero() bool {
	return m.value.Cmp(big.NewFloat(0.0)) == 0.0
}

func (m PlcBREAL) isEqualsOrGreaterZero() bool {
	return m.isGreaterOrEqual(0)
}

func (m PlcBREAL) isGreaterOrEqual(other float64) bool {
	return m.value.Cmp(big.NewFloat(other)) >= 0
}

func (m PlcBREAL) isLowerOrEqual(other float64) bool {
	return m.value.Cmp(big.NewFloat(other)) <= 0
}

func (m PlcBREAL) Serialize(writeBuffer utils.WriteBuffer) error {
	// TODO: fix this a insert a valid bit length calculation
	return writeBuffer.WriteBigFloat("PlcBREAL", uint8(m.value.MinPrec()), m.value)
}

func (m PlcBREAL) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPLCValueType(), m.value.MinPrec(), m.value)
}
