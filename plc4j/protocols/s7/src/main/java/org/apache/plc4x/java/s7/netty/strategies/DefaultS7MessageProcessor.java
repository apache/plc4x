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
package org.apache.plc4x.java.s7.netty.strategies;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.base.messages.PlcProtocolMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7ResponseMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.apache.plc4x.java.s7.netty.util.S7RequestSizeCalculator;
import org.apache.plc4x.java.s7.netty.util.S7ResponseSizeEstimator;

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

    private AtomicInteger tpduRefGen;

    public DefaultS7MessageProcessor() {
        this.tpduRefGen = new AtomicInteger(1);
    }

    @Override
    public Collection<S7RequestMessage> processRequest(S7RequestMessage request, int pduSize) throws PlcException {
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

        Optional<VarParameter> varParameterOptional = request.getParameter(VarParameter.class);
        if (varParameterOptional.isPresent()) {
            VarParameter varParameter = varParameterOptional.get();

            // If this is a read operation, try to get as many items in as possible.
            if(varParameter.getType() == ParameterType.READ_VAR) {
                return processReadVarParameter(request, varParameter, pduSize).getRequestMessages();
            }

            // If this is a write operation, split up every array item into single value items
            // and every item into a separate message.
            else if(varParameter.getType() == ParameterType.WRITE_VAR) {
                return processWriteVarParameter(request, varParameter).getRequestMessages();
            }
        }

        return Collections.singletonList(request);
    }

    private S7CompositeRequestMessage processReadVarParameter(S7RequestMessage request, VarParameter varParameter, int pduSize) {
        // Create a new composite request message.
        S7CompositeRequestMessage compositeRequestMessage = new S7CompositeRequestMessage(request);

        // Create a var parameter without any items (yet).
        VarParameter subVarParameter = new VarParameter(varParameter.getType(), new LinkedList<>());

        // Create a sub message with only this empty parameter.
        S7RequestMessage subMessage = new S7RequestMessage(
            request.getMessageType(), (short) tpduRefGen.getAndIncrement(),
            Collections.singletonList(subVarParameter), Collections.emptyList(), compositeRequestMessage);

        // Add this sub-message to the composite.
        compositeRequestMessage.addRequestMessage(subMessage);

        // This calculates the size of the header for the request and response.
        int initialRequestSize = S7RequestSizeCalculator.getRequestMessageSize(subMessage);
        int curRequestSize = initialRequestSize;
        int initialResponseSize = S7ResponseSizeEstimator.getEstimatedResponseMessageSize(subMessage);
        int curResponseSize = initialResponseSize;

        VarParameter preProcessedVarParameter = new VarParameter(varParameter.getType(), new LinkedList<>());
        for (VarParameterItem varParameterItem : varParameter.getItems()) {
            // Use the S7RequestSizeCalculator to calculate the actual and estimated item sizes.
            int itemRequestSize = S7RequestSizeCalculator.getRequestItemTotalSize(
                varParameterItem, null);
            int itemResponseSize = S7ResponseSizeEstimator.getEstimatedResponseReadItemTotalSize(
                varParameterItem, null);

            // If the item would not fit into a separate message, we have to split it.
            if((initialRequestSize + itemRequestSize > pduSize) || (initialResponseSize + itemResponseSize > pduSize)) {
                // The max response size is the size of the empty response, plus the type and num-items (each one byte) of one VarParameter, plus the size of the header one VarPayloadItem
                int maxResponseSize = pduSize - (initialResponseSize + 2 + 4);

                S7AnyVarParameterItem s7AnyVarParameterItem = ((S7AnyVarParameterItem) varParameterItem);
                int maxNumElements = (int) Math.floor(
                    (double) maxResponseSize / (double) s7AnyVarParameterItem.getDataType().getSizeInBytes());
                int sizeMaxNumElementInBytes = maxNumElements * s7AnyVarParameterItem.getDataType().getSizeInBytes();
                int remainingNumElements = s7AnyVarParameterItem.getNumElements();
                int curByteOffset = s7AnyVarParameterItem.getByteOffset();

                while(remainingNumElements > 0) {
                    int numCurElements = Math.min(remainingNumElements, maxNumElements);
                    VarParameterItem subVarParameterItem = new S7AnyVarParameterItem(
                        s7AnyVarParameterItem.getSpecificationType(), s7AnyVarParameterItem.getMemoryArea(),
                        s7AnyVarParameterItem.getDataType(), numCurElements, s7AnyVarParameterItem.getDataBlockNumber(),
                        curByteOffset, (byte) 0);
                    preProcessedVarParameter.getItems().add(subVarParameterItem);

                    remainingNumElements -= maxNumElements;
                    curByteOffset += sizeMaxNumElementInBytes;
                }
            }
            // In all other cases, just forward the item.
            else {
                preProcessedVarParameter.getItems().add(varParameterItem);
            }
        }

        // For each var item of the original request, try adding them to the current sub-message
        // as long as it or the resulting response does not exceed the max PDU size.
        for (VarParameterItem varParameterItem : preProcessedVarParameter.getItems()) {
            // Use the S7RequestSizeCalculator to calculate the actual and estimated item sizes.
            int itemRequestSize = S7RequestSizeCalculator.getRequestItemTotalSize(
                varParameterItem, null);
            int itemResponseSize = S7ResponseSizeEstimator.getEstimatedResponseReadItemTotalSize(
                varParameterItem, null);

            // If adding this item, would exceed either the request or response size,
            // create a new sub-message and add this item to that.
            if ((curRequestSize + itemRequestSize > pduSize) || (curResponseSize + itemResponseSize > pduSize)) {
                // Create a new var parameter without any items (yet).
                subVarParameter = new VarParameter(varParameter.getType(), new LinkedList<>());

                // Create a new sub message with only this empty parameter.
                subMessage = new S7RequestMessage(
                    request.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                    Collections.singletonList(subVarParameter),
                    Collections.emptyList(), compositeRequestMessage);

                // Reset the message size
                curRequestSize = S7RequestSizeCalculator.getRequestMessageSize(subMessage) + itemRequestSize;
                curResponseSize = S7ResponseSizeEstimator.getEstimatedResponseMessageSize(subMessage) + itemResponseSize;

                // Add this new sub-message to the composite.
                compositeRequestMessage.addRequestMessage(subMessage);
            } else {
                // Increment the current sizes.
                curRequestSize += itemRequestSize;
                curResponseSize += itemResponseSize;
            }

            subVarParameter.getItems().add(varParameterItem);
        }
        return compositeRequestMessage;
    }

    private S7CompositeRequestMessage processWriteVarParameter(S7RequestMessage request, VarParameter varParameter)
            throws PlcProtocolException {
        // Create a new composite request message.
        S7CompositeRequestMessage compositeRequestMessage = new S7CompositeRequestMessage(request);

        VarPayload varPayload = request.getPayload(VarPayload.class)
            .orElseThrow(() -> new PlcProtocolException("Expecting payloads for a write request"));
        if(varParameter.getItems().size() != varPayload.getItems().size()) {
            throw new PlcProtocolException("Number of items in parameter and payload don't match");
        }
        List<VarParameterItem> parameterItems = varParameter.getItems();
        List<VarPayloadItem> payloadItems = varPayload.getItems();

        for (int i = 0; i < parameterItems.size(); i++) {
            VarParameterItem varParameterItem = parameterItems.get(i);
            VarPayloadItem varPayloadItem = payloadItems.get(i);

            if (varParameterItem instanceof S7AnyVarParameterItem) {
                S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItem;
                int byteOffset = s7AnyVarParameterItem.getByteOffset();
                if (s7AnyVarParameterItem.getDataType() == TransportSize.BOOL) {
                    processBooleanWriteVarParameter(request, varParameter, varPayload, s7AnyVarParameterItem,
                        varPayloadItem, byteOffset, compositeRequestMessage);
                } else {
                    processNonBooleanWriteVarParameter(request, varParameter, varPayload, s7AnyVarParameterItem,
                        varPayloadItem, byteOffset, compositeRequestMessage);
                }
            } else {
                throw new NotImplementedException("Handling of other element types not implemented.");
            }
        }
        return compositeRequestMessage;
    }

    private void processBooleanWriteVarParameter(S7RequestMessage request, VarParameter varParameter, VarPayload varPayload,
                                         S7AnyVarParameterItem s7AnyVarParameterItem, VarPayloadItem varPayloadItem,
                                         int byteOffset, S7CompositeRequestMessage compositeRequestMessage) {
        int curParameterByteOffset = byteOffset;
        byte curParameterBitOffset = s7AnyVarParameterItem.getBitOffset();
        byte curPayloadBitOffset = 0;
        for (int i = 0; i < s7AnyVarParameterItem.getNumElements(); i++) {
            // Create a new message with only one single value item in the var parameter.
            VarParameterItem item = new S7AnyVarParameterItem(
                s7AnyVarParameterItem.getSpecificationType(),
                s7AnyVarParameterItem.getMemoryArea(),
                s7AnyVarParameterItem.getDataType(), (short) 1,
                s7AnyVarParameterItem.getDataBlockNumber(),
                curParameterByteOffset, curParameterBitOffset);
            S7Parameter subVarParameter = new VarParameter(varParameter.getType(),
                Collections.singletonList(item));

            // Create a one-byte byte array and set it to 0x01, if the corresponding bit
            // was 1 else set it to 0x00.
            byte[] elementData = new byte[1];
            short curPayloadByteOffset = (short) ((short) i / 8);
            elementData[0] = (byte) ((varPayloadItem.getData()[curPayloadByteOffset] >> curPayloadBitOffset) & 0x01);

            // Create the new payload item.
            VarPayloadItem itemPayload = new VarPayloadItem(
                varPayloadItem.getReturnCode(),
                varPayloadItem.getDataTransportSize(), elementData);
            S7Payload subVarPayload = new VarPayload(varPayload.getType(),
                Collections.singletonList(itemPayload));

            // Create a new sub message.
            S7RequestMessage subMessage = new S7RequestMessage(
                request.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                Collections.singletonList(subVarParameter),
                Collections.singletonList(subVarPayload),
                compositeRequestMessage);

            // Add the new message to the composite.
            compositeRequestMessage.addRequestMessage(subMessage);

            // Calculate the new memory and bit offset.
            curParameterBitOffset++;
            if ((i > 0) && ((curParameterBitOffset % 8) == 0)) {
                curParameterByteOffset++;
                curParameterBitOffset = 0;
            }
            curPayloadBitOffset++;
            if ((i > 0) && ((curPayloadBitOffset % 8) == 0)) {
                curPayloadBitOffset = 0;
            }
        }
    }

    private void processNonBooleanWriteVarParameter(S7RequestMessage request, VarParameter varParameter, VarPayload varPayload,
                                            S7AnyVarParameterItem s7AnyVarParameterItem, VarPayloadItem varPayloadItem,
                                            int byteOffset, S7CompositeRequestMessage compositeRequestMessage) {
        int curByteOffset = byteOffset;
        int payloadPosition = 0;
        for (int i = 0; i < s7AnyVarParameterItem.getNumElements(); i++) {
            int elementSize = s7AnyVarParameterItem.getDataType().getSizeInBytes();

            // Create a new message with only one single value item in the var parameter.
            VarParameterItem itemParameter = new S7AnyVarParameterItem(
                s7AnyVarParameterItem.getSpecificationType(),
                s7AnyVarParameterItem.getMemoryArea(),
                s7AnyVarParameterItem.getDataType(), (short) 1,
                s7AnyVarParameterItem.getDataBlockNumber(),
                curByteOffset, (byte) 0);
            S7Parameter subVarParameter = new VarParameter(varParameter.getType(),
                Collections.singletonList(itemParameter));

            // Split up the big array into a separate byte-array that contains a single element.
            byte[] elementData = new byte[elementSize];
            System.arraycopy(varPayloadItem.getData(), payloadPosition, elementData, 0, elementSize);
            payloadPosition += elementSize;

            // Create the new payload item.
            VarPayloadItem itemPayload = new VarPayloadItem(
                varPayloadItem.getReturnCode(),
                varPayloadItem.getDataTransportSize(), elementData);
            S7Payload subVarPayload = new VarPayload(varPayload.getType(),
                Collections.singletonList(itemPayload));

            // Create a new sub message.
            S7RequestMessage subMessage = new S7RequestMessage(
                request.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                Collections.singletonList(subVarParameter),
                Collections.singletonList(subVarPayload),
                compositeRequestMessage);

            // Add the new message to the composite.
            compositeRequestMessage.addRequestMessage(subMessage);

            // Calculate the new memory offset.
            curByteOffset += elementSize;
        }
    }

    @Override
    public S7ResponseMessage processResponse(S7RequestMessage request, S7ResponseMessage response) {
        // If it's a split-up message, check if all parts are now acknowledged.
        if (request.getParent() instanceof S7CompositeRequestMessage) {
            S7CompositeRequestMessage parent = (S7CompositeRequestMessage) request.getParent();

            // Add the response to the container so we can add it's information to the composite response.
            parent.addResponseMessage(response);

            // If all parts of this split-up message are now acknowledged, create a unified
            // response object and pass that up to the higher layers.
            if (parent.isAcknowledged()) {
                return getMergedResponseMessage(parent.originalRequest, parent.getResponseMessages());
            } else {
                return null;
            }
        }
        // If it's an ordinary request, simply forward it  to the higher layers.
        else {
            return response;
        }
    }

    private S7ResponseMessage getMergedResponseMessage(S7RequestMessage requestMessage,
                                                       Collection<? extends S7ResponseMessage> responses) {
        MessageType messageType = null;
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
        return new S7ResponseMessage(messageType, tpduReference, s7Parameters, s7Payloads, (byte) 0xFF, (byte) 0xFF);
    }

    static class S7CompositeRequestMessage implements PlcProtocolMessage {

        private S7RequestMessage originalRequest;
        private List<S7RequestMessage> requestMessages;
        private List<S7ResponseMessage> responseMessages;

        S7CompositeRequestMessage(S7RequestMessage originalRequest) {
            this.originalRequest = originalRequest;
            this.requestMessages = new LinkedList<>();
            this.responseMessages = new LinkedList<>();
        }

        @Override
        public PlcProtocolMessage getParent() {
            return originalRequest;
        }

        /**
         * A {@link S7CompositeRequestMessage} is only acknowledged, if all children are acknowledged.
         *
         * @return true if all children are acknowledged.
         */
        private boolean isAcknowledged() {
            for (S7RequestMessage requestMessage : requestMessages) {
                if(!requestMessage.isAcknowledged()) {
                    return false;
                }
            }
            return true;
        }

        void addRequestMessage(S7RequestMessage requestMessage) {
            requestMessages.add(requestMessage);
        }

        public List<S7RequestMessage> getRequestMessages() {
            return requestMessages;
        }

        private void addResponseMessage(S7ResponseMessage responseMessage) {
            responseMessages.add(responseMessage);
        }

        public List<S7ResponseMessage> getResponseMessages() {
            return responseMessages;
        }
    }

}
