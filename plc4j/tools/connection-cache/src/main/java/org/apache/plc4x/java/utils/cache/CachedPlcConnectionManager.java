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

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.utils.cache.exceptions.PlcConnectionManagerClosedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;

public class CachedPlcConnectionManager implements PlcConnectionManager, AutoCloseable {

    private static final Logger LOG = LoggerFactory.getLogger(CachedPlcConnectionManager.class);

    private final PlcConnectionManager connectionManager;
    private final Duration maxLeaseTime;
    private final Duration maxWaitTime;
    private final Duration maxIdleTime;

    private final Map<String, ConnectionContainer> connectionContainers;

    private final AtomicBoolean closed = new AtomicBoolean(false);

    public static Builder getBuilder() {
        return new Builder(new DefaultPlcDriverManager());
    }

    public static Builder getBuilder(PlcConnectionManager connectionManager) {
        return new Builder(connectionManager);
    }

    public CachedPlcConnectionManager(PlcConnectionManager connectionManager, Duration maxLeaseTime, Duration maxWaitTime, Duration maxIdleTime) {
        this.connectionManager = connectionManager;
        this.maxLeaseTime = maxLeaseTime;
        this.maxWaitTime = maxWaitTime;
        this.maxIdleTime = maxIdleTime;
        this.connectionContainers = new HashMap<>();
    }

    /**
     * @return set of connection-urls the CachedPlcConnectionManager is currently managing.
     */
    public Set<String> getCachedConnections() {
        synchronized (connectionContainers) {
            return connectionContainers.keySet();
        }
    }

    /**
     * Removes a given connection from the cache (Should only be used in order to remove somehow broken connections).
     * @param url url of the connection that should be removed.
     */
    public void removeCachedConnection(String url) {
        synchronized (connectionContainers) {
            connectionContainers.remove(url);
        }
    }

    public PlcConnection getConnection(String url) throws PlcConnectionException {
        // If the connection manager is already closed, abort.
        if(closed.get()) {
            throw new PlcConnectionManagerClosedException();
        }

        // Get a connection container for the given url.
        ConnectionContainer connectionContainer;
        synchronized (connectionContainers) {
            connectionContainer = connectionContainers.get(url);
            if (connectionContainer == null) {
                LOG.debug("Creating new connection");

                // Crate a connection container to manage handling this connection
                connectionContainer = new ConnectionContainer(connectionManager, url, maxLeaseTime, maxIdleTime,
                    closeConnection -> {
                        removeCachedConnection(closeConnection);
                        return null;
                    });
                connectionContainers.put(url, connectionContainer);
            } else {
                LOG.debug("Reusing exising connection");
            }
        }

        // Get a lease (a future for a connection)
        Future<PlcConnection> leaseFuture = connectionContainer.lease();
        try {
            return leaseFuture.get(this.maxWaitTime.toMillis(), TimeUnit.MILLISECONDS);
        } catch (ExecutionException | InterruptedException | TimeoutException e) {
            throw new PlcConnectionException("Error acquiring lease for connection", e);
        }
    }

    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("the cached driver manager currently doesn't support authentication");
    }

    @Override
    public void close() throws Exception {
        // Set the cache to "closed" so no new connections can be requested.
        closed.set(true);

        // Tell all connections to close themselves.
        connectionContainers.forEach((connectionString, connectionContainer) -> {
            connectionContainer.close();
        });
    }

    public static class Builder {

        private final PlcConnectionManager connectionManager;
        private Duration maxLeaseTime;
        private Duration maxWaitTime;
        private Duration maxIdleTime;

        public Builder(PlcConnectionManager connectionManager) {
            this.connectionManager = connectionManager;
            this.maxLeaseTime = Duration.ofSeconds(4);
            this.maxWaitTime = Duration.ofSeconds(20);
            this.maxIdleTime = Duration.ofMinutes(5);
        }

        public CachedPlcConnectionManager build() {
            return new CachedPlcConnectionManager(
                this.connectionManager, this.maxLeaseTime, this.maxWaitTime, this.maxIdleTime);
        }

        public CachedPlcConnectionManager.Builder withMaxLeaseTime(Duration maxLeaseTime) {
            this.maxLeaseTime = maxLeaseTime;
            return this;
        }

        public CachedPlcConnectionManager.Builder withMaxWaitTime(Duration maxWaitTime) {
            this.maxWaitTime = maxWaitTime;
            return this;
        }

        public CachedPlcConnectionManager.Builder withMaxIdleTime(Duration maxIdleTime) {
            this.maxIdleTime = maxIdleTime;
            return this;
        }
    }

}
