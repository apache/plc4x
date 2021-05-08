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

import (
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"time"
)

// Dummy structure
type PlcValueAdapter struct {
}

// Simple Types
func (m PlcValueAdapter) IsSimple() bool {
	return false
}
func (m PlcValueAdapter) IsNullable() bool {
	return false
}
func (m PlcValueAdapter) IsNull() bool {
	return false
}

// Boolean
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
	return m.IsUint8()
}

func (m PlcValueAdapter) GetByte() byte {
	return m.GetUint8()
}

// Integer
func (m PlcValueAdapter) IsUint8() bool {
	return false
}
func (m PlcValueAdapter) GetUint8() uint8 {
	return 0
}
func (m PlcValueAdapter) IsUint16() bool {
	return false
}
func (m PlcValueAdapter) GetUint16() uint16 {
	return 0
}
func (m PlcValueAdapter) IsUint32() bool {
	return false
}
func (m PlcValueAdapter) GetUint32() uint32 {
	return 0
}
func (m PlcValueAdapter) IsUint64() bool {
	return false
}
func (m PlcValueAdapter) GetUint64() uint64 {
	return 0
}
func (m PlcValueAdapter) IsInt8() bool {
	return false
}
func (m PlcValueAdapter) GetInt8() int8 {
	return 0
}
func (m PlcValueAdapter) IsInt16() bool {
	return false
}
func (m PlcValueAdapter) GetInt16() int16 {
	return 0
}
func (m PlcValueAdapter) IsInt32() bool {
	return false
}
func (m PlcValueAdapter) GetInt32() int32 {
	return 0
}
func (m PlcValueAdapter) IsInt64() bool {
	return false
}
func (m PlcValueAdapter) GetInt64() int64 {
	return 0
}

// Floating Point
func (m PlcValueAdapter) IsFloat32() bool {
	return false
}
func (m PlcValueAdapter) GetFloat32() float32 {
	return 0.0
}
func (m PlcValueAdapter) IsFloat64() bool {
	return false
}
func (m PlcValueAdapter) GetFloat64() float64 {
	return 0.0
}

// String
func (m PlcValueAdapter) IsString() bool {
	return false
}
func (m PlcValueAdapter) GetString() string {
	return ""
}

// Time
func (m PlcValueAdapter) IsTime() bool {
	return false
}
func (m PlcValueAdapter) GetTime() time.Time {
	return time.Time{}
}
func (m PlcValueAdapter) IsDuration() bool {
	return false
}
func (m PlcValueAdapter) GetDuration() time.Duration {
	return 0
}

// Raw Access
func (m PlcValueAdapter) GetRaw() []byte {
	return nil
}

// List Methods
func (m PlcValueAdapter) IsList() bool {
	return true
}
func (m PlcValueAdapter) GetLength() uint32 {
	return 1
}
func (m PlcValueAdapter) GetIndex(i uint32) api.PlcValue {
	return nil
}
func (m PlcValueAdapter) GetList() []api.PlcValue {
	return []api.PlcValue{}
}

// Struct Methods
func (m PlcValueAdapter) IsStruct() bool {
	return false
}
func (m PlcValueAdapter) GetKeys() []string {
	return []string{}
}
func (m PlcValueAdapter) HasKey(key string) bool {
	return false
}
func (m PlcValueAdapter) GetValue(key string) api.PlcValue {
	return nil
}
func (m PlcValueAdapter) GetStruct() map[string]api.PlcValue {
	return nil
}
func (m PlcValueAdapter) IsDate() bool {
	return false
}
func (m PlcValueAdapter) GetDate() time.Time {
	return time.Time{}
}
func (m PlcValueAdapter) IsDateTime() bool {
	return false
}
func (m PlcValueAdapter) GetDateTime() time.Time {
	return time.Time{}
}
