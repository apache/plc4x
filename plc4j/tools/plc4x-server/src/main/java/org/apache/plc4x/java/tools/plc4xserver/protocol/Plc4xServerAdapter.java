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
package org.apache.plc4x.java.tools.plc4xserver.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.plc4x.readwrite.*;
import org.apache.plc4x.java.utils.connectionpool2.PooledDriverManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4xServerAdapter extends ChannelInboundHandlerAdapter {

    private final Logger logger = LoggerFactory.getLogger(Plc4xServerAdapter.class);

    private final PooledDriverManager driverManager;
    private final AtomicInteger connectionIdGenerator;
    private final ConcurrentHashMap<Integer, String> connectionUrls;

    public Plc4xServerAdapter() {
        driverManager = new PooledDriverManager();
        connectionIdGenerator = new AtomicInteger(1);
        connectionUrls = new ConcurrentHashMap<>();
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) {
        if (msg instanceof Plc4xMessage) {
            final Plc4xMessage plc4xMessage = (Plc4xMessage) msg;
            switch (plc4xMessage.getRequestType()) {
                case CONNECT_REQUEST: {
                    Plc4xConnectRequest request = (Plc4xConnectRequest) plc4xMessage;
                    try (final PlcConnection ignored = driverManager.getConnection(request.getConnectionString())) {
                        //connection.ping().get();
                        final int connectionId = connectionIdGenerator.getAndIncrement();
                        connectionUrls.put(connectionId, request.getConnectionString());
                        Plc4xConnectResponse response = new Plc4xConnectResponse(
                            request.getRequestId(), connectionId, Plc4xResponseCode.OK);
                        ctx.writeAndFlush(response);
                    } catch (Exception e) {
                        Plc4xConnectResponse response = new Plc4xConnectResponse(
                            request.getRequestId(), 0, Plc4xResponseCode.INVALID_ADDRESS);
                        ctx.writeAndFlush(response);
                    }
                    break;
                }
                case READ_REQUEST: {
                    final Plc4xReadRequest request = (Plc4xReadRequest) plc4xMessage;
                    String connectionUrl = connectionUrls.get(request.getConnectionId());
                    try (final PlcConnection connection = driverManager.getConnection(connectionUrl)) {
                        // Build a read request for all fields in the request.
                        final PlcReadRequest.Builder builder = connection.readRequestBuilder();
                        for (Plc4xFieldRequest requestField : request.getFields()) {
                            builder.addFieldAddress(requestField.getField().getName(), requestField.getField().getFieldQuery());
                        }
                        final PlcReadRequest rr = builder.build();

                        // Execute the query.
                        // (It has to be synchronously when working with the connection cache)
                        final PlcReadResponse apiReadResponse = rr.execute().get();

                        // Create the response.
                        List<Plc4xFieldValueResponse> fields = new ArrayList<>(apiReadResponse.getFieldNames().size());
                        for (Plc4xFieldRequest plc4xRequestField : request.getFields()) {
                            final PlcResponseCode responseCode = apiReadResponse.getResponseCode(plc4xRequestField.getField().getName());
                            Plc4xResponseCode resCode;
                            Plc4xValueType valueType;
                            PlcValue value;
                            if(responseCode == PlcResponseCode.OK) {
                                resCode = Plc4xResponseCode.OK;
                                value = apiReadResponse.getPlcValue(plc4xRequestField.getField().getName());
                                final String valueTypeName = value.getClass().getSimpleName();
                                // Cut off the "Plc" prefix to get the name of the PlcValueType.
                                valueType = Plc4xValueType.valueOf(valueTypeName.substring(3));
                            } else {
                                resCode = Plc4xResponseCode.INVALID_ADDRESS;
                                value = null;
                                valueType = Plc4xValueType.NULL;
                            }
                            fields.add(new Plc4xFieldValueResponse(
                                plc4xRequestField.getField(), resCode, valueType, value));
                        }
                        Plc4xReadResponse response = new Plc4xReadResponse(
                            request.getRequestId(), request.getConnectionId(), Plc4xResponseCode.OK, fields);

                        // Send the response.
                        ctx.writeAndFlush(response);
                    } catch (Exception e) {
                        logger.error("Error executing request", e);
                        Plc4xReadResponse response = new Plc4xReadResponse(
                            request.getRequestId(), request.getConnectionId(),
                            Plc4xResponseCode.INVALID_ADDRESS, Collections.emptyList());
                        ctx.writeAndFlush(response);
                    }
                    break;
                }

                case WRITE_REQUEST:
                    final Plc4xWriteRequest plc4xWriteRequest = (Plc4xWriteRequest) plc4xMessage;
                    String connectionUrl = connectionUrls.get(plc4xWriteRequest.getConnectionId());
                    try (final PlcConnection connection = driverManager.getConnection(connectionUrl)) {
                        // Build a write request for all fields in the request.
                        final PlcWriteRequest.Builder builder = connection.writeRequestBuilder();
                        for (Plc4xFieldValueRequest plc4xRequestField : plc4xWriteRequest.getFields()) {
                            builder.addFieldAddress(plc4xRequestField.getField().getName(),
                                plc4xRequestField.getField().getFieldQuery(), plc4xRequestField.getValue().getObject());
                        }
                        final PlcWriteRequest apiWriteRequest = builder.build();

                        // Execute the query
                        // (It has to be synchronously when working with the connection cache)
                        final PlcWriteResponse apiWriteResponse = apiWriteRequest.execute().get();

                        // Create the response.
                        List<Plc4xFieldResponse> plc4xFields =
                            new ArrayList<>(apiWriteResponse.getFieldNames().size());
                        for (Plc4xFieldValueRequest plc4xRequestField : plc4xWriteRequest.getFields()) {
                            final PlcResponseCode apiResponseCode =
                                apiWriteResponse.getResponseCode(plc4xRequestField.getField().getName());
                            Plc4xResponseCode resCode;
                            if(apiResponseCode == PlcResponseCode.OK) {
                                resCode = Plc4xResponseCode.OK;
                            } else {
                                resCode = Plc4xResponseCode.INVALID_ADDRESS;
                            }
                            plc4xFields.add(new Plc4xFieldResponse(plc4xRequestField.getField(), resCode));
                        }
                        Plc4xWriteResponse plc4xWriteResponse = new Plc4xWriteResponse(
                            plc4xWriteRequest.getRequestId(), plc4xWriteRequest.getConnectionId(),
                            Plc4xResponseCode.OK, plc4xFields);

                        // Send the response.
                        ctx.writeAndFlush(plc4xWriteResponse);
                    } catch (Exception e) {
                        logger.error("Error executing request", e);
                        Plc4xWriteResponse response = new Plc4xWriteResponse(
                            plc4xWriteRequest.getRequestId(), plc4xWriteRequest.getConnectionId(),
                            Plc4xResponseCode.INVALID_ADDRESS, Collections.emptyList());
                        ctx.writeAndFlush(response);
                    }
                    break;
            }
        }
    }

}
