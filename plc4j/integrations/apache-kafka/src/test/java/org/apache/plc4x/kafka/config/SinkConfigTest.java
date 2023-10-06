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

import org.apache.plc4x.kafka.Plc4xSinkConnector;
import org.apache.plc4x.kafka.Plc4xSinkTask;

import static org.junit.jupiter.api.Assertions.assertThrows;
import org.junit.jupiter.api.Test;

import org.apache.kafka.common.config.*;

import java.io.StringReader;
import java.nio.file.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.lang.NullPointerException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SinkConfigTest {

    private static final Logger log = LoggerFactory.getLogger(SinkConfigTest.class);
    private static final String TEST_PATH = "src/test/java/org/apache/plc4x/kafka/properties/";

    @Test
    public void validateExampleConfigFile() throws Exception{
        log.info("-----------------Sink Validate Example Config----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath(TEST_PATH, "sink_task_no_error.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSinkConnector sinkConnector = new Plc4xSinkConnector();
        Config config = sinkConnector.validate(map);

        for (ConfigValue value : config.configValues()) {
            log.info(value.name() + " = " + value.value());
            assertEquals(value.errorMessages(), new ArrayList());
        }
    }

    @Test
    public void validateExampleConfigFile2() throws Exception{
        log.info("-----------------Sink Validate Example Config 2----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath(TEST_PATH, "sink_task_no_error.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSinkConnector sinkConnector = new Plc4xSinkConnector();
        ConfigDef defs = sinkConnector.config();
        log.info(defs.toString());
        List<ConfigValue> configs = defs.validate(map);

        for (ConfigValue cValue : configs) {
            log.info(cValue.toString());
        }
    }

    @Test
    public void checkConnectorClass() throws Exception {
        log.info("-----------------CheckTaskClass----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath(TEST_PATH, "sink_task_no_error.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSinkConnector sinkConnector = new Plc4xSinkConnector();
        assertEquals(Plc4xSinkTask.class, sinkConnector.taskClass());
    }

    @Test
    public void checkConnectorStartStop() throws Exception {
        log.info("-----------------CheckConnectorStartStop----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath(TEST_PATH, "sink_task_no_error.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSinkConnector sinkConnector = new Plc4xSinkConnector();
        sinkConnector.start(map);
        sinkConnector.toString();
        sinkConnector.stop();
        assertThrows(NullPointerException.class, sinkConnector::toString);
    }

    private static Map<String, String> toStringMap(Properties properties) {
        Map<String, String> map = new HashMap<>();
        for (String stringPropertyName : properties.stringPropertyNames()) {
            map.put(stringPropertyName, properties.getProperty(stringPropertyName));
        }
        return map;
    }

}
