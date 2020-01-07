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
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.Message;

import java.util.function.Consumer;
import java.util.function.Function;

/**
 * Used to Build Instances of {@link SingleProtocolStackConfigurer}.
 *
 * @param <BASE_PACKET_CLASS> Type of Created Message that is Exchanged.
 */
public final class SingleProtocolStackBuilder<BASE_PACKET_CLASS extends Message> {

    private Class<BASE_PACKET_CLASS> basePacketClass;
    private Class<? extends Plc4xProtocolBase<BASE_PACKET_CLASS>> protocol;
    private Class<? extends Function<ByteBuf, Integer>> packetSizeEstimator;
    private Class<? extends Consumer<ByteBuf>> corruptPacketRemover;

    SingleProtocolStackBuilder(Class<BASE_PACKET_CLASS> basePacketClass) {
        this.basePacketClass = basePacketClass;
    }

    public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withProtocol(Class<? extends Plc4xProtocolBase<BASE_PACKET_CLASS>> protocol) {
        this.protocol = protocol;
        return this;
    }

    public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withPacketSizeEstimator(Class<? extends Function<ByteBuf, Integer>> packetSizeEstimator) {
        this.packetSizeEstimator = packetSizeEstimator;
        return this;
    }

    public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withCorruptPacketRemover(Class<? extends Consumer<ByteBuf>> corruptPacketRemover) {
        this.corruptPacketRemover = corruptPacketRemover;
        return this;
    }

    public SingleProtocolStackConfigurer<BASE_PACKET_CLASS> build() {
        assert this.protocol != null;
        return new SingleProtocolStackConfigurer<>(basePacketClass, protocol, packetSizeEstimator, corruptPacketRemover);
    }

}
