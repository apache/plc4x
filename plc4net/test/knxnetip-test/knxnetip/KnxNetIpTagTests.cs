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
using org.apache.plc4net.drivers.knxnetip;
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.exceptions;
using Xunit;

namespace org.apache.plc4net.test.knxnetip
{
    public class KnxNetIpTagTests
    {
        [Theory]
        [InlineData("1/2/3", 3, "1", "2", "3")]
        [InlineData("31/7/255", 3, "31", "7", "255")]
        [InlineData("1/2", 2, "1", "", "2")]
        [InlineData("42", 1, "42", "", "42")]
        public void Parses_the_three_group_address_notations(
            string address, int levels, string main, string middle, string sub)
        {
            var tag = KnxNetIpTag.Parse(address);
            Assert.Equal(levels, tag.Levels);
            Assert.Equal(main, tag.MainGroup);
            Assert.Equal(middle.Length == 0 ? null : middle, tag.MiddleGroup);
            Assert.Equal(sub, tag.SubGroup);
        }

        [Fact]
        public void Keeps_the_dpt_hint_from_the_suffix()
        {
            var tag = KnxNetIpTag.Parse("1/2/3:DPT9.001");
            Assert.Equal("DPT9.001", tag.DptId);
        }

        [Theory]
        [InlineData("*/*/*")]
        [InlineData("1/*/3")]
        public void Recognises_wildcards(string address)
        {
            Assert.True(KnxNetIpTag.Parse(address).HasWildcard);
        }

        [Theory]
        [InlineData("not an address")]
        [InlineData("1/2/3/4")]
        [InlineData("1/2/3:NOPE")]
        [InlineData("١/٢/٣")]          // Arabic-Indic digits: \d would match, [0-9] does not
        [InlineData("１/２/３")]        // fullwidth digits
        [InlineData("1/2/3\n")]        // trailing newline (\z, not $)
        public void Rejects_a_string_that_is_not_a_group_address(string address)
        {
            Assert.Throws<ArgumentException>(() => KnxNetIpTag.Parse(address));
        }

        [Fact]
        public void Encodes_a_three_level_address_to_its_two_wire_bytes()
        {
            // 2/2/12 packs as 5|3|8 bits -> 0b00010_010_00001100 = 0x120C, the
            // destination address in the shared Tunneling Request test vector.
            var wire = KnxNetIpTag.Parse("2/2/12").ToWireAddress(3);
            Assert.Equal(new byte[] { 0x12, 0x0C }, wire);
        }

        [Fact]
        public void Encodes_a_two_level_and_a_free_level_address()
        {
            // 3/1801 -> main 3 (5 bits), sub 1801 (11 bits) = 0b00011_00111001001
            Assert.Equal(new byte[] { 0x1F, 0x09 }, KnxNetIpTag.Parse("3/1801").ToWireAddress(2));
            // 42 -> 16 bits
            Assert.Equal(new byte[] { 0x00, 0x2A }, KnxNetIpTag.Parse("42").ToWireAddress(1));
        }

        [Theory]
        [InlineData("32/1/1", 3)]   // main > 31
        [InlineData("1/8/1", 3)]    // middle > 7
        [InlineData("1/1/256", 3)]  // sub > 255
        [InlineData("31/2048", 2)]  // 2-level sub > 2047
        public void Rejects_an_out_of_range_segment_instead_of_truncating_it(string address, int levels)
        {
            var tag = KnxNetIpTag.Parse(address);
            Assert.Throws<PlcInvalidFieldException>(() => tag.ToWireAddress(levels));
        }

        [Fact]
        public void Rejects_a_tag_whose_level_count_differs_from_the_connection()
        {
            // A 3-level tag on a 2-level connection would silently drop the middle group.
            Assert.Throws<PlcInvalidFieldException>(() => KnxNetIpTag.Parse("1/2/3").ToWireAddress(2));
            Assert.Throws<PlcInvalidFieldException>(() => KnxNetIpTag.Parse("1/2").ToWireAddress(3));
        }

        [Fact]
        public void Encoding_a_wildcard_address_is_an_error()
        {
            Assert.Throws<PlcInvalidFieldException>(() => KnxNetIpTag.Parse("1/*/3").ToWireAddress(3));
        }

        [Fact]
        public void Matches_an_inbound_wire_address_through_a_wildcard()
        {
            var tag = KnxNetIpTag.Parse("2/*/*");
            Assert.True(tag.MatchesWireAddress(new byte[] { 0x12, 0x0C }, 3));   // 2/2/12
            Assert.False(tag.MatchesWireAddress(new byte[] { 0x1A, 0x0C }, 3));  // 3/2/12
        }

        [Fact]
        public void Matches_an_inbound_wire_address_written_with_leading_zeros()
        {
            Assert.True(KnxNetIpTag.Parse("1/2/03").MatchesWireAddress(new byte[] { 0x0A, 0x03 }, 3));
        }

        [Fact]
        public void MatchesWireAddress_tolerates_a_short_buffer()
        {
            Assert.False(KnxNetIpTag.Parse("1/2/3").MatchesWireAddress(new byte[] { 0x0A }, 3));
            Assert.False(KnxNetIpTag.Parse("1/2/3").MatchesWireAddress(null, 3));
        }

        [Fact]
        public void Resolves_an_exact_dpt_to_the_generated_enum_constant()
        {
            // DPST-1-1 is DPT_Switch.
            var dpt = KnxNetIpTag.Parse("1/2/3:DPT1.001").ResolveDatapointType();
            Assert.Equal(KnxDatapointType.DPT_Switch, dpt);
        }

        [Fact]
        public void Resolves_a_bare_main_dpt_to_the_first_of_that_type()
        {
            var dpt = KnxNetIpTag.Parse("1/2/3:DPT1").ResolveDatapointType();
            Assert.NotNull(dpt);
            Assert.StartsWith("DPST-1-", dpt.Value.GetId());
        }

        [Fact]
        public void An_address_without_a_hint_resolves_to_no_datapoint_type()
        {
            Assert.Null(KnxNetIpTag.Parse("1/2/3").ResolveDatapointType());
        }
    }
}
