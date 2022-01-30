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

import org.apache.commons.lang3.NotImplementedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.Socket;
import java.net.SocketAddress;
import java.net.SocketOption;
import java.nio.ByteBuffer;
import java.nio.channels.SocketChannel;
import java.nio.channels.spi.SelectorProvider;
import java.util.Set;

class SerialSocketChannel extends SocketChannel {

    private static final Logger logger = LoggerFactory.getLogger(SerialSocketChannel.class);

    private SerialChannel child =  null;

    /**
     * Initializes a new instance of this class.
     *
     * @param provider The provider that created this channel
     */
    protected SerialSocketChannel(SelectorProvider provider) {
        super(provider);
    }

    @Override
    public SocketChannel bind(SocketAddress local) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public <T> SocketChannel setOption(SocketOption<T> name, T value) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public SocketChannel shutdownInput() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public SocketChannel shutdownOutput() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public Socket socket() {
        throw new NotImplementedException("");
    }

    @Override
    public boolean isConnected() {
        throw new NotImplementedException("");
    }

    @Override
    public boolean isConnectionPending() {
        throw new NotImplementedException("");
    }

    @Override
    public boolean connect(SocketAddress remote) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public boolean finishConnect() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public SocketAddress getRemoteAddress() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public int read(ByteBuffer dst) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public long read(ByteBuffer[] dsts, int offset, int length) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public int write(ByteBuffer src) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public long write(ByteBuffer[] srcs, int offset, int length) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public SocketAddress getLocalAddress() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public <T> T getOption(SocketOption<T> name) throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    public Set<SocketOption<?>> supportedOptions() {
        throw new NotImplementedException("");
    }

    @Override
    protected void implCloseSelectableChannel() throws IOException {
        throw new NotImplementedException("");
    }

    @Override
    protected void implConfigureBlocking(boolean block) throws IOException {
        logger.debug("Requesting Blocking mode to '{}'", block ? "blocking" : "non blocking");
    }

    public SerialChannel getChild() {
        return child;
    }

    public void setChild(SerialChannel child) {
        this.child = child;
    }
}
