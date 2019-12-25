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
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelPipeline;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.Message;

import java.util.function.Consumer;
import java.util.function.Function;

/**
 * Builds a Protocol Stack.
 */
public class SingleProtocolStackConfigurer<BASE_PAKET_CLASS extends Message> implements ProtocolStackConfigurer<BASE_PAKET_CLASS> {

    private final Class<BASE_PAKET_CLASS> basePaketClass;
    private final Plc4xProtocolBase<BASE_PAKET_CLASS> protocol;
    private final Function<ByteBuf, Integer> packetSizeEstimator;
    private final Consumer<ByteBuf> corruptPacketRemover;

    /** Only accessible via Builder */
    SingleProtocolStackConfigurer(Class<BASE_PAKET_CLASS> basePaketClass, Plc4xProtocolBase<BASE_PAKET_CLASS> protocol,
                                  Function<ByteBuf, Integer> packetSizeEstimator,
                                  Consumer<ByteBuf> corruptPacketRemover) {
        this.basePaketClass = basePaketClass;
        this.protocol = protocol;
        this.packetSizeEstimator = packetSizeEstimator;
        this.corruptPacketRemover = corruptPacketRemover;
    }

    public static <BPC extends Message> SingleProtocolStackBuilder<BPC> builder(Class<BPC> basePaketClass) {
        return new SingleProtocolStackBuilder<>(basePaketClass);
    }

    private ChannelHandler getMessageCodec() {
        ReflectionBasedIo<BASE_PAKET_CLASS> io = new ReflectionBasedIo<>(basePaketClass);
        return new GeneratedProtocolMessageCodec<>(basePaketClass, io, io, packetSizeEstimator, corruptPacketRemover);
    }

    /** Applies the given Stack to the Pipeline */
    @Override public Plc4xProtocolBase<BASE_PAKET_CLASS> apply(ChannelPipeline pipeline) {
        pipeline.addLast(getMessageCodec());
        Plc4xNettyWrapper<BASE_PAKET_CLASS> context = new Plc4xNettyWrapper<>(pipeline, protocol, basePaketClass);
        pipeline.addLast(context);
        return protocol;
    }

}
