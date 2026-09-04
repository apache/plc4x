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
    public partial class PnDcp_Block_IpParameter : PnDcp_Block
    {
        public override PnDcp_BlockOptions Option => PnDcp_BlockOptions.IP_OPTION;
        public override byte Suboption => (byte) (2);

        public bool IpConflictDetected { get; }
        public bool SetViaDhcp { get; }
        public bool SetManually { get; }
        public byte[] IpAddress { get; }
        public byte[] SubnetMask { get; }
        public byte[] StandardGateway { get; }

        public PnDcp_Block_IpParameter(bool ipConflictDetected, bool setViaDhcp, bool setManually, byte[] ipAddress, byte[] subnetMask, byte[] standardGateway)
        {
            IpConflictDetected = ipConflictDetected;
            SetViaDhcp = setViaDhcp;
            SetManually = setManually;
            IpAddress = ipAddress;
            SubnetMask = subnetMask;
            StandardGateway = standardGateway;
        }

        public static new PnDcp_Block_IpParameter StaticParse(ReadBuffer readBuffer)
        {
            {
                var reserved = readBuffer.ReadByte("reserved", 8);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var ipConflictDetected = readBuffer.ReadBit("ipConflictDetected");
            {
                var reserved = readBuffer.ReadByte("reserved", 5);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var setViaDhcp = readBuffer.ReadBit("setViaDhcp");
            var setManually = readBuffer.ReadBit("setManually");
            var ipAddress = readBuffer.ReadByteArray("ipAddress", (int) (4) * 8);
            var subnetMask = readBuffer.ReadByteArray("subnetMask", (int) (4) * 8);
            var standardGateway = readBuffer.ReadByteArray("standardGateway", (int) (4) * 8);
            return new PnDcp_Block_IpParameter(ipConflictDetected, setViaDhcp, setManually, ipAddress, subnetMask, standardGateway);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("reserved", 8, (byte) (0x00));
            writeBuffer.WriteBit("ipConflictDetected", IpConflictDetected);
            writeBuffer.WriteByte("reserved", 5, (byte) (0x00));
            writeBuffer.WriteBit("setViaDhcp", SetViaDhcp);
            writeBuffer.WriteBit("setManually", SetManually);
            writeBuffer.WriteByteArray("ipAddress", IpAddress);
            writeBuffer.WriteByteArray("subnetMask", SubnetMask);
            writeBuffer.WriteByteArray("standardGateway", StandardGateway);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 1;
            lengthInBits += 5;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += (IpAddress.Length * 8);
            lengthInBits += (SubnetMask.Length * 8);
            lengthInBits += (StandardGateway.Length * 8);
            return lengthInBits;
        }

    }
}
