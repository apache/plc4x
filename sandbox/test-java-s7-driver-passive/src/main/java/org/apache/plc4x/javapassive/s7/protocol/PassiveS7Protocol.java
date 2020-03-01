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
package org.apache.plc4x.javapassive.s7.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.base.PlcByteToMessageCodec;
import org.apache.plc4x.java.s7.passive.TPKTPacket;
import org.apache.plc4x.java.s7.passive.io.TPKTPacketIO;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class PassiveS7Protocol extends PlcByteToMessageCodec<TPKTPacket> {

    private static final Logger logger = LoggerFactory.getLogger(PassiveS7Protocol.class);

    private TPKTPacketIO io;

    public PassiveS7Protocol() {
        io = new TPKTPacketIO();
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, TPKTPacket tpktPacket, ByteBuf byteBuf) throws Exception {
        System.out.println(tpktPacket);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> out) throws Exception {
        byte[] bytes = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(bytes);
        ReadBuffer readBuffer = new ReadBuffer(bytes);
        while(readBuffer.getPos() < bytes.length) {
            try {
                TPKTPacket packet = io.parse(readBuffer);
                out.add(packet);
            } catch(Exception e) {
                logger.warn("Error decoding package: " + e.getMessage());
            }
        }
    }

}
