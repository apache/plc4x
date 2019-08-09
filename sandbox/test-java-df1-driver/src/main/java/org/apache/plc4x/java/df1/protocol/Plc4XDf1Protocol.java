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
package org.apache.plc4x.java.df1.protocol;

import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.df1.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4XDf1Protocol extends PlcMessageToMessageCodec<DF1Command, PlcRequestContainer> {

    private static final Logger logger = LoggerFactory.getLogger(Plc4XDf1Protocol.class);

    private final AtomicInteger transactionId = new AtomicInteger(0);

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
        if (msg.getRequest() instanceof PlcReadRequest) {
            for (PlcField field : ((PlcReadRequest) msg.getRequest()).getFields()) {
                if (!(field instanceof Df1Field)) {
                    throw new IllegalArgumentException("Invalid field type found inside Df1 Request");
                }
                int address = ((Df1Field) field).getAddress();
                short size = (short) ((Df1Field) field).getAddress();
                int transactionId = this.transactionId.getAndIncrement();
                logger.debug("Creating request for offset {}, with length {} and transaction id {}", address, size, transactionId);
                // TODO: differentiate commands
                out.add(new DF1UnprotectedReadRequest((short) 0x00, transactionId, address, size));
            }
        } else {
            throw new IllegalStateException("This should not happen!");
        }
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, DF1Command msg, List<Object> out) throws Exception {
        System.out.println("Hello");
    }

}
