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
package org.apache.plc4x.java.profinet;

import io.netty.buffer.ByteBuf;
import org.apache.commons.net.util.SubnetUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.profinet.channel.ProfinetChannel;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.config.ProfinetRawSocketTransportConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.discovery.ProfinetDiscoverer;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.connection.SingleProtocolStackConfigurer;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;
import org.pcap4j.core.*;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.*;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.function.ToIntFunction;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetDriver extends GeneratedDriverBase<Ethernet_Frame> {

    private final Logger logger = LoggerFactory.getLogger(ProfinetDriver.class);

    public static final Pattern MAC_ADDRESS = Pattern.compile(
        "^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})?");

    public static final String DRIVER_CODE = "profinet";

    @Override
    public String getProtocolCode() {
        return DRIVER_CODE;
    }

    @Override
    public String getProtocolName() {
        return "Profinet";
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        // TODO: This should actually happen in the execute method of the discoveryRequest and not here ...
        try {
            ProfinetChannel channel = new ProfinetChannel(Pcaps.findAllDevs());
            ProfinetDiscoverer discoverer = new ProfinetDiscoverer(channel);
            return new DefaultPlcDiscoveryRequest.Builder(discoverer);
        } catch (PcapNativeException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    protected Class<? extends PlcConnectionConfiguration> getConfigurationClass() {
        return ProfinetConfiguration.class;
    }

    @Override
    protected Optional<Class<? extends PlcTransportConfiguration>> getTransportConfigurationClass(String transportCode) {
        switch (transportCode) {
            case "raw":
                return Optional.of(ProfinetRawSocketTransportConfiguration.class);
        }
        return Optional.empty();
    }

    @Override
    protected Optional<String> getDefaultTransportCode() {
        return Optional.of("raw");
    }

    @Override
    protected List<String> getSupportedTransportCodes() {
        return Collections.singletonList("raw");
    }

    @Override
    protected boolean awaitSetupComplete() {
        return true;
    }

    /**
     * This protocol doesn't have a disconnect procedure, so there is no need to wait for a login to finish.
     *
     * @return false
     */
    @Override
    protected boolean awaitDisconnectComplete() {
        return false;
    }

    @Override
    protected boolean canRead() {
        return false;
    }

    @Override
    protected boolean canWrite() {
        return false;
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected boolean canBrowse() {
        return true;
    }

    @Override
    protected boolean canDiscover() {
        return true;
    }

    @Override
    protected BaseOptimizer getOptimizer() {
        return new SingleTagOptimizer();
    }

    @Override
    protected ProtocolStackConfigurer<Ethernet_Frame> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(Ethernet_Frame.class, Ethernet_Frame::staticParse)
            .withProtocol(ProfinetProtocolLogic.class)
            .withDriverContext(ProfinetDriverContext.class)
            .build();
    }

    /**
     * Estimate the Length of a Packet
     */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 6) {
                return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 4) + 6;
            }
            return -1;
        }

    }

    @Override
    public ProfinetTag prepareTag(String query) {
        return ProfinetTag.of(query);
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication) throws PlcConnectionException {
        // Check if this is a connection string with a MAC address and "assign-ip" in the options.
        Matcher matcher = URI_PATTERN.matcher(connectionString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-config}(?{parameter-string)?'");
        }
        final String protocolCode = matcher.group("protocolCode");
        String transportCodeMatch = matcher.group("transportCode");
        final String transportCode = (transportCodeMatch != null) ? transportCodeMatch : getDefaultTransportCode().get();
        final String transportConfig = matcher.group("transportConfig");
        final String paramString = matcher.group("paramString");
        Matcher macMatcher = MAC_ADDRESS.matcher(transportConfig);
        if (macMatcher.matches()) {
            logger.info("Setting remote PROFINET device IP using DCP");
            ConfigurationFactory configurationFactory = new ConfigurationFactory();
            ProfinetConfiguration configuration = (ProfinetConfiguration) configurationFactory
                .createConfiguration(getConfigurationClass(), protocolCode, transportCode, transportConfig, paramString);
            if (configuration == null) {
                throw new PlcConnectionException("Unsupported configuration");
            }

            // Check if the paramString, contains an "ip-address" option, as this is required for connecting.
            if(configuration.ipAddress == null) {
                throw new PlcConnectionException("When using mac-address connection string, the parameter 'ip-address' is required");
            }

            // Find which network device could communicate with a device on the given IP address.
            try {
                Inet4Address inet4Address = (Inet4Address) Inet4Address.getByName(configuration.ipAddress);
                deviceLoop: for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                    // We're only interested in real running network interfaces, skip the rest.
                    if (dev.isLoopBack()) {
                        continue;
                    }

                    for (PcapAddress curAddress : dev.getAddresses()) {
                        if((curAddress.getAddress() == null) || (curAddress.getNetmask() == null)) {
                            continue;
                        }
                        if(!(curAddress instanceof PcapIpV4Address)) {
                            continue;
                        }
                        final SubnetUtils.SubnetInfo subnetInfo = new SubnetUtils(curAddress.getAddress().getHostAddress(), curAddress.getNetmask().getHostAddress()).getInfo();
                        if(subnetInfo.isInRange(inet4Address.getHostAddress())) {
                            ProfinetRawSocketTransportConfiguration profinetRawSocketTransportConfiguration = new ProfinetRawSocketTransportConfiguration();
                            InetSocketAddress remoteAddress = new InetSocketAddress(configuration.ipAddress, profinetRawSocketTransportConfiguration.getDefaultPort());
                            MacAddress remoteMacAddress = MacAddress.getByName(transportConfig);
                            MacAddress localMacAddress = dev.getLinkLayerAddresses().stream()
                                .filter(linkLayerAddress -> linkLayerAddress instanceof MacAddress).map(linkLayerAddress -> (MacAddress) linkLayerAddress)
                                .findFirst().orElse(null);

                            //  TODO: Use DCP to assign an IP address to the device (Send one packet to every network device)
                            Ethernet_Frame frame = new Ethernet_Frame(
                                new org.apache.plc4x.java.profinet.readwrite.MacAddress(remoteMacAddress.getAddress()),
                                new org.apache.plc4x.java.profinet.readwrite.MacAddress(localMacAddress.getAddress()),
                                new Ethernet_FramePayload_VirtualLan(
                                    VirtualLanPriority.INTERNETWORK_CONTROL,
                                    false,
                                    (short) 0,
                                    new Ethernet_FramePayload_PnDcp(
                                        new PcDcp_GetSet_Pdu(
                                            PnDcp_FrameId.DCP_GetSet_PDU.getValue(),
                                            false,
                                            false,
                                            0x10000001L,
                                            Collections.singletonList(
                                                new PnDcp_Block_IpParameter(
                                                    false,
                                                    false,
                                                    true,
                                                    remoteAddress.getAddress().getAddress(),
                                                    new byte[] {(byte) 255, (byte) 255, (byte) 255, (byte) 0},
                                                    new byte[] {(byte) 0, (byte) 0, (byte) 0, (byte) 0}
                                                )
                                            )
                                        )
                                    )
                                )
                            );
                            // Open a raw socket and send out the packet to the device.
                            try (PcapHandle handle = dev.openLive(65536, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, 10)) {
                                handle.setFilter(
                                    // Profinet (0x8892) packets received from the remote mac address targeted at our own.
                                    "(ether proto 0x8892) and (ether dst " + Pcaps.toBpfString(localMacAddress) + ") and (ether src " + Pcaps.toBpfString(remoteMacAddress) + ")",
                                    BpfProgram.BpfCompileMode.OPTIMIZE);
                                WriteBufferByteBased buffer = new WriteBufferByteBased(frame.getLengthInBytes());
                                try {
                                    frame.serialize(buffer);
                                    Packet packet = EthernetPacket.newPacket(buffer.getBytes(), 0, frame.getLengthInBytes());
                                    handle.sendPacket(packet);
                                } catch (PcapNativeException | NotOpenException | SerializationException |
                                         IllegalRawDataException e) {
                                    throw new RuntimeException(e);
                                }

                                // Now wait for a short while for the device to respond.
                                CompletableFuture<Boolean> future = new CompletableFuture<>();
                                handle.loop(1, (PacketListener) packet -> {
                                    try {
                                        ReadBufferByteBased readBuffer = new ReadBufferByteBased(packet.getRawData());
                                        Ethernet_Frame ethernetFrame = Ethernet_Frame.staticParse(readBuffer);
                                        if(ethernetFrame.getPayload() instanceof Ethernet_FramePayload_PnDcp) {
                                            Ethernet_FramePayload_PnDcp payloadPnDcp = (Ethernet_FramePayload_PnDcp) ethernetFrame.getPayload();
                                            if(payloadPnDcp.getPdu() instanceof PcDcp_GetSet_Pdu) {
                                                // TODO: Possibly have a look if this operation was successful
                                                future.complete(true);
                                                return;
                                            }
                                        }
                                        future.completeExceptionally(new PlcConnectionException("Unexpected response"));
                                    } catch (ParseException e) {
                                        future.completeExceptionally(new PlcConnectionException("Error setting ip address", e));
                                    }
                                });
                                future.get(1000, TimeUnit.MILLISECONDS);

                                logger.info("Finished setting remote PROFINET device IP");
                            } catch (InterruptedException | TimeoutException | ExecutionException e) {
                                throw new PlcConnectionException(e);
                            }
                            break deviceLoop;
                        }
                    }
                }
            } catch (PcapNativeException | UnknownHostException |NotOpenException e) {
                throw new PlcConnectionException(e);
            }

            // Use a re-written connection string using the IP to really connect.
            return super.getConnection(String.format("%s:%s://%s?%s", protocolCode, transportCode, configuration.ipAddress, paramString), authentication);
        }

        return super.getConnection(connectionString, authentication);
    }

}
