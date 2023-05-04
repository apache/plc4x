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
package org.apache.plc4x.kafka.config;

import org.apache.plc4x.kafka.Plc4xSourceConnector;
import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.Config;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigException;
import org.apache.kafka.common.config.ConfigValue;
import java.util.Map;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class SourceConfig extends AbstractConfig{

    private static final Logger log = LoggerFactory.getLogger(SourceConfig.class);

    private final String defaultTopic;
    private final List<Source> sources;
    private final List<Job> jobs;

    public SourceConfig(Map originals) {
        super(configDef(), originals);
        defaultTopic = getString(Constants.DEFAULT_TOPIC_CONFIG);

        sources = new ArrayList<>(getList(Constants.SOURCES_CONFIG).size());
        for (String source : getList(Constants.SOURCES_CONFIG)) {
            sources.add(new Source(source, defaultTopic, originalsWithPrefix(Constants.SOURCES_CONFIG + "." + source + ".")));
        }

        jobs = new ArrayList<>(getList(Constants.JOBS_CONFIG).size());
        for (String job : getList(Constants.JOBS_CONFIG)) {
            jobs.add(new Job(job, originalsWithPrefix(Constants.JOBS_CONFIG + "." + job + ".")));
        }
    }

    public void validate() throws ConfigException {
        for (Source source : sources) {
            source.validate();
            for (JobReference jobReference : source.getJobReferences()) {
                Boolean found = false;
                for (Job job : jobs) {
                    if (jobReference.getName().equals(job.getName())) {
                        found = true;
                        break;
                    }
                }
                if (!found) {
                    throw new ConfigException(
                        String.format("Couldn't find a matching job for job reference %s in source %s ", jobReference.getName(), source.getName()));
                }
            }
        }
        for (Job job : jobs) {
            job.validate();
        }
    }

    public List<Source> getSources() {
        return sources;
    }

    public Source getSource(String sourceName) {
        if(sources == null) {
            return null;
        }
        return sources.stream().filter(source -> source.getName().equals(sourceName)).findFirst().orElse(null);
    }

    public List<Job> getJobs() {
        return jobs;
    }

    public Job getJob(String jobName) {
        if(jobs == null) {
            return null;
        }
        return jobs.stream().filter(job -> job.getName().equals(jobName)).findFirst().orElse(null);
    }

    public static ConfigDef configDef() {
        return new ConfigDef()
            .define(Constants.DEFAULT_TOPIC_CONFIG,
                    ConfigDef.Type.STRING,
                    ConfigDef.Importance.LOW,
                    Constants.DEFAULT_TOPIC_DOC)
            .define(Constants.SOURCES_CONFIG,
                    ConfigDef.Type.LIST,
                    Constants.SOURCES_DEFAULT,
                    ConfigDef.Importance.LOW,
                    Constants.SOURCES_DOC)
            .define(Constants.JOBS_CONFIG,
                    ConfigDef.Type.LIST,
                    Constants.JOBS_DEFAULT,
                    ConfigDef.Importance.LOW,
                    Constants.JOBS_DOC);
    }

    @Override
    public String toString() {
        StringBuilder query = new StringBuilder();
        query.append(Constants.DEFAULT_TOPIC_CONFIG + "=" + defaultTopic + ",\n");
        for (Source source : sources) {
            query.append(source.toString());
        }
        for (Job job : jobs) {
            query.append(job.toString());
        }
        return query.toString();
    }

}
