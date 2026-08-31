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

    public class ReadBuffer
    {
        private readonly byte[] _data;
        private readonly BitReader _reader;

        public ReadBuffer(byte[] data)
        {
            this._data = data;
            _reader = new BitReader(data);
        }

        public void Reset()
        {
            _reader.Position = 0;
        }

        public int GetPos()
        {
            return _reader.Position;
        }

        public byte[] GetBytes()
        {
            return _data;
        }

        public int GetTotalBytes()
        {
            return _data.Length;
        }

        public bool HasMore(int bitLength)
        {
            // A negative bit length is nonsensical and should not silently return true
            // just because the comparison below would hold for any non-negative remaining.
            if (bitLength < 0)
            {
                return false;
            }
            // Inclusive: a buffer with exactly bitLength bits left can still serve
            // a read of that size.
            return bitLength <= _reader.Remaining;
        }

        public byte PeekByte(int offset)
        {
            return _data[(_reader.Position / 8) + offset];
        }

        public bool ReadBit(String logicalName)
        {
            return _reader.ReadBool();
        }

        public byte ReadByte(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            return (byte) _reader.ReadBits(bitLength);
        }

        public ushort ReadUshort(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            return (ushort) _reader.ReadBits(bitLength);
        }

        public uint ReadUint(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            return (uint) _reader.ReadBits(bitLength);
        }

        public ulong ReadUlong(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 64))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            // The previous implementation read (bitLength - 32) bits and then a
            // further bitLength bits, consuming 2*bitLength-32 bits instead of
            // bitLength, so it never round-tripped against WriteBuffer.WriteUlong.
            return _reader.ReadBits(bitLength);
        }

        public sbyte ReadSbyte(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            return (sbyte) _reader.ReadSignedBits(bitLength);
        }

        public short ReadShort(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            return (short) _reader.ReadSignedBits(bitLength);
        }

        public int ReadInt(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            return (int) _reader.ReadSignedBits(bitLength);
        }

        public long ReadLong(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 64))
            {
                throw new ArgumentOutOfRangeException(nameof(bitLength));
            }
            // Same over-read as ReadUlong had.
            return _reader.ReadSignedBits(bitLength);
        }

        public float ReadFloat(String logicalName, int bitLength)
        {
            if (bitLength == 32)
            {
                // Reinterpret the 32 raw bits as IEEE-754. Going via
                // BitConverter.GetBytes would reorder them on a little-endian host.
                return BitConverter.Int32BitsToSingle((int) _reader.ReadBits(32));
            }
            // KNX DPT 9.x: a 16 bit fixed-point form, not IEEE-754 half precision.
            if (bitLength == 16)
            {
                var sign = _reader.ReadBool();
                var exp = (int) _reader.ReadBits(4);
                var mantissa = (int) _reader.ReadBits(11);
                if (sign)
                {
                    // Negative values are stored as a two's complement mantissa.
                    mantissa -= 2048;
                }
                return (float) (0.01 * mantissa * Math.Pow(2, exp));
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public double ReadDouble(String logicalName, int bitLength)
        {
            if (bitLength == 32)
            {
                // Was BitConverter.ToDouble over a 4 byte array, which always threw.
                return ReadFloat(logicalName, 32);
            }
            if (bitLength == 64)
            {
                return BitConverter.Int64BitsToDouble((long) _reader.ReadBits(64));
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public string ReadString(String logicalName, int bitLength, Encoding encoding)
        {
            if (encoding == null)
            {
                throw new ArgumentNullException(nameof(encoding));
            }
            if (bitLength < 0 || bitLength % 8 != 0)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitLength), bitLength, "String lengths must be a whole number of bytes.");
            }

            var raw = _reader.ReadBytes(bitLength / 8);
            // PLC strings are commonly padded to a fixed field width with NULs.
            return encoding.GetString(raw).TrimEnd('\0');
        }

        public byte[] ReadByteArray(String logicalName, int bitLength)
        {
            if (bitLength < 0 || bitLength % 8 != 0)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitLength), bitLength, "Byte arrays must be a whole number of bytes.");
            }
            return _reader.ReadBytes(bitLength / 8);
        }

    }
}
