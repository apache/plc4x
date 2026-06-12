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
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcStruct;

import java.util.List;
import java.util.Map;

/**
 * Mirror of {@link ValueDecoder}: walks a {@link PlcValue} guided by a
 * {@link ResolvedAdsTag} (and the data type table) and emits bytes in the layout the PLC
 * expects.
 */
public final class ValueEncoder {

    private final Map<String, AdsDataTypeTableEntry> dataTypeTable;

    public ValueEncoder(Map<String, AdsDataTypeTableEntry> dataTypeTable) {
        this.dataTypeTable = dataTypeTable;
    }

    public void encode(WriteBuffer wb, ResolvedAdsTag tag, PlcValue value) throws BufferException {
        // Partial-array write: caller supplied a nested PlcList for remaining dims.
        if (!tag.remainingArrayInfo().isEmpty()) {
            AdsDataTypeTableEntry dataType = dataTypeTable.get(tag.dataTypeName());
            if (dataType == null) {
                throw new BufferException("Unknown array data type: " + tag.dataTypeName());
            }
            if (!(value instanceof PlcList list)) {
                throw new PlcRuntimeException("Expected PlcList for array tag " + tag.dataTypeName());
            }
            encodePartialArray(wb, dataType, tag.remainingArrayInfo(), list);
            return;
        }
        if (tag.plcValueType() == PlcValueType.Struct) {
            AdsDataTypeTableEntry dataType = dataTypeTable.get(tag.dataTypeName());
            if (dataType == null) {
                throw new BufferException("Unknown struct data type: " + tag.dataTypeName());
            }
            if (!(value instanceof PlcStruct s)) {
                throw new PlcRuntimeException("Expected PlcStruct for struct tag " + tag.dataTypeName());
            }
            encodeStruct(wb, dataType, s);
            return;
        }
        if (tag.plcValueType() == PlcValueType.List) {
            AdsDataTypeTableEntry dataType = dataTypeTable.get(tag.dataTypeName());
            if (dataType == null) {
                throw new BufferException("Unknown array data type: " + tag.dataTypeName());
            }
            if (!(value instanceof PlcList list)) {
                throw new PlcRuntimeException("Expected PlcList for array tag " + tag.dataTypeName());
            }
            encodePartialArray(wb, dataType, dataType.getArrayInfo(), list);
            return;
        }
        DataItem.staticSerialize(wb, value, tag.plcValueType(), tag.stringLength());
    }

    private void encodePartialArray(WriteBuffer wb,
                                    AdsDataTypeTableEntry arrayDataType,
                                    List<AdsDataTypeArrayInfo> dims,
                                    PlcList list) throws BufferException {
        if (dims.isEmpty()) {
            encodeArrayElement(wb, arrayDataType, list);
            return;
        }
        AdsDataTypeArrayInfo cur = dims.get(0);
        long expected = cur.getNumElements();
        if (list.getList().size() != expected) {
            throw new PlcRuntimeException("PlcList size mismatch for array dimension: expected "
                + expected + ", got " + list.getList().size());
        }
        List<AdsDataTypeArrayInfo> rest = dims.subList(1, dims.size());
        for (PlcValue v : list.getList()) {
            if (rest.isEmpty()) {
                if (v instanceof PlcList inner) {
                    encodeArrayElement(wb, arrayDataType, inner);
                } else {
                    encodeLeafElement(wb, arrayDataType, v);
                }
            } else {
                if (!(v instanceof PlcList inner)) {
                    throw new PlcRuntimeException("Expected nested PlcList for multi-dim array");
                }
                encodePartialArray(wb, arrayDataType, rest, inner);
            }
        }
    }

    /**
     * For the inner-most array level when the element is itself a PlcList (struct etc. is handled
     * via {@link #encodeLeafElement}). This shouldn't actually be reached in practice — kept as a
     * safety net.
     */
    private void encodeArrayElement(WriteBuffer wb, AdsDataTypeTableEntry arrayDataType,
                                    PlcList listValue) throws BufferException {
        for (PlcValue v : listValue.getList()) {
            encodeLeafElement(wb, arrayDataType, v);
        }
    }

    /** Encode a single element of {@code arrayDataType}; element type from secondary name. */
    private void encodeLeafElement(WriteBuffer wb, AdsDataTypeTableEntry arrayDataType,
                                   PlcValue value) throws BufferException {
        String elementTypeName = arrayDataType.getSecondaryName();
        AdsDataTypeTableEntry elementType = dataTypeTable.get(elementTypeName);
        if (elementType != null) {
            if (!elementType.getChildren().isEmpty()) {
                if (!(value instanceof PlcStruct s)) {
                    throw new PlcRuntimeException("Expected PlcStruct for struct array element");
                }
                encodeStruct(wb, elementType, s);
                return;
            }
            int stringLen = stringLengthForLeaf(elementType);
            PlcValueType plcType = TagResolver.plcValueTypeForName(elementType.getMainName(), elementType);
            DataItem.staticSerialize(wb, value, plcType, stringLen);
            return;
        }
        int stringLen = TagResolver.extractStringLength(elementTypeName);
        if (stringLen == 0 && (elementTypeName.equals("STRING") || elementTypeName.equals("WSTRING"))) {
            long totalElements = 1;
            for (AdsDataTypeArrayInfo dim : arrayDataType.getArrayInfo()) {
                totalElements *= dim.getNumElements();
            }
            long perElement = arrayDataType.getSize() / totalElements;
            stringLen = elementTypeName.equals("WSTRING") ? (int) (perElement / 2 - 1) : (int) (perElement - 1);
        }
        PlcValueType plcType = TagResolver.plcValueTypeForName(elementTypeName, null);
        DataItem.staticSerialize(wb, value, plcType, stringLen);
    }

    private void encodeStruct(WriteBuffer wb, AdsDataTypeTableEntry dataType, PlcStruct s) throws BufferException {
        int startBytePos = wb.getPositionInBits() / 8;
        int totalSize = (int) dataType.getSize();
        int curPos = 0;
        for (AdsDataTypeTableEntry child : dataType.getChildren()) {
            if (child.getOffset() == 0xFFFFFFFFL) continue;
            if (child.getOffset() < curPos) continue;
            if (child.getOffset() + child.getSize() > totalSize) continue;

            if (child.getOffset() > curPos) {
                long padBytes = child.getOffset() - curPos;
                writeZeroBytes(wb, (int) padBytes);
                curPos = (wb.getPositionInBits() / 8) - startBytePos;
            }

            PlcValue v = s.getValue(child.getMainName());
            if (v == null) {
                // Case-insensitive fallback (TC2/TC3 casing differences and user-supplied keys).
                for (String key : s.getKeys()) {
                    if (key.equalsIgnoreCase(child.getMainName())) {
                        v = s.getValue(key);
                        break;
                    }
                }
            }
            if (v == null) {
                throw new PlcRuntimeException("Missing struct field: " + child.getMainName());
            }

            String childTypeName = child.getSecondaryName();
            AdsDataTypeTableEntry childType = dataTypeTable.get(childTypeName);
            if (childType == null) {
                int stringLen = TagResolver.extractStringLength(childTypeName);
                if (stringLen == 0) {
                    if (childTypeName.equals("STRING")) stringLen = (int) (child.getSize() - 1);
                    else if (childTypeName.equals("WSTRING")) stringLen = (int) (child.getSize() / 2 - 1);
                }
                PlcValueType plcType = TagResolver.plcValueTypeForName(childTypeName, null);
                DataItem.staticSerialize(wb, v, plcType, stringLen);
            } else if (!childType.getArrayInfo().isEmpty()) {
                if (!(v instanceof PlcList list)) {
                    throw new PlcRuntimeException("Expected PlcList for array field " + child.getMainName());
                }
                encodePartialArray(wb, childType, childType.getArrayInfo(), list);
            } else if (!childType.getChildren().isEmpty()) {
                if (!(v instanceof PlcStruct ns)) {
                    throw new PlcRuntimeException("Expected PlcStruct for struct field " + child.getMainName());
                }
                encodeStruct(wb, childType, ns);
            } else {
                int stringLen = stringLengthForLeaf(childType);
                PlcValueType plcType = TagResolver.plcValueTypeForName(childType.getMainName(), childType);
                DataItem.staticSerialize(wb, v, plcType, stringLen);
            }
            curPos = (wb.getPositionInBits() / 8) - startBytePos;
        }
        if (curPos < totalSize) {
            writeZeroBytes(wb, totalSize - curPos);
        }
    }

    private static void writeZeroBytes(WriteBuffer wb, int n) throws BufferException {
        if (n <= 0) return;
        wb.writeBits(n * 8, new byte[n]);
    }

    private static int stringLengthForLeaf(AdsDataTypeTableEntry leaf) {
        int n = TagResolver.extractStringLength(leaf.getMainName());
        if (n != 0) return n;
        n = TagResolver.extractStringLength(leaf.getSecondaryName());
        if (n != 0) return n;
        if (leaf.getMainName().equals("STRING")) return (int) (leaf.getSize() - 1);
        if (leaf.getMainName().equals("WSTRING")) return (int) (leaf.getSize() / 2 - 1);
        return 0;
    }
}
