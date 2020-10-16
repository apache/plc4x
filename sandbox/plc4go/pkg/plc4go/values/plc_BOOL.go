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
package values

type PlcBOOL struct {
	value bool
	plcSimpleNumericValueAdapter
}

func NewPlcBOOL(value bool) PlcBOOL {
	return PlcBOOL{
		value: value,
	}
}

func (m PlcBOOL) GetBoolean() bool {
	return m.value
}

func (m PlcBOOL) GetUint8() uint8 {
	if m.value == true {
		return 1
	}
	return 0
}

func (m PlcBOOL) GetUint16() uint16 {
	return uint16(m.GetUint8())
}

func (m PlcBOOL) GetUint32() uint32 {
	return uint32(m.GetUint8())
}

func (m PlcBOOL) GetUint64() uint64 {
	return uint64(m.GetUint8())
}

func (m PlcBOOL) GetInt8() int8 {
	return int8(m.GetUint8())
}

func (m PlcBOOL) GetInt16() int16 {
	return int16(m.GetUint8())
}

func (m PlcBOOL) GetInt32() int32 {
	return int32(m.GetUint8())
}

func (m PlcBOOL) GetInt64() int64 {
	return int64(m.GetUint8())
}

func (m PlcBOOL) GetFloat32() float32 {
	//TODO: Check if this is ok
	return float32(m.GetUint8())
}

func (m PlcBOOL) GetFloat64() float64 {
	//TODO: Check if this is ok
	return float64(m.GetUint8())
}
