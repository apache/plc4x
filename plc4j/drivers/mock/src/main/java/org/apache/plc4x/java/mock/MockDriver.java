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
package org.apache.plc4x.java.mock;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.field.MockField;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Mocking Driver that keeps a Map of references to Connections so that you can fetch a reference to a connection
 * which will be acquired by someone else (via the connection string).
 * This allows for efficient Mocking.
 */
public class MockDriver implements PlcDriver {

    private Map<String, PlcConnection> connectionMap = new ConcurrentHashMap<>();

    @Override
    public String getProtocolCode() {
        return "mock";
    }

    @Override
    public String getProtocolName() {
        return "Mock Protocol Implementation";
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        return getConnection(url, null);
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        String deviceName = url.substring(5);
        if (deviceName.isEmpty()) {
            throw new PlcConnectionException("Invalid URL: no device name given.");
        }
        return connectionMap.computeIfAbsent(deviceName, name -> new MockConnection(authentication));
    }

    @Override
    public MockField prepareField(String query){
        return MockField.of(query);
    }

}
