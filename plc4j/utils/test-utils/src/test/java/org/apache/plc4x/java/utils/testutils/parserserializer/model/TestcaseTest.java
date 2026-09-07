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
package org.apache.plc4x.java.utils.testutils.parserserializer.model;

import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class TestcaseTest {

    @Test
    void testTestcaseCreation() {
        Element xml = DocumentHelper.createElement("test");
        List<String> parserArgs = Arrays.asList("arg1", "arg2");
        byte[] raw = new byte[]{0x01, 0x02, 0x03};

        Testcase testcase = new Testcase(
            "SuiteName",
            "ProtocolName",
            "OutputFlavor",
            "TestCase1",
            "Description",
            raw,
            "RootType",
            parserArgs,
            xml
        );

        assertNotNull(testcase);
        assertEquals("SuiteName", testcase.getTestSuiteName());
        assertEquals("ProtocolName", testcase.getProtocolName());
        assertEquals("OutputFlavor", testcase.getOutputFlavor());
        assertEquals("TestCase1", testcase.getName());
        assertEquals("Description", testcase.getDescription());
        assertArrayEquals(raw, testcase.getRaw());
        assertEquals("RootType", testcase.getRootType());
        assertEquals(parserArgs, testcase.getParserArguments());
        assertEquals(xml, testcase.getXml());
    }

    @Test
    void testSetAndGetLocation() {
        Element xml = DocumentHelper.createElement("test");
        List<String> parserArgs = new ArrayList<>();
        byte[] raw = new byte[]{0x01};

        Testcase testcase = new Testcase(
            "Suite",
            "Protocol",
            "Flavor",
            "Test",
            "Desc",
            raw,
            "Type",
            parserArgs,
            xml
        );

        // Initially no location
        assertTrue(testcase.getLocation().isEmpty());

        // Set location
        Location location = new Location(10, 5);
        testcase.setLocation(location);

        assertEquals(location, testcase.getLocation().orElse(null));
    }

    @Test
    void testToString() {
        Element xml = DocumentHelper.createElement("test");
        List<String> parserArgs = new ArrayList<>();
        byte[] raw = new byte[]{};

        Testcase testcase = new Testcase(
            "Suite",
            "Protocol",
            "Flavor",
            "Test",
            "Desc",
            raw,
            "Type",
            parserArgs,
            xml
        );

        String result = testcase.toString();
        assertTrue(result.contains("Suite"));
        assertTrue(result.contains("Protocol"));
        assertTrue(result.contains("Flavor"));
        assertTrue(result.contains("Test"));
    }
}
