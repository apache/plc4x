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
using System.Collections.Generic;

namespace org.apache.plc4net.spi.generation
{
    /// <summary>
    /// Writes big-endian (MSB-first) bit fields into a growing byte buffer.
    /// Mirror image of <see cref="BitReader"/>; see that type for why the
    /// third-party bit codec was dropped.
    /// </summary>
    public class BitWriter
    {
        private readonly List<byte> _data = new List<byte>();
        private int _bitPosition;

        /// <summary>Number of bits written so far.</summary>
        public int Length => _bitPosition;

        public void WriteBool(bool value)
        {
            WriteBits(value ? 1UL : 0UL, 1);
        }

        /// <summary>
        /// Writes the low <paramref name="bitCount"/> bits of <paramref name="value"/>,
        /// most significant bit first.
        /// </summary>
        public void WriteBits(ulong value, int bitCount)
        {
            if (bitCount < 0 || bitCount > 64)
            {
                throw new ArgumentOutOfRangeException(
                    nameof(bitCount), bitCount, "Can write between 0 and 64 bits at a time.");
            }

            for (var i = bitCount - 1; i >= 0; i--)
            {
                var bit = (value >> i) & 1;

                if ((_bitPosition & 7) == 0)
                {
                    _data.Add(0);
                }
                if (bit != 0)
                {
                    var byteIndex = _bitPosition >> 3;
                    var bitOffset = 7 - (_bitPosition & 7);
                    _data[byteIndex] |= (byte) (1 << bitOffset);
                }
                _bitPosition++;
            }
        }

        public void WriteBytes(byte[] value)
        {
            if (value == null)
            {
                throw new ArgumentNullException(nameof(value));
            }
            foreach (var b in value)
            {
                WriteBits(b, 8);
            }
        }

        /// <summary>
        /// The bytes written so far. A trailing partial byte is zero padded.
        /// </summary>
        public byte[] GetBytes()
        {
            return _data.ToArray();
        }
    }
}
