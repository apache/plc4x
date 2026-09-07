package org.apache.plc4x.java.utils.testutils.driver;

import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderBigEndian;
import org.apache.plc4x.java.utils.testutils.driver.internal.DriverTestsuiteConfiguration;
import org.junit.jupiter.api.Test;

import java.net.URI;
import java.util.HashMap;
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
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class DriverTestsuiteConfigurationTest {

    @Test
    void testConfigurationCreation() throws Exception {
        URI uri = new URI("file:///test/suite.xml");
        Map<String, String> options = new HashMap<>();
        options.put("option1", "value1");

        Map<String, String> params = new HashMap<>();
        params.put("param1", "value1");

        DriverTestsuiteConfiguration config = new DriverTestsuiteConfiguration(
            uri, "Test Suite", "s7", "read-write", "s7",
            options, params, false, false, ByteOrderBigEndian.NAME
        );

        assertEquals(uri, config.getSuiteUri());
        assertEquals("Test Suite", config.getTestsuiteName());
        assertEquals("s7", config.getProtocolName());
        assertEquals("read-write", config.getOutputFlavor());
        assertEquals("s7", config.getDriverName());
        assertEquals(ByteOrderBigEndian.NAME, config.getByteOrder());
        assertFalse(config.isAutoMigrate());

        // Options should include the merged values
        assertTrue(config.getOptions().containsKey("protocolName"));
        assertEquals("s7", config.getOptions().get("protocolName"));
    }
}
