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

import io.netty.channel.*;
import io.netty.channel.oio.OioByteStreamChannel;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.utils.rawsockets.RawIpSocket;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;

import java.net.SocketAddress;

public class RawSocketChannel extends OioByteStreamChannel {

    private static final ChannelMetadata METADATA = new ChannelMetadata(false, 16);

    PcapHandle handle;

    protected class RawByteUnsafe extends AbstractChannel.AbstractUnsafe {
        @Override
        public void connect(SocketAddress remoteAddress, SocketAddress localAddress, ChannelPromise promise) {
            // Connect?!
            try {
                doConnect(remoteAddress, localAddress);
                pipeline().fireChannelActive();
                promise.setSuccess();
            } catch (Exception e) {
                promise.setFailure(e);
            }
        }

    }

    public RawSocketChannel() {
        super(null);
    }

    @Override
    protected AbstractUnsafe newUnsafe() {
        return new RawByteUnsafe();
    }

    @Override
    protected boolean isCompatible(EventLoop loop) {
        return true;
    }

    @Override
    protected void doConnect(SocketAddress remoteAddress, SocketAddress localAddress) throws Exception {
        PcapNetworkInterface nif = Pcaps.findAllDevs().get(0);
        this.handle = nif.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
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
        System.out.println(localAddress);
    }

    @Override
    protected void doDisconnect() throws Exception {
        throw new NotImplementedException("");
    }

    @Override
    protected void doClose() throws Exception {
        if (this.handle != null) {
            this.handle.close();
        }
    }

    @Override
    protected void doBeginRead() throws Exception {
        this.handle.getNextRawPacketEx();
    }

    @Override
    protected void doWrite(ChannelOutboundBuffer in) throws Exception {
        System.out.println(in);
    }

    @Override
    public ChannelConfig config() {
        return new RawSocketChannelConfig(this);
    }

    @Override
    public boolean isOpen() {
        return true;
    }

    @Override
    public boolean isActive() {
        return true;
    }

    @Override
    public ChannelMetadata metadata() {
        return METADATA;
    }

    @Override
    protected boolean isInputShutdown() {
        return false;
    }

    @Override
    protected ChannelFuture shutdownInput() {
        return null;
    }

}
