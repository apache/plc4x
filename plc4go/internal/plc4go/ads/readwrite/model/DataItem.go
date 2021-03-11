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
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

// Code generated by build-utils. DO NOT EDIT.

func DataItemParse(io *utils.ReadBuffer, dataFormatName string, stringLength int32) (api.PlcValue, error) {
	switch {
	case dataFormatName == "IEC61131_BOOL": // BOOL

		// Reserved Field (Just skip the bytes)
		if _, _err := io.ReadUint8(7); _err != nil {
			return nil, errors.New("Error parsing reserved field " + _err.Error())
		}

		// Simple Field (value)
		value, _valueErr := io.ReadBit()
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcBOOL(value), nil
	case dataFormatName == "IEC61131_BYTE": // BitString

		// Simple Field (value)
		value, _valueErr := io.ReadUint8(8)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcBitString(value), nil
	case dataFormatName == "IEC61131_WORD": // BitString

		// Simple Field (value)
		value, _valueErr := io.ReadUint16(16)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcBitString(value), nil
	case dataFormatName == "IEC61131_DWORD": // BitString

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcBitString(value), nil
	case dataFormatName == "IEC61131_SINT": // SINT

		// Simple Field (value)
		value, _valueErr := io.ReadInt8(8)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcSINT(value), nil
	case dataFormatName == "IEC61131_USINT": // USINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint8(8)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcUSINT(value), nil
	case dataFormatName == "IEC61131_INT": // INT

		// Simple Field (value)
		value, _valueErr := io.ReadInt16(16)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcINT(value), nil
	case dataFormatName == "IEC61131_UINT": // UINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint16(16)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcUINT(value), nil
	case dataFormatName == "IEC61131_DINT": // DINT

		// Simple Field (value)
		value, _valueErr := io.ReadInt32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcDINT(value), nil
	case dataFormatName == "IEC61131_UDINT": // UDINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcUDINT(value), nil
	case dataFormatName == "IEC61131_LINT": // LINT

		// Simple Field (value)
		value, _valueErr := io.ReadInt64(64)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcLINT(value), nil
	case dataFormatName == "IEC61131_ULINT": // ULINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint64(64)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcULINT(value), nil
	case dataFormatName == "IEC61131_REAL": // REAL

		// Simple Field (value)
		value, _valueErr := io.ReadFloat32(true, 8, 23)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcREAL(value), nil
	case dataFormatName == "IEC61131_LREAL": // LREAL

		// Simple Field (value)
		value, _valueErr := io.ReadFloat64(true, 11, 52)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcLREAL(value), nil
	case dataFormatName == "IEC61131_CHAR": // STRING
	case dataFormatName == "IEC61131_WCHAR": // STRING
	case dataFormatName == "IEC61131_STRING": // STRING

		// Manual Field (value)
		value, _valueErr := StaticHelperParseAmsString(io, stringLength, "UTF-8")
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcSTRING(value), nil
	case dataFormatName == "IEC61131_WSTRING": // STRING

		// Manual Field (value)
		value, _valueErr := StaticHelperParseAmsString(io, stringLength, "UTF-16")
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcSTRING(value), nil
	case dataFormatName == "IEC61131_TIME": // TIME

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	case dataFormatName == "IEC61131_LTIME": // LTIME

		// Simple Field (value)
		value, _valueErr := io.ReadInt64(64)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcLTIME(value), nil
	case dataFormatName == "IEC61131_DATE": // DATE

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcDATE(value), nil
	case dataFormatName == "IEC61131_TIME_OF_DAY": // TIME_OF_DAY

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME_OF_DAY(value), nil
	case dataFormatName == "IEC61131_DATE_AND_TIME": // DATE_AND_TIME

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcDATE_AND_TIME(value), nil
	}
	return nil, errors.New("unsupported type")
}

func DataItemSerialize(io *utils.WriteBuffer, value api.PlcValue, dataFormatName string, stringLength int32) error {
	switch {
	case dataFormatName == "IEC61131_BOOL": // BOOL

		// Reserved Field (Just skip the bytes)
		if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
			return errors.New("Error serializing reserved field " + _err.Error())
		}

		// Simple Field (value)
		if _err := io.WriteBit(value.GetBool()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_BYTE": // BitString

		// Simple Field (value)
		if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_WORD": // BitString

		// Simple Field (value)
		if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_DWORD": // BitString

		// Simple Field (value)
		if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_SINT": // SINT

		// Simple Field (value)
		if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_USINT": // USINT

		// Simple Field (value)
		if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_INT": // INT

		// Simple Field (value)
		if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_UINT": // UINT

		// Simple Field (value)
		if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_DINT": // DINT

		// Simple Field (value)
		if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_UDINT": // UDINT

		// Simple Field (value)
		if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_LINT": // LINT

		// Simple Field (value)
		if _err := io.WriteInt64(64, value.GetInt64()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_ULINT": // ULINT

		// Simple Field (value)
		if _err := io.WriteUint64(64, value.GetUint64()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_REAL": // REAL

		// Simple Field (value)
		if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_LREAL": // LREAL

		// Simple Field (value)
		if _err := io.WriteFloat64(64, value.GetFloat64()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_CHAR": // STRING
	case dataFormatName == "IEC61131_WCHAR": // STRING
	case dataFormatName == "IEC61131_STRING": // STRING

		// Manual Field (value)
		_valueErr := StaticHelperSerializeAmsString(io, value, stringLength, "UTF-8")
		if _valueErr != nil {
			return errors.New("Error serializing 'value' field " + _valueErr.Error())
		}
	case dataFormatName == "IEC61131_WSTRING": // STRING

		// Manual Field (value)
		_valueErr := StaticHelperSerializeAmsString(io, value, stringLength, "UTF-16")
		if _valueErr != nil {
			return errors.New("Error serializing 'value' field " + _valueErr.Error())
		}
	case dataFormatName == "IEC61131_TIME": // TIME

		// Simple Field (value)
		if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_LTIME": // LTIME

		// Simple Field (value)
		if _err := io.WriteInt64(64, value.GetInt64()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_DATE": // DATE

		// Simple Field (value)
		if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_TIME_OF_DAY": // TIME_OF_DAY

		// Simple Field (value)
		if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	case dataFormatName == "IEC61131_DATE_AND_TIME": // DATE_AND_TIME

		// Simple Field (value)
		if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
			return errors.New("Error serializing 'value' field " + _err.Error())
		}
	default:

		return errors.New("unsupported type")
	}
	return nil
}
