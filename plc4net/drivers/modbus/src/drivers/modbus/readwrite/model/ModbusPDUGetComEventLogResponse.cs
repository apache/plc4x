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

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    public partial class ModbusPDUGetComEventLogResponse : ModbusPDU
    {
        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x0C);
        public override bool Response => true;

        public ushort Status { get; }
        public ushort EventCount { get; }
        public ushort MessageCount { get; }
        public byte[] Events { get; }

        public ModbusPDUGetComEventLogResponse(ushort status, ushort eventCount, ushort messageCount, byte[] events)
        {
            Status = status;
            EventCount = eventCount;
            MessageCount = messageCount;
            Events = events;
        }

        public static new ModbusPDUGetComEventLogResponse StaticParse(ReadBuffer readBuffer, bool response)
        {
            var byteCount = readBuffer.ReadByte("byteCount", 8);
            var status = readBuffer.ReadUshort("status", 16);
            var eventCount = readBuffer.ReadUshort("eventCount", 16);
            var messageCount = readBuffer.ReadUshort("messageCount", 16);
            var events = readBuffer.ReadByteArray("events", (int) ((byteCount - 6)) * 8);
            return new ModbusPDUGetComEventLogResponse(status, eventCount, messageCount, events);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("byteCount", 8, (byte) ((Events.Length + 6)));
            writeBuffer.WriteUshort("status", 16, Status);
            writeBuffer.WriteUshort("eventCount", 16, EventCount);
            writeBuffer.WriteUshort("messageCount", 16, MessageCount);
            writeBuffer.WriteByteArray("events", Events);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += (Events.Length * 8);
            return lengthInBits;
        }

    }
}
