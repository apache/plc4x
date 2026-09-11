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
    public partial class S7PayloadDiagnosticMessage : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x04);
        public override byte CpuFunctionType => (byte) (0x00);
        public override byte CpuSubfunction => (byte) (0x03);

        public ushort EventId { get; }
        public byte PriorityClass { get; }
        public byte ObNumber { get; }
        public ushort DatId { get; }
        public ushort Info1 { get; }
        public uint Info2 { get; }
        public DateAndTime TimeStamp { get; }

        public S7PayloadDiagnosticMessage(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, ushort eventId, byte priorityClass, byte obNumber, ushort datId, ushort info1, uint info2, DateAndTime timeStamp) : base(returnCode, transportSize, dataLength)
        {
            EventId = eventId;
            PriorityClass = priorityClass;
            ObNumber = obNumber;
            DatId = datId;
            Info1 = info1;
            Info2 = info2;
            TimeStamp = timeStamp;
        }

        public static S7PayloadDiagnosticMessage StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var eventId = readBuffer.ReadUshort("eventId", 16);
            var priorityClass = readBuffer.ReadByte("priorityClass", 8);
            var obNumber = readBuffer.ReadByte("obNumber", 8);
            var datId = readBuffer.ReadUshort("datId", 16);
            var info1 = readBuffer.ReadUshort("info1", 16);
            var info2 = readBuffer.ReadUint("info2", 32);
            var timeStamp = DateAndTime.StaticParse(readBuffer);
            return new S7PayloadDiagnosticMessage(returnCode, transportSize, dataLength, eventId, priorityClass, obNumber, datId, info1, info2, timeStamp);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("eventId", 16, EventId);
            writeBuffer.WriteByte("priorityClass", 8, PriorityClass);
            writeBuffer.WriteByte("obNumber", 8, ObNumber);
            writeBuffer.WriteUshort("datId", 16, DatId);
            writeBuffer.WriteUshort("info1", 16, Info1);
            writeBuffer.WriteUint("info2", 32, Info2);
            TimeStamp.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 32;
            lengthInBits += TimeStamp.GetLengthInBits();
            return lengthInBits;
        }

    }
}
