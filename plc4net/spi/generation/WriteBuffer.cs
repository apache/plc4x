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
using System.Text;

namespace org.apache.plc4net.spi.generation
{
    public class WriteBuffer
    {
        private readonly BitWriter _writer;

        public WriteBuffer()
        {
            _writer = new BitWriter();
        }

        public int GetPos()
        {
            return _writer.Length;
        }

        public byte[] GetBytes()
        {
            return _writer.GetBytes();
        }

        public int GetTotalBytes()
        {
            return _writer.Length / 8 + (_writer.Length % 8 != 0 ? 1 : 0);
        }

        public void WriteBit(String logicalName, bool value)
        {
            _writer.WriteBool(value);
        }

        public void WriteByte(String logicalName, int bitLength, byte value)
        {
            if ((bitLength < 0) || (bitLength > 8))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(value, bitLength);
        }

        public void WriteUshort(String logicalName, int bitLength, ushort value)
        {
            if ((bitLength < 0) || (bitLength > 16))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(value, bitLength);
        }

        public void WriteUint(String logicalName, int bitLength, uint value)
        {
            if ((bitLength < 0) || (bitLength > 32))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(value, bitLength);
        }

        public void WriteUlong(String logicalName, int bitLength, ulong value)
        {
            if ((bitLength < 0) || (bitLength > 64))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(value, bitLength);
        }

        public void WriteSbyte(String logicalName, int bitLength, sbyte value)
        {
            if ((bitLength < 0) || (bitLength > 8))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(unchecked((ulong) value), bitLength);
        }

        public void WriteShort(String logicalName, int bitLength, short value)
        {
            if ((bitLength < 0) || (bitLength > 16))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(unchecked((ulong) value), bitLength);
        }

        public void WriteInt(String logicalName, int bitLength, int value)
        {
            if ((bitLength < 0) || (bitLength > 32))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(unchecked((ulong) value), bitLength);
        }

        public void WriteLong(String logicalName, int bitLength, long value)
        {
            if ((bitLength < 0) || (bitLength > 64))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            _writer.WriteBits(unchecked((ulong) value), bitLength);
        }

        public void WriteFloat(String logicalName, int bitLength, float value)
        {
            if (bitLength == 32)
            {
                // The previous version called WriteByte(8, bytes[i]) against a
                // (value, bitLength) signature, so it wrote the constant 8 with a
                // bit width taken from the data.
                _writer.WriteBits(unchecked((uint) BitConverter.SingleToInt32Bits(value)), 32);
                return;
            }
            if (bitLength == 16)
            {
                // KNX DPT 9.x fixed point. The previous version had an empty branch
                // here and silently emitted nothing.
                WriteKnxFloat16(value);
                return;
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        /// <summary>
        /// Encodes a float in the KNX 16 bit fixed-point form: sign, 4 bit exponent,
        /// 11 bit two's complement mantissa, where value = 0.01 * mantissa * 2^exp.
        /// </summary>
        private void WriteKnxFloat16(float value)
        {
            var scaled = value / 0.01;
            var exp = 0;
            while (scaled < -2048 || scaled > 2047)
            {
                scaled /= 2;
                exp++;
                if (exp > 15)
                {
                    throw new ArgumentOutOfRangeException(
                        nameof(value), value, "Value cannot be represented as a KNX 16 bit float.");
                }
            }

            var mantissa = (int) Math.Round(scaled);
            var sign = mantissa < 0;
            if (sign)
            {
                mantissa += 2048;
            }

            _writer.WriteBool(sign);
            _writer.WriteBits((ulong) exp, 4);
            _writer.WriteBits((ulong) (mantissa & 0x7FF), 11);
        }

        public void WriteDouble(String logicalName, int bitLength, double value)
        {
            if (bitLength == 32)
            {
                WriteFloat(logicalName, 32, (float) value);
                return;
            }
            if (bitLength == 64)
            {
                _writer.WriteBits(unchecked((ulong) BitConverter.DoubleToInt64Bits(value)), 64);
                return;
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public void WriteString(String logicalName, int bitLength, string encoding, string value)
        {
            if (bitLength < 0 || bitLength % 8 != 0)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitLength), bitLength, "String lengths must be a whole number of bytes.");
            }

            var raw = ResolveEncoding(encoding).GetBytes(value ?? string.Empty);
            var fieldBytes = bitLength / 8;

            // Truncate or NUL-pad to the declared field width.
            for (var i = 0; i < fieldBytes; i++)
            {
                _writer.WriteBits(i < raw.Length ? raw[i] : (byte) 0, 8);
            }
        }

        public void WriteByteArray(String logicalName, byte[] value)
        {
            _writer.WriteBytes(value);
        }

        private static Encoding ResolveEncoding(string encoding)
        {
            // mspec spells these without the separator, e.g. "UTF8" rather than "UTF-8".
            switch ((encoding ?? "UTF8").Replace("-", "").ToUpperInvariant())
            {
                case "":
                case "UTF8":
                    return Encoding.UTF8;
                case "ASCII":
                    return Encoding.ASCII;
                case "UTF16":
                case "UTF16LE":
                    return Encoding.Unicode;
                case "UTF16BE":
                    return Encoding.BigEndianUnicode;
                default:
                    throw new NotImplementedException($"Unsupported string encoding '{encoding}'");
            }
        }

    }
}
