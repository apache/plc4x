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
    /// Frames Modbus RTU responses. An RTU frame carries no length field, so the
    /// length is derived from the function code (and, for reads, the byte count).
    /// Once that many bytes are there, the frame is accepted only if the generated
    /// <see cref="ModbusRtuADU"/> parses it, which includes the CRC; otherwise one byte is
    /// dropped and the search goes on, so noise in front of a frame, a corrupted frame, or the
    /// echo of the request that a half-duplex adapter sends back cannot hide the good frame
    /// behind it.
    /// </summary>
    internal sealed class ModbusRtuMessageCodec : MessageCodecBase<ModbusRtuADU>
    {
        // Address, function code, and the byte that is a byte count for reads
        // and an exception code for errors.
        private const int HeaderSize = 3;

        // A read response carries at most 250 bytes of data (125 registers), so
        // the whole frame stays within the 256-byte RTU limit.
        private const int MaxByteCount = 250;

        // How much of the buffer is searched for a good frame behind a header that
        // announces more than has arrived. The bound also caps the resync cost on a
        // persistently noisy line, where every offset re-parses a candidate frame. A
        // valid answer that starts further into the buffer is not found and ends as
        // a request timeout.
        private const int LookaheadBytes = 512;

        public ModbusRtuMessageCodec(ITransportInstance transport, Action<ModbusRtuADU> handler)
            : base("Modbus RTU", transport, handler)
        {
        }

        protected override int GetMinimumHeaderSize() => HeaderSize;

        protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
        {
            var total = ExpectedLength(header);
            if (availableBytes < total)
            {
                // The frame this header announces has not arrived in full, so wait for it -
                // unless a whole good frame is already buffered behind it. The header was
                // noise then, and waiting would keep that frame from being looked at (noise
                // that happens to announce a long read response would hold up the real
                // answer until the timeout).
                if (HasCompleteFrameBehind(availableBytes))
                {
                    throw new MessageCodecException("Noise in front of a frame.");
                }
                return total;
            }

            if (!IsFrame(TransportInstance.PeekReadableBytes(total)))
            {
                // Not a frame after all (a failed CRC, most likely): the base class
                // drops a byte and tries again from the next one.
                throw new MessageCodecException("Not a valid Modbus RTU frame.");
            }
            return total;
        }

        private bool HasCompleteFrameBehind(int availableBytes)
        {
            var buffered = TransportInstance.PeekReadableBytes(Math.Min(availableBytes, LookaheadBytes));
            for (var offset = 1; offset + HeaderSize <= buffered.Length; offset++)
            {
                int length;
                try
                {
                    length = ExpectedLength(new[] { buffered[offset], buffered[offset + 1], buffered[offset + 2] });
                }
                catch (MessageCodecException)
                {
                    continue; // nothing that could be a frame starts here
                }

                if (offset + length <= buffered.Length
                    && IsFrame(buffered[offset..(offset + length)]))
                {
                    return true;
                }
            }
            return false;
        }

        /// <summary>Whether the bytes are one whole frame that the generated parser accepts, CRC included.</summary>
        private static bool IsFrame(byte[] bytes)
        {
            try
            {
                ModbusRtuADU.StaticParse(new ReadBuffer(bytes), DriverType.MODBUS_RTU, true);
                return true;
            }
            catch (ParseException)
            {
                return false;
            }
        }

        private static int ExpectedLength(byte[] header)
        {
            var functionCode = header[1];

            // Exception response: address, function, exception code, CRC.
            if ((functionCode & 0x80) != 0)
            {
                return 5;
            }

            switch (functionCode)
            {
                case 0x01: // Read Coils
                case 0x02: // Read Discrete Inputs
                case 0x03: // Read Holding Registers
                case 0x04: // Read Input Registers
                    if (header[2] > MaxByteCount)
                    {
                        throw new MessageCodecException(
                            $"Implausible Modbus RTU byte count {header[2]}.");
                    }
                    // Address, function, byte count, data, CRC.
                    return 3 + header[2] + 2;
                case 0x05: // Write Single Coil
                case 0x06: // Write Single Register
                case 0x0F: // Write Multiple Coils
                case 0x10: // Write Multiple Registers
                    // Address, function, the echoed address and value or quantity, CRC.
                    return 8;
                default:
                    // Not the start of a frame: let the base class resync.
                    throw new MessageCodecException(
                        $"Unsupported Modbus RTU function code 0x{functionCode:X2}.");
            }
        }

        protected override ModbusRtuADU ParseMessage(ReadBuffer readBuffer)
            => ModbusRtuADU.StaticParse(readBuffer, DriverType.MODBUS_RTU, true);
    }
}