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

using System;
using System.Globalization;

namespace org.apache.plc4net.spi.model.values
{
    /// <summary>
    /// Base for the IEC bit-string values (<c>BYTE</c>, <c>WORD</c>, <c>DWORD</c>,
    /// <c>LWORD</c>). Keeps the pattern as a right-aligned <see cref="ulong"/> and
    /// serves it both as individual bits and as the (unsigned) integer of the
    /// matching width, so a value survives a parse/serialize round trip and a
    /// consumer holding it as <see cref="IPlcValue"/> can still read a number.
    /// </summary>
    public class PlcBitString : PlcSimpleValueAdapter
    {
        private readonly ulong _bits;
        private readonly int _byteWidth;

        public PlcBitString()
        {
        }

        public PlcBitString(byte value)
        {
            _bits = value;
            _byteWidth = 1;
        }

        public PlcBitString(ushort value)
        {
            _bits = value;
            _byteWidth = 2;
        }

        public PlcBitString(uint value)
        {
            _bits = value;
            _byteWidth = 4;
        }

        public PlcBitString(ulong value)
        {
            _bits = value;
            _byteWidth = 8;
        }

        public override bool IsByte()
        {
            return _byteWidth == 1;
        }

        public override byte GetByte()
        {
            return (byte) _bits;
        }

        public override bool IsUshort()
        {
            return _byteWidth <= 2;
        }

        public override ushort GetUshort()
        {
            return (ushort) _bits;
        }

        public override bool IsUint()
        {
            return _byteWidth <= 4;
        }

        public override uint GetUint()
        {
            return (uint) _bits;
        }

        public override bool IsUlong()
        {
            return true;
        }

        public override ulong GetUlong()
        {
            return _bits;
        }

        // ── signed-integer + string views (a BYTE/WORD is also a plain number) ──

        public override bool IsShort()
        {
            return _bits <= (ulong) short.MaxValue;
        }

        public override short GetShort()
        {
            if (!IsShort())
            {
                throw new ArgumentOutOfRangeException(nameof(GetShort),
                    $"bit-string value {_bits} does not fit in a short");
            }
            return (short) _bits;
        }

        public override bool IsInt()
        {
            return _bits <= int.MaxValue;
        }

        public override int GetInt()
        {
            if (!IsInt())
            {
                throw new ArgumentOutOfRangeException(nameof(GetInt),
                    $"bit-string value {_bits} does not fit in an int");
            }
            return (int) _bits;
        }

        public override bool IsLong()
        {
            return _bits <= long.MaxValue;
        }

        public override long GetLong()
        {
            if (!IsLong())
            {
                throw new ArgumentOutOfRangeException(nameof(GetLong),
                    $"bit-string value {_bits} does not fit in a long");
            }
            return (long) _bits;
        }

        public override bool IsString()
        {
            return true;
        }

        public override string GetString()
        {
            return _bits.ToString(CultureInfo.InvariantCulture);
        }

        public override int GetBoolLength()
        {
            return _byteWidth * 8;
        }

        public override bool GetBoolAt(int index)
        {
            if (index < 0 || index >= _byteWidth * 8)
            {
                return false;
            }

            return ((_bits >> index) & 1UL) == 1UL;
        }

        public override bool[] GetBoolArray()
        {
            var array = new bool[_byteWidth * 8];
            for (var i = 0; i < array.Length; i++)
            {
                array[i] = GetBoolAt(i);
            }

            return array;
        }
    }
}
