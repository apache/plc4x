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

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    public partial class ModbusPDUReadFileRecordRequest : ModbusPDU
    {
        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x14);
        public override bool Response => false;

        public System.Collections.Generic.List<ModbusPDUReadFileRecordRequestItem> Items { get; }

        public ModbusPDUReadFileRecordRequest(System.Collections.Generic.List<ModbusPDUReadFileRecordRequestItem> items)
        {
            Items = items;
        }

        public static new ModbusPDUReadFileRecordRequest StaticParse(ReadBuffer readBuffer, bool response)
        {
            var byteCount = readBuffer.ReadByte("byteCount", 8);
            var items = new System.Collections.Generic.List<ModbusPDUReadFileRecordRequestItem>();
            var _itemsEnd = readBuffer.GetPos() + (int) (byteCount) * 8;
            while (readBuffer.GetPos() < _itemsEnd)
            {
                items.Add(ModbusPDUReadFileRecordRequestItem.StaticParse(readBuffer));
            }
            return new ModbusPDUReadFileRecordRequest(items);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteByte("byteCount", 8, (byte) (ModbusStaticHelper.ArraySizeInBytes(Items)));
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
