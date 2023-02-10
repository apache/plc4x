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
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.eip.readwrite.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.readwrite.tag.EipTag;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
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

public class EipProtocolLogic extends Plc4xProtocolBase<EipPacket> implements HasConfiguration<EIPConfiguration> {

    private static final Logger logger = LoggerFactory.getLogger(EipProtocolLogic.class);
    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);

    private static final List<Short> emptySenderContext = Arrays.asList((short) 0x00, (short) 0x00, (short) 0x00,
        (short) 0x00, (short) 0x00, (short) 0x00, (short) 0x00, (short) 0x00);
    private List<Short> senderContext;
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
        logger.debug("Sending RegisterSession EIP Package");
        EipConnectionRequest connectionRequest =
            new EipConnectionRequest(0L, 0L, emptySenderContext, 0L);
        context.sendRequest(connectionRequest)
            .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
            .check(p -> p instanceof EipConnectionRequest)
            .unwrap(p -> (EipConnectionRequest) p)
            .handle(p -> {
                if (p.getStatus() == 0L) {
                    sessionHandle = p.getSessionHandle();
                    senderContext = p.getSenderContext();
                    logger.debug("Got assigned with Session {}", sessionHandle);
                    // Send an event that connection setup is complete.
                    context.fireConnected();
                } else {
                    logger.warn("Got status code [{}]", p.getStatus());
                }

            });
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        DefaultPlcReadRequest request = (DefaultPlcReadRequest) readRequest;
        List<CipReadRequest> requests = new ArrayList<>(request.getNumberOfTags());
        for (PlcTag tag : request.getTags()) {
            EipTag eipTag = (EipTag) tag;
            String tagName = eipTag.getTag();
            int elements = 1;
            if (eipTag.getElementNb() > 1) {
                elements = eipTag.getElementNb();
            }
            CipReadRequest req = new CipReadRequest(getRequestSize(tagName), toAnsi(tagName), elements);
            requests.add(req);
        }
        return toPlcReadResponse(readRequest, readInternal(requests));
    }

    private byte getRequestSize(String tag) {
        //We need the size of the request in words (0x91, tagLength, ... tag + possible pad)
        // Taking half to get word size
        boolean isArray = false;
        boolean isStruct = false;
        String tagIsolated = tag;
        if (tag.contains("[")) {
            isArray = true;
            tagIsolated = tag.substring(0, tag.indexOf("["));
        }

        if (tag.contains(".")) {
            isStruct = true;
            tagIsolated = tagIsolated.replace(".", "");
        }
        int dataLength = (tagIsolated.length() + 2)
            + (tagIsolated.length() % 2)
            + (isArray ? 2 : 0)
            + (isStruct ? 2 : 0);
        byte requestPathSize = (byte) (dataLength / 2);
        return requestPathSize;
    }

    private byte[] toAnsi(String tag) {
        int arrayIndex = 0;
        boolean isArray = false;
        boolean isStruct = false;
        String tagFinal = tag;
        if (tag.contains("[")) {
            isArray = true;
            String index = tag.substring(tag.indexOf("[") + 1, tag.indexOf("]"));
            arrayIndex = Integer.parseInt(index);
            tagFinal = tag.substring(0, tag.indexOf("["));
        }
        if (tag.contains(".")) {
            tagFinal = tag.substring(0, tag.indexOf("."));
            isStruct = true;
        }
        boolean isPadded = tagFinal.length() % 2 != 0;
        int dataSegLength = 2 + tagFinal.length()
            + (isPadded ? 1 : 0)
            + (isArray ? 2 : 0);

        if (isStruct) {
            for (String subStr : tag.substring(tag.indexOf(".") + 1).split("\\.", -1)) {
                dataSegLength += 2 + subStr.length() + subStr.length() % 2;
            }
        }

        ByteBuffer buffer = ByteBuffer.allocate(dataSegLength).order(ByteOrder.LITTLE_ENDIAN);

        buffer.put((byte) 0x91);
        buffer.put((byte) tagFinal.length());
        byte[] tagBytes = null;
        tagBytes = tagFinal.getBytes(StandardCharsets.US_ASCII);

        buffer.put(tagBytes);
        buffer.position(2 + tagBytes.length);


        if (isPadded) {
            buffer.put((byte) 0x00);
        }

        if (isArray) {
            buffer.put((byte) 0x28);
            buffer.put((byte) arrayIndex);
        }
        if (isStruct) {
            buffer.put(toAnsi(tag.substring(tag.indexOf(".") + 1, tag.length())));
        }
        return buffer.array();
    }

    private CompletableFuture<PlcReadResponse> toPlcReadResponse(PlcReadRequest readRequest, CompletableFuture<CipService> response) {
        return response
            .thenApply(p -> {
                return ((PlcReadResponse) decodeReadResponse(p, readRequest));
            });
    }

    private CompletableFuture<CipService> readInternal(List<CipReadRequest> request) {
        CompletableFuture<CipService> future = new CompletableFuture<>();
        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if (request.size() > 1) {

            short nb = (short) request.size();
            List<Integer> offsets = new ArrayList<>(nb);
            int offset = 2 + nb * 2;
            for (int i = 0; i < nb; i++) {
                offsets.add(offset);
                offset += request.get(i).getLengthInBytes();
            }

            List<CipService> serviceArr = new ArrayList<>(nb);
            for (int i = 0; i < nb; i++) {
                serviceArr.add(request.get(i));
            }
            Services data = new Services(nb, offsets, serviceArr);
            //Encapsulate the data

            CipRRData pkt = new CipRRData(sessionHandle, 0L, emptySenderContext, 0L,
                new CipExchange(
                    new CipUnconnectedRequest(
                        new MultipleServiceRequest(data),
                        (byte) configuration.getBackplane(),
                        (byte) configuration.getSlot()
                    )
                )
            );


            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> p.getExchange().getService()).check(p -> p instanceof MultipleServiceResponse)
                .unwrap(p -> (MultipleServiceResponse) p)
                .check(p -> p.getServiceNb() == nb)
                .handle(p -> {
                    future.complete(p);
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        } else if (request.size() == 1) {
            CipExchange exchange = new CipExchange(
                new CipUnconnectedRequest(
                    request.get(0), (byte) configuration.getBackplane(), (byte) configuration.getSlot()
                )
            );
            CipRRData pkt = new CipRRData(sessionHandle, 0L, emptySenderContext, 0L, exchange);
            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> p.getExchange().getService())
                .check(p -> p instanceof CipReadResponse)
                .unwrap(p -> (CipReadResponse) p)
                .handle(p -> {
                    future.complete(p);
                    // Finish the request-transaction.
                    transaction.endRequest();
                }));
        }
        return future;
    }

    private PlcResponse decodeReadResponse(CipService p, PlcReadRequest readRequest) {
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        // only 1 tag
        if (p instanceof CipReadResponse) {
            CipReadResponse resp = (CipReadResponse) p;
            String tagName = readRequest.getTagNames().iterator().next();
            EipTag tag = (EipTag) readRequest.getTag(tagName);
            PlcResponseCode code = decodeResponseCode(resp.getStatus());
            PlcValue plcValue = null;
            CIPDataTypeCode type = resp.getDataType();
            ByteBuf data = Unpooled.wrappedBuffer(resp.getData());
            if (code == PlcResponseCode.OK) {
                plcValue = parsePlcValue(tag, data, type);
            }
            ResponseItem<PlcValue> result = new ResponseItem<>(code, plcValue);
            values.put(tagName, result);
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
                ReadBuffer serviceBuf = new ReadBufferByteBased(read.getBytes(offset, offset + length), org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN);
                CipService service = null;
                try {
                    service = CipService.staticParse(read, length);
                    arr.add(service);
                } catch (ParseException e) {
                    throw new PlcRuntimeException(e);
                }
            }
            Services services = new Services(nb, responses.getOffsets(), arr);
            Iterator<String> it = readRequest.getTagNames().iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String tagName = it.next();
                EipTag tag = (EipTag) readRequest.getTag(tagName);
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
                        plcValue = parsePlcValue(tag, data, type);
                    }
                    ResponseItem<PlcValue> result = new ResponseItem<>(code, plcValue);
                    values.put(tagName, result);
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest, values);
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
                    case Struct: {
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
                case STRING:
                case Struct: {
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
        List<CipWriteRequest> items = new ArrayList<>(writeRequest.getNumberOfTags());
        for (String tagName : request.getTagNames()) {
            final EipTag tag = (EipTag) request.getTag(tagName);
            final PlcValue value = request.getPlcValue(tagName);
            String tagTag = tag.getTag();
            int elements = 1;
            if (tag.getElementNb() > 1) {
                elements = tag.getElementNb();
            }

            //We need the size of the request in words (0x91, tagLength, ... tag + possible pad)
            // Taking half to get word size
            boolean isArray = false;
            boolean isStruct = false;
            String tagIsolated = tagTag;
            if (tagTag.contains("[")) {
                isArray = true;
                tagIsolated = tagTag.substring(0, tagTag.indexOf("["));
            }

            if (tagTag.contains(".")) {
                isStruct = true;
                tagIsolated = tagIsolated.replace(".", "");
            }

            int dataLength = (tagIsolated.length() + 2 + ((tagIsolated.length() % 2) * 2) + (isArray ? 2 : 0) + (isStruct ? 2 : 0));
            byte requestPathSize = (byte) (dataLength / 2);
            byte[] data = encodeValue(value, tag.getType(), (short) elements);
            CipWriteRequest writeReq = new CipWriteRequest(requestPathSize, toAnsi(tagTag), tag.getType(), elements, data);
            items.add(writeReq);
        }

        RequestTransactionManager.RequestTransaction transaction = tm.startRequest();
        if (items.size() == 1) {
            tm.startRequest();
            CipRRData rrdata = new CipRRData(sessionHandle, 0L, senderContext, 0L,
                new CipExchange(
                    new CipUnconnectedRequest(
                        items.get(0), (byte) configuration.getBackplane(), (byte) configuration.getSlot()
                    )
                )
            );
            transaction.submit(() -> context.sendRequest(rrdata)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData).unwrap(p -> (CipRRData) p)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .check(p -> p.getExchange().getService() instanceof CipWriteResponse)
                .unwrap(p -> (CipWriteResponse) p.getExchange().getService())
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
            Services data = new Services(nb, offsets, serviceArr);
            //Encapsulate the data

            CipRRData pkt = new CipRRData(sessionHandle, 0L, emptySenderContext, 0L,
                new CipExchange(
                    new CipUnconnectedRequest(
                        new MultipleServiceRequest(data),
                        (byte) configuration.getBackplane(),
                        (byte) configuration.getSlot()
                    )
                )
            );


            transaction.submit(() -> context.sendRequest(pkt)
                .expectResponse(EipPacket.class, REQUEST_TIMEOUT)
                .onTimeout(future::completeExceptionally)
                .onError((p, e) -> future.completeExceptionally(e))
                .check(p -> p instanceof CipRRData)
                .check(p -> p.getSessionHandle() == sessionHandle)
                //.check(p -> p.getSenderContext() == senderContext)
                .unwrap(p -> (CipRRData) p)
                .unwrap(p -> p.getExchange().getService()).check(p -> p instanceof MultipleServiceResponse)
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
            String tagName = writeRequest.getTagNames().iterator().next();
            EipTag tag = (EipTag) writeRequest.getTag(tagName);
            responses.put(tagName, decodeResponseCode(resp.getStatus()));
            return new DefaultPlcWriteResponse(writeRequest, responses);
        } else if (p instanceof MultipleServiceResponse) {
            MultipleServiceResponse resp = (MultipleServiceResponse) p;
            int nb = resp.getServiceNb();
            List<CipService> arr = new ArrayList<>(nb);
            ReadBufferByteBased read = new ReadBufferByteBased(resp.getServicesData());
            int total = read.getTotalBytes();
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
            Services services = new Services(nb, resp.getOffsets(), arr);
            Iterator<String> it = writeRequest.getTagNames().iterator();
            for (int i = 0; i < nb && it.hasNext(); i++) {
                String tagName = it.next();
                EipTag tag = (EipTag) writeRequest.getTag(tagName);
                PlcValue plcValue = null;
                if (services.getServices().get(i) instanceof CipWriteResponse) {
                    CipWriteResponse writeResponse = (CipWriteResponse) services.getServices().get(i);
                    PlcResponseCode code = decodeResponseCode(writeResponse.getStatus());
                    responses.put(tagName, code);
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
                buffer.putInt(value.getString().length());
                buffer.put(value.getString().getBytes(), 0, value.getString().length());
                break;
            case Struct:
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
        logger.debug("Sending UnregisterSession EIP Packet");
        context.sendRequest(new EipDisconnectRequest(sessionHandle, 0L, emptySenderContext, 0L)); //Unregister gets no response
        logger.debug("Unregistered Session {}", sessionHandle);
    }

    @Override
    protected void decode(ConversationContext<EipPacket> context, EipPacket msg) throws Exception {
        super.decode(context, msg);
    }

}
