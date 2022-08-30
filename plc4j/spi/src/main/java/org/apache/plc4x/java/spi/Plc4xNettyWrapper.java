/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi;

import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.MessageToMessageCodec;
import io.vavr.control.Either;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.events.*;
import org.apache.plc4x.java.spi.internal.DefaultExpectRequestContext;
import org.apache.plc4x.java.spi.internal.DefaultSendRequestContext;
import org.apache.plc4x.java.spi.internal.HandlerRegistration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.time.Instant;
import java.util.Deque;
import java.util.Iterator;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;

public class Plc4xNettyWrapper<T> extends MessageToMessageCodec<T, Object> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xNettyWrapper.class);

    private final Plc4xProtocolBase<T> protocolBase;

    private final PlcAuthentication authentication;

    private final Queue<HandlerRegistration> registeredHandlers;
    private final ChannelPipeline pipeline;
    private final boolean passive;

    public Plc4xNettyWrapper(ChannelPipeline pipeline, boolean passive, Plc4xProtocolBase<T> protocol,
                             PlcAuthentication authentication, Class<T> clazz) {
        super(clazz, Object.class);
        this.pipeline = pipeline;
        this.passive = passive;
        this.registeredHandlers = new ConcurrentLinkedQueue<>();
        this.protocolBase = protocol;
        this.authentication = authentication;
        this.protocolBase.setContext(new ConversationContext<T>() {

            @Override
            public PlcAuthentication getAuthentication() {
                return authentication;
            }

            @Override
            public Channel getChannel() {
                return pipeline.channel();
            }

            @Override
            public boolean isPassive() {
                return passive;
            }

            @Override
            public void sendToWire(T msg) {
                pipeline.writeAndFlush(msg);
            }

            @Override
            public void fireConnected() {
                pipeline.fireUserEventTriggered(ConnectedEvent.class);
            }

            @Override
            public void fireDisconnected() {
                pipeline.fireUserEventTriggered(DisconnectedEvent.class);
            }

            @Override
            public void fireDiscovered(Configuration c) {
                pipeline.fireUserEventTriggered(DiscoveredEvent.class);
            }

            @Override
            public SendRequestContext<T> sendRequest(T packet) {
                return new DefaultSendRequestContext<>(handler -> {
                    logger.trace("Adding Response Handler ...");
                    registeredHandlers.add(handler);
                }, packet, this);
            }

            @Override
            public ExpectRequestContext<T> expectRequest(Class<T> clazz, Duration timeout) {
                return new DefaultExpectRequestContext<>(handler -> {
                    logger.trace("Adding Request Handler ...");
                    registeredHandlers.add(handler);
                }, clazz, timeout, this);
            }

        });
    }

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, Object msg, List<Object> list) throws Exception {
//        logger.trace("Encoding {}", plcRequestContainer);
//        protocolBase.encode(new DefaultConversationContext<T>(channelHandlerContext) {
//            @Override
//            public void sendToWire(T msg) {
//                logger.trace("Sending to wire {}", msg);
//                list.add(msg);
//            }
//        }, plcRequestContainer);
        // NOOP
        logger.debug("Forwarding request to plc {}", msg);
        list.add(msg);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, T t, List<Object> list) throws Exception {
        logger.trace("Decoding {}", t);
        // Just iterate the list to find a suitable  Handler

        registrations:
        for (Iterator<HandlerRegistration> iter = this.registeredHandlers.iterator(); iter.hasNext(); ) {
            HandlerRegistration registration = iter.next();
            // Check if the handler can still be used or should be removed
            // Was cancelled?
            if (registration.isCancelled()) {
                logger.debug("Removing {} as it was cancelled!", registration);
                iter.remove();
                continue;
            }
            // Timeout?
            final Instant now = Instant.now();
            if (registration.getTimeoutAt().isBefore(now)) {
                logger.debug("Removing {} as its timed out (timeout of {} was set till {} and now is {})",
                    registration, registration.getTimeout(), registration.getTimeoutAt(), now);
                // pass timeout back to caller so it can do ie. transaction compensation
                registration.getOnTimeoutConsumer().accept(new TimeoutException());
                iter.remove();
                continue;
            }
            logger.trace("Checking handler {} for Object of type {}", registration, t.getClass().getSimpleName());
            if (registration.getExpectClazz().isInstance(t)) {
                logger.trace("Handler {} has right expected type {}, checking condition", registration, registration.getExpectClazz().getSimpleName());
                // Check all Commands / Functions
                Deque<Either<Function<?, ?>, Predicate<?>>> commands = registration.getCommands();
                Object instance = t;
                for (Iterator<Either<Function<?, ?>, Predicate<?>>> iterator = commands.iterator(); iterator.hasNext(); ) {
                    Either<Function<?, ?>, Predicate<?>> either = iterator.next();
                    if (either.isLeft()) {
                        Function unwrap = either.getLeft();
                        instance = unwrap.apply(instance);
                    } else {
                        Predicate predicate = either.get();
                        if (!predicate.test(instance)) {
                            // We do not match -> cannot handle
                            logger.trace("Registration {} with predicate {} does not match object {} (currently wrapped to {})", registration, predicate,
                                t.getClass().getSimpleName(), instance.getClass().getSimpleName());
                            continue registrations;
                        }
                    }
                }
                logger.trace("Handler {} accepts element {}, calling handle method", registration, t);
                this.registeredHandlers.remove(registration);
                Consumer handler = registration.getPacketConsumer();
                handler.accept(instance);
                // Confirm that it was handled!
                registration.confirmHandled();
                return;
            }
        }
        logger.trace("None of {} registered handlers could handle message {}, using default decode method", this.registeredHandlers.size(), t);
        protocolBase.decode(new DefaultConversationContext<>(channelHandlerContext, authentication, passive), t);
    }

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        logger.debug("User Event triggered {}", evt);
        if (evt instanceof ConnectEvent) {
            this.protocolBase.onConnect(new DefaultConversationContext<>(ctx, authentication, passive));
        } else if (evt instanceof DisconnectEvent) {
            this.protocolBase.onDisconnect(new DefaultConversationContext<>(ctx, authentication, passive));
        } else if (evt instanceof DiscoverEvent) {
            this.protocolBase.onDiscover(new DefaultConversationContext<>(ctx, authentication, passive));
        } else if (evt instanceof CloseConnectionEvent) {
            this.protocolBase.close(new DefaultConversationContext<>(ctx, authentication, passive));
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    public class DefaultConversationContext<T1> implements ConversationContext<T1> {
        private final ChannelHandlerContext channelHandlerContext;

        private final PlcAuthentication authentication;
        private final boolean passive;

        public DefaultConversationContext(ChannelHandlerContext channelHandlerContext,
                                          PlcAuthentication authentication,
                                          boolean passive) {
            this.channelHandlerContext = channelHandlerContext;
            this.authentication = authentication;
            this.passive = passive;
        }

        @Override
        public Channel getChannel() {
            return channelHandlerContext.channel();
        }

        public PlcAuthentication getAuthentication() {
            return authentication;
        }

        @Override
        public boolean isPassive() {
            return passive;
        }

        @Override
        public void sendToWire(T1 msg) {
            logger.trace("Sending to wire {}", msg);
            channelHandlerContext.channel().writeAndFlush(msg);
        }

        @Override
        public void fireConnected() {
            logger.trace("Firing Connected!");
            channelHandlerContext.pipeline().fireUserEventTriggered(new ConnectedEvent());
        }

        @Override
        public void fireDisconnected() {
            logger.trace("Firing Disconnected!");
            channelHandlerContext.pipeline().fireUserEventTriggered(new DisconnectedEvent());
        }

        @Override
        public void fireDiscovered(Configuration c) {
            logger.trace("Firing Discovered!");
            channelHandlerContext.pipeline().fireUserEventTriggered(new DiscoveredEvent(c));
        }

        @Override
        public SendRequestContext<T1> sendRequest(T1 packet) {
            return new DefaultSendRequestContext<>(handler -> {
                logger.trace("Adding Response Handler ...");
                registeredHandlers.add(handler);
            }, packet, this);
        }

        @Override
        public ExpectRequestContext<T1> expectRequest(Class<T1> clazz, Duration timeout) {
            return new DefaultExpectRequestContext<>(handler -> {
                logger.trace("Adding Request Handler ...");
                registeredHandlers.add(handler);
            }, clazz, timeout, this);
        }
    }

}
