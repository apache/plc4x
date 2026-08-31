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
    /// Round-trips the generated Modbus <c>DataItem</c>. Modbus keys the same
    /// <c>[dataIo]</c> emitter on three discriminators - an <c>enum</c>
    /// (<c>dataType</c>), a count (<c>numberOfValues</c>) and a byte-order bit
    /// (<c>bigEndian</c>) - so the cases compile to an <c>else if</c> chain
    /// rather than S7's flat string switch.
    /// </summary>
    /// <remarks>
    /// The single-value scalar cases round-trip; the <c>numberOfValues &gt; 1</c>
    /// list cases (a <c>PlcList</c> of primitives) are a throwing stub, the same
    /// GAP-8 boundary the S7 date / time cases sit behind.
    /// </remarks>
    public class ModbusDataItemRoundTripTests
    {
        public static TheoryData<string, string, string, ushort, bool> Vectors() => new()
        {
            // name              dataType  hex               numberOfValues  bigEndian
            { "BOOL big",         "BOOL",   "0001",           1,              true },
            { "BOOL little",      "BOOL",   "0100",           1,              false },
            { "BYTE",             "BYTE",   "00a5",           1,              true },
            { "WORD",             "WORD",   "1234",           1,              true },
            { "DWORD",            "DWORD",  "cafebabe",        1,              true },
            { "INT negative",     "INT",    "fffe",           1,              true },
            { "DINT negative",    "DINT",   "fffffffe",       1,              true },
            { "UDINT",            "UDINT",  "01020304",       1,              true },
            { "REAL",             "REAL",   "3fc00000",       1,              true },
            { "LREAL",            "LREAL",  "3ff8000000000000", 1,            true },
            { "CHAR",             "CHAR",   "41",             1,              true },
            { "WCHAR",            "WCHAR",  "0041",           1,              true },
            // numberOfValues != 1 -> a PlcList, one element per count (the
            // `…,'1',…` scalar cases take precedence, so these need count != 1).
            { "BOOL x3",          "BOOL",   "a0",             3,              true },
            { "BYTE x2 (16 bits)","BYTE",   "a5c3",           2,              true },
            { "INT x3",           "INT",    "0001fffe0102",   3,              true },
            { "UINT x2",          "UINT",   "12345678",       2,              true },
            { "DINT x2",          "DINT",   "00000001fffffffe", 2,            true },
            { "REAL x2",          "REAL",   "3fc00000bfc00000", 2,            true },
            { "CHAR x2",          "CHAR",   "4142",           2,              true },
            { "WCHAR x2",         "WCHAR",  "00410042",       2,              true },
        };

        [Theory]
        [MemberData(nameof(Vectors))]
        public void Round_trips(string name, string dataType, string hex, ushort numberOfValues, bool bigEndian)
        {
            _ = name;
            var type = Enum.Parse<ModbusDataType>(dataType);
            var expected = FromHex(hex);

            var value = DataItem.StaticParse(new ReadBuffer(expected), type, numberOfValues, bigEndian);
            Assert.NotNull(value);

            var writeBuffer = new WriteBuffer();
            DataItem.StaticSerialize(writeBuffer, value, type, numberOfValues, bigEndian);

            Assert.Equal(hex, ToHex(writeBuffer.GetBytes()));
        }

        [Fact]
        public void Scalars_expose_the_expected_value()
        {
            Assert.True(DataItem.StaticParse(new ReadBuffer(FromHex("0001")), ModbusDataType.BOOL, 1, true).GetBool());
            Assert.Equal((ushort) 0x1234,
                DataItem.StaticParse(new ReadBuffer(FromHex("1234")), ModbusDataType.WORD, 1, true).GetUshort());
            Assert.Equal((short) (-2),
                DataItem.StaticParse(new ReadBuffer(FromHex("fffe")), ModbusDataType.INT, 1, true).GetShort());
            Assert.Equal(1.5f,
                DataItem.StaticParse(new ReadBuffer(FromHex("3fc00000")), ModbusDataType.REAL, 1, true).GetFloat());
        }

        [Fact]
        public void Multi_value_lists_parse_to_a_PlcList_of_wrapped_primitives()
        {
            var ints = DataItem.StaticParse(new ReadBuffer(FromHex("0001fffe0102")), ModbusDataType.INT, 3, true);
            Assert.True(ints.IsList());
            Assert.Equal(3, ints.GetLength());
            Assert.Equal(new[] { (short) 1, (short) -2, (short) 258 },
                ints.GetList().Select(v => v.GetShort()));

            var bools = DataItem.StaticParse(new ReadBuffer(FromHex("a0")), ModbusDataType.BOOL, 3, true);
            Assert.Equal(new[] { true, false, true }, bools.GetList().Select(v => v.GetBool()));

            var chars = DataItem.StaticParse(new ReadBuffer(FromHex("4142")), ModbusDataType.CHAR, 2, true);
            Assert.Equal("AB", string.Concat(chars.GetList().Select(v => v.GetString())));
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
