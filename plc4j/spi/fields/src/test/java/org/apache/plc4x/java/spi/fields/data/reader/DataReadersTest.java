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

package org.apache.plc4x.java.spi.fields.data.reader;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.TestBuffers.DummyReadBuffer;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.math.BigInteger;

import static org.junit.jupiter.api.Assertions.*;

class DataReadersTest {

    @Test
    void booleanReaderDelegatesToReadBit() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.bitValue = true;
        assertTrue(DataReaderFactory.readBoolean(rb).read());
    }

    @Test
    void byteReaderHappyAndInvalidBitLength() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.bitsValue = new byte[] { (byte) 0xA5 };
        assertEquals((byte)0xA5, DataReaderFactory.readByte(rb, 8).read());

        DataReaderSimpleByte invalid = new DataReaderSimpleByte(rb, 7);
        assertThrows(BufferException.class, invalid::read);
    }

    @Test
    void unsignedReadersTests() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.unsignedByteValue = (byte) 250;
        assertEquals((byte)250, DataReaderFactory.readUnsignedByte(rb, 8).read());

        rb.unsignedShortValue = (short) 0xBEEF;
        assertEquals((short)0xBEEF, DataReaderFactory.readUnsignedShort(rb, 16).read());

        rb.unsignedIntValue = 0xCAFEBABE;
        assertEquals(0xCAFEBABE, DataReaderFactory.readUnsignedInt(rb, 32).read());

        rb.unsignedLongValue = 0x0123456789ABCDEFL;
        assertEquals(0x0123456789ABCDEFL, DataReaderFactory.readUnsignedLong(rb, 64).read());

        rb.unsignedBigIntegerValue = new BigInteger("12345678901234567890");
        assertEquals(rb.unsignedBigIntegerValue, DataReaderFactory.readUnsignedBigInteger(rb, 80).read());
    }

    @Test
    void signedReadersTests() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.signedByteValue = (byte) -100;
        assertEquals((byte)-100, DataReaderFactory.readSignedByte(rb, 8).read());

        rb.signedShortValue = (short) -12345;
        assertEquals((short)-12345, DataReaderFactory.readSignedShort(rb, 16).read());

        rb.signedIntValue = -123456789;
        assertEquals(-123456789, DataReaderFactory.readSignedInt(rb, 32).read());

        rb.signedLongValue = -1234567890123456789L;
        assertEquals(-1234567890123456789L, DataReaderFactory.readSignedLong(rb, 64).read());

        rb.signedBigIntegerValue = new BigInteger("-98765432109876543210");
        assertEquals(rb.signedBigIntegerValue, DataReaderFactory.readSignedBigInteger(rb, 88).read());
    }

    @Test
    void floatingPointReadersTests() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.floatValue = 3.14159f;
        assertEquals(3.14159f, DataReaderFactory.readFloat(rb, 32).read(), 1e-6);

        rb.doubleValue = Math.E;
        assertEquals(Math.E, DataReaderFactory.readDouble(rb, 64).read(), 1e-10);

        rb.bigDecimalValue = new BigDecimal("123.456");
        // BigDecimal reader is not in factory; ensure class exists and delegates via direct constructor
        DataReaderSimpleBigDecimal bigDecReader = new DataReaderSimpleBigDecimal(rb, 128);
        assertEquals(rb.bigDecimalValue, bigDecReader.read());
    }

    @Test
    void stringReadersTests() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.stringValue = "Hello";
        assertEquals("Hello", DataReaderFactory.readString(rb, 40).read());
    }

    @Test
    void byteArrayReaderDelegatesToReadBits() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        byte[] payload = new byte[] {1,2,3,4,5};
        rb.bitsValue = payload;
        assertArrayEquals(payload, DataReaderFactory.readByteArray(rb, payload.length).read());
    }

    @Test
    void enumReaderDelegatesToWrappedReader() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        rb.unsignedIntValue = 2;
        DataReaderEnumDefault<String,Integer> enumReader = DataReaderFactory.readEnum(i -> i == 2 ? "TWO" : "OTHER", DataReaderFactory.readUnsignedInt(rb, 8));
        assertEquals("TWO", enumReader.read());

        rb.setPositionInBits(42);
        int positionInBits = enumReader.getPositionInBits();
        assertEquals(42, positionInBits);

        enumReader.setPositionInBits(positionInBits + 1);
        assertEquals(43, rb.getPositionInBits());

        assertTrue(rb.context.isEmpty());
        enumReader.pushContext(WithOption.WithName("test"));
        assertFalse(rb.context.isEmpty());
        enumReader.popContext();
        assertTrue(rb.context.isEmpty());

        assertEquals(rb, enumReader.getReadBuffer());
    }

    @Test
    void complexReaderUsesSupplier() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        // We will simulate reading by setting bitsValue and using two byte reads from rb via DataReaderFactory
        rb.bitsValue = new byte[] { 0x12, 0x34 };
        DataReaderComplexDefault<Integer> complex = DataReaderFactory.readComplex(() -> {
            try {
                byte b1 = DataReaderFactory.readByte(rb, 8).read();
                byte b2 = DataReaderFactory.readByte(rb, 8).read();
                return ((b1 & 0xFF) << 8) | (b2 & 0xFF);
            } catch (BufferException e) { throw new RuntimeException(e); }
        }, rb);
        assertEquals(0x1234, complex.read());

        rb.setPositionInBits(42);
        int positionInBits = complex.getPositionInBits();
        assertEquals(42, positionInBits);

        complex.setPositionInBits(positionInBits + 1);
        assertEquals(43, rb.getPositionInBits());

        assertTrue(rb.context.isEmpty());
        complex.pushContext(WithOption.WithName("test"));
        assertFalse(rb.context.isEmpty());
        complex.popContext();
        assertTrue(rb.context.isEmpty());

        assertEquals(rb, complex.getReadBuffer());
    }

    @Test
    void dateTimeReadersExist() throws Exception {
        DummyReadBuffer rb = new DummyReadBuffer();
        // The date/time readers delegate into ReadBuffer via custom implementations;
        // Here we assert the factory constructs instances and calling read doesn't throw with dummy defaults.
        assertDoesNotThrow(() -> DataReaderFactory.readDate(rb).read());
        assertDoesNotThrow(() -> DataReaderFactory.readDateTime(rb).read());
        assertDoesNotThrow(() -> DataReaderFactory.readTime(rb).read());
    }
}
