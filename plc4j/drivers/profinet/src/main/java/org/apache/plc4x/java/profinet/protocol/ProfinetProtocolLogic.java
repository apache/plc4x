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
package org.apache.plc4x.java.profinet.protocol;

import io.netty.channel.Channel;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.device.ProfinetDeviceMessageHandler;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.pcap4j.core.PcapAddress;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.*;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> implements HasConfiguration<ProfinetConfiguration>, PlcSubscriber {

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);
    private final Logger logger = LoggerFactory.getLogger(ProfinetProtocolLogic.class);

    private ProfinetDriverContext profinetDriverContext;
    private boolean connected = false;

    private DatagramSocket udpSocket;
    private RawSocketChannel rawSocketChannel;
    private Channel channel;

    private ProfinetDeviceMessageHandler handler = new ProfinetDeviceMessageHandler();

    private ProfinetConfiguration configuration;

    private static final Uuid ARUUID;

    static {
        try {
            ARUUID = new Uuid(Hex.decodeHex("654519352df3b6428f874371217c2b51"));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
        this.handler.setConfiguredDevices(configuration.configuredDevices);
    }

    @Override
    public void setContext(ConversationContext<Ethernet_Frame> context) {
        super.setContext(context);
        this.profinetDriverContext = (ProfinetDriverContext) driverContext;
        for (Map.Entry<MacAddress, ProfinetDevice> device : configuration.configuredDevices.entrySet()) {
            device.getValue().setContext(context);
        }
        try {
            onDeviceDiscovery();
        } catch (InterruptedException e) {
        }
    }

    private void onDeviceDiscovery() throws InterruptedException {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer();
        DefaultPlcDiscoveryRequest request = new DefaultPlcDiscoveryRequest(
            discoverer,
            new LinkedHashMap<>()
        );

        // TODO:- Add handler for un-requested messages
        discoverer.discoverWithHandler(
            request,
            handler
        );
        waitForDeviceDiscovery();
    }

    private void waitForDeviceDiscovery() throws InterruptedException {
        // Once we receive an LLDP and PN-DCP message for each device move on.
        boolean discovered = false;
        int count = 0;
        while (!discovered) {
            discovered = true;
            for (Map.Entry<MacAddress, ProfinetDevice> device : configuration.configuredDevices.entrySet()) {
                if (!device.getValue().hasLldpPdu() || !device.getValue().hasDcpPdu()) {
                    discovered = false;
                }
            }
            if (!discovered) {
                Thread.sleep(3000L);
                count += 1;
            }
            if (count > 5) {
                break;
            }
        }
    }

    @Override
    public void onConnect(ConversationContext<Ethernet_Frame> context) {

        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Initialize some important datastructures, that will be used a lot.
        // TODO: Possibly we can remove the ARP lookup and simply use the mac address in the connection-response.
        // Local connectivity attributes
        profinetDriverContext.setLocalMacAddress(new MacAddress(rawSocketChannel.getLocalMacAddress().getAddress()));
        final InetSocketAddress localAddress = (InetSocketAddress) rawSocketChannel.getLocalAddress();
        Inet4Address localIpAddress = (Inet4Address) localAddress.getAddress();
        profinetDriverContext.setLocalIpAddress(new IpAddress(localIpAddress.getAddress()));
        // Use the port of the udp socket
        profinetDriverContext.setLocalUdpPort(udpSocket.getPort());

        for (Map.Entry<MacAddress, ProfinetDevice> device : configuration.configuredDevices.entrySet()) {
            device.getValue().onConnect();
        }

        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Create the connection request.
        try {
            // Create the packet
            final DceRpc_Packet profinetConnectionRequest = createProfinetConnectionRequest();
            // Serialize it to a byte-payload
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(profinetConnectionRequest.getLengthInBytes());
            profinetConnectionRequest.serialize(writeBuffer);
            // Create a udp packet.
            DatagramPacket connectRequestPacket = new DatagramPacket(writeBuffer.getData(), writeBuffer.getData().length);
            connectRequestPacket.setAddress(remoteAddress.getAddress());
            connectRequestPacket.setPort(remoteAddress.getPort());
            // Send it.

            udpSocket.send(connectRequestPacket);

            // Receive the response.
            byte[] resultBuffer = new byte[profinetConnectionRequest.getLengthInBytes()];
            DatagramPacket connectResponsePacket = new DatagramPacket(resultBuffer, resultBuffer.length);
            udpSocket.receive(connectResponsePacket);
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(resultBuffer);
            final DceRpc_Packet dceRpc_packet = DceRpc_Packet.staticParse(readBuffer);
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.CONNECT) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {
                    // Get the remote MAC address and store it in the context.
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if ((connectResponse.getBlocks().size() > 0) && (connectResponse.getBlocks().get(0) instanceof PnIoCm_Block_ArRes)) {
                        final PnIoCm_Block_ArRes pnIoCm_block_arRes = (PnIoCm_Block_ArRes) connectResponse.getBlocks().get(0);
                        profinetDriverContext.setRemoteMacAddress(pnIoCm_block_arRes.getCmResponderMacAddr());

                        // Update the raw-socket transports filter expression.
                        ((RawSocketChannel) channel).setRemoteMacAddress(org.pcap4j.util.MacAddress.getByAddress(profinetDriverContext.getRemoteMacAddress().getAddress()));
                    } else {
                        throw new PlcException("Unexpected type of first block.");
                    }
                } else {
                    throw new PlcException("Unexpected response");
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                throw new PlcException("Device rejected connection request");
            } else {
                throw new PlcException("Unexpected response");
            }

            // Create the packet
            final DceRpc_Packet profinetAdvancedConnectionWriteRequest = createProfinetAdvancedConnectionWriteRequest();
            // Serialize it to a byte-payload
            writeBuffer = new WriteBufferByteBased(profinetAdvancedConnectionWriteRequest.getLengthInBytes());
            profinetAdvancedConnectionWriteRequest.serialize(writeBuffer);
            // Create a udp packet.
            connectRequestPacket = new DatagramPacket(writeBuffer.getData(), writeBuffer.getData().length);
            connectRequestPacket.setAddress(remoteAddress.getAddress());
            connectRequestPacket.setPort(remoteAddress.getPort());
            // Send it.

            udpSocket.send(connectRequestPacket);

            // Receive the response.
            resultBuffer = new byte[profinetAdvancedConnectionWriteRequest.getLengthInBytes()];
            connectResponsePacket = new DatagramPacket(resultBuffer, resultBuffer.length);
            udpSocket.receive(connectResponsePacket);


            // Create the packet
            final DceRpc_Packet profinetAdvancedConnectionParameterEnd = createProfinetAdvancedConnectionParameterEnd();
            // Serialize it to a byte-payload
            writeBuffer = new WriteBufferByteBased(profinetAdvancedConnectionParameterEnd.getLengthInBytes());
            profinetAdvancedConnectionParameterEnd.serialize(writeBuffer);
            // Create a udp packet.
            connectRequestPacket = new DatagramPacket(writeBuffer.getData(), writeBuffer.getData().length);
            connectRequestPacket.setAddress(remoteAddress.getAddress());
            connectRequestPacket.setPort(remoteAddress.getPort());
            // Send it.

            udpSocket.send(connectRequestPacket);

            // Receive the response.
            resultBuffer = new byte[profinetAdvancedConnectionParameterEnd.getLengthInBytes()];
            connectResponsePacket = new DatagramPacket(resultBuffer, resultBuffer.length);
            udpSocket.receive(connectResponsePacket);

            // Create the packet
            final DceRpc_Packet profinetAdvancedConnectionApplicationReady = createProfinetAdvancedConnectionApplicationReady();
            // Serialize it to a byte-payload
            writeBuffer = new WriteBufferByteBased(profinetAdvancedConnectionApplicationReady.getLengthInBytes());
            profinetAdvancedConnectionApplicationReady.serialize(writeBuffer);
            // Create a udp packet.
            connectRequestPacket = new DatagramPacket(writeBuffer.getData(), writeBuffer.getData().length);
            connectRequestPacket.setAddress(remoteAddress.getAddress());
            connectRequestPacket.setPort(remoteAddress.getPort());
            // Send it.

            udpSocket.send(connectRequestPacket);

            // Receive the response.
            resultBuffer = new byte[profinetAdvancedConnectionApplicationReady.getLengthInBytes()];
            connectResponsePacket = new DatagramPacket(resultBuffer, resultBuffer.length);
            udpSocket.receive(connectResponsePacket);
            context.fireConnected();
            connected = true;

        } catch (SerializationException | IOException | PlcException | ParseException e) {
            logger.error("Error", e);
        }
    }

    @Override
    public void close(ConversationContext<Ethernet_Frame> context) {
        // Nothing to do here ...
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        future.completeExceptionally(new NotImplementedException());
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        future.completeExceptionally(new NotImplementedException());
        return future;
    }

    private Ethernet_FramePayload_PnDcp createProfinetCyclicDataRequest() {
        return new Ethernet_FramePayload_PnDcp(
            new PnDcp_Pdu_RealTimeCyclic(
                0x8000,
                new PnIo_CyclicServiceDataUnit((short) 0,(short) 0, (short) 0),
                16696,
                false,
                false,
                false,
                false,
                false,
                false));
    }


    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        if (!connected) {
            throw new RuntimeException("Not Connected");
        }

        final InetSocketAddress remoteAddress = (InetSocketAddress) rawSocketChannel.getRemoteAddress();

        try {
            // Create the packet
            final Ethernet_FramePayload_PnDcp profinetConnectionRequest = createProfinetCyclicDataRequest();
            // Serialize it to a byte-payload
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(profinetConnectionRequest.getLengthInBytes());
            profinetConnectionRequest.serialize(writeBuffer);
            // Create a udp packet.
            DatagramPacket connectRequestPacket = new DatagramPacket(writeBuffer.getData(), writeBuffer.getData().length);
            connectRequestPacket.setAddress(remoteAddress.getAddress());
            connectRequestPacket.setPort(remoteAddress.getPort());
            // Send it.

            udpSocket.send(connectRequestPacket);

            // Receive the response.
            byte[] resultBuffer = new byte[profinetConnectionRequest.getLengthInBytes()];
            DatagramPacket connectResponsePacket = new DatagramPacket(resultBuffer, resultBuffer.length);
            udpSocket.receive(connectResponsePacket);
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(resultBuffer);
            final DceRpc_Packet dceRpc_packet = DceRpc_Packet.staticParse(readBuffer);
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.CONNECT) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {
                    // Get the remote MAC address and store it in the context.
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if ((connectResponse.getBlocks().size() > 0) && (connectResponse.getBlocks().get(0) instanceof PnIoCm_Block_ArRes)) {
                        final PnIoCm_Block_ArRes pnIoCm_block_arRes = (PnIoCm_Block_ArRes) connectResponse.getBlocks().get(0);
                        profinetDriverContext.setRemoteMacAddress(pnIoCm_block_arRes.getCmResponderMacAddr());

                        // Update the raw-socket transports filter expression.
                        ((RawSocketChannel) channel).setRemoteMacAddress(org.pcap4j.util.MacAddress.getByAddress(profinetDriverContext.getRemoteMacAddress().getAddress()));
                    } else {
                        throw new PlcException("Unexpected type of first block.");
                    }
                } else {
                    throw new PlcException("Unexpected response");
                }
            } else if (dceRpc_packet.getPacketType() == DceRpc_PacketType.REJECT) {
                throw new PlcException("Device rejected connection request");
            } else {
                throw new PlcException("Unexpected response");
            }
        } catch (SerializationException e) {
            throw new RuntimeException(e);
        } catch (ParseException e) {
            throw new RuntimeException(e);
        } catch (PlcException e) {
            throw new RuntimeException(e);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        return future;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        return null;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {

    }

    @Override
    protected void decode(ConversationContext<Ethernet_Frame> context, Ethernet_Frame msg) throws Exception {
        super.decode(context, msg);
    }


    private Optional<PcapNetworkInterface> getNetworkInterfaceForConnection(InetAddress address) {
        try {
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                // We're only interested in real running network interfaces, skip the rest.
                if (dev.isLoopBack() || !dev.isRunning() || dev.isUp()) {
                    continue;
                }

                for (PcapAddress curAddress : dev.getAddresses()) {

                }
            }
        } catch (PcapNativeException e) {
            logger.warn(String.format("Error finding network device for connection to %s", address.toString()), e);
        }
        return Optional.empty();
    }

    private DceRpc_Packet createProfinetConnectionRequest() throws PlcException {
        try {
            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                profinetDriverContext.getDceRpcActivityUuid(),
                0, 0, DceRpc_Operation.CONNECT,
                new PnIoCm_Packet_Req(16696, 16696, 0, 0,
                    Arrays.asList(
                        new PnIoCm_Block_ArReq((short) 1, (short) 0, PnIoCm_ArType.IO_CONTROLLER,
                            new Uuid(Hex.decodeHex("654519352df3b6428f874371217c2b51")),
                            profinetDriverContext.getSessionKey(),
                            profinetDriverContext.getLocalMacAddress(),
                            new Uuid(Hex.decodeHex("dea000006c9711d1827100640008002a")),
                            false, true, false,
                            false, PnIoCm_CompanionArType.SINGLE_AR, false,
                            true, false, PnIoCm_State.ACTIVE,
                            600,
                            // This actually needs to be set to this value and not the real port number.
                            0x8892,
                            // It seems that it must be set to this value, or it won't work.
                            "plc4x"),
                        new PnIoCm_Block_IoCrReq((short) 1, (short) 0, PnIoCm_IoCrType.INPUT_CR,
                            0x0001,
                            0x8892,
                            false, false,
                            false, false, PnIoCm_RtClass.RT_CLASS_2, 40,
                            0xBBF0, 128, 8, 1, 0, 0xffffffff,
                            50, 50, 0xC000,
                            new org.apache.plc4x.java.profinet.readwrite.MacAddress(Hex.decodeHex("000000000000")),
                            Collections.singletonList(
                                new PnIoCm_IoCrBlockReqApi(
                                    Arrays.asList(
                                        new PnIoCm_IoDataObject(0, 0x0001, 0),
                                        new PnIoCm_IoDataObject(0, 0x8000, 1),
                                        new PnIoCm_IoDataObject(0, 0x8001, 2)
                                    ),
                                    new ArrayList<PnIoCm_IoCs>(0))
                            )),
                        new PnIoCm_Block_IoCrReq((short) 1, (short) 0, PnIoCm_IoCrType.OUTPUT_CR,
                            0x0002, 0x8892, false, false,
                            false, false, PnIoCm_RtClass.RT_CLASS_2, 40,
                            0xFFFF, 128, 8, 1, 0, 0xffffffff,
                            50, 50, 0xC000,
                            new MacAddress(Hex.decodeHex("000000000000")),
                            Collections.singletonList(
                                new PnIoCm_IoCrBlockReqApi(
                                    new ArrayList<PnIoCm_IoDataObject>(0),
                                    Arrays.asList(
                                        new PnIoCm_IoCs(0, 0x0001, 0),
                                        new PnIoCm_IoCs(0, 0x8000, 1),
                                        new PnIoCm_IoCs(0, 0x8001, 2)
                                    )
                                )
                            )
                        ),
                        new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                            Collections.singletonList(
                                new PnIoCm_ExpectedSubmoduleBlockReqApi(0,
                                    0x00000001, 0x00000000,
                                    Arrays.asList(
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x0001,
                                            0x00000001, false, false,
                                            false, false),
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x8000,
                                            0x00008000, false, false,
                                            false, false),
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x8001,
                                            0x00008001, false, false,
                                            false, false)
                                    )
                                )
                            )
                        ),
                        new PnIoCm_Block_AlarmCrReq((short) 1, (short) 0,
                            PnIoCm_AlarmCrType.ALARM_CR, 0x8892, false, false, 1, 3,
                            0x0000, 200, 0xC000, 0xA000)
                    ))
            );

            /*// Build the UDP/IP/EthernetFrame to transport the package.
            return new Ethernet_Frame(profinetDriverContext.getRemoteMacAddress(), profinetDriverContext.getLocalMacAddress(),
                new Ethernet_FramePayload_IPv4(ThreadLocalRandom.current().nextInt(0, Integer.MAX_VALUE), (short) 64,
                    profinetDriverContext.getLocalIpAddress(), profinetDriverContext.getRemoteIpAddress(),
                    new Udp_Packet(profinetDriverContext.getLocalUdpPort(), profinetDriverContext.getRemoteUdpPort(),
                        dceRpcConnectionRequest)));*/
        } catch (DecoderException e) {
            throw new PlcException("Error creating connection request", e);
        }
    }

    private DceRpc_Packet createProfinetAdvancedConnectionWriteRequest() throws PlcException {

        return new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            profinetDriverContext.getDceRpcActivityUuid(),
            0, 1, DceRpc_Operation.WRITE,
            new PnIoCm_Packet_Req(16696, 16696, 0, 244,
                Arrays.asList(
                    new IODWriteRequestHeader(
                        (short) 1,
                        (short) 0,
                        0,
                        ARUUID,
                        0x00000000,
                        0x0000,
                        0x0000,
                        0xe040,
                        180
                        ),
                    new IODWriteRequestHeader(
                        (short) 1,
                        (short) 0,
                        1,
                        ARUUID,
                        0x00000000,
                        0x0000,
                        0x8000,
                        0x8071,
                        12
                    ),
                    new PDInterfaceAdjust(
                        (short) 1,
                        (short) 0,
                        MultipleInterfaceModeNameOfDevice.NAME_PROVIDED_BY_LLDP
                    )
                ))
        );
    }

    private DceRpc_Packet createProfinetAdvancedConnectionParameterEnd() throws PlcException {

        return new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            profinetDriverContext.getDceRpcActivityUuid(),
            0, 1, DceRpc_Operation.CONTROL,
            new PnIoCm_Packet_Req(16696, 16696, 0, 244,
                Arrays.asList(
                    new PnIoCm_Control_Request(
                        (short) 1,
                        (short) 0,
                        ARUUID,
                        0x0001,
                        0x0001
                    )
                ))
        );
    }



    private DceRpc_Packet createProfinetAdvancedConnectionApplicationReady() throws PlcException {

        return new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            profinetDriverContext.getDceRpcActivityUuid(),
            0, 1, DceRpc_Operation.CONTROL,
            new PnIoCm_Packet_Req(16696, 16696, 0, 244,
                Arrays.asList(
                    new PnIoCM_Block_Request(
                        (short) 1,
                        (short) 0,
                        ARUUID,
                        0x0001,
                        0x0002,
                        0x0000
                    )
                ))
        );
    }



}
