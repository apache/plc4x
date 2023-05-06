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
package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector;

import org.apache.commons.lang3.concurrent.BasicThreadFactory;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.apache.plc4x.java.utils.cache.CachedPlcConnectionManager;
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
    private static final int FUTURE_TIMEOUT                     = 2000;
    private static final int READ_REQUEST_TIMEOUT               = 2000;

    private final PlcConnectionManager plcConnectionManager;
    private final Map<String,RequestElement> currentRequestElements;
    private long schedulerInterval;
    private final long futureTimeout;

    private final ScheduledExecutorService scheduledExecutorService;
    private final ExecutorService executorService;

    public TriggerCollectorImpl(PlcConnectionManager plcConnectionManager, long schedulerInterval, long futureTimeout, int poolSizeScheduler, int poolSizeExecutor) {
        if (!(plcConnectionManager instanceof CachedPlcConnectionManager)) {
            logger.warn("The Triggered Scraper is intended to be used with a Cached Connection-Manager. In other situations leaks could occur!");
        }
        this.plcConnectionManager = plcConnectionManager;
        this.currentRequestElements = new ConcurrentHashMap<>();
        this.schedulerInterval = schedulerInterval;
        this.futureTimeout = futureTimeout;

        this.scheduledExecutorService = Executors.newScheduledThreadPool(poolSizeScheduler,
            new BasicThreadFactory.Builder()
                .namingPattern("triggercollector-scheduledExecutorService-thread-%d")
                .daemon(false)
                .build()
        );
        this.executorService = Executors.newFixedThreadPool(poolSizeExecutor,
            new BasicThreadFactory.Builder()
                .namingPattern("triggercollector-executerService-thread-%d")
                .daemon(true)
                .build()
        );

    }

    public TriggerCollectorImpl(PlcConnectionManager plcConnectionManager, long schedulerInterval, long futureTimeout) {
        this(plcConnectionManager,schedulerInterval,futureTimeout,10,20);
    }

    public TriggerCollectorImpl(PlcConnectionManager plcConnectionManager) {
        this(plcConnectionManager,DEFAULT_SCHEDULED_TRIGGER_INTERVAL, FUTURE_TIMEOUT);
    }

    /**
     * submits a trigger request to TriggerCollector
     *
     * @param tag              a (plc) tag that is used for triggering procedure
     * @param plcConnectionString   the connection string to the regarding source
     * @param interval              max awaiting time until request shall be submitted
     * @return a uuid under that the request is handled internally
     */
    @Override
    public String submitTrigger(String tag, String plcConnectionString, long interval) throws ScraperException {
        String uuid = UUID.randomUUID().toString();

        if(this.schedulerInterval>interval){
            this.schedulerInterval=interval;
        }

        RequestElement requestElement = new RequestElement(plcConnectionString, tag,interval, uuid);
        if(!currentRequestElements.containsValue(requestElement)){
            currentRequestElements.put(uuid,requestElement);
            if(logger.isDebugEnabled()) {
                logger.debug("Received request to: {} for PLC: {}", tag, plcConnectionString);
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
            throw new ScraperException(String.format("Could not evaluate UUID for given trigger (%s,%s). Should not happen please report!", tag,plcConnectionString));
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
                        plcConnection = TriggeredScraperImpl.getPlcConnection(plcConnectionManager,plcConnectionString,executorService,futureTimeout,info);
                        plcConnectionList.add(plcConnection);
                        plcReadRequestBuilderMap.put(plcConnectionString,plcConnection.readRequestBuilder());
                        plcReadRequestBuilderMap.get(plcConnectionString).addTagAddress(entry.getKey(),entry.getValue().getPlcTag());
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
                    plcReadRequestBuilderMap.get(plcConnectionString).addTagAddress(entry.getKey(),entry.getValue().getPlcTag());
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
    public Object requestResult(String uuid, long timeout){
        return currentRequestElements.get(uuid).getResult();
    }

    /**
     * starts the acquirement of triggers
     */
    @Override
    public void start() {
        this.scheduledExecutorService.scheduleAtFixedRate(this::processActiveTrigger, 1_000, this.schedulerInterval, TimeUnit.MILLISECONDS);
    }

    /**
     * stops acquirement of triggers
     */
    @Override
    public void stop() {
        this.scheduledExecutorService.shutdown();
        this.executorService.shutdown();
    }


    static class RequestElement{
        private final String plcConnectionString;
        private final String plcTag;
        private LocalDateTime lastAcquirement;
        private Object result;
        private final String uuid;
        private long scanIntervalMs;


        RequestElement(String plcConnectionString, String plcTag, long scanIntervalMs, String uuid) {
            this.plcConnectionString = plcConnectionString;
            this.plcTag = plcTag;
            this.uuid = uuid;
            this.scanIntervalMs = scanIntervalMs;
            //set initial acquirement to a long time ago
            this.lastAcquirement = LocalDateTime.of(1,1,1,1,1,1);
        }

        String getPlcConnectionString() {
            return plcConnectionString;
        }

        String getPlcTag() {
            return plcTag;
        }

        public Object getResult() {
            return result;
        }

        public void setResult(Object result) {
            this.result = result;
        }

        String getUuid() {
            return uuid;
        }

        long getScanIntervalMs() {
            return scanIntervalMs;
        }

        void setScanIntervalMs(long scanIntervalMs) {
            this.scanIntervalMs = scanIntervalMs;
        }

        LocalDateTime getLastAcquirement() {
            return lastAcquirement;
        }

        void setLastAcquirement(LocalDateTime lastAcquirement) {
            this.lastAcquirement = lastAcquirement;
        }

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;
            RequestElement that = (RequestElement) o;
            return Objects.equals(plcConnectionString, that.plcConnectionString) &&
                Objects.equals(plcTag, that.plcTag);
        }

        @Override
        public int hashCode() {
            return Objects.hash(plcConnectionString, plcTag);
        }

        @Override
        public String toString() {
            return "RequestElement{" +
                "plcConnectionString='" + plcConnectionString + '\'' +
                ", plcTag='" + plcTag + '\'' +
                ", lastAcquirement=" + lastAcquirement +
                ", result=" + result +
                ", uuid='" + uuid + '\'' +
                ", scanIntervalMs=" + scanIntervalMs +
                '}';
        }
    }

}
