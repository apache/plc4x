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
package org.apache.plc4x.java.umas.readwrite.context;

import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.umas.readwrite.UmasArrayDimension;
import org.apache.plc4x.java.umas.readwrite.UmasDataType;
import org.apache.plc4x.java.umas.readwrite.UmasUDTDefinition;
import org.apache.plc4x.java.umas.readwrite.UmasUnlocatedVariableReference;
import org.apache.plc4x.java.umas.readwrite.configuration.UmasConfiguration;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Maintains the runtime context for a UMAS driver connection.
 * Stores connection-specific state populated during the handshake
 * and used for subsequent read/write/browse operations.
 */
public class UmasDriverContext implements DriverContext, HasConfiguration<UmasConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(UmasDriverContext.class);

    private UmasConfiguration configuration;

    // Modbus transaction ID counter, wraps at 0xFFFF
    private final AtomicInteger transactionIdGenerator = new AtomicInteger(1);

    // PLC identification (populated during PlcIdent handshake step)
    private volatile String plcHostname;
    private volatile int plcModel;
    private volatile int plcFirmwareVersion;

    // Negotiated protocol parameters
    private volatile int maxFrameSize = 65535;
    private volatile short pairingKey;
    private volatile long hardwareId;
    private volatile long projectCrc;

    // Custom type definitions: type index -> field definitions
    private final Map<Integer, List<UmasUDTDefinition>> customTypeFields = new ConcurrentHashMap<>();
    private final Map<Integer, String> customTypeNames = new ConcurrentHashMap<>();
    private final Map<Integer, Integer> customTypeElementTypeIds = new ConcurrentHashMap<>();
    private final Map<Integer, List<UmasArrayDimension>> customTypeDimensions = new ConcurrentHashMap<>();

    // Symbol table: symbolic name -> variable reference
    private final Map<String, UmasUnlocatedVariableReference> symbolTable = new ConcurrentHashMap<>();

    // Data type table: type ID -> UmasDataType
    private final Map<Integer, UmasDataType> dataTypeTable = new ConcurrentHashMap<>();

    // Data type sizes from DD03: type ID -> allocated byte size
    private final Map<Integer, Integer> dataTypeSizes = new ConcurrentHashMap<>();

    // Per-symbol allocated sizes computed from symbol table layout (gap between adjacent symbols)
    private final Map<String, Integer> symbolSizes = new ConcurrentHashMap<>();

    @Override
    public void setConfiguration(UmasConfiguration configuration) {
        this.configuration = configuration;
        this.maxFrameSize = configuration.getMaxFrameSize();
    }

    public UmasConfiguration getConfiguration() {
        return configuration;
    }

    public int getNextTransactionId() {
        int id = transactionIdGenerator.getAndIncrement();
        if (id > 0xFFFF) {
            transactionIdGenerator.compareAndSet(id + 1, 1);
            return id & 0xFFFF;
        }
        return id;
    }

    // --- PLC identification ---

    public String getPlcHostname() {
        return plcHostname;
    }

    public void setPlcHostname(String plcHostname) {
        this.plcHostname = plcHostname;
    }

    public int getPlcModel() {
        return plcModel;
    }

    public void setPlcModel(int plcModel) {
        this.plcModel = plcModel;
    }

    public int getPlcFirmwareVersion() {
        return plcFirmwareVersion;
    }

    public void setPlcFirmwareVersion(int plcFirmwareVersion) {
        this.plcFirmwareVersion = plcFirmwareVersion;
    }

    // --- Negotiated parameters ---

    public int getMaxFrameSize() {
        return maxFrameSize;
    }

    public void setMaxFrameSize(int maxFrameSize) {
        this.maxFrameSize = maxFrameSize;
    }

    public short getPairingKey() {
        return pairingKey;
    }

    public void setPairingKey(short pairingKey) {
        this.pairingKey = pairingKey;
    }

    public long getHardwareId() {
        return hardwareId;
    }

    public void setHardwareId(long hardwareId) {
        this.hardwareId = hardwareId;
    }

    public long getProjectCrc() {
        return projectCrc;
    }

    public void setProjectCrc(long projectCrc) {
        this.projectCrc = projectCrc;
    }

    // --- Custom type operations ---

    public void addCustomType(int typeIndex, String typeName, List<UmasUDTDefinition> fields) {
        customTypeNames.put(typeIndex, typeName);
        customTypeFields.put(typeIndex, fields);
    }

    public void addArrayType(int typeIndex, String typeName, int elementTypeId, List<UmasArrayDimension> dimensions) {
        customTypeNames.put(typeIndex, typeName);
        customTypeFields.put(typeIndex, Collections.emptyList());
        customTypeElementTypeIds.put(typeIndex, elementTypeId);
        customTypeDimensions.put(typeIndex, dimensions);
    }

    public Optional<Integer> getArrayElementTypeId(int typeIndex) {
        return Optional.ofNullable(customTypeElementTypeIds.get(typeIndex));
    }

    public Optional<List<UmasArrayDimension>> getArrayDimensions(int typeIndex) {
        return Optional.ofNullable(customTypeDimensions.get(typeIndex));
    }

    public Optional<List<UmasUDTDefinition>> getCustomTypeFields(int typeIndex) {
        return Optional.ofNullable(customTypeFields.get(typeIndex));
    }

    public Optional<String> getCustomTypeName(int typeIndex) {
        return Optional.ofNullable(customTypeNames.get(typeIndex));
    }

    // --- Symbol table operations ---

    public void addSymbol(String name, UmasUnlocatedVariableReference reference) {
        symbolTable.put(name.toLowerCase(), reference);
    }

    public Optional<UmasUnlocatedVariableReference> getSymbol(String name) {
        return Optional.ofNullable(symbolTable.get(name.toLowerCase()));
    }

    public Map<String, UmasUnlocatedVariableReference> getSymbolTable() {
        return Collections.unmodifiableMap(symbolTable);
    }

    public int getSymbolCount() {
        return symbolTable.size();
    }

    // --- Data type table operations ---

    public void addDataType(int typeId, UmasDataType dataType) {
        dataTypeTable.put(typeId, dataType);
    }

    public Optional<UmasDataType> getDataType(int typeId) {
        return Optional.ofNullable(dataTypeTable.get(typeId));
    }

    public Map<Integer, UmasDataType> getDataTypeTable() {
        return Collections.unmodifiableMap(dataTypeTable);
    }

    // --- Data type size operations (from DD03) ---

    /**
     * Stores the allocated byte size for a data type from the DD03 data dictionary.
     * This is the total memory footprint, important for STRING and custom types
     * where the size is not derivable from the UmasDataType enum alone.
     *
     * @param typeId   the data type identifier
     * @param dataSize the allocated byte size from UmasDatatypeReference.dataSize
     */
    public void addDataTypeSize(int typeId, int dataSize) {
        dataTypeSizes.put(typeId, dataSize);
    }

    /**
     * Returns the allocated byte size for a data type, or empty if not registered.
     *
     * @param typeId the data type identifier
     * @return the allocated byte size
     */
    public Optional<Integer> getDataTypeSize(int typeId) {
        return Optional.ofNullable(dataTypeSizes.get(typeId));
    }

    // --- Per-symbol size operations ---

    /**
     * Returns the computed allocated byte size for a symbol, or empty if not available.
     *
     * @param name the symbolic name (case-insensitive)
     * @return the byte size if computed
     */
    public Optional<Integer> getSymbolSize(String name) {
        return Optional.ofNullable(symbolSizes.get(name.toLowerCase()));
    }

    /**
     * Computes the allocated byte size for each symbol from the symbol table layout.
     * For each memory block, symbols are sorted by offset and the size is computed
     * as the gap to the next symbol. The last symbol in each block gets no size entry.
     * Must be called after all symbols have been loaded via {@link #addSymbol}.
     */
    public void computeSymbolSizes() {
        symbolSizes.clear();

        // Group symbols by block
        java.util.Map<Integer, java.util.List<java.util.Map.Entry<String, UmasUnlocatedVariableReference>>> byBlock =
            new java.util.HashMap<>();
        for (var entry : symbolTable.entrySet()) {
            byBlock.computeIfAbsent(entry.getValue().getBlock(), k -> new java.util.ArrayList<>()).add(entry);
        }

        // For each block, sort by offset and compute gaps
        for (var blockEntry : byBlock.entrySet()) {
            var symbols = blockEntry.getValue();
            symbols.sort(java.util.Comparator.comparingLong(e -> e.getValue().getOffset()));

            for (int i = 0; i < symbols.size() - 1; i++) {
                var current = symbols.get(i);
                var next = symbols.get(i + 1);
                int size = (int) (next.getValue().getOffset() - current.getValue().getOffset());
                if (size > 0) {
                    symbolSizes.put(current.getKey(), size);
                }
            }
        }

        LOGGER.debug("Computed sizes for {} of {} symbols", symbolSizes.size(), symbolTable.size());
    }

}
