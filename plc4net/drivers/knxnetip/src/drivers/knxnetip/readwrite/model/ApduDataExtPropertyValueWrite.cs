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
using System.Linq;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{
    public partial class ApduDataExtPropertyValueWrite : ApduDataExt
    {
        public override byte ExtApciType => (byte) (0x17);

        public byte ObjectIndex { get; }
        public byte PropertyId { get; }
        public byte Count { get; }
        public ushort Index { get; }
        public byte[] Data { get; }

        public ApduDataExtPropertyValueWrite(byte objectIndex, byte propertyId, byte count, ushort index, byte[] data)
        {
            ObjectIndex = objectIndex;
            PropertyId = propertyId;
            Count = count;
            Index = index;
            Data = data;
        }

        public static new ApduDataExtPropertyValueWrite StaticParse(ReadBuffer readBuffer, byte length)
        {
            var objectIndex = readBuffer.ReadByte("objectIndex", 8);
            var propertyId = readBuffer.ReadByte("propertyId", 8);
            var count = readBuffer.ReadByte("count", 4);
            var index = readBuffer.ReadUshort("index", 12);
            var data = readBuffer.ReadByteArray("data", (int) ((length - 5)) * 8);
            return new ApduDataExtPropertyValueWrite(objectIndex, propertyId, count, index, data);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("objectIndex", 8, ObjectIndex);
            writeBuffer.WriteByte("propertyId", 8, PropertyId);
            writeBuffer.WriteByte("count", 4, Count);
            writeBuffer.WriteUshort("index", 12, Index);
            writeBuffer.WriteByteArray("data", Data);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 4;
            lengthInBits += 12;
            lengthInBits += (Data.Length * 8);
            return lengthInBits;
        }

    }
}
