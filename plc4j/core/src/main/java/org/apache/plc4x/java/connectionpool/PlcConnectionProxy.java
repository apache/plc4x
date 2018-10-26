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


import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.connection.PlcWriter;

import java.util.Optional;

/**
 * wrapper for PlcConnections to make access thread safe
 * Tells ConnectionManager when the connection is not used anymore
 */
public class PlcConnectionProxy implements PlcConnection {
    private final PlcConnectionManager parent;
    private final String connectionString;
    private final PlcConnection plcConnection;
    private boolean closed;


    public PlcConnectionProxy(PlcConnectionManager parent, String connectionString, PlcConnection plcConnection) {
        this.parent = parent;
        this.connectionString = connectionString;
        this.plcConnection = plcConnection;
        closed = false;
    }

    /**
     * connect should already be handled in the connection manager
     * returns without doing anything
     */
    @Override
    public void connect() {
    }

    @Override
    public boolean isConnected() {
        if(closed){
            return false;
        }else {
            return plcConnection.isConnected();
        }
    }

    /**
     * tell PlcConnectionManager that the connection is free again
     */
    @Override
    public void close() {
        closed = true;
        parent.returnConnection(connectionString);
    }

    @Override
    public Optional<PlcReader> getReader() {
        if(closed){
            return Optional.empty();
        }else {
            return plcConnection.getReader();
        }
    }

    @Override
    public Optional<PlcWriter> getWriter() {
        if(closed) {
            return Optional.empty();
        }else{
            return plcConnection.getWriter();
        }
    }

    @Override
    public Optional<PlcSubscriber> getSubscriber() {
        if(closed){
            return Optional.empty();
        }else {
            return plcConnection.getSubscriber();
        }
    }


}
