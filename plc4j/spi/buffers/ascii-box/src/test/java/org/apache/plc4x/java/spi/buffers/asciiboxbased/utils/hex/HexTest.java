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
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.hex;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class HexTest {

    @Test
    void testDumpNull() {
        String result = Hex.dump(null);
        assertEquals("", result);
    }

    @Test
    void testDumpEmpty() {
        String result = Hex.dump(new byte[0]);
        assertEquals("", result);
    }

    @Test
    void testDumpSingleByte() {
        byte[] data = new byte[]{0x41}; // 'A'
        String result = Hex.dump(data);

        assertNotNull(result);
        assertTrue(result.contains("41"));
        // String representation is quoted
        assertTrue(result.contains("'"));
    }

    @Test
    void testDumpMultipleBytes() {
        byte[] data = new byte[]{0x48, 0x65, 0x6C, 0x6C, 0x6F}; // "Hello"
        String result = Hex.dump(data);

        assertNotNull(result);
        assertTrue(result.contains("48"));
        assertTrue(result.contains("65"));
        assertTrue(result.contains("6c"));
        // Contains string representation with quotes
        assertTrue(result.contains("'"));
    }

    @Test
    void testDumpWithNonPrintableChars() {
        byte[] data = new byte[]{0x00, 0x1F, 0x7F}; // Non-printable chars
        String result = Hex.dump(data);

        assertNotNull(result);
        assertTrue(result.contains("00"));
        assertTrue(result.contains("1f"));
        assertTrue(result.contains("7f"));
        // Non-printable chars should be masked as '.'
        assertTrue(result.contains("."));
    }

    @Test
    void testDumpWithHighlights() {
        byte[] data = new byte[]{0x41, 0x42, 0x43}; // "ABC"
        String result = Hex.dump(data, Hex.DefaultWidth, 1);

        // Should contain ANSI escape codes for highlight
        assertTrue(result.contains("\033[0;31m"));
        assertTrue(result.contains("\033[0m"));
    }

    @Test
    void testDumpWithCustomWidth() {
        byte[] data = new byte[]{0x41, 0x42, 0x43, 0x44, 0x45};
        String result = Hex.dump(data, 20);

        assertNotNull(result);
        assertTrue(result.contains("41"));
    }

    @Test
    void testDumpLargeData() {
        byte[] data = new byte[100];
        for (int i = 0; i < data.length; i++) {
            data[i] = (byte) (i % 256);
        }

        String result = Hex.dump(data);
        assertNotNull(result);
        // Should have multiple lines
        assertTrue(result.contains("\n"));
    }

    @Test
    void testDumpWithMinimalWidth() {
        byte[] data = new byte[]{0x41, 0x42};
        // Very small width should still work
        String result = Hex.dump(data, 10);

        assertNotNull(result);
        assertTrue(result.contains("41"));
    }

    @Test
    void testConstants() {
        assertEquals(46, Hex.DefaultWidth);
        assertEquals(2, Hex.boxLineOverheat);
        assertEquals(1, Hex.blankWidth);
        assertEquals(3, Hex.byteWidth);
        assertEquals(1, Hex.pipeWidth);
    }

    @Test
    void testDumpWithDebugOff() {
        boolean originalDebug = Hex.DebugHex;
        try {
            Hex.DebugHex = false;
            byte[] data = new byte[]{0x41, 0x42, 0x43};
            String result = Hex.dump(data);
            assertNotNull(result);
        } finally {
            Hex.DebugHex = originalDebug;
        }
    }

    @Test
    void testDumpWithDebugOn() {
        boolean originalDebug = Hex.DebugHex;
        try {
            Hex.DebugHex = true;
            byte[] data = new byte[]{0x41, 0x42, 0x43};
            String result = Hex.dump(data);
            assertNotNull(result);
        } finally {
            Hex.DebugHex = originalDebug;
        }
    }

    @Test
    void testDumpPreservesOriginalArray() {
        byte[] data = new byte[]{0x41, 0x42, 0x43};
        byte[] originalCopy = data.clone();

        Hex.dump(data);

        assertArrayEquals(originalCopy, data);
    }

    @Test
    void testDumpWithLargeArrayIndex() {
        // Create data large enough to have multi-digit index
        byte[] data = new byte[1000];
        for (int i = 0; i < data.length; i++) {
            data[i] = (byte) (i % 256);
        }

        String result = Hex.dump(data);
        assertNotNull(result);
        // Should contain index like "000|" format
        assertTrue(result.contains("|"));
    }

    @Test
    void testDumpPrintableRange() {
        // Test characters at boundaries of printable range (32-126)
        byte[] data = new byte[]{31, 32, 126, 127};
        String result = Hex.dump(data);

        // 31 and 127 should be masked, 32 (' ') and 126 ('~') should be visible
        assertTrue(result.contains("1f"));
        assertTrue(result.contains("20"));
        assertTrue(result.contains("7e"));
        assertTrue(result.contains("7f"));
    }
}
