/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.amsads.attic.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.amsads.protocol.exception.AdsException;
import org.apache.plc4x.java.amsads.readwrite.AmsPacket;
import org.apache.plc4x.java.amsads.readwrite.AmsTCPPacket;
import org.apache.plc4x.java.amsads.readwrite.io.AmsPacketIO;
import org.apache.plc4x.java.amsads.readwrite.io.AmsTCPPacketIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

@Deprecated
public class Payload2TcpProtocol extends MessageToMessageCodec<ByteBuf, ByteBuf> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Payload2TcpProtocol.class);

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, ByteBuf amsPacket, List<Object> out) throws AdsException {
        LOGGER.trace("(<--OUT): {}, {}, {}", channelHandlerContext, amsPacket, out);

        byte[] bytes = amsPacket.array();
        AmsPacket amsPacketSer;
        try {
            amsPacketSer = AmsPacketIO.staticParse(new ReadBuffer(bytes, true));
        } catch (ParseException e) {
            throw new AdsException(-1L, e);
        }

        WriteBuffer writeBuffer = new WriteBuffer(amsPacketSer.getLengthInBytes(), true);
        try {
            AmsTCPPacketIO.staticSerialize(writeBuffer, new AmsTCPPacket(amsPacketSer));
        } catch (ParseException e) {
            throw new AdsException(amsPacketSer.getInvokeId(), e);
        }
        out.add(writeBuffer.getData());
    }

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) {
        if (byteBuf == Unpooled.EMPTY_BUFFER) {
            // Cleanup...
            return;
        }
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, byteBuf, out);
        byte[] bytes = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(bytes);
        ReadBuffer readBuffer = new ReadBuffer(bytes);
        while (readBuffer.getPos() < bytes.length) {
            try {
                AmsTCPPacket amsTCPPacket = AmsTCPPacketIO.staticParse(readBuffer);
                AmsPacket amsPacket = amsTCPPacket.getUserdata();

                WriteBuffer writeBuffer = new WriteBuffer(amsPacket.getLengthInBytes(), true);
                try {
                    AmsPacketIO.staticSerialize(writeBuffer, amsPacket);
                } catch (ParseException e) {
                    throw new AdsException(amsPacket.getInvokeId(), e);
                }
                out.add(writeBuffer.getData());
            } catch (Exception e) {
                LOGGER.warn("Error decoding package: " + e.getMessage());
            }
        }
    }

}
