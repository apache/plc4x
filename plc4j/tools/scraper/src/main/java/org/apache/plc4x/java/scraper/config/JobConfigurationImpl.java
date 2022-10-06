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
package org.apache.plc4x.java.scraper.config;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;

import java.util.List;
import java.util.Map;

/**
 * abstract configuration for scrape-job configuration
 */
public class JobConfigurationImpl implements JobConfiguration {
    protected final String name;
    protected final String triggerConfig;
    protected final Integer scrapeRate;
    protected final List<String> sources;
    protected final Map<String, String> fields;

    /**
     * Default constructor
     * @param name Job Name / identifier
     * @param triggerConfig configuration string for triggered jobs
     * @param scrapeRate    rate in which the data should be acquired
     * @param sources source alias (<b>not</b> connection string but the alias (from @{@link ScraperConfigurationClassicImpl}).
     * @param fields Map from field alias (how it is named in the result map) to plc4x field query
     */
    @JsonCreator
    public JobConfigurationImpl(@JsonProperty(value = "name", required = true) String name,
                                @JsonProperty(value = "triggerConfig") String triggerConfig,
                                @JsonProperty(value = "scrapeRate") Integer scrapeRate,
                                @JsonProperty(value = "sources", required = true) List<String> sources,
                                @JsonProperty(value = "fields", required = true) Map<String, String> fields) {
        this.name = name;
        this.triggerConfig = triggerConfig;
        this.scrapeRate = scrapeRate;
        this.sources = sources;
        this.fields = fields;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public String getTriggerConfig() {
        return triggerConfig;
    }

    @Override
    public List<String> getSources() {
        return sources;
    }

    @Override
    public Map<String, String> getFields() {
        return fields;
    }

    @Override
    public Integer getScrapeRate() {
        return scrapeRate;
    }
}
