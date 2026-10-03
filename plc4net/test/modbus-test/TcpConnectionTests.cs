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
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The Modbus TCP connection against a scripted transport: what goes on the wire
    /// for each kind of tag, and how each kind of answer becomes a response code.
    /// The tag address counts from 1; the request carries it minus 1.
    /// </summary>
    public class TcpConnectionTests
    {
        // The first request on a fresh connection carries transaction identifier 1.

        [Theory]
        [InlineData("holding-register:1", 0x03, 0x00, 0x00)]
        [InlineData("holding-register:259", 0x03, 0x01, 0x02)]
        [InlineData("holding-register:65536", 0x03, 0xFF, 0xFF)]
        [InlineData("input-register:8", 0x04, 0x00, 0x07)]
        public async Task A_register_read_asks_for_one_register_at_the_wire_address_and_returns_its_value(
            string tag, byte function, byte addressHigh, byte addressLow)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, function, 0x02, 0x12, 0x34));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, tag);

            Assert.Equal((short)0x1234, ValueOf(response).GetShort());
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, function, addressHigh, addressLow, 0x00, 0x01),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_register_is_a_signed_INT_unless_the_tag_says_otherwise()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x03, 0x02, 0xFF, 0xFE));
            var connection = NewTcp(transport);

            var value = ValueOf(await ReadAsync(connection, "holding-register:1"));

            Assert.IsType<PlcINT>(value);
            Assert.Equal((short)-2, value.GetShort());
        }

        [Theory]
        [InlineData("holding-register:1:INT", typeof(PlcINT))]
        [InlineData("holding-register:1:UINT", typeof(PlcUINT))]
        [InlineData("holding-register:1:WORD", typeof(PlcWORD))]
        [InlineData("holding-register:1:uint", typeof(PlcUINT))]
        public async Task The_data_type_of_the_tag_decides_the_type_of_the_value(string tag, Type valueType)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x03, 0x02, 0xBE, 0xEF));
            var connection = NewTcp(transport);

            var value = ValueOf(await ReadAsync(connection, tag));

            Assert.IsType(valueType, value);
            if (valueType != typeof(PlcINT))
            {
                Assert.Equal((ushort)0xBEEF, value.GetUshort());
            }
        }

        [Theory]
        [InlineData("coil:4", 0x01, 0x01, true)]
        [InlineData("coil:4", 0x01, 0x00, false)]
        [InlineData("discrete-input:4", 0x02, 0x01, true)]
        [InlineData("discrete-input:4", 0x02, 0x00, false)]
        public async Task A_bit_read_asks_for_one_bit_and_returns_it(
            string tag, byte function, byte data, bool expected)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, function, 0x01, data));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, tag);

            Assert.IsType<PlcBOOL>(ValueOf(response));
            Assert.Equal(expected, ValueOf(response).GetBool());
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, function, 0x00, 0x03, 0x00, 0x01),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_device_that_answers_with_more_data_than_asked_for_still_yields_the_first_element()
        {
            // Some devices answer a one-register read with every register they have.
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1,
                0x03, 0x08, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:1");

            Assert.Equal((short)1, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task A_write_single_register_sends_the_value_and_succeeds_on_the_echo()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x06, 0x00, 0x05, 0x12, 0x34));
            var connection = NewTcp(transport);

            var response = await WriteAsync<ushort>(connection, "holding-register:6:UINT", 0x1234);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, 0x06, 0x00, 0x05, 0x12, 0x34),
                transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("holding-register:3", -2, 0xFF, 0xFE)] // INT: two's complement
        [InlineData("holding-register:3:INT", -32768, 0x80, 0x00)]
        [InlineData("holding-register:3:INT", 32767, 0x7F, 0xFF)]
        [InlineData("holding-register:3:UINT", 65535, 0xFF, 0xFF)]
        [InlineData("holding-register:3:WORD", 43981, 0xAB, 0xCD)]
        public async Task A_register_value_is_encoded_for_the_data_type_of_the_tag(
            string tag, int value, byte high, byte low)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x06, 0x00, 0x02, high, low));
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, tag, value);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, 0x06, 0x00, 0x02, high, low),
                transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData(true, 0xFF)]
        [InlineData(false, 0x00)]
        public async Task A_write_single_coil_sends_FF00_for_on_and_0000_for_off(bool on, byte high)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x05, 0x00, 0x0A, high, 0x00));
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, "coil:11", on);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Tcp(1, 1, 0x05, 0x00, 0x0A, high, 0x00),
                transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("discrete-input:2")]
        [InlineData("input-register:2")]
        public async Task Read_only_tags_cannot_be_written_and_nothing_goes_on_the_wire(string tag)
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, tag, 1);

            Assert.Equal(PlcResponseCode.AccessDenied, response.GetResponseCode("v"));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Theory]
        [InlineData("holding-register:2", 40000)]      // too big for an INT
        [InlineData("holding-register:2", -32769)]     // too small for an INT
        [InlineData("holding-register:2:UINT", -1)]    // a UINT is not negative
        [InlineData("holding-register:2:UINT", 65536)]
        [InlineData("holding-register:2:WORD", -1)]
        public async Task A_value_outside_the_range_of_the_tag_is_an_invalid_datatype_and_nothing_goes_on_the_wire(
            string tag, int value)
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, tag, value);

            Assert.Equal(PlcResponseCode.InvalidDatatype, response.GetResponseCode("v"));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public async Task A_value_that_is_not_a_number_is_an_invalid_datatype_and_nothing_goes_on_the_wire()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);

            var text = await WriteAsync(connection, "holding-register:2", "banana");
            var fraction = await WriteAsync(connection, "holding-register:2", 1.5);
            var flag = await WriteAsync(connection, "holding-register:2", true);
            var notBool = await WriteAsync(connection, "coil:2", "banana");

            Assert.All(new[] { text, fraction, flag, notBool },
                response => Assert.Equal(PlcResponseCode.InvalidDatatype, response.GetResponseCode("v")));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public async Task A_whole_number_in_a_floating_point_type_is_accepted()
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x06, 0x00, 0x01, 0x00, 0x07));
            var connection = NewTcp(transport);

            var response = await WriteAsync(connection, "holding-register:2", 7.0);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
        }

        [Theory]
        [InlineData(0x01, PlcResponseCode.Unsupported)]
        [InlineData(0x02, PlcResponseCode.InvalidAddress)]
        [InlineData(0x03, PlcResponseCode.InvalidDatatype)]
        [InlineData(0x04, PlcResponseCode.InternalError)]
        [InlineData(0x05, PlcResponseCode.Ok)]
        [InlineData(0x06, PlcResponseCode.RequestTimeout)]
        [InlineData(0x07, PlcResponseCode.InternalError)]
        [InlineData(0x08, PlcResponseCode.InternalError)]
        [InlineData(0x09, PlcResponseCode.InternalError)] // not a defined exception code
        [InlineData(0x0A, PlcResponseCode.InternalError)]
        [InlineData(0x0B, PlcResponseCode.InternalError)]
        [InlineData(0x0C, PlcResponseCode.InternalError)] // not a defined exception code
        public async Task A_device_exception_maps_to_a_specific_code(byte exceptionCode, PlcResponseCode expected)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x83, exceptionCode));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:100");

            Assert.Equal(expected, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task The_connection_survives_a_device_exception()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);

            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x83, 0x02));
            var first = await ReadAsync(connection, "holding-register:100");
            // A device only answers a request that has been sent.
            transport.InjectTestData(ModbusFrames.Tcp(2, 1, 0x03, 0x02, 0xBE, 0xEF));
            var second = await ReadAsync(connection, "holding-register:1:UINT");

            Assert.Equal(PlcResponseCode.InvalidAddress, first.GetResponseCode("v"));
            Assert.Equal((ushort)0xBEEF, ValueOf(second).GetUshort());
        }

        [Fact]
        public async Task An_answer_of_the_wrong_kind_is_an_internal_error()
        {
            var transport = NewTransport();
            // A write echo where a register read was asked for.
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x06, 0x00, 0x05, 0x12, 0x34));
            var connection = NewTcp(transport);

            var response = await ReadAsync(connection, "holding-register:6");

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_well_framed_answer_that_does_not_parse_is_an_internal_error_without_waiting()
        {
            var transport = NewTransport();
            // A response with a function code that is not Modbus.
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x7F, 0x00, 0x00));
            var connection = NewTcp(transport, "?request-timeout-ms=5000");
            var clock = System.Diagnostics.Stopwatch.StartNew();

            var response = await ReadAsync(connection, "holding-register:1");

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
            Assert.True(clock.Elapsed < TimeSpan.FromSeconds(2), $"took {clock.Elapsed}");
        }

        [Fact]
        public async Task No_answer_within_the_timeout_is_a_request_timeout()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport, "?request-timeout-ms=100");

            var response = await ReadAsync(connection, "holding-register:1");

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_stale_response_from_a_timed_out_read_does_not_brick_the_connection()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport, "?request-timeout-ms=600");

            var first = await ReadAsync(connection, "holding-register:1"); // transaction 1, never answered
            Assert.Equal(PlcResponseCode.RequestTimeout, first.GetResponseCode("v"));
            transport.GetAllWrittenData(); // the first request, so that the helper below waits for the second

            // Transaction 1's answer turns up late. The answer to transaction 2 comes later
            // still, so the stale one is seen on its own first and must not be taken for it.
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x03, 0x02, 0x11, 0x11));
            var answer = InjectAfterRequest(transport, 150, ModbusFrames.Tcp(2, 1, 0x03, 0x02, 0x22, 0x22));
            var second = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)0x2222, ValueOf(second).GetShort());
        }

        [Fact]
        public async Task A_response_that_arrives_in_pieces_is_assembled()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);
            var frame = ModbusFrames.Tcp(1, 1, 0x03, 0x02, 0x00, 0x42);

            var read = ReadAsync(connection, "holding-register:1");
            transport.InjectTestData(frame[..4]);
            await Task.Delay(30);
            transport.InjectTestData(frame[4..9]);
            await Task.Delay(30);
            transport.InjectTestData(frame[9..]);

            Assert.Equal((short)0x42, ValueOf(await read).GetShort());
        }

        [Fact]
        public async Task Concurrent_requests_are_written_one_after_the_other()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);

            // A device that answers each request as it arrives, echoing its transaction id.
            var requests = new List<byte[]>();
            var device = Task.Run(() =>
            {
                for (var n = 0; n < 2; n++)
                {
                    var request = transport.WaitForWrittenData(12, 5000);
                    requests.Add(request);
                    transport.InjectTestData(ModbusFrames.Tcp(
                        (ushort)((request[0] << 8) | request[1]), 1, 0x03, 0x02, 0x00, (byte)(0x0A + n)));
                }
            });

            var results = await Task.WhenAll(
                ReadAsync(connection, "holding-register:2"),
                ReadAsync(connection, "holding-register:3"));
            await device;

            // The requests were written whole and one after the other, so each got the
            // answer that was sent for it.
            Assert.Equal(2, requests.Count);
            Assert.All(requests, request => Assert.Equal(12, request.Length));
            Assert.Equal(1, requests[0][1]);
            Assert.Equal(2, requests[1][1]);
            var values = new[] { ValueOf(results[0]).GetShort(), ValueOf(results[1]).GetShort() };
            Assert.Equal(new short[] { 0x0A, 0x0B }, values.OrderBy(v => v));
        }

        [Fact]
        public async Task A_cancelled_request_throws_instead_of_becoming_a_tag_error()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);
            using var cancellation = new CancellationTokenSource();
            cancellation.Cancel();

            await Assert.ThrowsAnyAsync<OperationCanceledException>(
                () => ReadAsync(connection, "holding-register:1", cancellation.Token));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public async Task A_transport_that_has_been_closed_fails_the_read()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport);
            connection.Close();

            var response = await ReadAsync(connection, "holding-register:1");

            Assert.False(connection.IsConnected);
            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_connection_closed_while_waiting_fails_the_read_promptly()
        {
            var transport = NewTransport();
            var connection = NewTcp(transport, "?request-timeout-ms=5000");
            var clock = System.Diagnostics.Stopwatch.StartNew();

            var read = ReadAsync(connection, "holding-register:1");
            await Task.Delay(50);
            connection.Close();
            var response = await read;

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
            Assert.True(clock.Elapsed < TimeSpan.FromSeconds(2), $"took {clock.Elapsed}");
        }

        // ── Transaction identifiers ────────────────────────────

        [Theory]
        [InlineData(0, 1)]
        [InlineData(1, 2)]
        [InlineData(0xFFFD, 0xFFFE)]
        [InlineData(0xFFFE, 1)] // 0xFFFF is never used
        [InlineData(0xFFFF, 1)]
        public void Transaction_identifiers_run_from_1_to_FFFE_and_start_again(int last, int expected)
        {
            Assert.Equal((ushort)expected, ModbusTcpConnection.NextTransactionId((ushort)last));
        }

        // ── Connection string ──────────────────────────────────

        [Theory]
        [InlineData("?default-unit-identifier=7", 7)]
        [InlineData("?unit-identifier=7", 7)] // the older spelling
        [InlineData("?default-unit-identifier=9&unit-identifier=7", 9)] // plc4j's name wins
        [InlineData("", 1)]
        [InlineData("?default-unit-identifier=0", 0)]
        [InlineData("?default-unit-identifier=255", 255)]
        public async Task The_unit_identifier_comes_from_the_connection_string_and_is_sent(
            string parameters, int expected)
        {
            var transport = NewTransport();
            transport.InjectTestData(ModbusFrames.Tcp(1, (byte)expected, 0x03, 0x02, 0x00, 0x01));
            var connection = NewTcp(transport, parameters);

            await ReadAsync(connection, "holding-register:1");

            Assert.Equal((byte)expected, connection.UnitIdentifier);
            Assert.Equal(
                ModbusFrames.Tcp(1, (byte)expected, 0x03, 0x00, 0x00, 0x00, 0x01),
                transport.GetAllWrittenData());
        }

        [Theory]
        [InlineData("?default-unit-identifier=-1")]
        [InlineData("?default-unit-identifier=256")]
        [InlineData("?default-unit-identifier=seven")]
        [InlineData("?unit-identifier=")]
        [InlineData("?request-timeout-ms=0")]
        [InlineData("?request-timeout-ms=-5")]
        [InlineData("?request-timeout-ms=soon")]
        public void A_bad_parameter_is_rejected_instead_of_replaced_by_a_default(string parameters)
        {
            Assert.Throws<PlcConnectionException>(() => NewTcp(NewTransport(), parameters));
        }

        [Theory]
        [InlineData("?request-timeout-ms=100")]
        [InlineData("?request-timeout=100")] // the older spelling
        public async Task The_request_timeout_comes_from_the_connection_string(string parameters)
        {
            var transport = NewTransport();
            var connection = NewTcp(transport, parameters);
            var clock = System.Diagnostics.Stopwatch.StartNew();

            var response = await ReadAsync(connection, "holding-register:1");

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
            Assert.True(clock.Elapsed < TimeSpan.FromSeconds(3), $"took {clock.Elapsed}");
        }
    }
}