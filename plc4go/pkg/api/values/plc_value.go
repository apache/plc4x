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

import "time"

type PlcValue interface {
	////
	// Simple Types

	IsSimple() bool
	IsNullable() bool
	IsNull() bool
	//
	///

	////
	// Boolean

	IsBool() bool
	GetBoolLength() uint32
	GetBool() bool
	GetBoolAt(index uint32) bool
	GetBoolArray() []bool
	//
	///

	////
	// Byte

	IsByte() bool
	GetByte() byte
	//
	///

	////
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
	//
	///

	////
	// Floating Point

	IsFloat32() bool
	GetFloat32() float32
	IsFloat64() bool
	GetFloat64() float64
	//
	///

	////
	// String

	IsString() bool
	GetString() string
	//
	///

	////
	// Time

	IsTime() bool
	GetTime() time.Time
	IsDuration() bool
	GetDuration() time.Duration
	IsDate() bool
	GetDate() time.Time
	IsDateTime() bool
	GetDateTime() time.Time
	//
	///

	////
	// Raw Access

	GetRaw() []byte
	//
	///

	////
	// List Methods

	IsList() bool
	GetLength() uint32
	GetIndex(i uint32) PlcValue
	GetList() []PlcValue
	//
	///

	////
	// Struct Methods

	IsStruct() bool
	GetKeys() []string
	HasKey(key string) bool
	GetValue(key string) PlcValue
	GetStruct() map[string]PlcValue
	//
	///

	GetPlcValueType() PlcValueType
}

// RawPlcValue This type is used in cases where the driver doesn't have access to type information and therefore can't decode
// the payload yet. This allows an application to take the raw plc-value and have the payload decoded later.
type RawPlcValue interface {
	// PlcValue the base value
	PlcValue

	// RawDecodeValue Read the internal buffer and parse a value of given type
	RawDecodeValue(typeName string) PlcValue
	// RawHasMore If the internal read-buffer has not yet reached the end
	RawHasMore() bool
	// RawReset Reset the internal read-buffer (For the case that a raw plc-value has to be parsed multiple times)
	RawReset()
}

type PlcValueType uint8

const (
	BINT PlcValueType = iota
	BIT_STRING
	BOOL
	BREAL
	BYTE
	BYTE_ARRAY
	CHAR
	DATE
	DATE_AND_TIME
	DINT
	DWORD
	INT
	LINT
	LIST
	LREAL
	LTIME
	LWORD
	NULL
	RAW_PLC_VALUE
	REAL
	STRUCT
	SINT
	STRING
	TIME
	TIME_OF_DAY
	UDINT
	UINT
	ULINT
	USINT
	WCHAR
	WORD
	WSTRING
)

func (p PlcValueType) String() string {
	switch {
	case p == BINT:
		return "BINT"
	case p == BIT_STRING:
		return "BIT_STRING"
	case p == BOOL:
		return "BOOL"
	case p == BREAL:
		return "BREAL"
	case p == BYTE:
		return "BYTE"
	case p == BYTE_ARRAY:
		return "BYTE_ARRAY"
	case p == CHAR:
		return "CHAR"
	case p == DATE:
		return "DATE"
	case p == DATE_AND_TIME:
		return "DATE_AND_TIME"
	case p == DINT:
		return "DINT"
	case p == DWORD:
		return "DWORD"
	case p == INT:
		return "INT"
	case p == LINT:
		return "LINT"
	case p == LIST:
		return "LIST"
	case p == LREAL:
		return "LREAL"
	case p == LTIME:
		return "LTIME"
	case p == LWORD:
		return "LWORD"
	case p == NULL:
		return "NULL"
	case p == RAW_PLC_VALUE:
		return "RAW_PLC_VALUE"
	case p == REAL:
		return "REAL"
	case p == STRUCT:
		return "STRUCT"
	case p == SINT:
		return "SINT"
	case p == STRING:
		return "STRING"
	case p == TIME:
		return "TIME"
	case p == TIME_OF_DAY:
		return "TIME_OF_DAY"
	case p == UDINT:
		return "UDINT"
	case p == UINT:
		return "UINT"
	case p == ULINT:
		return "ULINT"
	case p == USINT:
		return "USINT"
	case p == WCHAR:
		return "WCHAR"
	case p == WORD:
		return "WORD"
	case p == WSTRING:
		return "WSTRING"
	}
	return "Unknown"
}
