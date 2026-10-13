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
using System.Globalization;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.types;
using Xunit;
using static org.apache.plc4net.drivers.s7.test.ConnectionTestSupport;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>
    /// The S7 connection against a transport that plays the CPU: what goes on the wire for each kind
    /// of tag, and how each kind of answer becomes a value or a response code. The first request
    /// after the connection carries TPDU reference 2, since Setup Communication took 1.
    /// </summary>
    public class S7ConnectionTests
    {
        private const byte Inputs = 0x81;
        private const byte Outputs = 0x82;
        private const byte Flags = 0x83;
        private const byte DataBlocks = 0x84;

        private static byte[] Hex(string hex)
            => Enumerable.Range(0, hex.Length / 2)
                .Select(i => byte.Parse(hex.Substring(i * 2, 2), NumberStyles.HexNumber, CultureInfo.InvariantCulture))
                .ToArray();

        // ── Connect ────────────────────────────────────────────

        [Fact]
        public void Connect_sends_the_plc4j_TSAPs_in_the_connection_request()
        {
            // remote = PG_OR_PC (0x01), rack 0, slot 1 -> 0x0101; local = OTHERS (0x03), rack 1, slot 1 ->
            // 0x0311. A wrong shift, a byte swap or a renamed parameter would turn this red.
            var transport = Ready();

            Connect(transport);

            var cr = Frames(transport.GetAllWrittenData())[0];
            Assert.Equal(0xC1, cr[14]); // calling TSAP
            Assert.Equal(0x03, cr[16]);
            Assert.Equal(0x11, cr[17]);
            Assert.Equal(0xC2, cr[18]); // called TSAP
            Assert.Equal(0x01, cr[20]);
            Assert.Equal(0x01, cr[21]);
        }

        [Theory]
        [InlineData("s7://192.168.0.1?remote-rack=2&remote-slot=3", 0x01, 0x23)]
        [InlineData("s7://192.168.0.1?remote-rack=0&remote-slot=1&remote-device-group=OTHERS", 0x03, 0x01)]
        [InlineData("s7://192.168.0.1?remote-rack=0&remote-slot=1&remote-device-group=OS", 0x02, 0x01)]
        [InlineData("s7://192.168.0.1?remote-rack=0&remote-slot=1&remote-device-group=0x02", 0x02, 0x01)]
        [InlineData("s7://192.168.0.1?remote-rack=3&remote-slot=4&remote-tsap=0x0300", 0x03, 0x00)] // wins over rack and slot
        public void The_remote_TSAP_comes_from_the_rack_slot_and_device_group_or_from_the_parameter(
            string uri, int high, int low)
        {
            var transport = Ready();

            Connect(transport, uri);

            var cr = Frames(transport.GetAllWrittenData())[0];
            Assert.Equal(0xC2, cr[18]);
            Assert.Equal(high, cr[20]);
            Assert.Equal(low, cr[21]);
        }

        [Fact]
        public void The_local_TSAP_can_be_given_too()
        {
            var transport = Ready();

            Connect(transport, "s7://192.168.0.1?local-tsap=0x4D57");

            var cr = Frames(transport.GetAllWrittenData())[0];
            Assert.Equal(0xC1, cr[14]);
            Assert.Equal(0x4D, cr[16]);
            Assert.Equal(0x57, cr[17]);
        }

        [Fact]
        public void Connect_runs_Setup_Communication_and_records_the_PDU_length_the_CPU_grants()
        {
            var transport = Ready(240);

            var connection = Connect(transport);

            var frames = Frames(transport.GetAllWrittenData());
            Assert.Equal(2, frames.Count); // the Connection Request, then Setup Communication
            Assert.Equal(S7Frames.SetupRequest(1, 1024), frames[1]);
            Assert.Equal(240, connection.NegotiatedPduLength);
        }

        [Fact]
        public void The_PDU_length_that_is_asked_for_comes_from_the_connection_string()
        {
            var transport = Ready(480);

            var connection = Connect(transport, S7Uri + "&pdu-size=480");

            Assert.Equal(S7Frames.SetupRequest(1, 480), Frames(transport.GetAllWrittenData())[1]);
            Assert.Equal(480, connection.NegotiatedPduLength);
        }

        [Fact]
        public void Connect_fails_cleanly_when_the_CPU_does_not_answer_Setup_Communication()
        {
            var transport = NewTransport();
            transport.InjectTestData(S7Frames.ConnectionConfirm());

            var exception = Assert.Throws<PlcConnectionException>(() => Connect(transport, S7Uri + "&read-timeout-ms=100"));

            Assert.Contains("Setup Communication", exception.Message);
        }

        [Fact]
        public void Connect_reports_a_refusal_of_Setup_Communication_with_the_hint_for_PUT_GET()
        {
            var transport = NewTransport();
            transport.InjectTestData(S7Frames.ConnectionConfirm());
            transport.InjectTestData(S7Frames.Dt(S7Frames.Refusal(1, 0x81, 0x04)));

            var exception = Assert.Throws<PlcConnectionException>(() => Connect(transport));

            Assert.Contains("0x8104", exception.Message);
            Assert.Contains("PUT/GET", exception.Message);
        }

        [Fact]
        public void Connect_over_a_transport_that_is_not_cotp_is_rejected()
        {
            var connection = new S7Connection(ConnectionString.Parse(S7Uri), NewTransport());

            var exception = Assert.Throws<PlcConnectionException>(() => connection.Connect());

            Assert.Contains("cotp", exception.Message);
        }

        [Theory]
        [InlineData("remote-rack=16")]
        [InlineData("remote-rack=-1")]
        [InlineData("remote-slot=16")]
        [InlineData("local-rack=x")]
        [InlineData("local-slot=1a")]
        [InlineData("remote-tsap=zz")]
        [InlineData("remote-tsap=0x10000")]
        [InlineData("local-tsap=-5")]
        [InlineData("remote-device-group=BANANA")]
        [InlineData("local-device-group=256")]
        [InlineData("pdu-size=0")]
        [InlineData("pdu-size=70000")]
        [InlineData("max-amq-caller=0")]
        [InlineData("max-amq-callee=many")]
        [InlineData("read-timeout-ms=0")]
        [InlineData("read-timeout-ms=soon")]
        public void A_bad_parameter_is_rejected_instead_of_replaced_by_a_default(string parameter)
        {
            Assert.Throws<PlcConnectionException>(
                () => new S7Connection(ConnectionString.Parse("s7://192.168.0.1?" + parameter), NewTransport()));
        }

        [Theory]
        [InlineData("remote-rack=15&remote-slot=15")]
        [InlineData("remote-device-group=PG_OR_PC&local-device-group=OTHERS")]
        [InlineData("remote-tsap=0x0100&local-tsap=0x0200")]
        [InlineData("pdu-size=480&max-amq-caller=4&max-amq-callee=4&read-timeout-ms=2000")]
        public void The_parameters_that_plc4j_knows_are_accepted(string parameters)
        {
            Assert.NotNull(new S7Connection(ConnectionString.Parse("s7://192.168.0.1?" + parameters), NewTransport()));
        }

        // ── Reading ────────────────────────────────────────────

        [Theory]
        [InlineData("%DB1.DBX0.0:BOOL", 0x01, 1, DataBlocks, 0, 0, 1, 0x03)]
        [InlineData("%M0.0:BOOL", 0x01, 0, Flags, 0, 0, 1, 0x03)]
        [InlineData("%I0.3:BOOL", 0x01, 0, Inputs, 0, 3, 1, 0x03)]
        [InlineData("%Q7.7:BOOL", 0x01, 0, Outputs, 7, 7, 1, 0x03)]
        [InlineData("%DB1.DBB1:BYTE", 0x02, 1, DataBlocks, 1, 0, 1, 0x04)]
        [InlineData("%DB1.DBB1:SINT", 0x02, 1, DataBlocks, 1, 0, 1, 0x04)]
        [InlineData("%DB1.DBB1:USINT", 0x02, 1, DataBlocks, 1, 0, 1, 0x04)]
        [InlineData("%DB1.DBW2:WORD", 0x04, 1, DataBlocks, 2, 0, 2, 0x04)]
        [InlineData("%DB1.DBW10:INT", 0x05, 1, DataBlocks, 10, 0, 2, 0x05)]
        [InlineData("%DB1.DBW2:UINT", 0x05, 1, DataBlocks, 2, 0, 2, 0x05)]
        [InlineData("%DB1.DBD4:DWORD", 0x06, 1, DataBlocks, 4, 0, 4, 0x04)]
        [InlineData("%DB1.DBD4:DINT", 0x07, 1, DataBlocks, 4, 0, 4, 0x04)]
        [InlineData("%DB1.DBD4:UDINT", 0x07, 1, DataBlocks, 4, 0, 4, 0x04)]
        [InlineData("%DB5.DBD8:REAL", 0x08, 5, DataBlocks, 8, 0, 4, 0x07)]
        [InlineData("%MW30:WORD", 0x04, 0, Flags, 30, 0, 2, 0x04)]
        [InlineData("%QD4:DWORD", 0x06, 0, Outputs, 4, 0, 4, 0x04)]
        [InlineData("%IB2:BYTE", 0x02, 0, Inputs, 2, 0, 1, 0x04)]
        public async Task A_read_asks_for_the_address_with_the_transport_size_of_the_type(
            string address, byte transportSize, ushort block, byte area, int byteOffset, int bitOffset,
            int dataBytes, byte answerSize)
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, answerSize, new byte[dataBytes])));

            var response = await ReadAsync(connection, address);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                S7Frames.ReadRequest(2, S7Frames.AddressItem(transportSize, block, area, byteOffset, bitOffset)),
                Written(transport));
        }

        [Theory]
        [InlineData("%DB1.DBX0.0:BOOL", 0x03, "01", "PlcBOOL", "True")]
        [InlineData("%DB1.DBX0.1:BOOL", 0x03, "00", "PlcBOOL", "False")]
        [InlineData("%DB1.DBB1:BYTE", 0x04, "A5", "PlcBYTE", "165")]
        [InlineData("%DB1.DBB1:SINT", 0x04, "FE", "PlcSINT", "-2")]
        [InlineData("%DB1.DBB1:USINT", 0x04, "FE", "PlcUSINT", "254")]
        [InlineData("%DB1.DBW2:WORD", 0x04, "BEEF", "PlcWORD", "48879")]
        [InlineData("%DB1.DBW2:INT", 0x05, "FFFE", "PlcINT", "-2")]
        [InlineData("%DB1.DBW2:INT", 0x04, "FFFE", "PlcINT", "-2")] // the CPU may answer an INT as a word
        [InlineData("%DB1.DBW2:UINT", 0x05, "FFFE", "PlcUINT", "65534")]
        [InlineData("%DB1.DBD4:DWORD", 0x04, "DEADBEEF", "PlcDWORD", "3735928559")]
        [InlineData("%DB1.DBD4:DINT", 0x04, "FFFFFF38", "PlcDINT", "-200")]
        [InlineData("%DB1.DBD4:UDINT", 0x04, "FFFFFF38", "PlcUDINT", "4294967096")]
        public async Task A_value_is_decoded_as_the_type_of_the_tag(
            string address, byte answerSize, string answerHex, string valueType, string expected)
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, answerSize, Hex(answerHex))));

            var value = ValueOf(await ReadAsync(connection, address));

            Assert.Equal(valueType, value.GetType().Name);
            Assert.Equal(expected, value is org.apache.plc4net.spi.model.values.PlcBOOL
                ? value.GetBool().ToString()
                : value.GetLong().ToString(CultureInfo.InvariantCulture));
        }

        [Fact]
        public async Task A_REAL_is_decoded_as_a_single_precision_float()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, 0x07, Hex("40490FD0"))));

            var value = ValueOf(await ReadAsync(connection, "%DB5.DBD8:REAL"));

            Assert.Equal("PlcREAL", value.GetType().Name);
            Assert.Equal(3.14159, value.GetFloat(), 5);
        }

        [Fact]
        public async Task A_multi_item_read_is_decoded_item_by_item_in_the_order_of_the_request()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(
                2,
                (0xFF, 0x03, Hex("01")), // an odd item that is not the last: the CPU pads it
                (0xFF, 0x05, Hex("002A")),
                (0xFF, 0x07, Hex("40490FD0"))));

            var response = await ReadAsync(
                connection, ("flag", "%M0.0:BOOL"), ("count", "%DB1.DBW10:INT"), ("setpoint", "%DB1.DBD12:REAL"));

            Assert.True(ValueOf(response, "flag").GetBool());
            Assert.Equal((short)42, ValueOf(response, "count").GetShort());
            Assert.Equal(3.14159, ValueOf(response, "setpoint").GetFloat(), 5);
            Assert.Equal(
                S7Frames.ReadRequest(
                    2,
                    S7Frames.AddressItem(0x01, 0, Flags, 0, 0),
                    S7Frames.AddressItem(0x05, 1, DataBlocks, 10, 0),
                    S7Frames.AddressItem(0x08, 1, DataBlocks, 12, 0)),
                Written(transport));
        }

        [Theory]
        [InlineData(0x0A, PlcResponseCode.NotFound)]
        [InlineData(0x05, PlcResponseCode.InvalidAddress)]
        [InlineData(0x03, PlcResponseCode.AccessDenied)]
        [InlineData(0x06, PlcResponseCode.InvalidDatatype)]
        [InlineData(0x01, PlcResponseCode.InternalError)]
        [InlineData(0x07, PlcResponseCode.InternalError)]
        public async Task An_item_error_of_the_CPU_becomes_a_response_code_for_that_tag_only(
            byte returnCode, PlcResponseCode expected)
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (returnCode, 0x00, Array.Empty<byte>()), (0xFF, 0x05, Hex("0007"))));

            var response = await ReadAsync(connection, ("bad", "%DB99.DBW0:INT"), ("good", "%DB1.DBW0:INT"));

            Assert.Equal(expected, response.GetResponseCode("bad"));
            Assert.Equal((short)7, ValueOf(response, "good").GetShort());
        }

        [Fact]
        public async Task A_refusal_of_the_whole_job_is_access_denied_for_every_tag()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.Dt(S7Frames.Refusal(2, 0x81, 0x04)));

            var response = await ReadAsync(connection, ("a", "%DB1.DBW0:INT"), ("b", "%DB1.DBW2:INT"));

            Assert.Equal(PlcResponseCode.AccessDenied, response.GetResponseCode("a"));
            Assert.Equal(PlcResponseCode.AccessDenied, response.GetResponseCode("b"));
        }

        [Fact]
        public async Task A_bare_acknowledgement_does_not_shift_the_answer_to_the_next_request()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.Dt(S7Frames.Refusal(2, 0x85, 0x00)));

            var refused = await ReadAsync(connection, "%DB1.DBW0:INT");
            transport.InjectTestData(S7Frames.ReadAck(3, (0xFF, 0x05, Hex("BEEF"))));
            var read = await ReadAsync(connection, "%DB1.DBW0:INT");

            Assert.Equal(PlcResponseCode.AccessDenied, refused.GetResponseCode("v"));
            Assert.Equal(unchecked((short)0xBEEF), ValueOf(read).GetShort());
        }

        [Fact]
        public async Task An_answer_with_another_TPDU_reference_is_not_the_answer()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            // What a request that timed out earlier got as an answer is there first; the answer to
            // this one comes later still, so the stale one is seen on its own and must not be taken.
            transport.InjectTestData(S7Frames.ReadAck(99, (0xFF, 0x05, Hex("1111"))));
            var cpu = new FakeCpu(transport, 100, S7Frames.ReadAck(2, (0xFF, 0x05, Hex("2222"))));

            var value = ValueOf(await ReadAsync(connection, "%DB1.DBW0:INT"));
            await cpu.Running;

            Assert.Equal((short)0x2222, value.GetShort());
        }

        [Fact]
        public async Task No_answer_within_the_timeout_is_a_request_timeout_and_the_late_answer_is_not_taken_for_the_next()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport, S7Uri + "&read-timeout-ms=300");

            var first = await ReadAsync(connection, "%DB1.DBW0:INT"); // reference 2, never answered
            Assert.Equal(PlcResponseCode.RequestTimeout, first.GetResponseCode("v"));
            transport.GetAllWrittenData();

            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, 0x05, Hex("1111")))); // the late answer
            var cpu = new FakeCpu(transport, 100, S7Frames.ReadAck(3, (0xFF, 0x05, Hex("2222"))));
            var second = await ReadAsync(connection, "%DB1.DBW0:INT");
            await cpu.Running;

            Assert.Equal((short)0x2222, ValueOf(second).GetShort());
        }

        [Fact]
        public async Task An_answer_that_arrives_in_pieces_is_assembled()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            var answer = S7Frames.ReadAck(2, (0xFF, 0x05, Hex("0042")));

            var read = ReadAsync(connection, "%DB1.DBW0:INT");
            transport.InjectTestData(answer[..9]);
            await Task.Delay(30);
            transport.InjectTestData(answer[9..20]);
            await Task.Delay(30);
            transport.InjectTestData(answer[20..]);

            Assert.Equal((short)0x42, ValueOf(await read).GetShort());
        }

        [Fact]
        public async Task An_answer_with_fewer_items_than_asked_for_leaves_the_missing_ones_as_an_internal_error()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, 0x05, Hex("0001"))));

            var response = await ReadAsync(connection, ("a", "%DB1.DBW0:INT"), ("b", "%DB1.DBW2:INT"));

            Assert.Equal((short)1, ValueOf(response, "a").GetShort());
            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("b"));
        }

        [Fact]
        public async Task Data_that_is_shorter_than_the_type_is_an_internal_error_not_a_guess()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, 0x04, Hex("01"))));

            var response = await ReadAsync(connection, "%DB1.DBW0:INT");

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_CPU_that_answers_with_more_data_than_asked_for_still_yields_the_first_element()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.ReadAck(2, (0xFF, 0x04, Hex("0007FFFF"))));

            var value = ValueOf(await ReadAsync(connection, "%DB1.DBW0:INT"));

            Assert.Equal((short)7, value.GetShort());
        }

        [Fact]
        public async Task An_answer_to_a_write_is_not_the_answer_to_a_read()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.WriteAck(2, 0xFF));

            var response = await ReadAsync(connection, "%DB1.DBW0:INT");

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_closed_transport_fails_the_read_instead_of_hanging()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.Close();

            var response = await ReadAsync(connection, "%DB1.DBW0:INT");

            Assert.Equal(PlcResponseCode.InternalError, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_cancelled_read_throws_and_sends_nothing()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            var builder = connection.ReadRequestBuilder!;
            builder.AddTagAddress("v", "%DB1.DBW0:INT");
            var request = (org.apache.plc4net.spi.drivers.messages.DefaultPlcReadRequest)builder.Build();
            using var cancelled = new CancellationTokenSource();
            cancelled.Cancel();

            await Assert.ThrowsAnyAsync<OperationCanceledException>(
                () => connection.Read(request, cancelled.Token));

            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public async Task A_read_of_many_tags_is_cut_into_jobs_that_fit_the_PDU()
        {
            var transport = Ready(240);
            var connection = ConnectQuietly(transport);
            var cpu = new FakeCpu(
                transport,
                S7Frames.ReadAck(
                    2, Enumerable.Range(0, 19).Select(i => ((byte)0xFF, (byte)0x04, new[] { (byte)i })).ToArray()),
                S7Frames.ReadAck(
                    3, Enumerable.Range(19, 6).Select(i => ((byte)0xFF, (byte)0x04, new[] { (byte)i })).ToArray()));

            var response = await ReadAsync(
                connection, Enumerable.Range(0, 25).Select(i => ($"t{i}", $"%MB{i}:BYTE")).ToArray());
            await cpu.Running;

            for (var i = 0; i < 25; i++)
            {
                Assert.Equal((byte)i, ValueOf(response, $"t{i}").GetByte());
            }
            Assert.Equal(2, cpu.Requests.Count);
            Assert.Equal((byte)(2 + 12 * 19), cpu.Requests[0][7 + 7]); // the parameter length of the first job
            Assert.Equal((byte)(2 + 12 * 6), cpu.Requests[1][7 + 7]);
        }

        [Fact]
        public async Task A_job_that_fails_fails_its_own_tags_and_the_rest_go_on()
        {
            var transport = Ready(240);
            var connection = ConnectQuietly(transport, S7Uri + "&read-timeout-ms=100");
            // The first job (19 tags) gets no answer; the second (2 tags) does. A timeout throws away what
            // is buffered, so that answer must not be there before the second job has been sent.
            var cpu = new FakeCpu(
                transport, null, S7Frames.ReadAck(3, (0xFF, 0x04, Hex("07")), (0xFF, 0x04, Hex("08"))));

            var response = await ReadAsync(
                connection, Enumerable.Range(0, 21).Select(i => ($"t{i}", $"%MB{i}:BYTE")).ToArray());
            await cpu.Running;

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("t0"));
            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("t18"));
            Assert.Equal((byte)7, ValueOf(response, "t19").GetByte());
            Assert.Equal((byte)8, ValueOf(response, "t20").GetByte());
        }

        // ── Planning the jobs ──────────────────────────────────

        private static List<S7Tag> Tags(int count, string type = "BYTE", string area = "M")
            => Enumerable.Range(0, count)
                .Select(i => S7Tag.Parse($"%{area}{(type == "REAL" ? "D" : type == "BYTE" ? "B" : "W")}{i * 4}:{type}"))
                .ToList();

        [Theory]
        [InlineData(25, 240, "19,6")] // 12 + 12 * 19 = 240
        [InlineData(19, 240, "19")]
        [InlineData(20, 240, "19,1")]
        [InlineData(100, 1024, "84,16")] // 12 + 12 * 84 = 1020
        [InlineData(25, 1024, "25")]
        [InlineData(1, 240, "1")]
        public void Reads_are_planned_by_the_size_of_the_request(int tags, int pdu, string expected)
        {
            var jobs = S7Connection.PlanJobs(Tags(tags), false, pdu);

            Assert.Equal(expected, string.Join(",", jobs.Select(j => j.Count)));
            Assert.Equal(Enumerable.Range(0, tags), jobs.SelectMany(j => j));
        }

        [Fact]
        public void A_read_of_wider_tags_is_planned_by_the_size_of_the_request()
        {
            // 12 + 12 per tag: four of them fill 60 bytes. The answer needs less, 14 + 8 per REAL.
            var jobs = S7Connection.PlanJobs(Tags(12, "REAL"), false, 60);

            Assert.Equal("4,4,4", string.Join(",", jobs.Select(j => j.Count)));
        }

        [Fact]
        public void A_write_is_planned_by_the_size_of_the_request_with_its_data()
        {
            // 12 + 20 per REAL (twelve address bytes, four header bytes, four data bytes).
            var jobs = S7Connection.PlanJobs(Tags(5, "REAL"), true, 60);

            Assert.Equal("2,2,1", string.Join(",", jobs.Select(j => j.Count)));
        }

        [Fact]
        public void Nothing_to_do_is_no_job()
        {
            Assert.Empty(S7Connection.PlanJobs(new List<S7Tag>(), false, 240));
        }

        [Fact]
        public void A_tag_that_cannot_fit_a_job_still_gets_one_of_its_own()
        {
            var jobs = S7Connection.PlanJobs(Tags(3, "REAL"), true, 10);

            Assert.Equal("1,1,1", string.Join(",", jobs.Select(j => j.Count)));
        }

        // ── Writing ────────────────────────────────────────────

        [Theory]
        // address, value, the transport size of the address item, the one of the data item, the data
        [InlineData("%DB1.DBW10:INT", 258, 0x05, 0x05, "0102")]
        [InlineData("%DB1.DBW10:INT", -2, 0x05, 0x05, "FFFE")]
        [InlineData("%DB1.DBW10:UINT", 65534, 0x05, 0x05, "FFFE")]
        [InlineData("%M0.0:BOOL", true, 0x01, 0x03, "01")]
        [InlineData("%M0.0:BOOL", false, 0x01, 0x03, "00")]
        [InlineData("%DB1.DBB1:BYTE", 165, 0x02, 0x04, "A5")]
        [InlineData("%DB1.DBB1:SINT", -2, 0x02, 0x04, "FE")]
        [InlineData("%DB1.DBB1:USINT", 254, 0x02, 0x04, "FE")]
        [InlineData("%QW4:WORD", 48879, 0x04, 0x04, "BEEF")]
        [InlineData("%DB1.DBD4:DWORD", 3735928559L, 0x06, 0x04, "DEADBEEF")]
        [InlineData("%DB1.DBD4:DINT", -200, 0x07, 0x05, "FFFFFF38")]
        [InlineData("%DB1.DBD4:UDINT", 4294967096L, 0x07, 0x05, "FFFFFF38")]
        [InlineData("%DB1.DBD8:REAL", 3.14159f, 0x08, 0x07, "40490FD0")]
        public async Task A_value_is_written_as_the_type_of_the_tag(
            string address, object value, byte addressSize, byte dataSize, string dataHex)
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.WriteAck(2, 0xFF));
            var tag = S7Tag.Parse(address);
            var area = tag.Area == MemoryArea.DATA_BLOCKS ? DataBlocks
                : tag.Area == MemoryArea.OUTPUTS ? Outputs
                : tag.Area == MemoryArea.INPUTS ? Inputs : Flags;

            var response = await WriteAsync(connection, address, value);

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("v"));
            Assert.Equal(
                S7Frames.WriteRequest(
                    2,
                    new[] { S7Frames.AddressItem(addressSize, (ushort)tag.DbNumber, area, tag.ByteOffset, tag.BitOffset) },
                    (dataSize, Hex(dataHex))),
                Written(transport));
        }

        [Fact]
        public async Task A_write_to_a_flag_word_puts_the_value_big_endian_on_the_wire()
        {
            // Spelled out in bytes, not rebuilt from the tag: job header, parameter, one data item.
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.WriteAck(2, 0xFF));

            await WriteAsync(connection, "%MW6:WORD", 0x1234);

            Assert.Equal(
                S7Frames.Dt(Hex("32010000000200 0E 0006 0501 120A10040001000083000030 FF040010 1234".Replace(" ", ""))),
                Written(transport));
        }

        [Theory]
        [InlineData(0xFF, PlcResponseCode.Ok)]
        [InlineData(0x03, PlcResponseCode.AccessDenied)]
        [InlineData(0x05, PlcResponseCode.InvalidAddress)]
        [InlineData(0x0A, PlcResponseCode.NotFound)]
        [InlineData(0x06, PlcResponseCode.InvalidDatatype)]
        [InlineData(0x01, PlcResponseCode.InternalError)]
        public async Task The_status_of_each_item_becomes_the_response_code_of_its_tag(
            byte status, PlcResponseCode expected)
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.WriteAck(2, status, 0xFF));

            var response = await WriteAsync(connection, ("a", "%DB1.DBW0:INT", 1), ("b", "%DB1.DBW2:INT", 2));

            Assert.Equal(expected, response.GetResponseCode("a"));
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("b"));
        }

        [Fact]
        public async Task A_refusal_of_the_whole_write_is_access_denied_for_every_tag()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.Dt(S7Frames.Refusal(2, 0x81, 0x04)));

            var response = await WriteAsync(connection, ("a", "%DB1.DBW0:INT", 1), ("b", "%DB1.DBW2:INT", 2));

            Assert.Equal(PlcResponseCode.AccessDenied, response.GetResponseCode("a"));
            Assert.Equal(PlcResponseCode.AccessDenied, response.GetResponseCode("b"));
        }

        [Theory]
        [InlineData("%DB1.DBW0:INT", 40000)]
        [InlineData("%DB1.DBW0:INT", -32769)]
        [InlineData("%DB1.DBW0:INT", 1.5)]
        [InlineData("%DB1.DBW0:INT", "many")]
        [InlineData("%DB1.DBW0:INT", true)]
        [InlineData("%DB1.DBW0:UINT", -1)]
        [InlineData("%DB1.DBB0:BYTE", 256)]
        [InlineData("%DB1.DBB0:BYTE", -1)]
        [InlineData("%DB1.DBB0:SINT", 128)]
        [InlineData("%DB1.DBB0:USINT", 256)]
        [InlineData("%DB1.DBD0:DINT", 2147483648L)]
        [InlineData("%DB1.DBD0:UDINT", -1L)]
        [InlineData("%DB1.DBD0:DWORD", 4294967296L)]
        [InlineData("%DB1.DBD0:REAL", "not a number")]
        [InlineData("%DB1.DBD0:REAL", 1e300)]
        [InlineData("%DB1.DBD0:REAL", true)]
        [InlineData("%M0.0:BOOL", "maybe")]
        public async Task A_value_that_does_not_fit_the_type_is_an_invalid_datatype_and_nothing_goes_on_the_wire(
            string address, object value)
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);

            var response = await WriteAsync(connection, address, value);

            Assert.Equal(PlcResponseCode.InvalidDatatype, response.GetResponseCode("v"));
            Assert.Equal(0, transport.GetNumBytesWritten());
        }

        [Fact]
        public async Task A_value_that_does_not_fit_fails_its_own_tag_and_the_others_are_written()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport);
            transport.InjectTestData(S7Frames.WriteAck(2, 0xFF));

            var response = await WriteAsync(connection, ("bad", "%DB1.DBW0:INT", 70000), ("good", "%DB1.DBW2:INT", 5));

            Assert.Equal(PlcResponseCode.InvalidDatatype, response.GetResponseCode("bad"));
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("good"));
        }

        [Fact]
        public async Task A_write_without_an_answer_is_a_request_timeout()
        {
            var transport = Ready();
            var connection = ConnectQuietly(transport, S7Uri + "&read-timeout-ms=100");

            var response = await WriteAsync(connection, "%DB1.DBW0:INT", 1);

            Assert.Equal(PlcResponseCode.RequestTimeout, response.GetResponseCode("v"));
        }

        [Fact]
        public async Task A_write_of_many_tags_is_cut_into_jobs_that_fit_the_PDU()
        {
            // 12 + 18 per tag (twelve address bytes, four header bytes, two data bytes with the pad):
            // twelve tags fit 240 bytes.
            var transport = Ready(240);
            var connection = ConnectQuietly(transport);
            var cpu = new FakeCpu(
                transport,
                S7Frames.WriteAck(2, Enumerable.Repeat((byte)0xFF, 12).ToArray()),
                S7Frames.WriteAck(3, Enumerable.Repeat((byte)0xFF, 11).ToArray()));

            var response = await WriteAsync(
                connection, Enumerable.Range(0, 23).Select(i => ($"t{i}", $"%MB{i}:BYTE", (object)i)).ToArray());
            await cpu.Running;

            for (var i = 0; i < 23; i++)
            {
                Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode($"t{i}"));
            }
            Assert.Equal(2, cpu.Requests.Count);
        }

        // ── The tag the connection hands out ───────────────────

        [Fact]
        public void The_connection_parses_tags_with_the_s7_syntax_and_says_what_it_can_do()
        {
            var connection = new S7Connection(ConnectionString.Parse(S7Uri), NewTransport());

            Assert.IsType<S7Tag>(connection.Parse("%DB1.DBW10:INT"));
            Assert.Throws<PlcInvalidFieldException>(() => connection.Parse("holding-register:1"));
            Assert.True(connection.PlcConnectionMetadata.CanRead);
            Assert.True(connection.PlcConnectionMetadata.CanWrite);
            Assert.False(connection.PlcConnectionMetadata.CanSubscribe);
        }
    }
}
