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
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

using System;
using System.Linq;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.slmp;
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.transports.test;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.test.slmp
{
    public class SlmpMultiBlockReadTests
    {
        [Fact]
        public async Task Word_and_bit_ranges_share_one_multi_block_read_and_keep_request_order()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString("3412785681800900")));
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("bits", "M0[0..31]:BOOL")
                .AddTagAddress("words", "D0[0..1]:WORD")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.Equal(new[] { "bits", "words" }, response.TagNames);
            var words = Assert.IsType<PlcList>(response.GetValue("words"));
            Assert.Equal((ushort)0x1234, words.GetIndex(0).GetUshort());
            Assert.Equal((ushort)0x5678, words.GetIndex(1).GetUshort());
            var bits = Assert.IsType<PlcList>(response.GetValue("bits"));
            Assert.Equal(32, bits.GetLength());
            Assert.True(bits.GetIndex(0).GetBool());
            Assert.False(bits.GetIndex(1).GetBool());
            Assert.True(bits.GetIndex(7).GetBool());
            Assert.False(bits.GetIndex(8).GetBool());
            Assert.True(bits.GetIndex(15).GetBool());
            Assert.True(bits.GetIndex(16).GetBool());
            Assert.True(bits.GetIndex(19).GetBool());
            Assert.Equal(Convert.FromHexString(
                "500000FFFF030014001000060400000101000000A80200000000900200"),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_partial_bit_block_falls_back_to_batch_read()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var multiBlock = transport.WaitForWrittenData(29, 5000);
                transport.InjectTestData(Response(0, new byte[8]));
                var batch = transport.WaitForWrittenData(21, 5000);
                transport.InjectTestData(Response(0, new byte[10]));
                return multiBlock.Concat(batch).ToArray();
            });
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("bits", "M0[0..19]:BOOL")
                .AddTagAddress("first", "D0[0..1]:WORD")
                .AddTagAddress("second", "D10[0..1]:WORD")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(2, written[15]);
            Assert.Equal(0, written[16]);
            var secondFrame = 9 + (written[7] | written[8] << 8);
            Assert.Equal(0x0401, written[secondFrame + 11] | written[secondFrame + 12] << 8);
            Assert.Equal(1, written[secondFrame + 13] | written[secondFrame + 14] << 8);
            Assert.Equal(20, written[secondFrame + 19] | written[secondFrame + 20] << 8);
        }

        [Theory]
        [InlineData(0x1234, "", PlcResponseCode.RemoteError)]
        [InlineData(0, "3412", PlcResponseCode.InvalidData)]
        public async Task Multi_block_read_applies_frame_errors_to_every_grouped_tag(
            ushort endCode, string dataHex, PlcResponseCode expected)
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(endCode, Convert.FromHexString(dataHex)));
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("first", "D0[0..1]:WORD")
                .AddTagAddress("second", "D10[0..1]:WORD")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.Equal(expected, response.GetResponseCode("first"));
            Assert.Equal(expected, response.GetResponseCode("second"));
            Assert.Null(response.GetValue("first"));
            Assert.Null(response.GetValue("second"));
        }

        [Fact]
        public async Task More_than_120_blocks_split_into_multi_block_then_batch_read()
        {
            using var transport = NewTransport(4096);
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var multiBlock = transport.WaitForWrittenData(737, 5000);
                transport.InjectTestData(Response(0, new byte[480]));
                var batch = transport.WaitForWrittenData(21, 5000);
                transport.InjectTestData(Response(0, new byte[4]));
                return multiBlock.Concat(batch).ToArray();
            });
            var builder = connection.ReadRequestBuilder!;
            for (var i = 0; i < 121; i++)
                builder.AddTagAddress("d" + i, "D" + i * 2 + "[0..1]:WORD");

            var response = Assert.IsType<DefaultPlcReadResponse>(await builder.Build().ExecuteAsync());

            Assert.Equal(121, response.TagNames.Count());
            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(0x0406, written[11] | written[12] << 8);
            Assert.Equal(120, written[15]);
            Assert.Equal(0, written[16]);
            var secondFrame = 9 + (written[7] | written[8] << 8);
            Assert.Equal(0x0401, written[secondFrame + 11] | written[secondFrame + 12] << 8);
            Assert.Equal(Convert.FromHexString("F00000A80200"),
                written.AsSpan(secondFrame + 15, 6).ToArray());
        }

        [Fact]
        public async Task More_than_960_points_split_at_the_combined_point_limit()
        {
            using var transport = NewTransport(4096);
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var multiBlock = transport.WaitForWrittenData(29, 5000);
                transport.InjectTestData(Response(0, new byte[1920]));
                var batch = transport.WaitForWrittenData(21, 5000);
                transport.InjectTestData(Response(0, new byte[4]));
                return multiBlock.Concat(batch).ToArray();
            });
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("first", "D0[0..599]:WORD")
                .AddTagAddress("second", "D600[0..359]:WORD")
                .AddTagAddress("third", "D960[0..1]:WORD")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(2, written[15]);
            Assert.Equal(0, written[16]);
            var secondFrame = 9 + (written[7] | written[8] << 8);
            Assert.Equal(0x0401, written[secondFrame + 11] | written[secondFrame + 12] << 8);
            Assert.Equal(Convert.FromHexString("C00300A80200"),
                written.AsSpan(secondFrame + 15, 6).ToArray());
        }

        private static TestTransportInstance NewTransport(int bufferSize = 1024)
            => new TestTransportInstance(new TestTransportConfiguration { ReceiveBufferSize = bufferSize });

        private static SlmpConnection NewConnection(TestTransportInstance transport)
            => new SlmpConnection(ConnectionString.Parse("slmp://localhost"), transport);

        private static byte[] Response(ushort endCode, byte[] data)
        {
            var length = checked((ushort)(2 + data.Length));
            return new byte[]
            {
                0xD0, 0x00, 0x00, 0xFF, 0xFF, 0x03, 0x00,
                (byte)length, (byte)(length >> 8),
                (byte)endCode, (byte)(endCode >> 8)
            }.Concat(data).ToArray();
        }
    }
}
