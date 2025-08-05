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
package org.apache.plc4x.java.s7light.readwrite.optimizer;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.tag.*;
import org.apache.plc4x.java.s7light.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.utils.*;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.apache.plc4x.java.spi.values.PlcValueAdapter;

import java.util.*;
import java.util.stream.Collectors;

public class S7Optimizer extends BaseOptimizer {

    public static final int EMPTY_READ_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterReadVarRequest(
        Collections.emptyList()), null).getLengthInBytes();
    public static final int EMPTY_READ_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterReadVarResponse(
        (short) 0), new S7PayloadReadVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();
    public static final int EMPTY_WRITE_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterWriteVarRequest(
        Collections.emptyList()), new S7PayloadWriteVarRequest(Collections.emptyList())).getLengthInBytes();
    public static final int EMPTY_WRITE_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterWriteVarResponse(
        (short) 0), new S7PayloadWriteVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();
    public static final int S7_ADDRESS_ANY_SIZE = 2 +
        new S7AddressAny(TransportSize.INT, 1, 1, MemoryArea.DATA_BLOCKS, 1, (byte) 0).getLengthInBytes();

    @Override
    protected int getNanosDelay() {
        return 500000;
    }

    @Override
    protected List<PlcReadRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        S7DriverContext s7DriverContext = (S7DriverContext) driverContext;
        List<PlcReadRequest> processedRequests = new LinkedList<>();

        // This calculates the size of the header for the request and response.
        int curRequestSize = EMPTY_READ_REQUEST_SIZE;
        // An empty response has the same size as an empty request.
        int curResponseSize = EMPTY_READ_RESPONSE_SIZE;

        // List of all items in the current request.

        LinkedHashMap<String, PlcTagItem<PlcTag>> curTagItems = new LinkedHashMap<>();

        for (String tagName : readRequest.getTagNames()) {

            //TODO: Individual processing of these types of tags. like S7StringTag
            if ((readRequest.getTag(tagName) instanceof S7SzlTag) ||
                (readRequest.getTag(tagName) instanceof S7ClkTag)) {
                // We are only expecting valid tagValueItems being passed in.
                curTagItems.put(tagName, new DefaultPlcTagItem<>(readRequest.getTag(tagName)));
                continue;
            }

            ////////////////////////////////////////////////////////////////////////////////////////////////////////////
            // Processing of normal tags.
            ////////////////////////////////////////////////////////////////////////////////////////////////////////////
            S7Tag tag = (S7Tag) readRequest.getTag(tagName);

            int readRequestItemSize = S7_ADDRESS_ANY_SIZE;
            // If we're reading var-length strings, then we'll read the sizes of the strings instead
            // and the S7ProtocolLogic will handle reading the actual strings in an additional request.
            int readResponseItemElementSize = tag.getDataType().getSizeInBytes();
            if(readRequest.getTag(tagName) instanceof S7StringVarLengthTag) {
                if(((S7StringVarLengthTag) readRequest.getTag(tagName)).getDataType() == TransportSize.STRING) {
                    readResponseItemElementSize = 2;
                } else if(((S7StringVarLengthTag) readRequest.getTag(tagName)).getDataType() == TransportSize.WSTRING) {
                    readResponseItemElementSize = 4;
                }
            }
            int readResponseItemSize = 4 + (tag.getNumberOfElements() * readResponseItemElementSize);
            // If it's an odd number of bytes, add one to make it even
            if (readResponseItemSize % 2 == 1) {
                readResponseItemSize++;
            }

            // If adding the item would not exceed the sizes, add it to the current request.
            if (((curRequestSize + readRequestItemSize) <= s7DriverContext.getPduSize()) &&
                ((curResponseSize + readResponseItemSize) <= s7DriverContext.getPduSize())) {
                // Increase the current request sizes.
                curRequestSize += readRequestItemSize;
                curResponseSize += readResponseItemSize;

                // Add the tag to the current request
                // We are only expecting valid tagValueItems being passed in.
                curTagItems.put(tagName, new DefaultPlcTagItem<>(tag));
            }
            // If the current item would exceed the PDU size in the response, even if this is a new request, split
            // up the array to fill up the last request and create so many sub-requests that read all in total.
            else if (EMPTY_READ_RESPONSE_SIZE + readResponseItemSize > s7DriverContext.getPduSize()){
                int maxPayloadSize = s7DriverContext.getPduSize() - (EMPTY_READ_RESPONSE_SIZE + 4);
                int numRequests = (int) Math.ceil(((double) tag.getNumberOfElements() * (double) tag.getDataType().getSizeInBytes()) / (double) maxPayloadSize);
                int curByteOffset = tag.getByteOffset();
                int numItemsPerRequest = (int) Math.ceil((double) tag.getNumberOfElements() / (double) numRequests);
                int itemsLeft = tag.getNumberOfElements();
                for(int curRequest = 0; curRequest < numRequests; curRequest++) {
                    int numCurRequestItems = Math.min(numItemsPerRequest, itemsLeft);
                    S7Tag tagFragment = new S7Tag(tag.getDataType(), tag.getMemoryArea(), tag.getBlockNumber(), curByteOffset, (byte) 0, numCurRequestItems);
                    LinkedHashMap<String, PlcTagItem<PlcTag>> tagFragments = new LinkedHashMap<>();
                    tagFragments.put(tagName, new DefaultPlcTagItem<>(tagFragment));
                    processedRequests.add(new DefaultPlcReadRequest(((DefaultPlcReadRequest) readRequest).getReader(), tagFragments));
                    curByteOffset += numItemsPerRequest * tag.getDataType().getSizeInBytes();
                    itemsLeft -= numCurRequestItems;
                }
            }
            // If adding the current item would exceed the PDU size in the request or the response, start a new request.
            else {
                // Create a new PlcReadRequest containing the current tag item.
                processedRequests.add(new DefaultPlcReadRequest(
                    ((DefaultPlcReadRequest) readRequest).getReader(), curTagItems));

                // Reset the size and item lists.
                curRequestSize = EMPTY_READ_REQUEST_SIZE + readRequestItemSize;
                curResponseSize = EMPTY_READ_RESPONSE_SIZE + readResponseItemSize;
                curTagItems = new LinkedHashMap<>();

                // Splitting of huge tags not yet implemented, throw an exception instead.
                if (((curRequestSize + readRequestItemSize) > s7DriverContext.getPduSize()) &&
                    ((curResponseSize + readResponseItemSize) > s7DriverContext.getPduSize())) {
                    throw new PlcRuntimeException("Tag size exceeds maximum payload for one item.");
                }

                // Add the tag to the new current request
                // We are only expecting valid tagValueItems being passed in.
                curTagItems.put(tagName, new DefaultPlcTagItem<>(tag));
            }
        }

        // Create a new PlcReadRequest from the remaining tag items.
        if (!curTagItems.isEmpty()) {
            processedRequests.add(new DefaultPlcReadRequest(
                ((DefaultPlcReadRequest) readRequest).getReader(), curTagItems));
        }

        return processedRequests;
    }

    protected PlcReadResponse processReadResponses(PlcReadRequest readRequest, Map<PlcReadRequest, SubResponse<PlcReadResponse>> readResponses, DriverContext driverContext) {
        Map<String, PlcResponseItem<PlcValue>> tagValues = new HashMap<>();
        for (Map.Entry<PlcReadRequest, SubResponse<PlcReadResponse>> requestsEntries : readResponses.entrySet()) {
            PlcReadRequest curRequest = requestsEntries.getKey();
            SubResponse<PlcReadResponse> readResponse = requestsEntries.getValue();
            for (String tagName : curRequest.getTagNames()) {
                if (readResponse.isSuccess()) {
                    PlcReadResponse subReadResponse = readResponse.getResponse();
                    PlcResponseCode responseCode = subReadResponse.getResponseCode(tagName);
                    PlcValue value = subReadResponse.getAsPlcValue().getValue(tagName);

                    if (readRequest.getTag(tagName) instanceof S7SzlTag) {
                        tagValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, value));
                    } else {                                        
                        // If the number of elements in the current tag differs from the global number,
                        // then this is a split-up array and needs to be explicitly handled.
                        S7Tag globalS7Tag = (S7Tag) readRequest.getTag(tagName);
                        S7Tag currentS7Tag = (S7Tag) curRequest.getTag(tagName);
                        if(currentS7Tag.getNumberOfElements() != globalS7Tag.getNumberOfElements()) {
                            if(responseCode == PlcResponseCode.OK) {
                                // We have to merge byte-array data differently.
                                if (globalS7Tag.getDataType() == TransportSize.BYTE) {
                                    PlcRawByteArray existingItem;
                                    if (tagValues.containsKey(tagName)) {
                                        PlcResponseItem<PlcValue> existingTagItem = tagValues.get(tagName);
                                        // If a previous response was invalid, we'll discard this part of it too.
                                        if(existingTagItem.getResponseCode() != PlcResponseCode.OK) {
                                            continue;
                                        }
                                        existingItem = (PlcRawByteArray) existingTagItem.getValue();
                                    } else {
                                        existingItem = new PlcRawByteArray(new byte[globalS7Tag.getNumberOfElements()]);
                                        tagValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, existingItem));
                                    }
                                    byte[] currentItemByteArray = value.getRaw();
                                    int elementOffset = (currentS7Tag.getByteOffset() - globalS7Tag.getByteOffset()) / globalS7Tag.getDataType().getSizeInBytes();
                                    byte[] existingByteArray = existingItem.getRaw();
                                    System.arraycopy(currentItemByteArray, 0, existingByteArray, elementOffset, currentItemByteArray.length);
                                }
                                // Merge typical list-type responses.
                                else {
                                    List<PlcValue> existingItems;
                                    if (tagValues.containsKey(tagName)) {
                                        PlcResponseItem<PlcValue> existingTagItem = tagValues.get(tagName);
                                        // If a previous response was invalid, we'll discard this part of it too.
                                        if(existingTagItem.getResponseCode() != PlcResponseCode.OK) {
                                            continue;
                                        }
                                        existingItems = (List<PlcValue>) existingTagItem.getValue().getList();
                                    } else {
                                        existingItems = new ArrayList<>(Arrays.asList(new PlcValue[globalS7Tag.getNumberOfElements()]));
                                        PlcListModifiable mergedValue = new PlcListModifiable(existingItems);
                                        tagValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, mergedValue));
                                    }
                                    List<? extends PlcValue> currentItems = value.getList();
                                    int elementOffset = (currentS7Tag.getByteOffset() - globalS7Tag.getByteOffset()) / globalS7Tag.getDataType().getSizeInBytes();
                                    for (int i = 0; i < currentS7Tag.getNumberOfElements(); i++) {
                                        existingItems.set(i + elementOffset, currentItems.get(i));
                                    }
                                }
                            } else {
                                // Even if a previous part was ok, if at least one part is invalid, we'll
                                // treat all parts equally broken.
                                tagValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, new PlcNull()));
                            }
                        } else {
                            tagValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, value));
                        }
                    
                    }
                } else {
                    tagValues.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
                
            }
        }
        return new DefaultPlcReadResponse(readRequest, tagValues);
    }


    @Override
    protected List<PlcWriteRequest> processWriteRequest(PlcWriteRequest writeRequest, DriverContext driverContext) {
        for (String tagName : writeRequest.getTagNames()) {
            if (writeRequest.getTag(tagName) instanceof S7ClkTag) {
                return Collections.singletonList(writeRequest);
            }
        }

        S7DriverContext s7DriverContext = (S7DriverContext) driverContext;
        List<PlcWriteRequest> processedRequests = new LinkedList<>();

        // This calculates the size of the header for the request and response.
        int curRequestSize = EMPTY_WRITE_REQUEST_SIZE;
        // An empty response has the same size as an empty request.
        int curResponseSize = EMPTY_WRITE_RESPONSE_SIZE;

        // List of all items in the current request.
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> curTags = new LinkedHashMap<>();

        for (String tagName : writeRequest.getTagNames()) {
            S7Tag tag = (S7Tag) writeRequest.getTag(tagName);
            PlcValue value = writeRequest.getPlcValue(tagName);

            int writeRequestItemSize = S7_ADDRESS_ANY_SIZE + 4/* Size of Payload item header*/;
            if (tag.getDataType() == TransportSize.BOOL) {
                writeRequestItemSize += (int) Math.ceil((double) tag.getNumberOfElements() / 8);
            }
            // Handle fixed length strings differently.
            else if (tag instanceof S7StringFixedLengthTag) {
                S7StringFixedLengthTag stringFixedLengthTag = (S7StringFixedLengthTag) tag;
                if(tag.getDataType() == TransportSize.WSTRING) {
                    writeRequestItemSize += tag.getNumberOfElements() * (stringFixedLengthTag.getStringLength() + 2) * 2;
                } else {
                    writeRequestItemSize += tag.getNumberOfElements() * (stringFixedLengthTag.getStringLength() + 2);
                }
            }
            // With var-length strings, we need to get the length of the string first and use that.
            else if (tag instanceof S7StringVarLengthTag) {
                PlcValue plcValue = writeRequest.getPlcValue(tagName);
                int length = plcValue.getString().length();
                if(tag.getDataType() == TransportSize.WSTRING) {
                    writeRequestItemSize += tag.getNumberOfElements() * (length + 2) * 2;
                } else {
                    writeRequestItemSize += tag.getNumberOfElements() * (length + 2);
                }
            }
            else {
                writeRequestItemSize += (tag.getNumberOfElements() * tag.getDataType().getSizeInBytes());
            }
            // If it's an odd number of bytes, add one to make it even
            if (writeRequestItemSize % 2 == 1) {
                writeRequestItemSize++;
            }
            // The response for one item is just one byte containing the return code.
            int writeResponseItemSize = 1;

            // If adding the item would not exceed the sizes, add it to the current request.
            if (((curRequestSize + writeRequestItemSize) <= s7DriverContext.getPduSize()) &&
                ((curResponseSize + writeResponseItemSize) <= s7DriverContext.getPduSize())) {
                // Increase the current request sizes.
                curRequestSize += writeRequestItemSize;
                curResponseSize += writeResponseItemSize;

                // Add the item.
            }
            // If adding them would exceed, start a new request.
            else {
                // Create a new PlcWriteRequest containing the current tag item.
                processedRequests.add(new DefaultPlcWriteRequest(
                    ((DefaultPlcWriteRequest) writeRequest).getWriter(), curTags));

                // Reset the size and item lists.
                curRequestSize = EMPTY_WRITE_REQUEST_SIZE + writeRequestItemSize;
                curResponseSize = EMPTY_WRITE_RESPONSE_SIZE + writeResponseItemSize;
                curTags = new LinkedHashMap<>();

                // Splitting of huge tags not yet implemented, throw an exception instead.
                if (((curRequestSize + writeRequestItemSize) > s7DriverContext.getPduSize()) &&
                    ((curResponseSize + writeResponseItemSize) > s7DriverContext.getPduSize())) {
                    throw new PlcRuntimeException("Tag size exceeds maximum payload for one item.");
                }
            }
            // We are only expecting valid tagValueItems being passed in.
            curTags.put(tagName, new DefaultPlcTagValueItem<>(tag, value));
        }

        // Create a new PlcWriteRequest from the remaining tag items.
        if (!curTags.isEmpty()) {
            processedRequests.add(new DefaultPlcWriteRequest(
                ((DefaultPlcWriteRequest) writeRequest).getWriter(), curTags));
        }

        return processedRequests;
    }

    /**
     * Little helper that helps avoid the problem, that the elements of a PlcList
     * cannot be modified after being created.
     */
    private static class PlcListModifiable extends PlcValueAdapter {

        private final List<PlcValue> listItems;

        public PlcListModifiable(List<PlcValue> listItems) {
            // This line is the only difference to the original one.
            this.listItems = listItems;
        }

        @Override
        public PlcValueType getPlcValueType() {
            return PlcValueType.List;
        }

        public void add(PlcValue value) {
            listItems.add(value);
        }

        @Override
        public Object getObject() {
            return getList();
        }

        @Override
        public boolean isList() {
            return true;
        }

        @Override
        public int getLength() {
            return listItems.size();
        }

        @Override
        public PlcValue getIndex(int i) {
            return listItems.get(i);
        }

        @Override
        public List<PlcValue> getList() {
            return listItems;
        }

        @Override
        public String toString() {
            return "[" + listItems.stream().map(PlcValue::toString).collect(Collectors.joining(",")) + "]";
        }

        @Override
        public void serialize(WriteBuffer writeBuffer) throws SerializationException {
            writeBuffer.pushContext("PlcListModifiable");
            for (PlcValue listItem : listItems) {
                if (!(listItem instanceof Serializable)) {
                    throw new PlcRuntimeException("Error serializing. List item doesn't implement XmlSerializable");
                }
                ((Serializable) listItem).serialize(writeBuffer);
            }
            writeBuffer.popContext("PlcListModifiable");
        }
    }

}
