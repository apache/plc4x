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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"math"
	"strconv"
)

type PlcUSINT struct {
	value uint8
	PlcSimpleNumericValueAdapter
}

func NewPlcUSINT(value uint8) PlcUSINT {
	return PlcUSINT{
		value: value,
	}
}

func (m PlcUSINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcUSINT) GetUint8() uint8 {
	return m.value
}

func (m PlcUSINT) GetUint16() uint16 {
	return uint16(m.GetUint8())
}

func (m PlcUSINT) GetUint32() uint32 {
	return uint32(m.GetUint8())
}

func (m PlcUSINT) GetUint64() uint64 {
	return uint64(m.GetUint8())
}

func (m PlcUSINT) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcUSINT) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint8())
	}
	return 0
}

func (m PlcUSINT) GetInt16() int16 {
	return int16(m.GetUint8())
}

func (m PlcUSINT) GetInt32() int32 {
	return int32(m.GetUint8())
}

func (m PlcUSINT) GetInt64() int64 {
	return int64(m.GetUint8())
}

func (m PlcUSINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint8())
}

func (m PlcUSINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint8())
}

func (m PlcUSINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcUSINT) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint8("PlcUSINT", 64, m.value)
}
