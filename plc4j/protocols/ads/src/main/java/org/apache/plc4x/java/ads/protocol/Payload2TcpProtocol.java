/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ads.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.ads.api.tcp.AmsTCPPacket;
import org.apache.plc4x.java.ads.api.tcp.AmsTcpHeader;
import org.apache.plc4x.java.ads.api.tcp.types.TcpLength;
import org.apache.plc4x.java.ads.api.tcp.types.UserData;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class Payload2TcpProtocol extends MessageToMessageCodec<ByteBuf, ByteBuf> {
    private ByteBuf retainingBuf = Unpooled.EMPTY_BUFFER;

    private static final Logger LOGGER = LoggerFactory.getLogger(Payload2TcpProtocol.class);

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, ByteBuf amsPacket, List<Object> out) {
        LOGGER.trace("(<--OUT): {}, {}, {}", channelHandlerContext, amsPacket, out);
        out.add(AmsTCPPacket.of(UserData.of(amsPacket)).getByteBuf());
    }

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) {
        if (byteBuf == Unpooled.EMPTY_BUFFER) {
            // Cleanup...
            return;
        }
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, byteBuf, out);

        retainingBuf = Unpooled.wrappedBuffer(retainingBuf, byteBuf.retain());

        while (retainingBuf.readableBytes() >= AmsTcpHeader.Reserved.NUM_BYTES + TcpLength.NUM_BYTES
                && retainingBuf.readableBytes() >= retainingBuf.getUnsignedIntLE(retainingBuf.readerIndex()
                        + AmsTcpHeader.Reserved.NUM_BYTES) + AmsTcpHeader.Reserved.NUM_BYTES + TcpLength.NUM_BYTES) {
            // Reserved
            retainingBuf.skipBytes(AmsTcpHeader.Reserved.NUM_BYTES);
            TcpLength packetLength = TcpLength.of(retainingBuf);
            AmsTcpHeader amsTcpHeader = AmsTcpHeader.of(packetLength);
            LOGGER.debug("AMS TCP Header {}", amsTcpHeader);

            out.add(retainingBuf.readBytes((int) packetLength.getAsLong()));
        }
        if (retainingBuf.readableBytes() == 0) {
            retainingBuf.release();
            retainingBuf = Unpooled.EMPTY_BUFFER;
        }
    }

}
