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

import org.dom4j.Document;
import org.dom4j.DocumentFactory;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

class LocationAwareSAXReaderTest {

    @Test
    void testReadXmlDocument() throws Exception {
        String xml = "<root><child>test</child></root>";
        ByteArrayInputStream inputStream = new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8));

        LocationAwareSAXReader reader = new LocationAwareSAXReader();
        Document document = reader.read(inputStream);

        assertNotNull(document);
        assertNotNull(document.getRootElement());
        assertEquals("root", document.getRootElement().getName());
    }

    @Test
    void testSetDocumentFactory() {
        LocationAwareSAXReader reader = new LocationAwareSAXReader();
        DocumentFactory factory = new LocationAwareDocumentFactory();

        assertDoesNotThrow(() -> reader.setDocumentFactory(factory));
    }
}
