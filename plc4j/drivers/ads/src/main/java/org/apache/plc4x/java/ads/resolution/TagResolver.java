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
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.ArrayList;
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

    private final boolean symbolTablesLoaded;

    public TagResolver(Map<String, AdsSymbolTableEntry> symbolTable,
                       Map<String, AdsDataTypeTableEntry> dataTypeTable) {
        this(symbolTable, dataTypeTable, true);
    }

    /**
     * @param symbolTablesLoaded whether the symbol- and data-type tables were loaded at all. When
     *                           they weren't (the connection option
     *                           {@code load-symbol-and-data-type-tables} is disabled), symbolic
     *                           addresses simply cannot be resolved and saying so is more useful
     *                           than reporting every symbol as unknown - see GH-1626.
     */
    public TagResolver(Map<String, AdsSymbolTableEntry> symbolTable,
                       Map<String, AdsDataTypeTableEntry> dataTypeTable,
                       boolean symbolTablesLoaded) {
        this.symbolTable = symbolTable;
        this.dataTypeTable = dataTypeTable;
        this.symbolTablesLoaded = symbolTablesLoaded;
    }

    public ResolvedAdsTag resolve(SymbolicAdsTag tag) {
        if (!symbolTablesLoaded) {
            throw new PlcInvalidTagException(
                "Cannot resolve the symbolic address '" + tag.getSymbolicAddress()
                    + "': the symbol and data-type tables were not loaded because the connection"
                    + " option 'load-symbol-and-data-type-tables' is disabled. Either enable that"
                    + " option or address the value directly as"
                    + " '{IndexGroup}/{IndexOffset}:{TYPE}'.");
        }
        // The trailing selection is not part of the symbolic path; it says which elements of the
        // resolved location to read. Its first index is appended to the path's own indices so the
        // existing bounds checking and lower-bound arithmetic apply to it unchanged.
        String path = ArrayNotationParser.addressPart(tag.getSymbolicAddress());
        List<ArrayInfo> selection = tag.getSelection();
        AddressParser.AddressPart root = withSelectionStart(AddressParser.parse(path), selection);

        AdsSymbolTableEntry symbol = symbolTable.get(root.baseSegment());
        if (symbol == null) {
            throw new PlcInvalidTagException("Unknown symbol: " + root.baseSegment());
        }
        AdsDataTypeTableEntry dataType = dataTypeTable.get(symbol.getDataTypeName());
        if (dataType == null) {
            throw new PlcInvalidTagException(
                "Unknown data type for symbol " + root.baseSegment() + ": " + symbol.getDataTypeName());
        }
        verifyDeclaredBase(tag, symbol, dataType);

        ResolvedAdsTag resolved = resolvePart(symbol.getGroup(), symbol.getOffset(), dataType,
            root.arrayIndices(), root.child());
        return scaleToSelection(resolved, selection);
    }

    /**
     * Appends the first index of each selected dimension to the deepest segment of the path, so
     * that the location the read starts at is resolved by the same code that resolves an index
     * written in the path itself.
     */
    private static AddressParser.AddressPart withSelectionStart(AddressParser.AddressPart part,
                                                                List<ArrayInfo> selection) {
        if (selection.isEmpty()) {
            return part;
        }
        if (part.child() != null) {
            return new AddressParser.AddressPart(part.baseSegment(), part.arrayIndices(),
                withSelectionStart(part.child(), selection));
        }
        List<Integer> indices = new ArrayList<>(part.arrayIndices());
        for (ArrayInfo dimension : selection) {
            indices.add(dimension.getLowerBound());
        }
        return new AddressParser.AddressPart(part.baseSegment(), indices, null);
    }

    /**
     * Checks a declared lower bound written in the address against the one the symbol table
     * declares. The table is authoritative; a base in the address is the user's statement of
     * intent, and a disagreement means the address was written against a different layout than
     * the PLC has - which would otherwise read silently shifted data.
     *
     * <p>This is the one rule of the notation that cannot be checked while the address is
     * parsed, because the symbol table is not loaded then.
     */
    private void verifyDeclaredBase(SymbolicAdsTag tag, AdsSymbolTableEntry symbol,
                                    AdsDataTypeTableEntry dataType) {
        Integer declared = tag.getDeclaredBase();
        if (declared == null || dataType.getArrayInfo().isEmpty()) {
            return;
        }
        long actual = dataType.getArrayInfo().get(dataType.getArrayInfo().size() - 1).getLowerBound();
        if (declared != actual) {
            throw new PlcInvalidTagException(String.format(
                "Address '%s' declares the array to start at %d, but %s declares it to start at %d",
                tag.getSymbolicAddress(), declared, symbol.getName(), actual));
        }
    }

    /**
     * Widens a location resolved for a single element to cover the whole selection: the same
     * start, as many bytes as the selection spans, decoded as a list.
     */
    private static ResolvedAdsTag scaleToSelection(ResolvedAdsTag resolved, List<ArrayInfo> selection) {
        int elements = 1;
        for (ArrayInfo dimension : selection) {
            elements *= dimension.getSize();
        }
        if (elements <= 1) {
            return resolved;
        }
        // The decoder builds its lists from the ADS array-info shape, so the selection is
        // restated in those terms: one dimension, starting where the user asked, as many
        // elements as it spans.
        List<AdsDataTypeArrayInfo> dimensions = new ArrayList<>(selection.size());
        for (ArrayInfo dimension : selection) {
            dimensions.add(new AdsDataTypeArrayInfo(
                (long) dimension.getLowerBound(), (long) dimension.getSize()));
        }
        return new ResolvedAdsTag(resolved.indexGroup(), resolved.indexOffset(),
            resolved.sizeInBytes() * elements, resolved.dataTypeName(), PlcValueType.List,
            resolved.stringLength(), dimensions);
    }

    private ResolvedAdsTag resolvePart(long indexGroup, long indexOffset,
                                       AdsDataTypeTableEntry dataType,
                                       List<Integer> arrayIndices,
                                       AddressParser.AddressPart child) {
        if (!arrayIndices.isEmpty()) {
            return resolveArray(indexGroup, indexOffset, dataType, arrayIndices, child);
        }
        if (child != null) {
            if (!dataType.getArrayInfo().isEmpty()) {
                // Omitting the brackets asks for the whole array, so a member access after one
                // asks for that member of every element - which is not one contiguous read. The
                // partially indexed form is refused in resolveArray for the same reason; this is
                // the same rule where no index was given at all. Without it the address would
                // silently resolve against the first element.
                throw new PlcInvalidTagException(
                    "Field access requires an array element to be specified for "
                        + dataType.getMainName() + ": an address that omits the index asks for the"
                        + " whole array, and a member of every element is not a single read");
            }
            return resolveChild(indexGroup, indexOffset, dataType, child);
        }
        // remainingArrayInfo stays empty for a whole-array read: it means "dimensions still to
        // be applied", and the decoder reads the full shape from the type table itself. This is
        // an internal signal, not the caller-facing report - see SymbolicAdsTag#getArrayInfo.
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
