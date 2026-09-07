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
import org.dom4j.Element;
import org.dom4j.QName;
import org.junit.jupiter.api.Test;
import org.xml.sax.Locator;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class LocationAwareDocumentFactoryTest {

    @Test
    void testCreateElementWithoutLocator() {
        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        QName qname = new QName("test");

        Element element = factory.createElement(qname);

        assertNotNull(element);
        assertEquals("test", element.getName());
    }

    @Test
    void testCreateElementWithLocator() {
        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        Locator locator = mock(Locator.class);
        when(locator.getLineNumber()).thenReturn(10);
        when(locator.getColumnNumber()).thenReturn(5);

        factory.setLocator(locator);

        QName qname = new QName("test");
        Element element = factory.createElement(qname);

        assertNotNull(element);
        assertInstanceOf(LocationAwareElement.class, element);
        LocationAwareElement locationAwareElement = (LocationAwareElement) element;
        Location location = locationAwareElement.getLocation();
        assertEquals(10, location.line());
        assertEquals(5, location.column());
    }

    @Test
    void testSetLocator() {
        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        Locator locator = mock(Locator.class);

        assertDoesNotThrow(() -> factory.setLocator(locator));
    }

    @Test
    void testMultipleElementsWithSameLocator() {
        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        Locator locator = mock(Locator.class);
        when(locator.getLineNumber()).thenReturn(20);
        when(locator.getColumnNumber()).thenReturn(10);

        factory.setLocator(locator);

        Element element1 = factory.createElement(new QName("elem1"));
        Element element2 = factory.createElement(new QName("elem2"));

        assertInstanceOf(LocationAwareElement.class, element1);
        assertInstanceOf(LocationAwareElement.class, element2);
    }
}
