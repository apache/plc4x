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
package org.apache.plc4x.java.base;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.utils.Message;
import org.apache.plc4x.java.utils.MessageIO;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public abstract class GeneratedDriverByteToMessageCodec<T extends Message> extends PlcByteToMessageCodec<T> {

    private static final Logger logger = LoggerFactory.getLogger(GeneratedDriverByteToMessageCodec.class);

    private MessageIO<T, T> io;

    public GeneratedDriverByteToMessageCodec(MessageIO<T, T> io) {
        this.io = io;
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, T packet, ByteBuf byteBuf) throws Exception {
        WriteBuffer buffer = new WriteBuffer(packet.getLengthInBytes());
        io.serialize(buffer, packet);
        byteBuf.writeBytes(buffer.getData());
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf byteBuf, List<Object> out) throws Exception {
        // Check if enough data is present to process the entire package.
        int packetSize = getPacketSize(byteBuf);
        if(packetSize == -1 || packetSize > byteBuf.readableBytes()) {
            return;
        }

        byte[] bytes = new byte[packetSize];
        byteBuf.readBytes(bytes);

        ReadBuffer readBuffer = new ReadBuffer(bytes);
        while (readBuffer.getPos() < bytes.length) {
            try {
                T packet = io.parse(readBuffer);
                out.add(packet);
            } catch (Exception e) {
                logger.warn("Error decoding package with content [" + Hex.encodeHexString(bytes) + "]: "
                    + e.getMessage(), e);
                // Just remove any trailing junk ... if there is any.
                removeRestOfCorruptPackage(byteBuf);
            }
        }
    }

    abstract protected int getPacketSize(ByteBuf byteBuf);

    abstract protected void removeRestOfCorruptPackage(ByteBuf byteBuf);

}
