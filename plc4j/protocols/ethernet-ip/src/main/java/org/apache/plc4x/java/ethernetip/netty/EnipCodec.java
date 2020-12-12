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
package org.apache.plc4x.java.ethernetip.netty;

import com.digitalpetri.enip.EnipPacket;
import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.ByteToMessageCodec;

import java.nio.ByteOrder;
import java.util.List;

public class EnipCodec extends ByteToMessageCodec<EnipPacket> {

    private static final int HEADER_SIZE = 24;
    private static final int LENGTH_OFFSET = 2;

    @Override
    protected void encode(ChannelHandlerContext ctx, EnipPacket packet, ByteBuf out) {
        EnipPacket.encode(packet, out.order(ByteOrder.LITTLE_ENDIAN));
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) {
        ByteBuf buffer = in.order(ByteOrder.LITTLE_ENDIAN);

        int startIndex = buffer.readerIndex();

        while (buffer.readableBytes() >= HEADER_SIZE &&
            buffer.readableBytes() >= HEADER_SIZE + getLength(buffer, startIndex)) {

            out.add(EnipPacket.decode(buffer));

            startIndex = buffer.readerIndex();
        }
    }

    private int getLength(ByteBuf buffer, int startIndex) {
        return buffer.getUnsignedShort(startIndex + LENGTH_OFFSET);
    }

}