package org.apache.plc4x.java.base;/*
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

import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.commons.lang3.reflect.FieldUtils;

import java.lang.reflect.Field;

public abstract class PlcMessageToMessageCodec<INBOUND_IN, OUTBOUND_IN>
    extends MessageToMessageCodec<INBOUND_IN, OUTBOUND_IN> {

    private volatile ChannelHandler prevChannelHandler = null;

    public PlcMessageToMessageCodec() {
    }

    public PlcMessageToMessageCodec(Class<? extends INBOUND_IN> inboundMessageType,
                                    Class<? extends OUTBOUND_IN> outboundMessageType) {
        super(inboundMessageType, outboundMessageType);
    }



    protected ChannelHandler getPrevChannelHandler(ChannelHandlerContext ctx) {
        if(prevChannelHandler == null) {
            try {
                Field prevField = FieldUtils.getField(ctx.getClass(), "prev", true);
                if(prevField != null) {
                    ChannelHandlerContext prevContext = (ChannelHandlerContext) prevField.get(ctx);
                    prevChannelHandler = prevContext.handler();
                }
            } catch(Exception e) {
                e.printStackTrace();
            }
        }
        return prevChannelHandler;
    }

}
