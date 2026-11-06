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
using org.apache.plc4net.drivers.slmp;
using org.apache.plc4net.drivers.slmp.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.transports.test;
using Xunit;

namespace org.apache.plc4net.test.slmp
{
    public class SlmpMessageCodecTests
    {
        private const string Response = "D00000FFFF0300040000003412";

        [Theory]
        [InlineData(1)]
        [InlineData(7)]
        [InlineData(8)]
        [InlineData(9)]
        [InlineData(10)]
        [InlineData(12)]
        public void Partial_frames_wait_without_consuming_bytes(int split)
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var messages = new List<SlmpMessage>();
            var codec = new SlmpMessageCodec(transport, messages.Add);
            var wire = Convert.FromHexString(Response);
            Assert.Equal(split, transport.InjectTestData(wire.Take(split).ToArray()));
            codec.ProcessIncomingData();
            Assert.Empty(messages);
            Assert.Equal(split, transport.GetNumBytesAvailable());
            transport.InjectTestData(wire.Skip(split).ToArray());
            codec.ProcessIncomingData();
            var response = Assert.IsType<SlmpResponseFrame3E>(Assert.Single(messages));
            Assert.Equal((ushort)0, response.EndCode);
            Assert.Equal(new byte[] { 0x34, 0x12 }, response.ResponseData);
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        [Fact]
        public void Coalesced_frames_are_delivered_in_order_including_error_responses()
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var messages = new List<SlmpMessage>();
            var codec = new SlmpMessageCodec(transport, messages.Add);
            var wire = Convert.FromHexString(Response + "D00000FFFF030002005BC0");
            transport.InjectTestData(wire);
            codec.ProcessIncomingData();
            Assert.Equal(2, messages.Count);
            Assert.Equal(new byte[] { 0x34, 0x12 }, Assert.IsType<SlmpResponseFrame3E>(messages[0]).ResponseData);
            var error = Assert.IsType<SlmpResponseFrame3E>(messages[1]);
            Assert.Equal((ushort)0xC05B, error.EndCode);
            Assert.Empty(error.ResponseData);
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        [Theory]
        [InlineData("D00000FFFF03000000")]
        [InlineData("D00000FFFF03000100")]
        [InlineData("500000FFFF03000200")]
        [InlineData("D00100FFFF03000200")]
        [InlineData("D00001FFFF03000200")]
        [InlineData("010203040506070809")]
        public void Invalid_headers_resynchronize_to_the_next_supported_frame(string prefix)
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var messages = new List<SlmpMessage>();
            var codec = new SlmpMessageCodec(transport, messages.Add);
            transport.InjectTestData(Convert.FromHexString(prefix + Response));
            codec.ProcessIncomingData();
            Assert.IsType<SlmpResponseFrame3E>(Assert.Single(messages));
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        [Fact]
        public void Length_field_is_little_endian_and_can_span_more_than_one_byte()
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var messages = new List<SlmpMessage>();
            var codec = new SlmpMessageCodec(transport, messages.Add);
            var data = Enumerable.Range(0, 256).Select(i => (byte)i).ToArray();
            var wire = Convert.FromHexString("D00000FFFF030002010000").Concat(data).ToArray();
            transport.InjectTestData(wire);
            codec.ProcessIncomingData();
            Assert.Equal(data, Assert.IsType<SlmpResponseFrame3E>(Assert.Single(messages)).ResponseData);
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        [Fact]
        public void Maximum_response_length_is_representable_on_the_wire()
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            MessageCodecBase<SlmpMessage> codec = new SlmpMessageCodec(transport, _ => { });
            codec.Send(new SlmpResponseFrame3E(0, new byte[65533]));
            var wire = transport.GetAllWrittenData();
            Assert.Equal(65544, wire.Length);
            Assert.Equal(0xFF, wire[7]);
            Assert.Equal(0xFF, wire[8]);
        }

        [Theory]
        [InlineData(65534)]
        [InlineData(65535)]
        [InlineData(65536)]
        public void Oversized_responses_fail_before_any_transport_write(int payloadLength)
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            // The guard must also run through a reference to the shared base codec.
            MessageCodecBase<SlmpMessage> codec = new SlmpMessageCodec(transport, _ => { });
            Assert.Throws<MessageCodecException>(() => codec.Send(new SlmpResponseFrame3E(0, new byte[payloadLength])));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public void Oversized_write_requests_fail_before_any_transport_write()
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var codec = new SlmpMessageCodec(transport, _ => { });
            var data = new SlmpWriteRequest(100, SlmpDeviceCode.D, 32762, new byte[65524]);
            var request = new SlmpRequestFrame3E(16, 0x1401, 0, data);
            Assert.Throws<MessageCodecException>(() => codec.Send(request));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public void Read_requests_match_fixed_wire_bytes()
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var codec = new SlmpMessageCodec(transport, _ => { });
            codec.Send(new SlmpRequestFrame3E(16, 0x0401, 0, new SlmpReadRequest(100, SlmpDeviceCode.D, 1)));
            Assert.Equal(Convert.FromHexString("500000FFFF03000C00100001040000640000A80100"), transport.GetAllWrittenData());
        }

        [Fact]
        public void Sending_serializes_the_generated_message_to_fixed_wire_bytes()
        {
            using var transport = new TestTransportInstance(new TestTransportConfiguration());
            var codec = new SlmpMessageCodec(transport, _ => { });
            codec.Send(new SlmpResponseFrame3E(0, new byte[] { 0x34, 0x12 }));
            Assert.Equal(Convert.FromHexString(Response), transport.GetAllWrittenData());
        }
    }
}
