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
import org.pcap4j.core.*;
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
import org.streampipes.model.staticproperty.FileStaticProperty;
import org.streampipes.sdk.StaticProperties;
import org.streampipes.sdk.builder.PrimitivePropertyBuilder;
import org.streampipes.sdk.builder.adapter.SpecificDataStreamAdapterBuilder;
import org.streampipes.sdk.helpers.*;
import org.streampipes.sdk.utils.Datatypes;

import java.io.File;
import java.nio.ByteBuffer;
import java.util.*;

public class BacNetIpAdapter extends SpecificDataStreamAdapter {

    public static final String ID = "http://plc4x.apache.org/streampipes/adapter/bacnetip";

    private static final Logger logger = LoggerFactory.getLogger(BacNetIpAdapter.class);

    private NettyPlcConnection connection;

    public BacNetIpAdapter() {
        super();
    }

    public BacNetIpAdapter(SpecificAdapterStreamDescription adapterDescription) {
        super(adapterDescription);
    }

    @Override
    public SpecificAdapterStreamDescription declareModel() {
        Label fileLabel = Labels.from("pcap-file", "PCAP File", "File containing the network traffic recording");

        Tuple2<String, String>[] deviceList = null;
        try {
            final List<PcapNetworkInterface> allDevs = Pcaps.findAllDevs();
            deviceList = new Tuple2[allDevs.size()];
            for (int i = 0; i < allDevs.size(); i++) {
                final PcapNetworkInterface pcapNetworkInterface = allDevs.get(i);
                StringBuilder deviceName = new StringBuilder((pcapNetworkInterface.getDescription() != null) ? pcapNetworkInterface.getDescription() : pcapNetworkInterface.getName());
                deviceName.append(" (");
                for (PcapAddress address : pcapNetworkInterface.getAddresses()) {
                    if(address instanceof PcapIpV4Address) {
                        deviceName.append(address.getAddress().toString()).append("/").append(address.getNetmask().toString()).append(", ");
                    }
                }
                String name = deviceName.toString();
                name = name.substring(0, name.length() - 2) + ((name.endsWith(", ")) ? ")": "");
                deviceList[i] = new Tuple2<>(pcapNetworkInterface.getName(), name);
            }
        } catch (PcapNativeException e) {
            logger.error("Error getting the list of installed network devices");
        }

        SpecificAdapterStreamDescription description = SpecificDataStreamAdapterBuilder.create(ID, "BACnet/IP", "")
            .iconUrl("bacnetip.png")
            .category(AdapterType.Manufacturing)
            .requiredAlternatives(Labels.from("source", "Source", "Select the source, where the data is read from"),
                Alternatives.from(Labels.from("device", "Network", "Capture data via network device"),
                    StaticProperties.group(Labels.withId("device-group"),
                        StaticProperties.singleValueSelection(Labels.from("network-device", "Network Device", "Network device used for capturing"),
                            Options.from(deviceList)))),
                Alternatives.from(Labels.from("file", "File", "Capture data from a PCAP network recording"),
                    StaticProperties.group(Labels.withId("file-group"),
                        new FileStaticProperty(fileLabel.getInternalId(), fileLabel.getLabel(), fileLabel.getDescription()))))
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
                .create(Datatypes.Long, "time")
                .label("Time")
                .description("The time the event was processed in the BACnet adapter")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, "objectType")
                .label("Object type")
                .description("Type of BACnet object emitting the event (usually 'device')")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, "objectId")
                .label("Object id")
                .description("Id of the BACnet object emitting the event (usually 'device id')")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, "notificationType")
                .label("Notification type")
                .description("The type of notification this event resembles (usually some type of input)")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, "notificationInstanceNumber")
                .label("Notification instance number")
                .description("The instance number of the component emitting the event (usually the id of the property changed on a device)")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Integer, "valueType")
                .label("Value type")
                .description("The type the value has (real, uint, int, bit-string, ...)")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.String, "value")
                .label("Value")
                .description("This is the actual payload of the event.")
                .build());

        allProperties.add(
            PrimitivePropertyBuilder
                .create(Datatypes.Sequence, "status")
                .label("Status")
                .description("Some times an array of status bits are passed along.")
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
                new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/BacNET/Captures/Merck/BACnet.pcapng"), null,
                PassiveBacNetIpDriver.BACNET_IP_PORT, PcapSocketAddress.ALL_PROTOCOLS,
                PcapSocketChannelConfig.SPEED_REALTIME, new UdpIpPacketHandler()), "",
                new PlcMessageToMessageCodec<BVLC, PlcRequestContainer>() {

                @Override
                protected void decode(ChannelHandlerContext channelHandlerContext, BVLC packet, List<Object> list) throws Exception {
                    NPDU npdu = null;
                    if(packet instanceof BVLCOriginalUnicastNPDU) {
                        npdu = ((BVLCOriginalUnicastNPDU) packet).getNpdu();
                    } else if(packet instanceof BVLCForwardedNPDU) {
                        npdu = ((BVLCForwardedNPDU) packet).getNpdu();
                    } else if(packet instanceof BVLCOriginalBroadcastNPDU) {
                        npdu = ((BVLCOriginalBroadcastNPDU) packet).getNpdu();
                    } else {
                        throw new RuntimeException("Unexpected type of packet");
                    }
                    final APDU apdu = npdu.getApdu();
                    if(apdu instanceof APDUConfirmedRequest) {
                        APDUConfirmedRequest request = (APDUConfirmedRequest) apdu;
                        final BACnetConfirmedServiceRequest serviceRequest = request.getServiceRequest();
                        if(serviceRequest instanceof BACnetConfirmedServiceRequestConfirmedCOVNotification) {
                            BACnetConfirmedServiceRequestConfirmedCOVNotification covNotification = (BACnetConfirmedServiceRequestConfirmedCOVNotification) serviceRequest;
                            final BACnetTagWithContent[] notifications = covNotification.getNotifications();

                            String objectType = Integer.toString(covNotification.getMonitoredObjectType());
                            String objectId = Long.toString(covNotification.getMonitoredObjectInstanceNumber());

                            String notificationType = Integer.toString(covNotification.getIssueConfirmedNotificationsType());
                            String notificationInstanceNumber = Long.toString(covNotification.getMonitoredObjectInstanceNumber());

                            String type = null;
                            Object value = null;
                            boolean[] status = null;
                            for (BACnetTagWithContent notification : notifications) {
                                // Id of the property that changed
                                short propertyId = notification.getPropertyIdentifier()[0];

                                // Present-Value has the property id 85
                                // (This is the actual value to which a given property has changed)
                                if(propertyId == 85) {
                                    // Depending on the type of object, parse the data accordingly.
                                    if (notification.getVal() instanceof BACnetTagApplicationBoolean) {
                                        type = "boolean";
                                        final BACnetTagApplicationBoolean val = (BACnetTagApplicationBoolean) notification.getVal();

                                    } else if (notification.getVal() instanceof BACnetTagApplicationUnsignedInteger) {
                                        type = "uint";
                                        final BACnetTagApplicationUnsignedInteger val = (BACnetTagApplicationUnsignedInteger) notification.getVal();
                                        // Convert any number of bytes into an unsigned integer.
                                        switch (val.getData().length) {
                                            case 1:
                                                value = Byte.toString(val.getData()[0]);
                                                break;
                                            case 2:
                                                value = Short.toString(ByteBuffer.wrap(val.getData()).getShort());
                                                break;
                                            case 3:
                                                byte[] extValues = new byte[4];
                                                extValues[0] = 0x00;
                                                for(int i = 0; i < 3; i++) {
                                                    extValues[i+1] = val.getData()[i];
                                                }
                                                value = ByteBuffer.wrap(extValues).getInt();
                                                break;
                                            default:
                                                value = "Hurz";
                                                break;
                                        }
                                    } else if (notification.getVal() instanceof BACnetTagApplicationSignedInteger) {
                                        type = "int";
                                        final BACnetTagApplicationSignedInteger val = (BACnetTagApplicationSignedInteger) notification.getVal();

                                    } else if (notification.getVal() instanceof BACnetTagApplicationReal) {
                                        type = "real";
                                        final BACnetTagApplicationReal val = (BACnetTagApplicationReal) notification.getVal();
                                        value = Float.intBitsToFloat(ByteBuffer.wrap(val.getData()).getInt());
                                    } else if (notification.getVal() instanceof BACnetTagApplicationDouble) {
                                        type = "double";
                                        final BACnetTagApplicationDouble val = (BACnetTagApplicationDouble) notification.getVal();

                                    } else if (notification.getVal() instanceof BACnetTagApplicationBitString) {
                                        type = "bit-string";
                                        final BACnetTagApplicationBitString val = (BACnetTagApplicationBitString) notification.getVal();
                                        int numBits = (val.getData().length * 8) - val.getUnusedBits();
                                        BitSet bitSet = BitSet.valueOf(val.getData());
                                        boolean[] bits = new boolean[numBits];
                                        for (int i = 0; i < numBits; i++) {
                                            bits[i] = bitSet.get(i);
                                        }
                                        value = bits;
                                    } else if (notification.getVal() instanceof BACnetTagApplicationEnumerated) {
                                        type = "enumeration";
                                        final BACnetTagApplicationEnumerated val = (BACnetTagApplicationEnumerated) notification.getVal();

                                    }
                                }

                                // Status-Flags have the property id 111
                                // (This is some additional information passed along)
                                else if(propertyId == 111) {
                                    final BACnetTagApplicationBitString val = (BACnetTagApplicationBitString) notification.getVal();
                                    int numBits = (val.getData().length * 8) - val.getUnusedBits();
                                    BitSet bitSet = BitSet.valueOf(val.getData());
                                    boolean[] bits = new boolean[numBits];
                                    for (int i = 0; i < numBits; i++) {
                                        bits[i] = bitSet.get(i);
                                    }
                                    status = bits;
                                }
                            }

                            if(value != null) {
                                // Create the event object.
                                Map<String, Object> event = new HashMap<>();
                                event.put("time", System.currentTimeMillis());

                                event.put("objectType", objectType);
                                event.put("objectId", objectId);
                                event.put("notificationType", notificationType);
                                event.put("notificationInstanceNumber", notificationInstanceNumber);

                                event.put("valueType", type);
                                event.put("value", value);
                                event.put("status", status);

                                // Send it to StreamPipes
                                adapterPipeline.process(event);
                            }
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
        return new BacNetIpAdapter(specificAdapterStreamDescription);
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
