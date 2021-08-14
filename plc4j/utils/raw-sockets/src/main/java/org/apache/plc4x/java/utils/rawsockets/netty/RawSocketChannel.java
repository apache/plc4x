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
package org.apache.plc4x.java.utils.rawsockets.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelConfig;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelPromise;
import io.netty.channel.oio.OioByteStreamChannel;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.utils.pcap.netty.exception.PcapException;
import org.apache.plc4x.java.utils.rawsockets.netty.address.RawSocketAddress;
import org.apache.plc4x.java.utils.rawsockets.netty.config.RawSocketChannelConfig;
import org.pcap4j.core.*;
import org.pcap4j.packet.Packet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.SocketAddress;
import java.net.SocketTimeoutException;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-16
 */
public class RawSocketChannel extends OioByteStreamChannel {

    private static final Logger logger = LoggerFactory.getLogger(RawSocketChannel.class);

    private final RawSocketChannelConfig config;

    private RawSocketAddress remoteRawSocketAddress;
    private SocketAddress localAddress;
    private PcapHandle handle;
    private Thread loopThread;

    public RawSocketChannel() {
        super(null);
        config = new RawSocketChannelConfig(this);
    }

    @Override
    protected boolean isInputShutdown() {
        return false;
    }

    @Override
    protected ChannelFuture shutdownInput() {
        throw new NotImplementedException("");
    }

    @Override
    protected void doConnect(SocketAddress remoteAddress, SocketAddress localAddress) throws Exception {
        if(!(remoteAddress instanceof RawSocketAddress)) {
            logger.error("Expecting remote address of type RawSocketAddress");
            pipeline().fireExceptionCaught(
                new PcapException("Expecting remote address of type RawSocketAddress"));
            return;
        }
        this.localAddress = localAddress;
        remoteRawSocketAddress = (RawSocketAddress) remoteAddress;

        // Try to get the device name of the network interface that we want to open.
        String deviceName = getDeviceName(remoteRawSocketAddress);
        if(deviceName == null) {
            logger.error("Network device not specified and couldn't detect it automatically");
            pipeline().fireExceptionCaught(
                new PcapException("Network device not specified and couldn't detect it automatically"));
            return;
        }

        // Get a handle to the network-device and open it.
        PcapNetworkInterface nif = Pcaps.getDevByName(deviceName);
        handle = nif.openLive(65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
        if(logger.isDebugEnabled()) {
            logger.debug(String.format("Listening on device %s", deviceName));
        }

        // If the address allows fine tuning which packets to process, set a filter to reduce the load.
        String filter = config.getFilterString(localAddress, remoteAddress);
        if(filter.length() > 0) {
            handle.setFilter(filter, BpfProgram.BpfCompileMode.OPTIMIZE);
        }

        // Create a buffer where the raw socket worker can send data to.
        ByteBuf buffer = Unpooled.buffer();

        // Start a thread that processes the callbacks from the raw socket and simply
        // forwards the bytes read to the buffer.
        loopThread = new Thread(() -> {
            try {
                handle.loop(-1, new PacketListener() {
                    @Override
                    public void gotPacket(Packet packet) {
                        buffer.writeBytes(config.getPacketHandler().getData(packet));
                    }
                });
            } catch (PcapNativeException | NotOpenException e) {
                // TODO this should close everything automatically
                logger.error("Pcap4j loop thread died!", e);
                pipeline().fireExceptionCaught(e);
            } catch (InterruptedException e) {
                logger.warn("PCAP Loop Thread was interrupted (hopefully intentionally)", e);
                Thread.currentThread().interrupt();
            }
        });
        loopThread.start();

        // Right now we're using an output stream that simply discards everything.
        // This is ok while implementing passive drivers for protocols, however as
        // soon as we start implementing ethernet layer protocols, we'll have to also
        // be able to actually send data. The PcapInputStream simply acts as a
        // breaking point if no packets are coming in and the read operation would
        // simply block indefinitely.
        activate(new PcapInputStream(buffer), new DiscardingOutputStream());
    }

    @Override
    protected SocketAddress localAddress0() {
        return localAddress;
    }

    @Override
    protected SocketAddress remoteAddress0() {
        return remoteRawSocketAddress;
    }

    @Override
    protected void doBind(SocketAddress localAddress) {
        throw new UnsupportedOperationException("Not implemented");
    }

    @Override
    protected void doDisconnect() {
        this.loopThread.interrupt();
        if (this.handle != null) {
            this.handle.close();
        }
    }

    @Override
    protected int doReadBytes(ByteBuf buf) throws Exception {
        if (handle == null || !handle.isOpen()) {
            return -1;
        }
        try {
            return super.doReadBytes(buf);
        } catch (SocketTimeoutException ignored) {
            return 0;
        }
    }

    @Override
    public ChannelConfig config() {
        return this.config;
    }

    @Override
    public boolean isOpen() {
        return true;
    }

    @Override
    protected AbstractUnsafe newUnsafe() {
        return new RawSocketUnsafe();
    }

    private String getDeviceName(RawSocketAddress rawSocketAddress) {
        // If the device name is provided, simply use this.
        if(rawSocketAddress.getDeviceName() != null) {
            return rawSocketAddress.getDeviceName();
        }

        // TODO: Implement this ...
        return null;
    }

    /**
     * This output stream simply discards anything it should send.
     */
    private static class DiscardingOutputStream extends OutputStream {
        @Override
        public void write(int b) throws IOException {
            // discard
            logger.debug("Discarding {}", b);
        }

        @Override
        public void write(byte[] b, int off, int len) {
            logger.debug("Discarding {}", b);
        }
    }

    /**
     * InputStream that fulfills the contract of Netty for read operations to timeout.
     * Without this the InputStream would simply block indefinitely which would block
     * the entire IO module.
     */
    private static class PcapInputStream extends InputStream {
        final ByteBuf buf;

        private PcapInputStream(ByteBuf buf) {
            this.buf = buf;
        }

        @Override
        public int available() {
            return buf.readableBytes();
        }

        @Override
        public int read() throws IOException {
            // Timeout 10 ms
            final long timeout = System.nanoTime() + 10_000;
            // Yes, this could make the thread go nuts in case of no data,
            // but the Netty guys are doing it the same way and there probably
            // is a reason for it ;-)
            while (System.nanoTime() < timeout) {
                if (buf.readableBytes() > 0) {
                    return buf.readByte() & 0xFF;
                }
            }
            throw new SocketTimeoutException();
        }
    }

    /**
     * Internal helper to wrap access to unsafe operations (Only used internally by netty)
     */
    private class RawSocketUnsafe extends AbstractUnsafe {
        @Override
        public void connect(SocketAddress remoteAddress, SocketAddress localAddress, ChannelPromise promise) {
            try {
                doConnect(remoteAddress, localAddress);
                pipeline().fireChannelActive();
                promise.setSuccess();
            } catch (Exception e) {
                promise.setFailure(e);
            }
        }
    }
}
