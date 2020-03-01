/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.kafka.config;

import org.apache.plc4x.kafka.Plc4xSourceConnector;

import java.util.*;

public class SourceConfig {

    private static final String CONNECTION_STRING_CONFIG = "connectionString";
    private static final String JOB_REFERENCES_CONFIG = "jobReferences";
    private static final String TOPIC_CONFIG = "topic";
    private static final String INTERVAL_CONFIG = "interval";
    private static final String FIELDS_CONFIG = "fields";

    private final List<Source> sources;
    private final List<Job> jobs;

    public static SourceConfig fromPropertyMap(Map<String, String> properties) {
        String defaultTopic = properties.getOrDefault(Plc4xSourceConnector.DEFAULT_TOPIC_CONFIG, null);

        String[] sourceNames = properties.getOrDefault(Plc4xSourceConnector.SOURCES_CONFIG, "").split(",");
        List<Source> sources = new ArrayList<>(sourceNames.length);
        for (String sourceName : sourceNames) {
            String connectionString = properties.get(
                Plc4xSourceConnector.SOURCES_CONFIG + "." + sourceName + "." + CONNECTION_STRING_CONFIG);
            String sourceTopic = properties.getOrDefault(
                Plc4xSourceConnector.SOURCES_CONFIG + "." + sourceName + "." + TOPIC_CONFIG, defaultTopic);
            String[] jobReferenceNames = properties.get(
                Plc4xSourceConnector.SOURCES_CONFIG + "." + sourceName + "." + JOB_REFERENCES_CONFIG).split(",");
            JobReference[] jobReferences = new JobReference[jobReferenceNames.length];
            for (int i = 0; i < jobReferenceNames.length; i++) {
                String jobReferenceName = jobReferenceNames[i];
                String topic = properties.getOrDefault(Plc4xSourceConnector.SOURCES_CONFIG + "." + sourceName + "." + JOB_REFERENCES_CONFIG +
                    "." + jobReferenceName + "." + TOPIC_CONFIG, sourceTopic);
                JobReference jobReference = new JobReference(jobReferenceName, topic);
                jobReferences[i] = jobReference;
            }
            Source source = new Source(sourceName, connectionString, jobReferences);
            sources.add(source);
        }

        String[] jobNames = properties.getOrDefault(Plc4xSourceConnector.JOBS_CONFIG, "").split(",");
        List<Job> jobs = new ArrayList<>(jobNames.length);
        for (String jobName : jobNames) {
            int interval = Integer.parseInt(properties.get(
                Plc4xSourceConnector.JOBS_CONFIG + "." + jobName + "." + INTERVAL_CONFIG));
            String[] fieldNames = properties.get(
                Plc4xSourceConnector.JOBS_CONFIG + "." + jobName + "." + FIELDS_CONFIG).split(",");
            Map<String, String> fields = new HashMap<>();
            for (String fieldName : fieldNames) {
                String fieldAddress = properties.get(
                    Plc4xSourceConnector.JOBS_CONFIG + "." + jobName + "." + FIELDS_CONFIG + "." + fieldName);
                fields.put(fieldName, fieldAddress);
            }
            Job job = new Job(jobName, interval, fields);
            jobs.add(job);
        }

        return new SourceConfig(sources, jobs);
    }

    public SourceConfig(List<Source> sources, List<Job> jobs) {
        this.sources = sources;
        this.jobs = jobs;
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

}
