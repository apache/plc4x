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

package org.apache.plc4x.java.spi.utils.json;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.io.File;
import java.io.IOException;
import java.io.Reader;
import java.io.StringReader;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.net.URI;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.time.Duration;
import java.time.Instant;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class JsonTest {

    // ----------------------------------------------------------------------------------------
    // Parsing
    // ----------------------------------------------------------------------------------------

    @Test
    void parsesScalars() throws Exception {
        assertEquals("abc", Json.parse("\"abc\""));
        assertEquals(Boolean.TRUE, Json.parse("true"));
        assertEquals(Boolean.FALSE, Json.parse(" false "));
        assertNull(Json.parse("null"));
        assertEquals(0L, Json.parse("0"));
        assertEquals(-42L, Json.parse("-42"));
        assertEquals(1.5, Json.parse("1.5"));
        assertEquals(-2.5e-3, Json.parse("-2.5e-3"));
        assertEquals(1e10, Json.parse("1E+10"));
        assertEquals(Long.MAX_VALUE, Json.parse(Long.toString(Long.MAX_VALUE)));
        assertEquals(Long.MIN_VALUE, Json.parse(Long.toString(Long.MIN_VALUE)));
        assertEquals(new BigInteger("9223372036854775808"), Json.parse("9223372036854775808"));
    }

    @Test
    void parsesNestedStructuresKeepingMemberOrder() throws Exception {
        Object parsed = Json.parse("{\"b\": [1, {\"x\": null}, []], \"a\": {}, \"c\": \"d\"}");
        Map<String, Object> expected = new LinkedHashMap<>();
        Map<String, Object> inner = new LinkedHashMap<>();
        inner.put("x", null);
        expected.put("b", Arrays.asList(1L, inner, List.of()));
        expected.put("a", Map.of());
        expected.put("c", "d");
        assertEquals(expected, parsed);
        assertEquals(List.of("b", "a", "c"), new ArrayList<>(((Map<?, ?>) parsed).keySet()));
    }

    @Test
    void parsesEscapes() throws Exception {
        assertEquals("\"\\/\b\f\n\r\t\u00e9\uD83D\uDE00end",
            Json.parse("\"\\\"\\\\\\/\\b\\f\\n\\r\\t\\u00E9\\ud83d\\ude00end\""));
        assertEquals("a\"", Json.parse("\"a\\\"\""));
    }

    @Test
    void parsesFromReader() throws Exception {
        assertEquals(List.of(1L, 2L), Json.parse(new StringReader("[1,2]")));
    }

    @ParameterizedTest
    @ValueSource(strings = {
        "", " ", "{", "[", "[1", "[1,]", "[1 2]", "{\"a\"}", "{\"a\":1,}", "{\"a\":1 \"b\":2}", "{a:1}",
        "\"abc", "\"a\nb\"", "\"\\x\"", "\"\\u12\"", "\"\\u12G4\"", "tru", "nul", "falsy", "01x", "-",
        "-a", "1.", "1.e5", "1e", "1e+", "+1", "1 2", "[]]", "@", "\"\\"
    })
    void rejectsInvalidJson(String json) {
        assertThrows(JsonException.class, () -> Json.parse(json));
    }

    @Test
    void rejectsTooDeepNesting() {
        String deep = "[".repeat(Json.MAX_DEPTH + 1) + "]".repeat(Json.MAX_DEPTH + 1);
        JsonException e = assertThrows(JsonException.class, () -> Json.parse(deep));
        assertTrue(e.getMessage().contains("Nesting"), e.getMessage());
        String deepObject = "{\"a\":".repeat(Json.MAX_DEPTH + 1);
        assertThrows(JsonException.class, () -> Json.parse(deepObject));
        String ok = "[".repeat(Json.MAX_DEPTH) + "]".repeat(Json.MAX_DEPTH);
        assertTrue(() -> {
            try {
                return Json.parse(ok) instanceof List;
            } catch (JsonException ex) {
                return false;
            }
        });
    }

    @Test
    void reportsReaderFailures() {
        Reader broken = new Reader() {
            @Override
            public int read(char[] cbuf, int off, int len) throws IOException {
                throw new IOException("boom");
            }

            @Override
            public void close() {
            }
        };
        IOException e = assertThrows(IOException.class, () -> Json.parse(broken));
        assertEquals("boom", e.getMessage());
    }

    // ----------------------------------------------------------------------------------------
    // Writing
    // ----------------------------------------------------------------------------------------

    enum Color {RED}

    record Point(int x, int y) {
    }

    @SuppressWarnings("unused")
    public static class Base {
        private final String name = "base";

        public String getName() {
            return name;
        }
    }

    @SuppressWarnings("unused")
    private static class Bean extends Base {
        private final int value = 42;
        private final boolean enabled = true;

        public int getValue() {
            return value;
        }

        public boolean isEnabled() {
            return enabled;
        }

        public String getURL() {
            return "u";
        }

        // Not getters:
        public Boolean isBoxed() {
            return Boolean.TRUE;
        }

        public String isNotABoolean() {
            return "x";
        }

        public String getter() {
            return "x";
        }

        public String get() {
            return "x";
        }

        public void getNothing() {
        }

        public String getWithArg(int i) {
            return "x";
        }

        public static String getStatic() {
            return "x";
        }

        String getPackagePrivate() {
            return "x";
        }
    }

    private static class Empty {
    }

    private static class Throwing {
        public String getBoom() {
            throw new IllegalStateException("boom");
        }
    }

    private static class SelfReferencing {
        public SelfReferencing getSelf() {
            return this;
        }
    }

    @Test
    void writesScalars() throws Exception {
        assertEquals("null", Json.toJson(null));
        assertEquals("\"abc\"", Json.toJson("abc"));
        assertEquals("\"abc\"", Json.toJson(new StringBuilder("abc")));
        assertEquals("\"c\"", Json.toJson('c'));
        assertEquals("true", Json.toJson(true));
        assertEquals("42", Json.toJson(42));
        assertEquals("1.5", Json.toJson(1.5));
        assertEquals("2.5", Json.toJson(2.5f));
        assertEquals("12345678901234567890", Json.toJson(new BigInteger("12345678901234567890")));
        assertEquals("1.10", Json.toJson(new BigDecimal("1.10")));
        assertEquals("\"NaN\"", Json.toJson(Double.NaN));
        assertEquals("\"Infinity\"", Json.toJson(Double.POSITIVE_INFINITY));
        assertEquals("\"-Infinity\"", Json.toJson(Float.NEGATIVE_INFINITY));
        assertEquals("\"NaN\"", Json.toJson(Float.NaN));
        assertEquals("\"RED\"", Json.toJson(Color.RED));
        assertEquals("\"x\"", Json.toJson(Optional.of("x")));
        assertEquals("null", Json.toJson(Optional.empty()));
        assertEquals("\"AQID\"", Json.toJson(new byte[]{1, 2, 3}));
    }

    @Test
    void writesValueTypesAsStrings() throws Exception {
        assertEquals("\"1970-01-01T00:00:00Z\"", Json.toJson(Instant.EPOCH));
        assertEquals("\"PT1S\"", Json.toJson(Duration.ofSeconds(1)));
        assertEquals("\"Europe/Berlin\"", Json.toJson(ZoneId.of("Europe/Berlin")));
        UUID uuid = UUID.randomUUID();
        assertEquals("\"" + uuid + "\"", Json.toJson(uuid));
        assertEquals("\"http://example.com\"", Json.toJson(URI.create("http://example.com")));
        URL url = URI.create("http://example.com").toURL();
        assertEquals("\"http://example.com\"", Json.toJson(url));
        assertEquals("\"a\"", Json.toJson(Path.of("a")));
        assertEquals("\"a\"", Json.toJson(new File("a")));
        assertEquals("\"de-DE\"", Json.toJson(Locale.GERMANY));
        assertEquals("\"UTF-8\"", Json.toJson(StandardCharsets.UTF_8));
    }

    @Test
    void escapesStrings() throws Exception {
        assertEquals("\"\\\"\\\\\\b\\f\\n\\r\\t\\u0001/\u00e9\"", Json.toJson("\"\\\b\f\n\r\t\u0001/\u00e9"));
    }

    @Test
    void writesContainers() throws Exception {
        Map<Object, Object> map = new LinkedHashMap<>();
        map.put("a", List.of(1, 2));
        map.put(3, new int[]{4, 5});
        map.put("e", new ArrayList<>());
        map.put("m", Map.of());
        map.put("n", null);
        assertEquals("{\"a\":[1,2],\"3\":[4,5],\"e\":[],\"m\":{},\"n\":null}", Json.toJson(map));
        assertEquals("[\"x\",null]", Json.toJson(new String[]{"x", null}));
        assertEquals("[]", Json.toJson(new long[0]));
    }

    @Test
    void writesRecordsAndBeans() throws Exception {
        assertEquals("{\"x\":1,\"y\":2}", Json.toJson(new Point(1, 2)));
        assertEquals("{\"name\":\"base\",\"value\":42,\"enabled\":true,\"boxed\":true,\"url\":\"u\"}",
            Json.toJson(new Bean()));
        assertEquals("{}", Json.toJson(new Empty()));
    }

    @Test
    void writesPretty() throws Exception {
        Map<String, Object> map = new LinkedHashMap<>();
        map.put("a", List.of(1, Map.of("b", true)));
        map.put("e", List.of());
        assertEquals("""
            {
              "a": [
                1,
                {
                  "b": true
                }
              ],
              "e": []
            }""", Json.toPrettyJson(map));
    }

    @Test
    void roundTrips() throws Exception {
        String json = "{\"s\":\"a\\nb\",\"n\":[1,-2.5,12345678901234567890],\"o\":{\"t\":true,\"f\":false,\"z\":null}}";
        assertEquals(json, Json.toJson(Json.parse(json)));
        assertEquals(Json.parse(json), Json.parse(Json.toPrettyJson(Json.parse(json))));
    }

    @Test
    void rejectsCyclesAndFailingGetters() {
        List<Object> list = new ArrayList<>();
        list.add(list);
        assertThrows(JsonException.class, () -> Json.toJson(list));
        assertThrows(JsonException.class, () -> Json.toJson(new SelfReferencing()));
        JsonException e = assertThrows(JsonException.class, () -> Json.toJson(new Throwing()));
        assertEquals("boom", e.getCause().getMessage());
    }

    @Test
    void sameObjectTwiceIsNotACycle() throws Exception {
        List<Integer> shared = List.of(1);
        assertEquals("[[1],[1]]", Json.toJson(List.of(shared, shared)));
    }

    @Test
    void rejectsTooDeepValues() {
        List<Object> root = new ArrayList<>();
        List<Object> current = root;
        for (int i = 0; i < Json.MAX_DEPTH; i++) {
            List<Object> next = new ArrayList<>();
            current.add(next);
            current = next;
        }
        assertThrows(JsonException.class, () -> Json.toJson(root));
    }

    @Test
    void exceptionConstructors() {
        JsonException plain = new JsonException("m");
        assertEquals("m", plain.getMessage());
        assertNull(plain.getCause());
    }

}
