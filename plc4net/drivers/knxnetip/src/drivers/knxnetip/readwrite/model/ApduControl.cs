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
    public abstract partial class ApduControl : IMessage
    {
        public abstract byte ControlType { get; }

        public static ApduControl StaticParse(ReadBuffer readBuffer)
        {
            var controlType = readBuffer.ReadByte("controlType", 2);
            if (Equals(controlType, (byte) (0x0)))
            {
                return ApduControlConnect.StaticParse(readBuffer);
            }
            if (Equals(controlType, (byte) (0x1)))
            {
                return ApduControlDisconnect.StaticParse(readBuffer);
            }
            if (Equals(controlType, (byte) (0x2)))
            {
                return ApduControlAck.StaticParse(readBuffer);
            }
            if (Equals(controlType, (byte) (0x3)))
            {
                return ApduControlNack.StaticParse(readBuffer);
            }
            throw new ParseException("No matching subtype found for ApduControl");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("controlType", 2, ControlType);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 2;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
