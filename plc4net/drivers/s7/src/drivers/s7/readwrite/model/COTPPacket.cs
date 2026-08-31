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
    public abstract partial class COTPPacket : IMessage
    {
        public abstract byte TpduCode { get; }

        public System.Collections.Generic.List<COTPParameter> Parameters { get; }
        public S7Message? Payload { get; }

        protected COTPPacket(System.Collections.Generic.List<COTPParameter> parameters, S7Message? payload)
        {
            Parameters = parameters;
            Payload = payload;
        }

        public static COTPPacket StaticParse(ReadBuffer readBuffer, ushort cotpLen)
        {
            var _startPos = readBuffer.GetPos();
            var headerLength = readBuffer.ReadByte("headerLength", 8);
            var tpduCode = readBuffer.ReadByte("tpduCode", 8);
            if (Equals(tpduCode, (byte) (0xF0)))
            {
                return COTPPacketData.StaticParse(readBuffer, cotpLen, headerLength, _startPos);
            }
            if (Equals(tpduCode, (byte) (0xE0)))
            {
                return COTPPacketConnectionRequest.StaticParse(readBuffer, cotpLen, headerLength, _startPos);
            }
            if (Equals(tpduCode, (byte) (0xD0)))
            {
                return COTPPacketConnectionResponse.StaticParse(readBuffer, cotpLen, headerLength, _startPos);
            }
            if (Equals(tpduCode, (byte) (0x80)))
            {
                return COTPPacketDisconnectRequest.StaticParse(readBuffer, cotpLen, headerLength, _startPos);
            }
            if (Equals(tpduCode, (byte) (0xC0)))
            {
                return COTPPacketDisconnectResponse.StaticParse(readBuffer, cotpLen, headerLength, _startPos);
            }
            if (Equals(tpduCode, (byte) (0x70)))
            {
                return COTPPacketTpduError.StaticParse(readBuffer, cotpLen, headerLength, _startPos);
            }
            throw new ParseException("No matching subtype found for COTPPacket");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("headerLength", 8, (byte) ((GetLengthInBytes() - (((Payload != null) ? Payload.GetLengthInBytes() : 0) + 1))));
            writeBuffer.WriteByte("tpduCode", 8, TpduCode);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
