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
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.drivers.modbus
{
    /// <summary>
    /// Frames Modbus TCP responses: the 6-byte MBAP header up to and including the
    /// length field says how many bytes follow, and the generated
    /// <see cref="ModbusTcpADU"/> parses the whole frame.
    /// </summary>
    internal sealed class ModbusTcpMessageCodec : MessageCodecBase<ModbusTcpADU>
    {
        // Transaction identifier (2) + protocol identifier (2) + length (2).
        private const int HeaderSize = 6;

        // The length field counts the unit identifier and the PDU. A PDU is at
        // most 253 bytes, and the smallest valid response is an exception
        // response: unit identifier, function code, exception code.
        private const int MinLength = 3;
        private const int MaxLength = 254;

        public ModbusTcpMessageCodec(ITransportInstance transport, Action<ModbusTcpADU> handler)
            : base("Modbus TCP", transport, handler)
        {
        }

        protected override int GetMinimumHeaderSize() => HeaderSize;

        protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
        {
            var protocolIdentifier = (header[2] << 8) | header[3];
            var length = (header[4] << 8) | header[5];
            if (protocolIdentifier != 0 || length < MinLength || length > MaxLength)
            {
                // Not the start of a frame: let the base class resync.
                throw new MessageCodecException(
                    $"Not a Modbus TCP header (protocol identifier {protocolIdentifier}, length {length}).");
            }
            return HeaderSize + length;
        }

        protected override ModbusTcpADU ParseMessage(ReadBuffer readBuffer)
            => ModbusTcpADU.StaticParse(readBuffer, DriverType.MODBUS_TCP, true);
    }
}