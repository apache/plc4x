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

import org.apache.plc4x.java.ads.readwrite.AdsDataTypeTableEntry;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.values.*;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import static org.apache.plc4x.java.ads.resolution.AdsTableFixtures.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests the read decoder ↔ write encoder pair: encode a {@link PlcValue} into bytes via
 * {@link ValueEncoder}, then decode the same bytes back via {@link ValueDecoder} and check
 * that we land back where we started. Uses small, hand-built data type tables.
 */
class ValueCodecTest {

    private final Map<String, AdsDataTypeTableEntry> types = new HashMap<>();
    private final ValueDecoder decoder;
    private final ValueEncoder encoder;

    private ValueCodecTest() {
        types.put("BOOL", scalar("BOOL", 1));
        types.put("SINT", scalar("SINT", 1));
        types.put("INT", scalar("INT", 2));
        types.put("DINT", scalar("DINT", 4));
        types.put("REAL", scalar("REAL", 4));
        types.put("STRING(10)", stringType(10));
        decoder = new ValueDecoder(types);
        encoder = new ValueEncoder(types);
    }

    /** Build a {@link ResolvedAdsTag} that points at the resolved leaf for round-tripping. */
    private static ResolvedAdsTag tag(String dataTypeName, long size, PlcValueType pvt, int stringLen) {
        return new ResolvedAdsTag(0, 0, size, dataTypeName, pvt, stringLen, Collections.emptyList());
    }

    private static ResolvedAdsTag tagPartial(String dataTypeName, long size,
                                             org.apache.plc4x.java.ads.readwrite.AdsDataTypeArrayInfo... dims) {
        return new ResolvedAdsTag(0, 0, size, dataTypeName, PlcValueType.List, 0, List.of(dims));
    }

    private byte[] encode(ResolvedAdsTag t, PlcValue v) throws Exception {
        WriteBufferByteBased wb = new WriteBufferByteBased(new byte[(int) t.sizeInBytes()],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
        encoder.encode(wb, t, v);
        return wb.getBytes();
    }

    private PlcValue decode(ResolvedAdsTag t, byte[] data) throws Exception {
        ReadBufferByteBased rb = new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"));
        return decoder.decode(rb, t);
    }

    @Test
    void scalarBool_roundtrip() throws Exception {
        ResolvedAdsTag t = tag("BOOL", 1, PlcValueType.BOOL, 0);
        assertEquals(new PlcBOOL(true), decode(t, encode(t, new PlcBOOL(true))));
        assertEquals(new PlcBOOL(false), decode(t, encode(t, new PlcBOOL(false))));
    }

    @Test
    void scalarDint_negativeRoundtrip() throws Exception {
        ResolvedAdsTag t = tag("DINT", 4, PlcValueType.DINT, 0);
        PlcValue v = new PlcDINT(-12345678);
        byte[] bytes = encode(t, v);
        assertEquals(4, bytes.length);
        assertEquals(v, decode(t, bytes));
    }

    @Test
    void scalarReal_roundtrip() throws Exception {
        ResolvedAdsTag t = tag("REAL", 4, PlcValueType.REAL, 0);
        PlcValue v = new PlcREAL(3.14159f);
        assertEquals(v, decode(t, encode(t, v)));
    }

    @Test
    void scalarString_roundtrip() throws Exception {
        // Total size = 11 bytes (10 chars + NUL).
        ResolvedAdsTag t = tag("STRING(10)", 11, PlcValueType.STRING, 10);
        PlcValue v = new PlcSTRING("Hello");
        byte[] bytes = encode(t, v);
        assertEquals(11, bytes.length);
        // Decoder strips trailing NULs and gives us back the original string.
        assertEquals(v, decode(t, bytes));
    }

    @Test
    void oneDimArrayOfInt_roundtrip() throws Exception {
        // ARRAY [1..3] OF INT — 6 bytes total, 1 dim of 3 elements lower-bound 1.
        types.put("ARR_INT3",
            array("ARR_INT3", "INT", 6, List.of(dim(1, 3))));
        ResolvedAdsTag t = tag("ARR_INT3", 6, PlcValueType.List, 0);
        PlcValue v = new PlcList(List.of(new PlcINT(-3), new PlcINT(0), new PlcINT(7)));
        assertEquals(v, decode(t, encode(t, v)));
    }

    @Test
    void twoDimMatrix_partialRowRead_returnsPlcList() throws Exception {
        // ARRAY [1..2, 1..3] OF INT — 12 bytes; partial read of one row (3 elements, 6 bytes).
        types.put("MAT_INT_2x3",
            array("MAT_INT_2x3", "INT", 12, List.of(dim(1, 2), dim(1, 3))));
        // Partial-dim resolved tag: only the inner dim remains (3 INTs).
        ResolvedAdsTag t = tagPartial("MAT_INT_2x3", 6, dim(1, 3));
        PlcValue row = new PlcList(List.of(new PlcINT(10), new PlcINT(11), new PlcINT(12)));
        assertEquals(row, decode(t, encode(t, row)));
    }

    @Test
    void wholeMatrix_returnsNestedPlcList() throws Exception {
        types.put("MAT_INT_2x3",
            array("MAT_INT_2x3", "INT", 12, List.of(dim(1, 2), dim(1, 3))));
        ResolvedAdsTag t = tag("MAT_INT_2x3", 12, PlcValueType.List, 0);
        PlcValue full = new PlcList(List.of(
            new PlcList(List.of(new PlcINT(10), new PlcINT(11), new PlcINT(12))),
            new PlcList(List.of(new PlcINT(-10), new PlcINT(-11), new PlcINT(-12)))
        ));
        assertEquals(full, decode(t, encode(t, full)));
    }

    @Test
    void struct_roundtrip_withPaddingAndStringField() throws Exception {
        // TPair { a: SINT @0, str: STRING(5) @1 }; size = 7 bytes (1 + 6).
        types.put("STRING(5)", stringType(5));
        AdsDataTypeTableEntry pair = struct("TPair", 7, List.of(
            field("a", "SINT", 0, 1),
            field("str", "STRING(5)", 1, 6)
        ));
        types.put("TPair", pair);
        ResolvedAdsTag t = tag("TPair", 7, PlcValueType.Struct, 0);
        Map<String, PlcValue> props = new LinkedHashMap<>();
        props.put("a", new PlcSINT(-7));
        props.put("str", new PlcSTRING("Hi"));
        PlcStruct s = new PlcStruct(props);
        PlcValue back = decode(t, encode(t, s));
        assertInstanceOf(PlcStruct.class, back);
        PlcStruct decoded = (PlcStruct) back;
        assertEquals(new PlcSINT(-7), decoded.getValue("a"));
        assertEquals(new PlcSTRING("Hi"), decoded.getValue("str"));
    }

    @Test
    void struct_skipsTrailingPaddingOnDecode() throws Exception {
        // TWithPad { a: SINT @0 } stored in 4 bytes (3 trailing padding) — decoder must
        // still consume the full 4 bytes so subsequent reads stay aligned.
        AdsDataTypeTableEntry t = struct("TWithPad", 4, List.of(field("a", "SINT", 0, 1)));
        types.put("TWithPad", t);
        ResolvedAdsTag tag = tag("TWithPad", 4, PlcValueType.Struct, 0);
        byte[] bytes = encode(tag, new PlcStruct(Map.of("a", new PlcSINT(5))));
        assertEquals(4, bytes.length);
        PlcStruct decoded = (PlcStruct) decode(tag, bytes);
        assertEquals(new PlcSINT(5), decoded.getValue("a"));
    }

    @Test
    void encode_missingStructField_isReportedClearly() {
        AdsDataTypeTableEntry t = struct("TPair", 2, List.of(
            field("a", "SINT", 0, 1),
            field("b", "SINT", 1, 1)
        ));
        types.put("TPair", t);
        ResolvedAdsTag tag = tag("TPair", 2, PlcValueType.Struct, 0);
        Map<String, PlcValue> props = new LinkedHashMap<>();
        props.put("a", new PlcSINT(1)); // missing "b"
        Exception e = assertThrows(Exception.class,
            () -> encode(tag, new PlcStruct(props)));
        assertTrue(e.getMessage().contains("b"), "expected mention of missing field 'b': " + e.getMessage());
    }

    @Test
    void encode_caseInsensitiveStructFieldFallback() throws Exception {
        // Struct field is "value"; user supplies "Value" — encoder should still resolve it.
        AdsDataTypeTableEntry t = struct("TWrap", 2, List.of(field("value", "INT", 0, 2)));
        types.put("TWrap", t);
        ResolvedAdsTag tag = tag("TWrap", 2, PlcValueType.Struct, 0);
        Map<String, PlcValue> props = new LinkedHashMap<>();
        props.put("Value", new PlcINT(42));
        PlcStruct decoded = (PlcStruct) decode(tag, encode(tag, new PlcStruct(props)));
        assertEquals(new PlcINT(42), decoded.getValue("value"));
    }

    @Test
    void encode_arraySizeMismatch_isReportedClearly() {
        types.put("ARR_INT3",
            array("ARR_INT3", "INT", 6, List.of(dim(1, 3))));
        ResolvedAdsTag tag = tag("ARR_INT3", 6, PlcValueType.List, 0);
        // Only 2 elements supplied for a 3-element array.
        PlcValue v = new PlcList(List.of(new PlcINT(1), new PlcINT(2)));
        Exception e = assertThrows(Exception.class, () -> encode(tag, v));
        assertTrue(e.getMessage().toLowerCase().contains("size"),
            "expected size-mismatch message: " + e.getMessage());
    }

    @Test
    void resolvedAdsTag_record_accessors() {
        ResolvedAdsTag t = new ResolvedAdsTag(1, 2, 3, "DINT", PlcValueType.DINT, 0, List.of());
        assertEquals(1, t.indexGroup());
        assertEquals(2, t.indexOffset());
        assertEquals(3, t.sizeInBytes());
        assertEquals("DINT", t.dataTypeName());
        assertEquals(PlcValueType.DINT, t.plcValueType());
        assertEquals(0, t.stringLength());
        assertTrue(t.remainingArrayInfo().isEmpty());
    }
}
