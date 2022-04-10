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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import org.apache.plc4x.java.scraper.ScrapeJob;
import org.apache.plc4x.java.scraper.ScrapeJobImpl;
import org.apache.plc4x.java.scraper.config.JobConfigurationImpl;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.config.JobConfiguration;
import org.apache.plc4x.java.scraper.exception.ScraperConfigurationException;
import org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScrapeJobImpl;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.function.Function;
import java.util.stream.Collectors;

/**
 * Configuration class for {@link org.apache.plc4x.java.scraper.triggeredscraper.TriggeredScraperImpl}.
 */
public class ScraperConfigurationTriggeredImpl implements ScraperConfiguration {
    private static final Logger logger = LoggerFactory.getLogger( ScraperConfigurationTriggeredImpl.class );

    private final Map<String, String> sources;
    private final List<JobConfigurationImpl> jobConfigurations;

    /**
     * Default constructor.
     *
     * @param sources           Map from connection alias to connection string
     * @param jobConfigurations List of configurations one for each Job
     */
    @JsonCreator
    public ScraperConfigurationTriggeredImpl(@JsonProperty(value = "sources", required = true) Map<String, String> sources,
                                             @JsonProperty(value = "jobs", required = true) List<JobConfigurationImpl> jobConfigurations) {
        checkNoUnreferencedSources(sources, jobConfigurations);
        this.sources = sources;
        this.jobConfigurations = jobConfigurations;
    }

    private void checkNoUnreferencedSources(Map<String, String> sources, List<JobConfigurationImpl> jobConfigurations) {
        Set<String> unreferencedSources = jobConfigurations.stream()
            .flatMap(job -> job.getSources().stream())
            .filter(source -> !sources.containsKey(source))
            .collect(Collectors.toSet());
        if (!unreferencedSources.isEmpty()) {
            throw new ScraperConfigurationException("There are the following unreferenced sources: " + unreferencedSources);
        }
    }

    @Override
    public Map<String, String> getSources() {
        return sources;
    }

    @Override
    public List<JobConfigurationImpl> getJobConfigurations() {
        return jobConfigurations;
    }

    @Override
    public List<ScrapeJob> getJobs() throws ScraperException {
        return getJobs(jobConfigurations,sources);
    }

    public static List<ScrapeJob> getJobs(List<JobConfigurationImpl> jobConfigurations, Map<String, String> sources) throws ScraperConfigurationException {
        List<ScrapeJob> scrapeJobs = new ArrayList<>();
        for(JobConfiguration jobConfiguration:jobConfigurations){
            if(jobConfiguration.getTriggerConfig()!=null){
                logger.info("Assuming job as triggered job because triggerConfig has been set");
                scrapeJobs.add(new TriggeredScrapeJobImpl(jobConfiguration.getName(),
                    jobConfiguration.getTriggerConfig(),
                    getSourcesForAliases(jobConfiguration.getSources(),sources),
                    jobConfiguration.getFields()));
            }
            else {
                if(jobConfiguration.getScrapeRate()!=null){
                    logger.info("Assuming job as classic job because triggerConfig has NOT been set but scrapeRate has.");
                    scrapeJobs.add(new ScrapeJobImpl(
                        jobConfiguration.getName(),
                        jobConfiguration.getScrapeRate(),
                        getSourcesForAliases(jobConfiguration.getSources(),sources),
                        jobConfiguration.getFields()));
                }
                else {
                    logger.info("Job has lack of trigger/scheduled config");
                    throw new ScraperConfigurationException(
                        String.format("Job %s was intended to be o triggered annotation, but no triggerConfig-Field could be found. Canceling!",jobConfiguration.getName()));
                }
            }
        }
        return scrapeJobs;
    }

    private static Map<String, String> getSourcesForAliases(List<String> aliases, Map<String, String> sources) {
        return aliases.stream()
            .collect(Collectors.toMap(
                Function.identity(),
                sources::get
            ));
    }
}
