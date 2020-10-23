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

func DataItemParse(io *utils.ReadBuffer, dataType string, numberOfValues uint16) (api.PlcValue, error) {
    switch {
        case dataType == "IEC61131_BOOL" && numberOfValues == 1: // BOOL

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(7)

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case dataType == "IEC61131_BOOL": // BOOL

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_BYTE" && numberOfValues == 1: // BYTE

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBYTE(value), nil
        case dataType == "IEC61131_BYTE": // BYTE

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_WORD" && numberOfValues == 1: // WORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcWORD(value), nil
        case dataType == "IEC61131_WORD": // WORD

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_DWORD" && numberOfValues == 1: // DWORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDWORD(value), nil
        case dataType == "IEC61131_DWORD": // DWORD

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_LWORD" && numberOfValues == 1: // LWORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLWORD(value), nil
        case dataType == "IEC61131_LWORD": // LWORD

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_SINT" && numberOfValues == 1: // SINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case dataType == "IEC61131_SINT": // SINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_INT" && numberOfValues == 1: // INT

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case dataType == "IEC61131_INT": // INT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_DINT" && numberOfValues == 1: // DINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case dataType == "IEC61131_DINT": // DINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_LINT" && numberOfValues == 1: // LINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLINT(value), nil
        case dataType == "IEC61131_LINT": // LINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_USINT" && numberOfValues == 1: // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case dataType == "IEC61131_USINT": // USINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_UINT" && numberOfValues == 1: // UINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case dataType == "IEC61131_UINT": // UINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_UDINT" && numberOfValues == 1: // UDINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case dataType == "IEC61131_UDINT": // UDINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_ULINT" && numberOfValues == 1: // ULINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcULINT(value), nil
        case dataType == "IEC61131_ULINT": // ULINT

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_REAL" && numberOfValues == 1: // REAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case dataType == "IEC61131_REAL": // REAL

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_LREAL" && numberOfValues == 1: // LREAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLREAL(value), nil
        case dataType == "IEC61131_LREAL": // LREAL

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_CHAR" && numberOfValues == 1: // CHAR

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcCHAR(value), nil
        case dataType == "IEC61131_CHAR": // CHAR

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
        case dataType == "IEC61131_WCHAR" && numberOfValues == 1: // WCHAR

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcWCHAR(value), nil
        case dataType == "IEC61131_WCHAR": // WCHAR

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(numberOfValues); i++ {
                _item, _itemErr := DataItemParse(io, dataType, uint16(1))
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, _item)
            }
            return values.NewPlcList(value), nil
    }
    return nil, errors.New("unsupported type")
}

