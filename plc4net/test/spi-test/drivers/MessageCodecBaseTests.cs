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
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;
using Xunit;

namespace org.apache.plc4net.spi.test.drivers
{
    /// <summary>
    /// The framing loop shared by every codec: partial frames, several frames in one
    /// read, and - the point of these tests - resynchronising after the byte stream is
    /// left misaligned (a dropped datagram, a truncated or malformed frame).
    /// </summary>
    public class MessageCodecBaseTests
    {
        // 4-byte header: 0xAA 0xBB then a big-endian total length (header + body).
        private sealed class ToyCodec : MessageCodecBase<byte[]>
        {
            public ToyCodec(ITransportInstance transport, Action<byte[]> handler)
                : base("toy", transport, handler)
            {
            }

            protected override int GetMinimumHeaderSize() => 4;

            protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
            {
                if (header[0] != 0xAA || header[1] != 0xBB)
                {
                    throw new MessageCodecException("bad magic");
                }
                var total = (header[2] << 8) | header[3];
                if (total < 4)
                {
                    throw new MessageCodecException("length below header size");
                }
                return total;
            }

            protected override byte[] ParseMessage(ReadBuffer readBuffer) => readBuffer.GetBytes();
        }

        private static byte[] Frame(params byte[] body)
        {
            var total = 4 + body.Length;
            var frame = new byte[total];
            frame[0] = 0xAA;
            frame[1] = 0xBB;
            frame[2] = (byte) (total >> 8);
            frame[3] = (byte) total;
            Array.Copy(body, 0, frame, 4, body.Length);
            return frame;
        }

        private static (ToyCodec codec, List<byte[]> received) Build(byte[] buffered)
        {
            var transport = new BufferBackedTransport(buffered);
            var received = new List<byte[]>();
            return (new ToyCodec(transport, received.Add), received);
        }

        [Fact]
        public void Delivers_a_whole_frame()
        {
            var (codec, received) = Build(Frame(1, 2, 3));
            codec.ProcessIncomingData();
            Assert.Single(received);
            Assert.Equal(new byte[] { 0xAA, 0xBB, 0x00, 0x07, 1, 2, 3 }, received[0]);
        }

        [Fact]
        public void Leaves_a_partial_frame_buffered()
        {
            var whole = Frame(1, 2, 3);
            var (codec, received) = Build(new[] { whole[0], whole[1], whole[2], whole[3], whole[4] });
            codec.ProcessIncomingData();
            Assert.Empty(received);
        }

        [Fact]
        public void Resyncs_past_leading_garbage_to_the_next_valid_frame()
        {
            var junk = new byte[] { 0x01, 0x02, 0x03, 0x04, 0x05 };
            var good = Frame(9, 9);
            var buffer = new byte[junk.Length + good.Length];
            Array.Copy(junk, buffer, junk.Length);
            Array.Copy(good, 0, buffer, junk.Length, good.Length);

            var (codec, received) = Build(buffer);
            codec.ProcessIncomingData();

            Assert.Single(received);
            Assert.Equal(good, received[0]);
        }

        [Fact]
        public void Resyncs_past_a_frame_with_an_impossible_length()
        {
            // A valid magic but a length of 2 (below the 4-byte header): the stream is
            // misaligned. The codec must skip it, not consume 2 bytes and stay off.
            var bad = new byte[] { 0xAA, 0xBB, 0x00, 0x02 };
            var good = Frame(7);
            var buffer = new byte[bad.Length + good.Length];
            Array.Copy(bad, buffer, bad.Length);
            Array.Copy(good, 0, buffer, bad.Length, good.Length);

            var (codec, received) = Build(buffer);
            codec.ProcessIncomingData();

            Assert.Single(received);
            Assert.Equal(good, received[0]);
        }

        /// <summary>An <see cref="ITransportInstance"/> that just serves a fixed byte buffer.</summary>
        private sealed class BufferBackedTransport : ITransportInstance
        {
            private readonly List<byte> _buffer;

            public BufferBackedTransport(byte[] initial) => _buffer = new List<byte>(initial);

            public ITransportConfiguration Configuration => null;
            public string DriverConfig => string.Empty;
            public bool IsOpen => true;

            public int GetNumBytesAvailable() => _buffer.Count;

            public byte[] PeekReadableBytes(int numBytes)
            {
                if (_buffer.Count < numBytes)
                {
                    throw new TransportException("not enough buffered");
                }
                return _buffer.GetRange(0, numBytes).ToArray();
            }

            public byte[] Read(int numBytes)
            {
                var slice = PeekReadableBytes(numBytes);
                _buffer.RemoveRange(0, numBytes);
                return slice;
            }

            public void Write(byte[] bytes)
            {
            }

            public void Close()
            {
            }

            public void Dispose()
            {
            }
        }
    }
}
