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
	"strconv"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcSINT struct {
	value int8
	PlcSimpleNumericValueAdapter
}

func NewPlcSINT(value int8) PlcSINT {
	return PlcSINT{
		value: value,
	}
}

func (m PlcSINT) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcSINT) GetBoolean() bool {
	if m.value == 0 {
		return false
	}
	return true
}

func (m PlcSINT) IsByte() bool {
	return m.IsUint8()
}

func (m PlcSINT) GetByte() byte {
	return m.GetUint8()
}

func (m PlcSINT) IsUint8() bool {
	return m.value >= 0
}

func (m PlcSINT) GetUint8() uint8 {
	if m.IsUint8() {
		return uint8(m.GetInt8())
	}
	return 0
}

func (m PlcSINT) IsUint16() bool {
	return m.value >= 0
}

func (m PlcSINT) GetUint16() uint16 {
	if m.IsUint16() {
		return uint16(m.GetInt8())
	}
	return 0
}

func (m PlcSINT) IsUint32() bool {
	return m.value >= 0
}

func (m PlcSINT) GetUint32() uint32 {
	if m.IsUint32() {
		return uint32(m.GetInt8())
	}
	return 0
}

func (m PlcSINT) IsUint64() bool {
	return m.value >= 0
}

func (m PlcSINT) GetUint64() uint64 {
	if m.IsUint64() {
		return uint64(m.GetInt8())
	}
	return 0
}

func (m PlcSINT) GetInt8() int8 {
	return m.value
}

func (m PlcSINT) GetInt16() int16 {
	return int16(m.GetInt8())
}

func (m PlcSINT) GetInt32() int32 {
	return int32(m.GetInt8())
}

func (m PlcSINT) GetInt64() int64 {
	return int64(m.GetInt8())
}

func (m PlcSINT) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetInt8())
}

func (m PlcSINT) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetInt8())
}

func (m PlcSINT) GetString() string {
	return strconv.Itoa(int(m.GetInt64()))
}

func (m PlcSINT) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.SINT
}

func (m PlcSINT) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcSINT) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteInt8("PlcINT", 8, m.value)
}

func (m PlcSINT) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 8, m.value)
}
