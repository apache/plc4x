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

package org.apache.plc4x.java.spi.config;

import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;

class ConfigurationParameterConverterTest {

    @Test
    void testDurationConverter() {
        DurationConverter converter = new DurationConverter();

        assertEquals(Duration.class, converter.getType());
        assertEquals(Duration.ofSeconds(30), converter.convert("30s"));
        assertEquals(Duration.ofMinutes(5), converter.convert("5m"));
        assertEquals(Duration.ofHours(2), converter.convert("2h"));
        assertEquals(Duration.ofMillis(500), converter.convert("500ms"));
    }

    @Test
    void testDurationConverter_invalidFormat() {
        DurationConverter converter = new DurationConverter();

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
            () -> converter.convert("invalid"));

        assertTrue(exception.getMessage().contains("Invalid duration format"));
    }

    @Test
    void testTimeUnitConverter() {
        TimeUnitConverter converter = new TimeUnitConverter();

        assertEquals(TimeUnit.class, converter.getType());
        assertEquals(TimeUnit.SECONDS, converter.convert("SECONDS"));
        assertEquals(TimeUnit.MILLISECONDS, converter.convert("MILLISECONDS"));
        assertEquals(TimeUnit.MINUTES, converter.convert("MINUTES"));
    }

    @Test
    void testTimeUnitConverter_invalidValue() {
        TimeUnitConverter converter = new TimeUnitConverter();

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
            () -> converter.convert("INVALID_UNIT"));

        assertTrue(exception.getMessage().contains("Invalid TimeUnit"));
    }

    @Test
    void testIntegerArrayConverter() {
        IntegerArrayConverter converter = new IntegerArrayConverter();

        assertEquals(int[].class, converter.getType());
        assertArrayEquals(new int[]{1, 2, 3, 4}, converter.convert("1,2,3,4"));
        assertArrayEquals(new int[]{42}, converter.convert("42"));
        assertArrayEquals(new int[]{}, converter.convert(""));
    }

    @Test
    void testIntegerArrayConverter_invalidNumbers() {
        IntegerArrayConverter converter = new IntegerArrayConverter();

        NumberFormatException exception = assertThrows(NumberFormatException.class,
            () -> converter.convert("1,abc,3"));
    }

    @Test
    void testStringArrayConverter() {
        StringArrayConverter converter = new StringArrayConverter();

        assertEquals(String[].class, converter.getType());
        assertArrayEquals(new String[]{"a", "b", "c"}, converter.convert("a,b,c"));
        assertArrayEquals(new String[]{"single"}, converter.convert("single"));
        assertArrayEquals(new String[]{}, converter.convert(""));
    }

    @Test
    void testStringArrayConverter_withSpaces() {
        StringArrayConverter converter = new StringArrayConverter();

        assertArrayEquals(new String[]{"a", "b", "c"}, converter.convert("a, b, c"));
        assertArrayEquals(new String[]{"hello world", "test"}, converter.convert("hello world,test"));
    }

    // Sample converter implementations for testing

    public static class DurationConverter implements ConfigurationParameterConverter<Duration> {
        @Override
        public Class<Duration> getType() {
            return Duration.class;
        }

        @Override
        public Duration convert(String value) {
            if (value.endsWith("ms")) {
                return Duration.ofMillis(Long.parseLong(value.substring(0, value.length() - 2)));
            } else if (value.endsWith("s")) {
                return Duration.ofSeconds(Long.parseLong(value.substring(0, value.length() - 1)));
            } else if (value.endsWith("m")) {
                return Duration.ofMinutes(Long.parseLong(value.substring(0, value.length() - 1)));
            } else if (value.endsWith("h")) {
                return Duration.ofHours(Long.parseLong(value.substring(0, value.length() - 1)));
            } else {
                throw new IllegalArgumentException("Invalid duration format: " + value);
            }
        }
    }

    public static class TimeUnitConverter implements ConfigurationParameterConverter<TimeUnit> {
        @Override
        public Class<TimeUnit> getType() {
            return TimeUnit.class;
        }

        @Override
        public TimeUnit convert(String value) {
            try {
                return TimeUnit.valueOf(value);
            } catch (IllegalArgumentException e) {
                throw new IllegalArgumentException("Invalid TimeUnit: " + value, e);
            }
        }
    }

    public static class IntegerArrayConverter implements ConfigurationParameterConverter<int[]> {
        @Override
        public Class<int[]> getType() {
            return int[].class;
        }

        @Override
        public int[] convert(String value) {
            if (value.trim().isEmpty()) {
                return new int[0];
            }
            String[] parts = value.split(",");
            int[] result = new int[parts.length];
            for (int i = 0; i < parts.length; i++) {
                result[i] = Integer.parseInt(parts[i].trim());
            }
            return result;
        }
    }

    public static class StringArrayConverter implements ConfigurationParameterConverter<String[]> {
        @Override
        public Class<String[]> getType() {
            return String[].class;
        }

        @Override
        public String[] convert(String value) {
            if (value.trim().isEmpty()) {
                return new String[0];
            }
            String[] parts = value.split(",");
            for (int i = 0; i < parts.length; i++) {
                parts[i] = parts[i].trim();
            }
            return parts;
        }
    }
}