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
    public partial class ConnectionRequest : KnxNetIpMessage
    {
        public override ushort MsgType => (ushort) (0x0205);

        public HPAIDiscoveryEndpoint HpaiDiscoveryEndpoint { get; }
        public HPAIDataEndpoint HpaiDataEndpoint { get; }
        public ConnectionRequestInformation ConnectionRequestInformation { get; }

        public ConnectionRequest(HPAIDiscoveryEndpoint hpaiDiscoveryEndpoint, HPAIDataEndpoint hpaiDataEndpoint, ConnectionRequestInformation connectionRequestInformation)
        {
            HpaiDiscoveryEndpoint = hpaiDiscoveryEndpoint;
            HpaiDataEndpoint = hpaiDataEndpoint;
            ConnectionRequestInformation = connectionRequestInformation;
        }

        public static new ConnectionRequest StaticParse(ReadBuffer readBuffer)
        {
            var hpaiDiscoveryEndpoint = HPAIDiscoveryEndpoint.StaticParse(readBuffer);
            var hpaiDataEndpoint = HPAIDataEndpoint.StaticParse(readBuffer);
            var connectionRequestInformation = ConnectionRequestInformation.StaticParse(readBuffer);
            return new ConnectionRequest(hpaiDiscoveryEndpoint, hpaiDataEndpoint, connectionRequestInformation);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            HpaiDiscoveryEndpoint.Serialize(writeBuffer);
            HpaiDataEndpoint.Serialize(writeBuffer);
            ConnectionRequestInformation.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += HpaiDiscoveryEndpoint.GetLengthInBits();
            lengthInBits += HpaiDataEndpoint.GetLengthInBits();
            lengthInBits += ConnectionRequestInformation.GetLengthInBits();
            return lengthInBits;
        }

    }
}
