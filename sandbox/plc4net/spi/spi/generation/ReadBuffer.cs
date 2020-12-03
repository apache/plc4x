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
using System.Text;
using Ayx.BitIO;

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
            return bitLength < _reader.Remain;
        }

        public byte PeekByte(int offset)
        {
            return _data[(_reader.Position / 8) + offset];
        }

        public bool ReadBit()
        {
            return _reader.ReadBool();
        }

        public byte ReadByte(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (byte) _reader.ReadInt(bitLength);
        }
        
        public ushort ReadUshort(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (ushort) _reader.ReadInt(bitLength);
        }

        public uint ReadUint(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (uint) _reader.ReadInt(bitLength);
        }

        public ulong ReadUlong(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 64)) 
            {
                throw new ArgumentOutOfRangeException();
            }

            ulong firstInt = 0;
            if (bitLength > 32)
            {
                firstInt = (ulong) _reader.ReadInt(bitLength - 32) << 32;
            }
            return firstInt | (ulong) _reader.ReadInt(bitLength);
        }

        public sbyte ReadSbyte(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (sbyte) _reader.ReadInt(bitLength);
        }
        
        public short ReadShort(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (short) _reader.ReadInt(bitLength);
        }

        public int ReadInt(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return _reader.ReadInt(bitLength);
        }

        public long ReadLong(int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 64)) 
            {
                throw new ArgumentOutOfRangeException();
            }

            long firstInt = 0;
            if (bitLength > 32)
            {
                firstInt = (long) _reader.ReadInt(bitLength - 32) << 32;
            }
            return firstInt | (long) _reader.ReadInt(bitLength);
        }

        public float ReadFloat(bool signed, int exponentBitLength, int mantissaBitLength)
        {
            if (signed && exponentBitLength == 8 && mantissaBitLength == 23)
            {
                return BitConverter.ToSingle(BitConverter.GetBytes(ReadInt(32)), 0);
            }
            // This is the format as described in the KNX spec ... it's not a real half precision floating point.
            if (signed && exponentBitLength == 4 && mantissaBitLength == 11)
            {
                bool sign = true;
                if (signed) {
                    sign = _reader.ReadBool();
                }

                var exp = _reader.ReadInt(exponentBitLength);
                var mantissa = _reader.ReadInt(mantissaBitLength);
                // In the mantissa notation actually the first bit is omitted, we need to add it back
                var f = 0.01 * mantissa * Math.Pow(2, exp);
                if (sign)
                {
                    return (float) f * -1;
                }
                return (float) f;
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public double ReadDouble(bool signed, int exponentBitLength, int mantissaBitLength)
        {
            if (signed && exponentBitLength == 8 && mantissaBitLength == 23)
            {
                return BitConverter.ToDouble(BitConverter.GetBytes(ReadInt(32)), 0);
            }
            if (signed && exponentBitLength == 11 && mantissaBitLength == 52)
            {
                return BitConverter.ToDouble(BitConverter.GetBytes(ReadLong(64)), 0);
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public string ReadString(int bitLength, Encoding encoding)
        {
            throw new NotImplementedException("This encoding is currently not supported");
        }

    }
}