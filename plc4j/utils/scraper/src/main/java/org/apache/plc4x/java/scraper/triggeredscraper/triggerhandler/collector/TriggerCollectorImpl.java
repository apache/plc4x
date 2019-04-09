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
package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector;

import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.scraper.Scraper;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.concurrent.*;

/**
 * default implementation for TriggerCollector
 */
public class TriggerCollectorImpl implements TriggerCollector {
    private static final Logger logger = LoggerFactory.getLogger( TriggerCollectorImpl.class );

    private static final int DEFAULT_SCHEDULED_TRIGGER_INTERVAL = 1000;
    private static final int FUTURE_TIMEOUT = 2000;
    private static final int READ_REQUEST_TIMEOUT = 2000;
    private final PlcDriverManager plcDriverManager;
    private final Map<String,RequestElement> currentRequestElements;
    private long schedulerInterval;
    private final long futureTimeout;

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(3,
        new BasicThreadFactory.Builder()
            .namingPattern("triggercollector-scheduler-thread-%d")
            .daemon(false)
            .build()
    );
    private final ExecutorService executorService = Executors.newFixedThreadPool(10,
        new BasicThreadFactory.Builder()
            .namingPattern("triggercollector-executer-thread-%d")
            .daemon(true)
            .build()
    );


    public TriggerCollectorImpl(PlcDriverManager plcDriverManager, long schedulerInterval, long futureTimeout) {
        this.plcDriverManager = plcDriverManager;
        this.currentRequestElements = new ConcurrentHashMap<>();
        this.schedulerInterval = schedulerInterval;
        this.futureTimeout = futureTimeout;
    }

    public TriggerCollectorImpl(PlcDriverManager plcDriverManager) {
        this(plcDriverManager,DEFAULT_SCHEDULED_TRIGGER_INTERVAL, FUTURE_TIMEOUT);
    }

    /**
     * submits a trigger request to TriggerCollector
     *
     * @param plcField              a (plc) field that is used for triggering procedure
     * @param plcConnectionString   the connection string to the regarding source
     * @param interval              max awaiting time until request shall be submitted
     * @return a uuid under that the request is handled internally
     */
    @Override
    public String submitTrigger(String plcField, String plcConnectionString, long interval) throws ScraperException {
        String uuid = UUID.randomUUID().toString();

        if(this.schedulerInterval>interval){
            this.schedulerInterval=interval;
        }

        RequestElement requestElement = new RequestElement(plcConnectionString,plcField,interval, uuid);
        if(!currentRequestElements.containsValue(requestElement)){
            currentRequestElements.put(uuid,requestElement);
            if(logger.isInfoEnabled()) {
                logger.info("Received request to: {} for PLC: {}", plcField, plcConnectionString);
            }
            return uuid;
        }
        else{
            if(logger.isTraceEnabled()) {
                logger.trace("Received a placed trigger");
            }
            for(RequestElement requestElementFromMap:currentRequestElements.values()){
                if(requestElementFromMap.equals(requestElement)){
                    //detect shortest interval if trigger used more than once
                    if(requestElementFromMap.getScanIntervalMs()>interval){
                        requestElementFromMap.setScanIntervalMs(interval);
                    }
                    return requestElementFromMap.getUuid();
                }
            }

            //should not happen
            throw new ScraperException(String.format("Could not evaluate UUID for given trigger (%s,%s). Should not happen please report!",plcField,plcConnectionString));
        }

    }

    /**
     * acquire all triggers within given interval from definition
     */
    private void processActiveTrigger(){
        LocalDateTime currentTimestamp = LocalDateTime.now();
        Map<String, PlcReadRequest.Builder> plcReadRequestBuilderMap = new HashMap<>();
        Map<String, PlcReadResponse> plcReadResponseMap = new HashMap<>();
        List<RequestElement> activeRequestElements = new ArrayList<>();
        List<PlcConnection> plcConnectionList = new ArrayList<>();
        PlcConnection plcConnection=null;
        for(Map.Entry<String,RequestElement> entry:currentRequestElements.entrySet()){
            if(entry.getValue().getLastAcquirement().isBefore(
                currentTimestamp
                    .minus(entry.getValue().scanIntervalMs,ChronoUnit.MILLIS))
            ){
                String plcConnectionString = entry.getValue().plcConnectionString;
                if(!plcReadRequestBuilderMap.containsKey(plcConnectionString)){
                    try {
                        String info = "";
                        if(logger.isTraceEnabled()) {
                            info = String.format("acquiring trigger connection to (%s)", plcConnectionString);
                            logger.trace("acquiring trigger connection to ({})", plcConnectionString);
                        }
                        plcConnection = TriggeredScraperImpl.getPlcConnection(plcDriverManager,plcConnectionString,executorService,futureTimeout,info);
                        plcConnectionList.add(plcConnection);
                        plcReadRequestBuilderMap.put(plcConnectionString,plcConnection.readRequestBuilder());
                        plcReadRequestBuilderMap.get(plcConnectionString).addItem(entry.getKey(),entry.getValue().getPlcField());
                        activeRequestElements.add(entry.getValue());
                    } catch (InterruptedException e) {
                        logger.warn("Acquirement of PLC-Connection was interrupted",e);
                        Thread.currentThread().interrupt();
                    } catch (ExecutionException e) {
                        logger.warn("Acquirement of PLC-Connection could not be executed",e);
                    } catch (TimeoutException e) {
                        logger.warn("Acquirement of PLC-Connection was timeouted",e);
                    }
                }
                else{
                    plcReadRequestBuilderMap.get(plcConnectionString).addItem(entry.getKey(),entry.getValue().getPlcField());
                    activeRequestElements.add(entry.getValue());
                }
            }
        }

        for(Map.Entry<String,PlcReadRequest.Builder> entry: plcReadRequestBuilderMap.entrySet()){
            try {
                PlcReadResponse plcReadResponse = entry.getValue().build().execute().get(futureTimeout, TimeUnit.MILLISECONDS);
                plcReadResponseMap.put(entry.getKey(), plcReadResponse);
            } catch (InterruptedException e) {
                logger.warn("Extraction of PlcResponse was interrupted",e);
                Thread.currentThread().interrupt();
            } catch (ExecutionException e) {
                logger.warn("Extraction of PlcResponse could not be executed",e);
            } catch (TimeoutException e) {
                logger.warn("Extraction of PlcResponse was timeouted",e);
            }
        }

        LocalDateTime currentTime = LocalDateTime.now();
        for(RequestElement requestElement:activeRequestElements){
            requestElement.setResult(plcReadResponseMap.get(requestElement.getPlcConnectionString()).getObject(requestElement.getUuid()));
            requestElement.setLastAcquirement(currentTime);
        }
        for(PlcConnection plcConnectionFromList:plcConnectionList){
            if(plcConnectionFromList!=null){
                try {
                    plcConnectionFromList.close();
                } catch (Exception e) {
                    logger.warn("Could not close connection ...");
                }
            }
        }

    }

    /**
     * requests the result of submitted plc request with default timeout
     *
     * @param uuid uuid that represents the request
     * @return the object acquired by requesting plc instance
     */
    @Override
    public Object requestResult(String uuid) throws ScraperException {
        return requestResult(uuid, READ_REQUEST_TIMEOUT);
    }

    /**
     * requests the result of submitted plc request
     *
     * @param uuid uuid that represents the request
     * @return the object acquired by requesting plc instance
     */
    @Override
    public Object requestResult(String uuid, long timeout) throws ScraperException {
        Object result = currentRequestElements.get(uuid).getResult();
        /*
        if(result==null){
            throw new ScraperException("No result acquired yet");
        }
        */
        return result;
    }

    /**
     * starts the acquirement of triggers
     */
    @Override
    public void start() {
        this.scheduler.scheduleAtFixedRate(() -> processActiveTrigger(), 1_000, this.schedulerInterval, TimeUnit.MILLISECONDS);
    }

    /**
     * stops acquirement of triggers
     */
    @Override
    public void stop() {
        //ToDo stop everything the right way
    }

    class RequestElement{
        private String plcConnectionString;
        private String plcField;
        private LocalDateTime submitTimeOut;
        private LocalDateTime lastAcquirement;
        private Object result;
        private String uuid;
        private boolean submitted;
        private CompletableFuture<PlcReadResponse> plcReadResponse;
        private long scanIntervalMs;


        public RequestElement(String plcConnectionString, String plcField, long scanIntervalMs, String uuid) {
            this.plcConnectionString = plcConnectionString;
            this.plcField = plcField;
            this.submitted = false;
            this.submitTimeOut = LocalDateTime.now().plus(scanIntervalMs, ChronoUnit.MILLIS);
            this.uuid = uuid;
            this.scanIntervalMs = scanIntervalMs;
            //set initial acquirement to a long time ago
            this.lastAcquirement = LocalDateTime.of(1,1,1,1,1,1);
        }

        public String getPlcConnectionString() {
            return plcConnectionString;
        }

        public String getPlcField() {
            return plcField;
        }

        public LocalDateTime getSubmitTimeOut() {
            return submitTimeOut;
        }

        public boolean isSubmitted() {
            return submitted;
        }

        public void setSubmitted(boolean submitted) {
            this.submitted = submitted;
        }

        public Object getResult() {
            return result;
        }

        public void setResult(Object result) {
            this.result = result;
        }

        public String getUuid() {
            return uuid;
        }

        public CompletableFuture<PlcReadResponse> getPlcReadResponse() {
            return plcReadResponse;
        }

        public void setPlcReadResponse(CompletableFuture<PlcReadResponse> plcReadResponse) {
            this.plcReadResponse = plcReadResponse;
        }

        public long getScanIntervalMs() {
            return scanIntervalMs;
        }

        public void setScanIntervalMs(long scanIntervalMs) {
            this.scanIntervalMs = scanIntervalMs;
        }

        public LocalDateTime getLastAcquirement() {
            return lastAcquirement;
        }

        public void setLastAcquirement(LocalDateTime lastAcquirement) {
            this.lastAcquirement = lastAcquirement;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            RequestElement that = (RequestElement) o;
            return Objects.equals(plcConnectionString, that.plcConnectionString) &&
                Objects.equals(plcField, that.plcField);
        }

        @Override
        public int hashCode() {
            return Objects.hash(plcConnectionString, plcField);
        }
    }

}
