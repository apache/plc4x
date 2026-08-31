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
    public partial class DIBSuppSvcFamilies : IMessage
    {
        public byte DescriptionType { get; }
        public System.Collections.Generic.List<ServiceId> ServiceIds { get; }

        public DIBSuppSvcFamilies(byte descriptionType, System.Collections.Generic.List<ServiceId> serviceIds)
        {
            DescriptionType = descriptionType;
            ServiceIds = serviceIds;
        }

        public static DIBSuppSvcFamilies StaticParse(ReadBuffer readBuffer)
        {
            var structureLength = readBuffer.ReadByte("structureLength", 8);
            var descriptionType = readBuffer.ReadByte("descriptionType", 8);
            var serviceIds = new System.Collections.Generic.List<ServiceId>();
            var _serviceIdsEnd = readBuffer.GetPos() + (int) ((structureLength - 2)) * 8;
            while (readBuffer.GetPos() < _serviceIdsEnd)
            {
                serviceIds.Add(ServiceId.StaticParse(readBuffer));
            }
            return new DIBSuppSvcFamilies(descriptionType, serviceIds);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("structureLength", 8, (byte) (GetLengthInBytes()));
            writeBuffer.WriteByte("descriptionType", 8, DescriptionType);
            foreach (var _e in ServiceIds)
            {
                _e.Serialize(writeBuffer);
            }
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += ServiceIds.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
