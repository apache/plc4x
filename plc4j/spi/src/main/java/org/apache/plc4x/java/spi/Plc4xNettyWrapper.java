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

package org.apache.plc4x.java.spi;

import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import io.vavr.control.Either;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.internal.DefaultSendRequestContext;
import org.apache.plc4x.java.spi.internal.HandlerRegistration;
import org.apache.plc4x.java.spi.messages.PlcRequestContainer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Deque;
import java.util.Iterator;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class Plc4xNettyWrapper<T> extends MessageToMessageCodec<T, PlcRequestContainer> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xNettyWrapper.class);

    private final Plc4xProtocolBase<T> protocolBase;
    private final Queue<HandlerRegistration> registeredHandlers;

    public Plc4xNettyWrapper(Plc4xProtocolBase<T> parent, Class<T> clazz) {
        super(clazz, PlcRequestContainer.class);
        this.registeredHandlers = new ConcurrentLinkedQueue<>();
        this.protocolBase = parent;
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, PlcRequestContainer plcRequestContainer, List<Object> list) throws Exception {
        logger.info("Encoding {}", plcRequestContainer);
        protocolBase.encode(new DefaultConversationContext<T>(channelHandlerContext) {
            @Override
            public void sendToWire(T msg) {
                logger.info("Sending to wire {}", msg);
                list.add(msg);
            }
        }, plcRequestContainer);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, T t, List<Object> list) throws Exception {
        logger.info("Decoding {}", t);
        // Just iterate the list to find a suitable  Handler
        registrations:
        for (HandlerRegistration registration : this.registeredHandlers) {
            logger.info("Checking handler {} for Object of type {}", registration, t.getClass().getSimpleName());
            if (registration.getExpectClazz().isInstance(t)) {
                logger.info("Handler {} has right expected type {}, checking condition", registration, registration.getExpectClazz().getSimpleName());
                // Check all Commands / Functions
                Deque<Either<Function<?, ?>, Predicate<?>>> commands = registration.getCommands();
                Object instance = t;
                for (Iterator<Either<Function<?, ?>, Predicate<?>>> iterator = commands.iterator(); iterator.hasNext();) {
                    Either<Function<?, ?>, Predicate<?>> either = iterator.next();
                    if (either.isLeft()) {
                        Function unwrap = either.getLeft();
                        instance = unwrap.apply(instance);
                    } else {
                        Predicate predicate = either.get();
                        if (predicate.test(instance) == false) {
                            // We do not match -> cannot handle
                            logger.info("Registration {} does not match object {} (currently wrapped to {})", registration, t.getClass().getSimpleName(), instance.getClass().getSimpleName());
                            continue registrations;
                        }
                    }
                }
                logger.info("Handler {} accepts element {}, calling handle method", registration, t);
                this.registeredHandlers.remove(registration);
                Consumer handler = registration.getPacketConsumer();
                handler.accept(instance);
                return;
            }
        }
        logger.info("No registered handler found for message {}, using default decode method", t);
        protocolBase.decode(new DefaultConversationContext<>(channelHandlerContext), t);
    }

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        if (evt instanceof ConnectEvent) {
            this.protocolBase.onConnect(new DefaultConversationContext<>(ctx));
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    public class DefaultConversationContext<T1> implements ConversationContext<T1> {
        private final ChannelHandlerContext channelHandlerContext;

        public DefaultConversationContext(ChannelHandlerContext channelHandlerContext) {
            this.channelHandlerContext = channelHandlerContext;
        }

        @Override
        public void sendToWire(T1 msg) {
            logger.info("Sending to wire {}", msg);
            channelHandlerContext.channel().writeAndFlush(msg);
        }

        @Override
        public void fireConnected() {
            logger.info("Firing Connected!");
            channelHandlerContext.pipeline().fireUserEventTriggered(new ConnectedEvent());
        }

        @Override
        public SendRequestContext<T1> sendRequest(T1 packet) {
            return new DefaultSendRequestContext<T1>(handler -> {
                logger.info("Adding Response Handler...");
                registeredHandlers.add(handler);
            }, packet, (DefaultConversationContext)this);
        }
    }
}
