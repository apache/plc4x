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
        // announces more than has arrived. The bound also caps the cost of that search
        // on a persistently noisy line. A valid answer that starts further into the
        // buffer is not found and ends as a request timeout.
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
            => HasFrameBehindFirstByte(
                TransportInstance.PeekReadableBytes(Math.Min(availableBytes, LookaheadBytes)));

        /// <summary>
        /// Whether a whole good frame starts anywhere in the buffer after its first byte. On a
        /// noisy line nearly every offset is a miss, so a candidate has to pass two cheap checks
        /// - a function code that can start a frame, and a CRC that fits - before the generated
        /// parser, which throws on every miss, is asked.
        /// </summary>
        internal static bool HasFrameBehindFirstByte(byte[] buffered)
        {
            for (var offset = 1; offset + HeaderSize <= buffered.Length; offset++)
            {
                if (!TryExpectedLength(buffered[offset + 1], buffered[offset + 2], out var length)
                    || offset + length > buffered.Length)
                {
                    continue; // nothing that could be a frame starts here, or not all of it has arrived
                }

                var candidate = buffered.AsSpan(offset, length);
                if (CrcMatches(candidate) && IsFrame(buffered, offset, length))
                {
                    return true;
                }
            }
            return false;
        }

        /// <summary>Whether the last two bytes are the CRC-16/MODBUS of the rest, low byte first.</summary>
        internal static bool CrcMatches(ReadOnlySpan<byte> frame)
        {
            ushort crc = 0xFFFF;
            for (var i = 0; i < frame.Length - 2; i++)
            {
                crc ^= frame[i];
                for (var bit = 0; bit < 8; bit++)
                {
                    crc = (crc & 1) != 0 ? (ushort)((crc >> 1) ^ 0xA001) : (ushort)(crc >> 1);
                }
            }
            return frame[frame.Length - 2] == (byte)crc && frame[frame.Length - 1] == (byte)(crc >> 8);
        }

        /// <summary>Whether the bytes are one whole frame that the generated parser accepts, CRC included.</summary>
        private static bool IsFrame(byte[] bytes, int offset = 0, int length = -1)
        {
            try
            {
                ModbusRtuADU.StaticParse(
                    new ReadBuffer(bytes, offset, length < 0 ? bytes.Length - offset : length),
                    DriverType.MODBUS_RTU,
                    true);
                return true;
            }
            catch (ParseException)
            {
                return false;
            }
        }

        private static int ExpectedLength(byte[] header)
        {
            if (!TryExpectedLength(header[1], header[2], out var length))
            {
                // Not the start of a frame: let the base class resync.
                throw new MessageCodecException(
                    $"Not the start of a Modbus RTU frame (function code 0x{header[1]:X2}, third byte {header[2]}).");
            }
            return length;
        }

        /// <summary>
        /// The length of the frame that starts with an address byte, then this function code and
        /// this third byte - or false when nothing that can be a response starts like that.
        /// </summary>
        private static bool TryExpectedLength(byte functionCode, byte third, out int length)
        {
            // Exception response: address, function, exception code, CRC.
            if ((functionCode & 0x80) != 0)
            {
                length = 5;
                return true;
            }

            switch (functionCode)
            {
                case 0x01: // Read Coils
                case 0x02: // Read Discrete Inputs
                case 0x03: // Read Holding Registers
                case 0x04: // Read Input Registers
                    // Address, function, byte count, data, CRC.
                    length = 3 + third + 2;
                    return third <= MaxByteCount;
                case 0x05: // Write Single Coil
                case 0x06: // Write Single Register
                case 0x0F: // Write Multiple Coils
                case 0x10: // Write Multiple Registers
                    // Address, function, the echoed address and value or quantity, CRC.
                    length = 8;
                    return true;
                default:
                    length = 0;
                    return false;
            }
        }

        protected override ModbusRtuADU ParseMessage(ReadBuffer readBuffer)
            => ModbusRtuADU.StaticParse(readBuffer, DriverType.MODBUS_RTU, true);
    }
}