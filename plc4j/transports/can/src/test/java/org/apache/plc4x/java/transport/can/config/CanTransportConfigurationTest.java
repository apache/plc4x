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
package org.apache.plc4x.java.transport.can.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.can.CanIdFilter;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class CanTransportConfigurationTest {

    @Test
    void implementsTransportConfiguration() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void defaultValues() {
        CanTransportConfiguration config = new CanTransportConfiguration();

        assertNull(config.filterIds);
        assertEquals(-1, config.filterRangeStart);
        assertEquals(-1, config.filterRangeEnd);
    }

    @Test
    void fieldsAreSettable() {
        CanTransportConfiguration config = new CanTransportConfiguration();

        config.filterIds = "0x100,0x200";
        config.filterRangeStart = 0x300;
        config.filterRangeEnd = 0x3FF;

        assertEquals("0x100,0x200", config.filterIds);
        assertEquals(0x300, config.filterRangeStart);
        assertEquals(0x3FF, config.filterRangeEnd);
    }

    @Test
    void buildFilterWithNoParams() {
        CanTransportConfiguration config = new CanTransportConfiguration();

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.isEmpty());
        assertTrue(filter.matches(0x100));
    }

    @Test
    void buildFilterWithIds() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "0x100,0x200,768"; // Mix of hex and decimal

        CanIdFilter filter = config.buildFilter();

        assertFalse(filter.isEmpty());
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x200));
        assertTrue(filter.matches(768));    // 0x300
        assertFalse(filter.matches(0x150));
    }

    @Test
    void buildFilterWithRange() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterRangeStart = 0x100;
        config.filterRangeEnd = 0x1FF;

        CanIdFilter filter = config.buildFilter();

        assertFalse(filter.isEmpty());
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x150));
        assertTrue(filter.matches(0x1FF));
        assertFalse(filter.matches(0x200));
    }

    @Test
    void buildFilterWithIdsAndRange() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "0x000";
        config.filterRangeStart = 0x100;
        config.filterRangeEnd = 0x1FF;

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.matches(0x000));  // explicit ID
        assertTrue(filter.matches(0x150));  // in range
        assertFalse(filter.matches(0x050)); // neither
    }

    @Test
    void buildFilterWithEmptyFilterIds() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "";

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.isEmpty());
    }

    @Test
    void buildFilterWithWhitespaceFilterIds() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "  ";

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.isEmpty());
    }

    @Test
    void buildFilterWithSpacesInIds() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = " 0x100 , 0x200 ";

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x200));
    }

    @Test
    void buildFilterWithInvalidIdThrows() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "notANumber";

        assertThrows(IllegalArgumentException.class, config::buildFilter);
    }

    @Test
    void buildFilterOnlyRangeStartSetDoesNotApplyRange() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterRangeStart = 0x100;
        // filterRangeEnd stays at default -1

        CanIdFilter filter = config.buildFilter();

        // Neither rangeStart nor rangeEnd are both >= 0, so no range applied
        assertTrue(filter.isEmpty());
    }

    @Test
    void buildFilterOnlyRangeEndSetDoesNotApplyRange() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterRangeEnd = 0x200;
        // filterRangeStart stays at default -1

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.isEmpty());
    }

    @Test
    void buildFilterWithDecimalIds() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "256,512";

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.matches(256));
        assertTrue(filter.matches(512));
        assertFalse(filter.matches(300));
    }

    @Test
    void buildFilterWithUpperCaseHexPrefix() {
        CanTransportConfiguration config = new CanTransportConfiguration();
        config.filterIds = "0X100";

        CanIdFilter filter = config.buildFilter();

        assertTrue(filter.matches(0x100));
    }
}
