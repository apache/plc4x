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

import org.apache.plc4x.kafka.Plc4xSinkConnector;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SinkConfig {

    private static final Logger log = LoggerFactory.getLogger(SinkConfig.class);

    private static final String CONNECTION_STRING_CONFIG = "connectionString";
    private static final String TOPIC_CONFIG = "topic";
    private static final String RETRIES_CONFIG = "retries";
    private static final String TIMEOUT_CONFIG = "timeout";
    private static final String FIELDS_CONFIG = "fields";

    private final List<Sink> sinks;

    public static SinkConfig fromPropertyMap(Map<String, String> properties) {
        String[] sinkNames = properties.getOrDefault(Plc4xSinkConnector.SINK_CONFIG, "").split(",");
        List<Sink> sinks = new ArrayList<>(sinkNames.length);
        for (String sinkName : sinkNames) {
            String connectionString = properties.get(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + CONNECTION_STRING_CONFIG);
            String sinkTopic = properties.get(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + TOPIC_CONFIG);
            String sinkRetries = properties.get(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + RETRIES_CONFIG);
            String sinkTimeout = properties.get(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + TIMEOUT_CONFIG);
            String[] fieldNames = properties.get(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + FIELDS_CONFIG).split(",");
            Map<String, String> fields = new HashMap<>();
            for (String fieldName : fieldNames) {
                String fieldAddress = properties.get(
                    Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + FIELDS_CONFIG + "." + fieldName);
                fields.put(fieldName, fieldAddress);
            }
            Sink sink = new Sink(sinkName, connectionString, sinkTopic, fields, Integer.parseInt(sinkRetries), Integer.parseInt(sinkTimeout));
            sinks.add(sink);
        }

        return new SinkConfig(sinks);
    }

    public SinkConfig(List<Sink> sinks) {
        this.sinks = sinks;
    }

    public List<Sink> getSinks() {
        return sinks;
    }

    public Sink getSink(String sinkName) {
        if(sinks == null) {
            return null;
        }
        return sinks.stream().filter(sink -> sink.getName().equals(sinkName)).findFirst().orElse(null);
    }
}
