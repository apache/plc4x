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

import io.netty.channel.*;
import org.apache.commons.configuration2.Configuration;
import org.apache.commons.configuration2.SystemConfiguration;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.base.connection.AbstractPlcConnection;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.isoontcp.netty.IsoOnTcpProtocol;
import org.apache.plc4x.java.isotp.netty.IsoTPProtocol;
import org.apache.plc4x.java.isotp.netty.model.tpdus.DisconnectRequestTpdu;
import org.apache.plc4x.java.isotp.netty.model.types.DeviceGroup;
import org.apache.plc4x.java.isotp.netty.model.types.DisconnectReason;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.apache.plc4x.java.s7.netty.Plc4XS7Protocol;
import org.apache.plc4x.java.s7.netty.S7Protocol;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.util.S7PlcFieldHandler;
import org.apache.plc4x.java.s7.utils.S7TsapIdEncoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Collections;
import java.util.Optional;
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
public class S7PlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    private static final int ISO_ON_TCP_PORT = 102;

    // Fetch values from configuration
    private static final Configuration CONF = new SystemConfiguration();
    private static final long CLOSE_DEVICE_TIMEOUT_MS = CONF.getLong("plc4x.s7connection.close.device,timeout", 1_000);

    private static final Logger logger = LoggerFactory.getLogger(S7PlcConnection.class);

    private final int rack;
    private final int slot;

    private final TpduSize paramPduSize;
    private final short paramMaxAmqCaller;
    private final short paramMaxAmqCallee;

    public S7PlcConnection(InetAddress address, int rack, int slot, String params) {
        this(new TcpSocketChannelFactory(address, ISO_ON_TCP_PORT), rack, slot, params);

        logger.info("Setting up S7cConnection with: host-name {}, rack {}, slot {}, pdu-size {}, max-amq-caller {}, " +
                "max-amq-callee {}", address.getHostAddress(), rack, slot,
            paramPduSize.getValue(), paramMaxAmqCaller, paramMaxAmqCallee);
    }

    public S7PlcConnection(ChannelFactory channelFactory, int rack, int slot, String params) {
        super(channelFactory, true);

        this.rack = rack;
        this.slot = slot;

        int paramPduSize = 1024;
        short paramMaxAmqCaller = 8;
        short paramMaxAmqCallee = 8;

        if (!StringUtils.isEmpty(params)) {
            for (String param : params.split("&")) {
                String[] paramElements = param.split("=");
                String paramName = paramElements[0];
                if (paramElements.length == 2) {
                    String paramValue = paramElements[1];
                    switch (paramName) {
                        case "pdu-size":
                            paramPduSize = Integer.parseInt(paramValue);
                            break;
                        case "max-amq-caller":
                            paramMaxAmqCaller = Short.parseShort(paramValue);
                            break;
                        case "max-amq-callee":
                            paramMaxAmqCallee = Short.parseShort(paramValue);
                            break;
                        default:
                            logger.debug("Unknown parameter {} with value {}", paramName, paramValue);
                    }
                } else {
                    logger.debug("Unknown no-value parameter {}", paramName);
                }
            }
        }

        // IsoTP uses pre defined sizes. Find the smallest box,
        // that would be able to contain the requested pdu size.
        this.paramPduSize = TpduSize.valueForGivenSize(paramPduSize);
        this.paramMaxAmqCaller = paramMaxAmqCaller;
        this.paramMaxAmqCallee = paramMaxAmqCallee;
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        short calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, 0, 0);
        short callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OTHERS, rack, slot);

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
                pipeline.addLast(new IsoTPProtocol(callingTsapId, calledTsapId, paramPduSize));
                pipeline.addLast(new S7Protocol(paramMaxAmqCaller, paramMaxAmqCallee, (short) paramPduSize.getValue()));
                pipeline.addLast(new Plc4XS7Protocol());
            }
        };
    }

    @Override
    protected void sendChannelCreatedEvent() {
        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new ConnectEvent());
    }

    public int getRack() {
        return rack;
    }

    public int getSlot() {
        return slot;
    }

    public TpduSize getParamPduSize() {
        return paramPduSize;
    }

    public int getParamMaxAmqCaller() {
        return paramMaxAmqCaller;
    }

    public int getParamMaxAmqCallee() {
        return paramMaxAmqCallee;
    }

    @Override
    public void close() throws PlcConnectionException {
        if ((channel != null) && channel.isOpen()) {
            // Send the PLC a message that the connection is being closed.
            DisconnectRequestTpdu disconnectRequest = new DisconnectRequestTpdu(
                (short) 0x0000, (short) 0x000F, DisconnectReason.NORMAL, Collections.emptyList(),
                null);

            // In case of an ISO TP Class 0 connection, the remote is usually expected to actively
            // close the connection. So we add a listener waiting for this to happen.
            CompletableFuture<Void> disconnectFuture = new CompletableFuture<>();
            channel.closeFuture().addListener(
                (ChannelFutureListener) future -> disconnectFuture.complete(null));

            // Send the disconnect request.
            channel.writeAndFlush(disconnectRequest);
            // Wait for the configured time for the remote to close the session.
            try {
                disconnectFuture.get(CLOSE_DEVICE_TIMEOUT_MS, TimeUnit.MILLISECONDS);
            }
            // If the remote didn't close the connection within the given time-frame, we have to take
            // care of closing the connection.
            catch (TimeoutException e) {
                logger.info("Remote didn't close connection within the configured timeout of {}ms, shutting down actively.", CLOSE_DEVICE_TIMEOUT_MS, e);
                channel.close();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } catch (ExecutionException e) {
                throw new PlcConnectionException(e);
            }

            // Do some additional cleanup operations ...
            // In normal operation, the channels event loop has a parent, however when running with
            // the embedded channel for unit tests, parent is null.
            if (channel.eventLoop().parent() != null) {
                channel.eventLoop().parent().shutdownGracefully();
            }
        }
        super.close();
    }

    @Override
    public Optional<PlcReadRequest.Builder> readRequestBuilder() {
        return Optional.of(new DefaultPlcReadRequest.Builder(this, new S7PlcFieldHandler()));
    }

    @Override
    public Optional<PlcWriteRequest.Builder> writeRequestBuilder() {
        return Optional.of(new DefaultPlcWriteRequest.Builder(this, new S7PlcFieldHandler()));
    }

    @Override
    public Optional<PlcSubscriptionRequest.Builder> subscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public Optional<PlcUnsubscriptionRequest.Builder> unsubscriptionRequestBuilder() {
        return Optional.empty();
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<InternalPlcReadResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcReadRequest, InternalPlcReadResponse> container =
            new PlcRequestContainer<>((InternalPlcReadRequest) readRequest, future);
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
        CompletableFuture<InternalPlcWriteResponse> future = new CompletableFuture<>();
        PlcRequestContainer<InternalPlcWriteRequest, InternalPlcWriteResponse> container =
            new PlcRequestContainer<>((InternalPlcWriteRequest) writeRequest, future);
        channel.writeAndFlush(container).addListener(f -> {
            if (!f.isSuccess()) {
                future.completeExceptionally(f.cause());
            }
        });
        return future
            .thenApply(PlcWriteResponse.class::cast);
    }

}
