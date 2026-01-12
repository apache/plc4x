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

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Test class for PlcRawByteArray - Raw byte array representation
 */
class PlcRawByteArrayTest {

    @Test
    void testConstructor() {
        byte[] data = {1, 2, 3, 4, 5};
        PlcRawByteArray value = new PlcRawByteArray(data);
        assertArrayEquals(data, value.getRaw());
        assertFalse(value.isNullable());
    }

    @Test
    void testOfWithPlcRawByteArray() {
        byte[] data = {1, 2, 3};
        PlcRawByteArray original = new PlcRawByteArray(data);
        PlcRawByteArray result = PlcRawByteArray.of(original);
        assertSame(original, result);
    }

    @Test
    void testOfWithByteArray() {
        byte[] data = {10, 20, 30};
        PlcRawByteArray result = PlcRawByteArray.of(data);
        assertArrayEquals(data, result.getRaw());
    }

    @Test
    void testOfWithInvalidType() {
        assertThrows(IllegalArgumentException.class, () -> PlcRawByteArray.of("string"));
        assertThrows(IllegalArgumentException.class, () -> PlcRawByteArray.of(123));
    }

    @Test
    void testGetPlcValueType() {
        PlcRawByteArray value = new PlcRawByteArray(new byte[]{1, 2});
        assertEquals(PlcValueType.RAW_BYTE_ARRAY, value.getPlcValueType());
    }

    @Test
    void testGetRaw() {
        byte[] data = {-128, 0, 127};
        PlcRawByteArray value = new PlcRawByteArray(data);
        assertArrayEquals(data, value.getRaw());
    }

    @Test
    void testToString() {
        byte[] data = {(byte) 0xAB, (byte) 0xCD, (byte) 0xEF};
        PlcRawByteArray value = new PlcRawByteArray(data);
        String result = value.toString();
        assertNotNull(result);
        // Should be a hex string representation
        assertTrue(result.length() > 0);
    }

    @Test
    void testEquals() {
        byte[] data1 = {1, 2, 3};
        byte[] data2 = {1, 2, 3};
        byte[] data3 = {1, 2, 4};

        PlcRawByteArray value1 = new PlcRawByteArray(data1);
        PlcRawByteArray value2 = new PlcRawByteArray(data2);
        PlcRawByteArray value3 = new PlcRawByteArray(data3);

        // Same instance
        assertEquals(value1, value1);

        // Different instances, same data
        assertEquals(value1, value2);
        assertEquals(value2, value1);

        // Different data
        assertNotEquals(value1, value3);

        // Null
        assertNotEquals(value1, null);

        // Different type
        assertNotEquals(value1, new PlcBOOL(true));
    }

    @Test
    void testHashCode() {
        byte[] data1 = {1, 2, 3};
        byte[] data2 = {1, 2, 3};
        byte[] data3 = {1, 2, 4};

        PlcRawByteArray value1 = new PlcRawByteArray(data1);
        PlcRawByteArray value2 = new PlcRawByteArray(data2);
        PlcRawByteArray value3 = new PlcRawByteArray(data3);

        // Same data should have same hash code
        assertEquals(value1.hashCode(), value2.hashCode());

        // Different data should (likely) have different hash code
        assertNotEquals(value1.hashCode(), value3.hashCode());
    }

    @Test
    void testIsList() {
        PlcRawByteArray value = new PlcRawByteArray(new byte[]{1, 2, 3});
        assertTrue(value.isList());
    }

    @Test
    void testGetList() {
        byte[] data = {1, 2, -1};
        PlcRawByteArray value = new PlcRawByteArray(data);
        List<PlcValue> list = value.getList();

        assertNotNull(list);
        assertEquals(3, list.size());

        // Each byte should be wrapped as PlcSINT
        assertTrue(list.get(0) instanceof PlcSINT);
        assertTrue(list.get(1) instanceof PlcSINT);
        assertTrue(list.get(2) instanceof PlcSINT);

        assertEquals((short) 1, ((PlcSINT) list.get(0)).getShort());
        assertEquals((short) 2, ((PlcSINT) list.get(1)).getShort());
        assertEquals((short) -1, ((PlcSINT) list.get(2)).getShort());
    }

    @Test
    void testGetListWithEmptyArray() {
        PlcRawByteArray value = new PlcRawByteArray(new byte[0]);
        List<PlcValue> list = value.getList();
        assertNotNull(list);
        assertTrue(list.isEmpty());
    }

    @Test
    void testGetListWithLargeArray() {
        byte[] data = new byte[100];
        for (int i = 0; i < 100; i++) {
            data[i] = (byte) i;
        }
        PlcRawByteArray value = new PlcRawByteArray(data);
        List<PlcValue> list = value.getList();
        assertEquals(100, list.size());
    }
}
