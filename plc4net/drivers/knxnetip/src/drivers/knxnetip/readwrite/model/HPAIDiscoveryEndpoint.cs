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
    public partial class HPAIDiscoveryEndpoint : IMessage
    {
        public HostProtocolCode HostProtocolCode { get; }
        public IPAddress IpAddress { get; }
        public ushort IpPort { get; }

        public HPAIDiscoveryEndpoint(HostProtocolCode hostProtocolCode, IPAddress ipAddress, ushort ipPort)
        {
            HostProtocolCode = hostProtocolCode;
            IpAddress = ipAddress;
            IpPort = ipPort;
        }

        public static HPAIDiscoveryEndpoint StaticParse(ReadBuffer readBuffer)
        {
            var structureLength = readBuffer.ReadByte("structureLength", 8);
            var hostProtocolCode = (HostProtocolCode) readBuffer.ReadByte("hostProtocolCode", 8);
            var ipAddress = IPAddress.StaticParse(readBuffer);
            var ipPort = readBuffer.ReadUshort("ipPort", 16);
            return new HPAIDiscoveryEndpoint(hostProtocolCode, ipAddress, ipPort);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("structureLength", 8, (byte) (GetLengthInBytes()));
            writeBuffer.WriteByte("hostProtocolCode", 8, (byte) HostProtocolCode);
            IpAddress.Serialize(writeBuffer);
            writeBuffer.WriteUshort("ipPort", 16, IpPort);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += IpAddress.GetLengthInBits();
            lengthInBits += 16;
            return lengthInBits;
        }

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
