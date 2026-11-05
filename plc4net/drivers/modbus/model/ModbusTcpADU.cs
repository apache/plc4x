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
    public partial class ModbusTcpADU : ModbusADU
    {
        public const ushort ProtocolIdentifier = 0x0000;

        public override DriverType DriverType => DriverType.MODBUS_TCP;

        public ushort TransactionIdentifier { get; }
        public byte UnitIdentifier { get; }
        public ModbusPDU Pdu { get; }

        public ModbusTcpADU(ushort transactionIdentifier, byte unitIdentifier, ModbusPDU pdu)
        {
            TransactionIdentifier = transactionIdentifier;
            UnitIdentifier = unitIdentifier;
            Pdu = pdu;
        }

        public static new ModbusTcpADU StaticParse(ReadBuffer readBuffer, DriverType driverType, bool response)
        {
            var transactionIdentifier = readBuffer.ReadUshort("transactionIdentifier", 16);
            var protocolIdentifier = readBuffer.ReadUshort("protocolIdentifier", 16);
            if (!Equals(protocolIdentifier, (ushort) (0x0000)))
                throw new ParseException($"Expected constant {ProtocolIdentifier} for 'protocolIdentifier' but got {protocolIdentifier}");
            var length = readBuffer.ReadUshort("length", 16);
            var unitIdentifier = readBuffer.ReadByte("unitIdentifier", 8);
            var pdu = ModbusPDU.StaticParse(readBuffer, response);
            return new ModbusTcpADU(transactionIdentifier, unitIdentifier, pdu);
        }

        protected override void SerializeChild(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("transactionIdentifier", 16, TransactionIdentifier);
            writeBuffer.WriteUshort("protocolIdentifier", 16, ProtocolIdentifier);
            writeBuffer.WriteUshort("length", 16, (ushort) ((Pdu.GetLengthInBytes() + 1)));
            writeBuffer.WriteByte("unitIdentifier", 8, UnitIdentifier);
            Pdu.Serialize(writeBuffer);
        }

        protected override int GetLengthInBitsChild()
        {
            var lengthInBits = 0;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 16;
            lengthInBits += 8;
            lengthInBits += Pdu.GetLengthInBits();
            return lengthInBits;
        }

    }
}
