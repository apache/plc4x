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
package org.apache.plc4x.java.bacnetip.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.bacnetip.readwrite.io.BVLCIO;
import org.apache.plc4x.java.base.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.utils.MessageIO;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;

public class BacNetIpProtocol extends GeneratedDriverByteToMessageCodec<BVLC> {

    public BacNetIpProtocol() {
        super(new MessageIO<BVLC, BVLC>() {
            @Override
            public BVLC parse(ReadBuffer io) throws ParseException {
                return BVLCIO.parse(io);
            }

            @Override
            public void serialize(WriteBuffer io, BVLC value) throws ParseException {
                BVLCIO.serialize(io, value);
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
        while (byteBuf.getUnsignedByte(0) != BVLC.BACNETTYPE) {
            // Just consume the bytes till the next possible start position.
            byteBuf.readByte();
        }
    }

}
