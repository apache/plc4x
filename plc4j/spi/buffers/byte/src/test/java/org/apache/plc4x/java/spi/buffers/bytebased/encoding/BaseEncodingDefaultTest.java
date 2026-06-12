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

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.junit.jupiter.api.function.Executable;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.function.BiFunction;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Base class for encoding tests providing common test utilities.
 */
public abstract class BaseEncodingDefaultTest {

    @FunctionalInterface
    protected interface ThrowingDecoder<T> {
        T apply(int numBits, byte[] bytes) throws BufferException;
    }

    @FunctionalInterface
    protected interface ThrowingEncoder<T> {
        byte[] apply(int numBits, T value) throws BufferException;
    }

    protected <T> BiFunction<Integer, byte[], T> wrapDecoder(ThrowingDecoder<T> d) {
        return (n, b) -> {
            try {
                return d.apply(n, b);
            } catch (BufferException e) {
                throw new RuntimeException(e);
            }
        };
    }

    protected <T> BiFunction<Integer, T, byte[]> wrapEncoder(ThrowingEncoder<T> e) {
        return (n, v) -> {
            try {
                return e.apply(n, v);
            } catch (BufferException ex) {
                throw new RuntimeException(ex);
            }
        };
    }

    protected void assertEncodeDecode(byte value, int numBits, BiFunction<Integer, byte[], Byte> decoder, BiFunction<Integer, Byte, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertEncodeDecode(short value, int numBits, BiFunction<Integer, byte[], Short> decoder, BiFunction<Integer, Short, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertEncodeDecode(int value, int numBits, BiFunction<Integer, byte[], Integer> decoder, BiFunction<Integer, Integer, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertEncodeDecode(long value, int numBits, BiFunction<Integer, byte[], Long> decoder, BiFunction<Integer, Long, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertEncodeDecode(BigInteger value, int numBits, BiFunction<Integer, byte[], BigInteger> decoder, BiFunction<Integer, BigInteger, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertEncodeDecode(float value, int numBits, BiFunction<Integer, byte[], Float> decoder, BiFunction<Integer, Float, byte[]> encoder, float delta) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded), delta);
    }

    protected void assertEncodeDecode(double value, int numBits, BiFunction<Integer, byte[], Double> decoder, BiFunction<Integer, Double, byte[]> encoder, double delta) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded), delta);
    }

    protected void assertEncodeDecode(BigDecimal value, int numBits, BiFunction<Integer, byte[], BigDecimal> decoder, BiFunction<Integer, BigDecimal, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertEncodeDecode(String value, int numBits, BiFunction<Integer, byte[], String> decoder, BiFunction<Integer, String, byte[]> encoder) {
        byte[] encoded = encoder.apply(numBits, value);
        assertNotNull(encoded);
        assertEquals(value, decoder.apply(numBits, encoded));
    }

    protected void assertUnsupportedOperation(Executable executable) {
        assertThrows(UnsupportedOperationException.class, executable);
    }

}