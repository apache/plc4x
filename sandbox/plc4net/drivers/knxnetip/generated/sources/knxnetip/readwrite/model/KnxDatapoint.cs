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

using System;
using System.Collections.Generic;
using System.Text;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.model.values;

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public class KnxDatapoint
    {

        public IPlcValue Parse(ReadBuffer io, string formatName)
        {
            if (formatName == "B1") { // BOOL

                // Reserved Field (Just skip the bits)
                io.ReadByte(7);

                // Simple Field (value)
                var value = io.ReadBit();

                return new PlcBOOL(value);
            } else
            if (formatName == "B2") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(6);

                // Simple Field (control)
                var control = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(control);

                // Simple Field (value)
                var value = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(value);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "B1U3") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (control)
                var control = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(control);

                // Simple Field (value)
                var value = io.ReadByte(3);
                internalMap["Struct"] = new PlcUSINT(value);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "A8_ASCII") { // STRING

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadString(8, Encoding.GetEncoding("ASCII"));

                return new PlcSTRING(value);
            } else
            if (formatName == "A8_8859_1") { // STRING

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadString(8, Encoding.GetEncoding("ISO-8859-1"));

                return new PlcSTRING(value);
            } else
            if (formatName == "U8") { // USINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadByte(8);

                return new PlcUSINT(value);
            } else
            if (formatName == "V8") { // SINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadSbyte(8);

                return new PlcSINT(value);
            } else
            if (formatName == "B5N3") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(3);

                // Simple Field (a)
                var a = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(a);

                // Simple Field (b)
                var b = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(b);

                // Simple Field (c)
                var c = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(c);

                // Simple Field (d)
                var d = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(d);

                // Simple Field (e)
                var e = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(e);

                // Simple Field (value)
                var value = io.ReadSbyte(8);
                internalMap["Struct"] = new PlcSINT(value);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16") { // UINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadUshort(16);

                return new PlcUINT(value);
            } else
            if (formatName == "V16") { // INT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadShort(16);

                return new PlcINT(value);
            } else
            if (formatName == "F16") { // REAL

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadFloat(true, 4, 11);

                return new PlcREAL(value);
            } else
            if (formatName == "N3N5r2N6r2N6") { // TIME_OF_DAY

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (day)
                var day = io.ReadByte(3);

                // Simple Field (hour)
                var hour = io.ReadByte(5);

                // Reserved Field (Just skip the bits)
                io.ReadByte(2);

                // Simple Field (minutes)
                var minutes = io.ReadByte(6);

                // Reserved Field (Just skip the bits)
                io.ReadByte(2);

                // Simple Field (seconds)
                var seconds = io.ReadByte(6);

                var value = new DateTime(0,0,0, hour, minutes, seconds);
                return new PlcTIME_OF_DAY(value);
            } else
            if (formatName == "r3N5r4N4r1U7") { // DATE

                // Reserved Field (Just skip the bits)
                io.ReadByte(3);

                // Simple Field (day)
                var day = io.ReadByte(5);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (month)
                var month = io.ReadByte(4);

                // Reserved Field (Just skip the bits)
                io.ReadByte(1);

                // Simple Field (year)
                var year = io.ReadByte(7);

                var value = new DateTime(year, month, day, 0, 0, 0);
                return new PlcDATE(value);
            } else
            if (formatName == "U32") { // UDINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadUint(32);

                return new PlcUDINT(value);
            } else
            if (formatName == "V32") { // DINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadInt(32);

                return new PlcDINT(value);
            } else
            if (formatName == "F32") { // REAL

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadFloat(true, 8, 23);

                return new PlcREAL(value);
            } else
            if (formatName == "U4U4U4U4U4U4B4N4") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (d6)
                var d6 = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(d6);

                // Simple Field (d5)
                var d5 = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(d5);

                // Simple Field (d4)
                var d4 = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(d4);

                // Simple Field (d3)
                var d3 = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(d3);

                // Simple Field (d2)
                var d2 = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(d2);

                // Simple Field (d1)
                var d1 = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(d1);

                // Simple Field (e)
                var e = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(e);

                // Simple Field (p)
                var p = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(p);

                // Simple Field (d)
                var d = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(d);

                // Simple Field (c)
                var c = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(c);

                // Simple Field (index)
                var index = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(index);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "A112_ASCII") { // STRING

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadString(112, Encoding.GetEncoding("ASCII"));

                return new PlcSTRING(value);
            } else
            if (formatName == "A112_8859_1") { // STRING

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadString(112, Encoding.GetEncoding("ISO-8859-1"));

                return new PlcSTRING(value);
            } else
            if (formatName == "r2U6") { // USINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(2);

                // Simple Field (value)
                var value = io.ReadByte(6);

                return new PlcUSINT(value);
            } else
            if (formatName == "B1r1U6") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (learn)
                var learn = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(learn);

                // Reserved Field (Just skip the bits)
                io.ReadByte(1);

                // Simple Field (sceneNumber)
                var sceneNumber = io.ReadByte(6);
                internalMap["Struct"] = new PlcUSINT(sceneNumber);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8r4U4r3U5U3U5r2U6r2U6B16") { // DATE_AND_TIME

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (year)
                var year = io.ReadByte(8);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (month)
                var month = io.ReadByte(4);

                // Reserved Field (Just skip the bits)
                io.ReadByte(3);

                // Simple Field (day)
                var day = io.ReadByte(5);

                // Simple Field (dayOfWeek)
                var dayOfWeek = io.ReadByte(3);

                // Simple Field (hour)
                var hour = io.ReadByte(5);

                // Reserved Field (Just skip the bits)
                io.ReadByte(2);

                // Simple Field (minutes)
                var minutes = io.ReadByte(6);

                // Reserved Field (Just skip the bits)
                io.ReadByte(2);

                // Simple Field (seconds)
                var seconds = io.ReadByte(6);

                // Simple Field (fault)
                var fault = io.ReadBit();

                // Simple Field (workingDay)
                var workingDay = io.ReadBit();

                // Simple Field (noWorkingDay)
                var noWorkingDay = io.ReadBit();

                // Simple Field (noYear)
                var noYear = io.ReadBit();

                // Simple Field (noMonthAndDay)
                var noMonthAndDay = io.ReadBit();

                // Simple Field (noDayOfWeek)
                var noDayOfWeek = io.ReadBit();

                // Simple Field (noTime)
                var noTime = io.ReadBit();

                // Simple Field (standardSummerTime)
                var standardSummerTime = io.ReadBit();

                // Simple Field (clockWithSyncSignal)
                var clockWithSyncSignal = io.ReadBit();

                var value = new DateTime(year, month, day, hour, minutes, seconds);
                return new PlcDATE_AND_TIME(value);
            } else
            if (formatName == "N8") { // USINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadByte(8);

                return new PlcUSINT(value);
            } else
            if (formatName == "B8") { // BYTE

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadByte(8);

                return new PlcBYTE(value);
            } else
            if (formatName == "B16") { // WORD

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadUshort(16);

                return new PlcWORD(value);
            } else
            if (formatName == "U4U4") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (busy)
                var busy = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(busy);

                // Simple Field (nak)
                var nak = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(nak);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "r1b1U6") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Reserved Field (Just skip the bits)
                io.ReadByte(1);

                // Simple Field (sceneActive)
                var sceneActive = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(sceneActive);

                // Simple Field (sceneNumber)
                var sceneNumber = io.ReadByte(6);
                internalMap["Struct"] = new PlcUSINT(sceneNumber);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "B32") { // DWORD

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadUint(32);

                return new PlcDWORD(value);
            } else
            if (formatName == "V64") { // LINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (value)
                var value = io.ReadLong(64);

                return new PlcLINT(value);
            } else
            if (formatName == "B24") { // List

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Array Field (value);
                var value = new List<IPlcValue>();
                for (int i = 0; i < (24); i++) {
                    var internalItem = io.ReadBit();
                    value.Add(new PlcBOOL(internalItem));
                }

                return new PlcList(value);
            } else
            if (formatName == "N3") { // USINT

                // Reserved Field (Just skip the bits)
                io.ReadByte(5);

                // Simple Field (value)
                var value = io.ReadByte(3);

                return new PlcUSINT(value);
            } else
            if (formatName == "B1Z8HeatingOrCoolingZ") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(7);

                // Simple Field (heating)
                var heating = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(heating);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "B1Z8BinaryValueZ") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(7);

                // Simple Field (high)
                var high = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(high);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8HvacOperatingMode") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (hvacOperatingMode)
                var hvacOperatingMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacOperatingMode);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8DhwMode") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (dhwMode)
                var dhwMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(dhwMode);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8HvacControllingMode") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (hvacControllingMode)
                var hvacControllingMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacControllingMode);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8EnableHeatingOrCoolingStage") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (enableHeatingOrCoolingStage)
                var enableHeatingOrCoolingStage = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(enableHeatingOrCoolingStage);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8BuildingMode") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (buildingMode)
                var buildingMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(buildingMode);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8OccupancyMode") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (occupancyMode)
                var occupancyMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(occupancyMode);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N8Z8EmergencyMode") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (hvacEmergencyMode)
                var hvacEmergencyMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacEmergencyMode);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8Z8Rel") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (relValue)
                var relValue = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relValue);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8Z8Counter") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (counterValue)
                var counterValue = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(counterValue);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8TimePeriod") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (timePeriod)
                var timePeriod = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(timePeriod);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8FlowRate") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (flowRate)
                var flowRate = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(flowRate);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8Counter") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (counterValue)
                var counterValue = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(counterValue);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8ElectricCurrent") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (electricalCurrent)
                var electricalCurrent = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(electricalCurrent);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8Power") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (power)
                var power = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(power);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8AtmPressure") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (atmPressure)
                var atmPressure = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(atmPressure);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8PercentValue") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (percentValue)
                var percentValue = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(percentValue);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8HvacAirQuality") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (ppmResolution)
                var ppmResolution = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(ppmResolution);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8WindSpeed") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (windSpeed)
                var windSpeed = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(windSpeed);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8SunIntensity") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (sunIntensity)
                var sunIntensity = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(sunIntensity);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16Z8HvacAirFlow") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (airFlow)
                var airFlow = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(airFlow);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V8Z8RelSignedValue") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (relSignedValue)
                var relSignedValue = io.ReadSbyte(8);
                internalMap["Struct"] = new PlcSINT(relSignedValue);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16Z8DeltaTime") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (deltaTime)
                var deltaTime = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(deltaTime);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16Z8RelSignedValue") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (relSignedValue)
                var relSignedValue = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(relSignedValue);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16N8HvacModeAndTimeDelay") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (delayTime)
                var delayTime = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(delayTime);

                // Simple Field (hvacMode)
                var hvacMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16N8DhwModeAndTimeDelay") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (delayTime)
                var delayTime = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(delayTime);

                // Simple Field (dhwMode)
                var dhwMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(dhwMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16N8OccupancyModeAndTimeDelay") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (delayTime)
                var delayTime = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(delayTime);

                // Simple Field (occupationMode)
                var occupationMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(occupationMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16N8BuildingModeAndTimeDelay") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (delayTime)
                var delayTime = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(delayTime);

                // Simple Field (buildingMode)
                var buildingMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(buildingMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8B8StatusBurnerController") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (actualRelativePower)
                var actualRelativePower = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(actualRelativePower);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (stage2Active)
                var stage2Active = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage2Active);

                // Simple Field (stage1Active)
                var stage1Active = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage1Active);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                // Simple Field (actualRelativePowerValid)
                var actualRelativePowerValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(actualRelativePowerValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8B8LockingSignal") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (requestedPowerReduction)
                var requestedPowerReduction = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(requestedPowerReduction);

                // Reserved Field (Just skip the bits)
                io.ReadByte(6);

                // Simple Field (critical)
                var critical = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(critical);

                // Simple Field (requestedPowerReductionValid)
                var requestedPowerReductionValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(requestedPowerReductionValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8B8BoilerControllerDemandSignal") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (relativeDemand)
                var relativeDemand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relativeDemand);

                // Reserved Field (Just skip the bits)
                io.ReadByte(6);

                // Simple Field (controlsOperationStage2)
                var controlsOperationStage2 = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(controlsOperationStage2);

                // Simple Field (controlsOperationStage1)
                var controlsOperationStage1 = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(controlsOperationStage1);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8B8ActuatorPositionDemand") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (actuatorPositionDemand)
                var actuatorPositionDemand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(actuatorPositionDemand);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (emergencyDemand)
                var emergencyDemand = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(emergencyDemand);

                // Simple Field (shiftLoadPriority)
                var shiftLoadPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(shiftLoadPriority);

                // Simple Field (absoluteLoadPriority)
                var absoluteLoadPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(absoluteLoadPriority);

                // Simple Field (actuatorPositionDemandValid)
                var actuatorPositionDemandValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(actuatorPositionDemandValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8B8ActuatorPositionStatus") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (actualActuatorPosition)
                var actualActuatorPosition = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(actualActuatorPosition);

                // Reserved Field (Just skip the bits)
                io.ReadByte(3);

                // Simple Field (synchronizationMode)
                var synchronizationMode = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(synchronizationMode);

                // Simple Field (valveKick)
                var valveKick = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(valveKick);

                // Simple Field (callibrationMode)
                var callibrationMode = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(callibrationMode);

                // Simple Field (positionManuallyOverridden)
                var positionManuallyOverridden = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(positionManuallyOverridden);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8B8StatusLightingActuator") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (lightingLevel)
                var lightingLevel = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(lightingLevel);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                // Simple Field (localOverride)
                var localOverride = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(localOverride);

                // Simple Field (dimming)
                var dimming = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(dimming);

                // Simple Field (staircaseLightingFunction)
                var staircaseLightingFunction = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(staircaseLightingFunction);

                // Simple Field (nightMode)
                var nightMode = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(nightMode);

                // Simple Field (forced)
                var forced = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(forced);

                // Simple Field (locked)
                var locked = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(locked);

                // Simple Field (lightingLevelValid)
                var lightingLevelValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(lightingLevelValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16B8HeatProducerManagerStatus") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (tempFlowProdSegmH)
                var tempFlowProdSegmH = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(tempFlowProdSegmH);

                // Reserved Field (Just skip the bits)
                io.ReadByte(3);

                // Simple Field (temporarilyOff)
                var temporarilyOff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(temporarilyOff);

                // Simple Field (permanentlyOff)
                var permanentlyOff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(permanentlyOff);

                // Simple Field (switchedOffSummerMode)
                var switchedOffSummerMode = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(switchedOffSummerMode);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                // Simple Field (tempFlowProdSegmHValid)
                var tempFlowProdSegmHValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(tempFlowProdSegmHValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16B8RoomTemperatureDemand") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (roomTemperatureDemand)
                var roomTemperatureDemand = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(roomTemperatureDemand);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (emergencyDemand)
                var emergencyDemand = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(emergencyDemand);

                // Simple Field (shiftLoadPriority)
                var shiftLoadPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(shiftLoadPriority);

                // Simple Field (absoluteLoadPriority)
                var absoluteLoadPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(absoluteLoadPriority);

                // Simple Field (roomTemperatureDemandValid)
                var roomTemperatureDemandValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(roomTemperatureDemandValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16B8ColdWaterProducerManagerStatus") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (flowTemperatureProdSegmC)
                var flowTemperatureProdSegmC = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(flowTemperatureProdSegmC);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (temporarilyOff)
                var temporarilyOff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(temporarilyOff);

                // Simple Field (permanentlyOff)
                var permanentlyOff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(permanentlyOff);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                // Simple Field (flowTemperatureProdSegmCValid)
                var flowTemperatureProdSegmCValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(flowTemperatureProdSegmCValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16B8WaterTemperatureControllerStatus") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (actualTemperature)
                var actualTemperature = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(actualTemperature);

                // Reserved Field (Just skip the bits)
                io.ReadByte(5);

                // Simple Field (controllerWorking)
                var controllerWorking = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(controllerWorking);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                // Simple Field (actualTemperatureValid)
                var actualTemperatureValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(actualTemperatureValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16B16") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (flowTemperatureDemand)
                var flowTemperatureDemand = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(flowTemperatureDemand);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (demandFromDhwWhileLegionellaFunctionIsActive)
                var demandFromDhwWhileLegionellaFunctionIsActive = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(demandFromDhwWhileLegionellaFunctionIsActive);

                // Simple Field (emergencyDemandForFrostProtection)
                var emergencyDemandForFrostProtection = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(emergencyDemandForFrostProtection);

                // Simple Field (requestForWaterCirculationInPrimaryDistributionSegment)
                var requestForWaterCirculationInPrimaryDistributionSegment = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(requestForWaterCirculationInPrimaryDistributionSegment);

                // Simple Field (demandFromAuxillaryHeatOrCoolConsumer)
                var demandFromAuxillaryHeatOrCoolConsumer = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(demandFromAuxillaryHeatOrCoolConsumer);

                // Simple Field (demandFromVentilation)
                var demandFromVentilation = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(demandFromVentilation);

                // Simple Field (demandForRoomHeatingOrCooling)
                var demandForRoomHeatingOrCooling = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(demandForRoomHeatingOrCooling);

                // Simple Field (heatDemandFromDhw)
                var heatDemandFromDhw = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(heatDemandFromDhw);

                // Simple Field (flowTemperatureDemandIsMin)
                var flowTemperatureDemandIsMin = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(flowTemperatureDemandIsMin);

                // Simple Field (flowTemperatureDemandIsMax)
                var flowTemperatureDemandIsMax = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(flowTemperatureDemandIsMax);

                // Simple Field (shiftLoadPriority)
                var shiftLoadPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(shiftLoadPriority);

                // Simple Field (absoluteLoadPriority)
                var absoluteLoadPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(absoluteLoadPriority);

                // Simple Field (flowTemperatureDemandValid)
                var flowTemperatureDemandValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(flowTemperatureDemandValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8N8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (energyDemand)
                var energyDemand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(energyDemand);

                // Simple Field (actualControllerMode)
                var actualControllerMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(actualControllerMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16V16V16RoomTemperature") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (temperatureSetpointComfort)
                var temperatureSetpointComfort = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointComfort);

                // Simple Field (temperatureSetpointStandby)
                var temperatureSetpointStandby = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointStandby);

                // Simple Field (temperatureSetpointEco)
                var temperatureSetpointEco = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointEco);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16V16V16RoomTemperatureShift") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (temperatureSetpointShiftComfort)
                var temperatureSetpointShiftComfort = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftComfort);

                // Simple Field (temperatureSetpointShiftStandby)
                var temperatureSetpointShiftStandby = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftStandby);

                // Simple Field (temperatureSetpointShiftEco)
                var temperatureSetpointShiftEco = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftEco);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16V16V16V16RoomTemperature") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (temperatureSetpointComfort)
                var temperatureSetpointComfort = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointComfort);

                // Simple Field (temperatureSetpointStandby)
                var temperatureSetpointStandby = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointStandby);

                // Simple Field (temperatureSetpointEco)
                var temperatureSetpointEco = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointEco);

                // Simple Field (temperatureSetpointBProt)
                var temperatureSetpointBProt = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointBProt);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16V16V16V16DhwtTemperature") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (temperatureSetpointLegioProtect)
                var temperatureSetpointLegioProtect = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointLegioProtect);

                // Simple Field (temperatureSetpointNormal)
                var temperatureSetpointNormal = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointNormal);

                // Simple Field (temperatureSetpointReduced)
                var temperatureSetpointReduced = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointReduced);

                // Simple Field (temperatureSetpointFrostProtect)
                var temperatureSetpointFrostProtect = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointFrostProtect);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16V16V16V16RoomTemperatureShift") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (temperatureSetpointShiftComfort)
                var temperatureSetpointShiftComfort = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftComfort);

                // Simple Field (temperatureSetpointShiftStandby)
                var temperatureSetpointShiftStandby = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftStandby);

                // Simple Field (temperatureSetpointShiftEco)
                var temperatureSetpointShiftEco = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftEco);

                // Simple Field (temperatureSetpointShiftBProt)
                var temperatureSetpointShiftBProt = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftBProt);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16U8B8Heat") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (flowTemperatureDemand)
                var flowTemperatureDemand = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(flowTemperatureDemand);

                // Simple Field (relativePower)
                var relativePower = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relativePower);

                // Reserved Field (Just skip the bits)
                io.ReadByte(2);

                // Simple Field (boilerEnabled)
                var boilerEnabled = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(boilerEnabled);

                // Simple Field (stage2Forced)
                var stage2Forced = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage2Forced);

                // Simple Field (stage2Enabled)
                var stage2Enabled = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage2Enabled);

                // Simple Field (stage1Forced)
                var stage1Forced = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage1Forced);

                // Simple Field (stage1Enabled)
                var stage1Enabled = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage1Enabled);

                // Simple Field (flowTemperatureDemandValid)
                var flowTemperatureDemandValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(flowTemperatureDemandValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16U8B8ChilledWater") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (chilledWaterFlowTemperatureDemand)
                var chilledWaterFlowTemperatureDemand = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(chilledWaterFlowTemperatureDemand);

                // Simple Field (relativePower)
                var relativePower = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relativePower);

                // Reserved Field (Just skip the bits)
                io.ReadByte(5);

                // Simple Field (chilledWaterPumpEnabled)
                var chilledWaterPumpEnabled = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(chilledWaterPumpEnabled);

                // Simple Field (relativePowerValid)
                var relativePowerValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(relativePowerValid);

                // Simple Field (chilledWaterFlowTemperatureDemandValid)
                var chilledWaterFlowTemperatureDemandValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(chilledWaterFlowTemperatureDemandValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16U8B16Boiler") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (tempBoiler)
                var tempBoiler = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(tempBoiler);

                // Simple Field (relativePower)
                var relativePower = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relativePower);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (chimneySweepFunctionActive)
                var chimneySweepFunctionActive = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(chimneySweepFunctionActive);

                // Simple Field (reducedAvailability)
                var reducedAvailability = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(reducedAvailability);

                // Simple Field (powerLimitBoilerReached)
                var powerLimitBoilerReached = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(powerLimitBoilerReached);

                // Simple Field (powerLimitStage1Reached)
                var powerLimitStage1Reached = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(powerLimitStage1Reached);

                // Simple Field (stage2Enabled)
                var stage2Enabled = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage2Enabled);

                // Simple Field (stage1Enabled)
                var stage1Enabled = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(stage1Enabled);

                // Simple Field (boilerTemporarilyNotProvidingHeat)
                var boilerTemporarilyNotProvidingHeat = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(boilerTemporarilyNotProvidingHeat);

                // Simple Field (permanentlyOff)
                var permanentlyOff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(permanentlyOff);

                // Simple Field (boilerSwitchedOffWinterSummerMode)
                var boilerSwitchedOffWinterSummerMode = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(boilerSwitchedOffWinterSummerMode);

                // Simple Field (boilerFailure)
                var boilerFailure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(boilerFailure);

                // Simple Field (relativePowerValid)
                var relativePowerValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(relativePowerValid);

                // Simple Field (tempBoilerValid)
                var tempBoilerValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(tempBoilerValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16U8B16Chiller") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (tempChiller)
                var tempChiller = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(tempChiller);

                // Simple Field (relativePower)
                var relativePower = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relativePower);

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (reducedAvailability)
                var reducedAvailability = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(reducedAvailability);

                // Simple Field (powerLimitChillerReached)
                var powerLimitChillerReached = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(powerLimitChillerReached);

                // Simple Field (powerLimitCurrentStageReached)
                var powerLimitCurrentStageReached = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(powerLimitCurrentStageReached);

                // Simple Field (permanentlyOff)
                var permanentlyOff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(permanentlyOff);

                // Simple Field (chillerFailure)
                var chillerFailure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(chillerFailure);

                // Simple Field (chillerRunningStatus)
                var chillerRunningStatus = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(chillerRunningStatus);

                // Simple Field (relativePowerValid)
                var relativePowerValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(relativePowerValid);

                // Simple Field (tempChillerValid)
                var tempChillerValid = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(tempChillerValid);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16U8N8B8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (nominalPower)
                var nominalPower = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(nominalPower);

                // Simple Field (relativePowerLimit)
                var relativePowerLimit = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(relativePowerLimit);

                // Simple Field (burnerType)
                var burnerType = io.ReadSbyte(8);
                internalMap["Struct"] = new PlcSINT(burnerType);

                // Reserved Field (Just skip the bits)
                io.ReadByte(5);

                // Simple Field (solidState)
                var solidState = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(solidState);

                // Simple Field (gas)
                var gas = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(gas);

                // Simple Field (oil)
                var oil = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(oil);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U5U5U6") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (magicNumber)
                var magicNumber = io.ReadByte(5);
                internalMap["Struct"] = new PlcUSINT(magicNumber);

                // Simple Field (versionNumber)
                var versionNumber = io.ReadByte(5);
                internalMap["Struct"] = new PlcUSINT(versionNumber);

                // Simple Field (revisionNumber)
                var revisionNumber = io.ReadByte(6);
                internalMap["Struct"] = new PlcUSINT(revisionNumber);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V32Z8VolumeLiter") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (volumeLiter)
                var volumeLiter = io.ReadInt(32);
                internalMap["Struct"] = new PlcDINT(volumeLiter);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V32Z8FlowRate") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (flowRate)
                var flowRate = io.ReadInt(32);
                internalMap["Struct"] = new PlcDINT(flowRate);

                // Simple Field (statusCommand)
                var statusCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8N8N8N8B8B8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (logNumber)
                var logNumber = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(logNumber);

                // Simple Field (alarmPriority)
                var alarmPriority = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(alarmPriority);

                // Simple Field (applicationArea)
                var applicationArea = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(applicationArea);

                // Simple Field (errorClass)
                var errorClass = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(errorClass);

                // Reserved Field (Just skip the bits)
                io.ReadByte(4);

                // Simple Field (errorCode_Sup)
                var errorCode_Sup = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(errorCode_Sup);

                // Simple Field (alarmText_Sup)
                var alarmText_Sup = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(alarmText_Sup);

                // Simple Field (timeStamp_Sup)
                var timeStamp_Sup = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(timeStamp_Sup);

                // Simple Field (ack_Sup)
                var ack_Sup = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(ack_Sup);

                // Reserved Field (Just skip the bits)
                io.ReadByte(5);

                // Simple Field (alarmUnAck)
                var alarmUnAck = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(alarmUnAck);

                // Simple Field (locked)
                var locked = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(locked);

                // Simple Field (inAlarm)
                var inAlarm = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(inAlarm);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16V16") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (delayTime)
                var delayTime = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(delayTime);

                // Simple Field (temperature)
                var temperature = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(temperature);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "N16U32") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (manufacturerCode)
                var manufacturerCode = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(manufacturerCode);

                // Simple Field (incrementedNumber)
                var incrementedNumber = io.ReadUint(32);
                internalMap["Struct"] = new PlcUDINT(incrementedNumber);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "F16F16F16") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (temperatureSetpointComfort)
                var temperatureSetpointComfort = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointComfort);

                // Simple Field (temperatureSetpointShiftStandby)
                var temperatureSetpointShiftStandby = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftStandby);

                // Simple Field (temperatureSetpointShiftEco)
                var temperatureSetpointShiftEco = io.ReadFloat(true, 4, 11);
                internalMap["Struct"] = new PlcREAL(temperatureSetpointShiftEco);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V8N8N8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (energyDemand)
                var energyDemand = io.ReadSbyte(8);
                internalMap["Struct"] = new PlcSINT(energyDemand);

                // Simple Field (hvacControllerMode)
                var hvacControllerMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacControllerMode);

                // Simple Field (hvacEmergencyMode)
                var hvacEmergencyMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacEmergencyMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V16V16N8N8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (tempSetpointCooling)
                var tempSetpointCooling = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(tempSetpointCooling);

                // Simple Field (tempSetpointHeating)
                var tempSetpointHeating = io.ReadShort(16);
                internalMap["Struct"] = new PlcINT(tempSetpointHeating);

                // Simple Field (hvacControllerMode)
                var hvacControllerMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacControllerMode);

                // Simple Field (hvacEmergencyMode)
                var hvacEmergencyMode = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(hvacEmergencyMode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16U8Scaling") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (timePeriod)
                var timePeriod = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(timePeriod);

                // Simple Field (percent)
                var percent = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(percent);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16U8TariffNext") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (delayTime)
                var delayTime = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(delayTime);

                // Simple Field (tariff)
                var tariff = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(tariff);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V32N8Z8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (countVal)
                var countVal = io.ReadInt(32);
                internalMap["Struct"] = new PlcDINT(countVal);

                // Simple Field (valInfField)
                var valInfField = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(valInfField);

                // Simple Field (statusOrCommand)
                var statusOrCommand = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(statusOrCommand);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U16U32U8N8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (manufacturerId)
                var manufacturerId = io.ReadUshort(16);
                internalMap["Struct"] = new PlcUINT(manufacturerId);

                // Simple Field (identNumber)
                var identNumber = io.ReadUint(32);
                internalMap["Struct"] = new PlcUDINT(identNumber);

                // Simple Field (version)
                var version = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(version);

                // Simple Field (medium)
                var medium = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(medium);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "A8A8A8A8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (languageCode)
                var languageCode = io.ReadString(16, Encoding.GetEncoding("ASCII"));
                internalMap["Struct"] = new PlcSTRING(languageCode);

                // Simple Field (regionCode)
                var regionCode = io.ReadString(16, Encoding.GetEncoding("ASCII"));
                internalMap["Struct"] = new PlcSTRING(regionCode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8U8U8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (red)
                var red = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(red);

                // Simple Field (green)
                var green = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(green);

                // Simple Field (blue)
                var blue = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(blue);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "A8A8Language") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (languageCode)
                var languageCode = io.ReadString(16, Encoding.GetEncoding("ASCII"));
                internalMap["Struct"] = new PlcSTRING(languageCode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "A8A8Region") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (regionCode)
                var regionCode = io.ReadString(16, Encoding.GetEncoding("ASCII"));
                internalMap["Struct"] = new PlcSTRING(regionCode);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "V32U8B8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (activeElectricalEnergy)
                var activeElectricalEnergy = io.ReadInt(32);
                internalMap["Struct"] = new PlcDINT(activeElectricalEnergy);

                // Simple Field (tariff)
                var tariff = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(tariff);

                // Reserved Field (Just skip the bits)
                io.ReadByte(6);

                // Simple Field (noTariff)
                var noTariff = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(noTariff);

                // Simple Field (noActiveElectricalEnergy)
                var noActiveElectricalEnergy = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(noActiveElectricalEnergy);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "B1N3N4") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (deactivationOfPriority)
                var deactivationOfPriority = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(deactivationOfPriority);

                // Simple Field (priorityLevel)
                var priorityLevel = io.ReadByte(3);
                internalMap["Struct"] = new PlcUSINT(priorityLevel);

                // Simple Field (modeLevel)
                var modeLevel = io.ReadByte(4);
                internalMap["Struct"] = new PlcUSINT(modeLevel);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "B10U6") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(5);

                // Simple Field (convertorError)
                var convertorError = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(convertorError);

                // Simple Field (ballastFailure)
                var ballastFailure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(ballastFailure);

                // Simple Field (lampError)
                var lampError = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(lampError);

                // Simple Field (read)
                var read = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(read);

                // Simple Field (groupAddress)
                var groupAddress = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(groupAddress);

                // Simple Field (address)
                var address = io.ReadByte(6);
                internalMap["Struct"] = new PlcUSINT(address);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "B2U6") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (sceneActivationInactive)
                var sceneActivationInactive = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(sceneActivationInactive);

                // Simple Field (storageFunctionDisable)
                var storageFunctionDisable = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(storageFunctionDisable);

                // Simple Field (sceneNumber)
                var sceneNumber = io.ReadByte(6);
                internalMap["Struct"] = new PlcUSINT(sceneNumber);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8r7B1") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (setValue)
                var setValue = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(setValue);

                // Reserved Field (Just skip the bits)
                io.ReadByte(7);

                // Simple Field (channelActivationActive)
                var channelActivationActive = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(channelActivationActive);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8U8B8") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (heightPosition)
                var heightPosition = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(heightPosition);

                // Simple Field (slatsPosition)
                var slatsPosition = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(slatsPosition);

                // Reserved Field (Just skip the bits)
                io.ReadByte(6);

                // Simple Field (validSlatsPos)
                var validSlatsPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(validSlatsPos);

                // Simple Field (validHeightPos)
                var validHeightPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(validHeightPos);

                return new PlcStruct(internalMap);
            } else
            if (formatName == "U8U8B16") { // Struct
                var internalMap = new Dictionary<string, IPlcValue>();

                // Reserved Field (Just skip the bits)
                io.ReadByte(8);

                // Simple Field (heightPosition)
                var heightPosition = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(heightPosition);

                // Simple Field (slatsPosition)
                var slatsPosition = io.ReadByte(8);
                internalMap["Struct"] = new PlcUSINT(slatsPosition);

                // Simple Field (validSlatsPos)
                var validSlatsPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(validSlatsPos);

                // Simple Field (validHeightPos)
                var validHeightPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(validHeightPos);

                // Reserved Field (Just skip the bits)
                io.ReadByte(3);

                // Simple Field (failure)
                var failure = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(failure);

                // Simple Field (localOverride)
                var localOverride = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(localOverride);

                // Simple Field (locked)
                var locked = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(locked);

                // Simple Field (forced)
                var forced = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(forced);

                // Simple Field (weatherAlarm)
                var weatherAlarm = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(weatherAlarm);

                // Simple Field (targetSPosRestrict)
                var targetSPosRestrict = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(targetSPosRestrict);

                // Simple Field (targetHPosRestrict)
                var targetHPosRestrict = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(targetHPosRestrict);

                // Simple Field (driveState)
                var driveState = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(driveState);

                // Simple Field (lowerPredefPos)
                var lowerPredefPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(lowerPredefPos);

                // Simple Field (lowerEndPos)
                var lowerEndPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(lowerEndPos);

                // Simple Field (upperEndPos)
                var upperEndPos = io.ReadBit();
                internalMap["Struct"] = new PlcBOOL(upperEndPos);

                return new PlcStruct(internalMap);
            } 
            return null;
        }

        public void Serialize(WriteBuffer io, IPlcValue value, string formatName)
        {
            if (formatName == "B1") { // BOOL

                // Reserved Field (Just skip the bytes)
                io.WriteByte(7, 0x0);

                // Simple Field (value)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "B2") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(6, 0x0);

                // Simple Field (control)
                io.WriteBit(value.GetBool());

                // Simple Field (value)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "B1U3") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x0);

                // Simple Field (control)
                io.WriteBit(value.GetBool());

                // Simple Field (value)
                io.WriteByte(3, value.GetByte());
            } else
            if (formatName == "A8_ASCII") { // STRING

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteString(8, "ASCII", value.GetString());
            } else
            if (formatName == "A8_8859_1") { // STRING

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteString(8, "ISO-8859-1", value.GetString());
            } else
            if (formatName == "U8") { // USINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V8") { // SINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteSbyte(8, value.GetSbyte());
            } else
            if (formatName == "B5N3") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(3, 0x0);

                // Simple Field (a)
                io.WriteBit(value.GetBool());

                // Simple Field (b)
                io.WriteBit(value.GetBool());

                // Simple Field (c)
                io.WriteBit(value.GetBool());

                // Simple Field (d)
                io.WriteBit(value.GetBool());

                // Simple Field (e)
                io.WriteBit(value.GetBool());

                // Simple Field (value)
                io.WriteSbyte(8, value.GetSbyte());
            } else
            if (formatName == "U16") { // UINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteUshort(16, value.GetUshort());
            } else
            if (formatName == "V16") { // INT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteShort(16, value.GetShort());
            } else
            if (formatName == "F16") { // REAL

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "N3N5r2N6r2N6") { // TIME_OF_DAY

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (day)
                io.WriteByte(3, value.GetByte());

                // Simple Field (hour)
                io.WriteByte(5, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(2, 0x00);

                // Simple Field (minutes)
                io.WriteByte(6, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(2, 0x00);

                // Simple Field (seconds)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "r3N5r4N4r1U7") { // DATE

                // Reserved Field (Just skip the bytes)
                io.WriteByte(3, 0x00);

                // Simple Field (day)
                io.WriteByte(5, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (month)
                io.WriteByte(4, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(1, 0x00);

                // Simple Field (year)
                io.WriteByte(7, value.GetByte());
            } else
            if (formatName == "U32") { // UDINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteUint(32, value.GetUint());
            } else
            if (formatName == "V32") { // DINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteInt(32, value.GetInt());
            } else
            if (formatName == "F32") { // REAL

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteFloat(32, value.GetFloat());
            } else
            if (formatName == "U4U4U4U4U4U4B4N4") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (d6)
                io.WriteByte(4, value.GetByte());

                // Simple Field (d5)
                io.WriteByte(4, value.GetByte());

                // Simple Field (d4)
                io.WriteByte(4, value.GetByte());

                // Simple Field (d3)
                io.WriteByte(4, value.GetByte());

                // Simple Field (d2)
                io.WriteByte(4, value.GetByte());

                // Simple Field (d1)
                io.WriteByte(4, value.GetByte());

                // Simple Field (e)
                io.WriteBit(value.GetBool());

                // Simple Field (p)
                io.WriteBit(value.GetBool());

                // Simple Field (d)
                io.WriteBit(value.GetBool());

                // Simple Field (c)
                io.WriteBit(value.GetBool());

                // Simple Field (index)
                io.WriteByte(4, value.GetByte());
            } else
            if (formatName == "A112_ASCII") { // STRING

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteString(112, "ASCII", value.GetString());
            } else
            if (formatName == "A112_8859_1") { // STRING

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteString(112, "ISO-8859-1", value.GetString());
            } else
            if (formatName == "r2U6") { // USINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(2, 0x00);

                // Simple Field (value)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "B1r1U6") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (learn)
                io.WriteBit(value.GetBool());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(1, 0x00);

                // Simple Field (sceneNumber)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "U8r4U4r3U5U3U5r2U6r2U6B16") { // DATE_AND_TIME

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (year)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (month)
                io.WriteByte(4, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(3, 0x00);

                // Simple Field (day)
                io.WriteByte(5, value.GetByte());

                // Simple Field (dayOfWeek)
                io.WriteByte(3, value.GetByte());

                // Simple Field (hour)
                io.WriteByte(5, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(2, 0x00);

                // Simple Field (minutes)
                io.WriteByte(6, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(2, 0x00);

                // Simple Field (seconds)
                io.WriteByte(6, value.GetByte());

                // Simple Field (fault)
                io.WriteBit(value.GetBool());

                // Simple Field (workingDay)
                io.WriteBit(value.GetBool());

                // Simple Field (noWorkingDay)
                io.WriteBit(value.GetBool());

                // Simple Field (noYear)
                io.WriteBit(value.GetBool());

                // Simple Field (noMonthAndDay)
                io.WriteBit(value.GetBool());

                // Simple Field (noDayOfWeek)
                io.WriteBit(value.GetBool());

                // Simple Field (noTime)
                io.WriteBit(value.GetBool());

                // Simple Field (standardSummerTime)
                io.WriteBit(value.GetBool());

                // Simple Field (clockWithSyncSignal)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "N8") { // USINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "B8") { // BYTE

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "B16") { // WORD

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteUshort(16, value.GetUshort());
            } else
            if (formatName == "U4U4") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (busy)
                io.WriteByte(4, value.GetByte());

                // Simple Field (nak)
                io.WriteByte(4, value.GetByte());
            } else
            if (formatName == "r1b1U6") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Reserved Field (Just skip the bytes)
                io.WriteByte(1, 0x00);

                // Simple Field (sceneActive)
                io.WriteBit(value.GetBool());

                // Simple Field (sceneNumber)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "B32") { // DWORD

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteUint(32, value.GetUint());
            } else
            if (formatName == "V64") { // LINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (value)
                io.WriteLong(64, value.GetLong());
            } else
            if (formatName == "B24") { // List

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Array Field (value)
                for (int i = 0; i < (24); i++) {
                    io.WriteBit(value.GetIndex(i).GetBool());
                }
            } else
            if (formatName == "N3") { // USINT

                // Reserved Field (Just skip the bytes)
                io.WriteByte(5, 0x00);

                // Simple Field (value)
                io.WriteByte(3, value.GetByte());
            } else
            if (formatName == "B1Z8HeatingOrCoolingZ") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(7, 0x00);

                // Simple Field (heating)
                io.WriteBit(value.GetBool());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "B1Z8BinaryValueZ") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(7, 0x00);

                // Simple Field (high)
                io.WriteBit(value.GetBool());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8HvacOperatingMode") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (hvacOperatingMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8DhwMode") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (dhwMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8HvacControllingMode") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (hvacControllingMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8EnableHeatingOrCoolingStage") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (enableHeatingOrCoolingStage)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8BuildingMode") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (buildingMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8OccupancyMode") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (occupancyMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "N8Z8EmergencyMode") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (hvacEmergencyMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U8Z8Rel") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (relValue)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U8Z8Counter") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (counterValue)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8TimePeriod") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (timePeriod)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8FlowRate") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (flowRate)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8Counter") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (counterValue)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8ElectricCurrent") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (electricalCurrent)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8Power") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (power)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8AtmPressure") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (atmPressure)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8PercentValue") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (percentValue)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8HvacAirQuality") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (ppmResolution)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8WindSpeed") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (windSpeed)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8SunIntensity") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (sunIntensity)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16Z8HvacAirFlow") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (airFlow)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V8Z8RelSignedValue") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (relSignedValue)
                io.WriteSbyte(8, value.GetSbyte());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V16Z8DeltaTime") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (deltaTime)
                io.WriteShort(16, value.GetShort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V16Z8RelSignedValue") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (relSignedValue)
                io.WriteShort(16, value.GetShort());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16N8HvacModeAndTimeDelay") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (delayTime)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (hvacMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16N8DhwModeAndTimeDelay") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (delayTime)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (dhwMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16N8OccupancyModeAndTimeDelay") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (delayTime)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (occupationMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16N8BuildingModeAndTimeDelay") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (delayTime)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (buildingMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U8B8StatusBurnerController") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (actualRelativePower)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (stage2Active)
                io.WriteBit(value.GetBool());

                // Simple Field (stage1Active)
                io.WriteBit(value.GetBool());

                // Simple Field (failure)
                io.WriteBit(value.GetBool());

                // Simple Field (actualRelativePowerValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8B8LockingSignal") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (requestedPowerReduction)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(6, 0x00);

                // Simple Field (critical)
                io.WriteBit(value.GetBool());

                // Simple Field (requestedPowerReductionValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8B8BoilerControllerDemandSignal") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (relativeDemand)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(6, 0x00);

                // Simple Field (controlsOperationStage2)
                io.WriteBit(value.GetBool());

                // Simple Field (controlsOperationStage1)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8B8ActuatorPositionDemand") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (actuatorPositionDemand)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (emergencyDemand)
                io.WriteBit(value.GetBool());

                // Simple Field (shiftLoadPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (absoluteLoadPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (actuatorPositionDemandValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8B8ActuatorPositionStatus") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (actualActuatorPosition)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(3, 0x00);

                // Simple Field (synchronizationMode)
                io.WriteBit(value.GetBool());

                // Simple Field (valveKick)
                io.WriteBit(value.GetBool());

                // Simple Field (callibrationMode)
                io.WriteBit(value.GetBool());

                // Simple Field (positionManuallyOverridden)
                io.WriteBit(value.GetBool());

                // Simple Field (failure)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8B8StatusLightingActuator") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (lightingLevel)
                io.WriteByte(8, value.GetByte());

                // Simple Field (failure)
                io.WriteBit(value.GetBool());

                // Simple Field (localOverride)
                io.WriteBit(value.GetBool());

                // Simple Field (dimming)
                io.WriteBit(value.GetBool());

                // Simple Field (staircaseLightingFunction)
                io.WriteBit(value.GetBool());

                // Simple Field (nightMode)
                io.WriteBit(value.GetBool());

                // Simple Field (forced)
                io.WriteBit(value.GetBool());

                // Simple Field (locked)
                io.WriteBit(value.GetBool());

                // Simple Field (lightingLevelValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16B8HeatProducerManagerStatus") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (tempFlowProdSegmH)
                io.WriteFloat(16, value.GetFloat());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(3, 0x00);

                // Simple Field (temporarilyOff)
                io.WriteBit(value.GetBool());

                // Simple Field (permanentlyOff)
                io.WriteBit(value.GetBool());

                // Simple Field (switchedOffSummerMode)
                io.WriteBit(value.GetBool());

                // Simple Field (failure)
                io.WriteBit(value.GetBool());

                // Simple Field (tempFlowProdSegmHValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16B8RoomTemperatureDemand") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (roomTemperatureDemand)
                io.WriteFloat(16, value.GetFloat());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (emergencyDemand)
                io.WriteBit(value.GetBool());

                // Simple Field (shiftLoadPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (absoluteLoadPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (roomTemperatureDemandValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16B8ColdWaterProducerManagerStatus") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (flowTemperatureProdSegmC)
                io.WriteFloat(16, value.GetFloat());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (temporarilyOff)
                io.WriteBit(value.GetBool());

                // Simple Field (permanentlyOff)
                io.WriteBit(value.GetBool());

                // Simple Field (failure)
                io.WriteBit(value.GetBool());

                // Simple Field (flowTemperatureProdSegmCValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16B8WaterTemperatureControllerStatus") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (actualTemperature)
                io.WriteFloat(16, value.GetFloat());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(5, 0x00);

                // Simple Field (controllerWorking)
                io.WriteBit(value.GetBool());

                // Simple Field (failure)
                io.WriteBit(value.GetBool());

                // Simple Field (actualTemperatureValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16B16") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (flowTemperatureDemand)
                io.WriteFloat(16, value.GetFloat());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (demandFromDhwWhileLegionellaFunctionIsActive)
                io.WriteBit(value.GetBool());

                // Simple Field (emergencyDemandForFrostProtection)
                io.WriteBit(value.GetBool());

                // Simple Field (requestForWaterCirculationInPrimaryDistributionSegment)
                io.WriteBit(value.GetBool());

                // Simple Field (demandFromAuxillaryHeatOrCoolConsumer)
                io.WriteBit(value.GetBool());

                // Simple Field (demandFromVentilation)
                io.WriteBit(value.GetBool());

                // Simple Field (demandForRoomHeatingOrCooling)
                io.WriteBit(value.GetBool());

                // Simple Field (heatDemandFromDhw)
                io.WriteBit(value.GetBool());

                // Simple Field (flowTemperatureDemandIsMin)
                io.WriteBit(value.GetBool());

                // Simple Field (flowTemperatureDemandIsMax)
                io.WriteBit(value.GetBool());

                // Simple Field (shiftLoadPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (absoluteLoadPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (flowTemperatureDemandValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8N8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (energyDemand)
                io.WriteByte(8, value.GetByte());

                // Simple Field (actualControllerMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V16V16V16RoomTemperature") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (temperatureSetpointComfort)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointStandby)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointEco)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "V16V16V16RoomTemperatureShift") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (temperatureSetpointShiftComfort)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftStandby)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftEco)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "V16V16V16V16RoomTemperature") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (temperatureSetpointComfort)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointStandby)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointEco)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointBProt)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "V16V16V16V16DhwtTemperature") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (temperatureSetpointLegioProtect)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointNormal)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointReduced)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointFrostProtect)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "V16V16V16V16RoomTemperatureShift") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (temperatureSetpointShiftComfort)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftStandby)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftEco)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftBProt)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "V16U8B8Heat") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (flowTemperatureDemand)
                io.WriteShort(16, value.GetShort());

                // Simple Field (relativePower)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(2, 0x00);

                // Simple Field (boilerEnabled)
                io.WriteBit(value.GetBool());

                // Simple Field (stage2Forced)
                io.WriteBit(value.GetBool());

                // Simple Field (stage2Enabled)
                io.WriteBit(value.GetBool());

                // Simple Field (stage1Forced)
                io.WriteBit(value.GetBool());

                // Simple Field (stage1Enabled)
                io.WriteBit(value.GetBool());

                // Simple Field (flowTemperatureDemandValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16U8B8ChilledWater") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (chilledWaterFlowTemperatureDemand)
                io.WriteShort(16, value.GetShort());

                // Simple Field (relativePower)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(5, 0x00);

                // Simple Field (chilledWaterPumpEnabled)
                io.WriteBit(value.GetBool());

                // Simple Field (relativePowerValid)
                io.WriteBit(value.GetBool());

                // Simple Field (chilledWaterFlowTemperatureDemandValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16U8B16Boiler") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (tempBoiler)
                io.WriteShort(16, value.GetShort());

                // Simple Field (relativePower)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (chimneySweepFunctionActive)
                io.WriteBit(value.GetBool());

                // Simple Field (reducedAvailability)
                io.WriteBit(value.GetBool());

                // Simple Field (powerLimitBoilerReached)
                io.WriteBit(value.GetBool());

                // Simple Field (powerLimitStage1Reached)
                io.WriteBit(value.GetBool());

                // Simple Field (stage2Enabled)
                io.WriteBit(value.GetBool());

                // Simple Field (stage1Enabled)
                io.WriteBit(value.GetBool());

                // Simple Field (boilerTemporarilyNotProvidingHeat)
                io.WriteBit(value.GetBool());

                // Simple Field (permanentlyOff)
                io.WriteBit(value.GetBool());

                // Simple Field (boilerSwitchedOffWinterSummerMode)
                io.WriteBit(value.GetBool());

                // Simple Field (boilerFailure)
                io.WriteBit(value.GetBool());

                // Simple Field (relativePowerValid)
                io.WriteBit(value.GetBool());

                // Simple Field (tempBoilerValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "V16U8B16Chiller") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (tempChiller)
                io.WriteShort(16, value.GetShort());

                // Simple Field (relativePower)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x00);

                // Simple Field (reducedAvailability)
                io.WriteBit(value.GetBool());

                // Simple Field (powerLimitChillerReached)
                io.WriteBit(value.GetBool());

                // Simple Field (powerLimitCurrentStageReached)
                io.WriteBit(value.GetBool());

                // Simple Field (permanentlyOff)
                io.WriteBit(value.GetBool());

                // Simple Field (chillerFailure)
                io.WriteBit(value.GetBool());

                // Simple Field (chillerRunningStatus)
                io.WriteBit(value.GetBool());

                // Simple Field (relativePowerValid)
                io.WriteBit(value.GetBool());

                // Simple Field (tempChillerValid)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U16U8N8B8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (nominalPower)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (relativePowerLimit)
                io.WriteByte(8, value.GetByte());

                // Simple Field (burnerType)
                io.WriteSbyte(8, value.GetSbyte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(5, 0x00);

                // Simple Field (solidState)
                io.WriteBit(value.GetBool());

                // Simple Field (gas)
                io.WriteBit(value.GetBool());

                // Simple Field (oil)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U5U5U6") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (magicNumber)
                io.WriteByte(5, value.GetByte());

                // Simple Field (versionNumber)
                io.WriteByte(5, value.GetByte());

                // Simple Field (revisionNumber)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "V32Z8VolumeLiter") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (volumeLiter)
                io.WriteInt(32, value.GetInt());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V32Z8FlowRate") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (flowRate)
                io.WriteInt(32, value.GetInt());

                // Simple Field (statusCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U8N8N8N8B8B8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (logNumber)
                io.WriteByte(8, value.GetByte());

                // Simple Field (alarmPriority)
                io.WriteByte(8, value.GetByte());

                // Simple Field (applicationArea)
                io.WriteByte(8, value.GetByte());

                // Simple Field (errorClass)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(4, 0x00);

                // Simple Field (errorCode_Sup)
                io.WriteBit(value.GetBool());

                // Simple Field (alarmText_Sup)
                io.WriteBit(value.GetBool());

                // Simple Field (timeStamp_Sup)
                io.WriteBit(value.GetBool());

                // Simple Field (ack_Sup)
                io.WriteBit(value.GetBool());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(5, 0x00);

                // Simple Field (alarmUnAck)
                io.WriteBit(value.GetBool());

                // Simple Field (locked)
                io.WriteBit(value.GetBool());

                // Simple Field (inAlarm)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U16V16") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (delayTime)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (temperature)
                io.WriteShort(16, value.GetShort());
            } else
            if (formatName == "N16U32") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (manufacturerCode)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (incrementedNumber)
                io.WriteUint(32, value.GetUint());
            } else
            if (formatName == "F16F16F16") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (temperatureSetpointComfort)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftStandby)
                io.WriteFloat(16, value.GetFloat());

                // Simple Field (temperatureSetpointShiftEco)
                io.WriteFloat(16, value.GetFloat());
            } else
            if (formatName == "V8N8N8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (energyDemand)
                io.WriteSbyte(8, value.GetSbyte());

                // Simple Field (hvacControllerMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (hvacEmergencyMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V16V16N8N8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (tempSetpointCooling)
                io.WriteShort(16, value.GetShort());

                // Simple Field (tempSetpointHeating)
                io.WriteShort(16, value.GetShort());

                // Simple Field (hvacControllerMode)
                io.WriteByte(8, value.GetByte());

                // Simple Field (hvacEmergencyMode)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16U8Scaling") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (timePeriod)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (percent)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16U8TariffNext") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (delayTime)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (tariff)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "V32N8Z8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (countVal)
                io.WriteInt(32, value.GetInt());

                // Simple Field (valInfField)
                io.WriteByte(8, value.GetByte());

                // Simple Field (statusOrCommand)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "U16U32U8N8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (manufacturerId)
                io.WriteUshort(16, value.GetUshort());

                // Simple Field (identNumber)
                io.WriteUint(32, value.GetUint());

                // Simple Field (version)
                io.WriteByte(8, value.GetByte());

                // Simple Field (medium)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "A8A8A8A8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (languageCode)
                io.WriteString(16, "ASCII", value.GetString());

                // Simple Field (regionCode)
                io.WriteString(16, "ASCII", value.GetString());
            } else
            if (formatName == "U8U8U8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (red)
                io.WriteByte(8, value.GetByte());

                // Simple Field (green)
                io.WriteByte(8, value.GetByte());

                // Simple Field (blue)
                io.WriteByte(8, value.GetByte());
            } else
            if (formatName == "A8A8Language") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (languageCode)
                io.WriteString(16, "ASCII", value.GetString());
            } else
            if (formatName == "A8A8Region") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (regionCode)
                io.WriteString(16, "ASCII", value.GetString());
            } else
            if (formatName == "V32U8B8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (activeElectricalEnergy)
                io.WriteInt(32, value.GetInt());

                // Simple Field (tariff)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(6, 0x00);

                // Simple Field (noTariff)
                io.WriteBit(value.GetBool());

                // Simple Field (noActiveElectricalEnergy)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "B1N3N4") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (deactivationOfPriority)
                io.WriteBit(value.GetBool());

                // Simple Field (priorityLevel)
                io.WriteByte(3, value.GetByte());

                // Simple Field (modeLevel)
                io.WriteByte(4, value.GetByte());
            } else
            if (formatName == "B10U6") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(5, 0x00);

                // Simple Field (convertorError)
                io.WriteBit(value.GetBool());

                // Simple Field (ballastFailure)
                io.WriteBit(value.GetBool());

                // Simple Field (lampError)
                io.WriteBit(value.GetBool());

                // Simple Field (read)
                io.WriteBit(value.GetBool());

                // Simple Field (groupAddress)
                io.WriteBit(value.GetBool());

                // Simple Field (address)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "B2U6") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (sceneActivationInactive)
                io.WriteBit(value.GetBool());

                // Simple Field (storageFunctionDisable)
                io.WriteBit(value.GetBool());

                // Simple Field (sceneNumber)
                io.WriteByte(6, value.GetByte());
            } else
            if (formatName == "U8r7B1") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (setValue)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(7, 0x00);

                // Simple Field (channelActivationActive)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8U8B8") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (heightPosition)
                io.WriteByte(8, value.GetByte());

                // Simple Field (slatsPosition)
                io.WriteByte(8, value.GetByte());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(6, 0x00);

                // Simple Field (validSlatsPos)
                io.WriteBit(value.GetBool());

                // Simple Field (validHeightPos)
                io.WriteBit(value.GetBool());
            } else
            if (formatName == "U8U8B16") { // Struct

                // Reserved Field (Just skip the bytes)
                io.WriteByte(8, 0x0);

                // Simple Field (heightPosition)
                io.WriteByte(8, value.GetByte());

                // Simple Field (slatsPosition)
                io.WriteByte(8, value.GetByte());

                // Simple Field (validSlatsPos)
                io.WriteBit(value.GetBool());

                // Simple Field (validHeightPos)
                io.WriteBit(value.GetBool());

                // Reserved Field (Just skip the bytes)
                io.WriteByte(3, 0x00);

                // Simple Field (failure)
                io.WriteBit(value.GetBool());

                // Simple Field (localOverride)
                io.WriteBit(value.GetBool());

                // Simple Field (locked)
                io.WriteBit(value.GetBool());

                // Simple Field (forced)
                io.WriteBit(value.GetBool());

                // Simple Field (weatherAlarm)
                io.WriteBit(value.GetBool());

                // Simple Field (targetSPosRestrict)
                io.WriteBit(value.GetBool());

                // Simple Field (targetHPosRestrict)
                io.WriteBit(value.GetBool());

                // Simple Field (driveState)
                io.WriteBit(value.GetBool());

                // Simple Field (lowerPredefPos)
                io.WriteBit(value.GetBool());

                // Simple Field (lowerEndPos)
                io.WriteBit(value.GetBool());

                // Simple Field (upperEndPos)
                io.WriteBit(value.GetBool());
            }         }

    }

}

