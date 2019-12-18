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
package org.apache.plc4x.java.s7.readwrite.connection;

import io.netty.channel.*;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.protocol.Plc4xProtocolBase;
import org.apache.plc4x.java.s7.readwrite.protocol.Plc4xNettyWrapper;
import org.apache.plc4x.java.s7.readwrite.protocol.Plc4xS7Protocol;
import org.apache.plc4x.java.s7.readwrite.protocol.S7Protocol;
import org.apache.plc4x.java.s7.readwrite.types.COTPTpduSize;
import org.apache.plc4x.java.s7.readwrite.types.DeviceGroup;
import org.apache.plc4x.java.s7.readwrite.types.S7ControllerType;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcFieldHandler;
import org.apache.plc4x.java.s7.readwrite.utils.S7TsapIdEncoder;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.concurrent.CompletableFuture;

public class S7Connection extends NettyPlcConnection implements PlcReader, PlcWriter {

    private static final int ISO_ON_TCP_PORT = 102;

    private static final Logger logger = LoggerFactory.getLogger(S7Connection.class);

    private final short rack;
    private final short slot;
    private final COTPTpduSize tpduSize;
    private final short maxAmqCaller;
    private final short maxAmqCallee;
    private final S7ControllerType controllerType;

    public S7Connection(InetAddress address, String params) {
        this(new TcpSocketChannelFactory(address, ISO_ON_TCP_PORT), params);
    }

    public S7Connection(ChannelFactory channelFactory, String params) {
        super(channelFactory, true);

        short curRack = 1;
        short curSlot = 1;
        short curParamPduSize = 1024;
        short curParamMaxAmqCaller = 8;
        short curParamMaxAmqCallee = 8;
        S7ControllerType curParamControllerType = S7ControllerType.ANY;

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        case "rack":
                            curRack = Short.parseShort(paramValue);
                            break;
                        case "slot":
                            curSlot = Short.parseShort(paramValue);
                            break;
                        case "pdu-size":
                            curParamPduSize = Short.parseShort(paramValue);
                            break;
                        case "max-amq-caller":
                            curParamMaxAmqCaller = Short.parseShort(paramValue);
                            break;
                        case "max-amq-callee":
                            curParamMaxAmqCallee = Short.parseShort(paramValue);
                            break;
                        case "controller-type":
                            curParamControllerType = S7ControllerType.valueOf(paramValue);
                            break;
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }

        // It seems that the LOGO devices are a little picky about the pdu-size.
        // Instead of handling this out, they just hang up without any error message.
        // So in case of a LOGO controller, set this to a known working value.
        if(curParamControllerType == S7ControllerType.LOGO && curParamPduSize == 1024) {
            curParamPduSize = 480;
        }

        this.tpduSize = getNearestMatchingTpduSize(curParamPduSize);

        this.rack = curRack;
        this.slot = curSlot;
        this.maxAmqCallee = curParamMaxAmqCallee;
        this.maxAmqCaller = curParamMaxAmqCaller;
        this.controllerType = curParamControllerType;

        logger.info("Setting up S7 Connection with: rack {}, slot {}, tpdu-size {}, max-amq-caller {}, " +
                "max-amq-callee {}", rack, slot, tpduSize, maxAmqCaller, maxAmqCallee);
    }

    @Override
    public boolean canRead() {
        return true;
    }

    @Override
    public boolean canWrite() {
        return true;
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OS, 0, 0);
        short callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, rack, slot);

        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                // Build the protocol stack for communicating with the s7 protocol.
                ChannelPipeline pipeline = channel.pipeline();
                pipeline.addLast(new ChannelInboundHandlerAdapter() {
                    @Override
                    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
                        if (evt instanceof ConnectedEvent) {
                            sessionSetupCompleteFuture.complete(null);
                        } else {
                            super.userEventTriggered(ctx, evt);
                        }
                    }
                });
                pipeline.addLast(new S7Protocol());
                Plc4xProtocolBase<TPKTPacket> plc4xS7Protocol = new Plc4xS7Protocol(callingTsapId, calledTsapId, tpduSize,
                    maxAmqCaller, maxAmqCallee, controllerType);
                Plc4xNettyWrapper<TPKTPacket> context = plc4xS7Protocol.getContext();
                pipeline.addLast(context);
            }
        };
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new S7PlcFieldHandler());
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        return new DefaultPlcWriteRequest.Builder(this, new S7PlcFieldHandler());
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        InternalPlcReadRequest internalReadRequest = checkInternal(readRequest, InternalPlcReadRequest.class);
        CompletableFuture<InternalPlcReadResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>(internalReadRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcReadResponse.class::cast);
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        InternalPlcWriteRequest internalWriteRequest = checkInternal(writeRequest, InternalPlcWriteRequest.class);
        CompletableFuture<InternalPlcWriteResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcWriteRequest, InternalPlcWriteResponse> container =
            new PlcRequestContainer<>(internalWriteRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcWriteResponse.class::cast);
    }

    @Override
    protected void sendChannelCreatedEvent() {
        logger.trace("Channel was created, firing ChannelCreated Event");
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

    /**
     * Iterate over all values until one is found that the given tpdu size will fit.
     * @param tpduSizeParameter requested tpdu size.
     * @return smallest {@link COTPTpduSize} which will fit a given size of tpdu.
     */
    protected COTPTpduSize getNearestMatchingTpduSize(short tpduSizeParameter) {
        for (COTPTpduSize value : COTPTpduSize.values()) {
            if(value.getSizeInBytes() >= tpduSizeParameter) {
                return value;
            }
        }
        return null;
    }

}
