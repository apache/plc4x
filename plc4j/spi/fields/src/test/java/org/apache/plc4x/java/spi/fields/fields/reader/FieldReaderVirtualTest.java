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

package org.apache.plc4x.java.spi.fields.fields.reader;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.fields.fields.reader.FieldReaderVirtual;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderVirtualTest {

    @Test
    void readVirtualFieldBooleanPrimitive() throws Exception {
        FieldReaderVirtual<Boolean> reader = new FieldReaderVirtual<>();

        Boolean result = reader.readVirtualField(boolean.class, true);

        assertEquals(true, result);
    }

    @Test
    void readVirtualFieldBytePrimitive() throws Exception {
        FieldReaderVirtual<Byte> reader = new FieldReaderVirtual<>();

        Byte result = reader.readVirtualField(byte.class, (byte) 42);

        assertEquals((byte) 42, result);
    }

    @Test
    void readVirtualFieldByteOverflow() {
        FieldReaderVirtual<Byte> reader = new FieldReaderVirtual<>();

        assertThrows(ArithmeticException.class, () ->
            reader.readVirtualField(byte.class, 256L)
        );
    }

    @Test
    void readVirtualFieldByteSpecialCase0xFF() throws Exception {
        FieldReaderVirtual<Byte> reader = new FieldReaderVirtual<>();

        Byte result = reader.readVirtualField(byte.class, 0xFFL);

        assertEquals((byte) 0xFF, result);
    }

    @Test
    void readVirtualFieldShortPrimitive() throws Exception {
        FieldReaderVirtual<Short> reader = new FieldReaderVirtual<>();

        Short result = reader.readVirtualField(short.class, (short) 12345);

        assertEquals((short) 12345, result);
    }

    @Test
    void readVirtualFieldShortOverflow() {
        FieldReaderVirtual<Short> reader = new FieldReaderVirtual<>();

        assertThrows(ArithmeticException.class, () ->
            reader.readVirtualField(short.class, 65536L)
        );
    }

    @Test
    void readVirtualFieldIntPrimitive() throws Exception {
        FieldReaderVirtual<Integer> reader = new FieldReaderVirtual<>();

        Integer result = reader.readVirtualField(int.class, 123456789);

        assertEquals(123456789, result);
    }

    @Test
    void readVirtualFieldIntOverflow() {
        FieldReaderVirtual<Integer> reader = new FieldReaderVirtual<>();

        assertThrows(ArithmeticException.class, () ->
            reader.readVirtualField(int.class, (long) Integer.MAX_VALUE + 1)
        );
    }

    @Test
    void readVirtualFieldLongPrimitive() throws Exception {
        FieldReaderVirtual<Long> reader = new FieldReaderVirtual<>();

        Long result = reader.readVirtualField(long.class, 9876543210L);

        assertEquals(9876543210L, result);
    }

    @Test
    void readVirtualFieldCharPrimitive() throws Exception {
        FieldReaderVirtual<Character> reader = new FieldReaderVirtual<>();

        Character result = reader.readVirtualField(char.class, 'X');

        assertEquals('X', result);
    }

    @Test
    void readVirtualFieldFloatPrimitive() throws Exception {
        FieldReaderVirtual<Float> reader = new FieldReaderVirtual<>();

        Float result = reader.readVirtualField(float.class, 3.14159f);

        assertEquals(3.14159f, result, 0.0001f);
    }

    @Test
    void readVirtualFieldDoublePrimitive() throws Exception {
        FieldReaderVirtual<Double> reader = new FieldReaderVirtual<>();

        Double result = reader.readVirtualField(double.class, Math.E);

        assertEquals(Math.E, result, 0.0000001);
    }

    @Test
    void readVirtualFieldString() throws Exception {
        FieldReaderVirtual<String> reader = new FieldReaderVirtual<>();

        String result = reader.readVirtualField(String.class, "test value");

        assertEquals("test value", result);
    }

    @Test
    void readVirtualFieldStringFromNumber() throws Exception {
        FieldReaderVirtual<String> reader = new FieldReaderVirtual<>();

        String result = reader.readVirtualField(String.class, 42);

        assertEquals("42", result);
    }

    @Test
    void readVirtualFieldBigIntegerFromLong() throws Exception {
        FieldReaderVirtual<BigInteger> reader = new FieldReaderVirtual<>();

        BigInteger result = reader.readVirtualField(BigInteger.class, 12345L);

        assertEquals(BigInteger.valueOf(12345L), result);
    }

    @Test
    void readVirtualFieldBigIntegerFromNumber() throws Exception {
        FieldReaderVirtual<BigInteger> reader = new FieldReaderVirtual<>();

        BigInteger result = reader.readVirtualField(BigInteger.class, 67890);

        assertEquals(BigInteger.valueOf(67890L), result);
    }

    @Test
    void readVirtualFieldDirectCast() throws Exception {
        FieldReaderVirtual<Integer> reader = new FieldReaderVirtual<>();

        Integer result = reader.readVirtualField(Integer.class, 999);

        assertEquals(999, result);
    }

    @Test
    void readVirtualFieldWithOptions() throws Exception {
        FieldReaderVirtual<String> reader = new FieldReaderVirtual<>();

        WithOption option = WithOption.WithName("virtualField");
        String result = reader.readVirtualField(String.class, "virtual", option);

        assertEquals("virtual", result);
    }

    /*@Test
    void readVirtualFieldUnmappedPrimitive() {
        FieldReaderVirtual<Object> reader = new FieldReaderVirtual<>();

        // Using void.class as an unmapped primitive type
        assertThrows(IllegalStateException.class, () ->
            reader.readVirtualField(void.class, null)
        );
    }*/
}