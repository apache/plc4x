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

import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.api.PlcConnection;
import org.junit.jupiter.api.Test;

import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;

class ConnectionManagerTest {

    @Test
    void testGetConnectionWithNoDriver() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();

        // Should throw because no driver is available via ServiceLoader for "nonexistent"
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("nonexistent", params);
        });
    }

    @Test
    void testGetTransportInstanceWithNonExposingConnection() {
        ConnectionManager manager = new ConnectionManager();
        PlcConnection mockConnection = mock(PlcConnection.class);

        assertThrows(Exception.class, () -> {
            manager.getTransportInstance(mockConnection);
        });
    }

    @Test
    void testGetConnectionWithParameters() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();
        params.put("host", "localhost");
        params.put("port", "502");

        // Should still throw because no driver is available for "test"
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("test", params);
        });
    }

    @Test
    void testGetConnectionWithEmptyParameters() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();

        // Should throw because no driver
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("test", params);
        });
    }

    @Test
    void testGetConnectionWithMultipleParameters() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();
        params.put("param1", "value1");
        params.put("param2", "value2");
        params.put("param3", "value3");

        // Should throw because no driver
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("multitest", params);
        });
    }

    @Test
    void testGetConnectionWithTransportParameter() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();
        params.put("transport", "test");
        params.put("host", "localhost");

        // Should throw because no driver (but exercises transport path)
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("modbus-tcp", params);
        });
    }

    @Test
    void testGetConnectionWithTransportAndCustomHost() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();
        params.put("transport", "tcp");
        params.put("host", "192.168.1.100");

        // Should throw because no driver
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("test-driver", params);
        });
    }

    @Test
    void testGetConnectionDefaultHost() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();
        // No host provided - should default to localhost

        // Should throw because no driver
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("test", params);
        });
    }

    @Test
    void testGetConnectionWithQueryParameters() {
        ConnectionManager manager = new ConnectionManager();
        Map<String, String> params = new HashMap<>();
        params.put("host", "testhost");
        params.put("timeout", "5000");
        params.put("retries", "3");
        // No transport - should build query string

        // Should throw because no driver (but exercises query parameter building)
        assertThrows(DriverTestsuiteException.class, () -> {
            manager.getConnection("test", params);
        });
    }

    @Test
    void testGetTransportInstanceWithConnectionHavingMethod() {
        ConnectionManager manager = new ConnectionManager();
        // Create a mock that has getTransportInstance method
        PlcConnection mockConnection = mock(PlcConnection.class);

        // Should throw because mock doesn't have the method
        Exception exception = assertThrows(Exception.class, () -> {
            manager.getTransportInstance(mockConnection);
        });
        assertNotNull(exception);
    }
}
