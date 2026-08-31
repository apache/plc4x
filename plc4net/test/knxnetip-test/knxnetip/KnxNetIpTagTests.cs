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
        public void Encoding_a_wildcard_address_is_an_error()
        {
            Assert.Throws<InvalidOperationException>(() => KnxNetIpTag.Parse("1/*/3").ToWireAddress(3));
        }

        [Fact]
        public void Matches_an_inbound_wire_address_through_a_wildcard()
        {
            var tag = KnxNetIpTag.Parse("2/*/*");
            Assert.True(tag.MatchesWireAddress(new byte[] { 0x12, 0x0C }, 3));   // 2/2/12
            Assert.False(tag.MatchesWireAddress(new byte[] { 0x1A, 0x0C }, 3));  // 3/2/12
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
