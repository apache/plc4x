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
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;

import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Date;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.CompletableFuture;

public class LeasedPlcConnection implements PlcConnection {

    private final ConnectionContainer connectionContainer;
    private PlcConnection connection;

    public LeasedPlcConnection(ConnectionContainer connectionContainer, PlcConnection connection, Duration maxUseTime) {
        this.connectionContainer = connectionContainer;
        this.connection = connection;
        Timer usageTimer = new Timer();
        usageTimer.schedule(new TimerTask() {
            @Override
            public void run() {
                close();
            }
        }, Date.from(LocalDateTime.now().plusNanos(maxUseTime.toNanos()).atZone(ZoneId.systemDefault()).toInstant()));
    }

    @Override
    public synchronized void close() {
        // Make the connection unusable.
        connection = null;

        // Tell the connection container that the connection is free to be reused.
        connectionContainer.returnConnection(this);
    }

    @Override
    public void connect() throws PlcConnectionException {
        throw new PlcConnectionException("Error connecting leased connection");
    }

    @Override
    public boolean isConnected() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.isConnected();
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.getMetadata();
    }

    @Override
    public CompletableFuture<Void> ping() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.ping();
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.readRequestBuilder();
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.writeRequestBuilder();
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.subscriptionRequestBuilder();
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.unsubscriptionRequestBuilder();
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.browseRequestBuilder();
    }

}
