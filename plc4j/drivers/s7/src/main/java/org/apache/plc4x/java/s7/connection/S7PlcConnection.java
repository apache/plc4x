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
package org.apache.plc4x.java.s7.connection;

import io.netty.buffer.Unpooled;
import io.netty.channel.*;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.isoontcp.protocol.IsoOnTcpProtocol;
import org.apache.plc4x.java.isotp.protocol.IsoTPProtocol;
import org.apache.plc4x.java.isotp.protocol.model.tpdus.DisconnectRequestTpdu;
import org.apache.plc4x.java.isotp.protocol.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.protocol.model.types.DisconnectReason;
import org.apache.plc4x.java.isotp.protocol.model.types.TpduSize;
import org.apache.plc4x.java.s7.model.S7Field;
import org.apache.plc4x.java.s7.netty.Plc4XS7Protocol;
import org.apache.plc4x.java.s7.netty.S7Protocol;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.strategies.DefaultS7MessageProcessor;
import org.apache.plc4x.java.s7.netty.util.S7PlcFieldHandler;
import org.apache.plc4x.java.s7.types.S7ControllerType;
import org.apache.plc4x.java.s7.utils.S7TsapIdEncoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Collections;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * Class implementing the Connection handling for Siemens S7.
 * The adressing of Values in S7 works as follows:
 * <p>
 * For adressing values from Datablocks the following syntax is used:
 * <pre>
 *     DATA_BLOCKS/{blockNumer}/{byteOffset}
 * </pre>
 * <p>
 * For adressing data from other memory segments like I/O, Markers, ...
 * <pre>
 *     {memory area}/{byte offset}
 *     or
 *     {memory area}/{byte offset}/{bit offset}
 * </pre>
 * where the {bit-offset} is optional.
 * All Available Memory Areas for this mode are defined in the {@link MemoryArea} enum.
 */
public class S7PlcConnection extends NettyPlcConnection implements PlcReader, PlcWriter {

    private static final int ISO_ON_TCP_PORT = 102;

    // Fetch values from configuration
    private static final Configuration CONF = new SystemConfiguration();
    private static final long CLOSE_DEVICE_TIMEOUT_MS = CONF.getLong("plc4x.s7connection.close.device,timeout", 10_000);

    private static final Logger logger = LoggerFactory.getLogger(S7PlcConnection.class);

    private final int rack;
    private final int slot;

    private final short paramPduSize;
    private final short paramMaxAmqCaller;
    private final short paramMaxAmqCallee;
    private final S7ControllerType paramControllerType;
    private DeviceGroup deviceGroup = DeviceGroup.OS;

    public S7PlcConnection(InetAddress address, int rack, int slot, String params) {
        this(new TcpSocketChannelFactory(address, ISO_ON_TCP_PORT), rack, slot, params);

        logger.info("Setting up S7 Connection with: host-name {}, rack {}, slot {}, pdu-size {}, max-amq-caller {}, " +
                "max-amq-callee {}", address.getHostAddress(), rack, slot,
            paramPduSize, paramMaxAmqCaller, paramMaxAmqCallee);
    }

    public S7PlcConnection(ChannelFactory channelFactory, int rack, int slot, String params) {
        super(channelFactory, true);

        this.rack = rack;
        this.slot = slot;

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
                        case "device-group":
                            if (Objects.equals(paramValue, "others")) {
                                this.deviceGroup = DeviceGroup.OTHERS;
                            } else {
                                throw new RuntimeException("Only OTHERS as Device Group is allowed!");
                            }
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

        // IsoTP uses pre defined sizes. Find the smallest box,
        // that would be able to contain the requested pdu size.
        this.paramPduSize = curParamPduSize;
        this.paramMaxAmqCaller = curParamMaxAmqCaller;
        this.paramMaxAmqCallee = curParamMaxAmqCallee;
        this.paramControllerType = curParamControllerType;
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
        short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(this.deviceGroup, rack, slot);
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
                pipeline.addLast(new IsoOnTcpProtocol());
                pipeline.addLast(new IsoTPProtocol(callingTsapId, calledTsapId, TpduSize.valueForGivenSize(paramPduSize)));
                pipeline.addLast(new S7Protocol(paramMaxAmqCaller, paramMaxAmqCallee, paramPduSize, paramControllerType,
                    new DefaultS7MessageProcessor()));
                pipeline.addLast(new Plc4XS7Protocol());
            }
        };
    }

    @Override
    protected void sendChannelCreatedEvent() {
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

    @Override
    public PlcField prepareField(String fieldQuery) throws PlcInvalidFieldException {
        return S7Field.of(fieldQuery);
    }

    public int getRack() {
        return rack;
    }

    public int getSlot() {
        return slot;
    }

    public short getParamPduSize() {
        return paramPduSize;
    }

    public int getParamMaxAmqCaller() {
        return paramMaxAmqCaller;
    }

    public int getParamMaxAmqCallee() {
        return paramMaxAmqCallee;
    }

    public S7ControllerType getParamControllerType() {
        return paramControllerType;
    }

    @Override
    public void close() throws PlcConnectionException {
        logger.debug("S7 Connection's close method is triggered");
        if ((channel != null) && channel.isOpen()) {
            // Send the PLC a message that the connection is being closed.
            DisconnectRequestTpdu disconnectRequest = new DisconnectRequestTpdu(
                (short) 0x0000, (short) 0x000F, DisconnectReason.NORMAL, Collections.emptyList(),
                Unpooled.EMPTY_BUFFER);

            // In case of an ISO TP Class 0 connection, the remote is usually expected to actively
            // close the connection. So we add a listener waiting for this to happen.
            CompletableFuture<Void> disconnectFuture = new CompletableFuture<>();
            channel.closeFuture().addListener(
                (ChannelFutureListener) future -> disconnectFuture.complete(null));

            // Send the disconnect request.
            logger.trace("Sending disconnect request to PLC");
            channel.writeAndFlush(disconnectRequest);
            // Wait for the configured time for the remote to close the session.
            try {
                disconnectFuture.get(CLOSE_DEVICE_TIMEOUT_MS, TimeUnit.MILLISECONDS);
                logger.trace("Got disconnect response from PLC, can close channel now.");
            }
            // If the remote didn't close the connection within the given time-frame, we have to take
            // care of closing the connection.
            catch (TimeoutException e) {
                logger.debug("Remote didn't close connection within the configured timeout of {} ms, shutting down actively.", CLOSE_DEVICE_TIMEOUT_MS, e);
                channel.close();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } catch (ExecutionException e) {
                throw new PlcConnectionException(e);
            } finally {
                // Do some additional cleanup operations ...
                // In normal operation, the channels event loop has a parent, however when running with
                // the embedded channel for unit tests, parent is null.
                if (channel.eventLoop().parent() != null) {
                    channel.eventLoop().parent().shutdownGracefully();
                }
            }
        }
        super.close();
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

}
