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
package org.apache.plc4x.java.modbus.base.optimizer;

import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.*;
import org.apache.plc4x.java.modbus.base.ModbusRegisterCodec;
import org.apache.plc4x.java.modbus.readwrite.DataItem;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

/**
 * Optimizer for Modbus read requests that merges adjacent tags into larger block reads.
 * <p>
 * Tags of the same type (coils, holding registers, etc.) that are in close proximity
 * are merged into a single larger read request. The response data is then split apart
 * to serve each original tag.
 */
public class ModbusReadOptimizer {

    private static final Logger LOGGER = LoggerFactory.getLogger(ModbusReadOptimizer.class);

    private final int maxCoilsPerRequest;
    private final int maxRegistersPerRequest;
    private final ModbusByteOrder defaultByteOrder;

    public ModbusReadOptimizer(int maxCoilsPerRequest, int maxRegistersPerRequest, ModbusByteOrder defaultByteOrder) {
        this.maxCoilsPerRequest = maxCoilsPerRequest;
        this.maxRegistersPerRequest = maxRegistersPerRequest;
        this.defaultByteOrder = defaultByteOrder;
    }

    /**
     * Groups tags by type and unit-id and merges adjacent ones into optimized block reads.
     * Returns a list of OptimizedRead objects, each containing a merged tag to read
     * and the original tag names it covers.
     * <p>
     * A block read is addressed to exactly one unit, so tags carrying different unit-ids must
     * never end up in the same request, and the unit-id of a group has to be passed on to the
     * merged tag (otherwise the connection would fall back to its default unit-id).
     */
    public List<OptimizedRead> optimizeReads(Map<String, ModbusTag> tagsByName) {
        // Sort tags by type and unit-id
        Map<Short, TreeMap<String, ModbusTag>> coils = newUnitIdGroups();
        Map<Short, TreeMap<String, ModbusTag>> holdingRegisters = newUnitIdGroups();
        Map<Short, TreeMap<String, ModbusTag>> inputRegisters = newUnitIdGroups();
        Map<Short, TreeMap<String, ModbusTag>> extendedRegisters = newUnitIdGroups();
        Map<Short, TreeMap<String, ModbusTag>> discreteInputs = newUnitIdGroups();

        for (Map.Entry<String, ModbusTag> entry : tagsByName.entrySet()) {
            ModbusTag tag = entry.getValue();
            if (tag instanceof ModbusTagCoil) addToUnitIdGroup(coils, entry);
            else if (tag instanceof ModbusTagHoldingRegister) addToUnitIdGroup(holdingRegisters, entry);
            else if (tag instanceof ModbusTagInputRegister) addToUnitIdGroup(inputRegisters, entry);
            else if (tag instanceof ModbusTagExtendedRegister) addToUnitIdGroup(extendedRegisters, entry);
            else if (tag instanceof ModbusTagDiscreteInput) addToUnitIdGroup(discreteInputs, entry);
        }

        List<OptimizedRead> result = new ArrayList<>();
        for (Map.Entry<Short, TreeMap<String, ModbusTag>> group : coils.entrySet()) {
            result.addAll(optimizeCoils(group.getKey(), group.getValue()));
        }
        for (Map.Entry<Short, TreeMap<String, ModbusTag>> group : holdingRegisters.entrySet()) {
            result.addAll(optimizeRegisters(group.getKey(), group.getValue(), ModbusReadOptimizer::createHoldingRegister));
        }
        for (Map.Entry<Short, TreeMap<String, ModbusTag>> group : inputRegisters.entrySet()) {
            result.addAll(optimizeRegisters(group.getKey(), group.getValue(), ModbusReadOptimizer::createInputRegister));
        }
        for (Map.Entry<Short, TreeMap<String, ModbusTag>> group : extendedRegisters.entrySet()) {
            result.addAll(optimizeRegisters(group.getKey(), group.getValue(), ModbusReadOptimizer::createExtendedRegister));
        }
        for (Map.Entry<Short, TreeMap<String, ModbusTag>> group : discreteInputs.entrySet()) {
            result.addAll(optimizeCoils(group.getKey(), group.getValue()));
        }
        return result;
    }

    /**
     * Given the raw response data for an optimized block read, extracts the values for each original tag.
     */
    public Map<String, PlcResponseItem<PlcValue>> splitResponse(
            OptimizedRead optimizedRead, PlcResponseCode blockResponseCode, byte[] blockData) {

        Map<String, PlcResponseItem<PlcValue>> result = new LinkedHashMap<>();

        if (blockResponseCode != PlcResponseCode.OK) {
            for (String tagName : optimizedRead.originalTagNames.keySet()) {
                result.put(tagName, new DefaultPlcResponseItem<>(blockResponseCode, null));
            }
            return result;
        }

        ModbusTag blockTag = optimizedRead.mergedTag;
        for (Map.Entry<String, ModbusTag> entry : optimizedRead.originalTagNames.entrySet()) {
            String tagName = entry.getKey();
            ModbusTag originalTag = entry.getValue();

            try {
                if (blockTag instanceof ModbusTagCoil || blockTag instanceof ModbusTagDiscreteInput) {
                    // Coils/discrete inputs: bit-level extraction. An array tag (BOOL[n]) occupies
                    // n consecutive coils - the block read already covers all of them (see
                    // optimizeCoils), so every element has to be extracted, not just the first.
                    if (originalTag.getDataType() != ModbusDataType.BOOL) {
                        // A coil carries a single bit. Assembling coils into wider types is not
                        // implemented - report that instead of silently returning the first bit.
                        LOGGER.warn("Reading coils/discrete inputs as {} is not supported (tag '{}'), only BOOL is.",
                            originalTag.getDataType(), tagName);
                        result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
                        continue;
                    }
                    int firstBitPosition = originalTag.getAddress() - blockTag.getAddress();
                    int numberOfElements = originalTag.getNumberOfElements();
                    int lastBytePosition = (firstBitPosition + numberOfElements - 1) / 8;
                    if (lastBytePosition >= blockData.length) {
                        // The device returned fewer coils than we asked for.
                        result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                        continue;
                    }
                    List<PlcValue> values = new ArrayList<>(numberOfElements);
                    for (int i = 0; i < numberOfElements; i++) {
                        int bitPosition = firstBitPosition + i;
                        int bytePosition = bitPosition / 8;
                        int bitPositionInByte = bitPosition % 8;
                        boolean isBitSet = (blockData[bytePosition] & (1 << bitPositionInByte)) != 0;
                        values.add(new PlcBOOL(isBitSet));
                    }
                    PlcValue plcValue = numberOfElements == 1 ? values.getFirst() : new PlcList(values);
                    result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue));
                } else {
                    // Registers: byte-level extraction
                    int byteOffset = (originalTag.getAddress() - blockTag.getAddress()) * 2;
                    int byteLength = originalTag.getLengthBytes();
                    byte[] tagData = new byte[byteLength];
                    System.arraycopy(blockData, byteOffset, tagData, 0, byteLength);

                    ModbusByteOrder byteOrder = originalTag.getByteOrder() != null ? originalTag.getByteOrder() : defaultByteOrder;
                    if (byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
                        tagData = byteSwap(tagData);
                    }
                    boolean bigEndian = (byteOrder == ModbusByteOrder.BIG_ENDIAN || byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);

                    ReadBufferByteBased readBuffer;
                    if (!bigEndian) {
                        readBuffer = new ReadBufferByteBased(tagData,
                            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
                            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                            WithOption.WithSignedIntegerEncoding("twos-complement"),
                            WithOption.WithFloatEncoding("IEEE754"),
                            WithOption.WithStringEncoding("UTF8"));
                    } else {
                        readBuffer = new ReadBufferByteBased(tagData,
                            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                            WithOption.WithSignedIntegerEncoding("twos-complement"),
                            WithOption.WithFloatEncoding("IEEE754"),
                            WithOption.WithStringEncoding("UTF8"));
                    }
                    PlcValue plcValue = ModbusRegisterCodec.parse(readBuffer, originalTag.getDataType(),
                        originalTag.getNumberOfElements(), bigEndian, originalTag.getStringLength());
                    result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, plcValue));
                }
            } catch (BufferException e) {
                result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
            }
        }
        return result;
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Internal
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    /**
     * Creates a map for grouping tags by their unit-id. Tags without an explicit unit-id (which
     * will be read using the connection's default unit-id) are grouped under the {@code null} key.
     */
    private static Map<Short, TreeMap<String, ModbusTag>> newUnitIdGroups() {
        return new TreeMap<>(Comparator.nullsFirst(Comparator.naturalOrder()));
    }

    private static void addToUnitIdGroup(Map<Short, TreeMap<String, ModbusTag>> groups,
                                         Map.Entry<String, ModbusTag> entry) {
        groups.computeIfAbsent(entry.getValue().getUnitId(), unitId -> new TreeMap<>())
            .put(entry.getKey(), entry.getValue());
    }

    /**
     * Builds the config of a merged tag, so that the unit-id of the tags it was built from is
     * used when sending the request.
     */
    private static Map<String, String> mergedTagConfig(Short unitId) {
        return unitId == null ? Collections.emptyMap() : Collections.singletonMap("unit-id", unitId.toString());
    }

    private List<OptimizedRead> optimizeCoils(Short unitId, Map<String, ModbusTag> tagsByName) {
        // Sort by address
        List<Map.Entry<String, ModbusTag>> sorted = new ArrayList<>(tagsByName.entrySet());
        sorted.sort(Comparator.comparingInt(e -> e.getValue().getAddress()));

        List<OptimizedRead> result = new ArrayList<>();
        int firstAddress = -1;
        int lastAddress = -1;
        int maxAddress = -1;
        LinkedHashMap<String, ModbusTag> currentGroup = new LinkedHashMap<>();

        for (Map.Entry<String, ModbusTag> entry : sorted) {
            ModbusTag tag = entry.getValue();
            int sizeInCoils = tag.getDataType() == ModbusDataType.BOOL ? 1 : tag.getDataType().getDataTypeSize() * 8;
            int tagEnd = tag.getAddress() + (sizeInCoils * tag.getNumberOfElements());

            if (firstAddress == -1) {
                firstAddress = tag.getAddress();
                lastAddress = tagEnd;
                maxAddress = tag.getAddress() + maxCoilsPerRequest;
                currentGroup.put(entry.getKey(), tag);
                continue;
            }

            if (tagEnd > maxAddress) {
                // Finish current group
                boolean isDiscreteInput = currentGroup.values().iterator().next() instanceof ModbusTagDiscreteInput;
                ModbusTag mergedTag = isDiscreteInput
                    ? new ModbusTagDiscreteInput(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, mergedTagConfig(unitId))
                    : new ModbusTagCoil(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, mergedTagConfig(unitId));
                result.add(new OptimizedRead(mergedTag, currentGroup));

                // Start new group
                currentGroup = new LinkedHashMap<>();
                firstAddress = tag.getAddress();
                lastAddress = tagEnd;
                maxAddress = tag.getAddress() + maxCoilsPerRequest;
            } else {
                lastAddress = Math.max(lastAddress, tagEnd);
            }
            currentGroup.put(entry.getKey(), tag);
        }

        // Finish last group
        if (!currentGroup.isEmpty()) {
            boolean isDiscreteInput = currentGroup.values().iterator().next() instanceof ModbusTagDiscreteInput;
            ModbusTag mergedTag = isDiscreteInput
                ? new ModbusTagDiscreteInput(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, mergedTagConfig(unitId))
                : new ModbusTagCoil(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, mergedTagConfig(unitId));
            result.add(new OptimizedRead(mergedTag, currentGroup));
        }
        return result;
    }

    private List<OptimizedRead> optimizeRegisters(Short unitId, Map<String, ModbusTag> tagsByName, TagFactory tagFactory) {
        List<Map.Entry<String, ModbusTag>> sorted = new ArrayList<>(tagsByName.entrySet());
        sorted.sort(Comparator.comparingInt(e -> e.getValue().getAddress()));

        List<OptimizedRead> result = new ArrayList<>();
        int firstRegister = -1;
        int lastRegister = -1;
        int maxRegister = -1;
        LinkedHashMap<String, ModbusTag> currentGroup = new LinkedHashMap<>();

        for (Map.Entry<String, ModbusTag> entry : sorted) {
            ModbusTag tag = entry.getValue();
            int sizeInRegisters = (int) Math.ceil((double) tag.getDataType().getDataTypeSize() / 2);
            int tagEnd = tag.getAddress() + (sizeInRegisters * tag.getNumberOfElements());

            if (firstRegister == -1) {
                firstRegister = tag.getAddress();
                lastRegister = tagEnd;
                maxRegister = tag.getAddress() + maxRegistersPerRequest;
                currentGroup.put(entry.getKey(), tag);
                continue;
            }

            if (tagEnd > maxRegister) {
                // Finish current group
                result.add(new OptimizedRead(
                    tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD, mergedTagConfig(unitId)),
                    currentGroup));

                // Start new group
                currentGroup = new LinkedHashMap<>();
                firstRegister = tag.getAddress();
                lastRegister = tagEnd;
                maxRegister = tag.getAddress() + maxRegistersPerRequest;
            } else {
                lastRegister = Math.max(lastRegister, tagEnd);
            }
            currentGroup.put(entry.getKey(), tag);
        }

        if (!currentGroup.isEmpty()) {
            result.add(new OptimizedRead(
                tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD, mergedTagConfig(unitId)),
                currentGroup));
        }
        return result;
    }

    private static ModbusTag createHoldingRegister(int address, int count, ModbusDataType dataType, Map<String, String> config) {
        return new ModbusTagHoldingRegister(address, count, dataType, config);
    }

    private static ModbusTag createInputRegister(int address, int count, ModbusDataType dataType, Map<String, String> config) {
        return new ModbusTagInputRegister(address, count, dataType, config);
    }

    private static ModbusTag createExtendedRegister(int address, int count, ModbusDataType dataType, Map<String, String> config) {
        return new ModbusTagExtendedRegister(address, count, dataType, config);
    }

    private static byte[] byteSwap(byte[] in) {
        byte[] out = new byte[in.length];
        for (int i = 0; i < out.length - 1; i += 2) {
            out[i] = in[i + 1];
            out[i + 1] = in[i];
        }
        if (in.length % 2 != 0) {
            out[in.length - 1] = in[in.length - 1];
        }
        return out;
    }

    /**
     * Represents an optimized read: a merged tag covering multiple original tags.
     */
    public static class OptimizedRead {
        public final ModbusTag mergedTag;
        public final LinkedHashMap<String, ModbusTag> originalTagNames;

        public OptimizedRead(ModbusTag mergedTag, LinkedHashMap<String, ModbusTag> originalTagNames) {
            this.mergedTag = mergedTag;
            this.originalTagNames = originalTagNames;
        }
    }

    @FunctionalInterface
    private interface TagFactory {
        ModbusTag createTag(int address, int count, ModbusDataType dataType, Map<String, String> config);
    }

}
