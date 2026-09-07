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
package org.apache.plc4x.java.utils.testutils.utils.migration;

import org.junit.jupiter.api.Test;

import java.util.regex.Pattern;

import static org.junit.jupiter.api.Assertions.*;

class TestCasePatcherTest {

    @Test
    void testIndent() {
        String xml = "<root>\n<child/>\n</root>";
        String result = TestCasePatcher.indent(xml, "  ");

        assertTrue(result.contains("  <root>"));
        assertTrue(result.contains("  <child/>"));
        assertTrue(result.contains("  </root>"));
    }

    @Test
    void testIndentSingleLine() {
        String xml = "<root/>";
        String result = TestCasePatcher.indent(xml, "    ");

        assertEquals("    <root/>", result);
    }

    @Test
    void testIndentEmptyString() {
        String result = TestCasePatcher.indent("", "  ");
        assertEquals("  ", result);
    }

    @Test
    void testIndentWithTabs() {
        String xml = "<root>\n<child/>\n</root>";
        String result = TestCasePatcher.indent(xml, "\t");

        assertTrue(result.contains("\t<root>"));
        assertTrue(result.contains("\t<child/>"));
        assertTrue(result.contains("\t</root>"));
    }

    @Test
    void testDetermineIndent() {
        String document = "  <root>\n  <child/>\n  </root>\n";
        String fragment = "<root>\n<child/>\n</root>";

        String indent = TestCasePatcher.determineIndent(document, fragment);

        assertEquals("  ", indent);
    }

    @Test
    void testDetermineIndentWithDifferentIndents() {
        String document = "    <root>\n    <child/>\n    </root>\n";
        String fragment = "<root>\n<child/>\n</root>";

        String indent = TestCasePatcher.determineIndent(document, fragment);

        assertEquals("    ", indent);
    }

    @Test
    void testDetermineIndentWithNoIndent() {
        String document = "<root>\n<child/>\n</root>\n";
        String fragment = "<root>\n<child/>\n</root>";

        String indent = TestCasePatcher.determineIndent(document, fragment);

        assertEquals("", indent);
    }

    @Test
    void testDetermineIndentNotFound() {
        String document = "<root>\n<other/>\n</root>";
        String fragment = "<missing>\n<child/>\n</missing>";

        assertThrows(RuntimeException.class, () -> {
            TestCasePatcher.determineIndent(document, fragment);
        });
    }

    @Test
    void testGetPatternForFragment() {
        String fragment = "<root>\n<child/>\n</root>";

        Pattern pattern = TestCasePatcher.getPatternForFragment(fragment);

        assertNotNull(pattern);

        // Test that the pattern matches correctly indented fragments
        assertTrue(pattern.matcher("  <root>\n  <child/>\n  </root>\n").find());
        assertTrue(pattern.matcher("    <root>\n    <child/>\n    </root>\n").find());
        assertTrue(pattern.matcher("<root>\n<child/>\n</root>\n").find());
    }

    @Test
    void testGetPatternForFragmentWithSpecialCharacters() {
        String fragment = "<test name=\"value\">\n<data>1.5</data>\n</test>";

        Pattern pattern = TestCasePatcher.getPatternForFragment(fragment);

        assertNotNull(pattern);

        // Test that special characters are properly escaped
        assertTrue(pattern.matcher("  <test name=\"value\">\n  <data>1.5</data>\n  </test>\n").find());
    }

    @Test
    void testGetPatternForSingleLine() {
        String fragment = "<root/>";

        Pattern pattern = TestCasePatcher.getPatternForFragment(fragment);

        assertNotNull(pattern);
        assertTrue(pattern.matcher("  <root/>\n").find());
    }
}
