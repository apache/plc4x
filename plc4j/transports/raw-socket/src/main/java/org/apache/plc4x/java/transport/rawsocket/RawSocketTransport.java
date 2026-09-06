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
package org.apache.plc4x.java.transport.rawsocket;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.rawsocket.config.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class RawSocketTransport implements Transport<RawSocketTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(RawSocketTransport.class);

    private final SharedRawSocketManager sharedRawSocketManager = new SharedRawSocketManager();

    @Override
    public String getTransportCode() {
        return "raw-socket";
    }

    @Override
    public String getTransportName() {
        return "Raw Socket (Ethernet)";
    }

    @Override
    public Class<RawSocketTransportConfiguration> getTransportConfigType() {
        return RawSocketTransportConfiguration.class;
    }

    @Override
    public TransportInstance<RawSocketTransportConfiguration> createTransportInstance(
            String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof RawSocketTransportConfiguration rawSocketTransportConfiguration)) {
            throw new IllegalArgumentException(String.format("Expected configuration of type %s but got %s",
                RawSocketTransportConfiguration.class.getSimpleName(), configuration.getClass().getSimpleName()));
        }

        // The address segment of the connection string names the network interface -
        // "raw-socket://en0" - the way every other transport is addressed by it. The
        // "interface-name" parameter stays as the alternative for callers that assemble a
        // connection string out of options alone; naming the interface in both places is
        // contradictory, and the address segment wins. Naming it in neither is still fine: the
        // transport then falls back to the first interface it finds.
        if (transportUrl != null && !transportUrl.trim().isEmpty()) {
            rawSocketTransportConfiguration.interfaceName = transportUrl.trim();
        }

        LOGGER.debug("Creating raw socket transport instance on interface {} for protocol 0x{} (reuseInterface={})",
        rawSocketTransportConfiguration.interfaceName,
        String.format("%04X", rawSocketTransportConfiguration.protocolId), rawSocketTransportConfiguration.reuseInterface);
        return new RawSocketTransportInstance(sharedRawSocketManager, rawSocketTransportConfiguration, auditLog);
    }

}
