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

import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.sink.SinkConnector;
import org.apache.plc4x.kafka.util.VersionUtil;

import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;

public class Plc4xSinkConnector extends SinkConnector {
    static final String URL_CONFIG = "url";
    private static final String URL_DOC = "Connection string used by PLC4X to connect to the PLC";

    static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(URL_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.HIGH, URL_DOC);

    private String url;
    private String query;

    @Override
    public Class<? extends Task> taskClass() {
        return Plc4xSinkTask.class;
    }

    @Override
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        List<Map<String, String>> configs = new LinkedList<>();
        for (int i = 0; i < maxTasks; i++) {
            Map<String, String> taskConfig = new HashMap<>();
            taskConfig.put(URL_CONFIG, url);
            configs.add(taskConfig);
        }
        return configs;
    }

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(Plc4xSinkConnector.CONFIG_DEF, props);
        url = config.getString(URL_CONFIG);
    }

    @Override
    public void stop() {}

    @Override
    public ConfigDef config() {
        return CONFIG_DEF;
    }

    @Override
    public String version() {
        return VersionUtil.getVersion();
    }

}
