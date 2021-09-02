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

import org.apache.plc4x.java.scraper.config.triggeredscraper.ScraperConfigurationTriggeredImplBuilder;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * @deprecated Scraper is deprecated please use {@link ScraperConfigurationTriggeredImplBuilder} instead all functions are supplied as well see java-doc of {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl}
 */
@Deprecated
public class ScraperConfigurationClassicImplBuilder {

    private final Map<String, String> sources = new HashMap<>();
    private final List<JobConfigurationImpl> jobConfigurations = new ArrayList<>();

    public ScraperConfigurationClassicImplBuilder addSource(String alias, String connectionString) {
        sources.put(alias, connectionString);
        return this;
    }

    public JobConfigurationClassicImplBuilder job(String name, int scrapeRateMs) {
        return new JobConfigurationClassicImplBuilder(this, name, scrapeRateMs);
    }

    public ScraperConfiguration build() {
        return new ScraperConfigurationClassicImpl(sources, jobConfigurations);
    }

    public void addJobConfiguration(JobConfigurationClassicImpl configuration) {
        this.jobConfigurations.add(configuration);
    }
}
