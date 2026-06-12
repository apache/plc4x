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
package org.apache.plc4x.java.transport.test;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.test.config.TestTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Test transport implementation for testing purposes.
 * This transport uses in-memory byte streams for communication,
 * allowing testing without actual network connections.
 */
public class TestTransport implements Transport<TestTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TestTransport.class);

    @Override
    public String getTransportCode() {
        return "test";
    }

    @Override
    public String getTransportName() {
        return "Test Transport";
    }

    @Override
    public Class<TestTransportConfiguration> getTransportConfigType() {
        return TestTransportConfiguration.class;
    }

    @Override
    public TransportInstance<TestTransportConfiguration> createTransportInstance(
            String transportUrl,
            TransportConfiguration configuration,
            AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof TestTransportConfiguration testTransportConfiguration)) {
            throw new IllegalArgumentException(String.format("Expected configuration of type %s but got %s",
                TestTransportConfiguration.class.getSimpleName(), configuration.getClass().getSimpleName()));
        }

        LOGGER.debug("Creating test transport instance for URL: {}", transportUrl);
        return new TestTransportInstance(testTransportConfiguration, auditLog);
    }

}
