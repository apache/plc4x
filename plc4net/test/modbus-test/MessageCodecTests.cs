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
using org.apache.plc4net.drivers.modbus.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.transports.test;
using Xunit;

namespace org.apache.plc4net.drivers.modbus.test
{
    /// <summary>
    /// Frame boundaries: how the TCP and RTU codecs cut the incoming byte stream
    /// into whole Modbus ADUs.
    /// </summary>
    public class MessageCodecTests
    {
        private static TestTransportInstance NewTransport()
            => new TestTransportInstance(new TestTransportConfiguration());

        // ── CRC reference ──────────────────────────────────────

        [Fact]
        public void The_reference_CRC_matches_the_shared_suite()
        {
            // 01 03 00 00 00 0A is the request of the shared RTU suite; its frame
            // ends in c5 cd (CRC low byte first).
            var frame = ModbusFrames.Rtu(1, 0x03, 0x00, 0x00, 0x00, 0x0A);
            Assert.Equal(new byte[] { 0x01, 0x03, 0x00, 0x00, 0x00, 0x0A, 0xC5, 0xCD }, frame);
        }

        // ── Modbus TCP ─────────────────────────────────────────

        private static (ModbusTcpMessageCodec Codec, List<ModbusTcpADU> Received) NewTcp(
            TestTransportInstance transport)
        {
            var received = new List<ModbusTcpADU>();
            return (new ModbusTcpMessageCodec(transport, received.Add), received);
        }

        [Fact]
        public void Tcp_a_whole_frame_is_delivered()
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            transport.InjectTestData(ModbusFrames.Tcp(0x0102, 7, 0x03, 0x02, 0x00, 0x42));

            codec.ProcessIncomingData();

            var adu = Assert.Single(received);
            Assert.Equal((ushort)0x0102, adu.TransactionIdentifier);
            Assert.Equal((byte)7, adu.UnitIdentifier);
            var pdu = Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(adu.Pdu);
            Assert.Equal(new byte[] { 0x00, 0x42 }, pdu.Value);
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        [Fact]
        public void Tcp_a_frame_that_arrives_in_pieces_is_delivered_once_complete()
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            var frame = ModbusFrames.Tcp(1, 1, 0x03, 0x02, 0x00, 0x42);

            transport.InjectTestData(frame[..3]); // not even the whole header
            codec.ProcessIncomingData();
            Assert.Empty(received);

            transport.InjectTestData(frame[3..9]); // header done, PDU incomplete
            codec.ProcessIncomingData();
            Assert.Empty(received);
            Assert.Equal(9, transport.GetNumBytesAvailable()); // nothing was consumed

            transport.InjectTestData(frame[9..]);
            codec.ProcessIncomingData();
            Assert.Single(received);
        }

        [Fact]
        public void Tcp_two_frames_in_one_chunk_are_delivered_in_order()
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            var first = ModbusFrames.Tcp(1, 1, 0x03, 0x02, 0x11, 0x11);
            var second = ModbusFrames.Tcp(2, 1, 0x04, 0x02, 0x22, 0x22);
            transport.InjectTestData(Concat(first, second));

            codec.ProcessIncomingData();

            Assert.Equal(new ushort[] { 1, 2 },
                new[] { received[0].TransactionIdentifier, received[1].TransactionIdentifier });
            Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(received[0].Pdu);
            Assert.IsType<ModbusPDUReadInputRegistersResponse>(received[1].Pdu);
        }

        [Fact]
        public void Tcp_an_exception_response_is_delivered_as_an_error_pdu()
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            transport.InjectTestData(ModbusFrames.Tcp(5, 1, 0x83, 0x02));

            codec.ProcessIncomingData();

            var error = Assert.IsType<ModbusPDUError>(Assert.Single(received).Pdu);
            Assert.Equal(ModbusErrorCode.ILLEGAL_DATA_ADDRESS, error.ExceptionCode);
        }

        [Fact]
        public void Tcp_junk_in_front_of_a_frame_is_skipped_byte_by_byte()
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            var junk = new byte[] { 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF };
            transport.InjectTestData(Concat(junk, ModbusFrames.Tcp(9, 1, 0x03, 0x02, 0x00, 0x42)));

            codec.ProcessIncomingData();

            Assert.Equal((ushort)9, Assert.Single(received).TransactionIdentifier);
        }

        [Theory]
        [InlineData(0x00, 0x01, 0x00, 0x06)] // protocol identifier 1 is not Modbus
        [InlineData(0x00, 0x00, 0x00, 0x00)] // length 0
        [InlineData(0x00, 0x00, 0x00, 0x02)] // too short to hold even an exception PDU
        [InlineData(0x00, 0x00, 0x01, 0x00)] // length 256 exceeds the 253-byte PDU limit
        public void Tcp_a_header_that_cannot_start_a_frame_is_skipped(
            byte protocolHigh, byte protocolLow, byte lengthHigh, byte lengthLow)
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            // Six header bytes and nothing behind them: resync drops one byte at a time
            // until fewer than six remain, and never delivers a message.
            transport.InjectTestData(new byte[]
            {
                0x00, 0x01, protocolHigh, protocolLow, lengthHigh, lengthLow, 0x01, 0x03,
            });

            codec.ProcessIncomingData();

            Assert.Empty(received);
        }

        [Fact]
        public void Tcp_a_frame_that_does_not_parse_is_consumed_and_reported()
        {
            var transport = NewTransport();
            var (codec, received) = NewTcp(transport);
            // A well-formed MBAP frame whose PDU carries an unknown function code.
            transport.InjectTestData(ModbusFrames.Tcp(1, 1, 0x7F, 0x00, 0x00));

            Assert.Throws<MessageCodecException>(() => codec.ProcessIncomingData());

            Assert.Empty(received);
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        // ── Modbus RTU ─────────────────────────────────────────

        private static (ModbusRtuMessageCodec Codec, List<ModbusRtuADU> Received) NewRtu(
            TestTransportInstance transport)
        {
            var received = new List<ModbusRtuADU>();
            return (new ModbusRtuMessageCodec(transport, received.Add), received);
        }

        [Fact]
        public void Rtu_a_read_response_is_framed_by_its_byte_count()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            // Four registers: 8 data bytes.
            transport.InjectTestData(ModbusFrames.Rtu(1,
                0x03, 0x08, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04));

            codec.ProcessIncomingData();

            var adu = Assert.Single(received);
            Assert.Equal((byte)1, adu.Address);
            var pdu = Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(adu.Pdu);
            Assert.Equal(8, pdu.Value.Length);
        }

        [Fact]
        public void Rtu_a_frame_that_arrives_in_pieces_is_delivered_once_complete()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            var frame = ModbusFrames.Rtu(1, 0x03, 0x02, 0x12, 0x34);

            transport.InjectTestData(frame[..2]); // address + function only
            codec.ProcessIncomingData();
            Assert.Empty(received);

            transport.InjectTestData(frame[2..6]); // byte count + data, CRC not yet complete
            codec.ProcessIncomingData();
            Assert.Empty(received);

            transport.InjectTestData(frame[6..]);
            codec.ProcessIncomingData();
            Assert.Single(received);
        }

        [Fact]
        public void Rtu_an_exception_response_is_five_bytes()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            transport.InjectTestData(ModbusFrames.Rtu(1, 0x83, 0x02));

            codec.ProcessIncomingData();

            var error = Assert.IsType<ModbusPDUError>(Assert.Single(received).Pdu);
            Assert.Equal(ModbusErrorCode.ILLEGAL_DATA_ADDRESS, error.ExceptionCode);
        }

        [Theory]
        [InlineData(0x05, 0x00, 0x0A, 0xFF, 0x00)] // Write Single Coil echo
        [InlineData(0x06, 0x00, 0x05, 0x12, 0x34)] // Write Single Register echo
        public void Rtu_a_write_echo_is_eight_bytes(byte function, byte a, byte b, byte c, byte d)
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            transport.InjectTestData(ModbusFrames.Rtu(1, function, a, b, c, d));

            codec.ProcessIncomingData();

            Assert.Single(received);
        }

        [Fact]
        public void Rtu_two_frames_in_one_chunk_are_delivered_in_order()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            transport.InjectTestData(Concat(
                ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x01),
                ModbusFrames.Rtu(1, 0x04, 0x02, 0x00, 0x02)));

            codec.ProcessIncomingData();

            Assert.Equal(2, received.Count);
            Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(received[0].Pdu);
            Assert.IsType<ModbusPDUReadInputRegistersResponse>(received[1].Pdu);
        }

        [Fact]
        public void Rtu_a_frame_with_a_bad_crc_is_not_delivered()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            var frame = ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42);
            frame[^1] ^= 0xFF;
            transport.InjectTestData(frame);

            codec.ProcessIncomingData(); // skips what is not a frame; it does not throw

            Assert.Empty(received);
        }

        [Fact]
        public void Rtu_a_corrupted_frame_does_not_hide_the_good_frame_behind_it()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            var corrupted = ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42);
            corrupted[^1] ^= 0xFF;
            transport.InjectTestData(Concat(corrupted, ModbusFrames.Rtu(1, 0x03, 0x02, 0x12, 0x34)));

            codec.ProcessIncomingData();

            var pdu = Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(Assert.Single(received).Pdu);
            Assert.Equal(new byte[] { 0x12, 0x34 }, pdu.Value);
        }

        [Fact]
        public void Rtu_noise_in_front_of_a_frame_is_skipped()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            transport.InjectTestData(Concat(
                new byte[] { 0x55, 0x55, 0x55 }, ModbusFrames.Rtu(1, 0x03, 0x02, 0x12, 0x34)));

            codec.ProcessIncomingData();

            var pdu = Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(Assert.Single(received).Pdu);
            Assert.Equal(new byte[] { 0x12, 0x34 }, pdu.Value);
        }

        [Fact]
        public void Rtu_noise_that_announces_a_long_frame_does_not_hide_the_good_frame_behind_it()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            // 01 03 C8 looks like the start of a 205-byte read response, and the real answer
            // is right behind it: waiting for those 205 bytes would never look at it.
            transport.InjectTestData(Concat(
                new byte[] { 0x01, 0x03, 0xC8 }, ModbusFrames.Rtu(1, 0x03, 0x02, 0x12, 0x34)));

            codec.ProcessIncomingData();

            var pdu = Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(Assert.Single(received).Pdu);
            Assert.Equal(new byte[] { 0x12, 0x34 }, pdu.Value);
        }

        [Fact]
        public void Rtu_a_long_frame_that_is_still_arriving_is_waited_for_when_nothing_good_follows_it()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            var pdu = new byte[2 + 100];
            pdu[0] = 0x03;
            pdu[1] = 100;
            var frame = ModbusFrames.Rtu(1, pdu);

            transport.InjectTestData(frame[..60]);
            codec.ProcessIncomingData();
            Assert.Empty(received);
            Assert.Equal(60, transport.GetNumBytesAvailable()); // kept, not dropped byte by byte

            transport.InjectTestData(frame[60..]);
            codec.ProcessIncomingData();
            Assert.Equal(100, Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(
                Assert.Single(received).Pdu).Value.Length);
        }

        [Fact]
        public void Rtu_a_complete_frame_is_checked_before_it_is_consumed()
        {
            // An incomplete frame is never judged: the bytes stay where they are until the
            // rest arrives, whatever the CRC will turn out to be.
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            var frame = ModbusFrames.Rtu(1, 0x03, 0x02, 0x00, 0x42);
            transport.InjectTestData(frame[..6]);

            codec.ProcessIncomingData();

            Assert.Empty(received);
            Assert.Equal(6, transport.GetNumBytesAvailable());
        }

        [Fact]
        public void Rtu_unsupported_function_codes_are_skipped_byte_by_byte()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            // Every three-byte window has an unsupported function code in the middle.
            transport.InjectTestData(new byte[] { 0x00, 0x55, 0x00, 0x00 });

            codec.ProcessIncomingData();

            Assert.Empty(received);
            Assert.Equal(2, transport.GetNumBytesAvailable()); // fewer than a header remain
        }

        [Fact]
        public void Rtu_an_implausible_read_byte_count_is_not_framed()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            // Byte count 251 exceeds the 250 bytes a response can carry. The frame is
            // complete and its CRC is right, so only the byte count stands in the way.
            var pdu = new byte[2 + 251];
            pdu[0] = 0x03;
            pdu[1] = 0xFB;
            transport.InjectTestData(ModbusFrames.Rtu(1, pdu));

            codec.ProcessIncomingData();

            Assert.Empty(received);
        }

        [Fact]
        public void Rtu_the_largest_read_response_is_framed()
        {
            var transport = NewTransport();
            var (codec, received) = NewRtu(transport);
            // 125 registers: 250 bytes of data, the most a response carries.
            var pdu = new byte[2 + 250];
            pdu[0] = 0x03;
            pdu[1] = 0xFA;
            transport.InjectTestData(ModbusFrames.Rtu(1, pdu));

            codec.ProcessIncomingData();

            var response = Assert.IsType<ModbusPDUReadHoldingRegistersResponse>(Assert.Single(received).Pdu);
            Assert.Equal(250, response.Value.Length);
        }

        private static byte[] Concat(byte[] first, byte[] second)
        {
            var result = new byte[first.Length + second.Length];
            Array.Copy(first, result, first.Length);
            Array.Copy(second, 0, result, first.Length, second.Length);
            return result;
        }
    }
}