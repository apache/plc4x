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
    public partial class CycServiceItemAnyType : CycServiceItemType
    {
        public TransportSize TransportSize { get; }
        public ushort Length { get; }
        public ushort DbNumber { get; }
        public MemoryArea MemoryArea { get; }
        public uint Address { get; }

        public CycServiceItemAnyType(byte byteLength, byte syntaxId, TransportSize transportSize, ushort length, ushort dbNumber, MemoryArea memoryArea, uint address) : base(byteLength, syntaxId)
        {
            TransportSize = transportSize;
            Length = length;
            DbNumber = dbNumber;
            MemoryArea = memoryArea;
            Address = address;
        }

        public static CycServiceItemAnyType StaticParse(ReadBuffer readBuffer, byte byteLength, byte syntaxId)
        {
            var transportSize = TransportSizeExtensions.FirstEnumForFieldCode(readBuffer.ReadByte("transportSize", 8));
            var length = readBuffer.ReadUshort("length", 16);
            var dbNumber = readBuffer.ReadUshort("dbNumber", 16);
            var memoryArea = (MemoryArea) readBuffer.ReadByte("memoryArea", 8);
            var address = readBuffer.ReadUint("address", 24);
            return new CycServiceItemAnyType(byteLength, syntaxId, transportSize, length, dbNumber, memoryArea, address);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("transportSize", 8, (byte) TransportSize.GetCode());
            writeBuffer.WriteUshort("length", 16, Length);
            writeBuffer.WriteUshort("dbNumber", 16, DbNumber);
            writeBuffer.WriteByte("memoryArea", 8, (byte) MemoryArea);
            writeBuffer.WriteUint("address", 24, Address);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 24;
            return lengthInBits;
        }

    }
}
