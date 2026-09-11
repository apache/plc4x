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

using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.spi.generation;
using Xunit;

namespace org.apache.plc4net.test.knxnetip.readwrite.model
{
    /// <summary>
    /// Parse -> serialize -> byte-identical for a spread of datapoint shapes the
    /// generated <c>KnxDatapoint</c> dataIo has to cover: a bit, a scaled byte, the
    /// KNX 16-bit float, and a <c>PlcStruct</c>. This is the round-trip rigor the
    /// S7 <c>DataItem</c> tests already have.
    /// </summary>
    public class KnxDatapointRoundTripTests
    {
        [Theory]
        [InlineData(KnxDatapointType.DPT_Switch, new byte[] { 0x01 })]            // 1.001 bit
        [InlineData(KnxDatapointType.DPT_Switch, new byte[] { 0x00 })]
        [InlineData(KnxDatapointType.DPT_Scaling, new byte[] { 0x00, 0x80 })]     // 5.001 8-bit
        [InlineData(KnxDatapointType.DPT_Value_Temp, new byte[] { 0x00, 0x00, 0x00 })]  // 9.001 F16 = 0.0
        [InlineData(KnxDatapointType.DPT_Value_Temp, new byte[] { 0x00, 0x0C, 0x1A })]  // 9.001 F16 = 21.0
        [InlineData(KnxDatapointType.DPT_Colour_RGB, new byte[] { 0x00, 0xFF, 0x80, 0x40 })]  // 232.600 struct
        public void A_datapoint_round_trips_byte_identical(KnxDatapointType dpt, byte[] wire)
        {
            var value = KnxDatapoint.StaticParse(new ReadBuffer(wire), dpt);

            var buffer = new WriteBuffer();
            KnxDatapoint.StaticSerialize(buffer, value, dpt);

            Assert.Equal(wire, buffer.GetBytes());
        }

        [Fact]
        public void The_16bit_float_decodes_to_its_engineering_value()
        {
            var value = KnxDatapoint.StaticParse(
                new ReadBuffer(new byte[] { 0x00, 0x0C, 0x1A }), KnxDatapointType.DPT_Value_Temp);
            Assert.Equal(21.0f, value.GetFloat(), 3);
        }

        [Fact]
        public void The_rgb_struct_exposes_its_three_channels()
        {
            var value = KnxDatapoint.StaticParse(
                new ReadBuffer(new byte[] { 0x00, 0xFF, 0x80, 0x40 }), KnxDatapointType.DPT_Colour_RGB);

            Assert.True(value.IsStruct());
            Assert.Equal((byte) 0xFF, value.GetValue("r").GetByte());
            Assert.Equal((byte) 0x80, value.GetValue("g").GetByte());
            Assert.Equal((byte) 0x40, value.GetValue("b").GetByte());
        }
    }
}
