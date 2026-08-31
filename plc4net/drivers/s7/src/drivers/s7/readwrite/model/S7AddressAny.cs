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

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public partial class S7AddressAny : S7Address
    {
        public override byte AddressType => (byte) (0x10);

        public TransportSize TransportSize { get; }
        public ushort NumberOfElements { get; }
        public ushort DbNumber { get; }
        public MemoryArea Area { get; }
        public ushort ByteAddress { get; }
        public byte BitAddress { get; }

        public S7AddressAny(TransportSize transportSize, ushort numberOfElements, ushort dbNumber, MemoryArea area, ushort byteAddress, byte bitAddress)
        {
            TransportSize = transportSize;
            NumberOfElements = numberOfElements;
            DbNumber = dbNumber;
            Area = area;
            ByteAddress = byteAddress;
            BitAddress = bitAddress;
        }

        public static new S7AddressAny StaticParse(ReadBuffer readBuffer)
        {
            var transportSize = TransportSizeExtensions.FirstEnumForFieldCode(readBuffer.ReadByte("transportSize", 8));
            var numberOfElements = readBuffer.ReadUshort("numberOfElements", 16);
            var dbNumber = readBuffer.ReadUshort("dbNumber", 16);
            var area = (MemoryArea) readBuffer.ReadByte("area", 8);
            {
                var reserved = readBuffer.ReadByte("reserved", 5);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var byteAddress = readBuffer.ReadUshort("byteAddress", 16);
            var bitAddress = readBuffer.ReadByte("bitAddress", 3);
            return new S7AddressAny(transportSize, numberOfElements, dbNumber, area, byteAddress, bitAddress);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("transportSize", 8, (byte) TransportSize.GetCode());
            writeBuffer.WriteUshort("numberOfElements", 16, NumberOfElements);
            writeBuffer.WriteUshort("dbNumber", 16, DbNumber);
            writeBuffer.WriteByte("area", 8, (byte) Area);
            writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
            writeBuffer.WriteUshort("byteAddress", 16, ByteAddress);
            writeBuffer.WriteByte("bitAddress", 3, BitAddress);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 5;
            lengthInBits += 16;
            lengthInBits += 3;
            return lengthInBits;
        }

    }
}
