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
    public partial class S7MessageResponseData : S7Message
    {
        public override byte MessageType => (byte) (0x03);

        public byte ErrorClass { get; }
        public byte ErrorCode { get; }

        public S7MessageResponseData(ushort tpduReference, byte errorClass, byte errorCode, S7Parameter? parameter, S7Payload? payload) : base(tpduReference, parameter, payload)
        {
            ErrorClass = errorClass;
            ErrorCode = errorCode;
        }

        public static S7MessageResponseData StaticParse(ReadBuffer readBuffer, byte messageType, ushort tpduReference, ushort parameterLength, ushort payloadLength)
        {
            var errorClass = readBuffer.ReadByte("errorClass", 8);
            var errorCode = readBuffer.ReadByte("errorCode", 8);
            S7Parameter? parameter = null;
            if ((parameterLength > 0))
            {
                parameter = S7Parameter.StaticParse(readBuffer, (byte) (messageType));
            }
            S7Payload? payload = null;
            if ((payloadLength > 0))
            {
                payload = S7Payload.StaticParse(readBuffer, (byte) (messageType), parameter);
            }
            return new S7MessageResponseData(tpduReference, errorClass, errorCode, parameter, payload);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("errorClass", 8, ErrorClass);
            writeBuffer.WriteByte("errorCode", 8, ErrorCode);
            if (Parameter != null)
            {
                Parameter.Serialize(writeBuffer);
            }
            if (Payload != null)
            {
                Payload.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += (Parameter?.GetLengthInBits() ?? 0);
            lengthInBits += (Payload?.GetLengthInBits() ?? 0);
            return lengthInBits;
        }

    }
}
