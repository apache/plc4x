/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.connectionpool;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.PlcDriverManagerInterface;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.exceptions.NoConnectionAvailableException;
import org.apache.plc4x.java.exceptions.NotConnectedException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.annotation.PreDestroy;
import java.util.HashMap;
import java.util.Map;


/**
 * Plc Connection Manager to act as a controller to give threads access to a connection in the pool only if it is not
 * used by a different thread
 */
public class PlcConnectionManager implements PlcDriverManagerInterface {
    private static final Logger logger = LoggerFactory.getLogger(PlcConnectionManager.class);
    //TODO: implement connection Pool with more than one connection
    private Map<String, PlcConnection> connections = new HashMap<>();
    private Map<String, Boolean> isBlockedMap = new HashMap<>();
    private PlcDriverManager plcDriverManager;


    public PlcConnectionManager(Map<String, PlcConnection> connectionProxyMap, Map<String, Boolean> isBlockedMap) {
        this.connections = connectionProxyMap;
        this.isBlockedMap = isBlockedMap;
        plcDriverManager = new PlcDriverManager();
    }


    /**
     * Connects to a PLC using the given plc connection string when no connection was requested jet and builds the pool.
     * Else it return the connection from the pool in a connection proxy.
     *
     * @param connectionString plc url with protocol prefix to select correct connection
     * @return PlcConnection in a proxy
     * @throws PlcConnectionException if connection was unsuccessfull, or the connection all other connections are
     * being used by different threads
     */
    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        if (!connections.containsKey(connectionString)) {
            connections.put(connectionString,plcDriverManager.getConnection(connectionString));
            isBlockedMap.put(connectionString, false);
        }
        if (connections.containsKey(connectionString) && isBlockedMap.containsKey(connectionString)) {
            // TODO: await set time till connection maybe available
            if (!isBlockedMap.get(connectionString)) {
                return tryConnect(connectionString);
            } else {
                throw new NoConnectionAvailableException("Connection not available");
            }
        } else {
            throw new NotConnectedException("Connection was not possible in the first place");
        }
    }

    /**
     * Connects to a PLC using the given plc connection string and authentification when no connection was requested jet and builds the pool.
     * Else it return the connection from the pool in a connection proxy.
     *
     * @param connectionString plc url with protocol prefix to select correct connection
     * @param authentication authentication credentials.
     * @return PlcConnection in a proxy
     * @throws PlcConnectionException if connection was unsuccessfull, or the connection all other connections are
     * being used by different threads
     */
    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication) throws PlcConnectionException {
        if (!connections.containsKey(connectionString)) {
            connections.put(connectionString,plcDriverManager.getConnection(connectionString,authentication));
            isBlockedMap.put(connectionString, false);
        }
        if (connections.containsKey(connectionString) && isBlockedMap.containsKey(connectionString)) {
            // TODO: await set time till connection maybe available
            if (!isBlockedMap.get(connectionString)) {
                return tryConnect(connectionString);
            } else {
                throw new NoConnectionAvailableException("Connection not available");
            }
        } else {
            throw new NotConnectedException("Connection was not possible in the first place");
        }
    }

    /**
     * helper function to check if the connection can be established if it is not already established
     * and get the available connection from map and return a proxy that tells this manager when the connection is closed
     *
     * @param connectionString of source
     * @return PlcConnection to return in getConnection
     * @throws NotConnectedException if connection could not be established
     */
    private PlcConnectionProxy tryConnect(String connectionString) throws NotConnectedException {
        try {
            if (connections.get(connectionString).isConnected()) {
                isBlockedMap.put(connectionString, true);
                return new PlcConnectionProxy(this, connectionString, connections.get(connectionString));
            } else {
                connections.get(connectionString).connect();
                isBlockedMap.put(connectionString, true);
                return new PlcConnectionProxy(this, connectionString, connections.get(connectionString));
            }
        } catch (PlcConnectionException ex) {
            throw new NotConnectedException(ex);
        }
    }

    /**
     * After using the connection received from getReader or getConnection to release the connection again
     * @param connectionString for which the connection should be released again
     */
    public void returnConnection(String connectionString) {
        isBlockedMap.put(connectionString, false);
    }

    /**
     * close all connections before closing
     */
    @PreDestroy
    public void close() {
        for (PlcConnection plcConnection : connections.values()) {
            try {
                plcConnection.close();
            } catch (Exception ex) {
                logger.error("Error disconnecting", ex);
            }
        }
    }
}
