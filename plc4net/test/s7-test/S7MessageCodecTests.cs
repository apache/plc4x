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
using System.Linq;
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.transports.test;
using Xunit;

namespace org.apache.plc4net.drivers.s7.test
{
    /// <summary>
    /// How the codec cuts the payload that the COTP transport hands up into whole S7 messages.
    /// The frames are built by hand in <see cref="S7Frames"/>, not by the generated model.
    /// </summary>
    public class S7MessageCodecTests
    {
        private static (S7MessageCodec Codec, List<S7Message> Received, TestTransportInstance Transport) New(
            int maximumMessageSize = 1024)
        {
            var transport = new TestTransportInstance(new TestTransportConfiguration());
            var received = new List<S7Message>();
            return (new S7MessageCodec(transport, received.Add, maximumMessageSize), received, transport);
        }

        // The S7 message inside a frame that S7Frames wraps for the transport.
        private static byte[] Unwrapped(byte[] frame) => frame[7..];

        private static byte[] ReadAckMessage(ushort reference, byte[]? data = null)
            => Unwrapped(S7Frames.ReadAck(reference, (0xFF, 0x04, data ?? new byte[] { 0x12, 0x34 })));

        [Fact]
        public void An_acknowledgement_with_data_is_framed_by_the_lengths_in_its_header()
        {
            var (codec, received, transport) = New();
            transport.InjectTestData(ReadAckMessage(7));

            codec.ProcessIncomingData();

            var message = Assert.IsType<S7MessageResponseData>(Assert.Single(received));
            Assert.Equal((ushort)7, message.TpduReference);
            Assert.Equal(0, transport.GetNumBytesAvailable());
        }

        [Fact]
        public void A_message_that_has_not_arrived_in_full_is_waited_for()
        {
            var (codec, received, transport) = New();
            var message = ReadAckMessage(7);

            transport.InjectTestData(message[..15]);
            codec.ProcessIncomingData();
            Assert.Empty(received);
            Assert.Equal(15, transport.GetNumBytesAvailable()); // kept, not dropped byte by byte

            transport.InjectTestData(message[15..]);
            codec.ProcessIncomingData();
            Assert.Single(received);
        }

        [Fact]
        public void Several_messages_in_one_read_are_all_delivered()
        {
            var (codec, received, transport) = New();
            transport.InjectTestData(ReadAckMessage(1).Concat(ReadAckMessage(2)).Concat(ReadAckMessage(3)).ToArray());

            codec.ProcessIncomingData();

            Assert.Equal(new ushort[] { 1, 2, 3 }, received.Select(m => ((S7MessageResponseData)m).TpduReference));
        }

        [Fact]
        public void A_bare_acknowledgement_has_twelve_header_bytes_and_does_not_shift_the_message_after_it()
        {
            // The error class and code follow the ten common bytes in a bare acknowledgement just
            // as in one with data. Framing it as ten bytes would leave two behind, and the next
            // message would be read two bytes off.
            var (codec, received, transport) = New();
            transport.InjectTestData(S7Frames.Refusal(3, 0x81, 0x04).Concat(ReadAckMessage(4)).ToArray());

            codec.ProcessIncomingData();

            Assert.Equal(2, received.Count);
            var refusal = Assert.IsType<S7MessageResponse>(received[0]);
            Assert.Equal((ushort)3, refusal.TpduReference);
            Assert.Equal((byte)0x81, refusal.ErrorClass);
            Assert.Equal((byte)0x04, refusal.ErrorCode);
            Assert.Equal((ushort)4, Assert.IsType<S7MessageResponseData>(received[1]).TpduReference);
        }

        [Fact]
        public void A_job_has_ten_header_bytes()
        {
            var (codec, received, transport) = New();
            transport.InjectTestData(Unwrapped(S7Frames.ReadRequest(
                9, S7Frames.AddressItem(0x05, 1, 0x84, 10, 0))));

            codec.ProcessIncomingData();

            Assert.Equal((ushort)9, Assert.IsType<S7MessageRequest>(Assert.Single(received)).TpduReference);
        }

        [Fact]
        public void Bytes_that_do_not_start_a_message_are_dropped_one_at_a_time()
        {
            var (codec, received, transport) = New();
            transport.InjectTestData(new byte[] { 0x00, 0x99, 0x17 }.Concat(ReadAckMessage(5)).ToArray());

            codec.ProcessIncomingData();

            Assert.Equal((ushort)5, Assert.IsType<S7MessageResponseData>(Assert.Single(received)).TpduReference);
        }

        [Fact]
        public void A_length_that_no_cpu_announces_is_not_waited_for()
        {
            // 32 03 ... announcing 65535 bytes of parameter and of data would hold up the real
            // message behind it until the timeout.
            var (codec, received, transport) = New();
            var noise = new byte[] { 0x32, 0x03, 0x00, 0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF };
            transport.InjectTestData(noise.Concat(ReadAckMessage(6)).ToArray());

            codec.ProcessIncomingData();

            Assert.Equal((ushort)6, Assert.IsType<S7MessageResponseData>(Assert.Single(received)).TpduReference);
        }

        [Fact]
        public void A_message_larger_than_the_negotiated_limit_is_not_waited_for()
        {
            var message = ReadAckMessage(7);
            var announcedSize = message.Length + 1;
            var (codec, received, transport) = New(message.Length);
            var noise = new byte[]
            {
                0x32, 0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00,
                (byte)((announcedSize - 12) >> 8), (byte)(announcedSize - 12),
            };
            transport.InjectTestData(noise.Concat(message).ToArray());

            codec.ProcessIncomingData();

            Assert.Equal((ushort)7, Assert.IsType<S7MessageResponseData>(Assert.Single(received)).TpduReference);
        }

        [Fact]
        public void A_message_that_does_not_parse_is_an_error_not_a_resync()
        {
            // Well framed, but a parameter type that s7.mspec does not know.
            var (codec, received, transport) = New();
            transport.InjectTestData(S7Frames.AckData(8, new byte[] { 0xEE, 0x00 }, Array.Empty<byte>()));

            Assert.Throws<MessageCodecException>(() => codec.ProcessIncomingData());
            Assert.Empty(received);
        }
    }
}
