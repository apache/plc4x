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
package org.apache.plc4x.camel;

import org.apache.camel.AsyncCallback;
import org.apache.camel.Exchange;
import org.apache.camel.Message;
import org.apache.camel.impl.DefaultAsyncProducer;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.model.Address;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class PLC4XProducer extends DefaultAsyncProducer {
    @SuppressWarnings("unused")
    private PLC4XEndpoint endpoint;
    private PlcConnection plcConnection;
    private AtomicInteger openRequests;

    public PLC4XProducer(PLC4XEndpoint endpoint) {
        super(endpoint);
        this.endpoint = endpoint;
        try {
            String plc4xURI = endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", "");
            plcConnection = endpoint.plcDriverManager.getConnection(plc4xURI);
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
        openRequests = new AtomicInteger();
    }

    @SuppressWarnings("unchecked")
    @Override
    public void process(Exchange exchange) throws Exception {
        Message in = exchange.getIn();
        Address address = in.getHeader(Constants.ADDRESS_HEADER, Address.class);
        Object body = in.getBody();
        PlcWriteRequest.Builder builder = PlcWriteRequest.builder();
        if (body instanceof List) {
            List<?> bodyList = in.getBody(List.class);
            bodyList
                .stream()
                .map(o -> (WriteRequestItem<?>) new WriteRequestItem(o.getClass(), address, o))
                .forEach(builder::addItem);
        } else {
            Object value = in.getBody(Object.class);
            builder.addItem(address, value);
        }
        PlcWriter plcWriter = plcConnection.getWriter().orElseThrow(() -> new IllegalArgumentException("Writer for driver not found"));
        CompletableFuture<? extends PlcWriteResponse> completableFuture = plcWriter.write(builder.build());
        int currentlyOpenRequests = openRequests.incrementAndGet();
        try {
            log.debug("Currently open requests including {}:{}", exchange, currentlyOpenRequests);
            Object plcWriteResponse = completableFuture.get();
            if (exchange.getPattern().isOutCapable()) {
                Message out = exchange.getOut();
                out.copyFrom(exchange.getIn());
                out.setBody(plcWriteResponse);
            } else {
                in.setBody(plcWriteResponse);
            }
        } finally {
            int openRequestsAfterFinish = openRequests.decrementAndGet();
            log.trace("Open Requests after {}:{}", exchange, openRequestsAfterFinish);
        }
    }

    @Override
    public boolean process(Exchange exchange, AsyncCallback callback) {
        try {
            process(exchange);
            Message out = exchange.getOut();
            out.copyFrom(exchange.getIn());
        } catch (Exception e) {
            exchange.setOut(null);
            exchange.setException(e);
        }
        callback.done(true);
        return true;
    }

    @Override
    protected void doStop() throws Exception {
        int openRequestsAtStop = openRequests.get();
        log.debug("Stopping with {} open requests", openRequestsAtStop);
        if (openRequestsAtStop > 0) {
            log.warn("There are still {} open requests", openRequestsAtStop);
        }
        try {
            plcConnection.close();
        } catch (Exception ignore) {
        }
        super.doStop();
    }

}
