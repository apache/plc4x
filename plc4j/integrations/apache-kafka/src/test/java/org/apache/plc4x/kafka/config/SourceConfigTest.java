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

import org.apache.plc4x.kafka.Plc4xSourceConnector;
import org.apache.plc4x.kafka.Plc4xSourceTask;

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

public class SourceConfigTest {

    private static final Logger log = LoggerFactory.getLogger(SourceConfigTest.class);

    @Test
    public void parseConfig() throws Exception {
        Properties properties = new Properties();
        properties.load(new StringReader("name=plc-source-test\n" +
            "connector.class=org.apache.plc4x.kafka.Plc4xSourceConnector\n" +
            "\n" +
            "default.topic=some/default\n" +
            "\n" +
            "sources=machineA,machineB,machineC\n" +
            "sources.machineA.connectionString=s7://1.2.3.4/1/1\n" +
            "sources.machineA.pollReturnInterval=5000\n" +
            "sources.machineA.bufferSize=1000\n" +
            "sources.machineA.jobReferences=s7-dashboard,s7-heartbeat\n" +
            "sources.machineA.jobReferences.s7-heartbeat.topic=heartbeat\n" +
            "\n" +
            "sources.machineB.connectionString=s7://10.20.30.40/1/1\n" +
            "sources.machineB.topic=heartbeat\n" +
            "sources.machineB.jobReferences=s7-heartbeat\n" +
            "\n" +
            "sources.machineC.connectionString=modbus-tcp://192.168.1.1:502\n" +
            "sources.machineC.topic=heartbeat\n" +
            "sources.machineC.jobReferences=ads-heartbeat\n" +
            "\n" +
            "jobs=s7-dashboard,s7-heartbeat,ads-heartbeat\n" +
            "jobs.s7-dashboard.interval=500\n" +
            "jobs.s7-dashboard.fields=inputPreasure,outputPreasure,temperature\n" +
            "jobs.s7-dashboard.fields.inputPreasure=%DB.DB1.4:INT\n" +
            "jobs.s7-dashboard.fields.outputPreasure=%Q1:BYT\n" +
            "jobs.s7-dashboard.fields.temperature=%I3:INT\n" +
            "\n" +
            "jobs.s7-heartbeat.interval=1000\n" +
            "jobs.s7-heartbeat.fields=active\n" +
            "jobs.s7-heartbeat.fields.active=%I0.2:BOOL\n" +
            "\n" +
            "jobs.ads-heartbeat.interval=1000\n" +
            "jobs.ads-heartbeat.fields=active\n" +
            "jobs.ads-heartbeat.fields.active=Main.running\n"));
        SourceConfig sourceConfig = new SourceConfig(toStringMap(properties));

        assertNotNull(sourceConfig);
        assertEquals(3, sourceConfig.getSources().size(), "Expected 3 sources");
        assertEquals(3, sourceConfig.getJobs().size(), "Expected 3 jobs");
    }

    @Test
    public void validateConfig() throws Exception{
        log.info("-----------------Validate Config----------------");
        Properties properties = new Properties();
        properties.load(new StringReader("name=plc-source-test\n" +
            "connector.class=org.apache.plc4x.kafka.Plc4xSourceConnector\n" +
            "\n" +
            "default.topic=some/default\n" +
            "\n" +
            "sources=machineA,machineB,machineC\n" +
            "sources.machineA.connectionString=s7://1.2.3.4/1/1\n" +
            "sources.machineA.pollReturnInterval=5000\n" +
            "sources.machineA.bufferSize=1000\n" +
            "sources.machineA.jobReferences=s7-dashboard,s7-heartbeat\n" +
            "sources.machineA.jobReferences.s7-heartbeat.topic=heartbeat\n" +
            "\n" +
            "sources.machineB.connectionString=s7://10.20.30.40/1/1\n" +
            "sources.machineB.topic=heartbeat\n" +
            "sources.machineB.jobReferences=s7-heartbeat\n" +
            "\n" +
            "sources.machineC.connectionString=modbus-tcp://127.0.0.1\n" +
            "sources.machineC.topic=heartbeat\n" +
            "sources.machineC.jobReferences=ads-heartbeat\n" +
            "\n" +
            "jobs=s7-dashboard,s7-heartbeat,ads-heartbeat\n" +
            "jobs.s7-dashboard.interval=500\n" +
            "jobs.s7-dashboard.fields=inputPreasure,outputPreasure,temperature\n" +
            "jobs.s7-dashboard.fields.inputPreasure=%DB.DB1.4:INT\n" +
            "jobs.s7-dashboard.fields.outputPreasure=%Q1:BYT\n" +
            "jobs.s7-dashboard.fields.temperature=%I3:INT\n" +
            "\n" +
            "jobs.s7-heartbeat.interval=1000\n" +
            "jobs.s7-heartbeat.fields=active\n" +
            "jobs.s7-heartbeat.fields.active=%I0.2:BOOL\n" +
            "\n" +
            "jobs.ads-heartbeat.interval=1000\n" +
            "jobs.ads-heartbeat.fields=active\n" +
            "jobs.ads-heartbeat.fields.active=Main.running\n"));

        Map<String, String> map = new HashMap<String, String>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSourceConnector sourceConnector = new Plc4xSourceConnector();
        Config config = sourceConnector.validate(map);

        for (ConfigValue value : config.configValues()) {
            log.info(value.name() + " = " + value.value());
            assertEquals(value.errorMessages(), new ArrayList());
        }
    }

    @Test
    public void validateExampleConfigFile() throws Exception{
        log.info("-----------------Validate Example Config----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath("config/", "plc4x-source.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSourceConnector sourceConnector = new Plc4xSourceConnector();
        Config config = sourceConnector.validate(map);

        for (ConfigValue value : config.configValues()) {
            log.info(value.name() + " = " + value.value());
            assertEquals(value.errorMessages(), new ArrayList());
        }
    }

    @Test
    public void validateExampleConfigFile2() throws Exception{
        log.info("-----------------Validate Example Config 2----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath("config/", "plc4x-source.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<String, String>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSourceConnector sourceConnector = new Plc4xSourceConnector();
        ConfigDef defs = sourceConnector.config();
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
        Path path = FileSystems.getDefault().getPath("config/", "plc4x-source.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<String, String>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSourceConnector sourceConnector = new Plc4xSourceConnector();
        assertEquals(Plc4xSourceTask.class, sourceConnector.taskClass());
    }

    @Test
    public void checkConnectorStartStop() throws Exception {
        log.info("-----------------CheckConnectorStartStop----------------");
        Properties properties = new Properties();
        Path path = FileSystems.getDefault().getPath("config/", "plc4x-source.properties");
        properties.load((new StringReader(new String(Files.readAllBytes(path)))));

        Map<String, String> map = new HashMap<String, String>();
        for (final String name: properties.stringPropertyNames())
            map.put(name, properties.getProperty(name));

        Plc4xSourceConnector sourceConnector = new Plc4xSourceConnector();
        sourceConnector.start(map);
        sourceConnector.toString();
        sourceConnector.stop();
        assertThrows(NullPointerException.class, () -> sourceConnector.toString());
    }

    @Test
    public void checkVersionString() throws Exception {
        log.info("-----------------CheckVersionString----------------");
        Plc4xSourceConnector sourceConnector = new Plc4xSourceConnector();
        assertNotEquals("0.0.0.0", sourceConnector.version());
    }

    private static Map<String, String> toStringMap(Properties properties) {
        Map<String, String> map = new HashMap<>();
        for (String stringPropertyName : properties.stringPropertyNames()) {
            map.put(stringPropertyName, properties.getProperty(stringPropertyName));
        }
        return map;
    }

}
