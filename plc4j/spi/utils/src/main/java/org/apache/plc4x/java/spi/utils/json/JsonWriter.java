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

import java.io.File;
import java.lang.reflect.Array;
import java.lang.reflect.Field;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.lang.reflect.RecordComponent;
import java.net.URI;
import java.net.URL;
import java.nio.charset.Charset;
import java.nio.file.Path;
import java.time.ZoneId;
import java.time.temporal.TemporalAccessor;
import java.time.temporal.TemporalAmount;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.IdentityHashMap;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.UUID;

/**
 * Writes Java values as JSON, see {@link Json#toJson(Object)}.
 */
final class JsonWriter {

    private final StringBuilder out = new StringBuilder();
    private final boolean pretty;
    private final Set<Object> inProgress = Collections.newSetFromMap(new IdentityHashMap<>());
    private int depth;

    private JsonWriter(boolean pretty) {
        this.pretty = pretty;
    }

    static String write(Object value, boolean pretty) throws JsonException {
        JsonWriter writer = new JsonWriter(pretty);
        writer.writeValue(value);
        return writer.out.toString();
    }

    private void writeValue(Object value) throws JsonException {
        switch (value) {
            case null -> out.append("null");
            case CharSequence s -> writeString(s.toString());
            case Character c -> writeString(c.toString());
            case Boolean b -> out.append(b);
            case Double d when d.isNaN() || d.isInfinite() -> writeString(d.toString());
            case Float f when f.isNaN() || f.isInfinite() -> writeString(f.toString());
            case Number n -> out.append(n);
            case Enum<?> e -> writeString(e.name());
            case Optional<?> o -> writeValue(o.orElse(null));
            case byte[] bytes -> writeString(Base64.getEncoder().encodeToString(bytes));
            case TemporalAccessor t -> writeString(t.toString());
            case TemporalAmount t -> writeString(t.toString());
            case ZoneId z -> writeString(z.toString());
            case UUID u -> writeString(u.toString());
            case URI u -> writeString(u.toString());
            case URL u -> writeString(u.toString());
            case Path p -> writeString(p.toString());
            case File f -> writeString(f.getPath());
            case Locale l -> writeString(l.toLanguageTag());
            case Charset c -> writeString(c.name());
            case Map<?, ?> map -> writeContainer(map, () -> writeMap(map));
            case Iterable<?> iterable -> writeContainer(iterable, () -> writeIterable(iterable.iterator()));
            default -> {
                if (value.getClass().isArray()) {
                    writeContainer(value, () -> writeArray(value));
                } else {
                    writeContainer(value, () -> writeMap(properties(value)));
                }
            }
        }
    }

    private interface ContainerBody {
        void write() throws JsonException;
    }

    private void writeContainer(Object container, ContainerBody body) throws JsonException {
        if (!inProgress.add(container)) {
            throw new JsonException("Reference cycle at " + container.getClass().getName());
        }
        if (++depth > Json.MAX_DEPTH) {
            throw new JsonException("Nesting deeper than " + Json.MAX_DEPTH + " levels");
        }
        body.write();
        depth--;
        inProgress.remove(container);
    }

    private void writeMap(Map<?, ?> map) throws JsonException {
        out.append('{');
        boolean first = true;
        for (Map.Entry<?, ?> entry : map.entrySet()) {
            separator(first);
            first = false;
            writeString(String.valueOf(entry.getKey()));
            out.append(pretty ? ": " : ":");
            writeValue(entry.getValue());
        }
        close(first, '}');
    }

    private void writeIterable(Iterator<?> iterator) throws JsonException {
        out.append('[');
        boolean first = true;
        while (iterator.hasNext()) {
            separator(first);
            first = false;
            writeValue(iterator.next());
        }
        close(first, ']');
    }

    private void writeArray(Object array) throws JsonException {
        out.append('[');
        int length = Array.getLength(array);
        for (int i = 0; i < length; i++) {
            separator(i == 0);
            writeValue(Array.get(array, i));
        }
        close(length == 0, ']');
    }

    private void separator(boolean first) {
        if (!first) {
            out.append(',');
        }
        newline(depth);
    }

    private void close(boolean empty, char bracket) {
        if (!empty) {
            newline(depth - 1);
        }
        out.append(bracket);
    }

    private void newline(int level) {
        if (pretty) {
            out.append('\n').append("  ".repeat(level));
        }
    }

    private void writeString(String s) {
        out.append('"');
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '"' -> out.append("\\\"");
                case '\\' -> out.append("\\\\");
                case '\b' -> out.append("\\b");
                case '\f' -> out.append("\\f");
                case '\n' -> out.append("\\n");
                case '\r' -> out.append("\\r");
                case '\t' -> out.append("\\t");
                default -> {
                    if (c < 0x20) {
                        out.append(String.format("\\u%04x", (int) c));
                    } else {
                        out.append(c);
                    }
                }
            }
        }
        out.append('"');
    }

    /**
     * The JSON members of a record (its components) or bean (its public getters), in declaration
     * order.
     */
    private static Map<String, Object> properties(Object value) throws JsonException {
        Class<?> type = value.getClass();
        Map<String, Object> properties = new LinkedHashMap<>();
        if (type.isRecord()) {
            for (RecordComponent component : type.getRecordComponents()) {
                properties.put(component.getName(), invoke(component.getAccessor(), value));
            }
            return properties;
        }
        Map<String, Method> getters = new LinkedHashMap<>();
        for (Method method : type.getMethods()) {
            String name = propertyName(method);
            if (name != null) {
                getters.put(name, method);
            }
        }
        // Order like the fields are declared (superclass first), the rest by name.
        List<String> order = new ArrayList<>();
        for (Class<?> c : superclassesFirst(type)) {
            for (Field field : c.getDeclaredFields()) {
                if (getters.containsKey(field.getName()) && !order.contains(field.getName())) {
                    order.add(field.getName());
                }
            }
        }
        getters.keySet().stream().filter(name -> !order.contains(name)).sorted().forEach(order::add);
        for (String name : order) {
            properties.put(name, invoke(getters.get(name), value));
        }
        return properties;
    }

    private static List<Class<?>> superclassesFirst(Class<?> type) {
        List<Class<?>> classes = new ArrayList<>();
        for (Class<?> c = type; c != null && c != Object.class; c = c.getSuperclass()) {
            classes.addFirst(c);
        }
        return classes;
    }

    /**
     * Bean property name of a getter like Jackson derives it ("getURL" -> "url",
     * "getName" -> "name", "isEnabled" -> "enabled"), or null if the method isn't a getter.
     */
    static String propertyName(Method method) {
        if (Modifier.isStatic(method.getModifiers()) || method.getParameterCount() != 0
            || method.getReturnType() == void.class || method.getDeclaringClass() == Object.class) {
            return null;
        }
        String name = method.getName();
        String suffix;
        if (name.startsWith("get") && name.length() > 3) {
            suffix = name.substring(3);
        } else if (name.startsWith("is") && name.length() > 2
            && (method.getReturnType() == boolean.class || method.getReturnType() == Boolean.class)) {
            suffix = name.substring(2);
        } else {
            return null;
        }
        if (!Character.isUpperCase(suffix.charAt(0))) {
            return null;
        }
        // Lower-case the leading run of upper-case letters.
        char[] chars = suffix.toCharArray();
        for (int i = 0; i < chars.length && Character.isUpperCase(chars[i]); i++) {
            chars[i] = Character.toLowerCase(chars[i]);
        }
        return new String(chars);
    }

    private static Object invoke(Method method, Object target) throws JsonException {
        try {
            // Public methods of non-public classes (e.g. a private nested class) need this.
            method.trySetAccessible();
            return method.invoke(target);
        } catch (IllegalAccessException e) {
            throw new JsonException("Can't access " + method, e);
        } catch (InvocationTargetException e) {
            throw new JsonException(method + " failed", e.getCause());
        }
    }

}
