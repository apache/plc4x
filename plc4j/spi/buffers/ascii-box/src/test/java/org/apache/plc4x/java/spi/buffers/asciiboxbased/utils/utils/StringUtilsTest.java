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
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.utils;

import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.*;

class StringUtilsTest {

    @Test
    void testIsBlankNull() {
        assertTrue(StringUtils.isBlank(null));
    }

    @Test
    void testIsBlankEmpty() {
        assertTrue(StringUtils.isBlank(""));
    }

    @Test
    void testIsBlankWhitespace() {
        assertTrue(StringUtils.isBlank("   "));
        assertTrue(StringUtils.isBlank("\t\n"));
    }

    @Test
    void testIsBlankNotBlank() {
        assertFalse(StringUtils.isBlank("a"));
        assertFalse(StringUtils.isBlank("  a  "));
    }

    @Test
    void testIsBlankCharSequence() {
        // Test with StringBuilder (non-String CharSequence)
        assertTrue(StringUtils.isBlank(new StringBuilder()));
        assertTrue(StringUtils.isBlank(new StringBuilder("   ")));
        assertFalse(StringUtils.isBlank(new StringBuilder("abc")));
    }

    @Test
    void testRepeatNull() {
        assertNull(StringUtils.repeat(null, 3));
    }

    @Test
    void testRepeatZero() {
        assertEquals("", StringUtils.repeat("abc", 0));
        assertEquals("", StringUtils.repeat("abc", -1));
    }

    @Test
    void testRepeatEmpty() {
        assertEquals("", StringUtils.repeat("", 5));
    }

    @Test
    void testRepeatNormal() {
        assertEquals("abcabcabc", StringUtils.repeat("abc", 3));
        assertEquals("xxx", StringUtils.repeat("x", 3));
    }

    @Test
    void testJoinArrayNull() {
        assertNull(StringUtils.join((Object[]) null, ","));
    }

    @Test
    void testJoinArrayEmpty() {
        assertEquals("", StringUtils.join(new Object[]{}, ","));
    }

    @Test
    void testJoinArrayNormal() {
        assertEquals("a,b,c", StringUtils.join(new Object[]{"a", "b", "c"}, ","));
    }

    @Test
    void testJoinArrayNullDelimiter() {
        assertEquals("abc", StringUtils.join(new Object[]{"a", "b", "c"}, null));
    }

    @Test
    void testJoinIterableNull() {
        assertNull(StringUtils.join((Iterable<?>) null, ","));
    }

    @Test
    void testJoinIterableEmpty() {
        assertEquals("", StringUtils.join(Collections.emptyList(), ","));
    }

    @Test
    void testJoinIterableNormal() {
        assertEquals("a,b,c", StringUtils.join(Arrays.asList("a", "b", "c"), ","));
    }

    @Test
    void testJoinIterableNullDelimiter() {
        assertEquals("abc", StringUtils.join(Arrays.asList("a", "b", "c"), null));
    }

    @Test
    void testTrimNull() {
        assertNull(StringUtils.trim(null));
    }

    @Test
    void testTrimNormal() {
        assertEquals("abc", StringUtils.trim("  abc  "));
        assertEquals("abc", StringUtils.trim("abc"));
    }

    @Test
    void testContainsAnyNull() {
        assertFalse(StringUtils.containsAny(null, "abc"));
        assertFalse(StringUtils.containsAny("abc", null));
    }

    @Test
    void testContainsAnyEmpty() {
        assertFalse(StringUtils.containsAny("", "abc"));
        assertFalse(StringUtils.containsAny("abc", ""));
    }

    @Test
    void testContainsAnyFound() {
        assertTrue(StringUtils.containsAny("hello", "aeiou"));
        assertTrue(StringUtils.containsAny("xyz", "z"));
    }

    @Test
    void testContainsAnyNotFound() {
        assertFalse(StringUtils.containsAny("bcd", "xyz"));
    }
}
