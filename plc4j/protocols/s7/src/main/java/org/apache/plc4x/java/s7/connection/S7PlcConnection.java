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

import io.netty.bootstrap.Bootstrap;
import io.netty.channel.*;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioSocketChannel;
import org.apache.plc4x.java.api.connection.AbstractPlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.isoontcp.netty.IsoOnTcpProtocol;
import org.apache.plc4x.java.isotp.netty.IsoTPProtocol;
import org.apache.plc4x.java.isotp.netty.model.tpdus.DisconnectRequestTpdu;
import org.apache.plc4x.java.isotp.netty.model.types.DisconnectReason;
import org.apache.plc4x.java.isotp.netty.model.types.TpduSize;
import org.apache.plc4x.java.netty.events.S7ConnectionEvent;
import org.apache.plc4x.java.netty.events.S7ConnectionState;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.Plc4XS7Protocol;
import org.apache.plc4x.java.s7.netty.S7Protocol;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Collections;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class S7PlcConnection extends AbstractPlcConnection implements PlcReader, PlcWriter {

    private static final int ISO_ON_TCP_PORT = 102;

    private static final Pattern S7_DATABLOCK_ADDRESS_PATTERN =
        Pattern.compile("^DATA_BLOCKS/(?<blockNumber>\\d{1,4})/(?<byteOffset>\\d{1,4})");
    private static final Pattern S7_ADDRESS_PATTERN =
        Pattern.compile("^(?<memoryArea>.*?)/(?<byteOffset>\\d{1,4})(?:/(?<bitOffset>\\d))?");

    private final static Logger logger = LoggerFactory.getLogger(S7PlcConnection.class);

    private final String hostName;
    private final int rack;
    private final int slot;

    private int pduSize;

    private EventLoopGroup workerGroup;
    private Channel channel;

    public S7PlcConnection(String hostName, int rack, int slot) {
        this.hostName = hostName;
        this.rack = rack;
        this.slot = slot;
        this.pduSize = 1024;
    }

    public String getHostName() {
        return hostName;
    }

    public int getRack() {
        return rack;
    }

    public int getSlot() {
        return slot;
    }

    public int getPduSize() {
        return pduSize;
    }

    @Override
    public void connect() throws PlcConnectionException {
        workerGroup = new NioEventLoopGroup();

        try {
            // As we don't just want to wait till the connection is established,
            // define a future we can use to signal back that the s7 session is
            // finished initializing.
            CompletableFuture<Void> sessionSetupCompleteFuture = new CompletableFuture<>();

            InetAddress serverInetAddress = InetAddress.getByName(hostName);

            Bootstrap bootstrap = new Bootstrap();
            bootstrap.group(workerGroup);
            bootstrap.channel(NioSocketChannel.class);
            bootstrap.option(ChannelOption.SO_KEEPALIVE, true);
            bootstrap.option(ChannelOption.TCP_NODELAY, true);
            bootstrap.handler(new ChannelInitializer() {
                @Override
                protected void initChannel(Channel channel) throws Exception {
                    // Build the protocol stack for communicating with the s7 protocol.
                    ChannelPipeline pipeline = channel.pipeline();
                    pipeline.addLast(new ChannelInboundHandlerAdapter() {
                        @Override
                        public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
                            if(evt instanceof S7ConnectionEvent &&
                                ((S7ConnectionEvent) evt).getState() == S7ConnectionState.SETUP_COMPLETE) {
                                sessionSetupCompleteFuture.complete(null);
                            } else {
                                super.userEventTriggered(ctx, evt);
                            }
                        }
                    });
                    pipeline.addLast(new IsoOnTcpProtocol());
                    pipeline.addLast(new IsoTPProtocol((byte) rack, (byte) slot, TpduSize.SIZE_1024));
                    pipeline.addLast(new S7Protocol((short) 8, (short) 8, (short) 1024));
                    pipeline.addLast(new Plc4XS7Protocol());
                }
            });
            // Start the client.
            ChannelFuture f = bootstrap.connect(serverInetAddress, ISO_ON_TCP_PORT).sync();
            f.awaitUninterruptibly();
            // Wait till the session is finished initializing.
            channel = f.channel();

            // Send an event to the pipeline telling the Protocol filters what's going on.
            channel.pipeline().fireUserEventTriggered(new S7ConnectionEvent());

            sessionSetupCompleteFuture.get();
        } catch (UnknownHostException e) {
            throw new PlcConnectionException("Unknown Host " + hostName, e);
        } catch (InterruptedException | ExecutionException e) {
            throw new PlcConnectionException(e);
        }
    }

    @Override
    public void close() throws Exception {
        if((channel != null) && channel.isOpen()) {
            // Send the PLC a message that the connection is being closed.
            DisconnectRequestTpdu disconnectRequest = new DisconnectRequestTpdu(
                (short) 0x0000, (short) 0x000F, DisconnectReason.NORMAL, Collections.emptyList(),
                null);
            ChannelFuture sendDisconnectRequestFuture = channel.writeAndFlush(disconnectRequest);
            sendDisconnectRequestFuture.addListener((ChannelFutureListener) future -> {
                // Close the session itself.
                channel.closeFuture().await();
                workerGroup.shutdownGracefully();
            });
            sendDisconnectRequestFuture.awaitUninterruptibly();
        } else if (workerGroup != null) {
            workerGroup.shutdownGracefully();
        }
    }

    @Override
    public Address parseAddress(String addressString) throws PlcException {
        Matcher datablockAddressMatcher = S7_DATABLOCK_ADDRESS_PATTERN.matcher(addressString);
        if (datablockAddressMatcher.matches()) {
            int datablockNumber = Integer.valueOf(datablockAddressMatcher.group("blockNumber"));
            int datablockByteOffset = Integer.valueOf(datablockAddressMatcher.group("byteOffset"));
            return new S7DataBlockAddress((short) datablockNumber, (short) datablockByteOffset);
        }
        Matcher addressMatcher = S7_ADDRESS_PATTERN.matcher(addressString);
        if (!addressMatcher.matches()) {
            throw new PlcConnectionException(
                "Address string doesn't match the format '{area}/{offset}[/{bit-offset}]'");
        }
        MemoryArea memoryArea = MemoryArea.valueOf(addressMatcher.group("memoryArea"));
        int byteOffset = Integer.valueOf(addressMatcher.group("byteOffset"));
        String bitOffset =  addressMatcher.group("bitOffset");
        if (bitOffset != null) {
            return new S7BitAddress(memoryArea, (short) byteOffset, Byte.valueOf(bitOffset));
        }
        return new S7Address(memoryArea, (short) byteOffset);
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> readFuture = new CompletableFuture<>();
        PlcRequestContainer<PlcReadRequest, PlcReadResponse> container =
            new PlcRequestContainer<>(readRequest, readFuture);
        channel.writeAndFlush(container);
        return readFuture;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> writeFuture = new CompletableFuture<>();
        PlcRequestContainer<PlcWriteRequest, PlcWriteResponse> container =
            new PlcRequestContainer<>(writeRequest, writeFuture);
        channel.writeAndFlush(container);
        return writeFuture;
    }

}
