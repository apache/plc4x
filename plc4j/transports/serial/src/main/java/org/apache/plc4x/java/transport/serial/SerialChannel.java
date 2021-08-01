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
package org.apache.plc4x.java.transport.serial;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.ByteBufAllocator;
import io.netty.channel.*;
import io.netty.channel.nio.AbstractNioByteChannel;
import io.netty.channel.nio.AbstractNioChannel;
import io.netty.channel.nio.NioEventLoop;
import io.netty.channel.socket.DuplexChannel;
import io.netty.util.ReferenceCountUtil;
import org.apache.commons.lang3.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.lang.reflect.Constructor;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.net.SocketAddress;
import java.nio.channels.ClosedChannelException;
import java.nio.channels.SelectableChannel;
import java.util.concurrent.RejectedExecutionException;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-10
 */
public class SerialChannel extends AbstractNioByteChannel implements DuplexChannel {

    private static final Logger logger = LoggerFactory.getLogger(SerialChannel.class);
    private final SerialChannelConfig config;

    private final VoidChannelPromise unsafeVoidPromise = new VoidChannelPromise(this, false);
    private boolean readPending = false; // Did we receive an EOF?
    private SocketAddress remoteAddress;
    private boolean active = false;
    private SerialSelectionKey selectionKey;
    private SerialChannelHandler comPort;
    private final DefaultChannelPipeline pipeline; // Copied from AbstractChannel


    public SerialChannel() {
        this(null, new SerialSocketChannel(new SerialSelectorProvider()));
        ((SerialSocketChannel) javaChannel()).setChild(this);
    }

    /**
     * Create a new instance
     *
     * @param parent the parent {@link Channel} by which this instance was created. May be {@code null}
     * @param ch     the underlying {@link SelectableChannel} on which it operates
     */
    protected SerialChannel(Channel parent, SelectableChannel ch) {
        super(parent, ch);
        config = new SerialChannelConfig(this);
        pipeline = newChannelPipeline();
    }

    @Override
    public NioUnsafe unsafe() {
        return new SerialNioUnsafe();
    }

    @Override
    public boolean isInputShutdown() {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelFuture shutdownInput() {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelFuture shutdownInput(ChannelPromise promise) {
        throw new NotImplementedException("");
    }

    @Override
    public boolean isOutputShutdown() {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelFuture shutdownOutput() {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelFuture shutdownOutput(ChannelPromise promise) {
        throw new NotImplementedException("");
    }

    @Override
    public boolean isShutdown() {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelFuture shutdown() {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelFuture shutdown(ChannelPromise promise) {
        throw new NotImplementedException("");
    }

    @Override
    protected long doWriteFileRegion(FileRegion region) throws Exception {
        throw new NotImplementedException("");
    }

    @Override
    protected int doReadBytes(ByteBuf buf) throws Exception {
        if (!active) {
            return 0;
        }
        logger.debug("Trying to read bytes from wire...");
        int bytesRead = comPort.read(buf);
        logger.debug("Read {} bytes from the wire", bytesRead);
        return bytesRead;
    }

    @Override
    protected int doWriteBytes(ByteBuf buf) throws Exception {
        // Here we really write bytes to Socket!
        if (!active) {
            return 0;
        }
        logger.debug("Trying to write bytes to wire...");
        int bytesWritten = comPort.write(buf);
        logger.debug("Wrote {} bytes to wire!", bytesWritten);
        return bytesWritten;
    }

    @Override
    protected boolean doConnect(SocketAddress remoteAddress, SocketAddress localAddress) throws Exception {
        this.remoteAddress = remoteAddress;
        if (!(remoteAddress instanceof SerialSocketAddress)) {
            throw new IllegalArgumentException("Socket Address has to be of type " + SerialSocketAddress.class);
        }
        logger.debug("Connecting to Socket Address '{}'", ((SerialSocketAddress) remoteAddress).getIdentifier());

        try {
            // A bit hacky but to make a Test Connection start the String with TEST
            if (((SerialSocketAddress) remoteAddress).getIdentifier().startsWith("TEST")) {
                comPort = SerialChannelHandler.DummyHandler.INSTANCE;
            } else {
                comPort = new SerialChannelHandler.SerialPortHandler(remoteAddress, config);
            }
            logger.debug("Using Com Port {}, trying to open port", comPort.getIdentifier());
            if (comPort.open()) {
                logger.debug("Opened port successful to {}", comPort.getIdentifier());
                comPort.registerSelectionKey(selectionKey);

                this.active = true;
                return true;
            } else {
                logger.debug("Unable to open port {}", comPort.getIdentifier());
                return false;
            }
        } catch (Exception e) {
            logger.warn("exception caught", e);
            this.active = false;
            return false;
        }
    }

    @Override
    protected void doClose() throws Exception {
        if (this.comPort != null) {
            this.comPort.close();
        }
    }

    @Override
    protected void doFinishConnect() throws Exception {
        throw new NotImplementedException("");
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
        throw new NotImplementedException("");
    }

    @Override
    protected void doDisconnect() throws Exception {
        throw new NotImplementedException("");
    }

    @Override
    public ChannelConfig config() {
        return this.config;
    }

    @Override
    public boolean isActive() {
        return active;
    }

    private class SerialNioUnsafe implements NioUnsafe {

        private boolean inFlush0; // Copied from AbstractUnsafe
        private Throwable initialCloseCause; // Copied from AbstractUnsafe

        private ChannelOutboundBuffer outboundBuffer;

        private RecvByteBufAllocator.Handle recvHandle;

        public SerialNioUnsafe() {
            try {
                Constructor<ChannelOutboundBuffer> ctor = ChannelOutboundBuffer.class.getDeclaredConstructor(AbstractChannel.class);
                ctor.setAccessible(true);
                this.outboundBuffer = ctor.newInstance(SerialChannel.this);
            } catch (NoSuchMethodException | InstantiationException | IllegalAccessException | InvocationTargetException e) {
                logger.warn("Problem with reflection", e);
                throw new RuntimeException("Problem providing Buffer", e);
            }
        }

        @Override
        public SelectableChannel ch() {
            throw new NotImplementedException("");
        }

        @Override
        public void finishConnect() {
            throw new NotImplementedException("");
        }

        // See NioByteUnsafe#read()
        @Override
        public void read() {
            logger.debug("Reading...");
            // TODO we should read something here, okay?!
            final ChannelConfig config = config();
            final ChannelPipeline pipeline = pipeline();
            final ByteBufAllocator allocator = config.getAllocator();
            final RecvByteBufAllocator.Handle allocHandle = recvBufAllocHandle();
            allocHandle.reset(config);

            boolean close = false;
            try {
                do {
                    ByteBuf byteBuf = allocHandle.allocate(allocator);
                    allocHandle.lastBytesRead(doReadBytes(byteBuf));
                    if (allocHandle.lastBytesRead() <= 0) {
                        // nothing was read. release the buffer.
                        byteBuf.release();
                        byteBuf = null;
                        close = allocHandle.lastBytesRead() < 0;
                        if (close) {
                            // There is nothing left to read as we received an EOF.
                            readPending = false;
                        }
                        break;
                    }

                    allocHandle.incMessagesRead(1);
                    readPending = false;
                    pipeline.fireChannelRead(byteBuf);
                } while (allocHandle.continueReading());

                allocHandle.readComplete();
                pipeline.fireChannelReadComplete();

                if (close) {
                    // TODO
                    //closeOnRead(pipeline);
                }
            } catch (Throwable t) {
                // TODO
                // handleReadException(pipeline, byteBuf, t, close, allocHandle);
                t.printStackTrace();
            } finally {
                // Check if there is a readPending which was not processed yet.
                // This could be for two reasons:
                // * The user called Channel.read() or ChannelHandlerContext.read() in channelRead(...) method
                // * The user called Channel.read() or ChannelHandlerContext.read() in channelReadComplete(...) method
                //
                // See https://github.com/netty/netty/issues/2254
                if (!readPending && !config.isAutoRead()) {
                    // TODO
                }
            }
        }

        @Override
        public void forceFlush() {
            throw new NotImplementedException("");
        }

        @Override
        public RecvByteBufAllocator.Handle recvBufAllocHandle() {
            if (recvHandle == null) {
                recvHandle = config().getRecvByteBufAllocator().newHandle();
            }
            return recvHandle;
        }

        @Override
        public SocketAddress localAddress() {
            throw new NotImplementedException("");
        }

        @Override
        public SocketAddress remoteAddress() {
            return null;
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

                // Register the channel
                selectionKey = (SerialSelectionKey) ((SerialChannel) promise.channel()).javaChannel().register(selector, 0, SerialChannel.this);

                // Set selection key
                final Field selectionKeyField = AbstractNioChannel.class.getDeclaredField("selectionKey");
                selectionKeyField.setAccessible(true);
                selectionKeyField.set(SerialChannel.this, selectionKey);

                // Set event loop (again, via reflection)
                final Field loop = AbstractChannel.class.getDeclaredField("eventLoop");
                loop.setAccessible(true);
                loop.set(SerialChannel.this, eventLoop);

                // Register Pipeline, if necessary
                // Ensure we call handlerAdded(...) before we actually notify the promise. This is needed as the
                // user may already fire events through the pipeline in the ChannelFutureListener.
                if (!(pipeline() instanceof DefaultChannelPipeline)) {
                    throw new IllegalStateException("Pipeline should be of Type " + DefaultChannelPipeline.class);
                }
                // Again reflection, but has to be done in an event loop
                eventLoop().execute(() -> {
                    try {
                        final Method invokeHandlerAddedIfNeeded = DefaultChannelPipeline.class.getDeclaredMethod("invokeHandlerAddedIfNeeded");
                        invokeHandlerAddedIfNeeded.setAccessible(true);

                        invokeHandlerAddedIfNeeded.invoke(pipeline());

                        pipeline().fireChannelRegistered();
                    } catch (IllegalAccessException | InvocationTargetException | NoSuchMethodException e) {
                        logger.warn("Exception caught", e);
                    }
                });

                // Return promise
                promise.setSuccess();
            } catch (NoSuchMethodException | IllegalAccessException | InvocationTargetException | ClosedChannelException | NoSuchFieldException e) {
                logger.warn("Exception caught", e);
                throw new NotImplementedException("Should register channel to event loop!!!");
            }
        }

        @Override
        public void bind(SocketAddress localAddress, ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void connect(SocketAddress remoteAddress, SocketAddress localAddress, ChannelPromise promise) {
            SerialChannel.this.remoteAddress = remoteAddress;
            eventLoop().execute(() -> {
                try {
                    final boolean success = doConnect(remoteAddress, localAddress);
                    if (success) {
                        // Send a message to the pipeline
                        pipeline().fireChannelActive();
                        // Finally, close the promise
                        promise.setSuccess();
                    } else {
                        promise.setFailure(new RuntimeException("Unable to open the com port '" + ((SerialSocketAddress) remoteAddress).getIdentifier() + "'"));
                    }
                } catch (Exception e) {
                    promise.setFailure(e);
                }
            });
        }

        @Override
        public void disconnect(ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public void close(ChannelPromise promise) {
            logger.debug("Closing the Serial Port '{}'", ((SerialSocketAddress) SerialChannel.this.remoteAddress).getIdentifier());
            eventLoop().execute(() -> {
                try {
                    doClose();
                    promise.setSuccess();
                } catch (Exception e) {
                    logger.warn("Unable to close the connection", e);
                    promise.setFailure(e);
                }
            });
        }

        @Override
        public void closeForcibly() {
            //throw new NotImplementedException("");
        }

        @Override
        public void deregister(ChannelPromise promise) {
            throw new NotImplementedException("");
        }

        @Override
        public final void beginRead() {
            assert eventLoop().inEventLoop();

            if (!isActive()) {
                return;
            }

            try {
                doBeginRead();
            } catch (final Exception e) {
                invokeLater(new Runnable() {
                    @Override
                    public void run() {
                        pipeline().fireExceptionCaught(e);
                    }
                });
                close(voidPromise());
            }
        }

        private void invokeLater(Runnable task) {
            try {
                // This method is used by outbound operation implementations to trigger an inbound event later.
                // They do not trigger an inbound event immediately because an outbound operation might have been
                // triggered by another inbound event handler method.  If fired immediately, the call stack
                // will look like this for example:
                //
                //   handlerA.inboundBufferUpdated() - (1) an inbound handler method closes a connection.
                //   -> handlerA.ctx.close()
                //      -> channel.unsafe.close()
                //         -> handlerA.channelInactive() - (2) another inbound handler method called while in (1) yet
                //
                // which means the execution of two inbound handler methods of the same handler overlap undesirably.
                eventLoop().execute(task);
            } catch (RejectedExecutionException e) {
                logger.warn("Can't invoke task later as EventLoop rejected it", e);
            }
        }

        /**
         * Copied and adapted from {@link io.netty.channel.nio.AbstractNioChannel.AbstractNioUnsafe}
         */
        @Override
        public void write(Object msg, ChannelPromise promise) {
            assert eventLoop().inEventLoop();

            ChannelOutboundBuffer outboundBuffer = this.outboundBuffer;
            if (outboundBuffer == null) {
                // If the outboundBuffer is null we know the channel was closed and so
                // need to fail the future right away. If it is not null the handling of the rest
                // will be done in flush0()
                // See https://github.com/netty/netty/issues/2362
                close(voidPromise());
                ReferenceCountUtil.release(msg);
                throw new RuntimeException("Unable to write", initialCloseCause);
            }

            int size;
            try {
                msg = filterOutboundMessage(msg);
                // Reflection due to privacy
                Method estimatorHandle = DefaultChannelPipeline.class.getDeclaredMethod("estimatorHandle");
                estimatorHandle.setAccessible(true);
                MessageSizeEstimator.Handle handle = (MessageSizeEstimator.Handle) estimatorHandle.invoke(pipeline);
                // end of reflection
                size = handle.size(msg);
                if (size < 0) {
                    size = 0;
                }
            } catch (Throwable t) {
                close(voidPromise());
                ReferenceCountUtil.release(msg);
                logger.error("Problem during write", t);
                throw new RuntimeException("Problem during write", t);
            }

            outboundBuffer.addMessage(msg, size, promise);
        }

        //
        // Start Copied from AbstractUnsafe
        //

        @Override
        public final void flush() {
            assert eventLoop().inEventLoop();

            ChannelOutboundBuffer outboundBuffer = this.outboundBuffer;
            if (outboundBuffer == null) {
                return;
            }

            outboundBuffer.addFlush();
            flush0();
        }

        @SuppressWarnings("deprecation")
        protected void flush0() {
            if (inFlush0) {
                // Avoid re-entrance
                return;
            }

            final ChannelOutboundBuffer outboundBuffer = this.outboundBuffer;
            if (outboundBuffer == null || outboundBuffer.isEmpty()) {
                return;
            }

            inFlush0 = true;

            // Mark all pending write requests as failure if the channel is inactive.
            if (!isActive()) {
                try {
                    if (isOpen()) {
                        callFailFlushed(true);
                    } else {
                        // Do not trigger channelWritabilityChanged because the channel is closed already.
                        callFailFlushed(false);
                    }
                } finally {
                    inFlush0 = false;
                }
                return;
            }

            try {
                doWrite(outboundBuffer);
            } catch (Throwable t) {
                if (t instanceof IOException && config().isAutoClose()) {
                    /**
                     * Just call {@link #close(ChannelPromise, Throwable, boolean)} here which will take care of
                     * failing all flushed messages and also ensure the actual close of the underlying transport
                     * will happen before the promises are notified.
                     *
                     * This is needed as otherwise {@link #isActive()} , {@link #isOpen()} and {@link #isWritable()}
                     * may still return {@code true} even if the channel should be closed as result of the exception.
                     */
                    initialCloseCause = t;
                    close(voidPromise());
                    throw new RuntimeException("Unable to flush", t);
                } else {
                    try {
                        shutdownOutput(voidPromise());
                        throw new RuntimeException("Unable to flush", t);
                    } catch (Throwable t2) {
                        initialCloseCause = t;
                        close(voidPromise());
                        throw new RuntimeException("Unable to flush", t);
                    }
                }
            } finally {
                inFlush0 = false;
            }
        }

        private void callFailFlushed(boolean notify) {
            try {
                Method failFlushed = ChannelOutboundBuffer.class.getDeclaredMethod("failFlushed", Throwable.class, boolean.class);
                failFlushed.setAccessible(true);
                failFlushed.invoke(new RuntimeException("Unable to Flush!"), notify);
            } catch (Exception e) {
                throw new IllegalStateException("Unable to call Failed Flushed!");
            }
        }

        //
        // End Copied from AbstractUnsafe
        //

        @Override
        public ChannelPromise voidPromise() {
            return unsafeVoidPromise;
        }

        @Override
        public ChannelOutboundBuffer outboundBuffer() {
            return this.outboundBuffer;
        }
    }
}
