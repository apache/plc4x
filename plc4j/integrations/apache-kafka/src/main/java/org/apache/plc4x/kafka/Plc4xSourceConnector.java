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
package org.apache.plc4x.kafka;

import org.apache.commons.lang3.math.NumberUtils;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigValue;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.source.SourceConnector;
import org.apache.plc4x.kafka.config.Job;
import org.apache.plc4x.kafka.config.JobReference;
import org.apache.plc4x.kafka.config.Source;
import org.apache.plc4x.kafka.config.SourceConfig;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class Plc4xSourceConnector extends SourceConnector {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSourceConnector.class);

    public static final String DEFAULT_TOPIC_CONFIG = "default-topic";
    private static final String DEFAULT_TOPIC_DOC = "Default topic to be used, if not otherwise configured.";

    public static final String SOURCES_CONFIG = "sources";
    private static final String SOURCES_DOC = "List of source names that will be configured.";

    public static final String JOBS_CONFIG = "jobs";
    private static final String JOBS_DOC = "List of job names that will be configured.";

    private static final String CONNECTION_STRING_CONFIG = "connectionString";
    private static final String JOB_REFERENCES_CONFIG = "jobReferences";
    private static final String TOPIC_CONFIG = "topic";
    private static final String INTERVAL_CONFIG = "interval";
    private static final String FIELDS_CONFIG = "fields";

    private SourceConfig sourceConfig;

    @Override
    public void start(Map<String, String> props) {
        sourceConfig = SourceConfig.fromPropertyMap(props);
    }

    @Override
    public void stop() {
        sourceConfig = null;
    }

    @Override
    public Class<? extends Task> taskClass() {
        return Plc4xSourceTask.class;
    }

    @Override
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        // Initially we planned to have the simple assumption that one task maps to one PLC connection.
        // But we could easily say that one scraper instance maps to a task and one scraper task can
        // process multiple PLC connections. But I guess this would be an optimization as we have to
        // balance the load manually.
        if(sourceConfig.getJobs().size() > maxTasks) {
            // not enough tasks
            log.warn("NOT ENOUGH TASKS!");
            return Collections.emptyList();
        }

        // For each configured source we'll start a dedicated scraper instance collecting
        // all the scraper jobs enabled for this source.
        List<Map<String, String>> configs = new LinkedList<>();
        for (Source source : sourceConfig.getSources()) {
            // Build a list of job configurations only containing the ones referenced from
            // the current source.
            StringBuilder query = new StringBuilder();
            for (JobReference jobReference : source.getJobReferences()) {
                Job job = sourceConfig.getJob(jobReference.getName());
                if(job == null) {
                    log.warn(String.format("Couldn't find referenced job '%s'", jobReference.getName()));
                } else {
                    query.append(",").append(jobReference.getName()).append("|").append(jobReference.getTopic());
                    query.append("|").append(job.getInterval());
                    for (Map.Entry<String, String> field : job.getFields().entrySet()) {
                        String fieldName = field.getKey();
                        String fieldAddress = field.getValue();
                        query.append("|").append(fieldName).append("#").append(fieldAddress);
                    }
                }
            }

            // Create a new task configuration.
            Map<String, String> taskConfig = new HashMap<>();
            taskConfig.put(Plc4xSourceTask.CONNECTION_NAME_CONFIG, source.getName());
            taskConfig.put(Plc4xSourceTask.PLC4X_CONNECTION_STRING_CONFIG, source.getConnectionString());
            taskConfig.put(Plc4xSourceTask.QUERIES_CONFIG, query.toString().substring(1));
            configs.add(taskConfig);
        }
        return configs;
    }

    @Override
    public ConfigDef config() {
        return new ConfigDef() {
            @Override
            public Map<String, ConfigValue> validateAll(Map<String, String> props) {
                // Make sure all basic config options are validated.
                Map<String, ConfigValue> result = super.validateAll(props);

                final String[] jobNames = props.getOrDefault(JOBS_CONFIG, "").split(",");
                for (String jobName : jobNames) {
                    String jobIntervalConfig = JOBS_CONFIG + "." + jobName + "." + INTERVAL_CONFIG;
                    ConfigValue configValue = new ConfigValue(jobIntervalConfig);
                    result.put(jobIntervalConfig, new ConfigValue(jobIntervalConfig));
                    String jobIntervalString = props.get(jobIntervalConfig);
                    if(jobIntervalString == null) {
                        configValue.addErrorMessage(jobIntervalConfig + " is mandatory");
                    } else if(NumberUtils.isParsable(jobIntervalString)) {
                        int jobInterval = NumberUtils.toInt(jobIntervalString);
                        if(jobInterval > 0) {
                            configValue.value(jobInterval);
                        } else {
                            configValue.addErrorMessage(
                                jobIntervalConfig + " must be greater than 0");
                        }
                    } else {
                        configValue.addErrorMessage(jobIntervalConfig + " must be a numeric value greater than 0");
                    }

                    String jobFieldsConfig = JOBS_CONFIG + "." + jobName + "." + FIELDS_CONFIG;
                    configValue = new ConfigValue(jobFieldsConfig);
                    if(!props.containsKey(jobFieldsConfig)) {
                        configValue.addErrorMessage(jobFieldsConfig + " is mandatory");
                    } else {
                        String[] jobFieldNames = props.getOrDefault(jobFieldsConfig, "").split(",");
                        if (jobFieldNames.length == 0) {
                            configValue.addErrorMessage(jobFieldsConfig + " at least has to contain one field name");
                        } else {
                            for (String jobFieldName : jobFieldNames) {
                                String jobFieldAddressConfig =
                                    JOBS_CONFIG + "." + jobName + "." + FIELDS_CONFIG + "." + jobFieldName;
                                configValue = new ConfigValue(jobFieldAddressConfig);
                                String jobFieldAddress = props.get(jobFieldAddressConfig);
                                if((jobFieldAddress == null) || jobFieldAddress.isEmpty()) {
                                    configValue.addErrorMessage(jobFieldAddressConfig + " is mandatory");
                                } else {
                                    configValue.value(jobFieldAddress);
                                }
                            }
                        }
                    }
                }

                final String[] sourceNames = props.getOrDefault(SOURCES_CONFIG, "").split(",");
                for (String sourceName : sourceNames) {
                    String connectionStringConfig = SOURCES_CONFIG + "." + sourceName + "." + CONNECTION_STRING_CONFIG;
                    ConfigValue configValue = new ConfigValue(connectionStringConfig);
                    if(!props.containsKey(connectionStringConfig)) {
                        configValue.addErrorMessage(connectionStringConfig + " is mandatory");
                    } else {
                        String connectionString = props.get(connectionStringConfig);
                        // TODO: Check if the connection string is valid.
                        result.put(connectionStringConfig, new ConfigValue(connectionStringConfig));
                        result.get(connectionStringConfig).addErrorMessage("");

                        String connectionTopicConfig = SOURCES_CONFIG + "." + sourceName + "." + TOPIC_CONFIG;
                        String connectionTopic = props.get(connectionTopicConfig);

                        String jobReferenceNamesConfig = SOURCES_CONFIG + "." + sourceName + "." + JOB_REFERENCES_CONFIG;
                        String[] jobReferenceNames = props.getOrDefault(jobReferenceNamesConfig, "").split(",");
                        // TODO: Check at least one reference is provided
                        for (String jobReferenceName : jobReferenceNames) {
                        }
                        // TODO: Check that for all of these a job with the given name is provided
                        result.put(jobReferenceNamesConfig, new ConfigValue(jobReferenceNamesConfig));
                        result.get(jobReferenceNamesConfig).addErrorMessage("");
                    }
                }

                // TODO: Validate each combination of source and job to check if the addresses are valid for the given driver type.

                return result;
            }

        }
            .define(DEFAULT_TOPIC_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.LOW, DEFAULT_TOPIC_DOC)
            .define(SOURCES_CONFIG, ConfigDef.Type.LIST, ConfigDef.Importance.HIGH, SOURCES_DOC)
            .define(JOBS_CONFIG, ConfigDef.Type.LIST, ConfigDef.Importance.HIGH, JOBS_DOC);
    }

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

}
