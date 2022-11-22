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

package org.apache.plc4x.java.spi.utils.ascii;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collection;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class AsciiBoxWriterTest {

    @Nested
    class GetBoxName {
        @Test
        void simpleName() {
            assertEquals("someName", AsciiBoxWriter.DEFAULT.boxString("someName", "some content", 0).getBoxName());
        }

        @Test
        void noName() {
            assertEquals("", AsciiBoxWriter.DEFAULT.boxString("", "some content", 0).getBoxName());
        }

        @Test
        void longName() {
            assertEquals("veryLongName12_13", AsciiBoxWriter.DEFAULT.boxString("veryLongName12_13", "some content", 0).getBoxName());
        }

        @Test
        void nameWithSpacesAndSlashes() {
            assertEquals("payload / Message / Concrete Message", AsciiBoxWriter.DEFAULT.boxString("payload / Message / Concrete Message", "some content", 0).getBoxName());
        }
    }

    @Nested
    class ChangeBoxName {
        @Test
        void boxWithSimpleName() {
            AsciiBox asciiBox = AsciiBoxWriter.DEFAULT.boxString("simpleName", "some content", 0);
            asciiBox = asciiBox.changeBoxName("newSimpleName");
            assertEquals(AsciiBoxWriter.DEFAULT.boxString("newSimpleName", "some content", 0), asciiBox);
        }

        @Test
        void boxWithShorterName() {
            AsciiBox asciiBox = AsciiBoxWriter.DEFAULT.boxString("veryLongName", "some content", 0);
            asciiBox = asciiBox.changeBoxName("name");
            assertEquals(AsciiBoxWriter.DEFAULT.boxString("name", "some content", 0), asciiBox);
        }

        @Test
        void boxGettingDressed() {
            AsciiBox asciiBox = AsciiBoxWriter.DEFAULT.boxString("", "some content", 0);
            asciiBox = asciiBox.changeBoxName("name");
            assertEquals(AsciiBoxWriter.DEFAULT.boxString("name", "some content", 0), asciiBox);
        }
    }

    @Nested
    class IsEmpty {
        @Test
        void emptyBox() {
            assertTrue(new AsciiBox("").isEmpty());
        }

        @Test
        void nonEmptyBox() {
            assertFalse(new AsciiBox("a").isEmpty());
        }

        @Test
        void nameEmptyBox() {
            assertTrue(AsciiBoxWriter.DEFAULT.boxString("name", "", 0).isEmpty());
        }

        @Test
        void nameNonEmptyBox() {
            assertFalse(AsciiBoxWriter.DEFAULT.boxString("name", "a", 0).isEmpty());
        }
    }

    @Nested
    class SideBySide {
        @Test
        void Test2Boxes() {
            String box1 = "" +
                "000 0x: 31  32  33  34  35  36  37  38  '12345678'\n" +
                "008 0x: 39  30  61  62  63  64  65  66  '90abcdef'\n" +
                "016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'\n" +
                "024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'\n" +
                "032 0x: 77  78  79  7a                  'wxyz    '";
            String box2 = "" +
                "╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗\n" +
                "║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║\n" +
                "║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║\n" +
                "║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║\n" +
                "║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║\n" +
                "║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║\n" +
                "║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║\n" +
                "║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║\n" +
                "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝";
            String want = "" +
                "000 0x: 31  32  33  34  35  36  37  38  '12345678'╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗\n" +
                "008 0x: 39  30  61  62  63  64  65  66  '90abcdef'║  000 0x: 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║\n" +
                "016 0x: 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'║  024 0x: 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║\n" +
                "024 0x: 6f  70  71  72  73  74  75  76  'opqrstuv'║  048 0x: 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║\n" +
                "032 0x: 77  78  79  7a                  'wxyz    '║  072 0x: 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║\n" +
                "                                                  ║  096 0x: 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║\n" +
                "                                                  ║  120 0x: 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║\n" +
                "                                                  ║  144 0x: 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║\n" +
                "                                                  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxSideBySide(new AsciiBox(box1), new AsciiBox(box2)));
        }

        @Test
        void anotherTwoBoxes() {
            String box1 = "" +
                "╔═exampleInt╗\n" +
                "║     4     ║\n" +
                "╚═══════════╝";
            String box2 = "" +
                "╔═exampleInt╗\n" +
                "║     7     ║\n" +
                "╚═══════════╝";
            String want = "" +
                "╔═exampleInt╗╔═exampleInt╗\n" +
                "║     4     ║║     7     ║\n" +
                "╚═══════════╝╚═══════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxSideBySide(new AsciiBox(box1), new AsciiBox(box2)));
        }

        @Test
        void sizeDifferenceFirstBox() {
            String box1 = "" +
                "╔═exampleInt╗\n" +
                "║     4     ║\n" +
                "║     4     ║\n" +
                "╚═══════════╝";
            String box2 = "" +
                "╔═exampleInt╗\n" +
                "║     7     ║\n" +
                "╚═══════════╝";
            String want = "" +
                "╔═exampleInt╗╔═exampleInt╗\n" +
                "║     4     ║║     7     ║\n" +
                "║     4     ║╚═══════════╝\n" +
                "╚═══════════╝             ";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxSideBySide(new AsciiBox(box1), new AsciiBox(box2)));
        }

        @Test
        void sizeDifferenceSecondBox() {
            String box1 = "" +
                "╔═exampleInt╗\n" +
                "║     4     ║\n" +
                "╚═══════════╝";
            String box2 = "" +
                "╔═exampleInt╗\n" +
                "║     7     ║\n" +
                "║     7     ║\n" +
                "╚═══════════╝";
            String want = "" +
                "╔═exampleInt╗╔═exampleInt╗\n" +
                "║     4     ║║     7     ║\n" +
                "╚═══════════╝║     7     ║\n" +
                "             ╚═══════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxSideBySide(new AsciiBox(box1), new AsciiBox(box2)));
        }
    }

    @Nested
    class BoxBelowBox {
        @Test
        void testTwoBoxes() {
            String box1 = "" +
                "000 31  32  33  34  35  36  37  38  '12345678'\n" +
                "008 39  30  61  62  63  64  65  66  '90abcdef'\n" +
                "016 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'\n" +
                "024 6f  70  71  72  73  74  75  76  'opqrstuv'\n" +
                "032 77  78  79  7a                  'wxyz    '";
            String box2 = "" +
                "╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗\n" +
                "║  000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║\n" +
                "║  024 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║\n" +
                "║  048 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║\n" +
                "║  072 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║\n" +
                "║  096 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║\n" +
                "║  120 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║\n" +
                "║  144 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║\n" +
                "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝";
            String want = "" +
                "000 31  32  33  34  35  36  37  38  '12345678'                                                                                      \n" +
                "008 39  30  61  62  63  64  65  66  '90abcdef'                                                                                      \n" +
                "016 67  68  69  6a  6b  6c  6d  6e  'ghijklmn'                                                                                      \n" +
                "024 6f  70  71  72  73  74  75  76  'opqrstuv'                                                                                      \n" +
                "032 77  78  79  7a                  'wxyz    '                                                                                      \n" +
                "╔═super nice data══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗\n" +
                "║  000 31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  '1234567890abcdefghijklmn'  ║\n" +
                "║  024 6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  38  39  30  61  'opqrstuvwxyz.1234567890a'  ║\n" +
                "║  048 62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  'bcdefghijklmnopqrstuvwxy'  ║\n" +
                "║  072 7a  d3  61  61  31  32  33  34  35  36  37  38  39  30  61  62  63  64  65  66  67  68  69  6a  'z.aa1234567890abcdefghij'  ║\n" +
                "║  096 6b  6c  6d  6e  6f  70  71  72  73  74  75  76  77  78  79  7a  d3  31  32  33  34  35  36  37  'klmnopqrstuvwxyz.1234567'  ║\n" +
                "║  120 38  39  30  61  62  63  64  65  66  67  68  69  6a  6b  6c  6d  6e  6f  70  71  72  73  74  75  '890abcdefghijklmnopqrstu'  ║\n" +
                "║  144 76  77  78  79  7a  d3  61  61  62                                                              'vwxyz.aab               '  ║\n" +
                "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxBelowBox(new AsciiBox(box1), new AsciiBox(box2)));
        }

        @Test
        void testDifferentSizedBoxes() {
            String box1 = "" +
                "╔═sampleField════════════╗\n" +
                "║123123123123123123123123║\n" +
                "╚════════════════════════╝";
            String box2 = "" +
                "╔═sampleField╗\n" +
                "║123123123123║\n" +
                "╚════════════╝";
            String want = "" +
                "╔═sampleField════════════╗\n" +
                "║123123123123123123123123║\n" +
                "╚════════════════════════╝\n" +
                "╔═sampleField╗            \n" +
                "║123123123123║            \n" +
                "╚════════════╝            ";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxBelowBox(new AsciiBox(box1), new AsciiBox(box2)));
        }
    }

    @Nested
    class BoxString {
        @Test
        void simpleBox() {
            String name = "sampleField";
            String data = "123123123123";
            int charWidth = 1;
            String want = "" +
                "╔═sampleField╗\n" +
                "║123123123123║\n" +
                "╚════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxString(name, data, charWidth));
        }

        @Test
        void simpleBoxUnnamed() {
            String name = "";
            String data = "123123123123";
            int charWidth = 1;
            String want = "" +
                "╔════════════╗\n" +
                "║123123123123║\n" +
                "╚════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxString(name, data, charWidth));
        }

        @Test
        void simpleBox2() {
            String name = "sampleField";
            String data = "123123123123\n123123123123123123123123";
            int charWidth = 1;
            String want = "" +
                "╔═sampleField════════════╗\n" +
                "║      123123123123      ║\n" +
                "║123123123123123123123123║\n" +
                "╚════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxString(name, data, charWidth));
        }

        @Test
        void simpleBoxWithTooLongName() {
            String name = "sampleFieldsampleFieldsampleFieldsampleField";
            String data = "123123123123\n123123123123123123123123";
            int charWidth = 1;
            String want = "" +
                "╔═sampleFieldsampleFieldsampleFieldsampleField╗\n" +
                "║                123123123123                 ║\n" +
                "║          123123123123123123123123           ║\n" +
                "╚═════════════════════════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.boxString(name, data, charWidth));
        }
    }

    @Nested
    class AlignBoxes {
        @Test
        void enoughSpace() {
            Collection<AsciiBox> boxes = Arrays.asList(
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝")
            );
            int desiredWidth = 1000;
            String want = "" +
                "╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.alignBoxes(boxes, desiredWidth));
        }

        @Test
        void notEnoughSpace() {
            Collection<AsciiBox> boxes = Arrays.asList(
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝")
            );
            int desiredWidth = 0;
            String want = "" +
                "╔═sampleField════════════╗\n" +
                "║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║\n" +
                "╚════════════════════════╝\n" +
                "╔═sampleField════════════╗\n" +
                "║      123123123123      ║\n" +
                "║123123123123123123123123║\n" +
                "╚════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.alignBoxes(boxes, desiredWidth));
        }

        @Test
        void notEnoughSpaceShouldResultInMultipleRows() {
            Collection<AsciiBox> boxes = Arrays.asList(
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝")
            );
            int desiredWidth = 65;
            String want = "" +
                "╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝\n" +
                "╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝\n" +
                "╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝\n" +
                "╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.alignBoxes(boxes, desiredWidth));
        }

        @Test
        void notEnoughSpaceShouldResultInMultipleRowsThreeColums() {
            Collection<AsciiBox> boxes = Arrays.asList(
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123ABABABABABAB123123║\n" +
                    "╚════════════════════════╝"),
                new AsciiBox("" +
                    "╔═sampleField════════════╗\n" +
                    "║      123123123123      ║\n" +
                    "║123123123123123123123123║\n" +
                    "╚════════════════════════╝")
            );
            int desiredWidth = 78;
            String want = "" +
                "╔═sampleField════════════╗╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║║      123123123123      ║\n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║║123123ABABABABABAB123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝╚════════════════════════╝\n" +
                "╔═sampleField════════════╗╔═sampleField════════════╗╔═sampleField════════════╗\n" +
                "║      123123123123      ║║      123123123123      ║║      123123123123      ║\n" +
                "║123123123123123123123123║║123123ABABABABABAB123123║║123123123123123123123123║\n" +
                "╚════════════════════════╝╚════════════════════════╝╚════════════════════════╝\n" +
                "╔═sampleField════════════╗╔═sampleField════════════╗                          \n" +
                "║      123123123123      ║║      123123123123      ║                          \n" +
                "║123123ABABABABABAB123123║║123123123123123123123123║                          \n" +
                "╚════════════════════════╝╚════════════════════════╝                          ";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.alignBoxes(boxes, desiredWidth));
        }
    }

    @Nested
    class Width {
        @Test
        void sameWidth() {
            AsciiBox asciiBox = new AsciiBox("" +
                "123123123123123\n" +
                "123123123123123\n" +
                "123123123123123");
            assertEquals(15, asciiBox.width());
        }

        @Test
        void differentWidth() {
            AsciiBox asciiBox = new AsciiBox("" +
                "123123123123123\n" +
                "123123123123123123123123123123\n" +
                "123123123123123");
            assertEquals(30, asciiBox.width());
        }
    }

    @Nested
    class MergeHorizontal {
        @Test
        void threeSame() {
            List<AsciiBox> boxes = Arrays.asList(
                new AsciiBox("" +
                    "123123123\n" +
                    "123123123\n" +
                    "123123123"),
                new AsciiBox("" +
                    "abcabcabc\n" +
                    "abcabcabc\n" +
                    "abcabcabc"),
                new AsciiBox("" +
                    "zxyzxyzxy\n" +
                    "zxyzxyzxy\n" +
                    "zxyzxyzxy")
            );
            String want = "" +
                "123123123abcabcabczxyzxyzxy\n" +
                "123123123abcabcabczxyzxyzxy\n" +
                "123123123abcabcabczxyzxyzxy";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.mergeHorizontal(boxes));
        }

        @Test
        void threeDifferent() {
            List<AsciiBox> boxes = Arrays.asList(
                new AsciiBox("" +
                    "123123123\n" +
                    "123123123\n" +
                    "123123123"),
                new AsciiBox("" +
                    "abcabcabc\n" +
                    "abcabcabcabcabcabcabcabcabc\n" +
                    "abcabcabc"),
                new AsciiBox("" +
                    "zxyzxyzxy\n" +
                    "zxyzxyzxy\n" +
                    "zxyzxyzxy")
            );
            String want = "" +
                "123123123abcabcabc                  zxyzxyzxy\n" +
                "123123123abcabcabcabcabcabcabcabcabczxyzxyzxy\n" +
                "123123123abcabcabc                  zxyzxyzxy";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.DEFAULT.mergeHorizontal(boxes));
        }
    }

    @Nested
    class ExpandBox {
        @Test
        void smallExpand() {
            AsciiBox box = new AsciiBox("" +
                "123123123\n" +
                "123123123\n" +
                "123123123");
            AsciiBox want = new AsciiBox("" +
                "123123123                                                                                           \n" +
                "123123123                                                                                           \n" +
                "123123123                                                                                           ");
            assertEquals(want, AsciiBoxWriter.DEFAULT.expandBox(box, 100));
        }

        @Test
        void bigExpand() {
            AsciiBox box = new AsciiBox("" +
                "123123123\n" +
                "123123123\n" +
                "123123123");
            AsciiBox want = new AsciiBox("" +
                "123123123                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       \n" +
                "123123123                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       \n" +
                "123123123                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       ");
            assertEquals(want, AsciiBoxWriter.DEFAULT.expandBox(box, 10000));
        }
    }

    @Nested
    class LightBoxString {
        @Test
        void simpleBox() {
            String name = "sampleField";
            String data = "123123123123";
            int charWidth = 1;
            String want = "" +
                "╭┄sampleField╮\n" +
                "┆123123123123┆\n" +
                "╰┄┄┄┄┄┄┄┄┄┄┄┄╯";
            assertEquals(new AsciiBox(want), AsciiBoxWriter.LIGHT.boxString(name, data, charWidth));
        }
    }
}