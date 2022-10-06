/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.transport.socketcan.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelConfig;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelPromise;
import io.netty.channel.EventLoop;
import io.netty.channel.oio.OioByteStreamChannel;
import org.apache.plc4x.java.transport.socketcan.netty.address.SocketCANAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tel.schich.javacan.CanChannels;
import tel.schich.javacan.NetworkDevice;
import tel.schich.javacan.RawCanChannel;
import tel.schich.javacan.platform.linux.LinuxNetworkDevice;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.SocketAddress;
import java.net.SocketTimeoutException;
import java.nio.ByteBuffer;

/**
 * A naive copy of pcap channel implementation which allows to pass over data to PLC4X as well as drop it back to JavaCAN.
 * Sadly all this involves double parsing.
 */
public class SocketCANChannel extends OioByteStreamChannel {

    private final Logger logger = LoggerFactory.getLogger(SocketCANChannel.class);

    private final SocketCANChannelConfig config;

    private SocketCANAddress remoteRawSocketAddress;
    private SocketAddress localAddress;
    private RawCanChannel handle;
    private Thread loopThread;

    public SocketCANChannel() {
        super(null);
        config = new SocketCANChannelConfig(this);
    }

    @Override
    protected boolean isInputShutdown() {
        return false;
    }

    @Override
    protected ChannelFuture shutdownInput() {
        throw new UnsupportedOperationException("");
    }

    @Override
    protected void doConnect(SocketAddress remoteAddress, SocketAddress localAddress) throws Exception {
        if (!(remoteAddress instanceof SocketCANAddress)) {
            logger.error("Expecting remote address of type SocketCANAddress");
            pipeline().fireExceptionCaught(new IllegalArgumentException("Expecting remote address of type SocketCANAddress"));
            return;
        }
        this.localAddress = localAddress;
        remoteRawSocketAddress = (SocketCANAddress) remoteAddress;

        // Try to get the device name of the network interface that we want to open.
        String interfaceName = getInterfaceName(remoteRawSocketAddress);
        if(interfaceName == null) {
            logger.error("Interface name is not specified and couldn't detect it automatically");
            pipeline().fireExceptionCaught(new IllegalArgumentException("Interface name is not specified and couldn't detect it automatically"));
            return;
        }

        NetworkDevice device = LinuxNetworkDevice.lookup(interfaceName);
        // Get a handle to the network-device and open it.
        handle = CanChannels.newRawChannel(device);

        if(logger.isDebugEnabled()) {
            logger.debug(String.format("Listening on device %s", interfaceName));
        }

        // TODO If the address allows fine tuning which packets to process, set a filter to reduce the load.
//        String filter =
//        if(filter.length() > 0) {
//            handle.setFilter(filter, BpfProgram.BpfCompileMode.OPTIMIZE);
//        }

        // Create a buffer where the raw socket worker can send data to.
        ByteBuf buffer = Unpooled.buffer();

        // Start a thread that processes the callbacks from the raw socket and simply
        // forwards the bytes read to the buffer.
        loopThread = new Thread(() -> {
            try {
                ByteBuffer byteBuffer = ByteBuffer.allocateDirect(16);
                while (!isInputShutdown()) {
                    handle.readUnsafe(byteBuffer);
                    buffer.writeBytes(byteBuffer);
                    byteBuffer.rewind();
                }
            } catch (Exception e) {
                logger.warn("Could not read data", e);
                pipeline().fireExceptionCaught(e);
            } catch (Throwable e) {
                logger.warn("Fatal error while handling CAN communication", e);
                pipeline().fireExceptionCaught(e);
            }
        }, "javacan-reader");
        loopThread.start();

        activate(new CANInputStream(buffer), new CANOutputStream(handle));
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
            try {
                this.handle.close();
            } catch (IOException e) {
                logger.error("Failed to close CAN socket!");
            }
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
    protected boolean isCompatible(EventLoop loop) {
        return super.isCompatible(loop);
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

    private String getInterfaceName(SocketCANAddress address) {
        // If the device name is provided, simply use this.
        if(address.getInterfaceName() != null) {
            return address.getInterfaceName();
        }

        // TODO: Implement this ...
        return null;
    }

    /**
     * InputStream that fulfills the contract of Netty for read operations to timeout.
     * Without this the InputStream would simply block indefinitely which would block
     * the entire IO module.
     */
    private static class CANInputStream extends InputStream {
        final ByteBuf buf;

        private CANInputStream(ByteBuf buf) {
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
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    throw new IOException(e);
                }
            }
            throw new SocketTimeoutException();
        }

    }


    private static class CANOutputStream extends OutputStream {

        private final RawCanChannel rawCanChannel;
        private final ByteBuffer buffer = ByteBuffer.allocateDirect(16);

        public CANOutputStream(RawCanChannel rawCanChannel) {
            this.rawCanChannel = rawCanChannel;
        }

        @Override
        public void write(int b) throws IOException {
            throw new IOException("Appending single bytes is not permitted. Use write(byte[], int, int)");
        }

        @Override
        public void write(byte[] b, int off, int len) throws IOException {
            buffer.put(b, off, len);
            buffer.flip();
            rawCanChannel.writeUnsafe(buffer);
            // the NIO clear call does clear positions and limits, but does not erase contents of buffer.
            // since we write a complete frame with length any remaining data at the end of memory block
            // should be ignored. When we get payload of full length it will eventually override earlier data
            // which was not used
            buffer.clear();
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
