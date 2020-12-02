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

namespace org.apache.plc4net.spi.generation
{
    public class WriteBuffer
    {
        
        public int GetPos()
        {
            return 0;
        }

        public byte[] GetBytes()
        {
            return null;
        }
    
        public int GetTotalBytes() 
        {
            return 0;
        }

        public void WriteBit(bool value)
        {
        }

        public void WriteByte(int bitLength, byte value)
        {
        }
        
        public void WriteUshort(int bitLength, ushort value)
        {
        }

        public void WriteUint(int bitLength, uint value)
        {
        }

        public void WriteUlong(int bitLength, ulong value)
        {
        }

        public void WriteSbyte(int bitLength, sbyte value)
        {
        }
        
        public void WriteShort(int bitLength, short value)
        {
        }

        public void WriteInt(int bitLength, int value)
        {
        }

        public void WriteLong(int bitLength, long value)
        {
        }

        public void WriteFloat(int bitLength, float value)
        {
        }

        public void WriteDouble(int bitLength, double value)
        {
        }

        public void WriteString(int bitLength, string encoding, string value)
        {
        }

    }
}