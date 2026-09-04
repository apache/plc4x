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
    public partial class S7PayloadReadVarResponse : S7Payload
    {
        public override byte MessageType => (byte) (0x03);

        public System.Collections.Generic.List<S7VarPayloadDataItem> Items { get; }

        public S7PayloadReadVarResponse(System.Collections.Generic.List<S7VarPayloadDataItem> items)
        {
            Items = items;
        }

        public static new S7PayloadReadVarResponse StaticParse(ReadBuffer readBuffer, byte messageType, S7Parameter parameter)
        {
            var items = new System.Collections.Generic.List<S7VarPayloadDataItem>();
            var _itemsCnt = (int) (((S7ParameterReadVarResponse) parameter).NumItems);
            for (var _itemsI = 0; _itemsI < _itemsCnt; _itemsI++)
            {
                items.Add(S7VarPayloadDataItem.StaticParse(readBuffer, _itemsI == _itemsCnt - 1));
            }
            return new S7PayloadReadVarResponse(items);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            for (var _i = 0; _i < Items.Count; _i++)
            {
                Items[_i].Serialize(writeBuffer, _i == Items.Count - 1);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += System.Linq.Enumerable.Range(0, Items.Count).Sum(_i => Items[_i].GetLengthInBits(_i == Items.Count - 1));
            return lengthInBits;
        }

    }
}
