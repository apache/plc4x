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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

public class VintIoTest {

    @ParameterizedTest
    @ValueSource(longs = {0, 1, 23, 42, 126, 127})
    void oneByteTest(long value) throws  Exception {
        roundtrip(value, 1);
    }

    @ParameterizedTest
    @ValueSource(longs = {128, 1234, 12345})
    void twoByteTest(long value) throws  Exception {
        roundtrip(value, 2);
    }

    @ParameterizedTest
    @ValueSource(longs = {123456})
    void threeByteTest(long value) throws  Exception {
        roundtrip(value, 3);
    }

    private void roundtrip(long value, int expectedNumBytes) throws  Exception {
        Assertions.assertEquals(expectedNumBytes, VintIo.getLengthInBytes(value));
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(10);
        VintIo.serialize(value, writeBuffer);
        Assertions.assertEquals(expectedNumBytes, writeBuffer.getPos());
        ReadBuffer readBuffer = new ReadBufferByteBased(writeBuffer.getBytes());
        Assertions.assertEquals(value, VintIo.parse(readBuffer));
        // For the first bytes the top most bit must be set to 1.
        for(int i = 0; i < (expectedNumBytes - 1); i++) {
            Assertions.assertEquals(0x80, writeBuffer.getBytes()[i] & 0x80);
        }
        // For the last byte the top most bit must be set to 0.
        Assertions.assertEquals(0x00, writeBuffer.getBytes()[expectedNumBytes - 1] & 0x80);
    }

}
