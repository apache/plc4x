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
    }
}
