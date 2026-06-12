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

        assertEquals(target, cfg.getTargetAmsNetId());
        assertEquals(851, cfg.getTargetAmsPort());
        assertEquals(source, cfg.getSourceAmsNetId());
        assertEquals(852, cfg.getSourceAmsPort());
        assertEquals(1234, cfg.getTimeoutRequest());
        assertFalse(cfg.isLoadSymbolAndDataTypeTables());
        assertTrue(cfg.toString().contains("targetAmsPort=851"));
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
