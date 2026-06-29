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

package org.apache.plc4x.java.transport.cotp;

import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.cotp.config.CotpTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * COTP (Connection Oriented Transport Protocol) Transport.
 * <p>
 * This transport implements ISO 8073 (COTP) with TPKT framing (RFC 1006)
 * on top of TCP. It handles:
 * <ul>
 *   <li>TPKT framing (4-byte header: version, reserved, length)</li>
 *   <li>COTP connection establishment (CR/CC packets)</li>
 *   <li>COTP data transmission (DT packets)</li>
 *   <li>COTP disconnection (DR/DC packets)</li>
 * </ul>
 * <p>
 * The transport extracts the payload from COTP Data packets and presents
 * it to drivers, hiding the TPKT/COTP framing details.
 */
public class CotpTransport implements Transport<CotpTransportConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(CotpTransport.class);

    private static final Pattern TRANSPORT_COTP_PATTERN = Pattern.compile(
        "^((?<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?<hostname>[a-zA-Z0-9.\\-]+))(:(?<port>[0-9]{1,5}))?.*");

    @Override
    public String getTransportCode() {
        return "cotp";
    }

    @Override
    public String getTransportName() {
        return "COTP";
    }

    @Override
    public Class<CotpTransportConfiguration> getTransportConfigType() {
        return CotpTransportConfiguration.class;
    }

    @Override
    public TransportInstance<CotpTransportConfiguration> createTransportInstance(
            String transportUrl,
            TransportConfiguration configuration, AuditLog auditLog) throws TransportException {

        if (!(configuration instanceof CotpTransportConfiguration cotpConfig)) {
            throw new IllegalArgumentException(String.format(
                "Expected configuration of type %s but got %s",
                CotpTransportConfiguration.class.getSimpleName(),
                configuration.getClass().getSimpleName()));
        }

        final Matcher matcher = TRANSPORT_COTP_PATTERN.matcher(transportUrl);
        if (!matcher.matches()) {
            throw new TransportException("Invalid url for COTP transport: " + transportUrl);
        }

        String ip = matcher.group("ip");
        String hostname = matcher.group("hostname");
        String portString = matcher.group("port");

        // If the port wasn't specified, try to get a default port from the configuration.
        int port;
        if (portString != null) {
            port = Integer.parseInt(portString);
        } else if (cotpConfig.getDefaultPort() != CotpTransportConfiguration.NO_DEFAULT_PORT) {
            port = cotpConfig.getDefaultPort();
        } else {
            throw new TransportException("No port defined");
        }

        String host = (ip != null) ? ip : hostname;

        LOGGER.debug("Creating COTP transport instance for {}:{}", host, port);
        LOGGER.debug("COTP Configuration: localTsap=0x{}, remoteTsap=0x{}, tpduSize={}",
            Integer.toHexString(cotpConfig.getLocalTsap()),
            Integer.toHexString(cotpConfig.getRemoteTsap()),
            cotpConfig.cotpTpduSize);

        return new CotpTransportInstance(host, port, cotpConfig, auditLog);
    }
}
