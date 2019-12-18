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
package org.apache.plc4x.java.streampipes.adapters.source.knxnetip;

import io.netty.channel.ChannelHandlerContext;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.base.connection.UdpSocketChannelFactory;
import org.apache.plc4x.java.spi.PlcMessageToMessageCodec;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.NettyPlcConnection;
import org.apache.plc4x.java.spi.messages.PlcRequestContainer;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpConnection;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.apache.plc4x.java.streampipes.shared.source.knxnetip.Constants;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.streampipes.connect.adapter.Adapter;
import org.streampipes.connect.adapter.exception.AdapterException;
import org.streampipes.connect.adapter.exception.ParseException;
import org.streampipes.connect.adapter.model.specific.SpecificDataStreamAdapter;
import org.streampipes.connect.adapter.sdk.ParameterExtractor;
import org.streampipes.model.AdapterType;
import org.streampipes.model.connect.adapter.SpecificAdapterStreamDescription;
import org.streampipes.model.connect.guess.GuessSchema;
import org.streampipes.model.schema.EventProperty;
import org.streampipes.model.schema.EventSchema;
import org.streampipes.model.staticproperty.FreeTextStaticProperty;
import org.streampipes.sdk.builder.PrimitivePropertyBuilder;
import org.streampipes.sdk.builder.adapter.SpecificDataStreamAdapterBuilder;
import org.streampipes.sdk.helpers.Labels;
import org.streampipes.sdk.utils.Datatypes;
import org.streampipes.vocabulary.SO;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class KnxNetIpAdapter extends SpecificDataStreamAdapter {

    public static final String MAPPING_FIELD_TIME = "time";
    public static final String MAPPING_FIELD_SOURCE_ADDRESS = "sourceAddress";
    public static final String MAPPING_FIELD_DESTINATION_ADDRESS = "destinationAddress";
    public static final String MAPPING_FIELD_PAYLOAD = "payload";

    private static final Logger logger = LoggerFactory.getLogger(KnxNetIpAdapter.class);

    private String gatewayIp;
    private NettyPlcConnection connection;

    public KnxNetIpAdapter() {
        super();
    }

    public KnxNetIpAdapter(SpecificAdapterStreamDescription adapterDescription, String gatewayIp) {
        super(adapterDescription);
        this.gatewayIp = gatewayIp;
    }

    @Override
    public SpecificAdapterStreamDescription declareModel() {
        SpecificAdapterStreamDescription description = SpecificDataStreamAdapterBuilder.create(Constants.KNXNET_ID, "KNXnet/IP", "")
            .iconUrl("knxnetip.png")
            .category(AdapterType.Manufacturing)
            .requiredTextParameter(Labels.from("gatewayIp", "KNXnet/IP Gateway", "Ip of the KNX gateway."))
            .build();
        description.setAppId(Constants.KNXNET_ID);
        return description;
    }

    @Override
    public GuessSchema getSchema(SpecificAdapterStreamDescription specificAdapterStreamDescription) throws AdapterException, ParseException {
        List<EventProperty> allProperties = new ArrayList<>();
        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Long, MAPPING_FIELD_TIME)
                .domainProperty(SO.DateTime)
                .label("Time")
                .description("The time the event was processed in the KNXnet/IP driver.")
                .build());
        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, MAPPING_FIELD_SOURCE_ADDRESS)
                .domainProperty(Constants.KNXNET_ID_SOURCE_ADDRESS)
                .label("Source Address")
                .description("Source address from which the event originated.")
                .build());
        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, MAPPING_FIELD_DESTINATION_ADDRESS)
                .domainProperty(Constants.KNXNET_ID_DESTINATION_ADDRESS)
                .label("Destination Address")
                .description("Destination address to which the event is targeted.")
                .build());
        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.String, MAPPING_FIELD_PAYLOAD)
                .domainProperty(Constants.KNXNET_ID_PAYLOAD)
                .label("Payload")
                .description("Raw payload of the event.")
                .build());

        EventSchema eventSchema = new EventSchema();
        eventSchema.setEventProperties(allProperties);

        GuessSchema guessSchema = new GuessSchema();
        guessSchema.setEventSchema(eventSchema);
        guessSchema.setPropertyProbabilityList(new ArrayList<>());

        return guessSchema;
    }

    @Override
    public void startAdapter() throws AdapterException {
        if((connection != null) && (connection.isConnected())) {
            return;
        }
        try {
            InetAddress inetAddress = InetAddress.getByName(gatewayIp);
            ChannelFactory channelFactory = new UdpSocketChannelFactory(inetAddress, KnxNetIpConnection.KNXNET_IP_PORT);

            connection = new KnxNetIpConnection(channelFactory, "",
                new PlcMessageToMessageCodec<KNXNetIPMessage, PlcRequestContainer>() {

                @Override
                protected void decode(ChannelHandlerContext channelHandlerContext, KNXNetIPMessage packet, List<Object> list) throws Exception {
                    if(packet instanceof TunnelingRequest) {
                        TunnelingRequest request = (TunnelingRequest) packet;
                        CEMI cemiPayload = request.getCemi();
                        if(cemiPayload instanceof CEMIBusmonInd) {
                            CEMIBusmonInd cemiBusmonInd = (CEMIBusmonInd) cemiPayload;
                            if(cemiBusmonInd.getCemiFrame() instanceof CEMIFrameData) {
                                CEMIFrameData cemiDataFrame = (CEMIFrameData) cemiBusmonInd.getCemiFrame();

                                // The first byte is actually just 6 bit long, but we'll treat it as a full one.
                                // So here we create a byte array containing the first and all the following bytes.
                                byte[] payload = new byte[1 + cemiDataFrame.getData().length];
                                payload[0] = cemiDataFrame.getDataFirstByte();
                                System.arraycopy(cemiDataFrame.getData(), 0, payload, 1, cemiDataFrame.getData().length);

                                Map<String, Object> event = new HashMap<>();
                                event.put(MAPPING_FIELD_TIME, System.currentTimeMillis());
                                event.put(MAPPING_FIELD_SOURCE_ADDRESS, addressToString(cemiDataFrame.getSourceAddress()));
                                event.put(MAPPING_FIELD_DESTINATION_ADDRESS, Hex.encodeHex(cemiDataFrame.getDestinationAddress()));
                                // Encode the payload as Hex String.
                                event.put(MAPPING_FIELD_PAYLOAD, Hex.encodeHexString(payload));

                                // Send it to StreamPipes
                                adapterPipeline.process(event);
                            } else if (cemiBusmonInd.getCemiFrame() instanceof CEMIFrameAck){
                                // Just ignore this ...
                            } else {
                                System.out.println(packet);
                            }
                        } else {
                            System.out.println(packet);
                        }
                    } else {
                        System.out.println(packet);
                    }
                }

                @Override
                protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) {
                    // Ignore this as we don't send anything.
                }
            });
            connection.connect();
        } catch (PlcConnectionException e) {
            logger.error("An error occurred starting the BACnet/IP driver", e);
            throw new AdapterException("An error occurred starting the BACnet/IP driver");
        } catch (UnknownHostException e) {
            logger.error("Error connecting to host " + gatewayIp, e);
            throw new AdapterException("Error connecting to host " + gatewayIp);
        } catch (Exception e) {
            logger.error("Something strange went wrong.", e);
        }
    }

    @Override
    public void stopAdapter() throws AdapterException {
        if(connection != null) {
            try {
                connection.close();
            } catch (PlcConnectionException e) {
                logger.error("An error occurred stopping the BACnet/IP driver", e);
                throw new AdapterException("An error occurred stopping the BACnet/IP driver");
            }
        }
    }

    @Override
    @SuppressWarnings("unchecked")
    public Adapter getInstance(SpecificAdapterStreamDescription specificAdapterStreamDescription) {
        ParameterExtractor extractor = new ParameterExtractor(specificAdapterStreamDescription.getConfig());

        FreeTextStaticProperty gatewayIpProperty = (FreeTextStaticProperty) extractor.getStaticPropertyByName("gatewayIp");
        gatewayIp = gatewayIpProperty.getValue();
        return new KnxNetIpAdapter(specificAdapterStreamDescription, gatewayIp);
    }

    @Override
    public String getId() {
        return Constants.KNXNET_ID;
    }

    private String addressToString(KNXAddress knxAddress) {
        return knxAddress.getMainGroup() + "." + knxAddress.getMiddleGroup() + "." + knxAddress.getSubGroup();
    }

}
