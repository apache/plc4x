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
    public abstract partial class ApduDataExt : IMessage
    {
        public abstract byte ExtApciType { get; }

        public static ApduDataExt StaticParse(ReadBuffer readBuffer, byte length)
        {
            var extApciType = readBuffer.ReadByte("extApciType", 6);
            if (Equals(extApciType, (byte) (0x00)))
            {
                return ApduDataExtOpenRoutingTableRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x01)))
            {
                return ApduDataExtReadRoutingTableRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x02)))
            {
                return ApduDataExtReadRoutingTableResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x03)))
            {
                return ApduDataExtWriteRoutingTableRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x08)))
            {
                return ApduDataExtReadRouterMemoryRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x09)))
            {
                return ApduDataExtReadRouterMemoryResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x0A)))
            {
                return ApduDataExtWriteRouterMemoryRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x0D)))
            {
                return ApduDataExtReadRouterStatusRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x0E)))
            {
                return ApduDataExtReadRouterStatusResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x0F)))
            {
                return ApduDataExtWriteRouterStatusRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x10)))
            {
                return ApduDataExtMemoryBitWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x11)))
            {
                return ApduDataExtAuthorizeRequest.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x12)))
            {
                return ApduDataExtAuthorizeResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x13)))
            {
                return ApduDataExtKeyWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x14)))
            {
                return ApduDataExtKeyResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x15)))
            {
                return ApduDataExtPropertyValueRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x16)))
            {
                return ApduDataExtPropertyValueResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x17)))
            {
                return ApduDataExtPropertyValueWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x18)))
            {
                return ApduDataExtPropertyDescriptionRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x19)))
            {
                return ApduDataExtPropertyDescriptionResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x1A)))
            {
                return ApduDataExtNetworkParameterRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x1B)))
            {
                return ApduDataExtNetworkParameterResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x1C)))
            {
                return ApduDataExtIndividualAddressSerialNumberRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x1D)))
            {
                return ApduDataExtIndividualAddressSerialNumberResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x1E)))
            {
                return ApduDataExtIndividualAddressSerialNumberWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x20)))
            {
                return ApduDataExtDomainAddressWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x21)))
            {
                return ApduDataExtDomainAddressRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x22)))
            {
                return ApduDataExtDomainAddressResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x23)))
            {
                return ApduDataExtDomainAddressSelectiveRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x24)))
            {
                return ApduDataExtNetworkParameterWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x25)))
            {
                return ApduDataExtLinkRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x26)))
            {
                return ApduDataExtLinkResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x27)))
            {
                return ApduDataExtLinkWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x28)))
            {
                return ApduDataExtGroupPropertyValueRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x29)))
            {
                return ApduDataExtGroupPropertyValueResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x2A)))
            {
                return ApduDataExtGroupPropertyValueWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x2B)))
            {
                return ApduDataExtGroupPropertyValueInfoReport.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x2C)))
            {
                return ApduDataExtDomainAddressSerialNumberRead.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x2D)))
            {
                return ApduDataExtDomainAddressSerialNumberResponse.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x2E)))
            {
                return ApduDataExtDomainAddressSerialNumberWrite.StaticParse(readBuffer, length);
            }
            if (Equals(extApciType, (byte) (0x30)))
            {
                return ApduDataExtFileStreamInfoReport.StaticParse(readBuffer, length);
            }
            throw new ParseException("No matching subtype found for ApduDataExt");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("extApciType", 6, ExtApciType);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 6;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
