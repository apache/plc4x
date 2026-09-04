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
    public abstract partial class KnxNetIpMessage : IMessage
    {
        public const byte ProtocolVersion = 0x10;

        public abstract ushort MsgType { get; }

        public static KnxNetIpMessage StaticParse(ReadBuffer readBuffer)
        {
            var headerLength = readBuffer.ReadByte("headerLength", 8);
            var protocolVersion = readBuffer.ReadByte("protocolVersion", 8);
            if (!Equals(protocolVersion, (byte) (0x10)))
                throw new ParseException($"Expected constant {ProtocolVersion} for 'protocolVersion' but got {protocolVersion}");
            var msgType = readBuffer.ReadUshort("msgType", 16);
            var totalLength = readBuffer.ReadUshort("totalLength", 16);
            if (Equals(msgType, (ushort) (0x0201)))
            {
                return SearchRequest.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0202)))
            {
                return SearchResponse.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0203)))
            {
                return DescriptionRequest.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0204)))
            {
                return DescriptionResponse.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0205)))
            {
                return ConnectionRequest.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0206)))
            {
                return ConnectionResponse.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0207)))
            {
                return ConnectionStateRequest.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0208)))
            {
                return ConnectionStateResponse.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0209)))
            {
                return DisconnectRequest.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x020A)))
            {
                return DisconnectResponse.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x020B)))
            {
                return UnknownMessage.StaticParse(readBuffer, totalLength);
            }
            if (Equals(msgType, (ushort) (0x0310)))
            {
                return DeviceConfigurationRequest.StaticParse(readBuffer, totalLength);
            }
            if (Equals(msgType, (ushort) (0x0311)))
            {
                return DeviceConfigurationAck.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0420)))
            {
                return TunnelingRequest.StaticParse(readBuffer, totalLength);
            }
            if (Equals(msgType, (ushort) (0x0421)))
            {
                return TunnelingResponse.StaticParse(readBuffer);
            }
            if (Equals(msgType, (ushort) (0x0530)))
            {
                return RoutingIndication.StaticParse(readBuffer);
            }
            throw new ParseException("No matching subtype found for KnxNetIpMessage");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("headerLength", 8, (byte) (6));
            writeBuffer.WriteByte("protocolVersion", 8, ProtocolVersion);
            writeBuffer.WriteUshort("msgType", 16, MsgType);
            writeBuffer.WriteUshort("totalLength", 16, (ushort) (GetLengthInBytes()));
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
