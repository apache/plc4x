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

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.*;
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

import java.util.*;

/**
 * Optimizer for Modbus read requests that merges adjacent tags into larger block reads.
 * <p>
 * Tags of the same type (coils, holding registers, etc.) that are in close proximity
 * are merged into a single larger read request. The response data is then split apart
 * to serve each original tag.
 */
public class ModbusReadOptimizer {

    private final int maxCoilsPerRequest;
    private final int maxRegistersPerRequest;
    private final ModbusByteOrder defaultByteOrder;

    public ModbusReadOptimizer(int maxCoilsPerRequest, int maxRegistersPerRequest, ModbusByteOrder defaultByteOrder) {
        this.maxCoilsPerRequest = maxCoilsPerRequest;
        this.maxRegistersPerRequest = maxRegistersPerRequest;
        this.defaultByteOrder = defaultByteOrder;
    }

    /**
     * Groups tags by type and merges adjacent ones into optimized block reads.
     * Returns a list of OptimizedRead objects, each containing a merged tag to read
     * and the original tag names it covers.
     */
    public List<OptimizedRead> optimizeReads(Map<String, ModbusTag> tagsByName) {
        // Sort tags by type
        TreeMap<String, ModbusTag> coils = new TreeMap<>();
        TreeMap<String, ModbusTag> holdingRegisters = new TreeMap<>();
        TreeMap<String, ModbusTag> inputRegisters = new TreeMap<>();
        TreeMap<String, ModbusTag> extendedRegisters = new TreeMap<>();
        TreeMap<String, ModbusTag> discreteInputs = new TreeMap<>();

        for (Map.Entry<String, ModbusTag> entry : tagsByName.entrySet()) {
            ModbusTag tag = entry.getValue();
            if (tag instanceof ModbusTagCoil) coils.put(entry.getKey(), tag);
            else if (tag instanceof ModbusTagHoldingRegister) holdingRegisters.put(entry.getKey(), tag);
            else if (tag instanceof ModbusTagInputRegister) inputRegisters.put(entry.getKey(), tag);
            else if (tag instanceof ModbusTagExtendedRegister) extendedRegisters.put(entry.getKey(), tag);
            else if (tag instanceof ModbusTagDiscreteInput) discreteInputs.put(entry.getKey(), tag);
        }

        List<OptimizedRead> result = new ArrayList<>();
        if (!coils.isEmpty()) result.addAll(optimizeCoils(coils));
        if (!holdingRegisters.isEmpty()) result.addAll(optimizeRegisters(holdingRegisters, ModbusReadOptimizer::createHoldingRegister));
        if (!inputRegisters.isEmpty()) result.addAll(optimizeRegisters(inputRegisters, ModbusReadOptimizer::createInputRegister));
        if (!extendedRegisters.isEmpty()) result.addAll(optimizeRegisters(extendedRegisters, ModbusReadOptimizer::createExtendedRegister));
        if (!discreteInputs.isEmpty()) result.addAll(optimizeCoils(discreteInputs));
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
                    // Coils/discrete inputs: bit-level extraction
                    int bitPosition = originalTag.getAddress() - blockTag.getAddress();
                    int bytePosition = bitPosition / 8;
                    int bitPositionInByte = bitPosition % 8;
                    boolean isBitSet = (blockData[bytePosition] & (1 << bitPositionInByte)) != 0;
                    result.put(tagName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcBOOL(isBitSet)));
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
                    PlcValue plcValue = DataItem.staticParse(readBuffer, originalTag.getDataType(),
                        originalTag.getNumberOfElements(), bigEndian);
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

    private List<OptimizedRead> optimizeCoils(Map<String, ModbusTag> tagsByName) {
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
                    ? new ModbusTagDiscreteInput(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, Collections.emptyMap())
                    : new ModbusTagCoil(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, Collections.emptyMap());
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
                ? new ModbusTagDiscreteInput(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, Collections.emptyMap())
                : new ModbusTagCoil(firstAddress, lastAddress - firstAddress, ModbusDataType.BYTE, Collections.emptyMap());
            result.add(new OptimizedRead(mergedTag, currentGroup));
        }
        return result;
    }

    private List<OptimizedRead> optimizeRegisters(Map<String, ModbusTag> tagsByName, TagFactory tagFactory) {
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
                    tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD),
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
                tagFactory.createTag(firstRegister, lastRegister - firstRegister, ModbusDataType.WORD),
                currentGroup));
        }
        return result;
    }

    private static ModbusTag createHoldingRegister(int address, int count, ModbusDataType dataType) {
        return new ModbusTagHoldingRegister(address, count, dataType, Collections.emptyMap());
    }

    private static ModbusTag createInputRegister(int address, int count, ModbusDataType dataType) {
        return new ModbusTagInputRegister(address, count, dataType, Collections.emptyMap());
    }

    private static ModbusTag createExtendedRegister(int address, int count, ModbusDataType dataType) {
        return new ModbusTagExtendedRegister(address, count, dataType, Collections.emptyMap());
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
        ModbusTag createTag(int address, int count, ModbusDataType dataType);
    }

}
