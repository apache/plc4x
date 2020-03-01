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
package org.apache.plc4x.simulator.server.s7.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.base.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.io.TPKTPacketIO;
import org.apache.plc4x.java.utils.MessageIO;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;

public class S7Step7Protocol extends GeneratedDriverByteToMessageCodec<TPKTPacket> {

    public S7Step7Protocol() {
        super(new MessageIO<TPKTPacket, TPKTPacket>() {
            @Override
            public TPKTPacket parse(ReadBuffer io) throws ParseException {
                try {
                    return TPKTPacketIO.parse(io);
                } catch (Exception e) {
                    e.printStackTrace();
                    throw new ParseException("Error parsing message", e);
                }
            }

            @Override
            public void serialize(WriteBuffer io, TPKTPacket value) throws ParseException {
                try {
                    TPKTPacketIO.serialize(io, value);
                } catch (Exception e) {
                    e.printStackTrace();
                    throw new ParseException("Error serializing message", e);
                }
            }
        });
    }

    @Override
    protected int getPacketSize(ByteBuf byteBuf) {
        if(byteBuf.readableBytes() >= 4) {
            if (byteBuf.getByte(0) != TPKTPacket.PROTOCOLID) {
                return -1;
            }
            // Byte 1 is a reserved byte set to 0x00
            return byteBuf.getShort(2);
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
