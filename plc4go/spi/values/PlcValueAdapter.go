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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"time"
)

// PlcValueAdapter Dummy structure
type PlcValueAdapter struct {
	utils.Serializable
}

////////
////
// Simple Types
//

func (m PlcValueAdapter) IsSimple() bool {
	return false
}
func (m PlcValueAdapter) IsNullable() bool {
	return false
}
func (m PlcValueAdapter) IsNull() bool {
	return false
}

////////
////
// Boolean
//

func (m PlcValueAdapter) IsBool() bool {
	return false
}
func (m PlcValueAdapter) GetBoolLength() uint32 {
	return 1
}
func (m PlcValueAdapter) GetBool() bool {
	return false
}
func (m PlcValueAdapter) GetBoolAt(index uint32) bool {
	if index == 0 {
		return m.GetBool()
	}
	return false
}
func (m PlcValueAdapter) GetBoolArray() []bool {
	return nil
}

func (m PlcValueAdapter) IsByte() bool {
	return false
}

func (m PlcValueAdapter) GetByte() byte {
	panic("GetByte not implemented")
}

////////
////
// Integer
//

func (m PlcValueAdapter) IsUint8() bool {
	return false
}
func (m PlcValueAdapter) GetUint8() uint8 {
	panic("GetUint8 not implemented")
}
func (m PlcValueAdapter) IsUint16() bool {
	return false
}
func (m PlcValueAdapter) GetUint16() uint16 {
	panic("GetUint16 not implemented")
}
func (m PlcValueAdapter) IsUint32() bool {
	return false
}
func (m PlcValueAdapter) GetUint32() uint32 {
	panic("GetUint32 not implemented")
}
func (m PlcValueAdapter) IsUint64() bool {
	return false
}
func (m PlcValueAdapter) GetUint64() uint64 {
	panic("GetUint64 not implemented")
}
func (m PlcValueAdapter) IsInt8() bool {
	return false
}
func (m PlcValueAdapter) GetInt8() int8 {
	panic("GetInt8 not implemented")
}
func (m PlcValueAdapter) IsInt16() bool {
	return false
}
func (m PlcValueAdapter) GetInt16() int16 {
	panic("GetInt16 not implemented")
}
func (m PlcValueAdapter) IsInt32() bool {
	return false
}
func (m PlcValueAdapter) GetInt32() int32 {
	panic("GetInt32 not implemented")
}
func (m PlcValueAdapter) IsInt64() bool {
	return false
}
func (m PlcValueAdapter) GetInt64() int64 {
	panic("GetInt64 not implemented")
}

////////
////
// Floating Point
//

func (m PlcValueAdapter) IsFloat32() bool {
	return false
}
func (m PlcValueAdapter) GetFloat32() float32 {
	panic("GetFloat32 not implemented")
}
func (m PlcValueAdapter) IsFloat64() bool {
	return false
}
func (m PlcValueAdapter) GetFloat64() float64 {
	panic("GetFloat64 not implemented")
}

////////
////
// String
//

func (m PlcValueAdapter) IsString() bool {
	return false
}
func (m PlcValueAdapter) GetString() string {
	panic("GetString not implemented")
}

////////
////
// Time
//

func (m PlcValueAdapter) IsTime() bool {
	return false
}
func (m PlcValueAdapter) GetTime() time.Time {
	panic("GetTime not implemented")
}
func (m PlcValueAdapter) IsDuration() bool {
	return false
}
func (m PlcValueAdapter) GetDuration() time.Duration {
	panic("GetDuration not implemented")
}

////////
////
// Raw access
//

func (m PlcValueAdapter) IsRaw() bool {
	return false
}

func (m PlcValueAdapter) GetRaw() []byte {
	panic("GetRaw not implemented")
}

////////
////
// List Methods
//

func (m PlcValueAdapter) IsList() bool {
	return false
}
func (m PlcValueAdapter) GetLength() uint32 {
	panic("GetLength not implemented")
}
func (m PlcValueAdapter) GetIndex(i uint32) apiValues.PlcValue {
	return nil
}
func (m PlcValueAdapter) GetList() []apiValues.PlcValue {
	panic("GetList not implemented")
}

////////
////
// Struct Methods
//

func (m PlcValueAdapter) IsStruct() bool {
	return false
}
func (m PlcValueAdapter) GetKeys() []string {
	panic("GetKeys not implemented")
}
func (m PlcValueAdapter) HasKey(_ string) bool {
	return false
}
func (m PlcValueAdapter) GetValue(_ string) apiValues.PlcValue {
	panic("GetValue not implemented")
}
func (m PlcValueAdapter) GetStruct() map[string]apiValues.PlcValue {
	panic("GetStruct not implemented")
}
func (m PlcValueAdapter) IsDate() bool {
	return false
}
func (m PlcValueAdapter) GetDate() time.Time {
	panic("GetDate not implemented")
}
func (m PlcValueAdapter) IsDateTime() bool {
	return false
}
func (m PlcValueAdapter) GetDateTime() time.Time {
	panic("GetDateTime not implemented")
}

func (m PlcValueAdapter) GetPlcValueType() apiValues.PlcValueType {
	panic("GetPlcValueType not implemented")
}

func (m PlcValueAdapter) String() string {
	return "not implemented"
}

type PlcSimpleValueAdapter struct {
	PlcValueAdapter
}

func (m PlcSimpleValueAdapter) IsSimple() bool {
	return true
}

func (m PlcSimpleValueAdapter) GetLength() uint32 {
	return 1
}

func (m PlcSimpleValueAdapter) String() string {
	return "not implemented"
}

type PlcSimpleNumericValueAdapter struct {
	PlcSimpleValueAdapter
}

func (m PlcSimpleNumericValueAdapter) IsBool() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsByte() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint8() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint16() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint32() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsUint64() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsInt8() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsInt16() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsInt32() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsInt64() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsFloat32() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsFloat64() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) IsString() bool {
	return true
}

func (m PlcSimpleNumericValueAdapter) String() string {
	return "not implemented"
}
