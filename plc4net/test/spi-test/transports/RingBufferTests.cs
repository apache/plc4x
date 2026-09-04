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
using org.apache.plc4net.spi.transports;
using Xunit;

namespace org.apache.plc4net.spi.test.transports
{
    public class RingBufferTests
    {
        [Fact]
        public void Write_and_read_round_trip()
        {
            var buf = new RingBuffer(16);
            var data = new byte[] { 0x01, 0x02, 0x03, 0x04 };

            buf.Write(data);
            Assert.Equal(4, buf.AvailableForReading);
            Assert.Equal(12, buf.RemainingForWriting);

            var read = buf.Read(4);
            Assert.Equal(data, read);
        }

        [Fact]
        public void Peek_does_not_consume()
        {
            var buf = new RingBuffer(16);
            buf.Write(new byte[] { 0xAA, 0xBB });

            Assert.Equal(new byte[] { 0xAA, 0xBB }, buf.Peek(2));
            Assert.Equal(2, buf.AvailableForReading); // still there
            Assert.Equal(new byte[] { 0xAA, 0xBB }, buf.Read(2));
            Assert.Equal(0, buf.AvailableForReading);
        }

        [Fact]
        public void Wraps_around_the_capacity()
        {
            var buf = new RingBuffer(8);
            buf.Write(new byte[] { 1, 2, 3, 4, 5, 6 });
            buf.Read(3);
            Assert.Equal(3, buf.AvailableForReading);

            // This write should wrap around to the start.
            buf.Write(new byte[] { 7, 8, 9, 10 });
            Assert.Equal(7, buf.AvailableForReading);

            Assert.Equal(new byte[] { 4, 5, 6, 7, 8, 9, 10 }, buf.Read(7));
        }

        [Fact]
        public void Write_that_ends_exactly_on_the_capacity_boundary()
        {
            var buf = new RingBuffer(4);

            buf.Write(new byte[] { 1, 2, 3, 4 });
            Assert.Equal(0, buf.RemainingForWriting);
            Assert.Equal(new byte[] { 1, 2, 3, 4 }, buf.Read(4));

            // Both positions have just wrapped back to zero; the buffer stays usable.
            buf.Write(new byte[] { 5, 6, 7, 8 });
            Assert.Equal(new byte[] { 5, 6, 7, 8 }, buf.Read(4));
        }

        [Fact]
        public void Write_fills_the_whole_capacity_across_the_wrap()
        {
            var buf = new RingBuffer(8);
            buf.Write(new byte[] { 1, 2, 3, 4, 5, 6 });
            buf.Read(6); // both positions now sit at 6

            // 2 bytes fit before the end of the backing array, 6 wrap to the front.
            buf.Write(new byte[] { 10, 11, 12, 13, 14, 15, 16, 17 });
            Assert.Equal(0, buf.RemainingForWriting);

            Assert.Equal(new byte[] { 10, 11, 12, 13, 14, 15, 16, 17 }, buf.Read(8));
        }

        [Fact]
        public void Write_honours_the_offset_across_the_wrap()
        {
            var buf = new RingBuffer(4);
            buf.Write(new byte[] { 1, 2, 3 });
            buf.Read(3); // positions now sit at 3

            var source = new byte[] { 0xFF, 0xFF, 0xA0, 0xA1, 0xA2, 0xFF };
            buf.Write(source, 2, 3); // one byte before the wrap, two after

            Assert.Equal(new byte[] { 0xA0, 0xA1, 0xA2 }, buf.Read(3));
        }

        [Fact]
        public void Peek_spans_the_wrap_without_consuming()
        {
            var buf = new RingBuffer(4);
            buf.Write(new byte[] { 1, 2, 3 });
            buf.Read(3);
            buf.Write(new byte[] { 7, 8, 9 });

            Assert.Equal(new byte[] { 7, 8, 9 }, buf.Peek(3));
            Assert.Equal(3, buf.AvailableForReading);
            Assert.Equal(new byte[] { 7, 8, 9 }, buf.Read(3));
        }

        [Fact]
        public void Overflow_throws()
        {
            var buf = new RingBuffer(4);
            buf.Write(new byte[] { 1, 2, 3, 4 });

            Assert.Throws<TransportException>(() => buf.Write(new byte[] { 5 }));
        }

        [Fact]
        public void Read_more_than_available_throws()
        {
            var buf = new RingBuffer(8);
            buf.Write(new byte[] { 1, 2 });

            Assert.Throws<TransportException>(() => buf.Read(4));
        }

        [Fact]
        public void Zero_capacity_is_rejected()
        {
            Assert.Throws<ArgumentOutOfRangeException>(() => new RingBuffer(0));
            Assert.Throws<ArgumentOutOfRangeException>(() => new RingBuffer(-1));
        }

        [Fact]
        public void Clear_resets_the_buffer()
        {
            var buf = new RingBuffer(8);
            buf.Write(new byte[] { 1, 2, 3 });
            buf.Clear();

            Assert.Equal(0, buf.AvailableForReading);
            Assert.Equal(8, buf.RemainingForWriting);

            buf.Write(new byte[] { 0xFF });
            Assert.Equal(new byte[] { 0xFF }, buf.Read(1));
        }
    }
}
