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
    /// mspec <c>[dataIo KnxProperty]</c> - reads and writes one
    /// <see cref="IPlcValue"/> whose wire layout the parser arguments pick.
    /// </summary>
    public static class KnxProperty
    {
        public static IPlcValue StaticParse(ReadBuffer readBuffer, KnxPropertyDataType propertyType, byte dataLengthInBytes)
        {
            if (Equals(propertyType, KnxPropertyDataType.PDT_CONTROL))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_CHAR))
            {
                // SINT
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_CHAR))
            {
                // USINT
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_INT))
            {
                // INT
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_INT) && Equals(dataLengthInBytes, (byte) (4)))
            {
                // UDINT
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_INT))
            {
                // UINT
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_KNX_FLOAT))
            {
                // REAL
                var value = readBuffer.ReadFloat("value", 16);
                return new PlcREAL(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DATE))
            {
                // Struct
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_TIME))
            {
                // Struct
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_LONG))
            {
                // DINT
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_LONG))
            {
                // UDINT
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_FLOAT))
            {
                // REAL
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DOUBLE))
            {
                // LREAL
                var value = readBuffer.ReadDouble("value", 64);
                return new PlcLREAL(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_CHAR_BLOCK))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (10);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_POLL_GROUP_SETTINGS))
            {
                // Struct
                var groupAddress = readBuffer.ReadByteArray("groupAddress", (int) (2) * 8);
                var disable = readBuffer.ReadBit("disable");
                {
                    var reserved = readBuffer.ReadByte("reserved", 3);
                    if (!Equals(reserved, (byte) (0x0))) { /* mspec reserved: value differs from the spec default */ }
                }
                var pollingSoftNr = readBuffer.ReadByte("pollingSoftNr", 4);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["groupAddress"] = new PlcRawByteArray(groupAddress);
                _map["disable"] = new PlcBOOL((bool) disable);
                _map["pollingSoftNr"] = new PlcUSINT((byte) pollingSoftNr);
                return new PlcStruct(_map);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_SHORT_CHAR_BLOCK))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (5);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DATE_TIME))
            {
                // Struct
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
                var dayOfMonth = readBuffer.ReadByte("dayOfMonth", 5);
                var dayOfWeek = readBuffer.ReadByte("dayOfWeek", 3);
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
                _map["dayOfMonth"] = new PlcUSINT((byte) dayOfMonth);
                _map["dayOfWeek"] = new PlcUSINT((byte) dayOfWeek);
                _map["hour"] = new PlcUSINT((byte) hour);
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_01))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (1);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_02))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (2);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_03))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (3);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_04))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (4);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_05))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (5);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_06))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (6);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_07))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (7);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_08))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (8);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_09))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (9);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_10))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (10);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_11))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (11);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_12))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (12);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_13))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (13);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_14))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (14);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_15))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (15);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_16))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (16);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_17))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (17);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_18))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (18);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_19))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (19);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_20))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (20);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_VERSION))
            {
                // Struct
                var magicNumber = readBuffer.ReadByte("magicNumber", 5);
                var versionNumber = readBuffer.ReadByte("versionNumber", 5);
                var revisionNumber = readBuffer.ReadByte("revisionNumber", 6);
                var _map = new System.Collections.Generic.Dictionary<string, IPlcValue>();
                _map["magicNumber"] = new PlcUSINT((byte) magicNumber);
                _map["versionNumber"] = new PlcUSINT((byte) versionNumber);
                _map["revisionNumber"] = new PlcUSINT((byte) revisionNumber);
                return new PlcStruct(_map);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_ALARM_INFO))
            {
                // Struct
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BINARY_INFORMATION))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BITSET8))
            {
                // BYTE
                var value = readBuffer.ReadByte("value", 8);
                return new PlcBYTE(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BITSET16))
            {
                // WORD
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcWORD(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_ENUM8))
            {
                // USINT
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_SCALING))
            {
                // USINT
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (dataLengthInBytes);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
        }

        public static void StaticSerialize(WriteBuffer writeBuffer, IPlcValue _value, KnxPropertyDataType propertyType, byte dataLengthInBytes)
        {
            if (Equals(propertyType, KnxPropertyDataType.PDT_CONTROL))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_CHAR))
            {
                // SINT
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_CHAR))
            {
                // USINT
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_INT))
            {
                // INT
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_INT) && Equals(dataLengthInBytes, (byte) (4)))
            {
                // UDINT
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_INT))
            {
                // UINT
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_KNX_FLOAT))
            {
                // REAL
                writeBuffer.WriteFloat("value", 16, (float) _value.GetFloat());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DATE))
            {
                // Struct
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteByte("dayOfMonth", 5, (byte) _value.GetValue("dayOfMonth").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteByte("month", 4, (byte) _value.GetValue("month").GetByte());
                writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
                writeBuffer.WriteByte("year", 7, (byte) _value.GetValue("year").GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_TIME))
            {
                // Struct
                writeBuffer.WriteByte("day", 3, (byte) _value.GetValue("day").GetByte());
                writeBuffer.WriteByte("hour", 5, (byte) _value.GetValue("hour").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("minutes", 6, (byte) _value.GetValue("minutes").GetByte());
                writeBuffer.WriteByte("reserved", 2, (byte) (0x00));
                writeBuffer.WriteByte("seconds", 6, (byte) _value.GetValue("seconds").GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_LONG))
            {
                // DINT
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_LONG))
            {
                // UDINT
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_FLOAT))
            {
                // REAL
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DOUBLE))
            {
                // LREAL
                writeBuffer.WriteDouble("value", 64, (double) _value.GetDouble());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_CHAR_BLOCK))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_POLL_GROUP_SETTINGS))
            {
                // Struct
                writeBuffer.WriteByteArray("groupAddress", _value.GetValue("groupAddress").GetRaw());
                writeBuffer.WriteBit("disable", (bool) _value.GetValue("disable").GetBool());
                writeBuffer.WriteByte("reserved", 3, (byte) (0x0));
                writeBuffer.WriteByte("pollingSoftNr", 4, (byte) _value.GetValue("pollingSoftNr").GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_SHORT_CHAR_BLOCK))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DATE_TIME))
            {
                // Struct
                writeBuffer.WriteByte("year", 8, (byte) _value.GetValue("year").GetByte());
                writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
                writeBuffer.WriteByte("month", 4, (byte) _value.GetValue("month").GetByte());
                writeBuffer.WriteByte("reserved", 3, (byte) (0x00));
                writeBuffer.WriteByte("dayOfMonth", 5, (byte) _value.GetValue("dayOfMonth").GetByte());
                writeBuffer.WriteByte("dayOfWeek", 3, (byte) _value.GetValue("dayOfWeek").GetByte());
                writeBuffer.WriteByte("hour", 5, (byte) _value.GetValue("hour").GetByte());
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_01))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_02))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_03))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_04))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_05))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_06))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_07))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_08))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_09))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_10))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_11))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_12))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_13))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_14))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_15))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_16))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_17))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_18))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_19))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_20))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_VERSION))
            {
                // Struct
                writeBuffer.WriteByte("magicNumber", 5, (byte) _value.GetValue("magicNumber").GetByte());
                writeBuffer.WriteByte("versionNumber", 5, (byte) _value.GetValue("versionNumber").GetByte());
                writeBuffer.WriteByte("revisionNumber", 6, (byte) _value.GetValue("revisionNumber").GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_ALARM_INFO))
            {
                // Struct
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BINARY_INFORMATION))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BITSET8))
            {
                // BYTE
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BITSET16))
            {
                // WORD
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_ENUM8))
            {
                // USINT
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_SCALING))
            {
                // USINT
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
        }

        public static int GetLengthInBytes(IPlcValue _value, KnxPropertyDataType propertyType, byte dataLengthInBytes) =>
            (GetLengthInBits(_value, propertyType, dataLengthInBytes) + 7) / 8;

        public static int GetLengthInBits(IPlcValue _value, KnxPropertyDataType propertyType, byte dataLengthInBytes)
        {
            var lengthInBits = 0;
            if (Equals(propertyType, KnxPropertyDataType.PDT_CONTROL))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_CHAR))
            {
                // SINT
                lengthInBits += 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_CHAR))
            {
                // USINT
                lengthInBits += 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_INT))
            {
                // INT
                lengthInBits += 16;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_INT) && Equals(dataLengthInBytes, (byte) (4)))
            {
                // UDINT
                lengthInBits += 32;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_INT))
            {
                // UINT
                lengthInBits += 16;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_KNX_FLOAT))
            {
                // REAL
                lengthInBits += 16;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DATE))
            {
                // Struct
                lengthInBits += 3;
                lengthInBits += 5;
                lengthInBits += 4;
                lengthInBits += 4;
                lengthInBits += 1;
                lengthInBits += 7;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_TIME))
            {
                // Struct
                lengthInBits += 3;
                lengthInBits += 5;
                lengthInBits += 2;
                lengthInBits += 6;
                lengthInBits += 2;
                lengthInBits += 6;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_LONG))
            {
                // DINT
                lengthInBits += 32;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_UNSIGNED_LONG))
            {
                // UDINT
                lengthInBits += 32;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_FLOAT))
            {
                // REAL
                lengthInBits += 32;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DOUBLE))
            {
                // LREAL
                lengthInBits += 64;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_CHAR_BLOCK))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_POLL_GROUP_SETTINGS))
            {
                // Struct
                lengthInBits += ((2) * 8);
                lengthInBits += 1;
                lengthInBits += 3;
                lengthInBits += 4;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_SHORT_CHAR_BLOCK))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_DATE_TIME))
            {
                // Struct
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_01))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_02))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_03))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_04))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_05))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_06))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_07))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_08))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_09))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_10))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_11))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_12))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_13))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_14))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_15))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_16))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_17))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_18))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_19))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_GENERIC_20))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_VERSION))
            {
                // Struct
                lengthInBits += 5;
                lengthInBits += 5;
                lengthInBits += 6;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_ALARM_INFO))
            {
                // Struct
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
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BINARY_INFORMATION))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BITSET8))
            {
                // BYTE
                lengthInBits += 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_BITSET16))
            {
                // WORD
                lengthInBits += 16;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_ENUM8))
            {
                // USINT
                lengthInBits += 8;
            }
            else if (Equals(propertyType, KnxPropertyDataType.PDT_SCALING))
            {
                // USINT
                lengthInBits += 8;
            }
            else
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            return lengthInBits;
        }
    }
}
