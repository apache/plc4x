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
using org.apache.plc4net.transports.cotp;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    public class TpktFrameTests
    {
        [Fact]
        public void Wrap_prefixes_the_payload_with_a_header_describing_the_whole_frame()
        {
            var frame = TpktFrame.Wrap(new byte[] { 0xAA, 0xBB });

            // Version 3, reserved 0, then the big-endian length of header + payload.
            Assert.Equal(new byte[] { 0x03, 0x00, 0x00, 0x06, 0xAA, 0xBB }, frame);
        }

        [Fact]
        public void Wrap_accepts_the_largest_payload_the_length_field_can_describe()
        {
            // The TPKT length field is 16 bits and counts the 4-byte header, so the
            // largest payload is 65535 - 4 = 65531 bytes.
            var frame = TpktFrame.Wrap(new byte[65531]);

            Assert.Equal(65535, frame.Length);
            Assert.Equal(0xFF, frame[2]);
            Assert.Equal(0xFF, frame[3]);
        }

        [Theory]
        [InlineData(65532)]
        [InlineData(70000)]
        public void Wrap_rejects_a_payload_the_length_field_cannot_describe(int payloadLength)
        {
            var payload = new byte[payloadLength];

            // Without the check the length wraps modulo 65536 and the frame declares a size
            // that does not match its contents.
            var ex = Assert.Throws<ArgumentException>(() => TpktFrame.Wrap(payload));

            Assert.Equal("payload", ex.ParamName);
        }
    }
}