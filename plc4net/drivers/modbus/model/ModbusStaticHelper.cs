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
using System.Collections.Generic;
using System.Linq;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.modbus.readwrite.model
{
    /// <summary>
    /// Shared helpers required by the generated protocol model.
    /// </summary>
    public static partial class ModbusStaticHelper
    {
        /// <summary>Serialized byte length of a sequence of messages.</summary>
        public static int ArraySizeInBytes<T>(IEnumerable<T> items) where T : IMessage
            => items?.Sum(i => i.GetLengthInBytes()) ?? 0;

        public static byte BcdToBin(byte value) => (byte) ((value >> 4) * 10 + (value & 0x0F));
        public static ushort BcdToBin12(ushort value) => (ushort) (((value >> 8) & 0x0F) * 100 + ((value >> 4) & 0x0F) * 10 + (value & 0x0F));
        public static byte BinToBcd(int value) => (byte) ((value / 10) * 16 + (value % 10));
        public static ushort BinToBcd12(int value) => (ushort) ((value / 100) * 256 + ((value / 10) % 10) * 16 + (value % 10));

        public static ushort AsciiLrcCheck(byte address, IMessage pdu)
        {
            var buffer = new WriteBuffer();
            buffer.WriteByte("address", 8, address);
            pdu.Serialize(buffer);
            var sum = 0;
            foreach (var b in buffer.GetBytes()) sum = (sum + b) & 0xFF;
            return (ushort) ((-sum) & 0xFF);
        }

        public static ushort RtuCrcCheck(byte address, IMessage pdu)
        {
            var buffer = new WriteBuffer();
            buffer.WriteByte("address", 8, address);
            pdu.Serialize(buffer);
            ushort crc = 0xFFFF;
            foreach (var b in buffer.GetBytes())
            {
                crc ^= b;
                for (var bit = 0; bit < 8; bit++) crc = (crc & 1) != 0 ? (ushort) ((crc >> 1) ^ 0xA001) : (ushort) (crc >> 1);
            }
            return (ushort) ((crc << 8) | (crc >> 8));
        }
    }
}