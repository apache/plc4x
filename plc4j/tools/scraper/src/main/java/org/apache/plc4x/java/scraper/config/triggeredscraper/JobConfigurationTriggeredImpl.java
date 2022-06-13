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
package org.apache.plc4x.java.scraper.config.triggeredscraper;


import org.apache.plc4x.java.scraper.config.JobConfigurationImpl;
import org.apache.plc4x.java.scraper.config.ScraperConfigurationClassicImpl;

import java.util.List;
import java.util.Map;

/**
 * Configuration for one {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl} in the {@link ScraperConfigurationTriggeredImpl}.
 */
public class JobConfigurationTriggeredImpl extends JobConfigurationImpl {


    /**
     * Default constructor
     *
     * @param name          Job Name / identifier
     * @param triggerConfig configuration string for triggered jobs
     * @param scrapeRate    rate in which the data should be acquired
     * @param sources       source alias (<b>not</b> connection string but the alias (from @{@link ScraperConfigurationClassicImpl}).
     * @param fields        Map from field alias (how it is named in the result map) to plc4x field query
     */
    public JobConfigurationTriggeredImpl(String name, String triggerConfig, Integer scrapeRate, List<String> sources, Map<String, String> fields) {
        super(name, triggerConfig, scrapeRate, sources, fields);
    }
}
