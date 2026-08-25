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

import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.apache.plc4x.java.opcua.readwrite.Variant;
import org.apache.plc4x.java.opcua.readwrite.VariantInt32;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * A Variant carrying an array announces how many elements follow, and the announcement arrives
 * four bytes ahead of the elements. A server that says Integer.MAX_VALUE and then sends nothing
 * is describing an array that cannot be there, and the reading of it must not reserve room for
 * the claim before testing it.
 */
public class VariantArrayCountLimitTest {

    /** arrayLengthSpecified, no arrayDimensions, VariantType 6 (Int32). */
    private static final byte INT32_ARRAY_HEADER = (byte) 0x86;

    private static ReadBufferByteBased buffer(byte... bytes) {
        return new ReadBufferByteBased(bytes, PayloadConverter.LITTLE_ENDIAN);
    }

    @Test
    void anArrayCountOfIntegerMaxValueWithNoPayloadIsRefused() {
        // The reported payload: the four length bytes set to ff ff ff 7f and no elements.
        BufferException e = assertThrows(BufferException.class, () -> Variant.staticParse(
            buffer(INT32_ARRAY_HEADER, (byte) 0xFF, (byte) 0xFF, (byte) 0xFF, (byte) 0x7F)));
        assertNotNull(e);
    }

    @Test
    void theCountIsNotReservedForBeforeTheElementsAreRead() {
        // Reserving Integer.MAX_VALUE slots is an eight gigabyte array. The point is that this
        // returns a parse failure rather than an OutOfMemoryError, which is an Error and would
        // pass straight through any handler catching Exception.
        assertDoesNotThrow(() -> {
            try {
                Variant.staticParse(
                    buffer(INT32_ARRAY_HEADER, (byte) 0xFF, (byte) 0xFF, (byte) 0xFF, (byte) 0x7F));
                fail("expected the overlarge count to be refused");
            } catch (BufferException expected) {
                // The buffer runs out long before the count does.
            }
        });
    }

    @Test
    void anArrayWhoseCountMatchesItsPayloadStillParses() throws Exception {
        Variant variant = Variant.staticParse(buffer(INT32_ARRAY_HEADER,
            (byte) 0x02, (byte) 0x00, (byte) 0x00, (byte) 0x00,   // two elements
            (byte) 0x07, (byte) 0x00, (byte) 0x00, (byte) 0x00,
            (byte) 0x09, (byte) 0x00, (byte) 0x00, (byte) 0x00));
        assertInstanceOf(VariantInt32.class, variant);
        assertEquals(java.util.List.of(7, 9), ((VariantInt32) variant).getValue());
    }
}
