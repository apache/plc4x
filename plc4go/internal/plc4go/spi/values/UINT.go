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

type PlcUINT struct {
	value uint16
	PlcSimpleNumericValueAdapter
}

func NewPlcUINT(value uint16) PlcUINT {
	return PlcUINT{
		value: value,
	}
}

func (m PlcUINT) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	m.Serialize(buf)
	return buf.GetBytes()
}

func (m PlcUINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcUINT) IsUint8() bool {
	return m.value <= math.MaxUint8
}

func (m PlcUINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetUint16())
	}
	return 0
}

func (m PlcUINT) GetUint16() uint16 {
	return m.value
}

func (m PlcUINT) GetUint32() uint32 {
	return uint32(m.GetUint16())
}

func (m PlcUINT) GetUint64() uint64 {
	return uint64(m.GetUint16())
}

func (m PlcUINT) IsInt8() bool {
	return m.value < math.MaxInt8
}

func (m PlcUINT) GetInt8() int8 {
	if m.IsInt8() {
		return int8(m.GetUint16())
	}
	return 0
}

func (m PlcUINT) IsInt16() bool {
	return m.value < math.MaxInt16
}

func (m PlcUINT) GetInt16() int16 {
	if m.IsInt16() {
		return int16(m.GetUint16())
	}
	return 0
}

func (m PlcUINT) GetInt32() int32 {
	return int32(m.GetUint16())
}

func (m PlcUINT) GetInt64() int64 {
	return int64(m.GetUint16())
}

func (m PlcUINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint16())
}

func (m PlcUINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint16())
}

func (m PlcUINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcUINT) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint16("PlcUINT", 16, m.value)
}
