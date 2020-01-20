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
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcString;
import org.apache.plc4x.java.api.value.PlcStruct;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.ets5.Ets5Parser;
import org.apache.plc4x.java.knxnetip.ets5.model.Ets5Model;
import org.apache.plc4x.java.knxnetip.ets5.model.GroupAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddress2Level;
import org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddress3Level;
import org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddressFreeLevel;
import org.apache.plc4x.java.knxnetip.readwrite.io.KnxDatapointIO;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.apache.plc4x.java.knxnetip.readwrite.types.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.types.KnxLayer;
import org.apache.plc4x.java.knxnetip.readwrite.types.Status;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionHandle;
import org.apache.plc4x.java.spi.model.InternalPlcSubscriptionHandle;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.net.InetSocketAddress;
import java.time.Duration;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

public class KnxNetIpProtocolLogic extends Plc4xProtocolBase<KNXNetIPMessage> implements HasConfiguration<KnxNetIpConfiguration>, PlcSubscriber {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpProtocolLogic.class);

    private KNXAddress gatewayAddress;
    private String gatewayName;
    private IPAddress localIPAddress;
    private int localPort;
    private short communicationChannelId;

    private Timer connectionStateTimer;

    private byte groupAddressType;
    private Ets5Model ets5Model;

    private Map<Integer, Consumer<PlcSubscriptionEvent>> consumerIdMap = new ConcurrentHashMap<>();

    @Override
    public void setConfiguration(KnxNetIpConfiguration configuration) {
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
        } else {
            groupAddressType = (byte) configuration.groupAddressType;
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
        LOGGER.info("Sending KNXnet/IP Search Request.");
        SearchRequest searchRequest = new SearchRequest(
            new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort));
        context.sendRequest(searchRequest)
            .expectResponse(KNXNetIPMessage.class, Duration.ofMillis(1000))
            .check(p -> p instanceof SearchResponse)
            .unwrap(p -> (SearchResponse) p)
            .handle(searchResponse -> {
                LOGGER.info("Got KNXnet/IP Search Response.");
                // Check if this device supports tunneling services.
                final ServiceId tunnelingService = Arrays.stream(searchResponse.getDibSuppSvcFamilies().getServiceIds()).filter(serviceId -> serviceId instanceof KnxNetIpTunneling).findFirst().orElse(null);

                // If this device supports this type of service, tell the driver, we found a suitable device.
                if(tunnelingService != null) {
                    // Extract the required information form the search request.
                    gatewayAddress = searchResponse.getDibDeviceInfo().getKnxAddress();
                    gatewayName = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName()).trim();

                    LOGGER.info(String.format("Found KNXnet/IP Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
                        gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(), gatewayAddress.getSubGroup()));

                    // Next send a connection request to the gateway.
                    ConnectionRequest connectionRequest = new ConnectionRequest(
                        new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
                        new HPAIDataEndpoint(HostProtocolCode.IPV4_UDP, localIPAddress, localPort),
                        new ConnectionRequestInformationTunnelConnection(KnxLayer.TUNNEL_BUSMONITOR));
                    LOGGER.info("Sending KNXnet/IP Connection Request.");
                    context.sendRequest(connectionRequest)
                        .expectResponse(KNXNetIPMessage.class, Duration.ofMillis(1000))
                        .check(p -> p instanceof ConnectionResponse)
                        .unwrap(p -> (ConnectionResponse) p)
                        .handle(connectionResponse -> {
                            // Remember the communication channel id.
                            communicationChannelId = connectionResponse.getCommunicationChannelId();

                            LOGGER.info(String.format("Received KNXnet/IP Connection Response (Connection Id %s)",
                                communicationChannelId));

                            // Check if everything went well.
                            Status status = connectionResponse.getStatus();
                            if (status == Status.NO_ERROR) {
                                LOGGER.info(String.format("Successfully connected to KNXnet/IP Gateway '%s' with KNX address '%d.%d.%d'", gatewayName,
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
                                LOGGER.error(String.format(
                                    "Not connected to KNXnet/IP Gateway '%s' with KNX address '%d.%d.%d' got status: '%s'",
                                    gatewayName, gatewayAddress.getMainGroup(), gatewayAddress.getMiddleGroup(),
                                    gatewayAddress.getSubGroup(), status.toString()));
                                // TODO: Actively disconnect
                            }
                        });
                } else {
                    // This device doesn't support tunneling ... do some error handling.
                    LOGGER.error("Not connected to KNCnet/IP Gateway. The device doesn't support Tunneling.");
                    // TODO: Actively disconnect
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
                    final String destinationAddress =
                        Ets5Model.parseGroupAddress(groupAddressType, destinationGroupAddress);

                    // If there is an ETS5 model provided, continue decoding the payload.
                    if(ets5Model != null) {
                        final GroupAddress groupAddress = ets5Model.getGroupAddresses().get(destinationAddress);

                        if(groupAddress != null) {
                            LOGGER.info("Message from: '" + toString(sourceAddress) +
                                "' to: '" + destinationAddress + "'");

                            // Parse the payload depending on the type of the group-address.
                            ReadBuffer rawDataReader = new ReadBuffer(payload);
                            final PlcValue value = KnxDatapointIO.parse(rawDataReader,
                                groupAddress.getType().getMainType(), groupAddress.getType().getSubType());

                            // Assemble the plc4x return data-structure.
                            Map<String, PlcValue> dataPointMap = new HashMap<>();
                            dataPointMap.put("sourceAddress", new PlcString(toString(sourceAddress)));
                            dataPointMap.put("targetAddress", new PlcString(groupAddress.getGroupAddress()));
                            dataPointMap.put("name", new PlcString(groupAddress.getName()));
                            dataPointMap.put("type", new PlcString(groupAddress.getType().getName()));
                            dataPointMap.put("functionId", new PlcString(groupAddress.getFunction().getId()));
                            dataPointMap.put("functionName", new PlcString(groupAddress.getFunction().getName()));
                            dataPointMap.put("functionType", new PlcString(groupAddress.getFunction().getType()));
                            dataPointMap.put("functionSpace", new PlcString(groupAddress.getFunction().getSpaceName()));
                            dataPointMap.put("value", value);
                            final PlcStruct dataPoint = new PlcStruct(dataPointMap);

                            // Send the data-structure.
                            publishEvent("knxData", dataPoint);
                        } else {
                            LOGGER.warn("Message from: '" + toString(sourceAddress) + "'" +
                                " to unknown group address: '" + destinationAddress + "'" +
                                "\n payload: '" + Hex.encodeHexString(payload) + "'");
                        }
                    }
                    // Else just output the raw payload.
                    else {
                        LOGGER.info("Raw Message: '" + KnxNetIpProtocolLogic.toString(sourceAddress) + "'" +
                            " to: '" + destinationAddress + "'" +
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

    @Override
    public void close(ConversationContext<KNXNetIPMessage> context) {
        // TODO Implement Closing on Protocol Level
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> values = new HashMap<>();
        for (String fieldName : subscriptionRequest.getFieldNames()) {
            values.put(fieldName, new ImmutablePair<>(PlcResponseCode.OK, new DefaultPlcSubscriptionHandle(this)));
        }
        return CompletableFuture.completedFuture(
            new DefaultPlcSubscriptionResponse((InternalPlcSubscriptionRequest) subscriptionRequest, values));
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> collection) {
        final DefaultPlcConsumerRegistration consumerRegistration =
            new DefaultPlcConsumerRegistration(this, consumer, collection.toArray(new InternalPlcSubscriptionHandle[0]));
        consumerIdMap.put(consumerRegistration.getConsumerHash(), consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration plcConsumerRegistration) {
        DefaultPlcConsumerRegistration consumerRegistration = (DefaultPlcConsumerRegistration) plcConsumerRegistration;
        consumerIdMap.remove(consumerRegistration.getConsumerHash());
    }

    protected void publishEvent(String name, PlcValue plcValue) {
        // Create a subscription event from the input.
        final PlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(),
            Collections.singletonMap(name, Pair.of(PlcResponseCode.OK, plcValue)));

        // Send the subscription event to all listeners.
        for (Consumer<PlcSubscriptionEvent> consumer : consumerIdMap.values()) {
            consumer.accept(event);
        }
    }

    protected static String toString(KNXAddress knxAddress) {
        return knxAddress.getMainGroup() + "." + knxAddress.getMiddleGroup() + "." + knxAddress.getSubGroup();
    }

    protected static String toString(KNXGroupAddress groupAddress) {
        if (groupAddress instanceof KNXGroupAddress3Level) {
            KNXGroupAddress3Level level3 = (KNXGroupAddress3Level) groupAddress;
            return level3.getMainGroup() + "/" + level3.getMiddleGroup() + "/" + level3.getSubGroup();
        } else if (groupAddress instanceof KNXGroupAddress2Level) {
            KNXGroupAddress2Level level2 = (KNXGroupAddress2Level) groupAddress;
            return level2.getMainGroup() + "/" + level2.getSubGroup();
        } else if(groupAddress instanceof KNXGroupAddressFreeLevel) {
            KNXGroupAddressFreeLevel free = (KNXGroupAddressFreeLevel) groupAddress;
            return free.getSubGroup() + "";
        }
        throw new RuntimeException("Unsupported Group Address Type " + groupAddress.getClass().getName());
    }

}
