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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;

class ConnectionContainer {
    private static final Logger LOGGER = LoggerFactory.getLogger(ConnectionContainer.class);
    private final PlcConnectionManager connectionManager;
    private final String connectionUrl;
    private final Duration maxLeaseTime;
    private final Queue<CompletableFuture<PlcConnection>> queue;

    private PlcConnection connection;
    private LeasedPlcConnection leasedConnection;

    public ConnectionContainer(PlcConnectionManager connectionManager, String connectionUrl, Duration maxLeaseTime) {
        this.connectionManager = connectionManager;
        this.connectionUrl = connectionUrl;
        this.maxLeaseTime = maxLeaseTime;
        this.queue = new LinkedList<>();
        this.connection = null;
        this.leasedConnection = null;
    }

    public synchronized Future<PlcConnection> lease() {
        CompletableFuture<PlcConnection> connectionFuture = new CompletableFuture<>();

        // Try to get a new connection, if we haven't got one yet.
        if(connection == null) {
            try {
                connection = connectionManager.getConnection(connectionUrl);
            } catch (PlcConnectionException e) {
                LOGGER.warn("Exception while getting connection for lease", e);
                connectionFuture.completeExceptionally(e);
                return connectionFuture;
            }
        }

        // If the connection is currently idle, return the connection immediately.
        if (leasedConnection == null) {
            leasedConnection = new LeasedPlcConnection(this, connection, maxLeaseTime);
            connectionFuture.complete(leasedConnection);
        }
        // Otherwise queue the future up for completion as soon as the connection is returned.
        else {
            queue.add(connectionFuture);
        }
        return connectionFuture;
    }

    public synchronized void returnConnection(LeasedPlcConnection returnedLeasedConnection, boolean invalidateConnection) {
        if(returnedLeasedConnection != leasedConnection) {
            LOGGER.error("Error trying to return lease from invalid connection: returned={} leased={}",
                returnedLeasedConnection, leasedConnection);
            throw new PlcRuntimeException("Error trying to return lease from invalid connection");
        }

        // If something happened while using the connection, invalidate this one and create a new connection.
        if(invalidateConnection) {
            // Close the old connection.
            try {
                connection.close();
            } catch (Exception e) {
                // We're ignoring this as we have no idea, what state the connection is in.
                // Nevertheless, it is polite to say something in logs about this situation.
                LOGGER.warn("Exception while closing connection", e);
            }

            // Try to get a new connection.
            try {
                connection = connectionManager.getConnection(connectionUrl);
            } catch (PlcConnectionException e) {
                // If something goes wrong, close all waiting futures exceptionally.
                LOGGER.warn("Can't get connection for {} complete queue items exceptionally", connectionUrl, e);
                queue.forEach(future -> future.completeExceptionally(e));
                connection = null;
            }
        }

        // If the queue is empty, simply return.
        if(queue.isEmpty()) {
            leasedConnection = null;
            return;
        }

        // Create a new lease and complete the next future in the queue with this.
        leasedConnection = new LeasedPlcConnection(this, connection, maxLeaseTime);
        CompletableFuture<PlcConnection> leaseFuture = queue.poll();
        if(leaseFuture != null) {
            leaseFuture.complete(leasedConnection);
        }
    }

}
