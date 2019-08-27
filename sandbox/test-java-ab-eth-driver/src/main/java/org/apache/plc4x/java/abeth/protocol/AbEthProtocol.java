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
package org.apache.plc4x.java.abeth.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.abeth.CIPEncapsulationConnectionRequest;
import org.apache.plc4x.java.abeth.CIPEncapsulationConnectionResponse;
import org.apache.plc4x.java.abeth.CIPEncapsulationPacket;
import org.apache.plc4x.java.abeth.io.CIPEncapsulationPacketIO;
import org.apache.plc4x.java.base.PlcByteToMessageCodec;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class AbEthProtocol extends PlcByteToMessageCodec<CIPEncapsulationPacket> {

    private static final Logger logger = LoggerFactory.getLogger(AbEthProtocol.class);

    private CIPEncapsulationPacketIO io;

    public AbEthProtocol() {
        io = new CIPEncapsulationPacketIO();
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, CIPEncapsulationPacket cipEncapsulationPacket, ByteBuf byteBuf) throws Exception {
        WriteBuffer buffer = new WriteBuffer(cipEncapsulationPacket.getLengthInBytes());
        io.serialize(buffer, cipEncapsulationPacket);
        byteBuf.writeBytes(buffer.getData());
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf byteBuf, List<Object> out) throws Exception {
        byte[] bytes = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(bytes);
        ReadBuffer readBuffer = new ReadBuffer(bytes);
        while (readBuffer.getPos() < bytes.length) {
            try {
                CIPEncapsulationPacket packet = io.parse(readBuffer);
                out.add(packet);
            } catch (Exception e) {
                logger.warn("Error decoding package: " + e.getMessage());
            }
        }
    }

}
