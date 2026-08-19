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
import org.apache.plc4x.java.ads.readwrite.AdsSymbolTableEntry;
import org.apache.plc4x.java.ads.tag.SymbolicAdsTag;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.apache.plc4x.java.ads.resolution.AdsTableFixtures.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * The resolver tests build a small, hand-crafted symbol/data-type table that mirrors a
 * subset of the TC3 PLC layout used by ManualFactoryAdsDriverTestTC3:
 * <ul>
 *   <li>{@code MAIN.g_s32}  — INT (DINT, 4 bytes)</li>
 *   <li>{@code MAIN.g_arrInt} — {@code ARRAY [1..5] OF INT}</li>
 *   <li>{@code MAIN.g_matI16_2x3} — {@code ARRAY [1..2, 1..3] OF INT}</li>
 *   <li>{@code MAIN.g_simple} — struct with {@code s8: SINT, str: STRING(13)}</li>
 *   <li>{@code MAIN.g_plant} — struct containing an array of structs</li>
 * </ul>
 */
class TagResolverTest {

    private final Map<String, AdsSymbolTableEntry> symbols = new HashMap<>();
    private final Map<String, AdsDataTypeTableEntry> types = new HashMap<>();
    private final TagResolver resolver = new TagResolver(symbols, types);

    private TagResolverTest() {
        types.put("BOOL", scalar("BOOL", 1));
        types.put("SINT", scalar("SINT", 1));
        types.put("INT", scalar("INT", 2));
        types.put("DINT", scalar("DINT", 4));
        types.put("REAL", scalar("REAL", 4));
        types.put("STRING(13)", stringType(13));

        // 1-D INT array, 1-based, 5 elements → 10 bytes.
        types.put("ARRAY [1..5] OF INT",
            array("ARRAY [1..5] OF INT", "INT", 10, List.of(dim(1, 5))));

        // 2-D INT matrix, 1-based, 2x3 → 12 bytes.
        types.put("ARRAY [1..2, 1..3] OF INT",
            array("ARRAY [1..2, 1..3] OF INT", "INT", 12, List.of(dim(1, 2), dim(1, 3))));

        // Simple struct: SINT @0, STRING(13) @1 (size 14 chars + 1 NUL = 14), total 15.
        AdsDataTypeTableEntry simpleStruct = struct("TSimpleStruct", 15, List.of(
            field("s8", "SINT", 0, 1),
            field("str", "STRING(13)", 1, 14)
        ));
        types.put("TSimpleStruct", simpleStruct);

        // Setpoints array (REAL[1..4], 16 bytes); channel struct: id @0 (DINT), setpoints @4.
        types.put("ARRAY [1..4] OF REAL",
            array("ARRAY [1..4] OF REAL", "REAL", 16, List.of(dim(1, 4))));
        AdsDataTypeTableEntry channel = struct("TChannel", 20, List.of(
            field("id", "DINT", 0, 4),
            field("setpoints", "ARRAY [1..4] OF REAL", 4, 16)
        ));
        types.put("TChannel", channel);
        // Two channels in a 1-based array (40 bytes total).
        types.put("ARRAY [1..2] OF TChannel",
            array("ARRAY [1..2] OF TChannel", "TChannel", 40, List.of(dim(1, 2))));
        // Plant struct contains the channels array at offset 0.
        AdsDataTypeTableEntry plant = struct("TPlant", 40, List.of(
            field("channels", "ARRAY [1..2] OF TChannel", 0, 40)
        ));
        types.put("TPlant", plant);

        // Symbols.
        symbols.put("MAIN.g_s32", symbol("MAIN.g_s32", "DINT", 0x4040, 0x100, 4));
        symbols.put("MAIN.g_arrInt", symbol("MAIN.g_arrInt", "ARRAY [1..5] OF INT", 0x4040, 0x200, 10));
        symbols.put("MAIN.g_matI16_2x3", symbol("MAIN.g_matI16_2x3", "ARRAY [1..2, 1..3] OF INT", 0x4040, 0x300, 12));
        symbols.put("MAIN.g_simple", symbol("MAIN.g_simple", "TSimpleStruct", 0x4040, 0x400, 15));
        symbols.put("MAIN.g_plant", symbol("MAIN.g_plant", "TPlant", 0x4040, 0x500, 40));
    }

    @Test
    void scalarLeaf() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_s32", null, List.of()));
        assertEquals(0x4040, t.indexGroup());
        assertEquals(0x100, t.indexOffset());
        assertEquals(4, t.sizeInBytes());
        assertEquals(PlcValueType.DINT, t.plcValueType());
        assertEquals("DINT", t.dataTypeName());
        assertTrue(t.remainingArrayInfo().isEmpty());
    }

    @Test
    void singleIndexedElement_offsetUses1BasedIndex() {
        // arr[3] in 1..5 → element 2 (zero-based) at offset 0x200 + 2*2 = 0x204
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[3]", null, List.of()));
        assertEquals(0x200 + 4, t.indexOffset());
        assertEquals(2, t.sizeInBytes());
        assertEquals(PlcValueType.INT, t.plcValueType());
    }

    @Test
    void firstAndLastElementBoundsAreInclusive() {
        ResolvedAdsTag first = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[1]", null, List.of()));
        assertEquals(0x200, first.indexOffset());
        ResolvedAdsTag last = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[5]", null, List.of()));
        assertEquals(0x200 + 8, last.indexOffset());
    }

    @Test
    void indexOutOfBoundsRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[0]", null, List.of())));
        assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[6]", null, List.of())));
    }

    @Test
    void multiDim_partialIndex_returnsRowSlice() {
        // matrix[1] selects row 1 of the 2x3 INT matrix → 6 bytes, 1 remaining dim (size 3).
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_matI16_2x3[1]", null, List.of()));
        assertEquals(0x300, t.indexOffset());
        assertEquals(6, t.sizeInBytes());
        assertEquals(PlcValueType.List, t.plcValueType());
        assertEquals(1, t.remainingArrayInfo().size());
        AdsDataTypeArrayInfo dim = t.remainingArrayInfo().get(0);
        assertEquals(1, dim.getLowerBound());
        assertEquals(3, dim.getNumElements());
    }

    @Test
    void multiDim_fullIndex_returnsScalar() {
        // matrix[2][3] is the last element: row 1 (0-based) * 6 bytes + col 2 * 2 = 6 + 4 = 10.
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_matI16_2x3[2][3]", null, List.of()));
        assertEquals(0x300 + 10, t.indexOffset());
        assertEquals(2, t.sizeInBytes());
        assertEquals(PlcValueType.INT, t.plcValueType());
    }

    @Test
    void tooManyIndicesRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[1][2]", null, List.of())));
    }

    @Test
    void structField_offsetAddedToParent() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_simple.s8", null, List.of()));
        assertEquals(0x400, t.indexOffset());
        assertEquals(1, t.sizeInBytes());
        assertEquals(PlcValueType.SINT, t.plcValueType());
    }

    @Test
    void structFieldString_carriesStringLength() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_simple.str", null, List.of()));
        assertEquals(0x400 + 1, t.indexOffset());
        assertEquals(14, t.sizeInBytes());
        assertEquals(PlcValueType.STRING, t.plcValueType());
        assertEquals(13, t.stringLength());
    }

    @Test
    void unknownFieldRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_simple.nope", null, List.of())));
    }

    @Test
    void unknownSymbolRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.unknown", null, List.of())));
    }

    @Test
    void deepMixedPath_resolvesElementInsideArrayOfStructs() {
        // plant.channels[2].setpoints[3]:
        //   plant @ 0x500, channels offset 0, channel size 20, channel[2] -> +20
        //   setpoints offset 4 inside channel, REAL[1..4], setpoints[3] -> element 2 (0-based) * 4 = 8
        // Final: 0x500 + 0 + 20 + 4 + 8 = 0x520
        ResolvedAdsTag t = resolver.resolve(
            new SymbolicAdsTag("MAIN.g_plant.channels[2].setpoints[3]", null, List.of()));
        assertEquals(0x500 + 20 + 4 + 8, t.indexOffset());
        assertEquals(4, t.sizeInBytes());
        assertEquals(PlcValueType.REAL, t.plcValueType());
    }

    @Test
    void wholeStructRead() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_simple", null, List.of()));
        assertEquals(0x400, t.indexOffset());
        assertEquals(15, t.sizeInBytes());
        assertEquals(PlcValueType.Struct, t.plcValueType());
    }

    @Test
    void wholeArrayRead_noRemainingDims() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt", null, List.of()));
        assertEquals(0x200, t.indexOffset());
        assertEquals(10, t.sizeInBytes());
        assertEquals(PlcValueType.List, t.plcValueType());
        // remainingArrayInfo is empty here — the decoder will use the type table's full arrayInfo.
        assertTrue(t.remainingArrayInfo().isEmpty());
    }

    @Test
    void plcValueTypeForName_handlesStandardTypes() {
        assertEquals(PlcValueType.STRING, TagResolver.plcValueTypeForName("STRING(80)", null));
        assertEquals(PlcValueType.WSTRING, TagResolver.plcValueTypeForName("WSTRING(40)", null));
        assertEquals(PlcValueType.DINT, TagResolver.plcValueTypeForName("DINT", null));
        assertEquals(PlcValueType.Struct, TagResolver.plcValueTypeForName("UnknownType", null));
        assertEquals(PlcValueType.Struct, TagResolver.plcValueTypeForName(null, null));
    }

    /**
     * With 'load-symbol-and-data-type-tables' disabled the tables are never fetched, so symbolic
     * addresses cannot be resolved at all. Saying that is far more useful than reporting every
     * symbol as unknown, which is what a user saw before - see GH-1626.
     */
    @Test
    void symbolicAddressIsRejectedWhenTablesWereNotLoaded() {
        TagResolver withoutTables = new TagResolver(symbols, types, false);

        PlcInvalidTagException exception = assertThrows(PlcInvalidTagException.class,
            () -> withoutTables.resolve(new SymbolicAdsTag("MAIN.g_s32", null, List.of())));

        assertTrue(exception.getMessage().contains("MAIN.g_s32"), exception.getMessage());
        assertTrue(exception.getMessage().contains("load-symbol-and-data-type-tables"),
            "the message has to name the option that disabled the tables: " + exception.getMessage());
    }

    /**
     * A genuinely unknown symbol still reports as unknown when the tables *are* loaded - the new
     * message must not swallow that case.
     */
    @Test
    void unknownSymbolStillReportsAsUnknownWhenTablesAreLoaded() {
        PlcInvalidTagException exception = assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.doesNotExist", null, List.of())));

        assertTrue(exception.getMessage().contains("Unknown symbol"), exception.getMessage());
    }

    @Test
    void extractStringLength_parsesParens() {
        assertEquals(80, TagResolver.extractStringLength("STRING(80)"));
        assertEquals(40, TagResolver.extractStringLength("WSTRING(40)"));
        assertEquals(0, TagResolver.extractStringLength("DINT"));
        assertEquals(0, TagResolver.extractStringLength(null));
        assertEquals(0, TagResolver.extractStringLength("STRING(abc)"));
    }
}
