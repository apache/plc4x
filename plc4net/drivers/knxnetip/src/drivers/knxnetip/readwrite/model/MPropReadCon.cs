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
    public partial class MPropReadCon : CEMI
    {
        public override byte MessageCode => (byte) (0xFB);

        public ushort InterfaceObjectType { get; }
        public byte ObjectInstance { get; }
        public byte PropertyId { get; }
        public byte NumberOfElements { get; }
        public ushort StartIndex { get; }
        public ushort Data { get; }

        public MPropReadCon(ushort interfaceObjectType, byte objectInstance, byte propertyId, byte numberOfElements, ushort startIndex, ushort data)
        {
            InterfaceObjectType = interfaceObjectType;
            ObjectInstance = objectInstance;
            PropertyId = propertyId;
            NumberOfElements = numberOfElements;
            StartIndex = startIndex;
            Data = data;
        }

        public static new MPropReadCon StaticParse(ReadBuffer readBuffer, ushort size)
        {
            var interfaceObjectType = readBuffer.ReadUshort("interfaceObjectType", 16);
            var objectInstance = readBuffer.ReadByte("objectInstance", 8);
            var propertyId = readBuffer.ReadByte("propertyId", 8);
            var numberOfElements = readBuffer.ReadByte("numberOfElements", 4);
            var startIndex = readBuffer.ReadUshort("startIndex", 12);
            var data = readBuffer.ReadUshort("data", 16);
            return new MPropReadCon(interfaceObjectType, objectInstance, propertyId, numberOfElements, startIndex, data);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("interfaceObjectType", 16, InterfaceObjectType);
            writeBuffer.WriteByte("objectInstance", 8, ObjectInstance);
            writeBuffer.WriteByte("propertyId", 8, PropertyId);
            writeBuffer.WriteByte("numberOfElements", 4, NumberOfElements);
            writeBuffer.WriteUshort("startIndex", 12, StartIndex);
            writeBuffer.WriteUshort("data", 16, Data);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 4;
            lengthInBits += 12;
            lengthInBits += 16;
            return lengthInBits;
        }

    }
}
