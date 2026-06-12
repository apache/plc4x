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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.apache.plc4x.java.utils.testutils.driver.internal.utils.Synchronizer;
import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class TestcaseTest {

    @Test
    void testTestcaseCreation() {
        List<TestStep> steps = new ArrayList<>();
        Synchronizer sync = new Synchronizer();
        Location location = new Location(10, 5);

        Testcase testcase = new Testcase("Test1", "Description", steps, sync, location);

        assertNotNull(testcase);
        assertEquals("Test1", testcase.getName());
        assertEquals("Description", testcase.getDescription());
        assertEquals(steps, testcase.getSteps());
        assertEquals(location, testcase.getLocation().orElse(null));
    }

    @Test
    void testTestcaseWithNullLocation() {
        List<TestStep> steps = new ArrayList<>();
        Synchronizer sync = new Synchronizer();

        Testcase testcase = new Testcase("Test1", "Description", steps, sync, null);

        assertTrue(testcase.getLocation().isEmpty());
    }

    @Test
    void testSetTestsuite() {
        List<TestStep> steps = new ArrayList<>();
        Synchronizer sync = new Synchronizer();
        Location location = new Location(10, 5);

        Testcase testcase = new Testcase("Test1", "Description", steps, sync, location);

        // Initially null
        assertDoesNotThrow(() -> testcase.setTestsuite(null));
    }

    @Test
    void testGetters() {
        List<TestStep> steps = new ArrayList<>();
        Synchronizer sync = new Synchronizer();
        Location location = new Location(15, 20);

        Testcase testcase = new Testcase("MyTest", "MyDescription", steps, sync, location);

        assertEquals("MyTest", testcase.getName());
        assertEquals("MyDescription", testcase.getDescription());
        assertSame(steps, testcase.getSteps());
    }

    @Test
    void testEmptyDescription() {
        List<TestStep> steps = new ArrayList<>();
        Synchronizer sync = new Synchronizer();

        Testcase testcase = new Testcase("Test", "", steps, sync, null);

        assertEquals("", testcase.getDescription());
    }

    @Test
    void testNullDescription() {
        List<TestStep> steps = new ArrayList<>();
        Synchronizer sync = new Synchronizer();

        Testcase testcase = new Testcase("Test", null, steps, sync, null);

        assertNull(testcase.getDescription());
    }
}
