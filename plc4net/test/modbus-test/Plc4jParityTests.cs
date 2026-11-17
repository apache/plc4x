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
using System.Threading.Tasks;
using org.apache.plc4net.spi.drivers.functions;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The requests and answers of plc4j's own driver test data,
    /// <c>protocols/modbus/src/test/resources/protocols/modbus/tcp/ManualFactoryModbusTCPDriverTest-testsuite.xml</c>,
    /// for the tags this driver supports. The same tag has to put the same PDU on the wire
    /// and decode the same answer; the bytes below are copied from that file (the transaction
    /// identifier and the unit identifier are the connection's own).
    /// </summary>
    public class Plc4jParityTests
    {
        private static byte[] Pdu(string hex) => Convert.FromHexString(hex);

        [Theory]
        // Read holding-register:3:WORD   (request 000500000006 01 0300020001, answer 000500000005 01 0302A5B8)
        [InlineData("holding-register:3:WORD", "0300020001", "0302A5B8", 42424)]
        // Read holding-register:12:INT   (-2424)
        [InlineData("holding-register:12:INT", "03000B0001", "0302F688", -2424)]
        // Read holding-register:13:UINT
        [InlineData("holding-register:13:UINT", "03000C0001", "0302A5B8", 42424)]
        public async Task A_read_puts_the_same_pdu_on_the_wire_and_decodes_the_same_answer(
            string tag, string requestPdu, string answerPdu, int expected)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(answerPdu)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, tag);

            var value = ValueOf(response);
            var actual = expected < 0 ? value.GetShort() : (long)value.GetUshort();
            Assert.Equal(expected, actual);
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)), transport.GetAllWrittenData());
        }

        [Theory]
        // Write holding-register:3:WORD  (request 00060000 0006 01 060002A5B8, the answer echoes it)
        [InlineData("holding-register:3:WORD", 42424, "060002A5B8")]
        // Write holding-register:12:INT
        [InlineData("holding-register:12:INT", -2424, "06000BF688")]
        // Write holding-register:13:UINT
        [InlineData("holding-register:13:UINT", 42424, "06000CA5B8")]
        public async Task A_write_puts_the_same_pdu_on_the_wire(string tag, int value, string requestPdu)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)));
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, tag, value);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu(requestPdu)), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("BIG_ENDIAN", "11223344")]
        [InlineData("LITTLE_ENDIAN", "44332211")]
        [InlineData("BIG_ENDIAN_BYTE_SWAP", "22114433")]
        [InlineData("LITTLE_ENDIAN_BYTE_SWAP", "33441122")]
        public async Task A_wide_value_is_read_with_the_same_byte_orders_as_plc4j(
            string byteOrder, string wireValue)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0304" + wireValue)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(
                connection, $"holding-register:1:DWORD{{byte-order:'{byteOrder}'}}");

            Assert.Equal(0x11223344U, ValueOf(response).GetUint());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("BIG_ENDIAN", "11223344")]
        [InlineData("LITTLE_ENDIAN", "44332211")]
        [InlineData("BIG_ENDIAN_BYTE_SWAP", "22114433")]
        [InlineData("LITTLE_ENDIAN_BYTE_SWAP", "33441122")]
        public async Task A_wide_value_is_written_with_the_same_byte_orders_as_plc4j(
            string byteOrder, string wireValue)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("1000000002")));
            var connection = NewTcp(transport);

            var response = await WriteAsync(
                connection, $"holding-register:1:DWORD{{byte-order:'{byteOrder}'}}", 0x11223344U);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, Pdu("100000000204" + wireValue)), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_register_array_is_read_in_one_request()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0306000100020003")));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:1[0..2]:INT");

            Assert.Equal(new short[] { 1, 2, 3 }, ValueOf(response).GetList().ConvertAll(v => v.GetShort()));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000003")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_register_array_is_written_in_one_request()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("1000000003")));
            var connection = NewTcp(transport);

            var response = await WriteAsync(
                connection, "holding-register:1[0..2]:INT", new List<short> { 1, 2, 3 });

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("100000000306000100020003")), transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("LITTLE_ENDIAN", "4433221188776655")]
        [InlineData("LITTLE_ENDIAN_BYTE_SWAP", "3344112277885566")]
        public async Task Byte_order_is_applied_to_each_array_element(
            string byteOrder, string wireValue)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0308" + wireValue)));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection,
                $"holding-register:1[0..1]:DWORD{{byte-order:'{byteOrder}'}}");

            Assert.Equal(new uint[] { 0x11223344, 0x55667788 },
                ValueOf(response).GetList().ConvertAll(v => v.GetUint()));
        }

        [Fact]
        public async Task Public_builder_params_write_a_register_array()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("1000000003")));
            var connection = NewTcp(transport);
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag("v", "holding-register:1[0..2]:INT", (short)1, (short)2, (short)3);

            var response = (DefaultPlcWriteResponse)await ((PlcWriter)connection)
                .Write((DefaultPlcWriteRequest)builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("100000000306000100020003")),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Connection_default_byte_order_and_tag_unit_id_are_used()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 3, Pdu("030444332211")));
            var connection = NewTcp(transport, "?default-payload-byte-order=LITTLE_ENDIAN");

            var response = await ReadAsync(connection, "holding-register:1:DWORD{unit-id:3}");

            Assert.Equal(0x11223344U, ValueOf(response).GetUint());
            Assert.Equal(ModbusFrames.Tcp(1, 3, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Fact]
        public void Extended_register_keeps_a_wide_flat_address()
        {
            var tag = ModbusTag.Parse("extended-register:65536:INT");

            Assert.Equal(65536, tag.Address);
        }

        [Fact]
        public async Task A_125_register_write_is_split_at_the_123_register_pdu_limit()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("100000007B")));
            var secondResponse = Task.Run(async () =>
            {
                while (transport.GetNumBytesWritten() < 259) await Task.Delay(1);
                transport.InjectTestData(ModbusFrames.Tcp(2, 1, Pdu("10007B0002")));
            });
            var connection = NewTcp(transport);
            var values = new short[125];
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag("v", "holding-register:1[0..124]:INT", values);

            var response = (DefaultPlcWriteResponse)await ((PlcWriter)connection)
                .Write((DefaultPlcWriteRequest)builder.Build());
            await secondResponse;

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(276, transport.GetAllWrittenData().Length);
        }

        [Fact]
        public async Task A_2000_coil_write_is_split_at_the_1968_coil_pdu_limit()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("0F000007B0")));
            var secondResponse = Task.Run(async () =>
            {
                while (transport.GetNumBytesWritten() < 259) await Task.Delay(1);
                transport.InjectTestData(ModbusFrames.Tcp(2, 1, Pdu("0F07B00020")));
            });
            var connection = NewTcp(transport);
            var values = new bool[2000];
            var builder = connection.WriteRequestBuilder!;
            builder.AddTag("v", "coil:1[0..1999]:BOOL", values);

            var response = (DefaultPlcWriteResponse)await ((PlcWriter)connection)
                .Write((DefaultPlcWriteRequest)builder.Build());
            await secondResponse;

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(276, transport.GetAllWrittenData().Length);
        }

        [Fact]
        public async Task A_fixed_length_string_is_read_from_its_registers()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("030441424344")));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:1:STRING(4)");

            Assert.Equal("ABCD", ValueOf(response).GetString());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task An_extended_register_uses_read_file_record()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("140403061234")));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "extended-register:1:INT");

            Assert.Equal(0x1234, ValueOf(response).GetShort());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("140706000100010001")), transport.GetAllWrittenData());
        }

        [Fact]
        public async Task Adjacent_register_tags_are_read_as_one_block()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, Pdu("030400010002")));
            var connection = NewTcp(transport);
            var builder = connection.ReadRequestBuilder!;
            builder.AddTagAddress("a", "holding-register:1:INT");
            builder.AddTagAddress("b", "holding-register:2:INT");
            var request = (DefaultPlcReadRequest)builder.Build();

            var response = (DefaultPlcReadResponse)await ((PlcReader)connection).Read(request);

            Assert.Equal((short)1, response.GetValue("a")!.GetShort());
            Assert.Equal((short)2, response.GetValue("b")!.GetShort());
            Assert.Equal(ModbusFrames.Tcp(1, 1, Pdu("0300000002")), transport.GetAllWrittenData());
        }
    }
}
