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
    public partial class ModbusPDUGetComEventCounterResponse : ModbusPDU
    {
        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x0B);
        public override bool Response => true;

        public ushort Status { get; }
        public ushort EventCount { get; }

        public ModbusPDUGetComEventCounterResponse(ushort status, ushort eventCount)
        {
            Status = status;
            EventCount = eventCount;
        }

        public static new ModbusPDUGetComEventCounterResponse StaticParse(ReadBuffer readBuffer, bool response)
        {
            var status = readBuffer.ReadUshort("status", 16);
            var eventCount = readBuffer.ReadUshort("eventCount", 16);
            return new ModbusPDUGetComEventCounterResponse(status, eventCount);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("status", 16, Status);
            writeBuffer.WriteUshort("eventCount", 16, EventCount);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 16;
            return lengthInBits;
        }

    }
}
