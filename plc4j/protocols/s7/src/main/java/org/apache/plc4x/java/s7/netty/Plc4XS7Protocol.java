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
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcIoException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolPayloadTooBigException;
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

import java.io.IOException;
import java.math.BigInteger;
import java.nio.charset.Charset;
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

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
                throw new PlcException("The field should have been of type S7Field");
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

            // The number of elements provided in the request must match the number defined in the field, or
            // bad things are going to happen.
            if (writeRequest.getNumberOfValues(fieldName) != s7Field.getNumElements()) {
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
            /*switch(s7Field.getDataType()) {
                // -----------------------------------------
                // Bit
                // -----------------------------------------
                case BOOL: {
                    if()
                    byte byteValue = data.readByte();
                    fieldItem = new BooleanFieldItem(byteValue != 0x00);
                    break;
                }
                // -----------------------------------------
                // Bit-strings
                // -----------------------------------------
                case BYTE: { // 1 byte
                    Long longValue = (long) data.readByte();
                    // TODO: Implement this ...
                    break;
                }
                case WORD: { // 2 byte (16 bit)
                    Long longValue = (long) data.readShort();
                    // TODO: Implement this ...
                    break;
                }
                case DWORD: { // 4 byte (32 bit)
                    Long longValue = (long) data.readInt();
                    // TODO: Implement this ...
                    break;
                }
                case LWORD: { // 8 byte (64 bit)
                    Long longValue = data.readLong();
                    // TODO: Implement this ...
                    break;
                }
                // -----------------------------------------
                // Integers
                // -----------------------------------------
                // 8 bit:
                case SINT: {
                    Long longValue = (long) data.readShort();
                    fieldItem = new IntegerFieldItem(longValue);
                    break;
                }
                case USINT: {
                    Long longValue = (long) data.readUnsignedShort();
                    fieldItem = new IntegerFieldItem(longValue);
                    break;
                }
                // 16 bit:
                case INT: {
                    Long longValue = (long) data.readInt();
                    fieldItem = new IntegerFieldItem(longValue);
                    break;
                }
                case UINT: {
                    Long longValue = data.readUnsignedInt();
                    fieldItem = new IntegerFieldItem(longValue);
                    break;
                }
                // 32 bit:
                case DINT: {
                    Long longValue = data.readLong();
                    fieldItem = new IntegerFieldItem(longValue);
                    break;
                }

                case UDINT: {
                    BigInteger bigIntegerValue = readUnsignedLong(data);
                    fieldItem = new BigIntegerFieldItem(bigIntegerValue);
                    break;
                }
                // 64 bit:
                case LINT: {
                    BigInteger bigIntegerValue = readSigned64BitInteger(data);
                    fieldItem = new BigIntegerFieldItem(bigIntegerValue);
                    break;
                }
                case ULINT: {
                    BigInteger bigIntegerValue = readUnsigned64BitInteger(data);
                    fieldItem = new BigIntegerFieldItem(bigIntegerValue);
                    break;
                }
                // -----------------------------------------
                // Floating point values
                // -----------------------------------------
                case REAL: {
                    double doubleValue = data.readFloat();
                    fieldItem = new FloatingPointFieldItem(doubleValue);
                    break;
                }
                case LREAL: {
                    double doubleValue = data.readDouble();
                    fieldItem = new FloatingPointFieldItem(doubleValue);
                    break;
                }
                // -----------------------------------------
                // Characters & Strings
                // -----------------------------------------
                case CHAR: { // 1 byte (8 bit)
                    // TODO: Double check, if this is ok?
                    String stringValue = data.readCharSequence(1, Charset.forName("UTF-8")).toString();
                    fieldItem = new StringFieldItem(stringValue);
                    break;
                }
                case WCHAR: { // 2 byte
                    // TODO: Double check, if this is ok?
                    String stringValue = data.readCharSequence(2, Charset.forName("UTF-16")).toString();
                    fieldItem = new StringFieldItem(stringValue);
                    break;
                }
                case STRING: {
                    // Max length ... ignored.
                    data.readByte();
                    byte actualLength = data.readByte();
                    // TODO: Double check, if this is ok?
                    String stringValue = data.readCharSequence(actualLength, Charset.forName("UTF-8")).toString();
                    fieldItem = new StringFieldItem(stringValue);
                    break;
                }
                case WSTRING: {
                    // Max length ... ignored.
                    data.readByte();
                    byte actualLength = data.readByte();
                    // TODO: Double check, if this is ok?
                    String stringValue = data.readCharSequence(
                        actualLength * 2, Charset.forName("UTF-16")).toString();
                    fieldItem = new StringFieldItem(stringValue);
                    break;
                }
                default:
                    throw new PlcProtocolException("Unsupported type " + field.getDataType());
            }

            VarPayloadItem varPayloadItem = new VarPayloadItem(
                DataTransportErrorCode.RESERVED, dataTransportSize, writeRequest.getValues(fieldName)[0]);

            payloadItems.add(varPayloadItem);*/
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

        Map<String, Pair<PlcResponseCode, FieldItem>> values = new HashMap<>();
        List<VarPayloadItem> payloadItems = payload.getItems();
        int index = 0;
        for (String fieldName : plcReadRequest.getFieldNames()) {
            S7Field field = (S7Field) plcReadRequest.getField(fieldName);
            VarPayloadItem payloadItem = payloadItems.get(index);

            PlcResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());
            FieldItem fieldItem = null;
            ByteBuf data = Unpooled.wrappedBuffer(payloadItem.getData());
            if (responseCode == PlcResponseCode.OK) {
                switch (field.getDataType()) {
                    // -----------------------------------------
                    // Bit
                    // -----------------------------------------
                    case BOOL: {
                        byte byteValue = data.readByte();
                        fieldItem = new BooleanFieldItem(byteValue != 0x00);
                        break;
                    }
                    // -----------------------------------------
                    // Bit-strings
                    // -----------------------------------------
                    case BYTE: { // 1 byte
                        Long longValue = (long) data.readByte();
                        // TODO: Implement this ...
                        break;
                    }
                    case WORD: { // 2 byte (16 bit)
                        Long longValue = (long) data.readShort();
                        // TODO: Implement this ...
                        break;
                    }
                    case DWORD: { // 4 byte (32 bit)
                        Long longValue = (long) data.readInt();
                        // TODO: Implement this ...
                        break;
                    }
                    case LWORD: { // 8 byte (64 bit)
                        Long longValue = data.readLong();
                        // TODO: Implement this ...
                        break;
                    }
                    // -----------------------------------------
                    // Integers
                    // -----------------------------------------
                    // 8 bit:
                    case SINT: {
                        Long longValue = (long) data.readShort();
                        fieldItem = new IntegerFieldItem(longValue);
                        break;
                    }
                    case USINT: {
                        Long longValue = (long) data.readUnsignedShort();
                        fieldItem = new IntegerFieldItem(longValue);
                        break;
                    }
                    // 16 bit:
                    case INT: {
                        Long longValue = (long) data.readInt();
                        fieldItem = new IntegerFieldItem(longValue);
                        break;
                    }
                    case UINT: {
                        Long longValue = data.readUnsignedInt();
                        fieldItem = new IntegerFieldItem(longValue);
                        break;
                    }
                    // 32 bit:
                    case DINT: {
                        Long longValue = data.readLong();
                        fieldItem = new IntegerFieldItem(longValue);
                        break;
                    }
                    case UDINT: {
                        BigInteger bigIntegerValue = readUnsignedLong(data);
                        fieldItem = new BigIntegerFieldItem(bigIntegerValue);
                        break;
                    }
                    // 64 bit:
                    case LINT: {
                        BigInteger bigIntegerValue = readSigned64BitInteger(data);
                        fieldItem = new BigIntegerFieldItem(bigIntegerValue);
                        break;
                    }
                    case ULINT: {
                        BigInteger bigIntegerValue = readUnsigned64BitInteger(data);
                        fieldItem = new BigIntegerFieldItem(bigIntegerValue);
                        break;
                    }
                    // -----------------------------------------
                    // Floating point values
                    // -----------------------------------------
                    case REAL: {
                        double doubleValue = data.readFloat();
                        fieldItem = new FloatingPointFieldItem(doubleValue);
                        break;
                    }
                    case LREAL: {
                        double doubleValue = data.readDouble();
                        fieldItem = new FloatingPointFieldItem(doubleValue);
                        break;
                    }
                    // -----------------------------------------
                    // Characters & Strings
                    // -----------------------------------------
                    case CHAR: { // 1 byte (8 bit)
                        // TODO: Double check, if this is ok?
                        String stringValue = data.readCharSequence(1, Charset.forName("UTF-8")).toString();
                        fieldItem = new StringFieldItem(stringValue);
                        break;
                    }
                    case WCHAR: { // 2 byte
                        // TODO: Double check, if this is ok?
                        String stringValue = data.readCharSequence(2, Charset.forName("UTF-16")).toString();
                        fieldItem = new StringFieldItem(stringValue);
                        break;
                    }
                    case STRING: {
                        // Max length ... ignored.
                        data.readByte();
                        byte actualLength = data.readByte();
                        // TODO: Double check, if this is ok?
                        String stringValue = data.readCharSequence(actualLength, Charset.forName("UTF-8")).toString();
                        fieldItem = new StringFieldItem(stringValue);
                        break;
                    }
                    case WSTRING: {
                        // Max length ... ignored.
                        data.readByte();
                        byte actualLength = data.readByte();
                        // TODO: Double check, if this is ok?
                        String stringValue = data.readCharSequence(
                            actualLength * 2, Charset.forName("UTF-16")).toString();
                        fieldItem = new StringFieldItem(stringValue);
                        break;
                    }
                    default:
                        throw new PlcProtocolException("Unsupported type " + field.getDataType());
                }
            }
            Pair<PlcResponseCode, FieldItem> result = new ImmutablePair<>(responseCode, fieldItem);
            values.put(fieldName, result);
            index++;
        }

        return new DefaultPlcReadResponse(plcReadRequest, values);
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

}
