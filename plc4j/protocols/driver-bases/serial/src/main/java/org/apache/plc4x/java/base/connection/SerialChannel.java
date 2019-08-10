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

package org.apache.plc4x.java.base.connection;

import io.netty.buffer.ByteBuf;
import io.netty.channel.*;
import io.netty.channel.Channel;
import io.netty.channel.nio.AbstractNioByteChannel;
import io.netty.channel.nio.AbstractNioChannel;
import io.netty.channel.nio.NioEventLoop;
import io.netty.channel.socket.DuplexChannel;
import org.apache.commons.lang3.NotImplementedException;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.net.SocketAddress;
import java.nio.channels.*;
import java.nio.channels.spi.AbstractSelectableChannel;
import java.nio.channels.spi.AbstractSelector;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-10
 */
public class SerialChannel extends AbstractNioByteChannel implements DuplexChannel {


    public SerialChannel() {
        this(null, new SerialSocketChannel(new SerialSelectorProvider()));
    }

    /**
     * Create a new instance
     *
     * @param parent the parent {@link Channel} by which this instance was created. May be {@code null}
     * @param ch     the underlying {@link SelectableChannel} on which it operates
     */
    protected SerialChannel(Channel parent, SelectableChannel ch) {
        super(parent, ch);
    }

    @Override
    public NioUnsafe unsafe() {
        return new SerialNioUnsafe();
    }

    @Override
    public boolean isInputShutdown() {
        return false;
    }

    @Override
    public ChannelFuture shutdownInput() {
        return null;
    }

    @Override
    public ChannelFuture shutdownInput(ChannelPromise promise) {
        return null;
    }

    @Override
    public boolean isOutputShutdown() {
        return false;
    }

    @Override
    public ChannelFuture shutdownOutput() {
        return null;
    }

    @Override
    public ChannelFuture shutdownOutput(ChannelPromise promise) {
        return null;
    }

    @Override
    public boolean isShutdown() {
        return false;
    }

    @Override
    public ChannelFuture shutdown() {
        return null;
    }

    @Override
    public ChannelFuture shutdown(ChannelPromise promise) {
        return null;
    }

    @Override
    protected long doWriteFileRegion(FileRegion region) throws Exception {
        return 0;
    }

    @Override
    protected int doReadBytes(ByteBuf buf) throws Exception {
        return 0;
    }

    @Override
    protected int doWriteBytes(ByteBuf buf) throws Exception {
        return 0;
    }

    @Override
    protected boolean doConnect(SocketAddress remoteAddress, SocketAddress localAddress) throws Exception {
        return false;
    }

    @Override
    protected void doFinishConnect() throws Exception {

    }

    @Override
    protected SocketAddress localAddress0() {
        return null;
    }

    @Override
    protected SocketAddress remoteAddress0() {
        return null;
    }

    @Override
    protected void doBind(SocketAddress localAddress) throws Exception {

    }

    @Override
    protected void doDisconnect() throws Exception {

    }

    @Override
    public ChannelConfig config() {
        return null;
    }

    @Override
    public boolean isActive() {
        return false;
    }

    private static class SerialNioUnsafe implements NioUnsafe {
        @Override
        public SelectableChannel ch() {
            throw new NotImplementedException("");
        }

        @Override
        public void finishConnect() {
            throw new NotImplementedException("");
        }

        @Override
        public void read() {
            throw new NotImplementedException("");
        }

        @Override
        public void forceFlush() {
            throw new NotImplementedException("");
        }

        @Override
        public RecvByteBufAllocator.Handle recvBufAllocHandle() {
            throw new NotImplementedException("");
        }

        @Override
        public SocketAddress localAddress() {
            throw new NotImplementedException("");
        }

        @Override
        public SocketAddress remoteAddress() {
            throw new NotImplementedException("");
        }

        @Override
        public void register(EventLoop eventLoop, ChannelPromise promise) {
            // Register
            if (!(eventLoop instanceof NioEventLoop)) {
                throw new IllegalArgumentException("Only valid for NioEventLoop!");
            }
            if (!(promise.channel() instanceof SerialChannel)) {
                throw new IllegalArgumentException("Only valid for " + SerialChannel.class + " but is " + promise.channel().getClass());
            }
            // Register channel to event loop
            // We have to use reflection here, I fear
            try {
                Method method = NioEventLoop.class.getDeclaredMethod("unwrappedSelector");
                method.setAccessible(true);
                SerialPollingSelector selector = (SerialPollingSelector) method.invoke(eventLoop);

                selector.register((AbstractSelectableChannel) promise.channel(), 0, null);
            } catch (NoSuchMethodException | IllegalAccessException | InvocationTargetException e) {
                e.printStackTrace();
                throw new NotImplementedException("Should register channel to event loop!!!");
            }
        }

        @Override
        public void bind(SocketAddress localAddress, ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void connect(SocketAddress remoteAddress, SocketAddress localAddress, ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void disconnect(ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void close(ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void closeForcibly() {
            throw new NotImplementedException("");
        }

        @Override
        public void deregister(ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void beginRead() {
            throw new NotImplementedException("");
        }

        @Override
        public void write(Object msg, ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void flush() {
            throw new NotImplementedException("");
        }

        @Override
        public ChannelPromise voidPromise() {
            throw new NotImplementedException("");
        }

        @Override
        public ChannelOutboundBuffer outboundBuffer() {
            throw new NotImplementedException("");
        }
    }
}
