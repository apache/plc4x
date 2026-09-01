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
using org.apache.plc4net.drivers.knxnetip.readwrite.model;
using org.apache.plc4net.spi.drivers;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.transports;

namespace org.apache.plc4net.drivers.knxnetip
{
    /// <summary>
    /// Framing for the KNXnet/IP wire format.
    /// </summary>
    /// <remarks>
    /// Every frame opens with a 6-byte common header - <c>0x06 0x10</c>, a 16-bit
    /// service type, then a 16-bit total length that counts the whole frame. UDP
    /// datagrams are already message-bounded, but reading the length from the header
    /// keeps the codec transport-agnostic.
    /// </remarks>
    public sealed class KnxNetIpMessageCodec : MessageCodecBase<KnxNetIpMessage>
    {
        public KnxNetIpMessageCodec(ITransportInstance transportInstance, Action<KnxNetIpMessage> messageHandler)
            : base("KNXnet/IP", transportInstance, messageHandler)
        {
        }

        protected override int GetMinimumHeaderSize() => 6;

        protected override int CalculateTotalMessageSize(byte[] header, int availableBytes)
        {
            if (header[0] != 0x06 || header[1] != 0x10)
            {
                throw new MessageCodecException(
                    $"Not a KNXnet/IP frame: header 0x{header[0]:X2}{header[1]:X2} (expected 0x0610).");
            }
            var total = (header[4] << 8) | header[5];
            if (total < 6)
            {
                // A length shorter than the header itself means the stream is
                // misaligned; reject it so the codec resyncs rather than consuming
                // a partial frame and staying off by a few bytes forever.
                throw new MessageCodecException(
                    $"KNXnet/IP frame length {total} is shorter than the 6-byte header.");
            }
            return total;
        }

        protected override KnxNetIpMessage ParseMessage(ReadBuffer readBuffer) =>
            KnxNetIpMessage.StaticParse(readBuffer);
    }
}
