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
package org.apache.plc4x.java.s7.configuration;

import org.apache.plc4x.java.s7.readwrite.DeviceGroup;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * These tests deliberately exercise the real configuration path - parsing a connection-string
 * parameter map via {@link ConfigurationFactory}, which injects field values directly via
 * reflection (it never calls setters). The previous version of these tests called setters
 * directly and so never covered the path that actually runs in production; that is how
 * issue #2620 (remote-slot silently dropped, connection landing on slot 0) slipped through.
 */
class S7CotpTransportConfigurationTest {

    private S7CotpTransportConfiguration parse(String params) {
        return new ConfigurationFactory()
            .createConfiguration(S7CotpTransportConfiguration.class, params);
    }

    @Test
    void defaultPort() {
        assertEquals(102, new S7CotpTransportConfiguration().getDefaultPort());
    }

    @Test
    void defaultsDeriveLegacyTsaps() {
        // No addressing parameters: local defaults to OTHERS rack=1/slot=1 -> 0x0311,
        // remote defaults to PG_OR_PC rack=0/slot=0 -> 0x0100. These match the pre-SPI3 wire bytes.
        S7CotpTransportConfiguration cfg = parse("");
        assertEquals(0x0311, cfg.getLocalTsap());
        assertEquals(0x0100, cfg.getRemoteTsap());
    }

    @Test
    void remoteSlotIsHonoured() {
        // Regression test for #2620: remote-slot=3 must reach the called TSAP as 0x0103.
        S7CotpTransportConfiguration cfg = parse("remote-rack=0&remote-slot=3&controller-type=S7_400");
        assertEquals(0x0103, cfg.getRemoteTsap());
    }

    @Test
    void remoteRackAndSlotAreEncoded() {
        // (rack << 4) | slot in the low byte, PG_OR_PC (0x01) in the high byte.
        S7CotpTransportConfiguration cfg = parse("remote-rack=2&remote-slot=5");
        assertEquals(0x0125, cfg.getRemoteTsap());
    }

    @Test
    void localRackSlotAndDeviceGroupAreEncoded() {
        S7CotpTransportConfiguration cfg = parse("local-rack=2&local-slot=3&local-device-group=OS");
        assertEquals(DeviceGroup.OS, cfg.getLocalDeviceGroup());
        // OS (0x02) high byte, (rack=2 << 4) | slot=3 = 0x23 low byte.
        assertEquals(0x0223, cfg.getLocalTsap());
    }

    @Test
    void explicitTsapOverridesDerivedValue() {
        S7CotpTransportConfiguration cfg = parse("remote-slot=3&remote-tsap=4660&local-tsap=18193");
        // Explicit overrides win even though rack/slot would derive something else.
        assertEquals(0x1234, cfg.getRemoteTsap());
        assertEquals(0x4711, cfg.getLocalTsap());
    }

    @Test
    void zeroTsapKeepsDerivedValue() {
        // The 0 sentinel means "not set" - the rack/slot-derived value must be used.
        S7CotpTransportConfiguration cfg = parse("remote-slot=3&remote-tsap=0&local-tsap=0");
        assertEquals(0x0103, cfg.getRemoteTsap());
        assertEquals(0x0311, cfg.getLocalTsap());
    }

    @Test
    void invalidDeviceGroupInConnectionStringFailsLoudly() {
        assertThrows(IllegalArgumentException.class,
            () -> parse("local-device-group=NOT_A_GROUP"));
    }
}
