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

package org.apache.plc4x.java.spi.utils;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.buffers.api.Message;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Nested;
import org.mockito.Mockito;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.when;

class StaticHelperTest {

    @Nested
    class ArraySizeInBytesTests {

        @Test
        void testArraySizeInBytes_withMessageList() {
            Message msg1 = Mockito.mock(Message.class);
            Message msg2 = Mockito.mock(Message.class);
            when(msg1.getLengthInBytes()).thenReturn(10);
            when(msg2.getLengthInBytes()).thenReturn(20);

            List<Message> messages = Arrays.asList(msg1, msg2);
            assertEquals(30, StaticHelper.ARRAY_SIZE_IN_BYTES(messages));
        }

        @Test
        void testArraySizeInBytes_withMessageArray() {
            Message msg1 = Mockito.mock(Message.class);
            Message msg2 = Mockito.mock(Message.class);
            when(msg1.getLengthInBytes()).thenReturn(15);
            when(msg2.getLengthInBytes()).thenReturn(25);

            Message[] messages = new Message[]{msg1, msg2};
            assertEquals(40, StaticHelper.ARRAY_SIZE_IN_BYTES(messages));
        }

        @Test
        void testArraySizeInBytes_withEmptyList() {
            assertEquals(0, StaticHelper.ARRAY_SIZE_IN_BYTES(new ArrayList<>()));
        }

        @Test
        void testArraySizeInBytes_withEmptyArray() {
            assertEquals(0, StaticHelper.ARRAY_SIZE_IN_BYTES(new Message[0]));
        }

        @Test
        void testArraySizeInBytes_throwsExceptionForNonMessageElements() {
            List<String> invalidList = Arrays.asList("not", "messages");
            assertThrows(RuntimeException.class, () -> StaticHelper.ARRAY_SIZE_IN_BYTES(invalidList));
        }

        @Test
        void testArraySizeInBytes_throwsExceptionForUnsupportedType() {
            assertThrows(RuntimeException.class, () -> StaticHelper.ARRAY_SIZE_IN_BYTES("unsupported"));
        }
    }

    @Nested
    class CountTests {

        @Test
        void testCount_withNull() {
            assertEquals(0, StaticHelper.COUNT(null));
        }

        @Test
        void testCount_withBooleanArray() {
            assertEquals(3, StaticHelper.COUNT(new boolean[]{true, false, true}));
        }

        @Test
        void testCount_withByteArray() {
            assertEquals(4, StaticHelper.COUNT(new byte[]{1, 2, 3, 4}));
        }

        @Test
        void testCount_withShortArray() {
            assertEquals(2, StaticHelper.COUNT(new short[]{100, 200}));
        }

        @Test
        void testCount_withIntArray() {
            assertEquals(5, StaticHelper.COUNT(new int[]{1, 2, 3, 4, 5}));
        }

        @Test
        void testCount_withLongArray() {
            assertEquals(3, StaticHelper.COUNT(new long[]{1L, 2L, 3L}));
        }

        @Test
        void testCount_withFloatArray() {
            assertEquals(2, StaticHelper.COUNT(new float[]{1.0f, 2.0f}));
        }

        @Test
        void testCount_withDoubleArray() {
            assertEquals(4, StaticHelper.COUNT(new double[]{1.0, 2.0, 3.0, 4.0}));
        }

        @Test
        void testCount_withObjectArray() {
            assertEquals(3, StaticHelper.COUNT(new String[]{"a", "b", "c"}));
        }

        @Test
        void testCount_withCollection() {
            assertEquals(4, StaticHelper.COUNT(Arrays.asList("a", "b", "c", "d")));
        }

        @Test
        void testCount_withEmptyArray() {
            assertEquals(0, StaticHelper.COUNT(new int[0]));
        }

        @Test
        void testCount_withEmptyCollection() {
            assertEquals(0, StaticHelper.COUNT(new ArrayList<>()));
        }

        @Test
        void testCount_throwsExceptionForUnsupportedType() {
            assertThrows(PlcRuntimeException.class, () -> StaticHelper.COUNT("unsupported"));
        }
    }

    @Nested
    class StrLenTests {

        @Test
        void testStrLen_withNull() {
            assertEquals(0, StaticHelper.STR_LEN(null));
        }

        @Test
        void testStrLen_withString() {
            assertEquals(5, StaticHelper.STR_LEN("hello"));
        }

        @Test
        void testStrLen_withEmptyString() {
            assertEquals(0, StaticHelper.STR_LEN(""));
        }

        @Test
        void testStrLen_withInteger() {
            assertEquals(3, StaticHelper.STR_LEN(123));
        }
    }

    @Nested
    class CastTests {

        @Test
        void testCast_successfulCast() {
            Object obj = "test";
            String result = StaticHelper.CAST(obj, String.class);
            assertEquals("test", result);
        }

        @Test
        void testCast_throwsExceptionOnInvalidCast() {
            Object obj = "test";
            assertThrows(PlcRuntimeException.class, () -> StaticHelper.CAST(obj, Integer.class));
        }

        @Test
        void testCast_withNull() {
            assertNull(StaticHelper.CAST(null, String.class));
        }

        @Test
        void testCast_withSubclass() {
            Object obj = new ArrayList<String>();
            List<?> result = StaticHelper.CAST(obj, List.class);
            assertNotNull(result);
        }
    }

    @Nested
    class CeilTests {

        @Test
        void testCeil_withPositiveValue() {
            assertEquals(4, StaticHelper.CEIL(3.1));
            assertEquals(4, StaticHelper.CEIL(3.9));
        }

        @Test
        void testCeil_withNegativeValue() {
            assertEquals(-3, StaticHelper.CEIL(-3.1));
            assertEquals(-3, StaticHelper.CEIL(-3.9));
        }

        @Test
        void testCeil_withWholeNumber() {
            assertEquals(5, StaticHelper.CEIL(5.0));
        }

        @Test
        void testCeil_withZero() {
            assertEquals(0, StaticHelper.CEIL(0.0));
        }
    }

    @Nested
    class PadCountTests {

        @Test
        void testPadCount_whenHasNextIsTrue() {
            int[] array = new int[]{1, 2, 3};
            assertEquals(3, StaticHelper.PADCOUNT(array, true));
        }

        @Test
        void testPadCount_whenHasNextIsFalse() {
            int[] array = new int[]{1, 2, 3};
            assertEquals(0, StaticHelper.PADCOUNT(array, false));
        }

        @Test
        void testPadCount_withEmptyArrayAndTrue() {
            assertEquals(0, StaticHelper.PADCOUNT(new int[0], true));
        }

        @Test
        void testPadCount_withNullAndTrue() {
            assertEquals(0, StaticHelper.PADCOUNT(null, true));
        }
    }

    @Nested
    class VarDUintLengthTests {

        @Test
        void testGetVarDUintLengthInBits_withZero() {
            assertEquals(8, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(0L));
        }

        @Test
        void testGetVarDUintLengthInBits_withSmallValue() {
            assertEquals(8, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(127L));
        }

        @Test
        void testGetVarDUintLengthInBits_with2ByteValue() {
            assertEquals(16, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(128L));
            assertEquals(16, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(16383L));
        }

        @Test
        void testGetVarDUintLengthInBits_with3ByteValue() {
            assertEquals(24, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(16384L));
        }

        @Test
        void testGetVarDUintLengthInBits_withBigInteger() {
            assertEquals(8, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(BigInteger.ZERO));
            assertEquals(8, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(BigInteger.valueOf(100)));
            assertEquals(16, StaticHelper.GET_VARDUINT_LENGTH_IN_BITS(BigInteger.valueOf(200)));
        }
    }

    @Nested
    class VarDIntLengthTests {

        @Test
        void testGetVarDIntLengthInBits_withZero() {
            assertEquals(8, StaticHelper.GET_VARDINT_LENGTH_IN_BITS(0L));
        }

        @Test
        void testGetVarDIntLengthInBits_withPositiveValue() {
            assertEquals(8, StaticHelper.GET_VARDINT_LENGTH_IN_BITS(63L));
        }

        @Test
        void testGetVarDIntLengthInBits_withNegativeValue() {
            assertEquals(8, StaticHelper.GET_VARDINT_LENGTH_IN_BITS(-64L));
        }

        @Test
        void testGetVarDIntLengthInBits_with2ByteValue() {
            assertEquals(16, StaticHelper.GET_VARDINT_LENGTH_IN_BITS(64L));
            assertEquals(16, StaticHelper.GET_VARDINT_LENGTH_IN_BITS(-65L));
        }
    }

    @Nested
    class EncodeHexTests {

        @Test
        void testEncodeHex_withNullOrEmpty() {
            assertEquals("", StaticHelper.ENCODE_HEX(null));
            assertEquals("", StaticHelper.ENCODE_HEX(new byte[0]));
        }

        @Test
        void testEncodeHex_withSingleByte() {
            assertEquals("00", StaticHelper.ENCODE_HEX(new byte[]{0}));
            assertEquals("FF", StaticHelper.ENCODE_HEX(new byte[]{(byte) 0xFF}));
        }

        @Test
        void testEncodeHex_withMultipleBytes() {
            byte[] data = new byte[]{0x01, 0x23, 0x45, 0x67, (byte) 0x89, (byte) 0xAB, (byte) 0xCD, (byte) 0xEF};
            assertEquals("0123456789ABCDEF", StaticHelper.ENCODE_HEX(data));
        }

        @Test
        void testEncodeHex_withMixedBytes() {
            byte[] data = new byte[]{0x0A, 0x0B, 0x0C, 0x0D};
            assertEquals("0A0B0C0D", StaticHelper.ENCODE_HEX(data));
        }
    }

    @Nested
    class DecodeHexTests {

        @Test
        void testDecodeHex_withValidHex() {
            assertArrayEquals(new byte[]{0x01, 0x23}, StaticHelper.DECODE_HEX("0123"));
            assertArrayEquals(new byte[]{(byte) 0xFF}, StaticHelper.DECODE_HEX("FF"));
        }

        @Test
        void testDecodeHex_withLowercaseHex() {
            assertArrayEquals(new byte[]{(byte) 0xab, (byte) 0xcd}, StaticHelper.DECODE_HEX("abcd"));
        }

        @Test
        void testDecodeHex_withMixedCaseHex() {
            assertArrayEquals(new byte[]{(byte) 0xAB, (byte) 0xCD}, StaticHelper.DECODE_HEX("AbCd"));
        }

        @Test
        void testDecodeHex_throwsExceptionForOddLength() {
            assertThrows(IllegalArgumentException.class, () -> StaticHelper.DECODE_HEX("123"));
        }

        @Test
        void testDecodeHex_throwsExceptionForInvalidCharacters() {
            assertThrows(IllegalArgumentException.class, () -> StaticHelper.DECODE_HEX("GGGG"));
            assertThrows(IllegalArgumentException.class, () -> StaticHelper.DECODE_HEX("12XY"));
        }

        @Test
        void testDecodeHex_roundTrip() {
            byte[] original = new byte[]{0x01, 0x23, 0x45, 0x67, (byte) 0x89, (byte) 0xAB, (byte) 0xCD, (byte) 0xEF};
            String encoded = StaticHelper.ENCODE_HEX(original);
            byte[] decoded = StaticHelper.DECODE_HEX(encoded);
            assertArrayEquals(original, decoded);
        }
    }

    @Nested
    class CapitalizeTests {

        @Test
        void testCapitalize_withLowercaseString() {
            assertEquals("Hello", StaticHelper.CAPITALIZE("hello"));
        }

        @Test
        void testCapitalize_withUppercaseString() {
            assertEquals("HELLO", StaticHelper.CAPITALIZE("HELLO"));
        }

        @Test
        void testCapitalize_withSingleCharacter() {
            assertEquals("A", StaticHelper.CAPITALIZE("a"));
        }

        @Test
        void testCapitalize_withMixedCase() {
            assertEquals("Test", StaticHelper.CAPITALIZE("test"));
        }
    }

    @Nested
    class VarLengthSIntTests {

        @Test
        void testGetVarLengthSIntInBits_withZero() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(0L));
        }

        @Test
        void testGetVarLengthSIntInBits_withPositiveValue() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(63L));
            assertEquals(16, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(64L));
        }

        @Test
        void testGetVarLengthSIntInBits_withNegativeValue() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(-64L));
            assertEquals(16, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(-65L));
        }

        @Test
        void testGetVarLengthSIntInBits_withLargeValues() {
            assertEquals(24, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(8192L));
            assertEquals(24, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(-8193L));
        }

        @Test
        void testGetVarLengthSIntInBits_withBigInteger() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(BigInteger.ZERO));
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(BigInteger.valueOf(50)));
            assertEquals(16, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(BigInteger.valueOf(100)));
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_SINT_IN_BITS(BigInteger.valueOf(-50)));
        }
    }

    @Nested
    class VarLengthUIntTests {

        @Test
        void testGetVarLengthUIntInBits_withZero() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(0L));
        }

        @Test
        void testGetVarLengthUIntInBits_withSmallValue() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(127L));
        }

        @Test
        void testGetVarLengthUIntInBits_with2ByteValue() {
            assertEquals(16, StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(128L));
        }

        @Test
        void testGetVarLengthUIntInBits_throwsExceptionForNegativeValue() {
            assertThrows(IllegalArgumentException.class, () -> StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(-1L));
        }

        @Test
        void testGetVarLengthUIntInBits_withBigInteger() {
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(BigInteger.ZERO));
            assertEquals(8, StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(BigInteger.valueOf(100)));
            assertEquals(16, StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(BigInteger.valueOf(200)));
        }

        @Test
        void testGetVarLengthUIntInBits_throwsExceptionForNegativeBigInteger() {
            assertThrows(IllegalArgumentException.class,
                () -> StaticHelper.GET_VAR_LENGTH_UINT_IN_BITS(BigInteger.valueOf(-1)));
        }
    }
}
