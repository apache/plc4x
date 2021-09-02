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
package org.apache.plc4x.camel;

import org.apache.camel.Endpoint;
import org.apache.camel.Exchange;
import org.apache.camel.PollingConsumer;
import org.apache.camel.Processor;
import org.apache.camel.spi.ExceptionHandler;
import org.apache.camel.support.LoggingExceptionHandler;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollectorImpl;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

public class Plc4XPollingConsumer implements PollingConsumer {
    private static final Logger LOGGER = LoggerFactory.getLogger(Plc4XPollingConsumer.class);

    private Plc4XEndpoint plc4XEndpoint;
    private ExceptionHandler exceptionHandler;
    private PlcConnection plcConnection;
    private PlcReadRequest.Builder requestBuilder;
    private  Map<String,Object> tags;
    private String trigger;

//TODO Is this still needed with the scraper working?
    public Plc4XPollingConsumer(Plc4XEndpoint endpoint) throws PlcException {
        plc4XEndpoint=endpoint;
        this.exceptionHandler = new LoggingExceptionHandler(endpoint.getCamelContext(), getClass());
        String plc4xURI = endpoint.getEndpointUri().replaceFirst("plc4x:/?/?", "");
        this.plcConnection = endpoint.getConnection();
        this.tags = endpoint.getTags();
        this.trigger= endpoint.getTrigger();
    }

    @Override
    public String toString() {
        return "Plc4XConsumer[" + plc4XEndpoint + "]";
    }

    @Override
    public Endpoint getEndpoint() {
        return plc4XEndpoint;
    }

    public ExceptionHandler getExceptionHandler() {
        return exceptionHandler;
    }

    public void setExceptionHandler(ExceptionHandler exceptionHandler) {
        this.exceptionHandler = exceptionHandler;
    }

    @Override
    public Exchange receive() {/**
        Exchange exchange = plc4XEndpoint.createExchange();
        try {
            PlcReadResponse read = createReadRequest().execute().get();
            if(plc4XEndpoint.getTags().size()==1) {
                TagData tag = plc4XEndpoint.getTags().get(0);
                tag.setValue(read.getAllObjects(tag.getTagName()));
                exchange.getIn().setBody(tag);
            }
            else{
                List<TagData> values = new ArrayList<>();
                for(TagData tag : plc4XEndpoint.getTags()){
                    tag.setValue(read.getObject(tag.getTagName()));
                    values.add(tag);
                }
                exchange.getIn().setBody(values);
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            exchange.setException(e);
        } catch (ExecutionException e) {
            exchange.setException(e);
        }
        return exchange;*/
    return null;
    }

    @Override
    public Exchange receiveNoWait() {
        return null;
    }

    @Override
    public Exchange receive(long timeout) {/**
        Exchange exchange = plc4XEndpoint.createExchange();
        CompletableFuture<? extends PlcReadResponse> read = createReadRequest().execute();
        try {
            PlcReadResponse plcReadResponse = read.get(timeout, TimeUnit.MILLISECONDS);
            if(plc4XEndpoint.getTags().size()==1) {
                TagData tag = plc4XEndpoint.getTags().get(0);
                tag.setValue(plcReadResponse.getAllObjects(tag.getTagName()));
                exchange.getIn().setBody(tag);
            }
            else{
                List<TagData> values = new ArrayList<>();
                for(TagData tag : plc4XEndpoint.getTags()){
                    tag.setValue(plcReadResponse.getObject(tag.getTagName()));
                    values.add(tag);
                }
                exchange.getIn().setBody(values);
            }
            } catch(InterruptedException e){
                Thread.currentThread().interrupt();
                exchange.setException(e);
            } catch(ExecutionException | TimeoutException e){
                exchange.setException(e);
            }
        return exchange;*/
    return null;
    }


    private PlcReadRequest createReadRequest() {/**
        requestBuilder = plcConnection.readRequestBuilder();
        if (plc4XEndpoint.getTags().size()>1){
            for(TagData tag : plc4XEndpoint.getTags()){
                requestBuilder.addItem(tag.getTagName(),tag.getQuery());
            }
        }
        else{
            TagData tag = plc4XEndpoint.getTags().get(0);
            requestBuilder.addItem(tag.getTagName(),tag.getQuery());
        }
        return requestBuilder.build();
    */return null;}

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

    @Override
    public void start() {

    }

    @Override
    public void stop() {

    }
}
