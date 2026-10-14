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

using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.exceptions;
using Xunit;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>
    /// The tag syntax: plc4j's address with the data type in it, the subset that is supported,
    /// and a clear refusal for the rest.
    /// </summary>
    public class S7TagTests
    {
        [Theory]
        [InlineData("%DB1.DBW10:INT", MemoryArea.DATA_BLOCKS, 1, 10, 0, TransportSize.INT)]
        [InlineData("%DB1.DBX0.0:BOOL", MemoryArea.DATA_BLOCKS, 1, 0, 0, TransportSize.BOOL)]
        [InlineData("%DB100.DBX5.7:BOOL", MemoryArea.DATA_BLOCKS, 100, 5, 7, TransportSize.BOOL)]
        [InlineData("%DB1.DBB3:BYTE", MemoryArea.DATA_BLOCKS, 1, 3, 0, TransportSize.BYTE)]
        [InlineData("%DB1.DBB3:SINT", MemoryArea.DATA_BLOCKS, 1, 3, 0, TransportSize.SINT)]
        [InlineData("%DB1.DBB3:USINT", MemoryArea.DATA_BLOCKS, 1, 3, 0, TransportSize.USINT)]
        [InlineData("%DB1.DBW2:WORD", MemoryArea.DATA_BLOCKS, 1, 2, 0, TransportSize.WORD)]
        [InlineData("%DB1.DBW2:UINT", MemoryArea.DATA_BLOCKS, 1, 2, 0, TransportSize.UINT)]
        [InlineData("%DB2.DBD4:DWORD", MemoryArea.DATA_BLOCKS, 2, 4, 0, TransportSize.DWORD)]
        [InlineData("%DB2.DBD4:DINT", MemoryArea.DATA_BLOCKS, 2, 4, 0, TransportSize.DINT)]
        [InlineData("%DB2.DBD4:UDINT", MemoryArea.DATA_BLOCKS, 2, 4, 0, TransportSize.UDINT)]
        [InlineData("%DB2.DBD8:REAL", MemoryArea.DATA_BLOCKS, 2, 8, 0, TransportSize.REAL)]
        [InlineData("%DB1.DB10:INT", MemoryArea.DATA_BLOCKS, 1, 10, 0, TransportSize.INT)] // the size letter is optional
        [InlineData("%DB1:10:INT", MemoryArea.DATA_BLOCKS, 1, 10, 0, TransportSize.INT)]
        [InlineData("%DB1:0.1:BOOL", MemoryArea.DATA_BLOCKS, 1, 0, 1, TransportSize.BOOL)]
        [InlineData("%I0.0:BOOL", MemoryArea.INPUTS, 0, 0, 0, TransportSize.BOOL)]
        [InlineData("%IB2:BYTE", MemoryArea.INPUTS, 0, 2, 0, TransportSize.BYTE)]
        [InlineData("%QW4:WORD", MemoryArea.OUTPUTS, 0, 4, 0, TransportSize.WORD)]
        [InlineData("%M10.3:BOOL", MemoryArea.FLAGS_MARKERS, 0, 10, 3, TransportSize.BOOL)]
        [InlineData("%MD12:REAL", MemoryArea.FLAGS_MARKERS, 0, 12, 0, TransportSize.REAL)]
        [InlineData("%M20:INT", MemoryArea.FLAGS_MARKERS, 0, 20, 0, TransportSize.INT)]
        [InlineData("%DB65535.DBW65535:INT", MemoryArea.DATA_BLOCKS, 65535, 65535, 0, TransportSize.INT)]
        public void A_address_names_its_area_block_offset_and_type(
            string address, MemoryArea area, int block, int byteOffset, int bitOffset, TransportSize type)
        {
            var tag = S7Tag.Parse(address);

            Assert.Equal(area, tag.Area);
            Assert.Equal(block, tag.DbNumber);
            Assert.Equal(byteOffset, tag.ByteOffset);
            Assert.Equal(bitOffset, tag.BitOffset);
            Assert.Equal(type, tag.DataType);
        }

        [Theory]
        [InlineData("%db1.dbw10:int")]
        [InlineData("%Db1.dBw10:Int")]
        [InlineData("  %DB1.DBW10:INT  ")]
        [InlineData("%m0.0:bool")]
        public void Area_type_and_surrounding_blanks_are_not_significant(string address)
        {
            Assert.NotNull(S7Tag.Parse(address));
        }

        [Theory]
        [InlineData("%DB1.DBW10:INT")]
        [InlineData("%DB1.DBX0.5:BOOL")]
        [InlineData("%DB7.DBD100:REAL")]
        [InlineData("%I3.2:BOOL")]
        [InlineData("%QB9:BYTE")]
        [InlineData("%MW30:UINT")]
        public void The_string_form_parses_back_to_the_same_tag(string address)
        {
            var tag = S7Tag.Parse(address);

            var again = S7Tag.Parse(tag.ToString());

            Assert.Equal(tag.Area, again.Area);
            Assert.Equal(tag.DbNumber, again.DbNumber);
            Assert.Equal(tag.ByteOffset, again.ByteOffset);
            Assert.Equal(tag.BitOffset, again.BitOffset);
            Assert.Equal(tag.DataType, again.DataType);
        }

        [Theory]
        [InlineData("%DB1.DBW10:INT", 2)]
        [InlineData("%DB1.DBX0.0:BOOL", 1)]
        [InlineData("%DB1.DBB0:BYTE", 1)]
        [InlineData("%DB1.DBD0:REAL", 4)]
        [InlineData("%DB1.DBD0:DINT", 4)]
        public void A_tag_knows_how_many_bytes_its_value_takes(string address, int size)
        {
            Assert.Equal(size, S7Tag.Parse(address).SizeInBytes);
        }

        [Theory]
        [InlineData("")]
        [InlineData("   ")]
        [InlineData("DB1.DBW10:INT")] // no %
        [InlineData("%DB1.DBW10")] // plc4j needs the data type in the address
        [InlineData("%M0.0")]
        [InlineData("%DB1.DBW:INT")] // no offset
        [InlineData("%DB.DBW10:INT")] // no block
        [InlineData("%DB1.DBW10:")] // no type
        [InlineData("%DB1.DBW10.:INT")]
        [InlineData("%DB1.DBW10:INT:INT")]
        [InlineData("%DB1.DBW10:INT junk")]
        [InlineData("%DB1.DBB0:STRING(10)")]
        [InlineData("%DB100000.DBW0:INT")] // six digits
        [InlineData("%M12345678:BYTE")] // eight digits
        [InlineData("%M0.8:BOOL")] // bit 8
        [InlineData("holding-register:1")]
        public void An_address_that_is_not_an_s7_tag_is_rejected(string address)
        {
            Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
        }

        [Fact]
        public void A_non_ascii_digit_is_rejected_as_an_invalid_tag()
        {
            // \d in .NET matches every Unicode decimal digit; the Arabic-Indic one would pass such a
            // pattern and then fail in int.Parse with a FormatException instead of this.
            Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse("%DB١.DBW0:INT"));
            Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse("%DB1.DBW١٠:INT"));
        }

        [Theory]
        [InlineData("%M0:BOOL")]
        [InlineData("%DB1.DBX0:BOOL")]
        public void A_BOOL_needs_its_bit(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
            Assert.Contains("bit offset", exception.Message);
        }

        [Theory]
        [InlineData("%MW0.3:INT")]
        [InlineData("%DB1.DBB0.1:BYTE")]
        public void Only_a_BOOL_has_a_bit_other_than_0(string address)
        {
            Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
        }

        [Fact]
        public void Bit_0_of_a_wider_type_is_what_it_would_be_without_it()
        {
            Assert.Equal(10, S7Tag.Parse("%DB1.DBW10.0:INT").ByteOffset);
        }

        [Theory]
        [InlineData("%DB1.DBW10:REAL")] // W is an INT's letter
        [InlineData("%MD0:INT")]
        [InlineData("%DB1.DBX0.0:BYTE")]
        [InlineData("%QB0:WORD")]
        public void A_size_letter_has_to_be_the_one_of_the_type(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
            Assert.Contains("size letter", exception.Message);
        }

        [Theory]
        [InlineData("%DB65536.DBW0:INT")]
        [InlineData("%DB1.DBW65536:INT")]
        [InlineData("%M70000:BYTE")]
        public void A_block_or_offset_that_the_address_item_cannot_carry_is_rejected(string address)
        {
            Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
        }

        [Fact]
        public void Array_selections_are_refused_with_an_explanation()
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse("%DB1.DBW10[4]:INT"));
            Assert.Contains("array", exception.Message);
        }

        [Theory]
        [InlineData("%DB1.DBD0:LREAL")]
        [InlineData("%DB1.DBD0:LWORD")]
        [InlineData("%DB1.DBD0:LINT")]
        [InlineData("%DB1.DBD0:TIME")]
        [InlineData("%DB1.DBW0:DATE")]
        [InlineData("%DB1.DBW0:S5TIME")]
        [InlineData("%DB1.DBB0:CHAR")]
        [InlineData("%DB1.DBX0.0:COUNTER")]
        [InlineData("%DB1.DBB0:STRING")]
        public void The_types_that_are_not_supported_yet_are_refused_with_an_explanation(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
            Assert.Contains("not supported yet", exception.Message);
        }

        [Fact]
        public void An_unknown_type_is_rejected()
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse("%M0.0:FLOATY"));
            Assert.Contains("Unknown S7 data type", exception.Message);
        }

        [Theory]
        [InlineData("%C5:WORD")]
        [InlineData("%T0:WORD")]
        [InlineData("%PIW0:WORD")]
        [InlineData("%L0.0:BOOL")]
        [InlineData("%D10:INT")]
        public void The_other_memory_areas_are_refused_with_an_explanation(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => S7Tag.Parse(address));
            Assert.Contains("not supported yet", exception.Message);
        }

        [Fact]
        public void The_constructor_applies_the_same_rules_as_the_parser()
        {
            Assert.Throws<PlcInvalidFieldException>(
                () => new S7Tag(MemoryArea.COUNTERS, 0, 0, 0, TransportSize.WORD));
            Assert.Throws<PlcInvalidFieldException>(
                () => new S7Tag(MemoryArea.DATA_BLOCKS, 1, 0, 0, TransportSize.LREAL));
            Assert.Throws<PlcInvalidFieldException>(
                () => new S7Tag(MemoryArea.FLAGS_MARKERS, 5, 0, 0, TransportSize.BYTE)); // a block number outside a data block
            Assert.Throws<PlcInvalidFieldException>(
                () => new S7Tag(MemoryArea.FLAGS_MARKERS, 0, 0, 8, TransportSize.BOOL));
            Assert.Throws<PlcInvalidFieldException>(
                () => new S7Tag(MemoryArea.FLAGS_MARKERS, 0, 0, 1, TransportSize.BYTE));
        }
    }
}
