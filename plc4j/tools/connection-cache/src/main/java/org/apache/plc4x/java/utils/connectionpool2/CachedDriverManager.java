/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.management.ObjectName;
import java.lang.management.ManagementFactory;
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicReference;

/**
 * Driver Manager who Caches ONE Single Connection.
 * <p>
 * Usage Example:
 * <code>
 * PlcDriverManager manager = new PlcDriverManager();
 * PlcDriverManager cached = new CachedDriverManager(url, () -&gt; manager.getConnection(url));
 * </code>
 * Now you can use "cached" everywhere you need the corresponding connection.
 */
public class CachedDriverManager extends PlcDriverManager implements CachedDriverManagerMBean {

    private static final Logger logger = LoggerFactory.getLogger(CachedDriverManager.class);

    // Constants
    public static final int LONG_BORROW_WATCHDOG_TIMEOUT_MS = 5_000;

    // JMX

    private final AtomicInteger numberOfConnects = new AtomicInteger(0);
    private final AtomicInteger numberOfBorrows = new AtomicInteger(0);
    private final AtomicInteger numberOfRejections = new AtomicInteger(0);
    private final AtomicInteger numberOfWatchdogs = new AtomicInteger(0);

    // End of JMX

    private final String url;
    private final PlcConnectionFactory connectionFactory;

    // Add Queue for Requests
    private final Queue<CompletableFuture<PlcConnection>> queue = new LinkedList<>();
    private final int timeoutMillis;

    private final ScheduledExecutorService executorService = Executors.newSingleThreadScheduledExecutor();

    private AtomicReference<ConnectionState> state = new AtomicReference<>(ConnectionState.DISCONNECTED);
    private PlcConnection activeConnection;
    private CachedPlcConnection borrowedConnection;
    private ScheduledFuture<?> borrowWatchdog;

    public CachedDriverManager(String url, PlcConnectionFactory connectionFactory) {
        this(url, connectionFactory, 1000);
    }

    /**
     * @param url               Url that this connection is for
     * @param connectionFactory Factory to create a suitable connection.
     * @param timeoutMillis     Time out in milliseonds
     */
    public CachedDriverManager(String url, PlcConnectionFactory connectionFactory, int timeoutMillis) {
        logger.info("Creating new cached Connection for url {} with timeout {} ms", url, timeoutMillis);
        this.url = url;
        this.connectionFactory = connectionFactory;
        this.timeoutMillis = timeoutMillis;

        // MBean
        try {
            ManagementFactory.getPlatformMBeanServer().registerMBean(this, new ObjectName("org.pragmaticindustries.cockpit.plc:name=cached-driver-manager,url=\"" + url + "\""));
        } catch (Exception ignore) {
        }
    }

    public synchronized void returnConnection(PlcConnection activeConnection) {
        logger.debug("Borrowed Connection is closed and returned.");
        // Stop Watchdog
        cancelWatchdog();
        if (state.get() == ConnectionState.DISCONNECTED) {
            // Getting Disconnected Connection, nothing to do.
            logger.trace("Connection allready disconnected");
            return;
        }
        if (state.get() != ConnectionState.BORROWED) {
            logger.warn("Connection was returned, although it is not borrowed, currently.");
        }
        this.borrowedConnection = null;
        setState(ConnectionState.AVAILABLE);
        // Check the queue
        checkQueue();
        logger.trace("Connection successfully returned");
    }

    private void setState(ConnectionState available) {
        logger.trace("Setting State from {} to {}", this.state.get(), available);
        this.state.set(available);
    }

    /**
     * used to signal, that the connection does not work and has to be re-established.
     */
    public synchronized void handleBrokenConnection() {
        logger.debug("Connection was detected as broken and is invalidated in Cached Manager");
        // Stop Watchdog
        cancelWatchdog();
        if (state.get() != ConnectionState.BORROWED) {
            // Log at least a WARN???
            logger.warn("Broken Connection was returned, although it is not borrowed, currently.");
        }
        this.borrowedConnection = null;
        try {
            this.activeConnection.close();
        } catch (Exception e) {
            logger.debug("Unable to Close 'broken' Connection", e);
        }
        this.activeConnection = null;
        setState(ConnectionState.DISCONNECTED);
    }

    public boolean isConnectionAvailable() {
        return this.getState().equals(ConnectionState.AVAILABLE);
    }

    /**
     * This call now waits (with the timeout given in constructor) until it fails
     * or returns a valid connection in this window.
     *
     * @throws PlcConnectionException if connection cannot be established
     */
    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        if (!this.url.equals(url)) {
            throw new IllegalArgumentException("This Cached Driver Manager only supports the Connection " + url);
        }
        synchronized (this) {
            logger.trace("current queue size before check {}", queue.size());
            if (queue.isEmpty() && isConnectionAvailable()) {
                logger.trace("queue is empty and a connection is available");
                return getConnection_(url);
            } else {
                logger.trace("Getting a connection and instantly close it");
                // At least trigger a connection
                try {
                    getConnection_(url).close();
                } catch (Exception ignore) {
                }
            }
        }
        CompletableFuture<PlcConnection> future = new CompletableFuture<>();
        synchronized (this) {
            logger.trace("current queue size before add {}", queue.size());
            queue.add(future);
        }
        try {
            return future.get(timeoutMillis, TimeUnit.MILLISECONDS);
        } catch (ExecutionException | TimeoutException e) {
            throw new PlcConnectionException("No Connection Available, timed out while waiting in queue.", e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcConnectionException("No Connection Available, interrupted while waiting in queue.", e);
        } finally {
            future.cancel(true);
        }
    }

    /**
     * Private Impl.
     */
    private synchronized PlcConnection getConnection_(String url) throws PlcConnectionException {
        logger.trace("Current State {}", this.state.get());
        switch (state.get()) {
            case AVAILABLE:
                logger.debug("Connection was requested and is available, thus, returning Chached Connection for usage");
                setState(ConnectionState.BORROWED);
                this.numberOfBorrows.incrementAndGet();

                this.borrowedConnection = new CachedPlcConnection(this, activeConnection);

                // Set the Borrwed Counter
                startWatchdog(this.borrowedConnection);

                return this.borrowedConnection;
            case DISCONNECTED:
                logger.debug("Connection was requested but no connection is active, trying to establish a Connection");
                // Initialize Connection
                setState(ConnectionState.CONNECTING);
                this.numberOfConnects.incrementAndGet();

                // Start Connection in Background
                CompletableFuture.runAsync(() -> {
                    logger.debug("Starting to establish Connection");
                    try {
                        PlcConnection connection = this.connectionFactory.create();
                        logger.debug("Connection successfully established");
                        synchronized (this) {
                            this.activeConnection = connection;
                            setState(ConnectionState.AVAILABLE);
                            // Now See if there is someone waiting in the line
                            checkQueue();
                            logger.trace("Inline queue check succeeded");
                        }
                    } catch (Exception e) {
                        logger.warn("Unable to establish connection to PLC {}", url, e);
                        setState(ConnectionState.DISCONNECTED);
                    }
                });

                this.numberOfRejections.incrementAndGet();
                throw new PlcConnectionException("No Connection Available, Starting Connection");
            case CONNECTING:
                // We cannot give a Connection
                logger.debug("Connection was requsted, but currently establishing one, so none available");
                this.numberOfRejections.incrementAndGet();
                throw new PlcConnectionException("No Connection Available, Currently Connecting");
            case BORROWED:
                // We cannot give a Connection
                logger.debug("Connection was requsted, but Connection currently is borrowed, so none available");
                this.numberOfRejections.incrementAndGet();
                throw new PlcConnectionException("No Connection Available, its in Use");
        }
        throw new IllegalStateException();
    }

    /**
     * Checks if someone is waiting in the line to get the connection.
     */
    private synchronized void checkQueue() {
        logger.debug("Connection is available, checking if someone is waiting in the queue...");
        CompletableFuture<PlcConnection> next;
        logger.trace("current queue size before check queue {}", queue.size());
        while ((next = queue.poll()) != null) {
            if (next.isCancelled()) {
                logger.trace("Cleaning up already timed out connection...");
                continue;
            }
            // Not timed out, give this connection
            try {
                next.complete(getConnection_(url));
                return;
            } catch (PlcConnectionException e) {
                logger.debug("Got an Exception on fetching a connection", e);
            }
        }
        logger.trace("check queue ended");
    }

    private void startWatchdog(CachedPlcConnection connection) {
        borrowWatchdog = executorService.schedule(() -> {
            // Just close the borrowed connection
            logger.warn("Watchdog detected a long borrowed connection, will be forcefully closed!");
            this.numberOfWatchdogs.incrementAndGet();
            handleBrokenConnection();
            try {
                connection.close();
            } catch (Exception e) {
                logger.warn("Unable to close the borrowed Connection from Watchdog", e);
            }
        }, LONG_BORROW_WATCHDOG_TIMEOUT_MS, TimeUnit.MILLISECONDS);
    }

    private void cancelWatchdog() {
        borrowWatchdog.cancel(false);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new NotImplementedException("");
    }

    public ConnectionState getState() {
        return state.get();
    }

    enum ConnectionState {
        DISCONNECTED,
        CONNECTING,
        AVAILABLE,
        BORROWED;
    }

    // JMX

    @Override
    public String getStateString() {
        return getState().toString();
    }

    @Override
    public int getNumberOfConnects() {
        return numberOfConnects.get();
    }

    @Override
    public int getNumberOfBorrows() {
        return numberOfBorrows.get();
    }

    @Override
    public int getNumberOfRejections() {
        return this.numberOfRejections.get();
    }

    @Override
    public int getNumberOfWachtdogs() {
        return numberOfWatchdogs.get();
    }

    @Override
    public int getQueueSize() {
        return this.queue.size();
    }

    @Override
    public synchronized void triggerReconnect() {
        logger.info("Disconnecting current connection, was triggered from external via JMX");
        handleBrokenConnection();
        if (this.state.get() == ConnectionState.BORROWED) {
            try {
                borrowedConnection.close();
            } catch (Exception e) {
                logger.warn("Unable to close the borrowed Connection from JMX", e);
            }
        }
    }
}



