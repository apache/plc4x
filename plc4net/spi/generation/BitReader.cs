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

namespace org.apache.plc4net.spi.generation
{
    /// <summary>
    /// Reads big-endian (MSB-first) bit fields out of a byte array.
    /// </summary>
    /// <remarks>
    /// This replaces the third-party Ayx.BitIO package, which is a .NET Framework
    /// only assembly, unmaintained since 2016, and which threw
    /// ArgumentOutOfRangeException when asked for a 32 bit field because it
    /// round-tripped every read through a StringBuilder of '0'/'1' characters.
    ///
    /// MSB-first is the order every mspec-described protocol uses: the first bit of
    /// a frame is the high bit of byte 0.
    /// </remarks>
    public class BitReader
    {
        private readonly byte[] _data;
        private int _bitPosition;

        public BitReader(byte[] data)
        {
            _data = data ?? throw new ArgumentNullException(nameof(data));
        }

        /// <summary>Current read offset, in bits.</summary>
        public int Position
        {
            get => _bitPosition;
            set
            {
                if (value < 0 || value > TotalBits)
                {
                    throw new ArgumentOutOfRangeException(nameof(value));
                }
                _bitPosition = value;
            }
        }

        public int TotalBits => _data.Length * 8;

        /// <summary>Number of bits not yet consumed.</summary>
        public int Remaining => TotalBits - _bitPosition;

        public bool ReadBool()
        {
            return ReadBits(1) != 0;
        }

        /// <summary>
        /// Reads <paramref name="bitCount"/> bits and returns them right-aligned.
        /// </summary>
        public ulong ReadBits(int bitCount)
        {
            if (bitCount < 0 || bitCount > 64)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitCount), bitCount, "Can read between 0 and 64 bits at a time.");
            }
            if (bitCount > Remaining)
            {
                throw new ParseException(
                    $"Not enough data: tried to read {bitCount} bits, {Remaining} available.");
            }

            ulong result = 0;
            for (var i = 0; i < bitCount; i++)
            {
                var byteIndex = _bitPosition >> 3;
                var bitOffset = 7 - (_bitPosition & 7); // MSB-first within the byte
                var bit = (ulong) ((_data[byteIndex] >> bitOffset) & 1);
                result = (result << 1) | bit;
                _bitPosition++;
            }
            return result;
        }

        /// <summary>
        /// Reads <paramref name="bitCount"/> bits and sign-extends them from that width.
        /// </summary>
        public long ReadSignedBits(int bitCount)
        {
            var raw = ReadBits(bitCount);
            if (bitCount == 0 || bitCount == 64)
            {
                return unchecked((long) raw);
            }
            var signMask = 1UL << (bitCount - 1);
            if ((raw & signMask) != 0)
            {
                // Set every bit above the field width so the value stays negative.
                raw |= ulong.MaxValue << bitCount;
            }
            return unchecked((long) raw);
        }

        /// <summary>Reads whole bytes. Requires the reader to be byte aligned.</summary>
        public byte[] ReadBytes(int byteCount)
        {
            if (byteCount < 0)
            {
                throw new ArgumentOutOfRangeException(nameof(byteCount));
            }
            if (byteCount * 8 > Remaining)
            {
                throw new ParseException(
                    $"Not enough data: tried to read {byteCount} bytes, {Remaining / 8} available.");
            }

            var result = new byte[byteCount];
            if ((_bitPosition & 7) == 0)
            {
                Array.Copy(_data, _bitPosition >> 3, result, 0, byteCount);
                _bitPosition += byteCount * 8;
            }
            else
            {
                for (var i = 0; i < byteCount; i++)
                {
                    result[i] = (byte) ReadBits(8);
                }
            }
            return result;
        }
    }
}
