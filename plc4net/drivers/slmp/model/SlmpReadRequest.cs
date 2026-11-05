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
    public partial class SlmpReadRequest : SlmpRequestData
    {
        public override ushort Command => (ushort) (0x0401);

        public uint HeadDeviceNumber { get; }
        public SlmpDeviceCode DeviceCode { get; }
        public ushort NumberOfPoints { get; }

        public SlmpReadRequest(uint headDeviceNumber, SlmpDeviceCode deviceCode, ushort numberOfPoints)
        {
            HeadDeviceNumber = headDeviceNumber;
            DeviceCode = deviceCode;
            NumberOfPoints = numberOfPoints;
        }

        public static new SlmpReadRequest StaticParse(ReadBuffer readBuffer, ushort command)
        {
            var headDeviceNumber = readBuffer.ReadUint("headDeviceNumber", 24);
            var deviceCode = (SlmpDeviceCode) readBuffer.ReadByte("deviceCode", 8);
            var numberOfPoints = readBuffer.ReadUshort("numberOfPoints", 16);
            return new SlmpReadRequest(headDeviceNumber, deviceCode, numberOfPoints);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUint("headDeviceNumber", 24, HeadDeviceNumber);
            writeBuffer.WriteByte("deviceCode", 8, (byte) DeviceCode);
            writeBuffer.WriteUshort("numberOfPoints", 16, NumberOfPoints);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 24;
            lengthInBits += 8;
            lengthInBits += 16;
            return lengthInBits;
        }

    }
}
