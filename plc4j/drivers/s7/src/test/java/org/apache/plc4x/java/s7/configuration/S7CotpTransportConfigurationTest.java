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

class S7CotpTransportConfigurationTest {

    @Test
    void defaultPort() {
        assertEquals(102, new S7CotpTransportConfiguration().getDefaultPort());
    }

    @Test
    void defaultsDeriveTsapsFromOthersAndPgOrPc() {
        S7CotpTransportConfiguration cfg = new S7CotpTransportConfiguration();
        // Constructor seeds: local = PG_OR_PC rack=1 slot=1; remote = PG_OR_PC rack=0 slot=0.
        assertNotEquals(0, cfg.localTsap);
        assertNotEquals(0, cfg.remoteTsap);
    }

    @Test
    void rackSlotSettersRecomputeTsap() {
        S7CotpTransportConfiguration cfg = new S7CotpTransportConfiguration();
        int initial = cfg.localTsap;
        cfg.setLocalRack(2);
        cfg.setLocalSlot(3);
        // Rack/slot change should produce a different TSAP than the initial.
        assertNotEquals(initial, cfg.localTsap);
        assertEquals(2, cfg.getLocalRack());
        assertEquals(3, cfg.getLocalSlot());
    }

    @Test
    void deviceGroupSetterRecomputesTsap() {
        S7CotpTransportConfiguration cfg = new S7CotpTransportConfiguration();
        int initial = cfg.localTsap;
        cfg.setLocalDeviceGroup(DeviceGroup.OS);
        assertNotEquals(initial, cfg.localTsap);
        assertEquals(DeviceGroup.OS, cfg.getLocalDeviceGroup());
    }

    @Test
    void invalidDeviceGroupInConnectionStringFailsLoudly() {
        // Now that the field is enum-typed, ConfigurationFactory rejects unknown values at
        // parse time instead of silently keeping the default — what we want.
        assertThrows(IllegalArgumentException.class, () -> new ConfigurationFactory()
            .createConfiguration(S7CotpTransportConfiguration.class, "local-device-group=NOT_A_GROUP"));
    }

    @Test
    void remoteSettersWork() {
        S7CotpTransportConfiguration cfg = new S7CotpTransportConfiguration();
        int initial = cfg.remoteTsap;
        cfg.setRemoteRack(1);
        cfg.setRemoteSlot(2);
        cfg.setRemoteDeviceGroup(DeviceGroup.OTHERS);
        assertNotEquals(initial, cfg.remoteTsap);
        assertEquals(1, cfg.getRemoteRack());
        assertEquals(2, cfg.getRemoteSlot());
        assertEquals(DeviceGroup.OTHERS, cfg.getRemoteDeviceGroup());
    }

    @Test
    void explicitTsapOverridesDerived() {
        S7CotpTransportConfiguration cfg = new S7CotpTransportConfiguration();
        cfg.setLocalTsap(0x4711);
        cfg.setRemoteTsap(0x1234);
        assertEquals(0x4711, cfg.localTsap);
        assertEquals(0x1234, cfg.remoteTsap);
    }

    @Test
    void zeroTsapIsIgnoredKeepingDerivedValue() {
        S7CotpTransportConfiguration cfg = new S7CotpTransportConfiguration();
        int derivedLocal = cfg.localTsap;
        int derivedRemote = cfg.remoteTsap;
        cfg.setLocalTsap(0);
        cfg.setRemoteTsap(0);
        assertEquals(derivedLocal, cfg.localTsap);
        assertEquals(derivedRemote, cfg.remoteTsap);
    }
}
