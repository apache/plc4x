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

type PlcUDINT struct {
	value uint32
	PlcSimpleNumericValueAdapter
}

func NewPlcUDINT(value uint32) PlcUDINT {
	return PlcUDINT{
		value: value,
	}
}

func (m PlcUDINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcUDINT) IsUint8() bool {
	return m.value <= math.MaxUint8
}

func (m PlcUDINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetUint16())
	}
	return 0
}

func (m PlcUDINT) IsUint16() bool {
	return m.value <= math.MaxUint16
}

func (m PlcUDINT) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetUint32())
	}
	return 0
}

func (m PlcUDINT) GetUint32() uint32 {
	return m.value
}

func (m PlcUDINT) GetUint64() uint64 {
	return uint64(m.GetUint32())
}

func (m PlcUDINT) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcUDINT) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint32())
	}
	return 0
}

func (m PlcUDINT) IsInt16() bool {
	return m.value < math.MaxInt16
}

func (m PlcUDINT) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetUint32())
	}
	return 0
}

func (m PlcUDINT) IsInt32() bool {
	return m.value < math.MaxInt32
}

func (m PlcUDINT) GetInt32() int32 {
	if m.IsInt32() {
		return int32(m.GetUint32())
	}
	return 0
}

func (m PlcUDINT) GetInt64() int64 {
	return int64(m.GetUint32())
}

func (m PlcUDINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint32())
}

func (m PlcUDINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint32())
}

func (m PlcUDINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcUDINT) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint32("PlcUDINT", 32, m.value)
}
