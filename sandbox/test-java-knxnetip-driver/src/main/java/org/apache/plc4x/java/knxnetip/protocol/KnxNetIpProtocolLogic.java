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
import org.apache.plc4x.java.spi.PlcMessageToMessageCodec;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.spi.events.DisconnectEvent;
import org.apache.plc4x.java.spi.events.DisconnectedEvent;
import org.apache.plc4x.java.knxnetip.events.KnxGatewayFoundEvent;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.apache.plc4x.java.knxnetip.readwrite.types.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.types.KnxLayer;
import org.apache.plc4x.java.knxnetip.readwrite.types.Status;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetSocketAddress;
import java.util.*;
import java.util.concurrent.CompletableFuture;

public class KnxNetIpProtocolLogic extends PlcMessageToMessageCodec<KNXNetIPMessage, KNXNetIPMessage> {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpProtocolLogic.class);

    private KNXAddress gatewayAddress;
    private String gatewayName;
    private IPAddress localIPAddress;
    private int localPort;
    private short communicationChannelId;

    private Timer connectionStateTimer;
    private CompletableFuture<Void> disconnectFuture;

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) {
        if (evt instanceof ConnectEvent) {
            DatagramChannel channel = (DatagramChannel) ctx.pipeline().channel();
            final InetSocketAddress localSocketAddress = channel.localAddress();
            localIPAddress = new IPAddress(localSocketAddress.getAddress().getAddress());
            localPort = localSocketAddress.getPort();
            SearchRequest searchRequest = new SearchRequest(
                new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
            ctx.channel().writeAndFlush(searchRequest);
        } else if(evt instanceof KnxGatewayFoundEvent) {
            DatagramChannel channel = (DatagramChannel) ctx.pipeline().channel();
            final InetSocketAddress localSocketAddress = channel.localAddress();
            localIPAddress = new IPAddress(localSocketAddress.getAddress().getAddress());
            localPort = localSocketAddress.getPort();
            ConnectionRequest connectionRequest = new ConnectionRequest(
                new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
                new HPAIDataEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
                new ConnectionRequestInformationTunnelConnection(KnxLayer.TUNNEL_BUSMONITOR));
            ctx.channel().writeAndFlush(connectionRequest);
        } else if(evt instanceof DisconnectEvent) {
            DisconnectEvent disconnectEvent = (DisconnectEvent) evt;
            disconnectFuture = disconnectEvent.getFuture();
            DisconnectRequest disconnectRequest = new DisconnectRequest(communicationChannelId,
                new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
            ctx.channel().writeAndFlush(disconnectRequest);
        }
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, KNXNetIPMessage msg, List<Object> out) {
        out.add(msg);
    }

    @Override
    protected void decode(ChannelHandlerContext ctx, KNXNetIPMessage msg, List<Object> out) {
        // Handle search responses to find the device able to provide tunneling services
        if(msg instanceof SearchResponse) {
            SearchResponse searchResponse = (SearchResponse) msg;
            // Check if this device supports tunneling services.
            final ServiceId tunnelingService = Arrays.stream(searchResponse.getDibSuppSvcFamilies().getServiceIds()).filter(serviceId -> serviceId instanceof KnxNetIpTunneling).findFirst().orElse(null);
            // If this device supports this type of service, tell the driver, we found a suitable device.
            if(tunnelingService != null) {
                gatewayAddress = searchResponse.getDibDeviceInfo().getKnxAddress();
                gatewayName = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName()).trim();
                LOGGER.info(String.format("Found KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                    gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));
                ctx.channel().pipeline().fireUserEventTriggered(new KnxGatewayFoundEvent());
            }
        }

        // Handle the response to a connection request.
        else if(msg instanceof ConnectionResponse) {
            ConnectionResponse connectionResponse = (ConnectionResponse) msg;
            Status status = connectionResponse.getStatus();
            // Remember the communication channel id.
            communicationChannelId = connectionResponse.getCommunicationChannelId();
            if (status == Status.NO_ERROR) {
                LOGGER.info(String.format("Connected to KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                    gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));
                ctx.channel().pipeline().fireUserEventTriggered(new ConnectedEvent());
                // Start a timer to check the connection state every 60 seconds.
                connectionStateTimer = new Timer();
                connectionStateTimer.scheduleAtFixedRate(new TimerTask() {
                    @Override
                    public void run() {
                        ConnectionStateRequest connectionStateRequest =
                            new ConnectionStateRequest(communicationChannelId,
                                new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
                        ctx.channel().writeAndFlush(connectionStateRequest);
                    }
                }, 60000, 60000);
            } else {
                LOGGER.error(String.format("Error connecting to KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                    gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));
            }
        }

        // Handle the responses to the connection state requests.
        else if(msg instanceof ConnectionStateResponse) {
            ConnectionStateResponse connectionStateResponse = (ConnectionStateResponse) msg;
            if(connectionStateResponse.getStatus() != Status.NO_ERROR) {
                if(connectionStateResponse.getStatus() != null) {
                    LOGGER.error(String.format("Connection state problems. Got %s",
                        connectionStateResponse.getStatus().name()));
                } else {
                    LOGGER.error("Connection state problems. Got no status information.");
                }
            }
        }

        // Handle a normal tunneling request, which is delivering KNX data.
        else if(msg instanceof TunnelingRequest) {
            TunnelingRequest tunnelingRequest = (TunnelingRequest) msg;
            final short curCommunicationChannelId =
                tunnelingRequest.getTunnelingRequestDataBlock().getCommunicationChannelId();
            // Only if the communication channel id match, do anything with the request.
            if(curCommunicationChannelId == communicationChannelId) {
                final short sequenceCounter = tunnelingRequest.getTunnelingRequestDataBlock().getSequenceCounter();
                TunnelingResponse tunnelingResponse = new TunnelingResponse(
                    new TunnelingResponseDataBlock(communicationChannelId, sequenceCounter, Status.NO_ERROR));
                ctx.channel().writeAndFlush(tunnelingResponse);
                out.add(tunnelingRequest);
            }
        }

        // Handle the cleaning up after getting the response to a disconnect request.
        else if(msg instanceof DisconnectResponse) {
            // In general we should probably check if the disconnect was successful, but in
            // the end we couldn't do much if the disconnect would fail.
            ctx.channel().pipeline().fireUserEventTriggered(new DisconnectedEvent());
            LOGGER.info(String.format("Disconnected from KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));
            // Notify the closer, that we're done disconnecting.
            if(disconnectFuture != null) {
                disconnectFuture.complete(null);
            }
        }
    }

}
