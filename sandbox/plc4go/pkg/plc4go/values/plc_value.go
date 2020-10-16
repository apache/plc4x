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

import "time"

type PlcValue interface {

	// Simple Types
	IsSimple() bool
	IsNullable() bool
	IsNull() bool

	// Boolean
	IsBoolean() bool
	GetBoolean() bool

	// Integer
	IsUint8() bool
	GetUint8() uint8
	IsUint16() bool
	GetUint16() uint16
	IsUint32() bool
	GetUint32() uint32
	IsUint64() bool
	GetUint64() uint64
	IsInt8() bool
	GetInt8() int8
	IsInt16() bool
	GetInt16() int16
	IsInt32() bool
	GetInt32() int32
	IsInt64() bool
	GetInt64() int64

	// Floating Point
	IsFloat32() bool
	GetFloat32() float32
	IsFloat64() bool
	GetFloat64() float64

	// String
	IsString() bool
	GetString() string

	// Time
	IsTime() bool
	GetTime()
	IsDate() bool
	GetDate()
	IsDateTime() bool
	GetDateTime()

	// Raw Access
	GetRaw() []byte

	// List Methods
	IsList() bool
	GetLength() uint32
	GetIndex(i uint32) PlcValue
	GetList() []PlcValue

	// Struct Methods
	IsStruct() bool
	GetKeys() []string
	HasKey(key string) bool
	GetValue(key string) PlcValue
	GetStruct() map[string]PlcValue
}

// Dummy structure
type plcValueAdapter struct {
}

// Simple Types
func (m plcValueAdapter) IsSimple() bool {
	return false
}
func (m plcValueAdapter) IsNullable() bool {
	return false
}
func (m plcValueAdapter) IsNull() bool {
	return false
}

// Boolean
func (m plcValueAdapter) IsBoolean() bool {
	return false
}
func (m plcValueAdapter) GetBoolean() bool {
	return false
}

// Integer
func (m plcValueAdapter) IsUint8() bool {
	return false
}
func (m plcValueAdapter) GetUint8() uint8 {
	return 0
}
func (m plcValueAdapter) IsUint16() bool {
	return false
}
func (m plcValueAdapter) GetUint16() uint16 {
	return 0
}
func (m plcValueAdapter) IsUint32() bool {
	return false
}
func (m plcValueAdapter) GetUint32() uint32 {
	return 0
}
func (m plcValueAdapter) IsUint64() bool {
	return false
}
func (m plcValueAdapter) GetUint64() uint64 {
	return 0
}
func (m plcValueAdapter) IsInt8() bool {
	return false
}
func (m plcValueAdapter) GetInt8() uint8 {
	return 0
}
func (m plcValueAdapter) IsInt16() bool {
	return false
}
func (m plcValueAdapter) GetInt16() uint16 {
	return 0
}
func (m plcValueAdapter) IsInt32() bool {
	return false
}
func (m plcValueAdapter) GetInt32() uint32 {
	return 0
}
func (m plcValueAdapter) IsInt64() bool {
	return false
}
func (m plcValueAdapter) GetInt64() uint64 {
	return 0
}

// Floating Point
func (m plcValueAdapter) IsFloat32() bool {
	return false
}
func (m plcValueAdapter) GetFloat32() float32 {
	return 0.0
}
func (m plcValueAdapter) IsFloat64() bool {
	return false
}
func (m plcValueAdapter) GetFloat64() float64 {
	return 0.0
}

// String
func (m plcValueAdapter) IsString() bool {
	return false
}
func (m plcValueAdapter) GetString() string {
	return ""
}

// Time
func (m plcValueAdapter) IsTime() bool {
	return false
}
func (m plcValueAdapter) GetTime() time.Time {
	return time.Time{}
}
func (m plcValueAdapter) IsDuration() bool {
	return false
}
func (m plcValueAdapter) GetDuration() time.Duration {
	return 0
}

// Raw Access
func (m plcValueAdapter) GetRaw() []byte {
	return nil
}

// List Methods
func (m plcValueAdapter) IsList() bool {
	return false
}
func (m plcValueAdapter) GetLength() int {
	return 0
}
func (m plcValueAdapter) GetIndex(i int) PlcValue {
	return nil
}
func (m plcValueAdapter) GetList() []PlcValue {
	return nil
}

// Struct Methods
func (m plcValueAdapter) IsStruct() bool {
	return false
}
func (m plcValueAdapter) GetKeys() []string {
	return []string{}
}
func (m plcValueAdapter) HasKey(key string) bool {
	return false
}
func (m plcValueAdapter) GetValue(key string) PlcValue {
	return nil
}
func (m plcValueAdapter) GetStruct() map[string]PlcValue {
	return nil
}
