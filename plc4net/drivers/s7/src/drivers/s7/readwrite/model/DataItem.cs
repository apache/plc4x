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

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    /// <summary>
    /// mspec <c>[dataIo DataItem]</c> - reads and writes one
    /// <see cref="IPlcValue"/> whose wire layout the parser arguments pick.
    /// </summary>
    public static class DataItem
    {
        public static IPlcValue StaticParse(ReadBuffer readBuffer, string dataProtocolId, ControllerType controllerType, int stringLength)
        {
            if (Equals(dataProtocolId, "IEC61131_BOOL"))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_BYTE"))
            {
                // BYTE
                var value = readBuffer.ReadByte("value", 8);
                return new PlcBYTE(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_WORD"))
            {
                // WORD
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcWORD(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_DWORD"))
            {
                // DWORD
                var value = readBuffer.ReadUint("value", 32);
                return new PlcDWORD(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_LWORD"))
            {
                // LWORD
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcLWORD(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_SINT"))
            {
                // SINT
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_USINT"))
            {
                // USINT
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_INT"))
            {
                // INT
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_UINT"))
            {
                // UINT
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_DINT"))
            {
                // DINT
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_UDINT"))
            {
                // UDINT
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_LINT"))
            {
                // LINT
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_ULINT"))
            {
                // ULINT
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcULINT(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_REAL"))
            {
                // REAL
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_LREAL"))
            {
                // LREAL
                var value = readBuffer.ReadDouble("value", 64);
                return new PlcLREAL(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_CHAR"))
            {
                // CHAR
                var value = readBuffer.ReadString("value", 8, System.Text.Encoding.UTF8);
                return new PlcCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(dataProtocolId, "IEC61131_WCHAR"))
            {
                // CHAR
                var value = readBuffer.ReadString("value", 16, System.Text.Encoding.BigEndianUnicode);
                return new PlcCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(dataProtocolId, "IEC61131_STRING"))
            {
                // STRING
                var value = (string) (S7StaticHelper.ParseS7String(readBuffer, stringLength, "UTF8"));
                return new PlcSTRING(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_WSTRING"))
            {
                // STRING
                var value = (string) (S7StaticHelper.ParseS7String(readBuffer, stringLength, "UTF16BE"));
                return new PlcSTRING(value);
            }
            else if (Equals(dataProtocolId, "IEC61131_TIME"))
            {
                // TIME
                var milliseconds = readBuffer.ReadUint("milliseconds", 32);
                return PlcTIME.OfMilliseconds(milliseconds);
            }
            else if (Equals(dataProtocolId, "S7_S5TIME"))
            {
                // TIME
                var milliseconds = S7StaticHelper.ParseS5Time(readBuffer);
                return PlcTIME.OfMilliseconds(milliseconds);
            }
            else if (Equals(dataProtocolId, "IEC61131_LTIME"))
            {
                // LTIME
                var nanoseconds = readBuffer.ReadUlong("nanoseconds", 64);
                return PlcLTIME.OfNanoseconds(nanoseconds);
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE"))
            {
                // DATE
                var daysSinceEpoch = S7StaticHelper.ParseTiaDate(readBuffer);
                return PlcDATE.OfDaysSinceEpoch(daysSinceEpoch);
            }
            else if (Equals(dataProtocolId, "IEC61131_TIME_OF_DAY"))
            {
                // TIME_OF_DAY
                var millisecondsSinceMidnight = readBuffer.ReadUint("millisecondsSinceMidnight", 32);
                return PlcTIME_OF_DAY.OfMillisecondsSinceMidnight(millisecondsSinceMidnight);
            }
            else if (Equals(dataProtocolId, "IEC61131_LTIME_OF_DAY"))
            {
                // LTIME_OF_DAY
                var nanosecondsSinceMidnight = readBuffer.ReadUlong("nanosecondsSinceMidnight", 64);
                return PlcLTIME_OF_DAY.OfNanosecondsSinceMidnight(nanosecondsSinceMidnight);
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE_AND_TIME"))
            {
                // DATE_AND_TIME
                var year = S7StaticHelper.ParseSiemensYear(readBuffer);
                var month = S7StaticHelper.BcdToBin(readBuffer.ReadByte("month", 8));
                var day = S7StaticHelper.BcdToBin(readBuffer.ReadByte("day", 8));
                var hour = S7StaticHelper.BcdToBin(readBuffer.ReadByte("hour", 8));
                var minutes = S7StaticHelper.BcdToBin(readBuffer.ReadByte("minutes", 8));
                var seconds = S7StaticHelper.BcdToBin(readBuffer.ReadByte("seconds", 8));
                var millisecondsOfSecond = S7StaticHelper.BcdToBin12(readBuffer.ReadUshort("millisecondsOfSecond", 12));
                readBuffer.ReadByte("dayOfWeek", 4);
                return PlcDATE_AND_TIME.OfSegments(year, month, day, hour, minutes, seconds, millisecondsOfSecond * 1000000);
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE_AND_LTIME"))
            {
                // DATE_AND_LTIME
                var nanosecondsSinceEpoch = readBuffer.ReadUlong("nanosecondsSinceEpoch", 64);
                return PlcDATE_AND_LTIME.OfNanosecondsSinceEpoch(nanosecondsSinceEpoch);
            }
            else if (Equals(dataProtocolId, "IEC61131_DTL"))
            {
                // DATE_AND_LTIME
                var year = readBuffer.ReadUshort("year", 16);
                var month = readBuffer.ReadByte("month", 8);
                var day = readBuffer.ReadByte("day", 8);
                readBuffer.ReadByte("dayOfWeek", 8);
                var hour = readBuffer.ReadByte("hour", 8);
                var minutes = readBuffer.ReadByte("minutes", 8);
                var seconds = readBuffer.ReadByte("seconds", 8);
                var nanosecondsOfSecond = readBuffer.ReadUint("nanosecondsOfSecond", 32);
                return PlcDATE_AND_LTIME.OfSegments(year, month, day, hour, minutes, seconds, nanosecondsOfSecond);
            }
            return new PlcNULL();
        }

        public static void StaticSerialize(WriteBuffer writeBuffer, IPlcValue _value, string dataProtocolId, ControllerType controllerType, int stringLength)
        {
            if (Equals(dataProtocolId, "IEC61131_BOOL"))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(dataProtocolId, "IEC61131_BYTE"))
            {
                // BYTE
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(dataProtocolId, "IEC61131_WORD"))
            {
                // WORD
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(dataProtocolId, "IEC61131_DWORD"))
            {
                // DWORD
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(dataProtocolId, "IEC61131_LWORD"))
            {
                // LWORD
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(dataProtocolId, "IEC61131_SINT"))
            {
                // SINT
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(dataProtocolId, "IEC61131_USINT"))
            {
                // USINT
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(dataProtocolId, "IEC61131_INT"))
            {
                // INT
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(dataProtocolId, "IEC61131_UINT"))
            {
                // UINT
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(dataProtocolId, "IEC61131_DINT"))
            {
                // DINT
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(dataProtocolId, "IEC61131_UDINT"))
            {
                // UDINT
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(dataProtocolId, "IEC61131_LINT"))
            {
                // LINT
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(dataProtocolId, "IEC61131_ULINT"))
            {
                // ULINT
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(dataProtocolId, "IEC61131_REAL"))
            {
                // REAL
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(dataProtocolId, "IEC61131_LREAL"))
            {
                // LREAL
                writeBuffer.WriteDouble("value", 64, (double) _value.GetDouble());
            }
            else if (Equals(dataProtocolId, "IEC61131_CHAR"))
            {
                // CHAR
                writeBuffer.WriteString("value", 8, "UTF8", _value.GetString());
            }
            else if (Equals(dataProtocolId, "IEC61131_WCHAR"))
            {
                // CHAR
                writeBuffer.WriteString("value", 16, "UTF16BE", _value.GetString());
            }
            else if (Equals(dataProtocolId, "IEC61131_STRING"))
            {
                // STRING
                S7StaticHelper.SerializeS7String(writeBuffer, _value, stringLength, "UTF8");
            }
            else if (Equals(dataProtocolId, "IEC61131_WSTRING"))
            {
                // STRING
                S7StaticHelper.SerializeS7String(writeBuffer, _value, stringLength, "UTF16BE");
            }
            else if (Equals(dataProtocolId, "IEC61131_TIME"))
            {
                // TIME
                writeBuffer.WriteUint("milliseconds", 32, (uint) _value.GetDuration().TotalMilliseconds);
            }
            else if (Equals(dataProtocolId, "S7_S5TIME"))
            {
                // TIME
                S7StaticHelper.SerializeS5Time(writeBuffer, _value);
            }
            else if (Equals(dataProtocolId, "IEC61131_LTIME"))
            {
                // LTIME
                writeBuffer.WriteUlong("nanoseconds", 64, ((PlcLTIME) _value).GetNanoseconds());
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE"))
            {
                // DATE
                S7StaticHelper.SerializeTiaDate(writeBuffer, _value);
            }
            else if (Equals(dataProtocolId, "IEC61131_TIME_OF_DAY"))
            {
                // TIME_OF_DAY
                writeBuffer.WriteUint("millisecondsSinceMidnight", 32, (uint) _value.GetTime().ToTimeSpan().TotalMilliseconds);
            }
            else if (Equals(dataProtocolId, "IEC61131_LTIME_OF_DAY"))
            {
                // LTIME_OF_DAY
                writeBuffer.WriteUlong("nanosecondsSinceMidnight", 64, ((PlcLTIME_OF_DAY) _value).GetNanosecondsSinceMidnight());
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE_AND_TIME"))
            {
                // DATE_AND_TIME
                var _dt = _value.GetDateTime();
                S7StaticHelper.SerializeSiemensYear(writeBuffer, _value);
                writeBuffer.WriteByte("month", 8, S7StaticHelper.BinToBcd(_dt.Month));
                writeBuffer.WriteByte("day", 8, S7StaticHelper.BinToBcd(_dt.Day));
                writeBuffer.WriteByte("hour", 8, S7StaticHelper.BinToBcd(_dt.Hour));
                writeBuffer.WriteByte("minutes", 8, S7StaticHelper.BinToBcd(_dt.Minute));
                writeBuffer.WriteByte("seconds", 8, S7StaticHelper.BinToBcd(_dt.Second));
                writeBuffer.WriteUshort("millisecondsOfSecond", 12, S7StaticHelper.BinToBcd12(_dt.Millisecond));
                writeBuffer.WriteByte("dayOfWeek", 4, (byte) ((int) _dt.DayOfWeek + 1));
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE_AND_LTIME"))
            {
                // DATE_AND_LTIME
                writeBuffer.WriteUlong("nanosecondsSinceEpoch", 64, ((PlcDATE_AND_LTIME) _value).GetNanosecondsSinceEpoch());
            }
            else if (Equals(dataProtocolId, "IEC61131_DTL"))
            {
                // DATE_AND_LTIME
                var _dtl = (PlcDATE_AND_LTIME) _value;
                var _dt = _dtl.GetDateTime();
                writeBuffer.WriteUshort("year", 16, (ushort) _dt.Year);
                writeBuffer.WriteByte("month", 8, (byte) _dt.Month);
                writeBuffer.WriteByte("day", 8, (byte) _dt.Day);
                writeBuffer.WriteByte("dayOfWeek", 8, (byte) ((int) _dt.DayOfWeek + 1));
                writeBuffer.WriteByte("hour", 8, (byte) _dt.Hour);
                writeBuffer.WriteByte("minutes", 8, (byte) _dt.Minute);
                writeBuffer.WriteByte("seconds", 8, (byte) _dt.Second);
                writeBuffer.WriteUint("nanosecondsOfSecond", 32, _dtl.GetNanosecondsOfSecond());
            }
        }

        public static int GetLengthInBytes(IPlcValue _value, string dataProtocolId, ControllerType controllerType, int stringLength) =>
            (GetLengthInBits(_value, dataProtocolId, controllerType, stringLength) + 7) / 8;

        public static int GetLengthInBits(IPlcValue _value, string dataProtocolId, ControllerType controllerType, int stringLength)
        {
            var lengthInBits = 0;
            if (Equals(dataProtocolId, "IEC61131_BOOL"))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
            }
            else if (Equals(dataProtocolId, "IEC61131_BYTE"))
            {
                // BYTE
                lengthInBits += 8;
            }
            else if (Equals(dataProtocolId, "IEC61131_WORD"))
            {
                // WORD
                lengthInBits += 16;
            }
            else if (Equals(dataProtocolId, "IEC61131_DWORD"))
            {
                // DWORD
                lengthInBits += 32;
            }
            else if (Equals(dataProtocolId, "IEC61131_LWORD"))
            {
                // LWORD
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_SINT"))
            {
                // SINT
                lengthInBits += 8;
            }
            else if (Equals(dataProtocolId, "IEC61131_USINT"))
            {
                // USINT
                lengthInBits += 8;
            }
            else if (Equals(dataProtocolId, "IEC61131_INT"))
            {
                // INT
                lengthInBits += 16;
            }
            else if (Equals(dataProtocolId, "IEC61131_UINT"))
            {
                // UINT
                lengthInBits += 16;
            }
            else if (Equals(dataProtocolId, "IEC61131_DINT"))
            {
                // DINT
                lengthInBits += 32;
            }
            else if (Equals(dataProtocolId, "IEC61131_UDINT"))
            {
                // UDINT
                lengthInBits += 32;
            }
            else if (Equals(dataProtocolId, "IEC61131_LINT"))
            {
                // LINT
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_ULINT"))
            {
                // ULINT
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_REAL"))
            {
                // REAL
                lengthInBits += 32;
            }
            else if (Equals(dataProtocolId, "IEC61131_LREAL"))
            {
                // LREAL
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_CHAR"))
            {
                // CHAR
                lengthInBits += 8;
            }
            else if (Equals(dataProtocolId, "IEC61131_WCHAR"))
            {
                // CHAR
                lengthInBits += 16;
            }
            else if (Equals(dataProtocolId, "IEC61131_STRING"))
            {
                // STRING
                lengthInBits += (((stringLength * 8) + 16));
            }
            else if (Equals(dataProtocolId, "IEC61131_WSTRING"))
            {
                // STRING
                lengthInBits += (((stringLength * 16) + 32));
            }
            else if (Equals(dataProtocolId, "IEC61131_TIME"))
            {
                // TIME
                lengthInBits += 32;
            }
            else if (Equals(dataProtocolId, "S7_S5TIME"))
            {
                // TIME
                lengthInBits += (2);
            }
            else if (Equals(dataProtocolId, "IEC61131_LTIME"))
            {
                // LTIME
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE"))
            {
                // DATE
                lengthInBits += (16);
            }
            else if (Equals(dataProtocolId, "IEC61131_TIME_OF_DAY"))
            {
                // TIME_OF_DAY
                lengthInBits += 32;
            }
            else if (Equals(dataProtocolId, "IEC61131_LTIME_OF_DAY"))
            {
                // LTIME_OF_DAY
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE_AND_TIME"))
            {
                // DATE_AND_TIME
                lengthInBits += (8);
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 12;
                lengthInBits += 4;
            }
            else if (Equals(dataProtocolId, "IEC61131_DATE_AND_LTIME"))
            {
                // DATE_AND_LTIME
                lengthInBits += 64;
            }
            else if (Equals(dataProtocolId, "IEC61131_DTL"))
            {
                // DATE_AND_LTIME
                lengthInBits += 16;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 8;
                lengthInBits += 32;
            }
            return lengthInBits;
        }
    }
}
