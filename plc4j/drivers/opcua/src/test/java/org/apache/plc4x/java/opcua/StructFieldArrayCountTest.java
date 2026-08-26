/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.opcua;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.junit.jupiter.api.Assertions.*;

/**
 * An array field in a struct body is prefixed by its element count, which is the server's claim
 * about what follows rather than a fact about what is there. Reserving room for the claim before
 * testing it hands over the memory to whoever wrote the four bytes.
 */
public class StructFieldArrayCountTest {

    /** An array field of the given built-in type; rank 1 means "one dimension", i.e. an array. */
    private static StructureField arrayFieldOfType(int builtInId) {
        return new StructureField(
            new PascalString("values"), null,
            new NodeId(new NodeIdTwoByte((short) builtInId)),
            1, Collections.emptyList(), 0L, false);
    }

    private static ReadBufferByteBased body(byte[] bytes) {
        return new ReadBufferByteBased(bytes, PayloadConverter.LITTLE_ENDIAN);
    }

    @Test
    void aCountLargerThanTheBodyIsRefusedRatherThanReservedFor() {
        // Four bytes saying "two billion elements", followed by two bytes.
        byte[] raw = {0x00, 0x00, 0x00, 0x70, 0x01, 0x02};
        PlcRuntimeException e = assertThrows(PlcRuntimeException.class,
            () -> OpcuaConnection.decodeFieldForTest(body(raw), arrayFieldOfType(3)));
        assertTrue(e.getMessage().contains("bytes remain"), e.getMessage());
    }

    @Test
    void theCountIsCheckedBeforeAnythingIsAllocatedForIt() {
        // Integer.MAX_VALUE elements would be an eight gigabyte array on a six byte body. The
        // point of the check is that this returns rather than trying.
        byte[] raw = {(byte) 0xFF, (byte) 0xFF, (byte) 0xFF, 0x7F, 0x01, 0x02};
        assertThrows(PlcRuntimeException.class,
            () -> OpcuaConnection.decodeFieldForTest(body(raw), arrayFieldOfType(3)));
    }

    @Test
    void aCountMatchingTheBodyStillDecodes() {
        byte[] raw = {0x03, 0x00, 0x00, 0x00, 0x07, 0x08, 0x09};
        PlcValue value = assertDoesNotThrow(
            () -> OpcuaConnection.decodeFieldForTest(body(raw), arrayFieldOfType(3)));
        assertEquals(3, value.getList().size());
        assertEquals(7, value.getList().get(0).getInt());
        assertEquals(9, value.getList().get(2).getInt());
    }

    @Test
    void aNegativeCountIsStillTheEmptyArray() {
        byte[] raw = {(byte) 0xFF, (byte) 0xFF, (byte) 0xFF, (byte) 0xFF};
        PlcValue value = assertDoesNotThrow(
            () -> OpcuaConnection.decodeFieldForTest(body(raw), arrayFieldOfType(3)));
        assertTrue(value.getList().isEmpty());
    }
}
