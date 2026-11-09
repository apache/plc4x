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
using org.apache.plc4net.drivers.slmp;
using org.apache.plc4net.drivers.slmp.readwrite.model;
using org.apache.plc4net.exceptions;
using Xunit;

namespace org.apache.plc4net.test.slmp
{
    public class SlmpTagTests
    {
        [Theory]
        [InlineData("D350", SlmpDeviceCode.D, 350, SlmpDataType.WORD, 1, false, 1)]
        [InlineData("w1a[0..9]:word", SlmpDeviceCode.W, 26, SlmpDataType.WORD, 10, true, 10)]
        [InlineData("W0X1A:UINT", SlmpDeviceCode.W, 26, SlmpDataType.UINT, 1, false, 1)]
        [InlineData("R200[0..3]:REAL", SlmpDeviceCode.R, 200, SlmpDataType.REAL, 4, true, 8)]
        [InlineData("D100[4]:DINT", SlmpDeviceCode.D, 108, SlmpDataType.DINT, 1, false, 2)]
        [InlineData("D100[4..4]:INT", SlmpDeviceCode.D, 104, SlmpDataType.INT, 1, true, 1)]
        [InlineData("D100[4..5]:UDINT", SlmpDeviceCode.D, 108, SlmpDataType.UDINT, 2, true, 4)]
        [InlineData("D0[0..959]", SlmpDeviceCode.D, 0, SlmpDataType.WORD, 960, true, 960)]
        [InlineData("D0[0..479]:REAL", SlmpDeviceCode.D, 0, SlmpDataType.REAL, 480, true, 960)]
        [InlineData("WFFFFFF", SlmpDeviceCode.W, 16777215, SlmpDataType.WORD, 1, false, 1)]
        public void Parses_device_selection_and_word_footprint(string text, SlmpDeviceCode device, int number,
            SlmpDataType type, int quantity, bool range, int points)
        {
            var tag = SlmpTag.Parse(text);
            Assert.Equal(device, tag.DeviceCode);
            Assert.Equal(number, tag.DeviceNumber);
            Assert.Equal(type, tag.DataType);
            Assert.Equal(quantity, tag.Quantity);
            Assert.Equal(range, tag.ExplicitRange);
            Assert.Equal(points, tag.NumberOfPoints);
            var canonical = SlmpTag.Parse(tag.ToString());
            Assert.Equal(tag.ToString(), canonical.ToString());
            Assert.Equal(tag.DeviceNumber, canonical.DeviceNumber);
            Assert.Equal(tag.Quantity, canonical.Quantity);
            Assert.Equal(tag.ExplicitRange, canonical.ExplicitRange);
        }

        [Theory]
        [InlineData(null)]
        [InlineData("")]
        [InlineData("M100")]
        [InlineData("D0x10")]
        [InlineData("D1A")]
        [InlineData("D-1")]
        [InlineData("D\u0661")]
        [InlineData("D0:1")]
        [InlineData("D0:BOOL")]
        [InlineData("D0:REAL[4]")]
        [InlineData("D0[]")]
        [InlineData("D0[3..2]")]
        [InlineData("D0[-1]")]
        [InlineData("D0[0..960]")]
        [InlineData("D0[0..480]:REAL")]
        [InlineData("D16777216")]
        [InlineData("W1000000")]
        [InlineData("D16777215:REAL")]
        [InlineData("D16777215[0..1]")]
        [InlineData("D16777214[1]:REAL")]
        [InlineData("D0[0..4294967295]:REAL")]
        [InlineData("D0[4294967295]:REAL")]
        [InlineData("D0[99999999999999999999]")]
        [InlineData("WFFFFFFFFFFFFFFFFFFFFFFFF")]
        [InlineData("D0[0,1]")]
        [InlineData("D0:INT\nD1")]
        public void Invalid_or_unrepresentable_selections_are_rejected(string? text)
            => Assert.Throws<PlcInvalidFieldException>(() => SlmpTag.Parse(text!));

        [Fact]
        public void Constructor_validates_the_same_invariants()
        {
            Assert.Throws<PlcInvalidFieldException>(() => new SlmpTag(SlmpDeviceCode.M, 0, SlmpDataType.WORD));
            Assert.Throws<PlcInvalidFieldException>(() => new SlmpTag(SlmpDeviceCode.D, -1, SlmpDataType.WORD));
            Assert.Throws<PlcInvalidFieldException>(() => new SlmpTag(SlmpDeviceCode.D, 0, (SlmpDataType)99));
            Assert.Throws<PlcInvalidFieldException>(() => new SlmpTag(SlmpDeviceCode.D, 0, SlmpDataType.REAL, int.MaxValue));
            Assert.Throws<PlcInvalidFieldException>(() => new SlmpTag(SlmpDeviceCode.D, 0, SlmpDataType.WORD, 0));
            Assert.Throws<PlcInvalidFieldException>(() => new SlmpTag(SlmpDeviceCode.D, 0, SlmpDataType.WORD, 2, false));
            Assert.True(new SlmpTag(SlmpDeviceCode.D, 0, SlmpDataType.WORD, 2).ExplicitRange);
        }
    }
}
