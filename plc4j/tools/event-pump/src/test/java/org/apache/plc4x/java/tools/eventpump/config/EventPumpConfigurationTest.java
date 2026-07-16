/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.tools.eventpump.config;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.File;
import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for EventPumpConfiguration serialization/deserialization.
 */
class EventPumpConfigurationTest {

    @TempDir
    File tempDir;

    @Test
    void testJsonSerialization() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        File jsonFile = new File(tempDir, "config.json");

        // Act - Save
        config.toJson(jsonFile);

        // Assert - File exists
        assertTrue(jsonFile.exists());

        // Act - Load
        EventPumpConfiguration loaded = EventPumpConfiguration.fromJson(jsonFile);

        // Assert - Content matches
        assertConfigurationEquals(config, loaded);
    }

    @Test
    void testYamlSerialization() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        File yamlFile = new File(tempDir, "config.yaml");

        // Act - Save
        config.toYaml(yamlFile);

        // Assert - File exists
        assertTrue(yamlFile.exists());

        // Act - Load
        EventPumpConfiguration loaded = EventPumpConfiguration.fromYaml(yamlFile);

        // Assert - Content matches
        assertConfigurationEquals(config, loaded);
    }

    @Test
    void testXmlSerialization() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        File xmlFile = new File(tempDir, "config.xml");

        // Act - Save
        config.toXml(xmlFile);

        // Assert - File exists
        assertTrue(xmlFile.exists());

        // Act - Load
        EventPumpConfiguration loaded = EventPumpConfiguration.fromXml(xmlFile);

        // Assert - Content matches
        assertConfigurationEquals(config, loaded);
    }

    private EventPumpConfiguration createTestConfiguration() {
        EventPumpConfiguration config = new EventPumpConfiguration();

        // Add connection
        ConnectionConfiguration conn = new ConnectionConfiguration();
        conn.setId("plc1");
        conn.setUrl("ads:tcp://192.168.1.1");
        conn.setUsername("admin");
        conn.setPassword("secret");
        config.getConnections().add(conn);

        // Add batch with timer trigger
        BatchConfiguration batch1 = new BatchConfiguration();
        batch1.setId("batch1");
        batch1.setConnectionId("plc1");

        Map<String, String> tags1 = new HashMap<>();
        tags1.put("temperature", "MAIN.temperature");
        tags1.put("pressure", "MAIN.pressure");
        batch1.setSimpleTags(tags1);

        TriggerConfiguration trigger1 = new TriggerConfiguration();
        trigger1.setType("timer");
        trigger1.setIntervalSeconds(5L);
        trigger1.setInitialDelaySeconds(2L);
        batch1.setTrigger(trigger1);

        config.getBatches().add(batch1);

        // Add batch with subscription trigger
        BatchConfiguration batch2 = new BatchConfiguration();
        batch2.setId("batch2");
        batch2.setConnectionId("plc1");

        Map<String, String> tags2 = new HashMap<>();
        tags2.put("status", "MAIN.status");
        batch2.setSimpleTags(tags2);

        TriggerConfiguration trigger2 = new TriggerConfiguration();
        trigger2.setType("subscription");
        trigger2.setTagName("triggerTag");
        trigger2.setTagAddress("MAIN.trigger");
        trigger2.setCondition("value == true");
        batch2.setTrigger(trigger2);

        config.getBatches().add(batch2);

        return config;
    }

    private void assertConfigurationEquals(EventPumpConfiguration expected, EventPumpConfiguration actual) {
        // Check connections
        assertEquals(expected.getConnections().size(), actual.getConnections().size());

        ConnectionConfiguration expectedConn = expected.getConnections().get(0);
        ConnectionConfiguration actualConn = actual.getConnections().get(0);

        assertEquals(expectedConn.getId(), actualConn.getId());
        assertEquals(expectedConn.getUrl(), actualConn.getUrl());
        assertEquals(expectedConn.getUsername(), actualConn.getUsername());
        assertEquals(expectedConn.getPassword(), actualConn.getPassword());

        // Check batches
        assertEquals(expected.getBatches().size(), actual.getBatches().size());

        // Check batch 1 (timer trigger)
        BatchConfiguration expectedBatch1 = expected.getBatches().get(0);
        BatchConfiguration actualBatch1 = actual.getBatches().get(0);

        assertEquals(expectedBatch1.getId(), actualBatch1.getId());
        assertEquals(expectedBatch1.getConnectionId(), actualBatch1.getConnectionId());
        assertEquals(expectedBatch1.getSimpleTagMap(), actualBatch1.getSimpleTagMap());

        TriggerConfiguration expectedTrigger1 = expectedBatch1.getTrigger();
        TriggerConfiguration actualTrigger1 = actualBatch1.getTrigger();

        assertEquals(expectedTrigger1.getType(), actualTrigger1.getType());
        assertEquals(expectedTrigger1.getIntervalSeconds(), actualTrigger1.getIntervalSeconds());
        assertEquals(expectedTrigger1.getInitialDelaySeconds(), actualTrigger1.getInitialDelaySeconds());

        // Check batch 2 (subscription trigger)
        BatchConfiguration expectedBatch2 = expected.getBatches().get(1);
        BatchConfiguration actualBatch2 = actual.getBatches().get(1);

        assertEquals(expectedBatch2.getId(), actualBatch2.getId());
        assertEquals(expectedBatch2.getConnectionId(), actualBatch2.getConnectionId());
        assertEquals(expectedBatch2.getSimpleTagMap(), actualBatch2.getSimpleTagMap());

        TriggerConfiguration expectedTrigger2 = expectedBatch2.getTrigger();
        TriggerConfiguration actualTrigger2 = actualBatch2.getTrigger();

        assertEquals(expectedTrigger2.getType(), actualTrigger2.getType());
        assertEquals(expectedTrigger2.getTagName(), actualTrigger2.getTagName());
        assertEquals(expectedTrigger2.getTagAddress(), actualTrigger2.getTagAddress());
        assertEquals(expectedTrigger2.getCondition(), actualTrigger2.getCondition());
    }
}
