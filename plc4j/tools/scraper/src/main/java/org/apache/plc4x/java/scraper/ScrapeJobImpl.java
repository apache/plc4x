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

import org.apache.plc4x.java.scraper.config.JobConfigurationClassicImpl;
import org.apache.plc4x.java.scraper.config.ScraperConfigurationClassicImpl;

import java.util.Map;

/**
 * POJO Object to transport all Job information.
 * Is generated from {@link ScraperConfigurationClassicImpl} by
 * merging the sources and the {@link JobConfigurationClassicImpl}.
 *
 * @deprecated Scraper is deprecated please use {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl} instead all functions are supplied as well see java-doc of {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl}
 */
@Deprecated
public class ScrapeJobImpl implements ScrapeJob {

    private final String name;
    private final long scrapeRate;
    private final Map<String, String> connections;
    private final Map<String, String> fields;

    public ScrapeJobImpl(String name, long scrapeRate, Map<String, String> connections, Map<String, String> fields) {
        this.name = name;
        this.scrapeRate = scrapeRate;
        this.connections = connections;
        this.fields = fields;
    }

    @Override
    public String getJobName() {
        return name;
    }

    @Override
    public long getScrapeRate() {
        return scrapeRate;
    }

    @Override
    public Map<String, String> getSourceConnections() {
        return connections;
    }

    @Override
    public Map<String, String> getFields() {
        return fields;
    }
}
