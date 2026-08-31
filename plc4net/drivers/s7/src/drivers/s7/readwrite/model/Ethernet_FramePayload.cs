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
    public abstract partial class Ethernet_FramePayload : IMessage
    {
        public abstract ushort PacketType { get; }

        public static Ethernet_FramePayload StaticParse(ReadBuffer readBuffer)
        {
            var packetType = readBuffer.ReadUshort("packetType", 16);
            if (Equals(packetType, (ushort) (0x8100)))
            {
                return Ethernet_FramePayload_VirtualLan.StaticParse(readBuffer);
            }
            if (Equals(packetType, (ushort) (0x8892)))
            {
                return Ethernet_FramePayload_PnDcp.StaticParse(readBuffer);
            }
            throw new ParseException("No matching subtype found for Ethernet_FramePayload");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("packetType", 16, PacketType);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
