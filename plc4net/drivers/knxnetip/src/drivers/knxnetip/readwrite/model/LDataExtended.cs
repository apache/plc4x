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
    public partial class LDataExtended : LDataFrame
    {
        public override bool NotAckFrame => true;
        public override bool Polling => false;

        public bool GroupAddress { get; }
        public byte HopCount { get; }
        public byte ExtendedFrameFormat { get; }
        public KnxAddress SourceAddress { get; }
        public byte[] DestinationAddress { get; }
        public Apdu Apdu { get; }

        public LDataExtended(bool frameType, bool notRepeated, CEMIPriority priority, bool acknowledgeRequested, bool errorFlag, bool groupAddress, byte hopCount, byte extendedFrameFormat, KnxAddress sourceAddress, byte[] destinationAddress, Apdu apdu) : base(frameType, notRepeated, priority, acknowledgeRequested, errorFlag)
        {
            GroupAddress = groupAddress;
            HopCount = hopCount;
            ExtendedFrameFormat = extendedFrameFormat;
            SourceAddress = sourceAddress;
            DestinationAddress = destinationAddress;
            Apdu = apdu;
        }

        public static LDataExtended StaticParse(ReadBuffer readBuffer, bool frameType, bool notRepeated, CEMIPriority priority, bool acknowledgeRequested, bool errorFlag)
        {
            var groupAddress = readBuffer.ReadBit("groupAddress");
            var hopCount = readBuffer.ReadByte("hopCount", 3);
            var extendedFrameFormat = readBuffer.ReadByte("extendedFrameFormat", 4);
            var sourceAddress = KnxAddress.StaticParse(readBuffer);
            var destinationAddress = readBuffer.ReadByteArray("destinationAddress", (int) (2) * 8);
            var dataLength = readBuffer.ReadByte("dataLength", 8);
            var apdu = Apdu.StaticParse(readBuffer, (byte) (dataLength));
            return new LDataExtended(frameType, notRepeated, priority, acknowledgeRequested, errorFlag, groupAddress, hopCount, extendedFrameFormat, sourceAddress, destinationAddress, apdu);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteBit("groupAddress", GroupAddress);
            writeBuffer.WriteByte("hopCount", 3, HopCount);
            writeBuffer.WriteByte("extendedFrameFormat", 4, ExtendedFrameFormat);
            SourceAddress.Serialize(writeBuffer);
            writeBuffer.WriteByteArray("destinationAddress", DestinationAddress);
            writeBuffer.WriteByte("dataLength", 8, (byte) ((Apdu.GetLengthInBytes() - 1)));
            Apdu.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 1;
            lengthInBits += 3;
            lengthInBits += 4;
            lengthInBits += SourceAddress.GetLengthInBits();
            lengthInBits += (DestinationAddress.Length * 8);
            lengthInBits += 8;
            lengthInBits += Apdu.GetLengthInBits();
            return lengthInBits;
        }

    }
}
