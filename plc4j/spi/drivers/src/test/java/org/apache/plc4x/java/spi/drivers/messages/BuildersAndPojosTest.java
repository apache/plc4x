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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.OptionType;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagErrorItem;
import org.apache.plc4x.java.spi.drivers.messages.metadata.DefaultMetadata;
import org.apache.plc4x.java.spi.drivers.messages.metadata.DefaultOption;
import org.apache.plc4x.java.spi.drivers.messages.metadata.DefaultOptionMetadata;
import org.apache.plc4x.java.spi.drivers.functions.PlcReader;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.apache.plc4x.java.spi.drivers.functions.PlcWriter;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.drivers.tags.TagConfigParser;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class BuildersAndPojosTest {

    record TestTag(String address) implements PlcTag {
        @Override public String getAddressString() { return address; }
        @Override public PlcValueType getPlcValueType() { return PlcValueType.INT; }
    }

    /** A tag handler that maps any address to a {@link TestTag}; queries are not used here. */
    static class TestTagHandler implements PlcTagHandler {
        @Override public PlcTag parseTag(String tagAddress) { return new TestTag(tagAddress); }
        @Override public PlcQuery parseQuery(String query) { return null; }
    }

    private final PlcTagHandler handler = new TestTagHandler();

    // Non-null stubs - the builders require a non-null owner but build() never invokes it.
    private final PlcReader reader = req -> null;
    private final PlcWriter writer = req -> null;
    private final PlcSubscriber subscriber = new PlcSubscriber() {
        @Override public java.util.concurrent.CompletableFuture<org.apache.plc4x.java.api.messages.PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest r) { return null; }
        @Override public java.util.concurrent.CompletableFuture<org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse> unsubscribe(org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest r) { return null; }
        @Override public org.apache.plc4x.java.api.model.PlcConsumerRegistration registerConsumer(java.util.function.Consumer<org.apache.plc4x.java.api.messages.PlcSubscriptionEvent> c, java.util.Collection<org.apache.plc4x.java.api.model.PlcSubscriptionHandle> h) { return null; }
        @Override public void unregisterConsumer(org.apache.plc4x.java.api.model.PlcConsumerRegistration registration) { }
    };

    @Test
    void readRequestBuilderProducesRequest() {
        PlcReadRequest request = new DefaultPlcReadRequest.Builder(reader, handler)
            .addTagAddress("a", "addr-a")
            .addTagAddress("b", "addr-b")
            .build();
        assertEquals(2, request.getNumberOfTags());
        assertTrue(request.getTagNames().contains("a"));
        assertEquals("addr-b", request.getTag("b").getAddressString());
    }

    @Test
    void writeRequestBuilderProducesRequest() {
        PlcWriteRequest request = new DefaultPlcWriteRequest.Builder(writer, handler, new DefaultPlcValueHandler())
            .addTagAddress("a", "addr-a", 42)
            .build();
        assertEquals(1, request.getNumberOfTags());
        assertEquals(42, request.getPlcValue("a").getInt());
    }

    @Test
    void subscriptionRequestBuilderProducesRequest() {
        PlcSubscriptionRequest request = new DefaultPlcSubscriptionRequest.Builder(subscriber, handler)
            .addCyclicTagAddress("cyc", "addr-c", Duration.ofSeconds(1))
            .addChangeOfStateTagAddress("cos", "addr-d")
            .build();
        assertEquals(2, request.getNumberOfTags());
        assertTrue(request.getTagNames().contains("cyc"));
        assertTrue(request.getTagNames().contains("cos"));
    }

    @Test
    void tagErrorItemCarriesResponseCode() {
        DefaultPlcTagErrorItem item = new DefaultPlcTagErrorItem(PlcResponseCode.INVALID_ADDRESS);
        assertEquals(PlcResponseCode.INVALID_ADDRESS, item.getResponseCode());
    }

    @Test
    void optionMetadataExposesOptions() {
        DefaultOption option = new DefaultOption(
            "timeout", OptionType.LONG, "request timeout", false, 5000L, "1.0");
        assertEquals("timeout", option.getKey());
        assertEquals(OptionType.LONG, option.getType());
        assertEquals("request timeout", option.getDescription());

        DefaultOptionMetadata metadata = new DefaultOptionMetadata(List.of(option));
        assertEquals(1, metadata.getOptions().size());
    }

    @Test
    void metadataStoresAndReturnsValues() {
        DefaultMetadata metadata = new DefaultMetadata(Map.of());
        assertNotNull(metadata.keys());
        assertTrue(metadata.entries().isEmpty());
    }

    @Test
    void tagConfigParserExtractsTrailingConfig() {
        Map<String, String> config = TagConfigParser.parse("%DB1.DBW0:INT{poll-rate: 100, name: \"x\"}");
        assertEquals("100", config.get("poll-rate"));
        assertEquals("x", config.get("name"));

        // an address with no trailing config block yields an empty map
        assertTrue(TagConfigParser.parse("%DB1.DBW0:INT").isEmpty());
    }
}
