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
package org.apache.plc4x.java.transport.tls;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tls.config.TlsTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetSocketAddress;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * TLS Transport factory that creates TLS-secured transport instances.
 * Provides encrypted connections with optional certificate validation bypass
 * for development environments with self-signed certificates.
 */
public class TlsTransport implements Transport<TlsTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TlsTransport.class);

    private static final Pattern TRANSPORT_TLS_PATTERN = Pattern.compile(
        "^((?<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?<hostname>[a-zA-Z0-9.\\-]+))(:(?<port>[0-9]{1,5}))?(?<driverConfig>.*)");

    @Override
    public String getTransportCode() {
        return "tls";
    }

    @Override
    public String getTransportName() {
        return "TLS";
    }

    @Override
    public Class<TlsTransportConfiguration> getTransportConfigType() {
        return TlsTransportConfiguration.class;
    }

    @Override
    public TransportInstance<TlsTransportConfiguration> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof TlsTransportConfiguration tlsTransportConfiguration)) {
            throw new IllegalArgumentException(String.format("Expected configuration of type %s but got %s",
                TlsTransportConfiguration.class.getSimpleName(), configuration.getClass().getSimpleName()));
        }

        final Matcher matcher = TRANSPORT_TLS_PATTERN.matcher(transportUrl);
        if (!matcher.matches()) {
            throw new PlcRuntimeException("Invalid url for TLS transport: " + transportUrl);
        }
        String ip = matcher.group("ip");
        String hostname = matcher.group("hostname");
        String portString = matcher.group("port");
        // Everything after host:port is the driver-config (e.g. an OPC UA "/milo" path).
        String driverConfig = matcher.group("driverConfig");

        // If the port wasn't specified, try to get a default port from the configuration.
        int port;
        if (portString != null) {
            port = Integer.parseInt(portString);
        } else if (tlsTransportConfiguration.getDefaultPort() != TlsTransportConfiguration.NO_DEFAULT_PORT) {
            port = tlsTransportConfiguration.getDefaultPort();
        } else {
            throw new PlcRuntimeException("No port defined");
        }

        InetSocketAddress remoteAddress = new InetSocketAddress((ip != null) ? ip : hostname, port);

        LOGGER.debug("Creating TLS transport instance for {}:{} (verify={})",
            (ip != null) ? ip : hostname, port, tlsTransportConfiguration.isVerifySsl());

        TlsTransportInstance instance = new TlsTransportInstance(remoteAddress, tlsTransportConfiguration, auditLog);
        instance.setDriverConfig(driverConfig);
        return instance;
    }

}
