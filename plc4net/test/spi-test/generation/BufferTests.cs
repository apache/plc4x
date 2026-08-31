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
            Assert.Equal((ushort) 0x0ABC, rb.ReadUshort("", 12));
            Assert.Equal((byte) 0x05, rb.ReadByte("", 3));
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
    }
}
