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

using org.apache.plc4net.drivers.modbus.messages;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    /// <summary>
    /// Hand-written bodies for the <c>STATIC_CALL</c> targets in modbus.mspec.
    /// The generated <see cref="ModbusStaticHelper"/> declares
    /// <c>params object[]</c> stubs that throw; these typed overloads are the
    /// ones the generated <c>ModbusRtuADU</c> / <c>ModbusAsciiADU</c> bind to.
    /// Ported from
    /// <c>plc4j/.../modbus/readwrite/utils/StaticHelper.java</c>.
    /// </summary>
    public static partial class ModbusStaticHelper
    {
        private static byte[] AddressAndPdu(byte address, ModbusPDU pdu)
        {
            var wb = new WriteBuffer();
            wb.WriteByte("address", 8, address);
            pdu.Serialize(wb);
            return wb.GetBytes();
        }

        /// <summary>
        /// CRC-16 (poly 0xA001, init 0xFFFF) over the address + PDU, returned so
        /// that a big-endian 16-bit write lays it on the wire low byte first —
        /// the Modbus RTU convention.
        /// </summary>
        public static ushort RtuCrcCheck(byte address, ModbusPDU pdu)
        {
            var crc = ModbusCRC.Compute(AddressAndPdu(address, pdu)); // { low, high }
            return (ushort) ((crc[0] << 8) | crc[1]);
        }

        /// <summary>
        /// LRC-8 (two's complement of the byte sum) over the address + PDU, per
        /// the Modbus over Serial Line spec.
        /// </summary>
        public static byte AsciiLrcCheck(byte address, ModbusPDU pdu)
        {
            byte lrc = 0;
            foreach (var b in AddressAndPdu(address, pdu))
            {
                lrc += b;
            }
            return (byte) (-lrc);
        }
    }
}
