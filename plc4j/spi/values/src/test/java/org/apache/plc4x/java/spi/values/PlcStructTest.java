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

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;

import java.util.HashMap;
import java.util.Map;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcStruct - Structured PLC value (map/dictionary)
 */
class PlcStructTest {

    @Test
    void testConstructor() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("field1", new PlcINT(10));
        map.put("field2", new PlcBOOL(true));

        PlcStruct struct = new PlcStruct(map);
        assertNotNull(struct);
        assertTrue(struct.isStruct());
    }

    @Test
    void testGetPlcValueType() {
        Map<String, PlcValue> map = new HashMap<>();
        PlcStruct struct = new PlcStruct(map);
        assertEquals(PlcValueType.Struct, struct.getPlcValueType());
    }

    @Test
    void testGetObject() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("test", new PlcINT(1));
        PlcStruct struct = new PlcStruct(map);

        Object obj = struct.getObject();
        assertNotNull(obj);
        assertTrue(obj instanceof Map);
    }

    @Test
    void testGetLength() {
        Map<String, PlcValue> map = new HashMap<>();
        PlcStruct struct = new PlcStruct(map);
        // PlcStruct always returns length of 1
        assertEquals(1, struct.getLength());
    }

    @Test
    void testIsStruct() {
        Map<String, PlcValue> map = new HashMap<>();
        PlcStruct struct = new PlcStruct(map);
        assertTrue(struct.isStruct());
    }

    @Test
    void testGetKeys() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("key1", new PlcINT(1));
        map.put("key2", new PlcINT(2));
        map.put("key3", new PlcINT(3));

        PlcStruct struct = new PlcStruct(map);
        Set<String> keys = struct.getKeys();

        assertNotNull(keys);
        assertEquals(3, keys.size());
        assertTrue(keys.contains("key1"));
        assertTrue(keys.contains("key2"));
        assertTrue(keys.contains("key3"));
    }

    @Test
    void testHasKey() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("existingKey", new PlcINT(100));

        PlcStruct struct = new PlcStruct(map);
        assertTrue(struct.hasKey("existingKey"));
        assertFalse(struct.hasKey("nonExistentKey"));
    }

    @Test
    void testGetValue() {
        PlcINT value1 = new PlcINT(42);
        PlcBOOL value2 = new PlcBOOL(true);

        Map<String, PlcValue> map = new HashMap<>();
        map.put("number", value1);
        map.put("flag", value2);

        PlcStruct struct = new PlcStruct(map);

        assertEquals(value1, struct.getValue("number"));
        assertEquals(value2, struct.getValue("flag"));
        assertNull(struct.getValue("nonExistent"));
    }

    @Test
    void testGetStruct() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("a", new PlcINT(1));
        map.put("b", new PlcINT(2));

        PlcStruct struct = new PlcStruct(map);
        Map<String, ? extends PlcValue> result = struct.getStruct();

        assertNotNull(result);
        assertEquals(2, result.size());
        assertEquals(map.get("a"), result.get("a"));
        assertEquals(map.get("b"), result.get("b"));
    }

    @Test
    void testToString() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("field1", new PlcINT(10));

        PlcStruct struct = new PlcStruct(map);
        String str = struct.toString();

        assertNotNull(str);
        assertTrue(str.startsWith("{"));
        assertTrue(str.endsWith("}"));
        assertTrue(str.contains("field1"));
    }

    @Test
    void testToStringEmpty() {
        PlcStruct struct = new PlcStruct(new HashMap<>());
        String str = struct.toString();
        assertEquals("{}", str);
    }

    @Test
    void testEquals() {
        Map<String, PlcValue> map1 = new HashMap<>();
        map1.put("a", new PlcINT(1));
        map1.put("b", new PlcINT(2));

        Map<String, PlcValue> map2 = new HashMap<>();
        map2.put("a", new PlcINT(1));
        map2.put("b", new PlcINT(2));

        Map<String, PlcValue> map3 = new HashMap<>();
        map3.put("a", new PlcINT(1));
        map3.put("b", new PlcINT(3));

        PlcStruct struct1 = new PlcStruct(map1);
        PlcStruct struct2 = new PlcStruct(map2);
        PlcStruct struct3 = new PlcStruct(map3);

        // Same instance
        assertEquals(struct1, struct1);

        // Different instances, same values
        assertEquals(struct1, struct2);
        assertEquals(struct2, struct1);

        // Different values
        assertNotEquals(struct1, struct3);

        // Null
        assertNotEquals(struct1, null);

        // Different type
        assertNotEquals(struct1, new PlcBOOL(true));
    }

    @Test
    void testHashCode() {
        Map<String, PlcValue> map1 = new HashMap<>();
        map1.put("a", new PlcINT(1));
        map1.put("b", new PlcINT(2));

        Map<String, PlcValue> map2 = new HashMap<>();
        map2.put("a", new PlcINT(1));
        map2.put("b", new PlcINT(2));

        PlcStruct struct1 = new PlcStruct(map1);
        PlcStruct struct2 = new PlcStruct(map2);

        // Same values should have same hash code
        assertEquals(struct1.hashCode(), struct2.hashCode());
    }

    @Test
    void testMixedTypes() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("integer", new PlcINT(10));
        map.put("boolean", new PlcBOOL(true));
        map.put("string", new PlcSTRING("test"));

        PlcStruct struct = new PlcStruct(map);

        assertEquals(3, struct.getKeys().size());
        assertTrue(struct.getValue("integer") instanceof PlcINT);
        assertTrue(struct.getValue("boolean") instanceof PlcBOOL);
        assertTrue(struct.getValue("string") instanceof PlcSTRING);
    }

    @Test
    void testNestedStructs() {
        Map<String, PlcValue> innerMap = new HashMap<>();
        innerMap.put("x", new PlcINT(1));
        innerMap.put("y", new PlcINT(2));
        PlcStruct innerStruct = new PlcStruct(innerMap);

        Map<String, PlcValue> outerMap = new HashMap<>();
        outerMap.put("inner", innerStruct);
        outerMap.put("value", new PlcINT(3));

        PlcStruct outerStruct = new PlcStruct(outerMap);

        assertEquals(2, outerStruct.getKeys().size());
        assertTrue(outerStruct.getValue("inner") instanceof PlcStruct);

        PlcStruct retrieved = (PlcStruct) outerStruct.getValue("inner");
        assertEquals(2, retrieved.getKeys().size());
        assertTrue(retrieved.hasKey("x"));
        assertTrue(retrieved.hasKey("y"));
    }

    @Test
    void testMapIsWrappedAsUnmodifiable() {
        Map<String, PlcValue> map = new HashMap<>();
        map.put("original", new PlcINT(1));

        PlcStruct struct = new PlcStruct(map);
        Map<String, ? extends PlcValue> structMap = struct.getStruct();

        // The struct map should be unmodifiable
        assertThrows(UnsupportedOperationException.class, () -> {
            ((Map<String, PlcValue>) structMap).put("new", new PlcINT(2));
        });
    }
}
