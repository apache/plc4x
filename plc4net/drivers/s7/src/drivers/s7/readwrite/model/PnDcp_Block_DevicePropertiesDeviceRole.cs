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
    public partial class PnDcp_Block_DevicePropertiesDeviceRole : PnDcp_Block
    {
        public override PnDcp_BlockOptions Option => PnDcp_BlockOptions.DEVICE_PROPERTIES_OPTION;
        public override byte Suboption => (byte) (4);

        public bool PnioSupervisor { get; }
        public bool PnioMultidevive { get; }
        public bool PnioController { get; }
        public bool PnioDevice { get; }

        public PnDcp_Block_DevicePropertiesDeviceRole(bool pnioSupervisor, bool pnioMultidevive, bool pnioController, bool pnioDevice)
        {
            PnioSupervisor = pnioSupervisor;
            PnioMultidevive = pnioMultidevive;
            PnioController = pnioController;
            PnioDevice = pnioDevice;
        }

        public static new PnDcp_Block_DevicePropertiesDeviceRole StaticParse(ReadBuffer readBuffer)
        {
            {
                var reserved = readBuffer.ReadUint("reserved", 20);
                if (!Equals(reserved, (uint) (0x000000))) { /* mspec reserved: value differs from the spec default */ }
            }
            var pnioSupervisor = readBuffer.ReadBit("pnioSupervisor");
            var pnioMultidevive = readBuffer.ReadBit("pnioMultidevive");
            var pnioController = readBuffer.ReadBit("pnioController");
            var pnioDevice = readBuffer.ReadBit("pnioDevice");
            {
                var reserved = readBuffer.ReadByte("reserved", 8);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            return new PnDcp_Block_DevicePropertiesDeviceRole(pnioSupervisor, pnioMultidevive, pnioController, pnioDevice);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUint("reserved", 20, (uint) (0x000000));
            writeBuffer.WriteBit("pnioSupervisor", PnioSupervisor);
            writeBuffer.WriteBit("pnioMultidevive", PnioMultidevive);
            writeBuffer.WriteBit("pnioController", PnioController);
            writeBuffer.WriteBit("pnioDevice", PnioDevice);
            writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 20;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
