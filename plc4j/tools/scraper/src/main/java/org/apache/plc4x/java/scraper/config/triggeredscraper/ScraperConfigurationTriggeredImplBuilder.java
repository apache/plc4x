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
package org.apache.plc4x.java.scraper.config.triggeredscraper;

import org.apache.plc4x.java.scraper.config.JobConfigurationImpl;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ScraperConfigurationTriggeredImplBuilder {

    private final Map<String, String> sources = new HashMap<>();
    private final List<JobConfigurationImpl> jobConfigurations = new ArrayList<>();

    public ScraperConfigurationTriggeredImplBuilder addSource(String alias, String connectionString) {
        sources.put(alias, connectionString);
        return this;
    }

    public JobConfigurationTriggeredImplBuilder job(String name, String triggerConfig) {
        return new JobConfigurationTriggeredImplBuilder(this, name, triggerConfig);
    }

    public ScraperConfigurationTriggeredImpl build() {
        return new ScraperConfigurationTriggeredImpl(sources, jobConfigurations);
    }

    public void addJobConfiguration(JobConfigurationTriggeredImpl configuration) {
        this.jobConfigurations.add(configuration);
    }
}
