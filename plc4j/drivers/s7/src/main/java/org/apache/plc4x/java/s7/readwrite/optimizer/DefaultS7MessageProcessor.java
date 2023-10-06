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
package org.apache.plc4x.java.s7.readwrite.optimizer;

import io.vavr.control.Either;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.s7.readwrite.*;

import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * While a SetupCommunication message is no problem, when reading or writing data,
 * situations could arise that have to be handled. The following situations have to
 * be handled:
 * - The number of request items is so big, that the resulting PDU would exceed the
 *   agreed upon PDU size: The request has to be split up into multiple requests.
 * - If large blocks of data are requested by request items the result of a request
 *   could exceed the PDU size: The requests has to be split up into multiple requests
 *   where each requests response doesn't exceed the PDU size.
 *
 * The following optimizations should be implemented:
 * - If blocks are read which are in near proximity to each other it could be better
 *   to replace multiple requests by one that includes multiple blocks.
 * - Rearranging the order of request items could reduce the number of needed PDUs.
 */
public class DefaultS7MessageProcessor implements S7MessageProcessor {

    private final AtomicInteger tpduRefGen;

    public static final int EMPTY_READ_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterReadVarRequest(
        Collections.emptyList()), null).getLengthInBytes();
    public static final int EMPTY_READ_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterReadVarResponse(
        (short) 0), new S7PayloadReadVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();
    public static final int EMPTY_WRITE_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterWriteVarRequest(
        Collections.emptyList()), new S7PayloadWriteVarRequest(Collections.emptyList())).getLengthInBytes();
    public static final int EMPTY_WRITE_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterWriteVarResponse(
        (short) 0), new S7PayloadWriteVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();

    public DefaultS7MessageProcessor(AtomicInteger tpduGenerator) {
        this.tpduRefGen = tpduGenerator;
    }

    @Override
    public Collection<S7MessageRequest> processRequest(S7MessageRequest request, int pduSize) throws PlcException {
        // The following considerations have to be taken into account:
        // - The size of all parameters and payloads of a message cannot exceed the negotiated PDU size
        // - When reading data, the size of the returned data cannot exceed the negotiated PDU size
        //
        // Examples:
        // - Size of the request exceeds the maximum
        //  When having a negotiated max PDU size of 256, the maximum size of individual addresses can be at most 18
        //  If more are sent, the S7 will respond with a frame error.
        // - Size of the response exceeds the maximum
        //  When reading two Strings of each 200 bytes length, the size of the request is ok, however the PLC would
        //  have to send back 400 bytes of String data, which would exceed the PDU size. In this case the first String
        //  is correctly returned, but for the second item the PLC will return a code of 0x03 = Access Denied
        // - A S7 device doesn't seem to accept more than one write item. So if we are doing write operations, we
        //  have to split that up into one message per written item. This also seems to affect arrays. So if
        //  an array of values is written, we have to also split up that array into single writes.

        S7Parameter parameter = request.getParameter();
        if (parameter instanceof S7ParameterReadVarRequest) {
            return processReadVarParameter(request, pduSize);
        }
        // If this is a write operation, split up every array item into single value items
        // and every item into a separate message.
        else if(parameter instanceof S7ParameterWriteVarRequest) {
            return processWriteVarParameter(request, pduSize);
        }

        return Collections.singletonList(request);
    }

    private Collection<S7MessageRequest> processReadVarParameter(S7MessageRequest request, int pduSize) {
        final S7ParameterReadVarRequest readVarParameter = (S7ParameterReadVarRequest) request.getParameter();

        List<S7MessageRequest> result = new LinkedList<>();

        // Calculate the maximum size an item can consume.
        int maxResponseSize = pduSize - EMPTY_READ_RESPONSE_SIZE;

        // This calculates the size of the header for the request and response.
        int curRequestSize = EMPTY_READ_REQUEST_SIZE;
        // An empty response has the same size as an empty request.
        int curResponseSize = EMPTY_READ_RESPONSE_SIZE;
        // List of all items in the current request.
        List<S7VarRequestParameterItem> curRequestItems = new LinkedList<>();

        for (S7VarRequestParameterItem readVarParameterItem : readVarParameter.getItems()) {
            final S7AddressAny address = (S7AddressAny)
                ((S7VarRequestParameterItemAddress) readVarParameterItem).getAddress();
            // Calculate the sizes in the request and response adding this item to the current request would add.
            int readRequestItemSize = readVarParameterItem.getLengthInBytes();
            // Constant size of the parameter item in the response (0 bytes) + Constant size of the payload item +
            // payload data size.
            int readResponseItemSize = 4 + (address.getNumberOfElements() * address.getTransportSize().getSizeInBytes());
            // If it's an odd number of bytes, add one to make it even
            if(readResponseItemSize % 2 == 1) {
                readResponseItemSize++;
            }

            // If the item would not fit into a separate message, we have to split it.
            if (((curRequestSize + readRequestItemSize) > pduSize) || (curResponseSize + readResponseItemSize > pduSize)) {
                // Create a new sub message.
                S7MessageRequest subMessage = new S7MessageRequest((short) tpduRefGen.getAndIncrement(),
                    new S7ParameterReadVarRequest(
                        curRequestItems),
                    null);
                result.add(subMessage);

                // Reset the counters.
                curRequestSize = EMPTY_READ_REQUEST_SIZE;
                curResponseSize = EMPTY_READ_RESPONSE_SIZE;
                curRequestItems = new LinkedList<>();

                S7VarRequestParameterItemAddress addressItem = (S7VarRequestParameterItemAddress) readVarParameterItem;
                if (addressItem.getAddress() instanceof S7AddressAny) {
                    S7AddressAny anyAddress = (S7AddressAny) addressItem.getAddress();

                    // Calculate the maximum number of items that would fit in a single request.
                    int maxNumElements = (int) Math.floor(
                        (double) maxResponseSize / (double) anyAddress.getTransportSize().getSizeInBytes());
                    int sizeMaxNumElementInBytes = maxNumElements * anyAddress.getTransportSize().getSizeInBytes();

                    // Initialize the loop with the total number of elements and the original address.
                    int remainingNumElements = anyAddress.getNumberOfElements();
                    int curByteAddress = anyAddress.getByteAddress();

                    // Keep on adding chunks of the original address until all have been added.
                    while (remainingNumElements > 0) {
                        int numCurElements = Math.min(remainingNumElements, maxNumElements);
                        S7VarRequestParameterItemAddress subVarParameterItem = new S7VarRequestParameterItemAddress(
                            new S7AddressAny(anyAddress.getTransportSize(), numCurElements,
                                anyAddress.getDbNumber(), anyAddress.getArea(), curByteAddress, (byte) 0));

                        // Create a new sub message.
                        subMessage = new S7MessageRequest((short) tpduRefGen.getAndIncrement(),
                            new S7ParameterReadVarRequest(Collections.singletonList(subVarParameterItem)),
                            null);
                        result.add(subMessage);

                        remainingNumElements -= maxNumElements;
                        curByteAddress += sizeMaxNumElementInBytes;
                    }
                }
            }

            // If adding the item would not exceed the sizes, add it to the current request.
            else {
                // Increase the current request sizes.
                curRequestSize += readRequestItemSize;
                curResponseSize += readResponseItemSize;
                // Add the item.
                curRequestItems.add(readVarParameterItem);
            }
        }

        // Add the remaining items to a final sub-request.
        if(!curRequestItems.isEmpty()) {
            // Create a new sub message.
            S7MessageRequest subMessage = new S7MessageRequest((short) tpduRefGen.getAndIncrement(),
                new S7ParameterReadVarRequest(
                    curRequestItems),
                null);
            result.add(subMessage);
        }

        return result;
    }

    private Collection<S7MessageRequest> processWriteVarParameter(S7MessageRequest request, int pduSize)
            throws PlcProtocolException {
        // TODO: Really find out the constraints ... do S7 devices all just accept single element write requests?
        return Collections.singletonList(request);
    }

    @Override
    public S7MessageResponseData processResponse(S7MessageRequest originalRequest,
                                             Map<S7MessageRequest, Either<S7MessageResponseData, Throwable>> result) {
        /*MessageType messageType = null;
        short tpduReference = requestMessage.getTpduReference();
        List<S7Parameter> s7Parameters = new LinkedList<>();
        List<S7Payload> s7Payloads = new LinkedList<>();

        Optional<VarParameter> varParameterOptional = requestMessage.getParameter(VarParameter.class);

        // This is neither a read request nor a write request, just merge all parameters together.
        if(!varParameterOptional.isPresent()) {
            for (S7ResponseMessage response : responses) {
                messageType = response.getMessageType();
                s7Parameters.addAll(response.getParameters());
                s7Payloads.addAll(response.getPayloads());
            }
        }

        // This is a read or write request, we have to merge all the items in the var parameter.
        else {
            List<VarParameterItem> parameterItems = new LinkedList<>();
            List<VarPayloadItem> payloadItems = new LinkedList<>();
            for (S7ResponseMessage response : responses) {
                messageType = response.getMessageType();
                parameterItems.addAll(response.getParameter(VarParameter.class)
                    .orElseThrow(() -> new PlcRuntimeException(
                        "Every response of a Read message should have a VarParameter instance")).getItems());
                Optional<VarPayload> payload = response.getPayload(VarPayload.class);
                payload.ifPresent(varPayload -> payloadItems.addAll(varPayload.getItems()));
            }

            List<VarParameterItem> mergedParameterItems = new LinkedList<>();
            List<VarPayloadItem> mergedPayloadItems = new LinkedList<>();
            VarParameter varParameter = varParameterOptional.get();

            int responseOffset = 0;
            for(int i = 0; i < varParameter.getItems().size(); i++) {
                S7AnyVarParameterItem requestItem = (S7AnyVarParameterItem) varParameter.getItems().get(i);

                // Get the pairs of corresponding parameter and payload items.
                S7AnyVarParameterItem responseParameterItem = (S7AnyVarParameterItem) parameterItems.get(0);
                VarPayloadItem responsePayloadItem = payloadItems.get(i + responseOffset);

                if(responsePayloadItem.getReturnCode() == DataTransportErrorCode.OK) {
                    int dataOffset = (responsePayloadItem.getData() != null) ? responsePayloadItem.getData().length : 0;

                    // The resulting parameter items is identical to the request parameter item.
                    mergedParameterItems.add(requestItem);

                    // The payload will have to be merged and the return codes will have to be examined.
                    if (requestItem.getNumElements() != responseParameterItem.getNumElements()) {
                        int itemSizeInBytes = requestItem.getDataType().getSizeInBytes();
                        int totalSizeInBytes = requestItem.getNumElements() * itemSizeInBytes;

                        if (varParameter.getType() == ParameterType.READ_VAR) {
                            byte[] data = new byte[totalSizeInBytes];
                            System.arraycopy(responsePayloadItem.getData(), 0, data, 0, responsePayloadItem.getData().length);

                            // Now iterate over the succeeding pairs of parameters and payloads till we have
                            // found the original number of elements.
                            while (dataOffset < totalSizeInBytes) {
                                responseOffset++;

                                // Get the next payload item in the list.
                                responsePayloadItem = payloadItems.get(i + responseOffset);

                                // Copy the data of this item behind the previous content.
                                if (varParameter.getType() == ParameterType.READ_VAR) {
                                    System.arraycopy(responsePayloadItem.getData(), 0, data, dataOffset, responsePayloadItem.getData().length);
                                    dataOffset += responsePayloadItem.getData().length;
                                }
                            }

                            mergedPayloadItems.add(new VarPayloadItem(DataTransportErrorCode.OK,
                                responsePayloadItem.getDataTransportSize(), data));
                        }
                    } else {
                        mergedPayloadItems.add(responsePayloadItem);
                    }
                } else {
                    mergedPayloadItems.add(responsePayloadItem);
                }
            }

            s7Parameters.add(new VarParameter(varParameter.getType(), mergedParameterItems));
            s7Payloads.add(new VarPayload(varParameter.getType(), mergedPayloadItems));
        }
        // TODO: The error codes are wrong
        return new S7ResponseMessage(messageType, tpduReference, s7Parameters, s7Payloads, (byte) 0xFF, (byte) 0xFF);*/
        return null;
    }

}
