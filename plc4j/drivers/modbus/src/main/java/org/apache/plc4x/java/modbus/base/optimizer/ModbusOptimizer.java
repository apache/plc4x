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
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagCoil;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagDiscreteInput;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagExtendedRegister;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagInputRegister;
import org.apache.plc4x.java.modbus.readwrite.DataItem;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.optimizer.SingleTagOptimizer;

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
            subRequests.addAll(processCoilRequests(coils, reader));
        }
        if (holdingRegisters != null) {
            subRequests.addAll(processRegisterRequests(holdingRegisters, reader, (address, count, dataType) -> new ModbusTagHoldingRegister(address, count, dataType, Collections.emptyMap())));
        }
        if (inputRegisters != null) {
            subRequests.addAll(processRegisterRequests(inputRegisters, reader, (address, count, dataType) -> new ModbusTagInputRegister(address, count, dataType, Collections.emptyMap())));
        }
        if (extendedRegisters != null) {
            subRequests.addAll(processRegisterRequests(extendedRegisters, reader, (address, count, dataType) -> new ModbusTagExtendedRegister(address, count, dataType, Collections.emptyMap())));
        }
        if (discreteInputs != null) {
            subRequests.addAll(processRegisterRequests(discreteInputs, reader, (address, count, dataType) -> new ModbusTagDiscreteInput(address, count, dataType, Collections.emptyMap())));
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
    protected PlcReadResponse processReadResponses(PlcReadRequest readRequest, Map<PlcReadRequest, SubResponse<PlcReadResponse>> readResponses) {
        // Build an index of all the data returned by all requests.
        // This data should contain all the bits needed to create the response of the original request.
        Map<String, List<Response>> responses = new HashMap<>();
        for (PlcReadRequest optimizedReadRequest : readResponses.keySet()) {
            PlcReadResponse optimizedReadResponse = readResponses.get(optimizedReadRequest).getResponse();
            if(optimizedReadResponse == null) {
                continue;
            }
            // Optimized read requests only contain one ModbusTag.
            String tagName = optimizedReadRequest.getTagNames().stream().findFirst().orElse(null);
            if(tagName == null) {
                continue;
            }
            ModbusTag modbusTag = (ModbusTag) optimizedReadRequest.getTag(tagName);
            String tagType = modbusTag.getClass().getSimpleName().substring("ModbusTag".length());
            if(!responses.containsKey(tagType)) {
                responses.put(tagType, new ArrayList<>());
            }
            responses.get(tagType).add(new Response(modbusTag.getAddress(), optimizedReadResponse.getPlcValue(tagName).getRaw()));
        }

        // Now go through the original requests and try to answer them by using the raw data we now have.
        Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
        for (String tagName : readRequest.getTagNames()) {
            ModbusTag modbusTag = (ModbusTag) readRequest.getTag(tagName);
            String tagType = modbusTag.getClass().getSimpleName().substring("ModbusTag".length());
            if(!responses.containsKey(tagType)) {
                values.put(tagName, new ResponseItem<>(PlcResponseCode.NOT_FOUND, null));
                continue;
            }
            // Go through all responses till we find one where that contains the current tag's data.
            for (Response response : responses.get(tagType)) {
                if(response.matches(modbusTag)) {
                    byte[] responseData = response.getResponseData(modbusTag);
                    ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(responseData);
                    try {
                        PlcValue plcValue = DataItem.staticParse(readBufferByteBased, modbusTag.getDataType(), modbusTag.getNumberOfElements(), true);
                        values.put(tagName, new ResponseItem<>(PlcResponseCode.OK, plcValue));
                    } catch (ParseException e) {
                        values.put(tagName, new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                    }
                    break;
                }
            }
            // If no response was found that contains the data, that's probably something we need to fix.
            if(!values.containsKey(tagName)) {
                values.put(tagName, new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
            }
        }

        return new DefaultPlcReadResponse(readRequest, values);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Internal
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    protected List<PlcReadRequest> processCoilRequests(TreeSet<ModbusTag> tags, PlcReader reader) {
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
                maxCoilCurRequest = tag.getAddress() + 2000;
            }

            // If adding the current coil would exceed the maximum number of coils that can be read by one request,
            // finish this one and start a new one.
            if (tag.getAddress() + (sizeInCoils * tag.getNumberOfElements()) > maxCoilCurRequest) {
                // Finish the current sub-request
                LinkedHashMap<String, PlcTag> subTags = new LinkedHashMap<>();
                subTags.put("coils" + subRequests.size(), new ModbusTagCoil(firstCoil, lastCoil - firstCoil, ModbusDataType.BYTE, Collections.emptyMap()));
                subRequests.add(new DefaultPlcReadRequest(reader, subTags));

                // Re-initialize the structures for the next request.
                firstCoil = tag.getAddress();
                lastCoil = tag.getAddress() + (sizeInCoils * tag.getNumberOfElements());
                maxCoilCurRequest = tag.getAddress() + 2000;
            }
            // Otherwise update the end-marker for the current block.
            else {
                lastCoil = tag.getAddress() + tag.getNumberOfElements();
            }
        }

        // Finish the last sub-request
        LinkedHashMap<String, PlcTag> subTags = new LinkedHashMap<>();
        subTags.put("coils" + subRequests.size(), new ModbusTagCoil(firstCoil, lastCoil - firstCoil, ModbusDataType.BYTE, Collections.emptyMap()));
        subRequests.add(new DefaultPlcReadRequest(reader, subTags));
        return subRequests;
    }

    protected List<PlcReadRequest> processRegisterRequests(TreeSet<ModbusTag> tags, PlcReader reader, TagFactory tagFactory) {
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
                maxRegisterCurRequest = tag.getAddress() + 125;
            }

            // If adding the current coil would exceed the maximum number of coils that can be read by one request,
            // finish this one and start a new one.
            if ((tag.getAddress() - firstRegister) + (sizeInRegisters * tag.getNumberOfElements()) > maxRegisterCurRequest) {
                // Finish the current sub-request
                LinkedHashMap<String, PlcTag> subTags = new LinkedHashMap<>();
                subTags.put("registers" + subRequests.size(), tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD));
                subRequests.add(new DefaultPlcReadRequest(reader, subTags));

                // Re-initialize the structures for the next request.
                firstRegister = tag.getAddress();
                lastRegister = tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements());
                maxRegisterCurRequest = tag.getAddress() + 125;
            }
            // Otherwise update the end-marker for the current block.
            else {
                lastRegister = tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements());
            }
        }

        // Finish the last sub-request
        LinkedHashMap<String, PlcTag> subTags = new LinkedHashMap<>();
        subTags.put("registers" + subRequests.size(), tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD));
        subRequests.add(new DefaultPlcReadRequest(reader, subTags));
        return subRequests;
    }

    protected static class Response {
        private final int startingAddress;
        private final byte[] responseData;

        public Response(int startingAddress, byte[] responseData) {
            this.startingAddress = startingAddress;
            this.responseData = responseData;
        }

        public boolean matches(ModbusTag modbusTag) {
            //boolean isRegisterTag = !(modbusTag instanceof ModbusTagCoil);
            return (modbusTag.getAddress() >= startingAddress) && (((modbusTag.getAddress() - startingAddress) * 2) + modbusTag.getLengthBytes() <= startingAddress + responseData.length);
        }

        public byte[] getResponseData(ModbusTag modbusTag) {
            byte[] itemData = new byte[modbusTag.getLengthBytes()];
            System.arraycopy(responseData, (modbusTag.getAddress() - startingAddress) * 2, itemData, 0, modbusTag.getLengthBytes());
            return itemData;
        }
    }

    protected interface TagFactory {
        PlcTag createTag(int address, int count, ModbusDataType dataType);
    }

}
