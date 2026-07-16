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

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Custom deserializer for tag maps that supports both simple and extended formats.
 * <p>
 * Handles:
 * - Simple format: "tagName": "address"
 * - Extended format: "tagName": { "address": "...", "transform": "..." }
 */
public class TagMapDeserializer extends JsonDeserializer<Map<String, TagConfiguration>> {

    @Override
    public Map<String, TagConfiguration> deserialize(JsonParser p, DeserializationContext ctxt) throws IOException {
        Map<String, TagConfiguration> tags = new LinkedHashMap<>();
        ObjectMapper mapper = (ObjectMapper) p.getCodec();
        JsonNode node = mapper.readTree(p);

        Iterator<Map.Entry<String, JsonNode>> fields = node.fields();
        while (fields.hasNext()) {
            Map.Entry<String, JsonNode> entry = fields.next();
            String tagName = entry.getKey();
            JsonNode tagNode = entry.getValue();

            TagConfiguration tagConfig;
            if (tagNode.isTextual()) {
                // Simple format: just a string address
                tagConfig = new TagConfiguration(tagNode.asText());
            } else if (tagNode.isObject()) {
                // Extended format: object with address and optional transform
                tagConfig = mapper.treeToValue(tagNode, TagConfiguration.class);
            } else {
                throw new IOException("Invalid tag configuration for '" + tagName + "': expected string or object");
            }

            tags.put(tagName, tagConfig);
        }

        return tags;
    }
}
