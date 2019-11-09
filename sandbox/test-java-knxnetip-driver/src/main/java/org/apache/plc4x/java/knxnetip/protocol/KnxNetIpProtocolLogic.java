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
package org.apache.plc4x.java.knxnetip.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.socket.DatagramChannel;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.events.ConnectEvent;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.knxnetip.events.KnxGatewayFoundEvent;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.apache.plc4x.java.knxnetip.readwrite.types.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.types.KnxLayer;
import org.apache.plc4x.java.knxnetip.readwrite.types.Status;

import java.net.InetSocketAddress;
import java.util.Arrays;
import java.util.List;

public class KnxNetIpProtocolLogic extends PlcMessageToMessageCodec<KNXNetIPMessage, PlcRequestContainer> {

    private KNXAddress gatewayAddress;
    private String gatewayName;

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        if (evt instanceof ConnectEvent) {
            DatagramChannel channel = (DatagramChannel) ctx.pipeline().channel();
            final InetSocketAddress localSocketAddress = channel.localAddress();
            final IPAddress localAddress = new IPAddress(localSocketAddress.getAddress().getAddress());
            final int localPort = localSocketAddress.getPort();
            SearchRequest searchRequest = new SearchRequest(new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localAddress, localPort));
            ctx.channel().writeAndFlush(searchRequest);
        } else if(evt instanceof KnxGatewayFoundEvent) {
            DatagramChannel channel = (DatagramChannel) ctx.pipeline().channel();
            final InetSocketAddress localSocketAddress = channel.localAddress();
            final IPAddress localAddress = new IPAddress(localSocketAddress.getAddress().getAddress());
            final int localPort = localSocketAddress.getPort();
            ConnectionRequest connectionRequest = new ConnectionRequest(
                new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localAddress, localPort),
                new HPAIDataEndpoint(HostProtocolCode.IPV4_UDP, localAddress, localPort),
                new ConnectionRequestInformationTunnelConnection(KnxLayer.TUNNEL_BUSMONITOR));
            ctx.channel().writeAndFlush(connectionRequest);
        }
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) {
        // Ignore ...
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, KNXNetIPMessage msg, List<Object> out) {
        if(msg instanceof SearchResponse) {
            SearchResponse searchResponse = (SearchResponse) msg;
            final ServiceId tunnelingService = Arrays.stream(searchResponse.getDibSuppSvcFamilies().getServiceIds()).filter(serviceId -> serviceId instanceof KnxNetIpTunneling).findFirst().orElse(null);
            // If this service has the
            if(tunnelingService != null) {
                gatewayAddress = searchResponse.getDibDeviceInfo().getKnxAddress();
                gatewayName = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName());
                ctx.channel().pipeline().fireUserEventTriggered(new KnxGatewayFoundEvent());
            }
        } else if(msg instanceof ConnectionResponse) {
            ConnectionResponse connectionResponse = (ConnectionResponse) msg;
            Status status = connectionResponse.getStatus();
            if(status == Status.NO_ERROR) {
                ctx.channel().pipeline().fireUserEventTriggered(new ConnectedEvent());
            }
        } else if(msg instanceof TunnelingRequest) {
            TunnelingRequest tunnelingRequest = (TunnelingRequest) msg;
            final short communicationChannelId = tunnelingRequest.getTunnelingRequestDataBlock().getCommunicationChannelId();
            final short sequenceCounter = tunnelingRequest.getTunnelingRequestDataBlock().getSequenceCounter();
            TunnelingResponse tunnelingResponse = new TunnelingResponse(
                new TunnelingResponseDataBlock(communicationChannelId, sequenceCounter, Status.NO_ERROR));
            ctx.channel().writeAndFlush(tunnelingResponse);
            CEMIBusmonInd busmonInd = (CEMIBusmonInd) tunnelingRequest.getCemi();
            if(busmonInd.getCemiFrame() instanceof CEMIFrameData) {
                outputStringRepresentation((CEMIFrameData) busmonInd.getCemiFrame());
            }
        }
    }

    private void outputStringRepresentation(CEMIFrameData data) {
        final CEMIAddress sourceAddress = data.getSourceAddress();
        final CEMIAddress destinationAddress = data.getDestinationAddress();
        final boolean groupAddress = data.getGroupAddress();
        final byte[] payload = data.getData();
        String payloadString = Hex.encodeHexString(payload);
        if(groupAddress) {
            System.out.println(String.format("Telegram from %d.%d.%d to %d/%d/%d with payload %s",
                sourceAddress.getArea(), sourceAddress.getLine(), sourceAddress.getDevice(),
                destinationAddress.getArea(), destinationAddress.getLine(), destinationAddress.getDevice(),
                payloadString));
        } else {
            System.out.println(String.format("Telegram from %d.%d.%d to %d.%d.%d with payload %s",
                sourceAddress.getArea(), sourceAddress.getLine(), sourceAddress.getDevice(),
                destinationAddress.getArea(), destinationAddress.getLine(), destinationAddress.getDevice(),
                payloadString));
        }
    }

}
