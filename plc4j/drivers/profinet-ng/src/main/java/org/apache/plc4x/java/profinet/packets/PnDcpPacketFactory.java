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

package org.apache.plc4x.java.profinet.packets;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;

import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.time.Duration;
import java.util.Collections;
import java.util.Random;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;

public class PnDcpPacketFactory {

    public static final int DEFAULT_UDP_PORT = 34964;

    public static Ethernet_Frame createIdentificationRequest(MacAddress localMacAddress, MacAddress remoteMacAddress) {
        // Construct and send the search request.
        return new Ethernet_Frame(
            remoteMacAddress,
            localMacAddress,
            new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, 0,
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
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if(ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .check(ethernetFramePayload -> ethernetFramePayload instanceof Ethernet_FramePayload_PnDcp)
            .unwrap(ethernetFramePayload -> (Ethernet_FramePayload_PnDcp) ethernetFramePayload)
            .unwrap(Ethernet_FramePayload_PnDcp::getPdu)
            .check(pnDcpPdu -> pnDcpPdu instanceof PnDcp_Pdu_IdentifyRes)
            .unwrap(pnDcpPdu -> (PnDcp_Pdu_IdentifyRes) pnDcpPdu)
            .handle(future::complete);
        return future;
    }

    public static Ethernet_Frame createReadIAndM0BlockRequest(RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        Random rand = new Random();

        InetSocketAddress localAddress = (InetSocketAddress) pnChannel.getLocalAddress();
        InetSocketAddress remoteAddress = (InetSocketAddress) pnChannel.getRemoteAddress();

        DceRpc_Packet packet = new DceRpc_Packet(
            DceRpc_PacketType.REQUEST, true, false, false,
            IntegerEncoding.LITTLE_ENDIAN, CharacterEncoding.ASCII, FloatingPointEncoding.IEEE,
            new DceRpc_ObjectUuid((byte) 0x00, 0x0001, driverContext.getDeviceId(), driverContext.getVendorId()),
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

        // Serialize it to a byte-payload
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

    public static CompletableFuture<PnIoCm_Block_IAndM0> sendReadIAndM0BlockRequest(ConversationContext<Ethernet_Frame> context, RawSocketChannel pnChannel, ProfinetDriverContext driverContext) {
        CompletableFuture<PnIoCm_Block_IAndM0> future = new CompletableFuture<>();
        context.sendRequest(PnDcpPacketFactory.createReadIAndM0BlockRequest(pnChannel, driverContext))
            .expectResponse(Ethernet_Frame.class, Duration.ofMillis(6000))
            .onTimeout(future::completeExceptionally)
            .onError((ethernetFrame, throwable) -> future.completeExceptionally(throwable))
            .unwrap(ethernetFrame -> {
                if(ethernetFrame.getPayload() instanceof Ethernet_FramePayload_VirtualLan) {
                    return ((Ethernet_FramePayload_VirtualLan) ethernetFrame.getPayload()).getPayload();
                }
                return ethernetFrame.getPayload();
            })
            .check(ethernetFramePayload -> ethernetFramePayload instanceof Ethernet_FramePayload_IPv4)
            .unwrap(ethernetFramePayload -> (Ethernet_FramePayload_IPv4) ethernetFramePayload)
            .unwrap(Ethernet_FramePayload_IPv4::getPayload)
            .check((dceRpcPacket -> dceRpcPacket.getPayload() instanceof PnIoCm_Packet_Res))
            .unwrap(dceRpcPacket -> (PnIoCm_Packet_Res) dceRpcPacket.getPayload())
            .handle(dceRpcPacketRes -> {
                if(dceRpcPacketRes.getBlocks().size() != 2) {
                    future.completeExceptionally(new PlcRuntimeException("Expected 2 blocks in the response"));
                    return;
                }
                if(!(dceRpcPacketRes.getBlocks().get(0) instanceof IODReadResponseHeader)) {
                    future.completeExceptionally(new PlcRuntimeException("The first block was expected to be of type IODReadResponseHeader"));
                    return;
                }
                if(!(dceRpcPacketRes.getBlocks().get(1) instanceof PnIoCm_Block_IAndM0)) {
                    future.completeExceptionally(new PlcRuntimeException("The second block was expected to be of type PnIoCm_Block_IAndM0"));
                    return;
                }
                PnIoCm_Block_IAndM0 iAndM0 = (PnIoCm_Block_IAndM0) dceRpcPacketRes.getBlocks().get(1);
                future.complete(iAndM0);
            });
        return future;
    }

}
