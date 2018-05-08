package org.apache.plc4x.java.s7.netty.strategies;/*
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

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.s7.netty.model.messages.S7CompositeRequestMessage;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
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
    public Collection<? extends S7Message> process(S7Message s7Message, int pduSize) {
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

        if(s7Message instanceof S7RequestMessage) {
            S7RequestMessage originalRequestMessage = (S7RequestMessage) s7Message;
            Optional<VarParameter> varParameterOptional = originalRequestMessage.getParameter(VarParameter.class);
            if (varParameterOptional.isPresent()) {
                VarParameter varParameter = varParameterOptional.get();

                // Create a new composite request message.
                S7CompositeRequestMessage compositeRequestMessage =
                    new S7CompositeRequestMessage(originalRequestMessage);

                // If this is a read operation, try to get as many items in as possible.
                if(varParameter.getType() == ParameterType.READ_VAR) {
                    // Create a var parameter without any items (yet).
                    S7Parameter subVarParameter = new VarParameter(varParameter.getType(), new LinkedList<>());

                    // Create a sub message with only this empty parameter.
                    S7RequestMessage subMessage = new S7RequestMessage(
                        s7Message.getMessageType(), (short) tpduRefGen.getAndIncrement(),
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
                        Optional<VarPayloadItem> payloadItem = s7Message.getPayload(VarPayloadItem.class);
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
                                s7Message.getMessageType(), (short) tpduRefGen.getAndIncrement(),
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
                        ((VarParameter) subVarParameter).getItems().add(varParameterItem);
                    }
                }

                // If this is a write operation, split up every array item into single value items
                // and every item into a separate message.
                else if(varParameter.getType() == ParameterType.WRITE_VAR) {
                    for (VarParameterItem varParameterItem : varParameter.getItems()) {
                        if(varParameterItem instanceof S7AnyVarParameterItem) {
                            S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameterItem;
                            short byteOffset = s7AnyVarParameterItem.getByteOffset();
                            if(s7AnyVarParameterItem.getTransportSize() == TransportSize.BIT) {
                                byte bitOffset = 0;
                                for (int i = 0; i < s7AnyVarParameterItem.getNumElements(); i++) {
                                    // Calculate the new memory and bit offset.
                                    if((i > 0) && (i % 8) == 0) {
                                        byteOffset++;
                                        bitOffset = 0;
                                    } else {
                                        bitOffset++;
                                    }

                                    // Create a new message with only one single value item in the var parameter.
                                    VarParameterItem item = new S7AnyVarParameterItem(
                                        s7AnyVarParameterItem.getSpecificationType(),
                                        s7AnyVarParameterItem.getMemoryArea(),
                                        s7AnyVarParameterItem.getTransportSize(), (short) 1,
                                        s7AnyVarParameterItem.getDataBlockNumber(),
                                        byteOffset, bitOffset);
                                    S7Parameter subVarParameter = new VarParameter(varParameter.getType(),
                                        Collections.singletonList(item));
                                    S7RequestMessage subMessage = new S7RequestMessage(
                                        s7Message.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                                        Collections.singletonList(subVarParameter), Collections.emptyList(),
                                        compositeRequestMessage);

                                    // Add the new message to the composite.
                                    compositeRequestMessage.addRequestMessage(subMessage);
                                }
                            } else {
                                for (int i = 0; i < s7AnyVarParameterItem.getNumElements(); i++) {
                                    byteOffset += s7AnyVarParameterItem.getTransportSize().getSizeInBytes();

                                    // Create a new message with only one single value item in the var parameter.
                                    VarParameterItem item = new S7AnyVarParameterItem(
                                        s7AnyVarParameterItem.getSpecificationType(),
                                        s7AnyVarParameterItem.getMemoryArea(),
                                        s7AnyVarParameterItem.getTransportSize(), (short) 1,
                                        s7AnyVarParameterItem.getDataBlockNumber(),
                                        byteOffset, (byte) 0);
                                    S7Parameter subVarParameter = new VarParameter(varParameter.getType(),
                                        Collections.singletonList(item));
                                    S7RequestMessage subMessage = new S7RequestMessage(
                                        s7Message.getMessageType(), (short) tpduRefGen.getAndIncrement(),
                                        Collections.singletonList(subVarParameter), Collections.emptyList(),
                                        compositeRequestMessage);

                                    // Add the new message to the composite.
                                    compositeRequestMessage.addRequestMessage(subMessage);
                                }
                            }
                        } else {
                            throw new NotImplementedException("Handling of other element types not implemented.");
                        }
                    }
                }
                return compositeRequestMessage.getRequestMessages();
            }
        }

        return Collections.singletonList(s7Message);
    }

}
