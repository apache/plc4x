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
    public partial class DIBDeviceInfo : IMessage
    {
        public byte DescriptionType { get; }
        public KnxMedium KnxMedium { get; }
        public DeviceStatus DeviceStatus { get; }
        public KnxAddress KnxAddress { get; }
        public ProjectInstallationIdentifier ProjectInstallationIdentifier { get; }
        public byte[] KnxNetIpDeviceSerialNumber { get; }
        public IPAddress KnxNetIpDeviceMulticastAddress { get; }
        public MACAddress KnxNetIpDeviceMacAddress { get; }
        public byte[] DeviceFriendlyName { get; }

        public DIBDeviceInfo(byte descriptionType, KnxMedium knxMedium, DeviceStatus deviceStatus, KnxAddress knxAddress, ProjectInstallationIdentifier projectInstallationIdentifier, byte[] knxNetIpDeviceSerialNumber, IPAddress knxNetIpDeviceMulticastAddress, MACAddress knxNetIpDeviceMacAddress, byte[] deviceFriendlyName)
        {
            DescriptionType = descriptionType;
            KnxMedium = knxMedium;
            DeviceStatus = deviceStatus;
            KnxAddress = knxAddress;
            ProjectInstallationIdentifier = projectInstallationIdentifier;
            KnxNetIpDeviceSerialNumber = knxNetIpDeviceSerialNumber;
            KnxNetIpDeviceMulticastAddress = knxNetIpDeviceMulticastAddress;
            KnxNetIpDeviceMacAddress = knxNetIpDeviceMacAddress;
            DeviceFriendlyName = deviceFriendlyName;
        }

        public static DIBDeviceInfo StaticParse(ReadBuffer readBuffer)
        {
            var structureLength = readBuffer.ReadByte("structureLength", 8);
            var descriptionType = readBuffer.ReadByte("descriptionType", 8);
            var knxMedium = (KnxMedium) readBuffer.ReadByte("knxMedium", 8);
            var deviceStatus = DeviceStatus.StaticParse(readBuffer);
            var knxAddress = KnxAddress.StaticParse(readBuffer);
            var projectInstallationIdentifier = ProjectInstallationIdentifier.StaticParse(readBuffer);
            var knxNetIpDeviceSerialNumber = readBuffer.ReadByteArray("knxNetIpDeviceSerialNumber", (int) (6) * 8);
            var knxNetIpDeviceMulticastAddress = IPAddress.StaticParse(readBuffer);
            var knxNetIpDeviceMacAddress = MACAddress.StaticParse(readBuffer);
            var deviceFriendlyName = readBuffer.ReadByteArray("deviceFriendlyName", (int) (30) * 8);
            return new DIBDeviceInfo(descriptionType, knxMedium, deviceStatus, knxAddress, projectInstallationIdentifier, knxNetIpDeviceSerialNumber, knxNetIpDeviceMulticastAddress, knxNetIpDeviceMacAddress, deviceFriendlyName);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("structureLength", 8, (byte) (GetLengthInBytes()));
            writeBuffer.WriteByte("descriptionType", 8, DescriptionType);
            writeBuffer.WriteByte("knxMedium", 8, (byte) KnxMedium);
            DeviceStatus.Serialize(writeBuffer);
            KnxAddress.Serialize(writeBuffer);
            ProjectInstallationIdentifier.Serialize(writeBuffer);
            writeBuffer.WriteByteArray("knxNetIpDeviceSerialNumber", KnxNetIpDeviceSerialNumber);
            KnxNetIpDeviceMulticastAddress.Serialize(writeBuffer);
            KnxNetIpDeviceMacAddress.Serialize(writeBuffer);
            writeBuffer.WriteByteArray("deviceFriendlyName", DeviceFriendlyName);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += DeviceStatus.GetLengthInBits();
            lengthInBits += KnxAddress.GetLengthInBits();
            lengthInBits += ProjectInstallationIdentifier.GetLengthInBits();
            lengthInBits += (KnxNetIpDeviceSerialNumber.Length * 8);
            lengthInBits += KnxNetIpDeviceMulticastAddress.GetLengthInBits();
            lengthInBits += KnxNetIpDeviceMacAddress.GetLengthInBits();
            lengthInBits += (DeviceFriendlyName.Length * 8);
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
