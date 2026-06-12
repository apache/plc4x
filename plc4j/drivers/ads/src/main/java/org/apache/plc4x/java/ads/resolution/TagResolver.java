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
package org.apache.plc4x.java.ads.resolution;

import org.apache.plc4x.java.ads.readwrite.AdsDataType;
import org.apache.plc4x.java.ads.readwrite.AdsDataTypeArrayInfo;
import org.apache.plc4x.java.ads.readwrite.AdsDataTypeTableEntry;
import org.apache.plc4x.java.ads.readwrite.AdsSymbolTableEntry;
import org.apache.plc4x.java.ads.tag.SymbolicAdsTag;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * Walks the ADS symbol- and data-type tables to translate a {@link SymbolicAdsTag} into a
 * {@link ResolvedAdsTag} (absolute group/offset/size + decoder shape).
 */
public final class TagResolver {

    private final Map<String, AdsSymbolTableEntry> symbolTable;
    private final Map<String, AdsDataTypeTableEntry> dataTypeTable;

    public TagResolver(Map<String, AdsSymbolTableEntry> symbolTable,
                       Map<String, AdsDataTypeTableEntry> dataTypeTable) {
        this.symbolTable = symbolTable;
        this.dataTypeTable = dataTypeTable;
    }

    public ResolvedAdsTag resolve(SymbolicAdsTag tag) {
        AddressParser.AddressPart root = AddressParser.parse(tag.getSymbolicAddress());
        AdsSymbolTableEntry symbol = symbolTable.get(root.baseSegment());
        if (symbol == null) {
            throw new PlcInvalidTagException("Unknown symbol: " + root.baseSegment());
        }
        AdsDataTypeTableEntry dataType = dataTypeTable.get(symbol.getDataTypeName());
        if (dataType == null) {
            throw new PlcInvalidTagException(
                "Unknown data type for symbol " + root.baseSegment() + ": " + symbol.getDataTypeName());
        }
        return resolvePart(symbol.getGroup(), symbol.getOffset(), dataType,
            root.arrayIndices(), root.child());
    }

    private ResolvedAdsTag resolvePart(long indexGroup, long indexOffset,
                                       AdsDataTypeTableEntry dataType,
                                       List<Integer> arrayIndices,
                                       AddressParser.AddressPart child) {
        if (!arrayIndices.isEmpty()) {
            return resolveArray(indexGroup, indexOffset, dataType, arrayIndices, child);
        }
        if (child != null) {
            return resolveChild(indexGroup, indexOffset, dataType, child);
        }
        return finalizeLeaf(indexGroup, indexOffset, dataType, Collections.emptyList());
    }

    private ResolvedAdsTag resolveArray(long indexGroup, long indexOffset,
                                        AdsDataTypeTableEntry dataType,
                                        List<Integer> arrayIndices,
                                        AddressParser.AddressPart child) {
        if (dataType.getArrayInfo().isEmpty()) {
            throw new PlcInvalidTagException(
                "Array index applied to non-array type " + dataType.getMainName());
        }
        if (arrayIndices.size() > dataType.getArrayInfo().size()) {
            throw new PlcInvalidTagException("Too many array indices for type " + dataType.getMainName()
                + " (got " + arrayIndices.size() + ", type has " + dataType.getArrayInfo().size() + ")");
        }

        long elementSize = dataType.getSize();
        for (int i = 0; i < arrayIndices.size(); i++) {
            int idx = arrayIndices.get(i);
            AdsDataTypeArrayInfo dim = dataType.getArrayInfo().get(i);
            if (idx < dim.getLowerBound() || idx > dim.getUpperBound()) {
                throw new PlcInvalidTagException(String.format(
                    "Index %d out of bounds [%d..%d] for dimension %d of %s",
                    idx, dim.getLowerBound(), dim.getUpperBound(), i, dataType.getMainName()));
            }
            elementSize /= dim.getNumElements();
            indexOffset += (idx - dim.getLowerBound()) * elementSize;
        }

        boolean fullyIndexed = arrayIndices.size() == dataType.getArrayInfo().size();
        if (!fullyIndexed && child != null) {
            throw new PlcInvalidTagException(
                "Field access requires all array dimensions to be specified for "
                    + dataType.getMainName() + " (have " + arrayIndices.size() + " of "
                    + dataType.getArrayInfo().size() + ")");
        }

        if (fullyIndexed) {
            // Recurse into the element type with no array indices.
            String elementTypeName = dataType.getSecondaryName();
            AdsDataTypeTableEntry elementDataType = dataTypeTable.get(elementTypeName);
            if (elementDataType == null) {
                // Primitive element (incl. STRING(n)) — synthesize a leaf result.
                if (child != null) {
                    throw new PlcInvalidTagException(
                        "Cannot descend into primitive element type " + elementTypeName);
                }
                return primitiveLeaf(indexGroup, indexOffset, elementTypeName, elementSize);
            }
            return resolvePart(indexGroup, indexOffset, elementDataType,
                Collections.emptyList(), child);
        }

        // Partial-dim read: carry remaining dims so the decoder can produce nested PlcLists.
        List<AdsDataTypeArrayInfo> remaining = dataType.getArrayInfo()
            .subList(arrayIndices.size(), dataType.getArrayInfo().size());

        // Total bytes for the partial slice = (size of one fully-indexed element of the OUTER dims)
        // We've already divided elementSize for each indexed dim; what remains covers the rest.
        long sizeInBytes = elementSize;

        return new ResolvedAdsTag(indexGroup, indexOffset, sizeInBytes,
            dataType.getMainName(), PlcValueType.List, 0, remaining);
    }

    private ResolvedAdsTag resolveChild(long indexGroup, long indexOffset,
                                        AdsDataTypeTableEntry dataType,
                                        AddressParser.AddressPart child) {
        AdsDataTypeTableEntry fieldEntry = null;
        for (AdsDataTypeTableEntry c : dataType.getChildren()) {
            if (c.getMainName().equals(child.baseSegment())) {
                fieldEntry = c;
                break;
            }
        }
        if (fieldEntry == null) {
            throw new PlcInvalidTagException(
                "Field '" + child.baseSegment() + "' not found in struct " + dataType.getMainName());
        }
        AdsDataTypeTableEntry fieldType = dataTypeTable.get(fieldEntry.getSecondaryName());
        if (fieldType == null) {
            // Primitive field type (TC2 fallback). Only valid if there's nothing further down.
            if (!child.arrayIndices().isEmpty() || child.child() != null) {
                throw new PlcInvalidTagException(
                    "Cannot descend into primitive field type " + fieldEntry.getSecondaryName());
            }
            return primitiveLeaf(indexGroup, indexOffset + fieldEntry.getOffset(),
                fieldEntry.getSecondaryName(), fieldEntry.getSize());
        }
        return resolvePart(indexGroup, indexOffset + fieldEntry.getOffset(),
            fieldType, child.arrayIndices(), child.child());
    }

    private ResolvedAdsTag finalizeLeaf(long indexGroup, long indexOffset,
                                        AdsDataTypeTableEntry dataType,
                                        List<AdsDataTypeArrayInfo> remainingArrayInfo) {
        String name = dataType.getMainName();
        long size = dataType.getSize();
        int stringLen = extractStringLength(name);
        if (stringLen == 0 && (name.equals("STRING") || name.equals("WSTRING"))) {
            stringLen = name.equals("WSTRING") ? (int) (size / 2 - 1) : (int) (size - 1);
        }

        PlcValueType plcType;
        if (!dataType.getArrayInfo().isEmpty()) {
            plcType = PlcValueType.List;
        } else if (!dataType.getChildren().isEmpty()) {
            plcType = PlcValueType.Struct;
        } else {
            plcType = plcValueTypeForName(name, dataType);
        }
        return new ResolvedAdsTag(indexGroup, indexOffset, size, name, plcType,
            stringLen, remainingArrayInfo);
    }

    private ResolvedAdsTag primitiveLeaf(long indexGroup, long indexOffset,
                                         String typeName, long size) {
        int stringLen = extractStringLength(typeName);
        if (stringLen == 0 && (typeName.equals("STRING") || typeName.equals("WSTRING"))) {
            stringLen = typeName.equals("WSTRING") ? (int) (size / 2 - 1) : (int) (size - 1);
        }
        PlcValueType plcType = plcValueTypeForName(typeName, null);
        return new ResolvedAdsTag(indexGroup, indexOffset, size, typeName, plcType,
            stringLen, Collections.emptyList());
    }

    /** Extracts the {@code N} from {@code STRING(N)} / {@code WSTRING(N)}, or {@code 0}. */
    public static int extractStringLength(String typeName) {
        if (typeName == null) return 0;
        if (typeName.startsWith("STRING(") || typeName.startsWith("WSTRING(")) {
            int open = typeName.indexOf('(');
            int close = typeName.indexOf(')');
            if (open > 0 && close > open) {
                try {
                    return Integer.parseInt(typeName.substring(open + 1, close));
                } catch (NumberFormatException ignored) {
                    return 0;
                }
            }
        }
        return 0;
    }

    /** Resolve a primitive type name to its PlcValueType. Returns {@code Struct} for unknowns. */
    public static PlcValueType plcValueTypeForName(String typeName, AdsDataTypeTableEntry maybeEntry) {
        if (typeName == null || typeName.isEmpty()) {
            return PlcValueType.Struct;
        }
        if (typeName.startsWith("STRING")) return PlcValueType.STRING;
        if (typeName.startsWith("WSTRING")) return PlcValueType.WSTRING;
        try {
            return PlcValueType.valueOf(typeName.toUpperCase());
        } catch (IllegalArgumentException ignored) {
            // try via AdsDataType enum
        }
        try {
            AdsDataType ads = AdsDataType.valueOf(typeName.toUpperCase());
            return ads.getPlcValueType();
        } catch (IllegalArgumentException ignored) {
            // fall through
        }
        return PlcValueType.Struct;
    }
}
