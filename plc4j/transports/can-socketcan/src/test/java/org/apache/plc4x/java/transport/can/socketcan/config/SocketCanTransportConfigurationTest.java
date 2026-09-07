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
package org.apache.plc4x.java.transport.can.socketcan.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.can.config.CanTransportConfiguration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SocketCanTransportConfigurationTest {

    @Test
    void implementsTransportConfiguration() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void extendsCanTransportConfiguration() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        assertInstanceOf(CanTransportConfiguration.class, config);
    }

    @Test
    void defaultValues() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();

        assertNull(config.interfaceName);
        assertFalse(config.reuseInterface);
        assertEquals(1000, config.readTimeout);
    }

    @Test
    void fieldsAreSettable() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();

        config.interfaceName = "can0";
        config.reuseInterface = true;
        config.readTimeout = 5000;

        assertEquals("can0", config.interfaceName);
        assertTrue(config.reuseInterface);
        assertEquals(5000, config.readTimeout);
    }

    @Test
    void inheritsFilterParamsFromBase() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();

        config.filterIds = "0x100,0x200";
        config.filterRangeStart = 0x300;
        config.filterRangeEnd = 0x3FF;

        var filter = config.buildFilter();
        assertFalse(filter.isEmpty());
        assertTrue(filter.matches(0x100));
        assertTrue(filter.matches(0x350));
        assertFalse(filter.matches(0x050));
    }

    @Test
    void buildFilterWithNoParamsReturnsAcceptAll() {
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();

        var filter = config.buildFilter();
        assertTrue(filter.isEmpty());
        assertTrue(filter.matches(0x100));
    }
}
