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
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.time.Duration;
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;

class ConnectionContainer {
    private PlcConnection connection;
    private boolean closed = false;
    private final Duration maxLeaseTime;
    private final Queue<CompletableFuture<PlcConnection>> queue;

    private LeasedPlcConnection leasedConnection;

    public boolean isClosed() {
        return closed;
    }
    public PlcConnection getRawConnection() {
        return connection;
    }

    public ConnectionContainer(PlcConnection connection, Duration maxLeaseTime) {
        this.connection = connection;
        this.maxLeaseTime = maxLeaseTime;
        this.queue = new LinkedList<>();
        this.leasedConnection = null;
    }

    public synchronized Future<PlcConnection> lease() {
        CompletableFuture<PlcConnection> connectionFuture = new CompletableFuture<>();
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
    public synchronized void close(){
        CompletableFuture<PlcConnection> leaseFuture;
        while((leaseFuture = queue.poll())!=null){
            leaseFuture.complete(null);
        }
        leasedConnection = null;
        connection = null;
        closed = true;
    }
    public synchronized void returnConnection(LeasedPlcConnection returnedLeasedConnection) {
        if(closed){
            return;
        }
        if(returnedLeasedConnection != leasedConnection) {
            throw new PlcRuntimeException("Error trying to return lease from invalid connection");
        }

        // If the queue is empty, simply return.
        if(queue.isEmpty()) {
            leasedConnection = null;
            return;
        }

        // Create a new lease and complete the next future in the queue with this.
        leasedConnection = new LeasedPlcConnection(this, connection, maxLeaseTime);
        CompletableFuture<PlcConnection> leaseFuture = queue.poll();
        leaseFuture.complete(leasedConnection);
    }

}
