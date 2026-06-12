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
package org.apache.plc4x.java.utils.testutils.utils.dom4j;

import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.dom4j.QName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class LocationAwareElementTest {

    @Test
    void testLocationAwareElementCreation() {
        Location location = new Location(10, 5);
        QName qname = new QName("test");

        LocationAwareElement element = new LocationAwareElement(qname, location);

        assertNotNull(element);
        assertEquals(location, element.getLocation());
        assertEquals("test", element.getName());
    }

    @Test
    void testLocationAwareElementWithNullLocation() {
        QName qname = new QName("test");

        LocationAwareElement element = new LocationAwareElement(qname, null);

        assertNotNull(element);
        assertNull(element.getLocation());
    }
}
