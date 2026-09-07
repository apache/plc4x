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

import org.dom4j.DocumentFactory;
import org.junit.jupiter.api.Test;
import org.xml.sax.Locator;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class LocationAwareSAXContentHandlerTest {

    @Test
    void testSetDocumentLocator() {
        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        LocationAwareSAXContentHandler handler = new LocationAwareSAXContentHandler(factory, null);

        Locator locator = mock(Locator.class);
        when(locator.getLineNumber()).thenReturn(10);
        when(locator.getColumnNumber()).thenReturn(5);

        handler.setDocumentLocator(locator);

        assertNotNull(handler.getLocator());
        assertEquals(locator, handler.getLocator());
    }

    @Test
    void testSetDocumentLocatorWithNonLocationAwareFactory() {
        DocumentFactory factory = new DocumentFactory();
        LocationAwareSAXContentHandler handler = new LocationAwareSAXContentHandler(factory, null);

        Locator locator = mock(Locator.class);

        assertDoesNotThrow(() -> handler.setDocumentLocator(locator));
        assertEquals(locator, handler.getLocator());
    }

    @Test
    void testGetLocatorBeforeSet() {
        LocationAwareDocumentFactory factory = new LocationAwareDocumentFactory();
        LocationAwareSAXContentHandler handler = new LocationAwareSAXContentHandler(factory, null);

        assertNull(handler.getLocator());
    }
}
