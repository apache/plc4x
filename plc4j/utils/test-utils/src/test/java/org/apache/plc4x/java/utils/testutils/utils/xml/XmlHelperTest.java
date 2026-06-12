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
package org.apache.plc4x.java.utils.testutils.utils.xml;

import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class XmlHelperTest {

    @Test
    void testExtractText() {
        Element root = DocumentHelper.createElement("root");
        Element child = root.addElement("child");
        child.setText("test value");

        String result = XmlHelper.extractText(root, "child");
        assertEquals("test value", result);
    }

    @Test
    void testExtractTextWithWhitespace() {
        Element root = DocumentHelper.createElement("root");
        Element child = root.addElement("child");
        child.setText("  test value  ");

        String result = XmlHelper.extractText(root, "child");
        assertEquals("test value", result);
    }

    @Test
    void testExtractTextMissingElement() {
        Element root = DocumentHelper.createElement("root");

        RuntimeException exception = assertThrows(RuntimeException.class, () -> {
            XmlHelper.extractText(root, "missing");
        });
        assertTrue(exception.getMessage().contains("Required element missing not present"));
    }

    @Test
    void testParseParametersEmpty() {
        Map<String, String> result = XmlHelper.parseParameters(null);
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void testParseParametersSingle() {
        Element params = DocumentHelper.createElement("parameters");
        Element param1 = params.addElement("parameter");
        param1.addElement("name").setText("key1");
        param1.addElement("value").setText("value1");

        Map<String, String> result = XmlHelper.parseParameters(params);
        assertEquals(1, result.size());
        assertEquals("value1", result.get("key1"));
    }

    @Test
    void testParseParametersMultiple() {
        Element params = DocumentHelper.createElement("parameters");

        Element param1 = params.addElement("parameter");
        param1.addElement("name").setText("key1");
        param1.addElement("value").setText("value1");

        Element param2 = params.addElement("parameter");
        param2.addElement("name").setText("key2");
        param2.addElement("value").setText("value2");

        Map<String, String> result = XmlHelper.parseParameters(params);
        assertEquals(2, result.size());
        assertEquals("value1", result.get("key1"));
        assertEquals("value2", result.get("key2"));
    }

    @Test
    void testParseParametersNoElements() {
        Element params = DocumentHelper.createElement("parameters");

        Map<String, String> result = XmlHelper.parseParameters(params);
        assertNotNull(result);
        assertTrue(result.isEmpty());
    }
}
