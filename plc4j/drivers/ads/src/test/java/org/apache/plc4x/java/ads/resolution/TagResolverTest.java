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
import org.apache.plc4x.java.ads.tag.DirectAdsTag;
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

    /**
     * Omitting the selection asks for the whole array (FR-022). ADS knows its extent from the
     * symbol table, so it reads all five INTs and reports the declared dimension rather than
     * falling back to a single element.
     */
    @Test
    void aBareArrayAddressReadsTheWholeArray() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt", null, List.of()));

        assertEquals(0x200, t.indexOffset(), "the start of the array");
        assertEquals(10, t.sizeInBytes(), "five INTs");
        assertEquals(PlcValueType.List, t.plcValueType(), "decoded as a list, not a scalar");
        // remainingArrayInfo means "dimensions still to apply" and is empty for a whole-array
        // read; the decoder takes the shape from the type table. See wholeArrayRead_noRemainingDims.
        assertTrue(t.remainingArrayInfo().isEmpty());
    }

    /** A bare index is one element, so it reads one and is a scalar - the contrast to the above. */
    @Test
    void aBareIndexReadsOneElement() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[3]", null, List.of()));

        assertEquals(2, t.sizeInBytes(), "one INT");
        assertEquals(PlcValueType.INT, t.plcValueType());
        assertTrue(t.remainingArrayInfo().isEmpty());
    }

    /**
     * A range reads several elements from where it starts. MAIN.g_arrInt is declared 1..5 of INT,
     * so [2..4] begins one element in and covers three of them.
     */
    @Test
    void aRangeReadsFromItsStartForAsManyElementsAsItSpans() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[2..4]", null, List.of()));

        assertEquals(0x200 + 2, t.indexOffset(), "starts at the second element");
        assertEquals(6, t.sizeInBytes(), "three INTs");
        assertEquals(PlcValueType.List, t.plcValueType());
        assertEquals(3, t.remainingArrayInfo().get(0).getNumElements());
    }

    /** A single-element range still yields a list, unlike a bare index. */
    @Test
    void aSingleElementRangeIsStillARange() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[2..2]", null, List.of()));
        assertEquals(2, t.sizeInBytes());
    }

    /**
     * A declared lower bound written in the address is checked against the symbol table, which is
     * authoritative. Agreeing is redundant but harmless; disagreeing means the address was written
     * against a different layout than the PLC has.
     */
    @Test
    void aDeclaredBaseThatMatchesTheSymbolTableIsAccepted() {
        ResolvedAdsTag t = resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[3;1]", null, List.of()));
        assertEquals(0x200 + 4, t.indexOffset(), "resolved exactly as without the base");
    }

    @Test
    void aDeclaredBaseThatContradictsTheSymbolTableIsRejected() {
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_arrInt[3;0]", null, List.of())));
        assertTrue(thrown.getMessage().contains("start at 0"), thrown::getMessage);
        assertTrue(thrown.getMessage().contains("start at 1"), thrown::getMessage);
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

    /**
     * A range before the last dimension is refused while the address is parsed: ADS carries one
     * element count for the whole address, so [1..2,2..3] would ask for two elements of each of
     * two rows in a single count. This is the guarantee the resolver's own contiguity check
     * stands behind, for a tag built without going through the parser.
     */
    @Test
    void multiDim_rangeBeforeTheLastDimensionRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> SymbolicAdsTag.of("MAIN.g_matI16_2x3[1..2,2..3]"));
    }

    /** A bare index collapses, so row 2 in full is a flat list of three. */
    @Test
    void multiDim_oneWholeRow() {
        ResolvedAdsTag t = resolver.resolve(
            new SymbolicAdsTag("MAIN.g_matI16_2x3[2,1..3]", null, List.of()));

        assertEquals(0x300 + 6, t.indexOffset(), "the second row");
        assertEquals(6, t.sizeInBytes());
        assertEquals(1, t.remainingArrayInfo().size(), "the row index collapses");
        assertEquals(3, t.remainingArrayInfo().get(0).getNumElements());
    }

    /** Part of one row is still one run of memory. */
    @Test
    void multiDim_partOfOneRow() {
        ResolvedAdsTag t = resolver.resolve(
            new SymbolicAdsTag("MAIN.g_matI16_2x3[2,2..3]", null, List.of()));

        assertEquals(0x300 + 6 + 2, t.indexOffset(), "row 2, column 2");
        assertEquals(4, t.sizeInBytes(), "two elements");
        assertEquals(1, t.remainingArrayInfo().size());
        assertEquals(2, t.remainingArrayInfo().get(0).getNumElements());
    }

    /**
     * A selection may name only the outer dimensions; the rest are selected whole. The dimensions
     * it did not name stay in the shape - two rows of three, not a flat two, whose bytes would
     * have been transferred and then dropped.
     */
    @Test
    void multiDim_selectionOfRowsKeepsTheColumns() {
        ResolvedAdsTag t = resolver.resolve(
            new SymbolicAdsTag("MAIN.g_matI16_2x3[1..2]", null, List.of()));

        assertEquals(0x300, t.indexOffset());
        assertEquals(12, t.sizeInBytes(), "two whole rows");
        assertEquals(2, t.remainingArrayInfo().size(), "the unnamed dimension is still a dimension");
        assertEquals(2, t.remainingArrayInfo().get(0).getNumElements());
        assertEquals(3, t.remainingArrayInfo().get(1).getNumElements());
    }

    /** Bounds are held per dimension, not only on the first. */
    @Test
    void multiDim_outOfBoundsColumnRejected() {
        assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_matI16_2x3[1,2..4]", null, List.of())));
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

    /**
     * Omitting the brackets asks for the whole array, so a member access after one asks for that
     * member of every element - which is not a single read. Without this the address would
     * resolve silently against the first element and report data for one channel as though it
     * were the whole path.
     */
    @Test
    void aMemberOfAnUnindexedArrayIsRejected() {
        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_plant.channels.setpoints", null, List.of())));
        assertTrue(thrown.getMessage().contains("whole array"), thrown::getMessage);
    }

    /** The whole array itself is still addressable - it is only a member of it that is not. */
    @Test
    void theWholeArrayItselfIsStillAddressable() {
        assertDoesNotThrow(
            () -> resolver.resolve(new SymbolicAdsTag("MAIN.g_plant.channels", null, List.of())));
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

    /**
     * A direct array tag asks the device for every element it selected, so it has to decode every
     * one of them. The request size was already multiplied by the count while the decoder was
     * given no shape at all, so the extra elements were read from the wire and dropped - a short
     * value that looks like a valid one.
     */
    @Test
    void aDirectSelectionBecomesTheDecodedShape() {
        ResolvedAdsTag scalar = new ResolvedAdsTag(0x4020, 0, 16, "DINT", PlcValueType.DINT,
            0, List.of());

        ResolvedAdsTag shaped = TagResolver.withDirectSelection(scalar,
            DirectAdsTag.of("0x4020/0[0..3]:DINT").getArrayInfo());

        assertEquals(PlcValueType.List, shaped.plcValueType(), "decoded as a list");
        assertEquals(1, shaped.remainingArrayInfo().size(), "one dimension");
        assertEquals(4, shaped.remainingArrayInfo().get(0).getNumElements(), "all four elements");
        assertEquals(16, shaped.sizeInBytes(), "the request size is unchanged");
    }

    @Test
    void aDirectScalarKeepsItsScalarShape() {
        ResolvedAdsTag scalar = new ResolvedAdsTag(0x4020, 0, 4, "DINT", PlcValueType.DINT,
            0, List.of());

        ResolvedAdsTag shaped = TagResolver.withDirectSelection(scalar,
            DirectAdsTag.of("0x4020/0[3]:DINT").getArrayInfo());

        assertEquals(PlcValueType.DINT, shaped.plcValueType());
        assertTrue(shaped.remainingArrayInfo().isEmpty());
    }
}
