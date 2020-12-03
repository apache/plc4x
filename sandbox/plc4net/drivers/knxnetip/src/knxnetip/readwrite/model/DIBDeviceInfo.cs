//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

namespace org.apache.plc4net.drivers.knxnetip.readwrite.model
{

    public class DIBDeviceInfo
    {

        // Properties.
        public byte DescriptionType { get; }
        public KnxMedium KnxMedium { get; }
        public DeviceStatus DeviceStatus { get; }
        public KnxAddress KnxAddress { get; }
        public ProjectInstallationIdentifier ProjectInstallationIdentifier { get; }
        public sbyte[] KnxNetIpDeviceSerialNumber { get; }
        public IPAddress KnxNetIpDeviceMulticastAddress { get; }
        public MACAddress KnxNetIpDeviceMacAddress { get; }
        public sbyte[] DeviceFriendlyName { get; }

        public DIBDeviceInfo(byte descriptionType, KnxMedium knxMedium, DeviceStatus deviceStatus, KnxAddress knxAddress, ProjectInstallationIdentifier projectInstallationIdentifier, sbyte[] knxNetIpDeviceSerialNumber, IPAddress knxNetIpDeviceMulticastAddress, MACAddress knxNetIpDeviceMacAddress, sbyte[] deviceFriendlyName)
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

    }

}
