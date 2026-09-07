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
package org.apache.plc4x.java.spi.buffers.bytebased.encoding;

import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferRaw;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.function.Executable;
import org.mockito.Mockito;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.function.BiFunction;

import static org.junit.jupiter.api.Assertions.*;

public class BaseEncodingRawTest {

    @Test
    public void testUnsupportedOperation() {

    }

    protected void assertEncodeDecode(byte value, int numBits, BiFunction<Integer, ReadBufferRaw, Byte> decoder, BiFunction<Integer, Byte, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertEncodeDecode(short value, int numBits, BiFunction<Integer, ReadBufferRaw, Short> decoder, BiFunction<Integer, Short, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertEncodeDecode(int value, int numBits, BiFunction<Integer, ReadBufferRaw, Integer> decoder, BiFunction<Integer, Integer, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertEncodeDecode(long value, int numBits, BiFunction<Integer, ReadBufferRaw, Long> decoder, BiFunction<Integer, Long, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertEncodeDecode(BigInteger value, int numBits, BiFunction<Integer, ReadBufferRaw, BigInteger> decoder, BiFunction<Integer, BigInteger, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertEncodeDecode(float value, int numBits, BiFunction<Integer, ReadBufferRaw, Float> decoder, BiFunction<Integer, Float, byte[]> encoder, float delta) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read), delta);
    }

    protected void assertEncodeDecode(double value, int numBits, BiFunction<Integer, ReadBufferRaw, Double> decoder, BiFunction<Integer, Double, byte[]> encoder, double delta) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read), delta);
    }

    protected void assertEncodeDecode(BigDecimal value, int numBits, BiFunction<Integer, ReadBufferRaw, BigDecimal> decoder, BiFunction<Integer, BigDecimal, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertEncodeDecode(String value, int numBits, BiFunction<Integer, ReadBufferRaw, String> decoder, BiFunction<Integer, String, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        ReadBufferByteBased read = new ReadBufferByteBased(encoded);
        assertEquals(value, decoder.apply(numBits, read));
    }

    protected void assertUnsupportedOperation(Executable executable) {
        assertThrows(UnsupportedOperationException.class, executable);
    }

    @Test
    protected void assertUnsupportedOperations() {
        EncodingRaw encoding = new BaseEncodingRaw() {
            @Override
            public String getName() {
                return "dummy";
            }
        };
        ReadBufferRaw mockReadBuffer = Mockito.mock(ReadBufferRaw.class);

        assertUnsupportedOperation(() -> encoding.decodeByte(8, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeByte(8, (byte) 41));

        assertUnsupportedOperation(() -> encoding.decodeShort(16, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeShort(16, (short) 42));

        assertUnsupportedOperation(() -> encoding.decodeInt(32, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeInt(32, 42));

        assertUnsupportedOperation(() -> encoding.decodeLong(64, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeLong(64, 4L));

        assertUnsupportedOperation(() -> encoding.decodeBigInteger(128, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeBigInteger(128, BigInteger.valueOf(42)));

        assertUnsupportedOperation(() -> encoding.decodeFloat(32, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeFloat(32, 42.23F));

        assertUnsupportedOperation(() -> encoding.decodeDouble(64, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeDouble(64, 42.23));

        assertUnsupportedOperation(() -> encoding.decodeBigDecimal(128, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeBigDecimal(128, BigDecimal.valueOf(42.23)));

        assertUnsupportedOperation(() -> encoding.decodeString(64, mockReadBuffer));
        assertUnsupportedOperation(() -> encoding.encodeString(64, "Hurz"));
    }

}
