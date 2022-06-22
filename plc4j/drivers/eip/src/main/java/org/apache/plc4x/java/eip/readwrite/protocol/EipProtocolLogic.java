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
package org.apache.plc4x.java.eip.readwrite.protocol;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.eip.readwrite.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.readwrite.field.EipField;
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
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class EipProtocolLogic extends Plc4xProtocolBase<EipPacket> implements HasConfiguration<EIPConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(EipProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static final byte[] DEFAULT_SENDER_CONTEXT = "PLC4X   ".getBytes(StandardCharsets.US_ASCII);
    private static final long EMPTY_SESSION_HANDLE = 0L;
    private static final long EMPTY_INTERFACE_HANDLE = 0L;

    private byte[] senderContext;

    private long connectionId = 0L;
    private int sequenceCount = 1;
    private EIPConfiguration configuration;

    private final AtomicInteger transactionCounterGenerator = new AtomicInteger(10);
    private RequestTransactionManager tm;
    private long sessionHandle;

    @Override
    public void setConfiguration(EIPConfiguration configuration) {
        this.configuration = configuration;
        // Set the transaction manager to allow only one message at a time.
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void onConnect(ConversationContext<EipPacket> context) {
        logger.debug("Sending Register Session EIP Package");
        EipConnectionRequest connectionRequest =
            new EipConnectionRequest(
                EMPTY_SESSION_HANDLE,
                CIPStatus.Success.getValue(),
                DEFAULT_SENDER_CONTEXT,
                0L,
                this.configuration.getByteOrder());
        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT).unwrap(p -> p)
            .check(p -> p instanceof EipConnectionRequest)
            .handle(p -> {
                if (p.getStatus() == CIPStatus.Success.getValue()) {
                    sessionHandle = p.getSessionHandle();
                    senderContext = p.getSenderContext();
                    logger.debug("Got assigned with Session {}", sessionHandle);
                    onConnectOpenConnectionManager(context, p);

                } else {
                    logger.warn("Got status code [{}]", p.getStatus());
                }
            });
    }

    public void onConnectOpenConnectionManager(ConversationContext<EipPacket> context, EipPacket response) {
        logger.debug("Sending Open Connection Manager EIP Package");

        PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6, this.configuration.getByteOrder()), this.configuration.getByteOrder());
        PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1, this.configuration.getByteOrder()), this.configuration.getByteOrder());

        UnConnectedDataItem exchange = new UnConnectedDataItem(
            new CipConnectionManagerRequest(
                classSegment,
                instanceSegment,
                (byte) 0,
                (byte) 10,
                (short) 14,
                536870914L,
                33944L,
                8592,
                4919,
                42L,
                (short) 3,
                2101812L,
                new NetworkConnectionParameters(
                    4002,
                    false,
                    (byte) 2,
                    (byte) 0,
                    true,
                    this.configuration.getByteOrder()
                ),
                (long) 2113537,
                new NetworkConnectionParameters(
                    4002,
                    false,
                    (byte) 2,
                    (byte) 0,
                    true,
                    this.configuration.getByteOrder()
                ),
                new TransportType(true, (byte) 2, (byte) 3, this.configuration.getByteOrder()),
                (short) 3,
                new PortSegment(false, (byte) 1, (short) 0, this.configuration.getByteOrder()),
                new LogicalSegment(new ClassID((byte) 0, (short) 2, this.configuration.getByteOrder()), this.configuration.getByteOrder()),
                new LogicalSegment(new InstanceID((byte) 0, (short) 1, this.configuration.getByteOrder()), this.configuration.getByteOrder()),
                (int) 0,
                this.configuration.getByteOrder()
            ),
            this.configuration.getByteOrder()
        );

        List<TypeId> typeIds = new ArrayList<TypeId>() {
            {
                add(exchange);
            }
        };

        CipRRData eipWrapper = new CipRRData(
            sessionHandle,
            response.getStatus(),
            senderContext,
            0L,
            EMPTY_INTERFACE_HANDLE,
            0,
            2,
            typeIds,
            this.configuration.getByteOrder(),
            -1
        );


        context.sendRequest(eipWrapper)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT).unwrap(p -> p)
            .check(p -> p instanceof CipRRData)
            .handle(p -> {
                if (p.getStatus() == 0L) {
                    CipRRData rrData = (CipRRData) p;
                    sessionHandle = rrData.getSessionHandle();
                    senderContext = rrData.getSenderContext();
                    List<TypeId> connectionManagerExchange = rrData.getTypeId();
                    CipConnectionManagerResponse connectionManagerResponse = (CipConnectionManagerResponse) ((ConnectedDataItem) connectionManagerExchange.get(0)).getService();
                    this.connectionId = connectionManagerResponse.getOtConnectionId();

                    logger.debug("Got assigned with Session {}", sessionHandle);
                    // Send an event that connection setup is complete.
                    context.fireConnected();
                } else {
                    logger.warn("Got status code [{}]", p.getStatus());
                }

            });
    }


    @Override
    public void onDisconnect(ConversationContext<EipPacket> context) {
        if (this.connectionId != 0L) {
            logger.debug("Sending Connection Manager Close Event");
            PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6, this.configuration.getByteOrder()), this.configuration.getByteOrder());
            PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1, this.configuration.getByteOrder()), this.configuration.getByteOrder());

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipConnectionManagerCloseRequest(
                    (byte) 2,
                    classSegment,
                    instanceSegment,
                    (byte) 0,
                    (byte) 10,
                    (short) 14,
                    8592,
                    4919,
                    42L,
                    (short) 3,
                    new PortSegment(false, (byte) 1, (short) 0, this.configuration.getByteOrder()),
                    new LogicalSegment(new ClassID((byte) 0, (short) 2, this.configuration.getByteOrder()), this.configuration.getByteOrder()),
                    new LogicalSegment(new InstanceID((byte) 0, (short) 1, this.configuration.getByteOrder()), this.configuration.getByteOrder()),
                    0,
                    this.configuration.getByteOrder()
                ),
                this.configuration.getByteOrder()
            );

            List<TypeId> typeIds = new ArrayList<TypeId>() {
                {
                    add(exchange);
                }
            };

            CipRRData eipWrapper = new CipRRData(
                sessionHandle,
                0L,
                senderContext,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                2,
                typeIds,
                this.configuration.getByteOrder(),
                -1
            );


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
                0L,
                sessionHandle,
                DEFAULT_SENDER_CONTEXT,
                0L,
                this.configuration.getByteOrder());
        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, Duration.ofMillis(1))
            .onTimeout(p -> context.fireDisconnected())
            .handle(p -> context.fireDisconnected());
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<CipService> requests = new ArrayList<>(request.getNumberOfFields());
        for (PlcField field : request.getFields()) {
            EipField plcField = (EipField) field;
            String tag = plcField.getTag();
            int elements = 1;
            if (plcField.getElementNb() > 1) {
                elements = plcField.getElementNb();
            }
            try {
                CipReadRequest req = new CipReadRequest(
                    toAnsi(tag, this.configuration.getByteOrder()),
                    1,
                    -1,
                    this.configuration.getByteOrder());
                requests.add(req);
            } catch (SerializationException e) {
                e.printStackTrace();
            }
        }

        List<TypeId> typeIds =new ArrayList<>(2);
        typeIds.add(new ConnectedAddressItem(this.connectionId, this.configuration.getByteOrder()));
        if (requests.size() == 1) {
            typeIds.add(new ConnectedDataItem(this.sequenceCount, requests.get(0), this.configuration.getByteOrder()));
        } else {
            List<Integer> offsets = new ArrayList<>(requests.size());
            offsets.add(6);
            for (CipService cipRequest : requests) {
                if (requests.indexOf(cipRequest) != (requests.size() - 1)) {
                    offsets.add(offsets.get(requests.indexOf(cipRequest)) + cipRequest.getLengthInBytes());
                }

            }
            MultipleServiceRequest serviceRequest = new MultipleServiceRequest(new Services(requests.size(), offsets, requests, 0, this.configuration.getByteOrder()), 0, this.configuration.getByteOrder());
            typeIds.add(new ConnectedDataItem(this.sequenceCount, serviceRequest, this.configuration.getByteOrder()));
        }


        SendUnitData pkt = new SendUnitData(
            sessionHandle,
            CIPStatus.Success.getValue(),
            DEFAULT_SENDER_CONTEXT,
            0L,
            0,
            2,
            typeIds,
            this.configuration.getByteOrder()
        );

        transaction.submit(() -> context.sendRequest(pkt)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .onTimeout(future::completeExceptionally)
            .onError((p, e) -> future.completeExceptionally(e))
            .check(p -> p instanceof SendUnitData)
            .check(p -> p.getSessionHandle() == sessionHandle)
            //.check(p -> p.getSenderContext() == senderContext)
            .unwrap(p -> (SendUnitData) p)
            .handle(p -> {
                SendUnitData unitData = (SendUnitData) p;
                List<TypeId> responseTypeIds = unitData.getTypeId();
                ConnectedDataItem dataItem = (ConnectedDataItem) responseTypeIds.get(1);
                future.complete(decodeReadResponse(dataItem.getService(), request));
                // Finish the request-transaction.
                transaction.endRequest();
            }));

        return future;
    }

    /*
        Takes a Tag name e.g. ZZZ_ZZZ.XXX and returns a buffer containing an array of ANSI Extended Symbol Seqments
     */
    public static byte[] toAnsi(String tag, IntegerEncoding order) throws SerializationException {
        final Pattern RESOURCE_ADDRESS_PATTERN = Pattern.compile("([.\\[\\]])*([A-Za-z_0-9]+){1}");
        Matcher matcher = RESOURCE_ADDRESS_PATTERN.matcher(tag);
        List<PathSegment> segments = new LinkedList<>();
        String tagWithoutQualifiers = "";
        int lengthBytes = 0;
        while (matcher.find()) {
            String identifier = matcher.group(2);
            String qualifier = matcher.group(1);

            PathSegment newSegment;
            if (qualifier != null) {
                switch (qualifier) {
                    case "[":
                        newSegment = new LogicalSegment(new MemberID((byte) 0x00, (short) Short.parseShort(identifier), order), order);
                        segments.add(newSegment);
                        break;
                    default:
                        newSegment = new DataSegment(new AnsiExtendedSymbolSegment(identifier, (short) 0, order), order);
                        segments.add(newSegment);
                        tagWithoutQualifiers += identifier;
                }
            } else {
                newSegment = new DataSegment(new AnsiExtendedSymbolSegment(identifier, (short) 0, order), order);
                segments.add(newSegment);
                tagWithoutQualifiers += identifier;
            }

            lengthBytes += newSegment.getLengthInBytes();
        }
        WriteBufferByteBased buffer = new WriteBufferByteBased(lengthBytes, org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);

        for (PathSegment segment : segments) {
            segment.serialize(buffer);
        }
        return buffer.getData();
    }

    private PlcReadResponse decodeReadResponse(CipService p, PlcReadRequest readRequest) {
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        // only 1 field
        if (p instanceof CipReadResponse) {
            CipReadResponse resp = (CipReadResponse) p;
            String fieldName = readRequest.getFieldNames().iterator().next();
            EipField field = (EipField) readRequest.getField(fieldName);
            PlcResponseCode code = decodeResponseCode(resp.getStatus());
            PlcValue plcValue = null;
            CIPDataTypeCode type = resp.getDataType();
            ByteBuf data = Unpooled.wrappedBuffer(resp.getData());
            if (code == PlcResponseCode.OK) {
                plcValue = parsePlcValue(field, data, type);
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
            int total = (int) read.getTotalBytes();
            for (int i = 0; i < nb; i++) {
                int length = 0;
                int offset = responses.getOffsets().get(i) - responses.getOffsets().get(0); //Substract first offset as we only have the service in the buffer (not servicesNb and offsets)
                if (i == nb - 1) {
                    length = total - offset; //Get the rest if last
                } else {
                    length = responses.getOffsets().get(i + 1) - offset - responses.getOffsets().get(0); //Calculate length with offsets (substracting first offset)
                }
                ReadBuffer serviceBuf = new ReadBufferByteBased(read.getBytes(offset, offset + length), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                CipService service = null;
                try {
                    service = CipService.staticParse(read, length);
                    arr.add(service);
                } catch (ParseException e) {
                    throw new PlcRuntimeException(e);
                }
            }
            Services services = new Services(nb, responses.getOffsets(), arr, -1, this.configuration.getByteOrder());
            Iterator<String> it = readRequest.getFieldNames().iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String fieldName = it.next();
                EipField field = (EipField) readRequest.getField(fieldName);
                PlcValue plcValue = null;
                if (services.getServices().get(i) instanceof CipReadResponse) {
                    CipReadResponse readResponse = (CipReadResponse) services.getServices().get(i);
                    PlcResponseCode code;
                    if (readResponse.getStatus() == 0) {
                        code = PlcResponseCode.OK;
                    } else {
                        code = PlcResponseCode.INTERNAL_ERROR;
                    }
                    CIPDataTypeCode type = readResponse.getDataType();
                    ByteBuf data = Unpooled.wrappedBuffer(readResponse.getData());
                    if (code == PlcResponseCode.OK) {
                        plcValue = parsePlcValue(field, data, type);
                    }
                    ResponseItem<PlcValue> result = new ResponseItem<>(code, plcValue);
                    values.put(fieldName, result);
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest, values);
    }

    private PlcValue parsePlcValue(EipField field, ByteBuf data, CIPDataTypeCode type) {
        final int STRING_LEN_OFFSET = 2, STRING_DATA_OFFSET = 6;
        int nb = field.getElementNb();
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
                        Short structuredType = Short.reverseBytes(data.getShort(0));
                        Short structuredLen = Short.reverseBytes(data.getShort(STRING_LEN_OFFSET));
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
                case STRING36:
                case STRING:
                case STRUCTURED: {
                    Short structuredType = Short.reverseBytes(data.getShort(0));
                    Short structuredLen = Short.reverseBytes(data.getShort(STRING_LEN_OFFSET));
                    if (structuredType == CIPStructTypeCode.STRING.getValue()) {
                        // Length offset is 2, data offset is 6
                        return new PlcSTRING(StandardCharsets
                            .UTF_8.decode(data.nioBuffer(STRING_DATA_OFFSET, structuredLen)).toString());
                    }
                    else {
                        // This is a different type of STRUCTURED data
                    }
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

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        DefaultPlcWriteRequest request = (DefaultPlcWriteRequest) writeRequest;
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfFields());
        for (String fieldName : request.getFieldNames()) {
            final EipField field = (EipField) request.getField(fieldName);
            final PlcValue value = request.getPlcValue(fieldName);
            String tag = field.getTag();
            int elements = 1;
            if (field.getElementNb() > 1) {
                elements = field.getElementNb();
            }

            byte[] data = encodeValue(value, field.getType(), (short) elements);
            try {
                CipWriteRequest writeReq = new CipWriteRequest(toAnsi(tag, this.configuration.getByteOrder()), field.getType(), elements, data, -1, this.configuration.getByteOrder());
                items.add(writeReq);
            } catch (SerializationException e) {
                e.printStackTrace();
            }

        }

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if (items.size() == 1) {
            tm.startRequest();

            PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6, this.configuration.getByteOrder()), this.configuration.getByteOrder());
            PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1, this.configuration.getByteOrder()), this.configuration.getByteOrder());

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    items.get(0),
                    (byte) configuration.getBackplane(),
                    (byte) configuration.getSlot(),
                    -1,
                    this.configuration.getByteOrder()

                ),
                this.configuration.getByteOrder()
            );

            List<TypeId> typeIds = new ArrayList<TypeId>() {
                {
                    add(exchange);
                }
            };

            CipRRData rrdata = new CipRRData(
                sessionHandle,
                0L,
                senderContext,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                2,
                typeIds,
                this.configuration.getByteOrder(),
                -1
            );
            transaction.submit(() -> context.sendRequest(rrdata)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData).unwrap(p -> (CipRRData) p)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .check(p -> ((UnConnectedDataItem) p.getTypeId().get(0)).getService() instanceof CipWriteResponse)
                .unwrap(p -> (CipWriteResponse) ((UnConnectedDataItem) p.getTypeId().get(0)).getService())
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
            Services data = new Services(nb, offsets, serviceArr, -1, this.configuration.getByteOrder());
            //Encapsulate the data

            PathSegment classSegment = new LogicalSegment(new ClassID((byte) 0, (short) 6, this.configuration.getByteOrder()), this.configuration.getByteOrder());
            PathSegment instanceSegment = new LogicalSegment(new InstanceID((byte) 0, (short) 1, this.configuration.getByteOrder()), this.configuration.getByteOrder());

            UnConnectedDataItem exchange = new UnConnectedDataItem(
                new CipUnconnectedRequest(
                    classSegment,
                    instanceSegment,
                    new MultipleServiceRequest(data, -1, this.configuration.getByteOrder()),
                    (byte) configuration.getBackplane(),
                    (byte) configuration.getSlot(),
                    (Integer) 0,
                    this.configuration.getByteOrder()
                ),
                this.configuration.getByteOrder()
            );

            List<TypeId> typeIds = new ArrayList<TypeId>() {
                {
                    add(exchange);
                }
            };

            CipRRData pkt = new CipRRData(
                sessionHandle,
                0L,
                DEFAULT_SENDER_CONTEXT,
                0L,
                EMPTY_INTERFACE_HANDLE,
                0,
                2,
                typeIds,
                this.configuration.getByteOrder(),
                -1
            );


            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> ((UnConnectedDataItem) p.getTypeId().get(0)).getService()).check(p -> p instanceof MultipleServiceResponse)
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

    private PlcResponse decodeWriteResponse(CipService p, PlcWriteRequest writeRequest) {
        Map<String, PlcResponseCode> responses = new HashMap<>();

        if (p instanceof CipWriteResponse) {
            CipWriteResponse resp = (CipWriteResponse) p;
            String fieldName = writeRequest.getFieldNames().iterator().next();
            EipField field = (EipField) writeRequest.getField(fieldName);
            responses.put(fieldName, decodeResponseCode(resp.getStatus()));
            return new DefaultPlcWriteResponse(writeRequest, responses);
        } else if (p instanceof MultipleServiceResponse) {
            MultipleServiceResponse resp = (MultipleServiceResponse) p;
            int nb = resp.getServiceNb();
            List<CipService> arr = new ArrayList<>(nb);
            ReadBufferByteBased read = new ReadBufferByteBased(resp.getServicesData());
            int total = (int) read.getTotalBytes();
            for (int i = 0; i < nb; i++) {
                int length = 0;
                int offset = resp.getOffsets().get(i);
                if (offset == nb - 1) {
                    length = total - offset; //Get the rest if last
                } else {
                    length = resp.getOffsets().get(i + 1) - offset; //Calculate length with offsets
                }
                ReadBuffer serviceBuf = new ReadBufferByteBased(read.getBytes(offset, length), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                CipService service = null;
                try {
                    service = CipService.staticParse(read, length);
                    arr.add(service);
                } catch (ParseException e) {
                    throw new PlcRuntimeException(e);
                }
            }
            Services services = new Services(nb, resp.getOffsets(), arr, -1, this.configuration.getByteOrder());
            Iterator<String> it = writeRequest.getFieldNames().iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String fieldName = it.next();
                EipField field = (EipField) writeRequest.getField(fieldName);
                PlcValue plcValue = null;
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

    private byte[] encodeValue(PlcValue value, CIPDataTypeCode type, short elements) {
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
            case STRING36:
                buffer.putInt(value.getString().length());
                buffer.put(value.getString().getBytes(), 0, value.getString().length());
                break;
            case STRUCTURED:
                // Need to handle
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
        onDisconnect(context);
    }
}
