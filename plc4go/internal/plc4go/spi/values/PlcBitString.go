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
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

type PlcBitString struct {
	PlcList
}

func NewPlcBitString(value interface{}) PlcBitString {
	var bools []api.PlcValue
	switch value.(type) {
	case uint8:
		bools = make([]api.PlcValue, 8)
		for i := 0; i < 8; i++ {
			bools[i] = NewPlcBOOL(((value.(uint8) >> uint8((8-1)-i)) & 0x01) == 0x01)
		}
	case uint16:
		bools = make([]api.PlcValue, 16)
		for i := 0; i < 16; i++ {
			bools[i] = NewPlcBOOL(((value.(uint16) >> uint8((16-1)-i)) & 0x01) == 0x01)
		}
	case uint32:
		bools = make([]api.PlcValue, 32)
		for i := 0; i < 32; i++ {
			bools[i] = NewPlcBOOL(((value.(uint32) >> uint8((32-1)-i)) & 0x01) == 0x01)
		}
	case uint64:
		bools = make([]api.PlcValue, 64)
		for i := 0; i < 64; i++ {
			bools[i] = NewPlcBOOL(((value.(uint64) >> uint8((64-1)-i)) & 0x01) == 0x01)
		}
	}

	return PlcBitString{
		PlcList{Values: bools},
	}
}

func (m PlcBitString) IsList() bool {
	return true
}

func (m PlcBitString) GetLength() uint32 {
	return uint32(len(m.Values))
}

func (m PlcBitString) GetIndex(i uint32) api.PlcValue {
	return m.Values[i]
}

func (m PlcBitString) GetList() []api.PlcValue {
	return m.Values
}
