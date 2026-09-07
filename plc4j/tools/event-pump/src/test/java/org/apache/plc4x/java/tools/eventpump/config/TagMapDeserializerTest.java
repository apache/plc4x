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
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class TagMapDeserializerTest {

    private final ObjectMapper mapper = new ObjectMapper();

    @Test
    void testDeserializeSimpleFormat() throws Exception {
        String json = """
            {
                "temperature": "MAIN.temperature",
                "pressure": "MAIN.pressure"
            }
            """;

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();
        Map<String, TagConfiguration> result = deserializer.deserialize(parser, mapper.getDeserializationContext());

        assertNotNull(result);
        assertEquals(2, result.size());

        TagConfiguration tempConfig = result.get("temperature");
        assertNotNull(tempConfig);
        assertEquals("MAIN.temperature", tempConfig.getAddress());
        assertFalse(tempConfig.hasTransform());

        TagConfiguration pressConfig = result.get("pressure");
        assertNotNull(pressConfig);
        assertEquals("MAIN.pressure", pressConfig.getAddress());
        assertFalse(pressConfig.hasTransform());
    }

    @Test
    void testDeserializeExtendedFormat() throws Exception {
        String json = """
            {
                "temperature": {
                    "address": "MAIN.temperature",
                    "transform": "value * 1.8 + 32"
                }
            }
            """;

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();
        Map<String, TagConfiguration> result = deserializer.deserialize(parser, mapper.getDeserializationContext());

        assertNotNull(result);
        assertEquals(1, result.size());

        TagConfiguration tempConfig = result.get("temperature");
        assertNotNull(tempConfig);
        assertEquals("MAIN.temperature", tempConfig.getAddress());
        assertEquals("value * 1.8 + 32", tempConfig.getTransform());
        assertTrue(tempConfig.hasTransform());
    }

    @Test
    void testDeserializeMixedFormat() throws Exception {
        String json = """
            {
                "temperature": {
                    "address": "MAIN.temperature",
                    "transform": "value * 1.8 + 32"
                },
                "pressure": "MAIN.pressure",
                "humidity": {
                    "address": "MAIN.humidity"
                }
            }
            """;

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();
        Map<String, TagConfiguration> result = deserializer.deserialize(parser, mapper.getDeserializationContext());

        assertNotNull(result);
        assertEquals(3, result.size());

        // Temperature with transform
        TagConfiguration tempConfig = result.get("temperature");
        assertNotNull(tempConfig);
        assertEquals("MAIN.temperature", tempConfig.getAddress());
        assertEquals("value * 1.8 + 32", tempConfig.getTransform());
        assertTrue(tempConfig.hasTransform());

        // Pressure simple format
        TagConfiguration pressConfig = result.get("pressure");
        assertNotNull(pressConfig);
        assertEquals("MAIN.pressure", pressConfig.getAddress());
        assertFalse(pressConfig.hasTransform());

        // Humidity extended format without transform
        TagConfiguration humidConfig = result.get("humidity");
        assertNotNull(humidConfig);
        assertEquals("MAIN.humidity", humidConfig.getAddress());
        assertFalse(humidConfig.hasTransform());
    }

    @Test
    void testDeserializeInvalidFormat() throws Exception {
        String json = """
            {
                "temperature": 123
            }
            """;

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();

        assertThrows(Exception.class, () -> {
            deserializer.deserialize(parser, mapper.getDeserializationContext());
        });
    }

    @Test
    void testDeserializeEmptyMap() throws Exception {
        String json = "{}";

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();
        Map<String, TagConfiguration> result = deserializer.deserialize(parser, mapper.getDeserializationContext());

        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    void testDeserializeExtendedFormatWithoutTransform() throws Exception {
        String json = """
            {
                "temperature": {
                    "address": "MAIN.temperature"
                }
            }
            """;

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();
        Map<String, TagConfiguration> result = deserializer.deserialize(parser, mapper.getDeserializationContext());

        assertNotNull(result);
        assertEquals(1, result.size());

        TagConfiguration tempConfig = result.get("temperature");
        assertNotNull(tempConfig);
        assertEquals("MAIN.temperature", tempConfig.getAddress());
        assertFalse(tempConfig.hasTransform());
    }

    @Test
    void testDeserializePreservesOrder() throws Exception {
        String json = """
            {
                "tag1": "MAIN.tag1",
                "tag2": "MAIN.tag2",
                "tag3": "MAIN.tag3"
            }
            """;

        JsonParser parser = mapper.getFactory().createParser(json);
        parser.nextToken(); // Move to start of object

        TagMapDeserializer deserializer = new TagMapDeserializer();
        Map<String, TagConfiguration> result = deserializer.deserialize(parser, mapper.getDeserializationContext());

        assertNotNull(result);
        assertEquals(3, result.size());

        // LinkedHashMap should preserve insertion order
        var keys = result.keySet().toArray(new String[0]);
        assertEquals("tag1", keys[0]);
        assertEquals("tag2", keys[1]);
        assertEquals("tag3", keys[2]);
    }
}
