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
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class SimpleExpressionEvaluatorTest {

    private SimpleExpressionEvaluator evaluator;

    @BeforeEach
    void setUp() {
        evaluator = new SimpleExpressionEvaluator();
    }

    @Test
    void testGetName() {
        assertEquals("simple", evaluator.getName());
    }

    @Test
    void testSimpleVariable() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(42));

        PlcValue result = evaluator.transform("value", context);
        assertEquals(42, result.getInteger());
    }

    @Test
    void testArithmeticAddition() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("value + 5", context);
        assertEquals(15L, result.getLong());
    }

    @Test
    void testArithmeticSubtraction() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("value - 3", context);
        assertEquals(7L, result.getLong());
    }

    @Test
    void testArithmeticMultiplication() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("value * 2", context);
        assertEquals(20L, result.getLong());
    }

    @Test
    void testArithmeticDivision() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("value / 2", context);
        assertEquals(5.0, result.getDouble(), 0.001);
    }

    @Test
    void testArithmeticModulo() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("value % 3", context);
        assertEquals(1L, result.getLong());
    }

    @Test
    void testFloatingPointArithmetic() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcREAL(100.0));

        // Celsius to Fahrenheit: (C * 1.8) + 32
        PlcValue result = evaluator.transform("value * 1.8 + 32", context);
        assertEquals(212.0, result.getDouble(), 0.001);
    }

    @Test
    void testParentheses() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("(value + 5) * 2", context);
        assertEquals(30L, result.getLong());
    }

    @Test
    void testComparison() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result1 = evaluator.transform("value > 5", context);
        assertTrue(result1.getBoolean());

        PlcValue result2 = evaluator.transform("value < 5", context);
        assertFalse(result2.getBoolean());

        PlcValue result3 = evaluator.transform("value == 10", context);
        assertTrue(result3.getBoolean());

        PlcValue result4 = evaluator.transform("value != 10", context);
        assertFalse(result4.getBoolean());
    }

    @Test
    void testLogicalOperators() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result1 = evaluator.transform("value > 5 && value < 15", context);
        assertTrue(result1.getBoolean());

        PlcValue result2 = evaluator.transform("value > 15 || value < 5", context);
        assertFalse(result2.getBoolean());

        PlcValue result3 = evaluator.transform("!(value > 15)", context);
        assertTrue(result3.getBoolean());
    }

    @Test
    void testUnaryNegation() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("-value", context);
        assertEquals(-10L, result.getLong());
    }

    @Test
    void testBooleanLiterals() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();

        PlcValue result1 = evaluator.transform("true", context);
        assertTrue(result1.getBoolean());

        PlcValue result2 = evaluator.transform("false", context);
        assertFalse(result2.getBoolean());
    }

    @Test
    void testMultipleVariables() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("temperature", new PlcDINT(25));
        context.put("humidity", new PlcDINT(60));

        PlcValue result = evaluator.transform("temperature + humidity", context);
        assertEquals(85L, result.getLong());
    }

    @Test
    void testComplexExpression() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("pressure", new PlcREAL(1050.0));

        // Pressure deviation from standard: (P - 1013.25) / 1013.25 * 100
        PlcValue result = evaluator.transform("(pressure - 1013.25) / 1013.25 * 100", context);
        assertEquals(3.626, result.getDouble(), 0.001);
    }

    @Test
    void testNullExpression() {
        Map<String, PlcValue> context = new HashMap<>();

        assertThrows(TransformException.class, () -> {
            evaluator.transform(null, context);
        });
    }

    @Test
    void testEmptyExpression() {
        Map<String, PlcValue> context = new HashMap<>();

        assertThrows(TransformException.class, () -> {
            evaluator.transform("", context);
        });
    }

    @Test
    void testUnknownVariable() {
        Map<String, PlcValue> context = new HashMap<>();

        assertThrows(TransformException.class, () -> {
            evaluator.transform("unknownVar + 5", context);
        });
    }

    @Test
    void testDivisionByZero() {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        assertThrows(TransformException.class, () -> {
            evaluator.transform("value / 0", context);
        });
    }

    @Test
    void testModuloByZero() {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        assertThrows(TransformException.class, () -> {
            evaluator.transform("value % 0", context);
        });
    }

    @Test
    void testWhitespace() throws TransformException {
        Map<String, PlcValue> context = new HashMap<>();
        context.put("value", new PlcDINT(10));

        PlcValue result = evaluator.transform("  value   +   5  ", context);
        assertEquals(15L, result.getLong());
    }
}
