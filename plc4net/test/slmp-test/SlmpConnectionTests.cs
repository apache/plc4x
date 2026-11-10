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
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.slmp;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.transports.test;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.test.slmp
{
    public class SlmpConnectionTests
    {
        [Fact]
        public async Task Read_sends_a_batch_read_and_returns_a_word_value()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString("3412")));

            var response = await Read(connection, "value", "D100:WORD");

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("value"));
            Assert.Equal((ushort)0x1234, response.GetValue("value")!.GetUshort());
            Assert.Equal(Convert.FromHexString("500000FFFF03000C00100001040000640000A80100"),
                transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("D100:WORD", "FFFF", "word", 65535)]
        [InlineData("D100:INT", "FFFF", "int", -1)]
        [InlineData("D100:UINT", "FFFF", "uint", 65535)]
        [InlineData("D100:DINT", "78563412", "dint", 0x12345678)]
        [InlineData("D100:UDINT", "FFFFFFFF", "udint", -1)]
        [InlineData("D100:REAL", "0000803F", "real", 1)]
        public async Task Read_decodes_each_supported_scalar_type(
            string address, string dataHex, string type, long expected)
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString(dataHex)));

            var response = await Read(connection, "value", address);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("value"));
            var value = response.GetValue("value")!;
            switch (type)
            {
                case "word":
                    Assert.IsType<PlcWORD>(value);
                    Assert.Equal((ushort)expected, value.GetUshort());
                    break;
                case "int":
                    Assert.IsType<PlcINT>(value);
                    Assert.Equal((short)expected, value.GetShort());
                    break;
                case "uint":
                    Assert.IsType<PlcUINT>(value);
                    Assert.Equal((ushort)expected, value.GetUshort());
                    break;
                case "dint":
                    Assert.IsType<PlcDINT>(value);
                    Assert.Equal((int)expected, value.GetInt());
                    break;
                case "udint":
                    Assert.IsType<PlcUDINT>(value);
                    Assert.Equal(uint.MaxValue, value.GetUint());
                    break;
                case "real":
                    Assert.IsType<PlcREAL>(value);
                    Assert.Equal((float)expected, value.GetFloat());
                    break;
                default:
                    throw new ArgumentOutOfRangeException(nameof(type));
            }
        }

        [Fact]
        public async Task Read_returns_a_list_for_an_explicit_range_including_one_element()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString("34127856")));

            var response = await Read(connection, "values", "D100[4..5]:WORD");

            var values = Assert.IsType<PlcList>(response.GetValue("values"));
            Assert.Equal(2, values.GetLength());
            Assert.Equal((ushort)0x1234, values.GetIndex(0).GetUshort());
            Assert.Equal((ushort)0x5678, values.GetIndex(1).GetUshort());
        }

        [Fact]
        public async Task Read_maps_a_device_end_code_to_remote_error()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0xC05B, Array.Empty<byte>()));

            var response = await Read(connection, "value", "D100:WORD");

            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("value"));
            Assert.Null(response.GetValue("value"));
        }

        [Fact]
        public async Task Read_rejects_a_response_with_the_wrong_data_length()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString("34127856")));

            var response = await Read(connection, "value", "D100:WORD");

            Assert.Equal(PlcResponseCode.InvalidData, response.GetResponseCode("value"));
        }

        [Fact]
        public async Task Read_maps_timeout_and_does_not_reuse_an_uncorrelatable_connection()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport, "request-timeout-ms=10");

            var response = await Read(connection, "value", "D100:WORD");

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("value"));
            var nextResponse = await Read(connection, "next", "D101:WORD");
            Assert.Equal(PlcResponseCode.InternalError, nextResponse.GetResponseCode("next"));
        }

        [Fact]
        public async Task Read_cancellation_after_sending_does_not_reuse_the_connection()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport, "request-timeout-ms=1000");
            using var cancellation = new CancellationTokenSource();
            var pendingRead = connection.Read(
                Request(connection, "value", "D100:WORD"), cancellation.Token);
            await Task.Run(() => transport.WaitForWrittenData(21, 1000));
            cancellation.Cancel();

            await Assert.ThrowsAnyAsync<OperationCanceledException>(async () =>
                await pendingRead);

            var response = await Read(connection, "next", "D101:WORD");
            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("next"));
        }

        [Fact]
        public async Task Write_sends_a_little_endian_batch_write()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Array.Empty<byte>()));

            var response = await Write(connection, "value", "D100:WORD", (ushort)0x1234);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("value"));
            Assert.Equal(Convert.FromHexString("500000FFFF03000E00100001140000640000A801003412"),
                transport.GetAllWrittenData());
        }

        [Fact]
        public void Write_capability_is_exposed_through_connection_metadata()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);

            Assert.True(connection.PlcConnectionMetadata.CanWrite);
        }

        [Theory]
        [InlineData("D100:WORD", (ushort)0x1234, "3412")]
        [InlineData("D100:INT", (short)-2, "FEFF")]
        [InlineData("D100:UINT", (ushort)0xFEDC, "DCFE")]
        [InlineData("D100:DINT", 0x12345678, "78563412")]
        [InlineData("D100:UDINT", 0xFEDCBA98u, "98BADCFE")]
        [InlineData("D100:REAL", 1.5f, "0000C03F")]
        public async Task Write_encodes_each_supported_scalar_type(string address, object value, string dataHex)
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Array.Empty<byte>()));

            var response = await Write(connection, "value", address, value);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("value"));
            Assert.EndsWith(dataHex, Convert.ToHexString(transport.GetAllWrittenData()));
        }

        [Fact]
        public async Task Write_encodes_an_explicit_range_from_a_plc_list()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Array.Empty<byte>()));
            var values = new PlcList(new System.Collections.Generic.List<org.apache.plc4net.api.value.IPlcValue>
            {
                new PlcWORD(0x1234), new PlcWORD(0x5678)
            });

            var response = await Write(connection, "values", "D100[4..5]:WORD", values);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("values"));
            Assert.EndsWith("34127856", Convert.ToHexString(transport.GetAllWrittenData()));
        }

        [Fact]
        public async Task Write_builder_accepts_multiple_typed_values_for_a_range()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Array.Empty<byte>()));
            var request = Assert.IsType<DefaultPlcWriteRequest>(connection.WriteRequestBuilder!
                .AddTag("values", "D100[4..5]:WORD", (short)0x1234, (short)0x5678)
                .Build());

            var response = Assert.IsType<DefaultPlcWriteResponse>(await connection.Write(request));

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("values"));
            Assert.EndsWith("34127856", Convert.ToHexString(transport.GetAllWrittenData()));
        }

        [Fact]
        public async Task Write_builder_preserves_a_single_value_explicit_range()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Array.Empty<byte>()));
            var request = Assert.IsType<DefaultPlcWriteRequest>(connection.WriteRequestBuilder!
                .AddTag("values", "D100[4..4]:WORD", (ushort)0x1234)
                .Build());

            var response = Assert.IsType<DefaultPlcWriteResponse>(await connection.Write(request));

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("values"));
            Assert.EndsWith("3412", Convert.ToHexString(transport.GetAllWrittenData()));
        }

        [Fact]
        public async Task Write_rejects_multiple_values_for_a_scalar_tag_without_sending()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var request = Assert.IsType<DefaultPlcWriteRequest>(connection.WriteRequestBuilder!
                .AddTag("values", "D100:WORD", (short)1, (short)2)
                .Build());

            var response = Assert.IsType<DefaultPlcWriteResponse>(await connection.Write(request));

            Assert.Equal(PlcResponseCode.InvalidData, response.GetResponseCode("values"));
            Assert.Empty(transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Write_rejects_a_range_with_the_wrong_value_count_without_sending()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            var values = new PlcList(new System.Collections.Generic.List<org.apache.plc4net.api.value.IPlcValue>
            {
                new PlcWORD(0x1234)
            });

            var response = await Write(connection, "values", "D100[4..5]:WORD", values);

            Assert.Equal(PlcResponseCode.InvalidData, response.GetResponseCode("values"));
            Assert.Empty(transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Write_rejects_a_value_with_an_incompatible_plc_datatype_without_sending()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);

            var response = await Write(connection, "value", "D100:WORD", new PlcSTRING("not a number"));

            Assert.Equal(PlcResponseCode.InvalidData, response.GetResponseCode("value"));
            Assert.Empty(transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Write_maps_a_device_end_code_to_remote_error()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0xC05B, Array.Empty<byte>()));

            var response = await Write(connection, "value", "D100:WORD", (ushort)1);

            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("value"));
        }

        [Fact]
        public async Task Write_maps_timeout_and_does_not_reuse_an_uncorrelatable_connection()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport, "request-timeout-ms=10");

            var response = await Write(connection, "value", "D100:WORD", (ushort)1);

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("value"));
            var nextResponse = await Read(connection, "next", "D101:WORD");
            Assert.Equal(PlcResponseCode.InternalError, nextResponse.GetResponseCode("next"));
        }

        [Fact]
        public async Task Write_rejects_a_success_response_with_unexpected_data()
        {
            using var transport = NewTransport();
            using var connection = NewConnection(transport);
            transport.InjectTestData(Response(0, Convert.FromHexString("3412")));

            var response = await Write(connection, "value", "D100:WORD", (ushort)1);

            Assert.Equal(PlcResponseCode.RemoteError, response.GetResponseCode("value"));
        }

        private static TestTransportInstance NewTransport()
            => new TestTransportInstance(new TestTransportConfiguration());

        private static SlmpConnection NewConnection(TestTransportInstance transport, string parameters = "")
            => new SlmpConnection(ConnectionString.Parse(
                "slmp://localhost" + (parameters.Length == 0 ? "" : "?" + parameters)), transport);

        private static DefaultPlcReadRequest Request(SlmpConnection connection, string name, string address)
            => (DefaultPlcReadRequest)connection.ReadRequestBuilder!
                .AddTagAddress(name, address).Build();

        private static async Task<DefaultPlcReadResponse> Read(
            SlmpConnection connection, string name, string address)
            => Assert.IsType<DefaultPlcReadResponse>(
                await connection.Read(Request(connection, name, address)));

        private static async Task<DefaultPlcWriteResponse> Write(
            SlmpConnection connection, string name, string address, object value)
        {
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag<object>(name, address, value);
            return Assert.IsType<DefaultPlcWriteResponse>(await connection.Write(
                Assert.IsType<DefaultPlcWriteRequest>(builder.Build())));
        }

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
