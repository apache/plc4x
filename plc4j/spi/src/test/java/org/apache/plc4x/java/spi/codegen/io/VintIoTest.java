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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.junit.jupiter.params.provider.ValueSource;

import java.util.stream.Stream;

public class VintIoTest {

    private static Stream<Arguments> inputData() {
        return Stream.of(
            Arguments.of(0,      new byte[]{0x00}),
            Arguments.of(1,      new byte[]{0x01}),
            Arguments.of(23,     new byte[]{0x17}),
            Arguments.of(42,     new byte[]{0x2A}),
            Arguments.of(126,    new byte[]{0x7E}),
            Arguments.of(127,    new byte[]{0x7F}),
            Arguments.of(128,    new byte[]{(byte) 0x81, 0x00}),
            Arguments.of(1234,   new byte[]{(byte) 0x89, 0x52}),
            Arguments.of(12345,  new byte[]{(byte) 0xE0, 0x39}),
            Arguments.of(123456, new byte[]{(byte) 0x87, (byte) 0xC4, 0x40})
        );
    }

    @ParameterizedTest
    @MethodSource("inputData")
    void oneByteTest(long value, byte[] expectedBytes) throws  Exception {
        Assertions.assertEquals(expectedBytes.length, VintIo.getLengthInBytes(value));
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(10);
        VintIo.serialize(value, writeBuffer);
        final byte[] serializedBytes = writeBuffer.getBytes();
        Assertions.assertEquals(expectedBytes.length, serializedBytes.length);
        // Compare the content of the byte buffer
        for(int i = 0; i < expectedBytes.length; i++) {
            Assertions.assertEquals(expectedBytes[i], serializedBytes[i]);
        }
        ReadBuffer readBuffer = new ReadBufferByteBased(serializedBytes);
        Assertions.assertEquals(value, VintIo.parse(readBuffer));
    }

}
