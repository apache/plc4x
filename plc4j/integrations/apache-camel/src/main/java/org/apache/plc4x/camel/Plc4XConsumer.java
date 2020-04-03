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

import org.apache.camel.*;
import org.apache.camel.spi.ExceptionHandler;
import org.apache.camel.support.AsyncProcessorConverterHelper;
import org.apache.camel.support.LoggingExceptionHandler;
import org.apache.camel.support.service.ServiceSupport;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.*;

public class Plc4XConsumer extends ServiceSupport implements Consumer {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XConsumer.class);

    private Plc4XEndpoint endpoint;
    private AsyncProcessor processor;
    private ExceptionHandler exceptionHandler;
    private PlcConnection plcConnection;
    private  List<TagData> tags;
    private Map parameters;
    private PlcSubscriptionResponse subscriptionResponse;

    private ScheduledExecutorService executorService = Executors.newSingleThreadScheduledExecutor();
    private ScheduledFuture<?> future;

    public Plc4XConsumer(Plc4XEndpoint endpoint, Processor processor) throws PlcException {
        this.endpoint = endpoint;
        this.parameters = endpoint.getParameters();
        this.processor = AsyncProcessorConverterHelper.convert(processor);
        this.exceptionHandler = new LoggingExceptionHandler(endpoint.getCamelContext(), getClass());
        this.plcConnection = endpoint.getConnection();
        this.tags = endpoint.getTags();
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
    protected void doStart() throws InterruptedException, ExecutionException {
        PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
        if (tags.size()==1){
            TagData tag = tags.get(0);
            builder.addItem(tag.getTagName(),tag.getQuery());

        }
        else{
           for(TagData tag : tags){
               builder.addItem(tag.getTagName(),tag.getQuery());
           }
        }
        PlcReadRequest request = builder.build();
        future = executorService.schedule(() -> {
            request.execute().thenAccept(response -> {
                    try {
                        Exchange exchange = endpoint.createExchange();
                        if (tags.size()>1){
                            List<TagData> values = new ArrayList<>();
                            for(TagData tag : tags){
                                tag.setValue(response.getObject(tag.getTagName()));
                                values.add(tag);
                            }
                            exchange.getIn().setBody(values);
                        }
                        else {
                            TagData tag = tags.get(0);
                            tag.setValue(response.getAllObjects(tag.getTagName()));
                            exchange.getIn().setBody(tag);
                        }
                        processor.process(exchange);
                    } catch (Exception e) {
                        exceptionHandler.handleException(e);
                    }
                });
        }, 500, TimeUnit.MILLISECONDS);
    }

    @Override
    protected void doStop() throws InterruptedException, ExecutionException, TimeoutException {
        // First stop the polling process
        if (future != null) {
            future.cancel(true);
        }
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
        return this.processor;
    }
}