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
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcIECValue - base IEC value class
 */
class PlcIECValueTest {

    // Concrete implementation for testing the abstract PlcIECValue
    private static class TestPlcIECValue extends PlcIECValue<String> {
        public TestPlcIECValue(String value) {
            this.value = value;
            this.isNullable = value == null;
        }

        @Override
        public PlcValueType getPlcValueType() {
            return PlcValueType.STRING;
        }
    }

    // ========== Object and Basic Tests ==========

    @Test
    void testGetObject_returnsValue() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertEquals("test", value.getObject());
    }

    @Test
    void testGetObject_nullValue() {
        TestPlcIECValue value = new TestPlcIECValue(null);
        assertNull(value.getObject());
    }

    @Test
    void testGetLength_returnsOne() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertEquals(1, value.getLength());
    }

    @Test
    void testIsSimple_returnsTrue() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertTrue(value.isSimple());
    }

    // ========== Nullable and Null Tests ==========

    @Test
    void testIsNullable_withNonNullValue() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertFalse(value.isNullable());
    }

    @Test
    void testIsNullable_withNullValue() {
        TestPlcIECValue value = new TestPlcIECValue(null);
        assertTrue(value.isNullable());
    }

    @Test
    void testIsNull_withNonNullValue() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertFalse(value.isNull());
    }

    @Test
    void testIsNull_withNullValue() {
        TestPlcIECValue value = new TestPlcIECValue(null);
        assertTrue(value.isNull());
    }

    // ========== GetIndex Tests ==========

    @Test
    void testGetIndex_zero_returnsSelf() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertEquals(value, value.getIndex(0));
    }

    @Test
    void testGetIndex_nonZero_throwsException() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertThrows(PlcIncompatibleDatatypeException.class, () -> value.getIndex(1));
        assertThrows(PlcIncompatibleDatatypeException.class, () -> value.getIndex(-1));
        assertThrows(PlcIncompatibleDatatypeException.class, () -> value.getIndex(5));
    }

    // ========== Equals and HashCode Tests ==========

    @Test
    void testEquals_sameInstance() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertEquals(value, value);
    }

    @Test
    void testEquals_nullOther() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertNotEquals(null, value);
    }

    @Test
    void testEquals_differentClass() {
        TestPlcIECValue value = new TestPlcIECValue("test");
        assertNotEquals("test", value);
    }

    @Test
    void testEquals_sameValue() {
        TestPlcIECValue value1 = new TestPlcIECValue("test");
        TestPlcIECValue value2 = new TestPlcIECValue("test");
        assertEquals(value1, value2);
    }

    @Test
    void testEquals_differentValue() {
        TestPlcIECValue value1 = new TestPlcIECValue("test1");
        TestPlcIECValue value2 = new TestPlcIECValue("test2");
        assertNotEquals(value1, value2);
    }

    @Test
    void testEquals_bothNull() {
        TestPlcIECValue value1 = new TestPlcIECValue(null);
        TestPlcIECValue value2 = new TestPlcIECValue(null);
        assertEquals(value1, value2);
    }

    @Test
    void testEquals_oneNull() {
        TestPlcIECValue value1 = new TestPlcIECValue("test");
        TestPlcIECValue value2 = new TestPlcIECValue(null);
        assertNotEquals(value1, value2);
    }

    @Test
    void testHashCode_sameValue() {
        TestPlcIECValue value1 = new TestPlcIECValue("test");
        TestPlcIECValue value2 = new TestPlcIECValue("test");
        assertEquals(value1.hashCode(), value2.hashCode());
    }

    @Test
    void testHashCode_nullValue() {
        TestPlcIECValue value1 = new TestPlcIECValue(null);
        TestPlcIECValue value2 = new TestPlcIECValue(null);
        assertEquals(value1.hashCode(), value2.hashCode());
    }

    // ========== Default Constructor Test ==========

    @Test
    void testDefaultConstructor() {
        // Test the default no-args constructor through a subclass
        PlcIECValue<Integer> value = new PlcIECValue<Integer>() {
            @Override
            public PlcValueType getPlcValueType() {
                return PlcValueType.INT;
            }
        };

        assertNull(value.value);
        assertTrue(value.isNullable);
        assertTrue(value.isNull());
        assertTrue(value.isNullable());
    }
}
