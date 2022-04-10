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
package org.apache.plc4x.java.scraper.triggeredscraper;

import org.apache.plc4x.java.scraper.ScrapeJob;
import org.apache.plc4x.java.scraper.exception.ScraperConfigurationException;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.TriggerConfiguration;

import java.util.Map;

public class TriggeredScrapeJobImpl implements ScrapeJob {
    private final String jobName;
    private final Map<String, String> sourceConnections;
    private final Map<String, String> fields;
    private final String triggerConfig;
    private final TriggerConfiguration triggerConfiguration;


    public TriggeredScrapeJobImpl(String jobName, String triggerConfig, Map<String, String> connections, Map<String, String> fields) throws ScraperConfigurationException {
        this.jobName = jobName;
        this.triggerConfig = triggerConfig;
        this.sourceConnections = connections;
        this.fields = fields;
        this.triggerConfiguration = TriggerConfiguration.createConfiguration(triggerConfig,this);
    }

    @Override
    public String getJobName() {
        return this.jobName;
    }

    @Override
    public long getScrapeRate() {
        return triggerConfiguration.getScrapeInterval();
    }

    /**
     * alias -&gt; connection-string
     */
    @Override
    public Map<String, String> getSourceConnections() {
        return this.sourceConnections;
    }

    /**
     * alias -&gt; field-query
     */
    @Override
    public Map<String, String> getFields() {
        return fields;
    }

    public String getTriggerConfig() {
        return triggerConfig;
    }
}
