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

import java.time.Duration;
import java.util.Collections;
import java.util.Random;
import java.util.concurrent.CompletableFuture;
import java.util.function.Predicate;

/**
 * Builds the PROFINET request frames and wires up listener+future pairs to
 * await the matching responses. Ported from the old SPI's fluent
 * {@code ConversationContext.sendRequest().expectResponse()...handle()} chains
 * to the new SPI's {@link ProfinetConversation#expect(Predicate, Duration)} +
 * {@link ProfinetConversation#sendToWire(Ethernet_Frame)} primitives.
 */
public class PnDcpPacketFactory {

    private static final Duration DEFAULT_TIMEOUT = Duration.ofMillis(6000);

    public static Ethernet_Frame createIdentificationRequest(MacAddress localMacAddress, MacAddress remoteMacAddress) {
        return new Ethernet_Frame(
            remoteMacAddress,
            localMacAddress,
            new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, (short) 0,
                new Ethernet_FramePayload_PnDcp(
                    new PnDcp_Pdu_IdentifyReq(PnDcp_FrameId.DCP_Identify_ReqPDU.getValue(),
                        1L,
                        256,
                        Collections.singletonList(
                            new PnDcp_Block_ALLSelector()
                        )))));
    }

    public static CompletableFuture<PnDcp_Pdu_IdentifyRes> sendIdentificationRequest(
        ProfinetConversation conversation, MacAddress localMacAddress, MacAddress remoteMacAddress) {
        CompletableFuture<PnDcp_Pdu_IdentifyRes> future = new CompletableFuture<>();
        CompletableFuture<Ethernet_Frame> response = conversation.expect(
            frame -> {
                Ethernet_FramePayload payload = unwrapVlan(frame.getPayload());
                if (!(payload instanceof Ethernet_FramePayload_PnDcp)) {
                    return false;
                }
                PnDcp_Pdu pdu = ((Ethernet_FramePayload_PnDcp) payload).getPdu();
                return pdu instanceof PnDcp_Pdu_IdentifyRes;
            }, DEFAULT_TIMEOUT);
        conversation.sendToWire(createIdentificationRequest(localMacAddress, remoteMacAddress));

        response.whenComplete((frame, error) -> {
            if (error != null) {
                future.completeExceptionally(error);
                return;
            }
            Ethernet_FramePayload payload = unwrapVlan(frame.getPayload());
            future.complete((PnDcp_Pdu_IdentifyRes) ((Ethernet_FramePayload_PnDcp) payload).getPdu());
        });
        return future;
    }

    public static Ethernet_Frame createReadIAndM0BlockRequest(ProfinetConversation conv, ProfinetDriverContext ctx) {
        DceRpc_Packet packet = readImplicitRequest(ctx, 0xAFF0);
        return createEthernetFrame(conv, ctx, packet);
    }

    public static CompletableFuture<PnIoCm_Block_IAndM0> sendReadIAndM0BlockRequest(ProfinetConversation conv, ProfinetDriverContext ctx) {
        CompletableFuture<PnIoCm_Block_IAndM0> future = new CompletableFuture<>();
        CompletableFuture<Ethernet_Frame> response = conv.expect(
            frame -> isIPv4DceRpcResponse(frame, PnIoCm_Packet_Res.class), DEFAULT_TIMEOUT);
        conv.sendToWire(createReadIAndM0BlockRequest(conv, ctx));

        response.whenComplete((frame, error) -> {
            if (error != null) {
                future.completeExceptionally(error);
                return;
            }
            PnIoCm_Packet_Res res = (PnIoCm_Packet_Res) ((Ethernet_FramePayload_IPv4) unwrapVlan(frame.getPayload()))
                .getPayload().getPayload();
            if (res.getBlocks().size() != 2 || !(res.getBlocks().get(1) instanceof PnIoCm_Block_IAndM0)) {
                future.completeExceptionally(new PlcRuntimeException("Unexpected I&M0 response shape"));
                return;
            }
            future.complete((PnIoCm_Block_IAndM0) res.getBlocks().get(1));
        });
        return future;
    }

    public static Ethernet_Frame createReadIAndM1BlockRequest(ProfinetConversation conv, ProfinetDriverContext ctx) {
        DceRpc_Packet packet = readImplicitRequest(ctx, 0xAFF1);
        return createEthernetFrame(conv, ctx, packet);
    }

    public static CompletableFuture<PnIoCm_Block_IAndM1> sendReadIAndM1BlockRequest(ProfinetConversation conv, ProfinetDriverContext ctx) {
        CompletableFuture<PnIoCm_Block_IAndM1> future = new CompletableFuture<>();
        CompletableFuture<Ethernet_Frame> response = conv.expect(
            frame -> isIPv4DceRpcResponse(frame, PnIoCm_Packet_Res.class), DEFAULT_TIMEOUT);
        conv.sendToWire(createReadIAndM1BlockRequest(conv, ctx));

        response.whenComplete((frame, error) -> {
            if (error != null) {
                future.completeExceptionally(error);
                return;
            }
            PnIoCm_Packet_Res res = (PnIoCm_Packet_Res) ((Ethernet_FramePayload_IPv4) unwrapVlan(frame.getPayload()))
                .getPayload().getPayload();
            if (res.getBlocks().size() != 2 || !(res.getBlocks().get(1) instanceof PnIoCm_Block_IAndM1)) {
                future.completeExceptionally(new PlcRuntimeException("Unexpected I&M1 response shape"));
                return;
            }
            future.complete((PnIoCm_Block_IAndM1) res.getBlocks().get(1));
        });
        return future;
    }

    public static Ethernet_Frame createReadRealIdentificationDataRequest(ProfinetConversation conv, ProfinetDriverContext ctx) {
        DceRpc_Packet packet = readImplicitRequest(ctx, 0xF000);
        return createEthernetFrame(conv, ctx, packet);
    }

    public static CompletableFuture<PnIoCm_Block_RealIdentificationData> sendRealIdentificationDataRequest(
        ProfinetConversation conv, ProfinetDriverContext ctx) {
        CompletableFuture<PnIoCm_Block_RealIdentificationData> future = new CompletableFuture<>();
        CompletableFuture<Ethernet_Frame> response = conv.expect(
            frame -> {
                Ethernet_FramePayload payload = unwrapVlan(frame.getPayload());
                if (!(payload instanceof Ethernet_FramePayload_IPv4)) {
                    return false;
                }
                Object dceRpc = ((Ethernet_FramePayload_IPv4) payload).getPayload().getPayload();
                return dceRpc instanceof PnIoCm_Packet_Res || dceRpc instanceof PnIoCm_Packet_Rej;
            }, DEFAULT_TIMEOUT);
        conv.sendToWire(createReadRealIdentificationDataRequest(conv, ctx));

        response.whenComplete((frame, error) -> {
            if (error != null) {
                future.completeExceptionally(error);
                return;
            }
            Object dceRpc = ((Ethernet_FramePayload_IPv4) unwrapVlan(frame.getPayload()))
                .getPayload().getPayload();
            if (dceRpc instanceof PnIoCm_Packet_Rej) {
                future.completeExceptionally(new PlcRuntimeException("RealIdentificationData not supported"));
                return;
            }
            PnIoCm_Packet_Res res = (PnIoCm_Packet_Res) dceRpc;
            if (res.getBlocks().size() != 2 || !(res.getBlocks().get(1) instanceof PnIoCm_Block_RealIdentificationData)) {
                future.completeExceptionally(new PlcRuntimeException("Unexpected RealIdentificationData response shape"));
                return;
            }
            future.complete((PnIoCm_Block_RealIdentificationData) res.getBlocks().get(1));
        });
        return future;
    }

    public static Ethernet_Frame createParameterEndRequest(ProfinetConversation conv, ProfinetDriverContext ctx) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, ctx.getDeviceId(), ctx.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            ctx.getActivityUuid(),
            0L,
            (long) ctx.getAndIncrementIdentification(),
            DceRpc_Operation.CONTROL,
            (short) 0,
            new PnIoCm_Packet_Req(16696L, 16696L, 0L,
                Collections.singletonList(
                    new PnIoCm_Control_Request_ParameterEnd((short) 1, (short) 0,
                        new Uuid(ctx.getApplicationRelationUuid().getData()), 1, 1)
                ))
        );
        return createEthernetFrame(conv, ctx, packet);
    }

    public static CompletableFuture<PnIoCm_Control_Response_ParameterEnd> sendParameterEndRequest(
        ProfinetConversation conv, ProfinetDriverContext ctx) {
        CompletableFuture<PnIoCm_Control_Response_ParameterEnd> future = new CompletableFuture<>();
        CompletableFuture<Ethernet_Frame> response = conv.expect(
            frame -> {
                Ethernet_FramePayload payload = unwrapVlan(frame.getPayload());
                if (!(payload instanceof Ethernet_FramePayload_IPv4)) {
                    return false;
                }
                Object dceRpc = ((Ethernet_FramePayload_IPv4) payload).getPayload().getPayload();
                return dceRpc instanceof PnIoCm_Packet_Res || dceRpc instanceof PnIoCm_Packet_Rej;
            }, DEFAULT_TIMEOUT);
        conv.sendToWire(createParameterEndRequest(conv, ctx));

        response.whenComplete((frame, error) -> {
            if (error != null) {
                future.completeExceptionally(error);
                return;
            }
            Object dceRpc = ((Ethernet_FramePayload_IPv4) unwrapVlan(frame.getPayload()))
                .getPayload().getPayload();
            if (dceRpc instanceof PnIoCm_Packet_Rej) {
                future.completeExceptionally(new PlcRuntimeException("ParameterEnd not supported"));
                return;
            }
            PnIoCm_Packet_Res res = (PnIoCm_Packet_Res) dceRpc;
            if (res.getBlocks().size() != 1 || !(res.getBlocks().get(0) instanceof PnIoCm_Control_Response_ParameterEnd)) {
                future.completeExceptionally(new PlcRuntimeException("Unexpected ParameterEnd response shape"));
                return;
            }
            future.complete((PnIoCm_Control_Response_ParameterEnd) res.getBlocks().get(0));
        });
        return future;
    }

    public static Ethernet_Frame createApplicationReadyResponse(ProfinetConversation conv, ProfinetDriverContext ctx,
                                                                int sourcePort, DceRpc_ActivityUuid activityUuid,
                                                                Uuid arUuid, int sessionKey) {
        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.RESPONSE, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, ctx.getDeviceId(), ctx.getVendorId()),
            new DceRpc_InterfaceUuid_ControllerInterface(),
            activityUuid,
            0L,
            0L,
            DceRpc_Operation.CONTROL,
            (short) 0,
            new PnIoCm_Packet_Res((short) 0, (short) 0, (short) 0, (short) 0, 16696L, 0L,
                Collections.singletonList(
                    new PnIoCm_Control_Response_ApplicationReady((short) 1, (short) 0,
                        arUuid, sessionKey, 0x0008, 0x0000))
            )
        );
        return wrapAsUdpEthernet(conv, ctx.getLocalPort(), sourcePort, packet);
    }

    public static void sendApplicationReadyResponse(ProfinetConversation conv, ProfinetDriverContext ctx,
                                                    int sourcePort, DceRpc_ActivityUuid activityUuid,
                                                    Uuid arUuid, int sessionKey) {
        conv.sendToWire(createApplicationReadyResponse(conv, ctx, sourcePort, activityUuid, arUuid, sessionKey));
    }

    public static Ethernet_Frame createPingResponse(ProfinetConversation conv, ProfinetDriverContext ctx,
                                                    Ethernet_FramePayload_IPv4 payloadIPv4) {
        DceRpc_Packet pingRequest = payloadIPv4.getPayload();
        DceRpc_Packet packet = new DceRpc_Packet(DceRpc_PacketType.WORKING,
            false, false, false,
            IntegerEncoding.BIG_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            pingRequest.getObjectUuid(), pingRequest.getInterfaceUuid(), pingRequest.getActivityUuid(),
            0L, 0L, DceRpc_Operation.CONNECT, (short) 0, new PnIoCm_Packet_Working());
        return wrapAsUdpEthernet(conv, ctx.getLocalPort(), payloadIPv4.getSourcePort(), packet);
    }

    public static void sendPingResponse(ProfinetConversation conv, ProfinetDriverContext ctx,
                                        Ethernet_FramePayload_IPv4 payloadIPv4) {
        conv.sendToWire(createPingResponse(conv, ctx, payloadIPv4));
    }

    /////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    /////////////////////////////////////////////////////////////////////////////////////////////////

    private static DceRpc_Packet readImplicitRequest(ProfinetDriverContext ctx, int slotNumberToReadFromIndex) {
        return new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, (short) 0x0001, ctx.getDeviceId(), ctx.getVendorId()),
            new DceRpc_InterfaceUuid_DeviceInterface(),
            ctx.getActivityUuid(),
            0L,
            (long) ctx.getAndIncrementIdentification(),
            DceRpc_Operation.READ_IMPLICIT,
            (short) 0,
            new PnIoCm_Packet_Req(16696L, 16696L, 0L,
                Collections.singletonList(
                    new IODReadRequestHeader((short) 1, (short) 0, 0,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}),
                        0L, 0,
                        // Some endpoints use subslot 1, others 0; the original code differed
                        // per call site. We pass 0 here — both seem accepted by devices.
                        (slotNumberToReadFromIndex == 0xF000) ? 0 : 1,
                        slotNumberToReadFromIndex, 16696L,
                        new Uuid(new byte[]{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}))
                ))
        );
    }

    private static Ethernet_FramePayload unwrapVlan(Ethernet_FramePayload payload) {
        if (payload instanceof Ethernet_FramePayload_VirtualLan) {
            return ((Ethernet_FramePayload_VirtualLan) payload).getPayload();
        }
        return payload;
    }

    private static boolean isIPv4DceRpcResponse(Ethernet_Frame frame, Class<?> expected) {
        Ethernet_FramePayload payload = unwrapVlan(frame.getPayload());
        if (!(payload instanceof Ethernet_FramePayload_IPv4)) {
            return false;
        }
        Object dceRpc = ((Ethernet_FramePayload_IPv4) payload).getPayload().getPayload();
        return expected.isInstance(dceRpc);
    }

    private static Ethernet_Frame createEthernetFrame(ProfinetConversation conv, ProfinetDriverContext ctx,
                                                      DceRpc_Packet packet) {
        return wrapAsUdpEthernet(conv, ctx.getLocalPort(), ctx.getRemotePortImplicitCommunication(), packet);
    }

    private static Ethernet_Frame wrapAsUdpEthernet(ProfinetConversation conv, int srcPort, int dstPort,
                                                    DceRpc_Packet packet) {
        Random rand = new Random();
        Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
            rand.nextInt(65536),
            true,
            false,
            (short) 64,
            conv.getLocalIp(),
            conv.getRemoteIp(),
            srcPort,
            dstPort,
            packet
        );
        return new Ethernet_Frame(conv.getRemoteMacAddress(), conv.getLocalMacAddress(), udpFrame);
    }

}
