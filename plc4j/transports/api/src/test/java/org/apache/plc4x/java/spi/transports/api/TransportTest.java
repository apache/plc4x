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

package org.apache.plc4x.java.spi.transports.api;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TransportTest {

    @Test
    void testDefaultGetTransportConfigType() {
        // Create a minimal implementation to test the default method
        Transport<TransportConfiguration> transport = new Transport<TransportConfiguration>() {
            @Override
            public String getTransportCode() {
                return "test";
            }

            @Override
            public String getTransportName() {
                return "Test Transport";
            }

            @Override
            public TransportInstance<TransportConfiguration> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
                return null;
            }
        };

        // The default implementation returns null
        assertNull(transport.getTransportConfigType());
    }

    @Test
    void testTransportCodeAndName() {
        Transport<TransportConfiguration> transport = new Transport<TransportConfiguration>() {
            @Override
            public String getTransportCode() {
                return "tcp";
            }

            @Override
            public String getTransportName() {
                return "TCP Transport";
            }

            @Override
            public TransportInstance<TransportConfiguration> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
                return null;
            }
        };

        assertEquals("tcp", transport.getTransportCode());
        assertEquals("TCP Transport", transport.getTransportName());
    }

    @Test
    void testCustomGetTransportConfigType() {
        Transport<TransportConfiguration> transport = new Transport<TransportConfiguration>() {
            @Override
            public String getTransportCode() {
                return "test";
            }

            @Override
            public String getTransportName() {
                return "Test Transport";
            }

            @Override
            public Class<TransportConfiguration> getTransportConfigType() {
                return TransportConfiguration.class;
            }

            @Override
            public TransportInstance<TransportConfiguration> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
                return null;
            }
        };

        assertEquals(TransportConfiguration.class, transport.getTransportConfigType());
    }
}
