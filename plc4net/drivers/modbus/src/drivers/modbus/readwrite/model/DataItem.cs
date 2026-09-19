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

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    /// <summary>
    /// mspec <c>[dataIo DataItem]</c> - reads and writes one
    /// <see cref="IPlcValue"/> whose wire layout the parser arguments pick.
    /// </summary>
    public static class DataItem
    {
        public static IPlcValue StaticParse(ReadBuffer readBuffer, ModbusDataType dataType, ushort stringLength)
        {
            if (Equals(dataType, ModbusDataType.BOOL))
            {
                // BOOL
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(dataType, ModbusDataType.BYTE))
            {
                // BYTE
                var value = readBuffer.ReadByte("value", 8);
                return new PlcBYTE(value);
            }
            else if (Equals(dataType, ModbusDataType.WORD))
            {
                // WORD
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcWORD(value);
            }
            else if (Equals(dataType, ModbusDataType.DWORD))
            {
                // DWORD
                var value = readBuffer.ReadUint("value", 32);
                return new PlcDWORD(value);
            }
            else if (Equals(dataType, ModbusDataType.LWORD))
            {
                // LWORD
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcLWORD(value);
            }
            else if (Equals(dataType, ModbusDataType.SINT))
            {
                // SINT
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(dataType, ModbusDataType.INT))
            {
                // INT
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(dataType, ModbusDataType.DINT))
            {
                // DINT
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(dataType, ModbusDataType.LINT))
            {
                // LINT
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(dataType, ModbusDataType.USINT))
            {
                // USINT
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(dataType, ModbusDataType.UINT))
            {
                // UINT
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(dataType, ModbusDataType.UDINT))
            {
                // UDINT
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(dataType, ModbusDataType.ULINT))
            {
                // ULINT
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcULINT(value);
            }
            else if (Equals(dataType, ModbusDataType.REAL))
            {
                // REAL
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(dataType, ModbusDataType.LREAL))
            {
                // LREAL
                var value = readBuffer.ReadDouble("value", 64);
                return new PlcLREAL(value);
            }
            else if (Equals(dataType, ModbusDataType.CHAR))
            {
                // CHAR
                var value = readBuffer.ReadString("value", 8, System.Text.Encoding.UTF8);
                return new PlcCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(dataType, ModbusDataType.WCHAR))
            {
                // WCHAR
                var value = readBuffer.ReadString("value", 16, System.Text.Encoding.BigEndianUnicode);
                return new PlcWCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(dataType, ModbusDataType.STRING))
            {
                // STRING
                var value = readBuffer.ReadString("value", (int) ((stringLength * 8)), System.Text.Encoding.UTF8);
                return new PlcSTRING(value);
            }
            else if (Equals(dataType, ModbusDataType.WSTRING))
            {
                // WSTRING
                var value = readBuffer.ReadString("value", (int) ((stringLength * 16)), System.Text.Encoding.BigEndianUnicode);
                return new PlcWSTRING(value);
            }
            return new PlcNULL();
        }

        public static void StaticSerialize(WriteBuffer writeBuffer, IPlcValue _value, ModbusDataType dataType, ushort stringLength)
        {
            if (Equals(dataType, ModbusDataType.BOOL))
            {
                // BOOL
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(dataType, ModbusDataType.BYTE))
            {
                // BYTE
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(dataType, ModbusDataType.WORD))
            {
                // WORD
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(dataType, ModbusDataType.DWORD))
            {
                // DWORD
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(dataType, ModbusDataType.LWORD))
            {
                // LWORD
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(dataType, ModbusDataType.SINT))
            {
                // SINT
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(dataType, ModbusDataType.INT))
            {
                // INT
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(dataType, ModbusDataType.DINT))
            {
                // DINT
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(dataType, ModbusDataType.LINT))
            {
                // LINT
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(dataType, ModbusDataType.USINT))
            {
                // USINT
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(dataType, ModbusDataType.UINT))
            {
                // UINT
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(dataType, ModbusDataType.UDINT))
            {
                // UDINT
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(dataType, ModbusDataType.ULINT))
            {
                // ULINT
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(dataType, ModbusDataType.REAL))
            {
                // REAL
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(dataType, ModbusDataType.LREAL))
            {
                // LREAL
                writeBuffer.WriteDouble("value", 64, (double) _value.GetDouble());
            }
            else if (Equals(dataType, ModbusDataType.CHAR))
            {
                // CHAR
                writeBuffer.WriteString("value", 8, "UTF8", _value.GetString());
            }
            else if (Equals(dataType, ModbusDataType.WCHAR))
            {
                // WCHAR
                writeBuffer.WriteString("value", 16, "UTF16BE", _value.GetString());
            }
            else if (Equals(dataType, ModbusDataType.STRING))
            {
                // STRING
                writeBuffer.WriteString("value", (int) ((stringLength * 8)), "UTF8", _value.GetString());
            }
            else if (Equals(dataType, ModbusDataType.WSTRING))
            {
                // WSTRING
                writeBuffer.WriteString("value", (int) ((stringLength * 16)), "UTF16BE", _value.GetString());
            }
        }

        public static int GetLengthInBytes(IPlcValue _value, ModbusDataType dataType, ushort stringLength) =>
            (GetLengthInBits(_value, dataType, stringLength) + 7) / 8;

        public static int GetLengthInBits(IPlcValue _value, ModbusDataType dataType, ushort stringLength)
        {
            var lengthInBits = 0;
            if (Equals(dataType, ModbusDataType.BOOL))
            {
                // BOOL
                lengthInBits += 1;
            }
            else if (Equals(dataType, ModbusDataType.BYTE))
            {
                // BYTE
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.WORD))
            {
                // WORD
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.DWORD))
            {
                // DWORD
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.LWORD))
            {
                // LWORD
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.SINT))
            {
                // SINT
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.INT))
            {
                // INT
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.DINT))
            {
                // DINT
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.LINT))
            {
                // LINT
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.USINT))
            {
                // USINT
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.UINT))
            {
                // UINT
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.UDINT))
            {
                // UDINT
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.ULINT))
            {
                // ULINT
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.REAL))
            {
                // REAL
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.LREAL))
            {
                // LREAL
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.CHAR))
            {
                // CHAR
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.WCHAR))
            {
                // WCHAR
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.STRING))
            {
                // STRING
                lengthInBits += ((stringLength * 8));
            }
            else if (Equals(dataType, ModbusDataType.WSTRING))
            {
                // WSTRING
                lengthInBits += ((stringLength * 16));
            }
            return lengthInBits;
        }
    }
}
