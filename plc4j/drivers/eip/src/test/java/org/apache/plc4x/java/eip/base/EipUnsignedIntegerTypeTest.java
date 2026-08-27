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
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.spi.values.PlcULINT;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Decoding and encoding of the CIP unsigned integer family - USINT, UINT, UDINT and ULINT.
 *
 * <p>Like the bit-string types these are unsigned, and the matching {@code Plc*} types validate
 * against a range starting at zero. Handing them a sign-extended value from the signed
 * {@link ByteBuffer} getters is rejected with a PlcInvalidTagException, which would escape the
 * response handler and fail the whole read request rather than the single tag.
 */
class EipUnsignedIntegerTypeTest {

    // --- scalar reads, including the top of each range ---

    @Test
    void usintDecodesUnsigned() {
        PlcValue value = parse("%N40:USINT", new byte[]{(byte) 0xFF}, CIPDataTypeCode.USINT);

        assertInstanceOf(PlcUSINT.class, value);
        assertEquals(255L, value.getLong());
    }

    @Test
    void uintDecodesUnsigned() {
        PlcValue value = parse("%N40:UINT", bytes(2, b -> b.putShort((short) 0xFFFF)), CIPDataTypeCode.UINT);

        assertInstanceOf(PlcUINT.class, value);
        assertEquals(65535L, value.getLong());
    }

    @Test
    void udintDecodesUnsigned() {
        PlcValue value = parse("%N40:UDINT", bytes(4, b -> b.putInt(0xFFFFFFFF)), CIPDataTypeCode.UDINT);

        assertInstanceOf(PlcUDINT.class, value);
        assertEquals(4294967295L, value.getLong());
    }

    /** The value a signed {@code getInt()} turns negative - and that PlcUDINT then rejects. */
    @Test
    void udintAboveIntegerMaxDoesNotThrow() {
        PlcValue value = parse("%N40:UDINT", bytes(4, b -> b.putInt(0x80000000)), CIPDataTypeCode.UDINT);

        assertNotNull(value);
        assertEquals(2147483648L, value.getLong());
    }

    @Test
    void ulintDecodesUnsigned() {
        PlcValue value = parse("%N40:ULINT", bytes(8, b -> b.putLong(0xFFFFFFFFFFFFFFFFL)), CIPDataTypeCode.ULINT);

        assertInstanceOf(PlcULINT.class, value);
        assertEquals(new BigInteger("18446744073709551615"), value.getBigInteger());
    }

    @Test
    void ulintAboveLongMaxDoesNotThrow() {
        PlcValue value = parse("%N40:ULINT", bytes(8, b -> b.putLong(Long.MIN_VALUE)), CIPDataTypeCode.ULINT);

        assertNotNull(value);
        assertEquals(new BigInteger("9223372036854775808"), value.getBigInteger());
    }

    // --- array reads ---

    @Test
    void arrayOfUdintsIsFullyDecoded() {
        EipTag tag = EipTag.of("%N40[0]:UDINT:3");
        assertEquals(3, tag.getElementNb());

        byte[] raw = bytes(12, b -> {
            b.putInt(0x00000001);
            b.putInt(0xFFFFFFFF);
            b.putInt(0x80000000);
        });
        PlcValue value = EipTcpConnection.parsePlcValue(tag, raw, CIPDataTypeCode.UDINT);

        assertInstanceOf(PlcList.class, value);
        assertEquals(3, value.getLength());
        assertEquals(1L, value.getIndex(0).getLong());
        assertEquals(4294967295L, value.getIndex(1).getLong());
        assertEquals(2147483648L, value.getIndex(2).getLong());
    }

    @Test
    void arrayOfUlintsIsFullyDecoded() {
        byte[] raw = bytes(16, b -> {
            b.putLong(1L);
            b.putLong(0xFFFFFFFFFFFFFFFFL);
        });
        PlcValue value = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0]:ULINT:2"), raw, CIPDataTypeCode.ULINT);

        assertEquals(2, value.getLength());
        assertEquals(BigInteger.ONE, value.getIndex(0).getBigInteger());
        assertEquals(new BigInteger("18446744073709551615"), value.getIndex(1).getBigInteger());
    }

    @Test
    void arrayOfUsintsAndUintsIsFullyDecoded() {
        PlcValue usints = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0]:USINT:3"), new byte[]{0x01, (byte) 0x80, (byte) 0xFF}, CIPDataTypeCode.USINT);
        assertEquals(3, usints.getLength());
        assertEquals(128L, usints.getIndex(1).getLong());
        assertEquals(255L, usints.getIndex(2).getLong());

        PlcValue uints = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0]:UINT:2"),
            bytes(4, b -> {
                b.putShort((short) 0x8000);
                b.putShort((short) 0xFFFF);
            }),
            CIPDataTypeCode.UINT);
        assertEquals(32768L, uints.getIndex(0).getLong());
        assertEquals(65535L, uints.getIndex(1).getLong());
    }

    /** Fixed-size types, so a reply shorter than the declared element count is reported, not thrown. */
    @Test
    void shortReplyIsReportedInsteadOfThrowing() {
        assertNull(EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0]:UDINT:8"), bytes(4, b -> b.putInt(1)), CIPDataTypeCode.UDINT));
        assertNull(EipTcpConnection.parsePlcValue(
            EipTag.of("%N40[0]:ULINT:4"), bytes(8, b -> b.putLong(1)), CIPDataTypeCode.ULINT));
    }

    // --- writes ---

    @Test
    void unsignedIntegersRoundTripThroughEncodeValue() {
        assertRoundTrip("%N40:USINT", CIPDataTypeCode.USINT, new PlcUSINT(255));
        assertRoundTrip("%N40:UINT", CIPDataTypeCode.UINT, new PlcUINT(65535));
        assertRoundTrip("%N40:UDINT", CIPDataTypeCode.UDINT, new PlcUDINT(4294967295L));
        assertRoundTrip("%N40:ULINT", CIPDataTypeCode.ULINT, new PlcULINT(new BigInteger("18446744073709551615")));
    }

    @Test
    void encodeValueProducesOneElementPerType() {
        assertEquals(1, EipTcpConnection.encodeValue(new PlcUSINT(1), CIPDataTypeCode.USINT).length);
        assertEquals(2, EipTcpConnection.encodeValue(new PlcUINT(1), CIPDataTypeCode.UINT).length);
        assertEquals(4, EipTcpConnection.encodeValue(new PlcUDINT(1L), CIPDataTypeCode.UDINT).length);
        assertEquals(8, EipTcpConnection.encodeValue(new PlcULINT(BigInteger.ONE), CIPDataTypeCode.ULINT).length);
    }

    @Test
    void unsignedIntegerCodesMatchTheCipSpecification() {
        assertEquals(0x00C6, CIPDataTypeCode.USINT.getValue());
        assertEquals(0x00C7, CIPDataTypeCode.UINT.getValue());
        assertEquals(0x00C8, CIPDataTypeCode.UDINT.getValue());
        assertEquals(0x00C9, CIPDataTypeCode.ULINT.getValue());

        assertEquals(1, CIPDataTypeCode.USINT.getSize());
        assertEquals(2, CIPDataTypeCode.UINT.getSize());
        assertEquals(4, CIPDataTypeCode.UDINT.getSize());
        assertEquals(8, CIPDataTypeCode.ULINT.getSize());
    }

    /** The signed counterparts must keep decoding to their own types, not to the unsigned ones. */
    @Test
    void signedCounterpartsAreUnaffected() {
        assertEquals(-1L, parse("%N40:DINT", bytes(4, b -> b.putInt(0xFFFFFFFF)), CIPDataTypeCode.DINT).getLong());
        assertEquals(-1L, parse("%N40:INT", bytes(2, b -> b.putShort((short) 0xFFFF)), CIPDataTypeCode.INT).getLong());
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
