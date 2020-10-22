//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package iec61131

import "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/model/values"

type PlcWORD struct {
	value uint16
	values.PlcSimpleValueAdapter
}

func NewPlcWORD(value uint16) PlcWORD {
	return PlcWORD{
		value: value,
	}
}

func (m PlcWORD) IsBoolean() bool {
	return true
}

func (m PlcWORD) GetBooleanLength() uint32 {
	return 16
}

func (m PlcWORD) GetBoolean() bool {
	return m.value&1 == 1
}

func (m PlcWORD) GetBooleanAt(index uint32) bool {
	if index > 15 {
		return false
	}
	return m.value>>index&1 == 1
}

func (m PlcWORD) GetBooleanArray() []bool {
	return []bool{m.value&1 == 1, m.value>>1&1 == 1,
		m.value>>2&1 == 1, m.value>>3&1 == 1,
		m.value>>4&1 == 1, m.value>>5&1 == 1,
		m.value>>6&1 == 1, m.value>>7&1 == 1,
		m.value>>8&1 == 1, m.value>>9&1 == 1,
		m.value>>10&1 == 1, m.value>>11&1 == 1,
		m.value>>12&1 == 1, m.value>>13&1 == 1,
		m.value>>14&1 == 1, m.value>>15&1 == 1}
}
