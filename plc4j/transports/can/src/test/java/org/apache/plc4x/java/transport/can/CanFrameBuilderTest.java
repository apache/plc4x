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

class CanFrameBuilderTest {

    @Test
    void validStandardFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x123)
                .data(new byte[]{0x01, 0x02, 0x03})
                .build();

        assertNotNull(frame);
        assertEquals(0x123, frame.getIdentifier());
        assertFalse(frame.isExtended());
    }

    @Test
    void validExtendedFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x12345678)
                .data(new byte[]{0x01})
                .extended(true)
                .build();

        assertNotNull(frame);
        assertEquals(0x12345678, frame.getIdentifier());
        assertTrue(frame.isExtended());
    }

    @Test
    void dataExceeding8BytesRejected() {
        CanFrameBuilder builder = CanFrame.builder()
                .identifier(0x100)
                .data(new byte[9]);

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, builder::build);
        assertTrue(ex.getMessage().contains("8"), "Message should mention 8 byte limit");
    }

    @Test
    void standardIdExceeding0x7FFRejected() {
        CanFrameBuilder builder = CanFrame.builder()
                .identifier(0x800)
                .data(new byte[]{0x01});

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, builder::build);
        assertTrue(ex.getMessage().contains("7FF"), "Message should mention standard ID limit");
    }

    @Test
    void extendedIdExceeding0x1FFFFFFFRejected() {
        CanFrameBuilder builder = CanFrame.builder()
                .identifier(0x20000000)
                .data(new byte[]{0x01})
                .extended(true);

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, builder::build);
        assertTrue(ex.getMessage().contains("1FFFFFFF"), "Message should mention extended ID limit");
    }

    @Test
    void negativeIdRejected() {
        CanFrameBuilder builder = CanFrame.builder()
                .identifier(-1)
                .data(new byte[]{0x01});

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, builder::build);
        assertTrue(ex.getMessage().contains("negative"), "Message should mention negative ID");
    }

    @Test
    void nullDataRejected() {
        assertThrows(IllegalArgumentException.class, () ->
                CanFrame.builder().identifier(0x100).data(null));
    }

    @Test
    void rtrFrameWithEmptyData() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .rtr(true)
                .build();

        assertNotNull(frame);
        assertTrue(frame.isRtr());
        assertEquals(0, frame.getDataLength());
    }

    @Test
    void rtrFrameWithDataRejected() {
        CanFrameBuilder builder = CanFrame.builder()
                .identifier(0x100)
                .data(new byte[]{0x01})
                .rtr(true);

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, builder::build);
        assertTrue(ex.getMessage().contains("RTR"), "Message should mention RTR");
    }

    @Test
    void identifierNotSetRejected() {
        CanFrameBuilder builder = CanFrame.builder()
                .data(new byte[]{0x01});

        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class, builder::build);
        assertTrue(ex.getMessage().contains("Identifier"), "Message should mention identifier");
    }

    @Test
    void builderDefaultsToEmptyData() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .build();

        assertEquals(0, frame.getDataLength());
    }

    @Test
    void builderDefaultsToStandardFrame() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .build();

        assertFalse(frame.isExtended());
    }

    @Test
    void builderDefaultsToNonRtr() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .build();

        assertFalse(frame.isRtr());
    }

    @Test
    void builderReuse() {
        // Each build() call should create an independent frame
        CanFrameBuilder builder = CanFrame.builder().identifier(0x100);

        CanFrame frame1 = builder.build();
        // Modify builder after first build
        builder.data(new byte[]{0x01, 0x02});
        CanFrame frame2 = builder.build();

        assertNotEquals(frame1, frame2);
        assertEquals(0, frame1.getDataLength());
        assertEquals(2, frame2.getDataLength());
    }

    @Test
    void zeroIdentifierValid() {
        CanFrame frame = CanFrame.builder()
                .identifier(0)
                .build();

        assertEquals(0, frame.getIdentifier());
    }

    @Test
    void maxStandardIdValid() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x7FF)
                .build();

        assertEquals(0x7FF, frame.getIdentifier());
    }

    @Test
    void maxExtendedIdValid() {
        CanFrame frame = CanFrame.builder()
                .identifier(0x1FFFFFFF)
                .extended(true)
                .build();

        assertEquals(0x1FFFFFFF, frame.getIdentifier());
    }

    @Test
    void exactlyMaxDataLength() {
        byte[] data = new byte[8];
        CanFrame frame = CanFrame.builder()
                .identifier(0x100)
                .data(data)
                .build();

        assertEquals(8, frame.getDataLength());
    }
}
