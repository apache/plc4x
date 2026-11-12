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

package org.apache.plc4x.java.tools.eventpump.config;

import java.io.IOException;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Binds the plain tree {@link ConfigFiles} reads (maps, lists, scalars) to the configuration
 * classes, with errors that name the offending place, e.g. {@code batches[0].trigger.intervalMillis}.
 * <p>
 * Lenient where a format can't say better, strict where a typo would otherwise go unnoticed:
 * <ul>
 *     <li>A text value may be written as any scalar (YAML reads {@code tagName: 123} as a number).</li>
 *     <li>A whole number may be written as text (XML has nothing else).</li>
 *     <li>A list in XML is a wrapper element holding the entries (see {@link ConfigFiles}).</li>
 *     <li>A property the configuration doesn't know is an error.</li>
 * </ul>
 */
final class ConfigNode {

    private final Map<String, Object> map;
    private final String path;

    private ConfigNode(Map<String, Object> map, String path) {
        this.map = map;
        this.path = path;
    }

    /**
     * The object at the given place, which may only use the given property names.
     */
    static ConfigNode object(Object node, String path, String... allowedKeys) throws IOException {
        if (!(node instanceof Map<?, ?> raw)) {
            throw new IOException(where(path) + "expected an object, got " + describe(node));
        }
        Map<String, Object> map = new LinkedHashMap<>();
        Set<String> allowed = new HashSet<>(Arrays.asList(allowedKeys));
        for (Map.Entry<?, ?> entry : raw.entrySet()) {
            String key = String.valueOf(entry.getKey());
            if (!allowed.contains(key)) {
                throw new IOException(where(path) + "unknown property '" + key + "', expected one of "
                    + Arrays.toString(allowedKeys));
            }
            map.put(key, entry.getValue());
        }
        return new ConfigNode(map, path);
    }

    String path(String key) {
        return path.isEmpty() ? key : path + "." + key;
    }

    Object get(String key) {
        return map.get(key);
    }

    String string(String key) throws IOException {
        return scalarToString(map.get(key), path(key));
    }

    Long longValue(String key) throws IOException {
        Object value = map.get(key);
        String at = path(key);
        return switch (value) {
            case null -> null;
            case Long l -> l;
            case Integer i -> i.longValue();
            case Short s -> s.longValue();
            case Byte b -> b.longValue();
            case BigInteger big when big.bitLength() < 64 -> big.longValue();
            case String s -> {
                try {
                    yield Long.parseLong(s.trim());
                } catch (NumberFormatException e) {
                    throw new IOException(where(at) + "expected a whole number, got '" + s + "'");
                }
            }
            default -> throw new IOException(where(at) + "expected a whole number, got " + describe(value));
        };
    }

    /**
     * A list; absent or empty means an empty list. In XML, the wrapper element's entries.
     */
    List<Object> list(String key) throws IOException {
        Object value = map.get(key);
        String at = path(key);
        return switch (value) {
            case null -> new ArrayList<>();
            case List<?> list -> new ArrayList<>(list);
            // An XML wrapper: one entry is a single element, several come as RepeatedElements.
            case ConfigFiles.XmlObject wrapper when wrapper.size() == 1
                && wrapper.values().iterator().next() instanceof ConfigFiles.RepeatedElements repeated ->
                new ArrayList<>(repeated);
            case ConfigFiles.XmlObject wrapper when wrapper.size() == 1 ->
                new ArrayList<>(List.of(wrapper.values().iterator().next()));
            default -> throw new IOException(where(at) + "expected a list, got " + describe(value));
        };
    }

    /**
     * The members of a map-valued property, in file order; absent or empty means none.
     */
    Map<String, Object> members(String key) throws IOException {
        Object value = map.get(key);
        Map<String, Object> members = new LinkedHashMap<>();
        if (value == null) {
            return members;
        }
        if (!(value instanceof Map<?, ?> raw)) {
            throw new IOException(where(path(key)) + "expected an object, got " + describe(value));
        }
        for (Map.Entry<?, ?> entry : raw.entrySet()) {
            members.put(String.valueOf(entry.getKey()), entry.getValue());
        }
        return members;
    }

    static String scalarToString(Object value, String path) throws IOException {
        return switch (value) {
            case null -> null;
            case String s -> s;
            case Number n -> n.toString();
            case Boolean b -> b.toString();
            default -> throw new IOException(where(path) + "expected a text, got " + describe(value));
        };
    }

    static boolean isScalar(Object value) {
        return value instanceof String || value instanceof Number || value instanceof Boolean;
    }

    private static String where(String path) {
        return path.isEmpty() ? "" : path + ": ";
    }

    private static String describe(Object value) {
        return switch (value) {
            case null -> "nothing";
            case Map<?, ?> m -> "an object";
            case List<?> l -> "a list";
            default -> "'" + value + "'";
        };
    }

    /**
     * Collects the properties of an object for writing, leaving out those without a value.
     */
    static final class Writer {
        private final Map<String, Object> tree = new LinkedHashMap<>();

        Writer put(String key, Object value) {
            if (value != null) {
                tree.put(key, value);
            }
            return this;
        }

        Map<String, Object> tree() {
            return tree;
        }
    }

}
