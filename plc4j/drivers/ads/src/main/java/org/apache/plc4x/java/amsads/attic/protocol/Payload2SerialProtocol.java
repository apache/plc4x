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
import io.netty.util.concurrent.ScheduledFuture;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.amsads.attic.protocol.exception.AdsException;
import org.apache.plc4x.java.amsads.attic.protocol.util.DigestUtil;
import org.apache.plc4x.java.ads.readwrite.AmsPacket;
import org.apache.plc4x.java.ads.readwrite.AmsSerialFrame;
import org.apache.plc4x.java.ads.readwrite.AmsTCPPacket;
import org.apache.plc4x.java.ads.readwrite.io.AmsPacketIO;
import org.apache.plc4x.java.ads.readwrite.io.AmsSerialFrameIO;
import org.apache.plc4x.java.ads.readwrite.io.AmsTCPPacketIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

@Deprecated
public class Payload2SerialProtocol extends MessageToMessageCodec<ByteBuf, ByteBuf> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Payload2SerialProtocol.class);

    private final AtomicInteger fragmentCounter = new AtomicInteger(0);

    private AtomicReference<ScheduledFuture<?>> currentRetryer = new AtomicReference<>();

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


        AmsSerialFrame amsSerialFrame = new AmsSerialFrame(
            0x5A01,
            (byte) 0,
            (byte) 0,
            (byte) 0,
            (byte) amsPacketSer.getLengthInBytes(),
            amsPacketSer,
            DigestUtil.calculateCrc16(ArrayUtils.addAll(new byte[]{0x5A, 0x01, 0, 0, 0, (byte) amsPacketSer.getLengthInBytes()}, bytes))
        );

        WriteBuffer writeBuffer = new WriteBuffer(amsPacketSer.getLengthInBytes(), true);
        try {
            AmsSerialFrameIO.staticSerialize(writeBuffer, amsSerialFrame);
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
