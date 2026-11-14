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
    public partial class SlmpDoubleWordWriteItem : IMessage
    {
        public uint DeviceNumber { get; }
        public SlmpDeviceCode DeviceCode { get; }
        public uint Value { get; }

        public SlmpDoubleWordWriteItem(uint deviceNumber, SlmpDeviceCode deviceCode, uint value)
        {
            DeviceNumber = deviceNumber;
            DeviceCode = deviceCode;
            Value = value;
        }

        public static SlmpDoubleWordWriteItem StaticParse(ReadBuffer readBuffer)
        {
            var deviceNumber = readBuffer.ReadUint("deviceNumber", 24);
            var deviceCode = (SlmpDeviceCode) readBuffer.ReadByte("deviceCode", 8);
            var value = readBuffer.ReadUint("value", 32);
            return new SlmpDoubleWordWriteItem(deviceNumber, deviceCode, value);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUint("deviceNumber", 24, DeviceNumber);
            writeBuffer.WriteByte("deviceCode", 8, (byte) DeviceCode);
            writeBuffer.WriteUint("value", 32, Value);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 24;
            lengthInBits += 8;
            lengthInBits += 32;
            return lengthInBits;
        }

        public int GetLengthInBytes() => (GetLengthInBits() + 7) / 8;
    }
}
