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
    public partial class LPollData : LDataFrame
    {
        public override bool NotAckFrame => true;
        public override bool Polling => true;

        public KnxAddress SourceAddress { get; }
        public byte[] TargetAddress { get; }
        public byte NumberExpectedPollData { get; }

        public LPollData(bool frameType, bool notRepeated, CEMIPriority priority, bool acknowledgeRequested, bool errorFlag, KnxAddress sourceAddress, byte[] targetAddress, byte numberExpectedPollData) : base(frameType, notRepeated, priority, acknowledgeRequested, errorFlag)
        {
            SourceAddress = sourceAddress;
            TargetAddress = targetAddress;
            NumberExpectedPollData = numberExpectedPollData;
        }

        public static LPollData StaticParse(ReadBuffer readBuffer, bool frameType, bool notRepeated, CEMIPriority priority, bool acknowledgeRequested, bool errorFlag)
        {
            var sourceAddress = KnxAddress.StaticParse(readBuffer);
            var targetAddress = readBuffer.ReadByteArray("targetAddress", (int) (2) * 8);
            {
                var reserved = readBuffer.ReadByte("reserved", 4);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var numberExpectedPollData = readBuffer.ReadByte("numberExpectedPollData", 6);
            return new LPollData(frameType, notRepeated, priority, acknowledgeRequested, errorFlag, sourceAddress, targetAddress, numberExpectedPollData);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            SourceAddress.Serialize(writeBuffer);
            writeBuffer.WriteByteArray("targetAddress", TargetAddress);
            writeBuffer.WriteByte("reserved", 4, (byte) (0x00));
            writeBuffer.WriteByte("numberExpectedPollData", 6, NumberExpectedPollData);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += SourceAddress.GetLengthInBits();
            lengthInBits += (TargetAddress.Length * 8);
            lengthInBits += 4;
            lengthInBits += 6;
            return lengthInBits;
        }

    }
}
