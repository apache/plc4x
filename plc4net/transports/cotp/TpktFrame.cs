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

namespace org.apache.plc4net.transports.cotp
{
    /// <summary>
    /// RFC 1006 TPKT (ISO-on-TCP) frame.
    /// A thin 4-byte header that frames higher-level protocol data over TCP.
    /// </summary>
    public static class TpktFrame
    {
        public const byte Version = 3;
        public const int HeaderSize = 4;

        /// <summary>
        /// The TPKT length field is 16 bits wide and counts the header itself, so a frame
        /// can be at most 65535 bytes long.
        /// </summary>
        public const int MaxFrameSize = ushort.MaxValue;

        /// <summary>The most payload a single TPKT frame can carry.</summary>
        public const int MaxPayloadSize = MaxFrameSize - HeaderSize;

        /// <summary>
        /// Reads the payload length from a TPKT header (bytes 2-3, big-endian).
        /// The length includes the 4-byte header itself.
        /// </summary>
        public static int ReadPayloadLength(byte[] header)
        {
            if (header.Length < 4) return -1;
            // Bytes 2-3: total length including header.
            return ((header[2] << 8) | header[3]) - HeaderSize;
        }

        /// <summary>
        /// Wraps a payload in a TPKT frame and returns the complete byte array.
        /// </summary>
        /// <exception cref="ArgumentException">
        /// The payload is longer than <see cref="MaxPayloadSize"/>: the frame length would not
        /// fit the 16-bit length field and the header would describe a different size than the
        /// frame actually has.
        /// </exception>
        public static byte[] Wrap(byte[] payload)
        {
            var payloadLength = payload?.Length ?? 0;
            if (payloadLength > MaxPayloadSize)
            {
                throw new ArgumentException(
                    $"A TPKT frame carries at most {MaxPayloadSize} payload bytes, but the payload has {payloadLength}.",
                    nameof(payload));
            }

            var totalLength = HeaderSize + payloadLength;
            var frame = new byte[totalLength];
            frame[0] = Version;
            frame[1] = 0; // reserved
            frame[2] = (byte)(totalLength >> 8);
            frame[3] = (byte)(totalLength & 0xFF);
            if (payload != null)
                Array.Copy(payload, 0, frame, HeaderSize, payload.Length);
            return frame;
        }

        /// <summary>Extracts the payload from a TPKT frame.</summary>
        public static byte[] Unwrap(byte[] frame)
        {
            if (frame.Length < HeaderSize) return Array.Empty<byte>();
            var payloadLen = frame.Length - HeaderSize;
            var payload = new byte[payloadLen];
            Array.Copy(frame, HeaderSize, payload, 0, payloadLen);
            return payload;
        }
    }
}
