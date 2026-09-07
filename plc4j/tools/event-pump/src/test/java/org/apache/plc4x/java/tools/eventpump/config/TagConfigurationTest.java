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

package org.apache.plc4x.java.tools.eventpump.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TagConfigurationTest {

    @Test
    void testDefaultConstructor() {
        TagConfiguration config = new TagConfiguration();

        assertNull(config.getAddress());
        assertNull(config.getTransform());
        assertFalse(config.hasTransform());
    }

    @Test
    void testConstructorWithAddress() {
        TagConfiguration config = new TagConfiguration("MAIN.temperature");

        assertEquals("MAIN.temperature", config.getAddress());
        assertNull(config.getTransform());
        assertFalse(config.hasTransform());
    }

    @Test
    void testConstructorWithAddressAndTransform() {
        TagConfiguration config = new TagConfiguration("MAIN.temperature", "value * 1.8 + 32");

        assertEquals("MAIN.temperature", config.getAddress());
        assertEquals("value * 1.8 + 32", config.getTransform());
        assertTrue(config.hasTransform());
    }

    @Test
    void testConstructorWithAddressAndNullTransform() {
        TagConfiguration config = new TagConfiguration("MAIN.temperature", null);

        assertEquals("MAIN.temperature", config.getAddress());
        assertNull(config.getTransform());
        assertFalse(config.hasTransform());
    }

    @Test
    void testSetAndGetAddress() {
        TagConfiguration config = new TagConfiguration();

        config.setAddress("MAIN.pressure");

        assertEquals("MAIN.pressure", config.getAddress());
    }

    @Test
    void testSetAndGetTransform() {
        TagConfiguration config = new TagConfiguration();

        config.setTransform("value + 100");

        assertEquals("value + 100", config.getTransform());
        assertTrue(config.hasTransform());
    }

    @Test
    void testHasTransformWithEmptyString() {
        TagConfiguration config = new TagConfiguration("MAIN.temp", "");

        assertFalse(config.hasTransform());
    }

    @Test
    void testHasTransformWithWhitespace() {
        TagConfiguration config = new TagConfiguration("MAIN.temp", "   ");

        assertFalse(config.hasTransform());
    }

    @Test
    void testToStringWithoutTransform() {
        TagConfiguration config = new TagConfiguration("MAIN.temperature");

        String result = config.toString();

        assertEquals("TagConfiguration{address='MAIN.temperature'}", result);
    }

    @Test
    void testToStringWithTransform() {
        TagConfiguration config = new TagConfiguration("MAIN.temperature", "value * 1.8 + 32");

        String result = config.toString();

        assertEquals("TagConfiguration{address='MAIN.temperature', transform='value * 1.8 + 32'}", result);
    }

    @Test
    void testEquals_Same() {
        TagConfiguration config = new TagConfiguration("MAIN.temp", "value * 2");

        assertTrue(config.equals(config));
    }

    @Test
    void testEquals_Null() {
        TagConfiguration config = new TagConfiguration("MAIN.temp");

        assertFalse(config.equals(null));
    }

    @Test
    void testEquals_DifferentClass() {
        TagConfiguration config = new TagConfiguration("MAIN.temp");

        assertFalse(config.equals("not a TagConfiguration"));
    }

    @Test
    void testEquals_SameAddressNoTransform() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp");

        assertTrue(config1.equals(config2));
        assertTrue(config2.equals(config1));
    }

    @Test
    void testEquals_SameAddressAndTransform() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp", "value * 2");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp", "value * 2");

        assertTrue(config1.equals(config2));
        assertTrue(config2.equals(config1));
    }

    @Test
    void testEquals_DifferentAddress() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp1");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp2");

        assertFalse(config1.equals(config2));
        assertFalse(config2.equals(config1));
    }

    @Test
    void testEquals_DifferentTransform() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp", "value * 2");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp", "value * 3");

        assertFalse(config1.equals(config2));
        assertFalse(config2.equals(config1));
    }

    @Test
    void testEquals_OneHasTransformOtherDoesNot() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp", "value * 2");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp");

        assertFalse(config1.equals(config2));
        assertFalse(config2.equals(config1));
    }

    @Test
    void testEquals_BothAddressesNull() {
        TagConfiguration config1 = new TagConfiguration();
        TagConfiguration config2 = new TagConfiguration();

        assertTrue(config1.equals(config2));
    }

    @Test
    void testEquals_OneAddressNull() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp");
        TagConfiguration config2 = new TagConfiguration();

        assertFalse(config1.equals(config2));
    }

    @Test
    void testHashCode_Equal() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp", "value * 2");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp", "value * 2");

        assertEquals(config1.hashCode(), config2.hashCode());
    }

    @Test
    void testHashCode_DifferentAddress() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp1");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp2");

        assertNotEquals(config1.hashCode(), config2.hashCode());
    }

    @Test
    void testHashCode_DifferentTransform() {
        TagConfiguration config1 = new TagConfiguration("MAIN.temp", "value * 2");
        TagConfiguration config2 = new TagConfiguration("MAIN.temp", "value * 3");

        assertNotEquals(config1.hashCode(), config2.hashCode());
    }

    @Test
    void testHashCode_NullFields() {
        TagConfiguration config = new TagConfiguration();

        // Should not throw exception
        int hashCode = config.hashCode();

        assertTrue(hashCode == 0); // 31 * 0 + 0 = 0
    }
}
