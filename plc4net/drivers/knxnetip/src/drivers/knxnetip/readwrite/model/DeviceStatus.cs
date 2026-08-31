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
    public partial class DeviceStatus : IMessage
    {
        public bool ProgramMode { get; }

        public DeviceStatus(bool programMode)
        {
            ProgramMode = programMode;
        }

        public static DeviceStatus StaticParse(ReadBuffer readBuffer)
        {
            {
                var reserved = readBuffer.ReadByte("reserved", 7);
                if (!Equals(reserved, (byte) (0x00))) { /* mspec reserved: value differs from the spec default */ }
            }
            var programMode = readBuffer.ReadBit("programMode");
            return new DeviceStatus(programMode);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("reserved", 7, (byte) (0x00));
            writeBuffer.WriteBit("programMode", ProgramMode);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 7;
            lengthInBits += 1;
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
