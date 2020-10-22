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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/model/values/iec61131"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    api "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
)

func DataItemParse(io *spi.ReadBuffer, dataType uint8, numberOfValues uint8) (api.PlcValue, error) {
    switch {
        case dataType == 1 && numberOfValues == 1: // BOOL


            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcBOOL(value), nil
        case dataType == 1: // BOOL

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 10 && numberOfValues == 1: // BYTE

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcBYTE(value), nil
        case dataType == 10: // BYTE

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 11 && numberOfValues == 1: // WORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcWORD(value), nil
        case dataType == 11: // WORD

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 12 && numberOfValues == 1: // DWORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcDWORD(value), nil
        case dataType == 12: // DWORD

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 13 && numberOfValues == 1: // LWORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcLWORD(value), nil
        case dataType == 13: // LWORD

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 20 && numberOfValues == 1: // SINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcSINT(value), nil
        case dataType == 20: // SINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 21 && numberOfValues == 1: // INT

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcINT(value), nil
        case dataType == 21: // INT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 22 && numberOfValues == 1: // DINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcDINT(value), nil
        case dataType == 22: // DINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 23 && numberOfValues == 1: // LINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcLINT(value), nil
        case dataType == 23: // LINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 24 && numberOfValues == 1: // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcUSINT(value), nil
        case dataType == 24: // USINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 25 && numberOfValues == 1: // UINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcUINT(value), nil
        case dataType == 25: // UINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 26 && numberOfValues == 1: // UDINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcUDINT(value), nil
        case dataType == 26: // UDINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 27 && numberOfValues == 1: // ULINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcULINT(value), nil
        case dataType == 27: // ULINT

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 30 && numberOfValues == 1: // REAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcREAL(value), nil
        case dataType == 30: // REAL

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 31 && numberOfValues == 1: // LREAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcLREAL(value), nil
        case dataType == 31: // LREAL

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 80 && numberOfValues == 1: // CHAR

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcCHAR(value), nil
        case dataType == 80: // CHAR

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
        case dataType == 81 && numberOfValues == 1: // WCHAR

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return iec61131.NewPlcWCHAR(value), nil
        case dataType == 81: // WCHAR

            // Array Field (value)
            var value []api.PlcValue
            return values.NewPlcList(value), nil
    }
    return nil, errors.New("unsupported type")
}

