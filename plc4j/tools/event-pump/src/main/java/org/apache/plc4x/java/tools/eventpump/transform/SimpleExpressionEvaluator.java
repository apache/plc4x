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

package org.apache.plc4x.java.tools.eventpump.transform;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

/**
 * Simple expression evaluator supporting basic arithmetic and comparisons.
 * <p>
 * Supported operations:
 * - Arithmetic: +, -, *, /, % (modulo)
 * - Comparisons: ==, !=, <, >, <=, >=
 * - Logical: &&, ||, ! (not)
 * - Parentheses: ( )
 * - Variables: Any identifier (resolved from context)
 * - Literals: Numeric (integer and floating point)
 * <p>
 * Examples:
 * - "value * 1.8 + 32" - Convert Celsius to Fahrenheit
 * - "temperature > 100" - Temperature check
 * - "(pressure - 1013.25) / 1013.25 * 100" - Pressure deviation percentage
 * - "active && temperature > 50" - Combined condition
 */
public class SimpleExpressionEvaluator implements ValueTransformer {

    private static final Logger LOGGER = LoggerFactory.getLogger(SimpleExpressionEvaluator.class);

    @Override
    public String getName() {
        return "simple";
    }

    @Override
    public PlcValue transform(String expression, Map<String, PlcValue> context) throws TransformException {
        if (expression == null || expression.trim().isEmpty()) {
            throw new TransformException("Expression cannot be null or empty");
        }

        try {
            Parser parser = new Parser(expression, context);
            return parser.parse();
        } catch (Exception e) {
            throw new TransformException("Failed to evaluate expression: " + expression, e);
        }
    }

    /**
     * Simple recursive descent parser for expressions.
     */
    private static class Parser {
        private final String expression;
        private final Map<String, PlcValue> context;
        private int pos = 0;

        Parser(String expression, Map<String, PlcValue> context) {
            this.expression = expression.trim();
            this.context = context;
        }

        PlcValue parse() throws TransformException {
            PlcValue result = parseLogicalOr();
            if (pos < expression.length()) {
                throw new TransformException("Unexpected character at position " + pos + ": " + expression.charAt(pos));
            }
            return result;
        }

        // Logical OR (lowest precedence)
        private PlcValue parseLogicalOr() throws TransformException {
            PlcValue left = parseLogicalAnd();

            while (pos < expression.length() && peek("||")) {
                consume("||");
                PlcValue right = parseLogicalAnd();
                left = new PlcBOOL(toBoolean(left) || toBoolean(right));
            }

            return left;
        }

        // Logical AND
        private PlcValue parseLogicalAnd() throws TransformException {
            PlcValue left = parseComparison();

            while (pos < expression.length() && peek("&&")) {
                consume("&&");
                PlcValue right = parseComparison();
                left = new PlcBOOL(toBoolean(left) && toBoolean(right));
            }

            return left;
        }

        // Comparison operators
        private PlcValue parseComparison() throws TransformException {
            PlcValue left = parseAddSub();

            while (pos < expression.length()) {
                if (peek("==")) {
                    consume("==");
                    PlcValue right = parseAddSub();
                    left = new PlcBOOL(compare(left, right) == 0);
                } else if (peek("!=")) {
                    consume("!=");
                    PlcValue right = parseAddSub();
                    left = new PlcBOOL(compare(left, right) != 0);
                } else if (peek("<=")) {
                    consume("<=");
                    PlcValue right = parseAddSub();
                    left = new PlcBOOL(compare(left, right) <= 0);
                } else if (peek(">=")) {
                    consume(">=");
                    PlcValue right = parseAddSub();
                    left = new PlcBOOL(compare(left, right) >= 0);
                } else if (peek("<")) {
                    consume("<");
                    PlcValue right = parseAddSub();
                    left = new PlcBOOL(compare(left, right) < 0);
                } else if (peek(">")) {
                    consume(">");
                    PlcValue right = parseAddSub();
                    left = new PlcBOOL(compare(left, right) > 0);
                } else {
                    break;
                }
            }

            return left;
        }

        // Addition and subtraction
        private PlcValue parseAddSub() throws TransformException {
            PlcValue left = parseMulDiv();

            while (pos < expression.length()) {
                skipWhitespace();
                if (pos >= expression.length()) break;

                char op = expression.charAt(pos);
                if (op == '+') {
                    pos++;
                    PlcValue right = parseMulDiv();
                    left = add(left, right);
                } else if (op == '-' && (pos == 0 || !Character.isDigit(expression.charAt(pos - 1)))) {
                    // Only treat as subtraction if not part of a number
                    if (pos > 0 && Character.isWhitespace(expression.charAt(pos - 1))) {
                        pos++;
                        PlcValue right = parseMulDiv();
                        left = subtract(left, right);
                    } else {
                        break;
                    }
                } else {
                    break;
                }
            }

            return left;
        }

        // Multiplication, division, and modulo
        private PlcValue parseMulDiv() throws TransformException {
            PlcValue left = parseUnary();

            while (pos < expression.length()) {
                skipWhitespace();
                if (pos >= expression.length()) break;

                char op = expression.charAt(pos);
                if (op == '*') {
                    pos++;
                    PlcValue right = parseUnary();
                    left = multiply(left, right);
                } else if (op == '/') {
                    pos++;
                    PlcValue right = parseUnary();
                    left = divide(left, right);
                } else if (op == '%') {
                    pos++;
                    PlcValue right = parseUnary();
                    left = modulo(left, right);
                } else {
                    break;
                }
            }

            return left;
        }

        // Unary operators (-, !)
        private PlcValue parseUnary() throws TransformException {
            skipWhitespace();

            if (pos < expression.length()) {
                char ch = expression.charAt(pos);
                if (ch == '-') {
                    pos++;
                    PlcValue value = parseUnary();
                    return negate(value);
                } else if (ch == '!') {
                    pos++;
                    PlcValue value = parseUnary();
                    return new PlcBOOL(!toBoolean(value));
                }
            }

            return parsePrimary();
        }

        // Primary expressions (numbers, variables, parentheses)
        private PlcValue parsePrimary() throws TransformException {
            skipWhitespace();

            if (pos >= expression.length()) {
                throw new TransformException("Unexpected end of expression");
            }

            char ch = expression.charAt(pos);

            // Parentheses
            if (ch == '(') {
                pos++;
                PlcValue value = parseLogicalOr();
                skipWhitespace();
                if (pos >= expression.length() || expression.charAt(pos) != ')') {
                    throw new TransformException("Expected ')' at position " + pos);
                }
                pos++;
                return value;
            }

            // Number
            if (Character.isDigit(ch) || ch == '.') {
                return parseNumber();
            }

            // Variable
            if (Character.isLetter(ch) || ch == '_') {
                return parseVariable();
            }

            throw new TransformException("Unexpected character at position " + pos + ": " + ch);
        }

        private PlcValue parseNumber() throws TransformException {
            int start = pos;
            boolean hasDecimal = false;

            while (pos < expression.length()) {
                char ch = expression.charAt(pos);
                if (Character.isDigit(ch)) {
                    pos++;
                } else if (ch == '.' && !hasDecimal) {
                    hasDecimal = true;
                    pos++;
                } else {
                    break;
                }
            }

            String numStr = expression.substring(start, pos);
            try {
                if (hasDecimal) {
                    return new PlcREAL(Double.parseDouble(numStr));
                } else {
                    long value = Long.parseLong(numStr);
                    if (value >= Integer.MIN_VALUE && value <= Integer.MAX_VALUE) {
                        return new PlcDINT((int) value);
                    } else {
                        return new PlcLINT(value);
                    }
                }
            } catch (NumberFormatException e) {
                throw new TransformException("Invalid number: " + numStr);
            }
        }

        private PlcValue parseVariable() throws TransformException {
            int start = pos;

            while (pos < expression.length()) {
                char ch = expression.charAt(pos);
                if (Character.isLetterOrDigit(ch) || ch == '_') {
                    pos++;
                } else {
                    break;
                }
            }

            String varName = expression.substring(start, pos);

            // Check for boolean literals
            if ("true".equalsIgnoreCase(varName)) {
                return new PlcBOOL(true);
            } else if ("false".equalsIgnoreCase(varName)) {
                return new PlcBOOL(false);
            }

            PlcValue value = context.get(varName);
            if (value == null) {
                throw new TransformException("Unknown variable: " + varName);
            }

            return value;
        }

        private void skipWhitespace() {
            while (pos < expression.length() && Character.isWhitespace(expression.charAt(pos))) {
                pos++;
            }
        }

        private boolean peek(String str) {
            skipWhitespace();
            if (pos + str.length() > expression.length()) {
                return false;
            }
            return expression.startsWith(str, pos);
        }

        private void consume(String str) {
            skipWhitespace();
            if (!peek(str)) {
                throw new RuntimeException("Expected '" + str + "' at position " + pos);
            }
            pos += str.length();
        }

        // Arithmetic operations
        private PlcValue add(PlcValue a, PlcValue b) throws TransformException {
            if (isFloatingPoint(a) || isFloatingPoint(b)) {
                return new PlcREAL(toDouble(a) + toDouble(b));
            } else {
                return new PlcLINT(toLong(a) + toLong(b));
            }
        }

        private PlcValue subtract(PlcValue a, PlcValue b) throws TransformException {
            if (isFloatingPoint(a) || isFloatingPoint(b)) {
                return new PlcREAL(toDouble(a) - toDouble(b));
            } else {
                return new PlcLINT(toLong(a) - toLong(b));
            }
        }

        private PlcValue multiply(PlcValue a, PlcValue b) throws TransformException {
            if (isFloatingPoint(a) || isFloatingPoint(b)) {
                return new PlcREAL(toDouble(a) * toDouble(b));
            } else {
                return new PlcLINT(toLong(a) * toLong(b));
            }
        }

        private PlcValue divide(PlcValue a, PlcValue b) throws TransformException {
            double divisor = toDouble(b);
            if (divisor == 0) {
                throw new TransformException("Division by zero");
            }
            return new PlcREAL(toDouble(a) / divisor);
        }

        private PlcValue modulo(PlcValue a, PlcValue b) throws TransformException {
            long divisor = toLong(b);
            if (divisor == 0) {
                throw new TransformException("Modulo by zero");
            }
            return new PlcLINT(toLong(a) % divisor);
        }

        private PlcValue negate(PlcValue value) throws TransformException {
            if (isFloatingPoint(value)) {
                return new PlcREAL(-toDouble(value));
            } else {
                return new PlcLINT(-toLong(value));
            }
        }

        private int compare(PlcValue a, PlcValue b) throws TransformException {
            if (isFloatingPoint(a) || isFloatingPoint(b)) {
                return Double.compare(toDouble(a), toDouble(b));
            } else {
                return Long.compare(toLong(a), toLong(b));
            }
        }

        // Type conversion helpers
        private boolean isFloatingPoint(PlcValue value) {
            return value instanceof PlcREAL || value instanceof PlcLREAL;
        }

        private double toDouble(PlcValue value) throws TransformException {
            if (value.isDouble()) {
                return value.getDouble();
            } else if (value.isFloat()) {
                return value.getFloat();
            } else if (value.isLong()) {
                return value.getLong();
            } else if (value.isInteger()) {
                return value.getInteger();
            } else if (value.isShort()) {
                return value.getShort();
            } else if (value.isByte()) {
                return value.getByte();
            } else if (value.isBigInteger()) {
                return value.getBigInteger().doubleValue();
            } else if (value.isBigDecimal()) {
                return value.getBigDecimal().doubleValue();
            }
            throw new TransformException("Cannot convert to double: " + value.getClass().getSimpleName());
        }

        private long toLong(PlcValue value) throws TransformException {
            if (value.isLong()) {
                return value.getLong();
            } else if (value.isInteger()) {
                return value.getInteger();
            } else if (value.isShort()) {
                return value.getShort();
            } else if (value.isByte()) {
                return value.getByte();
            } else if (value.isBigInteger()) {
                return value.getBigInteger().longValue();
            } else if (value.isDouble()) {
                return (long) value.getDouble();
            } else if (value.isFloat()) {
                return (long) value.getFloat();
            }
            throw new TransformException("Cannot convert to long: " + value.getClass().getSimpleName());
        }

        private boolean toBoolean(PlcValue value) throws TransformException {
            if (value.isBoolean()) {
                return value.getBoolean();
            }
            // Non-zero numbers are true
            if (value.isLong() || value.isInteger() || value.isShort() || value.isByte()) {
                return toLong(value) != 0;
            }
            if (value.isDouble() || value.isFloat()) {
                return toDouble(value) != 0.0;
            }
            throw new TransformException("Cannot convert to boolean: " + value.getClass().getSimpleName());
        }
    }
}
