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

import java.io.IOException;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Configuration for a tag batch.
 * <p>
 * Supports both simple and extended tag formats for backward compatibility:
 * <p>
 * Simple format (YAML):
 * <pre>
 * batches:
 *   - id: batch1
 *     connectionId: plc1
 *     tags:
 *       temperature: "MAIN.temperature"
 *       pressure: "MAIN.pressure"
 * </pre>
 * <p>
 * Extended format with transformations (YAML):
 * <pre>
 * batches:
 *   - id: batch1
 *     connectionId: plc1
 *     tags:
 *       temperature:
 *         address: "MAIN.temperature"
 *         transform: "value * 1.8 + 32"
 *       pressure:
 *         address: "MAIN.pressure"
 * </pre>
 */
public class BatchConfiguration {

    private String id;

    private String connectionId;

    private Map<String, TagConfiguration> tags = new LinkedHashMap<>();

    private TriggerConfiguration trigger;

    /**
     * Watchdog bound for a single fetch cycle, in milliseconds. Null means "use the
     * default". This is not the request timeout — set that on the connection URL, e.g.
     * {@code "opcua:tcp://host:4840?request-timeout-ms=10000"}.
     */
    private Long fetchTimeoutMs;

    /**
     * Get the batch ID.
     *
     * @return The batch ID
     */
    public String getId() {
        return id;
    }

    /**
     * Set the batch ID.
     *
     * @param id The batch ID
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Get the connection ID.
     *
     * @return The connection ID
     */
    public String getConnectionId() {
        return connectionId;
    }

    /**
     * Set the connection ID.
     *
     * @param connectionId The connection ID
     */
    public void setConnectionId(String connectionId) {
        this.connectionId = connectionId;
    }

    /**
     * Get the tag configurations (name -> TagConfiguration).
     *
     * @return The tags map
     */
    public Map<String, TagConfiguration> getTags() {
        return tags;
    }

    /**
     * Set the tag configurations.
     *
     * @param tags The tags map
     */
    public void setTags(Map<String, TagConfiguration> tags) {
        this.tags = tags;
    }

    /**
     * Get a simple map of tag names to addresses (for backward compatibility).
     *
     * @return A map of tag name -> address
     */
    public Map<String, String> getSimpleTagMap() {
        Map<String, String> simpleMap = new LinkedHashMap<>();
        for (Map.Entry<String, TagConfiguration> entry : tags.entrySet()) {
            simpleMap.put(entry.getKey(), entry.getValue().getAddress());
        }
        return simpleMap;
    }

    /**
     * Set tags from a simple map (for backward compatibility and testing).
     *
     * @param simpleTags A map of tag name -> address
     */
    public void setSimpleTags(Map<String, String> simpleTags) {
        this.tags.clear();
        for (Map.Entry<String, String> entry : simpleTags.entrySet()) {
            this.tags.put(entry.getKey(), new TagConfiguration(entry.getValue()));
        }
    }

    /**
     * Get the trigger configuration.
     *
     * @return The trigger configuration
     */
    public TriggerConfiguration getTrigger() {
        return trigger;
    }

    /**
     * Set the trigger configuration.
     *
     * @param trigger The trigger configuration
     */
    public void setTrigger(TriggerConfiguration trigger) {
        this.trigger = trigger;
    }

    /**
     * Get the fetch watchdog timeout in milliseconds.
     *
     * @return The timeout, or null to use the default
     */
    public Long getFetchTimeoutMs() {
        return fetchTimeoutMs;
    }

    /**
     * Set the fetch watchdog timeout in milliseconds.
     *
     * @param fetchTimeoutMs The timeout, 0 or less to disable, or null to use the default
     */
    public void setFetchTimeoutMs(Long fetchTimeoutMs) {
        this.fetchTimeoutMs = fetchTimeoutMs;
    }

    static BatchConfiguration fromTree(Object node, String path) throws IOException {
        // "simpleTagMap" is what Jackson wrote for getSimpleTagMap(); it repeats the tags and is
        // ignored, as Jackson did. "simpleTags" is setSimpleTags().
        ConfigNode object = ConfigNode.object(node, path, "id", "connectionId", "tags", "trigger",
            "fetchTimeoutMs", "simpleTagMap", "simpleTags");
        BatchConfiguration batch = new BatchConfiguration();
        batch.setId(object.string("id"));
        batch.setConnectionId(object.string("connectionId"));
        Map<String, TagConfiguration> tags = new LinkedHashMap<>();
        for (Map.Entry<String, Object> tag : object.members("tags").entrySet()) {
            tags.put(tag.getKey(), TagConfiguration.fromTree(tag.getValue(), object.path("tags") + "." + tag.getKey()));
        }
        batch.setTags(tags);
        if (object.get("simpleTags") != null) {
            Map<String, String> simpleTags = new LinkedHashMap<>();
            for (Map.Entry<String, Object> tag : object.members("simpleTags").entrySet()) {
                simpleTags.put(tag.getKey(),
                    ConfigNode.scalarToString(tag.getValue(), object.path("simpleTags") + "." + tag.getKey()));
            }
            batch.setSimpleTags(simpleTags);
        }
        if (object.get("trigger") != null) {
            batch.setTrigger(TriggerConfiguration.fromTree(object.get("trigger"), object.path("trigger")));
        }
        batch.setFetchTimeoutMs(object.longValue("fetchTimeoutMs"));
        return batch;
    }

    Map<String, Object> toTree() {
        Map<String, Object> tagTree = new LinkedHashMap<>();
        for (Map.Entry<String, TagConfiguration> tag : tags.entrySet()) {
            tagTree.put(tag.getKey(), tag.getValue().toTree());
        }
        return new ConfigNode.Writer()
            .put("id", id)
            .put("connectionId", connectionId)
            .put("tags", tagTree)
            .put("trigger", trigger != null ? trigger.toTree() : null)
            .put("fetchTimeoutMs", fetchTimeoutMs)
            .tree();
    }

}
