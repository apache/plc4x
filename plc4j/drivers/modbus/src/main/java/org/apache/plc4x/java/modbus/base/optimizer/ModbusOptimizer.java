/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.modbus.base.optimizer;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.context.ModbusContext;
import org.apache.plc4x.java.modbus.base.protocol.ModbusProtocolLogic;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagCoil;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagDiscreteInput;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagExtendedRegister;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagInputRegister;
import org.apache.plc4x.java.modbus.readwrite.DataItem;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferXmlBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.TreeSet;

/**
 * In order to read more data more efficiently, this optimizer for modbus joins together individual items
 * and reads larger arrays of data.
 */
public class ModbusOptimizer extends SingleTagOptimizer {
    private final Logger logger = LoggerFactory.getLogger(ModbusOptimizer.class);

    /**
     * Per default the number of registers that can be read are 125 registers.
     * The number of coils that can be read in one request are 2000 coils.
     * This optimizer will sort all request items, split them up into coils and registers.
     * Then it will try to group them together to produce the least amount of traffic
     * on the wire by requesting larger chunks of words and serving multiple items
     * with one request. The matching processReadResponses will then handle the splitting
     * of the results.
     *
     * @param readRequest   the original read request
     * @param driverContext the driver context
     * @return a list of rewritten sub-requests
     */
    @Override
    protected List<PlcReadRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        ModbusContext modbusContext = (ModbusContext) driverContext;

        // Sort the different types of tags as all need to be requested separately anyway.
        TreeSet<ModbusTag> coils = null;
        TreeSet<ModbusTag> holdingRegisters = null;
        TreeSet<ModbusTag> inputRegisters = null;
        TreeSet<ModbusTag> extendedRegisters = null;
        TreeSet<ModbusTag> discreteInputs = null;
        for (PlcTag tag : readRequest.getTags()) {
            if (tag instanceof ModbusTagCoil) {
                if (coils == null) {
                    coils = new TreeSet<>(Comparator.comparingInt(ModbusTag::getAddress));
                }
                coils.add(((ModbusTagCoil) tag));
            } else if (tag instanceof ModbusTagHoldingRegister) {
                if (holdingRegisters == null) {
                    holdingRegisters = new TreeSet<>(Comparator.comparingInt(ModbusTag::getAddress));
                }
                holdingRegisters.add(((ModbusTagHoldingRegister) tag));
            } else if (tag instanceof ModbusTagInputRegister) {
                if (inputRegisters == null) {
                    inputRegisters = new TreeSet<>(Comparator.comparingInt(ModbusTag::getAddress));
                }
                inputRegisters.add(((ModbusTagInputRegister) tag));
            } else if (tag instanceof ModbusTagExtendedRegister) {
                if (extendedRegisters == null) {
                    extendedRegisters = new TreeSet<>(Comparator.comparingInt(ModbusTag::getAddress));
                }
                extendedRegisters.add(((ModbusTagExtendedRegister) tag));
            } else if (tag instanceof ModbusTagDiscreteInput) {
                if (discreteInputs == null) {
                    discreteInputs = new TreeSet<>(Comparator.comparingInt(ModbusTag::getAddress));
                }
                discreteInputs.add(((ModbusTagDiscreteInput) tag));
            }
        }

        // Add sub requests for every type of tag in this request.
        PlcReader reader = ((DefaultPlcReadRequest) readRequest).getReader();
        List<PlcReadRequest> subRequests = new ArrayList<>();
        if (coils != null) {
            subRequests.addAll(processCoilRequests(coils, reader, modbusContext));
        }
        if (holdingRegisters != null) {
            subRequests.addAll(processRegisterRequests(holdingRegisters, reader, (address, count, dataType) -> new ModbusTagHoldingRegister(address, count, dataType, Collections.emptyMap()), modbusContext));
        }
        if (inputRegisters != null) {
            subRequests.addAll(processRegisterRequests(inputRegisters, reader, (address, count, dataType) -> new ModbusTagInputRegister(address, count, dataType, Collections.emptyMap()), modbusContext));
        }
        if (extendedRegisters != null) {
            subRequests.addAll(processRegisterRequests(extendedRegisters, reader, (address, count, dataType) -> new ModbusTagExtendedRegister(address, count, dataType, Collections.emptyMap()), modbusContext));
        }
        if (discreteInputs != null) {
            subRequests.addAll(processRegisterRequests(discreteInputs, reader, (address, count, dataType) -> new ModbusTagDiscreteInput(address, count, dataType, Collections.emptyMap()), modbusContext));
        }
        return subRequests;
    }

    /**
     * When reading large chunks of data, here we need to read the parts that were originally requested.
     *
     * @param readRequest original read request
     * @param readResponses map of the sub-requests that were executed
     * @return post processed response reflecting the unmodified user request items
     */
    @Override
    protected PlcReadResponse processReadResponses(PlcReadRequest readRequest, Map<PlcReadRequest, SubResponse<PlcReadResponse>> readResponses, DriverContext driverContext) {
        ModbusContext modbusContext = (ModbusContext) driverContext;
        try {
            // Build an index of all the data returned by all requests.
            // This data should contain all the bits needed to create the response of the original request.
            Map<String, List<Response>> responses = new HashMap<>();
            for (PlcReadRequest optimizedReadRequest : readResponses.keySet()) {
                PlcReadResponse optimizedReadResponse = readResponses.get(optimizedReadRequest).getResponse();
                if (optimizedReadResponse == null) {
                    continue;
                }
                // Optimized read requests only contain one ModbusTag.
                String tagName = optimizedReadRequest.getTagNames().stream().findFirst().orElse(null);
                if (tagName == null) {
                    continue;
                }
                ModbusTag modbusTag = (ModbusTag) optimizedReadRequest.getTag(tagName);
                String tagType = modbusTag.getClass().getSimpleName().substring("ModbusTag".length());
                if (!responses.containsKey(tagType)) {
                    responses.put(tagType, new ArrayList<>());
                }
                PlcResponseCode responseCode = optimizedReadResponse.getResponseCode(tagName);
                int startingAddress = modbusTag.getAddress();
                int endingAddressRegister = modbusTag.getAddress() + modbusTag.getNumberOfElements();
                int endingAddressCoil = modbusTag.getAddress() + modbusTag.getNumberOfElements();
                byte[] responseData = (responseCode == PlcResponseCode.OK) ? optimizedReadResponse.getPlcValue(tagName).getRaw() : null;
                responses.get(tagType).add(new Response(responseCode, startingAddress,
                    endingAddressRegister, endingAddressCoil, responseData));
            }

            // Now go through the original requests and try to answer them by using the raw data we now have.
            Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
            for (String tagName : readRequest.getTagNames()) {
                ModbusTag modbusTag = (ModbusTag) readRequest.getTag(tagName);
                String tagType = modbusTag.getClass().getSimpleName().substring("ModbusTag".length());
                if (!responses.containsKey(tagType)) {
                    values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null));
                    continue;
                }
                // Go through all responses till we find one where that contains the current tag's data.
                for (Response response : responses.get(tagType)) {
                    if((modbusTag instanceof ModbusTagCoil) || (modbusTag instanceof ModbusTagDiscreteInput)) {
                        if(response.matchesCoil(modbusTag)) {
                            // If this response was invalid, return all associated addresses as equally invalid.
                            // TODO: Possibly it would be worth doing a single item request for each of these
                            //  tags in order to find out which ones are actually invalid as if one item in the
                            //  current request exceeds the address range, all items in this chunk will fail, even
                            //  if only one element was invalid.
                            if(response.getResponseCode() != PlcResponseCode.OK) {
                                values.put(tagName, new DefaultPlcResponseItem<>(response.getResponseCode(), null));
                                break;
                            }

                            // Coils and Discrete Inputs are read completely different from registers.

                            // Calculate the byte that contains the response for this Coil
                            byte[] responseData = response.getResponseData();
                            int bitPosition = modbusTag.getAddress() - response.startingAddress;
                            int bytePosition = bitPosition / 8;
                            int bitPositionInByte = bitPosition % 8;
                            boolean isBitSet = (responseData[bytePosition] & (1 << bitPositionInByte)) != 0;
                            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcBOOL(isBitSet)));
                            break;
                        }
                    }
                    // Read a normal register.
                    else if (response.matchesRegister(modbusTag)) {
                        // If this response was invalid, return all associated addresses as equally invalid.
                        // TODO: Possibly it would be worth doing a single item request for each of these
                        //  tags in order to find out which ones are actually invalid as if one item in the
                        //  current request exceeds the address range, all items in this chunk will fail, even
                        //  if only one element was invalid.
                        if(response.getResponseCode() != PlcResponseCode.OK) {
                            values.put(tagName, new DefaultPlcResponseItem<>(response.getResponseCode(), null));
                            break;
                        }

                        var byteOrder = modbusTag.getByteOrder() != null ? modbusTag.getByteOrder() : modbusContext.getByteOrder();

                        byte[] responseData = response.getResponseDataForTag(modbusTag);
                        ReadBuffer readBuffer = getReadBuffer(responseData, byteOrder);
                        try {
                            PlcValue plcValue = DataItem.staticParse(readBuffer, modbusTag.getDataType(),
                                modbusTag.getNumberOfElements(),
                                byteOrder == ModbusByteOrder.BIG_ENDIAN);
                            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue));
                        } catch (ParseException e) {
                            values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                        }
                        break;
                    }
                }
                // If no response was found that contains the data, that's probably something we need to fix.
                if (!values.containsKey(tagName)) {
                    values.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }

            return new DefaultPlcReadResponse(readRequest, values);
        }
        // TODO: Remove this part once the driver is more stable.
        catch (Exception e) {
            try {
                logger.error("Error processing response:", e);
                WriteBufferXmlBased wb = new WriteBufferXmlBased();
                ((DefaultPlcReadRequest) readRequest).serialize(wb);
                logger.error("Original Request:\n{}", wb.getXmlString(), e);
                for (PlcReadRequest subRequest : readResponses.keySet()) {
                    SubResponse<PlcReadResponse> subResponse = readResponses.get(subRequest);
                    WriteBufferXmlBased wbSubRequest = new WriteBufferXmlBased();
                    ((DefaultPlcReadRequest) subRequest).serialize(wbSubRequest);
                    logger.error("Sub Request:\n{}", wbSubRequest.getXmlString());
                    if(subResponse.isSuccess()) {
                        PlcReadResponse plcReadResponse = subResponse.getResponse();
                        WriteBufferXmlBased wbSubResponse = new WriteBufferXmlBased();
                        ((DefaultPlcReadResponse) plcReadResponse).serialize(wbSubResponse);
                        logger.error("Sub Response (Success):\n{}", wbSubResponse.getXmlString());
                    } else {
                        Throwable throwable = subResponse.getThrowable();
                        logger.error("Sub Response (Error):", throwable);
                    }
                }
            } catch (SerializationException ex) {
                throw new RuntimeException(ex);
            }
            return null;
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Internal
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    protected List<PlcReadRequest> processCoilRequests(TreeSet<ModbusTag> tags, PlcReader reader, ModbusContext modbusContext) {
        List<PlcReadRequest> subRequests = new ArrayList<>();
        int firstCoil = -1;
        int lastCoil = -1;
        int maxCoilCurRequest = -1;
        for (ModbusTag tag : tags) {
            int sizeInCoils = tag.getDataType().getDataTypeSize() * 8;
            if (tag.getDataType() == ModbusDataType.BOOL) {
                sizeInCoils = 1;
            }
            // Initialize for the first item.
            if (firstCoil == -1) {
                firstCoil = tag.getAddress();
                lastCoil = tag.getAddress() + (sizeInCoils * tag.getNumberOfElements());
                // 2000 coils/request is the modbus limit.
                maxCoilCurRequest = tag.getAddress() + modbusContext.getMaxCoilsPerRequest();
            }

            // If adding the current coil would exceed the maximum number of coils that can be read by one request,
            // finish this one and start a new one.
            if (tag.getAddress() + (sizeInCoils * tag.getNumberOfElements()) > maxCoilCurRequest) {
                // Finish the current sub-request
                LinkedHashMap<String, PlcTagItem<PlcTag>> subTags = new LinkedHashMap<>();
                subTags.put("coils" + subRequests.size(), new DefaultPlcTagItem<>(new ModbusTagCoil(firstCoil, lastCoil - firstCoil, ModbusDataType.BYTE, Collections.emptyMap())));
                subRequests.add(new DefaultPlcReadRequest(reader, subTags));

                // Re-initialize the structures for the next request.
                firstCoil = tag.getAddress();
                lastCoil = tag.getAddress() + (sizeInCoils * tag.getNumberOfElements());
                maxCoilCurRequest = tag.getAddress() + modbusContext.getMaxCoilsPerRequest();
            }
            // Otherwise update the end-marker for the current block.
            else {
                lastCoil = tag.getAddress() + tag.getNumberOfElements();
            }
        }

        // Finish the last sub-request
        LinkedHashMap<String, PlcTagItem<PlcTag>> subTags = new LinkedHashMap<>();
        subTags.put("coils" + subRequests.size(), new DefaultPlcTagItem<>(new ModbusTagCoil(firstCoil, lastCoil - firstCoil, ModbusDataType.BYTE, Collections.emptyMap())));
        subRequests.add(new DefaultPlcReadRequest(reader, subTags));
        return subRequests;
    }

    protected List<PlcReadRequest> processRegisterRequests(TreeSet<ModbusTag> tags, PlcReader reader, TagFactory tagFactory, ModbusContext modbusContext) {
        List<PlcReadRequest> subRequests = new ArrayList<>();
        int firstRegister = -1;
        int lastRegister = -1;
        int maxRegisterCurRequest = -1;
        for (ModbusTag tag : tags) {
            int sizeInRegisters = (int) Math.ceil((double) tag.getDataType().getDataTypeSize() / 2);
            // Initialize for the first item.
            if (firstRegister == -1) {
                firstRegister = tag.getAddress();
                lastRegister = tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements());
                // 2000 coils/request is the modbus limit.
                maxRegisterCurRequest = tag.getAddress() + modbusContext.getMaxRegistersPerRequest();
            }

            // If adding the current coil would exceed the maximum number of coils that can be read by one request,
            // finish this one and start a new one.
            if (tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements()) > maxRegisterCurRequest) {
                // Finish the current sub-request
                LinkedHashMap<String, PlcTagItem<PlcTag>> subTags = new LinkedHashMap<>();
                subTags.put("registers" + subRequests.size(), new DefaultPlcTagItem<>(tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD)));
                subRequests.add(new DefaultPlcReadRequest(reader, subTags));

                // Re-initialize the structures for the next request.
                firstRegister = tag.getAddress();
                lastRegister = tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements());
                maxRegisterCurRequest = tag.getAddress() + modbusContext.getMaxRegistersPerRequest();
            }
            // Otherwise update the end-marker for the current block.
            else {
                lastRegister = tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements());
            }
        }

        // Finish the last sub-request
        LinkedHashMap<String, PlcTagItem<PlcTag>> subTags = new LinkedHashMap<>();
        subTags.put("registers" + subRequests.size(), new DefaultPlcTagItem<>(tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD)));
        subRequests.add(new DefaultPlcReadRequest(reader, subTags));
        return subRequests;
    }

    protected static class Response {
        private final PlcResponseCode responseCode;
        private final int startingAddress;
        private final int endingAddressCoil;
        private final int endingAddressRegister;
        private final byte[] responseData;

        public Response(PlcResponseCode responseCode, int startingAddress, int endingAddressRegister, int endingAddressCoil, byte[] responseData) {
            this.responseCode = responseCode;
            this.startingAddress = startingAddress;
            this.responseData = responseData;
            this.endingAddressCoil = endingAddressRegister;
            // In general a "Celil(responseData.length / 2)" would have been more correct,
            // but the data returned from the device should already be an even number.
            this.endingAddressRegister = endingAddressCoil;
        }

        public boolean matchesCoil(ModbusTag modbusTag) {
            int tagStartingAddress = modbusTag.getAddress();
            int tagEndingAddress = modbusTag.getAddress() + modbusTag.getNumberOfElements();
            return ((tagStartingAddress >= this.startingAddress) && (tagEndingAddress <= this.endingAddressCoil));
        }

        public boolean matchesRegister(ModbusTag modbusTag) {
            int tagStartingAddress = modbusTag.getAddress();
            int tagEndingAddress = modbusTag.getAddress() + (int) Math.ceil((double) modbusTag.getLengthBytes() / 2.0f);
            return ((tagStartingAddress >= this.startingAddress) && (tagEndingAddress <= this.endingAddressRegister));
        }

        public PlcResponseCode getResponseCode() {
            return responseCode;
        }

        public byte[] getResponseData() {
            return responseData;
        }

        public byte[] getResponseDataForTag(ModbusTag modbusTag) {
            byte[] itemData = new byte[modbusTag.getLengthBytes()];
            System.arraycopy(responseData, (modbusTag.getAddress() - startingAddress) * 2, itemData, 0, modbusTag.getLengthBytes());
            return itemData;
        }
    }

    protected interface TagFactory {
        PlcTag createTag(int address, int count, ModbusDataType dataType);
    }

    private ReadBuffer getReadBuffer(byte[] data, ModbusByteOrder byteOrder) {
        switch (byteOrder) {
            case LITTLE_ENDIAN: {
                // [4, 3, 2, 1]
                // [8, 7, 6, 5, 4, 3, 2, 1]
                return new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN);
            }
            case BIG_ENDIAN_BYTE_SWAP: {
                // [2, 1, 4, 3]
                // [2, 1, 4, 3, 6, 5, 8, 7]
                byte[] reordered = ModbusProtocolLogic.byteSwap(data);
                return new ReadBufferByteBased(reordered, ByteOrder.BIG_ENDIAN);
            }
            case LITTLE_ENDIAN_BYTE_SWAP: {
                // [3, 4, 1, 2]
                // [7, 8, 5, 6, 3, 4, 1, 2]
                byte[] reordered = ModbusProtocolLogic.byteSwap(data);
                return new ReadBufferByteBased(reordered, ByteOrder.LITTLE_ENDIAN);
            }
            default:
                // 16909060
                // [1, 2, 3, 4]
                // 72623859790382856
                // [1, 2, 3, 4, 5, 6, 7, 8]
                return new ReadBufferByteBased(data, ByteOrder.BIG_ENDIAN);
        }
    }

}
