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
package org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler;

import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperTask;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

/**
 * holds the handler for the regarding trigger-scraper on rising-trigger edge
 */
public class TriggerHandlerImpl implements TriggerHandler {
    private static final Logger LOGGER = LoggerFactory.getLogger(TriggerHandlerImpl.class);

    private final TriggerConfiguration triggerConfiguration;

    private final TriggerCollector triggerCollector;

    //used to enable trigger only on rising edge
    private boolean lastTriggerState;

    public TriggerHandlerImpl(String triggerStrategy, TriggeredScrapeJobImpl triggeredScrapeJob, TriggeredScraperTask parentScraperTask, TriggerCollector triggerCollector) throws ScraperException {
        this.triggerConfiguration = TriggerConfiguration.createConfiguration(triggerStrategy,triggeredScrapeJob);

        //transmit needed trigger to triggerCollection
        for(TriggerConfiguration.TriggerElement triggerElement:triggerConfiguration.getTriggerElementList()){
            triggerElement.setPlcConnectionString(parentScraperTask.getConnectionString());
            triggerElement.setUuid(triggerCollector.submitTrigger(triggerElement.getPlcFieldString(),parentScraperTask.getConnectionString(),this.triggerConfiguration.getScrapeInterval()));
        }

        this.lastTriggerState = false;
        this.triggerCollector = triggerCollector;
    }

    /**
     * checks rising edge of trigger event
     * @return true on detection of rising edge, false otherwise
     */
    @Override
    public boolean checkTrigger(){
        switch (this.triggerConfiguration.getTriggerType()){
            case SCHEDULED:
                //used base scheduling -> trigger is always true
                return true;
            case S7_TRIGGER_VAR:
                return checkS7TriggerVariable();
            case TRIGGER_VAR:
                return checkGenericTrigger();
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

        List<Object> acquiredValuesList = new ArrayList<>();
        for(TriggerConfiguration.TriggerElement triggerElement:triggerConfiguration.getTriggerElementList()){
            try {
                Object result = triggerCollector.requestResult(triggerElement.getUuid());
                if(result==null){
                    return false;
                }
                acquiredValuesList.add(result);
            } catch (ScraperException e) {
                LOGGER.warn("Went wrong",e);
            }
        }

        //check if trigger condition from TriggerConfiguration is fulfilled
        boolean trigger = false;
        try {
            trigger = triggerConfiguration.evaluateTrigger(acquiredValuesList);
        } catch (ScraperException e) {
            LOGGER.warn("Could not evaluate trigger");
        }

        //only trigger scraping of data on rising edge of trigger
        if(trigger && !this.lastTriggerState){
            this.lastTriggerState = true;
            return true;
        }
        else{
            this.lastTriggerState = trigger;
            return false;
        }

    }

    private boolean checkGenericTrigger(){
        List<Object> acquiredValuesList = new ArrayList<>();
        for(TriggerConfiguration.TriggerElement triggerElement:triggerConfiguration.getTriggerElementList()){
            try {
                Object result = triggerCollector.requestResult(triggerElement.getUuid());
                if(result==null){
                    return false;
                }
                acquiredValuesList.add(result);
            } catch (ScraperException e) {
                LOGGER.warn("Went wrong",e);
            }
        }

        //check if trigger condition from TriggerConfiguration is fulfilled
        boolean trigger = false;
        try {
            trigger = triggerConfiguration.evaluateTrigger(acquiredValuesList);
        } catch (ScraperException e) {
            LOGGER.warn("Could not evaluate trigger");
        }

        //only trigger scraping of data on rising edge of trigger
        if(trigger && !this.lastTriggerState){
            this.lastTriggerState = true;
            return true;
        }
        else{
            this.lastTriggerState = trigger;
            return false;
        }

    }

}
