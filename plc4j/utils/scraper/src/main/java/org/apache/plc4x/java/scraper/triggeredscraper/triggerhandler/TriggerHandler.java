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
package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperTask;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * holds the handler for the regarding trigger-scraper on rising-trigger edge
 */
public class TriggerHandler {
    private static final Logger LOGGER = LoggerFactory.getLogger(TriggerHandler.class);
    private static final String TRIGGER = "trigger";

    private final TriggerConfiguration triggerConfiguration;
    private final TriggeredScraperTask parentScraperTask;

    //used to enable trigger only on rising edge
    private boolean lastTriggerState;

    public TriggerHandler(String triggerStrategy, TriggeredScrapeJobImpl triggeredScrapeJob,TriggeredScraperTask parentScraperTask) throws ScraperException {
        this.triggerConfiguration = TriggerConfiguration.createConfiguration(triggerStrategy,triggeredScrapeJob);
        this.parentScraperTask = parentScraperTask;
        this.lastTriggerState = false;
    }

    /**
     * checks rising edge of trigger event
     * @return true on detection of rising edge, false otherwise
     */
    public boolean checkTrigger(){
        switch (this.triggerConfiguration.getTriggerType()){
            case SCHEDULED:
                //used base scheduling -> trigger is always true
                return true;
            case S7_TRIGGER_VAR:
                return checkS7TriggerVariable();
            default:
                //should not happen
                return false;
        }
    }

    /**
     * acquires the given S7Field from S7 and evaluates if trigger is released
     * @return true if rising-edge of trigger is detected, false otherwise
     */
    private boolean checkS7TriggerVariable(){

        CompletableFuture<PlcConnection> future = CompletableFuture.supplyAsync(() -> {
            try {
                return parentScraperTask.getDriverManager().getConnection(parentScraperTask.getConnectionString());
            } catch (PlcConnectionException e) {
                LOGGER.warn("Unable to instantiate connection to " + parentScraperTask.getConnectionString(), e);
                throw new PlcRuntimeException(e);
            }
        }, parentScraperTask.getExecutorService());
        PlcConnection connection = null;
        try {
            connection = future.get(parentScraperTask.getRequestTimeoutMs(), TimeUnit.MILLISECONDS);
            LOGGER.trace("Connection to {} established: {}", parentScraperTask.getConnectionString(), connection);
            PlcReadRequest.Builder builder = connection.readRequestBuilder();
            builder.addItem(TRIGGER, triggerConfiguration.getTriggerVariable());
            PlcReadResponse response = builder
                .build()
                .execute()
                .get(parentScraperTask.getRequestTimeoutMs(), TimeUnit.MILLISECONDS);

            //check if trigger condition from TriggerConfiguration is fulfilled
            boolean trigger = triggerConfiguration.evaluateTrigger(response.getObject(TRIGGER));

            //only trigger scraping of data on rising edge of trigger
            if(trigger && !this.lastTriggerState){
                this.lastTriggerState = true;
                return true;
            }
            else{
                this.lastTriggerState = trigger;
                return false;
            }

        } catch (Exception e) {
            // Handle execution exception
            parentScraperTask.handleException(e);
            return false;
        }
        finally {
            if (connection != null) {
                try {
                    connection.close();
                } catch (Exception e) {
                    LOGGER.warn("Error on closing connection",e);
                    // intentionally do nothing
                }
            }
        }
    }






}
