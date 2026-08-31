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
    public abstract partial class PnDcp_Pdu : IMessage
    {

        public ushort FrameIdValue { get; }

        protected PnDcp_Pdu(ushort frameIdValue)
        {
            FrameIdValue = frameIdValue;
        }

        public static PnDcp_Pdu StaticParse(ReadBuffer readBuffer)
        {
            var frameIdValue = readBuffer.ReadUshort("frameIdValue", 16);
            if (Equals(frameIdValue, (ushort) (0xFEFE)))
            {
                return PnDcp_Pdu_IdentifyReq.StaticParse(readBuffer, frameIdValue);
            }
            if (Equals(frameIdValue, (ushort) (0xFEFF)))
            {
                return PnDcp_Pdu_IdentifyRes.StaticParse(readBuffer, frameIdValue);
            }
            throw new ParseException("No matching subtype found for PnDcp_Pdu");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("frameIdValue", 16, FrameIdValue);
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
