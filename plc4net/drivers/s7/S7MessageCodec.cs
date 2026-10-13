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
using org.apache.plc4net.drivers.s7.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.drivers.s7
{
    /// <summary>
    /// Frames S7 messages on the payload that the COTP transport hands up. An S7 message
    /// announces its own length: ten header bytes (protocol id, message type, reserved,
    /// PDU reference, parameter length, data length), two more for an acknowledgement (error
    /// class and code), then the parameter and the data. Once that many bytes are there the
    /// generated <see cref="S7Message"/> parses them; a first byte that is not the S7 protocol
    /// id or a length no CPU would announce makes the base class drop a byte and look again.
    /// </summary>
    internal sealed class S7MessageCodec : MessageCodecBase<S7Message>
    {
        internal const byte ProtocolId = 0x32;
        internal const byte MessageTypeAck = 0x02;
        internal const byte MessageTypeAckData = 0x03;

        private const int HeaderSize = 10;
        private const int AckHeaderSize = 12;

        // Before Setup Communication completes, the requested PDU length is the largest message
        // this connection can accept. S7Connection replaces it with the CPU's negotiated value.
        // An announced length beyond that limit is noise; waiting for it would hold up the real
        // answer until the timeout.
        private int _maximumMessageSize;

        public S7MessageCodec(ITransportInstance transport, Action<S7Message> handler, int maximumMessageSize)
            : base("S7", transport, handler)
        {
            MaximumMessageSize = maximumMessageSize;
        }

        internal int MaximumMessageSize
        {
            get => _maximumMessageSize;
            set
            {
                if (value <= 0)
                {
                    throw new ArgumentOutOfRangeException(nameof(value), "The S7 message-size limit must be positive.");
                }
                _maximumMessageSize = value;
            }
        }

        protected override int GetMinimumHeaderSize() => HeaderSize;

        protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
        {
            if (header[0] != ProtocolId)
            {
                throw new MessageCodecException(
                    $"Not an S7 message: the first byte is 0x{header[0]:X2}, not 0x{ProtocolId:X2}.");
            }

            // An acknowledgement (0x02) and an acknowledgement with data (0x03) carry the error
            // class and the error code after the ten common bytes; a job and user data do not.
            // Framing a bare acknowledgement as ten bytes would leave the two behind and shift
            // every message after it.
            var headerSize = header[1] == MessageTypeAck || header[1] == MessageTypeAckData
                ? AckHeaderSize
                : HeaderSize;
            var parameterLength = (header[6] << 8) | header[7];
            var dataLength = (header[8] << 8) | header[9];
            var total = headerSize + parameterLength + dataLength;
            if (total > MaximumMessageSize)
            {
                throw new MessageCodecException($"Implausible S7 message length {total}.");
            }
            return total;
        }

        protected override S7Message ParseMessage(ReadBuffer readBuffer)
            => S7Message.StaticParse(readBuffer);
    }
}
