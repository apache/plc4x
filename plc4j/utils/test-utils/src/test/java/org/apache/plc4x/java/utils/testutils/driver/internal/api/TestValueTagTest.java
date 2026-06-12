/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.testutils.driver.internal.api;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TestValueTagTest {

    @Test
    void testConstructorAndGetters() {
        String[] values = {"value1", "value2"};
        TestValueTag tag = new TestValueTag("myTag", "myAddress", values);

        assertEquals("myTag", tag.getName());
        assertEquals("myAddress", tag.getAddress());
        assertArrayEquals(new String[]{"value1", "value2"}, tag.getValues());
    }

    @Test
    void testEmptyValues() {
        String[] values = {};
        TestValueTag tag = new TestValueTag("myTag", "myAddress", values);

        assertEquals(0, tag.getValues().length);
    }

    @Test
    void testSingleValue() {
        String[] values = {"42"};
        TestValueTag tag = new TestValueTag("myTag", "myAddress", values);

        assertEquals(1, tag.getValues().length);
        assertEquals("42", tag.getValues()[0]);
    }

    @Test
    void testInheritance() {
        TestValueTag tag = new TestValueTag("myTag", "myAddress", new String[]{"val"});

        // Should inherit from TestTag
        assertTrue(tag instanceof TestTag);
        assertEquals(0, tag.getLengthInBytes());
        assertEquals(0, tag.getLengthInBits());
    }
}
