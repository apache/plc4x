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

package org.apache.plc4x.java.s7.readwrite.protocol;

import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPipeline;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.s7.readwrite.events.IsoTPConnectedEvent;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Field;
import java.util.List;
import java.util.function.Consumer;

public class Plc4xProtocolContext<T> extends MessageToMessageCodec<T, PlcRequestContainer> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4xProtocolContext.class);

    private final Plc4xProtocolBase<T> parent;
    private volatile ChannelHandler prevChannelHandler = null;
    private ChannelPipeline pipeline;

    // VERY DIRTY
    private Consumer<Object> handler;

    public Plc4xProtocolContext(Plc4xProtocolBase<T> parent, Class<T> clazz) {
        super(clazz, PlcRequestContainer.class);
        this.parent = parent;
    }


    @Override protected void encode(ChannelHandlerContext channelHandlerContext, PlcRequestContainer plcRequestContainer, List<Object> list) throws Exception {
        logger.info("Encoding {}", plcRequestContainer);
        parent.encode(plcRequestContainer, o -> list.add(o));
    }

    @Override protected void decode(ChannelHandlerContext channelHandlerContext, T t, List<Object> list) throws Exception {
        logger.info("Decoding {}", t);
        parent.decode(t);
    }

    public void send(T packet) {
        logger.info("Sending {}", packet);
        pipeline.writeAndFlush(packet);
    }

    public void fireConnected() {
        logger.info("Firing Connected!");
        pipeline.fireUserEventTriggered(new ConnectedEvent());
    }

    /**
     * If the S7 protocol layer is used over Iso TP, then after receiving a {@link IsoTPConnectedEvent} the
     * corresponding S7 setup communication message has to be sent in order to negotiate the S7 protocol layer.
     *
     * @param ctx the current protocol layers context
     * @param evt the event
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // If the connection has just been established, start setting up the connection
        // by sending a connection request to the plc.
        if (evt instanceof ConnectEvent) {
            this.parent.onConnect();
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    protected ChannelHandler getPrevChannelHandler(ChannelHandlerContext ctx) {
        if (prevChannelHandler == null) {
            try {
                Field prevField = FieldUtils.getField(ctx.getClass(), "prev", true);
                if (prevField != null) {
                    ChannelHandlerContext prevContext = (ChannelHandlerContext) prevField.get(ctx);
                    prevChannelHandler = prevContext.handler();
                }
            } catch (Exception e) {
                logger.error("Error accessing field 'prev'", e);
            }
        }
        return prevChannelHandler;
    }

    public void setPipeline(ChannelPipeline pipeline) {
        this.pipeline = pipeline;
    }

    public ChannelPipeline getPipeline() {
        return pipeline;
    }
}
