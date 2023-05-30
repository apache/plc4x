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
	"time"
)

type PlcValue interface {
	fmt.Stringer

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

	IsRaw() bool
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
	NULL           PlcValueType = 0x00
	BOOL           PlcValueType = 0x01
	BYTE           PlcValueType = 0x02
	WORD           PlcValueType = 0x03
	DWORD          PlcValueType = 0x04
	LWORD          PlcValueType = 0x05
	USINT          PlcValueType = 0x11
	UINT           PlcValueType = 0x12
	UDINT          PlcValueType = 0x13
	ULINT          PlcValueType = 0x14
	SINT           PlcValueType = 0x21
	INT            PlcValueType = 0x22
	DINT           PlcValueType = 0x23
	LINT           PlcValueType = 0x24
	REAL           PlcValueType = 0x31
	LREAL          PlcValueType = 0x32
	CHAR           PlcValueType = 0x41
	WCHAR          PlcValueType = 0x42
	STRING         PlcValueType = 0x43
	WSTRING        PlcValueType = 0x44
	TIME           PlcValueType = 0x51
	LTIME          PlcValueType = 0x52
	DATE           PlcValueType = 0x53
	LDATE          PlcValueType = 0x54
	TIME_OF_DAY    PlcValueType = 0x55
	LTIME_OF_DAY   PlcValueType = 0x56
	DATE_AND_TIME  PlcValueType = 0x57
	LDATE_AND_TIME PlcValueType = 0x58
	Struct         PlcValueType = 0x61
	List           PlcValueType = 0x62
	RAW_BYTE_ARRAY PlcValueType = 0x71
)

func (p PlcValueType) String() string {
	switch {
	case p == NULL:
		return "NULL"
	case p == BOOL:
		return "BOOL"
	case p == BYTE:
		return "BYTE"
	case p == WORD:
		return "WORD"
	case p == DWORD:
		return "DWORD"
	case p == LWORD:
		return "LWORD"
	case p == USINT:
		return "USINT"
	case p == UINT:
		return "UINT"
	case p == UDINT:
		return "UDINT"
	case p == ULINT:
		return "ULINT"
	case p == SINT:
		return "SINT"
	case p == INT:
		return "INT"
	case p == DINT:
		return "DINT"
	case p == LINT:
		return "LINT"
	case p == REAL:
		return "REAL"
	case p == LREAL:
		return "LREAL"
	case p == CHAR:
		return "CHAR"
	case p == WCHAR:
		return "WCHAR"
	case p == STRING:
		return "STRING"
	case p == WSTRING:
		return "WSTRING"
	case p == TIME:
		return "TIME"
	case p == LTIME:
		return "LTIME"
	case p == DATE:
		return "DATE"
	case p == LDATE:
		return "LDATE"
	case p == TIME_OF_DAY:
		return "TIME_OF_DAY"
	case p == LTIME_OF_DAY:
		return "LTIME_OF_DAY"
	case p == DATE_AND_TIME:
		return "DATE_AND_TIME"
	case p == LDATE_AND_TIME:
		return "LDATE_AND_TIME"
	case p == Struct:
		return "Struct"
	case p == List:
		return "List"
	case p == RAW_BYTE_ARRAY:
		return "RAW_BYTE_ARRAY"
	}
	return "Unknown"
}

func PlcValueByName(value string) (valueType PlcValueType, ok bool) {
	switch value {
	case "NULL":
		return NULL, true
	case "BOOL":
		return BOOL, true
	case "BYTE":
		return BYTE, true
	case "WORD":
		return WORD, true
	case "DWORD":
		return DWORD, true
	case "LWORD":
		return LWORD, true
	case "USINT":
		return USINT, true
	case "UINT":
		return UINT, true
	case "UDINT":
		return UDINT, true
	case "ULINT":
		return ULINT, true
	case "SINT":
		return SINT, true
	case "INT":
		return INT, true
	case "DINT":
		return DINT, true
	case "LINT":
		return LINT, true
	case "REAL":
		return REAL, true
	case "LREAL":
		return LREAL, true
	case "CHAR":
		return CHAR, true
	case "WCHAR":
		return WCHAR, true
	case "STRING":
		return STRING, true
	case "WSTRING":
		return WSTRING, true
	case "TIME":
		return TIME, true
	case "LTIME":
		return LTIME, true
	case "DATE":
		return DATE, true
	case "LDATE":
		return LDATE, true
	case "TIME_OF_DAY":
		return TIME_OF_DAY, true
	case "LTIME_OF_DAY":
		return LTIME_OF_DAY, true
	case "DATE_AND_TIME":
		return DATE_AND_TIME, true
	case "LDATE_AND_TIME":
		return LDATE_AND_TIME, true
	case "Struct":
		return Struct, true
	case "List":
		return List, true
	case "RAW_BYTE_ARRAY":
		return RAW_BYTE_ARRAY, true
	}
	return NULL, false
}
