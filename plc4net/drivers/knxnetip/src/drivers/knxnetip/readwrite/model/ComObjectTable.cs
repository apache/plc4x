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
    public abstract partial class ComObjectTable : IMessage
    {
        public abstract FirmwareType FirmwareType { get; }

        public static ComObjectTable StaticParse(ReadBuffer readBuffer, FirmwareType firmwareType)
        {
            if (Equals(firmwareType, FirmwareType.SYSTEM_1))
            {
                return ComObjectTableRealisationType1.StaticParse(readBuffer, firmwareType);
            }
            if (Equals(firmwareType, FirmwareType.SYSTEM_2))
            {
                return ComObjectTableRealisationType2.StaticParse(readBuffer, firmwareType);
            }
            if (Equals(firmwareType, FirmwareType.SYSTEM_300))
            {
                return ComObjectTableRealisationType6.StaticParse(readBuffer, firmwareType);
            }
            throw new ParseException("No matching subtype found for ComObjectTable");
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            SerializeChild(writeBuffer);
        }

        protected abstract void SerializeChild(WriteBuffer writeBuffer);

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += GetLengthInBitsChild();
            return lengthInBits;
        }

        protected abstract int GetLengthInBitsChild();

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
