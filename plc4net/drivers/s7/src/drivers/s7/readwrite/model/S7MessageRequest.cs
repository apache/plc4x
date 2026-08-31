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
    public partial class S7MessageRequest : S7Message
    {
        public override byte MessageType => (byte) (0x01);

        public S7MessageRequest(ushort tpduReference, S7Parameter? parameter, S7Payload? payload) : base(tpduReference, parameter, payload)
        {
        }

        public static S7MessageRequest StaticParse(ReadBuffer readBuffer, byte messageType, ushort tpduReference, ushort parameterLength, ushort payloadLength)
        {
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
            return new S7MessageRequest(tpduReference, parameter, payload);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
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
            lengthInBits += (Parameter?.GetLengthInBits() ?? 0);
            lengthInBits += (Payload?.GetLengthInBits() ?? 0);
            return lengthInBits;
        }

    }
}
