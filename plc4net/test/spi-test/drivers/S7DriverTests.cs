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

using System.Collections.Generic;
using System.Threading.Tasks;
using org.apache.plc4net.drivers.s7;
using org.apache.plc4net.drivers.s7.messages;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.exceptions;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.model.values;
using org.apache.plc4net.spi.test.transports;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.transports.cotp;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// Tests for the S7 driver — tag parsing, PDU encoding, TPKT framing.
    /// </summary>
    public class S7DriverTests
    {
        /// <summary>
        /// A Setup Communication AckData, wrapped in TPKT + COTP DT, that a scripted
        /// inner transport can serve so <see cref="S7Connection.Connect"/> completes.
        /// The CPU negotiates the PDU length down to 240.
        /// </summary>
        internal static byte[] SetupCommunicationAck() => new byte[]
        {
            0x03, 0x00, 0x00, 0x1B,             // TPKT: version, reserved, length = 27
            0x02, 0xF0, 0x80,                   // COTP DT: LI = 2, DT, TPDU-NR + EOT
            0x32, 0x03, 0x00, 0x00, 0x00, 0x02, // S7: protocol id, AckData, redundancy, ref
            0x00, 0x08, 0x00, 0x00,             // param length = 8, data length = 0
            0x00, 0x00,                         // error class 0, error code 0
            0xF0, 0x00,                         // function = Setup Communication, reserved
            0x00, 0x01, 0x00, 0x01,             // max AmQ caller / callee = 1
            0x00, 0xF0,                         // negotiated PDU length = 240
        };

        /// <summary>Injects the CC + Setup-Communication ack a full Connect() needs.</summary>
        private static ScriptedTransportInstance HandshakeReady()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm());
            inner.Inject(SetupCommunicationAck());
            return inner;
        }

        // ── S7 response builders (via the generated model, wrapped in TPKT + COTP DT) ──

        private static byte[] WrapS7(byte[] s7)
        {
            var total = 4 + 3 + s7.Length;
            var frame = new byte[total];
            frame[0] = 0x03;
            frame[2] = (byte)(total >> 8);
            frame[3] = (byte)total;
            frame[4] = 0x02;
            frame[5] = 0xF0;
            frame[6] = 0x80;
            System.Array.Copy(s7, 0, frame, 7, s7.Length);
            return frame;
        }

        private static byte[] Serialize(S7Message message)
        {
            var wb = new WriteBuffer();
            message.Serialize(wb);
            return wb.GetBytes();
        }

        private static byte[] WriteVarAck(int numItems)
        {
            var items = new List<S7VarPayloadStatusItem>();
            for (var i = 0; i < numItems; i++)
                items.Add(new S7VarPayloadStatusItem(DataTransportErrorCode.OK));
            return WrapS7(Serialize(new S7MessageResponseData(1, 0, 0,
                new S7ParameterWriteVarResponse((byte)numItems),
                new S7PayloadWriteVarResponse(items))));
        }

        private static byte[] ReadVarAck(params (DataTransportSize Size, byte[] Data)[] items)
        {
            var payloadItems = new List<S7VarPayloadDataItem>();
            foreach (var (size, data) in items)
                payloadItems.Add(new S7VarPayloadDataItem(DataTransportErrorCode.OK, size, data));
            return WrapS7(Serialize(new S7MessageResponseData(1, 0, 0,
                new S7ParameterReadVarResponse((byte)items.Length),
                new S7PayloadReadVarResponse(payloadItems))));
        }

        [Theory]
        [InlineData("%DB1.DBW10", S7Tag.AreaType.DataBlock, 1, 10, -1, 2)]
        [InlineData("%DB1.DBX0.0", S7Tag.AreaType.DataBlock, 1, 0, 0, 1)]
        [InlineData("%DB5.DBD20", S7Tag.AreaType.DataBlock, 5, 20, -1, 4)]
        [InlineData("%M0.0", S7Tag.AreaType.Merker, 0, 0, 0, 1)]
        [InlineData("%M10", S7Tag.AreaType.Merker, 0, 10, -1, 1)]
        [InlineData("%I0.0", S7Tag.AreaType.Input, 0, 0, 0, 1)]
        [InlineData("%Q4.5", S7Tag.AreaType.Output, 0, 4, 5, 1)]
        public void Tag_parsing_yields_correct_areas_and_offsets(
            string address, S7Tag.AreaType expectedArea,
            int expectedDb, int expectedOffset, int expectedBit, int expectedSize)
        {
            var tag = S7Tag.Parse(address);

            Assert.Equal(expectedArea, tag.Area);
            Assert.Equal(expectedDb, tag.DbNumber);
            Assert.Equal(expectedOffset, tag.ByteOffset);
            Assert.Equal(expectedBit, tag.BitOffset);
            Assert.Equal(expectedSize, tag.DataTypeSize);
        }

        [Fact]
        public void Tag_parsing_rejects_invalid_input()
        {
            Assert.Throws<S7DriverException>(() => S7Tag.Parse(""));
            Assert.Throws<S7DriverException>(() => S7Tag.Parse("garbage"));
            Assert.Throws<S7DriverException>(() => S7Tag.Parse("%X0.0"));
        }

        [Fact]
        public void Read_request_PDU_round_trips_through_the_generated_model()
        {
            var tags = new List<S7Tag>
            {
                S7Tag.Parse("%DB1.DBW10"),
                S7Tag.Parse("%M0.0"),
            };

            var pdu = S7Constants.BuildReadRequest(0x1234, tags);
            var msg = Assert.IsType<S7MessageRequest>(S7Message.StaticParse(new ReadBuffer(pdu)));

            Assert.Equal(0x1234, msg.TpduReference);
            var param = Assert.IsType<S7ParameterReadVarRequest>(msg.Parameter);
            Assert.Equal(2, param.Items.Count);

            var addr0 = Assert.IsType<S7AddressAny>(
                Assert.IsType<S7VarRequestParameterItemAddress>(param.Items[0]).Address);
            Assert.Equal(MemoryArea.DATA_BLOCKS, addr0.Area);   // %DB... -> 0x84, not 0x82
            Assert.Equal((ushort)1, addr0.DbNumber);
            Assert.Equal(TransportSize.WORD, addr0.TransportSize);
            Assert.Equal((ushort)10, addr0.ByteAddress);

            var addr1 = Assert.IsType<S7AddressAny>(
                Assert.IsType<S7VarRequestParameterItemAddress>(param.Items[1]).Address);
            Assert.Equal(MemoryArea.FLAGS_MARKERS, addr1.Area);
            Assert.Equal(TransportSize.BOOL, addr1.TransportSize);
            Assert.Equal((byte)0, addr1.BitAddress);
        }

        [Fact]
        public void Tpkt_wrap_and_unwrap_are_inverses()
        {
            var payload = new byte[] { 0x32, 0x01, 0x00, 0x00, 0x00, 0x01 };
            var frame = TpktFrame.Wrap(payload);

            Assert.Equal(TpktFrame.Version, frame[0]);
            Assert.Equal(0, frame[1]);
            Assert.Equal(payload.Length + 4, (frame[2] << 8) | frame[3]);

            var unwrapped = TpktFrame.Unwrap(frame);
            Assert.Equal(payload, unwrapped);
        }

        [Fact]
        public void Tpkt_read_payload_length_is_correct()
        {
            var payload = new byte[20];
            var frame = TpktFrame.Wrap(payload);

            var payloadLen = TpktFrame.ReadPayloadLength(frame);

            Assert.Equal(20, payloadLen);
        }

        [Theory]
        [InlineData("%DB1.DBB10", TransportSize.BYTE)]
        [InlineData("%DB1.DBW10", TransportSize.WORD)]
        [InlineData("%DB1.DBD10", TransportSize.DWORD)]
        [InlineData("%DB1.DBX0.0", TransportSize.BOOL)]
        public void The_address_transport_size_matches_the_tag_width(string address, TransportSize expected)
        {
            var pdu = S7Constants.BuildReadRequest(1, new List<S7Tag> { S7Tag.Parse(address) });
            var param = (S7ParameterReadVarRequest)((S7MessageRequest)
                S7Message.StaticParse(new ReadBuffer(pdu))).Parameter!;
            var addr = (S7AddressAny)((S7VarRequestParameterItemAddress)param.Items[0]).Address;
            Assert.Equal(expected, addr.TransportSize);
        }

        [Fact]
        public void Write_request_PDU_round_trips_through_the_generated_model()
        {
            var items = new List<(S7Tag tag, byte[] data)>
            {
                (S7Tag.Parse("%DB1.DBW10"), new byte[] { 0x12, 0x34 }),
            };

            var pdu = S7Constants.BuildWriteRequest(0x0002, items);
            var msg = Assert.IsType<S7MessageRequest>(S7Message.StaticParse(new ReadBuffer(pdu)));

            var param = Assert.IsType<S7ParameterWriteVarRequest>(msg.Parameter);
            Assert.Single(param.Items);

            var payload = Assert.IsType<S7PayloadWriteVarRequest>(msg.Payload);
            var dataItem = Assert.Single(payload.Items);
            Assert.Equal(DataTransportSize.BYTE_WORD_DWORD, dataItem.TransportSize);
            Assert.Equal(new byte[] { 0x12, 0x34 }, dataItem.Data);
        }

        [Fact]
        public void Connect_sends_the_java_parity_tsaps_in_the_connection_request()
        {
            // Pins the TSAP derivation (the headline fix of the second review round)
            // byte for byte: remote = PG_OR_PC(0x01), rack 0, slot 1 -> 0x0101;
            // local = OTHERS(0x03), rack 1, slot 1 -> 0x0311. A wrong shift, byte swap,
            // or renamed parameter would turn this red.
            var inner = HandshakeReady();
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));

            connection.Connect();

            // First frame written is the COTP CR; the second is Setup Communication.
            var cr = inner.Written[0];
            // Called TSAP parameter (0xC2): remote TSAP hi/lo = 0x01 0x01.
            Assert.Equal(0xC2, cr[18]);
            Assert.Equal(0x01, cr[20]);
            Assert.Equal(0x01, cr[21]);
            // Calling TSAP parameter (0xC1): local TSAP hi/lo = 0x03 0x11.
            Assert.Equal(0xC1, cr[14]);
            Assert.Equal(0x03, cr[16]);
            Assert.Equal(0x11, cr[17]);
        }

        [Fact]
        public void Connect_honours_an_explicit_remote_tsap_override()
        {
            var inner = HandshakeReady();
            // remote-tsap=0x0300 overrides the rack/slot derivation entirely.
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=3&remote-slot=4&remote-tsap=0x0300"),
                new CotpTransportInstance(inner));

            connection.Connect();

            var cr = inner.Written[0];
            Assert.Equal(0x03, cr[20]);
            Assert.Equal(0x00, cr[21]);
        }

        [Fact]
        public void Connect_runs_S7_Setup_Communication_and_records_the_negotiated_pdu_length()
        {
            var inner = HandshakeReady();
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));

            connection.Connect();

            // CR, then a Setup Communication job (S7 header 0x32 0x01, function 0xF0).
            Assert.Equal(2, inner.Written.Count);
            var setup = inner.Written[1];
            var s7 = setup[7..]; // strip TPKT(4) + COTP DT(3)
            Assert.Equal(0x32, s7[0]);
            Assert.Equal(S7Constants.JobRequest, s7[1]);
            Assert.Equal(S7Constants.SetupCommunication, s7[10]); // first parameter byte
            Assert.Equal(240, connection.NegotiatedPduLength);
        }

        [Fact]
        public void Connect_fails_cleanly_when_the_CPU_does_not_answer_Setup_Communication()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm()); // CC only, no Setup ack
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1&request-timeout=200"),
                new CotpTransportInstance(inner));

            var ex = Assert.Throws<PlcConnectionException>(() => connection.Connect());
            Assert.Contains("Setup Communication", ex.Message);
        }

        [Fact]
        public void Connect_rejects_a_non_numeric_rack_parameter()
        {
            // Java's Integer.parseInt fails on 'remote-rack=1a'; the shared
            // GetIntParameter silently falls back to 0, so the S7 driver must parse
            // strictly instead.
            var ex = Assert.Throws<PlcConnectionException>(() =>
                new S7Connection(
                    ConnectionString.Parse("s7://192.168.0.1?remote-rack=1a"),
                    new FakeTransportInstance(new FakeTransportConfiguration())));
            Assert.Contains("remote-rack", ex.Message);
        }

        [Fact]
        public void Connect_with_a_non_cotp_transport_is_rejected()
        {
            // The S7 protocol is ISO-on-TCP: without the COTP handshake there is no
            // session, so a connection on a raw transport must fail loudly instead of
            // returning a connection that silently cannot talk. (S7Driver declares only
            // the cotp transport; this guards the OnConnect hook against any future path
            // that hands it a different one.)
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new FakeTransportInstance(new FakeTransportConfiguration()));

            var ex = Assert.Throws<PlcConnectionException>(() => connection.Connect());
            Assert.Contains("cotp", ex.Message, System.StringComparison.OrdinalIgnoreCase);
        }

        [Fact]
        public void Connect_rejects_an_out_of_range_remote_slot()
        {
            // The TSAP reserves 4 bits for the slot; 16 cannot be represented, and the
            // value must be rejected rather than silently truncated.
            var ex = Assert.Throws<PlcConnectionException>(() =>
                new S7Connection(
                    ConnectionString.Parse("s7://192.168.0.1?remote-slot=16"),
                    new FakeTransportInstance(new FakeTransportConfiguration())));
            Assert.Contains("remote-slot", ex.Message);
        }

        // ── Read path (PlcReader) ───────────────────────────────

        [Fact]
        public async Task Read_decodes_a_multi_item_response_item_by_item()
        {
            var inner = HandshakeReady();
            inner.Inject(ReadVarAck(
                (DataTransportSize.BIT, new byte[] { 0x01 }),                      // %M0.0 -> true
                (DataTransportSize.BYTE_WORD_DWORD, new byte[] { 0x12, 0x34 })));   // %DB1.DBW0 -> 0x1234

            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));
            connection.Connect();

            var builder = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
            builder.AddTagAddress("flag", "%M0.0");
            builder.AddTagAddress("word", "%DB1.DBW0");
            var response = (DefaultPlcReadResponse)await connection.Read(
                (DefaultPlcReadRequest)builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("flag"));
            Assert.True(response.GetValue("flag").GetBool());
            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("word"));
            Assert.Equal((ushort)0x1234, response.GetValue("word").GetUshort());
        }

        [Fact]
        public async Task Read_maps_an_S7_item_error_to_a_response_code()
        {
            var inner = HandshakeReady();
            inner.Inject(WrapS7(Serialize(new S7MessageResponseData(1, 0, 0,
                new S7ParameterReadVarResponse(1),
                new S7PayloadReadVarResponse(new List<S7VarPayloadDataItem>
                {
                    new S7VarPayloadDataItem(DataTransportErrorCode.NOT_FOUND, DataTransportSize.NULL, System.Array.Empty<byte>()),
                })))));

            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));
            connection.Connect();

            var builder = (DefaultPlcReadRequestBuilder)connection.ReadRequestBuilder;
            builder.AddTagAddress("missing", "%DB99.DBW0");
            var response = (DefaultPlcReadResponse)await connection.Read(
                (DefaultPlcReadRequest)builder.Build());

            Assert.Equal(PlcResponseCode.NotFound, response.GetResponseCode("missing"));
        }

        // ── Write path (PlcWriter) ──────────────────────────────

        [Fact]
        public async Task Write_WORD_to_a_data_block_returns_OK()
        {
            // S7 Write Var response for 1 item: paramLen=2, dataLen=3.
            // S7 payload = 10 (header) + 2 (param) + 3 (data) = 15.
            // TPKT = 4 (TPKT hdr) + 3 (COTP DT) + 15 = 22 = 0x16.
            var inner = HandshakeReady();
            inner.Inject(WriteVarAck(1));

            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));
            connection.Connect();

            var builder = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder;
            builder.AddTag("data", "%DB1.DBW10", (ushort)0x1234);
            var response = (DefaultPlcWriteResponse)await connection.Write(
                (DefaultPlcWriteRequest)builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("data"));
        }

        [Fact]
        public async Task Write_BOOL_to_an_output_returns_OK()
        {
            var inner = HandshakeReady();
            inner.Inject(WriteVarAck(1));

            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));
            connection.Connect();

            var builder = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder;
            builder.AddTag("q", "%Q0.0", true);
            var response = (DefaultPlcWriteResponse)await connection.Write(
                (DefaultPlcWriteRequest)builder.Build());

            Assert.Equal(PlcResponseCode.Ok, response.GetResponseCode("q"));
        }

        [Fact]
        public async Task Write_byte_to_Merker_produces_the_correct_wire_PDU()
        {
            // Verify the Write Var PDU that reaches the transport
            // carries the expected S7 function code and data.
            var inner = HandshakeReady();
            inner.Inject(WriteVarAck(1));

            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));
            connection.Connect();

            var wb = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder;
            wb.AddTag("m0", "%M0", (byte)0x42);
            var wRsp = (DefaultPlcWriteResponse)await connection.Write(
                (DefaultPlcWriteRequest)wb.Build());
            Assert.Equal(PlcResponseCode.Ok, wRsp.GetResponseCode("m0"));

            // Writes to the inner transport: CR, Setup Communication, then the Write
            // Var DT frame. The last carries the S7 Write Var function code.
            Assert.Equal(3, inner.Written.Count);
            var writeVarPdu = inner.Written[2][7..]; // strip TPKT(4) + COTP DT(3)
            Assert.Equal(0x32, writeVarPdu[0]);
            Assert.Equal(S7Constants.WriteVar, writeVarPdu[10]); // first parameter byte
        }

    }
}
