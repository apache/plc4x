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
    public partial class CEMIAdditionalInformationBusmonitorInfo : CEMIAdditionalInformation
    {
        public const byte Len = 1;

        public override byte AdditionalInformationType => (byte) (0x03);

        public bool FrameErrorFlag { get; }
        public bool BitErrorFlag { get; }
        public bool ParityErrorFlag { get; }
        public bool UnknownFlag { get; }
        public bool LostFlag { get; }
        public byte SequenceNumber { get; }

        public CEMIAdditionalInformationBusmonitorInfo(bool frameErrorFlag, bool bitErrorFlag, bool parityErrorFlag, bool unknownFlag, bool lostFlag, byte sequenceNumber)
        {
            FrameErrorFlag = frameErrorFlag;
            BitErrorFlag = bitErrorFlag;
            ParityErrorFlag = parityErrorFlag;
            UnknownFlag = unknownFlag;
            LostFlag = lostFlag;
            SequenceNumber = sequenceNumber;
        }

        public static new CEMIAdditionalInformationBusmonitorInfo StaticParse(ReadBuffer readBuffer)
        {
            var len = readBuffer.ReadByte("len", 8);
            if (!Equals(len, (byte) (1)))
                throw new ParseException($"Expected constant {Len} for 'len' but got {len}");
            var frameErrorFlag = readBuffer.ReadBit("frameErrorFlag");
            var bitErrorFlag = readBuffer.ReadBit("bitErrorFlag");
            var parityErrorFlag = readBuffer.ReadBit("parityErrorFlag");
            var unknownFlag = readBuffer.ReadBit("unknownFlag");
            var lostFlag = readBuffer.ReadBit("lostFlag");
            var sequenceNumber = readBuffer.ReadByte("sequenceNumber", 3);
            return new CEMIAdditionalInformationBusmonitorInfo(frameErrorFlag, bitErrorFlag, parityErrorFlag, unknownFlag, lostFlag, sequenceNumber);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("len", 8, Len);
            writeBuffer.WriteBit("frameErrorFlag", FrameErrorFlag);
            writeBuffer.WriteBit("bitErrorFlag", BitErrorFlag);
            writeBuffer.WriteBit("parityErrorFlag", ParityErrorFlag);
            writeBuffer.WriteBit("unknownFlag", UnknownFlag);
            writeBuffer.WriteBit("lostFlag", LostFlag);
            writeBuffer.WriteByte("sequenceNumber", 3, SequenceNumber);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 3;
            return lengthInBits;
        }

    }
}
