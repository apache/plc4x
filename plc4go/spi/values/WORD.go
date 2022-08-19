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
)

type PlcWORD struct {
	PlcSimpleValueAdapter
	value uint16
}

func NewPlcWORD(value uint16) PlcWORD {
	return PlcWORD{
		value: value,
	}
}

func (m PlcWORD) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	m.Serialize(buf)
	return buf.GetBytes()
}

func (m PlcWORD) IsBool() bool {
	return true
}

func (m PlcWORD) GetBoolLength() uint32 {
	return 16
}

func (m PlcWORD) GetBool() bool {
	return m.value&1 == 1
}

func (m PlcWORD) GetBoolAt(index uint32) bool {
	if index > 15 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcWORD) GetBoolArray() []bool {
	return []bool{m.value&1 == 1, m.value>>1&1 == 1,
		m.value>>2&1 == 1, m.value>>3&1 == 1,
		m.value>>4&1 == 1, m.value>>5&1 == 1,
		m.value>>6&1 == 1, m.value>>7&1 == 1,
		m.value>>8&1 == 1, m.value>>9&1 == 1,
		m.value>>10&1 == 1, m.value>>11&1 == 1,
		m.value>>12&1 == 1, m.value>>13&1 == 1,
		m.value>>14&1 == 1, m.value>>15&1 == 1}
}

func (m PlcWORD) IsString() bool {
	return true
}

func (m PlcWORD) GetString() string {
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

func (m PlcWORD) GetPLCValueType() apiValues.PLCValueType {
	return apiValues.WORD
}

func (m PlcWORD) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint16("PlcWORD", 16, m.value)
}

func (m PlcWORD) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPLCValueType(), 16, m.value)
}
