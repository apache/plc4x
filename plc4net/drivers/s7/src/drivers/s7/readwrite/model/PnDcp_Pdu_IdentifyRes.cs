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
    public partial class PnDcp_Pdu_IdentifyRes : PnDcp_Pdu
    {
        public const byte ServiceId = 0x05;
        public const bool Response = true;

        public bool NotSupported { get; }
        public uint Xid { get; }
        public System.Collections.Generic.List<PnDcp_Block> Blocks { get; }

        public PnDcp_Pdu_IdentifyRes(ushort frameIdValue, bool notSupported, uint xid, System.Collections.Generic.List<PnDcp_Block> blocks) : base(frameIdValue)
        {
            NotSupported = notSupported;
            Xid = xid;
            Blocks = blocks;
        }

        public static PnDcp_Pdu_IdentifyRes StaticParse(ReadBuffer readBuffer, ushort frameIdValue)
        {
            var serviceId = readBuffer.ReadByte("serviceId", 8);
            if (!Equals(serviceId, (byte) (0x05)))
                throw new ParseException($"Expected constant {ServiceId} for 'serviceId' but got {serviceId}");
            {
                var reserved = readBuffer.ReadByte("reserved", 5);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var notSupported = readBuffer.ReadBit("notSupported");
            {
                var reserved = readBuffer.ReadByte("reserved", 1);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var response = readBuffer.ReadBit("response");
            if (!Equals(response, true))
                throw new ParseException($"Expected constant {Response} for 'response' but got {response}");
            var xid = readBuffer.ReadUint("xid", 32);
            {
                var reserved = readBuffer.ReadUshort("reserved", 16);
                if (!Equals(reserved, (ushort) (0x0000))) { /* mspec reserved: value differs from the spec default */ }
            }
            var dcpDataLength = readBuffer.ReadUshort("dcpDataLength", 16);
            var blocks = new System.Collections.Generic.List<PnDcp_Block>();
            var _blocksEnd = readBuffer.GetPos() + (int) (dcpDataLength) * 8;
            while (readBuffer.GetPos() < _blocksEnd)
            {
                blocks.Add(PnDcp_Block.StaticParse(readBuffer));
            }
            return new PnDcp_Pdu_IdentifyRes(frameIdValue, notSupported, xid, blocks);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("serviceId", 8, ServiceId);
            writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
            writeBuffer.WriteBit("notSupported", NotSupported);
            writeBuffer.WriteByte("reserved", 1, (byte) (0x00));
            writeBuffer.WriteBit("response", Response);
            writeBuffer.WriteUint("xid", 32, Xid);
            writeBuffer.WriteUshort("reserved", 16, (ushort) (0x0000));
            writeBuffer.WriteUshort("dcpDataLength", 16, (ushort) ((GetLengthInBytes() - 12)));
            foreach (var _e in Blocks)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 5;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 32;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += Blocks.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
