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

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.common.config.AbstractConfig;
import org.apache.kafka.common.config.ConfigDef;
import org.apache.kafka.connect.connector.Task;
import org.apache.kafka.connect.source.SourceConnector;
import org.apache.plc4x.kafka.util.VersionUtil;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.stream.Collectors;

public class Plc4xSourceConnector extends SourceConnector {

    private static final Logger log = LoggerFactory.getLogger(Plc4xSourceConnector.class);

    private static final String TOPIC_CONFIG = "topic";
    private static final String TOPIC_DOC = "Kafka topic to publish to";

    private static final String QUERIES_CONFIG = "queries";
    private static final String QUERIES_DOC = "Field queries to be sent to the PLC";

    private static final String JSON_CONFIG = "json.url";
    private static final String JSON_DEFAULT = "";
    private static final String JSON_DOC = "JSON configuration";

    private static final String RATE_CONFIG = "rate";
    private static final Integer RATE_DEFAULT = 1000;
    private static final String RATE_DOC = "Polling rate";

    private static final ConfigDef CONFIG_DEF = new ConfigDef()
        .define(TOPIC_CONFIG, ConfigDef.Type.STRING, ConfigDef.Importance.HIGH, TOPIC_DOC)
        .define(QUERIES_CONFIG, ConfigDef.Type.LIST, new LinkedList<>(), ConfigDef.Importance.HIGH, QUERIES_DOC)
        .define(JSON_CONFIG, ConfigDef.Type.STRING, JSON_DEFAULT, ConfigDef.Importance.HIGH, JSON_DOC)
        .define(RATE_CONFIG, ConfigDef.Type.INT, RATE_DEFAULT, ConfigDef.Importance.MEDIUM, RATE_DOC);

    private String topic;
    private List<String> queries;
    private String json;
    private Integer rate;

    @Override
    public Class<? extends Task> taskClass() {
        return Plc4xSourceTask.class;
    }

    @Override
    @SuppressWarnings("unchecked")
    public List<Map<String, String>> taskConfigs(int maxTasks) {
        List<Map<String, String>> configs = new LinkedList<>();
        if (json.isEmpty()) {
            Map<String, List<String>> groupedByHost = new HashMap<>();
            queries.stream().map(query -> query.split("#", 2)).collect(Collectors.groupingBy(parts -> parts[0])).forEach((host, q) ->
                groupedByHost.put(host, q.stream().map(parts -> parts[1]).collect(Collectors.toList())));
            if (groupedByHost.size() > maxTasks) {
                // Not enough tasks
                // TODO: throw exception?
                return Collections.emptyList();
            }
            groupedByHost.forEach((host, qs) -> {
                Map<String, String> taskConfig = new HashMap<>();
                taskConfig.put(Plc4xSourceTask.TOPIC_CONFIG, topic);
                taskConfig.put(Plc4xSourceTask.URL_CONFIG, host);
                taskConfig.put(Plc4xSourceTask.QUERIES_CONFIG, String.join(",", qs));
                taskConfig.put(Plc4xSourceTask.RATE_CONFIG, rate.toString());
                configs.add(taskConfig);
            });
        } else {
            try {
                // TODO
                String config =  new Scanner(new URL(json).openStream(), StandardCharsets.UTF_8.name()).useDelimiter("\\A").next();
                ObjectMapper mapper = new ObjectMapper();
                Map<String, Object> values = mapper.readValue(config, new TypeReference<Map<String, Object>>() {});
                List<Map<String, Object>> plcs = (List<Map<String, Object>>) values.get("PLCs");
                log.info("TASKS REQUIRED: " + plcs.size());
                if (plcs.size() > maxTasks) {
                    // not enough tasks
                    log.warn("NOT ENOUGH TASKS!");
                    return Collections.emptyList();
                }
                for (Map<String, Object> plc : plcs) {
                    Map<String, String> taskConfig = new HashMap<>();
                    String ip = plc.get("IP").toString();
                    String topic = ip;
                    String url = "s7://" + ip + "/1/" + plc.get("Slot");
                    List<String> queries = new LinkedList<>();
                    for (Map<String, Object> operand : (List<Map<String, Object>>)plc.get("operands")) {
                        String query = "%" + operand.get("Operand") + ":" + operand.get("Datatype");
                        queries.add(query);
                    }
                    taskConfig.put(Plc4xSourceTask.TOPIC_CONFIG, topic);
                    taskConfig.put(Plc4xSourceTask.URL_CONFIG, url);
                    taskConfig.put(RATE_CONFIG, rate.toString());
                    taskConfig.put(Plc4xSourceTask.QUERIES_CONFIG, String.join(",", queries));
                    configs.add(taskConfig);
                }
            } catch (IOException e) {
                log.error("ERROR CONFIGURING TASK", e);
            }
        }
        return configs;
    }

    @Override
    public void start(Map<String, String> props) {
        AbstractConfig config = new AbstractConfig(CONFIG_DEF, props);
        topic = config.getString(TOPIC_CONFIG);
        queries = config.getList(QUERIES_CONFIG);
        rate = config.getInt(RATE_CONFIG);
        json = config.getString(JSON_CONFIG);
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
