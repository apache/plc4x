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
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.spi.generation;
using Xunit;

namespace org.apache.plc4net.spi.test.codegen
{
    /// <summary>
    /// Round-trips the generated Modbus <c>DataItem</c>. The mspec redesign
    /// (2026) reduced it to exactly one value at its natural width - register
    /// alignment, arrays and byte order are the driver's job now - so there is
    /// one case per <c>ModbusDataType</c> and the only argument that carries
    /// data is <c>stringLength</c> (for the two <c>vstring</c> cases).
    /// </summary>
    public class ModbusDataItemRoundTripTests
    {
        public static TheoryData<string, string, string, ushort> Vectors() => new()
        {
            // name             dataType   hex                    stringLength
            { "BOOL true",       "BOOL",    "80",                  0 },
            { "BOOL false",      "BOOL",    "00",                  0 },
            { "BYTE",            "BYTE",    "a5",                  0 },
            { "WORD",            "WORD",    "1234",                0 },
            { "DWORD",           "DWORD",   "cafebabe",            0 },
            { "LWORD",           "LWORD",   "0011223344556677",    0 },
            { "SINT negative",   "SINT",    "fe",                  0 },
            { "INT negative",    "INT",     "fffe",                0 },
            { "DINT negative",   "DINT",    "fffffffe",            0 },
            { "LINT negative",   "LINT",    "fffffffffffffffe",    0 },
            { "USINT",           "USINT",   "a5",                  0 },
            { "UINT",            "UINT",    "1234",                0 },
            { "UDINT",           "UDINT",   "01020304",            0 },
            { "ULINT",           "ULINT",   "0102030405060708",    0 },
            { "REAL",            "REAL",    "3fc00000",            0 },
            { "LREAL",           "LREAL",   "3ff8000000000000",    0 },
            { "CHAR",            "CHAR",    "41",                  0 },
            { "WCHAR",           "WCHAR",   "0041",                0 },
            { "STRING",          "STRING",  "48656c6c6f",          5 },
            { "WSTRING",         "WSTRING", "00480049",            2 },
        };

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Round_trips(string name, string dataType, string hex, ushort stringLength)
        {
            _ = name;
            var type = Enum.Parse<ModbusDataType>(dataType);
            var expected = FromHex(hex);

            var value = DataItem.StaticParse(new ReadBuffer(expected), type, stringLength);
            Assert.NotNull(value);

            var writeBuffer = new WriteBuffer();
            DataItem.StaticSerialize(writeBuffer, value, type, stringLength);

            Assert.Equal(hex, ToHex(writeBuffer.GetBytes()));
        }

        [Fact]
        public void Scalars_expose_the_expected_value()
        {
            Assert.True(DataItem.StaticParse(new ReadBuffer(FromHex("80")), ModbusDataType.BOOL, 0).GetBool());
            Assert.Equal((ushort) 0x1234,
                DataItem.StaticParse(new ReadBuffer(FromHex("1234")), ModbusDataType.WORD, 0).GetUshort());
            Assert.Equal((short) (-2),
                DataItem.StaticParse(new ReadBuffer(FromHex("fffe")), ModbusDataType.INT, 0).GetShort());
            Assert.Equal(1.5f,
                DataItem.StaticParse(new ReadBuffer(FromHex("3fc00000")), ModbusDataType.REAL, 0).GetFloat());
            Assert.Equal("Hello",
                DataItem.StaticParse(new ReadBuffer(FromHex("48656c6c6f")), ModbusDataType.STRING, 5).GetString());
        }

        // ── helpers ─────────────────────────────────────────────

        private static byte[] FromHex(string hex) =>
            Enumerable.Range(0, hex.Length / 2)
                .Select(i => byte.Parse(hex.Substring(i * 2, 2), NumberStyles.HexNumber, CultureInfo.InvariantCulture))
                .ToArray();

        private static string ToHex(byte[] bytes) =>
            string.Concat(bytes.Select(b => b.ToString("x2", CultureInfo.InvariantCulture)));
    }
}
