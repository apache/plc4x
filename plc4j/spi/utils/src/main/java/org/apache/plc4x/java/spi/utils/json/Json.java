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

import java.io.IOException;
import java.io.Reader;
import java.io.StringReader;

/**
 * Minimal JSON support without third-party libraries (RFC 8259).
 * <p>
 * Parsing produces a plain Java tree:
 * <ul>
 *     <li>object: {@link java.util.LinkedHashMap LinkedHashMap&lt;String, Object&gt;} (member order kept)</li>
 *     <li>array: {@link java.util.ArrayList ArrayList&lt;Object&gt;}</li>
 *     <li>string: {@link String}</li>
 *     <li>number: {@link Long} for integers that fit, {@link java.math.BigInteger} for larger
 *     integers, {@link Double} for numbers with a fraction or exponent</li>
 *     <li>{@code true}/{@code false}: {@link Boolean}</li>
 *     <li>{@code null}: {@code null}</li>
 * </ul>
 * Writing accepts such a tree, and also records, enums, arrays, {@link Iterable}s and plain Java
 * beans (public {@code getX()}/{@code isX()} getters), see {@link #toJson(Object)}.
 */
public final class Json {

    /**
     * Maximum nesting depth of arrays/objects, both when parsing and when writing. Keeps
     * hostile input (e.g. a network response) from overflowing the stack.
     */
    public static final int MAX_DEPTH = 512;

    private Json() {
    }

    /**
     * Parses a complete JSON text.
     *
     * @param json the JSON text
     * @return the parsed value, see the class documentation for the types
     * @throws JsonException if the text isn't valid JSON
     */
    public static Object parse(String json) throws JsonException {
        try {
            return parse(new StringReader(json));
        } catch (JsonException e) {
            throw e;
        } catch (IOException e) {
            // A StringReader doesn't throw.
            throw new JsonException("Error reading JSON", e);
        }
    }

    /**
     * Parses a complete JSON text read from the given reader. The reader is not closed.
     *
     * @param reader source of the JSON text
     * @return the parsed value, see the class documentation for the types
     * @throws JsonException if the text isn't valid JSON
     * @throws IOException   if reading fails
     */
    public static Object parse(Reader reader) throws IOException {
        return new JsonParser(reader).parseDocument();
    }

    /**
     * Writes a value as compact JSON.
     * <p>
     * Supported values: {@code null}, {@link CharSequence}, {@link Character}, {@link Number},
     * {@link Boolean}, {@link Enum} (its name), {@link java.util.Map} (keys via
     * {@link String#valueOf(Object)}), {@link Iterable}, arrays ({@code byte[]} as Base64, like
     * Jackson), {@link java.util.Optional}, records (their components), value types that are
     * best represented by their {@code toString()} (java.time types, UUID, URI, URL, Path, File,
     * Locale, Charset) and other objects as beans via their public getters. Non-finite floating
     * point numbers are written as strings, as JSON has no representation for them.
     *
     * @param value the value
     * @return the JSON text
     * @throws JsonException if the value contains a reference cycle, nests too deep or a getter
     *                       throws
     */
    public static String toJson(Object value) throws JsonException {
        return JsonWriter.write(value, false);
    }

    /**
     * Like {@link #toJson(Object)}, but indented with two spaces per level, for files meant to be
     * read by people.
     *
     * @param value the value
     * @return the JSON text
     * @throws JsonException see {@link #toJson(Object)}
     */
    public static String toPrettyJson(Object value) throws JsonException {
        return JsonWriter.write(value, true);
    }

}
