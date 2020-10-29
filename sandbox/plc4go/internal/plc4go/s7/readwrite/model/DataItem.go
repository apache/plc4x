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
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/model/values"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
	api "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
)

func DataItemParse(io *utils.ReadBuffer, dataProtocolId string, stringLength int32) (api.PlcValue, error) {
	switch {
	case dataProtocolId == "IEC61131_BOOL": // BOOL

		// Reserved Field (Just skip the bytes)
		io.ReadUint8(7)

		// Simple Field (value)
		value, _valueErr := io.ReadBit()
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcBOOL(value), nil
	case dataProtocolId == "IEC61131_BYTE": // BOOL

		// Array Field (value)
		var value []api.PlcValue
		for i := 0; i < int((8)); i++ {
			_item, _itemErr := DataItemParse(io, dataProtocolId, int32(1))
			if _itemErr != nil {
				return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
			}
			value = append(value, _item)
		}
		return values.NewPlcList(value), nil
	case dataProtocolId == "IEC61131_WORD": // BOOL

		// Array Field (value)
		var value []api.PlcValue
		for i := 0; i < int((16)); i++ {
			_item, _itemErr := DataItemParse(io, dataProtocolId, int32(1))
			if _itemErr != nil {
				return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
			}
			value = append(value, _item)
		}
		return values.NewPlcList(value), nil
	case dataProtocolId == "IEC61131_DWORD": // BOOL

		// Array Field (value)
		var value []api.PlcValue
		for i := 0; i < int((32)); i++ {
			_item, _itemErr := DataItemParse(io, dataProtocolId, int32(1))
			if _itemErr != nil {
				return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
			}
			value = append(value, _item)
		}
		return values.NewPlcList(value), nil
	case dataProtocolId == "IEC61131_LWORD": // BOOL

		// Array Field (value)
		var value []api.PlcValue
		for i := 0; i < int((64)); i++ {
			_item, _itemErr := DataItemParse(io, dataProtocolId, int32(1))
			if _itemErr != nil {
				return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
			}
			value = append(value, _item)
		}
		return values.NewPlcList(value), nil
	case dataProtocolId == "IEC61131_SINT": // SINT

		// Simple Field (value)
		value, _valueErr := io.ReadInt8(8)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcSINT(value), nil
	case dataProtocolId == "IEC61131_USINT": // USINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint8(8)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcUSINT(value), nil
	case dataProtocolId == "IEC61131_INT": // INT

		// Simple Field (value)
		value, _valueErr := io.ReadInt16(16)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcINT(value), nil
	case dataProtocolId == "IEC61131_UINT": // UINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint16(16)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcUINT(value), nil
	case dataProtocolId == "IEC61131_DINT": // DINT

		// Simple Field (value)
		value, _valueErr := io.ReadInt32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcDINT(value), nil
	case dataProtocolId == "IEC61131_UDINT": // UDINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcUDINT(value), nil
	case dataProtocolId == "IEC61131_LINT": // LINT

		// Simple Field (value)
		value, _valueErr := io.ReadInt64(64)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcLINT(value), nil
	case dataProtocolId == "IEC61131_ULINT": // ULINT

		// Simple Field (value)
		value, _valueErr := io.ReadUint64(64)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcULINT(value), nil
	case dataProtocolId == "IEC61131_REAL": // REAL

		// Simple Field (value)
		value, _valueErr := io.ReadFloat32(32)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcREAL(value), nil
	case dataProtocolId == "IEC61131_LREAL": // LREAL

		// Simple Field (value)
		value, _valueErr := io.ReadFloat64(64)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcLREAL(value), nil
	case dataProtocolId == "IEC61131_CHAR": // CHAR

		// Manual Field (value)
		value, _valueErr := StaticHelperParseS7Char(io, "UTF-8")
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcCHAR(value), nil
	case dataProtocolId == "IEC61131_WCHAR": // CHAR

		// Manual Field (value)
		value, _valueErr := StaticHelperParseS7Char(io, "UTF-16")
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcCHAR(value), nil
	case dataProtocolId == "IEC61131_STRING": // STRING

		// Manual Field (value)
		value, _valueErr := StaticHelperParseS7String(io, stringLength, "UTF-8")
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcSTRING(value), nil
	case dataProtocolId == "IEC61131_WSTRING": // STRING

		// Manual Field (value)
		value, _valueErr := StaticHelperParseS7String(io, stringLength, "UTF-16")
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcSTRING(value), nil
	case dataProtocolId == "IEC61131_TIME": // Time

		// Manual Field (value)
		value, _valueErr := StaticHelperParseTiaTime(io)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	case dataProtocolId == "S7_S5TIME": // Time

		// Manual Field (value)
		value, _valueErr := StaticHelperParseS5Time(io)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	case dataProtocolId == "IEC61131_LTIME": // Time

		// Manual Field (value)
		value, _valueErr := StaticHelperParseTiaLTime(io)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	case dataProtocolId == "IEC61131_DATE": // Date

		// Manual Field (value)
		value, _valueErr := StaticHelperParseTiaDate(io)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	case dataProtocolId == "IEC61131_TIME_OF_DAY": // Time

		// Manual Field (value)
		value, _valueErr := StaticHelperParseTiaTimeOfDay(io)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	case dataProtocolId == "IEC61131_DATE_AND_TIME": // DateTime

		// Manual Field (value)
		value, _valueErr := StaticHelperParseTiaDateTime(io)
		if _valueErr != nil {
			return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
		}
		return values.NewPlcTIME(value), nil
	}
	return nil, errors.New("unsupported type")
}
