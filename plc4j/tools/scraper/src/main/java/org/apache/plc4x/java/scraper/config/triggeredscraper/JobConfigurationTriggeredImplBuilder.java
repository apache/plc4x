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

import org.apache.plc4x.java.scraper.exception.ScraperConfigurationException;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class JobConfigurationTriggeredImplBuilder {

    private final ScraperConfigurationTriggeredImplBuilder parent;
    private final String name;
    private final String triggerConfig;

    private final List<String> sources = new ArrayList<>();
    private final Map<String, String> tags = new HashMap<>();

    public JobConfigurationTriggeredImplBuilder(ScraperConfigurationTriggeredImplBuilder parent, String name, String triggerConfig) {
        if(parent==null){
            throw new ScraperConfigurationException("parent builder cannot be null");
        }
        if (name == null || name.isEmpty()) {
            throw new ScraperConfigurationException("Job name must not be null or empty");
        }
        this.parent = parent;
        this.name = name;
        this.triggerConfig = triggerConfig;
    }

    public JobConfigurationTriggeredImplBuilder source(String alias) {
        if(alias==null || alias.isEmpty()){
            throw new ScraperConfigurationException("source alias cannot be null or empty");
        }
        this.sources.add(alias);
        return this;
    }

    public JobConfigurationTriggeredImplBuilder tag(String alias, String tagQuery) {
        this.tags.put(alias, tagQuery);
        return this;
    }

    private JobConfigurationTriggeredImpl buildInternal() {
        return new JobConfigurationTriggeredImpl(name, triggerConfig, null, sources, tags);
    }

    public ScraperConfigurationTriggeredImplBuilder build() {
        parent.addJobConfiguration(this.buildInternal());
        return this.parent;
    }
}
