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

package org.apache.plc4x.java.ads.configuration;

import org.apache.plc4x.java.ads.readwrite.AmsNetId;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class AdsConfigurationTest {

    @Test
    void gettersAndSetters() {
        AdsConfiguration cfg = new AdsConfiguration();
        AmsNetId target = AdsConfiguration.amsNetIdOf("1.2.3.4.5.6");
        AmsNetId source = AdsConfiguration.amsNetIdOf("7.8.9.10.11.12");
        cfg.setTargetAmsNetId(target);
        cfg.setTargetAmsPort(851);
        cfg.setSourceAmsNetId(source);
        cfg.setSourceAmsPort(852);
        cfg.setTimeoutRequest(1234);
        cfg.setLoadSymbolAndDataTypeTables(false);
        cfg.setMaxDataTypeTableDepth(64);
        cfg.setMaxSumCommandItems(2);

        assertEquals(target, cfg.getTargetAmsNetId());
        assertEquals(851, cfg.getTargetAmsPort());
        assertEquals(source, cfg.getSourceAmsNetId());
        assertEquals(852, cfg.getSourceAmsPort());
        assertEquals(1234, cfg.getTimeoutRequest());
        assertFalse(cfg.isLoadSymbolAndDataTypeTables());
        assertEquals(64, cfg.getMaxDataTypeTableDepth());
        assertEquals(2, cfg.getMaxSumCommandItems());
        assertTrue(cfg.toString().contains("targetAmsPort=851"));
    }

    /** The parameters AdsConfiguration insists on, so a test can vary only what it cares about. */
    private static final String REQUIRED_PARAMS =
        "target-ams-net-id=1.2.3.4.5.6&target-ams-port=851"
            + "&source-ams-net-id=7.8.9.10.11.12&source-ams-port=852";

    /**
     * The depth budget is handed straight to the table parser, and the parser rejects an entry once
     * the budget is spent. A default that failed to resolve would therefore arrive as 0 and reject
     * the very first entry of every upload - so the default is pinned here rather than assumed.
     */
    @Test
    void maxDataTypeTableDepthDefaultsToTwenty() throws Exception {
        AdsConfiguration defaults = new ConfigurationFactory()
            .createConfiguration(AdsConfiguration.class, REQUIRED_PARAMS);

        assertEquals(20, defaults.getMaxDataTypeTableDepth(),
            "an unconfigured connection must get a usable budget, never 0");
    }

    /**
     * The cap is Beckhoff's recommended ceiling for sub-ADS calls in one sum command. A default
     * that failed to resolve would arrive as 0, and chunk() rejects a non-positive chunk size, so
     * every sum read and write would fail outright - the default is pinned here rather than
     * assumed.
     */
    @Test
    void maxSumCommandItemsDefaultsToFiveHundred() throws Exception {
        AdsConfiguration defaults = new ConfigurationFactory()
            .createConfiguration(AdsConfiguration.class, REQUIRED_PARAMS);

        assertEquals(500, defaults.getMaxSumCommandItems(),
            "Beckhoff recommends at most 500 sub-ADS calls per sum command");
    }

    /**
     * Lowering the cap is what makes the split testable end-to-end: the driver testsuite drives
     * three tags with a cap of two rather than putting 501 tags in XML.
     */
    @Test
    void maxSumCommandItemsCanBeLowered() throws Exception {
        AdsConfiguration lowered = new ConfigurationFactory().createConfiguration(
            AdsConfiguration.class, REQUIRED_PARAMS + "&max-sum-command-items=2");

        assertEquals(2, lowered.getMaxSumCommandItems());
    }

    @Test
    void maxDataTypeTableDepthCanBeRaised() throws Exception {
        AdsConfiguration raised = new ConfigurationFactory().createConfiguration(
            AdsConfiguration.class, REQUIRED_PARAMS + "&max-data-type-table-depth=64");

        assertEquals(64, raised.getMaxDataTypeTableDepth(),
            "a device known to need deeper nesting must be accommodated without a patched jar");
    }

    @Test
    void amsNetIdOfParsesValid() {
        AmsNetId id = AdsConfiguration.amsNetIdOf("192.168.0.1.1.1");
        assertNotNull(id);
    }

    @Test
    void amsNetIdOfRejectsInvalid() {
        assertThrows(IllegalArgumentException.class, () -> AdsConfiguration.amsNetIdOf("1.2.3"));
    }

    @Test
    void converter() {
        AdsConfiguration.AmsNetIdConverter converter = new AdsConfiguration.AmsNetIdConverter();
        assertEquals(AmsNetId.class, converter.getType());
        assertNotNull(converter.convert("1.2.3.4.5.6"));
    }
}
