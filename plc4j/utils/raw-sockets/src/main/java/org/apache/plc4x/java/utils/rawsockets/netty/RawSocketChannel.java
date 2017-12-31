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
package org.apache.plc4x.java.utils.rawsockets.netty;

import com.savarese.rocksaw.net.RawSocket;
import io.netty.buffer.ByteBuf;
import io.netty.channel.*;
import io.netty.channel.nio.AbstractNioByteChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.util.internal.logging.InternalLogger;
import io.netty.util.internal.logging.InternalLoggerFactory;
import org.apache.plc4x.java.api.exceptions.PlcIoException;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.channels.SelectableChannel;

/**
 * Netty channel implementation that uses RockSaw to create a raw socket connection to implement
 * IP-socket based protocols not based on TCP or UDP.
 *
 * NOTE: This class is currently a WIP (Work in progress) it should only be used with great care.
 */
public class RawSocketChannel extends AbstractNioByteChannel {

    private static final InternalLogger logger = InternalLoggerFactory.getInstance(NioSocketChannel.class);

    // The protocol number is defined in the IP protocol and indicates the type of protocol the payload
    // the IP packet uses. This number is assigned by the IESG. A full list of the registered protocol
    // numbers can be found here: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
    private int protocolNumber;

    private RawSocket socket;
    private InetSocketAddress localAddress;
    private InetSocketAddress remoteAddress;

    /**
     * Initializes a raw socket that is able to communicate with raw IPv4 and IPv6 sockets, hereby
     * allowing to implement protocols below TCP and UDP.
     *
     * For a list of public known protocol numbers see:
     * https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
     *
     * @param parent
     * @param ch
     * @param protocolNumber protocol number identifying the protocol.
     * @throws PlcIoException
     */
    public RawSocketChannel(Channel parent, SelectableChannel ch, int protocolNumber) throws PlcIoException {
        super(parent, ch);

        this.protocolNumber = protocolNumber;

        try {
            socket = new RawSocket();
            socket.setIPHeaderInclude(true);
        } catch (IOException e) {
            throw new PlcIoException("Error setting up raw socket", e);
        }
    }

    /**
     * Opens a connection to the given remote address.
     *
     * @param remoteAddress
     * @param localAddress
     * @return
     * @throws Exception
     */
    @Override
    protected boolean doConnect(SocketAddress remoteAddress, SocketAddress localAddress) throws Exception {
        if(!(remoteAddress instanceof InetSocketAddress) || !(localAddress instanceof InetSocketAddress)) {
            throw new PlcIoException("Both remoteAddress and localAddress must be of type InetSocketAddress");
        }

        try {
            this.localAddress = (InetSocketAddress) localAddress;
            this.remoteAddress = (InetSocketAddress) remoteAddress;

            socket.open(RawSocket.PF_INET, protocolNumber);

            return socket.isOpen();
        } catch (IllegalStateException | IOException e) {
            return false;
        }
    }

    @Override
    protected void doFinishConnect() throws Exception {

    }

    /**
     * Opens a listening socket.
     *
     * @param localAddress
     * @throws Exception
     */
    @Override
    protected void doBind(SocketAddress localAddress) throws Exception {
        if(socket.isOpen()) {
            throw new PlcIoException("Raw socket already opened.");
        }
        if(localAddress instanceof InetSocketAddress) {
            this.localAddress = (InetSocketAddress) localAddress;
            socket.bind(this.localAddress.getAddress());
        } else {
            throw new PlcIoException("Unsupported type of local address. Only InetSocketAddress supported.");
        }
    }

    /**
     * Closes the connection.
     *
     * @throws Exception
     */
    @Override
    protected void doDisconnect() throws Exception {
        if(socket.isOpen()) {
            socket.close();
        }
    }

    @Override
    protected ChannelFuture shutdownInput() {
        return null;
    }

    @Override
    protected int doReadBytes(ByteBuf buf) throws Exception {
        byte[] byteBuf = new byte[1024];
        int readBytes = socket.read(byteBuf);
        buf.writeBytes(byteBuf, 0, readBytes);
        return readBytes;
    }

    @Override
    protected int doWriteBytes(ByteBuf buf) throws Exception {
        byte[] readableBytes = new byte[buf.readableBytes()];
        buf.readBytes(readableBytes);
        socket.write(remoteAddress.getAddress(), readableBytes);
        return readableBytes.length;
    }

    @Override
    protected long doWriteFileRegion(FileRegion region) throws Exception {
        throw new UnsupportedOperationException("doWriteFileRegion not implemented");
    }

    @Override
    protected SocketAddress localAddress0() {
        return localAddress;
    }

    @Override
    protected SocketAddress remoteAddress0() {
        return remoteAddress;
    }

    @Override
    public ChannelConfig config() {
        return null;
    }

    @Override
    public boolean isActive() {
        return socket.isOpen();
    }
}
