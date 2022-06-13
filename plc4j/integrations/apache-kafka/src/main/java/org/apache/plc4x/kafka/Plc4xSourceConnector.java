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
package org.apache.plc4x.kafka;

import org.apache.kafka.common.config.Config;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigValue;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.source.SourceConnector;
import org.apache.plc4x.kafka.config.Field;
import org.apache.plc4x.kafka.config.*;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class Plc4xSourceConnector extends SourceConnector {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSourceConnector.class);

    private SourceConfig sourceConfig;
    private Map<String, String> configProps;

    @Override
    public void start(Map<String, String> props) {
        this.sourceConfig = new SourceConfig(props);
        this.configProps = Collections.unmodifiableMap(props);
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
                    for (Field field : job.getFields()) {
                        String fieldName = field.getName();
                        String fieldAddress = field.getAddress();
                        query.append("|").append(fieldName).append("#").append(fieldAddress);
                    }
                }
            }

            // Create a new task configuration.
            Map<String, String> taskConfig = new HashMap<>();
            taskConfig.put(Constants.CONNECTION_NAME_CONFIG, source.getName());
            taskConfig.put(Constants.CONNECTION_STRING_CONFIG, source.getConnectionString());
            taskConfig.put(Constants.BUFFER_SIZE_CONFIG, source.getBufferSize().toString());
            taskConfig.put(Constants.KAFKA_POLL_RETURN_CONFIG, source.getPollReturnInterval().toString());
            taskConfig.put(Constants.QUERIES_CONFIG, query.toString().substring(1));
            configs.add(taskConfig);
        }
        return configs;
    }

    @Override
    @SuppressWarnings("unchecked")
    public Config validate(Map<String, String> connectorConfigs) {
        /////////////////////////////////////////////////////
        // Get the static part of the config
        Config config = super.validate(connectorConfigs);
        log.info("Validating PLC4X Source Connector Configuration");

        SourceConfig sourceConfigTemp;
        try {
            sourceConfigTemp = new SourceConfig(connectorConfigs);
        } catch (Exception e) {
            for (ConfigValue configValue : config.configValues()) {
                if (configValue.errorMessages().size() > 0) {
                    return config;
                }
            }
            throw e;
        }
        sourceConfigTemp.validate();
        return config;
    }

    @Override
    public ConfigDef config() {
        return sourceConfig.configDef();
    }

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public String toString(){
        return sourceConfig.toString();
    }

}
