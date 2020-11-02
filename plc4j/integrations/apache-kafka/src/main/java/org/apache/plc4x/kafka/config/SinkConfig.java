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

public class SinkConfig {

    private static final String CONNECTION_STRING_CONFIG = "connectionString";
    private static final String TOPIC_CONFIG = "topic";

    private final List<Sink> sinks;

    public static SinkConfig fromPropertyMap(Map<String, String> properties) {
        String defaultTopic = properties.getOrDefault(Plc4xSinkConnector.DEFAULT_TOPIC_CONFIG, null);

        String[] sinkNames = properties.getOrDefault(Plc4xSinkConnector.SINK_CONFIG, "").split(",");
        List<Sink> sinks = new ArrayList<>(sinkNames.length);
        for (String sinkName : sinkNames) {
            String connectionString = properties.get(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + CONNECTION_STRING_CONFIG);
            String sinkTopic = properties.getOrDefault(
                Plc4xSinkConnector.SINK_CONFIG + "." + sinkName + "." + TOPIC_CONFIG, defaultTopic);
            Sink sink = new Sink(sinkName, connectionString);
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
