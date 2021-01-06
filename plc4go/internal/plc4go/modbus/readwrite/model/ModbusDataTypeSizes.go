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
package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

type ModbusDataTypeSizes string

type IModbusDataTypeSizes interface {
	DataTypeSize() uint8
	Serialize(io utils.WriteBuffer) error
}

const (
	ModbusDataTypeSizes_BOOL           ModbusDataTypeSizes = "IEC61131_BOOL"
	ModbusDataTypeSizes_BYTE           ModbusDataTypeSizes = "IEC61131_BYTE"
	ModbusDataTypeSizes_WORD           ModbusDataTypeSizes = "IEC61131_WORD"
	ModbusDataTypeSizes_DWORD          ModbusDataTypeSizes = "IEC61131_DWORD"
	ModbusDataTypeSizes_LWORD          ModbusDataTypeSizes = "IEC61131_LWORD"
	ModbusDataTypeSizes_SINT           ModbusDataTypeSizes = "IEC61131_SINT"
	ModbusDataTypeSizes_INT            ModbusDataTypeSizes = "IEC61131_INT"
	ModbusDataTypeSizes_DINT           ModbusDataTypeSizes = "IEC61131_DINT"
	ModbusDataTypeSizes_LINT           ModbusDataTypeSizes = "IEC61131_LINT"
	ModbusDataTypeSizes_USINT          ModbusDataTypeSizes = "IEC61131_USINT"
	ModbusDataTypeSizes_UINT           ModbusDataTypeSizes = "IEC61131_UINT"
	ModbusDataTypeSizes_UDINT          ModbusDataTypeSizes = "IEC61131_UDINT"
	ModbusDataTypeSizes_ULINT          ModbusDataTypeSizes = "IEC61131_ULINT"
	ModbusDataTypeSizes_REAL           ModbusDataTypeSizes = "IEC61131_REAL"
	ModbusDataTypeSizes_LREAL          ModbusDataTypeSizes = "IEC61131_LREAL"
	ModbusDataTypeSizes_TIME           ModbusDataTypeSizes = "IEC61131_TIME"
	ModbusDataTypeSizes_LTIME          ModbusDataTypeSizes = "IEC61131_LTIME"
	ModbusDataTypeSizes_DATE           ModbusDataTypeSizes = "IEC61131_DATE"
	ModbusDataTypeSizes_LDATE          ModbusDataTypeSizes = "IEC61131_LDATE"
	ModbusDataTypeSizes_TIME_OF_DAY    ModbusDataTypeSizes = "IEC61131_TIME_OF_DAY"
	ModbusDataTypeSizes_LTIME_OF_DAY   ModbusDataTypeSizes = "IEC61131_LTIME_OF_DAY"
	ModbusDataTypeSizes_DATE_AND_TIME  ModbusDataTypeSizes = "IEC61131_DATE_AND_TIME"
	ModbusDataTypeSizes_LDATE_AND_TIME ModbusDataTypeSizes = "IEC61131_LDATE_AND_TIME"
	ModbusDataTypeSizes_CHAR           ModbusDataTypeSizes = "IEC61131_CHAR"
	ModbusDataTypeSizes_WCHAR          ModbusDataTypeSizes = "IEC61131_WCHAR"
	ModbusDataTypeSizes_STRING         ModbusDataTypeSizes = "IEC61131_STRING"
	ModbusDataTypeSizes_WSTRING        ModbusDataTypeSizes = "IEC61131_WSTRING"
)

func (e ModbusDataTypeSizes) DataTypeSize() uint8 {
	switch e {
	case "IEC61131_BOOL":
		{ /* 'IEC61131_BOOL' */
			return 1
		}
	case "IEC61131_BYTE":
		{ /* 'IEC61131_BYTE' */
			return 1
		}
	case "IEC61131_CHAR":
		{ /* 'IEC61131_CHAR' */
			return 1
		}
	case "IEC61131_DATE":
		{ /* 'IEC61131_DATE' */
			return 8
		}
	case "IEC61131_DATE_AND_TIME":
		{ /* 'IEC61131_DATE_AND_TIME' */
			return 8
		}
	case "IEC61131_DINT":
		{ /* 'IEC61131_DINT' */
			return 4
		}
	case "IEC61131_DWORD":
		{ /* 'IEC61131_DWORD' */
			return 4
		}
	case "IEC61131_INT":
		{ /* 'IEC61131_INT' */
			return 2
		}
	case "IEC61131_LDATE":
		{ /* 'IEC61131_LDATE' */
			return 8
		}
	case "IEC61131_LDATE_AND_TIME":
		{ /* 'IEC61131_LDATE_AND_TIME' */
			return 8
		}
	case "IEC61131_LINT":
		{ /* 'IEC61131_LINT' */
			return 8
		}
	case "IEC61131_LREAL":
		{ /* 'IEC61131_LREAL' */
			return 8
		}
	case "IEC61131_LTIME":
		{ /* 'IEC61131_LTIME' */
			return 8
		}
	case "IEC61131_LTIME_OF_DAY":
		{ /* 'IEC61131_LTIME_OF_DAY' */
			return 8
		}
	case "IEC61131_LWORD":
		{ /* 'IEC61131_LWORD' */
			return 8
		}
	case "IEC61131_REAL":
		{ /* 'IEC61131_REAL' */
			return 4
		}
	case "IEC61131_SINT":
		{ /* 'IEC61131_SINT' */
			return 1
		}
	case "IEC61131_STRING":
		{ /* 'IEC61131_STRING' */
			return 1
		}
	case "IEC61131_TIME":
		{ /* 'IEC61131_TIME' */
			return 8
		}
	case "IEC61131_TIME_OF_DAY":
		{ /* 'IEC61131_TIME_OF_DAY' */
			return 8
		}
	case "IEC61131_UDINT":
		{ /* 'IEC61131_UDINT' */
			return 4
		}
	case "IEC61131_UINT":
		{ /* 'IEC61131_UINT' */
			return 2
		}
	case "IEC61131_ULINT":
		{ /* 'IEC61131_ULINT' */
			return 8
		}
	case "IEC61131_USINT":
		{ /* 'IEC61131_USINT' */
			return 1
		}
	case "IEC61131_WCHAR":
		{ /* 'IEC61131_WCHAR' */
			return 2
		}
	case "IEC61131_WORD":
		{ /* 'IEC61131_WORD' */
			return 2
		}
	case "IEC61131_WSTRING":
		{ /* 'IEC61131_WSTRING' */
			return 2
		}
	default:
		{
			return 0
		}
	}
}
func ModbusDataTypeSizesByValue(value string) ModbusDataTypeSizes {
	switch value {
	case "IEC61131_BOOL":
		return ModbusDataTypeSizes_BOOL
	case "IEC61131_BYTE":
		return ModbusDataTypeSizes_BYTE
	case "IEC61131_CHAR":
		return ModbusDataTypeSizes_CHAR
	case "IEC61131_DATE":
		return ModbusDataTypeSizes_DATE
	case "IEC61131_DATE_AND_TIME":
		return ModbusDataTypeSizes_DATE_AND_TIME
	case "IEC61131_DINT":
		return ModbusDataTypeSizes_DINT
	case "IEC61131_DWORD":
		return ModbusDataTypeSizes_DWORD
	case "IEC61131_INT":
		return ModbusDataTypeSizes_INT
	case "IEC61131_LDATE":
		return ModbusDataTypeSizes_LDATE
	case "IEC61131_LDATE_AND_TIME":
		return ModbusDataTypeSizes_LDATE_AND_TIME
	case "IEC61131_LINT":
		return ModbusDataTypeSizes_LINT
	case "IEC61131_LREAL":
		return ModbusDataTypeSizes_LREAL
	case "IEC61131_LTIME":
		return ModbusDataTypeSizes_LTIME
	case "IEC61131_LTIME_OF_DAY":
		return ModbusDataTypeSizes_LTIME_OF_DAY
	case "IEC61131_LWORD":
		return ModbusDataTypeSizes_LWORD
	case "IEC61131_REAL":
		return ModbusDataTypeSizes_REAL
	case "IEC61131_SINT":
		return ModbusDataTypeSizes_SINT
	case "IEC61131_STRING":
		return ModbusDataTypeSizes_STRING
	case "IEC61131_TIME":
		return ModbusDataTypeSizes_TIME
	case "IEC61131_TIME_OF_DAY":
		return ModbusDataTypeSizes_TIME_OF_DAY
	case "IEC61131_UDINT":
		return ModbusDataTypeSizes_UDINT
	case "IEC61131_UINT":
		return ModbusDataTypeSizes_UINT
	case "IEC61131_ULINT":
		return ModbusDataTypeSizes_ULINT
	case "IEC61131_USINT":
		return ModbusDataTypeSizes_USINT
	case "IEC61131_WCHAR":
		return ModbusDataTypeSizes_WCHAR
	case "IEC61131_WORD":
		return ModbusDataTypeSizes_WORD
	case "IEC61131_WSTRING":
		return ModbusDataTypeSizes_WSTRING
	}
	return ""
}

func ModbusDataTypeSizesByName(value string) ModbusDataTypeSizes {
	switch value {
	case "BOOL":
		return ModbusDataTypeSizes_BOOL
	case "BYTE":
		return ModbusDataTypeSizes_BYTE
	case "CHAR":
		return ModbusDataTypeSizes_CHAR
	case "DATE":
		return ModbusDataTypeSizes_DATE
	case "DATE_AND_TIME":
		return ModbusDataTypeSizes_DATE_AND_TIME
	case "DINT":
		return ModbusDataTypeSizes_DINT
	case "DWORD":
		return ModbusDataTypeSizes_DWORD
	case "INT":
		return ModbusDataTypeSizes_INT
	case "LDATE":
		return ModbusDataTypeSizes_LDATE
	case "LDATE_AND_TIME":
		return ModbusDataTypeSizes_LDATE_AND_TIME
	case "LINT":
		return ModbusDataTypeSizes_LINT
	case "LREAL":
		return ModbusDataTypeSizes_LREAL
	case "LTIME":
		return ModbusDataTypeSizes_LTIME
	case "LTIME_OF_DAY":
		return ModbusDataTypeSizes_LTIME_OF_DAY
	case "LWORD":
		return ModbusDataTypeSizes_LWORD
	case "REAL":
		return ModbusDataTypeSizes_REAL
	case "SINT":
		return ModbusDataTypeSizes_SINT
	case "STRING":
		return ModbusDataTypeSizes_STRING
	case "TIME":
		return ModbusDataTypeSizes_TIME
	case "TIME_OF_DAY":
		return ModbusDataTypeSizes_TIME_OF_DAY
	case "UDINT":
		return ModbusDataTypeSizes_UDINT
	case "UINT":
		return ModbusDataTypeSizes_UINT
	case "ULINT":
		return ModbusDataTypeSizes_ULINT
	case "USINT":
		return ModbusDataTypeSizes_USINT
	case "WCHAR":
		return ModbusDataTypeSizes_WCHAR
	case "WORD":
		return ModbusDataTypeSizes_WORD
	case "WSTRING":
		return ModbusDataTypeSizes_WSTRING
	}
	return ""
}

func CastModbusDataTypeSizes(structType interface{}) ModbusDataTypeSizes {
	castFunc := func(typ interface{}) ModbusDataTypeSizes {
		if sModbusDataTypeSizes, ok := typ.(ModbusDataTypeSizes); ok {
			return sModbusDataTypeSizes
		}
		return ""
	}
	return castFunc(structType)
}

func (m ModbusDataTypeSizes) LengthInBits() uint16 {
	return 0
}

func (m ModbusDataTypeSizes) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func (e ModbusDataTypeSizes) String() string {
	switch e {
	case ModbusDataTypeSizes_BOOL:
		return "BOOL"
	case ModbusDataTypeSizes_BYTE:
		return "BYTE"
	case ModbusDataTypeSizes_CHAR:
		return "CHAR"
	case ModbusDataTypeSizes_DATE:
		return "DATE"
	case ModbusDataTypeSizes_DATE_AND_TIME:
		return "DATE_AND_TIME"
	case ModbusDataTypeSizes_DINT:
		return "DINT"
	case ModbusDataTypeSizes_DWORD:
		return "DWORD"
	case ModbusDataTypeSizes_INT:
		return "INT"
	case ModbusDataTypeSizes_LDATE:
		return "LDATE"
	case ModbusDataTypeSizes_LDATE_AND_TIME:
		return "LDATE_AND_TIME"
	case ModbusDataTypeSizes_LINT:
		return "LINT"
	case ModbusDataTypeSizes_LREAL:
		return "LREAL"
	case ModbusDataTypeSizes_LTIME:
		return "LTIME"
	case ModbusDataTypeSizes_LTIME_OF_DAY:
		return "LTIME_OF_DAY"
	case ModbusDataTypeSizes_LWORD:
		return "LWORD"
	case ModbusDataTypeSizes_REAL:
		return "REAL"
	case ModbusDataTypeSizes_SINT:
		return "SINT"
	case ModbusDataTypeSizes_STRING:
		return "STRING"
	case ModbusDataTypeSizes_TIME:
		return "TIME"
	case ModbusDataTypeSizes_TIME_OF_DAY:
		return "TIME_OF_DAY"
	case ModbusDataTypeSizes_UDINT:
		return "UDINT"
	case ModbusDataTypeSizes_UINT:
		return "UINT"
	case ModbusDataTypeSizes_ULINT:
		return "ULINT"
	case ModbusDataTypeSizes_USINT:
		return "USINT"
	case ModbusDataTypeSizes_WCHAR:
		return "WCHAR"
	case ModbusDataTypeSizes_WORD:
		return "WORD"
	case ModbusDataTypeSizes_WSTRING:
		return "WSTRING"
	}
	return ""
}
