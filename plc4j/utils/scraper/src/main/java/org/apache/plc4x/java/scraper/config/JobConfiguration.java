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

package org.apache.plc4x.java.scraper.config;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;
import java.util.Map;

/**
 * Configuration for one {@link org.apache.plc4x.java.scraper.ScrapeJob} in the @{@link ScraperConfiguration}.
 */
public class JobConfiguration {

    private final String name;
    private final int scrapeRate;
    private final List<String> sources;
    private final Map<String, String> fields;

    /**
     * Default construcotr
     * @param name Job Name / identifier
     * @param scrapeRate Scrape rate in ms
     * @param sources source alias (<b>not</b> connection string but the alias (from @{@link ScraperConfiguration}).
     * @param fields Map from field alias (how it is named in the result map) to plc4x field query
     */
    @JsonCreator
    public JobConfiguration(@JsonProperty(value = "name", required = true) String name,
                            @JsonProperty(value = "scrapeRate", required = true) int scrapeRate,
                            @JsonProperty(value = "sources", required = true) List<String> sources,
                            @JsonProperty(value = "fields", required = true) Map<String, String> fields) {
        this.name = name;
        this.scrapeRate = scrapeRate;
        this.sources = sources;
        this.fields = fields;
    }

    public String getName() {
        return name;
    }

    public int getScrapeRate() {
        return scrapeRate;
    }

    public List<String> getSources() {
        return sources;
    }

    public Map<String, String> getFields() {
        return fields;
    }
}
