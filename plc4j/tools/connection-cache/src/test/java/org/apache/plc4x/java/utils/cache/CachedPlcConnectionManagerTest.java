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
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.time.Duration;

public class CachedPlcConnectionManagerTest {

    /**
     * This is the simplest possible test. Here the ConnectionManager is used exactly once.
     * So not really much of the caching we can test, but it tests if we're creating connections the right way.
     *
     * @throws PlcConnectionException something went wrong
     */
    @Test
    public void testSingleConnectionRequestTest() throws PlcConnectionException {
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).build();

        // Get the connection for the first time.
        try (PlcConnection connection = connectionManager.getConnection("test")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
        } catch (Exception e) {
            Assertions.fail("Not expecting an exception here", e);
        }

        // Check getConnection was called on the mockConnectionManager instance exactly once.
        Mockito.verify(mockConnectionManager, Mockito.times(1)).getConnection("test");
    }

    /**
     * This test tries to get one connection two times after each other, in this case for the second time the
     * connection should not be created, but the initial one be reused.
     *
     * @throws PlcConnectionException something went wrong
     */
    @Test
    public void testDoubleConnectionRequestTest() throws PlcConnectionException {
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).build();

        // Get the connection for the first time.
        try (PlcConnection connection = connectionManager.getConnection("test")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
        } catch (Exception e) {
            Assertions.fail("Not expecting an exception here", e);
        }

        // Get the same connection a second time.
        try (PlcConnection connection = connectionManager.getConnection("test")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
        } catch (Exception e) {
            Assertions.fail("Not expecting an exception here", e);
        }

        // Check getConnection was called on the mockConnectionManager instance exactly once.
        Mockito.verify(mockConnectionManager, Mockito.times(1)).getConnection("test");
    }

    /**
     * In contrast to the previous test, in this case different connection strings are used and the
     * cache should create two different connections.
     *
     * @throws PlcConnectionException something went wrong
     */
    @Test
    public void testDoubleConnectionRequestWithDifferentConnectionsTest() throws PlcConnectionException {
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).build();

        // Get the connection for the first time.
        try (PlcConnection connection = connectionManager.getConnection("test")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
        } catch (Exception e) {
            Assertions.fail("Not expecting an exception here", e);
        }

        // Get the same connection a second time.
        try (PlcConnection connection = connectionManager.getConnection("test-other")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
        } catch (Exception e) {
            Assertions.fail("Not expecting an exception here", e);
        }

        // Check getConnection was called on the mockConnectionManager instance twice, as they are different connections.
        Mockito.verify(mockConnectionManager, Mockito.times(2)).getConnection(Mockito.any());
    }

    /**
     * This test is almost the same setup as the double connection test, but in this case the one usage exceeds
     * the maximum wait time, so the connection-cache gives up waiting and returns an exception.
     *
     * @throws PlcConnectionException something went wrong
     */
    @Test
    public void testDoubleConnectionRequestTimeoutTest() throws Exception {
        // Create a connectionManager with a maximum wait time of 50ms
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).withMaxWaitTime(Duration.ofMillis(50)).build();

        // Get the connection for the first time.
        (new Thread(() -> {
            try (PlcConnection connection = connectionManager.getConnection("test")) {
                Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
                // Sleep for a second.
                Thread.sleep(100L);
            } catch (Exception e) {
                Assertions.fail("Not expecting an exception here", e);
            }
        })).start();

        // This is needed as starting the previous thread seems to take a little-bit of time.
        Thread.sleep(10L);

        // Get the same connection a second time.
        try (PlcConnection ignored = connectionManager.getConnection("test")) {
            Assertions.fail("Was expecting an exception here");
        } catch (Exception e) {
            Assertions.assertInstanceOf(PlcConnectionException.class, e);
            Assertions.assertEquals("Error acquiring lease for connection", e.getMessage());
        }

        // Check getConnection was called on the mockConnectionManager instance exactly once.
        Mockito.verify(mockConnectionManager, Mockito.times(1)).getConnection("test");
    }

    /**
     * This is the simplest possible test. Here the ConnectionManager is used exactly once.
     * So not really much of the caching we can test, but it tests if we're creating connections the right way.
     *
     * @throws PlcConnectionException something went wrong
     */
    @Test
    public void testSingleConnectionRequestWithTimeoutTest() throws PlcConnectionException {
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).withMaxLeaseTime(Duration.ofMillis(10)).build();

        // Get the connection for the first time.
        try (PlcConnection connection = connectionManager.getConnection("test")) {
            Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
            Thread.sleep(100L);
        } catch (Exception e) {
            Assertions.assertInstanceOf(PlcRuntimeException.class, e);
            Assertions.assertEquals("Error trying to return lease from invalid connection", e.getMessage());
        }

        // Check getConnection was called on the mockConnectionManager instance exactly once.
        Mockito.verify(mockConnectionManager, Mockito.times(1)).getConnection("test");
    }

}
