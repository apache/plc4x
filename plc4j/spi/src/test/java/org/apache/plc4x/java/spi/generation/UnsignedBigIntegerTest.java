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
package org.apache.plc4x.java.spi.generation;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;

public class UnsignedBigIntegerTest {

    /**
     * In ASCII encoding the number is simply represented as a UTF-8 string.
     *
     * @throws Exception something went wrong
     */
    @Test
    public void testAsciiEncoding() throws Exception {
        String stringValue = "1234567891234567891";
        byte[] bytes = stringValue.getBytes(StandardCharsets.UTF_8);

        // Parse the input
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        BigInteger bigIntegerValue = readBufferByteBased.readUnsignedBigInteger("test", stringValue.length() * 8, WithReaderWriterArgs.WithEncoding("ASCII"));

        // Check the expected value matches
        Assertions.assertEquals(new BigInteger(stringValue), bigIntegerValue);

        // Serialize the value again
        WriteBufferByteBased writeBufferByteBased = new WriteBufferByteBased(stringValue.length());
        writeBufferByteBased.writeUnsignedBigInteger("test", stringValue.length() * 8, bigIntegerValue, WithReaderWriterArgs.WithEncoding("ASCII"));

        // Check the output matches the original one
        byte[] serializedBytes = writeBufferByteBased.getBytes();
        Assertions.assertArrayEquals(bytes, serializedBytes);
    }

    /**
     * In BCD encoding two digits are always encoded as the hex digits of a byte.
     *
     * @throws Exception something went wrong
     */
    @Test
    public void testBcdEncoding() throws Exception {
        byte[] bytes = new byte[] {(byte) 0x12, (byte) 0x34, (byte) 0x56, (byte) 0x78, (byte) 0x91, (byte) 0x23, (byte) 0x45, (byte) 0x67, (byte) 0x89, (byte) 0x10};

        // Parse the input
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        BigInteger bigIntegerValue = readBufferByteBased.readUnsignedBigInteger("test", bytes.length * 8, WithReaderWriterArgs.WithEncoding("BCD"));

        // Check the expected value matches
        Assertions.assertEquals(new BigInteger("12345678912345678910"), bigIntegerValue);

        // Serialize the value again
        WriteBufferByteBased writeBufferByteBased = new WriteBufferByteBased(bytes.length);
        writeBufferByteBased.writeUnsignedBigInteger("test", bytes.length * 8, bigIntegerValue, WithReaderWriterArgs.WithEncoding("BCD"));

        // Check the output matches the original one
        byte[] serializedBytes = writeBufferByteBased.getBytes();
        Assertions.assertArrayEquals(bytes, serializedBytes);
    }


    /**
     * In VARUDINT encoding is var length unsigned integer encoding.
     *
     * @throws Exception something went wrong
     */
    @Test
    public void testVarudintEncoding() throws Exception {
        // Decimal: 1234567891234567891
        // Binary:  1000100100010000100001111010011000000001000111011011011010011
        // Groups of 7 bits:
        //  0010001  0010001  0000100  0011110  1001100  0000001  0001110  1101101  1010011
        // Var length unsigned int binary:
        // 10010001 10010001 10000100 10011110 11001100 10000001 10001110 11101101 01010011
        // Var length unsigned int hex-bytes:
        // 91       91       84       9E       CC       81       8E       ED       53
        byte[] bytes = new byte[] {(byte) 0x91, (byte) 0x91, (byte) 0x84, (byte) 0x9E, (byte) 0xCC, (byte) 0x81, (byte) 0x8E, (byte) 0xED, (byte) 0x53};

        // Parse the input
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        BigInteger bigIntegerValue = readBufferByteBased.readUnsignedBigInteger("test", bytes.length * 8, WithReaderWriterArgs.WithEncoding("VARUDINT"));

        // Check the expected value matches
        Assertions.assertEquals(new BigInteger("1234567891234567891"), bigIntegerValue);

        // Serialize the value again
        WriteBufferByteBased writeBufferByteBased = new WriteBufferByteBased(bytes.length);
        writeBufferByteBased.writeUnsignedBigInteger("test", bytes.length * 8, bigIntegerValue, WithReaderWriterArgs.WithEncoding("VARUDINT"));

        // Check the output matches the original one
        byte[] serializedBytes = writeBufferByteBased.getBytes();
        Assertions.assertArrayEquals(bytes, serializedBytes);
    }

    /**
     * In default encoding
     *
     * @throws Exception something went wrong
     */
    @Test
    public void testDefaultEncoding() throws Exception {
        byte[] bytes = new byte[] {(byte) 0x11, (byte) 0x22, (byte) 0x10, (byte) 0xF4, (byte) 0xC0, (byte) 0x23, (byte) 0xB6, (byte) 0xD3};

        // Parse the input
        ReadBufferByteBased readBufferByteBased = new ReadBufferByteBased(bytes);
        BigInteger bigIntegerValue = readBufferByteBased.readUnsignedBigInteger("test", 64);

        // Check the expected value matches
        Assertions.assertEquals(new BigInteger("1234567891234567891"), bigIntegerValue);

        // Serialize the value again
        WriteBufferByteBased writeBufferByteBased = new WriteBufferByteBased(bytes.length);
        writeBufferByteBased.writeUnsignedBigInteger("test", 64, bigIntegerValue);

        // Check the output matches the original one
        byte[] serializedBytes = writeBufferByteBased.getBytes();
        Assertions.assertArrayEquals(bytes, serializedBytes);
    }


}
