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
package org.apache.plc4x.java.examples.dummydriver.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufUtil;
import io.netty.channel.ChannelHandlerContext;
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcRequestContainer;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class DummyProtocol extends MessageToMessageCodec<ByteBuf, PlcRequestContainer> {

    private static final Logger logger = LoggerFactory.getLogger(DummyProtocol.class);

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer in, List<Object> out) throws Exception {
        PlcRequest request = in.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(in, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(in, out);
        }
    }

    private void encodeWriteRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
    }

    private void encodeReadRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, ByteBuf in, List<Object> out) throws Exception {
        if(logger.isTraceEnabled()) {
            logger.trace("Got Data: {}", ByteBufUtil.hexDump(in));
        }
        // If at least 4 bytes are readable, peek into them (without changing the read position)
        // and get the packet length. Only if the available amount of readable bytes is larger or
        // equal to this, continue processing the rest.
        if(in.readableBytes() >= 4) {
        }
    }

}
