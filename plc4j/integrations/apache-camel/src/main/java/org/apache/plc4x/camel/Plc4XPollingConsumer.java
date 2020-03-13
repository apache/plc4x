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

import org.apache.camel.Endpoint;
import org.apache.camel.Exchange;
import org.apache.camel.PollingConsumer;
import org.apache.camel.Processor;
import org.apache.camel.spi.ExceptionHandler;
import org.apache.camel.support.LoggingExceptionHandler;
import org.apache.camel.support.service.ServiceSupport;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

public class Plc4XPollingConsumer extends ServiceSupport implements PollingConsumer {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XPollingConsumer.class);

    private Plc4XEndpoint endpoint;
    private ExceptionHandler exceptionHandler;
    private PlcConnection plcConnection;
    private PlcReadRequest.Builder requestBuilder;
    private Class dataType;

    //private int request =0;

    public Plc4XPollingConsumer(Plc4XEndpoint endpoint) throws PlcException {
        this.endpoint = endpoint;
        this.dataType = endpoint.getDataType();
        this.exceptionHandler = new LoggingExceptionHandler(endpoint.getCamelContext(), getClass());
        String plc4xURI = endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", "");
        this.plcConnection = endpoint.getConnection();
    }

    @Override
    public String toString() {
        return "Plc4XConsumer[" + endpoint + "]";
    }

    @Override
    public Endpoint getEndpoint() {
        return endpoint;
    }

    public ExceptionHandler getExceptionHandler() {
        return exceptionHandler;
    }

    public void setExceptionHandler(ExceptionHandler exceptionHandler) {
        this.exceptionHandler = exceptionHandler;
    }

    @Override
    public Exchange receive() {
        Exchange exchange = endpoint.createExchange();
        try {
            PlcReadResponse read = createReadRequest().execute().get();
            if(endpoint.getAddress().size()==1) {
                exchange.getIn().setBody(unwrapIfSingle(read.getAllObjects("default")));
            }
            else{
                List<Object> values = new ArrayList<>();
                for(String field : read.getFieldNames()){
                    values.add(read.getObject(field));
                }
                exchange.getIn().setBody(values);
            }        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            exchange.setException(e);
        } catch (ExecutionException e) {
            exchange.setException(e);
        }
        return exchange;
    }

    @Override
    public Exchange receiveNoWait() {
        return null;
    }

    @Override
    public Exchange receive(long timeout) {
        Exchange exchange = endpoint.createExchange();
        CompletableFuture<? extends PlcReadResponse> read = createReadRequest().execute();
        try {
            PlcReadResponse plcReadResponse = read.get(timeout, TimeUnit.MILLISECONDS);
            if (read.isDone()) {
                if (endpoint.getAddress().size() == 1) {
                    exchange.getIn().setBody(unwrapIfSingle(plcReadResponse.getAllObjects("default")));
                } else {
                    List<Object> values = new ArrayList<>();
                    for (String field : plcReadResponse.getFieldNames()) {
                        values.add(plcReadResponse.getObject(field));
                    }
                    exchange.getIn().setBody(values);
                }
            }
            } catch(InterruptedException e){
                Thread.currentThread().interrupt();
                exchange.setException(e);
            } catch(ExecutionException | TimeoutException e){
                exchange.setException(e);
            }
        return exchange;
    }

    @Override
    protected void doStart() {
    }

    @Override
    protected void doStop() throws Exception {
    }


    private PlcReadRequest createReadRequest() {
        requestBuilder = plcConnection.readRequestBuilder();
        int i=0;
        if (endpoint.getAddress().size()>1){
            for(String query : endpoint.getAddress()){
                requestBuilder.addItem(String.valueOf(i++),query);
            }
        }
        else{
            requestBuilder.addItem("default",endpoint.getAddress().get(0));
        }
        return requestBuilder.build();
    }

    private Object unwrapIfSingle(Collection collection) {
        if (collection.isEmpty()) {
            return null;
        }
        if (collection.size() == 1) {
            return collection.iterator().next();
        }
        return collection;
    }

    @Override
    public Processor getProcessor() {
        return null;
    }
}
