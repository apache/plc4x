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

import org.junit.jupiter.api.Test;

import java.io.ByteArrayInputStream;
import java.net.URI;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

class DriverTestsuiteTest {

    @Test
    void testParseTestsuite() throws Exception {
        String xml = """
            <?xml version="1.0" encoding="UTF-8"?>
            <driver-testsuite name="Test Suite">
                <protocol-name>test</protocol-name>
                <driver-name>test</driver-name>
                <byte-order>BIG_ENDIAN</byte-order>
                <testcase>
                    <name>Simple Test</name>
                    <description>A simple test</description>
                    <steps>
                        <delay>100</delay>
                    </steps>
                </testcase>
            </driver-testsuite>
            """;

        ByteArrayInputStream inputStream = new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8));
        URI uri = URI.create("test://suite.xml");

        DriverTestsuite testsuite = DriverTestsuite.parseTestsuite(uri, "", inputStream, false);

        assertNotNull(testsuite);
        assertNotNull(testsuite.getConfiguration());
        assertEquals("test", testsuite.getConfiguration().getDriverName());
        assertEquals(1, testsuite.getTestcases().size());
        assertEquals("Simple Test", testsuite.getTestcases().get(0).getName());
    }

    @Test
    void testParseTestsuiteWithSetupAndTeardown() throws Exception {
        String xml = """
            <?xml version="1.0" encoding="UTF-8"?>
            <driver-testsuite name="Test Suite">
                <protocol-name>test</protocol-name>
                <driver-name>test</driver-name>
                <byte-order>BIG_ENDIAN</byte-order>
                <setup>
                    <delay>10</delay>
                </setup>
                <testcase>
                    <name>Test 1</name>
                    <steps>
                        <delay>20</delay>
                    </steps>
                </testcase>
                <teardown>
                    <delay>10</delay>
                </teardown>
            </driver-testsuite>
            """;

        ByteArrayInputStream inputStream = new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8));
        URI uri = URI.create("test://suite.xml");

        DriverTestsuite testsuite = DriverTestsuite.parseTestsuite(uri, "", inputStream, false);

        assertNotNull(testsuite);
        assertEquals(1, testsuite.getSetupSteps().size());
        assertEquals(1, testsuite.getTeardownSteps().size());
        assertNotNull(testsuite.getSynchronizer());
    }

    @Test
    void testParseTestsuiteWithMultipleTestcases() throws Exception {
        String xml = """
            <?xml version="1.0" encoding="UTF-8"?>
            <driver-testsuite name="Test Suite">
                <protocol-name>test</protocol-name>
                <driver-name>test</driver-name>
                <byte-order>BIG_ENDIAN</byte-order>
                <testcase>
                    <name>Test 1</name>
                    <steps><delay>10</delay></steps>
                </testcase>
                <testcase>
                    <name>Test 2</name>
                    <steps><delay>20</delay></steps>
                </testcase>
                <testcase>
                    <name>Test 3</name>
                    <steps><delay>30</delay></steps>
                </testcase>
            </driver-testsuite>
            """;

        ByteArrayInputStream inputStream = new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8));
        URI uri = URI.create("test://suite.xml");

        DriverTestsuite testsuite = DriverTestsuite.parseTestsuite(uri, "", inputStream, false);

        assertNotNull(testsuite);
        assertEquals(3, testsuite.getTestcases().size());
        assertEquals("Test 1", testsuite.getTestcases().get(0).getName());
        assertEquals("Test 2", testsuite.getTestcases().get(1).getName());
        assertEquals("Test 3", testsuite.getTestcases().get(2).getName());
    }

    @Test
    void testGetConfiguration() throws Exception {
        String xml = """
            <?xml version="1.0" encoding="UTF-8"?>
            <driver-testsuite name="Test Suite">
                <protocol-name>protocol-test</protocol-name>
                <driver-name>driver-test</driver-name>
                <byte-order>LITTLE_ENDIAN</byte-order>
                <testcase>
                    <name>Test</name>
                    <steps><delay>1</delay></steps>
                </testcase>
            </driver-testsuite>
            """;

        ByteArrayInputStream inputStream = new ByteArrayInputStream(xml.getBytes(StandardCharsets.UTF_8));
        DriverTestsuite testsuite = DriverTestsuite.parseTestsuite(URI.create("test://test.xml"), "", inputStream, false);

        DriverTestsuiteConfiguration config = testsuite.getConfiguration();
        assertNotNull(config);
        assertEquals("driver-test", config.getDriverName());
    }
}
