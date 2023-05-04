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

import static org.junit.jupiter.api.Assertions.*;

import org.apache.kafka.common.record.TimestampType;

import org.apache.kafka.connect.sink.SinkRecord;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.data.SchemaBuilder;

import org.apache.plc4x.kafka.Plc4xSinkConnector;
import org.apache.plc4x.kafka.Plc4xSinkTask;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;

import java.io.StringReader;
import java.nio.file.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SinkTaskTest {

    private static final Logger log = LoggerFactory.getLogger(SinkTaskTest.class);

    private Plc4xSinkConnector sinkConnector;

    @BeforeEach
    public void setUp() throws Exception{
        log.info("-------------Setting Up SinkTaskTest----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath(
            "src/test/java/org/apache/plc4x/kafka/properties/",
            "sink_task_no_error.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        sinkConnector = new Plc4xSinkConnector();
        sinkConnector.start(toStringMap(properties));
    }

    @Test
    public void parseConfig() throws Exception {
        log.info("-----------------SinkTaskTest no Errors----------------");
        log.info(sinkConnector.toString());
        List<Map<String, String>> config = sinkConnector.taskConfigs(2);
        assertEquals(2, config.size());

        assertEquals("machineA", config.get(0).get(Constants.CONNECTION_NAME_CONFIG));
        assertEquals("simulated://127.0.0.1", config.get(0).get(Constants.CONNECTION_STRING_CONFIG));
        assertEquals("machineSinkA", config.get(0).get(Constants.TOPIC_CONFIG));
        assertEquals("5", config.get(0).get(Constants.RETRIES_CONFIG));
        assertEquals("5000", config.get(0).get(Constants.TIMEOUT_CONFIG));
        assertEquals("running#RANDOM/Temporary:Boolean|conveyorEntry#RANDOM/Temporary:Boolean|load#RANDOM/Temporary:Boolean|unload#RANDOM/Temporary:Boolean|transferLeft#RANDOM/Temporary:Boolean|transferRight#RANDOM/Temporary:Boolean|conveyorLeft#RANDOM/Temporary:Boolean|conveyorRight#RANDOM/Temporary:Boolean|numLargeBoxes#STATE/Temporary:Integer|numSmallBoxes#RANDOM/Temporary:Integer",
                    config.get(0).get(Constants.QUERIES_CONFIG));
    }

    @Test
    public void startTasks() throws Exception {
        log.info("-----------------SinkTaskTest.Put----------------");
        log.info(sinkConnector.toString());
        List<Map<String, String>> config = sinkConnector.taskConfigs(2);
        List<Plc4xSinkTask> sinkList = new ArrayList<>(config.size());
        for (Map<String, String> taskConfig : config) {
            log.info("Starting Sink Task");
            Plc4xSinkTask sinkTask = new Plc4xSinkTask();
            List<SinkRecord> records = new ArrayList<>(1);

            // Build the Schema for the result struct.
            Schema tagSchema = SchemaBuilder.struct()
                .name("org.apache.plc4x.kafka.schema.Field")
                .field("running", Schema.BOOLEAN_SCHEMA)
                .field("numLargeBoxes", Schema.INT32_SCHEMA)
                .build();


            Schema schema = SchemaBuilder.struct()
                            .name("org.apache.plc4x.kafka.schema.JobResult")
                            .field(Constants.TAGS_CONFIG, tagSchema)
                            .field(Constants.TIMESTAMP_CONFIG, Schema.INT64_SCHEMA)
                            .field("expires", Schema.OPTIONAL_INT64_SCHEMA)
                            .build();

            Struct tagsStruct = new Struct(tagSchema)
                                    .put("running", false)
                                    .put("numLargeBoxes", 765);

            Struct struct = new Struct(schema)
                                    .put("tags", tagsStruct)
                                    .put(Constants.TIMESTAMP_CONFIG, System.currentTimeMillis())
                                    .put("expires", 0L);

            records.add(new SinkRecord("machineSinkA",
                          1,
                          schema,
                          struct,
                          schema,
                          struct,
                          1,
                          0L,
                          TimestampType.CREATE_TIME));
                                      
            log.info("Sending Records to Sink task");
            sinkTask.start(taskConfig);
            sinkTask.put(records);
        }
    }

    private static Map<String, String> toStringMap(Properties properties) {
        Map<String, String> map = new HashMap<>();
        for (String stringPropertyName : properties.stringPropertyNames()) {
            map.put(stringPropertyName, properties.getProperty(stringPropertyName));
        }
        return map;
    }

}
