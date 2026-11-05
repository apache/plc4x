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
    public partial class SlmpMultiBlockReadRequest : SlmpRequestData
    {
        public override ushort Command => (ushort) (0x0406);

        public byte NumberOfWordDeviceBlocks { get; }
        public byte NumberOfBitDeviceBlocks { get; }
        public System.Collections.Generic.List<SlmpDeviceBlock> WordDeviceBlocks { get; }
        public System.Collections.Generic.List<SlmpDeviceBlock> BitDeviceBlocks { get; }

        public SlmpMultiBlockReadRequest(byte numberOfWordDeviceBlocks, byte numberOfBitDeviceBlocks, System.Collections.Generic.List<SlmpDeviceBlock> wordDeviceBlocks, System.Collections.Generic.List<SlmpDeviceBlock> bitDeviceBlocks)
        {
            NumberOfWordDeviceBlocks = numberOfWordDeviceBlocks;
            NumberOfBitDeviceBlocks = numberOfBitDeviceBlocks;
            WordDeviceBlocks = wordDeviceBlocks;
            BitDeviceBlocks = bitDeviceBlocks;
        }

        public static new SlmpMultiBlockReadRequest StaticParse(ReadBuffer readBuffer, ushort command)
        {
            var numberOfWordDeviceBlocks = readBuffer.ReadByte("numberOfWordDeviceBlocks", 8);
            var numberOfBitDeviceBlocks = readBuffer.ReadByte("numberOfBitDeviceBlocks", 8);
            var wordDeviceBlocks = new System.Collections.Generic.List<SlmpDeviceBlock>();
            var _wordDeviceBlocksCnt = (int) (numberOfWordDeviceBlocks);
            for (var _wordDeviceBlocksI = 0; _wordDeviceBlocksI < _wordDeviceBlocksCnt; _wordDeviceBlocksI++)
            {
                wordDeviceBlocks.Add(SlmpDeviceBlock.StaticParse(readBuffer));
            }
            var bitDeviceBlocks = new System.Collections.Generic.List<SlmpDeviceBlock>();
            var _bitDeviceBlocksCnt = (int) (numberOfBitDeviceBlocks);
            for (var _bitDeviceBlocksI = 0; _bitDeviceBlocksI < _bitDeviceBlocksCnt; _bitDeviceBlocksI++)
            {
                bitDeviceBlocks.Add(SlmpDeviceBlock.StaticParse(readBuffer));
            }
            return new SlmpMultiBlockReadRequest(numberOfWordDeviceBlocks, numberOfBitDeviceBlocks, wordDeviceBlocks, bitDeviceBlocks);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("numberOfWordDeviceBlocks", 8, NumberOfWordDeviceBlocks);
            writeBuffer.WriteByte("numberOfBitDeviceBlocks", 8, NumberOfBitDeviceBlocks);
            foreach (var _e in WordDeviceBlocks)
            {
                _e.Serialize(writeBuffer);
            }
            foreach (var _e in BitDeviceBlocks)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += WordDeviceBlocks.Sum(_e => _e.GetLengthInBits());
            lengthInBits += BitDeviceBlocks.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
