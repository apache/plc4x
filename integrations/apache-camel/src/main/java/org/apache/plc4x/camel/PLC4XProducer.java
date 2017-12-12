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
import org.apache.camel.ShutdownRunningTask;
import org.apache.camel.impl.DefaultAsyncProducer;
import org.apache.camel.spi.ShutdownAware;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.model.Address;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;

public class PLC4XProducer extends DefaultAsyncProducer implements ShutdownAware {
    private static final Logger LOG = LoggerFactory.getLogger(PLC4XProducer.class);

    private PLC4XEndpoint endpoint;
    private PlcConnection plcConnection;

    public PLC4XProducer(PLC4XEndpoint endpoint) {
        super(endpoint);
        this.endpoint = endpoint;
        try {
            String plc4xURI = endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", "");
            plcConnection = endpoint.plcDriverManager.getConnection(plc4xURI);
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
        PlcWriteRequest plcSimpleWriteRequest = new PlcWriteRequest(datatype, address, value);
        PlcWriter plcWriter = plcConnection.getWriter().orElseThrow(() -> new IllegalArgumentException("Writer for driver not found"));
        Object response = plcWriter.write(plcSimpleWriteRequest).get();
        if (exchange.getPattern().isOutCapable()) {
            Message out = exchange.getOut();
            out.copyFrom(exchange.getIn());
            out.setBody(response);
        } else {
            in.setBody(response);
        }
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
