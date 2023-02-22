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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

public class CachedPlcConnectionManager implements PlcConnectionManager {

    private static final Logger LOG = LoggerFactory.getLogger(CachedPlcConnectionManager.class);

    private final PlcConnectionManager connectionManager;
    private final Duration maxLeaseTime;
    private final Duration maxWaitTime;

    private final Map<String, ConnectionContainer> connectionContainers;

    public static Builder getBuilder() {
        return new Builder(new DefaultPlcDriverManager());
    }

    public static Builder getBuilder(PlcConnectionManager connectionManager) {
        return new Builder(connectionManager);
    }

    public CachedPlcConnectionManager(PlcConnectionManager connectionManager, Duration maxLeaseTime, Duration maxWaitTime) {
        this.connectionManager = connectionManager;
        this.maxLeaseTime = maxLeaseTime;
        this.maxWaitTime = maxWaitTime;
        this.connectionContainers = new HashMap<>();
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        return getConnection(url,null);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        ConnectionContainer connectionContainer;
        synchronized (connectionContainers) {
            connectionContainer = connectionContainers.get(url);
            if (connectionContainer == null || connectionContainer.isClosed()) {
                LOG.debug("Creating new connection");

                // Establish the real connection to the plc
                PlcConnection connection;
                if(authentication!=null) {
                    connection = connectionManager.getConnection(url,authentication);
                } else{
                    connection = connectionManager.getConnection(url);
                }
                connectionContainer = new ConnectionContainer(connection,maxLeaseTime);
                connectionContainers.put(url, connectionContainer);
            } else {
                LOG.debug("Reusing exising connection");
                if(connectionContainer.getRawConnection()!=null && !connectionContainer.getRawConnection().isConnected()){
                    connectionContainer.getRawConnection().connect();
                }
            }
        }

        // Get a lease (a future for a connection)
        Future<PlcConnection> leaseFuture = connectionContainer.lease();
        try {
            return leaseFuture.get(this.maxWaitTime.toMillis(), TimeUnit.MILLISECONDS);
        } catch (ExecutionException | InterruptedException | TimeoutException e) {
            connectionContainer.close();
            connectionContainers.remove(url);
            throw new PlcConnectionException("Error acquiring lease for connection", e);
        }
    }

    public void destroy(){
        connectionContainers.values().forEach(ConnectionContainer::close);
        connectionContainers.clear();
    }
    public static class Builder {

        private final PlcConnectionManager connectionManager;
        private Duration maxLeaseTime;
        private Duration maxWaitTime;

        public Builder(PlcConnectionManager connectionManager) {
            this.connectionManager = connectionManager;
            this.maxLeaseTime = Duration.ofSeconds(4);
            this.maxWaitTime = Duration.ofSeconds(20);
        }

        public CachedPlcConnectionManager build() {
            return new CachedPlcConnectionManager(this.connectionManager, this.maxLeaseTime, this.maxWaitTime);
        }

        public CachedPlcConnectionManager.Builder withMaxLeaseTime(Duration maxLeaseTime) {
            this.maxLeaseTime = maxLeaseTime;
            return this;
        }

        public CachedPlcConnectionManager.Builder withMaxWaitTime(Duration maxWaitTime) {
            this.maxWaitTime = maxWaitTime;
            return this;
        }
    }

}
