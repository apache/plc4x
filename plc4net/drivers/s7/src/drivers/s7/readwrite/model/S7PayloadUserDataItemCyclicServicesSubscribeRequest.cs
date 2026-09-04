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
    public partial class S7PayloadUserDataItemCyclicServicesSubscribeRequest : S7PayloadUserDataItem
    {
        public override byte CpuFunctionGroup => (byte) (0x02);
        public override byte CpuFunctionType => (byte) (0x04);
        public override byte CpuSubfunction => (byte) (0x01);

        public ushort ItemsCount { get; }
        public TimeBase TimeBase { get; }
        public byte TimeFactor { get; }
        public System.Collections.Generic.List<CycServiceItemType> Item { get; }

        public S7PayloadUserDataItemCyclicServicesSubscribeRequest(DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength, ushort itemsCount, TimeBase timeBase, byte timeFactor, System.Collections.Generic.List<CycServiceItemType> item) : base(returnCode, transportSize, dataLength)
        {
            ItemsCount = itemsCount;
            TimeBase = timeBase;
            TimeFactor = timeFactor;
            Item = item;
        }

        public static S7PayloadUserDataItemCyclicServicesSubscribeRequest StaticParse(ReadBuffer readBuffer, byte cpuFunctionGroup, byte cpuFunctionType, byte cpuSubfunction, DataTransportErrorCode returnCode, DataTransportSize transportSize, ushort dataLength)
        {
            var itemsCount = readBuffer.ReadUshort("itemsCount", 16);
            var timeBase = (TimeBase) readBuffer.ReadByte("timeBase", 8);
            var timeFactor = readBuffer.ReadByte("timeFactor", 8);
            var item = new System.Collections.Generic.List<CycServiceItemType>();
            var _itemCnt = (int) (itemsCount);
            for (var _itemI = 0; _itemI < _itemCnt; _itemI++)
            {
                item.Add(CycServiceItemType.StaticParse(readBuffer));
            }
            return new S7PayloadUserDataItemCyclicServicesSubscribeRequest(returnCode, transportSize, dataLength, itemsCount, timeBase, timeFactor, item);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("itemsCount", 16, ItemsCount);
            writeBuffer.WriteByte("timeBase", 8, (byte) TimeBase);
            writeBuffer.WriteByte("timeFactor", 8, TimeFactor);
            foreach (var _e in Item)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += 8;
            lengthInBits += Item.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
