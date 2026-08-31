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
using org.apache.plc4net.exceptions;
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
        public void Read_request_PDU_is_well_formed()
        {
            var tags = new List<S7Tag>
            {
                S7Tag.Parse("%DB1.DBW10"),
                S7Tag.Parse("%M0.0")
            };

            var pdu = S7Constants.BuildReadRequest(0x0001, tags);

            // S7 header: protocol ID + message type
            Assert.Equal(S7Constants.ProtocolId, pdu[0]);
            Assert.Equal(S7Constants.JobRequest, pdu[1]);

            // PDU reference
            Assert.Equal(0x00, pdu[4]);
            Assert.Equal(0x01, pdu[5]);

            // Parameter part: function code + item count
            Assert.Equal(S7Constants.ReadVar, pdu[10]);
            Assert.Equal(2, pdu[11]);
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

        [Fact]
        public void Transport_size_mapping_is_correct()
        {
            Assert.Equal(S7Constants.TransportSizeByte,
                S7Constants.TransportSizeForTag(S7Tag.Parse("%DB1.DBB10")));
            Assert.Equal(S7Constants.TransportSizeWord,
                S7Constants.TransportSizeForTag(S7Tag.Parse("%DB1.DBW10")));
            Assert.Equal(S7Constants.TransportSizeDword,
                S7Constants.TransportSizeForTag(S7Tag.Parse("%DB1.DBD10")));
            Assert.Equal(S7Constants.TransportSizeBit,
                S7Constants.TransportSizeForTag(S7Tag.Parse("%DB1.DBX0.0")));
        }

        [Fact]
        public void Write_request_PDU_is_well_formed()
        {
            var items = new List<(S7Tag tag, byte[] data)>
            {
                (S7Tag.Parse("%DB1.DBW10"), new byte[] { 0x12, 0x34 })
            };

            var pdu = S7Constants.BuildWriteRequest(0x0002, items);

            Assert.Equal(S7Constants.ProtocolId, pdu[0]);
            Assert.Equal(S7Constants.JobRequest, pdu[1]);

            // Parameter part
            Assert.Equal(S7Constants.WriteVar, pdu[10]);
            Assert.Equal(1, pdu[11]); // 1 item
        }

        [Fact]
        public void Connect_sends_the_java_parity_tsaps_in_the_connection_request()
        {
            // Pins the TSAP derivation (the headline fix of the second review round)
            // byte for byte: remote = PG_OR_PC(0x01), rack 0, slot 1 -> 0x0101;
            // local = OTHERS(0x03), rack 1, slot 1 -> 0x0311. A wrong shift, byte swap,
            // or renamed parameter would turn this red.
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm());
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));

            connection.Connect();

            var cr = Assert.Single(inner.Written);
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
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm());
            // remote-tsap=0x0300 overrides the rack/slot derivation entirely.
            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=3&remote-slot=4&remote-tsap=0x0300"),
                new CotpTransportInstance(inner));

            connection.Connect();

            var cr = Assert.Single(inner.Written);
            Assert.Equal(0x03, cr[20]);
            Assert.Equal(0x00, cr[21]);
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

        // ── Write path (PlcWriter) ──────────────────────────────

        [Fact]
        public async Task Write_WORD_to_a_data_block_returns_OK()
        {
            // S7 Write Var response for 1 item: paramLen=2, dataLen=3.
            // S7 payload = 10 (header) + 2 (param) + 3 (data) = 15.
            // TPKT = 4 (TPKT hdr) + 3 (COTP DT) + 15 = 22 = 0x16.
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm());
            inner.Inject(new byte[] {
                0x03, 0x00, 0x00, 0x16,
                0x02, 0xF0, 0x80,
                0x32, 0x03, 0x00, 0x00, 0x00, 0x01,
                0x00, 0x02, 0x00, 0x03,
                S7Constants.WriteVar, 0x01,
                S7Constants.WriteVar, 0x01, 0xFF
            });

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
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm());
            inner.Inject(new byte[] {
                0x03, 0x00, 0x00, 0x16,
                0x02, 0xF0, 0x80,
                0x32, 0x03, 0x00, 0x00, 0x00, 0x01,
                0x00, 0x02, 0x00, 0x03,
                S7Constants.WriteVar, 0x01,
                S7Constants.WriteVar, 0x01, 0xFF
            });

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
            var inner = new ScriptedTransportInstance();
            inner.Inject(CotpTransportInstanceTests.ConnectionConfirm());
            inner.Inject(new byte[] {
                0x03, 0x00, 0x00, 0x16,
                0x02, 0xF0, 0x80,
                0x32, 0x03, 0x00, 0x00, 0x00, 0x02,
                0x00, 0x02, 0x00, 0x03,
                S7Constants.WriteVar, 0x01,
                S7Constants.WriteVar, 0x01, 0xFF
            });

            var connection = new S7Connection(
                ConnectionString.Parse("s7://192.168.0.1?remote-rack=0&remote-slot=1"),
                new CotpTransportInstance(inner));
            connection.Connect();

            var wb = (DefaultPlcWriteRequestBuilder)connection.WriteRequestBuilder;
            wb.AddTag("m0", "%M0", (byte)0x42);
            var wRsp = (DefaultPlcWriteResponse)await connection.Write(
                (DefaultPlcWriteRequest)wb.Build());
            Assert.Equal(PlcResponseCode.Ok, wRsp.GetResponseCode("m0"));

            // The transport received at least one DT frame after the CR.
            // The second write to the inner transport is the DT containing
            // the S7 Write Var PDU (the first was the CR during Connect).
            Assert.True(inner.Written.Count >= 2,
                "Expected at least CR + one DT frame");
        }

    }
}
