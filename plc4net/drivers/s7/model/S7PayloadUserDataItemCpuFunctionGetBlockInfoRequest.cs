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

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    public partial class S7PayloadUserDataItemCpuFunctionGetBlockInfoRequest : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x03);
        public override byte CpuFunctionType => (byte) (0x04);
        public override byte CpuSubfunction => (byte) (0x03);

        public ushort BlockType { get; }
        public string BlockNumber { get; }
        public byte Filesystem { get; }

        public S7PayloadUserDataItemCpuFunctionGetBlockInfoRequest(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, ushort blockType, string blockNumber, byte filesystem) : base(returnCode, transportSize, dataLength)
        {
            BlockType = blockType;
            BlockNumber = blockNumber;
            Filesystem = filesystem;
        }

        public static S7PayloadUserDataItemCpuFunctionGetBlockInfoRequest StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var blockType = readBuffer.ReadUshort("blockType", 16);
            var blockNumber = readBuffer.ReadString("blockNumber", 40, System.Text.Encoding.UTF8);
            var filesystem = readBuffer.ReadByte("filesystem", 8);
            return new S7PayloadUserDataItemCpuFunctionGetBlockInfoRequest(returnCode, transportSize, dataLength, blockType, blockNumber, filesystem);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("blockType", 16, BlockType);
            writeBuffer.WriteString("blockNumber", 40, "UTF8", BlockNumber);
            writeBuffer.WriteByte("filesystem", 8, Filesystem);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 40;
            lengthInBits += 8;
            return lengthInBits;
        }

    }
}
