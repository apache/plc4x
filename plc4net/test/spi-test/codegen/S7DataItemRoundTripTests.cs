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
using System.Globalization;
using System.Linq;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.model.values;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Round-trips the generated S7 <c>DataItem</c> - the <c>[dataIo]</c> that
    /// reads and writes one <see cref="api.value.IPlcValue"/> keyed on the
    /// <c>dataProtocolId</c> parser argument. No shared test-suite vector
    /// exercises it (plc4j drives it straight from <c>S7ProtocolLogic</c>), so
    /// the byte vectors here are built from the IEC 61131 wire layout the mspec
    /// describes: a value parses, exposes the expected value, and serializes
    /// back byte-identical.
    /// </summary>
    /// <remarks>
    /// Every case the emitter generates is covered - the scalar and string
    /// families plus the TIA date / time family (<c>S5TIME</c>, the Siemens
    /// epoch, BCD in <c>DATE_AND_TIME</c>, nanosecond-exact <c>LTIME</c> /
    /// <c>DTL</c>).
    /// </remarks>
    public class S7DataItemRoundTripTests
    {
        private const int NoString = 0;

        public static TheoryData<string, string, string, int> Vectors() => new()
        {
            // name              protocolId              hex                          stringLength
            { "BOOL true",       "IEC61131_BOOL",        "01",                        NoString },
            { "BOOL false",      "IEC61131_BOOL",        "00",                        NoString },
            { "BYTE",            "IEC61131_BYTE",        "a5",                        NoString },
            { "WORD",            "IEC61131_WORD",        "1234",                      NoString },
            { "DWORD",           "IEC61131_DWORD",       "01020304",                  NoString },
            { "LWORD",           "IEC61131_LWORD",       "0102030405060708",          NoString },
            { "SINT negative",   "IEC61131_SINT",        "fe",                        NoString },
            { "USINT",           "IEC61131_USINT",       "a5",                        NoString },
            { "INT negative",    "IEC61131_INT",         "fffe",                      NoString },
            { "UINT",            "IEC61131_UINT",        "1234",                      NoString },
            { "DINT negative",   "IEC61131_DINT",        "fffffffe",                  NoString },
            { "UDINT",           "IEC61131_UDINT",       "01020304",                  NoString },
            { "LINT negative",   "IEC61131_LINT",        "fffffffffffffffe",          NoString },
            { "ULINT",           "IEC61131_ULINT",       "0102030405060708",          NoString },
            { "REAL",            "IEC61131_REAL",        "3fc00000",                  NoString },
            { "REAL negative",   "IEC61131_REAL",        "bfc00000",                  NoString },
            { "REAL NaN",        "IEC61131_REAL",        "7fc00000",                  NoString },
            { "REAL +Inf",       "IEC61131_REAL",        "7f800000",                  NoString },
            { "REAL -Inf",       "IEC61131_REAL",        "ff800000",                  NoString },
            { "LREAL",           "IEC61131_LREAL",       "3ff8000000000000",          NoString },
            { "LREAL NaN",       "IEC61131_LREAL",       "7ff8000000000000",          NoString },
            { "CHAR",            "IEC61131_CHAR",        "41",                        NoString },
            { "WCHAR",           "IEC61131_WCHAR",       "0041",                      NoString },
            { "STRING",          "IEC61131_STRING",      "02024869",                  2 },
            { "WSTRING",         "IEC61131_WSTRING",     "0002000200480069",          2 },
            { "TIME",            "IEC61131_TIME",        "000005dc",                  NoString },
            { "S5TIME",          "S7_S5TIME",            "0123",                      NoString },
            { "LTIME",           "IEC61131_LTIME",       "0000000059682f00",          NoString },
            { "DATE",            "IEC61131_DATE",        "3082",                      NoString },
            { "TIME_OF_DAY",     "IEC61131_TIME_OF_DAY", "02b32c95",                  NoString },
            { "LTIME_OF_DAY",    "IEC61131_LTIME_OF_DAY","0000011f71fb04cb",          NoString },
            { "DATE_AND_TIME",   "IEC61131_DATE_AND_TIME","2401011234567892",         NoString },
            { "DATE_AND_LTIME",  "IEC61131_DATE_AND_LTIME","17a6394954be2d15",        NoString },
            { "DTL",             "IEC61131_DTL",         "07e80101020c2238075bcd15",  NoString },
        };

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Round_trips(string name, string protocolId, string hex, int stringLength)
        {
            _ = name;
            var expected = FromHex(hex);

            var value = DataItem.StaticParse(
                new ReadBuffer(expected), protocolId, ControllerType.ANY, stringLength);
            Assert.NotNull(value);

            var writeBuffer = new WriteBuffer();
            DataItem.StaticSerialize(writeBuffer, value, protocolId, ControllerType.ANY, stringLength);

            Assert.Equal(hex, ToHex(writeBuffer.GetBytes()));
        }

        [Fact]
        public void Parsed_scalar_values_expose_the_expected_value()
        {
            Assert.True(Parse("IEC61131_BOOL", "01").GetBool());
            Assert.False(Parse("IEC61131_BOOL", "00").GetBool());
            Assert.Equal((byte) 0xA5, Parse("IEC61131_BYTE", "a5").GetByte());
            Assert.Equal((ushort) 0x1234, Parse("IEC61131_WORD", "1234").GetUshort());
            Assert.Equal(0x01020304u, Parse("IEC61131_DWORD", "01020304").GetUint());
            Assert.Equal(0x0102030405060708UL, Parse("IEC61131_LWORD", "0102030405060708").GetUlong());
            Assert.Equal((sbyte) (-2), Parse("IEC61131_SINT", "fe").GetSbyte());
            Assert.Equal((byte) 0xA5, Parse("IEC61131_USINT", "a5").GetByte());
            Assert.Equal((short) (-2), Parse("IEC61131_INT", "fffe").GetShort());
            Assert.Equal((ushort) 0x1234, Parse("IEC61131_UINT", "1234").GetUshort());
            Assert.Equal(-2, Parse("IEC61131_DINT", "fffffffe").GetInt());
            Assert.Equal(0x01020304u, Parse("IEC61131_UDINT", "01020304").GetUint());
            Assert.Equal(-2L, Parse("IEC61131_LINT", "fffffffffffffffe").GetLong());
            Assert.Equal(0x0102030405060708UL, Parse("IEC61131_ULINT", "0102030405060708").GetUlong());
            Assert.Equal(1.5f, Parse("IEC61131_REAL", "3fc00000").GetFloat());
            Assert.Equal(-1.5f, Parse("IEC61131_REAL", "bfc00000").GetFloat());
            Assert.Equal(1.5d, Parse("IEC61131_LREAL", "3ff8000000000000").GetDouble());
            Assert.Equal("A", Parse("IEC61131_CHAR", "41").GetString());
            Assert.Equal("A", Parse("IEC61131_WCHAR", "0041").GetString());
            Assert.Equal("Hi", Parse("IEC61131_STRING", "02024869", 2).GetString());
            Assert.Equal("Hi", Parse("IEC61131_WSTRING", "0002000200480069", 2).GetString());
        }

        [Fact]
        public void Non_finite_reals_round_trip_instead_of_throwing()
        {
            // PlcSimpleNumericValueAdapter.IsFloat() used to reject NaN / ±Inf,
            // so GetFloat() threw when the codec re-serialized them.
            Assert.True(float.IsNaN(Parse("IEC61131_REAL", "7fc00000").GetFloat()));
            Assert.True(float.IsPositiveInfinity(Parse("IEC61131_REAL", "7f800000").GetFloat()));
            Assert.True(double.IsNaN(Parse("IEC61131_LREAL", "7ff8000000000000").GetDouble()));
        }

        [Fact]
        public void Parsed_temporal_values_expose_the_expected_value()
        {
            Assert.Equal(TimeSpan.FromMilliseconds(1500), Parse("IEC61131_TIME", "000005dc").GetDuration());
            Assert.Equal(TimeSpan.FromMilliseconds(1230), Parse("S7_S5TIME", "0123").GetDuration());
            Assert.Equal(1500000000UL, ((PlcLTIME) Parse("IEC61131_LTIME", "0000000059682f00")).GetNanoseconds());
            Assert.Equal(new DateOnly(2024, 1, 1), Parse("IEC61131_DATE", "3082").GetDate());
            Assert.Equal(new TimeOnly(12, 34, 56, 789), Parse("IEC61131_TIME_OF_DAY", "02b32c95").GetTime());
            Assert.Equal(1234567890123UL,
                ((PlcLTIME_OF_DAY) Parse("IEC61131_LTIME_OF_DAY", "0000011f71fb04cb")).GetNanosecondsSinceMidnight());

            var dt = Parse("IEC61131_DATE_AND_TIME", "2401011234567892").GetDateTime();
            Assert.Equal(new DateTime(2024, 1, 1, 12, 34, 56, 789), dt);

            // DATE_AND_LTIME (single u64) and DTL are the same instant.
            var dandl = (PlcDATE_AND_LTIME) Parse("IEC61131_DATE_AND_LTIME", "17a6394954be2d15");
            Assert.Equal(1704112496123456789UL, dandl.GetNanosecondsSinceEpoch());
            Assert.Equal(123456789u, dandl.GetNanosecondsOfSecond());

            var dtl = (PlcDATE_AND_LTIME) Parse("IEC61131_DTL", "07e80101020c2238075bcd15");
            Assert.Equal(123456789u, dtl.GetNanosecondsOfSecond());
            Assert.Equal(1704112496123456789UL, dtl.GetNanosecondsSinceEpoch());
            Assert.Equal(2024, dtl.GetDateTime().Year);
            Assert.Equal(56, dtl.GetDateTime().Second);
            // Both construction paths must agree on Kind (they used to differ:
            // UnixEpoch is Utc, `new DateTime(y,m,d,…)` is Unspecified).
            Assert.Equal(dandl.GetDateTime().Kind, dtl.GetDateTime().Kind);
            Assert.Equal(DateTimeKind.Unspecified, dtl.GetDateTime().Kind);
        }

        [Fact]
        public void Temporal_values_report_as_simple()
        {
            // PlcDATE / PlcDATE_AND_TIME / PlcDATE_AND_LTIME extended
            // PlcValueAdapter, not PlcSimpleValueAdapter, so IsSimple() answered
            // false and a consumer branching on it treated a date as a struct.
            Assert.True(Parse("IEC61131_DATE", "3082").IsSimple());
            Assert.True(Parse("IEC61131_DATE_AND_TIME", "2401011234567892").IsSimple());
            Assert.True(Parse("IEC61131_DTL", "07e80101020c2238075bcd15").IsSimple());
            Assert.True(Parse("IEC61131_TIME", "000005dc").IsSimple());
            Assert.True(Parse("IEC61131_LTIME_OF_DAY", "0000011f71fb04cb").IsSimple());
        }

        [Fact]
        public void S5time_serializes_in_canonical_form()
        {
            // S5TIME is a non-unique encoding: 1000 ms is 0x1010 (base 100 ms)
            // or 0x0100 (base 10 ms). The codec always emits the smallest base,
            // matching plc4j's durationToS5Time, so a non-canonical frame does
            // NOT round-trip byte-identically.
            var value = Parse("S7_S5TIME", "1010");
            Assert.Equal(TimeSpan.FromMilliseconds(1000), value.GetDuration());

            var writeBuffer = new WriteBuffer();
            DataItem.StaticSerialize(writeBuffer, value, "S7_S5TIME", ControllerType.ANY, NoString);
            Assert.Equal("0100", ToHex(writeBuffer.GetBytes()));
        }

        [Fact]
        public void Bit_string_values_survive_the_round_trip_through_the_value_model()
        {
            // PlcBitString used to discard its constructor argument; a WORD that
            // could not read its own value back would serialize as zero.
            Assert.Equal("1234", RoundTrip("IEC61131_WORD", "1234"));
            Assert.Equal("c0debabe", RoundTrip("IEC61131_DWORD", "c0debabe"));
            Assert.Equal("0102030405060708", RoundTrip("IEC61131_LWORD", "0102030405060708"));
        }

        // ── helpers ─────────────────────────────────────────────

        private static api.value.IPlcValue Parse(string protocolId, string hex, int stringLength = NoString) =>
            DataItem.StaticParse(new ReadBuffer(FromHex(hex)), protocolId, ControllerType.ANY, stringLength);

        private static string RoundTrip(string protocolId, string hex, int stringLength = NoString)
        {
            var value = Parse(protocolId, hex, stringLength);
            var writeBuffer = new WriteBuffer();
            DataItem.StaticSerialize(writeBuffer, value, protocolId, ControllerType.ANY, stringLength);
            return ToHex(writeBuffer.GetBytes());
        }

        private static byte[] FromHex(string hex) =>
            Enumerable.Range(0, hex.Length / 2)
                .Select(i => byte.Parse(hex.Substring(i * 2, 2), NumberStyles.HexNumber, CultureInfo.InvariantCulture))
                .ToArray();

        private static string ToHex(byte[] bytes) =>
            string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}
