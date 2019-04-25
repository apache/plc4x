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

package org.apache.plc4x.interop.impl;

import org.apache.commons.lang.exception.ExceptionUtils;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.interop.*;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.thrift.TException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;

public class Handler implements InteropServer.Iface {

    private static final Logger LOGGER = LoggerFactory.getLogger(Handler.class);

    private final AtomicLong connectionCounter = new AtomicLong(0);

    private final PlcDriverManager driverManager;
    private final Map<Long, String> connectionStrings;
    private final Map<Long, PlcConnection> connections;

    public Handler(PlcDriverManager driverManager) {
        this.driverManager = driverManager;
        this.connections = new ConcurrentHashMap<>();
        this.connectionStrings = new ConcurrentHashMap<>();
    }


    @Override public ConnectionHandle connect(String connectionString) throws TException {
        LOGGER.debug("Receiving new connect request to '{}'", connectionString);
        try {
            long id = connectionCounter.getAndIncrement();

            final PlcConnection connection = driverManager.getConnection(connectionString);
            connection.connect();

            LOGGER.debug("Established connection to '{}' with handle {}", connectionString, id);

            this.connections.put(id, connection);
            this.connectionStrings.put(id, connectionString);
            return new ConnectionHandle(id);
        } catch (PlcConnectionException e) {
            LOGGER.warn("Unable to start a connection to url '" + connectionString + "'", e);
            throw new PlcException(connectionString, e.getMessage());
        }
    }

    @Override public void close(ConnectionHandle handle) throws TException {
        LOGGER.debug("Receiving new close request for handle {}", handle.getConnectionId());
        if (!connections.containsKey(handle)) {
            LOGGER.warn("Handle for close request {} does not exist. Perhaps already closed?", handle.getConnectionId());
            return;
        }
        try {
            connections.get(handle).close();
            connections.remove(handle);
        } catch (Exception e) {
            LOGGER.warn("Unable to close the conn / remove the handle", e);
        }

    }


    @Override public Response execute(ConnectionHandle handle, Request request) throws TException {
        LOGGER.debug("Executing " + request);
        if (request.getFields() == null) {
            throw new PlcException(
                connectionStrings.get(handle.getConnectionId()),
                "No fields given in the request!");
        }
        try {
            final PlcReadRequest.Builder builder = connections.get(handle.getConnectionId()).readRequestBuilder();
            for (Map.Entry<String, String> entry : request.getFields().entrySet()) {
                builder.addItem(entry.getKey(), entry.getValue());
            }
            final PlcReadResponse response = builder.build().execute().get(1_000L, TimeUnit.MILLISECONDS);

            final HashMap<String, FieldResponse> resultMap = new HashMap<>();
            for (String key : request.getFields().keySet()) {
                final PlcResponseCode responseCode = response.getResponseCode(key);

                final FieldResponse fieldResponse = new FieldResponse(convertResponseCode(responseCode));

                if (PlcResponseCode.OK.equals(responseCode)) {
                    if (response.isValidBoolean(key)) {
                        fieldResponse.setBoolValue(response.getBoolean(key));
                    }
                    if (response.isValidLong(key)) {
                        fieldResponse.setLongValue(response.getLong(key));
                    }
                    if (response.isValidDouble(key)) {
                        fieldResponse.setDoubleValue(response.getDouble(key));
                    }
                    if (response.isValidString(key)) {
                        fieldResponse.setStringValue(response.getString(key));
                    }
                }

                resultMap.put(key, fieldResponse);
            }
            return new Response(resultMap);
        } catch (Exception e) {
            LOGGER.warn("Exception during execution of request '" + request + "' for handle " + handle.getConnectionId(), e);
            throw new PlcException(
                connectionStrings.get(handle.getConnectionId()),
                ExceptionUtils.getStackTrace(e));
        }
    }

    private RESPONSE_CODE convertResponseCode(PlcResponseCode responseCode) {
        switch (responseCode) {
            case OK:
                return RESPONSE_CODE.OK;
            case NOT_FOUND:
                return RESPONSE_CODE.NOT_FOUND;
            case ACCESS_DENIED:
                return RESPONSE_CODE.ACCESS_DENIED;
            case INTERNAL_ERROR:
                return RESPONSE_CODE.INTERNAL_ERROR;
            case INVALID_ADDRESS:
                return RESPONSE_CODE.INVALID_ADDRESS;
            case INVALID_DATATYPE:
                return RESPONSE_CODE.INVALID_DATATYPE;
            case RESPONSE_PENDING:
                return RESPONSE_CODE.RESPONSE_PENDING;
            default:
                throw new NotImplementedException("This response code is not implemented!");
        }
    }

}
