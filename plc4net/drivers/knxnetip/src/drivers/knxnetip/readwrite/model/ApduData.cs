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
    public abstract partial class ApduData : IMessage
    {
        public abstract byte ApciType { get; }

        public static ApduData StaticParse(ReadBuffer readBuffer, byte dataLength)
        {
            var apciType = readBuffer.ReadByte("apciType", 4);
            if (Equals(apciType, (byte) (0x0)))
            {
                return ApduDataGroupValueRead.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x1)))
            {
                return ApduDataGroupValueResponse.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x2)))
            {
                return ApduDataGroupValueWrite.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x3)))
            {
                return ApduDataIndividualAddressWrite.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x4)))
            {
                return ApduDataIndividualAddressRead.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x5)))
            {
                return ApduDataIndividualAddressResponse.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x6)))
            {
                return ApduDataAdcRead.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x7)))
            {
                return ApduDataAdcResponse.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x8)))
            {
                return ApduDataMemoryRead.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0x9)))
            {
                return ApduDataMemoryResponse.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0xA)))
            {
                return ApduDataMemoryWrite.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0xB)))
            {
                return ApduDataUserMessage.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0xC)))
            {
                return ApduDataDeviceDescriptorRead.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0xD)))
            {
                return ApduDataDeviceDescriptorResponse.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0xE)))
            {
                return ApduDataRestart.StaticParse(readBuffer, dataLength);
            }
            if (Equals(apciType, (byte) (0xF)))
            {
                return ApduDataOther.StaticParse(readBuffer, dataLength);
            }
            throw new ParseException("No matching subtype found for ApduData");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("apciType", 4, ApciType);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 4;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
