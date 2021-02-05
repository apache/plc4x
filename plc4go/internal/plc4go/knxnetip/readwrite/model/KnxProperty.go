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
            "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
            "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
            api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

func KnxPropertyParse(io *utils.ReadBuffer, propertyType KnxPropertyDataType, dataLengthInBytes uint8) (api.PlcValue, error) {
    switch {
        case propertyType == KnxPropertyDataType_PDT_CONTROL: // BOOL

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
        case propertyType == KnxPropertyDataType_PDT_CHAR: // SINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcSINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_CHAR: // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_INT: // INT

            // Simple Field (value)
            value, _valueErr := io.ReadInt16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_INT && dataLengthInBytes == 4: // UDINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_INT: // UINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint16(16)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_KNX_FLOAT: // REAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 4, 11)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case propertyType == KnxPropertyDataType_PDT_DATE: // Struct
            _map := map[string]api.PlcValue{}

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (dayOfMonth)
            dayOfMonth, _dayOfMonthErr := io.ReadUint8(5)
            if _dayOfMonthErr != nil {
                return nil, errors.New("Error parsing 'dayOfMonth' field " + _dayOfMonthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(dayOfMonth)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (month)
            month, _monthErr := io.ReadUint8(4)
            if _monthErr != nil {
                return nil, errors.New("Error parsing 'month' field " + _monthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(month)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(1); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (year)
            year, _yearErr := io.ReadUint8(7)
            if _yearErr != nil {
                return nil, errors.New("Error parsing 'year' field " + _yearErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(year)
            return values.NewPlcStruct(_map), nil
        case propertyType == KnxPropertyDataType_PDT_TIME: // Struct
            _map := map[string]api.PlcValue{}

            // Simple Field (day)
            day, _dayErr := io.ReadUint8(3)
            if _dayErr != nil {
                return nil, errors.New("Error parsing 'day' field " + _dayErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(day)

            // Simple Field (hour)
            hour, _hourErr := io.ReadUint8(5)
            if _hourErr != nil {
                return nil, errors.New("Error parsing 'hour' field " + _hourErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(hour)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (minutes)
            minutes, _minutesErr := io.ReadUint8(6)
            if _minutesErr != nil {
                return nil, errors.New("Error parsing 'minutes' field " + _minutesErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(minutes)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (seconds)
            seconds, _secondsErr := io.ReadUint8(6)
            if _secondsErr != nil {
                return nil, errors.New("Error parsing 'seconds' field " + _secondsErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(seconds)
            return values.NewPlcStruct(_map), nil
        case propertyType == KnxPropertyDataType_PDT_LONG: // DINT

            // Simple Field (value)
            value, _valueErr := io.ReadInt32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcDINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_LONG: // UDINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint32(32)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUDINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_FLOAT: // REAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat32(true, 8, 23)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcREAL(value), nil
        case propertyType == KnxPropertyDataType_PDT_DOUBLE: // LREAL

            // Simple Field (value)
            value, _valueErr := io.ReadFloat64(true, 11, 52)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcLREAL(value), nil
        case propertyType == KnxPropertyDataType_PDT_CHAR_BLOCK: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((10)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS: // Struct
            _map := map[string]api.PlcValue{}

            // Array Field (groupAddress)
            var groupAddress []api.PlcValue
            for i := 0; i < int((2)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                groupAddress = append(groupAddress, values.NewPlcUSINT(_item))
            }

            // Simple Field (disable)
            disable, _disableErr := io.ReadBit()
            if _disableErr != nil {
                return nil, errors.New("Error parsing 'disable' field " + _disableErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(disable)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (pollingSoftNr)
            pollingSoftNr, _pollingSoftNrErr := io.ReadUint8(4)
            if _pollingSoftNrErr != nil {
                return nil, errors.New("Error parsing 'pollingSoftNr' field " + _pollingSoftNrErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(pollingSoftNr)
            return values.NewPlcStruct(_map), nil
        case propertyType == KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((5)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_DATE_TIME: // Struct
            _map := map[string]api.PlcValue{}

            // Simple Field (year)
            year, _yearErr := io.ReadUint8(8)
            if _yearErr != nil {
                return nil, errors.New("Error parsing 'year' field " + _yearErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(year)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (month)
            month, _monthErr := io.ReadUint8(4)
            if _monthErr != nil {
                return nil, errors.New("Error parsing 'month' field " + _monthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(month)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(3); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (dayofmonth)
            dayofmonth, _dayofmonthErr := io.ReadUint8(5)
            if _dayofmonthErr != nil {
                return nil, errors.New("Error parsing 'dayofmonth' field " + _dayofmonthErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(dayofmonth)

            // Simple Field (dayofweek)
            dayofweek, _dayofweekErr := io.ReadUint8(3)
            if _dayofweekErr != nil {
                return nil, errors.New("Error parsing 'dayofweek' field " + _dayofweekErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(dayofweek)

            // Simple Field (hourofday)
            hourofday, _hourofdayErr := io.ReadUint8(5)
            if _hourofdayErr != nil {
                return nil, errors.New("Error parsing 'hourofday' field " + _hourofdayErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(hourofday)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (minutes)
            minutes, _minutesErr := io.ReadUint8(6)
            if _minutesErr != nil {
                return nil, errors.New("Error parsing 'minutes' field " + _minutesErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(minutes)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(2); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (seconds)
            seconds, _secondsErr := io.ReadUint8(6)
            if _secondsErr != nil {
                return nil, errors.New("Error parsing 'seconds' field " + _secondsErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(seconds)

            // Simple Field (fault)
            fault, _faultErr := io.ReadBit()
            if _faultErr != nil {
                return nil, errors.New("Error parsing 'fault' field " + _faultErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(fault)

            // Simple Field (workingDay)
            workingDay, _workingDayErr := io.ReadBit()
            if _workingDayErr != nil {
                return nil, errors.New("Error parsing 'workingDay' field " + _workingDayErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(workingDay)

            // Simple Field (noWd)
            noWd, _noWdErr := io.ReadBit()
            if _noWdErr != nil {
                return nil, errors.New("Error parsing 'noWd' field " + _noWdErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(noWd)

            // Simple Field (noYear)
            noYear, _noYearErr := io.ReadBit()
            if _noYearErr != nil {
                return nil, errors.New("Error parsing 'noYear' field " + _noYearErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(noYear)

            // Simple Field (noDate)
            noDate, _noDateErr := io.ReadBit()
            if _noDateErr != nil {
                return nil, errors.New("Error parsing 'noDate' field " + _noDateErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(noDate)

            // Simple Field (noDayOfWeek)
            noDayOfWeek, _noDayOfWeekErr := io.ReadBit()
            if _noDayOfWeekErr != nil {
                return nil, errors.New("Error parsing 'noDayOfWeek' field " + _noDayOfWeekErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(noDayOfWeek)

            // Simple Field (noTime)
            noTime, _noTimeErr := io.ReadBit()
            if _noTimeErr != nil {
                return nil, errors.New("Error parsing 'noTime' field " + _noTimeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(noTime)

            // Simple Field (standardSummerTime)
            standardSummerTime, _standardSummerTimeErr := io.ReadBit()
            if _standardSummerTimeErr != nil {
                return nil, errors.New("Error parsing 'standardSummerTime' field " + _standardSummerTimeErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(standardSummerTime)

            // Simple Field (qualityOfClock)
            qualityOfClock, _qualityOfClockErr := io.ReadBit()
            if _qualityOfClockErr != nil {
                return nil, errors.New("Error parsing 'qualityOfClock' field " + _qualityOfClockErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(qualityOfClock)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(7); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }
            return values.NewPlcStruct(_map), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_01: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((1)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_02: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((2)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_03: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((3)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_04: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((4)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_05: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((5)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_06: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((6)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_07: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((7)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_08: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((8)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_09: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((9)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_10: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((10)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_11: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((11)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_12: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((12)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_13: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((13)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_14: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((14)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_15: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((15)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_16: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((16)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_17: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((17)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_18: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((18)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_19: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((19)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_GENERIC_20: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((20)); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_VERSION: // Struct
            _map := map[string]api.PlcValue{}

            // Simple Field (magicNumber)
            magicNumber, _magicNumberErr := io.ReadUint8(5)
            if _magicNumberErr != nil {
                return nil, errors.New("Error parsing 'magicNumber' field " + _magicNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(magicNumber)

            // Simple Field (versionNumber)
            versionNumber, _versionNumberErr := io.ReadUint8(5)
            if _versionNumberErr != nil {
                return nil, errors.New("Error parsing 'versionNumber' field " + _versionNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(versionNumber)

            // Simple Field (revisionNumber)
            revisionNumber, _revisionNumberErr := io.ReadUint8(6)
            if _revisionNumberErr != nil {
                return nil, errors.New("Error parsing 'revisionNumber' field " + _revisionNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(revisionNumber)
            return values.NewPlcStruct(_map), nil
        case propertyType == KnxPropertyDataType_PDT_ALARM_INFO: // Struct
            _map := map[string]api.PlcValue{}

            // Simple Field (logNumber)
            logNumber, _logNumberErr := io.ReadUint8(8)
            if _logNumberErr != nil {
                return nil, errors.New("Error parsing 'logNumber' field " + _logNumberErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(logNumber)

            // Simple Field (alarmPriority)
            alarmPriority, _alarmPriorityErr := io.ReadUint8(8)
            if _alarmPriorityErr != nil {
                return nil, errors.New("Error parsing 'alarmPriority' field " + _alarmPriorityErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(alarmPriority)

            // Simple Field (applicationArea)
            applicationArea, _applicationAreaErr := io.ReadUint8(8)
            if _applicationAreaErr != nil {
                return nil, errors.New("Error parsing 'applicationArea' field " + _applicationAreaErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(applicationArea)

            // Simple Field (errorClass)
            errorClass, _errorClassErr := io.ReadUint8(8)
            if _errorClassErr != nil {
                return nil, errors.New("Error parsing 'errorClass' field " + _errorClassErr.Error())
            }
            _map["Struct"] = values.NewPlcUSINT(errorClass)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(4); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (errorcodeSup)
            errorcodeSup, _errorcodeSupErr := io.ReadBit()
            if _errorcodeSupErr != nil {
                return nil, errors.New("Error parsing 'errorcodeSup' field " + _errorcodeSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(errorcodeSup)

            // Simple Field (alarmtextSup)
            alarmtextSup, _alarmtextSupErr := io.ReadBit()
            if _alarmtextSupErr != nil {
                return nil, errors.New("Error parsing 'alarmtextSup' field " + _alarmtextSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(alarmtextSup)

            // Simple Field (timestampSup)
            timestampSup, _timestampSupErr := io.ReadBit()
            if _timestampSupErr != nil {
                return nil, errors.New("Error parsing 'timestampSup' field " + _timestampSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(timestampSup)

            // Simple Field (ackSup)
            ackSup, _ackSupErr := io.ReadBit()
            if _ackSupErr != nil {
                return nil, errors.New("Error parsing 'ackSup' field " + _ackSupErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(ackSup)

            // Reserved Field (Just skip the bytes)
            if _, _err := io.ReadUint8(5); _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }

            // Simple Field (locked)
            locked, _lockedErr := io.ReadBit()
            if _lockedErr != nil {
                return nil, errors.New("Error parsing 'locked' field " + _lockedErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(locked)

            // Simple Field (alarmunack)
            alarmunack, _alarmunackErr := io.ReadBit()
            if _alarmunackErr != nil {
                return nil, errors.New("Error parsing 'alarmunack' field " + _alarmunackErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(alarmunack)

            // Simple Field (inalarm)
            inalarm, _inalarmErr := io.ReadBit()
            if _inalarmErr != nil {
                return nil, errors.New("Error parsing 'inalarm' field " + _inalarmErr.Error())
            }
            _map["Struct"] = values.NewPlcBOOL(inalarm)
            return values.NewPlcStruct(_map), nil
        case propertyType == KnxPropertyDataType_PDT_BINARY_INFORMATION: // BOOL

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
        case propertyType == KnxPropertyDataType_PDT_BITSET8: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((8)); i++ {
                _item, _itemErr := io.ReadBit()
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcBOOL(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_BITSET16: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int((16)); i++ {
                _item, _itemErr := io.ReadBit()
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcBOOL(_item))
            }
            return values.NewPlcList(value), nil
        case propertyType == KnxPropertyDataType_PDT_ENUM8: // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        case propertyType == KnxPropertyDataType_PDT_SCALING: // USINT

            // Simple Field (value)
            value, _valueErr := io.ReadUint8(8)
            if _valueErr != nil {
                return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
            }
            return values.NewPlcUSINT(value), nil
        default: // List

            // Array Field (value)
            var value []api.PlcValue
            for i := 0; i < int(dataLengthInBytes); i++ {
                _item, _itemErr := io.ReadUint8(8)
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                value = append(value, values.NewPlcUSINT(_item))
            }
            return values.NewPlcList(value), nil
    }
    return nil, errors.New("unsupported type")
}

func KnxPropertySerialize(io *utils.WriteBuffer, value api.PlcValue, propertyType KnxPropertyDataType, dataLengthInBytes uint8) error {
    switch {
        case propertyType == KnxPropertyDataType_PDT_CONTROL: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_CHAR: // SINT

            // Simple Field (value)
            if _err := io.WriteInt8(8, value.GetInt8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_CHAR: // USINT

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_INT: // INT

            // Simple Field (value)
            if _err := io.WriteInt16(16, value.GetInt16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_INT && dataLengthInBytes == 4: // UDINT

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_INT: // UINT

            // Simple Field (value)
            if _err := io.WriteUint16(16, value.GetUint16()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_KNX_FLOAT: // REAL

            // Simple Field (value)
            if _err := io.WriteFloat32(16, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_DATE: // Struct

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (dayOfMonth)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'dayOfMonth' field " + _err.Error())
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
        case propertyType == KnxPropertyDataType_PDT_TIME: // Struct

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
        case propertyType == KnxPropertyDataType_PDT_LONG: // DINT

            // Simple Field (value)
            if _err := io.WriteInt32(32, value.GetInt32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_UNSIGNED_LONG: // UDINT

            // Simple Field (value)
            if _err := io.WriteUint32(32, value.GetUint32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_FLOAT: // REAL

            // Simple Field (value)
            if _err := io.WriteFloat32(32, value.GetFloat32()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_DOUBLE: // LREAL

            // Simple Field (value)
            if _err := io.WriteFloat64(64, value.GetFloat64()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_CHAR_BLOCK: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((10)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_POLL_GROUP_SETTINGS: // Struct

            // Array Field (groupAddress)
            for i := uint32(0); i < uint32((2)); i++ {
                groupAddress := value.GetValue("groupAddress")
                _itemErr := io.WriteUint8(8, groupAddress.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }

            // Simple Field (disable)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'disable' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(3, uint8(0x0)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (pollingSoftNr)
            if _err := io.WriteUint8(4, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'pollingSoftNr' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_SHORT_CHAR_BLOCK: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((5)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_DATE_TIME: // Struct

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

            // Simple Field (dayofmonth)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'dayofmonth' field " + _err.Error())
            }

            // Simple Field (dayofweek)
            if _err := io.WriteUint8(3, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'dayofweek' field " + _err.Error())
            }

            // Simple Field (hourofday)
            if _err := io.WriteUint8(5, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'hourofday' field " + _err.Error())
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

            // Simple Field (noWd)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noWd' field " + _err.Error())
            }

            // Simple Field (noYear)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noYear' field " + _err.Error())
            }

            // Simple Field (noDate)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'noDate' field " + _err.Error())
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

            // Simple Field (qualityOfClock)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'qualityOfClock' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_01: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((1)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_02: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((2)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_03: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((3)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_04: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((4)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_05: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((5)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_06: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((6)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_07: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((7)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_08: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((8)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_09: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((9)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_10: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((10)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_11: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((11)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_12: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((12)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_13: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((13)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_14: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((14)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_15: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((15)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_16: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((16)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_17: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((17)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_18: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((18)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_19: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((19)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_GENERIC_20: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((20)); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_VERSION: // Struct

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
        case propertyType == KnxPropertyDataType_PDT_ALARM_INFO: // Struct

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

            // Simple Field (errorcodeSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'errorcodeSup' field " + _err.Error())
            }

            // Simple Field (alarmtextSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'alarmtextSup' field " + _err.Error())
            }

            // Simple Field (timestampSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'timestampSup' field " + _err.Error())
            }

            // Simple Field (ackSup)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'ackSup' field " + _err.Error())
            }

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(5, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (locked)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'locked' field " + _err.Error())
            }

            // Simple Field (alarmunack)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'alarmunack' field " + _err.Error())
            }

            // Simple Field (inalarm)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'inalarm' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_BINARY_INFORMATION: // BOOL

            // Reserved Field (Just skip the bytes)
            if _err := io.WriteUint8(7, uint8(0x00)); _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }

            // Simple Field (value)
            if _err := io.WriteBit(value.GetBool()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_BITSET8: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((8)); i++ {
                _itemErr := io.WriteBit(value.GetIndex(i).GetBool())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_BITSET16: // List

            // Array Field (value)
            for i := uint32(0); i < uint32((16)); i++ {
                _itemErr := io.WriteBit(value.GetIndex(i).GetBool())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
        case propertyType == KnxPropertyDataType_PDT_ENUM8: // USINT

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        case propertyType == KnxPropertyDataType_PDT_SCALING: // USINT

            // Simple Field (value)
            if _err := io.WriteUint8(8, value.GetUint8()); _err != nil {
                return errors.New("Error serializing 'value' field " + _err.Error())
            }
        default: // List

            // Array Field (value)
            for i := uint32(0); i < uint32(dataLengthInBytes); i++ {
                _itemErr := io.WriteUint8(8, value.GetIndex(i).GetUint8())
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
    }
    return nil
}

