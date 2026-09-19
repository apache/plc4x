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

namespace org.apache.plc4net.drivers.modbus.messages
{
    /// <summary>
    /// CRC-16-ANSI (Modbus CRC) as specified in the Modbus over Serial Line
    /// protocol.  Polynomial 0x8005, initial 0xFFFF, reflected, no XOR-out.
    /// </summary>
    public static class ModbusCRC
    {
        // Pre-computed lookup table for reflected polynomial 0xA001.
        private static readonly ushort[] _table = BuildTable();

        private static ushort[] BuildTable()
        {
            var table = new ushort[256];
            for (int i = 0; i < 256; i++)
            {
                ushort crc = (ushort)i;
                for (int j = 0; j < 8; j++)
                {
                    if ((crc & 0x0001) != 0)
                        crc = (ushort)((crc >> 1) ^ 0xA001);
                    else
                        crc >>= 1;
                }
                table[i] = crc;
            }
            return table;
        }

        /// <summary>
        /// Computes the CRC-16 over the given bytes and returns it as
        /// a little-endian pair (low byte first, per Modbus RTU convention).
        /// </summary>
        public static byte[] Compute(byte[] data, int offset, int length)
        {
            ushort crc = 0xFFFF;
            for (int i = offset; i < offset + length; i++)
                crc = (ushort)((crc >> 8) ^ _table[(crc ^ data[i]) & 0xFF]);
            // Modbus RTU puts CRC on the wire as low byte first.
            return new[] { (byte)(crc & 0xFF), (byte)(crc >> 8) };
        }

        public static byte[] Compute(byte[] data)
        {
            return Compute(data, 0, data.Length);
        }

        /// <summary>
        /// Validates that the two CRC bytes at the end of the frame are
        /// correct for the preceding bytes.
        /// </summary>
        public static bool Validate(byte[] frame, int length)
        {
            if (length < 2) return false;
            var computed = Compute(frame, 0, length - 2);
            return computed[0] == frame[length - 2]
                && computed[1] == frame[length - 1];
        }
    }
}
