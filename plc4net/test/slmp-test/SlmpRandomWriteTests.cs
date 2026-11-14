// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
// https://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.

using System;
using System.Linq;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.slmp;
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.transports.test;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.test.slmp
{
    public class SlmpRandomWriteTests
    {
        [Fact]
        public async Task Scalar_word_and_double_word_values_share_one_random_write()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0));
            var request = connection.WriteRequestBuilder!
                .AddTag("real", "D1500:REAL", 1.0f)
                .AddTag("word", "D0:WORD", (ushort)0x1234)
                .AddTag("uint", "W0x100:UINT", (ushort)0x5678)
                .Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.Equal(new[] { "real", "word", "uint" }, response.TagNames);
            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            Assert.Equal(Convert.FromHexString(
                "500000FFFF03001C001000021400000201" +
                "000000A83412000100B47856DC0500A80000803F"),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Bit_and_range_writes_keep_using_batch_write()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var bit = transport.WaitForWrittenData(22, 5000);
                transport.InjectTestData(Response(0));
                var range = transport.WaitForWrittenData(25, 5000);
                transport.InjectTestData(Response(0));
                return bit.Concat(range).ToArray();
            });
            var request = connection.WriteRequestBuilder!
                .AddTag("bit", "M0:BOOL", true)
                .AddTag("range", "D0[0..1]:WORD", (ushort)1, (ushort)2)
                .Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(0x1401, written[11] | written[12] << 8);
            Assert.Equal(1, written[13] | written[14] << 8);
            var second = 9 + (written[7] | written[8] << 8);
            Assert.Equal(0x1401, written[second + 11] | written[second + 12] << 8);
            Assert.Equal(0, written[second + 13] | written[second + 14] << 8);
        }

        [Fact]
        public async Task Random_write_applies_a_remote_error_to_every_grouped_tag()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0xC05B));
            var request = connection.WriteRequestBuilder!
                .AddTag("first", "D0:WORD", (ushort)1)
                .AddTag("second", "D1:WORD", (ushort)2)
                .Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("first"));
            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("second"));
            var written = transport.GetAllWrittenData();
            Assert.Equal(0x1402, written[11] | written[12] << 8);
        }

        [Fact]
        public async Task More_than_160_points_split_into_random_write_then_batch_write()
        {
            using var transport = NewTransport(2048);
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var random = transport.WaitForWrittenData(977, 5000);
                transport.InjectTestData(Response(0));
                var batch = transport.WaitForWrittenData(23, 5000);
                transport.InjectTestData(Response(0));
                return random.Concat(batch).ToArray();
            });
            var builder = connection.WriteRequestBuilder!;
            for (var i = 0; i < 161; i++)
                builder.AddTag("d" + i, "D" + i + ":WORD", (ushort)i);

            var response = Assert.IsType<DefaultPlcWriteResponse>(await builder.Build().ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(0x1402, written[11] | written[12] << 8);
            Assert.Equal(160, written[15]);
            var second = 9 + (written[7] | written[8] << 8);
            Assert.Equal(0x1401, written[second + 11] | written[second + 12] << 8);
        }

        private static TestTransportInstance NewTransport(int bufferSize = 1024)
            => new TestTransportInstance(new TestTransportConfiguration { ReceiveBufferSize = bufferSize });

        private static SlmpConnection NewConnection(TestTransportInstance transport)
            => new SlmpConnection(ConnectionString.Parse("slmp://localhost"), transport);

        private static byte[] Response(ushort endCode) => new byte[]
        {
            0xD0, 0x00, 0x00, 0xFF, 0xFF, 0x03, 0x00,
            0x02, 0x00, (byte)endCode, (byte)(endCode >> 8)
        };
    }
}
