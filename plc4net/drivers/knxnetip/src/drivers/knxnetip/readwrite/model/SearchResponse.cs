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
    public partial class SearchResponse : KnxNetIpMessage
    {
        public override ushort MsgType => (ushort) (0x0202);

        public HPAIControlEndpoint HpaiControlEndpoint { get; }
        public DIBDeviceInfo DibDeviceInfo { get; }
        public DIBSuppSvcFamilies DibSuppSvcFamilies { get; }

        public SearchResponse(HPAIControlEndpoint hpaiControlEndpoint, DIBDeviceInfo dibDeviceInfo, DIBSuppSvcFamilies dibSuppSvcFamilies)
        {
            HpaiControlEndpoint = hpaiControlEndpoint;
            DibDeviceInfo = dibDeviceInfo;
            DibSuppSvcFamilies = dibSuppSvcFamilies;
        }

        public static new SearchResponse StaticParse(ReadBuffer readBuffer)
        {
            var hpaiControlEndpoint = HPAIControlEndpoint.StaticParse(readBuffer);
            var dibDeviceInfo = DIBDeviceInfo.StaticParse(readBuffer);
            var dibSuppSvcFamilies = DIBSuppSvcFamilies.StaticParse(readBuffer);
            return new SearchResponse(hpaiControlEndpoint, dibDeviceInfo, dibSuppSvcFamilies);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            HpaiControlEndpoint.Serialize(writeBuffer);
            DibDeviceInfo.Serialize(writeBuffer);
            DibSuppSvcFamilies.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += HpaiControlEndpoint.GetLengthInBits();
            lengthInBits += DibDeviceInfo.GetLengthInBits();
            lengthInBits += DibSuppSvcFamilies.GetLengthInBits();
            return lengthInBits;
        }

    }
}
