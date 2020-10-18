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
	GetBooleanLength() uint32
	GetBoolean() bool
	GetBooleanAt(index uint32) bool
	GetBooleanArray() []bool

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
	GetTime() time.Time
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
