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
package org.apache.plc4x.java.tools.plc4xserver.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.plc4x.readwrite.*;
import org.apache.plc4x.java.utils.connectionpool2.PooledDriverManager;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4xServerAdapter extends ChannelInboundHandlerAdapter {

    private final PooledDriverManager driverManager;
    private final AtomicInteger connectionIdGenerator;
    private final ConcurrentHashMap<Integer, String> connectionUrls;

    public Plc4xServerAdapter() {
        driverManager = new PooledDriverManager();
        connectionIdGenerator = new AtomicInteger(1);
        connectionUrls = new ConcurrentHashMap<>();
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        if (msg instanceof Plc4xMessage) {
            final Plc4xMessage plc4xMessage = (Plc4xMessage) msg;
            switch (plc4xMessage.getRequestType()) {
                case CONNECT_REQUEST: {
                    Plc4xConnectRequest request = (Plc4xConnectRequest) plc4xMessage;
                    try (final PlcConnection connection = driverManager.getConnection(request.getConnectionString())) {
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
                        // Build a read request for all fields in the requet.
                        final PlcReadRequest.Builder builder = connection.readRequestBuilder();
                        for (Plc4xFieldRequest requestField : request.getFields()) {
                            builder.addItem(requestField.getField().getName(), requestField.getField().getFieldQuery());
                        }
                        final PlcReadRequest rr = builder.build();

                        // Execute the query.
                        final CompletableFuture<? extends PlcReadResponse> execute = rr.execute();

                        execute.whenComplete((plcReadResponse, throwable) -> {
                            if(throwable == null) {
                                // Create the response.
                                List<Plc4xFieldValueResponse> fields = new ArrayList<>(plcReadResponse.getFieldNames().size());
                                for (Plc4xFieldRequest requestField : request.getFields()) {
                                    final PlcResponseCode responseCode = plcReadResponse.getResponseCode(requestField.getField().getName());
                                    Plc4xResponseCode resCode;
                                    Plc4xValueType valueType;
                                    PlcValue value;
                                    if(responseCode == PlcResponseCode.OK) {
                                        resCode = Plc4xResponseCode.OK;
                                        // TODO: Get the real type.
                                        valueType = Plc4xValueType.BOOL;
                                        value = plcReadResponse.getPlcValue(requestField.getField().getName());
                                    } else {
                                        resCode = Plc4xResponseCode.INVALID_ADDRESS;
                                        valueType = Plc4xValueType.NULL;
                                        value = null;
                                    }
                                    fields.add(new Plc4xFieldValueResponse(
                                        requestField.getField(), resCode, valueType, value));
                                }
                                Plc4xReadResponse response = new Plc4xReadResponse(
                                    request.getRequestId(), request.getConnectionId(), Plc4xResponseCode.OK, fields);

                                // Send the response.
                                ctx.writeAndFlush(response);
                            } else {
                                Plc4xReadResponse response = new Plc4xReadResponse(
                                    request.getRequestId(), request.getConnectionId(), Plc4xResponseCode.NOT_FOUND,
                                    Collections.emptyList());
                                // Send the response.
                                ctx.writeAndFlush(response);
                            }
                        });
                    } catch (Exception e) {
                        Plc4xReadResponse response = new Plc4xReadResponse(
                            request.getRequestId(), request.getConnectionId(),
                            Plc4xResponseCode.INVALID_ADDRESS, Collections.emptyList());
                        ctx.writeAndFlush(response);
                    }
                    break;
                }

                case WRITE_REQUEST:
                    break;
            }
        }
    }

}
