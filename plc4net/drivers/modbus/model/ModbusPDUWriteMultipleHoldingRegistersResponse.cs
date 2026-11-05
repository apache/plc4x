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
    public partial class ModbusPDUWriteMultipleHoldingRegistersResponse : ModbusPDU
    {
        public override bool ErrorFlag => false;
        public override byte FunctionFlag => (byte) (0x10);
        public override bool Response => true;

        public ushort StartingAddress { get; }
        public ushort Quantity { get; }

        public ModbusPDUWriteMultipleHoldingRegistersResponse(ushort startingAddress, ushort quantity)
        {
            StartingAddress = startingAddress;
            Quantity = quantity;
        }

        public static new ModbusPDUWriteMultipleHoldingRegistersResponse StaticParse(ReadBuffer readBuffer, bool response)
        {
            var startingAddress = readBuffer.ReadUshort("startingAddress", 16);
            var quantity = readBuffer.ReadUshort("quantity", 16);
            return new ModbusPDUWriteMultipleHoldingRegistersResponse(startingAddress, quantity);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("startingAddress", 16, StartingAddress);
            writeBuffer.WriteUshort("quantity", 16, Quantity);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 16;
            return lengthInBits;
        }

    }
}
