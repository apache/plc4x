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
package org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.lang3;

import java.util.Arrays;
import java.util.Iterator;
import java.util.Objects;

public class StringUtils {

    public static final String EMPTY = "";

    private static final int PAD_LIMIT = 8192;

    public static final char[] EMPTY_CHAR_ARRAY = {};

    public static int length(final CharSequence cs) {
        return cs == null ? 0 : cs.length();
    }

    public static boolean isBlank(final CharSequence cs) {
        final int strLen = length(cs);
        if (strLen == 0) {
            return true;
        }
        for (int i = 0; i < strLen; i++) {
            if (!Character.isWhitespace(cs.charAt(i))) {
                return false;
            }
        }
        return true;
    }

    public static boolean isEmpty(final CharSequence cs) {
        return cs == null || cs.isEmpty();
    }

    public static String repeat(final String repeat, final int count) {
        // Performance tuned for 2.0 (JDK1.4)
        if (repeat == null) {
            return null;
        }
        if (count <= 0) {
            return EMPTY;
        }
        final int inputLength = repeat.length();
        if (count == 1 || inputLength == 0) {
            return repeat;
        }
        if (inputLength == 1 && count <= PAD_LIMIT) {
            return repeat(repeat.charAt(0), count);
        }
        final int outputLength = inputLength * count;
        switch (inputLength) {
            case 1:
                return repeat(repeat.charAt(0), count);
            case 2:
                final char ch0 = repeat.charAt(0);
                final char ch1 = repeat.charAt(1);
                final char[] output2 = new char[outputLength];
                for (int i = count * 2 - 2; i >= 0; i--, i--) {
                    output2[i] = ch0;
                    output2[i + 1] = ch1;
                }
                return new String(output2);
            default:
                return repeat.repeat(count);
        }
    }

    public static String repeat(final char repeat, final int count) {
        if (count <= 0) {
            return EMPTY;
        }
        char[] chars = new char[count];
        Arrays.fill(chars, repeat);
        return String.valueOf(chars);
    }

    public static String join(final Object[] array, final String delimiter) {
        return array != null ? join(array, Objects.toString(delimiter, EMPTY), 0, array.length) : null;
    }

    public static String join(final Object[] array, final String delimiter, final int startIndex, final int endIndex) {
        return array != null ? Streams.of(array).skip(startIndex).limit(Math.max(0, endIndex - startIndex))
            .collect(LangCollectors.joining(delimiter, EMPTY, EMPTY, del -> Objects.toString(del, EMPTY))) : null;
    }

    public static String join(final Iterable<?> iterable, final String separator) {
        return iterable != null ? join(iterable.iterator(), separator) : null;
    }

    public static String join(final Iterator<?> iterator, final String separator) {
        // handle null, zero and one elements before building a buffer
        if (iterator == null) {
            return null;
        }
        if (!iterator.hasNext()) {
            return EMPTY;
        }
        return Streams.of(iterator).collect(LangCollectors.joining(Objects.toString(separator, EMPTY), EMPTY, EMPTY, it ->  Objects.toString(it, EMPTY)));
    }

    public static String trim(final String str) {
        return str == null ? null : str.trim();
    }

    public static boolean containsAny(final CharSequence cs, final CharSequence searchChars) {
        if (searchChars == null) {
            return false;
        }
        return containsAny(cs, toCharArray(searchChars));
    }

    public static boolean containsAny(final CharSequence cs, final char... searchChars) {
        if (isEmpty(cs) || (searchChars.length == 0)) {
            return false;
        }
        final int csLength = cs.length();
        final int searchLength = searchChars.length;
        final int csLast = csLength - 1;
        final int searchLast = searchLength - 1;
        for (int i = 0; i < csLength; i++) {
            final char ch = cs.charAt(i);
            for (int j = 0; j < searchLength; j++) {
                if (searchChars[j] == ch) {
                    if (!Character.isHighSurrogate(ch) || j == searchLast || i < csLast && searchChars[j + 1] == cs.charAt(i + 1)) {
                        return true;
                    }
                }
            }
        }
        return false;
    }

    public static char[] toCharArray(final CharSequence source) {
        final int len = length(source);
        if (len == 0) {
            return EMPTY_CHAR_ARRAY;
        }
        if (source instanceof String) {
            return ((String) source).toCharArray();
        }
        final char[] array = new char[len];
        for (int i = 0; i < len; i++) {
            array[i] = source.charAt(i);
        }
        return array;
    }

}
