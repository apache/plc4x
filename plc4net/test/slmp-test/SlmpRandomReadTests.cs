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
using System.Collections.Generic;
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
    public class SlmpRandomReadTests
    {
        [Fact]
        public async Task Multiple_scalar_word_devices_share_one_random_read_and_keep_request_order()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString("341278560000803F")));
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("real", "D1500:REAL")
                .AddTagAddress("word", "D0:WORD")
                .AddTagAddress("uint", "W0x100:UINT")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.Equal(new[] { "real", "word", "uint" }, response.TagNames);
            Assert.Equal(1.0f, response.GetValue("real")!.GetFloat());
            Assert.Equal((ushort)0x1234, response.GetValue("word")!.GetUshort());
            Assert.Equal((ushort)0x5678, response.GetValue("uint")!.GetUshort());
            Assert.Equal(Convert.FromHexString(
                "500000FFFF030014001000030400000201000000A8000100B4DC0500A8"),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Random_read_and_batch_fallback_can_share_one_api_request()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var randomRead = transport.WaitForWrittenData(25, 5000);
                transport.InjectTestData(Response(0, Convert.FromHexString("34127856")));
                var batchRead = transport.WaitForWrittenData(21, 5000);
                transport.InjectTestData(Response(0, Convert.FromHexString("10")));
                return randomRead.Concat(batchRead).ToArray();
            });
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("first", "D0:WORD")
                .AddTagAddress("bit", "M100:BOOL")
                .AddTagAddress("second", "W0x100:UINT")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.Equal((ushort)0x1234, response.GetValue("first")!.GetUshort());
            Assert.True(response.GetValue("bit")!.GetBool());
            Assert.Equal((ushort)0x5678, response.GetValue("second")!.GetUshort());
            Assert.Equal(Convert.FromHexString(
                "500000FFFF030010001000030400000200000000A8000100B4" +
                "500000FFFF03000C00100001040100640000900100"),
                await frames);
        }

        [Fact]
        public async Task A_single_random_read_candidate_keeps_batch_request_order()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var first = transport.WaitForWrittenData(21, 5000);
                transport.InjectTestData(Response(0, Convert.FromHexString("3412")));
                var second = transport.WaitForWrittenData(21, 5000);
                transport.InjectTestData(Response(0, Convert.FromHexString("10")));
                return first.Concat(second).ToArray();
            });
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("word", "D0:WORD")
                .AddTagAddress("bit", "M0:BOOL")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.Equal((ushort)0x1234, response.GetValue("word")!.GetUshort());
            Assert.True(response.GetValue("bit")!.GetBool());
            Assert.Equal(Convert.FromHexString(
                "500000FFFF03000C00100001040000000000A80100" +
                "500000FFFF03000C00100001040100000000900100"),
                await frames);
        }

        [Theory]
        [InlineData(0x1234, "", PlcResponseCode.RemoteError)]
        [InlineData(0, "3412", PlcResponseCode.InvalidData)]
        public async Task Random_read_applies_frame_errors_to_every_grouped_tag(
            ushort endCode, string dataHex, PlcResponseCode expected)
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(endCode, Convert.FromHexString(dataHex)));
            var request = connection.ReadRequestBuilder!
                .AddTagAddress("first", "D0:WORD")
                .AddTagAddress("second", "D1:WORD")
                .Build();

            var response = Assert.IsType<DefaultPlcReadResponse>(await request.ExecuteAsync());

            Assert.Equal(expected, response.GetResponseCode("first"));
            Assert.Equal(expected, response.GetResponseCode("second"));
            Assert.Null(response.GetValue("first"));
            Assert.Null(response.GetValue("second"));
        }

        [Fact]
        public async Task More_than_192_mixed_points_split_into_random_read_then_batch_read()
        {
            using var transport = NewTransport(4096);
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, new byte[191 * 2 + 4]));
            var firstFrame = Task.Run(() =>
            {
                var written = transport.WaitForWrittenData(785, 5000);
                transport.InjectTestData(Response(0, new byte[4]));
                return written;
            });
            var builder = connection.ReadRequestBuilder!;
            for (var i = 0; i < 191; i++)
                builder.AddTagAddress("d" + i, "D" + i + ":WORD");
            builder.AddTagAddress("real1", "D1000:REAL");
            builder.AddTagAddress("real2", "D1002:REAL");

            var response = Assert.IsType<DefaultPlcReadResponse>(await builder.Build().ExecuteAsync());

            Assert.Equal(193, response.TagNames.Count());
            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = (await firstFrame).Concat(transport.GetAllWrittenData()).ToArray();
            Assert.Equal(0x0403, written[11] | written[12] << 8);
            Assert.Equal(191, written[15]);
            Assert.Equal(1, written[16]);
            var secondFrame = 9 + (written[7] | written[8] << 8);
            Assert.Equal(0x0401, written[secondFrame + 11] | written[secondFrame + 12] << 8);
            Assert.Equal(Convert.FromHexString("EA0300A80200"),
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
