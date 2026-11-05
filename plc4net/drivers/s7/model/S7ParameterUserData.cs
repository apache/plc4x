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
    public partial class S7ParameterUserData : S7Parameter
    {
        public override byte ParameterType => (byte) (0x00);
        public override byte MessageType => (byte) (0x07);

        public System.Collections.Generic.List<S7ParameterUserDataItem> Items { get; }

        public S7ParameterUserData(System.Collections.Generic.List<S7ParameterUserDataItem> items)
        {
            Items = items;
        }

        public static new S7ParameterUserData StaticParse(ReadBuffer readBuffer, byte messageType)
        {
            var numItems = readBuffer.ReadByte("numItems", 8);
            var items = new System.Collections.Generic.List<S7ParameterUserDataItem>();
            var _itemsCnt = (int) (numItems);
            for (var _itemsI = 0; _itemsI < _itemsCnt; _itemsI++)
            {
                items.Add(S7ParameterUserDataItem.StaticParse(readBuffer));
            }
            return new S7ParameterUserData(items);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("numItems", 8, (byte) (Items.Count));
            foreach (var _e in Items)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 8;
            lengthInBits += Items.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
