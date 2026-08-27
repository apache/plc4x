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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.apache.plc4x.java.spi.values.PlcBYTE;
import org.apache.plc4x.java.spi.values.PlcDWORD;
import org.apache.plc4x.java.spi.values.PlcLWORD;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcWORD;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.EnumSet;
import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Decoding and encoding of the CIP bit-string family - BYTE, WORD, DWORD and LWORD.
 *
 * <p>These are unsigned N-byte bit strings, so the signed {@link ByteBuffer} getters must not be
 * handed to the {@code Plc*} constructors directly: those validate against a range starting at
 * zero, and a sign-extended value is rejected with a PlcInvalidTagException. That exception would
 * escape the response handler and fail the whole read request rather than the single tag.
 */
class EipBitStringTypeTest {

    // --- scalar reads, including the all-bits-set value that sign extension breaks ---

    @Test
    void byteDecodesUnsigned() {
        PlcValue value = parse("%N40:BYTE", new byte[]{(byte) 0xFF}, CIPDataTypeCode.BYTE);

        assertInstanceOf(PlcBYTE.class, value);
        assertEquals(255L, value.getLong());
    }

    @Test
    void wordDecodesUnsigned() {
        PlcValue value = parse("%N40:WORD", bytes(2, b -> b.putShort((short) 0xFFFF)), CIPDataTypeCode.WORD);

        assertInstanceOf(PlcWORD.class, value);
        assertEquals(65535L, value.getLong());
    }

    @Test
    void dwordDecodesUnsigned() {
        PlcValue value = parse("%N40:DWORD", bytes(4, b -> b.putInt(0xFFFFFFFF)), CIPDataTypeCode.DWORD);

        assertInstanceOf(PlcDWORD.class, value);
        assertEquals(4294967295L, value.getLong());
    }

    /**
     * Bit 31 set is entirely ordinary for a 32 bit boolean array, and is the case a signed
     * {@code getInt()} turns into a negative number.
     */
    @Test
    void dwordWithHighBitSetDoesNotThrow() {
        PlcValue value = parse("%N40:DWORD", bytes(4, b -> b.putInt(0x80000000)), CIPDataTypeCode.DWORD);

        assertNotNull(value);
        assertEquals(2147483648L, value.getLong());
    }

    @Test
    void lwordDecodesUnsigned() {
        PlcValue value = parse("%N40:LWORD", bytes(8, b -> b.putLong(0xFFFFFFFFFFFFFFFFL)), CIPDataTypeCode.LWORD);

        assertInstanceOf(PlcLWORD.class, value);
        assertEquals(new BigInteger("18446744073709551615"), value.getBigInteger());
    }

    @Test
    void lwordWithHighBitSetDoesNotThrow() {
        PlcValue value = parse("%N40:LWORD", bytes(8, b -> b.putLong(Long.MIN_VALUE)), CIPDataTypeCode.LWORD);

        assertNotNull(value);
        assertEquals(new BigInteger("9223372036854775808"), value.getBigInteger());
    }

    // --- array reads ---

    @Test
    void arrayOfDwordsIsFullyDecoded() {
        EipTag tag = EipTag.of("%N40[0..2]:DWORD");
        assertEquals(3, tag.getElementNb());

        byte[] raw = bytes(12, b -> {
            b.putInt(0x00000001);
            b.putInt(0xFFFFFFFF);
            b.putInt(0x80000000);
        });
        PlcValue value = EipTcpConnection.parsePlcValue(tag, raw, CIPDataTypeCode.DWORD);

        assertInstanceOf(PlcList.class, value);
        assertEquals(3, value.getLength());
        assertEquals(1L, value.getIndex(0).getLong());
        assertEquals(4294967295L, value.getIndex(1).getLong());
        assertEquals(2147483648L, value.getIndex(2).getLong());
    }

    @Test
    void arrayOfLwordsIsFullyDecoded() {
        EipTag tag = EipTag.of("%N40[0..1]:LWORD");

        byte[] raw = bytes(16, b -> {
            b.putLong(1L);
            b.putLong(0xFFFFFFFFFFFFFFFFL);
        });
        PlcValue value = EipTcpConnection.parsePlcValue(tag, raw, CIPDataTypeCode.LWORD);

        assertEquals(2, value.getLength());
        assertEquals(BigInteger.ONE, value.getIndex(0).getBigInteger());
        assertEquals(new BigInteger("18446744073709551615"), value.getIndex(1).getBigInteger());
    }

    @Test
    void arrayOfBytesAndWordsIsFullyDecoded() {
        PlcValue bytesValue = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0..2]:BYTE"), new byte[]{0x01, (byte) 0x80, (byte) 0xFF}, CIPDataTypeCode.BYTE);
        assertEquals(3, bytesValue.getLength());
        assertEquals(128L, bytesValue.getIndex(1).getLong());
        assertEquals(255L, bytesValue.getIndex(2).getLong());

        PlcValue wordsValue = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0..1]:WORD"),
            bytes(4, b -> {
                b.putShort((short) 0x8000);
                b.putShort((short) 0xFFFF);
            }),
            CIPDataTypeCode.WORD);
        assertEquals(32768L, wordsValue.getIndex(0).getLong());
        assertEquals(65535L, wordsValue.getIndex(1).getLong());
    }

    /**
     * A reply carrying fewer elements than the tag declares is reported as undecodable rather
     * than thrown out of the response handler - the bit-string types are fixed size, so they
     * take part in the same guard as DINT and friends.
     */
    @Test
    void shortReplyIsReportedInsteadOfThrowing() {
        assertNull(EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0..7]:DWORD"), bytes(4, b -> b.putInt(1)), CIPDataTypeCode.DWORD));
        assertNull(EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0..3]:LWORD"), bytes(8, b -> b.putLong(1)), CIPDataTypeCode.LWORD));
    }

    // --- writes ---

    @Test
    void bitStringsRoundTripThroughEncodeValue() {
        assertRoundTrip("%N40:BYTE", CIPDataTypeCode.BYTE, new PlcBYTE(255));
        assertRoundTrip("%N40:WORD", CIPDataTypeCode.WORD, new PlcWORD(65535));
        assertRoundTrip("%N40:DWORD", CIPDataTypeCode.DWORD, new PlcDWORD(4294967295L));
        assertRoundTrip("%N40:LWORD", CIPDataTypeCode.LWORD, new PlcLWORD(new BigInteger("18446744073709551615")));
    }

    @Test
    void encodeValueProducesOneElementPerType() {
        assertEquals(1, EipTcpConnection.encodeValue(new PlcBYTE(1), CIPDataTypeCode.BYTE).length);
        assertEquals(2, EipTcpConnection.encodeValue(new PlcWORD(1), CIPDataTypeCode.WORD).length);
        assertEquals(4, EipTcpConnection.encodeValue(new PlcDWORD(1L), CIPDataTypeCode.DWORD).length);
        assertEquals(8, EipTcpConnection.encodeValue(new PlcLWORD(BigInteger.ONE), CIPDataTypeCode.LWORD).length);
    }

    // --- invariants that hold for every entry in the codec table ---

    /**
     * Every fixed-size type must decode a payload and encode it back to the identical bytes, in
     * exactly {@code getSize()} bytes. Driving this off the table itself means a type added later
     * is covered without anyone remembering to extend this test.
     */
    @Test
    void everyFixedSizeTypeRoundTripsItsOwnBytes() {
        for (CIPDataTypeCode type : EipTcpConnection.FIXED_SIZE_CODECS.keySet()) {
            if (type == CIPDataTypeCode.BOOL) {
                continue; // BOOL normalises any non-zero byte to 1 - see boolEncodesToZeroOrOne
            }
            for (byte filler : new byte[]{0x00, 0x7F, (byte) 0xFF}) {
                if (filler == (byte) 0xFF && (type == CIPDataTypeCode.REAL || type == CIPDataTypeCode.LREAL)) {
                    continue; // all bits set is NaN for the floats, which is not a value to compare
                }
                byte[] raw = new byte[type.getSize()];
                java.util.Arrays.fill(raw, filler);

                PlcValue decoded = EipTcpConnection.parsePlcValue(EipTag.of("%N40:" + type.name()), raw, type);
                assertNotNull(decoded, () -> type + " with filler " + filler + " did not decode");

                byte[] encoded = EipTcpConnection.encodeValue(decoded, type);
                assertEquals(type.getSize(), encoded.length, () -> type + " payload size");
                assertArrayEquals(raw, encoded, () -> type + " with filler " + filler + " did not round trip");
            }
        }
    }

    /** BOOL is the one type whose encoding is not the identity, so it is checked on its own. */
    @Test
    void boolEncodesToZeroOrOne() {
        assertArrayEquals(new byte[]{1}, EipTcpConnection.encodeValue(new org.apache.plc4x.java.spi.values.PlcBOOL(true), CIPDataTypeCode.BOOL));
        assertArrayEquals(new byte[]{0}, EipTcpConnection.encodeValue(new org.apache.plc4x.java.spi.values.PlcBOOL(false), CIPDataTypeCode.BOOL));
    }

    // --- the mspec defect behind all of this ---

    /**
     * LWORD used to share 0x00D3 with DWORD and STRINGI used to share 0x00DD with ENGUNIT, which
     * silently dropped one of each pair from the generated lookup tables.
     */
    @Test
    void dataTypeCodesAreUnique() {
        Map<Integer, CIPDataTypeCode> seen = new HashMap<>();
        for (CIPDataTypeCode code : EnumSet.allOf(CIPDataTypeCode.class)) {
            CIPDataTypeCode previous = seen.put(code.getValue(), code);
            assertNull(previous, () -> String.format("%s and %s share value 0x%04X", previous, code, code.getValue()));
        }
    }

    @Test
    void bitStringCodesMatchTheCipSpecification() {
        assertEquals(0x00D1, CIPDataTypeCode.BYTE.getValue());
        assertEquals(0x00D2, CIPDataTypeCode.WORD.getValue());
        assertEquals(0x00D3, CIPDataTypeCode.DWORD.getValue());
        assertEquals(0x00D4, CIPDataTypeCode.LWORD.getValue());

        assertEquals(1, CIPDataTypeCode.BYTE.getSize());
        assertEquals(2, CIPDataTypeCode.WORD.getSize());
        assertEquals(4, CIPDataTypeCode.DWORD.getSize());
        assertEquals(8, CIPDataTypeCode.LWORD.getSize());
    }

    // --- helpers ---

    private static PlcValue parse(String address, byte[] raw, CIPDataTypeCode type) {
        EipTag tag = EipTag.of(address);
        assertNotNull(tag, address);
        return EipTcpConnection.parsePlcValue(tag, raw, type);
    }

    private static void assertRoundTrip(String address, CIPDataTypeCode type, PlcValue value) {
        byte[] encoded = EipTcpConnection.encodeValue(value, type);
        assertEquals(type.getSize(), encoded.length, type + " payload size");
        assertEquals(value, parse(address, encoded, type), type + " round trip");
    }

    private static byte[] bytes(int size, java.util.function.Consumer<ByteBuffer> filler) {
        ByteBuffer buffer = ByteBuffer.allocate(size).order(ByteOrder.LITTLE_ENDIAN);
        filler.accept(buffer);
        return buffer.array();
    }
}
