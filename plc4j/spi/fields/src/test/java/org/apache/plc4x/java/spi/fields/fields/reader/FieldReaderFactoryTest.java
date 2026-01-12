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
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyReadBuffer;
import org.apache.plc4x.java.spi.fields.fields.TestFieldIoStubs.SimpleDataReader;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.util.List;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

class FieldReaderFactoryTest {

    // TODO: Implement some sensible test for this type.
    /*@Test
    void readAbstractFieldDelegates() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "abstract");
        String result = FieldReaderFactory.readAbstractField(dataReader);
        assertEquals("abstract", result);
    }*/

    @Test
    void readCountArrayFieldWithBigInteger() throws Exception {
        AtomicInteger counter = new AtomicInteger();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(counter::getAndIncrement);

        BigInteger count = BigInteger.valueOf(3);
        List<Integer> result = FieldReaderFactory.readCountArrayField(dataReader, count);

        assertEquals(List.of(0, 1, 2), result);
    }

    @Test
    void readCountArrayFieldWithLong() throws Exception {
        AtomicInteger counter = new AtomicInteger();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(counter::getAndIncrement);

        List<Integer> result = FieldReaderFactory.readCountArrayField(dataReader, 3L);

        assertEquals(List.of(0, 1, 2), result);
    }

    @Test
    void readCountArrayFieldBigIntegerOverflowThrows() {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 1);
        BigInteger hugeBigInteger = BigInteger.valueOf(2).pow(65); // More than 64 bits

        assertThrows(IllegalStateException.class, () ->
            FieldReaderFactory.readCountArrayField(dataReader, hugeBigInteger));
    }

    @Test
    void readLengthArrayFieldWithInt() throws Exception {
        AtomicInteger counter = new AtomicInteger();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(counter::getAndIncrement);

//        List<Integer> result = FieldReaderFactory.readLengthArrayField(dataReader, 2);

//        assertEquals(2, result.size());
    }

    @Test
    void readLengthArrayFieldWithLong() throws Exception {
        AtomicInteger counter = new AtomicInteger();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(counter::getAndIncrement);

        // TODO: Implement this ... Currently there are issues as the SimpleDataReader does not increment the positionInBits.
//        List<Integer> result = FieldReaderFactory.readLengthArrayField(dataReader, 2L);

        //assertEquals(2, result.size());
    }

    @Test
    void readTerminatedArrayFieldStopsOnCondition() throws Exception {
        AtomicInteger counter = new AtomicInteger();
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(counter::getAndIncrement);

        List<Integer> result = FieldReaderFactory.readTerminatedArrayField(
            dataReader,
            () -> counter.get() >= 2
        );

        assertEquals(List.of(0, 1), result);
    }

    @Test
    void readAssertFieldDelegates() throws Exception {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 42);
        Integer result = FieldReaderFactory.readAssertField(dataReader, 42);
        assertEquals(42, result);
    }

    @Test
    void readChecksumFieldDelegates() throws Exception {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 123);
        Integer result = FieldReaderFactory.readChecksumField(dataReader, 123);
        assertEquals(123, result);
    }

    @Test
    void readConstFieldDelegates() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "constant");
        String result = FieldReaderFactory.readConstField(dataReader, "constant");
        assertEquals("constant", result);
    }

    @Test
    void readDiscriminatorFieldDelegates() throws Exception {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 999);
        Integer result = FieldReaderFactory.readDiscriminatorField(dataReader);
        assertEquals(999, result);
    }

    @Test
    void readDiscriminatorEnumFieldDelegates() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "ENUM_VALUE");
        String result = FieldReaderFactory.readDiscriminatorEnumField(dataReader);
        assertEquals("ENUM_VALUE", result);
    }

    @Test
    void readEnumFieldDelegates() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "ENUM");
        String result = FieldReaderFactory.readEnumField(dataReader);
        assertEquals("ENUM", result);
    }

    @Test
    void readImplicitFieldDelegates() throws Exception {
        SimpleDataReader<Boolean> dataReader = new SimpleDataReader<>(() -> true);
        Boolean result = FieldReaderFactory.readImplicitField(dataReader);
        assertTrue(result);
    }

    @Test
    void readOptionalFieldWithDefaultCondition() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "optional");
        String result = FieldReaderFactory.readOptionalField(dataReader);
        assertEquals("optional", result);
    }

    @Test
    void readOptionalFieldWithCondition() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "conditional");
        String result = FieldReaderFactory.readOptionalField(dataReader, true);
        assertEquals("conditional", result);
    }

    @Test
    void readOptionalFieldWithFalseCondition() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "should_not_read");
        String result = FieldReaderFactory.readOptionalField(dataReader, false);
        assertNull(result);
        assertEquals(0, dataReader.getReadCount()); // Should not be called
    }

    @Test
    void readManualFieldDelegates() throws Exception {
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        String result = FieldReaderFactory.readManualField(readBuffer, () -> "manual");
        assertEquals("manual", result);
    }

    @Test
    void readManualArrayFieldDelegates() throws Exception {
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        AtomicInteger counter = new AtomicInteger();

        List<Integer> result = FieldReaderFactory.readManualArrayField(
            readBuffer,
            list -> list.size() >= 2,
            counter::getAndIncrement
        );

        assertEquals(List.of(0, 1), result);
    }

    @Test
    void readManualByteArrayFieldDelegates() throws Exception {
        DummyReadBuffer readBuffer = new DummyReadBuffer();
        AtomicInteger counter = new AtomicInteger();

        byte[] result = FieldReaderFactory.readManualByteArrayField(
            readBuffer,
            list -> list.size() >= 3,
            () -> (byte) counter.getAndIncrement()
        );

        assertArrayEquals(new byte[]{0, 1, 2}, result);
    }

    @Test
    void readPaddingFieldDelegates() throws Exception {
        SimpleDataReader<Byte> dataReader = new SimpleDataReader<>(() -> (byte) 0);
        assertDoesNotThrow(() ->
            FieldReaderFactory.readPaddingField(dataReader, 2)
        );
        assertEquals(2, dataReader.getReadCount());
    }

    @Test
    void readReservedFieldDelegates() throws Exception {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 0);
        Integer result = FieldReaderFactory.readReservedField(dataReader, 0);
        assertEquals(0, result);
    }

    @Test
    void readSimpleFieldDelegates() throws Exception {
        SimpleDataReader<Double> dataReader = new SimpleDataReader<>(() -> 3.14);
        Double result = FieldReaderFactory.readSimpleField(dataReader);
        assertEquals(3.14, result);
    }

    @Test
    void readUnknownFieldDelegates() throws Exception {
        SimpleDataReader<Object> dataReader = new SimpleDataReader<>(() -> "unknown");
        Object result = FieldReaderFactory.readUnknownField(dataReader);
        assertEquals("unknown", result);
    }

    @Test
    void readVirtualFieldDelegates() throws Exception {
        String result = FieldReaderFactory.readVirtualField(String.class, "virtual");
        assertEquals("virtual", result);
    }

    @Test
    void readPeekFieldWithDefaultOffset() throws Exception {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 555);
        Integer result = FieldReaderFactory.readPeekField(dataReader);
        assertEquals(555, result);
    }

    @Test
    void readPeekFieldWithOffset() throws Exception {
        SimpleDataReader<Integer> dataReader = new SimpleDataReader<>(() -> 666);
        Integer result = FieldReaderFactory.readPeekField(dataReader, 5);
        assertEquals(666, result);
    }

    @Test
    void factoryMethodsWithOptions() throws Exception {
        SimpleDataReader<String> dataReader = new SimpleDataReader<>(() -> "test");
        WithOption option = WithOption.WithName("testField");

        String result = FieldReaderFactory.readSimpleField(dataReader, option);
        assertEquals("test", result);
    }
}