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
package org.apache.plc4x.java.utils.testutils.utils.hex;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class HexDiffTest {

    @Test
    void testHexDiffWithIdenticalBytes() {
        String expected = "0102030405";
        String actual = "0102030405";

        String diff = HexDiff.hexDiff(expected, actual);
        assertNotNull(diff);
        assertTrue(diff.contains("expected") || diff.contains("actual"));
    }

    @Test
    void testHexDiffWithDifferentBytes() {
        String expected = "0102030405";
        String actual = "0102FF0405";

        String diff = HexDiff.hexDiff(expected, actual);
        assertNotNull(diff);
        assertTrue(diff.contains("expected") || diff.contains("actual"));
    }

    @Test
    void testHexDiffWithDifferentLengths() {
        String expected = "010203";
        String actual = "0102030405";

        String diff = HexDiff.hexDiff(expected, actual);
        assertNotNull(diff);
    }

    @Test
    void testHexDiffWithInvalidHex() {
        String expected = "ZZZZ";
        String actual = "0102";

        String diff = HexDiff.hexDiff(expected, actual);
        assertNotNull(diff);
        assertTrue(diff.contains("Failed to decode") || diff.contains("failed"));
    }

    @Test
    void testDiffHexWithByteArrays() {
        byte[] expected = new byte[]{0x01, 0x02, 0x03};
        byte[] actual = new byte[]{0x01, 0x02, 0x03};

        var result = HexDiff.diffHex(expected, actual);
        assertNotNull(result);
    }

    @Test
    void testDiffHexWithDifferentByteArrays() {
        byte[] expected = new byte[]{0x01, 0x02, 0x03};
        byte[] actual = new byte[]{0x01, (byte) 0xFF, 0x03};

        var result = HexDiff.diffHex(expected, actual);
        assertNotNull(result);
    }
}
