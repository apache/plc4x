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

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.tools.eventpump.EventPump;
import org.apache.plc4x.java.tools.eventpump.TagBatch;
import org.apache.plc4x.java.tools.eventpump.triggers.TimerTrigger;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.File;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for EventPumpFactory.
 */
class EventPumpFactoryTest {

    @TempDir
    File tempDir;

    private PlcConnectionManager connectionManager;

    @BeforeEach
    void setUp() throws Exception {
        // Create a stub connection manager
        connectionManager = new PlcConnectionManager() {
            @Override
            public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
                return new StubPlcConnection();
            }

            @Override
            public PlcConnection getConnection(String connectionString, org.apache.plc4x.java.api.authentication.PlcAuthentication authentication) throws PlcConnectionException {
                return new StubPlcConnection();
            }
        };
    }

    /**
     * Stub implementation of PlcConnection for testing.
     */
    private static class StubPlcConnection implements PlcConnection {
        @Override
        public void connect() throws PlcConnectionException {
        }

        @Override
        public void close() throws Exception {
        }

        @Override
        public boolean isConnected() {
            return true;
        }

        @Override
        public Optional<PlcTag> parseTagAddress(String tagAddress) {
            return Optional.empty();
        }

        @Override
        public Optional<PlcValue> parseTagValue(PlcTag tag, Object... values) {
            return Optional.empty();
        }

        @Override
        public CompletableFuture<? extends PlcPingResponse> ping() {
            return null;
        }

        @Override
        public PlcReadRequest.Builder readRequestBuilder() {
            return null;
        }

        @Override
        public PlcWriteRequest.Builder writeRequestBuilder() {
            return null;
        }

        @Override
        public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
            return null;
        }

        @Override
        public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
            return null;
        }

        @Override
        public PlcBrowseRequest.Builder browseRequestBuilder() {
            return null;
        }

        @Override
        public PlcConnectionMetadata getMetadata() {
            return null;
        }
    }

    @Test
    void testCreateFromConfiguration() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        TagBatch.TagBatchListener listener = (batch, response) -> {};

        // Act
        EventPump pump = EventPumpFactory.create(config, connectionManager, listener);

        // Assert
        assertNotNull(pump);
        assertEquals(2, pump.getAllBatches().size());
        assertTrue(pump.getAllBatches().containsKey("batch1"));
        assertTrue(pump.getAllBatches().containsKey("batch2"));

        // Cleanup
        pump.close();
    }

    @Test
    void testCreateWithoutListener() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();

        // Act
        EventPump pump = EventPumpFactory.create(config, connectionManager, null);

        // Assert
        assertNotNull(pump);
        assertEquals(2, pump.getAllBatches().size());

        // Cleanup
        pump.close();
    }

    @Test
    void testCreateWithMissingConnection() {
        // Arrange
        EventPumpConfiguration config = new EventPumpConfiguration();

        // Add batch without adding the connection it references
        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("nonexistent");

        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "address1");
        batch.setSimpleTags(tags);

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("timer");
        trigger.setIntervalSeconds(1L);
        batch.setTrigger(trigger);

        config.getBatches().add(batch);

        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> EventPumpFactory.create(config, connectionManager, null)
        );
        assertTrue(exception.getMessage().contains("Connection 'nonexistent' not found"));
    }

    /**
     * Polling faster than once per second has to be expressible - it is the common case for
     * PLC data collection, and the scraper this replaces configured its rate in milliseconds.
     */
    @Test
    void testCreateWithMillisecondInterval() throws Exception {
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration connection = new ConnectionConfiguration();
        connection.setId("conn1");
        connection.setUrl("mock:test");
        config.getConnections().add(connection);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");
        batch.setSimpleTags(Collections.singletonMap("tag1", "address1"));

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("timer");
        trigger.setIntervalMillis(250L);
        trigger.setInitialDelayMillis(50L);
        batch.setTrigger(trigger);
        config.getBatches().add(batch);

        EventPump pump = EventPumpFactory.create(config, connectionManager, null);

        TimerTrigger timerTrigger = (TimerTrigger) pump.getBatch("batch1").getTrigger();
        assertEquals(250, timerTrigger.getIntervalMs());
        assertEquals(50, timerTrigger.getInitialDelayMs());
        pump.close();
    }

    /**
     * Setting the interval in both units is ambiguous to whoever reads the config later, so it
     * has to fail loudly rather than silently picking one.
     */
    @Test
    void testRejectsIntervalInBothUnits() {
        EventPumpConfiguration config = configWithTrigger(trigger -> {
            trigger.setIntervalSeconds(10L);
            trigger.setIntervalMillis(100L);
        });

        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> EventPumpFactory.create(config, connectionManager, null)
        );
        assertTrue(exception.getMessage().contains("intervalSeconds"));
        assertTrue(exception.getMessage().contains("intervalMillis"));
        // The message has to name the offending values, otherwise finding them in a large
        // configuration file is guesswork.
        assertTrue(exception.getMessage().contains("10"));
        assertTrue(exception.getMessage().contains("100"));
    }

    @Test
    void testRejectsInitialDelayInBothUnits() {
        EventPumpConfiguration config = configWithTrigger(trigger -> {
            trigger.setIntervalMillis(100L);
            trigger.setInitialDelaySeconds(5L);
            trigger.setInitialDelayMillis(50L);
        });

        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> EventPumpFactory.create(config, connectionManager, null)
        );
        assertTrue(exception.getMessage().contains("initialDelaySeconds"));
        assertTrue(exception.getMessage().contains("initialDelayMillis"));
    }

    /**
     * Mixing the units across the two settings stays legal - only the same setting given twice
     * is ambiguous.
     */
    @Test
    void testAllowsIntervalAndInitialDelayInDifferentUnits() throws Exception {
        EventPumpConfiguration config = configWithTrigger(trigger -> {
            trigger.setIntervalMillis(250L);
            trigger.setInitialDelaySeconds(2L);
        });

        EventPump pump = EventPumpFactory.create(config, connectionManager, null);

        TimerTrigger timerTrigger = (TimerTrigger) pump.getBatch("batch1").getTrigger();
        assertEquals(250, timerTrigger.getIntervalMs());
        assertEquals(2000, timerTrigger.getInitialDelayMs());
        pump.close();
    }

    /**
     * Builds a single-batch configuration whose timer trigger is set up by the given callback.
     */
    private static EventPumpConfiguration configWithTrigger(java.util.function.Consumer<TriggerConfiguration> customizer) {
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration connection = new ConnectionConfiguration();
        connection.setId("conn1");
        connection.setUrl("mock:test");
        config.getConnections().add(connection);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");
        batch.setSimpleTags(Collections.singletonMap("tag1", "address1"));

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("timer");
        customizer.accept(trigger);
        batch.setTrigger(trigger);
        config.getBatches().add(batch);
        return config;
    }

    /**
     * A timer trigger without any interval used to fail with a NullPointerException while
     * unboxing the interval.
     */
    @Test
    void testCreateWithTimerTriggerWithoutInterval() {
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration connection = new ConnectionConfiguration();
        connection.setId("conn1");
        connection.setUrl("mock:test");
        config.getConnections().add(connection);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");
        batch.setSimpleTags(Collections.singletonMap("tag1", "address1"));

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("timer");
        batch.setTrigger(trigger);
        config.getBatches().add(batch);

        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> EventPumpFactory.create(config, connectionManager, null)
        );
        assertTrue(exception.getMessage().contains("intervalSeconds or intervalMillis"));
    }

    @Test
    void testCreateWithUnknownTriggerType() {
        // Arrange
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration conn = new ConnectionConfiguration();
        conn.setId("conn1");
        conn.setUrl("test://localhost");
        config.getConnections().add(conn);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");

        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "address1");
        batch.setSimpleTags(tags);

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("unknown");
        batch.setTrigger(trigger);

        config.getBatches().add(batch);

        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> EventPumpFactory.create(config, connectionManager, null)
        );
        assertTrue(exception.getMessage().contains("Unknown trigger type"));
    }

    @Test
    void testCreateSubscriptionTriggerWithMissingTagName() {
        // Arrange
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration conn = new ConnectionConfiguration();
        conn.setId("conn1");
        conn.setUrl("test://localhost");
        config.getConnections().add(conn);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");

        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "address1");
        batch.setSimpleTags(tags);

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("subscription");
        trigger.setTagAddress("MAIN.trigger");
        // Missing tagName
        batch.setTrigger(trigger);

        config.getBatches().add(batch);

        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> EventPumpFactory.create(config, connectionManager, null)
        );
        assertTrue(exception.getMessage().contains("Subscription trigger requires tagName and tagAddress"));
    }

    @Test
    void testFromYaml() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        File yamlFile = new File(tempDir, "config.yaml");
        config.toYaml(yamlFile);

        TagBatch.TagBatchListener listener = (batch, response) -> {};

        // Act
        EventPump pump = EventPumpFactory.fromYaml(yamlFile, connectionManager, listener);

        // Assert
        assertNotNull(pump);
        assertEquals(2, pump.getAllBatches().size());

        // Cleanup
        pump.close();
    }

    @Test
    void testFromJson() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        File jsonFile = new File(tempDir, "config.json");
        config.toJson(jsonFile);

        TagBatch.TagBatchListener listener = (batch, response) -> {};

        // Act
        EventPump pump = EventPumpFactory.fromJson(jsonFile, connectionManager, listener);

        // Assert
        assertNotNull(pump);
        assertEquals(2, pump.getAllBatches().size());

        // Cleanup
        pump.close();
    }

    @Test
    void testFromXml() throws Exception {
        // Arrange
        EventPumpConfiguration config = createTestConfiguration();
        File xmlFile = new File(tempDir, "config.xml");
        config.toXml(xmlFile);

        TagBatch.TagBatchListener listener = (batch, response) -> {};

        // Act
        EventPump pump = EventPumpFactory.fromXml(xmlFile, connectionManager, listener);

        // Assert
        assertNotNull(pump);
        assertEquals(2, pump.getAllBatches().size());

        // Cleanup
        pump.close();
    }

    @Test
    void testCreateTimerTriggerWithInitialDelay() throws Exception {
        // Arrange
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration conn = new ConnectionConfiguration();
        conn.setId("conn1");
        conn.setUrl("test://localhost");
        config.getConnections().add(conn);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");

        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "address1");
        batch.setSimpleTags(tags);

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("timer");
        trigger.setIntervalSeconds(5L);
        trigger.setInitialDelaySeconds(2L);
        batch.setTrigger(trigger);

        config.getBatches().add(batch);

        // Act
        EventPump pump = EventPumpFactory.create(config, connectionManager, null);

        // Assert
        assertNotNull(pump);
        TagBatch createdBatch = pump.getAllBatches().get("batch1");
        assertNotNull(createdBatch);

        // Cleanup
        pump.close();
    }

    @Test
    void testCreateTimerTriggerWithoutInitialDelay() throws Exception {
        // Arrange
        EventPumpConfiguration config = new EventPumpConfiguration();

        ConnectionConfiguration conn = new ConnectionConfiguration();
        conn.setId("conn1");
        conn.setUrl("test://localhost");
        config.getConnections().add(conn);

        BatchConfiguration batch = new BatchConfiguration();
        batch.setId("batch1");
        batch.setConnectionId("conn1");

        Map<String, String> tags = new HashMap<>();
        tags.put("tag1", "address1");
        batch.setSimpleTags(tags);

        TriggerConfiguration trigger = new TriggerConfiguration();
        trigger.setType("timer");
        trigger.setIntervalSeconds(5L);
        // No initial delay - should default to 0
        batch.setTrigger(trigger);

        config.getBatches().add(batch);

        // Act
        EventPump pump = EventPumpFactory.create(config, connectionManager, null);

        // Assert
        assertNotNull(pump);
        TagBatch createdBatch = pump.getAllBatches().get("batch1");
        assertNotNull(createdBatch);

        // Cleanup
        pump.close();
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
}
