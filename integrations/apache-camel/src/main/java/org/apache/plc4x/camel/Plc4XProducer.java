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
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;

public class Plc4XProducer extends DefaultAsyncProducer {
    private PlcConnection plcConnection;
    private AtomicInteger openRequests;

    public Plc4XProducer(Plc4XEndpoint endpoint) throws PlcException {
        super(endpoint);
        String plc4xURI = endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", "");
        plcConnection = endpoint.getPlcDriverManager().getConnection(plc4xURI);
        if (!plcConnection.writeRequestBuilder().isPresent()) {
            throw new PlcException("This connection (" + plc4xURI + ") doesn't support writing.");
        }
        openRequests = new AtomicInteger();
    }

    @Override
    public void process(Exchange exchange) throws Exception {
        Message in = exchange.getIn();
        String fieldName = in.getHeader(Constants.FIELD_NAME_HEADER, String.class);
        String fieldQuery = in.getHeader(Constants.FIELD_QUERY_HEADER, String.class);
        Object body = in.getBody();
        if (body instanceof List) {
            List<?> bodyList = in.getBody(List.class);
            Object[] values = bodyList.toArray();
//            builder.addItem(fieldName, fieldQuery, values);
        } else {
            Object value = in.getBody(Object.class);
//            builder.addItem(fieldName, fieldQuery, value);
        }
        PlcWriteRequest.Builder builder = plcConnection.writeRequestBuilder().orElseThrow(() -> new IllegalArgumentException("Writer for driver not found"));
        CompletableFuture<? extends PlcWriteResponse> completableFuture = builder.build().execute();
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
        } catch (Exception e) {
            log.warn("Could not close {}", plcConnection, e);
        }
        super.doStop();
    }

}
