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

import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.tools.eventpump.EventPump;
import org.apache.plc4x.java.tools.eventpump.TagBatch;
import org.apache.plc4x.java.tools.eventpump.transform.ValueTransformerRegistry;
import org.apache.plc4x.java.tools.eventpump.triggers.SubscriptionTrigger;
import org.apache.plc4x.java.tools.eventpump.triggers.TimerTrigger;
import org.apache.plc4x.java.tools.eventpump.triggers.Trigger;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

/**
 * Factory for creating EventPump instances from configuration.
 * <p>
 * This class provides methods to build an EventPump from:
 * - Configuration objects
 * - Configuration files (JSON, YAML, XML)
 * <p>
 * Example usage:
 * <pre>
 * // From a configuration file
 * PlcConnectionManager connectionManager = PlcDriverManager.getDefault();
 * EventPump pump = EventPumpFactory.fromYaml(new File("config.yaml"), connectionManager,
 *     (batch, response) -> {
 *         // Handle all batches
 *     });
 * pump.startAll();
 * </pre>
 */
public class EventPumpFactory {

    private static final Logger LOGGER = LoggerFactory.getLogger(EventPumpFactory.class);

    /**
     * Create an EventPump from configuration.
     *
     * @param config The configuration
     * @param connectionManager The connection manager to use for creating connections
     * @param defaultListener The default listener for all batches (can be null)
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump create(EventPumpConfiguration config,
                                   PlcConnectionManager connectionManager,
                                   TagBatch.TagBatchListener defaultListener) throws Exception {
        return create(config, connectionManager, defaultListener, null);
    }

    /**
     * Create an EventPump from configuration with a custom transformer registry.
     *
     * @param config The configuration
     * @param connectionManager The connection manager to use for creating connections
     * @param defaultListener The default listener for all batches (can be null)
     * @param transformerRegistry The transformer registry to use (or null to create default)
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump create(EventPumpConfiguration config,
                                    PlcConnectionManager connectionManager,
                                    TagBatch.TagBatchListener defaultListener,
                                    ValueTransformerRegistry transformerRegistry) throws Exception {

        EventPump pump = new EventPump();

        // Create or use a provided transformer registry (shared across all batches)
        ValueTransformerRegistry registry = transformerRegistry != null ?
            transformerRegistry : ValueTransformerRegistry.createDefault();

        LOGGER.debug("Using transformer registry with transformers: {}", registry.getRegisteredNames());

        // Store connection URLs (authentication info included in URL or stored separately if needed)
        Map<String, String> connectionUrls = new HashMap<>();
        for (ConnectionConfiguration connConfig : config.getConnections()) {
            // Store the connection URL for later use
            // Note: For authenticated connections, credentials would need to be handled separately
            // or embedded in the URL format if the protocol supports it
            connectionUrls.put(connConfig.getId(), connConfig.getUrl());
            LOGGER.info("Registered connection '{}' to '{}'", connConfig.getId(), connConfig.getUrl());
        }

        // Create batches
        for (BatchConfiguration batchConfig : config.getBatches()) {
            String connectionUrl = connectionUrls.get(batchConfig.getConnectionId());
            if (connectionUrl == null) {
                throw new IllegalArgumentException("Connection '" + batchConfig.getConnectionId() +
                    "' not found for batch '" + batchConfig.getId() + "'");
            }

            // Create trigger (no longer needs connection parameter for timer triggers)
            Trigger trigger = createTrigger(batchConfig.getTrigger());

            // Create batch using builder with ConnectionManager
            TagBatch.Builder batchBuilder = TagBatch.builder()
                .withBatchId(batchConfig.getId())
                .withConnectionManager(connectionManager)
                .withConnectionString(connectionUrl)
                .addTagAddresses(batchConfig.getSimpleTagMap())
                .withTrigger(trigger)
                .withTransformerRegistry(registry); // Share registry across all batches

            if (batchConfig.getFetchTimeoutMs() != null) {
                batchBuilder.withFetchTimeout(batchConfig.getFetchTimeoutMs(), TimeUnit.MILLISECONDS);
            }

            // Add transformations
            for (Map.Entry<String, TagConfiguration> entry : batchConfig.getTags().entrySet()) {
                if (entry.getValue().hasTransform()) {
                    batchBuilder.addTransform(entry.getKey(), entry.getValue().getTransform());
                }
            }

            // Set listener if provided
            if (defaultListener != null) {
                batchBuilder.withListener(defaultListener);
            }

            TagBatch batch = batchBuilder.build();

            pump.addBatch(batch);
            LOGGER.info("Created batch '{}' with {} tags", batchConfig.getId(), batchConfig.getTags().size());
        }

        return pump;
    }

    /**
     * Create a trigger from configuration.
     *
     * @param config The trigger configuration
     * @return The trigger
     */
    private static Trigger createTrigger(TriggerConfiguration config) {
        String type = config.getType();

        if ("timer".equalsIgnoreCase(type)) {
            // Timer trigger - the interval may be given in seconds or, for sub-second
            // polling, in milliseconds. Accepting both and silently picking one would leave
            // the reader of the config guessing at the actual rate, so it's rejected.
            if (config.getIntervalSeconds() != null && config.getIntervalMillis() != null) {
                throw new IllegalArgumentException(
                    "Timer trigger has both intervalSeconds (" + config.getIntervalSeconds() +
                        ") and intervalMillis (" + config.getIntervalMillis() + ") set - use one of them");
            }
            if (config.getInitialDelaySeconds() != null && config.getInitialDelayMillis() != null) {
                throw new IllegalArgumentException(
                    "Timer trigger has both initialDelaySeconds (" + config.getInitialDelaySeconds() +
                        ") and initialDelayMillis (" + config.getInitialDelayMillis() + ") set - use one of them");
            }
            if (config.getIntervalSeconds() == null && config.getIntervalMillis() == null) {
                throw new IllegalArgumentException(
                    "Timer trigger requires either intervalSeconds or intervalMillis");
            }

            if (config.getIntervalMillis() != null) {
                long initialDelay = config.getInitialDelayMillis() != null ? config.getInitialDelayMillis()
                    : TimeUnit.SECONDS.toMillis(config.getInitialDelaySeconds() != null ? config.getInitialDelaySeconds() : 0);
                return new TimerTrigger(config.getIntervalMillis(), initialDelay, TimeUnit.MILLISECONDS);
            }

            long initialDelay = config.getInitialDelayMillis() != null ? config.getInitialDelayMillis()
                : TimeUnit.SECONDS.toMillis(config.getInitialDelaySeconds() != null ? config.getInitialDelaySeconds() : 0);

            return new TimerTrigger(TimeUnit.SECONDS.toMillis(config.getIntervalSeconds()), initialDelay,
                TimeUnit.MILLISECONDS);

        } else if ("subscription".equalsIgnoreCase(type)) {
            // Subscription trigger (placeholder)
            String tagName = config.getTagName();
            String tagAddress = config.getTagAddress();

            if (tagName == null || tagAddress == null) {
                throw new IllegalArgumentException("Subscription trigger requires tagName and tagAddress");
            }

            // Note: SubscriptionTrigger is currently a placeholder
            // When implemented, it will need to manage its own connection or use a different approach
            return new SubscriptionTrigger(tagName, tagAddress);

        } else {
            throw new IllegalArgumentException("Unknown trigger type: " + type);
        }
    }

    /**
     * Create an EventPump from a YAML file.
     *
     * @param file The YAML file
     * @param connectionManager The connection manager
     * @param defaultListener The default listener
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump fromYaml(java.io.File file,
                                      PlcConnectionManager connectionManager,
                                      TagBatch.TagBatchListener defaultListener) throws Exception {
        return fromYaml(file, connectionManager, defaultListener, null);
    }

    /**
     * Create an EventPump from a YAML file with a custom transformer registry.
     *
     * @param file The YAML file
     * @param connectionManager The connection manager
     * @param defaultListener The default listener
     * @param transformerRegistry The transformer registry (or null to create default)
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump fromYaml(java.io.File file,
                                      PlcConnectionManager connectionManager,
                                      TagBatch.TagBatchListener defaultListener,
                                      ValueTransformerRegistry transformerRegistry) throws Exception {
        EventPumpConfiguration config = EventPumpConfiguration.fromYaml(file);
        return create(config, connectionManager, defaultListener, transformerRegistry);
    }

    /**
     * Create an EventPump from a JSON file.
     *
     * @param file The JSON file
     * @param connectionManager The connection manager
     * @param defaultListener The default listener
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump fromJson(java.io.File file,
                                      PlcConnectionManager connectionManager,
                                      TagBatch.TagBatchListener defaultListener) throws Exception {
        return fromJson(file, connectionManager, defaultListener, null);
    }

    /**
     * Create an EventPump from a JSON file with a custom transformer registry.
     *
     * @param file The JSON file
     * @param connectionManager The connection manager
     * @param defaultListener The default listener
     * @param transformerRegistry The transformer registry (or null to create default)
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump fromJson(java.io.File file,
                                      PlcConnectionManager connectionManager,
                                      TagBatch.TagBatchListener defaultListener,
                                      ValueTransformerRegistry transformerRegistry) throws Exception {
        EventPumpConfiguration config = EventPumpConfiguration.fromJson(file);
        return create(config, connectionManager, defaultListener, transformerRegistry);
    }

    /**
     * Create an EventPump from an XML file.
     *
     * @param file The XML file
     * @param connectionManager The connection manager
     * @param defaultListener The default listener
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump fromXml(java.io.File file,
                                     PlcConnectionManager connectionManager,
                                     TagBatch.TagBatchListener defaultListener) throws Exception {
        return fromXml(file, connectionManager, defaultListener, null);
    }

    /**
     * Create an EventPump from an XML file with a custom transformer registry.
     *
     * @param file The XML file
     * @param connectionManager The connection manager
     * @param defaultListener The default listener
     * @param transformerRegistry The transformer registry (or null to create default)
     * @return The configured EventPump
     * @throws Exception if creation fails
     */
    public static EventPump fromXml(java.io.File file,
                                     PlcConnectionManager connectionManager,
                                     TagBatch.TagBatchListener defaultListener,
                                     ValueTransformerRegistry transformerRegistry) throws Exception {
        EventPumpConfiguration config = EventPumpConfiguration.fromXml(file);
        return create(config, connectionManager, defaultListener, transformerRegistry);
    }
}
