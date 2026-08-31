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
    public partial class SzlDataTreeItem : IMessage
    {
        public ushort ItemIndex { get; }
        public byte[] Mlfb { get; }
        public ushort ModuleTypeId { get; }
        public ushort Ausbg { get; }
        public ushort Ausbe { get; }

        public SzlDataTreeItem(ushort itemIndex, byte[] mlfb, ushort moduleTypeId, ushort ausbg, ushort ausbe)
        {
            ItemIndex = itemIndex;
            Mlfb = mlfb;
            ModuleTypeId = moduleTypeId;
            Ausbg = ausbg;
            Ausbe = ausbe;
        }

        public static SzlDataTreeItem StaticParse(ReadBuffer readBuffer)
        {
            var itemIndex = readBuffer.ReadUshort("itemIndex", 16);
            var mlfb = readBuffer.ReadByteArray("mlfb", (int) (20) * 8);
            var moduleTypeId = readBuffer.ReadUshort("moduleTypeId", 16);
            var ausbg = readBuffer.ReadUshort("ausbg", 16);
            var ausbe = readBuffer.ReadUshort("ausbe", 16);
            return new SzlDataTreeItem(itemIndex, mlfb, moduleTypeId, ausbg, ausbe);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("itemIndex", 16, ItemIndex);
            writeBuffer.WriteByteArray("mlfb", Mlfb);
            writeBuffer.WriteUshort("moduleTypeId", 16, ModuleTypeId);
            writeBuffer.WriteUshort("ausbg", 16, Ausbg);
            writeBuffer.WriteUshort("ausbe", 16, Ausbe);
        }

        public int GetLengthInBits()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += (Mlfb.Length * 8);
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            return lengthInBits;
        }

        public int GetLengthInBytes() => GetLengthInBits() / 8;
    }
}
