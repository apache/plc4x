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

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.*;

import io.netty.buffer.ByteBuf;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelPipeline;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.MessageIO;

import java.lang.reflect.InvocationTargetException;
import java.util.function.Consumer;
import java.util.function.ToIntFunction;

/**
 * Builds a Protocol Stack.
 */
public class SingleProtocolStackConfigurer<BASE_PACKET_CLASS extends Message> implements ProtocolStackConfigurer<BASE_PACKET_CLASS> {

    private final Class<BASE_PACKET_CLASS> basePacketClass;
    private boolean bigEndian = true;
    private final Class<? extends Plc4xProtocolBase<BASE_PACKET_CLASS>> protocolClass;
    private final Class<? extends DriverContext> driverContextClass;
    private final MessageIO<BASE_PACKET_CLASS, BASE_PACKET_CLASS> protocolIO;
    private final Class<? extends ToIntFunction<ByteBuf>> packetSizeEstimatorClass;
    private final Class<? extends Consumer<ByteBuf>> corruptPacketRemoverClass;
    private final Object[] parserArgs;

    public static <BPC extends Message> SingleProtocolStackBuilder<BPC> builder(Class<BPC> basePacketClass, Class<? extends MessageIO<BPC, BPC>> messageIoClass) {
        return new SingleProtocolStackBuilder<>(basePacketClass, messageIoClass);
    }

    /** Only accessible via Builder */
    SingleProtocolStackConfigurer(Class<BASE_PACKET_CLASS> basePacketClass,
                                  boolean bigEndian,
                                  Object[] parserArgs,
                                  Class<? extends Plc4xProtocolBase<BASE_PACKET_CLASS>> protocol,
                                  Class<? extends DriverContext> driverContextClass,
                                  MessageIO<BASE_PACKET_CLASS, BASE_PACKET_CLASS> protocolIO,
                                  Class<? extends ToIntFunction<ByteBuf>> packetSizeEstimatorClass,
                                  Class<? extends Consumer<ByteBuf>> corruptPacketRemoverClass) {
        this.basePacketClass = basePacketClass;
        this.bigEndian = bigEndian;
        this.parserArgs = parserArgs;
        this.protocolClass = protocol;
        this.driverContextClass = driverContextClass;
        this.protocolIO = protocolIO;
        this.packetSizeEstimatorClass = packetSizeEstimatorClass;
        this.corruptPacketRemoverClass = corruptPacketRemoverClass;
    }

    private ChannelHandler getMessageCodec(Configuration configuration) {
        return new GeneratedProtocolMessageCodec<>(basePacketClass, protocolIO, bigEndian, parserArgs,
            packetSizeEstimatorClass != null ? configure(configuration, createInstance(packetSizeEstimatorClass)) : null,
            corruptPacketRemoverClass != null ? configure(configuration, createInstance(corruptPacketRemoverClass)) : null);
    }

    /** Applies the given Stack to the Pipeline */
    @Override
    public Plc4xProtocolBase<BASE_PACKET_CLASS> configurePipeline(
            Configuration configuration, ChannelPipeline pipeline, boolean passive) {
        pipeline.addLast(getMessageCodec(configuration));
        Plc4xProtocolBase<BASE_PACKET_CLASS> protocol = configure(configuration, createInstance(protocolClass));
        if(driverContextClass != null) {
            protocol.setDriverContext(configure(configuration, createInstance(driverContextClass)));
        }
        Plc4xNettyWrapper<BASE_PACKET_CLASS> context = new Plc4xNettyWrapper<>(pipeline, passive, protocol, basePacketClass);
        pipeline.addLast(context);
        return protocol;
    }

    private <T> T createInstance(Class<T> clazz, Object... args) {
        try {
            Class<?>[] parameterTypes = new Class<?>[args.length];
            for(int i = 0; i < args.length; i++) {
                parameterTypes[i] = args[i].getClass();
            }
            return clazz.getDeclaredConstructor(parameterTypes).newInstance(args);
        } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException  e) {
            throw new PlcRuntimeException("Error creating instance of class " + clazz.getName());
        }
    }

    /**
     * Used to Build Instances of {@link SingleProtocolStackConfigurer}.
     *
     * @param <BASE_PACKET_CLASS> Type of Created Message that is Exchanged.
     */
    public static final class SingleProtocolStackBuilder<BASE_PACKET_CLASS extends Message> {

        private final Class<BASE_PACKET_CLASS> basePacketClass;
        private final Class<? extends MessageIO<BASE_PACKET_CLASS, BASE_PACKET_CLASS>> messageIoClass;
        private Class<? extends DriverContext> driverContextClass;
        private boolean bigEndian = true;
        private Object[] parserArgs;
        private Class<? extends Plc4xProtocolBase<BASE_PACKET_CLASS>> protocol;
        private Class<? extends ToIntFunction<ByteBuf>> packetSizeEstimator;
        private Class<? extends Consumer<ByteBuf>> corruptPacketRemover;

        public SingleProtocolStackBuilder(Class<BASE_PACKET_CLASS> basePacketClass, Class<? extends MessageIO<BASE_PACKET_CLASS, BASE_PACKET_CLASS>> messageIoClass) {
            this.basePacketClass = basePacketClass;
            this.messageIoClass = messageIoClass;
        }

        public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withDriverContext(Class<? extends DriverContext> driverContextClass) {
            this.driverContextClass = driverContextClass;
            return this;
        }

        public SingleProtocolStackBuilder<BASE_PACKET_CLASS> littleEndian() {
            this.bigEndian = false;
            return this;
        }

        public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withParserArgs(Object... parserArgs) {
            this.parserArgs = parserArgs;
            return this;
        }

        public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withProtocol(Class<? extends Plc4xProtocolBase<BASE_PACKET_CLASS>> protocol) {
            this.protocol = protocol;
            return this;
        }

        public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withPacketSizeEstimator(Class<? extends ToIntFunction<ByteBuf>> packetSizeEstimator) {
            this.packetSizeEstimator = packetSizeEstimator;
            return this;
        }

        public SingleProtocolStackBuilder<BASE_PACKET_CLASS> withCorruptPacketRemover(Class<? extends Consumer<ByteBuf>> corruptPacketRemover) {
            this.corruptPacketRemover = corruptPacketRemover;
            return this;
        }

        public SingleProtocolStackConfigurer<BASE_PACKET_CLASS> build() {
            assert this.protocol != null;
            try {
                final MessageIO messageIo = messageIoClass.getDeclaredConstructor().newInstance();
                return new SingleProtocolStackConfigurer<>(
                    basePacketClass, bigEndian, parserArgs, protocol, driverContextClass, messageIo, packetSizeEstimator, corruptPacketRemover);
            } catch (InstantiationException | InvocationTargetException | NoSuchMethodException | IllegalAccessException e) {
                throw new PlcRuntimeException("Error initializing MessageIO instance", e);
            }
        }

    }

}
