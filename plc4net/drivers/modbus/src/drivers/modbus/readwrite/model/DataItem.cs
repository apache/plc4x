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
        public static IPlcValue StaticParse(ReadBuffer readBuffer, ModbusDataType dataType, ushort numberOfValues, bool bigEndian)
        {
            if (Equals(dataType, ModbusDataType.BOOL) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadUshort("reserved", 15);
                    if (!Equals(reserved, (ushort) (0x0000))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                return new PlcBOOL(value);
            }
            else if (Equals(dataType, ModbusDataType.BOOL) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // BOOL
                {
                    var reserved = readBuffer.ReadByte("reserved", 7);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadBit("value");
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                return new PlcBOOL(value);
            }
            else if (Equals(dataType, ModbusDataType.BOOL))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcBOOL(readBuffer.ReadBit("value")));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.BYTE) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // BYTE
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcBYTE(value);
            }
            else if (Equals(dataType, ModbusDataType.BYTE) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // BYTE
                var value = readBuffer.ReadByte("value", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                return new PlcBYTE(value);
            }
            else if (Equals(dataType, ModbusDataType.BYTE))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) ((numberOfValues * 8));
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcBOOL(readBuffer.ReadBit("value")));
                }
                return new PlcList(value);
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
            else if (Equals(dataType, ModbusDataType.SINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // SINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadSbyte("value", 8);
                return new PlcSINT(value);
            }
            else if (Equals(dataType, ModbusDataType.SINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // SINT
                var value = readBuffer.ReadSbyte("value", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                return new PlcSINT(value);
            }
            else if (Equals(dataType, ModbusDataType.SINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcSINT(readBuffer.ReadSbyte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.INT) && Equals(numberOfValues, (ushort) (1)))
            {
                // INT
                var value = readBuffer.ReadShort("value", 16);
                return new PlcINT(value);
            }
            else if (Equals(dataType, ModbusDataType.INT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcINT(readBuffer.ReadShort("value", 16)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.DINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // DINT
                var value = readBuffer.ReadInt("value", 32);
                return new PlcDINT(value);
            }
            else if (Equals(dataType, ModbusDataType.DINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcDINT(readBuffer.ReadInt("value", 32)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.LINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // LINT
                var value = readBuffer.ReadLong("value", 64);
                return new PlcLINT(value);
            }
            else if (Equals(dataType, ModbusDataType.LINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcLINT(readBuffer.ReadLong("value", 64)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.USINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // USINT
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                var value = readBuffer.ReadByte("value", 8);
                return new PlcUSINT(value);
            }
            else if (Equals(dataType, ModbusDataType.USINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // USINT
                var value = readBuffer.ReadByte("value", 8);
                {
                    var reserved = readBuffer.ReadByte("reserved", 8);
                    if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
                }
                return new PlcUSINT(value);
            }
            else if (Equals(dataType, ModbusDataType.USINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUSINT(readBuffer.ReadByte("value", 8)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.UINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // UINT
                var value = readBuffer.ReadUshort("value", 16);
                return new PlcUINT(value);
            }
            else if (Equals(dataType, ModbusDataType.UINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUINT(readBuffer.ReadUshort("value", 16)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.UDINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // UDINT
                var value = readBuffer.ReadUint("value", 32);
                return new PlcUDINT(value);
            }
            else if (Equals(dataType, ModbusDataType.UDINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcUDINT(readBuffer.ReadUint("value", 32)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.ULINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // ULINT
                var value = readBuffer.ReadUlong("value", 64);
                return new PlcULINT(value);
            }
            else if (Equals(dataType, ModbusDataType.ULINT))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcULINT(readBuffer.ReadUlong("value", 64)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.REAL) && Equals(numberOfValues, (ushort) (1)))
            {
                // REAL
                var value = readBuffer.ReadFloat("value", 32);
                return new PlcREAL(value);
            }
            else if (Equals(dataType, ModbusDataType.REAL))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcREAL(readBuffer.ReadFloat("value", 32)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.LREAL) && Equals(numberOfValues, (ushort) (1)))
            {
                // LREAL
                var value = readBuffer.ReadDouble("value", 64);
                return new PlcLREAL(value);
            }
            else if (Equals(dataType, ModbusDataType.LREAL))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    value.Add(new PlcLREAL(readBuffer.ReadDouble("value", 64)));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.CHAR) && Equals(numberOfValues, (ushort) (1)))
            {
                // CHAR
                var value = readBuffer.ReadString("value", 8, System.Text.Encoding.UTF8);
                return new PlcCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(dataType, ModbusDataType.CHAR))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    var _e = readBuffer.ReadString("value", 8, System.Text.Encoding.UTF8);
                    value.Add(new PlcCHAR(_e.Length > 0 ? _e[0] : '\0'));
                }
                return new PlcList(value);
            }
            else if (Equals(dataType, ModbusDataType.WCHAR) && Equals(numberOfValues, (ushort) (1)))
            {
                // WCHAR
                var value = readBuffer.ReadString("value", 16, System.Text.Encoding.BigEndianUnicode);
                return new PlcWCHAR(value.Length > 0 ? value[0] : '\0');
            }
            else if (Equals(dataType, ModbusDataType.WCHAR))
            {
                // List
                var value = new System.Collections.Generic.List<IPlcValue>();
                var _valueCnt = (int) (numberOfValues);
                for (var _i = 0; _i < _valueCnt; _i++)
                {
                    var _e = readBuffer.ReadString("value", 16, System.Text.Encoding.BigEndianUnicode);
                    value.Add(new PlcWCHAR(_e.Length > 0 ? _e[0] : '\0'));
                }
                return new PlcList(value);
            }
            return new PlcNULL();
        }

        public static void StaticSerialize(WriteBuffer writeBuffer, IPlcValue _value, ModbusDataType dataType, ushort numberOfValues, bool bigEndian)
        {
            if (Equals(dataType, ModbusDataType.BOOL) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // BOOL
                writeBuffer.WriteUshort("reserved", 15, (ushort) (0x0000));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
            }
            else if (Equals(dataType, ModbusDataType.BOOL) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // BOOL
                writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
                writeBuffer.WriteBit("value", (bool) _value.GetBool());
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            }
            else if (Equals(dataType, ModbusDataType.BOOL))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteBit("value", (bool) _e.GetBool());
                }
            }
            else if (Equals(dataType, ModbusDataType.BYTE) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // BYTE
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(dataType, ModbusDataType.BYTE) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // BYTE
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            }
            else if (Equals(dataType, ModbusDataType.BYTE))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteBit("value", (bool) _e.GetBool());
                }
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
            else if (Equals(dataType, ModbusDataType.SINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // SINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
            }
            else if (Equals(dataType, ModbusDataType.SINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // SINT
                writeBuffer.WriteSbyte("value", 8, (sbyte) _value.GetSbyte());
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            }
            else if (Equals(dataType, ModbusDataType.SINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteSbyte("value", 8, (sbyte) _e.GetSbyte());
                }
            }
            else if (Equals(dataType, ModbusDataType.INT) && Equals(numberOfValues, (ushort) (1)))
            {
                // INT
                writeBuffer.WriteShort("value", 16, (short) _value.GetShort());
            }
            else if (Equals(dataType, ModbusDataType.INT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteShort("value", 16, (short) _e.GetShort());
                }
            }
            else if (Equals(dataType, ModbusDataType.DINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // DINT
                writeBuffer.WriteInt("value", 32, (int) _value.GetInt());
            }
            else if (Equals(dataType, ModbusDataType.DINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteInt("value", 32, (int) _e.GetInt());
                }
            }
            else if (Equals(dataType, ModbusDataType.LINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // LINT
                writeBuffer.WriteLong("value", 64, (long) _value.GetLong());
            }
            else if (Equals(dataType, ModbusDataType.LINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteLong("value", 64, (long) _e.GetLong());
                }
            }
            else if (Equals(dataType, ModbusDataType.USINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // USINT
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
            }
            else if (Equals(dataType, ModbusDataType.USINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // USINT
                writeBuffer.WriteByte("value", 8, (byte) _value.GetByte());
                writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            }
            else if (Equals(dataType, ModbusDataType.USINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteByte("value", 8, (byte) _e.GetByte());
                }
            }
            else if (Equals(dataType, ModbusDataType.UINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // UINT
                writeBuffer.WriteUshort("value", 16, (ushort) _value.GetUshort());
            }
            else if (Equals(dataType, ModbusDataType.UINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteUshort("value", 16, (ushort) _e.GetUshort());
                }
            }
            else if (Equals(dataType, ModbusDataType.UDINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // UDINT
                writeBuffer.WriteUint("value", 32, (uint) _value.GetUint());
            }
            else if (Equals(dataType, ModbusDataType.UDINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteUint("value", 32, (uint) _e.GetUint());
                }
            }
            else if (Equals(dataType, ModbusDataType.ULINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // ULINT
                writeBuffer.WriteUlong("value", 64, (ulong) _value.GetUlong());
            }
            else if (Equals(dataType, ModbusDataType.ULINT))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteUlong("value", 64, (ulong) _e.GetUlong());
                }
            }
            else if (Equals(dataType, ModbusDataType.REAL) && Equals(numberOfValues, (ushort) (1)))
            {
                // REAL
                writeBuffer.WriteFloat("value", 32, (float) _value.GetFloat());
            }
            else if (Equals(dataType, ModbusDataType.REAL))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteFloat("value", 32, (float) _e.GetFloat());
                }
            }
            else if (Equals(dataType, ModbusDataType.LREAL) && Equals(numberOfValues, (ushort) (1)))
            {
                // LREAL
                writeBuffer.WriteDouble("value", 64, (double) _value.GetDouble());
            }
            else if (Equals(dataType, ModbusDataType.LREAL))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteDouble("value", 64, (double) _e.GetDouble());
                }
            }
            else if (Equals(dataType, ModbusDataType.CHAR) && Equals(numberOfValues, (ushort) (1)))
            {
                // CHAR
                writeBuffer.WriteString("value", 8, "UTF8", _value.GetString());
            }
            else if (Equals(dataType, ModbusDataType.CHAR))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteString("value", 8, "UTF8", _e.GetString());
                }
            }
            else if (Equals(dataType, ModbusDataType.WCHAR) && Equals(numberOfValues, (ushort) (1)))
            {
                // WCHAR
                writeBuffer.WriteString("value", 16, "UTF16BE", _value.GetString());
            }
            else if (Equals(dataType, ModbusDataType.WCHAR))
            {
                // List
                foreach (var _e in _value.GetList())
                {
                    writeBuffer.WriteString("value", 16, "UTF16BE", _e.GetString());
                }
            }
        }

        public static int GetLengthInBytes(IPlcValue _value, ModbusDataType dataType, ushort numberOfValues, bool bigEndian) =>
            (GetLengthInBits(_value, dataType, numberOfValues, bigEndian) + 7) / 8;

        public static int GetLengthInBits(IPlcValue _value, ModbusDataType dataType, ushort numberOfValues, bool bigEndian)
        {
            var lengthInBits = 0;
            if (Equals(dataType, ModbusDataType.BOOL) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // BOOL
                lengthInBits += 15;
                lengthInBits += 1;
            }
            else if (Equals(dataType, ModbusDataType.BOOL) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // BOOL
                lengthInBits += 7;
                lengthInBits += 1;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.BOOL))
            {
                // List
                lengthInBits += _value.GetLength() * 1;
            }
            else if (Equals(dataType, ModbusDataType.BYTE) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // BYTE
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.BYTE) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // BYTE
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.BYTE))
            {
                // List
                lengthInBits += _value.GetLength() * 1;
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
            else if (Equals(dataType, ModbusDataType.SINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // SINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.SINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // SINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.SINT))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(dataType, ModbusDataType.INT) && Equals(numberOfValues, (ushort) (1)))
            {
                // INT
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.INT))
            {
                // List
                lengthInBits += _value.GetLength() * 16;
            }
            else if (Equals(dataType, ModbusDataType.DINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // DINT
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.DINT))
            {
                // List
                lengthInBits += _value.GetLength() * 32;
            }
            else if (Equals(dataType, ModbusDataType.LINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // LINT
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.LINT))
            {
                // List
                lengthInBits += _value.GetLength() * 64;
            }
            else if (Equals(dataType, ModbusDataType.USINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, true))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.USINT) && Equals(numberOfValues, (ushort) (1)) && Equals(bigEndian, false))
            {
                // USINT
                lengthInBits += 8;
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.USINT))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(dataType, ModbusDataType.UINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // UINT
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.UINT))
            {
                // List
                lengthInBits += _value.GetLength() * 16;
            }
            else if (Equals(dataType, ModbusDataType.UDINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // UDINT
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.UDINT))
            {
                // List
                lengthInBits += _value.GetLength() * 32;
            }
            else if (Equals(dataType, ModbusDataType.ULINT) && Equals(numberOfValues, (ushort) (1)))
            {
                // ULINT
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.ULINT))
            {
                // List
                lengthInBits += _value.GetLength() * 64;
            }
            else if (Equals(dataType, ModbusDataType.REAL) && Equals(numberOfValues, (ushort) (1)))
            {
                // REAL
                lengthInBits += 32;
            }
            else if (Equals(dataType, ModbusDataType.REAL))
            {
                // List
                lengthInBits += _value.GetLength() * 32;
            }
            else if (Equals(dataType, ModbusDataType.LREAL) && Equals(numberOfValues, (ushort) (1)))
            {
                // LREAL
                lengthInBits += 64;
            }
            else if (Equals(dataType, ModbusDataType.LREAL))
            {
                // List
                lengthInBits += _value.GetLength() * 64;
            }
            else if (Equals(dataType, ModbusDataType.CHAR) && Equals(numberOfValues, (ushort) (1)))
            {
                // CHAR
                lengthInBits += 8;
            }
            else if (Equals(dataType, ModbusDataType.CHAR))
            {
                // List
                lengthInBits += _value.GetLength() * 8;
            }
            else if (Equals(dataType, ModbusDataType.WCHAR) && Equals(numberOfValues, (ushort) (1)))
            {
                // WCHAR
                lengthInBits += 16;
            }
            else if (Equals(dataType, ModbusDataType.WCHAR))
            {
                // List
                lengthInBits += _value.GetLength() * 16;
            }
            return lengthInBits;
        }
    }
}
