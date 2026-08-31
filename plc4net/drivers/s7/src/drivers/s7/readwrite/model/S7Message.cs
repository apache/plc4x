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
    public abstract partial class S7Message : IMessage
    {
        public const byte ProtocolId = 0x32;

        public abstract byte MessageType { get; }

        public ushort TpduReference { get; }
        public S7Parameter? Parameter { get; }
        public S7Payload? Payload { get; }

        protected S7Message(ushort tpduReference, S7Parameter? parameter, S7Payload? payload)
        {
            TpduReference = tpduReference;
            Parameter = parameter;
            Payload = payload;
        }

        public static S7Message StaticParse(ReadBuffer readBuffer)
        {
            var protocolId = readBuffer.ReadByte("protocolId", 8);
            if (!Equals(protocolId, (byte) (0x32)))
                throw new ParseException($"Expected constant {ProtocolId} for 'protocolId' but got {protocolId}");
            var messageType = readBuffer.ReadByte("messageType", 8);
            {
                var reserved = readBuffer.ReadUshort("reserved", 16);
                if (!Equals(reserved, (ushort) (0x0000))) { /* mspec reserved: value differs from the spec default */ }
            }
            var tpduReference = readBuffer.ReadUshort("tpduReference", 16);
            var parameterLength = readBuffer.ReadUshort("parameterLength", 16);
            var payloadLength = readBuffer.ReadUshort("payloadLength", 16);
            if (Equals(messageType, (byte) (0x01)))
            {
                return S7MessageRequest.StaticParse(readBuffer, messageType, tpduReference, parameterLength, payloadLength);
            }
            if (Equals(messageType, (byte) (0x02)))
            {
                return S7MessageResponse.StaticParse(readBuffer, messageType, tpduReference, parameterLength, payloadLength);
            }
            if (Equals(messageType, (byte) (0x03)))
            {
                return S7MessageResponseData.StaticParse(readBuffer, messageType, tpduReference, parameterLength, payloadLength);
            }
            if (Equals(messageType, (byte) (0x07)))
            {
                return S7MessageUserData.StaticParse(readBuffer, messageType, tpduReference, parameterLength, payloadLength);
            }
            throw new ParseException("No matching subtype found for S7Message");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("protocolId", 8, ProtocolId);
            writeBuffer.WriteByte("messageType", 8, MessageType);
            writeBuffer.WriteUshort("reserved", 16, (ushort) (0x0000));
            writeBuffer.WriteUshort("tpduReference", 16, TpduReference);
            writeBuffer.WriteUshort("parameterLength", 16, (ushort) (((Parameter != null) ? Parameter.GetLengthInBytes() : 0)));
            writeBuffer.WriteUshort("payloadLength", 16, (ushort) (((Payload != null) ? Payload.GetLengthInBytes() : 0)));
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
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
