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
import org.apache.plc4x.java.amsads.attic.protocol.exception.AdsException;
import org.apache.plc4x.java.ads.readwrite.AmsPacket;
import org.apache.plc4x.java.ads.readwrite.io.AmsPacketIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

@Deprecated
public class Ads2PayloadProtocol extends MessageToMessageCodec<ByteBuf, AmsPacket> {

    private static final Logger LOGGER = LoggerFactory.getLogger(Ads2PayloadProtocol.class);

    private final ConcurrentMap<Long, AmsPacket> requests;

    public Ads2PayloadProtocol() {
        this.requests = new ConcurrentHashMap<>();
    }

    /**
     * Resets this protocol and discard all send requests.
     */
    public void reset() {
        requests.clear();
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, AmsPacket amsPacket, List<Object> out) throws AdsException {
        LOGGER.trace("(<--OUT): {}, {}, {}", channelHandlerContext, amsPacket, out);
        Long invokeId = amsPacket.getInvokeId();
        if (invokeId != 0L) {
            requests.put(invokeId, amsPacket);
        }
        WriteBuffer writeBuffer = new WriteBuffer(amsPacket.getLengthInBytes(), true);
        try {
            AmsPacketIO.staticSerialize(writeBuffer, amsPacket);
        } catch (ParseException e) {
            throw new AdsException(invokeId, e);
        }
        out.add(writeBuffer.getData());
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws AdsException {
        if (byteBuf == Unpooled.EMPTY_BUFFER) {
            // Cleanup...
            reset();
            return;
        }
        LOGGER.trace("(-->IN): {}, {}, {}", channelHandlerContext, byteBuf, out);

        byte[] bytes = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(bytes);
        ReadBuffer readBuffer = new ReadBuffer(bytes);
        while (readBuffer.getPos() < bytes.length) {
            try {
                AmsPacket packet = AmsPacketIO.staticParse(readBuffer);
                out.add(packet);
            } catch (Exception e) {
                LOGGER.warn("Error decoding package: " + e.getMessage());
            }
        }
    }
}
