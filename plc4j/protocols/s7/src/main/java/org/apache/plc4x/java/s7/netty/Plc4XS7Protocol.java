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

import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.api.exceptions.*;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.messages.items.WriteResponseItem;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.apache.plc4x.java.api.types.ResponseCode;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.events.ConnectedEvent;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
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
import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

import static org.apache.plc4x.java.s7.netty.util.S7TypeDecoder.decodeData;
import static org.apache.plc4x.java.s7.netty.util.S7TypeEncoder.encodeData;

/**
 * This layer transforms between {@link PlcRequestContainer}s {@link S7Message}s.
 * And stores all "in-flight" requests in an internal structure ({@link Plc4XS7Protocol#requests}).
 *
 * While sending a request, a {@link S7RequestMessage} is generated and send downstream (to the {@link S7Protocol}.
 *
 * When a {@link S7ResponseMessage} is received it takes the existing request container from its Map and finishes
 * the {@link PlcRequestContainer}s future with the {@link PlcResponse}.
 */
public class Plc4XS7Protocol extends PlcMessageToMessageCodec<S7Message, PlcRequestContainer> {

    private static final AtomicInteger tpduGenerator = new AtomicInteger(1);

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
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
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
        encodeParameterItems(parameterItems, readRequest);
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
        for (WriteRequestItem requestItem : writeRequest.getRequestItems()) {
            // Try to get the correct S7 transport size for the given data type.
            // (Map PLC4X data type to S7 data type)
            TransportSize transportSize = encodeTransportSize(requestItem.getDatatype());
            if (transportSize == null) {
                throw new PlcException("Unknown transport size for datatype " + requestItem.getDatatype());
            }

            // Depending on the address type, generate the corresponding type of request item.
            VarParameterItem varParameterItem = encodeVarParameterItem(
                requestItem.getAddress(), transportSize, requestItem.getValues().size());
            parameterItems.add(varParameterItem);

            DataTransportSize dataTransportSize = encodeDataTransportSize(requestItem.getDatatype());
            if (dataTransportSize == null) {
                throw new PlcException("Unknown data transport size for datatype " + requestItem.getDatatype());
            }

            VarPayloadItem varPayloadItem = new VarPayloadItem(
                DataTransportErrorCode.RESERVED, dataTransportSize, encodeData(requestItem.getValues().toArray()));

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

    private void encodeParameterItems(List<VarParameterItem> parameterItems, PlcReadRequest readRequest) throws PlcException {
        for (ReadRequestItem requestItem : readRequest.getRequestItems()) {
            // Try to get the correct S7 transport size for the given data type.
            // (Map PLC4X data type to S7 data type)
            TransportSize transportSize = encodeTransportSize(requestItem.getDatatype());
            if (transportSize == null) {
                throw new PlcException("Unknown transport size for datatype " + requestItem.getDatatype());
            }

            // Depending on the address type, generate the corresponding type of request item.
            VarParameterItem varParameterItem = encodeVarParameterItem(requestItem.getAddress(), transportSize, requestItem.getSize());
            parameterItems.add(varParameterItem);
        }
    }

    private VarParameterItem encodeVarParameterItem(Address address, TransportSize transportSize, int size) throws PlcProtocolException {
        // Depending on the address type, generate the corresponding type of request item.
        if (!(address instanceof S7Address)) {
            throw new PlcProtocolException("Can only use S7Address types on S7 connection");
        }
        S7Address s7Address = (S7Address) address;
        if (s7Address instanceof S7DataBlockAddress) {
            S7DataBlockAddress s7DataBlockAddress = (S7DataBlockAddress) s7Address;
            return new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                transportSize, (short) size,
                s7DataBlockAddress.getDataBlockNumber(), s7DataBlockAddress.getByteOffset(), (byte) 0);
        } else if (s7Address instanceof S7BitAddress) {
            S7BitAddress s7BitAddress = (S7BitAddress) s7Address;
            return new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                transportSize, (short) size, (short) 0,
                s7Address.getByteOffset(), s7BitAddress.getBitOffset());
        } else {
            return new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, s7Address.getMemoryArea(),
                transportSize, (short) size, (short) 0,
                s7Address.getByteOffset(), (byte) 0);
        }
    }

    private TransportSize encodeTransportSize(Class<?> datatype) {
        if (datatype == Boolean.class) {
            return TransportSize.BIT;
        } else if (datatype == Byte.class) {
            return TransportSize.BYTE;
        } else if (datatype == Short.class) {
            return TransportSize.WORD;
        } else if (datatype == Calendar.class) {
            return TransportSize.DATE_AND_TIME;
        } else if (datatype == Float.class) {
            return TransportSize.REAL;
        } else if (datatype == Double.class) {
            return TransportSize.REAL;
        } else if (datatype == Integer.class) {
            return TransportSize.DWORD;
        } else if (datatype == String.class) {
            return TransportSize.CHAR;
        }
        return null;
    }

    private DataTransportSize encodeDataTransportSize(Class<?> datatype) {
        if (datatype == Boolean.class) {
            return DataTransportSize.BIT;
        } else if (datatype == Byte.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == Short.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == Calendar.class) {
            // TODO: Decide what to do here ...
            throw new PlcNotImplementedException("Calender support in S7 not yet implemented");
        } else if (datatype == Float.class) {
            return DataTransportSize.REAL;
        } else if (datatype == Double.class) {
            return DataTransportSize.REAL;
        } else if (datatype == Integer.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == String.class) {
            return DataTransportSize.OCTET_STRING;
        }
        return null;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @SuppressWarnings("unchecked")
    @Override
    protected void decode(ChannelHandlerContext ctx, S7Message msg, List<Object> out) throws Exception {
        if (!(msg instanceof S7ResponseMessage)) {
            return;
        }
        S7ResponseMessage responseMessage = (S7ResponseMessage) msg;
        short tpduReference = responseMessage.getTpduReference();
        if (requests.containsKey(tpduReference)) {
            PlcRequestContainer requestContainer = requests.remove(tpduReference);
            PlcRequest request = requestContainer.getRequest();
            PlcResponse response = null;

            // Handle the response to a read request.
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
        PlcResponse response;
        PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();

        List<ReadResponseItem<?>> responseItems = new LinkedList<>();
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));

        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getRequestItems().size() != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }

        List<VarPayloadItem> payloadItems = payload.getItems();
        final int noPayLoadItems = payloadItems.size();
        for (int i = 0; i < noPayLoadItems; i++) {
            VarPayloadItem payloadItem = payloadItems.get(i);

            // Get the request item for this payload item
            ReadRequestItem requestItem = plcReadRequest.getRequestItems().get(i);

            ResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());

            ReadResponseItem responseItem;
            // Something went wrong.
            if (responseCode != ResponseCode.OK) {
                responseItem = new ReadResponseItem<>(requestItem, responseCode);
            }
            // All Ok.
            else {
                byte[] data = payloadItem.getData();
                Class<?> datatype = requestItem.getDatatype();
                List<?> value = decodeData(datatype, data);
                responseItem = new ReadResponseItem(requestItem, responseCode, value);
            }
            responseItems.add(responseItem);
        }
        if (plcReadRequest instanceof TypeSafePlcReadRequest) {
            response = new TypeSafePlcReadResponse((TypeSafePlcReadRequest) plcReadRequest, responseItems);
        } else {
            response = new PlcReadResponse(plcReadRequest, responseItems);
        }
        return response;
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeWriteResponse(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        PlcResponse response;
        PlcWriteRequest plcWriteRequest = (PlcWriteRequest) requestContainer.getRequest();
        List<WriteResponseItem<?>> responseItems = new LinkedList<>();
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));
        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getRequestItems().size() != payload.getItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }
        List<VarPayloadItem> payloadItems = payload.getItems();
        final int noPayLoadItems = payloadItems.size();
        for (int i = 0; i < noPayLoadItems; i++) {
            VarPayloadItem payloadItem = payloadItems.get(i);

            // Get the request item for this payload item
            WriteRequestItem requestItem = plcWriteRequest.getRequestItems().get(i);

            // A write response contains only the return code for every item.
            ResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());

            WriteResponseItem responseItem = new WriteResponseItem(requestItem, responseCode);
            responseItems.add(responseItem);
        }

        if (plcWriteRequest instanceof TypeSafePlcWriteRequest) {
            response = new TypeSafePlcWriteResponse((TypeSafePlcWriteRequest) plcWriteRequest, responseItems);
        } else {
            response = new PlcWriteResponse(plcWriteRequest, responseItems);
        }
        return response;
    }

    private ResponseCode decodeResponseCode(DataTransportErrorCode dataTransportErrorCode) {
        if (dataTransportErrorCode == null) {
            return ResponseCode.INTERNAL_ERROR;
        }
        switch (dataTransportErrorCode) {
            case OK:
                return ResponseCode.OK;
            case NOT_FOUND:
                return ResponseCode.NOT_FOUND;
            case INVALID_ADDRESS:
                return ResponseCode.INVALID_ADDRESS;
            default:
                return ResponseCode.INTERNAL_ERROR;
        }
    }

}
