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

package org.apache.plc4x.java.profinet.packets;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;

import java.net.InetSocketAddress;
import java.time.Duration;
import java.util.Collections;
import java.util.Random;
import java.util.concurrent.CompletableFuture;

public class PnDcpPacketFactory {

    public static Ethernet_Frame createIdentificationRequest(MacAddress localMacAddress, MacAddress remoteMacAddress) {
        // Construct and send the search request.
        return new Ethernet_Frame(
            remoteMacAddress,
            localMacAddress,
            new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, (short) 0,
                new Ethernet_FramePayload_PnDcp(
                    new PnDcp_Pdu_IdentifyReq(PnDcp_FrameId.DCP_Identify_ReqPDU.getValue(),
                        1,
                        256,
                        Collections.singletonList(
                            new PnDcp_Block_ALLSelector()
                        )))));
    }

    public static CompletableFuture<PnDcp_Pdu_IdentifyRes> sendIdentificationRequest(ConversationContext<Ethernet_Frame> context, MacAddress localMacAddress, MacAddress remoteMacAddress) {
        CompletableFuture<PnDcp_Pdu_IdentifyRes> future = new CompletableFuture<>();
        context.sendRequest(PnDcpPacketFactory.createIdentificationRequest(localMacAddress, remoteMacAddress))
            .name("Expect Identification response")
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .check(ethernetFramePayload -> ethernetFramePayload instanceof Ethernet_FramePayload_PnDcp)
            .only(Ethernet_FramePayload_PnDcp.class)
            .unwrap(Ethernet_FramePayload_PnDcp::getPdu)
            .only(PnDcp_Pdu_IdentifyRes.class)
            .handle(future::complete);
        return future;
    }

    public static Ethernet_Frame createReadIAndM0BlockRequest(RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            driverContext.getActivityUuid(),
            0,
            driverContext.getAndIncrementIdentification(),
            DceRpc_Operation.READ_IMPLICIT,
            (short) 0,
            new PnIoCm_Packet_Req(16696, 16696, 0,
                Collections.singletonList(
                    new IODReadRequestHeader((short) 1, (short) 0, 0,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
                        0, 0, 1, 0xAFF0, 16696,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
                ))
        );

        return createEthernetFrame(pnChannel, driverContext, packet);
    }

    public static CompletableFuture<PnIoCm_Block_IAndM0> sendReadIAndM0BlockRequest(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        // TODO: Handle error responses quickly (If the device doesn't support PN CM, then we can abort a lot quicker.
        CompletableFuture<PnIoCm_Block_IAndM0> future = new CompletableFuture<>();
        context.sendRequest(PnDcpPacketFactory.createReadIAndM0BlockRequest(pnChannel, driverContext))
            .name("Expect ReadIAndM0Block response")
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .only(Ethernet_FramePayload_IPv4.class)
            .unwrap(Ethernet_FramePayload_IPv4::getPayload)
            .unwrap(DceRpc_Packet::getPayload)
            .only(PnIoCm_Packet_Res.class)
            .handle(dceRpcPacketRes -> {
                if (dceRpcPacketRes.getBlocks().size() != 2) {
                    future.completeExceptionally(new PlcRuntimeException("Expected 2 blocks in the response"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(0) instanceof IODReadResponseHeader)) {
                    future.completeExceptionally(new PlcRuntimeException("The first block was expected to be of type IODReadResponseHeader"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(1) instanceof PnIoCm_Block_IAndM0)) {
                    future.completeExceptionally(new PlcRuntimeException("The second block was expected to be of type PnIoCm_Block_IAndM0"));
                    return;
                }
                PnIoCm_Block_IAndM0 iAndM0 = (PnIoCm_Block_IAndM0) dceRpcPacketRes.getBlocks().get(1);
                future.complete(iAndM0);
            });
        return future;
    }

    public static Ethernet_Frame createReadIAndM1BlockRequest(RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            driverContext.getActivityUuid(),
            0,
            driverContext.getAndIncrementIdentification(),
            DceRpc_Operation.READ_IMPLICIT,
            (short) 0,
            new PnIoCm_Packet_Req(16696, 16696, 0,
                Collections.singletonList(
                    new IODReadRequestHeader((short) 1, (short) 0, 0,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
                        0, 0, 1, 0xAFF1, 16696,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
                ))
        );

        return createEthernetFrame(pnChannel, driverContext, packet);
    }

    public static CompletableFuture<PnIoCm_Block_IAndM1> sendReadIAndM1BlockRequest(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        // TODO: Handle error responses quickly (If the device doesn't support PN CM, then we can abort a lot quicker.
        CompletableFuture<PnIoCm_Block_IAndM1> future = new CompletableFuture<>();
        context.sendRequest(PnDcpPacketFactory.createReadIAndM1BlockRequest(pnChannel, driverContext))
            .name("Expect ReadIAndM1Block response")
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .only(Ethernet_FramePayload_IPv4.class)
            .unwrap(Ethernet_FramePayload_IPv4::getPayload)
            .unwrap(DceRpc_Packet::getPayload)
            .only(PnIoCm_Packet_Res.class)
            .handle(dceRpcPacketRes -> {
                if (dceRpcPacketRes.getBlocks().size() != 2) {
                    future.completeExceptionally(new PlcRuntimeException("Expected 2 blocks in the response"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(0) instanceof IODReadResponseHeader)) {
                    future.completeExceptionally(new PlcRuntimeException("The first block was expected to be of type IODReadResponseHeader"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(1) instanceof PnIoCm_Block_IAndM1)) {
                    future.completeExceptionally(new PlcRuntimeException("The second block was expected to be of type PnIoCm_Block_IAndM0"));
                    return;
                }
                PnIoCm_Block_IAndM1 iAndM1 = (PnIoCm_Block_IAndM1) dceRpcPacketRes.getBlocks().get(1);
                future.complete(iAndM1);
            });
        return future;
    }

    public static Ethernet_Frame createReadRealIdentificationDataRequest(RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            driverContext.getActivityUuid(),
            0,
            driverContext.getAndIncrementIdentification(),
            DceRpc_Operation.READ_IMPLICIT,
            (short) 0,
            new PnIoCm_Packet_Req(16696, 16696, 0,
                Collections.singletonList(
                    new IODReadRequestHeader((short) 1, (short) 0, 0,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
                        // * REMARK: It seems that it doesn't matter if we use subSlotNumber 0 or 1, the responses seem to be the same.
                        0, 0, 0, 0xF000, 16696,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
                ))
        );

        return createEthernetFrame(pnChannel, driverContext, packet);
    }

    public static CompletableFuture<PnIoCm_Block_RealIdentificationData> sendRealIdentificationDataRequest(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        CompletableFuture<PnIoCm_Block_RealIdentificationData> future = new CompletableFuture<>();
        context.sendRequest(createReadRealIdentificationDataRequest(pnChannel, driverContext))
            .name("Expect RealIdentificationData response")
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .only(Ethernet_FramePayload_IPv4.class)
            .unwrap(Ethernet_FramePayload_IPv4::getPayload)
            .unwrap(DceRpc_Packet::getPayload)
            .handle(dceRpcPacket -> {
                if (dceRpcPacket instanceof PnIoCm_Packet_Rej) {
                    future.completeExceptionally(new PlcRuntimeException("RealIdentificationData not supported"));
                    return;
                }
                if (!(dceRpcPacket instanceof PnIoCm_Packet_Res)) {
                    future.completeExceptionally(new PlcRuntimeException("Unexpected response type"));
                    return;
                }
                PnIoCm_Packet_Res dceRpcPacketRes = (PnIoCm_Packet_Res) dceRpcPacket;
                if (dceRpcPacketRes.getBlocks().size() != 2) {
                    future.completeExceptionally(new PlcRuntimeException("Expected 2 blocks in the response"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(0) instanceof IODReadResponseHeader)) {
                    future.completeExceptionally(new PlcRuntimeException("The first block was expected to be of type IODReadResponseHeader"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(1) instanceof PnIoCm_Block_RealIdentificationData)) {
                    future.completeExceptionally(new PlcRuntimeException("The second block was expected to be of type PnIoCm_Block_RealIdentificationData"));
                    return;
                }
                PnIoCm_Block_RealIdentificationData realIdentificationData = (PnIoCm_Block_RealIdentificationData) dceRpcPacketRes.getBlocks().get(1);
                future.complete(realIdentificationData);
            });
        return future;
    }

    public static Ethernet_Frame createParameterEndRequest(RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            driverContext.getActivityUuid(),
            0,
            driverContext.getAndIncrementIdentification(),
            DceRpc_Operation.CONTROL,
            (short) 0,
            new PnIoCm_Packet_Req(16696, 16696, 0,
                Collections.singletonList(
                    new PnIoCm_Control_Request_ParameterEnd((short) 1, (short) 0,
                        new Uuid(driverContext.getApplicationRelationUuid().getData()), 1, 1)
                ))
        );
        return createEthernetFrame(pnChannel, driverContext, packet);
    }

    public static CompletableFuture<PnIoCm_Control_Response_ParameterEnd> sendParameterEndRequest(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        CompletableFuture<PnIoCm_Control_Response_ParameterEnd> future = new CompletableFuture<>();
        context.sendRequest(createParameterEndRequest(pnChannel, driverContext))
            .name("Expect ParameterEnd response")
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if (ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .only(Ethernet_FramePayload_IPv4.class)
            .unwrap(Ethernet_FramePayload_IPv4::getPayload)
            .unwrap(DceRpc_Packet::getPayload)
            .check(pnIoCmPacket -> pnIoCmPacket instanceof PnIoCm_Packet_Rej || pnIoCmPacket instanceof PnIoCm_Packet_Res)
            .handle(dceRpcPacket -> {
                if (dceRpcPacket instanceof PnIoCm_Packet_Rej) {
                    future.completeExceptionally(new PlcRuntimeException("ParameterEnd not supported"));
                    return;
                }
                PnIoCm_Packet_Res dceRpcPacketRes = (PnIoCm_Packet_Res) dceRpcPacket;
                if (dceRpcPacketRes.getBlocks().size() != 1) {
                    future.completeExceptionally(new PlcRuntimeException("Expected 1 blocks in the response"));
                    return;
                }
                if (!(dceRpcPacketRes.getBlocks().get(0) instanceof PnIoCm_Control_Response_ParameterEnd)) {
                    future.completeExceptionally(new PlcRuntimeException("The block was expected to be of type PnIoCm_Control_Response_ParameterEnd"));
                    return;
                }
                PnIoCm_Control_Response_ParameterEnd controlResponse = (PnIoCm_Control_Response_ParameterEnd) dceRpcPacketRes.getBlocks().get(0);
                future.complete(controlResponse);
            });
        return future;
    }

    public static Ethernet_Frame createApplicationReadyResponse(RawSocketChannel pnChannel, ProfinetDriverContext driverContext, int sourcePort, DceRpc_ActivityUuid activityUuid, Uuid arUuid, int sessionKey) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.RESPONSE, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
            new DceRpc_InterfaceUuid_ControllerInterface(),
            activityUuid,
            0,
            0,
            DceRpc_Operation.CONTROL,
            (short) 0,
            new PnIoCm_Packet_Res((short) 0, (short) 0, (short) 0, (short) 0, 16696, (short) 0,
                Collections.singletonList(
                    new PnIoCm_Control_Response_ApplicationReady((short) 1, (short) 0,
                        arUuid, sessionKey, 0x0008, 0x0000))
            )
        );

        InetSocketAddress localAddress = (InetSocketAddress) pnChannel.getLocalAddress();
        InetSocketAddress remoteAddress = (InetSocketAddress) pnChannel.getRemoteAddress();

        // Serialize it to a byte-payload
        Random rand = new Random();
        Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
            rand.nextInt(65536),
            true,
            false,
            (short) 64,
            new IpAddress(localAddress.getAddress().getAddress()),
            new IpAddress(remoteAddress.getAddress().getAddress()),
            driverContext.getLocalPort(),
            sourcePort,
            packet
        );
        MacAddress srcAddress = new MacAddress(pnChannel.getLocalMacAddress().getAddress());
        MacAddress dstAddress = new MacAddress(pnChannel.getRemoteMacAddress().getAddress());
        return new Ethernet_Frame(
            dstAddress,
            srcAddress,
            udpFrame);
    }

    public static void sendApplicationReadyResponse(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext, int sourcePort, DceRpc_ActivityUuid activityUuid, Uuid arUuid, int sessionKey) {
        context.sendToWire(createApplicationReadyResponse(pnChannel, driverContext, sourcePort, activityUuid, arUuid, sessionKey));
    }

    public static Ethernet_Frame createPingResponse(RawSocketChannel pnChannel, ProfinetDriverContext driverContext, Ethernet_FramePayload_IPv4 payloadIPv4) {
        DceRpc_Packet pingRequest = payloadIPv4.getPayload();

        DceRpc_Packet packet = new DceRpc_Packet(DceRpc_PacketType.WORKING,
            false, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            pingRequest.getObjectUuid(), pingRequest.getInterfaceUuid(), pingRequest.getActivityUuid(),
            0L, 0L, DceRpc_Operation.CONNECT, (short) 0, new PnIoCm_Packet_Working());

        InetSocketAddress localAddress = (InetSocketAddress) pnChannel.getLocalAddress();
        InetSocketAddress remoteAddress = (InetSocketAddress) pnChannel.getRemoteAddress();

        // Serialize it to a byte-payload
        Random rand = new Random();
        Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
            rand.nextInt(65536),
            true,
            false,
            (short) 64,
            new IpAddress(localAddress.getAddress().getAddress()),
            new IpAddress(remoteAddress.getAddress().getAddress()),
            driverContext.getLocalPort(),
            payloadIPv4.getSourcePort(),
            packet
        );
        MacAddress srcAddress = new MacAddress(pnChannel.getLocalMacAddress().getAddress());
        MacAddress dstAddress = new MacAddress(pnChannel.getRemoteMacAddress().getAddress());
        return new Ethernet_Frame(
            dstAddress,
            srcAddress,
            udpFrame);
    }

    public static void sendPingResponse(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext, Ethernet_FramePayload_IPv4 payloadIPv4) {
        context.sendToWire(createPingResponse(pnChannel, driverContext, payloadIPv4));
    }

    /**
     * Simple helper that creates the UDP packet and Ethernet frame to transport the packet.
     *
     * @param pnChannel     the channel that contains the local and remote address information.
     * @param driverContext the context that contains the local and remote port information.
     * @param packet        the actual payload.
     * @return an Ethernet frame that we can send.
     */
    protected static Ethernet_Frame createEthernetFrame(RawSocketChannel pnChannel, ProfinetDriverContext driverContext, DceRpc_Packet packet) {
        InetSocketAddress localAddress = (InetSocketAddress) pnChannel.getLocalAddress();
        InetSocketAddress remoteAddress = (InetSocketAddress) pnChannel.getRemoteAddress();

        // Serialize it to a byte-payload
        Random rand = new Random();
        Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
            rand.nextInt(65536),
            true,
            false,
            (short) 64,
            new IpAddress(localAddress.getAddress().getAddress()),
            new IpAddress(remoteAddress.getAddress().getAddress()),
            driverContext.getLocalPort(),
            driverContext.getRemotePortImplicitCommunication(),
            packet
        );
        MacAddress srcAddress = new MacAddress(pnChannel.getLocalMacAddress().getAddress());
        MacAddress dstAddress = new MacAddress(pnChannel.getRemoteMacAddress().getAddress());
        return new Ethernet_Frame(
            dstAddress,
            srcAddress,
            udpFrame);
    }

}
