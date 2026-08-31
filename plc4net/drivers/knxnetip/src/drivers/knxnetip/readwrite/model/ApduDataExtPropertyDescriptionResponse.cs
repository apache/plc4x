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
    public partial class ApduDataExtPropertyDescriptionResponse : ApduDataExt
    {
        public override byte ExtApciType => (byte) (0x19);

        public byte ObjectIndex { get; }
        public byte PropertyId { get; }
        public byte Index { get; }
        public bool WriteEnabled { get; }
        public KnxPropertyDataType PropertyDataType { get; }
        public ushort MaxNrOfElements { get; }
        public AccessLevel ReadLevel { get; }
        public AccessLevel WriteLevel { get; }

        public ApduDataExtPropertyDescriptionResponse(byte objectIndex, byte propertyId, byte index, bool writeEnabled, KnxPropertyDataType propertyDataType, ushort maxNrOfElements, AccessLevel readLevel, AccessLevel writeLevel)
        {
            ObjectIndex = objectIndex;
            PropertyId = propertyId;
            Index = index;
            WriteEnabled = writeEnabled;
            PropertyDataType = propertyDataType;
            MaxNrOfElements = maxNrOfElements;
            ReadLevel = readLevel;
            WriteLevel = writeLevel;
        }

        public static new ApduDataExtPropertyDescriptionResponse StaticParse(ReadBuffer readBuffer, byte length)
        {
            var objectIndex = readBuffer.ReadByte("objectIndex", 8);
            var propertyId = readBuffer.ReadByte("propertyId", 8);
            var index = readBuffer.ReadByte("index", 8);
            var writeEnabled = readBuffer.ReadBit("writeEnabled");
            {
                var reserved = readBuffer.ReadByte("reserved", 1);
                if (!Equals(reserved, (byte) (0x0))) { /* mspec reserved: value differs from the spec default */ }
            }
            var propertyDataType = (KnxPropertyDataType) readBuffer.ReadByte("propertyDataType", 8);
            {
                var reserved = readBuffer.ReadByte("reserved", 4);
                if (!Equals(reserved, (byte) (0x0))) { /* mspec reserved: value differs from the spec default */ }
            }
            var maxNrOfElements = readBuffer.ReadUshort("maxNrOfElements", 12);
            var readLevel = (AccessLevel) readBuffer.ReadByte("readLevel", 4);
            var writeLevel = (AccessLevel) readBuffer.ReadByte("writeLevel", 4);
            return new ApduDataExtPropertyDescriptionResponse(objectIndex, propertyId, index, writeEnabled, propertyDataType, maxNrOfElements, readLevel, writeLevel);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("objectIndex", 8, ObjectIndex);
            writeBuffer.WriteByte("propertyId", 8, PropertyId);
            writeBuffer.WriteByte("index", 8, Index);
            writeBuffer.WriteBit("writeEnabled", WriteEnabled);
            writeBuffer.WriteByte("reserved", 1, (byte) (0x0));
            writeBuffer.WriteByte("propertyDataType", 8, (byte) PropertyDataType);
            writeBuffer.WriteByte("reserved", 4, (byte) (0x0));
            writeBuffer.WriteUshort("maxNrOfElements", 12, MaxNrOfElements);
            writeBuffer.WriteByte("readLevel", 4, (byte) ReadLevel);
            writeBuffer.WriteByte("writeLevel", 4, (byte) WriteLevel);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 8;
            lengthInBits += 4;
            lengthInBits += 12;
            lengthInBits += 4;
            lengthInBits += 4;
            return lengthInBits;
        }

    }
}
