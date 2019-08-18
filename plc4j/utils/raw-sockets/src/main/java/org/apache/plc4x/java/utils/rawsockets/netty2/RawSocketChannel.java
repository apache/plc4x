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

package org.apache.plc4x.java.utils.rawsockets.netty2;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelConfig;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelPromise;
import io.netty.channel.oio.OioByteStreamChannel;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannelConfig;
import org.pcap4j.core.NotOpenException;
import org.pcap4j.core.PacketListener;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.net.SocketTimeoutException;
import java.util.concurrent.TimeoutException;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-16
 */
public class RawSocketChannel extends OioByteStreamChannel {

    private static final Logger logger = LoggerFactory.getLogger(RawSocketChannel.class);

    private final RawSocketChannelConfig config;
    private PcapNetworkInterface nif;
    private PcapHandle handle;
    public static ByteBuf buffer;
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
        logger.debug("Connecting...");
        nif = Pcaps.getDevByName("en0");
        handle = nif.openLive(65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
        buffer = Unpooled.buffer();
        // Start loop in another Thread
        loopThread = new Thread(() -> {
            try {
                handle.loop(-1, (PacketListener) packet -> {
                    // logger.debug("Captured Packet from PCAP with length {} bytes", packet.getRawData().length);
                    buffer.writeBytes(packet.getRawData());
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
        activate(new PcapInputStream(buffer), new DiscardingOutputStream());
    }

    @Override
    protected SocketAddress localAddress0() {
        return new InetSocketAddress(1234);
    }

    @Override
    protected SocketAddress remoteAddress0() {
        return new InetSocketAddress(1234);
    }

    @Override
    protected void doBind(SocketAddress localAddress) throws Exception {
        throw new UnsupportedOperationException("");
    }

    @Override
    protected void doDisconnect() throws Exception {
        this.loopThread.interrupt();
        if (this.handle != null) {
            this.handle.close();
        }
    }

    @Override protected int doReadBytes(ByteBuf buf) throws Exception {
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

    private static class DiscardingOutputStream extends OutputStream {
        @Override
        public void write(int b) throws IOException {
            // discard
            logger.debug("Discarding {}", b);
        }
    }

    private static class PcapInputStream extends InputStream {

        final ByteBuf buf;

        private PcapInputStream(ByteBuf buf) {
            this.buf = buf;
        }

        @Override
        public int available() throws IOException {
            return buf.readableBytes();
        }

        @Override
        public int read() throws IOException {
            // Timeout 10 ms
            final long timeout = System.nanoTime() + 10_000;
            while (System.nanoTime() < timeout) {
                if (buf.readableBytes() > 0) {
                    return buf.readByte();
                }
            }
            throw new SocketTimeoutException();
        }

    }

    public class RawSocketUnsafe extends AbstractUnsafe {

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
