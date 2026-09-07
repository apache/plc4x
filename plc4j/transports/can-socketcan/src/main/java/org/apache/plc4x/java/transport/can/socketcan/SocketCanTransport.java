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
package org.apache.plc4x.java.transport.can.socketcan;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.can.socketcan.config.SocketCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Transport factory for the Linux SocketCAN transport.
 * <p>
 * Creates {@link SocketCanTransportInstance} instances for communicating over the
 * Linux kernel's SocketCAN subsystem. Each transport factory holds a single
 * {@link SharedCanManager} instance shared across all transport instances it creates,
 * enabling efficient resource sharing when multiple instances access the same CAN interface.
 * <p>
 * Transport code: {@code can-socketcan}
 */
public class SocketCanTransport implements Transport<SocketCanTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(SocketCanTransport.class);

    private final SharedCanManager sharedCanManager = new SharedCanManager();

    @Override
    public String getTransportCode() {
        return "can-socketcan";
    }

    @Override
    public String getTransportName() {
        return "SocketCAN";
    }

    @Override
    public Class<SocketCanTransportConfiguration> getTransportConfigType() {
        return SocketCanTransportConfiguration.class;
    }

    @Override
    public TransportInstance<SocketCanTransportConfiguration> createTransportInstance(
            String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof SocketCanTransportConfiguration socketCanConfig)) {
            throw new IllegalArgumentException(String.format(
                    "Expected configuration of type %s but got %s",
                    SocketCanTransportConfiguration.class.getSimpleName(),
                    configuration.getClass().getSimpleName()));
        }

        LOGGER.debug("Creating SocketCAN transport instance for interface {} (reuseInterface={})",
                socketCanConfig.interfaceName, socketCanConfig.reuseInterface);
        return new SocketCanTransportInstance(sharedCanManager, socketCanConfig, auditLog);
    }
}
