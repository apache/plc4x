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
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Recursive-descent parser for RFC 8259 JSON, see {@link Json#parse(Reader)}.
 */
final class JsonParser {

    private static final int EOF = -1;

    private final Reader reader;
    private int current;
    private long position;
    private int depth;

    JsonParser(Reader reader) throws IOException {
        this.reader = reader;
        this.position = -1;
        advance();
    }

    Object parseDocument() throws IOException {
        skipWhitespace();
        Object value = parseValue();
        skipWhitespace();
        if (current != EOF) {
            throw error("Unexpected content after the JSON value");
        }
        return value;
    }

    private Object parseValue() throws IOException {
        return switch (current) {
            case '{' -> parseObject();
            case '[' -> parseArray();
            case '"' -> parseString();
            case 't' -> parseLiteral("true", Boolean.TRUE);
            case 'f' -> parseLiteral("false", Boolean.FALSE);
            case 'n' -> parseLiteral("null", null);
            case EOF -> throw error("Unexpected end of input");
            default -> {
                if (current == '-' || isDigit(current)) {
                    yield parseNumber();
                }
                throw error("Unexpected character '" + (char) current + "'");
            }
        };
    }

    private Map<String, Object> parseObject() throws IOException {
        enter();
        advance(); // {
        Map<String, Object> object = new LinkedHashMap<>();
        skipWhitespace();
        if (current == '}') {
            advance();
            depth--;
            return object;
        }
        while (true) {
            if (current != '"') {
                throw error("Expected a member name");
            }
            String name = parseString();
            skipWhitespace();
            expect(':');
            skipWhitespace();
            object.put(name, parseValue());
            skipWhitespace();
            if (current == ',') {
                advance();
                skipWhitespace();
            } else if (current == '}') {
                advance();
                depth--;
                return object;
            } else {
                throw error("Expected ',' or '}'");
            }
        }
    }

    private List<Object> parseArray() throws IOException {
        enter();
        advance(); // [
        List<Object> array = new ArrayList<>();
        skipWhitespace();
        if (current == ']') {
            advance();
            depth--;
            return array;
        }
        while (true) {
            array.add(parseValue());
            skipWhitespace();
            if (current == ',') {
                advance();
                skipWhitespace();
            } else if (current == ']') {
                advance();
                depth--;
                return array;
            } else {
                throw error("Expected ',' or ']'");
            }
        }
    }

    private String parseString() throws IOException {
        advance(); // opening quote
        StringBuilder sb = new StringBuilder();
        while (true) {
            if (current == '"') {
                advance();
                return sb.toString();
            }
            if (current == EOF) {
                throw error("Unterminated string");
            }
            if (current < 0x20) {
                throw error("Unescaped control character in string");
            }
            if (current == '\\') {
                advance();
                switch (current) {
                    case '"' -> sb.append('"');
                    case '\\' -> sb.append('\\');
                    case '/' -> sb.append('/');
                    case 'b' -> sb.append('\b');
                    case 'f' -> sb.append('\f');
                    case 'n' -> sb.append('\n');
                    case 'r' -> sb.append('\r');
                    case 't' -> sb.append('\t');
                    case 'u' -> sb.append(parseUnicodeEscape());
                    default -> throw error("Invalid escape sequence");
                }
            } else {
                sb.append((char) current);
            }
            advance();
        }
    }

    /**
     * Reads the four hex digits of a unicode escape (backslash, 'u', four digits) and leaves the reader on the last of them.
     * Surrogate pairs come as two escapes and end up as two chars, which is what Java strings
     * want.
     */
    private char parseUnicodeEscape() throws IOException {
        int value = 0;
        for (int i = 0; i < 4; i++) {
            advance();
            int digit = current == EOF ? -1 : Character.digit(current, 16);
            if (digit < 0) {
                throw error("Invalid unicode escape");
            }
            value = (value << 4) | digit;
        }
        return (char) value;
    }

    private Number parseNumber() throws IOException {
        StringBuilder sb = new StringBuilder();
        boolean integer = true;
        if (current == '-') {
            sb.append('-');
            advance();
        }
        if (current == '0') {
            sb.append('0');
            advance();
        } else if (isDigit(current)) {
            appendDigits(sb);
        } else {
            throw error("Expected a digit");
        }
        if (current == '.') {
            integer = false;
            sb.append('.');
            advance();
            if (!isDigit(current)) {
                throw error("Expected a digit after the decimal point");
            }
            appendDigits(sb);
        }
        if (current == 'e' || current == 'E') {
            integer = false;
            sb.append('e');
            advance();
            if (current == '+' || current == '-') {
                sb.append((char) current);
                advance();
            }
            if (!isDigit(current)) {
                throw error("Expected a digit in the exponent");
            }
            appendDigits(sb);
        }
        String text = sb.toString();
        if (!integer) {
            return Double.parseDouble(text);
        }
        // 18 digits always fit into a long, only longer numbers need the slow check.
        if (text.length() <= 18) {
            return Long.parseLong(text);
        }
        BigInteger big = new BigInteger(text);
        return big.bitLength() < 64 ? (Number) big.longValue() : big;
    }

    private void appendDigits(StringBuilder sb) throws IOException {
        while (isDigit(current)) {
            sb.append((char) current);
            advance();
        }
    }

    private Object parseLiteral(String literal, Object value) throws IOException {
        for (int i = 0; i < literal.length(); i++) {
            if (current != literal.charAt(i)) {
                throw error("Invalid literal, expected '" + literal + "'");
            }
            advance();
        }
        return value;
    }

    private void enter() throws JsonException {
        if (++depth > Json.MAX_DEPTH) {
            throw error("Nesting deeper than " + Json.MAX_DEPTH + " levels");
        }
    }

    private void expect(char c) throws IOException {
        if (current != c) {
            throw error("Expected '" + c + "'");
        }
        advance();
    }

    private void skipWhitespace() throws IOException {
        while (current == ' ' || current == '\t' || current == '\n' || current == '\r') {
            advance();
        }
    }

    private void advance() throws IOException {
        current = reader.read();
        position++;
    }

    private static boolean isDigit(int c) {
        return c >= '0' && c <= '9';
    }

    private JsonException error(String message) {
        return new JsonException(message + " at offset " + position);
    }

}
