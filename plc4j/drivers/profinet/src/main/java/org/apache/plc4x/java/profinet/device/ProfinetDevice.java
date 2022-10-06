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

package org.apache.plc4x.java.profinet.device;

import io.netty.channel.Channel;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicInteger;

public class ProfinetDevice {

    private static final int DEFAULT_UDP_PORT = 34964;
    private final Logger logger = LoggerFactory.getLogger(ProfinetDevice.class);
    private final DceRpc_ActivityUuid uuid;

    private DatagramSocket udpSocket;
    private RawSocketChannel rawSocketChannel;
    private Channel channel;
    private final MacAddress macAddress;
    private ConversationContext<Ethernet_Frame> context;
    private ProfinetDeviceState state = ProfinetDeviceState.IDLE;
    private Lldp_Pdu lldpPdu = null;
    private PnDcp_Pdu dcpPdu = null;
    private String ipAddress;
    private String portId;

    private AtomicInteger sessionKeyGenerator = new AtomicInteger(1);

    private static final Uuid ARUUID;
    static {
        try {
            ARUUID = new Uuid(Hex.decodeHex("654519352df3b6428f874371217c2b51"));
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }


    private void closeUDPSocket() {
        // Handle the closing of the connection, might need to send some messages beforehand.
        if (udpSocket != null && !udpSocket.isConnected()) {
            udpSocket.close();
            context.getChannel().close();
        }
    }

    private boolean createUdpSocket() {
        if (state != ProfinetDeviceState.IDLE) {
            closeUDPSocket();
        }
        if (!(channel instanceof RawSocketChannel)) {
            logger.warn("Expected a 'raw' transport, closing channel...");
            closeUDPSocket();
            return false;
        }

        rawSocketChannel = (RawSocketChannel) channel;


        // Create an udp socket
        try {
            udpSocket = new DatagramSocket();
        } catch (SocketException e) {
            logger.warn("Unable to create udp socket " + e.getMessage());
            closeUDPSocket();
            return false;
        }
        return true;
    }

    public boolean onConnect() {
        if (!createUdpSocket()) {
            // Unable to create UDP connection
            return false;
        }

        ProfinetMessageWrapper.sendUdpMessage(
            new CreateConnection(),
            this
        );

        ProfinetMessageWrapper.sendUdpMessage(
            new WriteParameters(),
            this
        );

        ProfinetMessageWrapper.sendUdpMessage(
            new WriteParametersEnd(),
            this
        );



        return false;
    }

    private int generateSessionKey() {
        // Generate a new session key.
        Integer sessionKey = sessionKeyGenerator.getAndIncrement();
        // Reset the session key as soon as it reaches the max for a 16 bit uint
        if (sessionKeyGenerator.get() == 0xFFFF) {
            sessionKeyGenerator.set(1);
        }
        return sessionKey;
    }

    public boolean hasLldpPdu() {
        if (lldpPdu != null) {
            return true;
        }
        return false;
    }

    public boolean hasDcpPdu() {
        if (dcpPdu != null) {
            return true;
        }
        return false;
    }

    public void handle(PlcDiscoveryItem item) {
        logger.debug("Received Discovered item at device");
        if (item.getOptions().containsKey("IpAddress")) {
            this.ipAddress = item.getOptions().get("IpAddress");
        }
        if (item.getOptions().containsKey("PortId")) {
            this.portId = item.getOptions().get("PortId");
        }
    }

    public void setContext(ConversationContext<Ethernet_Frame> context) {
        this.context = context;
        channel = context.getChannel();
    }

    public ProfinetDevice(MacAddress macAddress) {
        this.macAddress = macAddress;
        // Generate a new Activity Id, which will be used throughout the connection.
        this.uuid = generateActivityUuid();
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

    public DatagramSocket getUdpSocket() {
        return this.udpSocket;
    }

    public RawSocketChannel getRawSocket() {
        return this.rawSocketChannel;
    }

    public InetAddress getIpAddress() throws UnknownHostException {
        return InetAddress.getByName(this.ipAddress);
    }

    public int getPort() {
        return DEFAULT_UDP_PORT;
    }

    public class CreateConnection implements ProfinetCallable<DceRpc_Packet> {

        public DceRpc_Packet create() throws PlcException {
            try {
                return new DceRpc_Packet(
                    DceRpc_PacketType.REQUEST, true, false, false,
                    IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                    new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
                    new DceRpc_InterfaceUuid_DeviceInterface(),
                    ProfinetDevice.this.uuid,
                    0, 0, DceRpc_Operation.CONNECT,
                    new PnIoCm_Packet_Req(16696, 16696, 0, 0,
                        Arrays.asList(
                            new PnIoCm_Block_ArReq((short) 1, (short) 0, PnIoCm_ArType.IO_CONTROLLER,
                                new Uuid(Hex.decodeHex("654519352df3b6428f874371217c2b51")),
                                ProfinetDevice.this.generateSessionKey(),
                                ProfinetDevice.this.macAddress,
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

        public void handle(DceRpc_Packet dceRpc_packet) throws PlcException {
            if ((dceRpc_packet.getOperation() == DceRpc_Operation.CONNECT) && (dceRpc_packet.getPacketType() == DceRpc_PacketType.RESPONSE)) {
                if (dceRpc_packet.getPayload().getPacketType() == DceRpc_PacketType.RESPONSE) {

                    // Get the remote MAC address and store it in the context.
                    final PnIoCm_Packet_Res connectResponse = (PnIoCm_Packet_Res) dceRpc_packet.getPayload();
                    if ((connectResponse.getBlocks().size() > 0) && (connectResponse.getBlocks().get(0) instanceof PnIoCm_Block_ArRes)) {
                        final PnIoCm_Block_ArRes pnIoCm_block_arRes = (PnIoCm_Block_ArRes) connectResponse.getBlocks().get(0);

                        // Update the raw-socket transports filter expression.
                        ((RawSocketChannel) channel).setRemoteMacAddress(org.pcap4j.util.MacAddress.getByAddress(macAddress.getAddress()));
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
        }
    }

    public class WriteParameters implements ProfinetCallable<DceRpc_Packet> {
        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                uuid,
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

        @Override
        public void handle(DceRpc_Packet packet) throws PlcException {
            logger.debug("Received a Write Parameter Response");
        }
    }

    public class WriteParametersEnd implements ProfinetCallable<DceRpc_Packet> {
        public DceRpc_Packet create() {
            return new DceRpc_Packet(
                DceRpc_PacketType.REQUEST, true, false, false,
                IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
                new DceRpc_ObjectUuid((byte) 0x00, 0x0001, 0x0904, 0x002A),
                new DceRpc_InterfaceUuid_DeviceInterface(),
                uuid,
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

        @Override
        public void handle(DceRpc_Packet packet) throws PlcException {
            logger.debug("Received a Write Parameter End Response");
        }
    }

    public class CyclicData implements ProfinetCallable<Ethernet_Frame> {
        public Ethernet_Frame create() {
            return new Ethernet_Frame(
                macAddress,
                macAddress,
                new Ethernet_FramePayload_PnDcp(
                new PnDcp_Pdu_RealTimeCyclic(
                    0x8000,
                    new PnIo_CyclicServiceDataUnit((short) 0,(short) 0, (short) 0),
                    16696,
                    false,
                    false,
                    false,
                    false,
                    false,
                    false)));
        }

        @Override
        public void handle(Ethernet_Frame packet) throws PlcException {
            logger.debug("Received a Write Parameter End Response");
        }
    }


}
