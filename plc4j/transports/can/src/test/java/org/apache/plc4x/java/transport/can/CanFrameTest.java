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
package org.apache.plc4x.java.transport.can;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class CanFrameTest {

    @Test
    void standardFrameConstruction() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x123)
                .data(new byte[]{0x01, 0x02, 0x03})
                .build();

        assertEquals(0x123, frame.getIdentifier());
        assertArrayEquals(new byte[]{0x01, 0x02, 0x03}, frame.getData());
        assertEquals(3, frame.getDataLength());
        assertFalse(frame.isExtended());
        assertFalse(frame.isRtr());
    }

    @Test
    void extendedFrameConstruction() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x12345678)
                .data(new byte[]{0x01})
                .extended(true)
                .build();

        assertEquals(0x12345678, frame.getIdentifier());
        assertTrue(frame.isExtended());
        assertFalse(frame.isRtr());
    }

    @Test
    void rtrFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .rtr(true)
                .build();

        assertTrue(frame.isRtr());
        assertEquals(0, frame.getDataLength());
    }

    @Test
    void getDataReturnsDefensiveCopy() {
        byte[] original = {0x01, 0x02, 0x03};
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .data(original)
                .build();

        byte[] retrieved = frame.getData();
        retrieved[0] = (byte) 0xFF;

        // Original frame data should not be affected
        assertArrayEquals(new byte[]{0x01, 0x02, 0x03}, frame.getData());
    }

    @Test
    void constructorDefensivelyCopiesData() {
        byte[] original = {0x01, 0x02};
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .data(original)
                .build();

        // Modify original after construction
        original[0] = (byte) 0xFF;

        // Frame data should not be affected
        assertArrayEquals(new byte[]{0x01, 0x02}, frame.getData());
    }

    @Test
    void equalsAndHashCodeContract() {
        CanFrame frame1 = CanFrame.builder().identifier(0x100).data(new byte[]{0x01}).build();
        CanFrame frame2 = CanFrame.builder().identifier(0x100).data(new byte[]{0x01}).build();
        CanFrame frame3 = CanFrame.builder().identifier(0x200).data(new byte[]{0x01}).build();

        // Reflexive
        assertEquals(frame1, frame1);

        // Symmetric
        assertEquals(frame1, frame2);
        assertEquals(frame2, frame1);

        // Consistent hashCode
        assertEquals(frame1.hashCode(), frame2.hashCode());

        // Not equal with different identifier
        assertNotEquals(frame1, frame3);

        // Not equal to null
        assertNotEquals(null, frame1);

        // Not equal to different type
        assertNotEquals("not a frame", frame1);
    }

    @Test
    void equalsDiffersByExtendedFlag() {
        CanFrame standard = CanFrame.builder().identifier(0x100).data(new byte[]{0x01}).build();
        CanFrame extended = CanFrame.builder().identifier(0x100).data(new byte[]{0x01}).extended(true).build();

        assertNotEquals(standard, extended);
    }

    @Test
    void equalsDiffersByData() {
        CanFrame frame1 = CanFrame.builder().identifier(0x100).data(new byte[]{0x01}).build();
        CanFrame frame2 = CanFrame.builder().identifier(0x100).data(new byte[]{0x02}).build();

        assertNotEquals(frame1, frame2);
    }

    @Test
    void equalsDiffersByRtrFlag() {
        CanFrame nonRtr = CanFrame.builder().identifier(0x100).build();
        CanFrame rtr = CanFrame.builder().identifier(0x100).rtr(true).build();

        assertNotEquals(nonRtr, rtr);
    }

    @Test
    void toStringContainsRelevantInfo() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x123)
                .data(new byte[]{0x01, 0x02})
                .build();

        String str = frame.toString();
        assertTrue(str.contains("123"), "Should contain identifier");
        assertTrue(str.contains("01 02"), "Should contain data hex");
        assertTrue(str.contains("len=2"), "Should contain length");
    }

    @Test
    void toStringExtendedFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x12345678)
                .data(new byte[]{})
                .extended(true)
                .build();

        String str = frame.toString();
        assertTrue(str.contains("EXT"), "Should indicate extended frame");
        assertTrue(str.contains("12345678"), "Should contain full extended ID");
    }

    @Test
    void toStringRtrFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .rtr(true)
                .build();

        String str = frame.toString();
        assertTrue(str.contains("RTR"), "Should indicate RTR frame");
    }

    @Test
    void emptyDataFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x000)
                .build();

        assertEquals(0, frame.getDataLength());
        assertArrayEquals(new byte[0], frame.getData());
    }

    @Test
    void maxDataLengthFrame() {
        byte[] maxData = new byte[]{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08};
        CanFrame frame = CanFrame.builder()
                .identifier(0x7FF)
                .data(maxData)
                .build();

        assertEquals(8, frame.getDataLength());
        assertArrayEquals(maxData, frame.getData());
    }

    @Test
    void zeroIdentifier() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x000)
                .data(new byte[]{0x01})
                .build();

        assertEquals(0, frame.getIdentifier());
    }

    @Test
    void maxStandardIdentifier() {
        CanFrame frame = CanFrame.builder()
                .identifier(CanFrame.MAX_STANDARD_ID)
                .data(new byte[]{0x01})
                .build();

        assertEquals(0x7FF, frame.getIdentifier());
    }

    @Test
    void maxExtendedIdentifier() {
        CanFrame frame = CanFrame.builder()
                .identifier(CanFrame.MAX_EXTENDED_ID)
                .data(new byte[]{0x01})
                .extended(true)
                .build();

        assertEquals(0x1FFFFFFF, frame.getIdentifier());
    }

    @Test
    void constants() {
        assertEquals(0x7FF, CanFrame.MAX_STANDARD_ID);
        assertEquals(0x1FFFFFFF, CanFrame.MAX_EXTENDED_ID);
        assertEquals(8, CanFrame.MAX_DATA_LENGTH);
    }
}
