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
package org.apache.plc4x.java.streampipes.bacnetip;

import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.bacnetip.PassiveBacNetIpDriver;
import org.apache.plc4x.java.bacnetip.connection.PassiveBacNetIpPlcConnection;
import org.apache.plc4x.java.bacnetip.readwrite.*;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.connection.PcapChannelFactory;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.streampipes.bacnetip.config.ConnectWorkerConfig;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketAddress;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketChannelConfig;
import org.apache.plc4x.java.utils.pcapsockets.netty.UdpIpPacketHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.streampipes.connect.adapter.Adapter;
import org.streampipes.connect.adapter.exception.AdapterException;
import org.streampipes.connect.adapter.exception.ParseException;
import org.streampipes.connect.adapter.model.specific.SpecificDataStreamAdapter;
import org.streampipes.connect.container.worker.init.AdapterWorkerContainer;
import org.streampipes.connect.init.AdapterDeclarerSingleton;
import org.streampipes.model.AdapterType;
import org.streampipes.model.connect.adapter.SpecificAdapterStreamDescription;
import org.streampipes.model.connect.guess.GuessSchema;
import org.streampipes.model.schema.EventProperty;
import org.streampipes.model.schema.EventSchema;
import org.streampipes.sdk.builder.PrimitivePropertyBuilder;
import org.streampipes.sdk.builder.adapter.SpecificDataStreamAdapterBuilder;
import org.streampipes.sdk.utils.Datatypes;

import java.io.File;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class BacNetIpAdapter extends SpecificDataStreamAdapter {

    public static final String ID = "http://plc4x.apache.org/streampipes/adapter/bacnetip";
    private static final Logger logger = LoggerFactory.getLogger(BacNetIpAdapter.class);

    private NettyPlcConnection connection;

    private Map<String, Object> event;
    private int numberProperties;

    public BacNetIpAdapter() {
        event = new HashMap<>();
        numberProperties = 0;
    }

    public BacNetIpAdapter(SpecificAdapterStreamDescription adapterDescription) {
        super(adapterDescription);
        event = new HashMap<>();
        numberProperties = 0;
    }

    @Override
    public SpecificAdapterStreamDescription declareModel() {
        SpecificAdapterStreamDescription description = SpecificDataStreamAdapterBuilder.create(ID, "BACnet/IP", "")
            .iconUrl("bacnetip.png")
            .category(AdapterType.Manufacturing)
            .build();
        description.setAppId(ID);
        return description;
    }

    @Override
    public GuessSchema getSchema(SpecificAdapterStreamDescription specificAdapterStreamDescription) throws AdapterException, ParseException {
        EventSchema eventSchema = new EventSchema();
        List<EventProperty> allProperties = new ArrayList<>();

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.String, "sourceId")
                .label("Source Id")
                .description("")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.String, "propertyId")
                .label("Property Id")
                .description("")
                .build());

        // We need to define the type of the value, I choose a numerical value
        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Float, "value")
                .label("Value")
                .description("")
                .build());

        eventSchema.setEventProperties(allProperties);

        GuessSchema guessSchema = new GuessSchema();
        guessSchema.setEventSchema(eventSchema);
        guessSchema.setPropertyProbabilityList(new ArrayList<>());

        return guessSchema;
    }

    @Override
    public void startAdapter() throws AdapterException {
        try {
            connection = new PassiveBacNetIpPlcConnection(new PcapChannelFactory(
                //new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/BacNET/Captures/Merck/BACnetWhoIsRouterToNetwork.pcapng"), null,
                new File("/Users/christofer.dutz/Downloads/20190906_udp.pcapng"), null,
                PassiveBacNetIpDriver.BACNET_IP_PORT, PcapSocketAddress.ALL_PROTOCOLS,
                PcapSocketChannelConfig.NO_THROTTLING, new UdpIpPacketHandler()), "",
                new PlcMessageToMessageCodec<BVLC, PlcRequestContainer>() {

                @Override
                protected void decode(ChannelHandlerContext channelHandlerContext, BVLC packet, List<Object> list) throws Exception {
                    final NPDU npdu = ((BVLCOriginalUnicastNPDU) packet).getNpdu();
                    final APDU apdu = npdu.getApdu();
                    if(apdu instanceof APDUConfirmedRequest) {
                        APDUConfirmedRequest request = (APDUConfirmedRequest) apdu;
                        final BACnetConfirmedServiceRequest serviceRequest = request.getServiceRequest();
                        if(serviceRequest instanceof BACnetConfirmedServiceRequestConfirmedCOVNotification) {
                            BACnetConfirmedServiceRequestConfirmedCOVNotification covNotification = (BACnetConfirmedServiceRequestConfirmedCOVNotification) serviceRequest;
                            final BACnetTagWithContent[] notifications = covNotification.getNotifications();

                            // TODO: Get the information from the decoded packet.
                            String key = ""; // Node-id + property-id
                            Float value = 1.0f; // Value

                            event.put(key, value);
                            if (event.keySet().size() >= numberProperties) {
                                adapterPipeline.process(event);
                            }

                            System.out.println("Simple-ACK(" + request.getInvokeId() + "): Confirmed COV Notification [" + notifications.length + "]");
                        }
                    }
                }

                @Override
                protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
                    // Ignore this as we don't send anything.
                }
            });
            connection.connect();
        } catch (PlcConnectionException e) {
            logger.error("An error occurred starting the BACnet/IP driver", e);
            throw new AdapterException("An error occurred starting the BACnet/IP driver");
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
    public Adapter getInstance(SpecificAdapterStreamDescription specificAdapterStreamDescription) {
        return new BacNetIpAdapter(adapterDescription);
    }

    @Override
    public String getId() {
        return ID;
    }

    public static class BacNetIpAdapterInit extends AdapterWorkerContainer {
        public static void main(String[] args) {
            AdapterDeclarerSingleton
                .getInstance()
                .add(new BacNetIpAdapter());

            String workerUrl = ConnectWorkerConfig.INSTANCE.getConnectContainerWorkerUrl();
            String masterUrl = ConnectWorkerConfig.INSTANCE.getConnectContainerMasterUrl();
            Integer workerPort = ConnectWorkerConfig.INSTANCE.getConnectContainerWorkerPort();

            new BacNetIpAdapterInit().init(workerUrl, masterUrl, workerPort);
        }
    }

}
