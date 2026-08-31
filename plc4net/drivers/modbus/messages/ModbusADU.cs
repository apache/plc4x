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

using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.modbus.messages
{
    /// <summary>
    /// A Modbus TCP Application Data Unit — the MBAP header plus a PDU.
    /// </summary>
    public class ModbusADU : IMessage
    {
        public ushort TransactionId { get; set; }
        public ushort ProtocolId { get; set; }
        public ushort Length { get; set; }
        public byte UnitId { get; set; }
        public byte FunctionCode { get; set; }
        public byte[] Data { get; set; } = System.Array.Empty<byte>();

        public int GetLengthInBytes()
        {
            return 7 + 1 + (Data?.Length ?? 0);
        }

        public void Serialize(WriteBuffer writeBuffer)
        {
            writeBuffer.WriteUshort("transactionId", TransactionId, 16);
            writeBuffer.WriteUshort("protocolId", ProtocolId, 16);
            writeBuffer.WriteUshort("length", Length, 16);
            writeBuffer.WriteByte("unitId", UnitId, 8);
            writeBuffer.WriteByte("functionCode", FunctionCode, 8);
            if (Data != null && Data.Length > 0)
                writeBuffer.WriteByteArray("data", Data);
        }

        /// <summary>Parses a Modbus TCP ADU from a buffer.</summary>
        public static ModbusADU StaticParse(ReadBuffer readBuffer, bool response)
        {
            var adu = new ModbusADU
            {
                TransactionId = readBuffer.ReadUshort("transactionId", 16),
                ProtocolId = readBuffer.ReadUshort("protocolId", 16)
            };
            adu.Length = readBuffer.ReadUshort("length", 16);
            adu.UnitId = readBuffer.ReadByte("unitId", 8);
            adu.FunctionCode = readBuffer.ReadByte("functionCode", 8);

            var dataLen = adu.Length - 1;
            if (dataLen > 0)
                adu.Data = readBuffer.ReadByteArray("data", dataLen * 8);

            return adu;
        }

        }
}
