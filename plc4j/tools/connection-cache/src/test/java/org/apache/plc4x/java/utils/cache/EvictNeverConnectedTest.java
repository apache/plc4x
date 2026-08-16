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
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

/**
 * Evicting a connection that never connected. Up to 0.13.1 ConnectionContainer.close() ended in an
 * unguarded connection.close() on the still-null connection, so removeCachedConnection() threw an
 * NPE - and because the map entry was only dropped afterwards, the dead container was never evicted
 * at all (apache/plc4x#2418).
 */
class EvictNeverConnectedTest {

    private static final String URL = "test:tcp://localhost";

    @Mock
    private PlcConnectionManager mockConnectionManager;

    private CachedPlcConnectionManager manager;
    private ScheduledExecutorService scheduler;
    private AutoCloseable mocks;

    @BeforeEach
    void setUp() {
        mocks = MockitoAnnotations.openMocks(this);
        scheduler = Executors.newScheduledThreadPool(2);
    }

    @AfterEach
    void tearDown() throws Exception {
        if (manager != null) {
            manager.close();
        }
        scheduler.shutdownNow();
        mocks.close();
    }

    @Test
    void evictsAContainerWhoseConnectNeverSucceeded() throws Exception {
        when(mockConnectionManager.getConnection(anyString()))
            .thenThrow(new PlcConnectionException("device is unreachable"));

        manager = CachedPlcConnectionManager.getBuilder()
            .withConnectionManager(mockConnectionManager)
            .withScheduler(scheduler)
            .build();

        // The device is down, so this fails and leaves a container that never held a connection.
        assertThrows(PlcConnectionException.class, () -> manager.getConnection(URL));

        // This is what the reporter does from the error handler of the read/write loop.
        assertDoesNotThrow(() -> manager.removeCachedConnection(URL));
        assertEquals(0, manager.getCachedConnectionCount(), "no dead container is left behind");
        assertEquals(0, manager.getActiveLeaseCount());
    }

    /**
     * The same thing along the path that actually caches a container first: connect once, then have
     * the device disappear so the connection is dropped, then evict.
     */
    @Test
    void evictsAfterTheConnectionWasAlreadyDropped() throws Exception {
        PlcConnection connection = org.mockito.Mockito.mock(PlcConnection.class);
        when(connection.isConnected()).thenReturn(true);
        when(mockConnectionManager.getConnection(anyString())).thenReturn(connection);

        manager = CachedPlcConnectionManager.getBuilder()
            .withConnectionManager(mockConnectionManager)
            .withScheduler(scheduler)
            .build();

        manager.getConnection(URL).close();
        assertEquals(1, manager.getCachedConnectionCount());

        assertDoesNotThrow(() -> manager.removeCachedConnection(URL));
        // A second eviction of the same, now absent, connection must not blow up either.
        assertDoesNotThrow(() -> manager.removeCachedConnection(URL));
        assertEquals(0, manager.getCachedConnectionCount());
    }
}
