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
            root.arrayIndices(), root.child(), selection);
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
     * Restates a direct tag's selection in the terms the decoder reads, so a direct array read
     * returns every element it asked the device for.
     *
     * <p>The size of the request is already multiplied by the element count, but the decoder
     * builds its lists from {@code remainingArrayInfo}: left empty, it read one element and
     * discarded the rest of the response - silently, because a shorter value is still a valid
     * one. This is what {@link #resolve} does for a symbolic tag; a direct tag names
     * its own location, so its selection is the whole of its shape.</p>
     */
    public static ResolvedAdsTag withDirectSelection(ResolvedAdsTag resolved, List<ArrayInfo> selection) {
        if (selection.isEmpty()) {
            return resolved;
        }
        List<AdsDataTypeArrayInfo> dimensions = new ArrayList<>(selection.size());
        for (ArrayInfo dimension : selection) {
            dimensions.add(new AdsDataTypeArrayInfo(
                (long) dimension.getLowerBound(), (long) dimension.getSize()));
        }
        return new ResolvedAdsTag(resolved.indexGroup(), resolved.indexOffset(), resolved.sizeInBytes(),
            resolved.dataTypeName(), PlcValueType.List, resolved.stringLength(), dimensions);
    }

    /**
     * Widens a location resolved for a single element to cover the whole selection: the same
     * start, as many bytes as the selection spans, decoded as a list.
     *
     * <p>The shape follows what the address wrote, dimension by dimension: a range contributes a
     * level of list, a bare index moves the start and collapses. So {@code grid[3,1..3]} is a
     * flat list of three and {@code grid[1..2,0..4]} is two lists of five, and a range spanning
     * one element is still a list of one.</p>
     */
    private static ResolvedAdsTag scaleToSelection(ResolvedAdsTag resolved, List<ArrayInfo> selection) {
        long elements = 1;
        // The decoder builds its lists from the ADS array-info shape, so the selection is
        // restated in those terms: the dimensions written as ranges, each starting where the
        // user asked and spanning as many elements.
        List<AdsDataTypeArrayInfo> dimensions = new ArrayList<>(selection.size());
        for (ArrayInfo dimension : selection) {
            elements *= dimension.getSize();
            if (dimension.isRange()) {
                dimensions.add(new AdsDataTypeArrayInfo(
                    (long) dimension.getLowerBound(), (long) dimension.getSize()));
            }
        }
        if (dimensions.isEmpty()) {
            // Every dimension the selection named was a bare index, so it named one element of
            // them: a scalar, or - where the selection named only the outer dimensions - the
            // whole of what lies inside one, which resolution has already shaped and sized.
            return resolved;
        }
        // A dimension the selection did not name is selected whole, and is still part of the
        // shape: grid[1..2] on an ARRAY [0..9,0..4] is two rows of five, not a flat two. Those
        // dimensions are what resolution left over, and their bytes are already in sizeInBytes.
        dimensions.addAll(resolved.remainingArrayInfo());
        return new ResolvedAdsTag(resolved.indexGroup(), resolved.indexOffset(),
            resolved.sizeInBytes() * elements, resolved.dataTypeName(), PlcValueType.List,
            resolved.stringLength(), dimensions);
    }

    /**
     * Holds a trailing selection to what the device declares and to what one read can express.
     *
     * <p>A read covers one contiguous run of memory. Scanning outwards from the innermost
     * dimension, every dimension inside a dimension selecting more than one element must be
     * selected whole: on an {@code ARRAY [0..9,0..4]}, {@code [0..9,1..3]} names ten separate
     * three-element runs, and the contiguous block of thirty starting at {@code [0,1]} that one
     * read returns is not what was asked for. This refuses it; before, that block was returned.</p>
     *
     * <p>A selection may name fewer dimensions than the array declares; the ones it does not name
     * are selected whole, which is what makes {@code grid[1..2]} two whole rows. Those dimensions
     * are contiguous by construction, so only the named ones are checked here.</p>
     */
    private static void verifySelectionIsOneRead(List<ArrayInfo> selection,
                                                 List<AdsDataTypeArrayInfo> declared,
                                                 String typeName) {
        if (selection.size() > declared.size()) {
            throw new PlcInvalidTagException(String.format(
                "A selection of %d dimension(s) on %s, which declares %d",
                selection.size(), typeName, declared.size()));
        }
        for (int dimension = 0; dimension < selection.size(); dimension++) {
            ArrayInfo selected = selection.get(dimension);
            AdsDataTypeArrayInfo available = declared.get(dimension);
            if (selected.getLowerBound() < available.getLowerBound()
                || selected.getUpperBound() > available.getUpperBound()) {
                throw new PlcInvalidTagException(String.format(
                    "Selection [%d..%d] is outside [%d..%d], which %s declares for dimension %d",
                    selected.getLowerBound(), selected.getUpperBound(),
                    available.getLowerBound(), available.getUpperBound(), typeName, dimension));
            }
            if (dimension == 0 || selected.getSize() == available.getNumElements()) {
                continue;
            }
            for (int outer = 0; outer < dimension; outer++) {
                if (selection.get(outer).getSize() > 1) {
                    throw new PlcInvalidTagException(String.format(
                        "A selection of part of dimension %d of %s, while dimension %d spans %d"
                            + " elements, is not one contiguous read - select the whole of the"
                            + " inner dimension, or one element of the outer one",
                        dimension, typeName, outer, selection.get(outer).getSize()));
                }
            }
        }
    }

    private ResolvedAdsTag resolvePart(long indexGroup, long indexOffset,
                                       AdsDataTypeTableEntry dataType,
                                       List<Integer> arrayIndices,
                                       AddressParser.AddressPart child,
                                       List<ArrayInfo> selection) {
        if (!arrayIndices.isEmpty()) {
            return resolveArray(indexGroup, indexOffset, dataType, arrayIndices, child, selection);
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
            return resolveChild(indexGroup, indexOffset, dataType, child, selection);
        }
        // remainingArrayInfo stays empty for a whole-array read: it means "dimensions still to
        // be applied", and the decoder reads the full shape from the type table itself. This is
        // an internal signal, not the caller-facing report - see SymbolicAdsTag#getArrayInfo.
        return finalizeLeaf(indexGroup, indexOffset, dataType, Collections.emptyList());
    }

    private ResolvedAdsTag resolveArray(long indexGroup, long indexOffset,
                                        AdsDataTypeTableEntry dataType,
                                        List<Integer> arrayIndices,
                                        AddressParser.AddressPart child,
                                        List<ArrayInfo> selection) {
        // The selection's own start indices were appended to the deepest segment of the path, so
        // this is the array it selects from, and its declared dimensions are the ones to hold it
        // to. Deeper segments carry it on; there is nothing to check against here.
        if (child == null && !selection.isEmpty()) {
            verifySelectionIsOneRead(selection,
                dataType.getArrayInfo().subList(
                    Math.max(0, arrayIndices.size() - selection.size()), dataType.getArrayInfo().size()),
                dataType.getMainName());
        }
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
                Collections.emptyList(), child, child == null ? Collections.emptyList() : selection);
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
                                        AddressParser.AddressPart child,
                                        List<ArrayInfo> selection) {
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
            fieldType, child.arrayIndices(), child.child(), selection);
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
