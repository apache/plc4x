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
using org.apache.plc4net.drivers.slmp.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.drivers.slmp
{
    /// <summary>Frames 3E binary responses using the generated SLMP model.</summary>
    internal sealed class SlmpMessageCodec : MessageCodecBase<SlmpMessage>
    {
        private const int HeaderSize = 9;

        public SlmpMessageCodec(ITransportInstance transport, Action<SlmpMessage> handler)
            : base("SLMP", transport, handler)
        {
        }

        protected override int GetMinimumHeaderSize() => HeaderSize;

        public override void Send(SlmpMessage message)
        {
            if (message == null)
            {
                throw new ArgumentNullException(nameof(message));
            }
            // Both 3E request and response lengths exclude the nine-byte header.
            // Validate before serialization narrows that length to ushort.
            var length = message.GetLengthInBytes() - HeaderSize;
            if (length < 2 || length > ushort.MaxValue)
            {
                throw new MessageCodecException("SLMP 3E message length exceeds the wire field range.");
            }
            base.Send(message);
        }

        protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
        {
            // The generated response model requires this fixed access route.
            // The little-endian length includes the two-byte end code.
            var length = header[7] | (header[8] << 8);
            if (header[0] != 0xD0 || header[1] != 0 ||
                header[2] != 0 || header[3] != 0xFF ||
                header[4] != 0xFF || header[5] != 3 || header[6] != 0 || length < 2)
            {
                throw new MessageCodecException("Not a supported SLMP 3E response header.");
            }
            return HeaderSize + length;
        }

        protected override SlmpMessage ParseMessage(ReadBuffer readBuffer)
        {
            var message = SlmpMessage.StaticParse(readBuffer);
            if (message is not SlmpResponseFrame3E || readBuffer.HasMore(1))
            {
                throw new ParseException("Expected exactly one SLMP 3E response.");
            }
            return message;
        }
    }
}
