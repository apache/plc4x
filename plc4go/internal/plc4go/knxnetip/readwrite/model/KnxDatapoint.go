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
            "time"
)

func KnxDatapointParse(io *utils.ReadBuffer, formatName string) (api.PlcValue, error) {
    switch {
        case formatName == "B1": // BOOL

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
        case formatName == "B2": // Struct
            _map := map[string]interface{}{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (control)
            control, _controlErr := io.ReadBit()
            if _controlErr != nil {
                return nil, errors.New("Error parsing 'control' field " + _controlErr.Error())
            }
            _map["Struct"] = control

            // Simple Field (value)
            value, _valueErr := io.ReadBit()
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            _map["Struct"] = value
            return values.NewPlcStruct(_map), nil
        case formatName == "B1U3": // Struct
            _map := map[string]interface{}{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (control)
            control, _controlErr := io.ReadBit()
            if _controlErr != nil {
                return nil, errors.New("Error parsing 'control' field " + _controlErr.Error())
            }
            _map["Struct"] = control

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(3)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            _map["Struct"] = value
            return values.NewPlcStruct(_map), nil
        case formatName == "A8_ASCII": // STRING

            // Simple Field (value)
            value, _valueErr := io.ReadString(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case formatName == "A8_8859_1": // STRING

            // Simple Field (value)
            value, _valueErr := io.ReadString(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case formatName == "U8": // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case formatName == "V8": // SINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case formatName == "F16": // REAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case formatName == "N3N5r2N6r2N6": // TIME_OF_DAY

            // Simple Field (day)
            _, _dayErr := io.ReadUint8(3)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }

            // Simple Field (hour)
            hour, _hourErr := io.ReadUint8(5)
            if _hourErr != nil {
                return nil, errors.New("Error parsing 'hour' field " + _hourErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (minutes)
            minutes, _minutesErr := io.ReadUint8(6)
            if _minutesErr != nil {
                return nil, errors.New("Error parsing 'minutes' field " + _minutesErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (seconds)
            seconds, _secondsErr := io.ReadUint8(6)
            if _secondsErr != nil {
                return nil, errors.New("Error parsing 'seconds' field " + _secondsErr.Error())
            }
            value := time.Date(0,0,0, int(hour), int(minutes), int(seconds), 0, nil)
            return values.NewPlcTIME_OF_DAY(value), nil
        case formatName == "r3N5r4N4r1U7": // DATE

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (day)
            day, _dayErr := io.ReadUint8(5)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (month)
            month, _monthErr := io.ReadUint8(4)
            if _monthErr != nil {
                return nil, errors.New("Error parsing 'month' field " + _monthErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (year)
            year, _yearErr := io.ReadUint8(7)
            if _yearErr != nil {
                return nil, errors.New("Error parsing 'year' field " + _yearErr.Error())
            }
            value := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, nil)
            return values.NewPlcDATE(value), nil
        case formatName == "U32": // UDINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case formatName == "V32": // DINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case formatName == "F32": // REAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case formatName == "U4U4U4U4U4U4B4N4": // Struct
            _map := map[string]interface{}{}

            // Simple Field (d6)
            d6, _d6Err := io.ReadUint8(4)
            if _d6Err != nil {
                return nil, errors.New("Error parsing 'd6' field " + _d6Err.Error())
            }
            _map["Struct"] = d6

            // Simple Field (d5)
            d5, _d5Err := io.ReadUint8(4)
            if _d5Err != nil {
                return nil, errors.New("Error parsing 'd5' field " + _d5Err.Error())
            }
            _map["Struct"] = d5

            // Simple Field (d4)
            d4, _d4Err := io.ReadUint8(4)
            if _d4Err != nil {
                return nil, errors.New("Error parsing 'd4' field " + _d4Err.Error())
            }
            _map["Struct"] = d4

            // Simple Field (d3)
            d3, _d3Err := io.ReadUint8(4)
            if _d3Err != nil {
                return nil, errors.New("Error parsing 'd3' field " + _d3Err.Error())
            }
            _map["Struct"] = d3

            // Simple Field (d2)
            d2, _d2Err := io.ReadUint8(4)
            if _d2Err != nil {
                return nil, errors.New("Error parsing 'd2' field " + _d2Err.Error())
            }
            _map["Struct"] = d2

            // Simple Field (d1)
            d1, _d1Err := io.ReadUint8(4)
            if _d1Err != nil {
                return nil, errors.New("Error parsing 'd1' field " + _d1Err.Error())
            }
            _map["Struct"] = d1

            // Simple Field (e)
            e, _eErr := io.ReadBit()
            if _eErr != nil {
                return nil, errors.New("Error parsing 'e' field " + _eErr.Error())
            }
            _map["Struct"] = e

            // Simple Field (p)
            p, _pErr := io.ReadBit()
            if _pErr != nil {
                return nil, errors.New("Error parsing 'p' field " + _pErr.Error())
            }
            _map["Struct"] = p

            // Simple Field (d)
            d, _dErr := io.ReadBit()
            if _dErr != nil {
                return nil, errors.New("Error parsing 'd' field " + _dErr.Error())
            }
            _map["Struct"] = d

            // Simple Field (c)
            c, _cErr := io.ReadBit()
            if _cErr != nil {
                return nil, errors.New("Error parsing 'c' field " + _cErr.Error())
            }
            _map["Struct"] = c

            // Simple Field (index)
            index, _indexErr := io.ReadUint8(4)
            if _indexErr != nil {
                return nil, errors.New("Error parsing 'index' field " + _indexErr.Error())
            }
            _map["Struct"] = index
            return values.NewPlcStruct(_map), nil
        case formatName == "A112_ASCII": // STRING

            // Simple Field (value)
            value, _valueErr := io.ReadString(112)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case formatName == "A112_8859_1": // STRING

            // Simple Field (value)
            value, _valueErr := io.ReadString(112)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSTRING(value), nil
        case formatName == "r2U6": // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(6)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case formatName == "B1r1U6": // Struct
            _map := map[string]interface{}{}

            // Simple Field (learn)
            learn, _learnErr := io.ReadBit()
            if _learnErr != nil {
                return nil, errors.New("Error parsing 'learn' field " + _learnErr.Error())
            }
            _map["Struct"] = learn

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (sceneNumber)
            sceneNumber, _sceneNumberErr := io.ReadUint8(6)
            if _sceneNumberErr != nil {
                return nil, errors.New("Error parsing 'sceneNumber' field " + _sceneNumberErr.Error())
            }
            _map["Struct"] = sceneNumber
            return values.NewPlcStruct(_map), nil
        case formatName == "U8r4U4r3U5U3U5r2U6r2U6B16": // DATE_AND_TIME

            // Simple Field (year)
            year, _yearErr := io.ReadUint8(8)
            if _yearErr != nil {
                return nil, errors.New("Error parsing 'year' field " + _yearErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (month)
            month, _monthErr := io.ReadUint8(4)
            if _monthErr != nil {
                return nil, errors.New("Error parsing 'month' field " + _monthErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (day)
            day, _dayErr := io.ReadUint8(5)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }

            // Simple Field (dayOfWeek)
            _, _dayOfWeekErr := io.ReadUint8(3)
            if _dayOfWeekErr != nil {
                return nil, errors.New("Error parsing 'dayOfWeek' field " + _dayOfWeekErr.Error())
            }

            // Simple Field (hour)
            hour, _hourErr := io.ReadUint8(5)
            if _hourErr != nil {
                return nil, errors.New("Error parsing 'hour' field " + _hourErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (minutes)
            minutes, _minutesErr := io.ReadUint8(6)
            if _minutesErr != nil {
                return nil, errors.New("Error parsing 'minutes' field " + _minutesErr.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (seconds)
            seconds, _secondsErr := io.ReadUint8(6)
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

            // Simple Field (noWorkingDay)
            _, _noWorkingDayErr := io.ReadBit()
            if _noWorkingDayErr != nil {
                return nil, errors.New("Error parsing 'noWorkingDay' field " + _noWorkingDayErr.Error())
            }

            // Simple Field (noYear)
            _, _noYearErr := io.ReadBit()
            if _noYearErr != nil {
                return nil, errors.New("Error parsing 'noYear' field " + _noYearErr.Error())
            }

            // Simple Field (noMonthAndDay)
            _, _noMonthAndDayErr := io.ReadBit()
            if _noMonthAndDayErr != nil {
                return nil, errors.New("Error parsing 'noMonthAndDay' field " + _noMonthAndDayErr.Error())
            }

            // Simple Field (noDayOfWeek)
            _, _noDayOfWeekErr := io.ReadBit()
            if _noDayOfWeekErr != nil {
                return nil, errors.New("Error parsing 'noDayOfWeek' field " + _noDayOfWeekErr.Error())
            }

            // Simple Field (noTime)
            _, _noTimeErr := io.ReadBit()
            if _noTimeErr != nil {
                return nil, errors.New("Error parsing 'noTime' field " + _noTimeErr.Error())
            }

            // Simple Field (standardSummerTime)
            _, _standardSummerTimeErr := io.ReadBit()
            if _standardSummerTimeErr != nil {
                return nil, errors.New("Error parsing 'standardSummerTime' field " + _standardSummerTimeErr.Error())
            }

            // Simple Field (clockWithSyncSignal)
            _, _clockWithSyncSignalErr := io.ReadBit()
            if _clockWithSyncSignalErr != nil {
                return nil, errors.New("Error parsing 'clockWithSyncSignal' field " + _clockWithSyncSignalErr.Error())
            }
            value := time.Date(int(year), time.Month(month), int(day), int(hour), int(minutes), int(seconds), 0, nil)
            return values.NewPlcDATE_AND_TIME(value), nil
        case formatName == "N8": // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case formatName == "B8": // BYTE

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcBYTE(value), nil
        case formatName == "B16": // WORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcWORD(value), nil
        case formatName == "U4U4": // Struct
            _map := map[string]interface{}{}

            // Simple Field (busy)
            busy, _busyErr := io.ReadUint8(4)
            if _busyErr != nil {
                return nil, errors.New("Error parsing 'busy' field " + _busyErr.Error())
            }
            _map["Struct"] = busy

            // Simple Field (nak)
            nak, _nakErr := io.ReadUint8(4)
            if _nakErr != nil {
                return nil, errors.New("Error parsing 'nak' field " + _nakErr.Error())
            }
            _map["Struct"] = nak
            return values.NewPlcStruct(_map), nil
        case formatName == "r1b1U6": // Struct
            _map := map[string]interface{}{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (sceneActive)
            sceneActive, _sceneActiveErr := io.ReadBit()
            if _sceneActiveErr != nil {
                return nil, errors.New("Error parsing 'sceneActive' field " + _sceneActiveErr.Error())
            }
            _map["Struct"] = sceneActive

            // Simple Field (sceneNumber)
            sceneNumber, _sceneNumberErr := io.ReadUint8(6)
            if _sceneNumberErr != nil {
                return nil, errors.New("Error parsing 'sceneNumber' field " + _sceneNumberErr.Error())
            }
            _map["Struct"] = sceneNumber
            return values.NewPlcStruct(_map), nil
        case formatName == "B32": // DWORD

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDWORD(value), nil
        case formatName == "V64": // LINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt64(64)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLINT(value), nil
        case formatName == "B24": // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((24)); i++ {
                _item, _itemErr := io.ReadBit()
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcBOOL(_item))
            }
            return values.NewPlcList(value), nil
        case formatName == "N3": // USINT

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(3)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case formatName == "B1Z8HeatingOrCoolingZ": // Struct
            _map := map[string]interface{}{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (heating)
            heating, _heatingErr := io.ReadBit()
            if _heatingErr != nil {
                return nil, errors.New("Error parsing 'heating' field " + _heatingErr.Error())
            }
            _map["Struct"] = heating

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "B1Z8BinaryValueZ": // Struct
            _map := map[string]interface{}{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (high)
            high, _highErr := io.ReadBit()
            if _highErr != nil {
                return nil, errors.New("Error parsing 'high' field " + _highErr.Error())
            }
            _map["Struct"] = high

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8HvacOperatingMode": // Struct
            _map := map[string]interface{}{}

            // Simple Field (hvacOperatingMode)
            hvacOperatingMode, _hvacOperatingModeErr := io.ReadUint8(8)
            if _hvacOperatingModeErr != nil {
                return nil, errors.New("Error parsing 'hvacOperatingMode' field " + _hvacOperatingModeErr.Error())
            }
            _map["Struct"] = hvacOperatingMode

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8DhwMode": // Struct
            _map := map[string]interface{}{}

            // Simple Field (dhwMode)
            dhwMode, _dhwModeErr := io.ReadUint8(8)
            if _dhwModeErr != nil {
                return nil, errors.New("Error parsing 'dhwMode' field " + _dhwModeErr.Error())
            }
            _map["Struct"] = dhwMode

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8HvacControllingMode": // Struct
            _map := map[string]interface{}{}

            // Simple Field (hvacControllingMode)
            hvacControllingMode, _hvacControllingModeErr := io.ReadUint8(8)
            if _hvacControllingModeErr != nil {
                return nil, errors.New("Error parsing 'hvacControllingMode' field " + _hvacControllingModeErr.Error())
            }
            _map["Struct"] = hvacControllingMode

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8EnableHeatingOrCoolingStage": // Struct
            _map := map[string]interface{}{}

            // Simple Field (enableHeatingOrCoolingStage)
            enableHeatingOrCoolingStage, _enableHeatingOrCoolingStageErr := io.ReadUint8(8)
            if _enableHeatingOrCoolingStageErr != nil {
                return nil, errors.New("Error parsing 'enableHeatingOrCoolingStage' field " + _enableHeatingOrCoolingStageErr.Error())
            }
            _map["Struct"] = enableHeatingOrCoolingStage

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8BuildingMode": // Struct
            _map := map[string]interface{}{}

            // Simple Field (buildingMode)
            buildingMode, _buildingModeErr := io.ReadUint8(8)
            if _buildingModeErr != nil {
                return nil, errors.New("Error parsing 'buildingMode' field " + _buildingModeErr.Error())
            }
            _map["Struct"] = buildingMode

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8OccupancyMode": // Struct
            _map := map[string]interface{}{}

            // Simple Field (occupancyMode)
            occupancyMode, _occupancyModeErr := io.ReadUint8(8)
            if _occupancyModeErr != nil {
                return nil, errors.New("Error parsing 'occupancyMode' field " + _occupancyModeErr.Error())
            }
            _map["Struct"] = occupancyMode

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "N8Z8EmergencyMode": // Struct
            _map := map[string]interface{}{}

            // Simple Field (hvacEmergencyMode)
            hvacEmergencyMode, _hvacEmergencyModeErr := io.ReadUint8(8)
            if _hvacEmergencyModeErr != nil {
                return nil, errors.New("Error parsing 'hvacEmergencyMode' field " + _hvacEmergencyModeErr.Error())
            }
            _map["Struct"] = hvacEmergencyMode

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U8Z8Rel": // Struct
            _map := map[string]interface{}{}

            // Simple Field (relValue)
            relValue, _relValueErr := io.ReadUint8(8)
            if _relValueErr != nil {
                return nil, errors.New("Error parsing 'relValue' field " + _relValueErr.Error())
            }
            _map["Struct"] = relValue

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U8Z8Counter": // Struct
            _map := map[string]interface{}{}

            // Simple Field (counterValue)
            counterValue, _counterValueErr := io.ReadUint8(8)
            if _counterValueErr != nil {
                return nil, errors.New("Error parsing 'counterValue' field " + _counterValueErr.Error())
            }
            _map["Struct"] = counterValue

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8TimePeriod": // Struct
            _map := map[string]interface{}{}

            // Simple Field (timePeriod)
            timePeriod, _timePeriodErr := io.ReadUint16(16)
            if _timePeriodErr != nil {
                return nil, errors.New("Error parsing 'timePeriod' field " + _timePeriodErr.Error())
            }
            _map["Struct"] = timePeriod

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8FlowRate": // Struct
            _map := map[string]interface{}{}

            // Simple Field (flowRate)
            flowRate, _flowRateErr := io.ReadUint16(16)
            if _flowRateErr != nil {
                return nil, errors.New("Error parsing 'flowRate' field " + _flowRateErr.Error())
            }
            _map["Struct"] = flowRate

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8Counter": // Struct
            _map := map[string]interface{}{}

            // Simple Field (counterValue)
            counterValue, _counterValueErr := io.ReadUint16(16)
            if _counterValueErr != nil {
                return nil, errors.New("Error parsing 'counterValue' field " + _counterValueErr.Error())
            }
            _map["Struct"] = counterValue

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8ElectricCurrent": // Struct
            _map := map[string]interface{}{}

            // Simple Field (electricalCurrent)
            electricalCurrent, _electricalCurrentErr := io.ReadUint16(16)
            if _electricalCurrentErr != nil {
                return nil, errors.New("Error parsing 'electricalCurrent' field " + _electricalCurrentErr.Error())
            }
            _map["Struct"] = electricalCurrent

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8Power": // Struct
            _map := map[string]interface{}{}

            // Simple Field (power)
            power, _powerErr := io.ReadUint16(16)
            if _powerErr != nil {
                return nil, errors.New("Error parsing 'power' field " + _powerErr.Error())
            }
            _map["Struct"] = power

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8AtmPressure": // Struct
            _map := map[string]interface{}{}

            // Simple Field (atmPressure)
            atmPressure, _atmPressureErr := io.ReadUint16(16)
            if _atmPressureErr != nil {
                return nil, errors.New("Error parsing 'atmPressure' field " + _atmPressureErr.Error())
            }
            _map["Struct"] = atmPressure

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8PercentValue": // Struct
            _map := map[string]interface{}{}

            // Simple Field (percentValue)
            percentValue, _percentValueErr := io.ReadUint16(16)
            if _percentValueErr != nil {
                return nil, errors.New("Error parsing 'percentValue' field " + _percentValueErr.Error())
            }
            _map["Struct"] = percentValue

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8HvacAirQuality": // Struct
            _map := map[string]interface{}{}

            // Simple Field (ppmResolution)
            ppmResolution, _ppmResolutionErr := io.ReadUint16(16)
            if _ppmResolutionErr != nil {
                return nil, errors.New("Error parsing 'ppmResolution' field " + _ppmResolutionErr.Error())
            }
            _map["Struct"] = ppmResolution

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8WindSpeed": // Struct
            _map := map[string]interface{}{}

            // Simple Field (windSpeed)
            windSpeed, _windSpeedErr := io.ReadUint16(16)
            if _windSpeedErr != nil {
                return nil, errors.New("Error parsing 'windSpeed' field " + _windSpeedErr.Error())
            }
            _map["Struct"] = windSpeed

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8SunIntensity": // Struct
            _map := map[string]interface{}{}

            // Simple Field (sunIntensity)
            sunIntensity, _sunIntensityErr := io.ReadUint16(16)
            if _sunIntensityErr != nil {
                return nil, errors.New("Error parsing 'sunIntensity' field " + _sunIntensityErr.Error())
            }
            _map["Struct"] = sunIntensity

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16Z8HvacAirFlow": // Struct
            _map := map[string]interface{}{}

            // Simple Field (airFlow)
            airFlow, _airFlowErr := io.ReadUint16(16)
            if _airFlowErr != nil {
                return nil, errors.New("Error parsing 'airFlow' field " + _airFlowErr.Error())
            }
            _map["Struct"] = airFlow

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "V8Z8RelSignedValue": // Struct
            _map := map[string]interface{}{}

            // Simple Field (relSignedValue)
            relSignedValue, _relSignedValueErr := io.ReadInt8(8)
            if _relSignedValueErr != nil {
                return nil, errors.New("Error parsing 'relSignedValue' field " + _relSignedValueErr.Error())
            }
            _map["Struct"] = relSignedValue

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "V16Z8DeltaTime": // Struct
            _map := map[string]interface{}{}

            // Simple Field (deltaTime)
            deltaTime, _deltaTimeErr := io.ReadInt16(16)
            if _deltaTimeErr != nil {
                return nil, errors.New("Error parsing 'deltaTime' field " + _deltaTimeErr.Error())
            }
            _map["Struct"] = deltaTime

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "V16Z8RelSignedValue": // Struct
            _map := map[string]interface{}{}

            // Simple Field (relSignedValue)
            relSignedValue, _relSignedValueErr := io.ReadInt16(16)
            if _relSignedValueErr != nil {
                return nil, errors.New("Error parsing 'relSignedValue' field " + _relSignedValueErr.Error())
            }
            _map["Struct"] = relSignedValue

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16N8HvacModeAndTimeDelay": // Struct
            _map := map[string]interface{}{}

            // Simple Field (delayTime)
            delayTime, _delayTimeErr := io.ReadUint16(16)
            if _delayTimeErr != nil {
                return nil, errors.New("Error parsing 'delayTime' field " + _delayTimeErr.Error())
            }
            _map["Struct"] = delayTime

            // Simple Field (hvacMode)
            hvacMode, _hvacModeErr := io.ReadUint8(8)
            if _hvacModeErr != nil {
                return nil, errors.New("Error parsing 'hvacMode' field " + _hvacModeErr.Error())
            }
            _map["Struct"] = hvacMode
            return values.NewPlcStruct(_map), nil
        case formatName == "U16N8DhwModeAndTimeDelay": // Struct
            _map := map[string]interface{}{}

            // Simple Field (delayTime)
            delayTime, _delayTimeErr := io.ReadUint16(16)
            if _delayTimeErr != nil {
                return nil, errors.New("Error parsing 'delayTime' field " + _delayTimeErr.Error())
            }
            _map["Struct"] = delayTime

            // Simple Field (dhwMode)
            dhwMode, _dhwModeErr := io.ReadUint8(8)
            if _dhwModeErr != nil {
                return nil, errors.New("Error parsing 'dhwMode' field " + _dhwModeErr.Error())
            }
            _map["Struct"] = dhwMode
            return values.NewPlcStruct(_map), nil
        case formatName == "U16N8OccupancyModeAndTimeDelay": // Struct
            _map := map[string]interface{}{}

            // Simple Field (delayTime)
            delayTime, _delayTimeErr := io.ReadUint16(16)
            if _delayTimeErr != nil {
                return nil, errors.New("Error parsing 'delayTime' field " + _delayTimeErr.Error())
            }
            _map["Struct"] = delayTime

            // Simple Field (occupationMode)
            occupationMode, _occupationModeErr := io.ReadUint8(8)
            if _occupationModeErr != nil {
                return nil, errors.New("Error parsing 'occupationMode' field " + _occupationModeErr.Error())
            }
            _map["Struct"] = occupationMode
            return values.NewPlcStruct(_map), nil
        case formatName == "U16N8BuildingModeAndTimeDelay": // Struct
            _map := map[string]interface{}{}

            // Simple Field (delayTime)
            delayTime, _delayTimeErr := io.ReadUint16(16)
            if _delayTimeErr != nil {
                return nil, errors.New("Error parsing 'delayTime' field " + _delayTimeErr.Error())
            }
            _map["Struct"] = delayTime

            // Simple Field (buildingMode)
            buildingMode, _buildingModeErr := io.ReadUint8(8)
            if _buildingModeErr != nil {
                return nil, errors.New("Error parsing 'buildingMode' field " + _buildingModeErr.Error())
            }
            _map["Struct"] = buildingMode
            return values.NewPlcStruct(_map), nil
        case formatName == "U8B8StatusBurnerController": // Struct
            _map := map[string]interface{}{}

            // Simple Field (actualRelativePower)
            actualRelativePower, _actualRelativePowerErr := io.ReadUint8(8)
            if _actualRelativePowerErr != nil {
                return nil, errors.New("Error parsing 'actualRelativePower' field " + _actualRelativePowerErr.Error())
            }
            _map["Struct"] = actualRelativePower

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (stage2Active)
            stage2Active, _stage2ActiveErr := io.ReadBit()
            if _stage2ActiveErr != nil {
                return nil, errors.New("Error parsing 'stage2Active' field " + _stage2ActiveErr.Error())
            }
            _map["Struct"] = stage2Active

            // Simple Field (stage1Active)
            stage1Active, _stage1ActiveErr := io.ReadBit()
            if _stage1ActiveErr != nil {
                return nil, errors.New("Error parsing 'stage1Active' field " + _stage1ActiveErr.Error())
            }
            _map["Struct"] = stage1Active

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure

            // Simple Field (actualRelativePowerValid)
            actualRelativePowerValid, _actualRelativePowerValidErr := io.ReadBit()
            if _actualRelativePowerValidErr != nil {
                return nil, errors.New("Error parsing 'actualRelativePowerValid' field " + _actualRelativePowerValidErr.Error())
            }
            _map["Struct"] = actualRelativePowerValid
            return values.NewPlcStruct(_map), nil
        case formatName == "U8B8LockingSignal": // Struct
            _map := map[string]interface{}{}

            // Simple Field (requestedPowerReduction)
            requestedPowerReduction, _requestedPowerReductionErr := io.ReadUint8(8)
            if _requestedPowerReductionErr != nil {
                return nil, errors.New("Error parsing 'requestedPowerReduction' field " + _requestedPowerReductionErr.Error())
            }
            _map["Struct"] = requestedPowerReduction

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (critical)
            critical, _criticalErr := io.ReadBit()
            if _criticalErr != nil {
                return nil, errors.New("Error parsing 'critical' field " + _criticalErr.Error())
            }
            _map["Struct"] = critical

            // Simple Field (requestedPowerReductionValid)
            requestedPowerReductionValid, _requestedPowerReductionValidErr := io.ReadBit()
            if _requestedPowerReductionValidErr != nil {
                return nil, errors.New("Error parsing 'requestedPowerReductionValid' field " + _requestedPowerReductionValidErr.Error())
            }
            _map["Struct"] = requestedPowerReductionValid
            return values.NewPlcStruct(_map), nil
        case formatName == "U8B8BoilerControllerDemandSignal": // Struct
            _map := map[string]interface{}{}

            // Simple Field (relativeDemand)
            relativeDemand, _relativeDemandErr := io.ReadUint8(8)
            if _relativeDemandErr != nil {
                return nil, errors.New("Error parsing 'relativeDemand' field " + _relativeDemandErr.Error())
            }
            _map["Struct"] = relativeDemand

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (controlsOperationStage2)
            controlsOperationStage2, _controlsOperationStage2Err := io.ReadBit()
            if _controlsOperationStage2Err != nil {
                return nil, errors.New("Error parsing 'controlsOperationStage2' field " + _controlsOperationStage2Err.Error())
            }
            _map["Struct"] = controlsOperationStage2

            // Simple Field (controlsOperationStage1)
            controlsOperationStage1, _controlsOperationStage1Err := io.ReadBit()
            if _controlsOperationStage1Err != nil {
                return nil, errors.New("Error parsing 'controlsOperationStage1' field " + _controlsOperationStage1Err.Error())
            }
            _map["Struct"] = controlsOperationStage1
            return values.NewPlcStruct(_map), nil
        case formatName == "U8B8ActuatorPositionDemand": // Struct
            _map := map[string]interface{}{}

            // Simple Field (actuatorPositionDemand)
            actuatorPositionDemand, _actuatorPositionDemandErr := io.ReadUint8(8)
            if _actuatorPositionDemandErr != nil {
                return nil, errors.New("Error parsing 'actuatorPositionDemand' field " + _actuatorPositionDemandErr.Error())
            }
            _map["Struct"] = actuatorPositionDemand

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (emergencyDemand)
            emergencyDemand, _emergencyDemandErr := io.ReadBit()
            if _emergencyDemandErr != nil {
                return nil, errors.New("Error parsing 'emergencyDemand' field " + _emergencyDemandErr.Error())
            }
            _map["Struct"] = emergencyDemand

            // Simple Field (shiftLoadPriority)
            shiftLoadPriority, _shiftLoadPriorityErr := io.ReadBit()
            if _shiftLoadPriorityErr != nil {
                return nil, errors.New("Error parsing 'shiftLoadPriority' field " + _shiftLoadPriorityErr.Error())
            }
            _map["Struct"] = shiftLoadPriority

            // Simple Field (absoluteLoadPriority)
            absoluteLoadPriority, _absoluteLoadPriorityErr := io.ReadBit()
            if _absoluteLoadPriorityErr != nil {
                return nil, errors.New("Error parsing 'absoluteLoadPriority' field " + _absoluteLoadPriorityErr.Error())
            }
            _map["Struct"] = absoluteLoadPriority

            // Simple Field (actuatorPositionDemandValid)
            actuatorPositionDemandValid, _actuatorPositionDemandValidErr := io.ReadBit()
            if _actuatorPositionDemandValidErr != nil {
                return nil, errors.New("Error parsing 'actuatorPositionDemandValid' field " + _actuatorPositionDemandValidErr.Error())
            }
            _map["Struct"] = actuatorPositionDemandValid
            return values.NewPlcStruct(_map), nil
        case formatName == "U8B8ActuatorPositionStatus": // Struct
            _map := map[string]interface{}{}

            // Simple Field (actualActuatorPosition)
            actualActuatorPosition, _actualActuatorPositionErr := io.ReadUint8(8)
            if _actualActuatorPositionErr != nil {
                return nil, errors.New("Error parsing 'actualActuatorPosition' field " + _actualActuatorPositionErr.Error())
            }
            _map["Struct"] = actualActuatorPosition

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (synchronizationMode)
            synchronizationMode, _synchronizationModeErr := io.ReadBit()
            if _synchronizationModeErr != nil {
                return nil, errors.New("Error parsing 'synchronizationMode' field " + _synchronizationModeErr.Error())
            }
            _map["Struct"] = synchronizationMode

            // Simple Field (valveKick)
            valveKick, _valveKickErr := io.ReadBit()
            if _valveKickErr != nil {
                return nil, errors.New("Error parsing 'valveKick' field " + _valveKickErr.Error())
            }
            _map["Struct"] = valveKick

            // Simple Field (callibrationMode)
            callibrationMode, _callibrationModeErr := io.ReadBit()
            if _callibrationModeErr != nil {
                return nil, errors.New("Error parsing 'callibrationMode' field " + _callibrationModeErr.Error())
            }
            _map["Struct"] = callibrationMode

            // Simple Field (positionManuallyOverridden)
            positionManuallyOverridden, _positionManuallyOverriddenErr := io.ReadBit()
            if _positionManuallyOverriddenErr != nil {
                return nil, errors.New("Error parsing 'positionManuallyOverridden' field " + _positionManuallyOverriddenErr.Error())
            }
            _map["Struct"] = positionManuallyOverridden

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure
            return values.NewPlcStruct(_map), nil
        case formatName == "U8B8StatusLightingActuator": // Struct
            _map := map[string]interface{}{}

            // Simple Field (lightingLevel)
            lightingLevel, _lightingLevelErr := io.ReadUint8(8)
            if _lightingLevelErr != nil {
                return nil, errors.New("Error parsing 'lightingLevel' field " + _lightingLevelErr.Error())
            }
            _map["Struct"] = lightingLevel

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure

            // Simple Field (localOverride)
            localOverride, _localOverrideErr := io.ReadBit()
            if _localOverrideErr != nil {
                return nil, errors.New("Error parsing 'localOverride' field " + _localOverrideErr.Error())
            }
            _map["Struct"] = localOverride

            // Simple Field (dimming)
            dimming, _dimmingErr := io.ReadBit()
            if _dimmingErr != nil {
                return nil, errors.New("Error parsing 'dimming' field " + _dimmingErr.Error())
            }
            _map["Struct"] = dimming

            // Simple Field (staircaseLightingFunction)
            staircaseLightingFunction, _staircaseLightingFunctionErr := io.ReadBit()
            if _staircaseLightingFunctionErr != nil {
                return nil, errors.New("Error parsing 'staircaseLightingFunction' field " + _staircaseLightingFunctionErr.Error())
            }
            _map["Struct"] = staircaseLightingFunction

            // Simple Field (nightMode)
            nightMode, _nightModeErr := io.ReadBit()
            if _nightModeErr != nil {
                return nil, errors.New("Error parsing 'nightMode' field " + _nightModeErr.Error())
            }
            _map["Struct"] = nightMode

            // Simple Field (forced)
            forced, _forcedErr := io.ReadBit()
            if _forcedErr != nil {
                return nil, errors.New("Error parsing 'forced' field " + _forcedErr.Error())
            }
            _map["Struct"] = forced

            // Simple Field (locked)
            locked, _lockedErr := io.ReadBit()
            if _lockedErr != nil {
                return nil, errors.New("Error parsing 'locked' field " + _lockedErr.Error())
            }
            _map["Struct"] = locked

            // Simple Field (lightingLevelValid)
            lightingLevelValid, _lightingLevelValidErr := io.ReadBit()
            if _lightingLevelValidErr != nil {
                return nil, errors.New("Error parsing 'lightingLevelValid' field " + _lightingLevelValidErr.Error())
            }
            _map["Struct"] = lightingLevelValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16B8HeatProducerManagerStatus": // Struct
            _map := map[string]interface{}{}

            // Simple Field (tempFlowProdSegmH)
            tempFlowProdSegmH, _tempFlowProdSegmHErr := io.ReadFloat32(16)
            if _tempFlowProdSegmHErr != nil {
                return nil, errors.New("Error parsing 'tempFlowProdSegmH' field " + _tempFlowProdSegmHErr.Error())
            }
            _map["Struct"] = tempFlowProdSegmH

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (temporarilyOff)
            temporarilyOff, _temporarilyOffErr := io.ReadBit()
            if _temporarilyOffErr != nil {
                return nil, errors.New("Error parsing 'temporarilyOff' field " + _temporarilyOffErr.Error())
            }
            _map["Struct"] = temporarilyOff

            // Simple Field (permanentlyOff)
            permanentlyOff, _permanentlyOffErr := io.ReadBit()
            if _permanentlyOffErr != nil {
                return nil, errors.New("Error parsing 'permanentlyOff' field " + _permanentlyOffErr.Error())
            }
            _map["Struct"] = permanentlyOff

            // Simple Field (switchedOffSummerMode)
            switchedOffSummerMode, _switchedOffSummerModeErr := io.ReadBit()
            if _switchedOffSummerModeErr != nil {
                return nil, errors.New("Error parsing 'switchedOffSummerMode' field " + _switchedOffSummerModeErr.Error())
            }
            _map["Struct"] = switchedOffSummerMode

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure

            // Simple Field (tempFlowProdSegmHValid)
            tempFlowProdSegmHValid, _tempFlowProdSegmHValidErr := io.ReadBit()
            if _tempFlowProdSegmHValidErr != nil {
                return nil, errors.New("Error parsing 'tempFlowProdSegmHValid' field " + _tempFlowProdSegmHValidErr.Error())
            }
            _map["Struct"] = tempFlowProdSegmHValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16B8RoomTemperatureDemand": // Struct
            _map := map[string]interface{}{}

            // Simple Field (roomTemperatureDemand)
            roomTemperatureDemand, _roomTemperatureDemandErr := io.ReadFloat32(16)
            if _roomTemperatureDemandErr != nil {
                return nil, errors.New("Error parsing 'roomTemperatureDemand' field " + _roomTemperatureDemandErr.Error())
            }
            _map["Struct"] = roomTemperatureDemand

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (emergencyDemand)
            emergencyDemand, _emergencyDemandErr := io.ReadBit()
            if _emergencyDemandErr != nil {
                return nil, errors.New("Error parsing 'emergencyDemand' field " + _emergencyDemandErr.Error())
            }
            _map["Struct"] = emergencyDemand

            // Simple Field (shiftLoadPriority)
            shiftLoadPriority, _shiftLoadPriorityErr := io.ReadBit()
            if _shiftLoadPriorityErr != nil {
                return nil, errors.New("Error parsing 'shiftLoadPriority' field " + _shiftLoadPriorityErr.Error())
            }
            _map["Struct"] = shiftLoadPriority

            // Simple Field (absoluteLoadPriority)
            absoluteLoadPriority, _absoluteLoadPriorityErr := io.ReadBit()
            if _absoluteLoadPriorityErr != nil {
                return nil, errors.New("Error parsing 'absoluteLoadPriority' field " + _absoluteLoadPriorityErr.Error())
            }
            _map["Struct"] = absoluteLoadPriority

            // Simple Field (roomTemperatureDemandValid)
            roomTemperatureDemandValid, _roomTemperatureDemandValidErr := io.ReadBit()
            if _roomTemperatureDemandValidErr != nil {
                return nil, errors.New("Error parsing 'roomTemperatureDemandValid' field " + _roomTemperatureDemandValidErr.Error())
            }
            _map["Struct"] = roomTemperatureDemandValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16B8ColdWaterProducerManagerStatus": // Struct
            _map := map[string]interface{}{}

            // Simple Field (flowTemperatureProdSegmC)
            flowTemperatureProdSegmC, _flowTemperatureProdSegmCErr := io.ReadFloat32(16)
            if _flowTemperatureProdSegmCErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureProdSegmC' field " + _flowTemperatureProdSegmCErr.Error())
            }
            _map["Struct"] = flowTemperatureProdSegmC

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (temporarilyOff)
            temporarilyOff, _temporarilyOffErr := io.ReadBit()
            if _temporarilyOffErr != nil {
                return nil, errors.New("Error parsing 'temporarilyOff' field " + _temporarilyOffErr.Error())
            }
            _map["Struct"] = temporarilyOff

            // Simple Field (permanentlyOff)
            permanentlyOff, _permanentlyOffErr := io.ReadBit()
            if _permanentlyOffErr != nil {
                return nil, errors.New("Error parsing 'permanentlyOff' field " + _permanentlyOffErr.Error())
            }
            _map["Struct"] = permanentlyOff

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure

            // Simple Field (flowTemperatureProdSegmCValid)
            flowTemperatureProdSegmCValid, _flowTemperatureProdSegmCValidErr := io.ReadBit()
            if _flowTemperatureProdSegmCValidErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureProdSegmCValid' field " + _flowTemperatureProdSegmCValidErr.Error())
            }
            _map["Struct"] = flowTemperatureProdSegmCValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16B8WaterTemperatureControllerStatus": // Struct
            _map := map[string]interface{}{}

            // Simple Field (actualTemperature)
            actualTemperature, _actualTemperatureErr := io.ReadFloat32(16)
            if _actualTemperatureErr != nil {
                return nil, errors.New("Error parsing 'actualTemperature' field " + _actualTemperatureErr.Error())
            }
            _map["Struct"] = actualTemperature

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (controllerWorking)
            controllerWorking, _controllerWorkingErr := io.ReadBit()
            if _controllerWorkingErr != nil {
                return nil, errors.New("Error parsing 'controllerWorking' field " + _controllerWorkingErr.Error())
            }
            _map["Struct"] = controllerWorking

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure

            // Simple Field (actualTemperatureValid)
            actualTemperatureValid, _actualTemperatureValidErr := io.ReadBit()
            if _actualTemperatureValidErr != nil {
                return nil, errors.New("Error parsing 'actualTemperatureValid' field " + _actualTemperatureValidErr.Error())
            }
            _map["Struct"] = actualTemperatureValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16B16": // Struct
            _map := map[string]interface{}{}

            // Simple Field (flowTemperatureDemand)
            flowTemperatureDemand, _flowTemperatureDemandErr := io.ReadFloat32(16)
            if _flowTemperatureDemandErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureDemand' field " + _flowTemperatureDemandErr.Error())
            }
            _map["Struct"] = flowTemperatureDemand

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (demandFromDhwWhileLegionellaFunctionIsActive)
            demandFromDhwWhileLegionellaFunctionIsActive, _demandFromDhwWhileLegionellaFunctionIsActiveErr := io.ReadBit()
            if _demandFromDhwWhileLegionellaFunctionIsActiveErr != nil {
                return nil, errors.New("Error parsing 'demandFromDhwWhileLegionellaFunctionIsActive' field " + _demandFromDhwWhileLegionellaFunctionIsActiveErr.Error())
            }
            _map["Struct"] = demandFromDhwWhileLegionellaFunctionIsActive

            // Simple Field (emergencyDemandForFrostProtection)
            emergencyDemandForFrostProtection, _emergencyDemandForFrostProtectionErr := io.ReadBit()
            if _emergencyDemandForFrostProtectionErr != nil {
                return nil, errors.New("Error parsing 'emergencyDemandForFrostProtection' field " + _emergencyDemandForFrostProtectionErr.Error())
            }
            _map["Struct"] = emergencyDemandForFrostProtection

            // Simple Field (requestForWaterCirculationInPrimaryDistributionSegment)
            requestForWaterCirculationInPrimaryDistributionSegment, _requestForWaterCirculationInPrimaryDistributionSegmentErr := io.ReadBit()
            if _requestForWaterCirculationInPrimaryDistributionSegmentErr != nil {
                return nil, errors.New("Error parsing 'requestForWaterCirculationInPrimaryDistributionSegment' field " + _requestForWaterCirculationInPrimaryDistributionSegmentErr.Error())
            }
            _map["Struct"] = requestForWaterCirculationInPrimaryDistributionSegment

            // Simple Field (demandFromAuxillaryHeatOrCoolConsumer)
            demandFromAuxillaryHeatOrCoolConsumer, _demandFromAuxillaryHeatOrCoolConsumerErr := io.ReadBit()
            if _demandFromAuxillaryHeatOrCoolConsumerErr != nil {
                return nil, errors.New("Error parsing 'demandFromAuxillaryHeatOrCoolConsumer' field " + _demandFromAuxillaryHeatOrCoolConsumerErr.Error())
            }
            _map["Struct"] = demandFromAuxillaryHeatOrCoolConsumer

            // Simple Field (demandFromVentilation)
            demandFromVentilation, _demandFromVentilationErr := io.ReadBit()
            if _demandFromVentilationErr != nil {
                return nil, errors.New("Error parsing 'demandFromVentilation' field " + _demandFromVentilationErr.Error())
            }
            _map["Struct"] = demandFromVentilation

            // Simple Field (demandForRoomHeatingOrCooling)
            demandForRoomHeatingOrCooling, _demandForRoomHeatingOrCoolingErr := io.ReadBit()
            if _demandForRoomHeatingOrCoolingErr != nil {
                return nil, errors.New("Error parsing 'demandForRoomHeatingOrCooling' field " + _demandForRoomHeatingOrCoolingErr.Error())
            }
            _map["Struct"] = demandForRoomHeatingOrCooling

            // Simple Field (heatDemandFromDhw)
            heatDemandFromDhw, _heatDemandFromDhwErr := io.ReadBit()
            if _heatDemandFromDhwErr != nil {
                return nil, errors.New("Error parsing 'heatDemandFromDhw' field " + _heatDemandFromDhwErr.Error())
            }
            _map["Struct"] = heatDemandFromDhw

            // Simple Field (flowTemperatureDemandIsMin)
            flowTemperatureDemandIsMin, _flowTemperatureDemandIsMinErr := io.ReadBit()
            if _flowTemperatureDemandIsMinErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureDemandIsMin' field " + _flowTemperatureDemandIsMinErr.Error())
            }
            _map["Struct"] = flowTemperatureDemandIsMin

            // Simple Field (flowTemperatureDemandIsMax)
            flowTemperatureDemandIsMax, _flowTemperatureDemandIsMaxErr := io.ReadBit()
            if _flowTemperatureDemandIsMaxErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureDemandIsMax' field " + _flowTemperatureDemandIsMaxErr.Error())
            }
            _map["Struct"] = flowTemperatureDemandIsMax

            // Simple Field (shiftLoadPriority)
            shiftLoadPriority, _shiftLoadPriorityErr := io.ReadBit()
            if _shiftLoadPriorityErr != nil {
                return nil, errors.New("Error parsing 'shiftLoadPriority' field " + _shiftLoadPriorityErr.Error())
            }
            _map["Struct"] = shiftLoadPriority

            // Simple Field (absoluteLoadPriority)
            absoluteLoadPriority, _absoluteLoadPriorityErr := io.ReadBit()
            if _absoluteLoadPriorityErr != nil {
                return nil, errors.New("Error parsing 'absoluteLoadPriority' field " + _absoluteLoadPriorityErr.Error())
            }
            _map["Struct"] = absoluteLoadPriority

            // Simple Field (flowTemperatureDemandValid)
            flowTemperatureDemandValid, _flowTemperatureDemandValidErr := io.ReadBit()
            if _flowTemperatureDemandValidErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureDemandValid' field " + _flowTemperatureDemandValidErr.Error())
            }
            _map["Struct"] = flowTemperatureDemandValid
            return values.NewPlcStruct(_map), nil
        case formatName == "U8N8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (energyDemand)
            energyDemand, _energyDemandErr := io.ReadUint8(8)
            if _energyDemandErr != nil {
                return nil, errors.New("Error parsing 'energyDemand' field " + _energyDemandErr.Error())
            }
            _map["Struct"] = energyDemand

            // Simple Field (actualControllerMode)
            actualControllerMode, _actualControllerModeErr := io.ReadUint8(8)
            if _actualControllerModeErr != nil {
                return nil, errors.New("Error parsing 'actualControllerMode' field " + _actualControllerModeErr.Error())
            }
            _map["Struct"] = actualControllerMode
            return values.NewPlcStruct(_map), nil
        case formatName == "V16V16V16RoomTemperature": // Struct
            _map := map[string]interface{}{}

            // Simple Field (temperatureSetpointComfort)
            temperatureSetpointComfort, _temperatureSetpointComfortErr := io.ReadFloat32(16)
            if _temperatureSetpointComfortErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointComfort' field " + _temperatureSetpointComfortErr.Error())
            }
            _map["Struct"] = temperatureSetpointComfort

            // Simple Field (temperatureSetpointStandby)
            temperatureSetpointStandby, _temperatureSetpointStandbyErr := io.ReadFloat32(16)
            if _temperatureSetpointStandbyErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointStandby' field " + _temperatureSetpointStandbyErr.Error())
            }
            _map["Struct"] = temperatureSetpointStandby

            // Simple Field (temperatureSetpointEco)
            temperatureSetpointEco, _temperatureSetpointEcoErr := io.ReadFloat32(16)
            if _temperatureSetpointEcoErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointEco' field " + _temperatureSetpointEcoErr.Error())
            }
            _map["Struct"] = temperatureSetpointEco
            return values.NewPlcStruct(_map), nil
        case formatName == "V16V16V16RoomTemperatureShift": // Struct
            _map := map[string]interface{}{}

            // Simple Field (temperatureSetpointShiftComfort)
            temperatureSetpointShiftComfort, _temperatureSetpointShiftComfortErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftComfortErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftComfort' field " + _temperatureSetpointShiftComfortErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftComfort

            // Simple Field (temperatureSetpointShiftStandby)
            temperatureSetpointShiftStandby, _temperatureSetpointShiftStandbyErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftStandbyErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftStandby' field " + _temperatureSetpointShiftStandbyErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftStandby

            // Simple Field (temperatureSetpointShiftEco)
            temperatureSetpointShiftEco, _temperatureSetpointShiftEcoErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftEcoErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftEco' field " + _temperatureSetpointShiftEcoErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftEco
            return values.NewPlcStruct(_map), nil
        case formatName == "V16V16V16V16RoomTemperature": // Struct
            _map := map[string]interface{}{}

            // Simple Field (temperatureSetpointComfort)
            temperatureSetpointComfort, _temperatureSetpointComfortErr := io.ReadFloat32(16)
            if _temperatureSetpointComfortErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointComfort' field " + _temperatureSetpointComfortErr.Error())
            }
            _map["Struct"] = temperatureSetpointComfort

            // Simple Field (temperatureSetpointStandby)
            temperatureSetpointStandby, _temperatureSetpointStandbyErr := io.ReadFloat32(16)
            if _temperatureSetpointStandbyErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointStandby' field " + _temperatureSetpointStandbyErr.Error())
            }
            _map["Struct"] = temperatureSetpointStandby

            // Simple Field (temperatureSetpointEco)
            temperatureSetpointEco, _temperatureSetpointEcoErr := io.ReadFloat32(16)
            if _temperatureSetpointEcoErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointEco' field " + _temperatureSetpointEcoErr.Error())
            }
            _map["Struct"] = temperatureSetpointEco

            // Simple Field (temperatureSetpointBProt)
            temperatureSetpointBProt, _temperatureSetpointBProtErr := io.ReadFloat32(16)
            if _temperatureSetpointBProtErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointBProt' field " + _temperatureSetpointBProtErr.Error())
            }
            _map["Struct"] = temperatureSetpointBProt
            return values.NewPlcStruct(_map), nil
        case formatName == "V16V16V16V16DhwtTemperature": // Struct
            _map := map[string]interface{}{}

            // Simple Field (temperatureSetpointLegioProtect)
            temperatureSetpointLegioProtect, _temperatureSetpointLegioProtectErr := io.ReadFloat32(16)
            if _temperatureSetpointLegioProtectErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointLegioProtect' field " + _temperatureSetpointLegioProtectErr.Error())
            }
            _map["Struct"] = temperatureSetpointLegioProtect

            // Simple Field (temperatureSetpointNormal)
            temperatureSetpointNormal, _temperatureSetpointNormalErr := io.ReadFloat32(16)
            if _temperatureSetpointNormalErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointNormal' field " + _temperatureSetpointNormalErr.Error())
            }
            _map["Struct"] = temperatureSetpointNormal

            // Simple Field (temperatureSetpointReduced)
            temperatureSetpointReduced, _temperatureSetpointReducedErr := io.ReadFloat32(16)
            if _temperatureSetpointReducedErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointReduced' field " + _temperatureSetpointReducedErr.Error())
            }
            _map["Struct"] = temperatureSetpointReduced

            // Simple Field (temperatureSetpointFrostProtect)
            temperatureSetpointFrostProtect, _temperatureSetpointFrostProtectErr := io.ReadFloat32(16)
            if _temperatureSetpointFrostProtectErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointFrostProtect' field " + _temperatureSetpointFrostProtectErr.Error())
            }
            _map["Struct"] = temperatureSetpointFrostProtect
            return values.NewPlcStruct(_map), nil
        case formatName == "V16V16V16V16RoomTemperatureShift": // Struct
            _map := map[string]interface{}{}

            // Simple Field (temperatureSetpointShiftComfort)
            temperatureSetpointShiftComfort, _temperatureSetpointShiftComfortErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftComfortErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftComfort' field " + _temperatureSetpointShiftComfortErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftComfort

            // Simple Field (temperatureSetpointShiftStandby)
            temperatureSetpointShiftStandby, _temperatureSetpointShiftStandbyErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftStandbyErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftStandby' field " + _temperatureSetpointShiftStandbyErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftStandby

            // Simple Field (temperatureSetpointShiftEco)
            temperatureSetpointShiftEco, _temperatureSetpointShiftEcoErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftEcoErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftEco' field " + _temperatureSetpointShiftEcoErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftEco

            // Simple Field (temperatureSetpointShiftBProt)
            temperatureSetpointShiftBProt, _temperatureSetpointShiftBProtErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftBProtErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftBProt' field " + _temperatureSetpointShiftBProtErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftBProt
            return values.NewPlcStruct(_map), nil
        case formatName == "V16U8B8Heat": // Struct
            _map := map[string]interface{}{}

            // Simple Field (flowTemperatureDemand)
            flowTemperatureDemand, _flowTemperatureDemandErr := io.ReadInt16(16)
            if _flowTemperatureDemandErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureDemand' field " + _flowTemperatureDemandErr.Error())
            }
            _map["Struct"] = flowTemperatureDemand

            // Simple Field (relativePower)
            relativePower, _relativePowerErr := io.ReadUint8(8)
            if _relativePowerErr != nil {
                return nil, errors.New("Error parsing 'relativePower' field " + _relativePowerErr.Error())
            }
            _map["Struct"] = relativePower

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (boilerEnabled)
            boilerEnabled, _boilerEnabledErr := io.ReadBit()
            if _boilerEnabledErr != nil {
                return nil, errors.New("Error parsing 'boilerEnabled' field " + _boilerEnabledErr.Error())
            }
            _map["Struct"] = boilerEnabled

            // Simple Field (stage2Forced)
            stage2Forced, _stage2ForcedErr := io.ReadBit()
            if _stage2ForcedErr != nil {
                return nil, errors.New("Error parsing 'stage2Forced' field " + _stage2ForcedErr.Error())
            }
            _map["Struct"] = stage2Forced

            // Simple Field (stage2Enabled)
            stage2Enabled, _stage2EnabledErr := io.ReadBit()
            if _stage2EnabledErr != nil {
                return nil, errors.New("Error parsing 'stage2Enabled' field " + _stage2EnabledErr.Error())
            }
            _map["Struct"] = stage2Enabled

            // Simple Field (stage1Forced)
            stage1Forced, _stage1ForcedErr := io.ReadBit()
            if _stage1ForcedErr != nil {
                return nil, errors.New("Error parsing 'stage1Forced' field " + _stage1ForcedErr.Error())
            }
            _map["Struct"] = stage1Forced

            // Simple Field (stage1Enabled)
            stage1Enabled, _stage1EnabledErr := io.ReadBit()
            if _stage1EnabledErr != nil {
                return nil, errors.New("Error parsing 'stage1Enabled' field " + _stage1EnabledErr.Error())
            }
            _map["Struct"] = stage1Enabled

            // Simple Field (flowTemperatureDemandValid)
            flowTemperatureDemandValid, _flowTemperatureDemandValidErr := io.ReadBit()
            if _flowTemperatureDemandValidErr != nil {
                return nil, errors.New("Error parsing 'flowTemperatureDemandValid' field " + _flowTemperatureDemandValidErr.Error())
            }
            _map["Struct"] = flowTemperatureDemandValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16U8B8ChilledWater": // Struct
            _map := map[string]interface{}{}

            // Simple Field (chilledWaterFlowTemperatureDemand)
            chilledWaterFlowTemperatureDemand, _chilledWaterFlowTemperatureDemandErr := io.ReadInt16(16)
            if _chilledWaterFlowTemperatureDemandErr != nil {
                return nil, errors.New("Error parsing 'chilledWaterFlowTemperatureDemand' field " + _chilledWaterFlowTemperatureDemandErr.Error())
            }
            _map["Struct"] = chilledWaterFlowTemperatureDemand

            // Simple Field (relativePower)
            relativePower, _relativePowerErr := io.ReadUint8(8)
            if _relativePowerErr != nil {
                return nil, errors.New("Error parsing 'relativePower' field " + _relativePowerErr.Error())
            }
            _map["Struct"] = relativePower

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (chilledWaterPumpEnabled)
            chilledWaterPumpEnabled, _chilledWaterPumpEnabledErr := io.ReadBit()
            if _chilledWaterPumpEnabledErr != nil {
                return nil, errors.New("Error parsing 'chilledWaterPumpEnabled' field " + _chilledWaterPumpEnabledErr.Error())
            }
            _map["Struct"] = chilledWaterPumpEnabled

            // Simple Field (relativePowerValid)
            relativePowerValid, _relativePowerValidErr := io.ReadBit()
            if _relativePowerValidErr != nil {
                return nil, errors.New("Error parsing 'relativePowerValid' field " + _relativePowerValidErr.Error())
            }
            _map["Struct"] = relativePowerValid

            // Simple Field (chilledWaterFlowTemperatureDemandValid)
            chilledWaterFlowTemperatureDemandValid, _chilledWaterFlowTemperatureDemandValidErr := io.ReadBit()
            if _chilledWaterFlowTemperatureDemandValidErr != nil {
                return nil, errors.New("Error parsing 'chilledWaterFlowTemperatureDemandValid' field " + _chilledWaterFlowTemperatureDemandValidErr.Error())
            }
            _map["Struct"] = chilledWaterFlowTemperatureDemandValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16U8B16Boiler": // Struct
            _map := map[string]interface{}{}

            // Simple Field (tempBoiler)
            tempBoiler, _tempBoilerErr := io.ReadInt16(16)
            if _tempBoilerErr != nil {
                return nil, errors.New("Error parsing 'tempBoiler' field " + _tempBoilerErr.Error())
            }
            _map["Struct"] = tempBoiler

            // Simple Field (relativePower)
            relativePower, _relativePowerErr := io.ReadUint8(8)
            if _relativePowerErr != nil {
                return nil, errors.New("Error parsing 'relativePower' field " + _relativePowerErr.Error())
            }
            _map["Struct"] = relativePower

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (chimneySweepFunctionActive)
            chimneySweepFunctionActive, _chimneySweepFunctionActiveErr := io.ReadBit()
            if _chimneySweepFunctionActiveErr != nil {
                return nil, errors.New("Error parsing 'chimneySweepFunctionActive' field " + _chimneySweepFunctionActiveErr.Error())
            }
            _map["Struct"] = chimneySweepFunctionActive

            // Simple Field (reducedAvailability)
            reducedAvailability, _reducedAvailabilityErr := io.ReadBit()
            if _reducedAvailabilityErr != nil {
                return nil, errors.New("Error parsing 'reducedAvailability' field " + _reducedAvailabilityErr.Error())
            }
            _map["Struct"] = reducedAvailability

            // Simple Field (powerLimitBoilerReached)
            powerLimitBoilerReached, _powerLimitBoilerReachedErr := io.ReadBit()
            if _powerLimitBoilerReachedErr != nil {
                return nil, errors.New("Error parsing 'powerLimitBoilerReached' field " + _powerLimitBoilerReachedErr.Error())
            }
            _map["Struct"] = powerLimitBoilerReached

            // Simple Field (powerLimitStage1Reached)
            powerLimitStage1Reached, _powerLimitStage1ReachedErr := io.ReadBit()
            if _powerLimitStage1ReachedErr != nil {
                return nil, errors.New("Error parsing 'powerLimitStage1Reached' field " + _powerLimitStage1ReachedErr.Error())
            }
            _map["Struct"] = powerLimitStage1Reached

            // Simple Field (stage2Enabled)
            stage2Enabled, _stage2EnabledErr := io.ReadBit()
            if _stage2EnabledErr != nil {
                return nil, errors.New("Error parsing 'stage2Enabled' field " + _stage2EnabledErr.Error())
            }
            _map["Struct"] = stage2Enabled

            // Simple Field (stage1Enabled)
            stage1Enabled, _stage1EnabledErr := io.ReadBit()
            if _stage1EnabledErr != nil {
                return nil, errors.New("Error parsing 'stage1Enabled' field " + _stage1EnabledErr.Error())
            }
            _map["Struct"] = stage1Enabled

            // Simple Field (boilerTemporarilyNotProvidingHeat)
            boilerTemporarilyNotProvidingHeat, _boilerTemporarilyNotProvidingHeatErr := io.ReadBit()
            if _boilerTemporarilyNotProvidingHeatErr != nil {
                return nil, errors.New("Error parsing 'boilerTemporarilyNotProvidingHeat' field " + _boilerTemporarilyNotProvidingHeatErr.Error())
            }
            _map["Struct"] = boilerTemporarilyNotProvidingHeat

            // Simple Field (permanentlyOff)
            permanentlyOff, _permanentlyOffErr := io.ReadBit()
            if _permanentlyOffErr != nil {
                return nil, errors.New("Error parsing 'permanentlyOff' field " + _permanentlyOffErr.Error())
            }
            _map["Struct"] = permanentlyOff

            // Simple Field (boilerSwitchedOffWinterSummerMode)
            boilerSwitchedOffWinterSummerMode, _boilerSwitchedOffWinterSummerModeErr := io.ReadBit()
            if _boilerSwitchedOffWinterSummerModeErr != nil {
                return nil, errors.New("Error parsing 'boilerSwitchedOffWinterSummerMode' field " + _boilerSwitchedOffWinterSummerModeErr.Error())
            }
            _map["Struct"] = boilerSwitchedOffWinterSummerMode

            // Simple Field (boilerFailure)
            boilerFailure, _boilerFailureErr := io.ReadBit()
            if _boilerFailureErr != nil {
                return nil, errors.New("Error parsing 'boilerFailure' field " + _boilerFailureErr.Error())
            }
            _map["Struct"] = boilerFailure

            // Simple Field (relativePowerValid)
            relativePowerValid, _relativePowerValidErr := io.ReadBit()
            if _relativePowerValidErr != nil {
                return nil, errors.New("Error parsing 'relativePowerValid' field " + _relativePowerValidErr.Error())
            }
            _map["Struct"] = relativePowerValid

            // Simple Field (tempBoilerValid)
            tempBoilerValid, _tempBoilerValidErr := io.ReadBit()
            if _tempBoilerValidErr != nil {
                return nil, errors.New("Error parsing 'tempBoilerValid' field " + _tempBoilerValidErr.Error())
            }
            _map["Struct"] = tempBoilerValid
            return values.NewPlcStruct(_map), nil
        case formatName == "V16U8B16Chiller": // Struct
            _map := map[string]interface{}{}

            // Simple Field (tempChiller)
            tempChiller, _tempChillerErr := io.ReadInt16(16)
            if _tempChillerErr != nil {
                return nil, errors.New("Error parsing 'tempChiller' field " + _tempChillerErr.Error())
            }
            _map["Struct"] = tempChiller

            // Simple Field (relativePower)
            relativePower, _relativePowerErr := io.ReadUint8(8)
            if _relativePowerErr != nil {
                return nil, errors.New("Error parsing 'relativePower' field " + _relativePowerErr.Error())
            }
            _map["Struct"] = relativePower

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(8); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (reducedAvailability)
            reducedAvailability, _reducedAvailabilityErr := io.ReadBit()
            if _reducedAvailabilityErr != nil {
                return nil, errors.New("Error parsing 'reducedAvailability' field " + _reducedAvailabilityErr.Error())
            }
            _map["Struct"] = reducedAvailability

            // Simple Field (powerLimitChillerReached)
            powerLimitChillerReached, _powerLimitChillerReachedErr := io.ReadBit()
            if _powerLimitChillerReachedErr != nil {
                return nil, errors.New("Error parsing 'powerLimitChillerReached' field " + _powerLimitChillerReachedErr.Error())
            }
            _map["Struct"] = powerLimitChillerReached

            // Simple Field (powerLimitCurrentStageReached)
            powerLimitCurrentStageReached, _powerLimitCurrentStageReachedErr := io.ReadBit()
            if _powerLimitCurrentStageReachedErr != nil {
                return nil, errors.New("Error parsing 'powerLimitCurrentStageReached' field " + _powerLimitCurrentStageReachedErr.Error())
            }
            _map["Struct"] = powerLimitCurrentStageReached

            // Simple Field (permanentlyOff)
            permanentlyOff, _permanentlyOffErr := io.ReadBit()
            if _permanentlyOffErr != nil {
                return nil, errors.New("Error parsing 'permanentlyOff' field " + _permanentlyOffErr.Error())
            }
            _map["Struct"] = permanentlyOff

            // Simple Field (chillerFailure)
            chillerFailure, _chillerFailureErr := io.ReadBit()
            if _chillerFailureErr != nil {
                return nil, errors.New("Error parsing 'chillerFailure' field " + _chillerFailureErr.Error())
            }
            _map["Struct"] = chillerFailure

            // Simple Field (chillerRunningStatus)
            chillerRunningStatus, _chillerRunningStatusErr := io.ReadBit()
            if _chillerRunningStatusErr != nil {
                return nil, errors.New("Error parsing 'chillerRunningStatus' field " + _chillerRunningStatusErr.Error())
            }
            _map["Struct"] = chillerRunningStatus

            // Simple Field (relativePowerValid)
            relativePowerValid, _relativePowerValidErr := io.ReadBit()
            if _relativePowerValidErr != nil {
                return nil, errors.New("Error parsing 'relativePowerValid' field " + _relativePowerValidErr.Error())
            }
            _map["Struct"] = relativePowerValid

            // Simple Field (tempChillerValid)
            tempChillerValid, _tempChillerValidErr := io.ReadBit()
            if _tempChillerValidErr != nil {
                return nil, errors.New("Error parsing 'tempChillerValid' field " + _tempChillerValidErr.Error())
            }
            _map["Struct"] = tempChillerValid
            return values.NewPlcStruct(_map), nil
        case formatName == "U16U8N8B8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (nominalPower)
            nominalPower, _nominalPowerErr := io.ReadUint16(16)
            if _nominalPowerErr != nil {
                return nil, errors.New("Error parsing 'nominalPower' field " + _nominalPowerErr.Error())
            }
            _map["Struct"] = nominalPower

            // Simple Field (relativePowerLimit)
            relativePowerLimit, _relativePowerLimitErr := io.ReadUint8(8)
            if _relativePowerLimitErr != nil {
                return nil, errors.New("Error parsing 'relativePowerLimit' field " + _relativePowerLimitErr.Error())
            }
            _map["Struct"] = relativePowerLimit

            // Simple Field (burnerType)
            burnerType, _burnerTypeErr := io.ReadInt8(8)
            if _burnerTypeErr != nil {
                return nil, errors.New("Error parsing 'burnerType' field " + _burnerTypeErr.Error())
            }
            _map["Struct"] = burnerType

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (solidState)
            solidState, _solidStateErr := io.ReadBit()
            if _solidStateErr != nil {
                return nil, errors.New("Error parsing 'solidState' field " + _solidStateErr.Error())
            }
            _map["Struct"] = solidState

            // Simple Field (gas)
            gas, _gasErr := io.ReadBit()
            if _gasErr != nil {
                return nil, errors.New("Error parsing 'gas' field " + _gasErr.Error())
            }
            _map["Struct"] = gas

            // Simple Field (oil)
            oil, _oilErr := io.ReadBit()
            if _oilErr != nil {
                return nil, errors.New("Error parsing 'oil' field " + _oilErr.Error())
            }
            _map["Struct"] = oil
            return values.NewPlcStruct(_map), nil
        case formatName == "U5U5U6": // Struct
            _map := map[string]interface{}{}

            // Simple Field (magicNumber)
            magicNumber, _magicNumberErr := io.ReadUint8(5)
            if _magicNumberErr != nil {
                return nil, errors.New("Error parsing 'magicNumber' field " + _magicNumberErr.Error())
            }
            _map["Struct"] = magicNumber

            // Simple Field (versionNumber)
            versionNumber, _versionNumberErr := io.ReadUint8(5)
            if _versionNumberErr != nil {
                return nil, errors.New("Error parsing 'versionNumber' field " + _versionNumberErr.Error())
            }
            _map["Struct"] = versionNumber

            // Simple Field (revisionNumber)
            revisionNumber, _revisionNumberErr := io.ReadUint8(6)
            if _revisionNumberErr != nil {
                return nil, errors.New("Error parsing 'revisionNumber' field " + _revisionNumberErr.Error())
            }
            _map["Struct"] = revisionNumber
            return values.NewPlcStruct(_map), nil
        case formatName == "V32Z8VolumeLiter": // Struct
            _map := map[string]interface{}{}

            // Simple Field (volumeLiter)
            volumeLiter, _volumeLiterErr := io.ReadInt32(32)
            if _volumeLiterErr != nil {
                return nil, errors.New("Error parsing 'volumeLiter' field " + _volumeLiterErr.Error())
            }
            _map["Struct"] = volumeLiter

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "V32Z8FlowRate": // Struct
            _map := map[string]interface{}{}

            // Simple Field (flowRate)
            flowRate, _flowRateErr := io.ReadInt32(32)
            if _flowRateErr != nil {
                return nil, errors.New("Error parsing 'flowRate' field " + _flowRateErr.Error())
            }
            _map["Struct"] = flowRate

            // Simple Field (statusCommand)
            statusCommand, _statusCommandErr := io.ReadUint8(8)
            if _statusCommandErr != nil {
                return nil, errors.New("Error parsing 'statusCommand' field " + _statusCommandErr.Error())
            }
            _map["Struct"] = statusCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U8N8N8N8B8B8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (logNumber)
            logNumber, _logNumberErr := io.ReadUint8(8)
            if _logNumberErr != nil {
                return nil, errors.New("Error parsing 'logNumber' field " + _logNumberErr.Error())
            }
            _map["Struct"] = logNumber

            // Simple Field (alarmPriority)
            alarmPriority, _alarmPriorityErr := io.ReadUint8(8)
            if _alarmPriorityErr != nil {
                return nil, errors.New("Error parsing 'alarmPriority' field " + _alarmPriorityErr.Error())
            }
            _map["Struct"] = alarmPriority

            // Simple Field (applicationArea)
            applicationArea, _applicationAreaErr := io.ReadUint8(8)
            if _applicationAreaErr != nil {
                return nil, errors.New("Error parsing 'applicationArea' field " + _applicationAreaErr.Error())
            }
            _map["Struct"] = applicationArea

            // Simple Field (errorClass)
            errorClass, _errorClassErr := io.ReadUint8(8)
            if _errorClassErr != nil {
                return nil, errors.New("Error parsing 'errorClass' field " + _errorClassErr.Error())
            }
            _map["Struct"] = errorClass

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (errorCode_Sup)
            errorCode_Sup, _errorCode_SupErr := io.ReadBit()
            if _errorCode_SupErr != nil {
                return nil, errors.New("Error parsing 'errorCode_Sup' field " + _errorCode_SupErr.Error())
            }
            _map["Struct"] = errorCode_Sup

            // Simple Field (alarmText_Sup)
            alarmText_Sup, _alarmText_SupErr := io.ReadBit()
            if _alarmText_SupErr != nil {
                return nil, errors.New("Error parsing 'alarmText_Sup' field " + _alarmText_SupErr.Error())
            }
            _map["Struct"] = alarmText_Sup

            // Simple Field (timeStamp_Sup)
            timeStamp_Sup, _timeStamp_SupErr := io.ReadBit()
            if _timeStamp_SupErr != nil {
                return nil, errors.New("Error parsing 'timeStamp_Sup' field " + _timeStamp_SupErr.Error())
            }
            _map["Struct"] = timeStamp_Sup

            // Simple Field (ack_Sup)
            ack_Sup, _ack_SupErr := io.ReadBit()
            if _ack_SupErr != nil {
                return nil, errors.New("Error parsing 'ack_Sup' field " + _ack_SupErr.Error())
            }
            _map["Struct"] = ack_Sup

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (alarmUnAck)
            alarmUnAck, _alarmUnAckErr := io.ReadBit()
            if _alarmUnAckErr != nil {
                return nil, errors.New("Error parsing 'alarmUnAck' field " + _alarmUnAckErr.Error())
            }
            _map["Struct"] = alarmUnAck

            // Simple Field (locked)
            locked, _lockedErr := io.ReadBit()
            if _lockedErr != nil {
                return nil, errors.New("Error parsing 'locked' field " + _lockedErr.Error())
            }
            _map["Struct"] = locked

            // Simple Field (inAlarm)
            inAlarm, _inAlarmErr := io.ReadBit()
            if _inAlarmErr != nil {
                return nil, errors.New("Error parsing 'inAlarm' field " + _inAlarmErr.Error())
            }
            _map["Struct"] = inAlarm
            return values.NewPlcStruct(_map), nil
        case formatName == "U16V16": // Struct
            _map := map[string]interface{}{}

            // Simple Field (delayTime)
            delayTime, _delayTimeErr := io.ReadUint16(16)
            if _delayTimeErr != nil {
                return nil, errors.New("Error parsing 'delayTime' field " + _delayTimeErr.Error())
            }
            _map["Struct"] = delayTime

            // Simple Field (temperature)
            temperature, _temperatureErr := io.ReadInt16(16)
            if _temperatureErr != nil {
                return nil, errors.New("Error parsing 'temperature' field " + _temperatureErr.Error())
            }
            _map["Struct"] = temperature
            return values.NewPlcStruct(_map), nil
        case formatName == "N16U32": // Struct
            _map := map[string]interface{}{}

            // Simple Field (manufacturerCode)
            manufacturerCode, _manufacturerCodeErr := io.ReadUint16(16)
            if _manufacturerCodeErr != nil {
                return nil, errors.New("Error parsing 'manufacturerCode' field " + _manufacturerCodeErr.Error())
            }
            _map["Struct"] = manufacturerCode

            // Simple Field (incrementedNumber)
            incrementedNumber, _incrementedNumberErr := io.ReadUint32(32)
            if _incrementedNumberErr != nil {
                return nil, errors.New("Error parsing 'incrementedNumber' field " + _incrementedNumberErr.Error())
            }
            _map["Struct"] = incrementedNumber
            return values.NewPlcStruct(_map), nil
        case formatName == "F16F16F16": // Struct
            _map := map[string]interface{}{}

            // Simple Field (temperatureSetpointComfort)
            temperatureSetpointComfort, _temperatureSetpointComfortErr := io.ReadFloat32(16)
            if _temperatureSetpointComfortErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointComfort' field " + _temperatureSetpointComfortErr.Error())
            }
            _map["Struct"] = temperatureSetpointComfort

            // Simple Field (temperatureSetpointShiftStandby)
            temperatureSetpointShiftStandby, _temperatureSetpointShiftStandbyErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftStandbyErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftStandby' field " + _temperatureSetpointShiftStandbyErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftStandby

            // Simple Field (temperatureSetpointShiftEco)
            temperatureSetpointShiftEco, _temperatureSetpointShiftEcoErr := io.ReadFloat32(16)
            if _temperatureSetpointShiftEcoErr != nil {
                return nil, errors.New("Error parsing 'temperatureSetpointShiftEco' field " + _temperatureSetpointShiftEcoErr.Error())
            }
            _map["Struct"] = temperatureSetpointShiftEco
            return values.NewPlcStruct(_map), nil
        case formatName == "V8N8N8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (energyDemand)
            energyDemand, _energyDemandErr := io.ReadInt8(8)
            if _energyDemandErr != nil {
                return nil, errors.New("Error parsing 'energyDemand' field " + _energyDemandErr.Error())
            }
            _map["Struct"] = energyDemand

            // Simple Field (hvacControllerMode)
            hvacControllerMode, _hvacControllerModeErr := io.ReadUint8(8)
            if _hvacControllerModeErr != nil {
                return nil, errors.New("Error parsing 'hvacControllerMode' field " + _hvacControllerModeErr.Error())
            }
            _map["Struct"] = hvacControllerMode

            // Simple Field (hvacEmergencyMode)
            hvacEmergencyMode, _hvacEmergencyModeErr := io.ReadUint8(8)
            if _hvacEmergencyModeErr != nil {
                return nil, errors.New("Error parsing 'hvacEmergencyMode' field " + _hvacEmergencyModeErr.Error())
            }
            _map["Struct"] = hvacEmergencyMode
            return values.NewPlcStruct(_map), nil
        case formatName == "V16V16N8N8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (tempSetpointCooling)
            tempSetpointCooling, _tempSetpointCoolingErr := io.ReadInt16(16)
            if _tempSetpointCoolingErr != nil {
                return nil, errors.New("Error parsing 'tempSetpointCooling' field " + _tempSetpointCoolingErr.Error())
            }
            _map["Struct"] = tempSetpointCooling

            // Simple Field (tempSetpointHeating)
            tempSetpointHeating, _tempSetpointHeatingErr := io.ReadInt16(16)
            if _tempSetpointHeatingErr != nil {
                return nil, errors.New("Error parsing 'tempSetpointHeating' field " + _tempSetpointHeatingErr.Error())
            }
            _map["Struct"] = tempSetpointHeating

            // Simple Field (hvacControllerMode)
            hvacControllerMode, _hvacControllerModeErr := io.ReadUint8(8)
            if _hvacControllerModeErr != nil {
                return nil, errors.New("Error parsing 'hvacControllerMode' field " + _hvacControllerModeErr.Error())
            }
            _map["Struct"] = hvacControllerMode

            // Simple Field (hvacEmergencyMode)
            hvacEmergencyMode, _hvacEmergencyModeErr := io.ReadUint8(8)
            if _hvacEmergencyModeErr != nil {
                return nil, errors.New("Error parsing 'hvacEmergencyMode' field " + _hvacEmergencyModeErr.Error())
            }
            _map["Struct"] = hvacEmergencyMode
            return values.NewPlcStruct(_map), nil
        case formatName == "U16U8Scaling": // Struct
            _map := map[string]interface{}{}

            // Simple Field (timePeriod)
            timePeriod, _timePeriodErr := io.ReadUint16(16)
            if _timePeriodErr != nil {
                return nil, errors.New("Error parsing 'timePeriod' field " + _timePeriodErr.Error())
            }
            _map["Struct"] = timePeriod

            // Simple Field (percent)
            percent, _percentErr := io.ReadUint8(8)
            if _percentErr != nil {
                return nil, errors.New("Error parsing 'percent' field " + _percentErr.Error())
            }
            _map["Struct"] = percent
            return values.NewPlcStruct(_map), nil
        case formatName == "U16U8TariffNext": // Struct
            _map := map[string]interface{}{}

            // Simple Field (delayTime)
            delayTime, _delayTimeErr := io.ReadUint16(16)
            if _delayTimeErr != nil {
                return nil, errors.New("Error parsing 'delayTime' field " + _delayTimeErr.Error())
            }
            _map["Struct"] = delayTime

            // Simple Field (tariff)
            tariff, _tariffErr := io.ReadUint8(8)
            if _tariffErr != nil {
                return nil, errors.New("Error parsing 'tariff' field " + _tariffErr.Error())
            }
            _map["Struct"] = tariff
            return values.NewPlcStruct(_map), nil
        case formatName == "V32N8Z8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (countVal)
            countVal, _countValErr := io.ReadInt32(32)
            if _countValErr != nil {
                return nil, errors.New("Error parsing 'countVal' field " + _countValErr.Error())
            }
            _map["Struct"] = countVal

            // Simple Field (valInfField)
            valInfField, _valInfFieldErr := io.ReadUint8(8)
            if _valInfFieldErr != nil {
                return nil, errors.New("Error parsing 'valInfField' field " + _valInfFieldErr.Error())
            }
            _map["Struct"] = valInfField

            // Simple Field (statusOrCommand)
            statusOrCommand, _statusOrCommandErr := io.ReadUint8(8)
            if _statusOrCommandErr != nil {
                return nil, errors.New("Error parsing 'statusOrCommand' field " + _statusOrCommandErr.Error())
            }
            _map["Struct"] = statusOrCommand
            return values.NewPlcStruct(_map), nil
        case formatName == "U16U32U8N8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (manufacturerId)
            manufacturerId, _manufacturerIdErr := io.ReadUint16(16)
            if _manufacturerIdErr != nil {
                return nil, errors.New("Error parsing 'manufacturerId' field " + _manufacturerIdErr.Error())
            }
            _map["Struct"] = manufacturerId

            // Simple Field (identNumber)
            identNumber, _identNumberErr := io.ReadUint32(32)
            if _identNumberErr != nil {
                return nil, errors.New("Error parsing 'identNumber' field " + _identNumberErr.Error())
            }
            _map["Struct"] = identNumber

            // Simple Field (version)
            version, _versionErr := io.ReadUint8(8)
            if _versionErr != nil {
                return nil, errors.New("Error parsing 'version' field " + _versionErr.Error())
            }
            _map["Struct"] = version

            // Simple Field (medium)
            medium, _mediumErr := io.ReadInt8(8)
            if _mediumErr != nil {
                return nil, errors.New("Error parsing 'medium' field " + _mediumErr.Error())
            }
            _map["Struct"] = medium
            return values.NewPlcStruct(_map), nil
        case formatName == "A8A8A8A8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (languageCode)
            languageCode, _languageCodeErr := io.ReadString(16)
            if _languageCodeErr != nil {
                return nil, errors.New("Error parsing 'languageCode' field " + _languageCodeErr.Error())
            }
            _map["Struct"] = languageCode

            // Simple Field (regionCode)
            regionCode, _regionCodeErr := io.ReadString(16)
            if _regionCodeErr != nil {
                return nil, errors.New("Error parsing 'regionCode' field " + _regionCodeErr.Error())
            }
            _map["Struct"] = regionCode
            return values.NewPlcStruct(_map), nil
        case formatName == "U8U8U8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (red)
            red, _redErr := io.ReadUint8(8)
            if _redErr != nil {
                return nil, errors.New("Error parsing 'red' field " + _redErr.Error())
            }
            _map["Struct"] = red

            // Simple Field (green)
            green, _greenErr := io.ReadUint8(8)
            if _greenErr != nil {
                return nil, errors.New("Error parsing 'green' field " + _greenErr.Error())
            }
            _map["Struct"] = green

            // Simple Field (blue)
            blue, _blueErr := io.ReadUint8(8)
            if _blueErr != nil {
                return nil, errors.New("Error parsing 'blue' field " + _blueErr.Error())
            }
            _map["Struct"] = blue
            return values.NewPlcStruct(_map), nil
        case formatName == "A8A8Language": // Struct
            _map := map[string]interface{}{}

            // Simple Field (languageCode)
            languageCode, _languageCodeErr := io.ReadString(16)
            if _languageCodeErr != nil {
                return nil, errors.New("Error parsing 'languageCode' field " + _languageCodeErr.Error())
            }
            _map["Struct"] = languageCode
            return values.NewPlcStruct(_map), nil
        case formatName == "A8A8Region": // Struct
            _map := map[string]interface{}{}

            // Simple Field (regionCode)
            regionCode, _regionCodeErr := io.ReadString(16)
            if _regionCodeErr != nil {
                return nil, errors.New("Error parsing 'regionCode' field " + _regionCodeErr.Error())
            }
            _map["Struct"] = regionCode
            return values.NewPlcStruct(_map), nil
        case formatName == "V32U8B8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (activeElectricalEnergy)
            activeElectricalEnergy, _activeElectricalEnergyErr := io.ReadInt32(32)
            if _activeElectricalEnergyErr != nil {
                return nil, errors.New("Error parsing 'activeElectricalEnergy' field " + _activeElectricalEnergyErr.Error())
            }
            _map["Struct"] = activeElectricalEnergy

            // Simple Field (tariff)
            tariff, _tariffErr := io.ReadUint8(8)
            if _tariffErr != nil {
                return nil, errors.New("Error parsing 'tariff' field " + _tariffErr.Error())
            }
            _map["Struct"] = tariff

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (noTariff)
            noTariff, _noTariffErr := io.ReadBit()
            if _noTariffErr != nil {
                return nil, errors.New("Error parsing 'noTariff' field " + _noTariffErr.Error())
            }
            _map["Struct"] = noTariff

            // Simple Field (noActiveElectricalEnergy)
            noActiveElectricalEnergy, _noActiveElectricalEnergyErr := io.ReadBit()
            if _noActiveElectricalEnergyErr != nil {
                return nil, errors.New("Error parsing 'noActiveElectricalEnergy' field " + _noActiveElectricalEnergyErr.Error())
            }
            _map["Struct"] = noActiveElectricalEnergy
            return values.NewPlcStruct(_map), nil
        case formatName == "B1N3N4": // Struct
            _map := map[string]interface{}{}

            // Simple Field (deactivationOfPriority)
            deactivationOfPriority, _deactivationOfPriorityErr := io.ReadBit()
            if _deactivationOfPriorityErr != nil {
                return nil, errors.New("Error parsing 'deactivationOfPriority' field " + _deactivationOfPriorityErr.Error())
            }
            _map["Struct"] = deactivationOfPriority

            // Simple Field (priorityLevel)
            priorityLevel, _priorityLevelErr := io.ReadUint8(3)
            if _priorityLevelErr != nil {
                return nil, errors.New("Error parsing 'priorityLevel' field " + _priorityLevelErr.Error())
            }
            _map["Struct"] = priorityLevel

            // Simple Field (modeLevel)
            modeLevel, _modeLevelErr := io.ReadUint8(4)
            if _modeLevelErr != nil {
                return nil, errors.New("Error parsing 'modeLevel' field " + _modeLevelErr.Error())
            }
            _map["Struct"] = modeLevel
            return values.NewPlcStruct(_map), nil
        case formatName == "B10U6": // Struct
            _map := map[string]interface{}{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (convertorError)
            convertorError, _convertorErrorErr := io.ReadBit()
            if _convertorErrorErr != nil {
                return nil, errors.New("Error parsing 'convertorError' field " + _convertorErrorErr.Error())
            }
            _map["Struct"] = convertorError

            // Simple Field (ballastFailure)
            ballastFailure, _ballastFailureErr := io.ReadBit()
            if _ballastFailureErr != nil {
                return nil, errors.New("Error parsing 'ballastFailure' field " + _ballastFailureErr.Error())
            }
            _map["Struct"] = ballastFailure

            // Simple Field (lampError)
            lampError, _lampErrorErr := io.ReadBit()
            if _lampErrorErr != nil {
                return nil, errors.New("Error parsing 'lampError' field " + _lampErrorErr.Error())
            }
            _map["Struct"] = lampError

            // Simple Field (read)
            read, _readErr := io.ReadBit()
            if _readErr != nil {
                return nil, errors.New("Error parsing 'read' field " + _readErr.Error())
            }
            _map["Struct"] = read

            // Simple Field (groupAddress)
            groupAddress, _groupAddressErr := io.ReadBit()
            if _groupAddressErr != nil {
                return nil, errors.New("Error parsing 'groupAddress' field " + _groupAddressErr.Error())
            }
            _map["Struct"] = groupAddress

            // Simple Field (address)
            address, _addressErr := io.ReadUint8(6)
            if _addressErr != nil {
                return nil, errors.New("Error parsing 'address' field " + _addressErr.Error())
            }
            _map["Struct"] = address
            return values.NewPlcStruct(_map), nil
        case formatName == "B2U6": // Struct
            _map := map[string]interface{}{}

            // Simple Field (sceneActivationInactive)
            sceneActivationInactive, _sceneActivationInactiveErr := io.ReadBit()
            if _sceneActivationInactiveErr != nil {
                return nil, errors.New("Error parsing 'sceneActivationInactive' field " + _sceneActivationInactiveErr.Error())
            }
            _map["Struct"] = sceneActivationInactive

            // Simple Field (storageFunctionDisable)
            storageFunctionDisable, _storageFunctionDisableErr := io.ReadBit()
            if _storageFunctionDisableErr != nil {
                return nil, errors.New("Error parsing 'storageFunctionDisable' field " + _storageFunctionDisableErr.Error())
            }
            _map["Struct"] = storageFunctionDisable

            // Simple Field (sceneNumber)
            sceneNumber, _sceneNumberErr := io.ReadUint8(6)
            if _sceneNumberErr != nil {
                return nil, errors.New("Error parsing 'sceneNumber' field " + _sceneNumberErr.Error())
            }
            _map["Struct"] = sceneNumber
            return values.NewPlcStruct(_map), nil
        case formatName == "U8r7B1": // Struct
            _map := map[string]interface{}{}

            // Simple Field (setValue)
            setValue, _setValueErr := io.ReadUint8(8)
            if _setValueErr != nil {
                return nil, errors.New("Error parsing 'setValue' field " + _setValueErr.Error())
            }
            _map["Struct"] = setValue

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (channelActivationActive)
            channelActivationActive, _channelActivationActiveErr := io.ReadBit()
            if _channelActivationActiveErr != nil {
                return nil, errors.New("Error parsing 'channelActivationActive' field " + _channelActivationActiveErr.Error())
            }
            _map["Struct"] = channelActivationActive
            return values.NewPlcStruct(_map), nil
        case formatName == "U8U8B8": // Struct
            _map := map[string]interface{}{}

            // Simple Field (heightPosition)
            heightPosition, _heightPositionErr := io.ReadUint8(8)
            if _heightPositionErr != nil {
                return nil, errors.New("Error parsing 'heightPosition' field " + _heightPositionErr.Error())
            }
            _map["Struct"] = heightPosition

            // Simple Field (slatsPosition)
            slatsPosition, _slatsPositionErr := io.ReadUint8(8)
            if _slatsPositionErr != nil {
                return nil, errors.New("Error parsing 'slatsPosition' field " + _slatsPositionErr.Error())
            }
            _map["Struct"] = slatsPosition

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(6); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (validSlatsPos)
            validSlatsPos, _validSlatsPosErr := io.ReadBit()
            if _validSlatsPosErr != nil {
                return nil, errors.New("Error parsing 'validSlatsPos' field " + _validSlatsPosErr.Error())
            }
            _map["Struct"] = validSlatsPos

            // Simple Field (validHeightPos)
            validHeightPos, _validHeightPosErr := io.ReadBit()
            if _validHeightPosErr != nil {
                return nil, errors.New("Error parsing 'validHeightPos' field " + _validHeightPosErr.Error())
            }
            _map["Struct"] = validHeightPos
            return values.NewPlcStruct(_map), nil
        case formatName == "U8U8B16": // Struct
            _map := map[string]interface{}{}

            // Simple Field (heightPosition)
            heightPosition, _heightPositionErr := io.ReadUint8(8)
            if _heightPositionErr != nil {
                return nil, errors.New("Error parsing 'heightPosition' field " + _heightPositionErr.Error())
            }
            _map["Struct"] = heightPosition

            // Simple Field (slatsPosition)
            slatsPosition, _slatsPositionErr := io.ReadUint8(8)
            if _slatsPositionErr != nil {
                return nil, errors.New("Error parsing 'slatsPosition' field " + _slatsPositionErr.Error())
            }
            _map["Struct"] = slatsPosition

            // Simple Field (validSlatsPos)
            validSlatsPos, _validSlatsPosErr := io.ReadBit()
            if _validSlatsPosErr != nil {
                return nil, errors.New("Error parsing 'validSlatsPos' field " + _validSlatsPosErr.Error())
            }
            _map["Struct"] = validSlatsPos

            // Simple Field (validHeightPos)
            validHeightPos, _validHeightPosErr := io.ReadBit()
            if _validHeightPosErr != nil {
                return nil, errors.New("Error parsing 'validHeightPos' field " + _validHeightPosErr.Error())
            }
            _map["Struct"] = validHeightPos

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (failure)
            failure, _failureErr := io.ReadBit()
            if _failureErr != nil {
                return nil, errors.New("Error parsing 'failure' field " + _failureErr.Error())
            }
            _map["Struct"] = failure

            // Simple Field (localOverride)
            localOverride, _localOverrideErr := io.ReadBit()
            if _localOverrideErr != nil {
                return nil, errors.New("Error parsing 'localOverride' field " + _localOverrideErr.Error())
            }
            _map["Struct"] = localOverride

            // Simple Field (locked)
            locked, _lockedErr := io.ReadBit()
            if _lockedErr != nil {
                return nil, errors.New("Error parsing 'locked' field " + _lockedErr.Error())
            }
            _map["Struct"] = locked

            // Simple Field (forced)
            forced, _forcedErr := io.ReadBit()
            if _forcedErr != nil {
                return nil, errors.New("Error parsing 'forced' field " + _forcedErr.Error())
            }
            _map["Struct"] = forced

            // Simple Field (weatherAlarm)
            weatherAlarm, _weatherAlarmErr := io.ReadBit()
            if _weatherAlarmErr != nil {
                return nil, errors.New("Error parsing 'weatherAlarm' field " + _weatherAlarmErr.Error())
            }
            _map["Struct"] = weatherAlarm

            // Simple Field (targetSPosRestrict)
            targetSPosRestrict, _targetSPosRestrictErr := io.ReadBit()
            if _targetSPosRestrictErr != nil {
                return nil, errors.New("Error parsing 'targetSPosRestrict' field " + _targetSPosRestrictErr.Error())
            }
            _map["Struct"] = targetSPosRestrict

            // Simple Field (targetHPosRestrict)
            targetHPosRestrict, _targetHPosRestrictErr := io.ReadBit()
            if _targetHPosRestrictErr != nil {
                return nil, errors.New("Error parsing 'targetHPosRestrict' field " + _targetHPosRestrictErr.Error())
            }
            _map["Struct"] = targetHPosRestrict

            // Simple Field (driveState)
            driveState, _driveStateErr := io.ReadBit()
            if _driveStateErr != nil {
                return nil, errors.New("Error parsing 'driveState' field " + _driveStateErr.Error())
            }
            _map["Struct"] = driveState

            // Simple Field (lowerPredefPos)
            lowerPredefPos, _lowerPredefPosErr := io.ReadBit()
            if _lowerPredefPosErr != nil {
                return nil, errors.New("Error parsing 'lowerPredefPos' field " + _lowerPredefPosErr.Error())
            }
            _map["Struct"] = lowerPredefPos

            // Simple Field (lowerEndPos)
            lowerEndPos, _lowerEndPosErr := io.ReadBit()
            if _lowerEndPosErr != nil {
                return nil, errors.New("Error parsing 'lowerEndPos' field " + _lowerEndPosErr.Error())
            }
            _map["Struct"] = lowerEndPos

            // Simple Field (upperEndPos)
            upperEndPos, _upperEndPosErr := io.ReadBit()
            if _upperEndPosErr != nil {
                return nil, errors.New("Error parsing 'upperEndPos' field " + _upperEndPosErr.Error())
            }
            _map["Struct"] = upperEndPos
            return values.NewPlcStruct(_map), nil
    }
    return nil, errors.New("unsupported type")
}

func KnxDatapointSerialize(io *utils.WriteBuffer, value api.PlcValue, formatName string) error {
    switch {
        case formatName == "B1": // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B2": // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'control' field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B1U3": // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (control)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'control' field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "A8_ASCII": // STRING

            // Simple Field (value)
            if _err := io.WriteString(8, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "A8_8859_1": // STRING

            // Simple Field (value)
            if _err := io.WriteString(8, "ISO-8859-1", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "U8": // USINT

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "V8": // SINT

            // Simple Field (value)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "F16": // REAL

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "N3N5r2N6r2N6": // TIME_OF_DAY

            // Simple Field (day)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'day' field " + _err.Error())
            }

            // Simple Field (hour)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hour' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (minutes)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'minutes' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (seconds)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'seconds' field " + _err.Error())
            }
        case formatName == "r3N5r4N4r1U7": // DATE

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (day)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'day' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (month)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'month' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (year)
            if _err := io.WriteUint8(7, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'year' field " + _err.Error())
            }
        case formatName == "U32": // UDINT

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "V32": // DINT

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "F32": // REAL

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "U4U4U4U4U4U4B4N4": // Struct

            // Simple Field (d6)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'd6' field " + _err.Error())
            }

            // Simple Field (d5)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'd5' field " + _err.Error())
            }

            // Simple Field (d4)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'd4' field " + _err.Error())
            }

            // Simple Field (d3)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'd3' field " + _err.Error())
            }

            // Simple Field (d2)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'd2' field " + _err.Error())
            }

            // Simple Field (d1)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'd1' field " + _err.Error())
            }

            // Simple Field (e)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'e' field " + _err.Error())
            }

            // Simple Field (p)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'p' field " + _err.Error())
            }

            // Simple Field (d)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'd' field " + _err.Error())
            }

            // Simple Field (c)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'c' field " + _err.Error())
            }

            // Simple Field (index)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'index' field " + _err.Error())
            }
        case formatName == "A112_ASCII": // STRING

            // Simple Field (value)
            if _err := io.WriteString(112, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "A112_8859_1": // STRING

            // Simple Field (value)
            if _err := io.WriteString(112, "ISO-8859-1", value.GetString()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "r2U6": // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B1r1U6": // Struct

            // Simple Field (learn)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'learn' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (sceneNumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'sceneNumber' field " + _err.Error())
            }
        case formatName == "U8r4U4r3U5U3U5r2U6r2U6B16": // DATE_AND_TIME

            // Simple Field (year)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'year' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (month)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'month' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (day)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'day' field " + _err.Error())
            }

            // Simple Field (dayOfWeek)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'dayOfWeek' field " + _err.Error())
            }

            // Simple Field (hour)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hour' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (minutes)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'minutes' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (seconds)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'seconds' field " + _err.Error())
            }

            // Simple Field (fault)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'fault' field " + _err.Error())
            }

            // Simple Field (workingDay)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'workingDay' field " + _err.Error())
            }

            // Simple Field (noWorkingDay)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noWorkingDay' field " + _err.Error())
            }

            // Simple Field (noYear)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noYear' field " + _err.Error())
            }

            // Simple Field (noMonthAndDay)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noMonthAndDay' field " + _err.Error())
            }

            // Simple Field (noDayOfWeek)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noDayOfWeek' field " + _err.Error())
            }

            // Simple Field (noTime)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noTime' field " + _err.Error())
            }

            // Simple Field (standardSummerTime)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'standardSummerTime' field " + _err.Error())
            }

            // Simple Field (clockWithSyncSignal)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'clockWithSyncSignal' field " + _err.Error())
            }
        case formatName == "N8": // USINT

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B8": // BYTE

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B16": // WORD

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "U4U4": // Struct

            // Simple Field (busy)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'busy' field " + _err.Error())
            }

            // Simple Field (nak)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'nak' field " + _err.Error())
            }
        case formatName == "r1b1U6": // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(1, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (sceneActive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'sceneActive' field " + _err.Error())
            }

            // Simple Field (sceneNumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'sceneNumber' field " + _err.Error())
            }
        case formatName == "B32": // DWORD

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "V64": // LINT

            // Simple Field (value)
            if _err := io.WriteInt64(64, value.GetInt64()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B24": // List

            // Array Field (value)
            for i := uint32(0); i < uint32((24)); i++ {
                _itemErr := io.WriteBit(value.GetIndex(i).GetBool())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case formatName == "N3": // USINT

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case formatName == "B1Z8HeatingOrCoolingZ": // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (heating)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'heating' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "B1Z8BinaryValueZ": // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (high)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'high' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8HvacOperatingMode": // Struct

            // Simple Field (hvacOperatingMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacOperatingMode' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8DhwMode": // Struct

            // Simple Field (dhwMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'dhwMode' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8HvacControllingMode": // Struct

            // Simple Field (hvacControllingMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacControllingMode' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8EnableHeatingOrCoolingStage": // Struct

            // Simple Field (enableHeatingOrCoolingStage)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'enableHeatingOrCoolingStage' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8BuildingMode": // Struct

            // Simple Field (buildingMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'buildingMode' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8OccupancyMode": // Struct

            // Simple Field (occupancyMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'occupancyMode' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "N8Z8EmergencyMode": // Struct

            // Simple Field (hvacEmergencyMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacEmergencyMode' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U8Z8Rel": // Struct

            // Simple Field (relValue)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relValue' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U8Z8Counter": // Struct

            // Simple Field (counterValue)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'counterValue' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8TimePeriod": // Struct

            // Simple Field (timePeriod)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'timePeriod' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8FlowRate": // Struct

            // Simple Field (flowRate)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'flowRate' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8Counter": // Struct

            // Simple Field (counterValue)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'counterValue' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8ElectricCurrent": // Struct

            // Simple Field (electricalCurrent)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'electricalCurrent' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8Power": // Struct

            // Simple Field (power)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'power' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8AtmPressure": // Struct

            // Simple Field (atmPressure)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'atmPressure' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8PercentValue": // Struct

            // Simple Field (percentValue)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'percentValue' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8HvacAirQuality": // Struct

            // Simple Field (ppmResolution)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'ppmResolution' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8WindSpeed": // Struct

            // Simple Field (windSpeed)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'windSpeed' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8SunIntensity": // Struct

            // Simple Field (sunIntensity)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'sunIntensity' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16Z8HvacAirFlow": // Struct

            // Simple Field (airFlow)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'airFlow' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "V8Z8RelSignedValue": // Struct

            // Simple Field (relSignedValue)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'relSignedValue' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "V16Z8DeltaTime": // Struct

            // Simple Field (deltaTime)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'deltaTime' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "V16Z8RelSignedValue": // Struct

            // Simple Field (relSignedValue)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'relSignedValue' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U16N8HvacModeAndTimeDelay": // Struct

            // Simple Field (delayTime)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'delayTime' field " + _err.Error())
            }

            // Simple Field (hvacMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacMode' field " + _err.Error())
            }
        case formatName == "U16N8DhwModeAndTimeDelay": // Struct

            // Simple Field (delayTime)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'delayTime' field " + _err.Error())
            }

            // Simple Field (dhwMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'dhwMode' field " + _err.Error())
            }
        case formatName == "U16N8OccupancyModeAndTimeDelay": // Struct

            // Simple Field (delayTime)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'delayTime' field " + _err.Error())
            }

            // Simple Field (occupationMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'occupationMode' field " + _err.Error())
            }
        case formatName == "U16N8BuildingModeAndTimeDelay": // Struct

            // Simple Field (delayTime)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'delayTime' field " + _err.Error())
            }

            // Simple Field (buildingMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'buildingMode' field " + _err.Error())
            }
        case formatName == "U8B8StatusBurnerController": // Struct

            // Simple Field (actualRelativePower)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'actualRelativePower' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (stage2Active)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage2Active' field " + _err.Error())
            }

            // Simple Field (stage1Active)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage1Active' field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }

            // Simple Field (actualRelativePowerValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'actualRelativePowerValid' field " + _err.Error())
            }
        case formatName == "U8B8LockingSignal": // Struct

            // Simple Field (requestedPowerReduction)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'requestedPowerReduction' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (critical)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'critical' field " + _err.Error())
            }

            // Simple Field (requestedPowerReductionValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'requestedPowerReductionValid' field " + _err.Error())
            }
        case formatName == "U8B8BoilerControllerDemandSignal": // Struct

            // Simple Field (relativeDemand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relativeDemand' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (controlsOperationStage2)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'controlsOperationStage2' field " + _err.Error())
            }

            // Simple Field (controlsOperationStage1)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'controlsOperationStage1' field " + _err.Error())
            }
        case formatName == "U8B8ActuatorPositionDemand": // Struct

            // Simple Field (actuatorPositionDemand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'actuatorPositionDemand' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (emergencyDemand)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'emergencyDemand' field " + _err.Error())
            }

            // Simple Field (shiftLoadPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'shiftLoadPriority' field " + _err.Error())
            }

            // Simple Field (absoluteLoadPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'absoluteLoadPriority' field " + _err.Error())
            }

            // Simple Field (actuatorPositionDemandValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'actuatorPositionDemandValid' field " + _err.Error())
            }
        case formatName == "U8B8ActuatorPositionStatus": // Struct

            // Simple Field (actualActuatorPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'actualActuatorPosition' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (synchronizationMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'synchronizationMode' field " + _err.Error())
            }

            // Simple Field (valveKick)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'valveKick' field " + _err.Error())
            }

            // Simple Field (callibrationMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'callibrationMode' field " + _err.Error())
            }

            // Simple Field (positionManuallyOverridden)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'positionManuallyOverridden' field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }
        case formatName == "U8B8StatusLightingActuator": // Struct

            // Simple Field (lightingLevel)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'lightingLevel' field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }

            // Simple Field (localOverride)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'localOverride' field " + _err.Error())
            }

            // Simple Field (dimming)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'dimming' field " + _err.Error())
            }

            // Simple Field (staircaseLightingFunction)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'staircaseLightingFunction' field " + _err.Error())
            }

            // Simple Field (nightMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'nightMode' field " + _err.Error())
            }

            // Simple Field (forced)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'forced' field " + _err.Error())
            }

            // Simple Field (locked)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'locked' field " + _err.Error())
            }

            // Simple Field (lightingLevelValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'lightingLevelValid' field " + _err.Error())
            }
        case formatName == "V16B8HeatProducerManagerStatus": // Struct

            // Simple Field (tempFlowProdSegmH)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'tempFlowProdSegmH' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (temporarilyOff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'temporarilyOff' field " + _err.Error())
            }

            // Simple Field (permanentlyOff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'permanentlyOff' field " + _err.Error())
            }

            // Simple Field (switchedOffSummerMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'switchedOffSummerMode' field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }

            // Simple Field (tempFlowProdSegmHValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'tempFlowProdSegmHValid' field " + _err.Error())
            }
        case formatName == "V16B8RoomTemperatureDemand": // Struct

            // Simple Field (roomTemperatureDemand)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'roomTemperatureDemand' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (emergencyDemand)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'emergencyDemand' field " + _err.Error())
            }

            // Simple Field (shiftLoadPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'shiftLoadPriority' field " + _err.Error())
            }

            // Simple Field (absoluteLoadPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'absoluteLoadPriority' field " + _err.Error())
            }

            // Simple Field (roomTemperatureDemandValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'roomTemperatureDemandValid' field " + _err.Error())
            }
        case formatName == "V16B8ColdWaterProducerManagerStatus": // Struct

            // Simple Field (flowTemperatureProdSegmC)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureProdSegmC' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (temporarilyOff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'temporarilyOff' field " + _err.Error())
            }

            // Simple Field (permanentlyOff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'permanentlyOff' field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }

            // Simple Field (flowTemperatureProdSegmCValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureProdSegmCValid' field " + _err.Error())
            }
        case formatName == "V16B8WaterTemperatureControllerStatus": // Struct

            // Simple Field (actualTemperature)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'actualTemperature' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (controllerWorking)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'controllerWorking' field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }

            // Simple Field (actualTemperatureValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'actualTemperatureValid' field " + _err.Error())
            }
        case formatName == "V16B16": // Struct

            // Simple Field (flowTemperatureDemand)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureDemand' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (demandFromDhwWhileLegionellaFunctionIsActive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'demandFromDhwWhileLegionellaFunctionIsActive' field " + _err.Error())
            }

            // Simple Field (emergencyDemandForFrostProtection)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'emergencyDemandForFrostProtection' field " + _err.Error())
            }

            // Simple Field (requestForWaterCirculationInPrimaryDistributionSegment)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'requestForWaterCirculationInPrimaryDistributionSegment' field " + _err.Error())
            }

            // Simple Field (demandFromAuxillaryHeatOrCoolConsumer)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'demandFromAuxillaryHeatOrCoolConsumer' field " + _err.Error())
            }

            // Simple Field (demandFromVentilation)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'demandFromVentilation' field " + _err.Error())
            }

            // Simple Field (demandForRoomHeatingOrCooling)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'demandForRoomHeatingOrCooling' field " + _err.Error())
            }

            // Simple Field (heatDemandFromDhw)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'heatDemandFromDhw' field " + _err.Error())
            }

            // Simple Field (flowTemperatureDemandIsMin)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureDemandIsMin' field " + _err.Error())
            }

            // Simple Field (flowTemperatureDemandIsMax)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureDemandIsMax' field " + _err.Error())
            }

            // Simple Field (shiftLoadPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'shiftLoadPriority' field " + _err.Error())
            }

            // Simple Field (absoluteLoadPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'absoluteLoadPriority' field " + _err.Error())
            }

            // Simple Field (flowTemperatureDemandValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureDemandValid' field " + _err.Error())
            }
        case formatName == "U8N8": // Struct

            // Simple Field (energyDemand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'energyDemand' field " + _err.Error())
            }

            // Simple Field (actualControllerMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'actualControllerMode' field " + _err.Error())
            }
        case formatName == "V16V16V16RoomTemperature": // Struct

            // Simple Field (temperatureSetpointComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointComfort' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointStandby' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointEco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointEco' field " + _err.Error())
            }
        case formatName == "V16V16V16RoomTemperatureShift": // Struct

            // Simple Field (temperatureSetpointShiftComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftComfort' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftStandby' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftEco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftEco' field " + _err.Error())
            }
        case formatName == "V16V16V16V16RoomTemperature": // Struct

            // Simple Field (temperatureSetpointComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointComfort' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointStandby' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointEco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointEco' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointBProt)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointBProt' field " + _err.Error())
            }
        case formatName == "V16V16V16V16DhwtTemperature": // Struct

            // Simple Field (temperatureSetpointLegioProtect)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointLegioProtect' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointNormal)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointNormal' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointReduced)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointReduced' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointFrostProtect)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointFrostProtect' field " + _err.Error())
            }
        case formatName == "V16V16V16V16RoomTemperatureShift": // Struct

            // Simple Field (temperatureSetpointShiftComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftComfort' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftStandby' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftEco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftEco' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftBProt)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftBProt' field " + _err.Error())
            }
        case formatName == "V16U8B8Heat": // Struct

            // Simple Field (flowTemperatureDemand)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureDemand' field " + _err.Error())
            }

            // Simple Field (relativePower)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relativePower' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(2, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (boilerEnabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'boilerEnabled' field " + _err.Error())
            }

            // Simple Field (stage2Forced)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage2Forced' field " + _err.Error())
            }

            // Simple Field (stage2Enabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage2Enabled' field " + _err.Error())
            }

            // Simple Field (stage1Forced)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage1Forced' field " + _err.Error())
            }

            // Simple Field (stage1Enabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage1Enabled' field " + _err.Error())
            }

            // Simple Field (flowTemperatureDemandValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'flowTemperatureDemandValid' field " + _err.Error())
            }
        case formatName == "V16U8B8ChilledWater": // Struct

            // Simple Field (chilledWaterFlowTemperatureDemand)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'chilledWaterFlowTemperatureDemand' field " + _err.Error())
            }

            // Simple Field (relativePower)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relativePower' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (chilledWaterPumpEnabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'chilledWaterPumpEnabled' field " + _err.Error())
            }

            // Simple Field (relativePowerValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'relativePowerValid' field " + _err.Error())
            }

            // Simple Field (chilledWaterFlowTemperatureDemandValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'chilledWaterFlowTemperatureDemandValid' field " + _err.Error())
            }
        case formatName == "V16U8B16Boiler": // Struct

            // Simple Field (tempBoiler)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'tempBoiler' field " + _err.Error())
            }

            // Simple Field (relativePower)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relativePower' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (chimneySweepFunctionActive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'chimneySweepFunctionActive' field " + _err.Error())
            }

            // Simple Field (reducedAvailability)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'reducedAvailability' field " + _err.Error())
            }

            // Simple Field (powerLimitBoilerReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'powerLimitBoilerReached' field " + _err.Error())
            }

            // Simple Field (powerLimitStage1Reached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'powerLimitStage1Reached' field " + _err.Error())
            }

            // Simple Field (stage2Enabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage2Enabled' field " + _err.Error())
            }

            // Simple Field (stage1Enabled)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'stage1Enabled' field " + _err.Error())
            }

            // Simple Field (boilerTemporarilyNotProvidingHeat)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'boilerTemporarilyNotProvidingHeat' field " + _err.Error())
            }

            // Simple Field (permanentlyOff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'permanentlyOff' field " + _err.Error())
            }

            // Simple Field (boilerSwitchedOffWinterSummerMode)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'boilerSwitchedOffWinterSummerMode' field " + _err.Error())
            }

            // Simple Field (boilerFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'boilerFailure' field " + _err.Error())
            }

            // Simple Field (relativePowerValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'relativePowerValid' field " + _err.Error())
            }

            // Simple Field (tempBoilerValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'tempBoilerValid' field " + _err.Error())
            }
        case formatName == "V16U8B16Chiller": // Struct

            // Simple Field (tempChiller)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'tempChiller' field " + _err.Error())
            }

            // Simple Field (relativePower)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relativePower' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(8, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (reducedAvailability)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'reducedAvailability' field " + _err.Error())
            }

            // Simple Field (powerLimitChillerReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'powerLimitChillerReached' field " + _err.Error())
            }

            // Simple Field (powerLimitCurrentStageReached)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'powerLimitCurrentStageReached' field " + _err.Error())
            }

            // Simple Field (permanentlyOff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'permanentlyOff' field " + _err.Error())
            }

            // Simple Field (chillerFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'chillerFailure' field " + _err.Error())
            }

            // Simple Field (chillerRunningStatus)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'chillerRunningStatus' field " + _err.Error())
            }

            // Simple Field (relativePowerValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'relativePowerValid' field " + _err.Error())
            }

            // Simple Field (tempChillerValid)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'tempChillerValid' field " + _err.Error())
            }
        case formatName == "U16U8N8B8": // Struct

            // Simple Field (nominalPower)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'nominalPower' field " + _err.Error())
            }

            // Simple Field (relativePowerLimit)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'relativePowerLimit' field " + _err.Error())
            }

            // Simple Field (burnerType)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'burnerType' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (solidState)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'solidState' field " + _err.Error())
            }

            // Simple Field (gas)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'gas' field " + _err.Error())
            }

            // Simple Field (oil)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'oil' field " + _err.Error())
            }
        case formatName == "U5U5U6": // Struct

            // Simple Field (magicNumber)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'magicNumber' field " + _err.Error())
            }

            // Simple Field (versionNumber)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'versionNumber' field " + _err.Error())
            }

            // Simple Field (revisionNumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'revisionNumber' field " + _err.Error())
            }
        case formatName == "V32Z8VolumeLiter": // Struct

            // Simple Field (volumeLiter)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'volumeLiter' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "V32Z8FlowRate": // Struct

            // Simple Field (flowRate)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'flowRate' field " + _err.Error())
            }

            // Simple Field (statusCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusCommand' field " + _err.Error())
            }
        case formatName == "U8N8N8N8B8B8": // Struct

            // Simple Field (logNumber)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'logNumber' field " + _err.Error())
            }

            // Simple Field (alarmPriority)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'alarmPriority' field " + _err.Error())
            }

            // Simple Field (applicationArea)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'applicationArea' field " + _err.Error())
            }

            // Simple Field (errorClass)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'errorClass' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(4, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (errorCode_Sup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'errorCode_Sup' field " + _err.Error())
            }

            // Simple Field (alarmText_Sup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'alarmText_Sup' field " + _err.Error())
            }

            // Simple Field (timeStamp_Sup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'timeStamp_Sup' field " + _err.Error())
            }

            // Simple Field (ack_Sup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ack_Sup' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (alarmUnAck)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'alarmUnAck' field " + _err.Error())
            }

            // Simple Field (locked)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'locked' field " + _err.Error())
            }

            // Simple Field (inAlarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'inAlarm' field " + _err.Error())
            }
        case formatName == "U16V16": // Struct

            // Simple Field (delayTime)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'delayTime' field " + _err.Error())
            }

            // Simple Field (temperature)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'temperature' field " + _err.Error())
            }
        case formatName == "N16U32": // Struct

            // Simple Field (manufacturerCode)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'manufacturerCode' field " + _err.Error())
            }

            // Simple Field (incrementedNumber)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'incrementedNumber' field " + _err.Error())
            }
        case formatName == "F16F16F16": // Struct

            // Simple Field (temperatureSetpointComfort)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointComfort' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftStandby)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftStandby' field " + _err.Error())
            }

            // Simple Field (temperatureSetpointShiftEco)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'temperatureSetpointShiftEco' field " + _err.Error())
            }
        case formatName == "V8N8N8": // Struct

            // Simple Field (energyDemand)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'energyDemand' field " + _err.Error())
            }

            // Simple Field (hvacControllerMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacControllerMode' field " + _err.Error())
            }

            // Simple Field (hvacEmergencyMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacEmergencyMode' field " + _err.Error())
            }
        case formatName == "V16V16N8N8": // Struct

            // Simple Field (tempSetpointCooling)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'tempSetpointCooling' field " + _err.Error())
            }

            // Simple Field (tempSetpointHeating)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'tempSetpointHeating' field " + _err.Error())
            }

            // Simple Field (hvacControllerMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacControllerMode' field " + _err.Error())
            }

            // Simple Field (hvacEmergencyMode)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hvacEmergencyMode' field " + _err.Error())
            }
        case formatName == "U16U8Scaling": // Struct

            // Simple Field (timePeriod)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'timePeriod' field " + _err.Error())
            }

            // Simple Field (percent)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'percent' field " + _err.Error())
            }
        case formatName == "U16U8TariffNext": // Struct

            // Simple Field (delayTime)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'delayTime' field " + _err.Error())
            }

            // Simple Field (tariff)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'tariff' field " + _err.Error())
            }
        case formatName == "V32N8Z8": // Struct

            // Simple Field (countVal)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'countVal' field " + _err.Error())
            }

            // Simple Field (valInfField)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'valInfField' field " + _err.Error())
            }

            // Simple Field (statusOrCommand)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'statusOrCommand' field " + _err.Error())
            }
        case formatName == "U16U32U8N8": // Struct

            // Simple Field (manufacturerId)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'manufacturerId' field " + _err.Error())
            }

            // Simple Field (identNumber)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'identNumber' field " + _err.Error())
            }

            // Simple Field (version)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'version' field " + _err.Error())
            }

            // Simple Field (medium)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'medium' field " + _err.Error())
            }
        case formatName == "A8A8A8A8": // Struct

            // Simple Field (languageCode)
            if _err := io.WriteString(16, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'languageCode' field " + _err.Error())
            }

            // Simple Field (regionCode)
            if _err := io.WriteString(16, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'regionCode' field " + _err.Error())
            }
        case formatName == "U8U8U8": // Struct

            // Simple Field (red)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'red' field " + _err.Error())
            }

            // Simple Field (green)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'green' field " + _err.Error())
            }

            // Simple Field (blue)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'blue' field " + _err.Error())
            }
        case formatName == "A8A8Language": // Struct

            // Simple Field (languageCode)
            if _err := io.WriteString(16, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'languageCode' field " + _err.Error())
            }
        case formatName == "A8A8Region": // Struct

            // Simple Field (regionCode)
            if _err := io.WriteString(16, "ASCII", value.GetString()); _err != nil {
                return errors.New("Error serializing 'regionCode' field " + _err.Error())
            }
        case formatName == "V32U8B8": // Struct

            // Simple Field (activeElectricalEnergy)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'activeElectricalEnergy' field " + _err.Error())
            }

            // Simple Field (tariff)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'tariff' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (noTariff)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noTariff' field " + _err.Error())
            }

            // Simple Field (noActiveElectricalEnergy)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noActiveElectricalEnergy' field " + _err.Error())
            }
        case formatName == "B1N3N4": // Struct

            // Simple Field (deactivationOfPriority)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'deactivationOfPriority' field " + _err.Error())
            }

            // Simple Field (priorityLevel)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'priorityLevel' field " + _err.Error())
            }

            // Simple Field (modeLevel)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'modeLevel' field " + _err.Error())
            }
        case formatName == "B10U6": // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (convertorError)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'convertorError' field " + _err.Error())
            }

            // Simple Field (ballastFailure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ballastFailure' field " + _err.Error())
            }

            // Simple Field (lampError)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'lampError' field " + _err.Error())
            }

            // Simple Field (read)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'read' field " + _err.Error())
            }

            // Simple Field (groupAddress)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'groupAddress' field " + _err.Error())
            }

            // Simple Field (address)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'address' field " + _err.Error())
            }
        case formatName == "B2U6": // Struct

            // Simple Field (sceneActivationInactive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'sceneActivationInactive' field " + _err.Error())
            }

            // Simple Field (storageFunctionDisable)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'storageFunctionDisable' field " + _err.Error())
            }

            // Simple Field (sceneNumber)
            if _err := io.WriteUint8(6, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'sceneNumber' field " + _err.Error())
            }
        case formatName == "U8r7B1": // Struct

            // Simple Field (setValue)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'setValue' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (channelActivationActive)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'channelActivationActive' field " + _err.Error())
            }
        case formatName == "U8U8B8": // Struct

            // Simple Field (heightPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'heightPosition' field " + _err.Error())
            }

            // Simple Field (slatsPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'slatsPosition' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(6, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (validSlatsPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'validSlatsPos' field " + _err.Error())
            }

            // Simple Field (validHeightPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'validHeightPos' field " + _err.Error())
            }
        case formatName == "U8U8B16": // Struct

            // Simple Field (heightPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'heightPosition' field " + _err.Error())
            }

            // Simple Field (slatsPosition)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'slatsPosition' field " + _err.Error())
            }

            // Simple Field (validSlatsPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'validSlatsPos' field " + _err.Error())
            }

            // Simple Field (validHeightPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'validHeightPos' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (failure)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'failure' field " + _err.Error())
            }

            // Simple Field (localOverride)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'localOverride' field " + _err.Error())
            }

            // Simple Field (locked)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'locked' field " + _err.Error())
            }

            // Simple Field (forced)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'forced' field " + _err.Error())
            }

            // Simple Field (weatherAlarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'weatherAlarm' field " + _err.Error())
            }

            // Simple Field (targetSPosRestrict)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'targetSPosRestrict' field " + _err.Error())
            }

            // Simple Field (targetHPosRestrict)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'targetHPosRestrict' field " + _err.Error())
            }

            // Simple Field (driveState)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'driveState' field " + _err.Error())
            }

            // Simple Field (lowerPredefPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'lowerPredefPos' field " + _err.Error())
            }

            // Simple Field (lowerEndPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'lowerEndPos' field " + _err.Error())
            }

            // Simple Field (upperEndPos)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'upperEndPos' field " + _err.Error())
            }
        default:

            return errors.New("unsupported type")
    }
    return nil
}

