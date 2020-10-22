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
package iec61131

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
	"reflect"
	"time"
)

const (
	// Bit & Bit-Strings
	IEC61131_BOOL  string = "IEC61131_BOOL"
	IEC61131_BYTE  string = "IEC61131_BYTE"
	IEC61131_WORD  string = "IEC61131_WORD"
	IEC61131_DWORD string = "IEC61131_DWORD"
	IEC61131_LWORD string = "IEC61131_LWORD"

	// Integers
	IEC61131_USINT string = "IEC61131_USINT"
	IEC61131_UINT  string = "IEC61131_UINT"
	IEC61131_UDINT string = "IEC61131_UDINT"
	IEC61131_ULINT string = "IEC61131_ULINT"
	IEC61131_SINT  string = "IEC61131_SINT"
	IEC61131_INT   string = "IEC61131_INT"
	IEC61131_DINT  string = "IEC61131_DINT"
	IEC61131_LINT  string = "IEC61131_LINT"

	// Floating Point Values
	IEC61131_REAL  string = "IEC61131_REAL"
	IEC61131_LREAL string = "IEC61131_LREAL"

	// Temporal Values
	IEC61131_TIME          string = "IEC61131_TIME"
	IEC61131_DATE          string = "IEC61131_DATE"
	IEC61131_TIME_OF_DAY   string = "IEC61131_TIME_OF_DAY"
	IEC61131_DATE_AND_TIME string = "IEC61131_DATE_AND_TIME"

	// Chars and Strings
	IEC61131_CHAR    string = "IEC61131_CHAR"
	IEC61131_WCHAR   string = "IEC61131_WCHAR"
	IEC61131_STRING  string = "IEC61131_STRING"
	IEC61131_WSTRING string = "IEC61131_WSTRING"
)

type IEC61131ValueHandler struct {
	spi.PlcValueHandler
}

func NewIEC61131ValueHandler() IEC61131ValueHandler {
	return IEC61131ValueHandler{}
}

func (m IEC61131ValueHandler) NewPlcValue(typeName string, value interface{}) (values.PlcValue, error) {
	switch typeName {
	// Bit & Bit-Strings
	case IEC61131_BOOL:
		casted, ok := value.(bool)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to bool")
		}
		return NewPlcBOOL(casted), nil
	case IEC61131_BYTE:
		casted, ok := value.(uint8)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint8")
		}
		return NewPlcBYTE(casted), nil
	case IEC61131_WORD:
		casted, ok := value.(uint16)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint16")
		}
		return NewPlcWORD(casted), nil
	case IEC61131_DWORD:
		casted, ok := value.(uint32)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint32")
		}
		return NewPlcDWORD(casted), nil
	case IEC61131_LWORD:
		casted, ok := value.(uint64)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint64")
		}
		return NewPlcLWORD(casted), nil

	// Integers
	case IEC61131_USINT:
		casted, ok := value.(uint8)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint8")
		}
		return NewPlcUSINT(casted), nil
	case IEC61131_UINT:
		casted, ok := value.(uint16)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint16")
		}
		return NewPlcUINT(casted), nil
	case IEC61131_UDINT:
		casted, ok := value.(uint32)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint32")
		}
		return NewPlcUDINT(casted), nil
	case IEC61131_ULINT:
		casted, ok := value.(uint64)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint64")
		}
		return NewPlcULINT(casted), nil
	case IEC61131_SINT:
		casted, ok := value.(int8)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int8")
		}
		return NewPlcSINT(casted), nil
	case IEC61131_INT:
		casted, ok := value.(int16)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int16")
		}
		return NewPlcINT(casted), nil
	case IEC61131_DINT:
		casted, ok := value.(int32)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int32")
		}
		return NewPlcDINT(casted), nil
	case IEC61131_LINT:
		casted, ok := value.(int64)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to int64")
		}
		return NewPlcLINT(casted), nil

	// Floating Point Values
	case IEC61131_REAL:
		casted, ok := value.(float32)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to float32")
		}
		return NewPlcREAL(casted), nil
	case IEC61131_LREAL:
		casted, ok := value.(float64)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to float64")
		}
		return NewPlcLREAL(casted), nil

	// Temporal Values
	case IEC61131_TIME:
		casted, ok := value.(uint32)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint32")
		}
		return NewPlcTIME(casted), nil
	case IEC61131_DATE:
		casted, ok := value.(time.Time)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
		}
		return NewPlcDATE(casted), nil
	case IEC61131_TIME_OF_DAY:
		casted, ok := value.(time.Time)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
		}
		return NewPlcTIMEOFDAY(casted), nil
	case IEC61131_DATE_AND_TIME:
		casted, ok := value.(time.Time)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to time.Time")
		}
		return NewPlcDATEANDTIME(casted), nil

	// Chars and Strings
	case IEC61131_CHAR:
		casted, ok := value.(uint8)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint8")
		}
		return NewPlcCHAR(casted), nil
	case IEC61131_WCHAR:
		casted, ok := value.(uint16)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to uint16")
		}
		return NewPlcWCHAR(casted), nil
	case IEC61131_STRING:
		/*casted, ok := value.([]uint8)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to []uint8")
		}
		return NewPlcSTRING(casted), nil*/
	case IEC61131_WSTRING:
		casted, ok := value.([]uint16)
		if !ok {
			return nil, errors.New("couldn't cast value of type " + reflect.TypeOf(value).Name() + " to []uint16")
		}
		return NewPlcWSTRING(casted), nil
	}

	return nil, errors.New("Unsupported type " + typeName)
}
