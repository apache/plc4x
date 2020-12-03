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
using Ayx.BitIO;

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

        public void WriteBit(bool value)
        {
            _writer.WriteBool(value);
        }

        public void WriteByte(int bitLength, byte value)
        {
            if (bitLength > 8)
            {
                throw new ArgumentOutOfRangeException();
            }
            _writer.WriteByte(value, bitLength);
        }
        
        public void WriteUshort(int bitLength, ushort value)
        {
            if (bitLength > 16)
            {
                throw new ArgumentOutOfRangeException();
            }
            _writer.WriteInt(value, bitLength);
        }

        public void WriteUint(int bitLength, uint value)
        {
            if (bitLength > 32)
            {
                throw new ArgumentOutOfRangeException();
            }
            _writer.WriteInt((int) value, bitLength);
        }

        public void WriteUlong(int bitLength, ulong value)
        {
            if (bitLength > 64)
            {
                throw new ArgumentOutOfRangeException();
            }
            var highInt = (int) ((value & 0xFFFFFFFF00000000) >> 32);
            var lowInt = (int) (value & 0xFFFFFFFF);
            if (bitLength > 32)
            {
                _writer.WriteInt(highInt, bitLength - 32);
                _writer.WriteInt(lowInt, 32);
            }
            else
            {
                _writer.WriteInt(lowInt, bitLength);
            }
        }

        public void WriteSbyte(int bitLength, sbyte value)
        {
            if (bitLength > 8)
            {
                throw new ArgumentOutOfRangeException();
            }
            _writer.WriteInt(value, bitLength);
        }
        
        public void WriteShort(int bitLength, short value)
        {
            if (bitLength > 16)
            {
                throw new ArgumentOutOfRangeException();
            }
            _writer.WriteInt(value, bitLength);
        }

        public void WriteInt(int bitLength, int value)
        {
            if (bitLength > 32)
            {
                throw new ArgumentOutOfRangeException();
            }
            _writer.WriteInt(value, bitLength);
        }

        public void WriteLong(int bitLength, long value)
        {
            if (bitLength > 64)
            {
                throw new ArgumentOutOfRangeException();
            }
            WriteUlong(bitLength, (ulong) value);
        }

        public void WriteFloat(int bitLength, float value)
        {
            if (bitLength == 32)
            {
                var bytes = BitConverter.GetBytes(value);
                for (var i = 0; i < bytes.Length; i++)
                {
                    _writer.WriteByte(8, bytes[i]);
                }
            }
            else if (bitLength == 16)
            {
                
            }
            else
            {
                throw new NotImplementedException("This encoding is currently not supported");
            }
        }

        public void WriteDouble(int bitLength, double value)
        {
            if (bitLength == 64)
            {
                var bytes = BitConverter.GetBytes(value);
                for (var i = 0; i < bytes.Length; i++)
                {
                    _writer.WriteByte(8, bytes[i]);
                }
            }
            else
            {
                throw new NotImplementedException("This encoding is currently not supported");
            }
        }

        public void WriteString(int bitLength, string encoding, string value)
        {
            throw new NotImplementedException("This encoding is currently not supported");
        }

    }
}