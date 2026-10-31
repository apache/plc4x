/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

namespace org.apache.plc4net.spi.generation
{
    /// <summary>
    /// Wire byte order for a multi-byte simple-type field. <c>BIG_ENDIAN</c> is
    /// the implicit default every buffer starts with. A type that declares
    /// <c>byteOrder='"LITTLE_ENDIAN"'</c> has its generated <c>StaticParse</c> /
    /// <c>Serialize</c> push <see cref="LITTLE_ENDIAN"/> before reading or
    /// writing its own fields and pop it again once done (mirrors the Java
    /// generator's <c>pushContext</c> / <c>popContext</c>): every nested type
    /// parsed or serialized through the same buffer in between - including one
    /// that declares no byte order of its own - inherits it.
    /// </summary>
    public enum ByteOrder
    {
        BIG_ENDIAN,
        LITTLE_ENDIAN
    }

    /// <summary>
    /// The byte-swap <see cref="ReadBuffer"/> / <see cref="WriteBuffer"/> apply
    /// to a multi-byte integer field when the current <see cref="ByteOrder"/> is
    /// <see cref="ByteOrder.LITTLE_ENDIAN"/>.
    /// </summary>
    internal static class ByteOrderSupport
    {
        /// <summary>
        /// Reverses the low <c>bitLength / 8</c> bytes of <paramref name="value"/>
        /// - e.g. a 24 bit field reverses exactly 3 bytes, leaving the unused
        /// high byte of the <see cref="ulong"/> at 0. A single byte
        /// (<paramref name="bitLength"/> &lt;= 8) or a field that is not a whole
        /// number of bytes has no defined byte order and is returned unchanged.
        /// </summary>
        internal static ulong Swap(ulong value, int bitLength, ByteOrder order)
        {
            if (order != ByteOrder.LITTLE_ENDIAN || bitLength <= 8 || bitLength % 8 != 0)
            {
                return value;
            }
            var byteCount = bitLength / 8;
            var result = 0UL;
            for (var i = 0; i < byteCount; i++)
            {
                result = (result << 8) | ((value >> (i * 8)) & 0xFF);
            }
            return result;
        }

        /// <summary>
        /// Sign-extends the low <paramref name="bitLength"/> bits of
        /// <paramref name="value"/> to a full 64 bit two's-complement
        /// <see cref="long"/>, the way <c>BitReader.ReadSignedBits</c> does for
        /// the big-endian path. Used after <see cref="Swap"/>, which operates on
        /// the raw unsigned bit pattern and leaves sign extension to the caller.
        /// </summary>
        internal static long SignExtend(ulong value, int bitLength)
        {
            if (bitLength >= 64)
            {
                return unchecked((long)value);
            }
            var signBit = 1UL << (bitLength - 1);
            return (value & signBit) != 0
                ? unchecked((long)(value | ~(signBit * 2 - 1)))
                : (long)value;
        }

        /// <summary>The bit mask covering the low <paramref name="bitLength"/> bits.</summary>
        internal static ulong Mask(int bitLength) =>
            bitLength >= 64 ? ulong.MaxValue : (1UL << bitLength) - 1;
    }
}