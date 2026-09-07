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
package org.apache.plc4x.java.transport.tls;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tls.config.PskTlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetSocketAddress;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * TLS-PSK Transport factory that creates TLS-PSK secured transport instances.
 * Uses Pre-Shared Key authentication via Bouncy Castle for devices that support TLS-PSK
 * (e.g., Beckhoff TwinCAT 3 Secure ADS).
 */
public class PskTlsTransport implements Transport<PskTlsTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(PskTlsTransport.class);

    private static final Pattern TRANSPORT_TLS_PATTERN = Pattern.compile(
        "^((?<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?<hostname>[a-zA-Z0-9.\\-]+))(:(?<port>[0-9]{1,5}))?.*");

    @Override
    public String getTransportCode() {
        return "tls-psk";
    }

    @Override
    public String getTransportName() {
        return "TLS-PSK";
    }

    @Override
    public Class<PskTlsTransportConfiguration> getTransportConfigType() {
        return PskTlsTransportConfiguration.class;
    }

    @Override
    public TransportInstance<PskTlsTransportConfiguration> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof PskTlsTransportConfiguration pskTlsTransportConfiguration)) {
            throw new IllegalArgumentException(String.format("Expected configuration of type %s but got %s",
                PskTlsTransportConfiguration.class.getSimpleName(), configuration.getClass().getSimpleName()));
        }

        final Matcher matcher = TRANSPORT_TLS_PATTERN.matcher(transportUrl);
        if (!matcher.matches()) {
            throw new PlcRuntimeException("Invalid url for TLS-PSK transport: " + transportUrl);
        }
        String ip = matcher.group("ip");
        String hostname = matcher.group("hostname");
        String portString = matcher.group("port");

        // If the port wasn't specified, try to get a default port from the configuration.
        int port;
        if (portString != null) {
            port = Integer.parseInt(portString);
        } else if (pskTlsTransportConfiguration.getDefaultPort() != PskTlsTransportConfiguration.NO_DEFAULT_PORT) {
            port = pskTlsTransportConfiguration.getDefaultPort();
        } else {
            throw new PlcRuntimeException("No port defined");
        }

        InetSocketAddress remoteAddress = new InetSocketAddress((ip != null) ? ip : hostname, port);

        LOGGER.debug("Creating TLS-PSK transport instance for {}:{}",
            (ip != null) ? ip : hostname, port);

        pskTlsTransportConfiguration.validatePskConfiguration();
        return new PskTlsTransportInstance(remoteAddress, pskTlsTransportConfiguration, auditLog);
    }

}
