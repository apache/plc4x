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

import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcList - List of PLC values
 */
class PlcListTest {

    @Test
    void testEmptyConstructor() {
        PlcList list = new PlcList();
        assertNotNull(list.getList());
        assertEquals(0, list.getLength());
        assertTrue(list.isList());
    }

    @Test
    void testConstructorWithList() {
        List<PlcValue> values = Arrays.asList(
            new PlcINT(10),
            new PlcINT(20),
            new PlcINT(30)
        );
        PlcList list = new PlcList(values);
        assertEquals(3, list.getLength());
        assertEquals(values.size(), list.getList().size());
    }

    @Test
    void testAdd() {
        PlcList list = new PlcList();
        PlcINT value1 = new PlcINT(1);
        PlcINT value2 = new PlcINT(2);

        list.add(value1);
        assertEquals(1, list.getLength());

        list.add(value2);
        assertEquals(2, list.getLength());
    }

    @Test
    void testGetPlcValueType() {
        PlcList list = new PlcList();
        assertEquals(PlcValueType.List, list.getPlcValueType());
    }

    @Test
    void testGetObject() {
        List<PlcValue> values = Arrays.asList(new PlcINT(1), new PlcINT(2));
        PlcList list = new PlcList(values);
        Object obj = list.getObject();
        assertNotNull(obj);
        assertTrue(obj instanceof List);
    }

    @Test
    void testIsList() {
        PlcList list = new PlcList();
        assertTrue(list.isList());
    }

    @Test
    void testGetLength() {
        PlcList list = new PlcList();
        assertEquals(0, list.getLength());

        list.add(new PlcINT(1));
        assertEquals(1, list.getLength());

        list.add(new PlcINT(2));
        assertEquals(2, list.getLength());
    }

    @Test
    void testGetIndex() {
        PlcINT value1 = new PlcINT(10);
        PlcINT value2 = new PlcINT(20);
        PlcList list = new PlcList(Arrays.asList(value1, value2));

        assertEquals(value1, list.getIndex(0));
        assertEquals(value2, list.getIndex(1));
    }

    @Test
    void testGetList() {
        List<PlcValue> values = Arrays.asList(
            new PlcINT(1),
            new PlcINT(2),
            new PlcINT(3)
        );
        PlcList list = new PlcList(values);
        List<PlcValue> result = list.getList();

        assertNotNull(result);
        assertEquals(3, result.size());
    }

    @Test
    void testToString() {
        PlcList list = new PlcList(Arrays.asList(
            new PlcINT(1),
            new PlcINT(2)
        ));
        String str = list.toString();
        assertNotNull(str);
        assertTrue(str.startsWith("["));
        assertTrue(str.endsWith("]"));
        assertTrue(str.contains(","));
    }

    @Test
    void testToStringEmpty() {
        PlcList list = new PlcList();
        String str = list.toString();
        assertEquals("[]", str);
    }

    @Test
    void testEquals() {
        List<PlcValue> values1 = Arrays.asList(new PlcINT(1), new PlcINT(2));
        List<PlcValue> values2 = Arrays.asList(new PlcINT(1), new PlcINT(2));
        List<PlcValue> values3 = Arrays.asList(new PlcINT(1), new PlcINT(3));

        PlcList list1 = new PlcList(values1);
        PlcList list2 = new PlcList(values2);
        PlcList list3 = new PlcList(values3);

        // Same instance
        assertEquals(list1, list1);

        // Different instances, same values
        assertEquals(list1, list2);
        assertEquals(list2, list1);

        // Different values
        assertNotEquals(list1, list3);

        // Null
        assertNotEquals(list1, null);

        // Different type
        assertNotEquals(list1, new PlcBOOL(true));
    }

    @Test
    void testHashCode() {
        List<PlcValue> values1 = Arrays.asList(new PlcINT(1), new PlcINT(2));
        List<PlcValue> values2 = Arrays.asList(new PlcINT(1), new PlcINT(2));

        PlcList list1 = new PlcList(values1);
        PlcList list2 = new PlcList(values2);

        // Same values should have same hash code
        assertEquals(list1.hashCode(), list2.hashCode());
    }

    @Test
    void testMixedTypes() {
        PlcList list = new PlcList();
        list.add(new PlcINT(10));
        list.add(new PlcBOOL(true));
        list.add(new PlcSTRING("test"));

        assertEquals(3, list.getLength());
        assertTrue(list.getIndex(0) instanceof PlcINT);
        assertTrue(list.getIndex(1) instanceof PlcBOOL);
        assertTrue(list.getIndex(2) instanceof PlcSTRING);
    }

    @Test
    void testNestedLists() {
        PlcList innerList = new PlcList(Arrays.asList(new PlcINT(1), new PlcINT(2)));
        PlcList outerList = new PlcList();
        outerList.add(innerList);
        outerList.add(new PlcINT(3));

        assertEquals(2, outerList.getLength());
        assertTrue(outerList.getIndex(0) instanceof PlcList);
        PlcList retrieved = (PlcList) outerList.getIndex(0);
        assertEquals(2, retrieved.getLength());
    }
}
