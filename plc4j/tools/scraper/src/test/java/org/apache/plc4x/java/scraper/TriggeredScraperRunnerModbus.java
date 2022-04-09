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
package org.apache.plc4x.java.scraper;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImpl;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollectorImpl;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.List;
import java.util.Map;

public class TriggeredScraperRunnerModbus {

    private static final Logger LOGGER = LoggerFactory.getLogger(TriggeredScraperRunnerModbus.class);

    /**
     * testing of TriggeredScraper vs real device (Modbus)
     */
    public static void main(String[] args) throws IOException, ScraperException {
        ScraperConfiguration configuration = ScraperConfiguration.fromFile("plc4j/utils/scraper/src/test/resources/example_triggered_scraper_modbus.yml", ScraperConfigurationTriggeredImpl.class);
        PlcDriverManager plcDriverManager = new PooledPlcDriverManager();
        TriggeredScraperImpl scraper = new TriggeredScraperImpl(
            configuration,
            plcDriverManager,
            (j, a, m) -> {
            LOGGER.info("Results from {}/{}: {}", j, a, m);
            for(Map.Entry<String, Object> entry:m.entrySet()){
                for(Object object:(List<Object>)entry.getValue()){
                    LOGGER.info("{}",object);
                }
            }
            },
            new TriggerCollectorImpl(plcDriverManager));

        scraper.start();
    }

}
