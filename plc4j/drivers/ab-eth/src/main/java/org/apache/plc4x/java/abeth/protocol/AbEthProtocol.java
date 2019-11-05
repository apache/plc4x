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
import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationPacket;
import org.apache.plc4x.java.abeth.readwrite.io.CIPEncapsulationPacketIO;
import org.apache.plc4x.java.base.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.utils.MessageIO;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class AbEthProtocol extends GeneratedDriverByteToMessageCodec<CIPEncapsulationPacket> {

    private static final Logger logger = LoggerFactory.getLogger(AbEthProtocol.class);

    public AbEthProtocol() {
        super(new MessageIO<CIPEncapsulationPacket, CIPEncapsulationPacket>() {
            @Override
            public CIPEncapsulationPacket parse(ReadBuffer io) throws ParseException {
                return CIPEncapsulationPacketIO.parse(io);
            }

            @Override
            public void serialize(WriteBuffer io, CIPEncapsulationPacket value) throws ParseException {
                CIPEncapsulationPacketIO.serialize(io, value);
            }
        });
        logger.trace("Created new AB-ETH protocol");
    }

    @Override
    protected int getPacketSize(ByteBuf byteBuf) {
        return byteBuf.readableBytes();
    }

    @Override
    protected void removeRestOfCorruptPackage(ByteBuf byteBuf) {
        // Nothing to do here ...
    }

}
