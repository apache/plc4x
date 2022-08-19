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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"math"
	"math/big"
	"strconv"
)

type PlcBINT struct {
	value *big.Int
	PlcSimpleNumericValueAdapter
}

func NewPlcBINT(value *big.Int) PlcBINT {
	return PlcBINT{
		value: value,
	}
}

func (m PlcBINT) GetRaw() []byte {
	return m.value.Bytes()
}

func (m PlcBINT) GetBoolean() bool {
	if m.isZero() {
		return false
	}
	return true
}

func (m PlcBINT) IsUint8() bool {
	return m.isEqualsOrGreaterZero() && m.value.Cmp(big.NewInt(math.MaxUint8)) <= 0
}

func (m PlcBINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetInt16())
	}
	return 0
}

func (m PlcBINT) IsUint16() bool {
	return m.isEqualsOrGreaterZero()
}

func (m PlcBINT) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetInt16())
	}
	return 0
}

func (m PlcBINT) IsUint32() bool {
	return m.isEqualsOrGreaterZero()
}

func (m PlcBINT) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetInt32())
	}
	return 0
}

func (m PlcBINT) IsUint64() bool {
	return m.isEqualsOrGreaterZero()
}

func (m PlcBINT) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.value.Int64())
	}
	return 0
}

func (m PlcBINT) IsInt8() bool {
	return m.isGreaterOrEqual(math.MinInt8) && m.isLowerOrEqual(math.MaxInt8)
}

func (m PlcBINT) GetInt8() int8 {
	return int8(m.value.Int64())
}

func (m PlcBINT) GetInt16() int16 {
	return int16(m.value.Int64())
}

func (m PlcBINT) GetInt32() int32 {
	return int32(m.value.Int64())
}

func (m PlcBINT) GetInt64() int64 {
	return m.value.Int64()
}

func (m PlcBINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.value.Int64())
}

func (m PlcBINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.value.Int64())
}

func (m PlcBINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcBINT) GetPLCValueType() apiValues.PLCValueType {
	return apiValues.BINT
}

func (m PlcBINT) isZero() bool {
	return m.value.Cmp(big.NewInt(0)) == 0
}

func (m PlcBINT) isEqualsOrGreaterZero() bool {
	return m.isGreaterOrEqual(0)
}

func (m PlcBINT) isGreaterOrEqual(other int64) bool {
	return m.value.Cmp(big.NewInt(other)) >= 0
}

func (m PlcBINT) isLowerOrEqual(other int64) bool {
	return m.value.Cmp(big.NewInt(other)) <= 0
}

func (m PlcBINT) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteBigInt("PlcBINT", uint8(m.value.BitLen()), m.value)
}
