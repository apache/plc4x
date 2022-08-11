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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"math"
	"strconv"
)

type PlcDINT struct {
	value int32
	PlcSimpleNumericValueAdapter
}

func NewPlcDINT(value int32) PlcDINT {
	return PlcDINT{
		value: value,
	}
}

func (m PlcDINT) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	m.Serialize(buf)
	return buf.GetBytes()
}

func (m PlcDINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcDINT) IsUint8() bool {
	return m.value >= 0 && m.value <= math.MaxUint8
}

func (m PlcDINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetInt32())
	}
	return 0
}

func (m PlcDINT) IsUint16() bool {
	return m.value >= 0 && m.value <= math.MaxUint16
}

func (m PlcDINT) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetInt32())
	}
	return 0
}

func (m PlcDINT) IsUint32() bool {
	return m.value >= 0
}

func (m PlcDINT) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetInt32())
	}
	return 0
}

func (m PlcDINT) IsUint64() bool {
	return m.value >= 0
}

func (m PlcDINT) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.GetInt32())
	}
	return 0
}

func (m PlcDINT) IsInt8() bool {
	return m.value >= math.MinInt8 && m.value <= math.MaxInt8
}

func (m PlcDINT) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetInt32())
	}
	return 0
}

func (m PlcDINT) IsInt16() bool {
	return m.value >= math.MinInt16 && m.value <= math.MaxInt16
}

func (m PlcDINT) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetInt32())
	}
	return 0
}

func (m PlcDINT) GetInt32() int32 {
	return m.value
}

func (m PlcDINT) GetInt64() int64 {
	return int64(m.GetInt32())
}

func (m PlcDINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetInt32())
}

func (m PlcDINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetInt32())
}

func (m PlcDINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcDINT) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteInt32("PlcDINT", 32, m.value)
}
