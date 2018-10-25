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
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.apache.plc4x.java.s7.netty.util.S7RequestSizeCalculator;
import org.apache.plc4x.java.s7.netty.util.S7ResponseSizeEstimator;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;

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
    public Collection<? extends S7RequestMessage> processRequest(S7RequestMessage request, int pduSize)
        throws PlcException {
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

            // Create a new composite request message.
            S7CompositeRequestMessage compositeRequestMessage =
                new S7CompositeRequestMessage(request);

            // If this is a read operation, try to get as many items in as possible.
            if(varParameter.getType() == ParameterType.READ_VAR) {
                // Create a var parameter without any items (yet).
                VarParameter subVarParameter = new VarParameter(varParameter.getType(), new LinkedList<>());

                // Create a sub message with only this empty parameter.
                S7RequestMessage subMessage = new S7RequestMessage(
                    request.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                    Collections.singletonList(subVarParameter), Collections.emptyList(), compositeRequestMessage);

                // Add this sub-message to the composite.
                compositeRequestMessage.addRequestMessage(subMessage);

                // This calculates the size of the header for the request and response.
                int curRequestSize = S7RequestSizeCalculator.getRequestMessageSize(subMessage);
                int curResponseSize = S7ResponseSizeEstimator.getEstimatedResponseMessageSize(subMessage);

                // For each var item of the original request, try adding them to the current sub-message
                // as long as it or the resulting response does not exceed the max PDU size.
                for (VarParameterItem varParameterItem : varParameter.getItems()) {
                    VarPayloadItem varPayloadItem = null;
                    Optional<VarPayloadItem> payloadItem = request.getPayload(VarPayloadItem.class);
                    if (payloadItem.isPresent()) {
                        varPayloadItem = payloadItem.get();
                    }

                    // Use the S7RequestSizeCalculator to calculate the actual and estimated item sizes.
                    int itemRequestSize = S7RequestSizeCalculator.getRequestItemTotalSize(
                        varParameterItem, varPayloadItem);
                    int itemResponseSize = S7ResponseSizeEstimator.getEstimatedResponseReadItemTotalSize(
                        varParameterItem, varPayloadItem);

                    // When adding this item to the request we would exceed the pdu size in
                    // the request or response, so we have to create a new sub-message.
                    if ((curRequestSize + itemRequestSize > pduSize) || (curResponseSize + itemResponseSize > pduSize)) {
                        // Create a new var parameter without any items (yet).
                        subVarParameter = new VarParameter(varParameter.getType(), new LinkedList<>());

                        // Create a new sub message with only this empty parameter.
                        subMessage = new S7RequestMessage(
                            request.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                            Collections.singletonList(subVarParameter),
                            Collections.emptyList(), compositeRequestMessage);

                        // Reset the message size
                        curRequestSize = S7RequestSizeCalculator.getRequestMessageSize(subMessage);
                        curResponseSize = S7ResponseSizeEstimator.getEstimatedResponseMessageSize(subMessage);

                        // Add this new sub-message to the composite.
                        compositeRequestMessage.addRequestMessage(subMessage);
                    } else {
                        // Increment the current sizes.
                        curRequestSize += itemRequestSize;
                        curResponseSize += itemResponseSize;
                    }

                    // Add the item to the current subVarParameter.
                    subVarParameter.getItems().add(varParameterItem);
                }
            }

            // If this is a write operation, split up every array item into single value items
            // and every item into a separate message.
            else if(varParameter.getType() == ParameterType.WRITE_VAR) {
                VarPayload varPayload = request.getPayload(VarPayload.class)
                    .orElseThrow(() -> new PlcProtocolException("Expecting payloads for a write request"));
                if(varParameter.getItems().size() != varPayload.getItems().size()) {
                    throw new PlcProtocolException("Number of items in parameter and payload don't match");
                }
                List<VarParameterItem> parameterItems = varParameter.getItems();
                List<VarPayloadItem> payloadItems = varPayload.getItems();

                for (int i1 = 0; i1 < parameterItems.size(); i1++) {
                    VarParameterItem varParameterItem = parameterItems.get(i1);
                    VarPayloadItem varPayloadItem = payloadItems.get(i1);
                    if (varParameterItem instanceof S7AnyVarParameterItem) {
                        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItem;
                        short byteOffset = s7AnyVarParameterItem.getByteOffset();
                        if (s7AnyVarParameterItem.getDataType() == TransportSize.BOOL) {
                            byte bitOffset = 0;
                            for (int i = 0; i < s7AnyVarParameterItem.getNumElements(); i++) {
                                // Create a new message with only one single value item in the var parameter.
                                VarParameterItem item = new S7AnyVarParameterItem(
                                    s7AnyVarParameterItem.getSpecificationType(),
                                    s7AnyVarParameterItem.getMemoryArea(),
                                    s7AnyVarParameterItem.getDataType(), (short) 1,
                                    s7AnyVarParameterItem.getDataBlockNumber(),
                                    byteOffset, bitOffset);
                                S7Parameter subVarParameter = new VarParameter(varParameter.getType(),
                                    Collections.singletonList(item));

                                // Create a one-byte byte array and set it to 0x01, if the corresponding bit
                                // was 1 else set it to 0x00.
                                byte[] elementData = new byte[1];
                                elementData[0] = (byte) ((varPayloadItem.getData()[byteOffset] >> bitOffset) & 0x01);

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
                                bitOffset++;
                                if ((i > 0) && ((bitOffset % 8) == 0)) {
                                    byteOffset++;
                                    bitOffset = 0;
                                }
                            }
                        } else {
                            int payloadPosition = 0;
                            for (int i = 0; i < s7AnyVarParameterItem.getNumElements(); i++) {
                                int elementSize = s7AnyVarParameterItem.getDataType().getSizeInBytes();

                                // Create a new message with only one single value item in the var parameter.
                                VarParameterItem itemParameter = new S7AnyVarParameterItem(
                                    s7AnyVarParameterItem.getSpecificationType(),
                                    s7AnyVarParameterItem.getMemoryArea(),
                                    s7AnyVarParameterItem.getDataType(), (short) 1,
                                    s7AnyVarParameterItem.getDataBlockNumber(),
                                    byteOffset, (byte) 0);
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
                                byteOffset += elementSize;
                            }
                        }
                    } else {
                        throw new NotImplementedException("Handling of other element types not implemented.");
                    }
                }
            }
            return compositeRequestMessage.getRequestMessages();
        }

        return Collections.singletonList(request);
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

        S7ResponseMessage firstResponse = null;
        short tpduReference = requestMessage.getTpduReference();
        List<S7Parameter> s7Parameters = new LinkedList<>();
        List<S7Payload> s7Payloads = new LinkedList<>();
        byte errorClass = 0;
        byte errorCode = 0;
        VarParameter readVarParameter = null;
        VarParameter writeVarParameter = null;
        VarPayload readVarPayload = null;
        VarPayload writeVarPayload = null;

        // TODO: We should change this code to not use the lists of the first parameter or payload as this can cause problems when using mutable lists.
        for (S7ResponseMessage response : responses) {
            if(firstResponse == null) {
                firstResponse = response;
            }
            // Some parameters have to be merged. In case of read and write parameters
            // their items have to be merged into one single parameter.
            for(S7Parameter parameter : response.getParameters()) {
                if (parameter.getType() == ParameterType.READ_VAR) {
                    if (readVarParameter == null) {
                        readVarParameter = (VarParameter) parameter;
                        s7Parameters.add(parameter);
                    } else {
                        readVarParameter.mergeParameter((VarParameter) parameter);
                    }
                } else if (parameter.getType() == ParameterType.WRITE_VAR) {
                    if (writeVarParameter == null) {
                        writeVarParameter = (VarParameter) parameter;
                        s7Parameters.add(parameter);
                    } else {
                        writeVarParameter.mergeParameter((VarParameter) parameter);
                    }
                } else {
                    s7Parameters.add(parameter);
                }
            }

            // Some payloads have to be merged. In case of read and write payloads
            // their items have to be merged into one single payload.
            for(S7Payload payload : response.getPayloads()) {
                if(payload.getType() == ParameterType.READ_VAR) {
                    if (readVarPayload == null) {
                        readVarPayload = (VarPayload) payload;
                        s7Payloads.add(payload);
                    } else {
                        readVarPayload.mergePayload((VarPayload) payload);
                    }
                } else if(payload.getType() == ParameterType.WRITE_VAR) {
                    if(writeVarPayload == null) {
                        writeVarPayload = (VarPayload) payload;
                        s7Payloads.add(payload);
                    } else {
                        writeVarPayload.mergePayload((VarPayload) payload);
                    }
                } else {
                    s7Payloads.add(payload);
                }
            }
        }
        if(firstResponse != null) {
            MessageType messageType = firstResponse.getMessageType();
            return new S7ResponseMessage(messageType, tpduReference, s7Parameters, s7Payloads, errorClass, errorCode);
        }
        return null;
    }

    static class S7CompositeRequestMessage implements PlcProtocolMessage {

        private S7RequestMessage originalRequest;
        private Collection<S7RequestMessage> requestMessages;
        private Collection<S7ResponseMessage> responseMessages;

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

        private Collection<S7RequestMessage> getRequestMessages() {
            return requestMessages;
        }

        private void addResponseMessage(S7ResponseMessage responseMessage) {
            responseMessages.add(responseMessage);
        }

        private Collection<S7ResponseMessage> getResponseMessages() {
            return responseMessages;
        }
    }

}
