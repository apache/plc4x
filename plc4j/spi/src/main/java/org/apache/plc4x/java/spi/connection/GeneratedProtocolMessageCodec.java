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
import org.apache.plc4x.java.spi.GeneratedDriverByteToMessageCodec;
import org.apache.plc4x.java.spi.generation.*;

import java.util.function.Consumer;
import java.util.function.ToIntFunction;

public class GeneratedProtocolMessageCodec<BASE_PACKET_CLASS extends Message> extends GeneratedDriverByteToMessageCodec<BASE_PACKET_CLASS> {

    private final ToIntFunction<ByteBuf> packetSizeEstimator;
    private final Consumer<ByteBuf> corruptPackageRemover;

    public GeneratedProtocolMessageCodec(
        Class<BASE_PACKET_CLASS> basePacketClass,
        MessageInput<BASE_PACKET_CLASS> messageIO,
        ByteOrder byteOrder,
        Object[] parserArgs,
        ToIntFunction<ByteBuf> packetSizeEstimator,
        Consumer<ByteBuf> corruptPackageRemover) {
        super(messageIO, basePacketClass, byteOrder, parserArgs);
        this.packetSizeEstimator = packetSizeEstimator;
        this.corruptPackageRemover = corruptPackageRemover;
    }

    @Override
    protected int getPacketSize(ByteBuf byteBuf) {
        if (this.packetSizeEstimator == null) {
            return byteBuf.readableBytes();
        }
        return packetSizeEstimator.applyAsInt(byteBuf);
    }

    @Override
    protected void removeRestOfCorruptPackage(ByteBuf byteBuf) {
        if (this.corruptPackageRemover != null) {
            this.corruptPackageRemover.accept(byteBuf);
        }
    }

}
