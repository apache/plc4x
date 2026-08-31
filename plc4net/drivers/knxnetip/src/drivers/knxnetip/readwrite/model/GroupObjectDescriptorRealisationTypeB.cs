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
    public partial class GroupObjectDescriptorRealisationTypeB : IMessage
    {
        public bool UpdateEnable { get; }
        public bool TransmitEnable { get; }
        public bool SegmentSelectorEnable { get; }
        public bool WriteEnable { get; }
        public bool ReadEnable { get; }
        public bool CommunicationEnable { get; }
        public CEMIPriority Priority { get; }
        public ComObjectValueType ValueType { get; }

        public GroupObjectDescriptorRealisationTypeB(bool updateEnable, bool transmitEnable, bool segmentSelectorEnable, bool writeEnable, bool readEnable, bool communicationEnable, CEMIPriority priority, ComObjectValueType valueType)
        {
            UpdateEnable = updateEnable;
            TransmitEnable = transmitEnable;
            SegmentSelectorEnable = segmentSelectorEnable;
            WriteEnable = writeEnable;
            ReadEnable = readEnable;
            CommunicationEnable = communicationEnable;
            Priority = priority;
            ValueType = valueType;
        }

        public static GroupObjectDescriptorRealisationTypeB StaticParse(ReadBuffer readBuffer)
        {
            var updateEnable = readBuffer.ReadBit("updateEnable");
            var transmitEnable = readBuffer.ReadBit("transmitEnable");
            var segmentSelectorEnable = readBuffer.ReadBit("segmentSelectorEnable");
            var writeEnable = readBuffer.ReadBit("writeEnable");
            var readEnable = readBuffer.ReadBit("readEnable");
            var communicationEnable = readBuffer.ReadBit("communicationEnable");
            var priority = (CEMIPriority) readBuffer.ReadByte("priority", 2);
            var valueType = (ComObjectValueType) readBuffer.ReadByte("valueType", 8);
            return new GroupObjectDescriptorRealisationTypeB(updateEnable, transmitEnable, segmentSelectorEnable, writeEnable, readEnable, communicationEnable, priority, valueType);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteBit("updateEnable", UpdateEnable);
            writeBuffer.WriteBit("transmitEnable", TransmitEnable);
            writeBuffer.WriteBit("segmentSelectorEnable", SegmentSelectorEnable);
            writeBuffer.WriteBit("writeEnable", WriteEnable);
            writeBuffer.WriteBit("readEnable", ReadEnable);
            writeBuffer.WriteBit("communicationEnable", CommunicationEnable);
            writeBuffer.WriteByte("priority", 2, (byte) Priority);
            writeBuffer.WriteByte("valueType", 8, (byte) ValueType);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 2;
            lengthInBits += 8;
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
