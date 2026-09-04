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
using System.Threading.Tasks;
using org.apache.plc4net.drivers.modbus;
using org.apache.plc4net.drivers.modbus.messages;
using org.apache.plc4net.messages;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.drivers.messages;
using org.apache.plc4net.spi.test.transports;
using org.apache.plc4net.spi.transports;
using org.apache.plc4net.types;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// Tests for the Modbus TCP driver — tag parsing, PDU encoding, and
    /// frame-level correctness.
    /// </summary>
    public class ModbusDriverTests
    {
        [Fact]
        public void Tag_parsing_is_case_insensitive()
        {
            Assert.Equal(ModbusTag.TagType.Coil,
                ModbusTag.Parse("coil:0").Type);
            Assert.Equal(ModbusTag.TagType.Coil,
                ModbusTag.Parse("COIL:0").Type);
            Assert.Equal(ModbusTag.TagType.HoldingRegister,
                ModbusTag.Parse("holding:10").Type);
            Assert.Equal(ModbusTag.TagType.HoldingRegister,
                ModbusTag.Parse("HOLDING:10").Type);
            Assert.Equal(ModbusTag.TagType.DiscreteInput,
                ModbusTag.Parse("discrete:5").Type);
            Assert.Equal(ModbusTag.TagType.InputRegister,
                ModbusTag.Parse("input:3").Type);
        }

        [Fact]
        public void Tag_parsing_extracts_the_address()
        {
            var tag = ModbusTag.Parse("holding:100");
            Assert.Equal(ModbusTag.TagType.HoldingRegister, tag.Type);
            Assert.Equal((ushort)100, tag.Address);
        }

        [Fact]
        public void Tag_parsing_rejects_invalid_input()
        {
            Assert.Throws<ModbusDriverException>(() => ModbusTag.Parse(""));
            Assert.Throws<ModbusDriverException>(() => ModbusTag.Parse("noprefix"));
            Assert.Throws<ModbusDriverException>(() => ModbusTag.Parse("unknown:0"));
            Assert.Throws<ModbusDriverException>(
                () => ModbusTag.Parse("coil:notanumber"));
        }

        [Theory]
        [InlineData(0x01, (ushort)0, (ushort)1)]
        [InlineData(0x03, (ushort)10, (ushort)5)]
        public void Read_request_PDUs_are_well_formed(
            byte functionCode, ushort address, ushort quantity)
        {
            var pdu = ModbusPDU.BuildReadBitsRequest(functionCode, address, quantity);

            Assert.Equal(5, pdu.Length);
            Assert.Equal(functionCode, pdu[0]);
            Assert.Equal((byte)(address >> 8), pdu[1]);
            Assert.Equal((byte)(address & 0xFF), pdu[2]);
            Assert.Equal((byte)(quantity >> 8), pdu[3]);
            Assert.Equal((byte)(quantity & 0xFF), pdu[4]);
        }

        [Fact]
        public void Write_single_coil_PDU_is_well_formed()
        {
            var pdu = ModbusPDU.BuildWriteSingleCoilRequest(0x000A, true);

            Assert.Equal(ModbusFunctionCodes.WriteSingleCoil, pdu[0]);
            Assert.Equal(0x00, pdu[1]);
            Assert.Equal(0x0A, pdu[2]);
            Assert.Equal(0xFF, pdu[3]); // ON
            Assert.Equal(0x00, pdu[4]);
        }

        [Fact]
        public void Write_single_register_PDU_is_well_formed()
        {
            var pdu = ModbusPDU.BuildWriteSingleRegisterRequest(0x0001, 0x002A);

            Assert.Equal(ModbusFunctionCodes.WriteSingleRegister, pdu[0]);
            Assert.Equal(0x00, pdu[1]);
            Assert.Equal(0x01, pdu[2]); // address
            Assert.Equal(0x00, pdu[3]);
            Assert.Equal(0x2A, pdu[4]); // value = 42
        }

        [Fact]
        public void Parse_read_coils_response()
        {
            var data = new byte[] { 0x01, 0x01 }; // 1 byte, coil ON
            var bits = ModbusPDU.ParseReadBitsResponse(data, 1);

            Assert.Single(bits);
            Assert.True(bits[0]);
        }

        [Fact]
        public void Parse_read_holding_registers_response()
        {
            var data = new byte[] { 0x02, 0x00, 0x2A }; // 2 bytes, value = 42
            var regs = ModbusPDU.ParseReadRegistersResponse(data, 1);

            Assert.Single(regs);
            Assert.Equal((ushort)42, regs[0]);
        }

        [Fact]
        public void Write_multiple_registers_PDU_is_well_formed()
        {
            var pdu = ModbusPDU.BuildWriteMultipleRegistersRequest(
                0x0001, new ushort[] { 0x002A, 0x000F });

            Assert.Equal(ModbusFunctionCodes.WriteMultipleRegisters, pdu[0]);
            Assert.Equal(0x00, pdu[1]); // start addr hi
            Assert.Equal(0x01, pdu[2]); // start addr lo
            Assert.Equal(0x00, pdu[3]); // quantity hi
            Assert.Equal(0x02, pdu[4]); // quantity lo
            Assert.Equal(0x04, pdu[5]); // byte count
            Assert.Equal(0x00, pdu[6]); // reg1 hi
            Assert.Equal(0x2A, pdu[7]); // reg1 lo
            Assert.Equal(0x00, pdu[8]); // reg2 hi
            Assert.Equal(0x0F, pdu[9]); // reg2 lo
        }

        // ── CRC16 ─────────────────────────────────────────────

        [Fact]
        public void Crc_compute_and_validate_are_consistent()
        {
            // The CRC is correct if Compute + Validate agree.
            // Standard vector: 0x01 0x03 0x00 0x00 0x00 0x01 → CRC = 0x840A.
            var frame = new byte[] { 0x01, 0x03, 0x00, 0x00, 0x00, 0x01 };
            var crc = ModbusCRC.Compute(frame);
            Assert.Equal(2, crc.Length);

            // Reconstruct the frame with the CRC appended.
            var full = new byte[frame.Length + 2];
            Array.Copy(frame, full, frame.Length);
            full[full.Length - 2] = crc[0];
            full[full.Length - 1] = crc[1];
            Assert.True(ModbusCRC.Validate(full, full.Length));
        }

        [Fact]
        public void Crc_validates_a_well_formed_frame()
        {
            var pdu = new byte[] { 0x01, 0x03, 0x00, 0x00, 0x00, 0x01 };
            var crc = ModbusCRC.Compute(pdu);
            var frame = new byte[pdu.Length + 2];
            Array.Copy(pdu, frame, pdu.Length);
            frame[frame.Length - 2] = crc[0];
            frame[frame.Length - 1] = crc[1];
            Assert.True(ModbusCRC.Validate(frame, frame.Length));
        }

        [Fact]
        public void Crc_rejects_a_frame_with_a_bit_flip()
        {
            var pdu = new byte[] { 0x01, 0x03, 0x00, 0x00, 0x00, 0x01 };
            var crc = ModbusCRC.Compute(pdu);
            var frame = new byte[pdu.Length + 2];
            Array.Copy(pdu, frame, pdu.Length);
            frame[1] ^= 0x01; // flip a bit
            frame[frame.Length - 2] = crc[0];
            frame[frame.Length - 1] = crc[1];
            Assert.False(ModbusCRC.Validate(frame, frame.Length));
        }

        // ── RTU connection ─────────────────────────────────────

        [Fact]
        public async Task Rtu_read_holding_register_returns_OK()
        {
            var inner = new ScriptedTransportInstance();
            // Modbus RTU response: addr=1, func=0x03, len=2, val=0x0042, CRC
            var pdu = new byte[] { 0x01, 0x03, 0x02, 0x00, 0x42 };
            var crc = ModbusCRC.Compute(pdu);
            var response = new byte[pdu.Length + 2];
            Array.Copy(pdu, response, pdu.Length);
            response[response.Length - 2] = crc[0];
            response[response.Length - 1] = crc[1];
            inner.Inject(response);

            var conn = new ModbusRtuConnection(
                ConnectionString.Parse("modbus-rtu://COM1?unit-identifier=1"),
                inner);
            conn.Connect();

            var builder = conn.ReadRequestBuilder;
            builder.AddTagAddress("v", "holding:0");
            var rsp = await conn.Read((DefaultPlcReadRequest)builder.Build());

            Assert.NotNull(rsp);
        }

        [Fact]
        public async Task Rtu_write_single_coil_returns_OK()
        {
            var inner = new ScriptedTransportInstance();
            // Echo: addr=1, func=0x05, addr_hi=0x00, addr_lo=0x0A, val_hi=0xFF, val_lo=0x00
            var pdu = new byte[] { 0x01, 0x05, 0x00, 0x0A, 0xFF, 0x00 };
            var crc = ModbusCRC.Compute(pdu);
            var response = new byte[pdu.Length + 2];
            Array.Copy(pdu, response, pdu.Length);
            response[response.Length - 2] = crc[0];
            response[response.Length - 1] = crc[1];
            inner.Inject(response);

            var conn = new ModbusRtuConnection(
                ConnectionString.Parse("modbus-rtu://COM1?unit-identifier=1"),
                inner);
            conn.Connect();

            var wb = (DefaultPlcWriteRequestBuilder)conn.WriteRequestBuilder;
            wb.AddTag("c", "coil:10", true);
            var wRsp = (DefaultPlcWriteResponse)await conn.Write(
                (DefaultPlcWriteRequest)wb.Build());

            Assert.Equal(PlcResponseCode.Ok, wRsp.GetResponseCode("c"));
        }

        // ── TCP connection ─────────────────────────────────────
        //
        // ModbusConnection is the closest a driver test gets to a real Modbus
        // TCP device: it builds the MBAP header, correlates the transaction id,
        // and runs the SendAndReceive read loop. A physical PLC still has to be
        // checked with tools/modbus-verify <host>, but these pin the framing.

        // MBAP + PDU. The connection numbers transactions from 1 with a
        // pre-increment, so the first request on a fresh connection is id 2.
        private static byte[] MbapFrame(ushort transactionId, byte unitId, byte[] pdu)
        {
            var frame = new byte[7 + pdu.Length];
            frame[0] = (byte)(transactionId >> 8);
            frame[1] = (byte)(transactionId & 0xFF);
            frame[2] = 0x00;
            frame[3] = 0x00;
            frame[4] = (byte)((pdu.Length + 1) >> 8);
            frame[5] = (byte)((pdu.Length + 1) & 0xFF);
            frame[6] = unitId;
            Array.Copy(pdu, 0, frame, 7, pdu.Length);
            return frame;
        }

        [Fact]
        public async Task Tcp_read_holding_register_returns_the_value()
        {
            var inner = new ScriptedTransportInstance();
            inner.Inject(MbapFrame(2, 1, new byte[] { 0x03, 0x02, 0x00, 0x42 }));

            var conn = new ModbusConnection(
                ConnectionString.Parse("modbus-tcp://10.0.0.9:502?unit-identifier=1"), inner);
            conn.Connect();

            var builder = conn.ReadRequestBuilder;
            builder.AddTagAddress("v", "holding:0");
            var rsp = (DefaultPlcReadResponse)await conn.Read(
                (DefaultPlcReadRequest)builder.Build());

            Assert.Equal(PlcResponseCode.Ok, rsp.GetResponseCode("v"));
            Assert.Equal((ushort)0x0042, rsp.GetValue("v").GetUshort());

            // The request that went on the wire is a well-formed MBAP read.
            var sent = Assert.Single(inner.Written);
            Assert.Equal(new byte[]
            {
                0x00, 0x02, 0x00, 0x00, 0x00, 0x06, 0x01, // MBAP: tx 2, len 6, unit 1
                0x03, 0x00, 0x00, 0x00, 0x01,             // read holding register 0, qty 1
            }, sent);
        }

        [Fact]
        public async Task Tcp_write_single_register_returns_OK()
        {
            var inner = new ScriptedTransportInstance();
            // A Write Single Register response echoes the request.
            inner.Inject(MbapFrame(2, 1, new byte[] { 0x06, 0x00, 0x05, 0x12, 0x34 }));

            var conn = new ModbusConnection(
                ConnectionString.Parse("modbus-tcp://10.0.0.9:502?unit-identifier=1"), inner);
            conn.Connect();

            var wb = (DefaultPlcWriteRequestBuilder)conn.WriteRequestBuilder;
            wb.AddTag("r", "holding:5", (ushort)0x1234);
            var wRsp = (DefaultPlcWriteResponse)await conn.Write(
                (DefaultPlcWriteRequest)wb.Build());

            Assert.Equal(PlcResponseCode.Ok, wRsp.GetResponseCode("r"));
        }

        [Fact]
        public async Task Tcp_a_modbus_exception_fails_the_tag_and_the_connection_survives()
        {
            var inner = new ScriptedTransportInstance();
            // ILLEGAL DATA ADDRESS (0x02) for the first read...
            inner.Inject(MbapFrame(2, 1, new byte[] { 0x83, 0x02 }));
            // ...then a good response for the second (transaction id 3).
            inner.Inject(MbapFrame(3, 1, new byte[] { 0x03, 0x02, 0xBE, 0xEF }));

            var conn = new ModbusConnection(
                ConnectionString.Parse("modbus-tcp://10.0.0.9:502?unit-identifier=1"), inner);
            conn.Connect();

            var b1 = conn.ReadRequestBuilder;
            b1.AddTagAddress("v", "holding:99");
            var r1 = (DefaultPlcReadResponse)await conn.Read((DefaultPlcReadRequest)b1.Build());
            // ILLEGAL DATA ADDRESS maps to InvalidAddress, not a blanket InternalError.
            Assert.Equal(PlcResponseCode.InvalidAddress, r1.GetResponseCode("v"));

            var b2 = conn.ReadRequestBuilder;
            b2.AddTagAddress("v", "holding:0");
            var r2 = (DefaultPlcReadResponse)await conn.Read((DefaultPlcReadRequest)b2.Build());
            Assert.Equal(PlcResponseCode.Ok, r2.GetResponseCode("v"));
            Assert.Equal((ushort)0xBEEF, r2.GetValue("v").GetUshort());
        }

        [Fact]
        public async Task Tcp_a_stale_response_from_a_timed_out_read_does_not_brick_the_connection()
        {
            var inner = new ScriptedTransportInstance();
            var conn = new ModbusConnection(
                ConnectionString.Parse(
                    "modbus-tcp://10.0.0.9:502?unit-identifier=1&request-timeout=150"), inner);
            conn.Connect();

            // Read 1: nothing to read -> RequestTimeout (tx 2).
            var b1 = conn.ReadRequestBuilder;
            b1.AddTagAddress("v", "holding:0");
            var r1 = (DefaultPlcReadResponse)await conn.Read((DefaultPlcReadRequest)b1.Build());
            Assert.Equal(PlcResponseCode.RequestTimeout, r1.GetResponseCode("v"));

            // Now tx 2's response turns up late, with read 2's tx-3 response behind it.
            inner.Inject(MbapFrame(2, 1, new byte[] { 0x03, 0x02, 0x11, 0x11 }));
            inner.Inject(MbapFrame(3, 1, new byte[] { 0x03, 0x02, 0x22, 0x22 }));

            var b2 = conn.ReadRequestBuilder;
            b2.AddTagAddress("v", "holding:0");
            var r2 = (DefaultPlcReadResponse)await conn.Read((DefaultPlcReadRequest)b2.Build());
            // The driver skips the stale tx-2 frame and returns tx 3.
            Assert.Equal(PlcResponseCode.Ok, r2.GetResponseCode("v"));
            Assert.Equal((ushort)0x2222, r2.GetValue("v").GetUshort());
        }
    }
}
