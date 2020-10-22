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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    api "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
)

func KnxDatapointParse(io *spi.ReadBuffer, mainNumber uint16, subNumber uint16) (api.PlcValue, error) {
    switch {
        case mainNumber == 1: // BOOL

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(7)

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case mainNumber == 2: // BOOL

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(6)

            // Simple Field (control)
            _, _controlErr := io.ReadBit()
            if _controlErr != nil {
                return nil, errors.New("Error parsing 'control' field " + _controlErr.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBOOL(value), nil
        case mainNumber == 21: // Struct
            _map := map[string]interface{}{}

            // Simple Field (b7)
            b7, _b7Err := io.ReadBit()
            if _b7Err != nil {
                return nil, errors.New("Error parsing 'b7' field " + _b7Err.Error())
            }
            _map["Struct"] = b7

            // Simple Field (b6)
            b6, _b6Err := io.ReadBit()
            if _b6Err != nil {
                return nil, errors.New("Error parsing 'b6' field " + _b6Err.Error())
            }
            _map["Struct"] = b6

            // Simple Field (b5)
            b5, _b5Err := io.ReadBit()
            if _b5Err != nil {
                return nil, errors.New("Error parsing 'b5' field " + _b5Err.Error())
            }
            _map["Struct"] = b5

            // Simple Field (b4)
            b4, _b4Err := io.ReadBit()
            if _b4Err != nil {
                return nil, errors.New("Error parsing 'b4' field " + _b4Err.Error())
            }
            _map["Struct"] = b4

            // Simple Field (b3)
            b3, _b3Err := io.ReadBit()
            if _b3Err != nil {
                return nil, errors.New("Error parsing 'b3' field " + _b3Err.Error())
            }
            _map["Struct"] = b3

            // Simple Field (b2)
            b2, _b2Err := io.ReadBit()
            if _b2Err != nil {
                return nil, errors.New("Error parsing 'b2' field " + _b2Err.Error())
            }
            _map["Struct"] = b2

            // Simple Field (b1)
            b1, _b1Err := io.ReadBit()
            if _b1Err != nil {
                return nil, errors.New("Error parsing 'b1' field " + _b1Err.Error())
            }
            _map["Struct"] = b1

            // Simple Field (b0)
            b0, _b0Err := io.ReadBit()
            if _b0Err != nil {
                return nil, errors.New("Error parsing 'b0' field " + _b0Err.Error())
            }
            _map["Struct"] = b0
        case mainNumber == 3: // USINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(4)

            // Simple Field (control)
            _, _controlErr := io.ReadBit()
            if _controlErr != nil {
                return nil, errors.New("Error parsing 'control' field " + _controlErr.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(3)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case mainNumber == 18: // USINT

            // Simple Field (control)
            _, _controlErr := io.ReadBit()
            if _controlErr != nil {
                return nil, errors.New("Error parsing 'control' field " + _controlErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(1)

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(6)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case mainNumber == 17: // USINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(2)

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(6)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case mainNumber == 5: // USINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case mainNumber == 7: // UINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case mainNumber == 12: // UDINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case mainNumber == 6 && subNumber == 20: // SINT

            // Simple Field (a)
            _, _aErr := io.ReadBit()
            if _aErr != nil {
                return nil, errors.New("Error parsing 'a' field " + _aErr.Error())
            }

            // Simple Field (b)
            _, _bErr := io.ReadBit()
            if _bErr != nil {
                return nil, errors.New("Error parsing 'b' field " + _bErr.Error())
            }

            // Simple Field (c)
            _, _cErr := io.ReadBit()
            if _cErr != nil {
                return nil, errors.New("Error parsing 'c' field " + _cErr.Error())
            }

            // Simple Field (d)
            _, _dErr := io.ReadBit()
            if _dErr != nil {
                return nil, errors.New("Error parsing 'd' field " + _dErr.Error())
            }

            // Simple Field (e)
            _, _eErr := io.ReadBit()
            if _eErr != nil {
                return nil, errors.New("Error parsing 'e' field " + _eErr.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case mainNumber == 6: // SINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case mainNumber == 8: // INT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case mainNumber == 13: // DINT

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case mainNumber == 9: // REAL

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Manual Field (value)
            value, _valueErr := KnxHelperBytesToF16(io)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case mainNumber == 14: // REAL

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case mainNumber == 4: // STRING

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadString(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case mainNumber == 16: // STRING

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(8)

            // Simple Field (value)
            value, _valueErr := io.ReadString(112)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case mainNumber == 10: // Time

            // Simple Field (day)
            _, _dayErr := io.ReadUint8(3)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }

            // Simple Field (hours)
            _, _hoursErr := io.ReadUint8(5)
            if _hoursErr != nil {
                return nil, errors.New("Error parsing 'hours' field " + _hoursErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(2)

            // Simple Field (minutes)
            _, _minutesErr := io.ReadUint8(6)
            if _minutesErr != nil {
                return nil, errors.New("Error parsing 'minutes' field " + _minutesErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(2)

            // Simple Field (seconds)
            _, _secondsErr := io.ReadUint8(6)
            if _secondsErr != nil {
                return nil, errors.New("Error parsing 'seconds' field " + _secondsErr.Error())
            }
        case mainNumber == 11: // Date

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(3)

            // Simple Field (day)
            _, _dayErr := io.ReadUint8(5)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(4)

            // Simple Field (month)
            _, _monthErr := io.ReadUint8(4)
            if _monthErr != nil {
                return nil, errors.New("Error parsing 'month' field " + _monthErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(1)

            // Simple Field (year)
            _, _yearErr := io.ReadUint8(6)
            if _yearErr != nil {
                return nil, errors.New("Error parsing 'year' field " + _yearErr.Error())
            }
        case mainNumber == 19: // DateTime

            // Simple Field (year)
            _, _yearErr := io.ReadUint8(8)
            if _yearErr != nil {
                return nil, errors.New("Error parsing 'year' field " + _yearErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(4)

            // Simple Field (month)
            _, _monthErr := io.ReadUint8(4)
            if _monthErr != nil {
                return nil, errors.New("Error parsing 'month' field " + _monthErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(3)

            // Simple Field (day)
            _, _dayErr := io.ReadUint8(5)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }

            // Simple Field (dayOfWeek)
            _, _dayOfWeekErr := io.ReadUint8(3)
            if _dayOfWeekErr != nil {
                return nil, errors.New("Error parsing 'dayOfWeek' field " + _dayOfWeekErr.Error())
            }

            // Simple Field (hours)
            _, _hoursErr := io.ReadUint8(5)
            if _hoursErr != nil {
                return nil, errors.New("Error parsing 'hours' field " + _hoursErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(2)

            // Simple Field (minutes)
            _, _minutesErr := io.ReadUint8(6)
            if _minutesErr != nil {
                return nil, errors.New("Error parsing 'minutes' field " + _minutesErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            io.ReadUint8(2)

            // Simple Field (seconds)
            _, _secondsErr := io.ReadUint8(6)
            if _secondsErr != nil {
                return nil, errors.New("Error parsing 'seconds' field " + _secondsErr.Error())
            }

            // Simple Field (fault)
            _, _faultErr := io.ReadBit()
            if _faultErr != nil {
                return nil, errors.New("Error parsing 'fault' field " + _faultErr.Error())
            }

            // Simple Field (workingDay)
            _, _workingDayErr := io.ReadBit()
            if _workingDayErr != nil {
                return nil, errors.New("Error parsing 'workingDay' field " + _workingDayErr.Error())
            }

            // Simple Field (workingDayValid)
            _, _workingDayValidErr := io.ReadBit()
            if _workingDayValidErr != nil {
                return nil, errors.New("Error parsing 'workingDayValid' field " + _workingDayValidErr.Error())
            }

            // Simple Field (yearValid)
            _, _yearValidErr := io.ReadBit()
            if _yearValidErr != nil {
                return nil, errors.New("Error parsing 'yearValid' field " + _yearValidErr.Error())
            }

            // Simple Field (dayAndMonthValid)
            _, _dayAndMonthValidErr := io.ReadBit()
            if _dayAndMonthValidErr != nil {
                return nil, errors.New("Error parsing 'dayAndMonthValid' field " + _dayAndMonthValidErr.Error())
            }

            // Simple Field (dayOfWeekValid)
            _, _dayOfWeekValidErr := io.ReadBit()
            if _dayOfWeekValidErr != nil {
                return nil, errors.New("Error parsing 'dayOfWeekValid' field " + _dayOfWeekValidErr.Error())
            }

            // Simple Field (timeValid)
            _, _timeValidErr := io.ReadBit()
            if _timeValidErr != nil {
                return nil, errors.New("Error parsing 'timeValid' field " + _timeValidErr.Error())
            }

            // Simple Field (standardSummerTime)
            _, _standardSummerTimeErr := io.ReadBit()
            if _standardSummerTimeErr != nil {
                return nil, errors.New("Error parsing 'standardSummerTime' field " + _standardSummerTimeErr.Error())
            }

            // Simple Field (clockQuality)
            _, _clockQualityErr := io.ReadBit()
            if _clockQualityErr != nil {
                return nil, errors.New("Error parsing 'clockQuality' field " + _clockQualityErr.Error())
            }
        case mainNumber == 15: // Struct
            _map := map[string]interface{}{}

            // Simple Field (D6)
            D6, _D6Err := io.ReadUint8(4)
            if _D6Err != nil {
                return nil, errors.New("Error parsing 'D6' field " + _D6Err.Error())
            }
            _map["Struct"] = D6

            // Simple Field (D5)
            D5, _D5Err := io.ReadUint8(4)
            if _D5Err != nil {
                return nil, errors.New("Error parsing 'D5' field " + _D5Err.Error())
            }
            _map["Struct"] = D5

            // Simple Field (D4)
            D4, _D4Err := io.ReadUint8(4)
            if _D4Err != nil {
                return nil, errors.New("Error parsing 'D4' field " + _D4Err.Error())
            }
            _map["Struct"] = D4

            // Simple Field (D3)
            D3, _D3Err := io.ReadUint8(4)
            if _D3Err != nil {
                return nil, errors.New("Error parsing 'D3' field " + _D3Err.Error())
            }
            _map["Struct"] = D3

            // Simple Field (D2)
            D2, _D2Err := io.ReadUint8(4)
            if _D2Err != nil {
                return nil, errors.New("Error parsing 'D2' field " + _D2Err.Error())
            }
            _map["Struct"] = D2

            // Simple Field (D1)
            D1, _D1Err := io.ReadUint8(4)
            if _D1Err != nil {
                return nil, errors.New("Error parsing 'D1' field " + _D1Err.Error())
            }
            _map["Struct"] = D1

            // Simple Field (BE)
            BE, _BEErr := io.ReadBit()
            if _BEErr != nil {
                return nil, errors.New("Error parsing 'BE' field " + _BEErr.Error())
            }
            _map["Struct"] = BE

            // Simple Field (BP)
            BP, _BPErr := io.ReadBit()
            if _BPErr != nil {
                return nil, errors.New("Error parsing 'BP' field " + _BPErr.Error())
            }
            _map["Struct"] = BP

            // Simple Field (BD)
            BD, _BDErr := io.ReadBit()
            if _BDErr != nil {
                return nil, errors.New("Error parsing 'BD' field " + _BDErr.Error())
            }
            _map["Struct"] = BD

            // Simple Field (BC)
            BC, _BCErr := io.ReadBit()
            if _BCErr != nil {
                return nil, errors.New("Error parsing 'BC' field " + _BCErr.Error())
            }
            _map["Struct"] = BC

            // Simple Field (index)
            index, _indexErr := io.ReadUint8(4)
            if _indexErr != nil {
                return nil, errors.New("Error parsing 'index' field " + _indexErr.Error())
            }
            _map["Struct"] = index
    }
    return nil, errors.New("unsupported type")
}

