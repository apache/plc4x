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
using System.Text;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.generation;

namespace org.apache.plc4net.drivers.s7.readwrite.model
{
    /// <summary>
    /// Hand-written bodies for the <c>STATIC_CALL</c> targets in s7.mspec.
    /// The generated <see cref="S7StaticHelper"/> declares <c>params object[]</c>
    /// stubs that throw; the typed overloads here are the ones the generated
    /// code actually binds to. Ported from
    /// <c>plc4j/.../s7/readwrite/utils/StaticHelper.java</c>.
    /// </summary>
    public static partial class S7StaticHelper
    {
        // ── S7 strings (S7STRING / S7WSTRING) ───────────────────
        //
        // Wire layout: [max length : 1][current length : 1][chars…], the
        // char width and the two length prefixes doubling for UTF16BE.

        public static string ParseS7String(ReadBuffer io, int stringLength, string encoding)
        {
            if (IsUtf8(encoding))
            {
                io.ReadByte("maxLength", 8);
                int totalStringLength = io.ReadByte("totalLength", 8);

                var bytes = new byte[Math.Max(totalStringLength, 0)];
                for (var i = 0; i < stringLength && io.HasMore(8); i++)
                {
                    var cur = (byte) io.ReadSbyte("char", 8);
                    if (i < totalStringLength)
                    {
                        bytes[i] = cur;
                    }
                    else
                    {
                        // consume the rest of the declared field, keep nothing
                        for (i++; i < stringLength && io.HasMore(8); i++)
                        {
                            io.ReadSbyte("pad", 8);
                        }

                        break;
                    }
                }

                return Encoding.UTF8.GetString(bytes);
            }

            if (IsUtf16Be(encoding))
            {
                io.ReadUshort("maxLength", 16);
                int totalStringLength = io.ReadUshort("totalLength", 16);

                var bytes = new byte[Math.Max(totalStringLength, 0) * 2];
                for (var i = 0; i < stringLength && io.HasMore(16); i++)
                {
                    var cur = (ushort) io.ReadShort("char", 16);
                    if (i < totalStringLength)
                    {
                        bytes[i * 2] = (byte) (cur >> 8);
                        bytes[(i * 2) + 1] = (byte) (cur & 0xFF);
                    }
                    else
                    {
                        for (i++; i < stringLength && io.HasMore(16); i++)
                        {
                            io.ReadShort("pad", 16);
                        }

                        break;
                    }
                }

                return Encoding.BigEndianUnicode.GetString(bytes);
            }

            throw new NotSupportedException($"Unsupported S7 string encoding '{encoding}'");
        }

        public static void SerializeS7String(WriteBuffer io, IPlcValue value, int stringLength, string encoding)
        {
            var text = value?.GetString() ?? string.Empty;

            if (IsUtf8(encoding))
            {
                var maxLength = Math.Min(stringLength, 254) & 0xFF;
                var actLength = Math.Min(maxLength, text.Length & 0xFF);

                var chars = new byte[maxLength];
                var actChars = Encoding.UTF8.GetBytes(text.Substring(0, actLength));
                Array.Copy(actChars, chars, Math.Min(actChars.Length, chars.Length));

                io.WriteByte("maxLength", 8, (byte) maxLength);
                io.WriteByte("totalLength", 8, (byte) actLength);
                io.WriteByteArray("chars", chars);
                return;
            }

            if (IsUtf16Be(encoding))
            {
                var maxLength = Math.Min(stringLength, 16382) & 0xFFFF;
                var actLength = Math.Min(maxLength, text.Length & 0xFFFF);

                var chars = new byte[maxLength * 2];
                var actChars = Encoding.BigEndianUnicode.GetBytes(text.Substring(0, actLength));
                Array.Copy(actChars, chars, Math.Min(actChars.Length, chars.Length));

                io.WriteUshort("maxLength", 16, (ushort) maxLength);
                io.WriteUshort("totalLength", 16, (ushort) actLength);
                io.WriteByteArray("chars", chars);
                return;
            }

            throw new NotSupportedException($"Unsupported S7 string encoding '{encoding}'");
        }

        private static bool IsUtf8(string encoding) =>
            string.Equals(encoding, "UTF8", StringComparison.OrdinalIgnoreCase);

        private static bool IsUtf16Be(string encoding) =>
            string.Equals(encoding, "UTF16BE", StringComparison.OrdinalIgnoreCase);

        // ── S7 "associated value" length coding ─────────────────

        /// <summary>The 16-bit length prefix of an event's associated value:
        /// a bit / real / octet-string length is a byte count, everything else
        /// is bit-shifted by 3.</summary>
        public static int RightShift3(ReadBuffer buffer, DataTransportSize transportSize)
        {
            int raw = buffer.ReadUshort("valueLength", 16);
            return transportSize is DataTransportSize.OCTET_STRING
                or DataTransportSize.REAL
                or DataTransportSize.BIT
                ? raw
                : raw >> 3;
        }

        public static void LeftShift3(WriteBuffer buffer, int valueLength) =>
            buffer.WriteUshort("valueLength", 16, (ushort) (valueLength << 3));

        /// <summary>An associated value is padded to an even byte count unless
        /// it is the last item in the buffer.</summary>
        public static int EventItemLength(ReadBuffer buffer, int valueLength) =>
            valueLength % 2 == 0 || !buffer.HasMore((valueLength + 1) * 8)
                ? valueLength
                : valueLength + 1;

        // ── TIA date / time codings ─────────────────────────────

        // S7-300/400 dates count from 1990-01-01, the rest of PLC4X from
        // 1970-01-01. 7305 days between them.
        private const int DaysBetweenUnixAndSiemensEpoch = 7305;

        /// <summary>S5TIME: a BCD mantissa (3 digits) and a 2-bit time base
        /// (10 ms .. 10 s) packed into 16 bits. Returns milliseconds.</summary>
        public static uint ParseS5Time(ReadBuffer io)
        {
            int t = io.ReadUshort("s5time", 16);
            long tv = (t & 0x000F) + ((t & 0x00F0) >> 4) * 10 + ((t & 0x0F00) >> 8) * 100;
            // S5TIME only defines time bases 0..3 (10 ms .. 10 s); clamp the
            // 4-bit field so a malformed frame can't overflow the multiply.
            var exponent = System.Math.Min((t & 0xF000) >> 12, 3);
            long tb = 10;
            for (var i = 0; i < exponent; i++)
            {
                tb *= 10;
            }

            var totalMs = tv * tb;
            return (uint) (totalMs <= 9990000 ? totalMs : 9990000);
        }

        public static void SerializeS5Time(WriteBuffer io, IPlcValue value)
        {
            var totalMs = (long) value.GetDuration().TotalMilliseconds;
            ushort s5time = 0;
            if (totalMs >= 0 && totalMs <= 9990000)
            {
                int tb, tv;
                if (totalMs <= 9990) { tb = 0; tv = (int) (totalMs / 10); }
                else if (totalMs <= 99900) { tb = 1; tv = (int) (totalMs / 100); }
                else if (totalMs <= 999000) { tb = 2; tv = (int) (totalMs / 1000); }
                else { tb = 3; tv = (int) (totalMs / 10000); }

                s5time = (ushort) (
                    (tb << 12) + ((tv / 100 % 10) << 8) + ((tv / 10 % 10) << 4) + (tv % 10));
            }

            io.WriteUshort("s5time", 16, s5time);
        }

        public static int ParseTiaDate(ReadBuffer io) =>
            io.ReadUshort("date", 16) + DaysBetweenUnixAndSiemensEpoch;

        public static void SerializeTiaDate(WriteBuffer io, IPlcValue value) =>
            io.WriteUshort("date", 16,
                (ushort) (DaysSinceUnixEpoch(value) - DaysBetweenUnixAndSiemensEpoch));

        private static int DaysSinceUnixEpoch(IPlcValue value)
        {
            var date = value.IsDate() ? value.GetDate() : DateOnly.FromDateTime(value.GetDateTime());
            return date.DayNumber - new DateOnly(1970, 1, 1).DayNumber;
        }

        /// <summary>DATE_AND_TIME stores the year as one BCD byte: 90-99 -&gt;
        /// 1990-1999, 00-89 -&gt; 2000-2089.</summary>
        public static int ParseSiemensYear(ReadBuffer io)
        {
            var y = BcdToBin(io.ReadByte("year", 8));
            return y < 90 ? 2000 + y : 1900 + y;
        }

        public static void SerializeSiemensYear(WriteBuffer io, IPlcValue value)
        {
            var year = value.GetDateTime().Year;
            io.WriteByte("year", 8, BinToBcd(year >= 2000 ? year - 2000 : year - 1900));
        }

        /// <summary>One binary-coded-decimal byte -&gt; its value (0x25 -&gt; 25).</summary>
        public static int BcdToBin(byte bcd) => ((bcd >> 4) & 0x0F) * 10 + (bcd & 0x0F);

        public static byte BinToBcd(int value) => (byte) ((((value / 10) % 10) << 4) | (value % 10));

        /// <summary>Three BCD nibbles (12 bits) -&gt; 0-999.</summary>
        public static int BcdToBin12(ushort bcd) =>
            ((bcd >> 8) & 0x0F) * 100 + ((bcd >> 4) & 0x0F) * 10 + (bcd & 0x0F);

        public static ushort BinToBcd12(int value) => (ushort) (
            (((value / 100) % 10) << 8) | (((value / 10) % 10) << 4) | (value % 10));

        // ── misc ────────────────────────────────────────────────

        public static int ArrayLength(byte[] array) => array?.Length ?? 0;
    }
}
