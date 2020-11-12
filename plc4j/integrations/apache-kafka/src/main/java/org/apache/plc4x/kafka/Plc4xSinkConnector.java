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
import org.apache.kafka.common.config.Config;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.common.config.ConfigValue;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.sink.SinkConnector;
import org.apache.plc4x.kafka.config.Job;
import org.apache.plc4x.kafka.config.JobReference;
import org.apache.plc4x.kafka.config.Sink;
import org.apache.plc4x.kafka.config.SinkConfig;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class Plc4xSinkConnector extends SinkConnector {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSinkConnector.class);

    public static final String SINK_CONFIG = "sinks";
    private static final String SINK_DOC = "List of sink names that will be configured.";

    private static final String CONNECTION_STRING_CONFIG = "connectionString";
    private static final String TOPIC_CONFIG = "topic";


    private SinkConfig sinkConfig;

    @Override
    public void start(Map<String, String> props) {
        sinkConfig = SinkConfig.fromPropertyMap(props);
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

        // For each configured source we'll start a dedicated scraper instance collecting
        // all the scraper jobs enabled for this source.
        List<Map<String, String>> configs = new LinkedList<>();

        for (Sink sink : sinkConfig.getSinks()) {
            // Create a new task configuration.
            Map<String, String> taskConfig = new HashMap<>();
            taskConfig.put(Plc4xSinkTask.CONNECTION_NAME_CONFIG, sink.getName());
            taskConfig.put(CONNECTION_STRING_CONFIG, sink.getConnectionString());
            taskConfig.put(TOPIC_CONFIG, sink.getTopic());
            configs.add(taskConfig);
        }
        return configs;
    }

    @Override
    @SuppressWarnings("unchecked")
    public Config validate(Map<String, String> connectorConfigs) {
        ////////////////////////////////////////////////////
        // Get the static part of the config
        Config config = super.validate(connectorConfigs);

        ////////////////////////////////////////////////////
        // Add the dynamic parts of the config

        // Find the important config elements
        ConfigValue sinks = null;

        for (ConfigValue configValue : config.configValues()) {
            switch (configValue.name()) {
                case SINK_CONFIG:
                    sinks = configValue;
                    break;
                default:
                    // Just ignore the others.
            }
        }

        // Configure the sinks
        if (sinks != null) {
            final List<String> sinkNames = (List<String>) sinks.value();

            if (sinkNames != null) {
                for (String sinkName : sinkNames) {
                    String connectionStringConfig = SINK_CONFIG + "." + sinkName + "." + CONNECTION_STRING_CONFIG;
                    final ConfigValue sinkConnectionStringConfigValue = new ConfigValue(connectionStringConfig);
                    config.configValues().add(sinkConnectionStringConfigValue);
                    String connectionString = connectorConfigs.get(connectionStringConfig);
                    sinkConnectionStringConfigValue.value(connectionString);
                    if (connectionString == null) {
                        sinkConnectionStringConfigValue.addErrorMessage(connectionStringConfig + " is mandatory");
                    } else {
                        // TODO: Check if the connection string is valid.
                        String sinkTopicConfig = SINK_CONFIG + "." + sinkName + "." + TOPIC_CONFIG;
                        final ConfigValue sinkTopicConfigValue = new ConfigValue(sinkTopicConfig);
                        config.configValues().add(sinkTopicConfigValue);
                        String sinkTopic = connectorConfigs.get(sinkTopicConfig);
                        sinkTopicConfigValue.value(sinkTopic);
                    }
                }
            }
        }

        return config;
    }

    @Override
    public ConfigDef config() {
        return new ConfigDef();
            //[BUG] - When including this, confluent fails when adding new connectors.
            //.define(SINK_CONFIG,
            //        ConfigDef.Type.LIST,
            //        new LinkedList<String>(),
            //        ConfigDef.Importance.LOW,
            //        SINK_DOC);
    }

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

}
