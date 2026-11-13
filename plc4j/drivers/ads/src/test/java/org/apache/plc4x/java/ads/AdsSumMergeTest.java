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
package org.apache.plc4x.java.ads;

import org.apache.plc4x.java.ads.tag.SymbolicAdsTag;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagValueItem;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * How the per-group results of a split request are merged back into one response.
 *
 * <p>{@link AdsSumChunkTest} covers the partitioning arithmetic, and the driver testsuite
 * ({@code AdsDriverIT} over {@code protocols/ads/DriverTestsuite.xml}) covers the connection
 * issuing one sum command per group, on the read and the write path. This covers the merge
 * itself: each tag keeps its own group's result, a tag that failed before anything was sent
 * keeps that failure code, and a tag no group answered for is reported rather than dropped.
 */
public class AdsSumMergeTest {

    private static List<String> names(int count) {
        List<String> names = new ArrayList<>(count);
        for (int i = 0; i < count; i++) {
            names.add("tag" + i);
        }
        return names;
    }

    /** A real request rather than a mock: only the tag names matter here. */
    private static DefaultPlcReadRequest requestFor(List<String> tagNames) {
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        for (String name : tagNames) {
            tags.put(name, new DefaultPlcTagItem<>(SymbolicAdsTag.of("Main." + name)));
        }
        return new DefaultPlcReadRequest(null, tags);
    }

    /** The write counterpart of {@link #requestFor}: again only the tag names matter. */
    private static DefaultPlcWriteRequest writeRequestFor(List<String> tagNames) {
        LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags = new LinkedHashMap<>();
        for (String name : tagNames) {
            tags.put(name, new DefaultPlcTagValueItem<>(SymbolicAdsTag.of("Main." + name), new PlcDINT(0)));
        }
        return new DefaultPlcWriteRequest(null, tags);
    }

    private static PlcResponseItem<PlcValue> ok(int value) {
        return new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcDINT(value));
    }

    /**
     * Asserting on {@code response.getTagNames()} would prove nothing: {@link DefaultPlcReadResponse}
     * answers that from the request, so it reads back in request order whatever the merge did. What
     * can actually break is the pairing of each name with its own group's value.
     */
    @Test
    void eachTagKeepsItsOwnValueWhateverOrderTheGroupsCameBackIn() {
        List<String> tagNames = names(3);
        // The second group came back first and put its result in ahead of the first group's.
        Map<String, PlcResponseItem<PlcValue>> accumulated = new LinkedHashMap<>();
        accumulated.put("tag2", ok(22));
        accumulated.put("tag0", ok(0));
        accumulated.put("tag1", ok(11));

        PlcReadResponse response =
            AdsTcpConnection.assembleReadResponse(requestFor(tagNames), Map.of(), accumulated);

        assertEquals(0, response.getInteger("tag0"));
        assertEquals(11, response.getInteger("tag1"));
        assertEquals(22, response.getInteger("tag2"));
    }

    @Test
    void aTagThatFailedBeforeSendingKeepsItsOwnCode() {
        List<String> tagNames = names(2);
        Map<String, PlcResponseCode> initialFailures = Map.of("tag0", PlcResponseCode.INVALID_ADDRESS);
        Map<String, PlcResponseItem<PlcValue>> accumulated = new LinkedHashMap<>();
        accumulated.put("tag1", ok(1));

        PlcReadResponse response =
            AdsTcpConnection.assembleReadResponse(requestFor(tagNames), initialFailures, accumulated);

        assertEquals(PlcResponseCode.INVALID_ADDRESS, response.getResponseCode("tag0"));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag1"));
    }

    @Test
    void aTagNoGroupAnsweredForIsReportedRatherThanDropped() {
        List<String> tagNames = names(2);
        Map<String, PlcResponseItem<PlcValue>> accumulated = new LinkedHashMap<>();
        accumulated.put("tag0", ok(0));

        PlcReadResponse response =
            AdsTcpConnection.assembleReadResponse(requestFor(tagNames), Map.of(), accumulated);

        assertEquals(PlcResponseCode.INTERNAL_ERROR, response.getResponseCode("tag1"));
    }

    /** The write side pairs names with codes the same way; see the read case above on why. */
    @Test
    void eachWrittenTagKeepsItsOwnCodeWhateverOrderTheGroupsCameBackIn() {
        List<String> tagNames = names(3);
        Map<String, PlcResponseCode> accumulated = new LinkedHashMap<>();
        accumulated.put("tag2", PlcResponseCode.OK);
        accumulated.put("tag0", PlcResponseCode.INVALID_DATA);
        accumulated.put("tag1", PlcResponseCode.ACCESS_DENIED);

        PlcWriteResponse response =
            AdsTcpConnection.assembleWriteResponse(writeRequestFor(tagNames), Map.of(), accumulated);

        assertEquals(PlcResponseCode.INVALID_DATA, response.getResponseCode("tag0"));
        assertEquals(PlcResponseCode.ACCESS_DENIED, response.getResponseCode("tag1"));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag2"));
    }

    /** A tag that never made it onto the wire keeps its own code; one no group answered for is reported. */
    @Test
    void writeKeepsPreSendFailuresAndReportsUnansweredTags() {
        List<String> tagNames = names(3);
        Map<String, PlcResponseCode> initialFailures = Map.of("tag0", PlcResponseCode.INVALID_DATA);
        Map<String, PlcResponseCode> accumulated = new LinkedHashMap<>();
        accumulated.put("tag1", PlcResponseCode.OK);

        PlcWriteResponse response = AdsTcpConnection.assembleWriteResponse(
            writeRequestFor(tagNames), initialFailures, accumulated);

        assertEquals(PlcResponseCode.INVALID_DATA, response.getResponseCode("tag0"));
        assertEquals(PlcResponseCode.OK, response.getResponseCode("tag1"));
        assertEquals(PlcResponseCode.INTERNAL_ERROR, response.getResponseCode("tag2"));
    }
}
