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
    public partial class Ethernet_FramePayload_VirtualLan : Ethernet_FramePayload
    {
        public override ushort PacketType => (ushort) (0x8100);

        public VirtualLanPriority Priority { get; }
        public bool Ineligible { get; }
        public ushort Id { get; }
        public Ethernet_FramePayload Payload { get; }

        public Ethernet_FramePayload_VirtualLan(VirtualLanPriority priority, bool ineligible, ushort id, Ethernet_FramePayload payload)
        {
            Priority = priority;
            Ineligible = ineligible;
            Id = id;
            Payload = payload;
        }

        public static new Ethernet_FramePayload_VirtualLan StaticParse(ReadBuffer readBuffer)
        {
            var priority = (VirtualLanPriority) readBuffer.ReadByte("priority", 3);
            var ineligible = readBuffer.ReadBit("ineligible");
            var id = readBuffer.ReadUshort("id", 12);
            var payload = Ethernet_FramePayload.StaticParse(readBuffer);
            return new Ethernet_FramePayload_VirtualLan(priority, ineligible, id, payload);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("priority", 3, (byte) Priority);
            writeBuffer.WriteBit("ineligible", Ineligible);
            writeBuffer.WriteUshort("id", 12, Id);
            Payload.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 3;
            lengthInBits += 1;
            lengthInBits += 12;
            lengthInBits += Payload.GetLengthInBits();
            return lengthInBits;
        }

    }
}
