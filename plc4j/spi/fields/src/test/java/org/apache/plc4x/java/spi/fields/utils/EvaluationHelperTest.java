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

package org.apache.plc4x.java.spi.fields.utils;

import org.apache.plc4x.java.spi.fields.utils.EvaluationHelper;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class EvaluationHelperTest {

    private enum Color { RED, GREEN }
    private enum Size { SMALL, LARGE }

    @Test
    void bothNullAreEqual() {
        assertTrue(EvaluationHelper.equals(null, null));
    }

    @Test
    void oneNullIsNotEqual() {
        assertFalse(EvaluationHelper.equals(null, 1));
        assertFalse(EvaluationHelper.equals("x", null));
    }

    @Test
    void numbersComparedByDoubleValue() {
        // Equal when numerically the same
        assertTrue(EvaluationHelper.equals(1, 1L));
        assertTrue(EvaluationHelper.equals(1, 1.0));
        assertTrue(EvaluationHelper.equals(0.0d, -0.0f));

        // Not equal when numerically different
        assertFalse(EvaluationHelper.equals(1, 2L));
        assertFalse(EvaluationHelper.equals(1.0, 1.1));
    }

    @Test
    void booleansComparedByEquality() {
        assertTrue(EvaluationHelper.equals(true, Boolean.TRUE));
        assertFalse(EvaluationHelper.equals(true, false));

        // Mismatch type should be false
        assertFalse(EvaluationHelper.equals(true, 1));
    }

    @Test
    void stringsComparedByEquals() {
        assertTrue(EvaluationHelper.equals("abc", "abc"));
        assertFalse(EvaluationHelper.equals("abc", "abd"));

        // Mismatch type should be false
        assertFalse(EvaluationHelper.equals("1", 1));
    }

    @Test
    void enumsComparedWhenBothEnums() {
        // Same enum type and same constant -> true
        assertTrue(EvaluationHelper.equals(Color.RED, Color.RED));
        // Same enum type different constant -> false
        assertFalse(EvaluationHelper.equals(Color.RED, Color.GREEN));
        // Different enum types (both enums) -> false
        assertFalse(EvaluationHelper.equals(Color.RED, Size.SMALL));
    }

    @Test
    void unsupportedTypesReturnFalse() {
        record Point(int x, int y) {}
        assertFalse(EvaluationHelper.equals(new Object(), new Object()));
        assertFalse(EvaluationHelper.equals(new Point(1, 2), new Point(1, 2)));
    }
}
