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
    public abstract partial class ServiceId : IMessage
    {
        public abstract byte ServiceType { get; }

        public static ServiceId StaticParse(ReadBuffer readBuffer)
        {
            var serviceType = readBuffer.ReadByte("serviceType", 8);
            if (Equals(serviceType, (byte) (0x02)))
            {
                return KnxNetIpCore.StaticParse(readBuffer);
            }
            if (Equals(serviceType, (byte) (0x03)))
            {
                return KnxNetIpDeviceManagement.StaticParse(readBuffer);
            }
            if (Equals(serviceType, (byte) (0x04)))
            {
                return KnxNetIpTunneling.StaticParse(readBuffer);
            }
            if (Equals(serviceType, (byte) (0x05)))
            {
                return KnxNetIpRouting.StaticParse(readBuffer);
            }
            if (Equals(serviceType, (byte) (0x06)))
            {
                return KnxNetRemoteLogging.StaticParse(readBuffer);
            }
            if (Equals(serviceType, (byte) (0x07)))
            {
                return KnxNetRemoteConfigurationAndDiagnosis.StaticParse(readBuffer);
            }
            if (Equals(serviceType, (byte) (0x08)))
            {
                return KnxNetObjectServer.StaticParse(readBuffer);
            }
            throw new ParseException("No matching subtype found for ServiceId");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("serviceType", 8, ServiceType);
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
