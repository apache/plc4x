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

import io.netty.channel.socket.DatagramChannel;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.builder.ToStringStyle;
import org.apache.plc4x.java.ets5.passive.*;
import org.apache.plc4x.java.ets5.passive.io.KNXGroupAddressIO;
import org.apache.plc4x.java.ets5.passive.io.KnxDatapointIO;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.ets5.Ets5Parser;
import org.apache.plc4x.java.knxnetip.ets5.model.Ets5Model;
import org.apache.plc4x.java.knxnetip.ets5.model.GroupAddress;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.HasConfiguration;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.apache.plc4x.java.knxnetip.readwrite.types.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.types.KnxLayer;
import org.apache.plc4x.java.knxnetip.readwrite.types.Status;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.net.InetSocketAddress;
import java.time.Duration;
import java.util.*;

public class KnxNetIpProtocolLogic extends Plc4xProtocolBase<KNXNetIPMessage> implements HasConfiguration<KnxNetIpConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpProtocolLogic.class);

    private KNXAddress gatewayAddress;
    private String gatewayName;
    private IPAddress localIPAddress;
    private int localPort;
    private short communicationChannelId;

    private Timer connectionStateTimer;

    private byte groupAddressType;
    private Ets5Model ets5Model;

    @Override public void setConfiguration(KnxNetIpConfiguration configuration) {
        if (configuration.knxprojFilePath != null) {
            File knxprojFile = new File(configuration.knxprojFilePath);
            if (knxprojFile.exists() && knxprojFile.isFile()) {
                ets5Model = new Ets5Parser().parse(knxprojFile);
                groupAddressType = ets5Model.getGroupAddressType();
            } else {
                throw new RuntimeException(String.format(
                    "File specified with 'knxproj-file-path' does not exist or is not a file: '%s'",
                    configuration.knxprojFilePath));
            }
        }
    }

    @Override
    public void onConnect(ConversationContext<KNXNetIPMessage> context) {
        DatagramChannel channel = (DatagramChannel) context.getChannel();
        final InetSocketAddress localSocketAddress = channel.localAddress();
        localIPAddress = new IPAddress(localSocketAddress.getAddress().getAddress());
        localPort = localSocketAddress.getPort();

        // First send out a search request
        // REMARK: This might be optional ... usually we would send a search request to ip 224.0.23.12
        // Any KNX Gateway will respond with a search response. We're currently directly sending to the
        // known gateway address, so it's sort of pointless, but at least only one device will respond.
        SearchRequest searchRequest = new SearchRequest(
            new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
        context.sendRequest(searchRequest)
            .expectResponse(KNXNetIPMessage.class, Duration.ofMillis(1000))
            .check(p -> p instanceof SearchResponse)
            .unwrap(p -> (SearchResponse) p)
            .handle(searchResponse -> {
                // Check if this device supports tunneling services.
                final ServiceId tunnelingService = Arrays.stream(searchResponse.getDibSuppSvcFamilies().getServiceIds()).filter(serviceId -> serviceId instanceof KnxNetIpTunneling).findFirst().orElse(null);

                // If this device supports this type of service, tell the driver, we found a suitable device.
                if(tunnelingService != null) {
                    // Extract the required information form the search request.
                    gatewayAddress = searchResponse.getDibDeviceInfo().getKnxAddress();
                    gatewayName = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName()).trim();

                    LOGGER.info(String.format("Found KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                        gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));

                    // Next send a connection request to the gateway.
                    ConnectionRequest connectionRequest = new ConnectionRequest(
                        new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
                        new HPAIDataEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
                        new ConnectionRequestInformationTunnelConnection(KnxLayer.TUNNEL_BUSMONITOR));
                    context.sendRequest(connectionRequest)
                        .expectResponse(KNXNetIPMessage.class, Duration.ofMillis(1000))
                        .check(p -> p instanceof ConnectionResponse)
                        .unwrap(p -> (ConnectionResponse) p)
                        .handle(connectionResponse -> {
                            // Remember the communication channel id.
                            communicationChannelId = connectionResponse.getCommunicationChannelId();

                            // Check if everything went well.
                            Status status = connectionResponse.getStatus();
                            if (status == Status.NO_ERROR) {
                                LOGGER.info(String.format("Connected to KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                                    gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));

                                // Send an event that connection setup is complete.
                                context.fireConnected();

                                // Start a timer to check the connection state every 60 seconds.
                                // This keeps the connection open if no data is transported.
                                // Otherwise the gateway will terminate the connection.
                                connectionStateTimer = new Timer();
                                connectionStateTimer.scheduleAtFixedRate(new TimerTask() {
                                    @Override
                                    public void run() {
                                        ConnectionStateRequest connectionStateRequest =
                                            new ConnectionStateRequest(communicationChannelId,
                                                new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
                                        context.sendRequest(connectionStateRequest)
                                            .expectResponse(KNXNetIPMessage.class, Duration.ofMillis(1000))
                                            .check(p -> p instanceof ConnectionStateResponse)
                                            .unwrap(p -> (ConnectionStateResponse) p)
                                            .handle(connectionStateResponse -> {
                                                if(connectionStateResponse.getStatus() != Status.NO_ERROR) {
                                                    if(connectionStateResponse.getStatus() != null) {
                                                        LOGGER.error(String.format("Connection state problems. Got %s",
                                                            connectionStateResponse.getStatus().name()));
                                                    } else {
                                                        LOGGER.error("Connection state problems. Got no status information.");
                                                    }
                                                }
                                            });
                                    }
                                }, 60000, 60000);
                            } else {
                                // The connection request wasn't successful.
                            }
                        });
                } else {
                    // This device doesn't support tunneling ... do some error handling.
                }
            });
    }

    @Override
    public void onDisconnect(ConversationContext<KNXNetIPMessage> context) {
        // Cancel the timer for sending connection state requests.
        connectionStateTimer.cancel();

        DisconnectRequest disconnectRequest = new DisconnectRequest(communicationChannelId,
            new HPAIControlEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
        context.sendRequest(disconnectRequest)
            .expectResponse(KNXNetIPMessage.class, Duration.ofMillis(1000))
            .check(p -> p instanceof DisconnectResponse)
            .unwrap(p -> (DisconnectResponse) p)
            .handle(disconnectResponse -> {
                // In general we should probably check if the disconnect was successful, but in
                // the end we couldn't do much if the disconnect would fail.
                LOGGER.info(String.format("Disconnected from KNX Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                    gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));

                // Send an event that connection disconnect is complete.
                context.fireDisconnected();
            });
    }

    @Override
    protected void decode(ConversationContext<KNXNetIPMessage> context, KNXNetIPMessage msg) throws Exception {
        // Handle a normal tunneling request, which is delivering KNX data.
        if(msg instanceof TunnelingRequest) {
            TunnelingRequest tunnelingRequest = (TunnelingRequest) msg;
            final short curCommunicationChannelId =
                tunnelingRequest.getTunnelingRequestDataBlock().getCommunicationChannelId();

            // Only if the communication channel id match, do anything with the request.
            if(curCommunicationChannelId == communicationChannelId) {
                CEMIBusmonInd busmonInd = (CEMIBusmonInd) tunnelingRequest.getCemi();
                if (busmonInd.getCemiFrame() instanceof CEMIFrameData) {
                    CEMIFrameData cemiDataFrame = (CEMIFrameData) busmonInd.getCemiFrame();

                    // The first byte is actually just 6 bit long, but we'll treat it as a full one.
                    // So here we create a byte array containing the first and all the following bytes.
                    byte[] payload = new byte[1 + cemiDataFrame.getData().length];
                    payload[0] = cemiDataFrame.getDataFirstByte();
                    System.arraycopy(cemiDataFrame.getData(), 0, payload, 1, cemiDataFrame.getData().length);

                    final KNXAddress sourceAddress = cemiDataFrame.getSourceAddress();
                    final byte[] destinationGroupAddress = cemiDataFrame.getDestinationAddress();

                    // Decode the group address depending on the project settings.
                    ReadBuffer addressReadBuffer = new ReadBuffer(destinationGroupAddress);
                    KNXGroupAddress destinationAddress =
                        KNXGroupAddressIO.parse(addressReadBuffer, groupAddressType);

                    // If there is an ETS5 model provided, continue decoding the payload.
                    if(ets5Model != null) {
                        final GroupAddress groupAddress = ets5Model.getGroupAddresses().get(destinationAddress);

                        ReadBuffer rawDataReader = new ReadBuffer(payload);

                        final KnxDatapoint datapoint = KnxDatapointIO.parse(rawDataReader, groupAddress.getType().getMainType(), groupAddress.getType().getSubType());
                        final String jsonDatapoint = datapoint.toString(ToStringStyle.JSON_STYLE);

                        LOGGER.info("Message from: '" + KnxNetIpProtocolLogic.toString(sourceAddress) + "'" +
                            " to: '" + KnxNetIpProtocolLogic.toString(destinationAddress) + "'" +
                            "\n location: '" + groupAddress.getFunction().getSpaceName() + "'" +
                            " function: '" + groupAddress.getFunction().getName() + "'" +
                            " meaning: '" + groupAddress.getName() + "'" +
                            " type: '" + groupAddress.getType().getName() + "'" +
                            "\n value: '" + jsonDatapoint + "'"
                        );
                    }
                    // Else just output the raw payload.
                    else {
                        LOGGER.info("Raw Message: '" + KnxNetIpProtocolLogic.toString(sourceAddress) + "'" +
                            " to: '" + KnxNetIpProtocolLogic.toString(destinationAddress) + "'" +
                            "\n payload: '" + Hex.encodeHexString(payload) + "'"
                        );
                    }
                }

                // Confirm receipt of the request.
                final short sequenceCounter = tunnelingRequest.getTunnelingRequestDataBlock().getSequenceCounter();
                TunnelingResponse tunnelingResponse = new TunnelingResponse(
                    new TunnelingResponseDataBlock(communicationChannelId, sequenceCounter, Status.NO_ERROR));
                context.sendToWire(tunnelingResponse);
            }
        }
    }

    protected static String toString(KNXAddress knxAddress) {
        return knxAddress.getMainGroup() + "." + knxAddress.getMiddleGroup() + "." + knxAddress.getSubGroup();
    }

    protected static String toString(KNXGroupAddress groupAddress) {
        if(groupAddress instanceof KNXGroupAddress3Level) {
            KNXGroupAddress3Level level3 = (KNXGroupAddress3Level) groupAddress;
            return level3.getMainGroup() + "/" + level3.getMiddleGroup() + "/" + level3.getSubGroup();
        } else if(groupAddress instanceof KNXGroupAddress2Level) {
            KNXGroupAddress2Level level2 = (KNXGroupAddress2Level) groupAddress;
            return level2.getMainGroup() + "/" + level2.getSubGroup();
        } else if(groupAddress instanceof KNXGroupAddressFreeLevel) {
            KNXGroupAddressFreeLevel free = (KNXGroupAddressFreeLevel) groupAddress;
            return free.getSubGroup() + "";
        }
        throw new RuntimeException("Unsupported Group Address Type " + groupAddress.getClass().getName());
    }
}
