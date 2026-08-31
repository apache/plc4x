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
    public partial class DeviceDescriptorType2 : IMessage
    {
        public ushort ManufacturerId { get; }
        public ushort DeviceType { get; }
        public byte Version { get; }
        public bool ReadSupported { get; }
        public bool WriteSupported { get; }
        public byte LogicalTagBase { get; }
        public ChannelInformation ChannelInfo1 { get; }
        public ChannelInformation ChannelInfo2 { get; }
        public ChannelInformation ChannelInfo3 { get; }
        public ChannelInformation ChannelInfo4 { get; }

        public DeviceDescriptorType2(ushort manufacturerId, ushort deviceType, byte version, bool readSupported, bool writeSupported, byte logicalTagBase, ChannelInformation channelInfo1, ChannelInformation channelInfo2, ChannelInformation channelInfo3, ChannelInformation channelInfo4)
        {
            ManufacturerId = manufacturerId;
            DeviceType = deviceType;
            Version = version;
            ReadSupported = readSupported;
            WriteSupported = writeSupported;
            LogicalTagBase = logicalTagBase;
            ChannelInfo1 = channelInfo1;
            ChannelInfo2 = channelInfo2;
            ChannelInfo3 = channelInfo3;
            ChannelInfo4 = channelInfo4;
        }

        public static DeviceDescriptorType2 StaticParse(ReadBuffer readBuffer)
        {
            var manufacturerId = readBuffer.ReadUshort("manufacturerId", 16);
            var deviceType = readBuffer.ReadUshort("deviceType", 16);
            var version = readBuffer.ReadByte("version", 8);
            var readSupported = readBuffer.ReadBit("readSupported");
            var writeSupported = readBuffer.ReadBit("writeSupported");
            var logicalTagBase = readBuffer.ReadByte("logicalTagBase", 6);
            var channelInfo1 = ChannelInformation.StaticParse(readBuffer);
            var channelInfo2 = ChannelInformation.StaticParse(readBuffer);
            var channelInfo3 = ChannelInformation.StaticParse(readBuffer);
            var channelInfo4 = ChannelInformation.StaticParse(readBuffer);
            return new DeviceDescriptorType2(manufacturerId, deviceType, version, readSupported, writeSupported, logicalTagBase, channelInfo1, channelInfo2, channelInfo3, channelInfo4);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("manufacturerId", 16, ManufacturerId);
            writeBuffer.WriteUshort("deviceType", 16, DeviceType);
            writeBuffer.WriteByte("version", 8, Version);
            writeBuffer.WriteBit("readSupported", ReadSupported);
            writeBuffer.WriteBit("writeSupported", WriteSupported);
            writeBuffer.WriteByte("logicalTagBase", 6, LogicalTagBase);
            ChannelInfo1.Serialize(writeBuffer);
            ChannelInfo2.Serialize(writeBuffer);
            ChannelInfo3.Serialize(writeBuffer);
            ChannelInfo4.Serialize(writeBuffer);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 1;
            lengthInBits += 1;
            lengthInBits += 6;
            lengthInBits += ChannelInfo1.GetLengthInBits();
            lengthInBits += ChannelInfo2.GetLengthInBits();
            lengthInBits += ChannelInfo3.GetLengthInBits();
            lengthInBits += ChannelInfo4.GetLengthInBits();
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
