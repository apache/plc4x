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
package org.apache.plc4x.java.transport.can.virtualcan;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.can.virtualcan.config.VirtualCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Factory for creating {@link VirtualCanTransportInstance} objects.
 * <p>
 * This transport is intended for testing: it provides a pure-Java, in-memory
 * virtual CAN bus that requires no hardware, OS support, or native libraries.
 * Multiple instances created with the same bus name can exchange frames in-memory.
 */
public class VirtualCanTransport implements Transport<VirtualCanTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(VirtualCanTransport.class);

    /**
     * {@inheritDoc}
     *
     * @return the transport code {@code "can-virtualcan"}
     */
    @Override
    public String getTransportCode() {
        return "can-virtualcan";
    }

    /**
     * {@inheritDoc}
     *
     * @return the human-readable name {@code "Virtual CAN"}
     */
    @Override
    public String getTransportName() {
        return "Virtual CAN";
    }

    /**
     * {@inheritDoc}
     *
     * @return {@link VirtualCanTransportConfiguration}{@code .class}
     */
    @Override
    public Class<VirtualCanTransportConfiguration> getTransportConfigType() {
        return VirtualCanTransportConfiguration.class;
    }

    /**
     * {@inheritDoc}
     * <p>
     * Creates a new {@link VirtualCanTransportInstance} connected to the virtual
     * bus specified in the configuration.
     *
     * @throws IllegalArgumentException if the configuration is not a
     *                                  {@link VirtualCanTransportConfiguration}
     */
    @Override
    public TransportInstance<VirtualCanTransportConfiguration> createTransportInstance(
            String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof VirtualCanTransportConfiguration virtualCanConfig)) {
            throw new IllegalArgumentException(String.format(
                    "Expected configuration of type %s but got %s",
                    VirtualCanTransportConfiguration.class.getSimpleName(),
                    configuration.getClass().getSimpleName()));
        }

        LOGGER.debug("Creating virtual CAN transport instance on bus '{}'", virtualCanConfig.busName);
        return new VirtualCanTransportInstance(virtualCanConfig, auditLog);
    }
}
