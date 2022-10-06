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

        public bool ReadBit(String logicalName)
        {
            return _reader.ReadBool();
        }

        public byte ReadByte(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (byte) _reader.ReadInt(bitLength);
        }
        
        public ushort ReadUshort(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (ushort) _reader.ReadInt(bitLength);
        }

        public uint ReadUint(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (uint) _reader.ReadInt(bitLength);
        }

        public ulong ReadUlong(String logicalName, int bitLength)
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

        public sbyte ReadSbyte(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 8)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (sbyte) _reader.ReadInt(bitLength);
        }
        
        public short ReadShort(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 16)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return (short) _reader.ReadInt(bitLength);
        }

        public int ReadInt(String logicalName, int bitLength)
        {
            if ((bitLength < 0) || (bitLength > 32)) 
            {
                throw new ArgumentOutOfRangeException();
            }
            return _reader.ReadInt(bitLength);
        }

        public long ReadLong(String logicalName, int bitLength)
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

        public float ReadFloat(String logicalName, int bitLength)
        {
            if (bitLength == 32)
            {
                return BitConverter.ToSingle(BitConverter.GetBytes(ReadInt(logicalName, 32)), 0);
            }
            // This is the format as described in the KNX spec ... it's not a real half precision floating point.
            if (bitLength == 16)
            {
                bool sign = true;
                sign = _reader.ReadBool();

                var exp = _reader.ReadInt(4);
                var mantissa = _reader.ReadInt(11);
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

        public double ReadDouble(String logicalName, int bitLength)
        {
            if (bitLength == 32)
            {
                return BitConverter.ToDouble(BitConverter.GetBytes(ReadInt(logicalName, 32)), 0);
            }
            if (bitLength == 64)
            {
                return BitConverter.ToDouble(BitConverter.GetBytes(ReadLong(logicalName, 64)), 0);
            }
            throw new NotImplementedException("This encoding is currently not supported");
        }

        public string ReadString(String logicalName, int bitLength, Encoding encoding)
        {
            throw new NotImplementedException("This encoding is currently not supported");
        }

    }
}