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
package org.apache.plc4x.java.s7.netty;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.*;
import org.apache.plc4x.java.s7.model.S7Field;
import org.apache.plc4x.java.s7.netty.events.S7ConnectedEvent;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.lang.reflect.Array;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

/**
 * This layer transforms between {@link PlcRequestContainer}s {@link S7Message}s.
 * And stores all "in-flight" requests in an internal structure ({@link Plc4XS7Protocol#requests}).
 * <p>
 * While sending a request, a {@link S7RequestMessage} is generated and send downstream (to the {@link S7Protocol}.
 * <p>
 * When a {@link S7ResponseMessage} is received it takes the existing request container from its Map and finishes
 * the {@link PlcRequestContainer}s future with the {@link PlcResponse}.
 */
public class Plc4XS7Protocol extends PlcMessageToMessageCodec<S7Message, PlcRequestContainer> {
    private static final Logger logger = LoggerFactory.getLogger( Plc4XS7Protocol.class );

    private static final AtomicInteger tpduGenerator = new AtomicInteger(10);

    private Map<Short, PlcRequestContainer> requests;

    public Plc4XS7Protocol() {
        this.requests = new HashMap<>();
    }

    /**
     * If this protocol layer catches an {@link S7ConnectedEvent} from the protocol layer beneath,
     * the connection establishment is finished.
     *
     * @param ctx the current protocol layers context
     * @param evt the event
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        if (evt instanceof S7ConnectedEvent) {
            ctx.channel().pipeline().fireUserEventTriggered(new ConnectedEvent());
        } else {
            super.userEventTriggered(ctx, evt);
        }
    }

    /**
     * When receiving an error inside the pipeline, we have to find out which {@link PlcRequestContainer}
     * correlates needs to be notified about the problem. If a container is found, we can relay the
     * exception to that by calling completeExceptionally and passing in the exception.
     *
     * @param ctx   the current protocol layers context
     * @param cause the exception that was caught
     * @throws Exception throws an exception if something goes wrong internally
     */
    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception {
        if (cause instanceof PlcProtocolPayloadTooBigException) {
            PlcProtocolPayloadTooBigException pptbe = (PlcProtocolPayloadTooBigException) cause;
            if (pptbe.getPayload() instanceof S7RequestMessage) {
                S7RequestMessage request = (S7RequestMessage) pptbe.getPayload();
                if (request.getParent() instanceof PlcRequestContainer) {
                    PlcRequestContainer requestContainer = (PlcRequestContainer) request.getParent();

                    // Remove the current request from the unconfirmed requests list.
                    requests.remove(request.getTpduReference());

                    requestContainer.getResponseFuture().completeExceptionally(cause);
                }
            }
        } else if ((cause instanceof IOException) && (cause.getMessage().contains("Connection reset by peer") ||
            cause.getMessage().contains("Operation timed out"))) {
            String reason = cause.getMessage().contains("Connection reset by peer") ?
                "Connection terminated unexpectedly" : "Remote host not responding";
            if (!requests.isEmpty()) {
                // If the connection is hung up, all still pending requests can be closed.
                for (PlcRequestContainer requestContainer : requests.values()) {
                    requestContainer.getResponseFuture().completeExceptionally(new PlcIoException(reason));
                }
                // Clear the list
                requests.clear();
            }
        } else {
            super.exceptionCaught(ctx, cause);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws PlcException {
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeReadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        }
    }

    private void encodeReadRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        List<VarParameterItem> parameterItems = new LinkedList<>();

        PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();
        for (String fieldName : readRequest.getFieldNames()) {
            PlcField field = readRequest.getField(fieldName);
            if (!(field instanceof S7Field)) {
                throw new PlcProtocolException("The field should have been of type S7Field");
            }
            S7Field s7Field = (S7Field) field;

            VarParameterItem varParameterItem = new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Field.getMemoryArea(),
                s7Field.getDataType(),
                s7Field.getNumElements(), s7Field.getBlockNumber(), s7Field.getByteOffset(), (byte) s7Field.getBitOffset());
            parameterItems.add(varParameterItem);
        }
        VarParameter readVarParameter = new VarParameter(ParameterType.READ_VAR, parameterItems);

        // Assemble the request.
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.JOB,
            (short) tpduGenerator.getAndIncrement(), Collections.singletonList(readVarParameter),
            Collections.emptyList(), msg);

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);
    }

    private void encodeWriteRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        List<VarParameterItem> parameterItems = new LinkedList<>();
        List<VarPayloadItem> payloadItems = new LinkedList<>();

        PlcWriteRequest writeRequest = (PlcWriteRequest) msg.getRequest();
        for (String fieldName : writeRequest.getFieldNames()) {
            PlcField field = writeRequest.getField(fieldName);
            if (!(field instanceof S7Field)) {
                throw new PlcException("The field should have been of type S7Field");
            }
            S7Field s7Field = (S7Field) field;
            if(!(writeRequest instanceof DefaultPlcWriteRequest)) {
                throw new PlcException("The writeRequest should have been of type DefaultPlcWriteRequest");
            }
            BaseDefaultFieldItem fieldItem = ((DefaultPlcWriteRequest) writeRequest).getFieldItem(fieldName);

            // The number of elements provided in the request must match the number defined in the field, or
            // bad things are going to happen.
            // An exception is STRINGS, as they are implemented as byte arrays
            if (s7Field.getDataType() != TransportSize.STRING &&
                writeRequest.getNumberOfValues(fieldName) != s7Field.getNumElements()) {
                throw new PlcException("The number of values provided doesn't match the number specified by the field.");
            }
            VarParameterItem varParameterItem = new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Field.getMemoryArea(),
                s7Field.getDataType(),
                s7Field.getNumElements(), s7Field.getBlockNumber(), s7Field.getByteOffset(), (byte) s7Field.getBitOffset());
            parameterItems.add(varParameterItem);

            DataTransportSize dataTransportSize = s7Field.getDataType().getDataTransportSize();

            // TODO: Checkout if the payload items are sort of a flatMap of all request items.
            byte[] byteData;
            switch(s7Field.getDataType()) {
                // -----------------------------------------
                // Bit
                // -----------------------------------------
                case BOOL:
                    byteData = encodeWriteRequestBitField(fieldItem);
                    break;
                // -----------------------------------------
                // Signed integer values
                // -----------------------------------------
                case BYTE:
                case SINT:
                case CHAR:  // 1 byte
                    byteData = encodeWriteRequestByteField(fieldItem, true);
                    break;
                case WORD:
                case INT:
                case WCHAR:  // 2 byte (16 bit)
                    byteData = encodeWriteRequestShortField(fieldItem, true);
                    break;
                case DWORD:
                case DINT:  // 4 byte (32 bit)
                    byteData = encodeWriteRequestIntegerField(fieldItem, true);
                    break;
                case LWORD:
                case LINT:  // 8 byte (64 bit)
                    byteData = encodeWriteRequestLongField(fieldItem, true);
                    break;
                // -----------------------------------------
                // Unsigned integer values
                // -----------------------------------------
                // 8 bit:
                case USINT:
                    byteData = encodeWriteRequestByteField(fieldItem, false);
                    break;
                // 16 bit:
                case UINT:
                    byteData = encodeWriteRequestShortField(fieldItem, false);
                    break;
                // 32 bit:
                case UDINT:
                    byteData = encodeWriteRequestIntegerField(fieldItem, false);
                    break;
                // 64 bit:
                case ULINT:
                    byteData = encodeWriteRequestLongField(fieldItem, false);
                    break;
                // -----------------------------------------
                // Floating point values
                // -----------------------------------------
                case REAL:
                    byteData = encodeWriteRequestFloatField(fieldItem);
                    break;
                case LREAL:
                    byteData = encodeWriteRequestDoubleField(fieldItem);
                    break;
                // -----------------------------------------
                // Characters & Strings
                // -----------------------------------------
                case STRING:
                    byteData = encodeWriteRequestStringField(fieldItem, false);
                    break;
                case WSTRING:
                    byteData = encodeWriteRequestStringField(fieldItem, true);
                    break;
                default:
                    throw new PlcProtocolException("Unsupported type " + s7Field.getDataType());
            }

            VarPayloadItem varPayloadItem = new VarPayloadItem(
                DataTransportErrorCode.RESERVED, dataTransportSize, byteData);

            payloadItems.add(varPayloadItem);
        }
        VarParameter writeVarParameter = new VarParameter(ParameterType.WRITE_VAR, parameterItems);
        VarPayload writeVarPayload = new VarPayload(ParameterType.WRITE_VAR, payloadItems);

        // Assemble the request.
        S7RequestMessage s7WriteRequest = new S7RequestMessage(MessageType.JOB,
            (short) tpduGenerator.getAndIncrement(), Collections.singletonList(writeVarParameter),
            Collections.singletonList(writeVarPayload), msg);

        requests.put(s7WriteRequest.getTpduReference(), msg);

        out.add(s7WriteRequest);
    }

    byte[] encodeWriteRequestBitField(BaseDefaultFieldItem fieldItem) {
        int numBytes = fieldItem.getNumberOfValues() >> 3 / 8;
        byte[] byteData = new byte[numBytes];
        BitSet bitSet = new BitSet();
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            bitSet.set(i, fieldItem.getBoolean(i));
        }
        byte[] src = bitSet.toByteArray();
        System.arraycopy(src, 0, byteData, 0, Math.min(src.length, numBytes));
        return byteData;
    }

    byte[] encodeWriteRequestByteField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues();
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.put(fieldItem.getByte(i));
            } else {
                buffer.put((byte) (short) fieldItem.getShort(i));
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestShortField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 2;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.putShort(fieldItem.getShort(i));
            } else {
                buffer.putShort((short) (int) fieldItem.getInteger(i));
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestIntegerField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 4;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.putInt(fieldItem.getInteger(i));
            } else {
                buffer.putInt((int) (long) fieldItem.getLong(i));
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestLongField(BaseDefaultFieldItem fieldItem, boolean signed) {
        int numBytes = fieldItem.getNumberOfValues() * 8;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            if(signed) {
                buffer.putLong(fieldItem.getLong(i));
            } else {
                // TODO: Implement this ...
            }
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestFloatField(BaseDefaultFieldItem fieldItem) {
        int numBytes = fieldItem.getNumberOfValues() * 4;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putFloat(fieldItem.getFloat(i));
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestDoubleField(BaseDefaultFieldItem fieldItem) {
        int numBytes = fieldItem.getNumberOfValues() * 8;
        ByteBuffer buffer = ByteBuffer.allocate(numBytes);
        for (int i = 0; i < fieldItem.getNumberOfValues(); i++) {
            buffer.putDouble(fieldItem.getDouble(i));
        }
        return buffer.array();
    }

    byte[] encodeWriteRequestStringField(BaseDefaultFieldItem fieldItem, boolean isUtf16) {
        // TODO: Implement this ...
        return new byte[0];
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext ctx, S7Message msg, List<Object> out) throws PlcException {
        // We're currently just expecting responses.
        if (!(msg instanceof S7ResponseMessage)) {
            return;
        }
        S7ResponseMessage responseMessage = (S7ResponseMessage) msg;
        short tpduReference = responseMessage.getTpduReference();
        if (requests.containsKey(tpduReference)) {
            // As every response has a matching request, get this request based on the tpdu.
            PlcRequestContainer requestContainer = requests.remove(tpduReference);
            PlcRequest request = requestContainer.getRequest();

            // Handle the response.
            PlcResponse response = null;
            if (request instanceof PlcReadRequest) {
                response = decodeReadResponse(responseMessage, requestContainer);
            } else if (request instanceof PlcWriteRequest) {
                response = decodeWriteResponse(responseMessage, requestContainer);
            }

            // Confirm the response being handled.
            if (response != null) {
                requestContainer.getResponseFuture().complete(response);
            }
        }
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeReadResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        InternalPlcReadRequest plcReadRequest = (InternalPlcReadRequest) requestContainer.getRequest();

        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));

        // TODO: Checkout if the payload items are sort of a flatMap of all request items.

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getNumberOfFields() != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> values = new HashMap<>();
        List<VarPayloadItem> payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
            VarPayloadItem payloadItem = payloadItems.get(index);

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            BaseDefaultFieldItem fieldItem = null;
            ByteBuf data = Unpooled.wrappedBuffer(payloadItem.getData());
            if (responseCode == PlcResponseCode.OK) {
                try {
                    switch (field.getDataType()) {
                        // -----------------------------------------
                        // Bit
                        // -----------------------------------------
                        case BOOL:
                            fieldItem = decodeReadResponseBitField(field, data);
                            break;
                        // -----------------------------------------
                        // Bit-strings
                        // -----------------------------------------
                        case BYTE:  // 1 byte
                            fieldItem = decodeReadResponseByteBitStringField(field, data);
                            break;
                        case WORD:  // 2 byte (16 bit)
                            fieldItem = decodeReadResponseShortBitStringField(field, data);
                            break;
                        case DWORD:  // 4 byte (32 bit)
                            fieldItem = decodeReadResponseIntegerBitStringField(field, data);
                            break;
                        case LWORD:  // 8 byte (64 bit)
                            fieldItem = decodeReadResponseLongBitStringField(field, data);
                            break;
                        // -----------------------------------------
                        // Integers
                        // -----------------------------------------
                        // 8 bit:
                        case SINT:
                            fieldItem = decodeReadResponseSignedByteField(field, data);
                            break;
                        case USINT:
                            fieldItem = decodeReadResponseUnsignedByteField(field, data);
                            break;
                        // 16 bit:
                        case INT:
                            fieldItem = decodeReadResponseSignedShortField(field, data);
                            break;
                        case UINT:
                            fieldItem = decodeReadResponseUnsignedShortField(field, data);
                            break;
                        // 32 bit:
                        case DINT:
                            fieldItem = decodeReadResponseSignedIntegerField(field, data);
                            break;
                        case UDINT:
                            fieldItem = decodeReadResponseUnsignedIntegerField(field, data);
                            break;
                        // 64 bit:
                        case LINT:
                            fieldItem = decodeReadResponseSignedLongField(field, data);
                            break;
                        case ULINT:
                            fieldItem = decodeReadResponseUnsignedLongField(field, data);
                            break;
                        // -----------------------------------------
                        // Floating point values
                        // -----------------------------------------
                        case REAL:
                            fieldItem = decodeReadResponseFloatField(field, data);
                            break;
                        case LREAL:
                            fieldItem = decodeReadResponseDoubleField(field, data);
                            break;
                        // -----------------------------------------
                        // Characters & Strings
                        // -----------------------------------------
                        case CHAR: // 1 byte (8 bit)
                            fieldItem = decodeReadResponseFixedLengthStringField(1, false, data);
                            break;
                        case WCHAR: // 2 byte
                            fieldItem = decodeReadResponseFixedLengthStringField(1, true, data);
                            break;
                        case STRING:
                            fieldItem = decodeReadResponseVarLengthStringField(false, data);
                            break;
                        case WSTRING:
                            fieldItem = decodeReadResponseVarLengthStringField(true, data);
                            break;
                        // -----------------------------------------
                        // TIA Date-Formats
                        // -----------------------------------------
                        case DATE_AND_TIME:
                            fieldItem = decodeReadResponseDateAndTime(field, data);
                            break;
                        case TIME_OF_DAY:
                            fieldItem = decodeReadResponseTimeOfDay(field, data);
                            break;
                        case DATE:
                            fieldItem = decodeReadResponseDate(field, data);
                            break;
                        default:
                            throw new PlcProtocolException("Unsupported type " + field.getDataType());
                    }
                }
                catch (Plc4XNettyException e){
                    logger.warn("Problem during casting of field {}: Exception: {}; FieldInformation: {}",fieldName,e.getMessage(),field);
                }
                catch (Exception e){
                    logger.warn("Some other error occurred casting field {}, FieldInformation: {}",fieldName, field,e);
                }
            }
            Pair<PlcResponseCode, BaseDefaultFieldItem> result = new ImmutablePair<>(responseCode, fieldItem);
            values.put(fieldName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
    }

    BaseDefaultFieldItem decodeReadResponseBitField(S7Field field, ByteBuf data) {
        Boolean[] booleans = readAllValues(Boolean.class, field, i -> data.readByte() != 0x00);
        return new DefaultBooleanFieldItem(booleans);
    }

    BaseDefaultFieldItem decodeReadResponseByteBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements()];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseShortBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 2];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseIntegerBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 4];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseLongBitStringField(S7Field field, ByteBuf data) {
        byte[] bytes = new byte[field.getNumElements() * 8];
        data.readBytes(bytes);
        return decodeBitStringField(bytes);
    }

    BaseDefaultFieldItem decodeBitStringField(byte[] bytes) {
        BitSet bitSet = BitSet.valueOf(bytes);
        Boolean[] booleanValues = new Boolean[8 * bytes.length];
        int k = 0;
        for(int i = bytes.length - 1; i >= 0; i--) {
            for(int j = 0; j < 8; j++) {
                booleanValues[k++] = bitSet.get(8 * i + j);
            }
        }
        return new DefaultBooleanFieldItem(booleanValues);
    }

    BaseDefaultFieldItem decodeReadResponseSignedByteField(S7Field field, ByteBuf data) {
        Byte[] bytes = readAllValues(Byte.class, field, i -> data.readByte());
        return new DefaultByteFieldItem(bytes);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedByteField(S7Field field, ByteBuf data) {
        Short[] shorts = readAllValues(Short.class, field, i -> data.readUnsignedByte());
        return new DefaultShortFieldItem(shorts);
    }

    BaseDefaultFieldItem decodeReadResponseSignedShortField(S7Field field, ByteBuf data) {
        Short[] shorts = readAllValues(Short.class, field, i -> data.readShort());
        return new DefaultShortFieldItem(shorts);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedShortField(S7Field field, ByteBuf data) {
        Integer[] ints = readAllValues(Integer.class, field, i -> data.readUnsignedShort());
        return new DefaultIntegerFieldItem(ints);
    }

    BaseDefaultFieldItem decodeReadResponseSignedIntegerField(S7Field field, ByteBuf data) {
        Integer[] ints = readAllValues(Integer.class, field, i -> data.readInt());
        return new DefaultIntegerFieldItem(ints);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedIntegerField(S7Field field, ByteBuf data) {
        Long[] longs = readAllValues(Long.class, field, i -> data.readUnsignedInt());
        return new DefaultLongFieldItem(longs);
    }

    BaseDefaultFieldItem decodeReadResponseSignedLongField(S7Field field, ByteBuf data) {
        Long[] longs = readAllValues(Long.class, field, i -> data.readLong());
        return new DefaultLongFieldItem(longs);
    }

    BaseDefaultFieldItem decodeReadResponseUnsignedLongField(S7Field field, ByteBuf data) {
        BigInteger[] bigIntegers = readAllValues(BigInteger.class, field, i -> readUnsigned64BitInteger(data));
        return new DefaultBigIntegerFieldItem(bigIntegers);
    }

    BaseDefaultFieldItem decodeReadResponseFloatField(S7Field field, ByteBuf data) {
        Float[] floats = readAllValues(Float.class, field, i -> data.readFloat());
        return new DefaultFloatFieldItem(floats);
    }

    BaseDefaultFieldItem decodeReadResponseDoubleField(S7Field field, ByteBuf data) {
        Double[] doubles = readAllValues(Double.class, field, i -> data.readDouble());
        return new DefaultDoubleFieldItem(doubles);
    }

    BaseDefaultFieldItem decodeReadResponseFixedLengthStringField(int numChars, boolean isUtf16, ByteBuf data) {
        int numBytes = isUtf16 ? numChars * 2 : numChars;
        String stringValue = data.readCharSequence(numBytes, StandardCharsets.UTF_8).toString();
        return new DefaultStringFieldItem(stringValue);
    }

    BaseDefaultFieldItem decodeReadResponseVarLengthStringField(boolean isUtf16, ByteBuf data) {
        // Max length ... ignored.
        data.skipBytes(1);

        //reading out byte and transforming that to an unsigned byte within an integer, otherwise longer strings are failing
        byte currentLengthByte = data.readByte();
        int currentLength = currentLengthByte & 0xFF;
        return decodeReadResponseFixedLengthStringField(currentLength, isUtf16, data);
    }

    BaseDefaultFieldItem decodeReadResponseDateAndTime(S7Field field,ByteBuf data) {
        LocalDateTime[] localDateTimes = readAllValues(LocalDateTime.class,field, i -> readDateAndTime(data));
        return new DefaultLocalDateTimeFieldItem(localDateTimes);
    }

    BaseDefaultFieldItem decodeReadResponseTimeOfDay(S7Field field,ByteBuf data) {
        LocalTime[] localTimes = readAllValues(LocalTime.class,field, i -> readTimeOfDay(data));
        return new DefaultLocalTimeFieldItem(localTimes);
    }

    BaseDefaultFieldItem decodeReadResponseDate(S7Field field,ByteBuf data) {
        LocalDate[] localTimes = readAllValues(LocalDate.class,field, i -> readDate(data));
        return new DefaultLocalDateFieldItem(localTimes);
    }

    // Returns a 32 bit unsigned value : from 0 to 4294967295 (2^32-1)
    public static int getUDIntAt(byte[] buffer, int pos) {
        int result;
        result = buffer[pos] & 0x0FF;
        result <<= 8;
        result += buffer[pos + 1] & 0x0FF;
        result <<= 8;
        result += buffer[pos + 2] & 0x0FF;
        result <<= 8;
        result += buffer[pos + 3] & 0x0FF;
        return result;
    }

    private static <T> T[] readAllValues(Class<T> clazz, S7Field field, Function<Integer, T> extract) {
        try {
            return IntStream.rangeClosed(1, field.getNumElements())
                .mapToObj(extract::apply)
                .collect(Collectors.toList())
                .toArray((T[])Array.newInstance(clazz, 0));
        } catch (IndexOutOfBoundsException e) {
            throw new PlcRuntimeException("To few bytes in the buffer to read requested type", e);
        }
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeWriteResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        InternalPlcWriteRequest plcWriteRequest = (InternalPlcWriteRequest) requestContainer.getRequest();
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));

        // TODO: Checkout if the payload items are sort of a flatMap of all request items.

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getNumberOfFields() != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        Map<String, PlcResponseCode> values = new HashMap<>();
        List<VarPayloadItem> payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcWriteRequest.getFieldNames()) {
            VarPayloadItem payloadItem = payloadItems.get(index);

            // A write response contains only the return code for every item.
            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());

            values.put(fieldName, responseCode);
            index++;
        }

        return new DefaultPlcWriteResponse(plcWriteRequest, values);
    }

    private PlcResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return PlcResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return PlcResponseCode.OK;
            case NOT_FOUND:
                return PlcResponseCode.NOT_FOUND;
            case INVALID_ADDRESS:
                return PlcResponseCode.INVALID_ADDRESS;
            case DATA_TYPE_NOT_SUPPORTED:
                return PlcResponseCode.INVALID_DATATYPE;
            default:
                return PlcResponseCode.INTERNAL_ERROR;
        }
    }

    private static BigInteger readUnsignedLong(ByteBuf data) {
        // as there is no unsigned long primitive, we have to switch to
        // BigDecimal and manually convert the bytes to a BigDecimal.
        // In order to be unsigned 4 bytes, we create an array of 5 bytes
        // where the 5th byte is set to 0. The most significant bit being
        // 0 we are guaranteed to interpret the input a positive value.
        byte[] bytes = new byte[5];
        // Set the first byte to 0
        bytes[0] = 0;
        // Read the next 4 bytes into the rest.
        data.readBytes(bytes, 1, 4);
        return new BigInteger(bytes);
    }

    private static BigInteger readSigned64BitInteger(ByteBuf data) {
        byte[] bytes = new byte[8];
        data.readBytes(bytes, 0, 8);
        return new BigInteger(bytes);
    }

    private static BigInteger readUnsigned64BitInteger(ByteBuf data) {
        byte[] bytes = new byte[9];
        // Set the first byte to 0
        bytes[0] = 0;
        // Read the next 8 bytes into the rest.
        data.readBytes(bytes, 1, 8);
        return new BigInteger(bytes);
    }

    LocalDateTime readDateAndTime(ByteBuf data) {
        //per definition for Date_And_Time only the first 6 bytes are used

        int year=convertByteToBcd(data.readByte());
        int month=convertByteToBcd(data.readByte());
        int day=convertByteToBcd(data.readByte());
        int hour=convertByteToBcd(data.readByte());
        int minute=convertByteToBcd(data.readByte());
        int second=convertByteToBcd(data.readByte());
        //skip the last 2 bytes no information present
        data.readByte();
        data.readByte();

        //data-type ranges from 1990 up to 2089
        if(year>=90){
            year+=1900;
        }
        else{
            year+=2000;
        }

        return LocalDateTime.of(year,month,day,hour,minute,second);
    }

    LocalTime readTimeOfDay(ByteBuf data) {
        //per definition for Date_And_Time only the first 6 bytes are used

        int millisSinsMidnight = data.readInt();


        return LocalTime.now().withHour(0).withMinute(0).withSecond(0).withNano(0).plus(millisSinsMidnight, ChronoUnit.MILLIS);

    }

    LocalDate readDate(ByteBuf data) {
        //per definition for Date_And_Time only the first 6 bytes are used

        int daysSince1990 = data.readUnsignedShort();

        System.out.println(daysSince1990);
        return LocalDate.now().withYear(1990).withDayOfMonth(1).withMonth(1).plus(daysSince1990, ChronoUnit.DAYS);

    }

    /**
     * converts incoming byte to an integer regarding used BCD format
     * @param incomingByte
     * @return converted BCD number
     */
    private static int convertByteToBcd(byte incomingByte) {
        int dec = (incomingByte >> 4) * 10;
        return dec + (incomingByte & 0x0f);
    }

}
