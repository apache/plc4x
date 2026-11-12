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

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * Configuration model for EventPump.
 * <p>
 * This class can be serialized/deserialized from XML, JSON, or YAML files.
 * <p>
 * Example YAML configuration:
 * <pre>
 * connections:
 *   - id: plc1
 *     url: "ads:tcp://192.168.1.1"
 *
 * batches:
 *   - id: batch1
 *     connectionId: plc1
 *     tags:
 *       temperature: "MAIN.temperature"
 *       pressure: "MAIN.pressure"
 *     trigger:
 *       type: timer
 *       intervalSeconds: 5
 *
 *   - id: batch2
 *     connectionId: plc1
 *     tags:
 *       status: "MAIN.status"
 *     trigger:
 *       type: subscription
 *       tagName: triggerTag
 *       tagAddress: "MAIN.trigger"
 * </pre>
 */
public class EventPumpConfiguration {

    private List<ConnectionConfiguration> connections = new ArrayList<>();

    private List<BatchConfiguration> batches = new ArrayList<>();

    /**
     * Get the connection configurations.
     *
     * @return The list of connection configurations
     */
    public List<ConnectionConfiguration> getConnections() {
        return connections;
    }

    /**
     * Set the connection configurations.
     *
     * @param connections The connection configurations
     */
    public void setConnections(List<ConnectionConfiguration> connections) {
        this.connections = connections;
    }

    /**
     * Get the batch configurations.
     *
     * @return The list of batch configurations
     */
    public List<BatchConfiguration> getBatches() {
        return batches;
    }

    /**
     * Set the batch configurations.
     *
     * @param batches The batch configurations
     */
    public void setBatches(List<BatchConfiguration> batches) {
        this.batches = batches;
    }

    /**
     * Load configuration from a JSON file.
     *
     * @param file The file to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromJson(File file) throws IOException {
        try (InputStream inputStream = new FileInputStream(file)) {
            return fromJson(inputStream);
        }
    }

    /**
     * Load configuration from a JSON input stream.
     *
     * @param inputStream The input stream to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromJson(InputStream inputStream) throws IOException {
        return fromTree(ConfigFiles.readJson(inputStream));
    }

    /**
     * Load configuration from a YAML file.
     *
     * @param file The file to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromYaml(File file) throws IOException {
        try (InputStream inputStream = new FileInputStream(file)) {
            return fromYaml(inputStream);
        }
    }

    /**
     * Load configuration from a YAML input stream.
     *
     * @param inputStream The input stream to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromYaml(InputStream inputStream) throws IOException {
        return fromTree(ConfigFiles.readYaml(inputStream));
    }

    /**
     * Load configuration from an XML file.
     *
     * @param file The file to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromXml(File file) throws IOException {
        try (InputStream inputStream = new FileInputStream(file)) {
            return fromXml(inputStream);
        }
    }

    /**
     * Load configuration from an XML input stream.
     *
     * @param inputStream The input stream to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromXml(InputStream inputStream) throws IOException {
        return fromTree(ConfigFiles.readXml(inputStream));
    }

    /**
     * Save configuration to a JSON file.
     *
     * @param file The file to save to
     * @throws IOException if saving fails
     */
    public void toJson(File file) throws IOException {
        Files.writeString(file.toPath(), ConfigFiles.writeJson(toTree()), StandardCharsets.UTF_8);
    }

    /**
     * Save configuration to a YAML file.
     *
     * @param file The file to save to
     * @throws IOException if saving fails
     */
    public void toYaml(File file) throws IOException {
        Files.writeString(file.toPath(), ConfigFiles.writeYaml(toTree()), StandardCharsets.UTF_8);
    }

    /**
     * Save configuration to an XML file.
     *
     * @param file The file to save to
     * @throws IOException if saving fails
     */
    public void toXml(File file) throws IOException {
        Files.writeString(file.toPath(), ConfigFiles.writeXml(toTree()), StandardCharsets.UTF_8);
    }

    static EventPumpConfiguration fromTree(Object tree) throws IOException {
        if (tree == null) {
            // An empty document: nothing configured.
            return new EventPumpConfiguration();
        }
        ConfigNode object = ConfigNode.object(tree, "", "connections", "batches");
        EventPumpConfiguration configuration = new EventPumpConfiguration();
        List<Object> connections = object.list("connections");
        for (int i = 0; i < connections.size(); i++) {
            configuration.connections.add(
                ConnectionConfiguration.fromTree(connections.get(i), "connections[" + i + "]"));
        }
        List<Object> batches = object.list("batches");
        for (int i = 0; i < batches.size(); i++) {
            configuration.batches.add(BatchConfiguration.fromTree(batches.get(i), "batches[" + i + "]"));
        }
        return configuration;
    }

    Map<String, Object> toTree() {
        List<Object> connectionTrees = new ArrayList<>();
        for (ConnectionConfiguration connection : connections) {
            connectionTrees.add(connection.toTree());
        }
        List<Object> batchTrees = new ArrayList<>();
        for (BatchConfiguration batch : batches) {
            batchTrees.add(batch.toTree());
        }
        return new ConfigNode.Writer()
            .put("connections", connectionTrees)
            .put("batches", batchTrees)
            .tree();
    }

}
