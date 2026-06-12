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
package org.apache.plc4x.java.transport.can.virtualcan.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.can.CanIdFilter;
import org.apache.plc4x.java.transport.can.config.CanTransportConfiguration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link VirtualCanTransportConfiguration} field defaults, type hierarchy,
 * and inherited CAN filter configuration.
 */
class VirtualCanTransportConfigurationTest {

    @Test
    void testImplementsTransportConfiguration() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void testExtendsCanTransportConfiguration() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        assertInstanceOf(CanTransportConfiguration.class, config);
    }

    @Test
    void testDefaultBusName() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        assertEquals("default", config.busName);
    }

    @Test
    void testBusNameIsSettable() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.busName = "my-test-bus";
        assertEquals("my-test-bus", config.busName);
    }

    @Test
    void testInheritsFilterDefaults() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();

        // Inherited filter fields should have their defaults
        assertNull(config.filterIds);
        assertEquals(-1, config.filterRangeStart);
        assertEquals(-1, config.filterRangeEnd);
    }

    @Test
    void testInheritedFilterIdsSettable() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.filterIds = "0x100,0x200";
        assertEquals("0x100,0x200", config.filterIds);
    }

    @Test
    void testInheritedFilterRangeSettable() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.filterRangeStart = 0x100;
        config.filterRangeEnd = 0x1FF;
        assertEquals(0x100, config.filterRangeStart);
        assertEquals(0x1FF, config.filterRangeEnd);
    }

    @Test
    void testBuildFilterWithNoFilters() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        CanIdFilter filter = config.buildFilter();

        // No filter configured means accept-all
        assertTrue(filter.isEmpty());
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x000));
    }

    @Test
    void testBuildFilterWithIds() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.filterIds = "0x100,0x200";
        CanIdFilter filter = config.buildFilter();

        assertFalse(filter.isEmpty());
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x200));
        assertFalse(filter.matches(0x300));
    }

    @Test
    void testBuildFilterWithRange() {
        VirtualCanTransportConfiguration config = new VirtualCanTransportConfiguration();
        config.filterRangeStart = 0x100;
        config.filterRangeEnd = 0x1FF;
        CanIdFilter filter = config.buildFilter();

        assertFalse(filter.isEmpty());
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x150));
        assertTrue(filter.matches(0x1FF));
        assertFalse(filter.matches(0x200));
        assertFalse(filter.matches(0x0FF));
    }
}
