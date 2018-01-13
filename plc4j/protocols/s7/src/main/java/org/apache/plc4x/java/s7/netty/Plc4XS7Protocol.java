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
import io.netty.handler.codec.MessageToMessageCodec;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
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
import org.apache.plc4x.java.s7.model.S7Address;
import org.apache.plc4x.java.s7.model.S7BitAddress;
import org.apache.plc4x.java.s7.model.S7DataBlockAddress;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;

import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4XS7Protocol extends MessageToMessageCodec<S7Message, PlcRequestContainer> {

    private static final AtomicInteger tpduGenerator = new AtomicInteger(1);

    private Map<Short, PlcRequestContainer> requests;

    public Plc4XS7Protocol() {
        this.requests = new HashMap<>();
    }

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
        PlcRequest request = msg.getRequest();
        if (request instanceof PlcReadRequest) {
            encodeREadRequest(msg, out);
        } else if (request instanceof PlcWriteRequest) {
            encodeWriteRequest(msg, out);
        }
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
            Collections.singletonList(writeVarPayload));

        requests.put(s7WriteRequest.getTpduReference(), msg);

        out.add(s7WriteRequest);
    }

    private void encodeREadRequest(PlcRequestContainer msg, List<Object> out) throws PlcException {
        List<VarParameterItem> parameterItems = new LinkedList<>();

        PlcReadRequest readRequest = (PlcReadRequest) msg.getRequest();
        encodeParameterItems(parameterItems, readRequest);
        VarParameter readVarParameter = new VarParameter(ParameterType.READ_VAR, parameterItems);

        // Assemble the request.
        S7RequestMessage s7ReadRequest = new S7RequestMessage(MessageType.JOB,
            (short) tpduGenerator.getAndIncrement(), Collections.singletonList(readVarParameter),
            Collections.emptyList());

        requests.put(s7ReadRequest.getTpduReference(), msg);

        out.add(s7ReadRequest);
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
                response = decodeReadRequest(responseMessage, requestContainer);
            }
            else if (request instanceof PlcWriteRequest) {
                response = decodeWriteRequest(responseMessage, requestContainer);
            }

            // Confirm the response being handled.
            if (response != null) {
                requestContainer.getResponseFuture().complete(response);
            }
        }
    }

    @SuppressWarnings("unchecked")
    private PlcResponse decodeWriteRequest(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        PlcResponse response;
        PlcWriteRequest plcWriteRequest = (PlcWriteRequest) requestContainer.getRequest();
        List<WriteResponseItem<?>> responseItems = new LinkedList<>();
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));
        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcWriteRequest.getRequestItems().size() != payload.getPayloadItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }
        List<VarPayloadItem> payloadItems = payload.getPayloadItems();
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

    @SuppressWarnings("unchecked")
    private PlcResponse decodeReadRequest(S7ResponseMessage responseMessage, PlcRequestContainer requestContainer) throws PlcProtocolException {
        PlcResponse response;
        PlcReadRequest plcReadRequest = (PlcReadRequest) requestContainer.getRequest();

        List<ReadResponseItem<?>> responseItems = new LinkedList<>();
        VarPayload payload = responseMessage.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("No VarPayload supplied"));
        // If the numbers of items don't match, we're in big trouble as the only
        // way to know how to interpret the responses is by aligning them with the
        // items from the request as this information is not returned by the PLC.
        if (plcReadRequest.getRequestItems().size() != payload.getPayloadItems().size()) {
            throw new PlcProtocolException(
                "The number of requested items doesn't match the number of returned items");
        }
        List<VarPayloadItem> payloadItems = payload.getPayloadItems();
        final int noPayLoadItems = payloadItems.size();
        for (int i = 0; i < noPayLoadItems; i++) {
            VarPayloadItem payloadItem = payloadItems.get(i);

            // Get the request item for this payload item
            ReadRequestItem requestItem = plcReadRequest.getRequestItems().get(i);

            ResponseCode responseCode = decodeResponseCode(payloadItem.getReturnCode());

            ReadResponseItem responseItem;
            // Something went wrong.
            if (responseCode != ResponseCode.OK) {
                responseItem = new ReadResponseItem<>(requestItem, responseCode, null);
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

    ////////////////////////////////////////////////////////////////////////////////
    // Encoding helpers.
    ////////////////////////////////////////////////////////////////////////////////

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
            return null;
        } else if (datatype == Float.class) {
            return DataTransportSize.REAL;
        } else if (datatype == Integer.class) {
            return DataTransportSize.BYTE_WORD_DWORD;
        } else if (datatype == String.class) {
            return DataTransportSize.OCTET_STRING;
        }
        return null;
    }

    private byte[] encodeData(Object[] values) {
        final int length = values.length;
        if (length == 0) {
            return new byte[]{};
        }
        byte[] result = null;
        Class valueType = values[0].getClass();
        if (valueType == Boolean.class) {
            result = encodeBoolean(values, length);
        } else if (valueType == Byte.class) {
            result = encodeByte(values, length);
        } else if (valueType == Short.class) {
            result = encodeShort(values, length);
        } else if (valueType == Integer.class) {
            result = encodeInteger(values, length);
        } else if (valueType == Calendar.class) {
            result = null;
        } else if (valueType == Float.class) {
            result = encodeFloat(values, length);
        } else if (valueType == String.class) {
            result = encodeString(values, length);
        }
        return result;
    }

    private byte[] encodeString(Object[] values, int length) {
        byte[] result;
        int size = 0;
        for (Object value : values) {
            size = size + ((String) value).length();
        }
        result = new byte[size + length];
        int j = 0;
        for (Object value : values) {
            String str = (String) value;
            for (int i = 0; i < str.length(); i++) {
                result[j++] = (byte) str.charAt(i);
            }
            result[j++] = (byte) 0x0;
        }
        return result;
    }

    private byte[] encodeFloat(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 4];
        for (int i = 0; i < length; i++) {
            float floatValue = (float) values[i];
            int intValue = Float.floatToIntBits(floatValue);
            result[i * 4] = (byte) ((intValue & 0xff000000) >> 24);
            result[(i * 4) + 1] = (byte) ((intValue & 0x00ff0000) >> 16);
            result[(i * 4) + 2] = (byte) ((intValue & 0x0000ff00) >> 8);
            result[(i * 4) + 3] = (byte) (intValue & 0xff);
        }
        return result;
    }

    private byte[] encodeInteger(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 4];
        for (int i = 0; i < length; i++) {
            int intValue = (int) values[i];
            result[i * 4] = (byte) ((intValue & 0xff000000) >> 24);
            result[(i * 4) + 1] = (byte) ((intValue & 0x00ff0000) >> 16);
            result[(i * 4) + 2] = (byte) ((intValue & 0x0000ff00) >> 8);
            result[(i * 4) + 3] = (byte) (intValue & 0xff);
        }
        return result;
    }

    private byte[] encodeShort(Object[] values, int length) {
        byte[] result;
        result = new byte[length * 2];
        for (int i = 0; i < length; i++) {
            short intValue = (short) values[i];
            result[i * 2] = (byte) ((intValue & 0xff00) >> 8);
            result[(i * 2) + 1] = (byte) (intValue & 0xff);
        }
        return result;
    }

    private byte[] encodeByte(Object[] values, int length) {
        byte[] result;
        result = new byte[length];
        for (int i = 0; i < length; i++) {
            result[i] = (byte) values[i];
        }
        return result;
    }

    private byte[] encodeBoolean(Object[] values, int length) {
        byte[] result;// TODO: Check if this is true and the result is not Math.ceil(values.lenght / 8)
        result = new byte[length];
        for (int i = 0; i < length; i++) {
            result[i] = (byte) (((Boolean) values[i]) ? 0x01 : 0x00);
        }
        return result;
    }

    ////////////////////////////////////////////////////////////////////////////////
    // Decoding helpers.
    ////////////////////////////////////////////////////////////////////////////////

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

    @SuppressWarnings("unchecked")
    private <T> List<T> decodeData(Class<T> datatype, byte[] s7Data) throws PlcProtocolException {

        List<Object> result = new LinkedList<>();
        int i = 0;
        final int length = s7Data.length;
        while (i < length) {
            if (datatype == Boolean.class) {
                result.add((s7Data[i] & 0x01) == 0x01);
                i += 1;
            } else if (datatype == Byte.class) {
                result.add(s7Data[i]);
                i += 1;
            } else if (datatype == Short.class) {
                result.add((short) (((s7Data[i] & 0xff) << 8) | (s7Data[i + 1] & 0xff)));
                i += 2;
            } else if (datatype == Integer.class) {
                result.add(((s7Data[i] & 0xff) << 24) | ((s7Data[i + 1] & 0xff) << 16) |
                    ((s7Data[i + 2] & 0xff) << 8) | (s7Data[i + 3] & 0xff));
                i += 4;
            } else if (datatype == Float.class) {
                // Description of the Real number format:
                // https://www.sps-lehrgang.de/zahlenformate-step7/#c144
                // https://de.wikipedia.org/wiki/IEEE_754
                int intValue = ((s7Data[i] & 0xff) << 24) | ((s7Data[i + 1] & 0xff) << 16) |
                    ((s7Data[i + 2] & 0xff) << 8) | (s7Data[i + 3] & 0xff);
                result.add(Float.intBitsToFloat(intValue));
                i += 4;
            } else if (datatype == String.class) {
                StringBuilder builder = new StringBuilder();
                while (s7Data[i] != (byte) 0x0 && i < length) {
                    builder.append((char) s7Data[i]);
                    i++;
                }
                i++; // skip terminating character
                result.add(builder.toString());
            } else {
                throw new PlcProtocolException("Unsupported datatype " + datatype.getSimpleName());
            }
        }
        return (List<T>) result;
    }

}
