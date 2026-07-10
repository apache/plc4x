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

import org.apache.plc4x.java.ads.readwrite.AdsDataTypeArrayInfo;
import org.apache.plc4x.java.ads.readwrite.AdsDataTypeTableEntry;
import org.apache.plc4x.java.ads.readwrite.DataItem;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcStruct;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Decodes ADS read responses into {@link PlcValue}s, driven by the shape baked into a
 * {@link ResolvedAdsTag} plus the data type table.
 *
 * <p>Read shape rules:
 * <ul>
 *   <li>{@code remainingArrayInfo} non-empty → produce nested {@link PlcList}, then recurse
 *   into the element type from the data type table.</li>
 *   <li>Leaf is a struct (children non-empty) → walk children, honoring offsets.</li>
 *   <li>Otherwise → scalar via {@code DataItem.staticParse}.</li>
 * </ul>
 */
public final class ValueDecoder {

    /**
     * Upper bound for the initial capacity of a list whose element count comes from the
     * (server-supplied) data type table. Prevents a bogus, attacker-controlled array dimension
     * from triggering an eager multi-GB allocation before any element is actually read. The list
     * still grows past this if that many elements really are present on the wire.
     */
    private static final int MAX_INITIAL_CAPACITY = 1024;

    private final Map<String, AdsDataTypeTableEntry> dataTypeTable;

    public ValueDecoder(Map<String, AdsDataTypeTableEntry> dataTypeTable) {
        this.dataTypeTable = dataTypeTable;
    }

    public PlcValue decode(ReadBuffer rb, ResolvedAdsTag tag) throws BufferException {
        // Partial-array result: produce nested PlcLists for the remaining dims, recurse into
        // the element type for the inner-most level.
        if (!tag.remainingArrayInfo().isEmpty()) {
            AdsDataTypeTableEntry dataType = dataTypeTable.get(tag.dataTypeName());
            if (dataType == null) {
                throw new BufferException("Unknown array data type: " + tag.dataTypeName());
            }
            return decodePartialArray(rb, dataType, tag.remainingArrayInfo());
        }

        // Full-leaf decode based on plcValueType.
        if (tag.plcValueType() == PlcValueType.Struct) {
            AdsDataTypeTableEntry dataType = dataTypeTable.get(tag.dataTypeName());
            if (dataType == null) {
                throw new BufferException("Unknown struct data type: " + tag.dataTypeName());
            }
            return decodeStruct(rb, dataType);
        }
        if (tag.plcValueType() == PlcValueType.List) {
            // Whole array — type table has the dimensions.
            AdsDataTypeTableEntry dataType = dataTypeTable.get(tag.dataTypeName());
            if (dataType == null) {
                throw new BufferException("Unknown array data type: " + tag.dataTypeName());
            }
            return decodePartialArray(rb, dataType, dataType.getArrayInfo());
        }
        // Scalar.
        return DataItem.staticParse(rb, tag.plcValueType(), tag.stringLength());
    }

    /**
     * Decode {@code dims.size()} levels of nested PlcList. The element type is taken from
     * {@code arrayDataType.getSecondaryName()}.
     */
    private PlcValue decodePartialArray(ReadBuffer rb,
                                        AdsDataTypeTableEntry arrayDataType,
                                        List<AdsDataTypeArrayInfo> dims) throws BufferException {
        if (dims.isEmpty()) {
            return decodeArrayElement(rb, arrayDataType);
        }
        AdsDataTypeArrayInfo cur = dims.get(0);
        List<AdsDataTypeArrayInfo> rest = dims.subList(1, dims.size());
        long count = cur.getNumElements();
        // numElements is a wire-supplied unsigned 32-bit value (up to 0xFFFFFFFF). Reject counts
        // that don't fit a positive int before the (int) cast - otherwise a large value either
        // overflows to a negative capacity (NegativeArraySizeException) or forces an eager
        // multi-GB allocation (OOM / DoS). The list itself is not pre-sized to the untrusted
        // count: it grows as elements are actually decoded, and dataReader.read() throws once the
        // buffer is exhausted.
        if (count > Integer.MAX_VALUE) {
            throw new BufferException("Array count of " + count + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
        }
        List<PlcValue> elements = new ArrayList<>((int) Math.min(count, MAX_INITIAL_CAPACITY));
        for (long i = 0; i < count; i++) {
            elements.add(decodePartialArray(rb, arrayDataType, rest));
        }
        return new PlcList(elements);
    }

    /**
     * Decode a single array element. Element type comes from the array's secondary name.
     */
    private PlcValue decodeArrayElement(ReadBuffer rb,
                                        AdsDataTypeTableEntry arrayDataType) throws BufferException {
        String elementTypeName = arrayDataType.getSecondaryName();
        AdsDataTypeTableEntry elementType = dataTypeTable.get(elementTypeName);
        if (elementType != null) {
            // Struct element → recurse; otherwise primitive via DataItem.
            if (!elementType.getChildren().isEmpty()) {
                return decodeStruct(rb, elementType);
            }
            // Primitive element described by a data type table entry (e.g. STRING(80))
            int stringLen = stringLengthForLeaf(elementType);
            PlcValueType plcType = TagResolver.plcValueTypeForName(elementType.getMainName(), elementType);
            return DataItem.staticParse(rb, plcType, stringLen);
        }
        // Primitive element NOT in the table — take type from secondary name (e.g. INT, BOOL,
        // STRING(80)). This is the common TC3 case for arrays of primitives.
        int stringLen = TagResolver.extractStringLength(elementTypeName);
        if (stringLen == 0 && (elementTypeName.equals("STRING") || elementTypeName.equals("WSTRING"))) {
            // Derive from total array bytes — only possible if we know element count.
            long totalElements = 1;
            for (AdsDataTypeArrayInfo dim : arrayDataType.getArrayInfo()) {
                totalElements *= dim.getNumElements();
            }
            long perElement = arrayDataType.getSize() / totalElements;
            stringLen = elementTypeName.equals("WSTRING") ? (int) (perElement / 2 - 1) : (int) (perElement - 1);
        }
        PlcValueType plcType = TagResolver.plcValueTypeForName(elementTypeName, null);
        return DataItem.staticParse(rb, plcType, stringLen);
    }

    private PlcValue decodeStruct(ReadBuffer rb, AdsDataTypeTableEntry dataType) throws BufferException {
        Map<String, PlcValue> properties = new LinkedHashMap<>();
        int startBytePos = rb.getPositionInBits() / 8;
        int totalSize = (int) dataType.getSize();
        int curPos = 0;

        for (AdsDataTypeTableEntry child : dataType.getChildren()) {
            // Skip TwinCAT compiler-generated / not-memory-mapped sentinel entries.
            if (child.getOffset() == 0xFFFFFFFFL) {
                continue;
            }
            // Skip overlapping (re-using same offset) and out-of-range children.
            if (child.getOffset() < curPos) {
                continue;
            }
            if (child.getOffset() + child.getSize() > totalSize) {
                continue;
            }
            // Pad to child offset.
            if (child.getOffset() > curPos) {
                long padBytes = child.getOffset() - curPos;
                rb.readBits((int) padBytes * 8);
                curPos = (rb.getPositionInBits() / 8) - startBytePos;
            }

            String childTypeName = child.getSecondaryName();
            AdsDataTypeTableEntry childType = dataTypeTable.get(childTypeName);
            PlcValue value;
            if (childType == null) {
                // Primitive field whose type isn't in the table.
                int stringLen = TagResolver.extractStringLength(childTypeName);
                if (stringLen == 0) {
                    if (childTypeName.equals("STRING")) stringLen = (int) (child.getSize() - 1);
                    else if (childTypeName.equals("WSTRING")) stringLen = (int) (child.getSize() / 2 - 1);
                }
                PlcValueType plcType = TagResolver.plcValueTypeForName(childTypeName, null);
                value = DataItem.staticParse(rb, plcType, stringLen);
            } else if (!childType.getArrayInfo().isEmpty()) {
                // Field is an array — fully read it.
                value = decodePartialArray(rb, childType, childType.getArrayInfo());
            } else if (!childType.getChildren().isEmpty()) {
                // Nested struct.
                value = decodeStruct(rb, childType);
            } else {
                // Primitive field with a data type table entry (e.g. STRING(80)).
                int stringLen = stringLengthForLeaf(childType);
                PlcValueType plcType = TagResolver.plcValueTypeForName(childType.getMainName(), childType);
                value = DataItem.staticParse(rb, plcType, stringLen);
            }
            properties.put(child.getMainName(), value);
            curPos = (rb.getPositionInBits() / 8) - startBytePos;
        }

        // Skip trailing padding to keep the cursor aligned to the end of the struct.
        if (curPos < totalSize) {
            rb.readBits((totalSize - curPos) * 8);
        }
        return new PlcStruct(properties);
    }

    private static int stringLengthForLeaf(AdsDataTypeTableEntry leaf) {
        int n = TagResolver.extractStringLength(leaf.getMainName());
        if (n != 0) return n;
        n = TagResolver.extractStringLength(leaf.getSecondaryName());
        if (n != 0) return n;
        // Last resort: derive from size.
        if (leaf.getMainName().equals("STRING")) return (int) (leaf.getSize() - 1);
        if (leaf.getMainName().equals("WSTRING")) return (int) (leaf.getSize() / 2 - 1);
        return 0;
    }

    /** Suppress unused-import warning in environments where WithOption isn't needed. */
    @SuppressWarnings("unused")
    private static final WithOption[] UNUSED_KEEP_IMPORT = new WithOption[0];
}
