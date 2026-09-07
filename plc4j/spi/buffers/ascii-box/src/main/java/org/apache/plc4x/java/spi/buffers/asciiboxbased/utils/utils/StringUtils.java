/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.utils;

import java.util.Arrays;
import java.util.Objects;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

/**
 * Simplified string utilities - uses Java built-in methods where possible.
 */
public final class StringUtils {

    private StringUtils() {
        // Utility class
    }

    /**
     * Checks if a CharSequence is empty (""), null or whitespace only.
     */
    public static boolean isBlank(final CharSequence cs) {
        if (cs == null || cs.isEmpty()) {
            return true;
        }
        // For String, use built-in isBlank() which is optimized
        if (cs instanceof String s) {
            return s.isBlank();
        }
        // Fallback for other CharSequence implementations
        for (int i = 0; i < cs.length(); i++) {
            if (!Character.isWhitespace(cs.charAt(i))) {
                return false;
            }
        }
        return true;
    }

    /**
     * Repeats a String n times to form a new String.
     */
    public static String repeat(final String str, final int count) {
        if (str == null) {
            return null;
        }
        if (count <= 0 || str.isEmpty()) {
            return "";
        }
        return str.repeat(count);
    }

    /**
     * Joins the elements of the provided array into a single String.
     */
    public static String join(final Object[] array, final String delimiter) {
        if (array == null) {
            return null;
        }
        return Arrays.stream(array)
            .map(Objects::toString)
            .collect(Collectors.joining(delimiter != null ? delimiter : ""));
    }

    /**
     * Joins the elements of the provided Iterable into a single String.
     */
    public static String join(final Iterable<?> iterable, final String delimiter) {
        if (iterable == null) {
            return null;
        }
        return StreamSupport.stream(iterable.spliterator(), false)
            .map(Objects::toString)
            .collect(Collectors.joining(delimiter != null ? delimiter : ""));
    }

    /**
     * Removes leading and trailing whitespace.
     */
    public static String trim(final String str) {
        return str == null ? null : str.trim();
    }

    /**
     * Checks if the CharSequence contains any character from the searchChars.
     */
    public static boolean containsAny(final CharSequence cs, final CharSequence searchChars) {
        if (cs == null || cs.isEmpty() || searchChars == null || searchChars.isEmpty()) {
            return false;
        }
        for (int i = 0; i < cs.length(); i++) {
            char ch = cs.charAt(i);
            for (int j = 0; j < searchChars.length(); j++) {
                if (ch == searchChars.charAt(j)) {
                    return true;
                }
            }
        }
        return false;
    }
}
