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
    public class SlmpMultiBlockWriteTests
    {
        [Fact]
        public async Task Word_and_bit_ranges_share_one_multi_block_write()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0));
            var bits = new bool[16];
            bits[0] = true;
            bits[15] = true;
            var request = connection.WriteRequestBuilder!
                .AddTag("words", "D0[0..1]:WORD", (ushort)0x1111, (ushort)0x2222)
                .AddTag("bits", "M0[0..15]:BOOL", bits.Cast<object>().ToArray())
                .Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            Assert.Equal(Convert.FromHexString(
                "500000FFFF03001A001000061400000101" +
                "000000A80200111122220000009001000180"),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Multi_block_remote_error_applies_to_every_grouped_tag()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0xC05B));
            var request = connection.WriteRequestBuilder!
                .AddTag("first", "D0[0..1]:WORD", (ushort)1, (ushort)2)
                .AddTag("second", "D10[0..1]:WORD", (ushort)3, (ushort)4)
                .Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("first"));
            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("second"));
            Assert.Equal(0x1406, CommandAt(transport.GetAllWrittenData(), 0));
        }

        [Fact]
        public async Task Scalar_write_keeps_its_position_between_multi_block_groups()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var firstGroup = transport.WaitForWrittenData(37, 5000);
                transport.InjectTestData(Response(0));
                var scalar = transport.WaitForWrittenData(23, 5000);
                transport.InjectTestData(Response(0));
                var secondGroup = transport.WaitForWrittenData(37, 5000);
                transport.InjectTestData(Response(0));
                return firstGroup.Concat(scalar).Concat(secondGroup).ToArray();
            });
            var request = connection.WriteRequestBuilder!
                .AddTag("first", "D0[0..1]:WORD", (ushort)1, (ushort)2)
                .AddTag("second", "D10[0..1]:WORD", (ushort)3, (ushort)4)
                .AddTag("scalar", "D20:WORD", (ushort)5)
                .AddTag("third", "D30[0..1]:WORD", (ushort)6, (ushort)7)
                .AddTag("fourth", "D40[0..1]:WORD", (ushort)8, (ushort)9)
                .Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(0x1406, CommandAt(written, 0));
            Assert.Equal(0x1401, CommandAt(written, 37));
            Assert.Equal(0x1406, CommandAt(written, 60));
        }

        [Fact]
        public async Task More_than_120_blocks_split_into_multi_block_then_batch_write()
        {
            using var transport = NewTransport(2048);
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var multiBlock = transport.WaitForWrittenData(1217, 5000);
                transport.InjectTestData(Response(0));
                var batch = transport.WaitForWrittenData(25, 5000);
                transport.InjectTestData(Response(0));
                return multiBlock.Concat(batch).ToArray();
            });
            var builder = connection.WriteRequestBuilder!;
            for (var i = 0; i < 121; i++)
                builder.AddTag("range" + i, "D" + (i * 2) + "[0..1]:WORD", (ushort)i, (ushort)(i + 1));

            var response = Assert.IsType<DefaultPlcWriteResponse>(await builder.Build().ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(0x1406, CommandAt(written, 0));
            Assert.Equal(120, written[15]);
            Assert.Equal(0x1401, CommandAt(written, 1217));
        }

        [Theory]
        [InlineData(476, 0x1406, 0)]
        [InlineData(477, 0x1401, 0x1401)]
        public async Task Cost_limit_splits_only_above_960_points(int secondPoints, int first, int second)
        {
            using var transport = NewTransport(4096);
            using var connection = NewConnection(transport);
            var frames = Task.Run(() =>
            {
                var firstFrame = transport.WaitForWrittenData(second == 0 ? 1933 : 973, 5000);
                transport.InjectTestData(Response(0));
                if (second == 0)
                    return firstFrame;
                var secondFrame = transport.WaitForWrittenData(21 + 2 * secondPoints, 5000);
                transport.InjectTestData(Response(0));
                return firstFrame.Concat(secondFrame).ToArray();
            });
            var builder = connection.WriteRequestBuilder!
                .AddTag("first", "D0[0..475]:WORD", Values(476))
                .AddTag("second", "D1000[0.." + (secondPoints - 1) + "]:WORD", Values(secondPoints));

            var response = Assert.IsType<DefaultPlcWriteResponse>(await builder.Build().ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(first, CommandAt(written, 0));
            if (second != 0)
                Assert.Equal(second, CommandAt(written, 973));
        }

        [Fact]
        public async Task Overlapping_ranges_keep_using_ordered_batch_writes()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var frames = Task.Run(() => CaptureFrames(transport, 25, 25, 25));
            var request = connection.WriteRequestBuilder!
                .AddTag("first", "D0[0..1]:WORD", (ushort)1, (ushort)2)
                .AddTag("overlap", "D1[0..1]:WORD", (ushort)3, (ushort)4)
                .AddTag("last", "D10[0..1]:WORD", (ushort)5, (ushort)6).Build();

            var response = Assert.IsType<DefaultPlcWriteResponse>(await request.ExecuteAsync());

            Assert.All(response.TagNames, name => Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode(name)));
            var written = await frames;
            Assert.Equal(new[] { 0x1401, 0x1401, 0x1401 },
                new[] { CommandAt(written, 0), CommandAt(written, 25), CommandAt(written, 50) });
        }

        [Fact]
        public async Task Seventeen_bits_span_two_little_endian_words()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0));
            var bits = new bool[17];
            bits[0] = bits[8] = bits[16] = true;
            var request = connection.WriteRequestBuilder!
                .AddTag("bits", "M0[0..16]:BOOL", bits.Cast<object>().ToArray())
                .AddTag("next", "M32[0..0]:BOOL", true).Build();

            await request.ExecuteAsync();

            Assert.Equal(Convert.FromHexString(
                "500000FFFF03001A001000061400000002" +
                "000000900200010101002000009001000100"), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Group_maps_unexpected_data_timeout_and_invalid_connection()
        {
            using (var transport = NewTransport())
            using (var connection = NewConnection(transport))
            {
                transport.InjectTestData(Response(0, 0x34));
                var response = await WriteTwoRanges(connection);
                Assert.All(response.TagNames,
                    name => Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode(name)));
            }
            using (var transport = NewTransport())
            using (var connection = NewConnection(transport, "request-timeout-ms=10"))
            {
                var timeout = await WriteTwoRanges(connection);
                Assert.All(timeout.TagNames,
                    name => Assert.Equal(PlcResponseCode.RequestTimeout, timeout.GetResponseCode(name)));
                var invalid = await WriteTwoRanges(connection);
                Assert.All(invalid.TagNames,
                    name => Assert.Equal(PlcResponseCode.InternalError, invalid.GetResponseCode(name)));
            }
        }

        private static int CommandAt(byte[] frame, int offset)
            => frame[offset + 11] | frame[offset + 12] << 8;

        private static object[] Values(int count)
            => Enumerable.Repeat<object>((ushort)1, count).ToArray();

        private static byte[] CaptureFrames(TestTransportInstance transport, params int[] lengths)
        {
            var written = Array.Empty<byte>();
            foreach (var length in lengths)
            {
                written = written.Concat(transport.WaitForWrittenData(length, 5000)).ToArray();
                transport.InjectTestData(Response(0));
            }
            return written;
        }

        private static async Task<DefaultPlcWriteResponse> WriteTwoRanges(SlmpConnection connection)
            => Assert.IsType<DefaultPlcWriteResponse>(await connection.WriteRequestBuilder!
                .AddTag("first", "D0[0..1]:WORD", (ushort)1, (ushort)2)
                .AddTag("second", "D10[0..1]:WORD", (ushort)3, (ushort)4)
                .Build().ExecuteAsync());

        private static TestTransportInstance NewTransport(int bufferSize = 1024)
            => new TestTransportInstance(new TestTransportConfiguration { ReceiveBufferSize = bufferSize });

        private static SlmpConnection NewConnection(TestTransportInstance transport, string parameters = "")
            => new SlmpConnection(ConnectionString.Parse(
                "slmp://localhost" + (parameters.Length == 0 ? "" : "?" + parameters)), transport);

        private static byte[] Response(ushort endCode, params byte[] data) => new byte[]
        {
            0xD0, 0x00, 0x00, 0xFF, 0xFF, 0x03, 0x00,
            (byte)(data.Length + 2), 0x00, (byte)endCode, (byte)(endCode >> 8)
        }.Concat(data).ToArray();
    }
}
