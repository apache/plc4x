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
using System.Text;
using org.apache.plc4net.spi.generation;
using Xunit;

namespace org.apache.plc4net.spi.test.generation
{
    /// <summary>
    /// Round-trip coverage for the bit-level codec every generated driver depends on.
    /// </summary>
    public class BufferTests
    {
        private static ReadBuffer Reader(Action<WriteBuffer> write)
        {
            var wb = new WriteBuffer();
            write(wb);
            return new ReadBuffer(wb.GetBytes());
        }

        [Fact]
        public void Bits_are_written_and_read_most_significant_first()
        {
            var wb = new WriteBuffer();
            wb.WriteBit("", true);
            wb.WriteBit("", false);
            wb.WriteBit("", true);
            wb.WriteBit("", false);
            wb.WriteBit("", false);
            wb.WriteBit("", false);
            wb.WriteBit("", false);
            wb.WriteBit("", true);

            Assert.Equal(new byte[] { 0xA1 }, wb.GetBytes());
        }

        [Theory]
        [InlineData(0UL)]
        [InlineData(1UL)]
        [InlineData(0x00000000FFFFFFFFUL)]
        [InlineData(0xFFFFFFFF00000000UL)]
        [InlineData(0x0123456789ABCDEFUL)]
        [InlineData(ulong.MaxValue)]
        public void Ulong_round_trips_at_64_bits(ulong value)
        {
            // The reader used to consume (bitLength - 32) + bitLength bits, i.e. 96
            // bits for a 64 bit field, so this never round-tripped.
            var rb = Reader(wb => wb.WriteUlong("", 64, value));

            Assert.Equal(value, rb.ReadUlong("", 64));
        }

        [Theory]
        [InlineData(0L)]
        [InlineData(-1L)]
        [InlineData(long.MinValue)]
        [InlineData(long.MaxValue)]
        public void Long_round_trips_at_64_bits(long value)
        {
            var rb = Reader(wb => wb.WriteLong("", 64, value));

            Assert.Equal(value, rb.ReadLong("", 64));
        }

        [Fact]
        public void Writer_emits_exactly_the_declared_width()
        {
            var wb = new WriteBuffer();
            wb.WriteUlong("", 64, ulong.MaxValue);

            Assert.Equal(64, wb.GetPos());
            Assert.Equal(8, wb.GetBytes().Length);
        }

        [Theory]
        [InlineData(-8, 4)]   // sign bit set in a narrow field
        [InlineData(7, 4)]
        [InlineData(-1, 12)]
        [InlineData(-32768, 16)]
        public void Signed_fields_are_sign_extended(int value, int bitLength)
        {
            var rb = Reader(wb => wb.WriteInt("", bitLength, value));

            Assert.Equal(value, rb.ReadInt("", bitLength));
        }

        [Fact]
        public void Unaligned_fields_pack_back_to_back()
        {
            var rb = Reader(wb =>
            {
                wb.WriteBit("", true);
                wb.WriteUshort("", 12, 0x0ABC);
                wb.WriteByte("", 3, 0x05);
            });

            Assert.True(rb.ReadBit(""));
            Assert.Equal((ushort)0x0ABC, rb.ReadUshort("", 12));
            Assert.Equal((byte)0x05, rb.ReadByte("", 3));
        }

        [Fact]
        public void Float32_round_trips()
        {
            var rb = Reader(wb => wb.WriteFloat("", 32, 22.0f));

            Assert.Equal(22.0f, rb.ReadFloat("", 32));
        }

        [Fact]
        public void Float32_is_read_big_endian_off_the_wire()
        {
            // 0x41B00000 is 22.0f in IEEE-754, most significant byte first.
            var rb = new ReadBuffer(new byte[] { 0x41, 0xB0, 0x00, 0x00 });

            Assert.Equal(22.0f, rb.ReadFloat("", 32));
        }

        [Fact]
        public void Double64_round_trips()
        {
            var rb = Reader(wb => wb.WriteDouble("", 64, -1234.5678d));

            Assert.Equal(-1234.5678d, rb.ReadDouble("", 64), 9);
        }

        [Fact]
        public void Double_at_32_bits_reads_as_a_float()
        {
            // Previously handed a 4 byte array to BitConverter.ToDouble and threw.
            var rb = Reader(wb => wb.WriteFloat("", 32, 1.5f));

            Assert.Equal(1.5d, rb.ReadDouble("", 32), 6);
        }

        [Theory]
        [InlineData(0.0f)]
        [InlineData(22.0f)]
        [InlineData(-5.12f)]
        [InlineData(20.48f)]
        [InlineData(-20.48f)]
        public void Knx_16bit_float_round_trips(float value)
        {
            // These all land exactly on the 0.01 * mantissa * 2^exp grid; see
            // Knx_16bit_float_is_coarse_at_the_top_of_its_range for values that
            // do not.
            // The 16 bit KNX branch of WriteFloat used to be an empty block that
            // emitted nothing at all.
            var rb = Reader(wb => wb.WriteFloat("", 16, value));

            Assert.Equal(value, rb.ReadFloat("", 16), 2);
        }

        [Fact]
        public void Knx_16bit_float_is_coarse_at_the_top_of_its_range()
        {
            // DPT 9.x stores 0.01 * mantissa * 2^exp with an 11 bit mantissa, so at
            // exp 15 the smallest step is 0.01 * 2^15 = 327.68. Values near the
            // 670760 maximum therefore snap to the nearest multiple of that step;
            // this is the format's resolution, not a rounding defect.
            var rb = Reader(wb => wb.WriteFloat("", 16, 670760.0f));

            Assert.Equal(670760.0f, rb.ReadFloat("", 16), 327.68f);
        }

        [Fact]
        public void String_round_trips_and_is_padded_to_the_field_width()
        {
            var rb = Reader(wb => wb.WriteString("", 64, "UTF8", "plc4x"));

            Assert.Equal("plc4x", rb.ReadString("", 64, Encoding.UTF8));
        }

        [Fact]
        public void Byte_array_round_trips()
        {
            var payload = new byte[] { 0xDE, 0xAD, 0xBE, 0xEF };
            var rb = Reader(wb => wb.WriteByteArray("", payload));

            Assert.Equal(payload, rb.ReadByteArray("", 32));
        }

        [Fact]
        public void HasMore_allows_a_read_that_exactly_fills_the_buffer()
        {
            var rb = new ReadBuffer(new byte[] { 0x00, 0x00 });

            Assert.True(rb.HasMore(16));  // was false: the check used < instead of <=
            Assert.False(rb.HasMore(17));
        }

        [Fact]
        public void HasMore_rejects_a_negative_bit_length()
        {
            // A negative bit length is nonsensical and should not silently return true
            // just because the comparison would hold for any non-negative remaining.
            var rb = new ReadBuffer(new byte[] { 0x00, 0x00 });

            Assert.False(rb.HasMore(-1));
            Assert.False(rb.HasMore(int.MinValue));
        }

        [Fact]
        public void Reading_past_the_end_reports_a_parse_error()
        {
            var rb = new ReadBuffer(new byte[] { 0x01 });

            var ex = Assert.Throws<ParseException>(() => rb.ReadUint("", 32));
            Assert.Contains("Not enough data", ex.Message);
        }

        // ── little-endian byte order ────────────────────────────

        private static ReadBuffer LittleEndianReader(Action<WriteBuffer> write)
        {
            var wb = new WriteBuffer();
            wb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            write(wb);
            wb.PopByteOrder();
            var rb = new ReadBuffer(wb.GetBytes());
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            return rb;
        }

        [Fact]
        public void Default_byte_order_is_still_big_endian()
        {
            // No push at all - every existing S7 / Modbus call site. Locks in
            // that adding the byte-order stack changed nothing by default.
            var rb = new ReadBuffer(new byte[] { 0x12, 0x34 });

            Assert.Equal((ushort)0x1234, rb.ReadUshort("", 16));
        }

        [Fact]
        public void Little_endian_16_bit_unsigned_is_read_off_the_wire_byte_reversed()
        {
            // 0x34, 0x12 little-endian is 0x1234, not 0x3412 - a wrong
            // direction swap or a no-op would both produce a different value.
            var rb = new ReadBuffer(new byte[] { 0x34, 0x12 });
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);

            Assert.Equal((ushort)0x1234, rb.ReadUshort("", 16));
        }

        [Fact]
        public void Little_endian_24_bit_unsigned_reverses_exactly_three_bytes()
        {
            // SLMP's device number: 3 bytes, little-endian. The unused high
            // byte of the uint must not end up in the swap.
            var rb = new ReadBuffer(new byte[] { 0x03, 0x02, 0x01 });
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);

            Assert.Equal(0x010203u, rb.ReadUint("", 24));
        }

        [Fact]
        public void Little_endian_32_bit_and_64_bit_unsigned_reverse_their_whole_width()
        {
            var rb32 = new ReadBuffer(new byte[] { 0x04, 0x03, 0x02, 0x01 });
            rb32.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            Assert.Equal(0x01020304u, rb32.ReadUint("", 32));

            var rb64 = new ReadBuffer(new byte[] { 0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01 });
            rb64.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            Assert.Equal(0x0102030405060708u, rb64.ReadUlong("", 64));
        }

        [Fact]
        public void Little_endian_8_bit_fields_are_unaffected()
        {
            // A single byte has no order to reverse.
            var rb = new ReadBuffer(new byte[] { 0xAB });
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);

            Assert.Equal((byte)0xAB, rb.ReadByte("", 8));
        }

        [Theory]
        [InlineData((ushort)0x1234)]
        [InlineData((ushort)0)]
        [InlineData(ushort.MaxValue)]
        public void Little_endian_16_bit_unsigned_round_trips(ushort value)
        {
            var rb = LittleEndianReader(wb => wb.WriteUshort("", 16, value));

            Assert.Equal(value, rb.ReadUshort("", 16));
        }

        [Theory]
        [InlineData(0x010203u)]
        [InlineData(0u)]
        [InlineData(0xFFFFFFu)]
        public void Little_endian_24_bit_unsigned_round_trips(uint value)
        {
            var rb = LittleEndianReader(wb => wb.WriteUint("", 24, value));

            Assert.Equal(value, rb.ReadUint("", 24));
        }

        [Fact]
        public void Little_endian_write_then_default_read_produces_the_reversed_wire_bytes()
        {
            // Proves the WriteBuffer side swaps too, independently of ReadBuffer:
            // write 0x1234 little-endian, then read the raw wire bytes back
            // big-endian (no push) and see them reversed.
            var wb = new WriteBuffer();
            wb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            wb.WriteUshort("", 16, 0x1234);
            wb.PopByteOrder();

            Assert.Equal(new byte[] { 0x34, 0x12 }, wb.GetBytes());
        }

        [Theory]
        [InlineData(-1, 16)]
        [InlineData(-2, 16)]
        [InlineData(short.MinValue, 16)]
        [InlineData(1, 24)]
        [InlineData(-1, 24)]
        public void Little_endian_signed_fields_round_trip_sign_extended(int value, int bitLength)
        {
            var rb = LittleEndianReader(wb => wb.WriteInt("", bitLength, value));

            Assert.Equal(value, rb.ReadInt("", bitLength));
        }

        [Fact]
        public void Little_endian_negative_16_bit_signed_is_read_off_the_wire_byte_reversed()
        {
            // -2 as a 16 bit two's-complement value is 0xFFFE; little-endian on
            // the wire that is [0xFE, 0xFF] - the low byte first.
            var rb = new ReadBuffer(new byte[] { 0xFE, 0xFF });
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);

            Assert.Equal(-2, rb.ReadInt("", 16));
        }

        [Fact]
        public void Pop_restores_the_byte_order_the_matching_push_replaced()
        {
            var rb = new ReadBuffer(new byte[] { 0x34, 0x12, 0x12, 0x34 });
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);
            Assert.Equal((ushort)0x1234, rb.ReadUshort("", 16)); // little-endian
            rb.PopByteOrder();

            Assert.Equal((ushort)0x1234, rb.ReadUshort("", 16)); // back to default big-endian
        }

        [Fact]
        public void Little_endian_floating_point_fails_loudly_instead_of_reading_big_endian()
        {
            var rb = new ReadBuffer(new byte[] { 0x41, 0xB0, 0x00, 0x00 });
            rb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);

            Assert.Throws<NotImplementedException>(() => rb.ReadFloat("", 32));
        }

        [Fact]
        public void Little_endian_floating_point_write_fails_loudly_too()
        {
            var wb = new WriteBuffer();
            wb.PushByteOrder(ByteOrder.LITTLE_ENDIAN);

            Assert.Throws<NotImplementedException>(() => wb.WriteFloat("", 32, 1.0f));
        }

        [Fact]
        public void A_window_reads_the_selected_bytes_without_crossing_its_end()
        {
            var rb = new ReadBuffer(new byte[] { 0xAA, 0x12, 0x34, 0xBB }, 1, 2);

            Assert.Equal(2, rb.GetTotalBytes());
            Assert.Equal(0, rb.GetPos());
            Assert.Equal((byte)0x12, rb.PeekByte(0));
            Assert.Equal((ushort)0x1234, rb.ReadUshort("", 16));
            Assert.Equal(16, rb.GetPos());
            Assert.False(rb.HasMore(1));
            Assert.Throws<ParseException>(() => rb.ReadByte("", 8));

            rb.Reset();
            Assert.Equal((byte)0x12, rb.ReadByte("", 8));
        }
    }
}