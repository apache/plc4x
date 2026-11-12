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

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The three configuration file formats. The {@code JACKSON_*} constants are files written by the
 * Jackson-based implementation this replaced, so files people already have keep loading.
 */
class ConfigFilesTest {

    // Written by EventPumpConfiguration.toJson/toYaml/toXml before Jackson was removed.
    private static final String JACKSON_JSON = """
        {
          "connections" : [ {
            "id" : "plc1",
            "url" : "opcua:tcp://1.2.3.4:4840?a=1&b=2",
            "username" : "u",
            "password" : "p<&>\\""
          }, {
            "id" : "plc2",
            "url" : "s7://1.2.3.5",
            "username" : null,
            "password" : null
          } ],
          "batches" : [ {
            "simpleTagMap" : {
              "temperature" : "ns=2;i=1001",
              "pressure" : "ns=2;i=1002"
            },
            "id" : "boiler",
            "connectionId" : "plc1",
            "tags" : {
              "temperature" : {
                "address" : "ns=2;i=1001",
                "transform" : "value * 1.8 + 32"
              },
              "pressure" : {
                "address" : "ns=2;i=1002",
                "transform" : null
              }
            },
            "trigger" : {
              "type" : "timer",
              "intervalSeconds" : null,
              "intervalMillis" : 250,
              "initialDelaySeconds" : 2,
              "initialDelayMillis" : null,
              "tagName" : null,
              "tagAddress" : null,
              "condition" : null
            },
            "fetchTimeoutMs" : 1500
          }, {
            "simpleTagMap" : {
              "x" : "%DB1:0:INT"
            },
            "id" : "events",
            "connectionId" : "plc2",
            "tags" : {
              "x" : {
                "address" : "%DB1:0:INT",
                "transform" : null
              }
            },
            "trigger" : {
              "type" : "subscription",
              "intervalSeconds" : null,
              "intervalMillis" : null,
              "initialDelaySeconds" : null,
              "initialDelayMillis" : null,
              "tagName" : "x",
              "tagAddress" : "%DB1:0:INT",
              "condition" : "value > 3"
            },
            "fetchTimeoutMs" : null
          } ]
        }
        """;

    private static final String JACKSON_YAML = """
        ---
        connections:
        - id: "plc1"
          url: "opcua:tcp://1.2.3.4:4840?a=1&b=2"
          username: "u"
          password: "p<&>\\""
        - id: "plc2"
          url: "s7://1.2.3.5"
          username: null
          password: null
        batches:
        - simpleTagMap:
            temperature: "ns=2;i=1001"
            pressure: "ns=2;i=1002"
          id: "boiler"
          connectionId: "plc1"
          tags:
            temperature:
              address: "ns=2;i=1001"
              transform: "value * 1.8 + 32"
            pressure:
              address: "ns=2;i=1002"
              transform: null
          trigger:
            type: "timer"
            intervalSeconds: null
            intervalMillis: 250
            initialDelaySeconds: 2
            initialDelayMillis: null
            tagName: null
            tagAddress: null
            condition: null
          fetchTimeoutMs: 1500
        - simpleTagMap:
            x: "%DB1:0:INT"
          id: "events"
          connectionId: "plc2"
          tags:
            x:
              address: "%DB1:0:INT"
              transform: null
          trigger:
            type: "subscription"
            intervalSeconds: null
            intervalMillis: null
            initialDelaySeconds: null
            initialDelayMillis: null
            tagName: "x"
            tagAddress: "%DB1:0:INT"
            condition: "value > 3"
          fetchTimeoutMs: null
        """;

    private static final String JACKSON_XML = """
        <EventPumpConfiguration>
          <connections>
            <connections>
              <id>plc1</id>
              <url>opcua:tcp://1.2.3.4:4840?a=1&amp;b=2</url>
              <username>u</username>
              <password>p&lt;&amp;>"</password>
            </connections>
            <connections>
              <id>plc2</id>
              <url>s7://1.2.3.5</url>
              <username/>
              <password/>
            </connections>
          </connections>
          <batches>
            <batches>
              <simpleTagMap>
                <temperature>ns=2;i=1001</temperature>
                <pressure>ns=2;i=1002</pressure>
              </simpleTagMap>
              <id>boiler</id>
              <connectionId>plc1</connectionId>
              <tags>
                <temperature>
                  <address>ns=2;i=1001</address>
                  <transform>value * 1.8 + 32</transform>
                </temperature>
                <pressure>
                  <address>ns=2;i=1002</address>
                  <transform/>
                </pressure>
              </tags>
              <trigger>
                <type>timer</type>
                <intervalSeconds/>
                <intervalMillis>250</intervalMillis>
                <initialDelaySeconds>2</initialDelaySeconds>
                <initialDelayMillis/>
                <tagName/>
                <tagAddress/>
                <condition/>
              </trigger>
              <fetchTimeoutMs>1500</fetchTimeoutMs>
            </batches>
            <batches>
              <simpleTagMap>
                <x>%DB1:0:INT</x>
              </simpleTagMap>
              <id>events</id>
              <connectionId>plc2</connectionId>
              <tags>
                <x>
                  <address>%DB1:0:INT</address>
                  <transform/>
                </x>
              </tags>
              <trigger>
                <type>subscription</type>
                <intervalSeconds/>
                <intervalMillis/>
                <initialDelaySeconds/>
                <initialDelayMillis/>
                <tagName>x</tagName>
                <tagAddress>%DB1:0:INT</tagAddress>
                <condition>value > 3</condition>
              </trigger>
              <fetchTimeoutMs/>
            </batches>
          </batches>
        </EventPumpConfiguration>
        """;

    private static final String JACKSON_XML_SINGLE_CONNECTION = """
        <EventPumpConfiguration>
          <connections>
            <connections>
              <id>plc2</id>
              <url>s7://1.2.3.5</url>
              <username/>
              <password/>
            </connections>
          </connections>
          <batches/>
        </EventPumpConfiguration>
        """;

    private static final String JACKSON_XML_EMPTY = """
        <EventPumpConfiguration>
          <connections/>
          <batches/>
        </EventPumpConfiguration>
        """;

    private static InputStream stream(String text) {
        return new ByteArrayInputStream(text.getBytes(StandardCharsets.UTF_8));
    }

    /** The configuration all JACKSON_* "full" files describe. */
    private static void assertFull(EventPumpConfiguration config) {
        assertEquals(2, config.getConnections().size());
        ConnectionConfiguration plc1 = config.getConnections().get(0);
        assertEquals("plc1", plc1.getId());
        assertEquals("opcua:tcp://1.2.3.4:4840?a=1&b=2", plc1.getUrl());
        assertEquals("u", plc1.getUsername());
        assertEquals("p<&>\"", plc1.getPassword());
        ConnectionConfiguration plc2 = config.getConnections().get(1);
        assertEquals("plc2", plc2.getId());
        // Jackson's XML reader made "" of these; nothing is nothing.
        assertNull(plc2.getUsername());
        assertNull(plc2.getPassword());

        assertEquals(2, config.getBatches().size());
        BatchConfiguration boiler = config.getBatches().get(0);
        assertEquals("boiler", boiler.getId());
        assertEquals("plc1", boiler.getConnectionId());
        assertEquals(1500L, boiler.getFetchTimeoutMs());
        assertEquals(List.of("temperature", "pressure"), List.copyOf(boiler.getTags().keySet()));
        assertEquals(new TagConfiguration("ns=2;i=1001", "value * 1.8 + 32"), boiler.getTags().get("temperature"));
        assertEquals("ns=2;i=1002", boiler.getTags().get("pressure").getAddress());
        assertFalse(boiler.getTags().get("pressure").hasTransform());
        TriggerConfiguration timer = boiler.getTrigger();
        assertEquals("timer", timer.getType());
        assertEquals(250L, timer.getIntervalMillis());
        assertEquals(2L, timer.getInitialDelaySeconds());
        assertNull(timer.getIntervalSeconds());
        assertNull(timer.getInitialDelayMillis());
        assertNull(timer.getTagName());

        BatchConfiguration events = config.getBatches().get(1);
        assertNull(events.getFetchTimeoutMs());
        TriggerConfiguration subscription = events.getTrigger();
        assertEquals("subscription", subscription.getType());
        assertEquals("x", subscription.getTagName());
        assertEquals("%DB1:0:INT", subscription.getTagAddress());
        assertEquals("value > 3", subscription.getCondition());
    }

    @Test
    void readsWhatJacksonWrote() throws IOException {
        assertFull(EventPumpConfiguration.fromJson(stream(JACKSON_JSON)));
        assertFull(EventPumpConfiguration.fromYaml(stream(JACKSON_YAML)));
        assertFull(EventPumpConfiguration.fromXml(stream(JACKSON_XML)));

        EventPumpConfiguration single = EventPumpConfiguration.fromXml(stream(JACKSON_XML_SINGLE_CONNECTION));
        assertEquals(1, single.getConnections().size());
        assertEquals("plc2", single.getConnections().get(0).getId());
        assertTrue(single.getBatches().isEmpty());

        EventPumpConfiguration empty = EventPumpConfiguration.fromXml(stream(JACKSON_XML_EMPTY));
        assertTrue(empty.getConnections().isEmpty());
        assertTrue(empty.getBatches().isEmpty());
    }

    @Test
    void writesWhatItReads(@TempDir Path dir) throws IOException {
        EventPumpConfiguration original = EventPumpConfiguration.fromJson(stream(JACKSON_JSON));
        File json = dir.resolve("c.json").toFile();
        File yaml = dir.resolve("c.yaml").toFile();
        File xml = dir.resolve("c.xml").toFile();
        original.toJson(json);
        original.toYaml(yaml);
        original.toXml(xml);
        assertFull(EventPumpConfiguration.fromJson(json));
        assertFull(EventPumpConfiguration.fromYaml(yaml));
        assertFull(EventPumpConfiguration.fromXml(xml));
        // Nulls and the derived simpleTagMap are left out; a tag without transform is just its address.
        String written = Files.readString(json.toPath());
        assertFalse(written.contains("null"), written);
        assertFalse(written.contains("simpleTagMap"), written);
        assertTrue(written.contains("\"pressure\": \"ns=2;i=1002\""), written);
        assertTrue(Files.readString(xml.toPath()).contains("<password>p&lt;&amp;&gt;\"</password>"));
    }

    @Test
    void readsTheDocumentedYaml() throws IOException {
        EventPumpConfiguration config = EventPumpConfiguration.fromYaml(stream("""
            connections:
              - id: plc1
                url: "opcua:tcp://192.168.1.1:4840?request-timeout-ms=10000"

            batches:
              - id: boiler
                connectionId: plc1
                tags:
                  temperature:
                    address: "ns=2;i=1001"
                    transform: "value * 1.8 + 32"
                  pressure: "ns=2;i=1002"
                  register: 40001
                trigger:
                  type: timer
                  intervalSeconds: 5
            """));
        BatchConfiguration batch = config.getBatches().get(0);
        assertEquals("value * 1.8 + 32", batch.getTags().get("temperature").getTransform());
        assertEquals("ns=2;i=1002", batch.getTags().get("pressure").getAddress());
        // YAML reads a bare number as a number; an address is text all the same.
        assertEquals("40001", batch.getTags().get("register").getAddress());
        assertEquals(5L, batch.getTrigger().getIntervalSeconds());
    }

    @Test
    void readsSimpleTagsAndNumbersAsText() throws IOException {
        EventPumpConfiguration config = EventPumpConfiguration.fromJson(stream("""
            {"batches": [{"id": "b", "simpleTags": {"a": "A", "b": 2},
              "trigger": {"type": "timer", "intervalMillis": "250", "tagName": true}}]}
            """));
        BatchConfiguration batch = config.getBatches().get(0);
        assertEquals("A", batch.getTags().get("a").getAddress());
        assertEquals("2", batch.getTags().get("b").getAddress());
        assertEquals(250L, batch.getTrigger().getIntervalMillis());
        assertEquals("true", batch.getTrigger().getTagName());
    }

    @Test
    void anEmptyFileIsAnEmptyConfiguration() throws IOException {
        EventPumpConfiguration config = EventPumpConfiguration.fromYaml(stream(""));
        assertTrue(config.getConnections().isEmpty());
        assertTrue(config.getBatches().isEmpty());
    }

    @Test
    void mistakesAreReportedWithTheirPlace() {
        assertError("batches[0].trigger: unknown property 'intervalMilis'",
            "{\"batches\": [{\"trigger\": {\"intervalMilis\": 5}}]}");
        assertError("batches[0].trigger.intervalMillis: expected a whole number, got 'soon'",
            "{\"batches\": [{\"trigger\": {\"intervalMillis\": \"soon\"}}]}");
        assertError("batches[0].trigger.intervalMillis: expected a whole number, got '2.5'",
            "{\"batches\": [{\"trigger\": {\"intervalMillis\": 2.5}}]}");
        assertError("batches[0].fetchTimeoutMs: expected a whole number, got '99999999999999999999'",
            "{\"batches\": [{\"fetchTimeoutMs\": 99999999999999999999}]}");
        assertError("batches[0].tags.t: expected an address or an object with an address, got nothing",
            "{\"batches\": [{\"tags\": {\"t\": null}}]}");
        assertError("batches[0].tags.t: unknown property 'adress'",
            "{\"batches\": [{\"tags\": {\"t\": {\"adress\": \"x\"}}}]}");
        assertError("batches[0].tags: expected an object, got a list",
            "{\"batches\": [{\"tags\": [\"x\"]}]}");
        assertError("connections: expected a list, got an object",
            "{\"connections\": {\"connections\": {\"id\": \"x\"}}}");
        assertError("connections[0].id: expected a text, got a list",
            "{\"connections\": [{\"id\": [\"x\"]}]}");
        assertError("connections[0]: expected an object, got 'x'", "{\"connections\": [\"x\"]}");
        assertError("expected an object, got a list", "[]");
    }

    private static void assertError(String expected, String json) {
        IOException e = assertThrows(IOException.class, () -> EventPumpConfiguration.fromJson(stream(json)));
        assertTrue(e.getMessage().startsWith(expected), e.getMessage());
    }

    @Test
    void longValuesAcceptAllIntegerTypes() throws IOException {
        Map<String, Object> trigger = new LinkedHashMap<>();
        trigger.put("intervalSeconds", 1);
        trigger.put("intervalMillis", (short) 2);
        trigger.put("initialDelaySeconds", (byte) 3);
        trigger.put("initialDelayMillis", java.math.BigInteger.valueOf(4));
        TriggerConfiguration config = TriggerConfiguration.fromTree(trigger, "t");
        assertEquals(1L, config.getIntervalSeconds());
        assertEquals(2L, config.getIntervalMillis());
        assertEquals(3L, config.getInitialDelaySeconds());
        assertEquals(4L, config.getInitialDelayMillis());
    }

    @Test
    void invalidFilesAreIOExceptions() {
        assertThrows(IOException.class, () -> EventPumpConfiguration.fromJson(stream("{")));
        assertThrows(IOException.class, () -> EventPumpConfiguration.fromYaml(stream("a: [")));
        assertThrows(IOException.class, () -> EventPumpConfiguration.fromXml(stream("<a>")));
    }

    @Test
    void xmlWithADoctypeIsRefused() {
        IOException e = assertThrows(IOException.class, () -> EventPumpConfiguration.fromXml(stream("""
            <?xml version="1.0"?>
            <!DOCTYPE c [<!ENTITY x SYSTEM "file:///etc/passwd">]>
            <EventPumpConfiguration><connections><connections><id>&x;</id></connections></connections></EventPumpConfiguration>
            """)));
        assertTrue(e.getMessage().contains("DOCTYPE"), e.getMessage());
    }

    @Test
    void xmlIgnoresTheRootNameAndNamespaces() throws IOException {
        EventPumpConfiguration config = EventPumpConfiguration.fromXml(stream("""
            <config xmlns="urn:x"><connections><connection><id>a</id></connection><connection><id>b</id></connection></connections></config>
            """));
        assertEquals(List.of("a", "b"), config.getConnections().stream().map(ConnectionConfiguration::getId).toList());
    }

    @Test
    void namesAndTextsXmlCantHoldAreRefused(@TempDir Path dir) {
        EventPumpConfiguration badName = new EventPumpConfiguration();
        BatchConfiguration batch = new BatchConfiguration();
        batch.setTags(new LinkedHashMap<>(Map.of("temp 1", new TagConfiguration("x"))));
        badName.getBatches().add(batch);
        assertThrows(IOException.class, () -> badName.toXml(dir.resolve("a.xml").toFile()));

        EventPumpConfiguration badText = new EventPumpConfiguration();
        ConnectionConfiguration connection = new ConnectionConfiguration();
        connection.setId("a\u0001");
        badText.getConnections().add(connection);
        assertThrows(IOException.class, () -> badText.toXml(dir.resolve("b.xml").toFile()));
    }

    @Test
    void xmlKeepsTabsAndNewlinesAndWritesEmptyObjects() throws IOException {
        Map<String, Object> tree = new LinkedHashMap<>();
        tree.put("text", "a\tb\nc");
        tree.put("empty", new LinkedHashMap<>());
        tree.put("none", null);
        String xml = ConfigFiles.writeXml(tree);
        assertTrue(xml.contains("<text>a\tb\nc</text>"), xml);
        assertTrue(xml.contains("<empty/>"), xml);
        assertTrue(xml.contains("<none/>"), xml);
    }

}
