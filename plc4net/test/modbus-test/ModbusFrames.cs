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

using System;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// Builds Modbus frames by hand, straight from the specification. The tests
    /// use these instead of the generated model so that a defect in the generated
    /// serializer cannot cancel itself out against a driver that uses the same code.
    /// </summary>
    internal static class ModbusFrames
    {
        /// <summary>MBAP header (transaction, protocol 0, length) + unit identifier + PDU.</summary>
        public static byte[] Tcp(ushort transactionId, byte unitId, params byte[] pdu)
        {
            var length = pdu.Length + 1;
            var frame = new byte[7 + pdu.Length];
            frame[0] = (byte)(transactionId >> 8);
            frame[1] = (byte)transactionId;
            frame[2] = 0;
            frame[3] = 0;
            frame[4] = (byte)(length >> 8);
            frame[5] = (byte)length;
            frame[6] = unitId;
            Array.Copy(pdu, 0, frame, 7, pdu.Length);
            return frame;
        }

        /// <summary>Address + PDU + CRC-16/MODBUS, the low CRC byte first.</summary>
        public static byte[] Rtu(byte address, params byte[] pdu)
        {
            var frame = new byte[1 + pdu.Length + 2];
            frame[0] = address;
            Array.Copy(pdu, 0, frame, 1, pdu.Length);
            var crc = Crc16(frame, frame.Length - 2);
            frame[^2] = (byte)crc;
            frame[^1] = (byte)(crc >> 8);
            return frame;
        }

        /// <summary>The reference bitwise CRC-16/MODBUS (polynomial 0xA001, initial value 0xFFFF).</summary>
        public static ushort Crc16(byte[] data, int count)
        {
            ushort crc = 0xFFFF;
            for (var i = 0; i < count; i++)
            {
                crc ^= data[i];
                for (var bit = 0; bit < 8; bit++)
                {
                    crc = (crc & 1) != 0 ? (ushort)((crc >> 1) ^ 0xA001) : (ushort)(crc >> 1);
                }
            }
            return crc;
        }
    }
}