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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.apache.plc4x.java.eip.readwrite.CIPStructTypeCode;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.junit.jupiter.api.Test;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Encoding of string writes. The payload has to match what the read path parses back: a 2-byte
 * structure handle, a 4-byte length and then the characters, padded out to the fixed size the
 * serializer emits for the type.
 */
class EipStringWriteTest {

    @Test
    void stringWriteMatchesTheStructureTheReadPathParses() {
        byte[] encoded = EipTcpConnection.encodeValue(new PlcSTRING("Hello"), CIPDataTypeCode.STRUCTURED);

        assertEquals(CIPDataTypeCode.STRUCTURED.getSize(), encoded.length, "payload must fill the serialized size");

        ByteBuffer buffer = ByteBuffer.wrap(encoded).order(ByteOrder.LITTLE_ENDIAN);
        assertEquals(CIPStructTypeCode.STRING.getValue(), buffer.getShort(0) & 0xFFFF, "structure handle");
        assertEquals(5, buffer.getInt(2), "length");
        assertEquals("Hello", new String(encoded, 6, 5, StandardCharsets.UTF_8));
    }

    /** The whole point: what we write must read back as the same value. */
    @Test
    void stringWriteRoundTripsThroughTheReadPath() {
        for (String text : new String[]{"", "a", "Hello, world", "0123456789"}) {
            byte[] encoded = EipTcpConnection.encodeValue(new PlcSTRING(text), CIPDataTypeCode.STRUCTURED);
            PlcValue decoded = EipTcpConnection.parsePlcValue(
                EipTag.of("%N40:STRUCTURED"), encoded, CIPDataTypeCode.STRUCTURED);

            assertNotNull(decoded, () -> "did not decode: '" + text + "'");
            assertEquals(text, decoded.getString(), () -> "round trip of '" + text + "'");
        }
    }

    /**
     * A multi-byte character is where counting characters instead of bytes goes wrong: the length
     * field and the number of bytes written have to agree, or the reader sees a truncated string.
     */
    @Test
    void nonAsciiStringIsMeasuredInBytesNotCharacters() {
        String text = "Grüße";  // 5 characters, 7 UTF-8 bytes
        assertEquals(5, text.length());
        assertEquals(7, text.getBytes(StandardCharsets.UTF_8).length);

        byte[] encoded = EipTcpConnection.encodeValue(new PlcSTRING(text), CIPDataTypeCode.STRUCTURED);

        ByteBuffer buffer = ByteBuffer.wrap(encoded).order(ByteOrder.LITTLE_ENDIAN);
        assertEquals(7, buffer.getInt(2), "length must be the byte count");

        PlcValue decoded = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40:STRUCTURED"), encoded, CIPDataTypeCode.STRUCTURED);
        assertEquals(text, decoded.getString());
    }

    /** The longest string that still fits: 88 bytes total minus the 6 byte header. */
    @Test
    void longestFittingStringIsAccepted() {
        String text = "x".repeat(CIPDataTypeCode.STRUCTURED.getSize() - 6);

        byte[] encoded = EipTcpConnection.encodeValue(new PlcSTRING(text), CIPDataTypeCode.STRUCTURED);
        PlcValue decoded = EipTcpConnection.parsePlcValue(
            EipTag.of("%N40:STRUCTURED"), encoded, CIPDataTypeCode.STRUCTURED);

        assertEquals(text, decoded.getString());
    }

    /**
     * A string that does not fit used to overflow the buffer with a BufferOverflowException, or
     * silently write a length that disagreed with the bytes that followed.
     */
    @Test
    void tooLongStringIsRejectedWithAClearError() {
        String text = "x".repeat(CIPDataTypeCode.STRUCTURED.getSize() - 5);

        PlcInvalidTagException thrown = assertThrows(PlcInvalidTagException.class,
            () -> EipTcpConnection.encodeValue(new PlcSTRING(text), CIPDataTypeCode.STRUCTURED));
        assertTrue(thrown.getMessage().contains("83"), () -> "message should name the size: " + thrown.getMessage());
    }

    /**
     * CIPDataTypeCode.STRING has a declared size of 0, so the CipWriteRequest serializer emits no
     * payload for it at all - see the 'count' expression in eip.mspec. That has to be reported
     * rather than silently sending an empty write.
     */
    @Test
    void stringTypeWithNoRoomIsRejected() {
        assertEquals(0, CIPDataTypeCode.STRING.getSize(), "precondition: STRING declares no size");

        assertThrows(PlcInvalidTagException.class,
            () -> EipTcpConnection.encodeValue(new PlcSTRING("Hello"), CIPDataTypeCode.STRING));
    }
}
