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
import org.apache.plc4x.java.utils.cache.exceptions.PlcConnectionManagerClosedException;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;

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
        Mockito.when(mockConnectionManager.getConnection("test")).thenReturn(Mockito.mock(PlcConnection.class));
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
        Mockito.when(mockConnectionManager.getConnection("test")).thenReturn(Mockito.mock(PlcConnection.class));
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).withMaxWaitTime(Duration.ofMillis(50)).build();
        CountDownLatch startSignal = new CountDownLatch(1);

        // Get the connection for the first time.
        (new Thread(() -> {
            try {
                PlcConnection connection = connectionManager.getConnection("test");
                startSignal.countDown();
                Assertions.assertInstanceOf(LeasedPlcConnection.class, connection);
            } catch (Exception e) {
                Assertions.fail("Not expecting an exception here", e);
            }
        })).start();

        // This is needed as starting the previous thread seems to take a little-bit of time.
        startSignal.await();

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

    @Test
    @Disabled("This test fails quite regularly when run on github actions")
    public void testClosingConnectionCache() throws Exception {
        PlcConnection mockConnection = Mockito.mock(PlcConnection.class);
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        Mockito.when(mockConnectionManager.getConnection("test")).thenReturn(mockConnection);

        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).withMaxLeaseTime(Duration.ofMillis(3000)).build();

        // Have multiple leases borrowed.
        // The first should get the lease directly but will hang on to it for some time.
        CompletableFuture<Void> firstFuture = new CompletableFuture<>();
        Thread thread = new Thread(() -> {
            try {
                PlcConnection connection = connectionManager.getConnection("test");
                Thread.sleep(1000L);
                connection.close();
                firstFuture.completeExceptionally(new Exception("First connection should have failed"));
            } catch (InterruptedException e) {
                // The thread will be blocking longer than the test execution, so we're sending an interrupt
                // and this is the way we find out we got it.
                firstFuture.completeExceptionally(e);
            } catch (Exception e) {
                firstFuture.completeExceptionally(e);
            }
        });
        thread.start();

        // Give the thread some time to start up.
        Thread.sleep(100L);

        // The second will wait in the queue.
        CompletableFuture<Void> secondFuture = new CompletableFuture<>();
        new Thread(() -> {
            try {
                connectionManager.getConnection("test");
                secondFuture.completeExceptionally(new Exception("Second connection should have failed"));
            } catch (PlcConnectionException e) {
                // This getConnection() call was waiting in line in order to get the connection
                // when the cache was closed, so we expect this to fail and to contain a reference to the closed
                // connection manager.
                secondFuture.completeExceptionally(e);
            }
        }).start();

        // Give the thread some time to start up.
        Thread.sleep(100L);

        // Check that only one connection has been requested from the PlcConnectionManager.
        Mockito.verify(mockConnectionManager, Mockito.times(1)).getConnection("test");
        // Check that the connection borrowed from the Mocked ConnectionManager has not been closed yet.
        Mockito.verify(mockConnection, Mockito.times(0)).close();

        // Close the connection manager.
        connectionManager.close();

        // Interrupt the first thread, that is just wasting our time waiting.
        thread.interrupt();

        // Borrowing after the connectionManager is closed, should result in exceptions.
        try {
            connectionManager.getConnection("test");
            Assertions.fail("This should have failed");
        } catch (PlcConnectionException e) {
            if(!(e instanceof PlcConnectionManagerClosedException)) {
                e.printStackTrace();
                Assertions.fail("Expected PlcConnectionManagerClosedException");
            }
        }

        // Check that the connection borrowed from the Mocked ConnectionManager has been closed.
        Mockito.verify(mockConnection, Mockito.times(1)).close();

        try {
            firstFuture.get();
            Assertions.fail("This should have failed");
        }
        catch (Exception e) {
            // In the case of the first thread, the thread was stuck waiting in the one-second pause
            // when we intentionally interrupted it after closing the cache. So we expect to see that
            // interrupt exception here.
            if(!(e instanceof ExecutionException)) {
                e.printStackTrace();
                Assertions.fail("Expected ExecutionException");
            }
            if(!(e.getCause() instanceof InterruptedException)) {
                e.printStackTrace();
                Assertions.fail("Expected InterruptedException");
            }
        }

        try {
            secondFuture.get();
            Assertions.fail("This should have failed");
        } catch (Exception e) {
            // In this case the process was waiting for getting the connection thread 1 was hogging.
            // When closing the cache, all waiting connection requests were instantly finished with
            // an exception.
            if(!(e instanceof ExecutionException)) {
                e.printStackTrace();
                Assertions.fail("Expected ExecutionException");
            }
            if(!(e.getCause() instanceof PlcConnectionException)) {
                e.printStackTrace();
                Assertions.fail("Expected PlcConnectionException");
            }
            if(!(e.getCause().getCause().getCause() instanceof PlcConnectionManagerClosedException)) {
                e.printStackTrace();
                Assertions.fail("Expected PlcConnectionManagerClosedException");
            }
        }
    }

    @Test
    public void testCloseAfterIdleTime() throws Exception {
        PlcConnectionManager mockConnectionManager = Mockito.mock(PlcConnectionManager.class);
        Mockito.when(mockConnectionManager.getConnection("test")).thenReturn(Mockito.mock(PlcConnection.class));
        CachedPlcConnectionManager connectionManager = CachedPlcConnectionManager.getBuilder(mockConnectionManager).withMaxWaitTime(Duration.ofMillis(50)).withMaxIdleTime(Duration.ofMillis(10)).build();

        // Get a connection and directly return it.
        PlcConnection connection = connectionManager.getConnection("test");
        connection.close();

        // Wait for longer than the max idle time.
        Thread.sleep(200);

        Assertions.assertEquals(0, connectionManager.getCachedConnections().size());
    }

}
