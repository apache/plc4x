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

using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.exceptions;
using Xunit;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The tag syntax: plc4j's <c>{area}:{address}[:{TYPE}]</c> with the address counted from 1,
    /// the subset of it that is supported, and a clear refusal for the rest.
    /// </summary>
    public class ModbusTagTests
    {
        [Theory]
        [InlineData("coil:1", ModbusTag.TagType.Coil, 1, ModbusDataType.BOOL)]
        [InlineData("discrete-input:7", ModbusTag.TagType.DiscreteInput, 7, ModbusDataType.BOOL)]
        [InlineData("holding-register:100", ModbusTag.TagType.HoldingRegister, 100, ModbusDataType.INT)]
        [InlineData("input-register:3", ModbusTag.TagType.InputRegister, 3, ModbusDataType.INT)]
        [InlineData("holding-register:5:UINT", ModbusTag.TagType.HoldingRegister, 5, ModbusDataType.UINT)]
        [InlineData("holding-register:5:WORD", ModbusTag.TagType.HoldingRegister, 5, ModbusDataType.WORD)]
        [InlineData("input-register:5:INT", ModbusTag.TagType.InputRegister, 5, ModbusDataType.INT)]
        [InlineData("coil:2:BOOL", ModbusTag.TagType.Coil, 2, ModbusDataType.BOOL)]
        [InlineData("holding-register:65536", ModbusTag.TagType.HoldingRegister, 65536, ModbusDataType.INT)]
        public void An_address_names_its_area_register_and_type(
            string address, ModbusTag.TagType type, int register, ModbusDataType dataType)
        {
            var tag = ModbusTag.Parse(address);

            Assert.Equal(type, tag.Type);
            Assert.Equal(register, tag.Address);
            Assert.Equal(dataType, tag.DataType);
        }

        [Theory]
        [InlineData("COIL:1")]
        [InlineData("Holding-Register:1:uint")]
        [InlineData("  holding-register:1  ")]
        public void Area_and_type_are_not_case_sensitive(string address)
        {
            Assert.NotNull(ModbusTag.Parse(address));
        }

        [Theory]
        [InlineData("holding-register:1", 0)]
        [InlineData("holding-register:2", 1)]
        [InlineData("holding-register:65536", 65535)]
        [InlineData("coil:11", 10)]
        public void The_request_carries_the_address_minus_one(string address, int wire)
        {
            Assert.Equal((ushort)wire, ModbusTag.Parse(address).WireAddress);
        }

        [Fact]
        public void A_bit_area_is_a_bit_and_a_register_area_is_not()
        {
            Assert.True(ModbusTag.Parse("coil:1").IsBit);
            Assert.True(ModbusTag.Parse("discrete-input:1").IsBit);
            Assert.False(ModbusTag.Parse("holding-register:1").IsBit);
            Assert.False(ModbusTag.Parse("input-register:1").IsBit);
        }

        [Theory]
        [InlineData("coil:1", "coil:1:BOOL")]
        [InlineData("holding-register:1", "holding-register:1:INT")]
        [InlineData("input-register:9:uint", "input-register:9:UINT")]
        public void The_string_form_is_the_complete_address(string address, string expected)
        {
            var tag = ModbusTag.Parse(address);

            Assert.Equal(expected, tag.ToString());
            // ...and it reads back as the same tag.
            var again = ModbusTag.Parse(tag.ToString());
            Assert.Equal(tag.Type, again.Type);
            Assert.Equal(tag.Address, again.Address);
            Assert.Equal(tag.DataType, again.DataType);
        }

        [Theory]
        [InlineData("")]
        [InlineData("   ")]
        [InlineData("holding-register")]
        [InlineData("holding-register:")]
        [InlineData("holding-register:abc")]
        [InlineData("holding-register:-1")]
        [InlineData("holding-register:+1")]
        [InlineData("holding-register:1.5")]
        [InlineData("holding-register: 1")]
        [InlineData("unknown:1")]
        [InlineData("holding:1")]   // revival-era spelling, not plc4j's
        [InlineData("input:1")]
        [InlineData("discrete:1")]
        public void An_address_that_is_not_a_modbus_tag_is_rejected(string address)
        {
            Assert.Throws<PlcInvalidFieldException>(() => ModbusTag.Parse(address));
        }

        [Theory]
        [InlineData("holding-register:0")]       // counting starts at 1
        [InlineData("holding-register:65537")]
        [InlineData("coil:0")]
        [InlineData("input-register:100000")]
        public void A_register_number_out_of_range_is_rejected(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => ModbusTag.Parse(address));
            Assert.Contains("between 1 and 65536", exception.Message);
        }

        [Theory]
        [InlineData("holding-register:1[0..3]:INT")]   // arrays
        [InlineData("holding-register:1:INT[4]")]
        [InlineData("holding-register:1:STRING(20)")]  // string lengths
        [InlineData("holding-register:1:INT{byte-order:'LITTLE_ENDIAN'}")] // tag configuration
        [InlineData("extended-register:1")]            // extended registers
        [InlineData("40001")]                          // short forms
        [InlineData("4x00001")]
        [InlineData("00001")]
        public void The_parts_of_the_syntax_that_are_not_supported_yet_are_refused_with_an_explanation(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => ModbusTag.Parse(address));
            Assert.Contains("not supported yet", exception.Message);
        }

        [Theory]
        [InlineData("holding-register:1:BOOL")]
        [InlineData("holding-register:1:REAL")]
        [InlineData("holding-register:1:DINT")]
        [InlineData("holding-register:1:BYTE")]
        [InlineData("input-register:1:LREAL")]
        [InlineData("coil:1:INT")]
        [InlineData("discrete-input:1:UINT")]
        public void A_data_type_that_does_not_fit_the_area_is_rejected(string address)
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(() => ModbusTag.Parse(address));
            Assert.Contains("not supported", exception.Message);
        }

        [Fact]
        public void An_unknown_data_type_is_rejected()
        {
            var exception = Assert.Throws<PlcInvalidFieldException>(
                () => ModbusTag.Parse("holding-register:1:FLOATY"));
            Assert.Contains("Unknown Modbus data type", exception.Message);
        }

        [Fact]
        public void A_number_is_not_a_data_type()
        {
            // Enum.TryParse would read "7" as ModbusDataType.INT.
            Assert.Throws<PlcInvalidFieldException>(() => ModbusTag.Parse("holding-register:1:7"));
        }
    }
}