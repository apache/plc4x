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
    public partial class ComObjectTableRealisationType2 : ComObjectTable
    {
        public override FirmwareType FirmwareType => FirmwareType.SYSTEM_2;

        public byte NumEntries { get; }
        public byte RamFlagsTablePointer { get; }
        public System.Collections.Generic.List<GroupObjectDescriptorRealisationType2> ComObjectDescriptors { get; }

        public ComObjectTableRealisationType2(byte numEntries, byte ramFlagsTablePointer, System.Collections.Generic.List<GroupObjectDescriptorRealisationType2> comObjectDescriptors)
        {
            NumEntries = numEntries;
            RamFlagsTablePointer = ramFlagsTablePointer;
            ComObjectDescriptors = comObjectDescriptors;
        }

        public static new ComObjectTableRealisationType2 StaticParse(ReadBuffer readBuffer, FirmwareType firmwareType)
        {
            var numEntries = readBuffer.ReadByte("numEntries", 8);
            var ramFlagsTablePointer = readBuffer.ReadByte("ramFlagsTablePointer", 8);
            var comObjectDescriptors = new System.Collections.Generic.List<GroupObjectDescriptorRealisationType2>();
            var _comObjectDescriptorsCnt = (int) (numEntries);
            for (var _comObjectDescriptorsI = 0; _comObjectDescriptorsI < _comObjectDescriptorsCnt; _comObjectDescriptorsI++)
            {
                comObjectDescriptors.Add(GroupObjectDescriptorRealisationType2.StaticParse(readBuffer));
            }
            return new ComObjectTableRealisationType2(numEntries, ramFlagsTablePointer, comObjectDescriptors);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("numEntries", 8, NumEntries);
            writeBuffer.WriteByte("ramFlagsTablePointer", 8, RamFlagsTablePointer);
            foreach (var _e in ComObjectDescriptors)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += ComObjectDescriptors.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
