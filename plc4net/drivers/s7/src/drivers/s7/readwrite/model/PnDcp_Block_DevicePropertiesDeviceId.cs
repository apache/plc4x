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
    public partial class PnDcp_Block_DevicePropertiesDeviceId : PnDcp_Block
    {
        public override PnDcp_BlockOptions Option => PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION;
        public override byte Suboption => (byte) (3);

        public ushort VendorId { get; }
        public ushort DeviceId { get; }

        public PnDcp_Block_DevicePropertiesDeviceId(ushort vendorId, ushort deviceId)
        {
            VendorId = vendorId;
            DeviceId = deviceId;
        }

        public static new PnDcp_Block_DevicePropertiesDeviceId StaticParse(ReadBuffer readBuffer)
        {
            {
                var reserved = readBuffer.ReadUshort("reserved", 16);
                if (!Equals(reserved, (ushort) (0x0000))) { /* mspec reserved: value differs from the spec default */ }
            }
            var vendorId = readBuffer.ReadUshort("vendorId", 16);
            var deviceId = readBuffer.ReadUshort("deviceId", 16);
            return new PnDcp_Block_DevicePropertiesDeviceId(vendorId, deviceId);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("reserved", 16, (ushort) (0x0000));
            writeBuffer.WriteUshort("vendorId", 16, VendorId);
            writeBuffer.WriteUshort("deviceId", 16, DeviceId);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            return lengthInBits;
        }

    }
}
