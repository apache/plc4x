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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class JobConfigurationBuilder {

    private final ScraperConfigurationBuilder parent;
    private final String name;
    private final int scrapeRateMs;

    private final List<String> sources = new ArrayList<>();
    private final Map<String, String> fields = new HashMap<>();

    public JobConfigurationBuilder(ScraperConfigurationBuilder parent, String name, int scrapeRateMs) {
        this.parent = parent;
        this.name = name;
        this.scrapeRateMs = scrapeRateMs;
    }

    public JobConfigurationBuilder source(String alias) {
        this.sources.add(alias);
        return this;
    }

    public JobConfigurationBuilder field(String alias, String fieldQuery) {
        this.fields.put(alias, fieldQuery);
        return this;
    }

    private JobConfiguration buildInternal() {
        return new JobConfiguration(name, scrapeRateMs, sources, fields);
    }

    public ScraperConfigurationBuilder build() {
        parent.addJobConfiguration(this.buildInternal());
        return this.parent;
    }
}
