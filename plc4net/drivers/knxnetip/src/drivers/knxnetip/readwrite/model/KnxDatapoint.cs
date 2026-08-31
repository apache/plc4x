//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Code generated from the mspec by plc4net-code-gen. DO NOT EDIT.

using System;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.model.values;

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    /// <summary>
    /// mspec <c>[dataIo KnxDatapoint]</c> - reads and writes one
    /// <see cref="IPlcValue"/> whose wire layout the parser arguments pick.
    /// </summary>
    public static class KnxDatapoint
    {
        public static IPlcValue StaticParse(ReadBuffer readBuffer, KnxDatapointType datapointType)
        {
            if (Equals(datapointType, KnxDatapointType.BOOL))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.BYTE))
            {
                // BYTE
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcBYTE(value);
            }
            else if (Equals(datapointType, KnxDatapointType.WORD))
            {
                // WORD
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcWORD(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DWORD))
            {
                // DWORD
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcDWORD(value);
            }
            else if (Equals(datapointType, KnxDatapointType.LWORD))
            {
                // LWORD
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcLWORD(value);
            }
            else if (Equals(datapointType, KnxDatapointType.USINT))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.SINT))
            {
                // SINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.UINT))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.INT))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.UDINT))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DINT))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.ULINT))
            {
                // ULINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcULINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.LINT))
            {
                // LINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.REAL))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.LREAL))
            {
                // LREAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadDouble("value", 64);
                return new PlcLREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.CHAR))
            {
                // CHAR
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 8, System.Text.Encoding.UTF8);
                return new PlcCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(datapointType, KnxDatapointType.WCHAR))
            {
                // WCHAR
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 16, System.Text.Encoding.Unicode);
                return new PlcWCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(datapointType, KnxDatapointType.TIME))
            {
                // TIME
                throw new NotImplementedException("KnxDatapoint 'TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.LTIME))
            {
                // LTIME
                throw new NotImplementedException("KnxDatapoint 'LTIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DATE))
            {
                // DATE
                throw new NotImplementedException("KnxDatapoint 'DATE' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.TIME_OF_DAY))
            {
                // TIME_OF_DAY
                throw new NotImplementedException("KnxDatapoint 'TIME_OF_DAY' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.TOD))
            {
                // TIME_OF_DAY
                throw new NotImplementedException("KnxDatapoint 'TIME_OF_DAY' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DATE_AND_TIME))
            {
                // DATE_AND_TIME
                throw new NotImplementedException("KnxDatapoint 'DATE_AND_TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DT))
            {
                // DATE_AND_TIME
                throw new NotImplementedException("KnxDatapoint 'DATE_AND_TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Switch))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Bool))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Enable))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ramp))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BinaryValue))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Step))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UpDown))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OpenClose))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Start))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_State))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Invert))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DimSendStyle))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_InputSource))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Reset))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ack))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Trigger))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Occupancy))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Window_Door))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LogicalFunction))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scene_AB))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ShutterBlinds_Mode))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DayNight))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Heat_Cool))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Switch_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var on = readBuffer.ReadBit("on");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["on"] = new PlcBOOL((bool) on);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Bool_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var valueTrue = readBuffer.ReadBit("valueTrue");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["valueTrue"] = new PlcBOOL((bool) valueTrue);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Enable_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var enable = readBuffer.ReadBit("enable");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["enable"] = new PlcBOOL((bool) enable);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ramp_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var ramp = readBuffer.ReadBit("ramp");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["ramp"] = new PlcBOOL((bool) ramp);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var alarm = readBuffer.ReadBit("alarm");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["alarm"] = new PlcBOOL((bool) alarm);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BinaryValue_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var high = readBuffer.ReadBit("high");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["high"] = new PlcBOOL((bool) high);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Step_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var increase = readBuffer.ReadBit("increase");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["increase"] = new PlcBOOL((bool) increase);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Direction1_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var down = readBuffer.ReadBit("down");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["down"] = new PlcBOOL((bool) down);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Direction2_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var close = readBuffer.ReadBit("close");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["close"] = new PlcBOOL((bool) close);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Start_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var start = readBuffer.ReadBit("start");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["start"] = new PlcBOOL((bool) start);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_State_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var active = readBuffer.ReadBit("active");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["active"] = new PlcBOOL((bool) active);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Invert_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var control = readBuffer.ReadBit("control");
                var inverted = readBuffer.ReadBit("inverted");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["control"] = new PlcBOOL((bool) control);
                _map["inverted"] = new PlcBOOL((bool) inverted);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Control_Dimming))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var increase = readBuffer.ReadBit("increase");
                var stepcode = readBuffer.ReadByte("stepcode", 3);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["increase"] = new PlcBOOL((bool) increase);
                _map["stepcode"] = new PlcUSINT((byte) stepcode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Control_Blinds))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var down = readBuffer.ReadBit("down");
                var stepcode = readBuffer.ReadByte("stepcode", 3);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["down"] = new PlcBOOL((bool) down);
                _map["stepcode"] = new PlcUSINT((byte) stepcode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Char_ASCII))
            {
                // STRING
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 8, System.Text.Encoding.ASCII);
                return new PlcSTRING(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Char_8859_1))
            {
                // STRING
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 8, System.Text.Encoding.UTF8);
                return new PlcSTRING(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Angle))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_U8))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DecimalFactor))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Tariff))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_1_Ucount))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FanStage))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_V8))
            {
                // SINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_1_Count))
            {
                // SINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Status_Mode3))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var statusA = readBuffer.ReadBit("statusA");
                var statusB = readBuffer.ReadBit("statusB");
                var statusC = readBuffer.ReadBit("statusC");
                var statusD = readBuffer.ReadBit("statusD");
                var statusE = readBuffer.ReadBit("statusE");
                var mode = readBuffer.ReadByte("mode", 3);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["statusA"] = new PlcBOOL((bool) statusA);
                _map["statusB"] = new PlcBOOL((bool) statusB);
                _map["statusC"] = new PlcBOOL((bool) statusC);
                _map["statusD"] = new PlcBOOL((bool) statusD);
                _map["statusE"] = new PlcBOOL((bool) statusE);
                _map["mode"] = new PlcUSINT((byte) mode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_2_Ucount))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodMsec))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriod10Msec))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriod100Msec))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodSec))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodMin))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodHrs))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PropDataType))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Length_mm))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UElCurrentmA))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Absolute_Colour_Temperature))
            {
                // UINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_2_Count))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeMsec))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTime10Msec))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTime100Msec))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeSec))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeMin))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeHrs))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_V16))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Rotation_Angle))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Length_m))
            {
                // INT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Temp))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Tempd))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Tempa))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Lux))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Wsp))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Pres))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Humidity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AirQuality))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AirFlow))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time1))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time2))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volt))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Curr))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PowerDensity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_KelvinPerPercent))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Power))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume_Flow))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Rain_Amount))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Temp_F))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Wsp_kmh))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Absolute_Humidity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Concentration_ygm3))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Coefficient))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimeOfDay))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var day = readBuffer.ReadByte("day", 3);
                var hour = readBuffer.ReadByte("hour", 5);
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var minutes = readBuffer.ReadByte("minutes", 6);
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var seconds = readBuffer.ReadByte("seconds", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["day"] = new PlcUSINT((byte) day);
                _map["hour"] = new PlcUSINT((byte) hour);
                _map["minutes"] = new PlcUSINT((byte) minutes);
                _map["seconds"] = new PlcUSINT((byte) seconds);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Date))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var dayOfMonth = readBuffer.ReadByte("dayOfMonth", 5);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var month = readBuffer.ReadByte("month", 4);
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var year = readBuffer.ReadByte("year", 7);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["dayOfMonth"] = new PlcUSINT((byte) dayOfMonth);
                _map["month"] = new PlcUSINT((byte) month);
                _map["year"] = new PlcUSINT((byte) year);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_4_Ucount))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Sec))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Min))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Hrs))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_VolumeLiquid_Litre))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_m_3))
            {
                // UDINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_4_Count))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FlowRate_m3h))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_kWh))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy_kVAh))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy_kVARh))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_MWh))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongDeltaTimeSec))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaVolumeLiquid_Litre))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaVolume_m_3))
            {
                // DINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Acceleration))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Acceleration_Angular))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Activation_Energy))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Activity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Mol))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Amplitude))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AngleRad))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AngleDeg))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Momentum))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Velocity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Area))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Capacitance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Charge_DensitySurface))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Charge_DensityVolume))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Compressibility))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Conductance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electrical_Conductivity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Density))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Charge))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Current))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_CurrentDensity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_DipoleMoment))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Displacement))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_FieldStrength))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Flux))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_FluxDensity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Polarization))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Potential))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_PotentialDifference))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ElectromagneticMoment))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electromotive_Force))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Energy))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Force))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Frequency))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Frequency))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_Capacity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_FlowRate))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_Quantity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Impedance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Length))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Light_Quantity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminous_Flux))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminous_Intensity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_FieldStrength))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Flux))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_FluxDensity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Moment))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Polarization))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetization))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_MagnetomotiveForce))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Mass))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_MassFlux))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Momentum))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Phase_AngleRad))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Phase_AngleDeg))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Power))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Power_Factor))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Pressure))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Reactance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Resistance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Resistivity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_SelfInductance))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_SolidAngle))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Sound_Intensity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Speed))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Stress))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Surface_Tension))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Common_Temperature))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Absolute_Temperature))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_TemperatureDifference))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Thermal_Capacity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Thermal_Conductivity))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ThermoelectricPower))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Torque))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume_Flux))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Weight))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Work))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ApparentPower))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_Flux_Meter))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_Flux_ls))
            {
                // REAL
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Access_Data))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var hurz = readBuffer.ReadByte("hurz", 4);
                var value1 = readBuffer.ReadByte("value1", 4);
                var value2 = readBuffer.ReadByte("value2", 4);
                var value3 = readBuffer.ReadByte("value3", 4);
                var value4 = readBuffer.ReadByte("value4", 4);
                var value5 = readBuffer.ReadByte("value5", 4);
                var detectionError = readBuffer.ReadBit("detectionError");
                var permission = readBuffer.ReadBit("permission");
                var readDirection = readBuffer.ReadBit("readDirection");
                var encryptionOfAccessInformation = readBuffer.ReadBit("encryptionOfAccessInformation");
                var indexOfAccessIdentificationCode = readBuffer.ReadByte("indexOfAccessIdentificationCode", 4);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["hurz"] = new PlcUSINT((byte) hurz);
                _map["value1"] = new PlcUSINT((byte) value1);
                _map["value2"] = new PlcUSINT((byte) value2);
                _map["value3"] = new PlcUSINT((byte) value3);
                _map["value4"] = new PlcUSINT((byte) value4);
                _map["value5"] = new PlcUSINT((byte) value5);
                _map["detectionError"] = new PlcBOOL((bool) detectionError);
                _map["permission"] = new PlcBOOL((bool) permission);
                _map["readDirection"] = new PlcBOOL((bool) readDirection);
                _map["encryptionOfAccessInformation"] = new PlcBOOL((bool) encryptionOfAccessInformation);
                _map["indexOfAccessIdentificationCode"] = new PlcUSINT((byte) indexOfAccessIdentificationCode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_String_ASCII))
            {
                // STRING
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 112, System.Text.Encoding.ASCII);
                return new PlcSTRING(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_String_8859_1))
            {
                // STRING
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 112, System.Text.Encoding.UTF8);
                return new PlcSTRING(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneNumber))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 6);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneControl))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var learnTheSceneCorrespondingToTheFieldSceneNumber = readBuffer.ReadBit("learnTheSceneCorrespondingToTheFieldSceneNumber");
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var sceneNumber = readBuffer.ReadByte("sceneNumber", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["learnTheSceneCorrespondingToTheFieldSceneNumber"] = new PlcBOOL((bool) learnTheSceneCorrespondingToTheFieldSceneNumber);
                _map["sceneNumber"] = new PlcUSINT((byte) sceneNumber);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DateTime))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var year = readBuffer.ReadByte("year", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var month = readBuffer.ReadByte("month", 4);
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var dayofmonth = readBuffer.ReadByte("dayofmonth", 5);
                var dayofweek = readBuffer.ReadByte("dayofweek", 3);
                var hourofday = readBuffer.ReadByte("hourofday", 5);
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var minutes = readBuffer.ReadByte("minutes", 6);
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var seconds = readBuffer.ReadByte("seconds", 6);
                var fault = readBuffer.ReadBit("fault");
                var workingDay = readBuffer.ReadBit("workingDay");
                var noWd = readBuffer.ReadBit("noWd");
                var noYear = readBuffer.ReadBit("noYear");
                var noDate = readBuffer.ReadBit("noDate");
                var noDayOfWeek = readBuffer.ReadBit("noDayOfWeek");
                var noTime = readBuffer.ReadBit("noTime");
                var standardSummerTime = readBuffer.ReadBit("standardSummerTime");
                var qualityOfClock = readBuffer.ReadBit("qualityOfClock");
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["year"] = new PlcUSINT((byte) year);
                _map["month"] = new PlcUSINT((byte) month);
                _map["dayofmonth"] = new PlcUSINT((byte) dayofmonth);
                _map["dayofweek"] = new PlcUSINT((byte) dayofweek);
                _map["hourofday"] = new PlcUSINT((byte) hourofday);
                _map["minutes"] = new PlcUSINT((byte) minutes);
                _map["seconds"] = new PlcUSINT((byte) seconds);
                _map["fault"] = new PlcBOOL((bool) fault);
                _map["workingDay"] = new PlcBOOL((bool) workingDay);
                _map["noWd"] = new PlcBOOL((bool) noWd);
                _map["noYear"] = new PlcBOOL((bool) noYear);
                _map["noDate"] = new PlcBOOL((bool) noDate);
                _map["noDayOfWeek"] = new PlcBOOL((bool) noDayOfWeek);
                _map["noTime"] = new PlcBOOL((bool) noTime);
                _map["standardSummerTime"] = new PlcBOOL((bool) standardSummerTime);
                _map["qualityOfClock"] = new PlcBOOL((bool) qualityOfClock);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SCLOMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BuildingMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OccMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Priority))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightApplicationMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApplicationArea))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AlarmClassType))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PSUMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ErrorClass_System))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ErrorClass_HVAC))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Time_Delay))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Beaufort_Wind_Force_Scale))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SensorSelect))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActuatorConnectType))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Cloud_Cover))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PowerReturnMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FuelType))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BurnerType))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DHWMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadPriority))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACContrMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACEmergMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ChangeoverMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ValveMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DamperMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HeaterMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FanMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MasterSlaveMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRoomSetp))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Metering_DeviceType))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HumDehumMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EnableHCStage))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ADAType))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BackupMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StartSynchronization))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Behaviour_Lock_Unlock))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Fade_Time))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BlinkingMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightControlMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SwitchPBModel))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PBAction))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DimmPBModel))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SwitchOnMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadTypeSet))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadTypeDetected))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Test_Control))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Control))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SABExcept_Behaviour))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SABBehaviour_Lock_Unlock))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SSSBMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BlindsControlMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CommMode))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AddInfoTypes))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_ModeSelect))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_FilterSelect))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_1))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_2))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_3))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusGen))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var alarmStatusOfCorrespondingDatapointIsNotAcknowledged = readBuffer.ReadBit("alarmStatusOfCorrespondingDatapointIsNotAcknowledged");
                var correspondingDatapointIsInAlarm = readBuffer.ReadBit("correspondingDatapointIsInAlarm");
                var correspondingDatapointMainValueIsOverridden = readBuffer.ReadBit("correspondingDatapointMainValueIsOverridden");
                var correspondingDatapointMainValueIsCorruptedDueToFailure = readBuffer.ReadBit("correspondingDatapointMainValueIsCorruptedDueToFailure");
                var correspondingDatapointValueIsOutOfService = readBuffer.ReadBit("correspondingDatapointValueIsOutOfService");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["alarmStatusOfCorrespondingDatapointIsNotAcknowledged"] = new PlcBOOL((bool) alarmStatusOfCorrespondingDatapointIsNotAcknowledged);
                _map["correspondingDatapointIsInAlarm"] = new PlcBOOL((bool) correspondingDatapointIsInAlarm);
                _map["correspondingDatapointMainValueIsOverridden"] = new PlcBOOL((bool) correspondingDatapointMainValueIsOverridden);
                _map["correspondingDatapointMainValueIsCorruptedDueToFailure"] = new PlcBOOL((bool) correspondingDatapointMainValueIsCorruptedDueToFailure);
                _map["correspondingDatapointValueIsOutOfService"] = new PlcBOOL((bool) correspondingDatapointValueIsOutOfService);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Device_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var verifyModeIsOn = readBuffer.ReadBit("verifyModeIsOn");
                var aDatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived = readBuffer.ReadBit("aDatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived");
                var theUserApplicationIsStopped = readBuffer.ReadBit("theUserApplicationIsStopped");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["verifyModeIsOn"] = new PlcBOOL((bool) verifyModeIsOn);
                _map["aDatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived"] = new PlcBOOL((bool) aDatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived);
                _map["theUserApplicationIsStopped"] = new PlcBOOL((bool) theUserApplicationIsStopped);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ForceSign))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var roomhmax = readBuffer.ReadBit("roomhmax");
                var roomhconf = readBuffer.ReadBit("roomhconf");
                var dhwlegio = readBuffer.ReadBit("dhwlegio");
                var dhwnorm = readBuffer.ReadBit("dhwnorm");
                var overrun = readBuffer.ReadBit("overrun");
                var oversupply = readBuffer.ReadBit("oversupply");
                var protection = readBuffer.ReadBit("protection");
                var forcerequest = readBuffer.ReadBit("forcerequest");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["roomhmax"] = new PlcBOOL((bool) roomhmax);
                _map["roomhconf"] = new PlcBOOL((bool) roomhconf);
                _map["dhwlegio"] = new PlcBOOL((bool) dhwlegio);
                _map["dhwnorm"] = new PlcBOOL((bool) dhwnorm);
                _map["overrun"] = new PlcBOOL((bool) overrun);
                _map["oversupply"] = new PlcBOOL((bool) oversupply);
                _map["protection"] = new PlcBOOL((bool) protection);
                _map["forcerequest"] = new PlcBOOL((bool) forcerequest);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ForceSignCool))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRHC))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var summermode = readBuffer.ReadBit("summermode");
                var statusstopoptim = readBuffer.ReadBit("statusstopoptim");
                var statusstartoptim = readBuffer.ReadBit("statusstartoptim");
                var statusmorningboost = readBuffer.ReadBit("statusmorningboost");
                var tempreturnlimit = readBuffer.ReadBit("tempreturnlimit");
                var tempflowlimit = readBuffer.ReadBit("tempflowlimit");
                var satuseco = readBuffer.ReadBit("satuseco");
                var fault = readBuffer.ReadBit("fault");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["summermode"] = new PlcBOOL((bool) summermode);
                _map["statusstopoptim"] = new PlcBOOL((bool) statusstopoptim);
                _map["statusstartoptim"] = new PlcBOOL((bool) statusstartoptim);
                _map["statusmorningboost"] = new PlcBOOL((bool) statusmorningboost);
                _map["tempreturnlimit"] = new PlcBOOL((bool) tempreturnlimit);
                _map["tempflowlimit"] = new PlcBOOL((bool) tempflowlimit);
                _map["satuseco"] = new PlcBOOL((bool) satuseco);
                _map["fault"] = new PlcBOOL((bool) fault);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusSDHWC))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var solarloadsufficient = readBuffer.ReadBit("solarloadsufficient");
                var sdhwloadactive = readBuffer.ReadBit("sdhwloadactive");
                var fault = readBuffer.ReadBit("fault");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["solarloadsufficient"] = new PlcBOOL((bool) solarloadsufficient);
                _map["sdhwloadactive"] = new PlcBOOL((bool) sdhwloadactive);
                _map["fault"] = new PlcBOOL((bool) fault);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FuelTypeSet))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var solidstate = readBuffer.ReadBit("solidstate");
                var gas = readBuffer.ReadBit("gas");
                var oil = readBuffer.ReadBit("oil");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["solidstate"] = new PlcBOOL((bool) solidstate);
                _map["gas"] = new PlcBOOL((bool) gas);
                _map["oil"] = new PlcBOOL((bool) oil);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRCC))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusAHU))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cool = readBuffer.ReadBit("cool");
                var heat = readBuffer.ReadBit("heat");
                var fanactive = readBuffer.ReadBit("fanactive");
                var fault = readBuffer.ReadBit("fault");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["cool"] = new PlcBOOL((bool) cool);
                _map["heat"] = new PlcBOOL((bool) heat);
                _map["fanactive"] = new PlcBOOL((bool) fanactive);
                _map["fault"] = new PlcBOOL((bool) fault);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_RTSM))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var statusOfHvacModeUser = readBuffer.ReadBit("statusOfHvacModeUser");
                var statusOfComfortProlongationUser = readBuffer.ReadBit("statusOfComfortProlongationUser");
                var effectiveValueOfTheComfortPushButton = readBuffer.ReadBit("effectiveValueOfTheComfortPushButton");
                var effectiveValueOfThePresenceStatus = readBuffer.ReadBit("effectiveValueOfThePresenceStatus");
                var effectiveValueOfTheWindowStatus = readBuffer.ReadBit("effectiveValueOfTheWindowStatus");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["statusOfHvacModeUser"] = new PlcBOOL((bool) statusOfHvacModeUser);
                _map["statusOfComfortProlongationUser"] = new PlcBOOL((bool) statusOfComfortProlongationUser);
                _map["effectiveValueOfTheComfortPushButton"] = new PlcBOOL((bool) effectiveValueOfTheComfortPushButton);
                _map["effectiveValueOfThePresenceStatus"] = new PlcBOOL((bool) effectiveValueOfThePresenceStatus);
                _map["effectiveValueOfTheWindowStatus"] = new PlcBOOL((bool) effectiveValueOfTheWindowStatus);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightActuatorErrorInfo))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var overheat = readBuffer.ReadBit("overheat");
                var lampfailure = readBuffer.ReadBit("lampfailure");
                var defectiveload = readBuffer.ReadBit("defectiveload");
                var underload = readBuffer.ReadBit("underload");
                var overcurrent = readBuffer.ReadBit("overcurrent");
                var undervoltage = readBuffer.ReadBit("undervoltage");
                var loaddetectionerror = readBuffer.ReadBit("loaddetectionerror");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["overheat"] = new PlcBOOL((bool) overheat);
                _map["lampfailure"] = new PlcBOOL((bool) lampfailure);
                _map["defectiveload"] = new PlcBOOL((bool) defectiveload);
                _map["underload"] = new PlcBOOL((bool) underload);
                _map["overcurrent"] = new PlcBOOL((bool) overcurrent);
                _map["undervoltage"] = new PlcBOOL((bool) undervoltage);
                _map["loaddetectionerror"] = new PlcBOOL((bool) loaddetectionerror);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_ModeInfo))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var bibatSlave = readBuffer.ReadBit("bibatSlave");
                var bibatMaster = readBuffer.ReadBit("bibatMaster");
                var asynchronous = readBuffer.ReadBit("asynchronous");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["bibatSlave"] = new PlcBOOL((bool) bibatSlave);
                _map["bibatMaster"] = new PlcBOOL((bool) bibatMaster);
                _map["asynchronous"] = new PlcBOOL((bool) asynchronous);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_FilterInfo))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var doa = readBuffer.ReadBit("doa");
                var knxSn = readBuffer.ReadBit("knxSn");
                var doaAndKnxSn = readBuffer.ReadBit("doaAndKnxSn");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["doa"] = new PlcBOOL((bool) doa);
                _map["knxSn"] = new PlcBOOL((bool) knxSn);
                _map["doaAndKnxSn"] = new PlcBOOL((bool) doaAndKnxSn);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_8))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var activationStateOfChannel1 = readBuffer.ReadBit("activationStateOfChannel1");
                var activationStateOfChannel2 = readBuffer.ReadBit("activationStateOfChannel2");
                var activationStateOfChannel3 = readBuffer.ReadBit("activationStateOfChannel3");
                var activationStateOfChannel4 = readBuffer.ReadBit("activationStateOfChannel4");
                var activationStateOfChannel5 = readBuffer.ReadBit("activationStateOfChannel5");
                var activationStateOfChannel6 = readBuffer.ReadBit("activationStateOfChannel6");
                var activationStateOfChannel7 = readBuffer.ReadBit("activationStateOfChannel7");
                var activationStateOfChannel8 = readBuffer.ReadBit("activationStateOfChannel8");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["activationStateOfChannel1"] = new PlcBOOL((bool) activationStateOfChannel1);
                _map["activationStateOfChannel2"] = new PlcBOOL((bool) activationStateOfChannel2);
                _map["activationStateOfChannel3"] = new PlcBOOL((bool) activationStateOfChannel3);
                _map["activationStateOfChannel4"] = new PlcBOOL((bool) activationStateOfChannel4);
                _map["activationStateOfChannel5"] = new PlcBOOL((bool) activationStateOfChannel5);
                _map["activationStateOfChannel6"] = new PlcBOOL((bool) activationStateOfChannel6);
                _map["activationStateOfChannel7"] = new PlcBOOL((bool) activationStateOfChannel7);
                _map["activationStateOfChannel8"] = new PlcBOOL((bool) activationStateOfChannel8);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusDHWC))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var tempoptimshiftactive = readBuffer.ReadBit("tempoptimshiftactive");
                var solarenergysupport = readBuffer.ReadBit("solarenergysupport");
                var solarenergyonly = readBuffer.ReadBit("solarenergyonly");
                var otherenergysourceactive = readBuffer.ReadBit("otherenergysourceactive");
                var dhwpushactive = readBuffer.ReadBit("dhwpushactive");
                var legioprotactive = readBuffer.ReadBit("legioprotactive");
                var dhwloadactive = readBuffer.ReadBit("dhwloadactive");
                var fault = readBuffer.ReadBit("fault");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["tempoptimshiftactive"] = new PlcBOOL((bool) tempoptimshiftactive);
                _map["solarenergysupport"] = new PlcBOOL((bool) solarenergysupport);
                _map["solarenergyonly"] = new PlcBOOL((bool) solarenergyonly);
                _map["otherenergysourceactive"] = new PlcBOOL((bool) otherenergysourceactive);
                _map["dhwpushactive"] = new PlcBOOL((bool) dhwpushactive);
                _map["legioprotactive"] = new PlcBOOL((bool) legioprotactive);
                _map["dhwloadactive"] = new PlcBOOL((bool) dhwloadactive);
                _map["fault"] = new PlcBOOL((bool) fault);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRHCC))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var overheatalarm = readBuffer.ReadBit("overheatalarm");
                var frostalarm = readBuffer.ReadBit("frostalarm");
                var dewpointstatus = readBuffer.ReadBit("dewpointstatus");
                var coolingdisabled = readBuffer.ReadBit("coolingdisabled");
                var statusprecool = readBuffer.ReadBit("statusprecool");
                var statusecoc = readBuffer.ReadBit("statusecoc");
                var heatcoolmode = readBuffer.ReadBit("heatcoolmode");
                var heatingdiabled = readBuffer.ReadBit("heatingdiabled");
                var statusstopoptim = readBuffer.ReadBit("statusstopoptim");
                var statusstartoptim = readBuffer.ReadBit("statusstartoptim");
                var statusmorningboosth = readBuffer.ReadBit("statusmorningboosth");
                var tempflowreturnlimit = readBuffer.ReadBit("tempflowreturnlimit");
                var tempflowlimit = readBuffer.ReadBit("tempflowlimit");
                var statusecoh = readBuffer.ReadBit("statusecoh");
                var fault = readBuffer.ReadBit("fault");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["overheatalarm"] = new PlcBOOL((bool) overheatalarm);
                _map["frostalarm"] = new PlcBOOL((bool) frostalarm);
                _map["dewpointstatus"] = new PlcBOOL((bool) dewpointstatus);
                _map["coolingdisabled"] = new PlcBOOL((bool) coolingdisabled);
                _map["statusprecool"] = new PlcBOOL((bool) statusprecool);
                _map["statusecoc"] = new PlcBOOL((bool) statusecoc);
                _map["heatcoolmode"] = new PlcBOOL((bool) heatcoolmode);
                _map["heatingdiabled"] = new PlcBOOL((bool) heatingdiabled);
                _map["statusstopoptim"] = new PlcBOOL((bool) statusstopoptim);
                _map["statusstartoptim"] = new PlcBOOL((bool) statusstartoptim);
                _map["statusmorningboosth"] = new PlcBOOL((bool) statusmorningboosth);
                _map["tempflowreturnlimit"] = new PlcBOOL((bool) tempflowreturnlimit);
                _map["tempflowlimit"] = new PlcBOOL((bool) tempflowlimit);
                _map["statusecoh"] = new PlcBOOL((bool) statusecoh);
                _map["fault"] = new PlcBOOL((bool) fault);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_HVA))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var calibrationMode = readBuffer.ReadBit("calibrationMode");
                var lockedPosition = readBuffer.ReadBit("lockedPosition");
                var forcedPosition = readBuffer.ReadBit("forcedPosition");
                var manuaOperationOverridden = readBuffer.ReadBit("manuaOperationOverridden");
                var serviceMode = readBuffer.ReadBit("serviceMode");
                var valveKick = readBuffer.ReadBit("valveKick");
                var overload = readBuffer.ReadBit("overload");
                var shortCircuit = readBuffer.ReadBit("shortCircuit");
                var currentValvePosition = readBuffer.ReadBit("currentValvePosition");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["calibrationMode"] = new PlcBOOL((bool) calibrationMode);
                _map["lockedPosition"] = new PlcBOOL((bool) lockedPosition);
                _map["forcedPosition"] = new PlcBOOL((bool) forcedPosition);
                _map["manuaOperationOverridden"] = new PlcBOOL((bool) manuaOperationOverridden);
                _map["serviceMode"] = new PlcBOOL((bool) serviceMode);
                _map["valveKick"] = new PlcBOOL((bool) valveKick);
                _map["overload"] = new PlcBOOL((bool) overload);
                _map["shortCircuit"] = new PlcBOOL((bool) shortCircuit);
                _map["currentValvePosition"] = new PlcBOOL((bool) currentValvePosition);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_RTC))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var coolingModeEnabled = readBuffer.ReadBit("coolingModeEnabled");
                var heatingModeEnabled = readBuffer.ReadBit("heatingModeEnabled");
                var additionalHeatingCoolingStage2Stage = readBuffer.ReadBit("additionalHeatingCoolingStage2Stage");
                var controllerInactive = readBuffer.ReadBit("controllerInactive");
                var overheatAlarm = readBuffer.ReadBit("overheatAlarm");
                var frostAlarm = readBuffer.ReadBit("frostAlarm");
                var dewPointStatus = readBuffer.ReadBit("dewPointStatus");
                var activeMode = readBuffer.ReadBit("activeMode");
                var generalFailureInformation = readBuffer.ReadBit("generalFailureInformation");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["coolingModeEnabled"] = new PlcBOOL((bool) coolingModeEnabled);
                _map["heatingModeEnabled"] = new PlcBOOL((bool) heatingModeEnabled);
                _map["additionalHeatingCoolingStage2Stage"] = new PlcBOOL((bool) additionalHeatingCoolingStage2Stage);
                _map["controllerInactive"] = new PlcBOOL((bool) controllerInactive);
                _map["overheatAlarm"] = new PlcBOOL((bool) overheatAlarm);
                _map["frostAlarm"] = new PlcBOOL((bool) frostAlarm);
                _map["dewPointStatus"] = new PlcBOOL((bool) dewPointStatus);
                _map["activeMode"] = new PlcBOOL((bool) activeMode);
                _map["generalFailureInformation"] = new PlcBOOL((bool) generalFailureInformation);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Media))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadUshort("reserved", 10);
                    if (!Equals(reserved, (ushort) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var knxIp = readBuffer.ReadBit("knxIp");
                var rf = readBuffer.ReadBit("rf");
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var pl110 = readBuffer.ReadBit("pl110");
                var tp1 = readBuffer.ReadBit("tp1");
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["knxIp"] = new PlcBOOL((bool) knxIp);
                _map["rf"] = new PlcBOOL((bool) rf);
                _map["pl110"] = new PlcBOOL((bool) pl110);
                _map["tp1"] = new PlcBOOL((bool) tp1);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_16))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var activationStateOfChannel1 = readBuffer.ReadBit("activationStateOfChannel1");
                var activationStateOfChannel2 = readBuffer.ReadBit("activationStateOfChannel2");
                var activationStateOfChannel3 = readBuffer.ReadBit("activationStateOfChannel3");
                var activationStateOfChannel4 = readBuffer.ReadBit("activationStateOfChannel4");
                var activationStateOfChannel5 = readBuffer.ReadBit("activationStateOfChannel5");
                var activationStateOfChannel6 = readBuffer.ReadBit("activationStateOfChannel6");
                var activationStateOfChannel7 = readBuffer.ReadBit("activationStateOfChannel7");
                var activationStateOfChannel8 = readBuffer.ReadBit("activationStateOfChannel8");
                var activationStateOfChannel9 = readBuffer.ReadBit("activationStateOfChannel9");
                var activationStateOfChannel10 = readBuffer.ReadBit("activationStateOfChannel10");
                var activationStateOfChannel11 = readBuffer.ReadBit("activationStateOfChannel11");
                var activationStateOfChannel12 = readBuffer.ReadBit("activationStateOfChannel12");
                var activationStateOfChannel13 = readBuffer.ReadBit("activationStateOfChannel13");
                var activationStateOfChannel14 = readBuffer.ReadBit("activationStateOfChannel14");
                var activationStateOfChannel15 = readBuffer.ReadBit("activationStateOfChannel15");
                var activationStateOfChannel16 = readBuffer.ReadBit("activationStateOfChannel16");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["activationStateOfChannel1"] = new PlcBOOL((bool) activationStateOfChannel1);
                _map["activationStateOfChannel2"] = new PlcBOOL((bool) activationStateOfChannel2);
                _map["activationStateOfChannel3"] = new PlcBOOL((bool) activationStateOfChannel3);
                _map["activationStateOfChannel4"] = new PlcBOOL((bool) activationStateOfChannel4);
                _map["activationStateOfChannel5"] = new PlcBOOL((bool) activationStateOfChannel5);
                _map["activationStateOfChannel6"] = new PlcBOOL((bool) activationStateOfChannel6);
                _map["activationStateOfChannel7"] = new PlcBOOL((bool) activationStateOfChannel7);
                _map["activationStateOfChannel8"] = new PlcBOOL((bool) activationStateOfChannel8);
                _map["activationStateOfChannel9"] = new PlcBOOL((bool) activationStateOfChannel9);
                _map["activationStateOfChannel10"] = new PlcBOOL((bool) activationStateOfChannel10);
                _map["activationStateOfChannel11"] = new PlcBOOL((bool) activationStateOfChannel11);
                _map["activationStateOfChannel12"] = new PlcBOOL((bool) activationStateOfChannel12);
                _map["activationStateOfChannel13"] = new PlcBOOL((bool) activationStateOfChannel13);
                _map["activationStateOfChannel14"] = new PlcBOOL((bool) activationStateOfChannel14);
                _map["activationStateOfChannel15"] = new PlcBOOL((bool) activationStateOfChannel15);
                _map["activationStateOfChannel16"] = new PlcBOOL((bool) activationStateOfChannel16);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OnOffAction))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 2);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm_Reaction))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 2);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UpDown_Action))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 2);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVAC_PB_Action))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 2);
                return new PlcUSINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DoubleNibble))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var busy = readBuffer.ReadByte("busy", 4);
                var nak = readBuffer.ReadByte("nak", 4);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["busy"] = new PlcUSINT((byte) busy);
                _map["nak"] = new PlcUSINT((byte) nak);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneInfo))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 1);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var sceneIsInactive = readBuffer.ReadBit("sceneIsInactive");
                var scenenumber = readBuffer.ReadByte("scenenumber", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["sceneIsInactive"] = new PlcBOOL((bool) sceneIsInactive);
                _map["scenenumber"] = new PlcUSINT((byte) scenenumber);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedInfoOnOff))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var maskBitInfoOnOffOutput16 = readBuffer.ReadBit("maskBitInfoOnOffOutput16");
                var maskBitInfoOnOffOutput15 = readBuffer.ReadBit("maskBitInfoOnOffOutput15");
                var maskBitInfoOnOffOutput14 = readBuffer.ReadBit("maskBitInfoOnOffOutput14");
                var maskBitInfoOnOffOutput13 = readBuffer.ReadBit("maskBitInfoOnOffOutput13");
                var maskBitInfoOnOffOutput12 = readBuffer.ReadBit("maskBitInfoOnOffOutput12");
                var maskBitInfoOnOffOutput11 = readBuffer.ReadBit("maskBitInfoOnOffOutput11");
                var maskBitInfoOnOffOutput10 = readBuffer.ReadBit("maskBitInfoOnOffOutput10");
                var maskBitInfoOnOffOutput9 = readBuffer.ReadBit("maskBitInfoOnOffOutput9");
                var maskBitInfoOnOffOutput8 = readBuffer.ReadBit("maskBitInfoOnOffOutput8");
                var maskBitInfoOnOffOutput7 = readBuffer.ReadBit("maskBitInfoOnOffOutput7");
                var maskBitInfoOnOffOutput6 = readBuffer.ReadBit("maskBitInfoOnOffOutput6");
                var maskBitInfoOnOffOutput5 = readBuffer.ReadBit("maskBitInfoOnOffOutput5");
                var maskBitInfoOnOffOutput4 = readBuffer.ReadBit("maskBitInfoOnOffOutput4");
                var maskBitInfoOnOffOutput3 = readBuffer.ReadBit("maskBitInfoOnOffOutput3");
                var maskBitInfoOnOffOutput2 = readBuffer.ReadBit("maskBitInfoOnOffOutput2");
                var maskBitInfoOnOffOutput1 = readBuffer.ReadBit("maskBitInfoOnOffOutput1");
                var infoOnOffOutput16 = readBuffer.ReadBit("infoOnOffOutput16");
                var infoOnOffOutput15 = readBuffer.ReadBit("infoOnOffOutput15");
                var infoOnOffOutput14 = readBuffer.ReadBit("infoOnOffOutput14");
                var infoOnOffOutput13 = readBuffer.ReadBit("infoOnOffOutput13");
                var infoOnOffOutput12 = readBuffer.ReadBit("infoOnOffOutput12");
                var infoOnOffOutput11 = readBuffer.ReadBit("infoOnOffOutput11");
                var infoOnOffOutput10 = readBuffer.ReadBit("infoOnOffOutput10");
                var infoOnOffOutput9 = readBuffer.ReadBit("infoOnOffOutput9");
                var infoOnOffOutput8 = readBuffer.ReadBit("infoOnOffOutput8");
                var infoOnOffOutput7 = readBuffer.ReadBit("infoOnOffOutput7");
                var infoOnOffOutput6 = readBuffer.ReadBit("infoOnOffOutput6");
                var infoOnOffOutput5 = readBuffer.ReadBit("infoOnOffOutput5");
                var infoOnOffOutput4 = readBuffer.ReadBit("infoOnOffOutput4");
                var infoOnOffOutput3 = readBuffer.ReadBit("infoOnOffOutput3");
                var infoOnOffOutput2 = readBuffer.ReadBit("infoOnOffOutput2");
                var infoOnOffOutput1 = readBuffer.ReadBit("infoOnOffOutput1");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["maskBitInfoOnOffOutput16"] = new PlcBOOL((bool) maskBitInfoOnOffOutput16);
                _map["maskBitInfoOnOffOutput15"] = new PlcBOOL((bool) maskBitInfoOnOffOutput15);
                _map["maskBitInfoOnOffOutput14"] = new PlcBOOL((bool) maskBitInfoOnOffOutput14);
                _map["maskBitInfoOnOffOutput13"] = new PlcBOOL((bool) maskBitInfoOnOffOutput13);
                _map["maskBitInfoOnOffOutput12"] = new PlcBOOL((bool) maskBitInfoOnOffOutput12);
                _map["maskBitInfoOnOffOutput11"] = new PlcBOOL((bool) maskBitInfoOnOffOutput11);
                _map["maskBitInfoOnOffOutput10"] = new PlcBOOL((bool) maskBitInfoOnOffOutput10);
                _map["maskBitInfoOnOffOutput9"] = new PlcBOOL((bool) maskBitInfoOnOffOutput9);
                _map["maskBitInfoOnOffOutput8"] = new PlcBOOL((bool) maskBitInfoOnOffOutput8);
                _map["maskBitInfoOnOffOutput7"] = new PlcBOOL((bool) maskBitInfoOnOffOutput7);
                _map["maskBitInfoOnOffOutput6"] = new PlcBOOL((bool) maskBitInfoOnOffOutput6);
                _map["maskBitInfoOnOffOutput5"] = new PlcBOOL((bool) maskBitInfoOnOffOutput5);
                _map["maskBitInfoOnOffOutput4"] = new PlcBOOL((bool) maskBitInfoOnOffOutput4);
                _map["maskBitInfoOnOffOutput3"] = new PlcBOOL((bool) maskBitInfoOnOffOutput3);
                _map["maskBitInfoOnOffOutput2"] = new PlcBOOL((bool) maskBitInfoOnOffOutput2);
                _map["maskBitInfoOnOffOutput1"] = new PlcBOOL((bool) maskBitInfoOnOffOutput1);
                _map["infoOnOffOutput16"] = new PlcBOOL((bool) infoOnOffOutput16);
                _map["infoOnOffOutput15"] = new PlcBOOL((bool) infoOnOffOutput15);
                _map["infoOnOffOutput14"] = new PlcBOOL((bool) infoOnOffOutput14);
                _map["infoOnOffOutput13"] = new PlcBOOL((bool) infoOnOffOutput13);
                _map["infoOnOffOutput12"] = new PlcBOOL((bool) infoOnOffOutput12);
                _map["infoOnOffOutput11"] = new PlcBOOL((bool) infoOnOffOutput11);
                _map["infoOnOffOutput10"] = new PlcBOOL((bool) infoOnOffOutput10);
                _map["infoOnOffOutput9"] = new PlcBOOL((bool) infoOnOffOutput9);
                _map["infoOnOffOutput8"] = new PlcBOOL((bool) infoOnOffOutput8);
                _map["infoOnOffOutput7"] = new PlcBOOL((bool) infoOnOffOutput7);
                _map["infoOnOffOutput6"] = new PlcBOOL((bool) infoOnOffOutput6);
                _map["infoOnOffOutput5"] = new PlcBOOL((bool) infoOnOffOutput5);
                _map["infoOnOffOutput4"] = new PlcBOOL((bool) infoOnOffOutput4);
                _map["infoOnOffOutput3"] = new PlcBOOL((bool) infoOnOffOutput3);
                _map["infoOnOffOutput2"] = new PlcBOOL((bool) infoOnOffOutput2);
                _map["infoOnOffOutput1"] = new PlcBOOL((bool) infoOnOffOutput1);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_V64))
            {
                // LINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy_V64))
            {
                // LINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy_V64))
            {
                // LINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_24))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var activationStateOfChannel1 = readBuffer.ReadBit("activationStateOfChannel1");
                var activationStateOfChannel2 = readBuffer.ReadBit("activationStateOfChannel2");
                var activationStateOfChannel3 = readBuffer.ReadBit("activationStateOfChannel3");
                var activationStateOfChannel4 = readBuffer.ReadBit("activationStateOfChannel4");
                var activationStateOfChannel5 = readBuffer.ReadBit("activationStateOfChannel5");
                var activationStateOfChannel6 = readBuffer.ReadBit("activationStateOfChannel6");
                var activationStateOfChannel7 = readBuffer.ReadBit("activationStateOfChannel7");
                var activationStateOfChannel8 = readBuffer.ReadBit("activationStateOfChannel8");
                var activationStateOfChannel9 = readBuffer.ReadBit("activationStateOfChannel9");
                var activationStateOfChannel10 = readBuffer.ReadBit("activationStateOfChannel10");
                var activationStateOfChannel11 = readBuffer.ReadBit("activationStateOfChannel11");
                var activationStateOfChannel12 = readBuffer.ReadBit("activationStateOfChannel12");
                var activationStateOfChannel13 = readBuffer.ReadBit("activationStateOfChannel13");
                var activationStateOfChannel14 = readBuffer.ReadBit("activationStateOfChannel14");
                var activationStateOfChannel15 = readBuffer.ReadBit("activationStateOfChannel15");
                var activationStateOfChannel16 = readBuffer.ReadBit("activationStateOfChannel16");
                var activationStateOfChannel17 = readBuffer.ReadBit("activationStateOfChannel17");
                var activationStateOfChannel18 = readBuffer.ReadBit("activationStateOfChannel18");
                var activationStateOfChannel19 = readBuffer.ReadBit("activationStateOfChannel19");
                var activationStateOfChannel20 = readBuffer.ReadBit("activationStateOfChannel20");
                var activationStateOfChannel21 = readBuffer.ReadBit("activationStateOfChannel21");
                var activationStateOfChannel22 = readBuffer.ReadBit("activationStateOfChannel22");
                var activationStateOfChannel23 = readBuffer.ReadBit("activationStateOfChannel23");
                var activationStateOfChannel24 = readBuffer.ReadBit("activationStateOfChannel24");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["activationStateOfChannel1"] = new PlcBOOL((bool) activationStateOfChannel1);
                _map["activationStateOfChannel2"] = new PlcBOOL((bool) activationStateOfChannel2);
                _map["activationStateOfChannel3"] = new PlcBOOL((bool) activationStateOfChannel3);
                _map["activationStateOfChannel4"] = new PlcBOOL((bool) activationStateOfChannel4);
                _map["activationStateOfChannel5"] = new PlcBOOL((bool) activationStateOfChannel5);
                _map["activationStateOfChannel6"] = new PlcBOOL((bool) activationStateOfChannel6);
                _map["activationStateOfChannel7"] = new PlcBOOL((bool) activationStateOfChannel7);
                _map["activationStateOfChannel8"] = new PlcBOOL((bool) activationStateOfChannel8);
                _map["activationStateOfChannel9"] = new PlcBOOL((bool) activationStateOfChannel9);
                _map["activationStateOfChannel10"] = new PlcBOOL((bool) activationStateOfChannel10);
                _map["activationStateOfChannel11"] = new PlcBOOL((bool) activationStateOfChannel11);
                _map["activationStateOfChannel12"] = new PlcBOOL((bool) activationStateOfChannel12);
                _map["activationStateOfChannel13"] = new PlcBOOL((bool) activationStateOfChannel13);
                _map["activationStateOfChannel14"] = new PlcBOOL((bool) activationStateOfChannel14);
                _map["activationStateOfChannel15"] = new PlcBOOL((bool) activationStateOfChannel15);
                _map["activationStateOfChannel16"] = new PlcBOOL((bool) activationStateOfChannel16);
                _map["activationStateOfChannel17"] = new PlcBOOL((bool) activationStateOfChannel17);
                _map["activationStateOfChannel18"] = new PlcBOOL((bool) activationStateOfChannel18);
                _map["activationStateOfChannel19"] = new PlcBOOL((bool) activationStateOfChannel19);
                _map["activationStateOfChannel20"] = new PlcBOOL((bool) activationStateOfChannel20);
                _map["activationStateOfChannel21"] = new PlcBOOL((bool) activationStateOfChannel21);
                _map["activationStateOfChannel22"] = new PlcBOOL((bool) activationStateOfChannel22);
                _map["activationStateOfChannel23"] = new PlcBOOL((bool) activationStateOfChannel23);
                _map["activationStateOfChannel24"] = new PlcBOOL((bool) activationStateOfChannel24);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACModeNext))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var delayTimeMin = readBuffer.ReadUshort("delayTimeMin", 16);
                var hvacMode = readBuffer.ReadByte("hvacMode", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["delayTimeMin"] = new PlcUINT((ushort) delayTimeMin);
                _map["hvacMode"] = new PlcUSINT((byte) hvacMode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DHWModeNext))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var delayTimeMin = readBuffer.ReadUshort("delayTimeMin", 16);
                var dhwMode = readBuffer.ReadByte("dhwMode", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["delayTimeMin"] = new PlcUINT((ushort) delayTimeMin);
                _map["dhwMode"] = new PlcUSINT((byte) dhwMode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OccModeNext))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var delayTimeMin = readBuffer.ReadUshort("delayTimeMin", 16);
                var occupancyMode = readBuffer.ReadByte("occupancyMode", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["delayTimeMin"] = new PlcUINT((ushort) delayTimeMin);
                _map["occupancyMode"] = new PlcUSINT((byte) occupancyMode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BuildingModeNext))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var delayTimeMin = readBuffer.ReadUshort("delayTimeMin", 16);
                var buildingMode = readBuffer.ReadByte("buildingMode", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["delayTimeMin"] = new PlcUINT((ushort) delayTimeMin);
                _map["buildingMode"] = new PlcUSINT((byte) buildingMode);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusLightingActuator))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var validactualvalue = readBuffer.ReadBit("validactualvalue");
                var locked = readBuffer.ReadBit("locked");
                var forced = readBuffer.ReadBit("forced");
                var nightmodeactive = readBuffer.ReadBit("nightmodeactive");
                var staircaselightingFunction = readBuffer.ReadBit("staircaselightingFunction");
                var dimming = readBuffer.ReadBit("dimming");
                var localoverride = readBuffer.ReadBit("localoverride");
                var failure = readBuffer.ReadBit("failure");
                var actualvalue = readBuffer.ReadByte("actualvalue", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["validactualvalue"] = new PlcBOOL((bool) validactualvalue);
                _map["locked"] = new PlcBOOL((bool) locked);
                _map["forced"] = new PlcBOOL((bool) forced);
                _map["nightmodeactive"] = new PlcBOOL((bool) nightmodeactive);
                _map["staircaselightingFunction"] = new PlcBOOL((bool) staircaselightingFunction);
                _map["dimming"] = new PlcBOOL((bool) dimming);
                _map["localoverride"] = new PlcBOOL((bool) localoverride);
                _map["failure"] = new PlcBOOL((bool) failure);
                _map["actualvalue"] = new PlcUSINT((byte) actualvalue);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Version))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var magicNumber = readBuffer.ReadByte("magicNumber", 5);
                var versionNumber = readBuffer.ReadByte("versionNumber", 5);
                var revisionNumber = readBuffer.ReadByte("revisionNumber", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["magicNumber"] = new PlcUSINT((byte) magicNumber);
                _map["versionNumber"] = new PlcUSINT((byte) versionNumber);
                _map["revisionNumber"] = new PlcUSINT((byte) revisionNumber);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AlarmInfo))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var logNumber = readBuffer.ReadByte("logNumber", 8);
                var alarmPriority = readBuffer.ReadByte("alarmPriority", 8);
                var applicationArea = readBuffer.ReadByte("applicationArea", 8);
                var errorClass = readBuffer.ReadByte("errorClass", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var errorcodeSup = readBuffer.ReadBit("errorcodeSup");
                var alarmtextSup = readBuffer.ReadBit("alarmtextSup");
                var timestampSup = readBuffer.ReadBit("timestampSup");
                var ackSup = readBuffer.ReadBit("ackSup");
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var locked = readBuffer.ReadBit("locked");
                var alarmunack = readBuffer.ReadBit("alarmunack");
                var inalarm = readBuffer.ReadBit("inalarm");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["logNumber"] = new PlcUSINT((byte) logNumber);
                _map["alarmPriority"] = new PlcUSINT((byte) alarmPriority);
                _map["applicationArea"] = new PlcUSINT((byte) applicationArea);
                _map["errorClass"] = new PlcUSINT((byte) errorClass);
                _map["errorcodeSup"] = new PlcBOOL((bool) errorcodeSup);
                _map["alarmtextSup"] = new PlcBOOL((bool) alarmtextSup);
                _map["timestampSup"] = new PlcBOOL((bool) timestampSup);
                _map["ackSup"] = new PlcBOOL((bool) ackSup);
                _map["locked"] = new PlcBOOL((bool) locked);
                _map["alarmunack"] = new PlcBOOL((bool) alarmunack);
                _map["inalarm"] = new PlcBOOL((bool) inalarm);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetF16_3))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var tempsetpcomf = readBuffer.ReadFloat("tempsetpcomf", 16);
                var tempsetpstdby = readBuffer.ReadFloat("tempsetpstdby", 16);
                var tempsetpeco = readBuffer.ReadFloat("tempsetpeco", 16);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["tempsetpcomf"] = new PlcREAL((float) tempsetpcomf);
                _map["tempsetpstdby"] = new PlcREAL((float) tempsetpstdby);
                _map["tempsetpeco"] = new PlcREAL((float) tempsetpeco);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetShiftF16_3))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var tempsetpshiftcomf = readBuffer.ReadFloat("tempsetpshiftcomf", 16);
                var tempsetpshiftstdby = readBuffer.ReadFloat("tempsetpshiftstdby", 16);
                var tempsetpshifteco = readBuffer.ReadFloat("tempsetpshifteco", 16);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["tempsetpshiftcomf"] = new PlcREAL((float) tempsetpshiftcomf);
                _map["tempsetpshiftstdby"] = new PlcREAL((float) tempsetpshiftstdby);
                _map["tempsetpshifteco"] = new PlcREAL((float) tempsetpshifteco);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling_Speed))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var timePeriod = readBuffer.ReadUshort("timePeriod", 16);
                var percent = readBuffer.ReadByte("percent", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["timePeriod"] = new PlcUINT((ushort) timePeriod);
                _map["percent"] = new PlcUSINT((byte) percent);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling_Step_Time))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var timePeriod = readBuffer.ReadUshort("timePeriod", 16);
                var percent = readBuffer.ReadByte("percent", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["timePeriod"] = new PlcUINT((ushort) timePeriod);
                _map["percent"] = new PlcUSINT((byte) percent);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MeteringValue))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var countval = readBuffer.ReadInt("countval", 32);
                var valinffield = readBuffer.ReadByte("valinffield", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var alarmunack = readBuffer.ReadBit("alarmunack");
                var inalarm = readBuffer.ReadBit("inalarm");
                var overridden = readBuffer.ReadBit("overridden");
                var fault = readBuffer.ReadBit("fault");
                var outofservice = readBuffer.ReadBit("outofservice");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["countval"] = new PlcDINT((int) countval);
                _map["valinffield"] = new PlcUSINT((byte) valinffield);
                _map["alarmunack"] = new PlcBOOL((bool) alarmunack);
                _map["inalarm"] = new PlcBOOL((bool) inalarm);
                _map["overridden"] = new PlcBOOL((bool) overridden);
                _map["fault"] = new PlcBOOL((bool) fault);
                _map["outofservice"] = new PlcBOOL((bool) outofservice);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MBus_Address))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var manufactid = readBuffer.ReadUshort("manufactid", 16);
                var identnumber = readBuffer.ReadUint("identnumber", 32);
                var version = readBuffer.ReadByte("version", 8);
                var medium = readBuffer.ReadByte("medium", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["manufactid"] = new PlcUINT((ushort) manufactid);
                _map["identnumber"] = new PlcUDINT((uint) identnumber);
                _map["version"] = new PlcUSINT((byte) version);
                _map["medium"] = new PlcUSINT((byte) medium);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_RGB))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var r = readBuffer.ReadByte("r", 8);
                var g = readBuffer.ReadByte("g", 8);
                var b = readBuffer.ReadByte("b", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["r"] = new PlcUSINT((byte) r);
                _map["g"] = new PlcUSINT((byte) g);
                _map["b"] = new PlcUSINT((byte) b);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII))
            {
                // STRING
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadString("value", 16, System.Text.Encoding.ASCII);
                return new PlcSTRING(value);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Tariff_ActiveEnergy))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var activeelectricalenergy = readBuffer.ReadInt("activeelectricalenergy", 32);
                var tariff = readBuffer.ReadByte("tariff", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var electricalengergyvalidity = readBuffer.ReadBit("electricalengergyvalidity");
                var tariffvalidity = readBuffer.ReadBit("tariffvalidity");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["activeelectricalenergy"] = new PlcDINT((int) activeelectricalenergy);
                _map["tariff"] = new PlcUSINT((byte) tariff);
                _map["electricalengergyvalidity"] = new PlcBOOL((bool) electricalengergyvalidity);
                _map["tariffvalidity"] = new PlcBOOL((bool) tariffvalidity);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Prioritised_Mode_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var deactivationOfPriority = readBuffer.ReadBit("deactivationOfPriority");
                var priorityLevel = readBuffer.ReadByte("priorityLevel", 3);
                var modeLevel = readBuffer.ReadByte("modeLevel", 4);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["deactivationOfPriority"] = new PlcBOOL((bool) deactivationOfPriority);
                _map["priorityLevel"] = new PlcUSINT((byte) priorityLevel);
                _map["modeLevel"] = new PlcUSINT((byte) modeLevel);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var convertorError = readBuffer.ReadBit("convertorError");
                var ballastFailure = readBuffer.ReadBit("ballastFailure");
                var lampFailure = readBuffer.ReadBit("lampFailure");
                var readOrResponse = readBuffer.ReadBit("readOrResponse");
                var addressIndicator = readBuffer.ReadBit("addressIndicator");
                var daliDeviceAddressOrDaliGroupAddress = readBuffer.ReadByte("daliDeviceAddressOrDaliGroupAddress", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["convertorError"] = new PlcBOOL((bool) convertorError);
                _map["ballastFailure"] = new PlcBOOL((bool) ballastFailure);
                _map["lampFailure"] = new PlcBOOL((bool) lampFailure);
                _map["readOrResponse"] = new PlcBOOL((bool) readOrResponse);
                _map["addressIndicator"] = new PlcBOOL((bool) addressIndicator);
                _map["daliDeviceAddressOrDaliGroupAddress"] = new PlcUSINT((byte) daliDeviceAddressOrDaliGroupAddress);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Diagnostics))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var ballastFailure = readBuffer.ReadBit("ballastFailure");
                var lampFailure = readBuffer.ReadBit("lampFailure");
                var deviceAddress = readBuffer.ReadByte("deviceAddress", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["ballastFailure"] = new PlcBOOL((bool) ballastFailure);
                _map["lampFailure"] = new PlcBOOL((bool) lampFailure);
                _map["deviceAddress"] = new PlcUSINT((byte) deviceAddress);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedPosition))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var heightPosition = readBuffer.ReadByte("heightPosition", 8);
                var slatsPosition = readBuffer.ReadByte("slatsPosition", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var validitySlatsPosition = readBuffer.ReadBit("validitySlatsPosition");
                var validityHeightPosition = readBuffer.ReadBit("validityHeightPosition");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["heightPosition"] = new PlcUSINT((byte) heightPosition);
                _map["slatsPosition"] = new PlcUSINT((byte) slatsPosition);
                _map["validitySlatsPosition"] = new PlcBOOL((bool) validitySlatsPosition);
                _map["validityHeightPosition"] = new PlcBOOL((bool) validityHeightPosition);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusSAB))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var heightPosition = readBuffer.ReadByte("heightPosition", 8);
                var slatsPosition = readBuffer.ReadByte("slatsPosition", 8);
                var upperEndPosReached = readBuffer.ReadBit("upperEndPosReached");
                var lowerEndPosReached = readBuffer.ReadBit("lowerEndPosReached");
                var lowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent = readBuffer.ReadBit("lowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent");
                var targetPosDrive = readBuffer.ReadBit("targetPosDrive");
                var restrictionOfTargetHeightPosPosCanNotBeReached = readBuffer.ReadBit("restrictionOfTargetHeightPosPosCanNotBeReached");
                var restrictionOfSlatsHeightPosPosCanNotBeReached = readBuffer.ReadBit("restrictionOfSlatsHeightPosPosCanNotBeReached");
                var atLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm = readBuffer.ReadBit("atLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm");
                var upDownPositionIsForcedByMoveupdownforcedInput = readBuffer.ReadBit("upDownPositionIsForcedByMoveupdownforcedInput");
                var movementIsLockedEGByDevicelockedInput = readBuffer.ReadBit("movementIsLockedEGByDevicelockedInput");
                var actuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface = readBuffer.ReadBit("actuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface");
                var generalFailureOfTheActuatorOrTheDrive = readBuffer.ReadBit("generalFailureOfTheActuatorOrTheDrive");
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var validityHeightPos = readBuffer.ReadBit("validityHeightPos");
                var validitySlatsPos = readBuffer.ReadBit("validitySlatsPos");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["heightPosition"] = new PlcUSINT((byte) heightPosition);
                _map["slatsPosition"] = new PlcUSINT((byte) slatsPosition);
                _map["upperEndPosReached"] = new PlcBOOL((bool) upperEndPosReached);
                _map["lowerEndPosReached"] = new PlcBOOL((bool) lowerEndPosReached);
                _map["lowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent"] = new PlcBOOL((bool) lowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent);
                _map["targetPosDrive"] = new PlcBOOL((bool) targetPosDrive);
                _map["restrictionOfTargetHeightPosPosCanNotBeReached"] = new PlcBOOL((bool) restrictionOfTargetHeightPosPosCanNotBeReached);
                _map["restrictionOfSlatsHeightPosPosCanNotBeReached"] = new PlcBOOL((bool) restrictionOfSlatsHeightPosPosCanNotBeReached);
                _map["atLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm"] = new PlcBOOL((bool) atLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm);
                _map["upDownPositionIsForcedByMoveupdownforcedInput"] = new PlcBOOL((bool) upDownPositionIsForcedByMoveupdownforcedInput);
                _map["movementIsLockedEGByDevicelockedInput"] = new PlcBOOL((bool) movementIsLockedEGByDevicelockedInput);
                _map["actuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface"] = new PlcBOOL((bool) actuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface);
                _map["generalFailureOfTheActuatorOrTheDrive"] = new PlcBOOL((bool) generalFailureOfTheActuatorOrTheDrive);
                _map["validityHeightPos"] = new PlcBOOL((bool) validityHeightPos);
                _map["validitySlatsPos"] = new PlcBOOL((bool) validitySlatsPos);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_xyY))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var xAxis = readBuffer.ReadUshort("xAxis", 16);
                var yAxis = readBuffer.ReadUshort("yAxis", 16);
                var brightness = readBuffer.ReadByte("brightness", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var validityXy = readBuffer.ReadBit("validityXy");
                var validityBrightness = readBuffer.ReadBit("validityBrightness");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["xAxis"] = new PlcUINT((ushort) xAxis);
                _map["yAxis"] = new PlcUINT((ushort) yAxis);
                _map["brightness"] = new PlcUSINT((byte) brightness);
                _map["validityXy"] = new PlcBOOL((bool) validityXy);
                _map["validityBrightness"] = new PlcBOOL((bool) validityBrightness);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Status))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var converterModeAccordingToTheDaliConverterStateMachine = readBuffer.ReadByte("converterModeAccordingToTheDaliConverterStateMachine", 4);
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var hardwiredSwitchIsActive = readBuffer.ReadBit("hardwiredSwitchIsActive");
                var hardwiredInhibitIsActive = readBuffer.ReadBit("hardwiredInhibitIsActive");
                var functionTestPending = readBuffer.ReadByte("functionTestPending", 2);
                var durationTestPending = readBuffer.ReadByte("durationTestPending", 2);
                var partialDurationTestPending = readBuffer.ReadByte("partialDurationTestPending", 2);
                var converterFailure = readBuffer.ReadByte("converterFailure", 2);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["converterModeAccordingToTheDaliConverterStateMachine"] = new PlcUSINT((byte) converterModeAccordingToTheDaliConverterStateMachine);
                _map["hardwiredSwitchIsActive"] = new PlcBOOL((bool) hardwiredSwitchIsActive);
                _map["hardwiredInhibitIsActive"] = new PlcBOOL((bool) hardwiredInhibitIsActive);
                _map["functionTestPending"] = new PlcUSINT((byte) functionTestPending);
                _map["durationTestPending"] = new PlcUSINT((byte) durationTestPending);
                _map["partialDurationTestPending"] = new PlcUSINT((byte) partialDurationTestPending);
                _map["converterFailure"] = new PlcUSINT((byte) converterFailure);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Test_Result))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var ltrf = readBuffer.ReadByte("ltrf", 4);
                var ltrd = readBuffer.ReadByte("ltrd", 4);
                var ltrp = readBuffer.ReadByte("ltrp", 4);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var sf = readBuffer.ReadByte("sf", 2);
                var sd = readBuffer.ReadByte("sd", 2);
                var sp = readBuffer.ReadByte("sp", 2);
                {
                    var reserved = readBuffer.ReadByte("reserved", 2);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var ldtr = readBuffer.ReadUshort("ldtr", 16);
                var lpdtr = readBuffer.ReadByte("lpdtr", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["ltrf"] = new PlcUSINT((byte) ltrf);
                _map["ltrd"] = new PlcUSINT((byte) ltrd);
                _map["ltrp"] = new PlcUSINT((byte) ltrp);
                _map["sf"] = new PlcUSINT((byte) sf);
                _map["sd"] = new PlcUSINT((byte) sd);
                _map["sp"] = new PlcUSINT((byte) sp);
                _map["ldtr"] = new PlcUINT((ushort) ldtr);
                _map["lpdtr"] = new PlcUSINT((byte) lpdtr);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Battery_Info))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var batteryFailure = readBuffer.ReadBit("batteryFailure");
                var batteryDurationFailure = readBuffer.ReadBit("batteryDurationFailure");
                var batteryFullyCharged = readBuffer.ReadBit("batteryFullyCharged");
                var batteryChargeLevel = readBuffer.ReadByte("batteryChargeLevel", 8);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["batteryFailure"] = new PlcBOOL((bool) batteryFailure);
                _map["batteryDurationFailure"] = new PlcBOOL((bool) batteryDurationFailure);
                _map["batteryFullyCharged"] = new PlcBOOL((bool) batteryFullyCharged);
                _map["batteryChargeLevel"] = new PlcUSINT((byte) batteryChargeLevel);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness_Colour_Temperature_Transition))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var ms = readBuffer.ReadUshort("ms", 16);
                var temperatureK = readBuffer.ReadUshort("temperatureK", 16);
                var percent = readBuffer.ReadByte("percent", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 5);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var validityOfTheTimePeriod = readBuffer.ReadBit("validityOfTheTimePeriod");
                var validityOfTheAbsoluteColourTemperature = readBuffer.ReadBit("validityOfTheAbsoluteColourTemperature");
                var validityOfTheAbsoluteBrightness = readBuffer.ReadBit("validityOfTheAbsoluteBrightness");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["ms"] = new PlcUINT((ushort) ms);
                _map["temperatureK"] = new PlcUINT((ushort) temperatureK);
                _map["percent"] = new PlcUSINT((byte) percent);
                _map["validityOfTheTimePeriod"] = new PlcBOOL((bool) validityOfTheTimePeriod);
                _map["validityOfTheAbsoluteColourTemperature"] = new PlcBOOL((bool) validityOfTheAbsoluteColourTemperature);
                _map["validityOfTheAbsoluteBrightness"] = new PlcBOOL((bool) validityOfTheAbsoluteBrightness);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness_Colour_Temperature_Control))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cct = readBuffer.ReadBit("cct");
                var stepCodeColourTemperature = readBuffer.ReadByte("stepCodeColourTemperature", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cb = readBuffer.ReadBit("cb");
                var stepCodeBrightness = readBuffer.ReadByte("stepCodeBrightness", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 6);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cctAndStepCodeColourValidity = readBuffer.ReadBit("cctAndStepCodeColourValidity");
                var cbAndStepCodeBrightnessValidity = readBuffer.ReadBit("cbAndStepCodeBrightnessValidity");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["cct"] = new PlcBOOL((bool) cct);
                _map["stepCodeColourTemperature"] = new PlcUSINT((byte) stepCodeColourTemperature);
                _map["cb"] = new PlcBOOL((bool) cb);
                _map["stepCodeBrightness"] = new PlcUSINT((byte) stepCodeBrightness);
                _map["cctAndStepCodeColourValidity"] = new PlcBOOL((bool) cctAndStepCodeColourValidity);
                _map["cbAndStepCodeBrightnessValidity"] = new PlcBOOL((bool) cbAndStepCodeBrightnessValidity);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_RGBW))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var colourLevelRed = readBuffer.ReadByte("colourLevelRed", 8);
                var colourLevelGreen = readBuffer.ReadByte("colourLevelGreen", 8);
                var colourLevelBlue = readBuffer.ReadByte("colourLevelBlue", 8);
                var colourLevelWhite = readBuffer.ReadByte("colourLevelWhite", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var mr = readBuffer.ReadBit("mr");
                var mg = readBuffer.ReadBit("mg");
                var mb = readBuffer.ReadBit("mb");
                var mw = readBuffer.ReadBit("mw");
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["colourLevelRed"] = new PlcUSINT((byte) colourLevelRed);
                _map["colourLevelGreen"] = new PlcUSINT((byte) colourLevelGreen);
                _map["colourLevelBlue"] = new PlcUSINT((byte) colourLevelBlue);
                _map["colourLevelWhite"] = new PlcUSINT((byte) colourLevelWhite);
                _map["mr"] = new PlcBOOL((bool) mr);
                _map["mg"] = new PlcBOOL((bool) mg);
                _map["mb"] = new PlcBOOL((bool) mb);
                _map["mw"] = new PlcBOOL((bool) mw);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Relative_Control_RGBW))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var maskcw = readBuffer.ReadBit("maskcw");
                var maskcb = readBuffer.ReadBit("maskcb");
                var maskcg = readBuffer.ReadBit("maskcg");
                var maskcr = readBuffer.ReadBit("maskcr");
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cw = readBuffer.ReadBit("cw");
                var stepCodeColourWhite = readBuffer.ReadByte("stepCodeColourWhite", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cb = readBuffer.ReadBit("cb");
                var stepCodeColourBlue = readBuffer.ReadByte("stepCodeColourBlue", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cg = readBuffer.ReadBit("cg");
                var stepCodeColourGreen = readBuffer.ReadByte("stepCodeColourGreen", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cr = readBuffer.ReadBit("cr");
                var stepCodeColourRed = readBuffer.ReadByte("stepCodeColourRed", 3);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["maskcw"] = new PlcBOOL((bool) maskcw);
                _map["maskcb"] = new PlcBOOL((bool) maskcb);
                _map["maskcg"] = new PlcBOOL((bool) maskcg);
                _map["maskcr"] = new PlcBOOL((bool) maskcr);
                _map["cw"] = new PlcBOOL((bool) cw);
                _map["stepCodeColourWhite"] = new PlcUSINT((byte) stepCodeColourWhite);
                _map["cb"] = new PlcBOOL((bool) cb);
                _map["stepCodeColourBlue"] = new PlcUSINT((byte) stepCodeColourBlue);
                _map["cg"] = new PlcBOOL((bool) cg);
                _map["stepCodeColourGreen"] = new PlcUSINT((byte) stepCodeColourGreen);
                _map["cr"] = new PlcBOOL((bool) cr);
                _map["stepCodeColourRed"] = new PlcUSINT((byte) stepCodeColourRed);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Relative_Control_RGB))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cb = readBuffer.ReadBit("cb");
                var stepCodeColourBlue = readBuffer.ReadByte("stepCodeColourBlue", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cg = readBuffer.ReadBit("cg");
                var stepCodeColourGreen = readBuffer.ReadByte("stepCodeColourGreen", 3);
                {
                    var reserved = readBuffer.ReadByte("reserved", 4);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var cr = readBuffer.ReadBit("cr");
                var stepCodeColourRed = readBuffer.ReadByte("stepCodeColourRed", 3);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["cb"] = new PlcBOOL((bool) cb);
                _map["stepCodeColourBlue"] = new PlcUSINT((byte) stepCodeColourBlue);
                _map["cg"] = new PlcBOOL((bool) cg);
                _map["stepCodeColourGreen"] = new PlcUSINT((byte) stepCodeColourGreen);
                _map["cr"] = new PlcBOOL((bool) cr);
                _map["stepCodeColourRed"] = new PlcUSINT((byte) stepCodeColourRed);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_GeographicalLocation))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var longitude = readBuffer.ReadFloat("longitude", 32);
                var latitude = readBuffer.ReadFloat("latitude", 32);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["longitude"] = new PlcREAL((float) longitude);
                _map["latitude"] = new PlcREAL((float) latitude);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetF16_4))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var roomTemperatureSetpointComfort = readBuffer.ReadFloat("roomTemperatureSetpointComfort", 16);
                var roomTemperatureSetpointStandby = readBuffer.ReadFloat("roomTemperatureSetpointStandby", 16);
                var roomTemperatureSetpointEconomy = readBuffer.ReadFloat("roomTemperatureSetpointEconomy", 16);
                var roomTemperatureSetpointBuildingProtection = readBuffer.ReadFloat("roomTemperatureSetpointBuildingProtection", 16);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["roomTemperatureSetpointComfort"] = new PlcREAL((float) roomTemperatureSetpointComfort);
                _map["roomTemperatureSetpointStandby"] = new PlcREAL((float) roomTemperatureSetpointStandby);
                _map["roomTemperatureSetpointEconomy"] = new PlcREAL((float) roomTemperatureSetpointEconomy);
                _map["roomTemperatureSetpointBuildingProtection"] = new PlcREAL((float) roomTemperatureSetpointBuildingProtection);
                return new PlcStruct(_map);
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetShiftF16_4))
            {
                // Struct
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var roomTemperatureSetpointShiftComfort = readBuffer.ReadFloat("roomTemperatureSetpointShiftComfort", 16);
                var roomTemperatureSetpointShiftStandby = readBuffer.ReadFloat("roomTemperatureSetpointShiftStandby", 16);
                var roomTemperatureSetpointShiftEconomy = readBuffer.ReadFloat("roomTemperatureSetpointShiftEconomy", 16);
                var roomTemperatureSetpointShiftBuildingProtection = readBuffer.ReadFloat("roomTemperatureSetpointShiftBuildingProtection", 16);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["roomTemperatureSetpointShiftComfort"] = new PlcREAL((float) roomTemperatureSetpointShiftComfort);
                _map["roomTemperatureSetpointShiftStandby"] = new PlcREAL((float) roomTemperatureSetpointShiftStandby);
                _map["roomTemperatureSetpointShiftEconomy"] = new PlcREAL((float) roomTemperatureSetpointShiftEconomy);
                _map["roomTemperatureSetpointShiftBuildingProtection"] = new PlcREAL((float) roomTemperatureSetpointShiftBuildingProtection);
                return new PlcStruct(_map);
            }
            return new PlcNULL();
        }

        public static void StaticSerialize(WriteBuffer writeBuffer, IPlcValue _value, KnxDatapointType datapointType)
        {
            if (Equals(datapointType, KnxDatapointType.BOOL))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.BYTE))
            {
                // BYTE
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.WORD))
            {
                // WORD
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DWORD))
            {
                // DWORD
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.LWORD))
            {
                // LWORD
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(datapointType, KnxDatapointType.USINT))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.SINT))
            {
                // SINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(datapointType, KnxDatapointType.UINT))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.INT))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.UDINT))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DINT))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.ULINT))
            {
                // ULINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(datapointType, KnxDatapointType.LINT))
            {
                // LINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(datapointType, KnxDatapointType.REAL))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.LREAL))
            {
                // LREAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteDouble("value", 64, (double) _value.GetDouble());
            }
            else if (Equals(datapointType, KnxDatapointType.CHAR))
            {
                // CHAR
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 8, "UTF8", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.WCHAR))
            {
                // WCHAR
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 16, "UTF16", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.TIME))
            {
                // TIME
                throw new NotImplementedException("KnxDatapoint 'TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.LTIME))
            {
                // LTIME
                throw new NotImplementedException("KnxDatapoint 'LTIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DATE))
            {
                // DATE
                throw new NotImplementedException("KnxDatapoint 'DATE' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.TIME_OF_DAY))
            {
                // TIME_OF_DAY
                throw new NotImplementedException("KnxDatapoint 'TIME_OF_DAY' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.TOD))
            {
                // TIME_OF_DAY
                throw new NotImplementedException("KnxDatapoint 'TIME_OF_DAY' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DATE_AND_TIME))
            {
                // DATE_AND_TIME
                throw new NotImplementedException("KnxDatapoint 'DATE_AND_TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DT))
            {
                // DATE_AND_TIME
                throw new NotImplementedException("KnxDatapoint 'DATE_AND_TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Switch))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Bool))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Enable))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ramp))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BinaryValue))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Step))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UpDown))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OpenClose))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Start))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_State))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Invert))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DimSendStyle))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_InputSource))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Reset))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ack))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Trigger))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Occupancy))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Window_Door))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LogicalFunction))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scene_AB))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ShutterBlinds_Mode))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DayNight))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Heat_Cool))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Switch_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("on", (bool) _value.GetValue("on").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Bool_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("valueTrue", (bool) _value.GetValue("valueTrue").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Enable_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("enable", (bool) _value.GetValue("enable").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ramp_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("ramp", (bool) _value.GetValue("ramp").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("alarm", (bool) _value.GetValue("alarm").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BinaryValue_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("high", (bool) _value.GetValue("high").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Step_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("increase", (bool) _value.GetValue("increase").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Direction1_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("down", (bool) _value.GetValue("down").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Direction2_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("close", (bool) _value.GetValue("close").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Start_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("start", (bool) _value.GetValue("start").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_State_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("active", (bool) _value.GetValue("active").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Invert_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("control", (bool) _value.GetValue("control").GetBool());
                writeBuffer.WriteBit("inverted", (bool) _value.GetValue("inverted").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Control_Dimming))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("increase", (bool) _value.GetValue("increase").GetBool());
                writeBuffer.WriteByte("stepcode", 3, (byte) _value.GetValue("stepcode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Control_Blinds))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("down", (bool) _value.GetValue("down").GetBool());
                writeBuffer.WriteByte("stepcode", 3, (byte) _value.GetValue("stepcode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Char_ASCII))
            {
                // STRING
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 8, "ASCII", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Char_8859_1))
            {
                // STRING
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 8, "ISO-8859-1", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Angle))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_U8))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DecimalFactor))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Tariff))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_1_Ucount))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FanStage))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_V8))
            {
                // SINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_1_Count))
            {
                // SINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Status_Mode3))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("statusA", (bool) _value.GetValue("statusA").GetBool());
                writeBuffer.WriteBit("statusB", (bool) _value.GetValue("statusB").GetBool());
                writeBuffer.WriteBit("statusC", (bool) _value.GetValue("statusC").GetBool());
                writeBuffer.WriteBit("statusD", (bool) _value.GetValue("statusD").GetBool());
                writeBuffer.WriteBit("statusE", (bool) _value.GetValue("statusE").GetBool());
                writeBuffer.WriteByte("mode", 3, (byte) _value.GetValue("mode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_2_Ucount))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodMsec))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriod10Msec))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriod100Msec))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodSec))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodMin))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodHrs))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PropDataType))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Length_mm))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UElCurrentmA))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Absolute_Colour_Temperature))
            {
                // UINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_2_Count))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeMsec))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTime10Msec))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTime100Msec))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeSec))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeMin))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeHrs))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_V16))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Rotation_Angle))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Length_m))
            {
                // INT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Temp))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Tempd))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Tempa))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Lux))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Wsp))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Pres))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Humidity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AirQuality))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AirFlow))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time1))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time2))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volt))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Curr))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PowerDensity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_KelvinPerPercent))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Power))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume_Flow))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Rain_Amount))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Temp_F))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Wsp_kmh))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Absolute_Humidity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Concentration_ygm3))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Coefficient))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimeOfDay))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("day", 3, (byte) _value.GetValue("day").GetByte());
                writeBuffer.WriteByte("hour", 5, (byte) _value.GetValue("hour").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("minutes", 6, (byte) _value.GetValue("minutes").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("seconds", 6, (byte) _value.GetValue("seconds").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Date))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteByte("dayOfMonth", 5, (byte) _value.GetValue("dayOfMonth").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteByte("month", 4, (byte) _value.GetValue("month").GetByte());
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteByte("year", 7, (byte) _value.GetValue("year").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_4_Ucount))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Sec))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Min))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Hrs))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_VolumeLiquid_Litre))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_m_3))
            {
                // UDINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_4_Count))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FlowRate_m3h))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_kWh))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy_kVAh))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy_kVARh))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_MWh))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongDeltaTimeSec))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaVolumeLiquid_Litre))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaVolume_m_3))
            {
                // DINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Acceleration))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Acceleration_Angular))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Activation_Energy))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Activity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Mol))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Amplitude))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AngleRad))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AngleDeg))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Momentum))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Velocity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Area))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Capacitance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Charge_DensitySurface))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Charge_DensityVolume))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Compressibility))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Conductance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electrical_Conductivity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Density))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Charge))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Current))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_CurrentDensity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_DipoleMoment))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Displacement))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_FieldStrength))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Flux))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_FluxDensity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Polarization))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Potential))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_PotentialDifference))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ElectromagneticMoment))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electromotive_Force))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Energy))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Force))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Frequency))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Frequency))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_Capacity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_FlowRate))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_Quantity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Impedance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Length))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Light_Quantity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminous_Flux))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminous_Intensity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_FieldStrength))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Flux))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_FluxDensity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Moment))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Polarization))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetization))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_MagnetomotiveForce))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Mass))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_MassFlux))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Momentum))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Phase_AngleRad))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Phase_AngleDeg))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Power))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Power_Factor))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Pressure))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Reactance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Resistance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Resistivity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_SelfInductance))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_SolidAngle))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Sound_Intensity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Speed))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Stress))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Surface_Tension))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Common_Temperature))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Absolute_Temperature))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_TemperatureDifference))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Thermal_Capacity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Thermal_Conductivity))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ThermoelectricPower))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Torque))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume_Flux))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Weight))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Work))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ApparentPower))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_Flux_Meter))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_Flux_ls))
            {
                // REAL
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Access_Data))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("hurz", 4, (byte) _value.GetValue("hurz").GetByte());
                writeBuffer.WriteByte("value1", 4, (byte) _value.GetValue("value1").GetByte());
                writeBuffer.WriteByte("value2", 4, (byte) _value.GetValue("value2").GetByte());
                writeBuffer.WriteByte("value3", 4, (byte) _value.GetValue("value3").GetByte());
                writeBuffer.WriteByte("value4", 4, (byte) _value.GetValue("value4").GetByte());
                writeBuffer.WriteByte("value5", 4, (byte) _value.GetValue("value5").GetByte());
                writeBuffer.WriteBit("detectionError", (bool) _value.GetValue("detectionError").GetBool());
                writeBuffer.WriteBit("permission", (bool) _value.GetValue("permission").GetBool());
                writeBuffer.WriteBit("readDirection", (bool) _value.GetValue("readDirection").GetBool());
                writeBuffer.WriteBit("encryptionOfAccessInformation", (bool) _value.GetValue("encryptionOfAccessInformation").GetBool());
                writeBuffer.WriteByte("indexOfAccessIdentificationCode", 4, (byte) _value.GetValue("indexOfAccessIdentificationCode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_String_ASCII))
            {
                // STRING
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 112, "ASCII", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_String_8859_1))
            {
                // STRING
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 112, "ISO-8859-1", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneNumber))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("value", 6, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneControl))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("learnTheSceneCorrespondingToTheFieldSceneNumber", (bool) _value.GetValue("learnTheSceneCorrespondingToTheFieldSceneNumber").GetBool());
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteByte("sceneNumber", 6, (byte) _value.GetValue("sceneNumber").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DateTime))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("year", 8, (byte) _value.GetValue("year").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteByte("month", 4, (byte) _value.GetValue("month").GetByte());
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteByte("dayofmonth", 5, (byte) _value.GetValue("dayofmonth").GetByte());
                writeBuffer.WriteByte("dayofweek", 3, (byte) _value.GetValue("dayofweek").GetByte());
                writeBuffer.WriteByte("hourofday", 5, (byte) _value.GetValue("hourofday").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("minutes", 6, (byte) _value.GetValue("minutes").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("seconds", 6, (byte) _value.GetValue("seconds").GetByte());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
                writeBuffer.WriteBit("workingDay", (bool) _value.GetValue("workingDay").GetBool());
                writeBuffer.WriteBit("noWd", (bool) _value.GetValue("noWd").GetBool());
                writeBuffer.WriteBit("noYear", (bool) _value.GetValue("noYear").GetBool());
                writeBuffer.WriteBit("noDate", (bool) _value.GetValue("noDate").GetBool());
                writeBuffer.WriteBit("noDayOfWeek", (bool) _value.GetValue("noDayOfWeek").GetBool());
                writeBuffer.WriteBit("noTime", (bool) _value.GetValue("noTime").GetBool());
                writeBuffer.WriteBit("standardSummerTime", (bool) _value.GetValue("standardSummerTime").GetBool());
                writeBuffer.WriteBit("qualityOfClock", (bool) _value.GetValue("qualityOfClock").GetBool());
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SCLOMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BuildingMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OccMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Priority))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightApplicationMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApplicationArea))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AlarmClassType))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PSUMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ErrorClass_System))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ErrorClass_HVAC))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Time_Delay))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Beaufort_Wind_Force_Scale))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SensorSelect))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActuatorConnectType))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Cloud_Cover))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PowerReturnMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FuelType))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BurnerType))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DHWMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadPriority))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACContrMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACEmergMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ChangeoverMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ValveMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DamperMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HeaterMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FanMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MasterSlaveMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRoomSetp))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Metering_DeviceType))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HumDehumMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EnableHCStage))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ADAType))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BackupMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StartSynchronization))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Behaviour_Lock_Unlock))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Fade_Time))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BlinkingMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightControlMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SwitchPBModel))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PBAction))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DimmPBModel))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SwitchOnMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadTypeSet))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadTypeDetected))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Test_Control))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Control))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SABExcept_Behaviour))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SABBehaviour_Lock_Unlock))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SSSBMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BlindsControlMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CommMode))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AddInfoTypes))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_ModeSelect))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_FilterSelect))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_1))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_2))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_3))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusGen))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteBit("alarmStatusOfCorrespondingDatapointIsNotAcknowledged", (bool) _value.GetValue("alarmStatusOfCorrespondingDatapointIsNotAcknowledged").GetBool());
                writeBuffer.WriteBit("correspondingDatapointIsInAlarm", (bool) _value.GetValue("correspondingDatapointIsInAlarm").GetBool());
                writeBuffer.WriteBit("correspondingDatapointMainValueIsOverridden", (bool) _value.GetValue("correspondingDatapointMainValueIsOverridden").GetBool());
                writeBuffer.WriteBit("correspondingDatapointMainValueIsCorruptedDueToFailure", (bool) _value.GetValue("correspondingDatapointMainValueIsCorruptedDueToFailure").GetBool());
                writeBuffer.WriteBit("correspondingDatapointValueIsOutOfService", (bool) _value.GetValue("correspondingDatapointValueIsOutOfService").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Device_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("verifyModeIsOn", (bool) _value.GetValue("verifyModeIsOn").GetBool());
                writeBuffer.WriteBit("aDatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived", (bool) _value.GetValue("aDatagramWithTheOwnIndividualAddressAsSourceAddressHasBeenReceived").GetBool());
                writeBuffer.WriteBit("theUserApplicationIsStopped", (bool) _value.GetValue("theUserApplicationIsStopped").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ForceSign))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("roomhmax", (bool) _value.GetValue("roomhmax").GetBool());
                writeBuffer.WriteBit("roomhconf", (bool) _value.GetValue("roomhconf").GetBool());
                writeBuffer.WriteBit("dhwlegio", (bool) _value.GetValue("dhwlegio").GetBool());
                writeBuffer.WriteBit("dhwnorm", (bool) _value.GetValue("dhwnorm").GetBool());
                writeBuffer.WriteBit("overrun", (bool) _value.GetValue("overrun").GetBool());
                writeBuffer.WriteBit("oversupply", (bool) _value.GetValue("oversupply").GetBool());
                writeBuffer.WriteBit("protection", (bool) _value.GetValue("protection").GetBool());
                writeBuffer.WriteBit("forcerequest", (bool) _value.GetValue("forcerequest").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ForceSignCool))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRHC))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("summermode", (bool) _value.GetValue("summermode").GetBool());
                writeBuffer.WriteBit("statusstopoptim", (bool) _value.GetValue("statusstopoptim").GetBool());
                writeBuffer.WriteBit("statusstartoptim", (bool) _value.GetValue("statusstartoptim").GetBool());
                writeBuffer.WriteBit("statusmorningboost", (bool) _value.GetValue("statusmorningboost").GetBool());
                writeBuffer.WriteBit("tempreturnlimit", (bool) _value.GetValue("tempreturnlimit").GetBool());
                writeBuffer.WriteBit("tempflowlimit", (bool) _value.GetValue("tempflowlimit").GetBool());
                writeBuffer.WriteBit("satuseco", (bool) _value.GetValue("satuseco").GetBool());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusSDHWC))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("solarloadsufficient", (bool) _value.GetValue("solarloadsufficient").GetBool());
                writeBuffer.WriteBit("sdhwloadactive", (bool) _value.GetValue("sdhwloadactive").GetBool());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FuelTypeSet))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("solidstate", (bool) _value.GetValue("solidstate").GetBool());
                writeBuffer.WriteBit("gas", (bool) _value.GetValue("gas").GetBool());
                writeBuffer.WriteBit("oil", (bool) _value.GetValue("oil").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRCC))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusAHU))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cool", (bool) _value.GetValue("cool").GetBool());
                writeBuffer.WriteBit("heat", (bool) _value.GetValue("heat").GetBool());
                writeBuffer.WriteBit("fanactive", (bool) _value.GetValue("fanactive").GetBool());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_RTSM))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteBit("statusOfHvacModeUser", (bool) _value.GetValue("statusOfHvacModeUser").GetBool());
                writeBuffer.WriteBit("statusOfComfortProlongationUser", (bool) _value.GetValue("statusOfComfortProlongationUser").GetBool());
                writeBuffer.WriteBit("effectiveValueOfTheComfortPushButton", (bool) _value.GetValue("effectiveValueOfTheComfortPushButton").GetBool());
                writeBuffer.WriteBit("effectiveValueOfThePresenceStatus", (bool) _value.GetValue("effectiveValueOfThePresenceStatus").GetBool());
                writeBuffer.WriteBit("effectiveValueOfTheWindowStatus", (bool) _value.GetValue("effectiveValueOfTheWindowStatus").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightActuatorErrorInfo))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteBit("overheat", (bool) _value.GetValue("overheat").GetBool());
                writeBuffer.WriteBit("lampfailure", (bool) _value.GetValue("lampfailure").GetBool());
                writeBuffer.WriteBit("defectiveload", (bool) _value.GetValue("defectiveload").GetBool());
                writeBuffer.WriteBit("underload", (bool) _value.GetValue("underload").GetBool());
                writeBuffer.WriteBit("overcurrent", (bool) _value.GetValue("overcurrent").GetBool());
                writeBuffer.WriteBit("undervoltage", (bool) _value.GetValue("undervoltage").GetBool());
                writeBuffer.WriteBit("loaddetectionerror", (bool) _value.GetValue("loaddetectionerror").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_ModeInfo))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("bibatSlave", (bool) _value.GetValue("bibatSlave").GetBool());
                writeBuffer.WriteBit("bibatMaster", (bool) _value.GetValue("bibatMaster").GetBool());
                writeBuffer.WriteBit("asynchronous", (bool) _value.GetValue("asynchronous").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_FilterInfo))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("doa", (bool) _value.GetValue("doa").GetBool());
                writeBuffer.WriteBit("knxSn", (bool) _value.GetValue("knxSn").GetBool());
                writeBuffer.WriteBit("doaAndKnxSn", (bool) _value.GetValue("doaAndKnxSn").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_8))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("activationStateOfChannel1", (bool) _value.GetValue("activationStateOfChannel1").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel2", (bool) _value.GetValue("activationStateOfChannel2").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel3", (bool) _value.GetValue("activationStateOfChannel3").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel4", (bool) _value.GetValue("activationStateOfChannel4").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel5", (bool) _value.GetValue("activationStateOfChannel5").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel6", (bool) _value.GetValue("activationStateOfChannel6").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel7", (bool) _value.GetValue("activationStateOfChannel7").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel8", (bool) _value.GetValue("activationStateOfChannel8").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusDHWC))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("tempoptimshiftactive", (bool) _value.GetValue("tempoptimshiftactive").GetBool());
                writeBuffer.WriteBit("solarenergysupport", (bool) _value.GetValue("solarenergysupport").GetBool());
                writeBuffer.WriteBit("solarenergyonly", (bool) _value.GetValue("solarenergyonly").GetBool());
                writeBuffer.WriteBit("otherenergysourceactive", (bool) _value.GetValue("otherenergysourceactive").GetBool());
                writeBuffer.WriteBit("dhwpushactive", (bool) _value.GetValue("dhwpushactive").GetBool());
                writeBuffer.WriteBit("legioprotactive", (bool) _value.GetValue("legioprotactive").GetBool());
                writeBuffer.WriteBit("dhwloadactive", (bool) _value.GetValue("dhwloadactive").GetBool());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRHCC))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteBit("overheatalarm", (bool) _value.GetValue("overheatalarm").GetBool());
                writeBuffer.WriteBit("frostalarm", (bool) _value.GetValue("frostalarm").GetBool());
                writeBuffer.WriteBit("dewpointstatus", (bool) _value.GetValue("dewpointstatus").GetBool());
                writeBuffer.WriteBit("coolingdisabled", (bool) _value.GetValue("coolingdisabled").GetBool());
                writeBuffer.WriteBit("statusprecool", (bool) _value.GetValue("statusprecool").GetBool());
                writeBuffer.WriteBit("statusecoc", (bool) _value.GetValue("statusecoc").GetBool());
                writeBuffer.WriteBit("heatcoolmode", (bool) _value.GetValue("heatcoolmode").GetBool());
                writeBuffer.WriteBit("heatingdiabled", (bool) _value.GetValue("heatingdiabled").GetBool());
                writeBuffer.WriteBit("statusstopoptim", (bool) _value.GetValue("statusstopoptim").GetBool());
                writeBuffer.WriteBit("statusstartoptim", (bool) _value.GetValue("statusstartoptim").GetBool());
                writeBuffer.WriteBit("statusmorningboosth", (bool) _value.GetValue("statusmorningboosth").GetBool());
                writeBuffer.WriteBit("tempflowreturnlimit", (bool) _value.GetValue("tempflowreturnlimit").GetBool());
                writeBuffer.WriteBit("tempflowlimit", (bool) _value.GetValue("tempflowlimit").GetBool());
                writeBuffer.WriteBit("statusecoh", (bool) _value.GetValue("statusecoh").GetBool());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_HVA))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("calibrationMode", (bool) _value.GetValue("calibrationMode").GetBool());
                writeBuffer.WriteBit("lockedPosition", (bool) _value.GetValue("lockedPosition").GetBool());
                writeBuffer.WriteBit("forcedPosition", (bool) _value.GetValue("forcedPosition").GetBool());
                writeBuffer.WriteBit("manuaOperationOverridden", (bool) _value.GetValue("manuaOperationOverridden").GetBool());
                writeBuffer.WriteBit("serviceMode", (bool) _value.GetValue("serviceMode").GetBool());
                writeBuffer.WriteBit("valveKick", (bool) _value.GetValue("valveKick").GetBool());
                writeBuffer.WriteBit("overload", (bool) _value.GetValue("overload").GetBool());
                writeBuffer.WriteBit("shortCircuit", (bool) _value.GetValue("shortCircuit").GetBool());
                writeBuffer.WriteBit("currentValvePosition", (bool) _value.GetValue("currentValvePosition").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_RTC))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("coolingModeEnabled", (bool) _value.GetValue("coolingModeEnabled").GetBool());
                writeBuffer.WriteBit("heatingModeEnabled", (bool) _value.GetValue("heatingModeEnabled").GetBool());
                writeBuffer.WriteBit("additionalHeatingCoolingStage2Stage", (bool) _value.GetValue("additionalHeatingCoolingStage2Stage").GetBool());
                writeBuffer.WriteBit("controllerInactive", (bool) _value.GetValue("controllerInactive").GetBool());
                writeBuffer.WriteBit("overheatAlarm", (bool) _value.GetValue("overheatAlarm").GetBool());
                writeBuffer.WriteBit("frostAlarm", (bool) _value.GetValue("frostAlarm").GetBool());
                writeBuffer.WriteBit("dewPointStatus", (bool) _value.GetValue("dewPointStatus").GetBool());
                writeBuffer.WriteBit("activeMode", (bool) _value.GetValue("activeMode").GetBool());
                writeBuffer.WriteBit("generalFailureInformation", (bool) _value.GetValue("generalFailureInformation").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Media))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("reserved", 10, (ushort) (0x00));
                writeBuffer.WriteBit("knxIp", (bool) _value.GetValue("knxIp").GetBool());
                writeBuffer.WriteBit("rf", (bool) _value.GetValue("rf").GetBool());
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteBit("pl110", (bool) _value.GetValue("pl110").GetBool());
                writeBuffer.WriteBit("tp1", (bool) _value.GetValue("tp1").GetBool());
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_16))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("activationStateOfChannel1", (bool) _value.GetValue("activationStateOfChannel1").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel2", (bool) _value.GetValue("activationStateOfChannel2").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel3", (bool) _value.GetValue("activationStateOfChannel3").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel4", (bool) _value.GetValue("activationStateOfChannel4").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel5", (bool) _value.GetValue("activationStateOfChannel5").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel6", (bool) _value.GetValue("activationStateOfChannel6").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel7", (bool) _value.GetValue("activationStateOfChannel7").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel8", (bool) _value.GetValue("activationStateOfChannel8").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel9", (bool) _value.GetValue("activationStateOfChannel9").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel10", (bool) _value.GetValue("activationStateOfChannel10").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel11", (bool) _value.GetValue("activationStateOfChannel11").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel12", (bool) _value.GetValue("activationStateOfChannel12").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel13", (bool) _value.GetValue("activationStateOfChannel13").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel14", (bool) _value.GetValue("activationStateOfChannel14").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel15", (bool) _value.GetValue("activationStateOfChannel15").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel16", (bool) _value.GetValue("activationStateOfChannel16").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OnOffAction))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteByte("value", 2, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm_Reaction))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteByte("value", 2, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UpDown_Action))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteByte("value", 2, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVAC_PB_Action))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteByte("value", 2, (byte) _value.GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DoubleNibble))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("busy", 4, (byte) _value.GetValue("busy").GetByte());
                writeBuffer.WriteByte("nak", 4, (byte) _value.GetValue("nak").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneInfo))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteBit("sceneIsInactive", (bool) _value.GetValue("sceneIsInactive").GetBool());
                writeBuffer.WriteByte("scenenumber", 6, (byte) _value.GetValue("scenenumber").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedInfoOnOff))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("maskBitInfoOnOffOutput16", (bool) _value.GetValue("maskBitInfoOnOffOutput16").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput15", (bool) _value.GetValue("maskBitInfoOnOffOutput15").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput14", (bool) _value.GetValue("maskBitInfoOnOffOutput14").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput13", (bool) _value.GetValue("maskBitInfoOnOffOutput13").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput12", (bool) _value.GetValue("maskBitInfoOnOffOutput12").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput11", (bool) _value.GetValue("maskBitInfoOnOffOutput11").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput10", (bool) _value.GetValue("maskBitInfoOnOffOutput10").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput9", (bool) _value.GetValue("maskBitInfoOnOffOutput9").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput8", (bool) _value.GetValue("maskBitInfoOnOffOutput8").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput7", (bool) _value.GetValue("maskBitInfoOnOffOutput7").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput6", (bool) _value.GetValue("maskBitInfoOnOffOutput6").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput5", (bool) _value.GetValue("maskBitInfoOnOffOutput5").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput4", (bool) _value.GetValue("maskBitInfoOnOffOutput4").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput3", (bool) _value.GetValue("maskBitInfoOnOffOutput3").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput2", (bool) _value.GetValue("maskBitInfoOnOffOutput2").GetBool());
                writeBuffer.WriteBit("maskBitInfoOnOffOutput1", (bool) _value.GetValue("maskBitInfoOnOffOutput1").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput16", (bool) _value.GetValue("infoOnOffOutput16").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput15", (bool) _value.GetValue("infoOnOffOutput15").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput14", (bool) _value.GetValue("infoOnOffOutput14").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput13", (bool) _value.GetValue("infoOnOffOutput13").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput12", (bool) _value.GetValue("infoOnOffOutput12").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput11", (bool) _value.GetValue("infoOnOffOutput11").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput10", (bool) _value.GetValue("infoOnOffOutput10").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput9", (bool) _value.GetValue("infoOnOffOutput9").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput8", (bool) _value.GetValue("infoOnOffOutput8").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput7", (bool) _value.GetValue("infoOnOffOutput7").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput6", (bool) _value.GetValue("infoOnOffOutput6").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput5", (bool) _value.GetValue("infoOnOffOutput5").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput4", (bool) _value.GetValue("infoOnOffOutput4").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput3", (bool) _value.GetValue("infoOnOffOutput3").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput2", (bool) _value.GetValue("infoOnOffOutput2").GetBool());
                writeBuffer.WriteBit("infoOnOffOutput1", (bool) _value.GetValue("infoOnOffOutput1").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_V64))
            {
                // LINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy_V64))
            {
                // LINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy_V64))
            {
                // LINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_24))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("activationStateOfChannel1", (bool) _value.GetValue("activationStateOfChannel1").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel2", (bool) _value.GetValue("activationStateOfChannel2").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel3", (bool) _value.GetValue("activationStateOfChannel3").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel4", (bool) _value.GetValue("activationStateOfChannel4").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel5", (bool) _value.GetValue("activationStateOfChannel5").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel6", (bool) _value.GetValue("activationStateOfChannel6").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel7", (bool) _value.GetValue("activationStateOfChannel7").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel8", (bool) _value.GetValue("activationStateOfChannel8").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel9", (bool) _value.GetValue("activationStateOfChannel9").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel10", (bool) _value.GetValue("activationStateOfChannel10").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel11", (bool) _value.GetValue("activationStateOfChannel11").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel12", (bool) _value.GetValue("activationStateOfChannel12").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel13", (bool) _value.GetValue("activationStateOfChannel13").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel14", (bool) _value.GetValue("activationStateOfChannel14").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel15", (bool) _value.GetValue("activationStateOfChannel15").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel16", (bool) _value.GetValue("activationStateOfChannel16").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel17", (bool) _value.GetValue("activationStateOfChannel17").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel18", (bool) _value.GetValue("activationStateOfChannel18").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel19", (bool) _value.GetValue("activationStateOfChannel19").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel20", (bool) _value.GetValue("activationStateOfChannel20").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel21", (bool) _value.GetValue("activationStateOfChannel21").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel22", (bool) _value.GetValue("activationStateOfChannel22").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel23", (bool) _value.GetValue("activationStateOfChannel23").GetBool());
                writeBuffer.WriteBit("activationStateOfChannel24", (bool) _value.GetValue("activationStateOfChannel24").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACModeNext))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("delayTimeMin", 16, (ushort) _value.GetValue("delayTimeMin").GetUshort());
                writeBuffer.WriteByte("hvacMode", 8, (byte) _value.GetValue("hvacMode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DHWModeNext))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("delayTimeMin", 16, (ushort) _value.GetValue("delayTimeMin").GetUshort());
                writeBuffer.WriteByte("dhwMode", 8, (byte) _value.GetValue("dhwMode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OccModeNext))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("delayTimeMin", 16, (ushort) _value.GetValue("delayTimeMin").GetUshort());
                writeBuffer.WriteByte("occupancyMode", 8, (byte) _value.GetValue("occupancyMode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BuildingModeNext))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("delayTimeMin", 16, (ushort) _value.GetValue("delayTimeMin").GetUshort());
                writeBuffer.WriteByte("buildingMode", 8, (byte) _value.GetValue("buildingMode").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusLightingActuator))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("validactualvalue", (bool) _value.GetValue("validactualvalue").GetBool());
                writeBuffer.WriteBit("locked", (bool) _value.GetValue("locked").GetBool());
                writeBuffer.WriteBit("forced", (bool) _value.GetValue("forced").GetBool());
                writeBuffer.WriteBit("nightmodeactive", (bool) _value.GetValue("nightmodeactive").GetBool());
                writeBuffer.WriteBit("staircaselightingFunction", (bool) _value.GetValue("staircaselightingFunction").GetBool());
                writeBuffer.WriteBit("dimming", (bool) _value.GetValue("dimming").GetBool());
                writeBuffer.WriteBit("localoverride", (bool) _value.GetValue("localoverride").GetBool());
                writeBuffer.WriteBit("failure", (bool) _value.GetValue("failure").GetBool());
                writeBuffer.WriteByte("actualvalue", 8, (byte) _value.GetValue("actualvalue").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Version))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("magicNumber", 5, (byte) _value.GetValue("magicNumber").GetByte());
                writeBuffer.WriteByte("versionNumber", 5, (byte) _value.GetValue("versionNumber").GetByte());
                writeBuffer.WriteByte("revisionNumber", 6, (byte) _value.GetValue("revisionNumber").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AlarmInfo))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("logNumber", 8, (byte) _value.GetValue("logNumber").GetByte());
                writeBuffer.WriteByte("alarmPriority", 8, (byte) _value.GetValue("alarmPriority").GetByte());
                writeBuffer.WriteByte("applicationArea", 8, (byte) _value.GetValue("applicationArea").GetByte());
                writeBuffer.WriteByte("errorClass", 8, (byte) _value.GetValue("errorClass").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("errorcodeSup", (bool) _value.GetValue("errorcodeSup").GetBool());
                writeBuffer.WriteBit("alarmtextSup", (bool) _value.GetValue("alarmtextSup").GetBool());
                writeBuffer.WriteBit("timestampSup", (bool) _value.GetValue("timestampSup").GetBool());
                writeBuffer.WriteBit("ackSup", (bool) _value.GetValue("ackSup").GetBool());
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("locked", (bool) _value.GetValue("locked").GetBool());
                writeBuffer.WriteBit("alarmunack", (bool) _value.GetValue("alarmunack").GetBool());
                writeBuffer.WriteBit("inalarm", (bool) _value.GetValue("inalarm").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetF16_3))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("tempsetpcomf", 16, (float) _value.GetValue("tempsetpcomf").GetFloat());
                writeBuffer.WriteFloat("tempsetpstdby", 16, (float) _value.GetValue("tempsetpstdby").GetFloat());
                writeBuffer.WriteFloat("tempsetpeco", 16, (float) _value.GetValue("tempsetpeco").GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetShiftF16_3))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("tempsetpshiftcomf", 16, (float) _value.GetValue("tempsetpshiftcomf").GetFloat());
                writeBuffer.WriteFloat("tempsetpshiftstdby", 16, (float) _value.GetValue("tempsetpshiftstdby").GetFloat());
                writeBuffer.WriteFloat("tempsetpshifteco", 16, (float) _value.GetValue("tempsetpshifteco").GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling_Speed))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("timePeriod", 16, (ushort) _value.GetValue("timePeriod").GetUshort());
                writeBuffer.WriteByte("percent", 8, (byte) _value.GetValue("percent").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling_Step_Time))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("timePeriod", 16, (ushort) _value.GetValue("timePeriod").GetUshort());
                writeBuffer.WriteByte("percent", 8, (byte) _value.GetValue("percent").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MeteringValue))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("countval", 32, (int) _value.GetValue("countval").GetInt());
                writeBuffer.WriteByte("valinffield", 8, (byte) _value.GetValue("valinffield").GetByte());
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteBit("alarmunack", (bool) _value.GetValue("alarmunack").GetBool());
                writeBuffer.WriteBit("inalarm", (bool) _value.GetValue("inalarm").GetBool());
                writeBuffer.WriteBit("overridden", (bool) _value.GetValue("overridden").GetBool());
                writeBuffer.WriteBit("fault", (bool) _value.GetValue("fault").GetBool());
                writeBuffer.WriteBit("outofservice", (bool) _value.GetValue("outofservice").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MBus_Address))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("manufactid", 16, (ushort) _value.GetValue("manufactid").GetUshort());
                writeBuffer.WriteUint("identnumber", 32, (uint) _value.GetValue("identnumber").GetUint());
                writeBuffer.WriteByte("version", 8, (byte) _value.GetValue("version").GetByte());
                writeBuffer.WriteByte("medium", 8, (byte) _value.GetValue("medium").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_RGB))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("r", 8, (byte) _value.GetValue("r").GetByte());
                writeBuffer.WriteByte("g", 8, (byte) _value.GetValue("g").GetByte());
                writeBuffer.WriteByte("b", 8, (byte) _value.GetValue("b").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII))
            {
                // STRING
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteString("value", 16, "ASCII", _value.GetString());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Tariff_ActiveEnergy))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteInt("activeelectricalenergy", 32, (int) _value.GetValue("activeelectricalenergy").GetInt());
                writeBuffer.WriteByte("tariff", 8, (byte) _value.GetValue("tariff").GetByte());
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("electricalengergyvalidity", (bool) _value.GetValue("electricalengergyvalidity").GetBool());
                writeBuffer.WriteBit("tariffvalidity", (bool) _value.GetValue("tariffvalidity").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Prioritised_Mode_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("deactivationOfPriority", (bool) _value.GetValue("deactivationOfPriority").GetBool());
                writeBuffer.WriteByte("priorityLevel", 3, (byte) _value.GetValue("priorityLevel").GetByte());
                writeBuffer.WriteByte("modeLevel", 4, (byte) _value.GetValue("modeLevel").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("convertorError", (bool) _value.GetValue("convertorError").GetBool());
                writeBuffer.WriteBit("ballastFailure", (bool) _value.GetValue("ballastFailure").GetBool());
                writeBuffer.WriteBit("lampFailure", (bool) _value.GetValue("lampFailure").GetBool());
                writeBuffer.WriteBit("readOrResponse", (bool) _value.GetValue("readOrResponse").GetBool());
                writeBuffer.WriteBit("addressIndicator", (bool) _value.GetValue("addressIndicator").GetBool());
                writeBuffer.WriteByte("daliDeviceAddressOrDaliGroupAddress", 6, (byte) _value.GetValue("daliDeviceAddressOrDaliGroupAddress").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Diagnostics))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteBit("ballastFailure", (bool) _value.GetValue("ballastFailure").GetBool());
                writeBuffer.WriteBit("lampFailure", (bool) _value.GetValue("lampFailure").GetBool());
                writeBuffer.WriteByte("deviceAddress", 6, (byte) _value.GetValue("deviceAddress").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedPosition))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("heightPosition", 8, (byte) _value.GetValue("heightPosition").GetByte());
                writeBuffer.WriteByte("slatsPosition", 8, (byte) _value.GetValue("slatsPosition").GetByte());
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("validitySlatsPosition", (bool) _value.GetValue("validitySlatsPosition").GetBool());
                writeBuffer.WriteBit("validityHeightPosition", (bool) _value.GetValue("validityHeightPosition").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusSAB))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("heightPosition", 8, (byte) _value.GetValue("heightPosition").GetByte());
                writeBuffer.WriteByte("slatsPosition", 8, (byte) _value.GetValue("slatsPosition").GetByte());
                writeBuffer.WriteBit("upperEndPosReached", (bool) _value.GetValue("upperEndPosReached").GetBool());
                writeBuffer.WriteBit("lowerEndPosReached", (bool) _value.GetValue("lowerEndPosReached").GetBool());
                writeBuffer.WriteBit("lowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent", (bool) _value.GetValue("lowerPredefPosReachedTypHeight100PercentSlatsAngle100Percent").GetBool());
                writeBuffer.WriteBit("targetPosDrive", (bool) _value.GetValue("targetPosDrive").GetBool());
                writeBuffer.WriteBit("restrictionOfTargetHeightPosPosCanNotBeReached", (bool) _value.GetValue("restrictionOfTargetHeightPosPosCanNotBeReached").GetBool());
                writeBuffer.WriteBit("restrictionOfSlatsHeightPosPosCanNotBeReached", (bool) _value.GetValue("restrictionOfSlatsHeightPosPosCanNotBeReached").GetBool());
                writeBuffer.WriteBit("atLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm", (bool) _value.GetValue("atLeastOneOfTheInputsWindRainFrostAlarmIsInAlarm").GetBool());
                writeBuffer.WriteBit("upDownPositionIsForcedByMoveupdownforcedInput", (bool) _value.GetValue("upDownPositionIsForcedByMoveupdownforcedInput").GetBool());
                writeBuffer.WriteBit("movementIsLockedEGByDevicelockedInput", (bool) _value.GetValue("movementIsLockedEGByDevicelockedInput").GetBool());
                writeBuffer.WriteBit("actuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface", (bool) _value.GetValue("actuatorSetvalueIsLocallyOverriddenEGViaALocalUserInterface").GetBool());
                writeBuffer.WriteBit("generalFailureOfTheActuatorOrTheDrive", (bool) _value.GetValue("generalFailureOfTheActuatorOrTheDrive").GetBool());
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteBit("validityHeightPos", (bool) _value.GetValue("validityHeightPos").GetBool());
                writeBuffer.WriteBit("validitySlatsPos", (bool) _value.GetValue("validitySlatsPos").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_xyY))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("xAxis", 16, (ushort) _value.GetValue("xAxis").GetUshort());
                writeBuffer.WriteUshort("yAxis", 16, (ushort) _value.GetValue("yAxis").GetUshort());
                writeBuffer.WriteByte("brightness", 8, (byte) _value.GetValue("brightness").GetByte());
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("validityXy", (bool) _value.GetValue("validityXy").GetBool());
                writeBuffer.WriteBit("validityBrightness", (bool) _value.GetValue("validityBrightness").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Status))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("converterModeAccordingToTheDaliConverterStateMachine", 4, (byte) _value.GetValue("converterModeAccordingToTheDaliConverterStateMachine").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteBit("hardwiredSwitchIsActive", (bool) _value.GetValue("hardwiredSwitchIsActive").GetBool());
                writeBuffer.WriteBit("hardwiredInhibitIsActive", (bool) _value.GetValue("hardwiredInhibitIsActive").GetBool());
                writeBuffer.WriteByte("functionTestPending", 2, (byte) _value.GetValue("functionTestPending").GetByte());
                writeBuffer.WriteByte("durationTestPending", 2, (byte) _value.GetValue("durationTestPending").GetByte());
                writeBuffer.WriteByte("partialDurationTestPending", 2, (byte) _value.GetValue("partialDurationTestPending").GetByte());
                writeBuffer.WriteByte("converterFailure", 2, (byte) _value.GetValue("converterFailure").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Test_Result))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("ltrf", 4, (byte) _value.GetValue("ltrf").GetByte());
                writeBuffer.WriteByte("ltrd", 4, (byte) _value.GetValue("ltrd").GetByte());
                writeBuffer.WriteByte("ltrp", 4, (byte) _value.GetValue("ltrp").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteByte("sf", 2, (byte) _value.GetValue("sf").GetByte());
                writeBuffer.WriteByte("sd", 2, (byte) _value.GetValue("sd").GetByte());
                writeBuffer.WriteByte("sp", 2, (byte) _value.GetValue("sp").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteUshort("ldtr", 16, (ushort) _value.GetValue("ldtr").GetUshort());
                writeBuffer.WriteByte("lpdtr", 8, (byte) _value.GetValue("lpdtr").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Battery_Info))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("batteryFailure", (bool) _value.GetValue("batteryFailure").GetBool());
                writeBuffer.WriteBit("batteryDurationFailure", (bool) _value.GetValue("batteryDurationFailure").GetBool());
                writeBuffer.WriteBit("batteryFullyCharged", (bool) _value.GetValue("batteryFullyCharged").GetBool());
                writeBuffer.WriteByte("batteryChargeLevel", 8, (byte) _value.GetValue("batteryChargeLevel").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness_Colour_Temperature_Transition))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteUshort("ms", 16, (ushort) _value.GetValue("ms").GetUshort());
                writeBuffer.WriteUshort("temperatureK", 16, (ushort) _value.GetValue("temperatureK").GetUshort());
                writeBuffer.WriteByte("percent", 8, (byte) _value.GetValue("percent").GetByte());
                writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
                writeBuffer.WriteBit("validityOfTheTimePeriod", (bool) _value.GetValue("validityOfTheTimePeriod").GetBool());
                writeBuffer.WriteBit("validityOfTheAbsoluteColourTemperature", (bool) _value.GetValue("validityOfTheAbsoluteColourTemperature").GetBool());
                writeBuffer.WriteBit("validityOfTheAbsoluteBrightness", (bool) _value.GetValue("validityOfTheAbsoluteBrightness").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness_Colour_Temperature_Control))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cct", (bool) _value.GetValue("cct").GetBool());
                writeBuffer.WriteByte("stepCodeColourTemperature", 3, (byte) _value.GetValue("stepCodeColourTemperature").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cb", (bool) _value.GetValue("cb").GetBool());
                writeBuffer.WriteByte("stepCodeBrightness", 3, (byte) _value.GetValue("stepCodeBrightness").GetByte());
                writeBuffer.WriteByte("reserved", 6, (byte) (0x00));
                writeBuffer.WriteBit("cctAndStepCodeColourValidity", (bool) _value.GetValue("cctAndStepCodeColourValidity").GetBool());
                writeBuffer.WriteBit("cbAndStepCodeBrightnessValidity", (bool) _value.GetValue("cbAndStepCodeBrightnessValidity").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_RGBW))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("colourLevelRed", 8, (byte) _value.GetValue("colourLevelRed").GetByte());
                writeBuffer.WriteByte("colourLevelGreen", 8, (byte) _value.GetValue("colourLevelGreen").GetByte());
                writeBuffer.WriteByte("colourLevelBlue", 8, (byte) _value.GetValue("colourLevelBlue").GetByte());
                writeBuffer.WriteByte("colourLevelWhite", 8, (byte) _value.GetValue("colourLevelWhite").GetByte());
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("mr", (bool) _value.GetValue("mr").GetBool());
                writeBuffer.WriteBit("mg", (bool) _value.GetValue("mg").GetBool());
                writeBuffer.WriteBit("mb", (bool) _value.GetValue("mb").GetBool());
                writeBuffer.WriteBit("mw", (bool) _value.GetValue("mw").GetBool());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Relative_Control_RGBW))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("maskcw", (bool) _value.GetValue("maskcw").GetBool());
                writeBuffer.WriteBit("maskcb", (bool) _value.GetValue("maskcb").GetBool());
                writeBuffer.WriteBit("maskcg", (bool) _value.GetValue("maskcg").GetBool());
                writeBuffer.WriteBit("maskcr", (bool) _value.GetValue("maskcr").GetBool());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cw", (bool) _value.GetValue("cw").GetBool());
                writeBuffer.WriteByte("stepCodeColourWhite", 3, (byte) _value.GetValue("stepCodeColourWhite").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cb", (bool) _value.GetValue("cb").GetBool());
                writeBuffer.WriteByte("stepCodeColourBlue", 3, (byte) _value.GetValue("stepCodeColourBlue").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cg", (bool) _value.GetValue("cg").GetBool());
                writeBuffer.WriteByte("stepCodeColourGreen", 3, (byte) _value.GetValue("stepCodeColourGreen").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cr", (bool) _value.GetValue("cr").GetBool());
                writeBuffer.WriteByte("stepCodeColourRed", 3, (byte) _value.GetValue("stepCodeColourRed").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Relative_Control_RGB))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cb", (bool) _value.GetValue("cb").GetBool());
                writeBuffer.WriteByte("stepCodeColourBlue", 3, (byte) _value.GetValue("stepCodeColourBlue").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cg", (bool) _value.GetValue("cg").GetBool());
                writeBuffer.WriteByte("stepCodeColourGreen", 3, (byte) _value.GetValue("stepCodeColourGreen").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteBit("cr", (bool) _value.GetValue("cr").GetBool());
                writeBuffer.WriteByte("stepCodeColourRed", 3, (byte) _value.GetValue("stepCodeColourRed").GetByte());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_GeographicalLocation))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("longitude", 32, (float) _value.GetValue("longitude").GetFloat());
                writeBuffer.WriteFloat("latitude", 32, (float) _value.GetValue("latitude").GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetF16_4))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("roomTemperatureSetpointComfort", 16, (float) _value.GetValue("roomTemperatureSetpointComfort").GetFloat());
                writeBuffer.WriteFloat("roomTemperatureSetpointStandby", 16, (float) _value.GetValue("roomTemperatureSetpointStandby").GetFloat());
                writeBuffer.WriteFloat("roomTemperatureSetpointEconomy", 16, (float) _value.GetValue("roomTemperatureSetpointEconomy").GetFloat());
                writeBuffer.WriteFloat("roomTemperatureSetpointBuildingProtection", 16, (float) _value.GetValue("roomTemperatureSetpointBuildingProtection").GetFloat());
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetShiftF16_4))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteFloat("roomTemperatureSetpointShiftComfort", 16, (float) _value.GetValue("roomTemperatureSetpointShiftComfort").GetFloat());
                writeBuffer.WriteFloat("roomTemperatureSetpointShiftStandby", 16, (float) _value.GetValue("roomTemperatureSetpointShiftStandby").GetFloat());
                writeBuffer.WriteFloat("roomTemperatureSetpointShiftEconomy", 16, (float) _value.GetValue("roomTemperatureSetpointShiftEconomy").GetFloat());
                writeBuffer.WriteFloat("roomTemperatureSetpointShiftBuildingProtection", 16, (float) _value.GetValue("roomTemperatureSetpointShiftBuildingProtection").GetFloat());
            }
        }

        public static int GetLengthInBytes(IPlcValue _value, KnxDatapointType datapointType) =>
            (GetLengthInBits(_value, datapointType) + 7) / 8;

        public static int GetLengthInBits(IPlcValue _value, KnxDatapointType datapointType)
        {
            var lengthInBits = 0;
            if (Equals(datapointType, KnxDatapointType.BOOL))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.BYTE))
            {
                // BYTE
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.WORD))
            {
                // WORD
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DWORD))
            {
                // DWORD
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.LWORD))
            {
                // LWORD
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.USINT))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.SINT))
            {
                // SINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.UINT))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.INT))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.UDINT))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DINT))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.ULINT))
            {
                // ULINT
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.LINT))
            {
                // LINT
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.REAL))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.LREAL))
            {
                // LREAL
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.CHAR))
            {
                // CHAR
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.WCHAR))
            {
                // WCHAR
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.TIME))
            {
                // TIME
                throw new NotImplementedException("KnxDatapoint 'TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.LTIME))
            {
                // LTIME
                throw new NotImplementedException("KnxDatapoint 'LTIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DATE))
            {
                // DATE
                throw new NotImplementedException("KnxDatapoint 'DATE' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.TIME_OF_DAY))
            {
                // TIME_OF_DAY
                throw new NotImplementedException("KnxDatapoint 'TIME_OF_DAY' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.TOD))
            {
                // TIME_OF_DAY
                throw new NotImplementedException("KnxDatapoint 'TIME_OF_DAY' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DATE_AND_TIME))
            {
                // DATE_AND_TIME
                throw new NotImplementedException("KnxDatapoint 'DATE_AND_TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DT))
            {
                // DATE_AND_TIME
                throw new NotImplementedException("KnxDatapoint 'DATE_AND_TIME' is not a shape the generator emits yet (design.md GAP-8)");
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Switch))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Bool))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Enable))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ramp))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BinaryValue))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Step))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UpDown))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OpenClose))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Start))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_State))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Invert))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DimSendStyle))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_InputSource))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Reset))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ack))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Trigger))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Occupancy))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Window_Door))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LogicalFunction))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scene_AB))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ShutterBlinds_Mode))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DayNight))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Heat_Cool))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Switch_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Bool_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Enable_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Ramp_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BinaryValue_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Step_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Direction1_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Direction2_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Start_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_State_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Invert_Control))
            {
                // Struct
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Control_Dimming))
            {
                // Struct
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Control_Blinds))
            {
                // Struct
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Char_ASCII))
            {
                // STRING
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Char_8859_1))
            {
                // STRING
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Angle))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_U8))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DecimalFactor))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Tariff))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_1_Ucount))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FanStage))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_V8))
            {
                // SINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_1_Count))
            {
                // SINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Status_Mode3))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 3;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_2_Ucount))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodMsec))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriod10Msec))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriod100Msec))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodSec))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodMin))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimePeriodHrs))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PropDataType))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Length_mm))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UElCurrentmA))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Absolute_Colour_Temperature))
            {
                // UINT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_2_Count))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeMsec))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTime10Msec))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTime100Msec))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeSec))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeMin))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaTimeHrs))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Percent_V16))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Rotation_Angle))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Length_m))
            {
                // INT
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Temp))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Tempd))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Tempa))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Lux))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Wsp))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Pres))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Humidity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AirQuality))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AirFlow))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time1))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time2))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volt))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Curr))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PowerDensity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_KelvinPerPercent))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Power))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume_Flow))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Rain_Amount))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Temp_F))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Wsp_kmh))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Absolute_Humidity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Concentration_ygm3))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Coefficient))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TimeOfDay))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 3;
                lengthInBits += 5;
                lengthInBits += 2;
                lengthInBits += 6;
                lengthInBits += 2;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Date))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 3;
                lengthInBits += 5;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 7;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_4_Ucount))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Sec))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Min))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongTimePeriod_Hrs))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_VolumeLiquid_Litre))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_m_3))
            {
                // UDINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_4_Count))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FlowRate_m3h))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_kWh))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy_kVAh))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy_kVARh))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_MWh))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LongDeltaTimeSec))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaVolumeLiquid_Litre))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DeltaVolume_m_3))
            {
                // DINT
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Acceleration))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Acceleration_Angular))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Activation_Energy))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Activity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Mol))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Amplitude))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AngleRad))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_AngleDeg))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Momentum))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Velocity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Area))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Capacitance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Charge_DensitySurface))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Charge_DensityVolume))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Compressibility))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Conductance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electrical_Conductivity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Density))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Charge))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Current))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_CurrentDensity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_DipoleMoment))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Displacement))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_FieldStrength))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Flux))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_FluxDensity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Polarization))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_Potential))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electric_PotentialDifference))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ElectromagneticMoment))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Electromotive_Force))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Energy))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Force))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Frequency))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Angular_Frequency))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_Capacity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_FlowRate))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Heat_Quantity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Impedance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Length))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Light_Quantity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminous_Flux))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Luminous_Intensity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_FieldStrength))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Flux))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_FluxDensity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Moment))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetic_Polarization))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Magnetization))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_MagnetomotiveForce))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Mass))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_MassFlux))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Momentum))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Phase_AngleRad))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Phase_AngleDeg))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Power))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Power_Factor))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Pressure))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Reactance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Resistance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Resistivity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_SelfInductance))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_SolidAngle))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Sound_Intensity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Speed))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Stress))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Surface_Tension))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Common_Temperature))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Absolute_Temperature))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_TemperatureDifference))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Thermal_Capacity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Thermal_Conductivity))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ThermoelectricPower))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Time))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Torque))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Volume_Flux))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Weight))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_Work))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Value_ApparentPower))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_Flux_Meter))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Volume_Flux_ls))
            {
                // REAL
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Access_Data))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 4;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_String_ASCII))
            {
                // STRING
                lengthInBits += 8;
                lengthInBits += 112;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_String_8859_1))
            {
                // STRING
                lengthInBits += 8;
                lengthInBits += 112;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneNumber))
            {
                // USINT
                lengthInBits += 2;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneControl))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DateTime))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 3;
                lengthInBits += 5;
                lengthInBits += 3;
                lengthInBits += 5;
                lengthInBits += 2;
                lengthInBits += 6;
                lengthInBits += 2;
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 7;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SCLOMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BuildingMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OccMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Priority))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightApplicationMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApplicationArea))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AlarmClassType))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PSUMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ErrorClass_System))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ErrorClass_HVAC))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Time_Delay))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Beaufort_Wind_Force_Scale))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SensorSelect))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActuatorConnectType))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Cloud_Cover))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PowerReturnMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FuelType))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BurnerType))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DHWMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadPriority))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACContrMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACEmergMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ChangeoverMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ValveMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DamperMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HeaterMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FanMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MasterSlaveMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRoomSetp))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Metering_DeviceType))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HumDehumMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EnableHCStage))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ADAType))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BackupMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StartSynchronization))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Behaviour_Lock_Unlock))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Behaviour_Bus_Power_Up_Down))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Fade_Time))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BlinkingMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightControlMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SwitchPBModel))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_PBAction))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DimmPBModel))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SwitchOnMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadTypeSet))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LoadTypeDetected))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Test_Control))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Control))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SABExcept_Behaviour))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SABBehaviour_Lock_Unlock))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SSSBMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BlindsControlMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CommMode))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AddInfoTypes))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_ModeSelect))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_FilterSelect))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_1))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_2))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_EVSEMode_3))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusGen))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 3;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Device_Control))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ForceSign))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ForceSignCool))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRHC))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusSDHWC))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_FuelTypeSet))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRCC))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusAHU))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_RTSM))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 3;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LightActuatorErrorInfo))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_ModeInfo))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_RF_FilterInfo))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_8))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusDHWC))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusRHCC))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_HVA))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 7;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedStatus_RTC))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 7;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Media))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 10;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_16))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OnOffAction))
            {
                // USINT
                lengthInBits += 6;
                lengthInBits += 2;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Alarm_Reaction))
            {
                // USINT
                lengthInBits += 6;
                lengthInBits += 2;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_UpDown_Action))
            {
                // USINT
                lengthInBits += 6;
                lengthInBits += 2;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVAC_PB_Action))
            {
                // USINT
                lengthInBits += 6;
                lengthInBits += 2;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DoubleNibble))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 4;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_SceneInfo))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedInfoOnOff))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ActiveEnergy_V64))
            {
                // LINT
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ApparentEnergy_V64))
            {
                // LINT
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_ReactiveEnergy_V64))
            {
                // LINT
                lengthInBits += 8;
                lengthInBits += 64;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Channel_Activation_24))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_HVACModeNext))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DHWModeNext))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_OccModeNext))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_BuildingModeNext))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusLightingActuator))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Version))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 5;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_AlarmInfo))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetF16_3))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetShiftF16_3))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling_Speed))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Scaling_Step_Time))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MeteringValue))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 32;
                lengthInBits += 8;
                lengthInBits += 3;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_MBus_Address))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 32;
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_RGB))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_LanguageCodeAlpha2_ASCII))
            {
                // STRING
                lengthInBits += 8;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Tariff_ActiveEnergy))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 32;
                lengthInBits += 8;
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Prioritised_Mode_Control))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Control_Gear_Diagnostic))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_DALI_Diagnostics))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 6;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_CombinedPosition))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_StatusSAB))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_xyY))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 8;
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Status))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 2;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 2;
                lengthInBits += 2;
                lengthInBits += 2;
                lengthInBits += 2;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Converter_Test_Result))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 2;
                lengthInBits += 2;
                lengthInBits += 2;
                lengthInBits += 2;
                lengthInBits += 16;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Battery_Info))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 8;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness_Colour_Temperature_Transition))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 8;
                lengthInBits += 5;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Brightness_Colour_Temperature_Control))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 6;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Colour_RGBW))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Relative_Control_RGBW))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 1;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_Relative_Control_RGB))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 3;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_GeographicalLocation))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 32;
                lengthInBits += 32;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetF16_4))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 16;
            }
            else if (Equals(datapointType, KnxDatapointType.DPT_TempRoomSetpSetShiftF16_4))
            {
                // Struct
                lengthInBits += 8;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 16;
                lengthInBits += 16;
            }
            return lengthInBits;
        }
    }
}
