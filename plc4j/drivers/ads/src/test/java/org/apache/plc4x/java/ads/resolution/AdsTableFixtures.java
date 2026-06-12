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
package org.apache.plc4x.java.ads.resolution;

import org.apache.plc4x.java.ads.readwrite.AdsDataTypeArrayInfo;
import org.apache.plc4x.java.ads.readwrite.AdsDataTypeTableEntry;
import org.apache.plc4x.java.ads.readwrite.AdsDatatypeId;
import org.apache.plc4x.java.ads.readwrite.AdsSymbolTableEntry;

import java.util.Collections;
import java.util.List;

/**
 * Tiny builders for the generated table-entry types so resolver/decoder tests don't have to
 * spell out 30+ ctor args every time. Only the fields used by the resolver/decoder are
 * surfaced — everything else gets reasonable zero/false defaults.
 */
final class AdsTableFixtures {

    private AdsTableFixtures() {
    }

    /** Primitive scalar (BOOL, DINT, REAL, …) — no array, no children. */
    static AdsDataTypeTableEntry scalar(String mainName, long size) {
        return entry(mainName, "", size, 0, Collections.emptyList(), Collections.emptyList());
    }

    /** STRING(n) leaf — same shape as a scalar, just with the bracketed length in the name. */
    static AdsDataTypeTableEntry stringType(int charCount) {
        return entry("STRING(" + charCount + ")", "", charCount + 1, 0,
            Collections.emptyList(), Collections.emptyList());
    }

    static AdsDataTypeTableEntry wstringType(int charCount) {
        return entry("WSTRING(" + charCount + ")", "", (charCount + 1L) * 2, 0,
            Collections.emptyList(), Collections.emptyList());
    }

    /** Array type — secondary name is the element type, size is total bytes. */
    static AdsDataTypeTableEntry array(String mainName, String elementTypeName,
                                       long totalSize, List<AdsDataTypeArrayInfo> dims) {
        return entry(mainName, elementTypeName, totalSize, dims.size(), dims, Collections.emptyList());
    }

    /** Struct type — children list (each {@link #field}) defines the layout. */
    static AdsDataTypeTableEntry struct(String mainName, long size, List<AdsDataTypeTableEntry> children) {
        return entry(mainName, "", size, 0, Collections.emptyList(), children);
    }

    /** A struct field: mainName is the field name, secondaryName is the field type's name. */
    static AdsDataTypeTableEntry field(String fieldName, String fieldTypeName, long offset, long size) {
        return entryAt(fieldName, fieldTypeName, size, offset, 0,
            Collections.emptyList(), Collections.emptyList());
    }

    static AdsDataTypeArrayInfo dim(long lowerBound, long numElements) {
        return new AdsDataTypeArrayInfo(lowerBound, numElements);
    }

    static AdsSymbolTableEntry symbol(String name, String dataTypeName, long group, long offset, long size) {
        return new AdsSymbolTableEntry(0L, group, offset, size, 0L,
            false, false, false, false, false, false, false, false, false, false, false, false,
            false, false, false, false, false,
            name, dataTypeName, "", 0L, new byte[0], null, new byte[0]);
    }

    private static AdsDataTypeTableEntry entry(String mainName, String secondaryName,
                                               long size, int arrayDims,
                                               List<AdsDataTypeArrayInfo> arrayInfo,
                                               List<AdsDataTypeTableEntry> children) {
        return entryAt(mainName, secondaryName, size, 0L, arrayDims, arrayInfo, children);
    }

    private static AdsDataTypeTableEntry entryAt(String mainName, String secondaryName,
                                                 long size, long offset, int arrayDims,
                                                 List<AdsDataTypeArrayInfo> arrayInfo,
                                                 List<AdsDataTypeTableEntry> children) {
        return new AdsDataTypeTableEntry(0L, 1L, 0L, 0L, size, offset,
            AdsDatatypeId.ADST_VOID,
            false, false, false, false, false, false, false, false, false, false, false, false,
            false, false, false, false, false, false, false, false, false, false, false,
            arrayDims, children.size(),
            mainName, secondaryName, "",
            arrayInfo, children, new byte[0],
            null, null, null, new byte[0]);
    }
}
