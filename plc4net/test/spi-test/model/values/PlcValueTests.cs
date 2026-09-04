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
using System.Linq;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.model.values;
using Xunit;

namespace org.apache.plc4net.spi.test.model.values
{
    /// <summary>
    /// Every assertion here goes through <see cref="IPlcValue"/> rather than through
    /// the concrete type.
    /// </summary>
    /// <remarks>
    /// That distinction is the whole point of these tests. The value classes used to
    /// declare their accessors with `new` instead of `override`, while the only class
    /// implementing IPlcValue (PlcValueAdapter) returned `default` from everything.
    /// Calling through the concrete type therefore looked correct, and calling through
    /// the interface -- which is how every consumer of the API reaches a value --
    /// silently returned 0/false/null.
    /// </remarks>
    public class PlcValueTests
    {
        [Fact]
        public void Dint_read_through_interface_returns_the_value()
        {
            IPlcValue value = new PlcDINT(42);

            Assert.Equal(42, value.GetInt());
            Assert.Equal(42L, value.GetLong());
            Assert.Equal("42", value.GetString());
        }

        [Fact]
        public void Numeric_value_narrows_to_byte_when_in_range()
        {
            // Previously threw InvalidCastException: the adapter unboxed the value with
            // a direct `(byte) value` cast, which is invalid on a boxed Int32.
            IPlcValue value = new PlcDINT(42);

            Assert.True(value.IsByte());
            Assert.Equal((byte) 42, value.GetByte());
        }

        [Fact]
        public void Numeric_value_out_of_byte_range_is_rejected()
        {
            IPlcValue value = new PlcDINT(300);

            Assert.False(value.IsByte());
            Assert.Throws<ArgumentOutOfRangeException>(() => value.GetByte());
        }

        [Fact]
        public void Negative_value_is_not_an_unsigned_long()
        {
            // Previously threw ArgumentException, because IComparable.CompareTo cannot
            // compare a boxed Int64 against a UInt64 bound.
            IPlcValue value = new PlcLINT(-1L);

            Assert.False(value.IsUlong());
        }

        [Theory]
        [InlineData(0, false)]
        [InlineData(1, true)]
        [InlineData(-7, true)]
        public void Numeric_value_coerces_to_bool_by_zero_test(int input, bool expected)
        {
            IPlcValue value = new PlcDINT(input);

            Assert.Equal(expected, value.GetBool());
        }

        [Fact]
        public void Bool_reads_back_through_interface()
        {
            Assert.True(((IPlcValue) new PlcBOOL(true)).GetBool());
            Assert.False(((IPlcValue) new PlcBOOL(false)).GetBool());
        }

        [Fact]
        public void String_reads_back_through_interface()
        {
            IPlcValue value = new PlcSTRING("hello");

            Assert.True(value.IsString());
            Assert.Equal("hello", value.GetString());
        }

        [Fact]
        public void Simple_values_report_themselves_as_simple()
        {
            // PlcSimpleValueAdapter.IsSimple() hid the base method rather than overriding
            // it, so the interface dispatched to the base and answered false.
            Assert.True(((IPlcValue) new PlcDINT(1)).IsSimple());
            Assert.True(((IPlcValue) new PlcSTRING("x")).IsSimple());
        }

        [Fact]
        public void Byte_exposes_its_individual_bits()
        {
            IPlcValue value = new PlcBYTE(0xA1); // 1010 0001

            Assert.Equal(8, value.GetBoolLength());
            Assert.True(value.GetBoolAt(0));
            Assert.False(value.GetBoolAt(1));
            Assert.True(value.GetBoolAt(5));
            Assert.True(value.GetBoolAt(7));
            Assert.Equal(
                new[] { true, false, false, false, false, true, false, true },
                value.GetBoolArray());
        }

        [Fact]
        public void Real_reads_back_through_interface()
        {
            IPlcValue value = new PlcREAL(1.5f);

            Assert.True(value.IsFloat());
            Assert.Equal(1.5f, value.GetFloat());
            Assert.Equal(1.5d, value.GetDouble(), 6);
        }

        [Fact]
        public void Struct_exposes_its_members_through_the_interface()
        {
            // PlcStruct stored the dictionary but overrode none of the struct
            // accessors, so IsStruct()/GetValue() answered off the base.
            IPlcValue value = new PlcStruct(new Dictionary<string, IPlcValue>
            {
                ["control"] = new PlcBOOL(true),
                ["level"] = new PlcUSINT(7),
            });

            Assert.True(value.IsStruct());
            Assert.True(value.HasKey("control"));
            Assert.False(value.HasKey("missing"));
            Assert.True(value.GetValue("control").GetBool());
            Assert.Equal((byte) 7, value.GetValue("level").GetByte());
            Assert.Equal(new[] { "control", "level" }, value.GetKeys().OrderBy(k => k));
        }

        [Fact]
        public void Structs_are_equal_regardless_of_member_order()
        {
            var a = new PlcStruct(new Dictionary<string, IPlcValue>
            {
                ["x"] = new PlcBOOL(true),
                ["y"] = new PlcUSINT(1),
            });
            var b = new PlcStruct(new Dictionary<string, IPlcValue>
            {
                ["y"] = new PlcUSINT(1),
                ["x"] = new PlcBOOL(true),
            });

            Assert.Equal(a, b);
            Assert.Equal(a.GetHashCode(), b.GetHashCode());
        }

        [Fact]
        public void Raw_byte_array_reads_back_and_compares_by_content()
        {
            // PlcRawByteArray discarded its bytes behind the base GetRaw() and
            // compared by reference.
            IPlcValue value = new PlcRawByteArray(new byte[] { 0x0A, 0x0B });

            Assert.Equal(new byte[] { 0x0A, 0x0B }, value.GetRaw());
            Assert.Equal(2, value.GetLength());
            Assert.Equal(
                new PlcRawByteArray(new byte[] { 0x0A, 0x0B }),
                new PlcRawByteArray(new byte[] { 0x0A, 0x0B }));
            Assert.NotEqual(
                new PlcRawByteArray(new byte[] { 0x0A, 0x0B }),
                new PlcRawByteArray(new byte[] { 0x0A, 0x0C }));
        }

        [Fact]
        public void Bit_string_is_also_a_plain_number_and_string_through_the_interface()
        {
            // PlcBitString overrode only the unsigned + bit accessors, so GetInt()
            // returned 0 and GetString() null on a value held as IPlcValue.
            IPlcValue w = new PlcWORD(0x1234);
            Assert.True(w.IsSimple());
            Assert.Equal(0x1234, w.GetInt());
            Assert.Equal((short) 0x1234, w.GetShort());
            Assert.Equal(0x1234L, w.GetLong());
            Assert.Equal("4660", w.GetString());
            Assert.Equal((ushort) 0x1234, w.GetUshort());

            IPlcValue b = new PlcBYTE(0xA1);
            Assert.Equal(0xA1, b.GetInt());
            Assert.Equal("161", b.GetString());
        }

        [Fact]
        public void Bool_is_also_a_number_and_lowercase_text_through_the_interface()
        {
            IPlcValue t = new PlcBOOL(true);
            Assert.Equal(1, t.GetInt());
            Assert.Equal((byte) 1, t.GetByte());
            Assert.Equal("true", t.GetString());
            Assert.Equal("false", ((IPlcValue) new PlcBOOL(false)).GetString());
        }

        [Fact]
        public void Out_of_range_temporal_wire_data_throws_ParseException_not_a_framework_exception()
        {
            // The generated DataItem parse path feeds these raw wire numbers
            // straight in; a framework ArgumentOutOfRangeException would escape
            // MessageCodecBase's catch and tear down the receive loop.
            Assert.Throws<org.apache.plc4net.spi.generation.ParseException>(
                () => PlcTIME_OF_DAY.OfMillisecondsSinceMidnight(0xFFFFFFFFL));
            Assert.Throws<org.apache.plc4net.spi.generation.ParseException>(
                () => PlcLTIME_OF_DAY.OfNanosecondsSinceMidnight(ulong.MaxValue));
            Assert.Throws<org.apache.plc4net.spi.generation.ParseException>(
                () => PlcDATE_AND_TIME.OfSegments(0, 0, 0, 0, 0, 0, 0));
            Assert.Throws<org.apache.plc4net.spi.generation.ParseException>(
                () => PlcDATE_AND_LTIME.OfSegments(70000, 1, 1, 0, 0, 0, 0));
        }

        [Fact]
        public void Null_value_is_null_equatable_and_hashable()
        {
            // PlcNULL.Equals()/GetHashCode() threw NotImplementedException and
            // IsNull() answered false off the base — it is the fall-through of
            // every generated DataItem.StaticParse.
            IPlcValue n = new PlcNULL();
            Assert.True(n.IsNull());
            Assert.True(n.IsNullable());
            Assert.Equal(new PlcNULL(), new PlcNULL());
            Assert.Equal(new PlcNULL().GetHashCode(), new PlcNULL().GetHashCode());
            _ = new HashSet<IPlcValue> { new PlcNULL(), new PlcNULL() };
        }
    }
}
