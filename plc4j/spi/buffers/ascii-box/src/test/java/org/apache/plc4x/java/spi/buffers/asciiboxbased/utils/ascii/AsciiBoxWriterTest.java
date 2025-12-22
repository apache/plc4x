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

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class AsciiBoxWriterTest {

    @Test
    void testDefaultWriter() {
        assertNotNull(AsciiBoxWriter.DEFAULT);
    }

    @Test
    void testLightWriter() {
        assertNotNull(AsciiBoxWriter.LIGHT);
    }

    @Test
    void testBoxStringSimple() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "hello", 0);

        assertNotNull(box);
        String result = box.toString();
        assertTrue(result.contains("test"));
        assertTrue(result.contains("hello"));
    }

    @Test
    void testBoxStringWithCharWidth() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "hi", 20);

        assertNotNull(box);
        assertTrue(box.width() >= 20);
    }

    @Test
    void testBoxStringWithMultipleLines() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "line1\nline2\nline3", 0);

        assertNotNull(box);
        String result = box.toString();
        assertTrue(result.contains("line1"));
        assertTrue(result.contains("line2"));
        assertTrue(result.contains("line3"));
    }

    @Test
    void testBoxStringConvertsLineEndings() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "line1\r\nline2", 0);

        assertNotNull(box);
        // Should not contain \r
        assertFalse(box.toString().contains("\r"));
    }

    @Test
    void testBoxStringConvertsTabs() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "col1\tcol2", 0);

        assertNotNull(box);
        // Tab should be converted to spaces
        assertTrue(box.toString().contains("col1  col2"));
    }

    @Test
    void testBoxBox() {
        AsciiBox innerBox = AsciiBoxWriter.DEFAULT.boxString("inner", "content", 0);
        AsciiBox outerBox = AsciiBoxWriter.DEFAULT.boxBox("outer", innerBox, 0);

        assertNotNull(outerBox);
        String result = outerBox.toString();
        assertTrue(result.contains("outer"));
        assertTrue(result.contains("inner"));
        assertTrue(result.contains("content"));
    }

    @Test
    void testBoxSideBySide() {
        AsciiBox box1 = AsciiBoxWriter.DEFAULT.boxString("left", "L", 0);
        AsciiBox box2 = AsciiBoxWriter.DEFAULT.boxString("right", "R", 0);

        AsciiBox combined = AsciiBoxWriter.DEFAULT.boxSideBySide(box1, box2);

        assertNotNull(combined);
        String result = combined.toString();
        assertTrue(result.contains("left"));
        assertTrue(result.contains("right"));
    }

    @Test
    void testBoxSideBySideDifferentHeights() {
        AsciiBox box1 = AsciiBoxWriter.DEFAULT.boxString("short", "A", 0);
        AsciiBox box2 = AsciiBoxWriter.DEFAULT.boxString("tall", "B\nC\nD", 0);

        AsciiBox combined = AsciiBoxWriter.DEFAULT.boxSideBySide(box1, box2);

        assertNotNull(combined);
        // Should have height of the taller box
        assertTrue(combined.height() >= box2.height());
    }

    @Test
    void testBoxBelowBox() {
        AsciiBox box1 = AsciiBoxWriter.DEFAULT.boxString("top", "T", 0);
        AsciiBox box2 = AsciiBoxWriter.DEFAULT.boxString("bottom", "B", 0);

        AsciiBox combined = AsciiBoxWriter.DEFAULT.boxBelowBox(box1, box2);

        assertNotNull(combined);
        String[] lines = combined.lines();
        // Should have lines from both boxes
        assertTrue(lines.length > box1.height());
    }

    @Test
    void testBoxBelowBoxDifferentWidths() {
        AsciiBox box1 = AsciiBoxWriter.DEFAULT.boxString("narrow", "X", 0);
        AsciiBox box2 = AsciiBoxWriter.DEFAULT.boxString("wide", "XXXXX XXXXX", 0);

        AsciiBox combined = AsciiBoxWriter.DEFAULT.boxBelowBox(box1, box2);

        assertNotNull(combined);
        // All lines should have same width
        String[] lines = combined.lines();
        int firstWidth = lines[0].length();
        for (String line : lines) {
            assertEquals(firstWidth, line.length());
        }
    }

    @Test
    void testAlignBoxesEmpty() {
        AsciiBox result = AsciiBoxWriter.DEFAULT.alignBoxes(Collections.emptyList(), 80);

        assertNotNull(result);
        assertEquals("", result.toString());
    }

    @Test
    void testAlignBoxesSingle() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("single", "content", 0);
        AsciiBox result = AsciiBoxWriter.DEFAULT.alignBoxes(Collections.singletonList(box), 80);

        assertNotNull(result);
        assertEquals(box.toString(), result.toString());
    }

    @Test
    void testAlignBoxesMultiple() {
        List<AsciiBox> boxes = Arrays.asList(
            AsciiBoxWriter.DEFAULT.boxString("box1", "A", 0),
            AsciiBoxWriter.DEFAULT.boxString("box2", "B", 0),
            AsciiBoxWriter.DEFAULT.boxString("box3", "C", 0)
        );

        AsciiBox result = AsciiBoxWriter.DEFAULT.alignBoxes(boxes, 200);

        assertNotNull(result);
        String resultStr = result.toString();
        assertTrue(resultStr.contains("box1"));
        assertTrue(resultStr.contains("box2"));
        assertTrue(resultStr.contains("box3"));
    }

    @Test
    void testAlignBoxesNarrowWidth() {
        List<AsciiBox> boxes = Arrays.asList(
            AsciiBoxWriter.DEFAULT.boxString("box1", "AAAAAAA", 0),
            AsciiBoxWriter.DEFAULT.boxString("box2", "BBBBBBB", 0)
        );

        // Narrow width forces boxes onto separate lines
        AsciiBox result = AsciiBoxWriter.DEFAULT.alignBoxes(boxes, 15);

        assertNotNull(result);
        // Should have more lines due to stacking
        assertTrue(result.height() > boxes.get(0).height());
    }

    @Test
    void testHasBordersTrue() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "content", 0);

        assertTrue(AsciiBoxWriter.DEFAULT.hasBorders(box));
    }

    @Test
    void testHasBordersEmpty() {
        AsciiBox box = new AsciiBox("");

        assertFalse(AsciiBoxWriter.DEFAULT.hasBorders(box));
    }

    @Test
    void testHasBordersNoBorder() {
        // Create a box without proper borders
        AsciiBox box = new AsciiBox("just text");

        assertFalse(AsciiBoxWriter.DEFAULT.hasBorders(box));
    }

    @Test
    void testUnwrap() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "content", 0);
        AsciiBox unwrapped = AsciiBoxWriter.DEFAULT.unwrap(box);

        assertNotNull(unwrapped);
        assertEquals("content", unwrapped.toString());
    }

    @Test
    void testUnwrapNoBorders() {
        AsciiBox box = new AsciiBox("no borders");
        AsciiBox unwrapped = AsciiBoxWriter.DEFAULT.unwrap(box);

        assertEquals(box, unwrapped);
    }

    @Test
    void testMergeHorizontalEmpty() {
        AsciiBox result = AsciiBoxWriter.DEFAULT.mergeHorizontal(Collections.emptyList());

        assertEquals("", result.toString());
    }

    @Test
    void testMergeHorizontalSingle() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "A", 0);
        AsciiBox result = AsciiBoxWriter.DEFAULT.mergeHorizontal(Collections.singletonList(box));

        assertEquals(box, result);
    }

    @Test
    void testMergeHorizontalTwo() {
        List<AsciiBox> boxes = Arrays.asList(
            AsciiBoxWriter.DEFAULT.boxString("box1", "A", 0),
            AsciiBoxWriter.DEFAULT.boxString("box2", "B", 0)
        );

        AsciiBox result = AsciiBoxWriter.DEFAULT.mergeHorizontal(boxes);

        assertNotNull(result);
        String resultStr = result.toString();
        assertTrue(resultStr.contains("box1"));
        assertTrue(resultStr.contains("box2"));
    }

    @Test
    void testMergeHorizontalThree() {
        List<AsciiBox> boxes = Arrays.asList(
            AsciiBoxWriter.DEFAULT.boxString("box1", "A", 0),
            AsciiBoxWriter.DEFAULT.boxString("box2", "B", 0),
            AsciiBoxWriter.DEFAULT.boxString("box3", "C", 0)
        );

        AsciiBox result = AsciiBoxWriter.DEFAULT.mergeHorizontal(boxes);

        assertNotNull(result);
        String resultStr = result.toString();
        assertTrue(resultStr.contains("box1"));
        assertTrue(resultStr.contains("box2"));
        assertTrue(resultStr.contains("box3"));
    }

    @Test
    void testExpandBox() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "A", 0);
        int originalWidth = box.width();

        AsciiBox expanded = AsciiBoxWriter.DEFAULT.expandBox(box, originalWidth + 10);

        assertEquals(originalWidth + 10, expanded.width());
    }

    @Test
    void testExpandBoxNoExpansionNeeded() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "A", 0);
        int originalWidth = box.width();

        AsciiBox expanded = AsciiBoxWriter.DEFAULT.expandBox(box, originalWidth - 5);

        assertEquals(originalWidth, expanded.width());
    }

    @Test
    void testCustomBoxWriter() {
        AsciiBoxWriter custom = new AsciiBoxWriter("+", "+", "-", "|", "+", "+");
        AsciiBox box = custom.boxString("test", "content", 0);

        String result = box.toString();
        assertTrue(result.contains("+"));
        assertTrue(result.contains("-"));
        assertTrue(result.contains("|"));
    }

    @Test
    void testBoxStringWithLongContent() {
        String longContent = "This is a very long content string that should cause the box to expand";
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", longContent, 10);

        // Box should expand to fit content
        assertTrue(box.width() > 10);
    }

    @Test
    void testBoxStringEmptyContent() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("test", "", 0);

        assertNotNull(box);
        assertTrue(box.toString().contains("test"));
    }

    @Test
    void testBoxStringEmptyName() {
        AsciiBox box = AsciiBoxWriter.DEFAULT.boxString("", "content", 0);

        assertNotNull(box);
        assertTrue(box.toString().contains("content"));
    }

    @Test
    void testLightWriterUsesLightChars() {
        AsciiBox box = AsciiBoxWriter.LIGHT.boxString("test", "content", 0);

        String result = box.toString();
        // Light writer uses different unicode characters
        assertTrue(result.contains("╭") || result.contains("┄"));
    }

    @Test
    void testAlignBoxesWithOverflow() {
        // Create boxes that exceed desired width
        List<AsciiBox> boxes = new ArrayList<>();
        for (int i = 0; i < 5; i++) {
            boxes.add(AsciiBoxWriter.DEFAULT.boxString("box" + i, "content" + i, 0));
        }

        AsciiBox result = AsciiBoxWriter.DEFAULT.alignBoxes(boxes, 50);

        assertNotNull(result);
        // Should contain all boxes
        String resultStr = result.toString();
        assertTrue(resultStr.contains("box0"));
        assertTrue(resultStr.contains("box4"));
    }
}
