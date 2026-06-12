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
package org.apache.plc4x.java.transport.tcp;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetSocketAddress;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class TcpTransport implements Transport<TcpTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(TcpTransport.class);

    private static final Pattern TRANSPORT_TCP_PATTERN = Pattern.compile(
        "^((?<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?<hostname>[a-zA-Z0-9.\\-]+))(:(?<port>[0-9]{1,5}))?.*");


    @Override
    public String getTransportCode() {
        return "tcp";
    }

    @Override
    public String getTransportName() {
        return "TCP";
    }

    @Override
    public Class<TcpTransportConfiguration> getTransportConfigType() {
        return TcpTransportConfiguration.class;
    }

    @Override
    public TransportInstance<TcpTransportConfiguration> createTransportInstance(String transportUrl, TransportConfiguration configuration, AuditLog auditLog) throws TransportException {
        if (!(configuration instanceof TcpTransportConfiguration tcpTransportConfiguration)) {
            throw new IllegalArgumentException(String.format("Expected configuration of type %s but got %s",
                TcpTransportConfiguration.class.getSimpleName(), configuration.getClass().getSimpleName()));
        }

        final Matcher matcher = TRANSPORT_TCP_PATTERN.matcher(transportUrl);
        if (!matcher.matches()) {
            throw new PlcRuntimeException("Invalid url for TCP transport: " + transportUrl);
        }
        String ip = matcher.group("ip");
        String hostname = matcher.group("hostname");
        String portString = matcher.group("port");

        // If the port wasn't specified, try to get a default port from the configuration.
        int port;
        if (portString != null) {
            port = Integer.parseInt(portString);
        } else if (tcpTransportConfiguration.getDefaultPort() != TcpTransportConfiguration.NO_DEFAULT_PORT) {
            port = tcpTransportConfiguration.getDefaultPort();
        } else {
            throw new PlcRuntimeException("No port defined");
        }

        InetSocketAddress remoteAddress = new InetSocketAddress((ip != null) ? ip : hostname, port);

        LOGGER.debug("Creating TCP transport instance for {}:{}", (ip != null) ? ip : hostname, port);
        return new TcpTransportInstance(remoteAddress, tcpTransportConfiguration, auditLog);
    }

}
