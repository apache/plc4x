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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcDWORD struct {
	PlcSimpleValueAdapter
	value uint32
}

func NewPlcDWORD(value uint32) PlcDWORD {
	return PlcDWORD{
		value: value,
	}
}

func (m PlcDWORD) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcDWORD) IsBool() bool {
	return true
}

func (m PlcDWORD) GetBoolLength() uint32 {
	return 32
}

func (m PlcDWORD) GetBool() bool {
	return m.value&1 == 1
}

func (m PlcDWORD) GetBoolAt(index uint32) bool {
	if index > 31 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcDWORD) GetBoolArray() []bool {
	return []bool{m.value&1 == 1, m.value>>1&1 == 1,
		m.value>>2&1 == 1, m.value>>3&1 == 1,
		m.value>>4&1 == 1, m.value>>5&1 == 1,
		m.value>>6&1 == 1, m.value>>7&1 == 1,
		m.value>>8&1 == 1, m.value>>9&1 == 1,
		m.value>>10&1 == 1, m.value>>11&1 == 1,
		m.value>>12&1 == 1, m.value>>13&1 == 1,
		m.value>>14&1 == 1, m.value>>15&1 == 1,
		m.value>>16&1 == 1, m.value>>17&1 == 1,
		m.value>>18&1 == 1, m.value>>19&1 == 1,
		m.value>>20&1 == 1, m.value>>21&1 == 1,
		m.value>>22&1 == 1, m.value>>23&1 == 1,
		m.value>>24&1 == 1, m.value>>25&1 == 1,
		m.value>>26&1 == 1, m.value>>27&1 == 1,
		m.value>>28&1 == 1, m.value>>29&1 == 1,
		m.value>>30&1 == 1, m.value>>31&1 == 1}
}

func (m PlcDWORD) IsString() bool {
	return true
}

func (m PlcDWORD) GetString() string {
	var strVal string
	for i, val := range m.GetBoolArray() {
		if i > 0 {
			strVal = strVal + ", "
		}
		if val {
			strVal = strVal + "true"
		} else {
			strVal = strVal + "false"
		}
	}
	return strVal
}

func (m PlcDWORD) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DWORD
}

func (m PlcDWORD) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcDWORD) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint32("PlcDINT", 32, m.value)
}

func (m PlcDWORD) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 32, m.value)
}
