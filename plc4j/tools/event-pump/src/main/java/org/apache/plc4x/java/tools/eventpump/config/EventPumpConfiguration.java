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

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;

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

    @JsonProperty("connections")
    private List<ConnectionConfiguration> connections = new ArrayList<>();

    @JsonProperty("batches")
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
        ObjectMapper mapper = new ObjectMapper();
        return mapper.readValue(file, EventPumpConfiguration.class);
    }

    /**
     * Load configuration from a JSON input stream.
     *
     * @param inputStream The input stream to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromJson(InputStream inputStream) throws IOException {
        ObjectMapper mapper = new ObjectMapper();
        return mapper.readValue(inputStream, EventPumpConfiguration.class);
    }

    /**
     * Load configuration from a YAML file.
     *
     * @param file The file to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromYaml(File file) throws IOException {
        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        return mapper.readValue(file, EventPumpConfiguration.class);
    }

    /**
     * Load configuration from a YAML input stream.
     *
     * @param inputStream The input stream to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromYaml(InputStream inputStream) throws IOException {
        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        return mapper.readValue(inputStream, EventPumpConfiguration.class);
    }

    /**
     * Load configuration from an XML file.
     *
     * @param file The file to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromXml(File file) throws IOException {
        XmlMapper mapper = new XmlMapper();
        return mapper.readValue(file, EventPumpConfiguration.class);
    }

    /**
     * Load configuration from an XML input stream.
     *
     * @param inputStream The input stream to load from
     * @return The configuration
     * @throws IOException if loading fails
     */
    public static EventPumpConfiguration fromXml(InputStream inputStream) throws IOException {
        XmlMapper mapper = new XmlMapper();
        return mapper.readValue(inputStream, EventPumpConfiguration.class);
    }

    /**
     * Save configuration to a JSON file.
     *
     * @param file The file to save to
     * @throws IOException if saving fails
     */
    public void toJson(File file) throws IOException {
        ObjectMapper mapper = new ObjectMapper();
        mapper.writerWithDefaultPrettyPrinter().writeValue(file, this);
    }

    /**
     * Save configuration to a YAML file.
     *
     * @param file The file to save to
     * @throws IOException if saving fails
     */
    public void toYaml(File file) throws IOException {
        ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
        mapper.writeValue(file, this);
    }

    /**
     * Save configuration to an XML file.
     *
     * @param file The file to save to
     * @throws IOException if saving fails
     */
    public void toXml(File file) throws IOException {
        XmlMapper mapper = new XmlMapper();
        mapper.writerWithDefaultPrettyPrinter().writeValue(file, this);
    }
}
