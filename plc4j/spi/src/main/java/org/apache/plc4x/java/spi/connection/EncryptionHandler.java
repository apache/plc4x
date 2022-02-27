/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.connection;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class EncryptionHandler extends MessageToMessageCodec<ByteBuf, ByteBuf> {

    private static final Logger logger = LoggerFactory.getLogger(DefaultNettyPlcConnection.class);

    public EncryptionHandler() {
        super(ByteBuf.class, ByteBuf.class);
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) {
        logger.debug("Encrypting outgoing message");
        in.retain();
        encrypt(ctx, in, out);
        out.add(in);
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) {
        logger.debug("Received Incoming message and decrypting");
        in.retain();
        decrypt(ctx, in, out);
        out.add(in);
    }

    /**
     * Overridable function used to encrypt an outgoing.
     *
     * @return ByteBuf the encrypted buffer should be returned.
     */
    protected ByteBuf encrypt(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) {
        return in;
    }

    /**
     * Overridable function used to decrypt the incoming message.
     *
     * @return ByteBuf the decrypted buffer should be returned.
     */
    protected ByteBuf decrypt(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) {
        return in;
    }

}
