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
package org.apache.plc4x.java.eip.base.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.eip.base.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.logix.configuration.LogixConfiguration;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EipProtocolLogic extends Plc4xProtocolBase<EipPacket> implements HasConfiguration<EIPConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(EipProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static final byte[] DEFAULT_SENDER_CONTEXT = "PLC4X   ".getBytes(StandardCharsets.US_ASCII);
    private static final long EMPTY_SESSION_HANDLE = 0L;
    private static final long EMPTY_INTERFACE_HANDLE = 0L;
    private NullAddressItem nullAddressItem;
    private byte[] senderContext;
    private long connectionId = 0L;
    private int sequenceCount = 1;
    private EIPConfiguration configuration;

    private final AtomicInteger transactionCounterGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;
    private long sessionHandle;

    private boolean useConnectionManager = false;

    private boolean cipEncapsulationAvailable = false;

    private boolean useMessageRouter = false;

    private final List<PathSegment> routingAddress = new ArrayList<>();
    short connectionPathSize = 0;
    private final int connectionSerialNumber = ThreadLocalRandom.current().nextInt();

    @Override
    public void setConfiguration(EIPConfiguration configuration) {
        this.configuration = configuration;
        this.nullAddressItem = new NullAddressItem();

        if (configuration instanceof LogixConfiguration) {
            // Use a connection path instead of the backplane and slot if it is available
            LogixConfiguration logixConfiguration = (LogixConfiguration) configuration;
            if (logixConfiguration.getCommunicationPath() != null) {
                String[] splitConnectionPath = logixConfiguration.getCommunicationPath().split(",");
                if (splitConnectionPath.length % 2 == 0) {
                    for (int i = 0; (i + 1) < splitConnectionPath.length; i += 2 ) {
                        switch(splitConnectionPath[i]) {
                            case "1":
                                int backplanePortId = Integer.parseInt(splitConnectionPath[i]);
                                int slot = Integer.parseInt(splitConnectionPath[i + 1]);
                                routingAddress.add(new PortSegment(new PortSegmentNormal((byte) backplanePortId, (short) slot)));
                                break;
                            case "2":
                                int ethernetPortId = Integer.parseInt(splitConnectionPath[i]);
                                String ipAddress = splitConnectionPath[i+1];
                                int lengthString = ipAddress.length();

                                if ((ipAddress.length() % 2) != 0) {
                                    ipAddress += "\0";
                                }

                                routingAddress.add(new PortSegment(new PortSegmentExtended((byte) ethernetPortId, (short) lengthString, ipAddress)));
                                break;
                            default:
                                logger.error("Only backplane or Ethernet module routing is supported");
                        }

                    }
                }
            } else {
                routingAddress.add(new PortSegment(new PortSegmentNormal((byte) 1, (short) this.configuration.getSlot())));
            }
        } else {
            routingAddress.add(new PortSegment(new PortSegmentNormal((byte) 1, (short)  this.configuration.getSlot())));
        }

        routingAddress.add(new LogicalSegment(new ClassID((byte) 0, (short) 2)));
        routingAddress.add(new LogicalSegment(new InstanceID((byte) 0, (short) 1)));

        for (PathSegment segment : this.routingAddress) {
            this.connectionPathSize += segment.getLengthInBytes();
        }
        if ((this.connectionPathSize % 2) != 0) {
            this.connectionPathSize += 1;
        }

        this.connectionPathSize = (short) (this.connectionPathSize / 2);

        // Set the transaction manager to allow only one message at a time.
        this.tm = new RequestTransactionManager(1);
    }

    public CompletableFuture<Boolean> detectEndianness(ConversationContext<EipPacket> context) {
        logger.debug("Sending Unknown Command to determine Endianess");
        CompletableFuture<Boolean> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        NullCommandRequest listServicesRequest = new NullCommandRequest(
            EMPTY_SESSION_HANDLE,
            CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT,
            0L);

        transaction.submit(() -> context.sendRequest(listServicesRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT).unwrap(p -> p)
            .onError((p,e) -> logger.warn("No response for initial packet. Suspect device uses Big endian"))
            .check(p -> p instanceof NullCommandRequest)
            .handle(p -> {
                logger.info("Device uses little endian");
                future.complete(true);
                transaction.endRequest();
            })
        );
        return future;
    }

    private void listServices(ConversationContext<EipPacket> context) {
        logger.debug("Sending List Services packet to confirm CIP Encapsulation is available");

        // TODO: It seems that we're only doing this request in order to find out, if we can do CIP encapsulation, however this value is never really used anywhere.
        ListServicesRequest listServicesRequest = new ListServicesRequest(
            EMPTY_SESSION_HANDLE,
            CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT,
            0L);

        context.sendRequest(listServicesRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .onError((p,e) -> {throw new PlcRuntimeException("List EIP Services failed");})
            .handle(p -> {
                if (p.getStatus() == CIPStatus.Success.getValue()) {
                    ServicesResponse listServicesResponse = (ServicesResponse) ((ListServicesResponse) p).getTypeIds().get(0);
                    if (listServicesResponse.getSupportsCIPEncapsulation()) {
                        logger.debug("Device is capable of CIP over EIP encapsulation");
                    }
                    this.cipEncapsulationAvailable = listServicesResponse.getSupportsCIPEncapsulation();
                } else if (p.getStatus() == CIPStatus.InvalidCommandWithWrongEndianess.getValue()) {
                    throw new PlcRuntimeException("The remote device doesn't seem to use " + configuration.getByteOrder().name() + " byte order.");
                } else {
                    throw new PlcRuntimeException("Got status code while polling for supported EIP services [" + p.getStatus() + "]");
                }
                onConnectRegisterSession(context);
            });
    }

    private void getAllAttributes(ConversationContext<EipPacket> context) {
        logger.debug("Requesting list of supported attributes");

        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 2));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

        UnConnectedDataItem exchange = new UnConnectedDataItem(
            new GetAttributeAllRequest(
                classSegment,
                instanceSegment)
        );

        List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

        CipRRData eipWrapper = new CipRRData(
            sessionHandle,
            CIPStatus.Success.getValue(),
            senderContext,
            0L,
            EMPTY_INTERFACE_HANDLE,
            0,
            typeIds
        );

        context.sendRequest(eipWrapper)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT).unwrap(p -> p)
            .check(p -> p instanceof CipRRData)
            .handle(p -> {
                if (p.getStatus() == CIPStatus.Success.getValue()) {
                    UnConnectedDataItem dataItem = (UnConnectedDataItem) ((CipRRData) p).getTypeIds().get(1);
                    GetAttributeAllResponse response = (GetAttributeAllResponse) dataItem.getService();
                    if ((long) response.getStatus() == CIPStatus.ServiceNotSupported.getValue()) {
                        context.fireConnected();
                        return;
                    } else if ( (long) response.getStatus() != CIPStatus.Success.getValue()) {
                        throw new PlcRuntimeException("Got status code while polling for supported CIP attributes [" + response.getStatus() + "]");
                    }
                    if (response.getAttributes() != null) {
                        for (Integer classId : response.getAttributes().getClassId()) {
                            if (CIPClassID.enumForValue(classId) == CIPClassID.MessageRouter) {
                                this.useMessageRouter = true;
                            }
                            if (CIPClassID.enumForValue(classId) == CIPClassID.ConnectionManager) {
                                this.useConnectionManager = true;
                            }
                        }
                    }

                    if (this.useConnectionManager) {
                        logger.debug("Device is using a Connection Manager");
                        onConnectOpenConnectionManager(context);
                    } else {
                        context.fireConnected();
                    }
                } else {
                    throw new PlcRuntimeException("Got status code while polling for supported CIP services [" + p.getStatus() + "]");
                }
            });
    }

    @Override
    public void onConnect(ConversationContext<EipPacket> context) {
        listServices(context);
    }

    private void onConnectRegisterSession(ConversationContext<EipPacket> context) {
        logger.debug("Sending Register Session EIP Package");

        EipConnectionRequest connectionRequest =
            new EipConnectionRequest(
                EMPTY_SESSION_HANDLE,
                CIPStatus.Success.getValue(),
                DEFAULT_SENDER_CONTEXT,
                0L);

        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .handle(p -> {
                if (p instanceof EipConnectionResponse) {
                    if (p.getStatus() == CIPStatus.Success.getValue()) {
                        sessionHandle = p.getSessionHandle();
                        senderContext = p.getSenderContext();
                        logger.debug("Got assigned with Session handle {}", sessionHandle);
                        getAllAttributes(context);
                    }  else {
                        throw new PlcRuntimeException("Got status code while polling for supported EIP services [" + p.getStatus() + "]");
                    }
                } else {
                    onConnectOpenConnectionManager(context);
                }
            });
    }

    public void onConnectOpenConnectionManager(ConversationContext<EipPacket> context) {
        logger.debug("Sending Open Connection Manager EIP Package");

        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

        UnConnectedDataItem exchange = new UnConnectedDataItem(
            new CipConnectionManagerRequest(
                classSegment,
                instanceSegment,
                (byte) 0,
                (byte) 10,
                (short) 14,
                536870914L,
                33944L,
                this.connectionSerialNumber,
                4919,
                42L,
                (short) 3,
                2101812L,
                new NetworkConnectionParameters(
                    4002,
                    false,
                    (byte) 2,
                    (byte) 0,
                    true),
                2113537L,
                new NetworkConnectionParameters(
                    4002,
                    false,
                    (byte) 2,
                    (byte) 0,
                    true),
                new TransportType(true, (byte) 2, (byte) 3),
                this.connectionPathSize,
                this.routingAddress)
        );

        List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

        CipRRData eipWrapper = new CipRRData(
            sessionHandle,
            CIPStatus.Success.getValue(),
            senderContext,
            0L,
            EMPTY_INTERFACE_HANDLE,
            0,
            typeIds
        );

        context.sendRequest(eipWrapper)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .check(p -> p instanceof CipRRData)
            .unwrap(p -> (CipRRData) p)
            .handle(p -> {
                if (p.getStatus() == 0L) {
                    CipRRData rrData = (CipRRData) p;
                    List<TypeId> connectionManagerExchange = rrData.getTypeIds();
                    CipConnectionManagerResponse connectionManagerResponse = (CipConnectionManagerResponse) ((UnConnectedDataItem) connectionManagerExchange.get(1)).getService();
                    this.connectionId = connectionManagerResponse.getOtConnectionId();

                    logger.debug("Got assigned with Connection Id {}", this.connectionId);
                    // Send an event that connection setup is complete.
                    context.fireConnected();
                } else {
                    throw new PlcRuntimeException("Got status code while opening Connection Manager[" + p.getStatus() + "]");
                }
            });
    }


    @Override
    public void onDisconnect(ConversationContext<EipPacket> context) {
        if (this.connectionId != 0L) {
            logger.debug("Sending Connection Manager Close Event");
            PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
            PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipConnectionManagerCloseRequest(
                    (byte) 2,
                    classSegment,
                    instanceSegment,
                    (byte) 0,
                    (byte) 10,
                    (short) 14,
                    this.connectionSerialNumber,
                    4919,
                    42L,
                    this.connectionPathSize,
                    this.routingAddress));

            List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

            CipRRData eipWrapper = new CipRRData(
                sessionHandle,
                0L,
                senderContext,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                typeIds);


            context.sendRequest(eipWrapper)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT).unwrap(p -> p)
                .check(p -> p instanceof CipRRData)
                .handle(p -> {
                    logger.debug("Un-Registering Session");
                    onDisconnectUnregisterSession(context);
                });
        } else {
            onDisconnectUnregisterSession(context);
        }
    }


    public void onDisconnectUnregisterSession(ConversationContext<EipPacket> context) {
        logger.debug("Sending Un RegisterSession EIP Package");

        EipDisconnectRequest connectionRequest =
            new EipDisconnectRequest(
                sessionHandle,
                0L,
                DEFAULT_SENDER_CONTEXT,
                0L);
        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, Duration.ofMillis(1))
            .onError((p,e) -> context.fireDisconnected())
            .onTimeout(p -> context.fireDisconnected())
            .handle(p -> context.fireDisconnected());

        context.fireDisconnected();
    }

    private CompletableFuture<PlcReadResponse> readWithoutMessageRouter(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();

        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        for (String tagName : request.getTagNames()) {
            CompletableFuture<Boolean> internalFuture = new CompletableFuture<>();
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
            EipTag eipTag = (EipTag) request.getTag(tagName);
            String tag = eipTag.getTag();

            try {
                CipReadRequest req = new CipReadRequest(
                    toAnsi(tag),
                    1);

                CipUnconnectedRequest requestItem = new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    req,
                    (byte) this.configuration.getBackplane(),
                    (byte) this.configuration.getSlot());

                List<TypeId> typeIds = new ArrayList<>(2);

                typeIds.add(nullAddressItem);
                typeIds.add(new UnConnectedDataItem(requestItem));

                CipRRData pkt = new CipRRData(
                    sessionHandle,
                    CIPStatus.Success.getValue(),
                    DEFAULT_SENDER_CONTEXT,
                    0L,
                    0L,
                    0,
                    typeIds);

                transaction.submit(() -> context.sendRequest(pkt)
                    .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                    .onTimeout(future::completeExceptionally)
                    .onError((p, e) -> future.completeExceptionally(e))
                    .check(p -> p instanceof CipRRData)
                    .unwrap(p -> (CipRRData) p)
                    .check(p -> p.getSessionHandle() == sessionHandle)
                    .handle(p -> {
                        List<TypeId> responseTypeIds = p.getTypeIds();
                        UnConnectedDataItem dataItem = (UnConnectedDataItem) responseTypeIds.get(1);
                        Map<String, ResponseItem<PlcValue>> readResponse = decodeSingleReadResponse(dataItem.getService(), tagName, eipTag);
                        values.putAll(readResponse);
                        internalFuture.complete(true);
                        // Finish the request-transaction.
                        transaction.endRequest();
                    }));
                // TODO: Remove this ...
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    throw new RuntimeException(e);
                }
            } catch (SerializationException e) {
                e.printStackTrace();
            }
        }

        // TODO: This seems to be blocking here ... we should probably do this asynchronously
        PlcReadResponse readResponse = new DefaultPlcReadResponse(readRequest, values);
        future.complete(readResponse);

        return future;
    }

    private CompletableFuture<PlcReadResponse> readWithoutConnectionManager(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<CipService> requests = new ArrayList<>(request.getNumberOfTags());
        for (PlcTag field : request.getTags()) {
            EipTag plcField = (EipTag) field;
            String tag = plcField.getTag();

            try {
                CipReadRequest req = new CipReadRequest(
                    toAnsi(tag),
                    1);

                CipUnconnectedRequest unreq = new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    req,
                    (byte) this.configuration.getBackplane(),
                    (byte) this.configuration.getSlot());
                requests.add(unreq);
            } catch (SerializationException e) {
                e.printStackTrace();
            }
        }

        List<TypeId> typeIds =new ArrayList<>(2);

        typeIds.add(nullAddressItem);
        if (requests.size() == 1) {
            typeIds.add(new UnConnectedDataItem(requests.get(0)));
        } else {
            List<Integer> offsets = new ArrayList<>(requests.size());
            offsets.add(8);
            for (CipService cipRequest : requests) {
                if (requests.indexOf(cipRequest) != (requests.size() - 1)) {
                    offsets.add(offsets.get(requests.indexOf(cipRequest)) + cipRequest.getLengthInBytes());
                }

            }
            MultipleServiceRequest serviceRequest = new MultipleServiceRequest(new Services(offsets, requests));
            typeIds.add(new UnConnectedDataItem(serviceRequest));
        }

        CipRRData pkt = new CipRRData(
            sessionHandle,
            CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT,
            0L,
            0L,
            0,
            typeIds
        );

        transaction.submit(() -> context.sendRequest(pkt)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p instanceof CipRRData)
            .unwrap(p -> (CipRRData) p)
            .check(p -> p.getSessionHandle() == sessionHandle)
            .handle(p -> {
                List<TypeId> responseTypeIds = p.getTypeIds();
                UnConnectedDataItem dataItem = (UnConnectedDataItem) responseTypeIds.get(1);
                PlcReadResponse readResponse = decodeReadResponse(dataItem.getService(), request);
                future.complete(readResponse);
                // Finish the request-transaction.
                transaction.endRequest();
            }));

        return future;
    }

    private CompletableFuture<PlcReadResponse> readWithConnectionManager(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<CipService> requests = new ArrayList<>(request.getNumberOfTags());

        for (PlcTag field : request.getTags()) {
            EipTag plcField = (EipTag) field;
            String tag = plcField.getTag();

            try {
                CipReadRequest req = new CipReadRequest(
                    toAnsi(tag),
                    1);

                requests.add(req);
            } catch (SerializationException e) {
                e.printStackTrace();
            }
        }

        ConnectedAddressItem addressItem = new ConnectedAddressItem(this.connectionId);

        List<TypeId> typeIds =new ArrayList<>(2);
        typeIds.add(addressItem);

        if (requests.size() == 1) {
            typeIds.add(new ConnectedDataItem(this.sequenceCount, requests.get(0)));
        } else {
            List<Integer> offsets = new ArrayList<>(requests.size());
            offsets.add(2 + 2 * request.getNumberOfTags());
            for (CipService cipRequest : requests) {
                if (requests.indexOf(cipRequest) != (requests.size() - 1)) {
                    offsets.add(offsets.get(requests.indexOf(cipRequest)) + cipRequest.getLengthInBytes());
                }

            }
            Services services = new Services(offsets, requests);
            MultipleServiceRequest serviceRequest = new MultipleServiceRequest(services);
            typeIds.add(new ConnectedDataItem(this.sequenceCount, serviceRequest));
        }

        SendUnitData pkt = new SendUnitData(
            sessionHandle,
            CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT,
            0L,
            0,
            typeIds
        );

        this.sequenceCount += 1;

        transaction.submit(() -> context.sendRequest(pkt)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p instanceof SendUnitData)
            .unwrap(p -> (SendUnitData) p)
            .check(p -> p.getSessionHandle() == sessionHandle)
            .handle(p -> {
                List<TypeId> responseTypeIds = p.getTypeIds();
                ConnectedDataItem dataItem = (ConnectedDataItem) responseTypeIds.get(1);
                PlcReadResponse readResponse = decodeReadResponse(dataItem.getService(), request);
                future.complete(readResponse);
                // Finish the request-transaction.
                transaction.endRequest();
            }));

        return future;
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future;
        if (!this.useMessageRouter && !this.useConnectionManager) {
            future = readWithoutMessageRouter(readRequest);
        } else if (this.useMessageRouter && !this.useConnectionManager) {
            future = readWithoutConnectionManager(readRequest);
        } else {
            future = readWithConnectionManager(readRequest);
        }
        return future;
    }

    /*
        Takes a Tag name e.g. ZZZ_ZZZ.XXX and returns a buffer containing an array of ANSI Extended Symbol Seqments
     */
    public static byte[] toAnsi(String tag) throws SerializationException {
        final Pattern RESOURCE_ADDRESS_PATTERN = Pattern.compile("([.\\[\\]])*([A-Za-z_0-9]+){1}");
        Matcher matcher = RESOURCE_ADDRESS_PATTERN.matcher(tag);
        List<PathSegment> segments = new LinkedList<>();
        int lengthBytes = 0;
        while (matcher.find()) {
            String identifier = matcher.group(2);
            String qualifier = matcher.group(1);

            PathSegment newSegment;
            if (qualifier != null) {
                if (qualifier.equals("[")) {
                    newSegment = new LogicalSegment(new MemberID((byte) 0x00, Short.parseShort(identifier)));
                    segments.add(newSegment);
                } else {
                    newSegment = new DataSegment(new AnsiExtendedSymbolSegment(identifier, (short) 0));
                    segments.add(newSegment);
                }
            } else {
                newSegment = new DataSegment(new AnsiExtendedSymbolSegment(identifier,
                    (identifier.length() % 2 == 0) ? null : (short) 0));
                segments.add(newSegment);
            }

            lengthBytes += newSegment.getLengthInBytes();
        }
        WriteBufferByteBased buffer = new WriteBufferByteBased(lengthBytes, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);

        for (PathSegment segment : segments) {
            segment.serialize(buffer);
        }
        return buffer.getBytes();
    }

    private PlcReadResponse decodeReadResponse(CipService p, PlcReadRequest readRequest) {
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        // only 1 field
        if (p instanceof CipReadResponse) {
            CipReadResponse resp = (CipReadResponse) p;
            String fieldName = readRequest.getTagNames().iterator().next();
            EipTag tag = (EipTag) readRequest.getTag(fieldName);
            PlcResponseCode code = decodeResponseCode(resp.getStatus());
            PlcValue plcValue = null;
            CIPDataTypeCode type = resp.getData().getDataType();
            ByteBuf data = Unpooled.wrappedBuffer(resp.getData().getData());
            if (code == PlcResponseCode.OK) {
                plcValue = parsePlcValue(tag, data, type);
            }
            ResponseItem<PlcValue> result = new ResponseItem<>(code, plcValue);
            values.put(fieldName, result);
        }
        //Multiple response
        else if (p instanceof MultipleServiceResponse) {
            MultipleServiceResponse responses = (MultipleServiceResponse) p;
            int nb = responses.getServiceNb();
            List<CipService> arr = new ArrayList<>(nb);
            ReadBufferByteBased read = new ReadBufferByteBased(responses.getServicesData(), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
            int total = read.getTotalBytes();
            for (int i = 0; i < nb; i++) {
                int length = 0;
                int offset = responses.getOffsets().get(i) - responses.getOffsets().get(0); //Substract first offset as we only have the service in the buffer (not servicesNb and offsets)
                if (i == nb - 1) {
                    length = total - offset; //Get the rest if last
                } else {
                    length = responses.getOffsets().get(i + 1) - offset - responses.getOffsets().get(0); //Calculate length with offsets (substracting first offset)
                }
                CipService service;
                try {
                    service = CipService.staticParse(read, false, length);
                    arr.add(service);
                } catch (ParseException e) {
                    throw new PlcRuntimeException(e);
                }
            }
            Services services = new Services(responses.getOffsets(), arr);
            Iterator<String> it = readRequest.getTagNames().iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String fieldName = it.next();
                EipTag tag = (EipTag) readRequest.getTag(fieldName);
                PlcValue plcValue = null;
                if (services.getServices().get(i) instanceof CipReadResponse) {
                    CipReadResponse readResponse = (CipReadResponse) services.getServices().get(i);
                    PlcResponseCode code;
                    if (readResponse.getStatus() == 0) {
                        code = PlcResponseCode.OK;
                    } else {
                        code = PlcResponseCode.INTERNAL_ERROR;
                    }
                    CIPDataTypeCode type = readResponse.getData().getDataType();
                    ByteBuf data = Unpooled.wrappedBuffer(readResponse.getData().getData());
                    if (code == PlcResponseCode.OK) {
                        plcValue = parsePlcValue(tag, data, type);
                    }
                    ResponseItem<PlcValue> result = new ResponseItem<>(code, plcValue);
                    values.put(fieldName, result);
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest, values);
    }

    private Map<String, ResponseItem<PlcValue>> decodeSingleReadResponse(CipService p, String tagName, PlcTag tag) {
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        CipReadResponse resp = (CipReadResponse) p;
        PlcResponseCode code = decodeResponseCode(resp.getStatus());
        PlcValue plcValue = null;
        CIPDataTypeCode type = resp.getData().getDataType();
        ByteBuf data = Unpooled.wrappedBuffer(resp.getData().getData());
        if (code == PlcResponseCode.OK) {
            plcValue = parsePlcValue((EipTag) tag, data, type);
        }
        ResponseItem<PlcValue> result = new ResponseItem<>(code, plcValue);
        values.put(tagName, result);
        return values;
    }

    private PlcValue parsePlcValue(EipTag tag, ByteBuf data, CIPDataTypeCode type) {
        final int STRING_LEN_OFFSET = 2, STRING_DATA_OFFSET = 6;
        int nb = tag.getElementNb();
        if (nb > 1) {
            int index = 0;
            List<PlcValue> list = new ArrayList<>();
            for (int i = 0; i < nb; i++) {
                switch (type) {
                    case DINT:
                        list.add(new PlcDINT(Integer.reverseBytes(data.getInt(index))));
                        index += type.getSize();
                        break;
                    case INT:
                        list.add(new PlcINT(Integer.reverseBytes(data.getInt(index))));
                        index += type.getSize();
                        break;
                    case SINT:
                        list.add(new PlcSINT(Integer.reverseBytes(data.getInt(index))));
                        index += type.getSize();
                        break;
                    case REAL:
                        list.add(new PlcLREAL(swap(data.getFloat(index))));
                        index += type.getSize();
                        break;
                    case LINT:
                        list.add(new PlcLINT(Long.reverseBytes(data.getLong(index))));
                        index += type.getSize();
                        break;							  
                    case BOOL:
                        list.add(new PlcBOOL(data.getBoolean(index)));
                        index += type.getSize();
                        break;
                    case STRUCTURED: {
                        short structuredType = Short.reverseBytes(data.getShort(0));
                        short structuredLen = Short.reverseBytes(data.getShort(STRING_LEN_OFFSET));
                        if (structuredType == CIPStructTypeCode.STRING.getValue()) {
                            // Length offset is 2, data offset is 6
                            list.add(new PlcSTRING(StandardCharsets
                                .UTF_8.decode(data.nioBuffer(STRING_DATA_OFFSET, structuredLen)).toString()));
                            index += type.getSize();
                        }
                        else {
                            // This is a different type of STRUCTURED data
                            // TODO: return as type STRUCT with structuredType to let user
                            // apps/progs handle it.
                        }
                        return null;
                    }
                    default:
                        return null;
                }
            }
            return new PlcList(list);
        } else {
            switch (type) {
                case SINT:
                    return new PlcSINT(data.getByte(0));
                case INT:
                    return new PlcINT(Short.reverseBytes(data.getShort(0)));
                case DINT:
                    return new PlcDINT(Integer.reverseBytes(data.getInt(0)));
                case LINT:
                    return new PlcLINT(Long.reverseBytes(data.getLong(0)));
                case REAL:
                    return new PlcREAL(swap(data.getFloat(0)));
                case BOOL:
                    return new PlcBOOL(data.getBoolean(0));
                case STRING:
                case STRUCTURED: {
                    short structuredType = Short.reverseBytes(data.getShort(0));
                    short structuredLen = Short.reverseBytes(data.getShort(STRING_LEN_OFFSET));
                    if (structuredType == CIPStructTypeCode.STRING.getValue()) {
                        // Length offset is 2, data offset is 6
                        return new PlcSTRING(StandardCharsets
                            .UTF_8.decode(data.nioBuffer(STRING_DATA_OFFSET, structuredLen)).toString());
                    }
                    else {
                        // This is a different type of STRUCTURED data
                    }
                    return null;
                }							  
                default:
                    return null;
            }
        }
    }

    public float swap(float value) {
        int bytes = Float.floatToIntBits(value);
        int b1 = (bytes) & 0xff;
        int b2 = (bytes >> 8) & 0xff;
        int b3 = (bytes >> 16) & 0xff;
        int b4 = (bytes >> 24) & 0xff;
        return Float.intBitsToFloat(b1 << 24 | b2 << 16 | b3 << 8 | b4);
    }

    public CompletableFuture<PlcWriteResponse> writeWithoutMessageRouter(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfTags());
        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));
        Map<String, PlcResponseCode> values = new HashMap<>();

        for (String fieldName : writeRequest.getTagNames()) {
            final EipTag field = (EipTag) request.getTag(fieldName);
            final PlcValue value = request.getPlcValue(fieldName);
            String tag = field.getTag();
            int elements = Math.max(field.getElementNb(), 1);

            byte[] data = encodeValue(value, field.getType());
            CipWriteRequest writeReq = null;
            try {
                writeReq = new CipWriteRequest(toAnsi(tag), field.getType(), elements, data);
            } catch (SerializationException e) {
                e.printStackTrace();
            }
            CompletableFuture<Boolean> internalFuture = new CompletableFuture<>();
            RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

            tm.startRequest();

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    writeReq,
                    (byte) configuration.getBackplane(),
                    (byte) configuration.getSlot()));

            List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

            CipRRData rrdata = new CipRRData(
                sessionHandle,
                0L,
                senderContext,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                typeIds);

            transaction.submit(() -> context.sendRequest(rrdata)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData).unwrap(p -> (CipRRData) p)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .check(p -> ((UnConnectedDataItem) p.getTypeIds().get(1)).getService() instanceof CipWriteResponse)
                .unwrap(p -> (CipWriteResponse) ((UnConnectedDataItem) p.getTypeIds().get(1)).getService())
                .handle(p -> {
                    Map<String, PlcResponseCode> responseItem = decodeSingleWriteResponse(p, fieldName);
                    values.putAll(responseItem);
                    internalFuture.complete(true);
                    transaction.endRequest();
                })
            );
            try {
                internalFuture.get(REQUEST_TIMEOUT.toMillis(), TimeUnit.MILLISECONDS);
            } catch (InterruptedException | ExecutionException | TimeoutException e) {
                future.completeExceptionally(new PlcRuntimeException("Failed to read field"));
            }

        }
        // TODO: This seems to be blocking here ... we should probably do this asynchronously
        PlcWriteResponse response = new DefaultPlcWriteResponse(writeRequest, values);
        future.complete(response);
        return future;
    }

    public CompletableFuture<PlcWriteResponse> writeWithoutConnectionManager(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String fieldName : request.getTagNames()) {
            final EipTag field = (EipTag) request.getTag(fieldName);
            final PlcValue value = request.getPlcValue(fieldName);
            String tag = field.getTag();
            int elements = Math.max(field.getElementNb(), 1);

            byte[] data = encodeValue(value, field.getType());
            try {
                CipWriteRequest writeReq = new CipWriteRequest(toAnsi(tag), field.getType(), elements, data);
                items.add(writeReq);
            } catch (SerializationException e) {
                e.printStackTrace();
            }

        }

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if (items.size() == 1) {
            tm.startRequest();

            PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
            PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    items.get(0),
                    (byte) configuration.getBackplane(),
                    (byte) configuration.getSlot()));

            List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

            CipRRData rrdata = new CipRRData(
                sessionHandle,
                0L,
                senderContext,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                typeIds);

            transaction.submit(() -> context.sendRequest(rrdata)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData).unwrap(p -> (CipRRData) p)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .check(p -> ((UnConnectedDataItem) p.getTypeIds().get(1)).getService() instanceof CipWriteResponse)
                .unwrap(p -> (CipWriteResponse) ((UnConnectedDataItem) p.getTypeIds().get(1)).getService())
                .handle(p -> {
                    future.complete((PlcWriteResponse) decodeWriteResponse(p, writeRequest));
                    transaction.endRequest();
                })
            );
        } else {
            tm.startRequest();
            short nb = (short) items.size();
            List<Integer> offsets = new ArrayList<>(nb);
            int offset = 2 + nb * 2;
            for (int i = 0; i < nb; i++) {
                offsets.add(offset);
                offset += items.get(i).getLengthInBytes();
            }

            List<CipService> serviceArr = new ArrayList<>(nb);
            for (int i = 0; i < nb; i++) {
                serviceArr.add(items.get(i));
            }
            Services data = new Services(offsets, serviceArr);
            //Encapsulate the data

            PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6));
            PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1));

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    new MultipleServiceRequest(data),
                    (byte) configuration.getBackplane(),
                    (byte) configuration.getSlot()));

            List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

            CipRRData pkt = new CipRRData(
                sessionHandle,
                0L,
                DEFAULT_SENDER_CONTEXT,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                typeIds);

            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> ((UnConnectedDataItem) p.getTypeIds().get(1)).getService()).check(p -> p instanceof MultipleServiceResponse)
                .unwrap(p -> (MultipleServiceResponse) p)
                .check(p -> p.getServiceNb() == nb)
                .handle(p -> {
                    future.complete((PlcWriteResponse) decodeWriteResponse(p, writeRequest));
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        }
        return future;
    }

    public CompletableFuture<PlcWriteResponse> writeWithConnectionManager(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String fieldName : request.getTagNames()) {
            final EipTag field = (EipTag) request.getTag(fieldName);
            final PlcValue value = request.getPlcValue(fieldName);
            String tag = field.getTag();
            int elements = Math.max(field.getElementNb(), 1);

            byte[] data = encodeValue(value, field.getType());
            try {
                CipWriteRequest writeReq = new CipWriteRequest(toAnsi(tag), field.getType(), elements, data);
                items.add(writeReq);
            } catch (SerializationException e) {
                e.printStackTrace();
            }

        }

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if (items.size() == 1) {
            tm.startRequest();

            ConnectedDataItem exchange = new ConnectedDataItem(
                this.sequenceCount,
                items.get(0));

            ConnectedAddressItem addressItem = new ConnectedAddressItem(this.connectionId);

            List<TypeId> typeIds = Arrays.asList(addressItem, exchange);

            SendUnitData rrdata = new SendUnitData(
                sessionHandle,
                CIPStatus.Success.getValue(),
                senderContext,
                0L,
                0,
                typeIds);

            transaction.submit(() -> context.sendRequest(rrdata)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof SendUnitData).unwrap(p -> (SendUnitData) p)
                .check(p -> p.getSessionHandle() == sessionHandle)
                .check(p -> ((ConnectedDataItem) p.getTypeIds().get(1)).getService() instanceof CipWriteResponse)
                .unwrap(p -> (CipWriteResponse) ((ConnectedDataItem) p.getTypeIds().get(1)).getService())
                .handle(p -> {
                    future.complete((PlcWriteResponse) decodeWriteResponse(p, writeRequest));
                    transaction.endRequest();
                })
            );
        } else {
            tm.startRequest();
            short nb = (short) items.size();
            List<Integer> offsets = new ArrayList<>(nb);
            int offset = 2 + nb * 2;
            for (int i = 0; i < nb; i++) {
                offsets.add(offset);
                offset += items.get(i).getLengthInBytes();
            }

            List<CipService> serviceArr = new ArrayList<>(nb);
            for (int i = 0; i < nb; i++) {
                serviceArr.add(items.get(i));
            }
            Services data = new Services(offsets, serviceArr);
            //Encapsulate the data

            ConnectedDataItem exchange = new ConnectedDataItem(
                this.sequenceCount,
                new MultipleServiceRequest(data));

            List<TypeId> typeIds = Arrays.asList(nullAddressItem, exchange);

            SendUnitData pkt = new SendUnitData(
                sessionHandle,
                0L,
                DEFAULT_SENDER_CONTEXT,
                0L,
                0,
                typeIds);

            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof SendUnitData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (SendUnitData) p)
                .unwrap(p -> ((ConnectedDataItem) p.getTypeIds().get(1)).getService()).check(p -> p instanceof MultipleServiceResponse)
                .unwrap(p -> (MultipleServiceResponse) p)
                .check(p -> p.getServiceNb() == nb)
                .handle(p -> {
                    future.complete((PlcWriteResponse) decodeWriteResponse(p, writeRequest));
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        }
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future;
        if (!this.useMessageRouter && !this.useConnectionManager) {
            future = writeWithoutMessageRouter(writeRequest);
        } else if (this.useMessageRouter && !this.useConnectionManager) {
            future = writeWithoutConnectionManager(writeRequest);
        } else {
            future = writeWithConnectionManager(writeRequest);
        }
        return future;
    }

    private Map<String, PlcResponseCode> decodeSingleWriteResponse(CipWriteResponse resp, String fieldName) {
        Map<String, PlcResponseCode> responses = new HashMap<>();
        responses.put(fieldName, decodeResponseCode(resp.getStatus()));
        return responses;
    }

    private PlcResponse decodeWriteResponse(CipService p, PlcWriteRequest writeRequest) {
        Map<String, PlcResponseCode> responses = new HashMap<>();

        if (p instanceof CipWriteResponse) {
            CipWriteResponse resp = (CipWriteResponse) p;
            String fieldName = writeRequest.getTagNames().iterator().next();
            responses.put(fieldName, decodeResponseCode(resp.getStatus()));
            return new DefaultPlcWriteResponse(writeRequest, responses);
        } else if (p instanceof MultipleServiceResponse) {
            MultipleServiceResponse resp = (MultipleServiceResponse) p;
            int nb = resp.getServiceNb();
            List<CipService> arr = new ArrayList<>(nb);
            ReadBufferByteBased read = new ReadBufferByteBased(resp.getServicesData());
            int total = read.getTotalBytes();
            for (int i = 0; i < nb; i++) {
                int length;
                int offset = resp.getOffsets().get(i);
                if (offset == nb - 1) {
                    length = total - offset; //Get the rest if last
                } else {
                    length = resp.getOffsets().get(i + 1) - offset; //Calculate length with offsets
                }
                CipService service;
                try {
                    service = CipService.staticParse(read, false, length);
                    arr.add(service);
                } catch (ParseException e) {
                    throw new PlcRuntimeException(e);
                }
            }
            Services services = new Services(resp.getOffsets(), arr);
            Iterator<String> it = writeRequest.getTagNames().iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String fieldName = it.next();
                if (services.getServices().get(i) instanceof CipWriteResponse) {
                    CipWriteResponse writeResponse = (CipWriteResponse) services.getServices().get(i);
                    PlcResponseCode code = decodeResponseCode(writeResponse.getStatus());
                    responses.put(fieldName, code);
                }
            }
            return new DefaultPlcWriteResponse(writeRequest, responses);
        }
        return null;
    }

    private byte[] encodeValue(PlcValue value, CIPDataTypeCode type) {
        //ByteBuffer buffer = ByteBuffer.allocate(4+type.getSize()).order(ByteOrder.LITTLE_ENDIAN);
        ByteBuffer buffer = ByteBuffer.allocate(type.getSize()).order(ByteOrder.LITTLE_ENDIAN);
        switch (type) {
            case BOOL:
                buffer.put(value.getByte());
                break;
            case SINT:
                buffer.put(value.getByte());
                break;
            case INT:
                buffer.putShort(value.getShort());
                break;
            case DINT:
                buffer.putInt(value.getInteger());
                break;
            case REAL:
                buffer.putDouble(value.getDouble());
                break;
            case LINT:
                buffer.putLong(value.getLong());
                break;
            case STRING:
            case STRUCTURED:
                buffer.putInt(value.getString().length());
                buffer.put(value.getString().getBytes(), 0, value.getString().length());
                break;
            default:
                break;
        }
        return buffer.array();

    }

    private PlcResponseCode decodeResponseCode(int status) {
        //TODO other status
        switch (status) {
            case 0:
                return PlcResponseCode.OK;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    @Override
    public void close(ConversationContext<EipPacket> context) {
    }

}
