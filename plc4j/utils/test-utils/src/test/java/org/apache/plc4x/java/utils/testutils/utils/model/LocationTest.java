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
package org.apache.plc4x.java.utils.testutils.utils.model;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class LocationTest {

    @Test
    void testLocationCreation() {
        Location location = new Location(10, 5);
        assertEquals(10, location.line());
        assertEquals(5, location.column());
    }

    @Test
    void testLocationToString() {
        Location location = new Location(10, 5);
        String result = location.toString();
        assertTrue(result.contains("line=10"));
        assertTrue(result.contains("column=5"));
    }

    @Test
    void testLocationEquality() {
        Location location1 = new Location(10, 5);
        Location location2 = new Location(10, 5);
        Location location3 = new Location(11, 5);

        assertEquals(location1, location2);
        assertNotEquals(location1, location3);
    }

    @Test
    void testLocationHashCode() {
        Location location1 = new Location(10, 5);
        Location location2 = new Location(10, 5);

        assertEquals(location1.hashCode(), location2.hashCode());
    }
}
