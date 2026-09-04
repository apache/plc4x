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
    public partial class ModbusPDUReadWriteMultipleHoldingRegistersRequest : ModbusPDU
    {
        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x17);
        public override bool Response => false;

        public ushort ReadStartingAddress { get; }
        public ushort ReadQuantity { get; }
        public ushort WriteStartingAddress { get; }
        public ushort WriteQuantity { get; }
        public byte[] Value { get; }

        public ModbusPDUReadWriteMultipleHoldingRegistersRequest(ushort readStartingAddress, ushort readQuantity, ushort writeStartingAddress, ushort writeQuantity, byte[] value)
        {
            ReadStartingAddress = readStartingAddress;
            ReadQuantity = readQuantity;
            WriteStartingAddress = writeStartingAddress;
            WriteQuantity = writeQuantity;
            Value = value;
        }

        public static new ModbusPDUReadWriteMultipleHoldingRegistersRequest StaticParse(ReadBuffer readBuffer, bool response)
        {
            var readStartingAddress = readBuffer.ReadUshort("readStartingAddress", 16);
            var readQuantity = readBuffer.ReadUshort("readQuantity", 16);
            var writeStartingAddress = readBuffer.ReadUshort("writeStartingAddress", 16);
            var writeQuantity = readBuffer.ReadUshort("writeQuantity", 16);
            var byteCount = readBuffer.ReadByte("byteCount", 8);
            var value = readBuffer.ReadByteArray("value", (int) (byteCount) * 8);
            return new ModbusPDUReadWriteMultipleHoldingRegistersRequest(readStartingAddress, readQuantity, writeStartingAddress, writeQuantity, value);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("readStartingAddress", 16, ReadStartingAddress);
            writeBuffer.WriteUshort("readQuantity", 16, ReadQuantity);
            writeBuffer.WriteUshort("writeStartingAddress", 16, WriteStartingAddress);
            writeBuffer.WriteUshort("writeQuantity", 16, WriteQuantity);
            writeBuffer.WriteByte("byteCount", 8, (byte) (Value.Length));
            writeBuffer.WriteByteArray("value", Value);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += (Value.Length * 8);
            return lengthInBits;
        }

    }
}
