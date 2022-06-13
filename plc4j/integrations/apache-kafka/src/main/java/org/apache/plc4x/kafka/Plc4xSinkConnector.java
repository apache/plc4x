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
import org.apache.kafka.connect.sink.SinkConnector;
import org.apache.plc4x.kafka.config.*;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class Plc4xSinkConnector extends SinkConnector {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSinkConnector.class);

    private SinkConfig sinkConfig;
    private Map<String, String> configProps;

    @Override
    public void start(Map<String, String> props) {
        this.sinkConfig = new SinkConfig(props);
        this.configProps = Collections.unmodifiableMap(props);
    }

    @Override
    public void stop() {
        sinkConfig = null;
    }

    @Override
    public Class<? extends Task> taskClass() {
        return Plc4xSinkTask.class;
    }

    @Override
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        List<Map<String, String>> configs = new LinkedList<>();

        for (Sink sink : sinkConfig.getSinks()) {

            StringBuilder query = new StringBuilder();

            for (Field field : sink.getFields()) {
                String fieldName = field.getName();
                String fieldAddress = field.getAddress();
                query.append("|").append(fieldName).append("#").append(fieldAddress);
            }

            // Create a new task configuration.
            Map<String, String> taskConfig = new HashMap<>();
            taskConfig.put(Constants.CONNECTION_NAME_CONFIG, sink.getName());
            taskConfig.put(Constants.CONNECTION_STRING_CONFIG, sink.getConnectionString());
            taskConfig.put(Constants.TOPIC_CONFIG, sink.getTopic());
            taskConfig.put(Constants.RETRIES_CONFIG, sink.getRetries().toString());
            taskConfig.put(Constants.TIMEOUT_CONFIG, sink.getTimeout().toString());
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
        log.info("Validating PLC4X Sink Connector Configuration");

        SinkConfig sinkConfigTemp;
        try {
            sinkConfigTemp = new SinkConfig(connectorConfigs);
        } catch (Exception e) {
            for (ConfigValue configValue : config.configValues()) {
                if (configValue.errorMessages().size() > 0) {
                    return config;
                }
            }
            throw e;
        }
        sinkConfigTemp.validate();
        return config;
    }


    @Override
    public ConfigDef config() {
        return sinkConfig.configDef();
    }

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

    @Override
    public String toString(){
        return sinkConfig.toString();
    }

}
