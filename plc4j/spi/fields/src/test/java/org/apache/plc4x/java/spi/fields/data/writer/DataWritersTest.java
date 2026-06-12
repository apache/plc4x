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

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyWriteBuffer;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

import static org.junit.jupiter.api.Assertions.*;

class DataWritersTest {

    @Test
    void booleanWriterDelegatesToWriteBit() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterFactory.writeBoolean(wb).write(true);
        assertEquals(Boolean.TRUE, wb.bitWritten);
    }

    @Test
    void byteWriterHappyAndInvalidBitLength() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterFactory.writeByte(wb, 8).write((byte) 0xA5);
        assertArrayEquals(new byte[] {(byte) 0xA5}, wb.bitsWritten);

        DataWriterSimpleByte invalid = new DataWriterSimpleByte(wb, 7);
        assertThrows(BufferException.class, () -> invalid.write((byte) 1));
    }

    @Test
    void unsignedWritersDelegate() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterFactory.writeUnsignedByte(wb, 8).write((byte) 250);
        assertEquals((byte)250, wb.unsignedByteWritten);

        DataWriterFactory.writeUnsignedShort(wb, 16).write((short) 0xBEEF);
        assertEquals((short)0xBEEF, wb.unsignedShortWritten);

        DataWriterFactory.writeUnsignedInt(wb, 32).write(0xCAFEBABE);
        assertEquals(0xCAFEBABE, wb.unsignedIntWritten);

        DataWriterFactory.writeUnsignedLong(wb, 64).write(0x0123456789ABCDEFL);
        assertEquals(0x0123456789ABCDEFL, wb.unsignedLongWritten);

        BigInteger big = new BigInteger("12345678901234567890");
        DataWriterFactory.writeUnsignedBigInteger(wb, 80).write(big);
        assertEquals(big, wb.unsignedBigIntegerWritten);
    }

    @Test
    void signedWritersDelegate() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterFactory.writeSignedByte(wb, 8).write((byte) -100);
        assertEquals((byte)-100, (byte) wb.signedByteWritten);

        DataWriterFactory.writeSignedShort(wb, 16).write((short) -12345);
        assertEquals((short)-12345, (short) wb.signedShortWritten);

        DataWriterFactory.writeSignedInt(wb, 32).write(-123456789);
        assertEquals(-123456789, (int) wb.signedIntWritten);

        DataWriterFactory.writeSignedLong(wb, 64).write(-1234567890123456789L);
        assertEquals(-1234567890123456789L, (long) wb.signedLongWritten);

        BigInteger sb = new BigInteger("-98765432109876543210");
        DataWriterFactory.writeSignedBigInteger(wb, 88).write(sb);
        assertEquals(sb, wb.signedBigIntegerWritten);
    }

    @Test
    void floatDoubleBigDecimalAndStringWritersDelegate() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterFactory.writeFloat(wb, 32).write(3.14159f);
        assertEquals(3.14159f, wb.floatWritten);

        DataWriterFactory.writeDouble(wb, 64).write(Math.E);
        assertEquals(Math.E, wb.doubleWritten);

        DataWriterSimpleBigDecimal bigDecWriter = new DataWriterSimpleBigDecimal(wb, 128);
        bigDecWriter.write(new BigDecimal("123.456"));
        assertEquals(new BigDecimal("123.456"), wb.bigDecimalWritten);

        DataWriterFactory.writeString(wb, 40).write("Hello");
        assertEquals("Hello", wb.stringWritten);
    }

    @Test
    void byteArrayWriterDelegatesToWriteBits() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        byte[] payload = new byte[] {1,2,3,4,5};
        DataWriterFactory.writeByteArray(wb, payload.length).write(payload);
        assertArrayEquals(payload, wb.bitsWritten);
    }

    enum SampleEnum { ONE, TWO }

    @Test
    void enumWriterDelegatesToWrappedWriter() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterEnumDefault<SampleEnum, Integer> enumWriter = DataWriterFactory.writeEnum(
            e -> e == SampleEnum.TWO ? 2 : 1,
            Enum::name,
            DataWriterFactory.writeUnsignedInt(wb, 8)
        );
        enumWriter.write(SampleEnum.TWO);
        assertEquals(2, (int) wb.unsignedIntWritten);

        assertTrue(wb.context.isEmpty());
        enumWriter.pushContext(WithOption.WithName("test"));
        assertFalse(wb.context.isEmpty());
        enumWriter.popContext();
        assertTrue(wb.context.isEmpty());

        assertEquals(wb, enumWriter.getWriteBuffer());
    }

    @Test
    void complexWriterUsesMessageSerializer() throws Exception {
        // Complex writer relies on Message#serialize, which we cannot easily emulate here.
        // Just ensure the factory can create it without exceptions.
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterComplexDefault<Message> complex = DataWriterFactory.writeComplex(wb);
        assertNotNull(complex);
        Message message = new Message() {
            @Override
            public int getLengthInBytes() {
                return 1;
            }

            @Override
            public int getLengthInBits() {
                return 8;
            }

            @Override
            public void serialize(WriteBuffer writeBuffer) throws BufferException {
                writeBuffer.writeBits(8 * 8, new byte[]{42});
            }
        };
        complex.write(message);
        assertEquals(message,  wb.messageWritten);

        assertTrue(wb.context.isEmpty());
        complex.pushContext(WithOption.WithName("test"));
        assertFalse(wb.context.isEmpty());
        complex.popContext();
        assertTrue(wb.context.isEmpty());

        assertEquals(wb, complex.getWriteBuffer());
    }

    @Test
    void dateTimeWritersDelegate() throws Exception {
        DummyWriteBuffer wb = new DummyWriteBuffer();
        DataWriterFactory.writeDate(wb).write(LocalDate.of(2024,12,31));
        assertNotNull(wb.unsignedLongWritten); // internal representation is written as unsigned int days

        DataWriterFactory.writeDateTime(wb).write(LocalDateTime.of(2024,12,31,23,59,58));
        assertNotNull(wb.unsignedLongWritten); // internal representation is written as unsigned long seconds

        DataWriterFactory.writeTime(wb).write(LocalTime.of(12,34,56));
        assertNotNull(wb.unsignedLongWritten); // internal representation is written as unsigned int seconds
    }
}
