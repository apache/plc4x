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
    public partial class S7PayloadUserData : S7Payload
    {
        public override byte MessageType => (byte) (0x07);

        public System.Collections.Generic.List<S7PayloadUserDataItem> Items { get; }

        public S7PayloadUserData(System.Collections.Generic.List<S7PayloadUserDataItem> items)
        {
            Items = items;
        }

        public static new S7PayloadUserData StaticParse(ReadBuffer readBuffer, byte messageType, S7Parameter parameter)
        {
            var items = new System.Collections.Generic.List<S7PayloadUserDataItem>();
            var _itemsCnt = (int) (((S7ParameterUserData) parameter).Items.Count);
            for (var _itemsI = 0; _itemsI < _itemsCnt; _itemsI++)
            {
                items.Add(S7PayloadUserDataItem.StaticParse(readBuffer, (byte) (((S7ParameterUserDataItemCPUFunctions) ((S7ParameterUserData) parameter).Items[0]).CpuFunctionGroup), (byte) (((S7ParameterUserDataItemCPUFunctions) ((S7ParameterUserData) parameter).Items[0]).CpuFunctionType), (byte) (((S7ParameterUserDataItemCPUFunctions) ((S7ParameterUserData) parameter).Items[0]).CpuSubfunction)));
            }
            return new S7PayloadUserData(items);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            foreach (var _e in Items)
            {
                _e.Serialize(writeBuffer);
            }
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += Items.Sum(_e => _e.GetLengthInBits());
            return lengthInBits;
        }

    }
}
