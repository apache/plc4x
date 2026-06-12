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
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcValueAdapter - base adapter class with default implementations
 */
class PlcValueAdapterTest {

    // Concrete implementation for testing the abstract PlcValueAdapter
    private static class TestPlcValueAdapter extends PlcValueAdapter {
        @Override
        public org.apache.plc4x.java.api.types.PlcValueType getPlcValueType() {
            return org.apache.plc4x.java.api.types.PlcValueType.NULL;
        }

        @Override
        public void serialize(WriteBuffer writeBuffer) throws org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException {
            // Test implementation
        }
    }

    private TestPlcValueAdapter adapter;

    @org.junit.jupiter.api.BeforeEach
    void setUp() {
        adapter = new TestPlcValueAdapter();
    }

    // ========== Basic Method Tests ==========

    @Test
    void testGetObject_returnsNull() {
        assertNull(adapter.getObject());
    }

    @Test
    void testIsSimple_returnsFalse() {
        assertFalse(adapter.isSimple());
    }

    @Test
    void testIsNullable_returnsFalse() {
        assertFalse(adapter.isNullable());
    }

    @Test
    void testIsNull_returnsFalse() {
        assertFalse(adapter.isNull());
    }

    @Test
    void testIs_returnsFalse() {
        assertFalse(adapter.is(String.class));
        assertFalse(adapter.is(Integer.class));
        assertFalse(adapter.is(Boolean.class));
    }

    @Test
    void testIsConvertibleTo_returnsFalse() {
        assertFalse(adapter.isConvertibleTo(String.class));
        assertFalse(adapter.isConvertibleTo(Integer.class));
        assertFalse(adapter.isConvertibleTo(Boolean.class));
    }

    // ========== Type Check Methods (is*) ==========

    @Test
    void testIsBoolean_returnsFalse() {
        assertFalse(adapter.isBoolean());
    }

    @Test
    void testIsByte_returnsFalse() {
        assertFalse(adapter.isByte());
    }

    @Test
    void testIsShort_returnsFalse() {
        assertFalse(adapter.isShort());
    }

    @Test
    void testIsInteger_returnsFalse() {
        assertFalse(adapter.isInteger());
    }

    @Test
    void testIsLong_returnsFalse() {
        assertFalse(adapter.isLong());
    }

    @Test
    void testIsBigInteger_returnsFalse() {
        assertFalse(adapter.isBigInteger());
    }

    @Test
    void testIsFloat_returnsFalse() {
        assertFalse(adapter.isFloat());
    }

    @Test
    void testIsDouble_returnsFalse() {
        assertFalse(adapter.isDouble());
    }

    @Test
    void testIsBigDecimal_returnsFalse() {
        assertFalse(adapter.isBigDecimal());
    }

    @Test
    void testIsString_returnsFalse() {
        assertFalse(adapter.isString());
    }

    @Test
    void testIsDuration_returnsFalse() {
        assertFalse(adapter.isDuration());
    }

    @Test
    void testIsTime_returnsFalse() {
        assertFalse(adapter.isTime());
    }

    @Test
    void testIsDate_returnsFalse() {
        assertFalse(adapter.isDate());
    }

    @Test
    void testIsDateTime_returnsFalse() {
        assertFalse(adapter.isDateTime());
    }

    @Test
    void testIsList_returnsFalse() {
        assertFalse(adapter.isList());
    }

    @Test
    void testIsStruct_returnsFalse() {
        assertFalse(adapter.isStruct());
    }

    // ========== Getter Methods (get*) - All Should Throw ==========

    @Test
    void testGet_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.get(String.class));
    }

    @Test
    void testGetBoolean_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getBoolean());
    }

    @Test
    void testGetByte_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getByte());
    }

    @Test
    void testGetShort_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getShort());
    }

    @Test
    void testGetInteger_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getInteger());
    }

    @Test
    void testGetInt_throwsException() {
        // getInt() calls getInteger(), so it should also throw
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getInt());
    }

    @Test
    void testGetLong_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getLong());
    }

    @Test
    void testGetBigInteger_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getBigInteger());
    }

    @Test
    void testGetFloat_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getFloat());
    }

    @Test
    void testGetDouble_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getDouble());
    }

    @Test
    void testGetBigDecimal_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getBigDecimal());
    }

    @Test
    void testGetString_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getString());
    }

    @Test
    void testGetDuration_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getDuration());
    }

    @Test
    void testGetTime_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getTime());
    }

    @Test
    void testGetDate_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getDate());
    }

    @Test
    void testGetDateTime_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getDateTime());
    }

    @Test
    void testGetRaw_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getRaw());
    }

    // ========== List Methods ==========

    @Test
    void testGetLength_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getLength());
    }

    @Test
    void testGetIndex_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getIndex(0));
    }

    @Test
    void testGetList_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getList());
    }

    // ========== Struct Methods ==========

    @Test
    void testGetKeys_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getKeys());
    }

    @Test
    void testHasKey_returnsFalse() {
        assertFalse(adapter.hasKey("anyKey"));
    }

    @Test
    void testGetValue_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getValue("anyKey"));
    }

    @Test
    void testGetStruct_throwsException() {
        assertThrows(PlcIncompatibleDatatypeException.class, () -> adapter.getStruct());
    }

    // ========== Metadata Tests ==========

    @Test
    void testAddMetaData_single() {
        PlcBOOL metaValue = new PlcBOOL(true);
        adapter.addMetaData("testKey", metaValue);

        assertTrue(adapter.hasMetaData("testKey"));
        assertEquals(metaValue, adapter.getMetaData("testKey"));
        assertTrue(adapter.getMetaDataNames().contains("testKey"));
    }

    @Test
    void testAddMetaData_multiple() {
        PlcBOOL metaValue1 = new PlcBOOL(true);
        PlcINT metaValue2 = new PlcINT(42);

        adapter.addMetaData("key1", metaValue1);
        adapter.addMetaData("key2", metaValue2);

        assertTrue(adapter.hasMetaData("key1"));
        assertTrue(adapter.hasMetaData("key2"));
        assertEquals(metaValue1, adapter.getMetaData("key1"));
        assertEquals(metaValue2, adapter.getMetaData("key2"));
        assertEquals(2, adapter.getMetaDataNames().size());
    }

    @Test
    void testHasMetaData_nonExistent_returnsFalse() {
        assertFalse(adapter.hasMetaData("nonExistentKey"));
    }

    @Test
    void testGetMetaData_nonExistent_returnsNull() {
        assertNull(adapter.getMetaData("nonExistentKey"));
    }

    @Test
    void testGetMetaDataNames_empty() {
        assertTrue(adapter.getMetaDataNames().isEmpty());
    }

    @Test
    void testAddMetaData_overwrite() {
        PlcBOOL metaValue1 = new PlcBOOL(true);
        PlcBOOL metaValue2 = new PlcBOOL(false);

        adapter.addMetaData("key", metaValue1);
        adapter.addMetaData("key", metaValue2);

        assertEquals(metaValue2, adapter.getMetaData("key"));
        assertEquals(1, adapter.getMetaDataNames().size());
    }

    // ========== Equals and HashCode Tests ==========

    @Test
    void testEquals_sameInstance() {
        assertEquals(adapter, adapter);
    }

    @Test
    void testEquals_nullOther() {
        assertNotEquals(null, adapter);
    }

    @Test
    void testEquals_differentClass() {
        assertNotEquals("string", adapter);
    }

    @Test
    void testEquals_sameMetadata() {
        TestPlcValueAdapter adapter1 = new TestPlcValueAdapter();
        TestPlcValueAdapter adapter2 = new TestPlcValueAdapter();

        PlcBOOL metaValue = new PlcBOOL(true);
        adapter1.addMetaData("key", metaValue);
        adapter2.addMetaData("key", metaValue);

        assertEquals(adapter1, adapter2);
    }

    @Test
    void testEquals_differentMetadata() {
        TestPlcValueAdapter adapter1 = new TestPlcValueAdapter();
        TestPlcValueAdapter adapter2 = new TestPlcValueAdapter();

        adapter1.addMetaData("key1", new PlcBOOL(true));
        adapter2.addMetaData("key2", new PlcBOOL(true));

        assertNotEquals(adapter1, adapter2);
    }

    @Test
    void testEquals_bothEmpty() {
        TestPlcValueAdapter adapter1 = new TestPlcValueAdapter();
        TestPlcValueAdapter adapter2 = new TestPlcValueAdapter();

        assertEquals(adapter1, adapter2);
    }

    @Test
    void testHashCode_sameMetadata() {
        TestPlcValueAdapter adapter1 = new TestPlcValueAdapter();
        TestPlcValueAdapter adapter2 = new TestPlcValueAdapter();

        PlcBOOL metaValue = new PlcBOOL(true);
        adapter1.addMetaData("key", metaValue);
        adapter2.addMetaData("key", metaValue);

        assertEquals(adapter1.hashCode(), adapter2.hashCode());
    }

    @Test
    void testHashCode_emptyMetadata() {
        TestPlcValueAdapter adapter1 = new TestPlcValueAdapter();
        TestPlcValueAdapter adapter2 = new TestPlcValueAdapter();

        assertEquals(adapter1.hashCode(), adapter2.hashCode());
    }
}
