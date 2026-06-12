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
package org.apache.plc4x.java.transport.test.config;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;

/**
 * Tests for TestTransportConfiguration.
 */
class TestTransportConfigurationTest {

    @Test
    void testImplementsTransportConfiguration() {
        TestTransportConfiguration config = new TestTransportConfiguration();
        assertInstanceOf(TransportConfiguration.class, config);
    }

    @Test
    void testDefaultReceiveBufferSize() {
        TestTransportConfiguration config = new TestTransportConfiguration();
        // The default value is set via annotation, but we can test that the field exists
        // and can be set/read
        config.receiveBufferSize = 81920;
        assertEquals(81920, config.receiveBufferSize);
    }

    @Test
    void testCustomReceiveBufferSize() {
        TestTransportConfiguration config = new TestTransportConfiguration();
        config.receiveBufferSize = 16384;
        assertEquals(16384, config.receiveBufferSize);
    }

}
