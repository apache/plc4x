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

package org.apache.plc4x.java.spi.fields.data.writer;

import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyWriteBuffer;
import org.apache.plc4x.java.spi.fields.data.reader.DataIoSerializerFunction;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DataWriterFactoryTest {

    @Test
    void factoryCreatesCorrectWriterTypes() {
        DummyWriteBuffer wb = new DummyWriteBuffer();

        assertTrue(DataWriterFactory.writeBoolean(wb) instanceof DataWriterSimpleBoolean);
        assertTrue(DataWriterFactory.writeByte(wb, 8) instanceof DataWriterSimpleByte);
        assertTrue(DataWriterFactory.writeByteArray(wb, 10) instanceof DataWriterSimpleByteArray);

        assertTrue(DataWriterFactory.writeUnsignedByte(wb, 8) instanceof DataWriterSimpleUnsignedByte);
        assertTrue(DataWriterFactory.writeUnsignedShort(wb, 16) instanceof DataWriterSimpleUnsignedShort);
        assertTrue(DataWriterFactory.writeUnsignedInt(wb, 32) instanceof DataWriterSimpleUnsignedInt);
        assertTrue(DataWriterFactory.writeUnsignedLong(wb, 64) instanceof DataWriterSimpleUnsignedLong);
        assertTrue(DataWriterFactory.writeUnsignedBigInteger(wb, 128) instanceof DataWriterSimpleUnsignedBigInteger);

        assertTrue(DataWriterFactory.writeSignedByte(wb, 8) instanceof DataWriterSimpleSignedByte);
        assertTrue(DataWriterFactory.writeSignedShort(wb, 16) instanceof DataWriterSimpleSignedShort);
        assertTrue(DataWriterFactory.writeSignedInt(wb, 32) instanceof DataWriterSimpleSignedInt);
        assertTrue(DataWriterFactory.writeSignedLong(wb, 64) instanceof DataWriterSimpleSignedLong);
        assertTrue(DataWriterFactory.writeSignedBigInteger(wb, 128) instanceof DataWriterSimpleSignedBigInteger);

        assertTrue(DataWriterFactory.writeFloat(wb, 32) instanceof DataWriterSimpleFloat);
        assertTrue(DataWriterFactory.writeDouble(wb, 64) instanceof DataWriterSimpleDouble);
        assertTrue(DataWriterFactory.writeString(wb, 80) instanceof DataWriterSimpleString);

        assertTrue(DataWriterFactory.writeComplex(wb) instanceof DataWriterComplexDefault);

        assertTrue(DataWriterFactory.writeDate(wb) instanceof DataWriterSimpleDateFromEpochDays);
        assertTrue(DataWriterFactory.writeDateTime(wb) instanceof DataWriterSimpleDateTimeFromEpochSeconds);
        assertTrue(DataWriterFactory.writeTime(wb) instanceof DataWriterSimpleTimeFromDaySeconds);
    }

    @Test
    void enumFactoryCreateEnumWriter() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriter<Integer> intWriter = DataWriterFactory.writeUnsignedInt(wb, 32);

        DataWriterEnumDefault<String, Integer> enumWriter = DataWriterFactory.writeEnum(
            s -> s.length(),
            s -> "String: " + s,
            intWriter
        );

        assertNotNull(enumWriter);
        assertTrue(enumWriter instanceof DataWriterEnumDefault);
    }

    @Test
    void dataIoFactoryCreatesDataIoWriter() {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataIoSerializerFunction serializer = (writeBuffer, value) -> {};

        DataWriterDataIoDefault writer = DataWriterFactory.writeDataIO(wb, serializer);

        assertNotNull(writer);
        assertTrue(writer instanceof DataWriterDataIoDefault);
    }

    @Test
    void factoryMethodsReturnNonNullInstances() {
        DummyWriteBuffer wb = new DummyWriteBuffer();

        assertNotNull(DataWriterFactory.writeBoolean(wb));
        assertNotNull(DataWriterFactory.writeByte(wb, 8));
        assertNotNull(DataWriterFactory.writeByteArray(wb, 10));

        assertNotNull(DataWriterFactory.writeUnsignedByte(wb, 8));
        assertNotNull(DataWriterFactory.writeUnsignedShort(wb, 16));
        assertNotNull(DataWriterFactory.writeUnsignedInt(wb, 32));
        assertNotNull(DataWriterFactory.writeUnsignedLong(wb, 64));
        assertNotNull(DataWriterFactory.writeUnsignedBigInteger(wb, 128));

        assertNotNull(DataWriterFactory.writeSignedByte(wb, 8));
        assertNotNull(DataWriterFactory.writeSignedShort(wb, 16));
        assertNotNull(DataWriterFactory.writeSignedInt(wb, 32));
        assertNotNull(DataWriterFactory.writeSignedLong(wb, 64));
        assertNotNull(DataWriterFactory.writeSignedBigInteger(wb, 128));

        assertNotNull(DataWriterFactory.writeFloat(wb, 32));
        assertNotNull(DataWriterFactory.writeDouble(wb, 64));
        assertNotNull(DataWriterFactory.writeString(wb, 80));

        assertNotNull(DataWriterFactory.writeComplex(wb));

        assertNotNull(DataWriterFactory.writeDate(wb));
        assertNotNull(DataWriterFactory.writeDateTime(wb));
        assertNotNull(DataWriterFactory.writeTime(wb));
    }

    @Test
    void factoryMethodsAcceptDifferentBitLengths() {
        DummyWriteBuffer wb = new DummyWriteBuffer();

        // Test that factory methods accept various bit lengths
        assertDoesNotThrow(() -> DataWriterFactory.writeByte(wb, 8));
        assertDoesNotThrow(() -> DataWriterFactory.writeUnsignedShort(wb, 12));
        assertDoesNotThrow(() -> DataWriterFactory.writeUnsignedInt(wb, 24));
        assertDoesNotThrow(() -> DataWriterFactory.writeUnsignedLong(wb, 48));
        assertDoesNotThrow(() -> DataWriterFactory.writeUnsignedBigInteger(wb, 256));

        assertDoesNotThrow(() -> DataWriterFactory.writeSignedByte(wb, 8));
        assertDoesNotThrow(() -> DataWriterFactory.writeSignedShort(wb, 16));
        assertDoesNotThrow(() -> DataWriterFactory.writeSignedInt(wb, 32));
        assertDoesNotThrow(() -> DataWriterFactory.writeSignedLong(wb, 64));
        assertDoesNotThrow(() -> DataWriterFactory.writeSignedBigInteger(wb, 128));

        assertDoesNotThrow(() -> DataWriterFactory.writeFloat(wb, 32));
        assertDoesNotThrow(() -> DataWriterFactory.writeDouble(wb, 64));
        assertDoesNotThrow(() -> DataWriterFactory.writeString(wb, 160));
        assertDoesNotThrow(() -> DataWriterFactory.writeByteArray(wb, 1024));
    }
}