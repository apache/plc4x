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
package org.apache.plc4x.java.knxnetip.protocol;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.spi.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.knxnetip.readwrite.KNXNetIPMessage;
import org.apache.plc4x.java.knxnetip.readwrite.io.KNXNetIPMessageIO;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class KnxNetIpProtocolMessage extends GeneratedDriverByteToMessageCodec<KNXNetIPMessage> {

    public KnxNetIpProtocolMessage() {
        super(new MessageIO<KNXNetIPMessage, KNXNetIPMessage>() {
            @Override
            public KNXNetIPMessage parse(ReadBuffer io) throws ParseException {
                return KNXNetIPMessageIO.parse(io);
            }

            @Override
            public void serialize(WriteBuffer io, KNXNetIPMessage value) throws ParseException {
                KNXNetIPMessageIO.serialize(io, value);
            }
        });
    }

    @Override
    protected int getPacketSize(ByteBuf byteBuf) {
        if(byteBuf.readableBytes() >= 6) {
            return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 4);
        }
        return -1;
    }

    @Override
    protected void removeRestOfCorruptPackage(ByteBuf byteBuf) {
        // Not really any marker byte to look for.
    }

}
