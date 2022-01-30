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
import org.apache.commons.net.util.SubnetUtils;
import org.apache.plc4x.java.utils.pcap.netty.exception.PcapException;
import org.apache.plc4x.java.utils.rawsockets.netty.address.RawSocketPassiveAddress;
import org.apache.plc4x.java.utils.rawsockets.netty.config.RawSocketChannelConfig;
import org.apache.plc4x.java.utils.rawsockets.netty.utils.ArpLookup;
import org.pcap4j.core.*;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.*;

/**
 * Channel based on a Pcap4J raw-socket, allowing to actively communicate based on ethernet frames.
 */
public class RawSocketChannel extends OioByteStreamChannel {

    private static final Logger logger = LoggerFactory.getLogger(RawSocketChannel.class);

    private final RawSocketChannelConfig config;

    private MacAddress remoteMacAddress;
    private SocketAddress remoteAddress;
    private MacAddress localMacAddress;
    private SocketAddress localAddress;
    private PcapHandle receiveHandle;
    private Thread loopThread;

    public RawSocketChannel() {
        super(null);
        config = new RawSocketChannelConfig(this);
    }

    public void setRemoteMacAddress(MacAddress remoteMacAddress) {
        this.remoteMacAddress = remoteMacAddress;

        // Update the filter expression, if the handle is already open.
        if((receiveHandle != null) && (receiveHandle.isOpen())){
            MacAddress tempRemoteMacAddress = remoteMacAddress != null ? remoteMacAddress : MacAddress.getByAddress(new byte[]{0,0,0,0,0,0});
            String filter = config.getMacBasedFilterString(localMacAddress, tempRemoteMacAddress);
            if(filter.length() > 0) {
                try {
                    receiveHandle.setFilter(filter, BpfProgram.BpfCompileMode.OPTIMIZE);
                } catch (NotOpenException | PcapNativeException e) {
                    throw new RuntimeException("Error updating filter expression");
                }
            }
        }
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
        if(!((remoteAddress instanceof RawSocketPassiveAddress) || (remoteAddress instanceof InetSocketAddress))) {
            logger.error("Expecting remote address of type RawSocketPassiveAddress or InetSocketAddress");
            pipeline().fireExceptionCaught(
                new PcapException("Expecting remote address of type RawSocketPassiveAddress or InetSocketAddress"));
            return;
        }
        this.localAddress = localAddress;
        this.remoteAddress = remoteAddress;

        PcapNetworkInterface nif = null;
        // Try to get the device name of the network interface that we want to open.
        if(remoteAddress instanceof RawSocketPassiveAddress) {
            RawSocketPassiveAddress rawSocketPassiveAddress = (RawSocketPassiveAddress) remoteAddress;
            String deviceName = getDeviceName(rawSocketPassiveAddress);
            if(deviceName == null) {
                logger.error("Network device not specified and couldn't detect it automatically");
                pipeline().fireExceptionCaught(
                    new PcapException("Network device not specified and couldn't detect it automatically"));
                return;
            }

            // Get a handle to the network-device and open it.
            nif = Pcaps.getDevByName(deviceName);
        } else {
            InetSocketAddress inetSocketAddress = (InetSocketAddress) remoteAddress;
            InetAddress address = inetSocketAddress.getAddress();
            deviceLoop:
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                // We're only interested in real running network interfaces, skip the rest.
                if (dev.isLoopBack()) {
                    continue;
                }

                // Go through all configured addresses and subnets and if we find a matching
                // device, use that as nif.
                for (PcapAddress curAddress : dev.getAddresses()) {
                    if((curAddress.getAddress() == null) || (curAddress.getNetmask() == null)) {
                        continue;
                    }
                    if((address instanceof Inet4Address) && !(curAddress instanceof PcapIpV4Address)) {
                        continue;
                    }
                    if((address instanceof Inet6Address) && !(curAddress instanceof PcapIpV6Address)) {
                        continue;
                    }
                    final SubnetUtils.SubnetInfo subnetInfo = new SubnetUtils(curAddress.getAddress().getHostAddress(), curAddress.getNetmask().getHostAddress()).getInfo();
                    if(subnetInfo.isInRange(address.getHostAddress())) {
                        nif = dev;
                        // Update the local address
                        this.localAddress = new InetSocketAddress(curAddress.getAddress(), 0);
                        this.localMacAddress = dev.getLinkLayerAddresses().stream()
                            .filter(linkLayerAddress -> linkLayerAddress instanceof MacAddress).map(linkLayerAddress -> (MacAddress) linkLayerAddress)
                            .findFirst().orElse(null);
                        break deviceLoop;
                    }
                }
            }
        }

        if (nif == null) {
            logger.error(String.format("Couldn't find network device for %s", remoteAddress));
            pipeline().fireExceptionCaught(
                new PcapException(String.format("Couldn't find network device for %s", remoteAddress)));
            return;
        }

        // If desired: Try to get the mac address of the remote station by sending an ARP request.
        if(config.isResolveMacAddress()) {
            if ((this.remoteMacAddress == null) && (this.remoteAddress instanceof InetSocketAddress)) {
                this.remoteMacAddress = ArpLookup.resolveMacAddress(nif, (InetSocketAddress) this.remoteAddress,
                    (InetSocketAddress) this.localAddress, localMacAddress).orElse(null);
            }
        }

        receiveHandle = nif.openLive(65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
        if(logger.isDebugEnabled()) {
            logger.debug(String.format("Listening on device %s", nif.getName()));
        }

        // If the address allows fine-tuning which packets to process, set a filter to reduce the load.
        // If no remote mac address is set, we'll set it to 00:00:00:00:00, which is an invalid address
        // and therefore no traffic will be incoming.
        // TODO: Make configurable, if we should accept from any source or not.
        MacAddress tempRemoteMacAddress = remoteMacAddress != null ? remoteMacAddress : MacAddress.getByAddress(new byte[]{0,0,0,0,0,0});
        String filter = config.getMacBasedFilterString(localMacAddress, tempRemoteMacAddress);
        if(filter.length() > 0) {
            receiveHandle.setFilter(filter, BpfProgram.BpfCompileMode.OPTIMIZE);
        }

        // Create a buffer where the raw socket worker can send data to.
        ByteBuf buffer = Unpooled.buffer();

        // Start a thread that processes the callbacks from the raw socket and simply
        // forwards the bytes read to the buffer.
        loopThread = new Thread(() -> {
            try {
                receiveHandle.loop(-1,
                    (PacketListener) packet -> buffer.writeBytes(config.getPacketHandler().getData(packet)));
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

        if(remoteAddress instanceof RawSocketPassiveAddress) {
            // Right now we're using an output stream that simply discards everything.
            // This is ok while implementing passive drivers for protocols, however as
            // soon as we start implementing ethernet layer protocols, we'll have to also
            // be able to actually send data. The PcapInputStream simply acts as a
            // breaking point if no packets are coming in and the read operation would
            // simply block indefinitely.
            activate(new PcapInputStream(buffer), new DiscardingOutputStream());
        } else {
            PcapHandle sendHandle = nif.openLive(65535, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10);
            activate(new PcapInputStream(buffer), new PcapOutputStream(sendHandle));
        }
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
    protected void doBind(SocketAddress localAddress) {
        throw new UnsupportedOperationException("Not implemented");
    }

    @Override
    protected void doDisconnect() {
        this.loopThread.interrupt();
        if (this.receiveHandle != null) {
            this.receiveHandle.close();
        }
    }

    @Override
    protected int doReadBytes(ByteBuf buf) throws Exception {
        if (receiveHandle == null || !receiveHandle.isOpen()) {
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

    public SocketAddress getRemoteAddress() {
        return remoteAddress;
    }

    public MacAddress getRemoteMacAddress() {
        return remoteMacAddress;
    }

    public SocketAddress getLocalAddress() {
        return localAddress;
    }

    public MacAddress getLocalMacAddress() {
        return localMacAddress;
    }

    private String getDeviceName(RawSocketPassiveAddress rawSocketAddress) {
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

    private static class PcapOutputStream extends OutputStream {

        private final PcapHandle sendHandle;

        public PcapOutputStream(PcapHandle sendHandle) {
            this.sendHandle = sendHandle;
        }

        @Override
        public void write(int b) throws IOException {
            // This should actually not be called, as in contrast to TCP, Raw sockets are not a stream
            // messages which should be sent should be processed in total by the write(byte[]) method.
            throw new RuntimeException("write(byte) should never be called in a RawSocketChannel");
        }

        @Override
        public void write(byte[] packetBytes, int offset, int len) throws IOException {
            if((offset < 0) || (len < 0) || (offset + len > packetBytes.length)) {
                throw new IndexOutOfBoundsException();
            }
            try {
                // Create a new EthernetPacket with the raw content of the packet bytes.
                Packet rawPacket = EthernetPacket.newPacket(packetBytes, offset, len);
                // Send the packet.
                sendHandle.sendPacket(rawPacket);
            } catch (IllegalRawDataException | NotOpenException | PcapNativeException e) {
                throw new IOException("Error sending packet", e);
            }
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
