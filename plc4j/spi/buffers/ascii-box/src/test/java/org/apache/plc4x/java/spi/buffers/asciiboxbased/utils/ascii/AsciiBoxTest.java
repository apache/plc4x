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
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.ascii;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class AsciiBoxTest {

    @Test
    void testWidth() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "hello", 0);

        assertTrue(box.width() > 0);
    }

    @Test
    void testWidthWithAnsiCodes() {
        // Create a box with ANSI color codes
        AsciiBox box = new AsciiBox("\033[31mred\033[0m");

        // Width should not count ANSI codes
        assertEquals(3, box.width());
    }

    @Test
    void testWidthMultipleLines() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "short\nlongerline\nx", 0);

        // Width should be the maximum line width
        assertTrue(box.width() >= "longerline".length());
    }

    @Test
    void testHeight() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "hello", 0);

        // Box should have at least 3 lines (top border, content, bottom border)
        assertTrue(box.height() >= 3);
    }

    @Test
    void testHeightMultipleLines() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "line1\nline2\nline3", 0);

        // Should have 5 lines (top border, 3 content lines, bottom border)
        assertEquals(5, box.height());
    }

    @Test
    void testLines() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "content", 0);

        String[] lines = box.lines();
        assertNotNull(lines);
        assertTrue(lines.length > 0);
    }

    @Test
    void testGetBoxName() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("myName", "content", 0);

        assertEquals("myName", box.getBoxName());
    }

    @Test
    void testGetBoxNameEmpty() {
        AsciiBox box = new AsciiBox("no name box");

        assertEquals("", box.getBoxName());
    }

    @Test
    void testChangeBoxName() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("oldName", "content", 0);
        AsciiBox renamed = box.changeBoxName("newName");

        assertEquals("newName", renamed.getBoxName());
    }

    @Test
    void testChangeBoxNameNoBorders() {
        AsciiBox box = new AsciiBox("simple text");
        AsciiBox renamed = box.changeBoxName("newName");

        assertEquals("newName", renamed.getBoxName());
    }

    @Test
    void testIsEmptyFalse() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "content", 0);

        assertFalse(box.isEmpty());
    }

    @Test
    void testIsEmptyTrue() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "", 0);

        assertTrue(box.isEmpty());
    }

    @Test
    void testIsEmptyNoBorders() {
        AsciiBox empty = new AsciiBox("");
        AsciiBox blank = new AsciiBox("   ");
        AsciiBox notEmpty = new AsciiBox("text");

        assertTrue(empty.isEmpty());
        assertTrue(blank.isEmpty());
        assertFalse(notEmpty.isEmpty());
    }

    @Test
    void testToString() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "hello", 0);

        String result = box.toString();
        assertNotNull(result);
        assertTrue(result.contains("test"));
        assertTrue(result.contains("hello"));
    }

    @Test
    void testEquals() {
        AsciiBox box1 = new AsciiBox("same data");
        AsciiBox box2 = new AsciiBox("same data");
        AsciiBox box3 = new AsciiBox("different data");

        assertEquals(box1, box2);
        assertNotEquals(box1, box3);
    }

    @Test
    void testEqualsItself() {
        AsciiBox box = new AsciiBox("data");

        assertEquals(box, box);
    }

    @Test
    void testEqualsNull() {
        AsciiBox box = new AsciiBox("data");

        assertNotEquals(null, box);
    }

    @Test
    void testEqualsDifferentType() {
        AsciiBox box = new AsciiBox("data");

        assertNotEquals("data", box);
    }

    @Test
    void testHashCode() {
        AsciiBox box1 = new AsciiBox("same data");
        AsciiBox box2 = new AsciiBox("same data");

        assertEquals(box1.hashCode(), box2.hashCode());
    }

    @Test
    void testConstructorWithNullThrows() {
        assertThrows(NullPointerException.class, () -> new AsciiBox(null));
    }

    @Test
    void testChangeBoxNameWithLongerName() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("a", "content", 0);
        AsciiBox renamed = box.changeBoxName("muchLongerName");

        assertEquals("muchLongerName", renamed.getBoxName());
        // Box should expand to fit new name
        assertTrue(renamed.width() >= "muchLongerName".length());
    }

    @Test
    void testChangeBoxNameWithShorterName() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("veryLongName", "content", 0);
        AsciiBox renamed = box.changeBoxName("x");

        assertEquals("x", renamed.getBoxName());
    }

    @Test
    void testWidthEmptyBox() {
        AsciiBox box = new AsciiBox("");

        assertEquals(0, box.width());
    }

    @Test
    void testHeightEmptyBox() {
        AsciiBox box = new AsciiBox("");

        assertEquals(1, box.height()); // Split on empty string returns array of length 1
    }

    @Test
    void testLinesEmptyBox() {
        AsciiBox box = new AsciiBox("");

        String[] lines = box.lines();
        assertEquals(1, lines.length);
        assertEquals("", lines[0]);
    }

    @Test
    void testBoxWithSpecialCharacters() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "special: @#$%^&*()", 0);

        String result = box.toString();
        assertTrue(result.contains("@#$%^&*()"));
    }

    @Test
    void testBoxWithUnicodeContent() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("unicode", "Hello 世界", 0);

        String result = box.toString();
        assertTrue(result.contains("世界"));
    }
}
