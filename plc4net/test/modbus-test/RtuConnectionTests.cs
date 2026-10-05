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
using System.Diagnostics;
using System.Threading.Tasks;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.modbus.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// The Modbus RTU connection against a scripted transport. The connection throws
    /// away whatever is buffered before it writes, so every answer is injected after
    /// the request has gone out.
    /// </summary>
    public class RtuConnectionTests
    {
        [Theory]
        [InlineData("holding-register:1", 0x03, 0x00, 0x00)]
        [InlineData("holding-register:259", 0x03, 0x01, 0x02)]
        [InlineData("input-register:8", 0x04, 0x00, 0x07)]
        public async Task A_register_read_asks_for_one_register_at_the_wire_address_and_returns_its_value(
            string tag, byte function, byte addressHigh, byte addressLow)
        {
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, function, 0x02, 0x00, 0x42));

            var response = await ReadAsync(connection, tag);
            await answer;

            Assert.Equal((short)0x42, ValueOf(response).GetShort());
            Assert.Equal(
                ModbusFrames.Rtu(1, function, addressHigh, addressLow, 0x00, 0x01),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task The_data_type_of_the_tag_decides_the_value()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x03, 0x02, 0xBE, 0xEF));

            var response = await ReadAsync(connection, "holding-register:1:UINT");
            await answer;

            Assert.Equal((ushort)0xBEEF, ValueOf(response).GetUshort());
        }

        [Theory]
        [InlineData("coil:4", 0x01, true)]
        [InlineData("discrete-input:4", 0x02, true)]
        [InlineData("coil:4", 0x01, false)]
        public async Task A_bit_read_asks_for_one_bit_and_returns_it(string tag, byte function, bool expected)
        {
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport,
                ModbusFrames.Rtu(1, function, 0x01, (byte)(expected ? 0x01 : 0x00)));

            var response = await ReadAsync(connection, tag);
            await answer;

            Assert.Equal(expected, ValueOf(response).GetBool());
            Assert.Equal(
                ModbusFrames.Rtu(1, function, 0x00, 0x03, 0x00, 0x01),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_device_that_answers_with_more_data_than_asked_for_still_yields_the_first_element()
        {
            // The Mitsubishi QJ71C24N the driver is verified against answers a
            // one-register read with four registers.
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1,
                0x03, 0x08, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)1, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task A_write_single_coil_succeeds_on_the_echo()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x05, 0x00, 0x0A, 0xFF, 0x00));

            var response = await WriteAsync(connection, "coil:11", true);
            await answer;

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Rtu(1, 0x05, 0x00, 0x0A, 0xFF, 0x00),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task A_write_single_register_succeeds_on_the_echo()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x06, 0x00, 0x05, 0x12, 0x34));

            var response = await WriteAsync<ushort>(connection, "holding-register:6:UINT", 0x1234);
            await answer;

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                ModbusFrames.Rtu(1, 0x06, 0x00, 0x05, 0x12, 0x34),
                transport.GetAllWrittenData());
        }

        [Fact]
        public async Task The_unit_identifier_is_the_slave_address_of_request_and_response()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?default-unit-identifier=17");
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(17, 0x03, 0x02, 0x00, 0x01));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)1, ValueOf(response).GetShort());
            Assert.Equal(17, transport.GetAllWrittenData()[0]);
        }

        [Theory]
        [InlineData(0x01, PlcResponseCode.Unsupported)]
        [InlineData(0x02, PlcResponseCode.InvalidAddress)]
        [InlineData(0x03, PlcResponseCode.InvalidDatatype)]
        [InlineData(0x04, PlcResponseCode.InternalError)]
        [InlineData(0x05, PlcResponseCode.Ok)]
        [InlineData(0x06, PlcResponseCode.RequestTimeout)]
        public async Task A_device_exception_maps_to_a_specific_code(byte exceptionCode, PlcResponseCode expected)
        {
            // The exception frame's CRC covers the function code with its error bit
            // set; the model has to keep that byte to check it.
            var transport = NewTransport();
            var connection = NewRtu(transport);
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x83, exceptionCode));

            var response = await ReadAsync(connection, "holding-register:10000");
            await answer;

            Assert.Equal(expected, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_half_duplex_echo_of_the_request_is_discarded_not_decoded_as_a_value()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            // Two-wire RS-485 adapters echo the transmitter: the request, then the answer.
            var echo = ModbusFrames.Rtu(1, 0x03, 0x00, 0x06, 0x00, 0x01);
            var answer = InjectAfterRequest(transport, echo, ModbusFrames.Rtu(1, 0x03, 0x02, 0xBE, 0xEF));

            var response = await ReadAsync(connection, "holding-register:7:UINT");
            await answer;

            Assert.Equal((ushort)0xBEEF, ValueOf(response).GetUshort());
        }

        [Fact]
        public async Task An_echo_that_arrives_before_the_answer_is_skipped_and_the_answer_still_counts()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            // The adapter echoes first, the device answers about 120 ms later, so the
            // codec chews through the echo before the answer is even on the wire.
            var echoTask = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x03, 0x00, 0x06, 0x00, 0x01));
            var answer = InjectAfterRequest(transport, 120, ModbusFrames.Rtu(1, 0x03, 0x02, 0xBE, 0xEF));

            var response = await ReadAsync(connection, "holding-register:7:UINT");
            await echoTask;
            await answer;

            Assert.Equal((ushort)0xBEEF, ValueOf(response).GetUshort());
        }

        // ── What is not an answer ──────────────────────────────

        [Fact]
        public async Task A_corrupted_answer_ends_as_a_timeout_and_the_next_request_still_works()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=300");
            var corrupted = ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42);
            corrupted[^1] ^= 0xFF;
            var firstAnswer = InjectAfterRequest(transport, corrupted);

            var first = await ReadAsync(connection, "holding-register:1");
            await firstAnswer;
            transport.GetAllWrittenData();

            var secondAnswer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x43));
            var second = await ReadAsync(connection, "holding-register:1");
            await secondAnswer;

            Assert.Equal(PlcResponseCode.RequestTimeout, first.GetResponseCode("v"));
            Assert.Equal((short)0x43, ValueOf(second).GetShort());
        }

        [Fact]
        public async Task A_corrupted_frame_does_not_hide_the_good_one_behind_it()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            var corrupted = ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42);
            corrupted[^1] ^= 0xFF;
            var answer = InjectAfterRequest(transport, corrupted, ModbusFrames.Rtu(1, 0x03, 0x02, 0x12, 0x34));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)0x1234, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task Noise_in_front_of_the_answer_is_skipped()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            var answer = InjectAfterRequest(transport,
                new byte[] { 0x55, 0x55, 0x55 },
                ModbusFrames.Rtu(1, 0x03, 0x02, 0x12, 0x34));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)0x1234, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task A_frame_from_another_slave_is_ignored_and_the_right_answer_still_counts()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            var answer = InjectAfterRequest(transport,
                ModbusFrames.Rtu(9, 0x03, 0x02, 0x99, 0x99),
                ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)0x42, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task Only_a_frame_from_another_slave_ends_as_a_timeout()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=200");
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(9, 0x03, 0x02, 0x00, 0x42));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_response_to_another_function_is_ignored_and_the_right_answer_still_counts()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            var answer = InjectAfterRequest(transport,
                ModbusFrames.Rtu(1, 0x04, 0x02, 0x99, 0x99), // an input-register answer
                ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)0x42, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task An_exception_for_another_function_is_ignored()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=200");
            var answer = InjectAfterRequest(transport, ModbusFrames.Rtu(1, 0x84, 0x02));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_frame_left_in_the_buffer_by_an_earlier_exchange_is_not_taken_for_the_answer()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=3000");
            transport.InjectTestData(ModbusFrames.Rtu(1, 0x03, 0x02, 0x11, 0x11)); // stale
            // The answer comes well after the request, so a stale frame that was not thrown
            // away would be seen, and taken for the answer, long before it.
            var answer = InjectAfterRequest(transport, 150, ModbusFrames.Rtu(1, 0x03, 0x02, 0x22, 0x22));

            var response = await ReadAsync(connection, "holding-register:1");
            await answer;

            Assert.Equal((short)0x2222, ValueOf(response).GetShort());
        }

        [Fact]
        public async Task No_answer_within_the_timeout_is_a_request_timeout()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=100");

            var response = await ReadAsync(connection, "holding-register:1");

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_connection_closed_while_waiting_fails_the_read_promptly()
        {
            var transport = NewTransport();
            var connection = NewRtu(transport, "?request-timeout-ms=5000");
            var clock = Stopwatch.StartNew();

            var read = ReadAsync(connection, "holding-register:1");
            await Task.Delay(100);
            connection.Close();
            var response = await read;

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
            Assert.True(clock.Elapsed < TimeSpan.FromSeconds(2), $"took {clock.Elapsed}");
        }

        [Theory]
        [InlineData("holding-register:2", 40000)]
        [InlineData("holding-register:2:UINT", -1)]
        public async Task A_value_outside_the_range_of_the_tag_is_an_invalid_datatype_and_nothing_goes_on_the_wire(
            string tag, int value)
        {
            var transport = NewTransport();
            var connection = NewRtu(transport);

            var response = await WriteAsync(connection, tag, value);

            Assert.Equal(PlcResponseCode.InvalidDatatype, response.GetResponseCode("v"));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        // ── Connection string ──────────────────────────────────

        [Theory]
        [InlineData("?default-unit-identifier=21", 21)]
        [InlineData("?unit-identifier=21", 21)] // the older spelling
        [InlineData("?default-unit-identifier=1&unit-identifier=21", 1)]
        [InlineData("?default-unit-identifier=247", 247)]
        public void The_unit_identifier_is_the_slave_address(string parameters, int expected)
        {
            Assert.Equal((byte)expected, NewRtu(NewTransport(), parameters).UnitIdentifier);
        }

        [Theory]
        [InlineData("?default-unit-identifier=0")]
        [InlineData("?default-unit-identifier=248")]
        [InlineData("?default-unit-identifier=-1")]
        [InlineData("?default-unit-identifier=one")]
        [InlineData("?request-timeout-ms=0")]
        public void A_bad_parameter_is_rejected_instead_of_replaced_by_a_default(string parameters)
        {
            Assert.Throws<PlcConnectionException>(() => NewRtu(NewTransport(), parameters));
        }
    }
}