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
    public abstract partial class PnDcp_Block : IMessage
    {
        public abstract PnDcp_BlockOptions Option { get; }
        public abstract byte Suboption { get; }

        public static PnDcp_Block StaticParse(ReadBuffer readBuffer)
        {
            var option = (PnDcp_BlockOptions) readBuffer.ReadByte("option", 8);
            var suboption = readBuffer.ReadByte("suboption", 8);
            var blockLength = readBuffer.ReadUshort("blockLength", 16);
            if (Equals(option, PnDcp_BlockOptions.IP_OPTION) && Equals(suboption, (byte) (1)))
            {
                return PnDcp_Block_IpMacAddress.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.IP_OPTION) && Equals(suboption, (byte) (2)))
            {
                return PnDcp_Block_IpParameter.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.IP_OPTION) && Equals(suboption, (byte) (3)))
            {
                return PnDcp_Block_FullIpSuite.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (1)))
            {
                return PnDcp_Block_DevicePropertiesDeviceVendor.StaticParse(readBuffer, blockLength);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (2)))
            {
                return PnDcp_Block_DevicePropertiesNameOfStation.StaticParse(readBuffer, blockLength);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (3)))
            {
                return PnDcp_Block_DevicePropertiesDeviceId.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (4)))
            {
                return PnDcp_Block_DevicePropertiesDeviceRole.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (5)))
            {
                return PnDcp_Block_DevicePropertiesDeviceOptions.StaticParse(readBuffer, blockLength);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (6)))
            {
                return PnDcp_Block_DevicePropertiesAliasName.StaticParse(readBuffer, blockLength);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (7)))
            {
                return PnDcp_Block_DevicePropertiesDeviceInstance.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (8)))
            {
                return PnDcp_Block_DevicePropertiesOemDeviceId.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION) && Equals(suboption, (byte) (9)))
            {
                return PnDcp_Block_DevicePropertiesStandardGateway.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (12)))
            {
                return PnDcp_Block_DhcpOptionHostName.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (43)))
            {
                return PnDcp_Block_DhcpOptionVendorSpecificInformation.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (54)))
            {
                return PnDcp_Block_DhcpOptionServerIdentifier.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (55)))
            {
                return PnDcp_Block_DhcpOptionParameterRequestList.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (60)))
            {
                return PnDcp_Block_DhcpOptionClassIdentifier.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (61)))
            {
                return PnDcp_Block_DhcpOptionDhcpClientIdentifier.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (81)))
            {
                return PnDcp_Block_DhcpOptionFullyQualifiedDomainName.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DCP_OPTION) && Equals(suboption, (byte) (97)))
            {
                return PnDcp_Block_DhcpOptionUuidBasedClient.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.CONTROL_OPTION) && Equals(suboption, (byte) (1)))
            {
                return PnDcp_Block_ControlOptionStart.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.CONTROL_OPTION) && Equals(suboption, (byte) (2)))
            {
                return PnDcp_Block_ControlOptionStop.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.CONTROL_OPTION) && Equals(suboption, (byte) (3)))
            {
                return PnDcp_Block_ControlOptionSignal.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.CONTROL_OPTION) && Equals(suboption, (byte) (4)))
            {
                return PnDcp_Block_ControlOptionResponse.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.CONTROL_OPTION) && Equals(suboption, (byte) (5)))
            {
                return PnDcp_Block_ControlOptionFactoryReset.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.CONTROL_OPTION) && Equals(suboption, (byte) (6)))
            {
                return PnDcp_Block_ControlOptionResetToFactory.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.DEVICE_INITIATIVE_OPTION) && Equals(suboption, (byte) (1)))
            {
                return PnDcp_Block_DeviceInitiativeOption.StaticParse(readBuffer);
            }
            if (Equals(option, PnDcp_BlockOptions.ALL_SELECTOR_OPTION) && Equals(suboption, (byte) (0xFF)))
            {
                return PnDcp_Block_ALLSelector.StaticParse(readBuffer);
            }
            throw new ParseException("No matching subtype found for PnDcp_Block");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("option", 8, (byte) Option);
            writeBuffer.WriteByte("suboption", 8, Suboption);
            writeBuffer.WriteUshort("blockLength", 16, (ushort) ((GetLengthInBytes() - 4)));
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 16;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
