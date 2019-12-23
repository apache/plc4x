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
package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.io.TPKTPacketIO;
import org.apache.plc4x.java.spi.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class S7ProtocolMessage extends GeneratedDriverByteToMessageCodec<TPKTPacket> {

    private static final Logger logger = LoggerFactory.getLogger(S7ProtocolMessage.class);

    public S7ProtocolMessage() {
        super(new MessageIO<TPKTPacket, TPKTPacket>() {
            @Override
            public TPKTPacket parse(ReadBuffer io) throws ParseException {
                return TPKTPacketIO.parse(io);
            }

            @Override
            public void serialize(WriteBuffer io, TPKTPacket value) throws ParseException {
                TPKTPacketIO.serialize(io, value);
            }
        });
    }

    @Override
    protected int getPacketSize(ByteBuf byteBuf) {
        if(byteBuf.readableBytes() >= 4) {
            return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 2);
        }
        return -1;
    }

    @Override
    protected void removeRestOfCorruptPackage(ByteBuf byteBuf) {
        while (byteBuf.getUnsignedByte(0) != TPKTPacket.PROTOCOLID) {
            // Just consume the bytes till the next possible start position.
            byteBuf.readByte();
        }
    }

}
