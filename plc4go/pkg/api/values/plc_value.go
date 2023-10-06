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

	// IsSimple tells if this value is a simple datatype
	IsSimple() bool
	// IsNullable tells if this value is nullable
	IsNullable() bool
	// IsNull tells if this value is null
	IsNull() bool
	//
	///

	////
	// Boolean

	// IsBool tells if this value is a bool
	IsBool() bool
	// GetBoolLength return the bool length. Attention: Before using check with IsBool otherwise it might panic.
	GetBoolLength() uint32
	// GetBool return the bool. Attention: Before using check with IsBool otherwise it might panic.
	GetBool() bool
	// GetBoolAt return the bool at specified index. Attention: Before using check with IsBool otherwise it might panic.
	GetBoolAt(index uint32) bool
	// GetBoolArray return an array of bool. Attention: Before using check with IsBool otherwise it might panic.
	GetBoolArray() []bool
	//
	///

	////
	// Byte

	// IsByte tells if this is a byte
	IsByte() bool
	// GetByte return the byte length. Attention: Before using check with IsByte otherwise it might panic.
	GetByte() byte
	//
	///

	////
	// Integer

	// IsUint8 tells if this is an uint8
	IsUint8() bool
	// GetUint8 return the uint8. Attention: Before using check with IsUint8 otherwise it might panic.
	GetUint8() uint8
	// IsUint16 tells if this is an uint16
	IsUint16() bool
	// GetUint16 return the uint16. Attention: Before using check with IsUint16 otherwise it might panic.
	GetUint16() uint16
	// IsUint32 tells if this is an uint32
	IsUint32() bool
	// GetUint32 return the uint32. Attention: Before using check with IsUint32 otherwise it might panic.
	GetUint32() uint32
	// IsUint64 tells if this is an uint64
	IsUint64() bool
	// GetUint64 return the uint64. Attention: Before using check with IsUint64 otherwise it might panic.
	GetUint64() uint64
	// IsInt8 tells if this is an int8
	IsInt8() bool
	// GetInt8 return the int8. Attention: Before using check with IsInt8 otherwise it might panic.
	GetInt8() int8
	// IsInt16 tells if this is an int16
	IsInt16() bool
	// GetInt16 return the int16. Attention: Before using check with IsInt16 otherwise it might panic.
	GetInt16() int16
	// IsInt32 tells if this is an int32
	IsInt32() bool
	// GetInt32 return the int32. Attention: Before using check with IsInt32 otherwise it might panic.
	GetInt32() int32
	// IsInt64 tells if this is an int64
	IsInt64() bool
	// GetInt64 return the int64. Attention: Before using check with IsInt64 otherwise it might panic.
	GetInt64() int64
	//
	///

	////
	// Floating Point

	// IsFloat32 tells if this is a float32
	IsFloat32() bool
	// GetFloat32 return the float32. Attention: Before using check with IsFloat32 otherwise it might panic.
	GetFloat32() float32
	// IsFloat64 tells if this is a float64
	IsFloat64() bool
	// GetFloat64 return the float64. Attention: Before using check with IsFloat64 otherwise it might panic.
	GetFloat64() float64
	//
	///

	////
	// String

	// IsString tells if this is a string
	IsString() bool
	// GetString return the string. Attention: Before using check with IsString otherwise it might panic.
	GetString() string
	//
	///

	////
	// Time

	// IsTime tells if this is a time.Time
	IsTime() bool
	// GetTime return the time.Time. Attention: Before using check with IsTime otherwise it might panic.
	GetTime() time.Time
	// IsDuration tells if this is a time.Duration
	IsDuration() bool
	// GetDuration return the time.Duration. Attention: Before using check with IsDuration otherwise it might panic.
	GetDuration() time.Duration
	// IsDate tells if this is a time.Time
	IsDate() bool
	// GetDate return the time.Time. Attention: Before using check with IsDate otherwise it might panic.
	GetDate() time.Time
	// IsDateTime tells if this is a time.Time
	IsDateTime() bool
	// GetDateTime return the time.Time. Attention: Before using check with IsDateTime otherwise it might panic.
	GetDateTime() time.Time
	//
	///

	////
	// Raw Access

	// IsRaw tells if this is a raw value
	IsRaw() bool
	// GetRaw return the []byte. Attention: Before using check with IsRaw otherwise it might panic.
	GetRaw() []byte
	//
	///

	////
	// List Methods

	// IsList tells if this is a list
	IsList() bool
	// GetLength return the length of list. Attention: Before using check with IsList otherwise it might panic.
	GetLength() uint32
	// GetIndex return the element at index or nil if not found. Attention: Before using check with IsList otherwise it might panic.
	GetIndex(i uint32) PlcValue
	// GetList return the list. Attention: Before using check with IsList otherwise it might panic.
	GetList() []PlcValue
	//
	///

	////
	// Struct Methods

	// IsStruct tells if this is a struct (map)
	IsStruct() bool
	// GetKeys return the keys of the struct. Attention: Before using check with IsStruct otherwise it might panic.
	GetKeys() []string
	// HasKey returns true if it has the key. Attention: Before using check with IsStruct otherwise it might panic.
	HasKey(key string) bool
	// GetValue return the value of the struct or nil if not found. Attention: Before using check with IsStruct otherwise it might panic.
	GetValue(key string) PlcValue
	// GetStruct return the struct map. Attention: Before using check with IsStruct otherwise it might panic.
	GetStruct() map[string]PlcValue
	//
	///

	// GetPlcValueType returns the PlcValueType
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
