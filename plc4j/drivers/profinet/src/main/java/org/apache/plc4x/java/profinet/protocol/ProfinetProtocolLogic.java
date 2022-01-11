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
package org.apache.plc4x.java.profinet.protocol;

import io.netty.channel.Channel;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.readwrite.io.DceRpc_PacketIO;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.generation.*;
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
import java.util.Arrays;
import java.util.Collections;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> {

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static AtomicInteger sessionKeyGenerator = new AtomicInteger(1);

    private final Logger logger = LoggerFactory.getLogger(ProfinetProtocolLogic.class);

    private ProfinetDriverContext profinetDriverContext;

    @Override
    public void setContext(ConversationContext<Ethernet_Frame> context) {
        super.setContext(context);
        this.profinetDriverContext = (ProfinetDriverContext) driverContext;
    }

    @Override
    public void onConnect(ConversationContext<Ethernet_Frame> context) {
        final Channel channel = context.getChannel();
        if(!(channel instanceof RawSocketChannel)) {
            logger.warn("Expected a 'raw' transport, closing channel...");
            context.getChannel().close();
            return;
        }

        RawSocketChannel rawSocketChannel = (RawSocketChannel) channel;

        // Create an udp socket
        DatagramSocket udpSocket;
        try {
            udpSocket = new DatagramSocket();
        } catch (SocketException e) {
            logger.warn("Unable to create udp socket " + e.getMessage());
            context.getChannel().close();
            return;
        }

        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Initialize some important datastructures, that will be used a lot.

        // Generate a new Activity Id, which will be used throughout the connection.
        profinetDriverContext.setDceRpcActivityUuid(generateActivityUuid());

        // TODO: Possibly we can remove the ARP lookup and simply use the mac address in the connection-response.

        // Local connectivity attributes
        profinetDriverContext.setLocalMacAddress(new MacAddress(rawSocketChannel.getLocalMacAddress().getAddress()));
        final InetSocketAddress localAddress = (InetSocketAddress) rawSocketChannel.getLocalAddress();
        Inet4Address localIpAddress = (Inet4Address) localAddress.getAddress();
        profinetDriverContext.setLocalIpAddress(new IpAddress(localIpAddress.getAddress()));
        // Use the port of the udp socket
        profinetDriverContext.setLocalUdpPort(udpSocket.getPort());

        // Remote connectivity attributes
        profinetDriverContext.setRemoteMacAddress(new MacAddress(rawSocketChannel.getRemoteMacAddress().getAddress()));
        final InetSocketAddress remoteAddress = (InetSocketAddress) rawSocketChannel.getRemoteAddress();
        Inet4Address remoteIpAddress = (Inet4Address) remoteAddress.getAddress();
        profinetDriverContext.setRemoteIpAddress(new IpAddress(remoteIpAddress.getAddress()));
        profinetDriverContext.setRemoteUdpPort(remoteAddress.getPort());

        // Generate a new session key.
        profinetDriverContext.setSessionKey(sessionKeyGenerator.getAndIncrement());
        // Reset the session key as soon as it reaches the max for a 16 bit uint
        if(sessionKeyGenerator.get() == 0xFFFF) {
            sessionKeyGenerator.set(1);
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
            final DceRpc_Packet dceRpc_packet = DceRpc_PacketIO.staticParse(readBuffer);
            if(dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE) {
                System.out.println(dceRpc_packet);
            }
        } catch (SerializationException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        } catch (PlcException e) {
            e.printStackTrace();
        } catch (ParseException e) {
            e.printStackTrace();
        }

        System.out.println(rawSocketChannel);
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
                new DceRpc_ObjectUuid(0x0001, 0x0904, 0x002A),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                profinetDriverContext.getDceRpcActivityUuid(),
                0, 0, DceRpc_Operation.CONNECT,
                new PnIoCm_Packet_Req(404, 404, 404,0, 404,
                    Arrays.asList(
                        new PnIoCm_Block_ArReq((short) 1, (short) 0, PnIoCm_ArType.IO_CONTROLLER,
                            new Uuid(Hex.decodeHex("654519352df3b6428f874371217c2b51")),
                            profinetDriverContext.getSessionKey(),
                            profinetDriverContext.getLocalMacAddress(),
                            new Uuid(Hex.decodeHex("dea000006c9711d1827100640008002a")),
                            false, false, false,
                            false, PnIoCm_CompanionArType.SINGLE_AR, false,
                            true, false, PnIoCm_State.ACTIVE,
                            600,
                            // This actually needs to be set to this value and not the real port number.
                            0x8892,
                            // It seems that it must be set to this value, or it won't work.
                            "profinetxadriver4933"),
                        new PnIoCm_Block_IoCrReq((short) 1, (short) 0, PnIoCm_IoCrType.INPUT_CR,
                            0x0001,
                            0x8892,
                            false, false,
                            false, false, PnIoCm_RtClass.RT_CLASS_2, 40,
                            0xBBF0, 128, 8, 1, 0, 0xffffffff,
                            3, 3, 0xC000,
                            new org.apache.plc4x.java.profinet.readwrite.MacAddress(Hex.decodeHex("000000000000")),
                            Collections.singletonList(
                                new PnIoCm_IoCrBlockReqApi(
                                    Arrays.asList(
                                        new PnIoCm_IoDataObject(0, 0x0001, 0),
                                        new PnIoCm_IoDataObject(0, 0x8000, 1),
                                        new PnIoCm_IoDataObject(0, 0x8001, 2),
                                        new PnIoCm_IoDataObject(0, 0x8002, 3),
                                        new PnIoCm_IoDataObject(1, 0x0001, 4)
                                    ),
                                    Collections.singletonList(
                                        new PnIoCm_IoCs(0x0001, 0x0001, 0x0019)
                                    ))
                            )),
                        new PnIoCm_Block_IoCrReq((short) 1, (short) 0, PnIoCm_IoCrType.OUTPUT_CR,
                            0x0002, 0x8892, false, false,
                            false, false, PnIoCm_RtClass.RT_CLASS_2, 40,
                            0x8000, 128, 8, 1, 0, 0xffffffff,
                            3, 3, 0xC000,
                            new MacAddress(Hex.decodeHex("000000000000")),
                            Collections.singletonList(
                                new PnIoCm_IoCrBlockReqApi(
                                    Collections.singletonList(
                                        new PnIoCm_IoDataObject(0x0001, 0x0001, 0x0005)
                                    ),
                                    Arrays.asList(
                                        new PnIoCm_IoCs(0, 0x0001, 0),
                                        new PnIoCm_IoCs(0, 0x8000, 1),
                                        new PnIoCm_IoCs(0, 0x8001, 2),
                                        new PnIoCm_IoCs(0, 0x8002, 3),
                                        new PnIoCm_IoCs(1, 0x0001, 4)
                                    )
                                )
                            )
                        ),
                        new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                            Collections.singletonList(
                                new PnIoCm_ExpectedSubmoduleBlockReqApi(0,
                                    0x00000010, 0x00000000,
                                    Arrays.asList(
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x0001,
                                            0x00000001, false, false,
                                            false, false),
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x8000,
                                            0x00000002, false, false,
                                            false, false),
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x8001,
                                            0x00000003, false, false,
                                            false, false),
                                        new PnIoCm_Submodule_NoInputNoOutputData(0x8002,
                                            0x00000003, false, false,
                                            false, false)
                                    )
                                )
                            )
                        ),
                        new PnIoCm_Block_ExpectedSubmoduleReq((short) 1, (short) 0,
                            Collections.singletonList(
                                new PnIoCm_ExpectedSubmoduleBlockReqApi(1,
                                    0x00000022, 0x00000000, Collections.singletonList(
                                    new PnIoCm_Submodule_InputAndOutputData(0x0001, 0x00000010,
                                        false, false, false,
                                        false, 20, (short) 1, (short) 1,
                                        6, (short) 1, (short) 1))
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

    protected static DceRpc_ActivityUuid generateActivityUuid() {
        UUID number = UUID.randomUUID();
        try {
            WriteBufferByteBased wb = new WriteBufferByteBased(128);
            wb.writeLong(64, number.getMostSignificantBits());
            wb.writeLong(64, number.getLeastSignificantBits());

            ReadBuffer rb = new ReadBufferByteBased(wb.getData());
            return new DceRpc_ActivityUuid(rb.readLong(32), rb.readInt(16), rb.readInt(16), rb.readByteArray(8));
        } catch (SerializationException | ParseException e) {
            // Ignore ... this should actually never happen.
        }
        return null;
    }

}
