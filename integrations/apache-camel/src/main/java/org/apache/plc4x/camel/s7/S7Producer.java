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
package org.apache.plc4x.camel.s7;

import org.apache.camel.AsyncCallback;
import org.apache.camel.Exchange;
import org.apache.camel.Message;
import org.apache.camel.ShutdownRunningTask;
import org.apache.camel.impl.DefaultAsyncProducer;
import org.apache.camel.spi.ShutdownAware;
import org.apache.plc4x.camel.util.StreamUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.GenericPlcWriteRequest;
import org.apache.plc4x.java.api.messages.GenericPlcWriteResponse;
import org.apache.plc4x.java.api.model.Address;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

/**
 * The Awesome producer.
 */
public class S7Producer extends DefaultAsyncProducer implements ShutdownAware {
    private static final Logger LOG = LoggerFactory.getLogger(S7Producer.class);
    private S7Endpoint endpoint;
    private PlcConnection plcConnection;

    public S7Producer(S7Endpoint endpoint) {
        super(endpoint);
        this.endpoint = endpoint;
        try {
            plcConnection = new PlcDriverManager().getConnection(endpoint.getEndpointUri());
            plcConnection.connect();
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
    }

    @SuppressWarnings("unchecked")
    public void process(Exchange exchange) throws Exception {
        Message in = exchange.getIn();
        Address address = in.getHeader(Constants.ADDRESS_HEADER, Address.class);
        Class<?> datatype = in.getHeader(Constants.DATATYPE_HEADER, Class.class);
        Object value = in.getBody(Object.class);
        GenericPlcWriteRequest plcSimpleWriteRequest = new GenericPlcWriteRequest(datatype, address, value);
        StreamUtils.streamOf(plcConnection.getWriter())
            .map(plcWriter -> plcWriter.write(plcSimpleWriteRequest))
            .forEach(plcWriteResponseCompletableFuture -> {
                try {
                    // FIXME: If I omit the cast to CompletableFuture the java compiler complains
                    GenericPlcWriteResponse response = (GenericPlcWriteResponse)
                        ((CompletableFuture) plcWriteResponseCompletableFuture).get();
                    in.setHeader(Constants.DATATYPE_HEADER, datatype);
                    in.setHeader(Constants.ADDRESS_HEADER, address);
                    in.setBody(response);
                } catch (InterruptedException | ExecutionException e) {
                    throw new RuntimeException(e);
                }
            });
    }

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
    public boolean deferShutdown(ShutdownRunningTask shutdownRunningTask) {
        switch (shutdownRunningTask) {
            case CompleteCurrentTaskOnly:
                break;
            case CompleteAllTasks:
                break;
        }
        try {
            plcConnection.close();
        } catch (Exception ignore) {
        }
        return false;
    }

    @Override
    public int getPendingExchangesSize() {
        return 0;
    }

    @Override
    public void prepareShutdown(boolean suspendOnly, boolean forced) {

    }
}
