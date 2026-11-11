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

namespace org.apache.plc4net.drivers.slmp.readwrite.model
{
    public partial class SlmpRandomReadRequest : SlmpRequestData
    {
        public override ushort Command => (ushort) (0x0403);

        public byte NumberOfWordAccessPoints { get; }
        public byte NumberOfDoubleWordAccessPoints { get; }
        public System.Collections.Generic.List<SlmpDeviceSpec> WordAccessDevices { get; }
        public System.Collections.Generic.List<SlmpDeviceSpec> DoubleWordAccessDevices { get; }

        public SlmpRandomReadRequest(byte numberOfWordAccessPoints, byte numberOfDoubleWordAccessPoints, System.Collections.Generic.List<SlmpDeviceSpec> wordAccessDevices, System.Collections.Generic.List<SlmpDeviceSpec> doubleWordAccessDevices)
        {
            NumberOfWordAccessPoints = numberOfWordAccessPoints;
            NumberOfDoubleWordAccessPoints = numberOfDoubleWordAccessPoints;
            WordAccessDevices = wordAccessDevices;
            DoubleWordAccessDevices = doubleWordAccessDevices;
        }

        public static new SlmpRandomReadRequest StaticParse(ReadBuffer readBuffer, ushort command, ushort subCommand)
        {
            var numberOfWordAccessPoints = readBuffer.ReadByte("numberOfWordAccessPoints", 8);
            var numberOfDoubleWordAccessPoints = readBuffer.ReadByte("numberOfDoubleWordAccessPoints", 8);
            var wordAccessDevices = new System.Collections.Generic.List<SlmpDeviceSpec>();
            var _wordAccessDevicesCnt = (int) (numberOfWordAccessPoints);
            for (var _wordAccessDevicesI = 0; _wordAccessDevicesI < _wordAccessDevicesCnt; _wordAccessDevicesI++)
            {
                wordAccessDevices.Add(SlmpDeviceSpec.StaticParse(readBuffer));
            }
            var doubleWordAccessDevices = new System.Collections.Generic.List<SlmpDeviceSpec>();
            var _doubleWordAccessDevicesCnt = (int) (numberOfDoubleWordAccessPoints);
            for (var _doubleWordAccessDevicesI = 0; _doubleWordAccessDevicesI < _doubleWordAccessDevicesCnt; _doubleWordAccessDevicesI++)
            {
                doubleWordAccessDevices.Add(SlmpDeviceSpec.StaticParse(readBuffer));
            }
            return new SlmpRandomReadRequest(numberOfWordAccessPoints, numberOfDoubleWordAccessPoints, wordAccessDevices, doubleWordAccessDevices);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("numberOfWordAccessPoints", 8, NumberOfWordAccessPoints);
            writeBuffer.WriteByte("numberOfDoubleWordAccessPoints", 8, NumberOfDoubleWordAccessPoints);
            foreach (var _e in WordAccessDevices)
            {
                _e.Serialize(writeBuffer);
            }
            foreach (var _e in DoubleWordAccessDevices)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += WordAccessDevices.Sum(_e => _e.GetLengthInBits());
            lengthInBits += DoubleWordAccessDevices.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
