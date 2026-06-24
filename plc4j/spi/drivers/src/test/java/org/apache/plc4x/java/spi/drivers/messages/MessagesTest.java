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

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class MessagesTest {

    /** Minimal {@link PlcTag} - the interface only requires an address string. */
    record TestTag(String address, PlcValueType valueType) implements PlcTag {
        @Override public String getAddressString() { return address; }
        @Override public PlcValueType getPlcValueType() { return valueType; }
    }

    private static LinkedHashMap<String, PlcTagItem<PlcTag>> tagItems(String name, PlcValueType type) {
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put(name, new DefaultPlcTagItem<>(new TestTag(name + "-addr", type)));
        return tags;
    }

    private static Map<String, PlcResponseItem<PlcValue>> ok(String name, PlcValue value) {
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
        values.put(name, new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
        return values;
    }

    @Test
    void readRequestExposesTags() {
        DefaultPlcReadRequest request = new DefaultPlcReadRequest(null, tagItems("a", PlcValueType.INT));
        assertEquals(1, request.getNumberOfTags());
        assertTrue(request.getTagNames().contains("a"));
        assertEquals(1, request.getTags().size());
        assertEquals("a-addr", request.getTag("a").getAddressString());
    }

    @Test
    void readResponseNumericGettersDelegateToValue() {
        DefaultPlcReadRequest request = new DefaultPlcReadRequest(null, tagItems("a", PlcValueType.INT));
        DefaultPlcReadResponse response = new DefaultPlcReadResponse(request, ok("a", new PlcINT(42)));

        assertEquals(PlcResponseCode.OK, response.getResponseCode("a"));
        assertEquals(42, response.getInteger("a"));
        assertEquals(42L, response.getLong("a"));
        assertEquals((short) 42, response.getShort("a"));
        assertEquals((byte) 42, response.getByte("a"));
        assertEquals(42.0f, response.getFloat("a"));
        assertEquals(42.0, response.getDouble("a"));
        assertEquals(java.math.BigInteger.valueOf(42), response.getBigInteger("a"));
        assertEquals(42, ((Number) response.getObject("a")).intValue());
        assertEquals("a-addr", response.getTag("a").getAddressString());
        assertNotNull(response.getPlcValue("a"));
        assertEquals(1, response.getNumberOfValues("a"));
        assertEquals(1, response.getTagNames().size());
        assertNotNull(response.getRequest());

        // isValid* must never throw and reflect the value type
        assertTrue(response.isValidInteger("a"));
        assertFalse(response.isValidDate("a"));
        assertFalse(response.isValidTime("a"));
        assertFalse(response.isValidDateTime("a"));

        // bulk accessors
        assertEquals(1, response.getAllIntegers("a").size());
        assertEquals(1, response.getAllObjects("a").size());
    }

    @Test
    void readResponseStringAndBooleanAndReal() {
        DefaultPlcReadRequest request = new DefaultPlcReadRequest(null, tagItems("s", PlcValueType.STRING));
        DefaultPlcReadResponse strResp = new DefaultPlcReadResponse(request, ok("s", new PlcSTRING("hi")));
        assertEquals("hi", strResp.getString("s"));
        assertTrue(strResp.isValidString("s"));
        assertEquals("hi", strResp.getObject("s"));

        DefaultPlcReadResponse boolResp = new DefaultPlcReadResponse(
            new DefaultPlcReadRequest(null, tagItems("b", PlcValueType.BOOL)), ok("b", new PlcBOOL(true)));
        assertTrue(boolResp.getBoolean("b"));
        assertTrue(boolResp.isValidBoolean("b"));

        DefaultPlcReadResponse realResp = new DefaultPlcReadResponse(
            new DefaultPlcReadRequest(null, tagItems("r", PlcValueType.INT)), ok("r", new PlcREAL(1.5f)));
        assertEquals(1.5f, realResp.getFloat("r"));
        assertEquals(1.5, realResp.getDouble("r"), 0.0001);
    }

    @Test
    void readResponseReportsErrorCodePerTag() {
        DefaultPlcReadRequest request = new DefaultPlcReadRequest(null, tagItems("a", PlcValueType.INT));
        Map<String, PlcResponseItem<PlcValue>> values = new LinkedHashMap<>();
        values.put("a", new DefaultPlcResponseItem<>(PlcResponseCode.NOT_FOUND, null));
        DefaultPlcReadResponse response = new DefaultPlcReadResponse(request, values);
        assertEquals(PlcResponseCode.NOT_FOUND, response.getResponseCode("a"));
    }

    @Test
    void subscriptionEventExposesTimestampAndValues() {
        Instant now = Instant.ofEpochMilli(1_700_000_000_000L);
        DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(now, ok("a", new PlcINT(7)));
        assertEquals(now, event.getTimestamp());
        assertEquals(PlcResponseCode.OK, event.getResponseCode("a"));
        assertEquals(7, event.getInteger("a"));
        assertEquals("7", event.getString("a"));
        assertTrue(event.isValidInteger("a"));
        assertTrue(event.getTagNames().contains("a"));
        assertNotNull(event.getPlcValue("a"));
    }

    @Test
    void writeRequestAndResponse() {
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        tags.put("a", new DefaultPlcTagValueItem<>(new TestTag("a-addr", PlcValueType.INT), new PlcINT(5)));
        DefaultPlcWriteRequest request = new DefaultPlcWriteRequest(null, tags);
        assertEquals(1, request.getNumberOfTags());
        assertEquals(5, request.getPlcValue("a").getInt());
        assertEquals("a-addr", request.getTag("a").getAddressString());

        DefaultPlcWriteResponse response = new DefaultPlcWriteResponse(request, Map.of("a", PlcResponseCode.OK));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("a"));
        assertTrue(response.getTagNames().contains("a"));
        assertNotNull(response.getRequest());
    }

    @Test
    void responseAndTagValueItems() {
        DefaultPlcResponseItem<PlcValue> item = new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcINT(1));
        assertEquals(PlcResponseCode.OK, item.getResponseCode());
        assertEquals(1, item.getValue().getInt());

        DefaultPlcTagItem<PlcTag> tagItem = new DefaultPlcTagItem<>(new TestTag("x", PlcValueType.INT));
        assertEquals("x", tagItem.getTag().getAddressString());

        DefaultPlcTagValueItem<PlcTag> valueItem =
            new DefaultPlcTagValueItem<>(new TestTag("y", PlcValueType.INT), new PlcINT(9));
        assertEquals("y", valueItem.getTag().getAddressString());
        assertEquals(9, valueItem.getValue().getInt());
    }

    @Test
    void subscriptionTagAndArrayInfo() {
        DefaultPlcSubscriptionTag tag = new DefaultPlcSubscriptionTag(
            PlcSubscriptionType.CYCLIC, new TestTag("a", PlcValueType.INT), Duration.ofSeconds(1));
        assertEquals(PlcSubscriptionType.CYCLIC, tag.getPlcSubscriptionType());
        assertTrue(tag.getDuration().isPresent());
        assertEquals("a", tag.getTag().getAddressString());

        DefaultArrayInfo arrayInfo = new DefaultArrayInfo(0, 9);
        assertEquals(0, arrayInfo.getLowerBound());
        assertEquals(9, arrayInfo.getUpperBound());
        assertEquals(10, arrayInfo.getSize());
    }
}
