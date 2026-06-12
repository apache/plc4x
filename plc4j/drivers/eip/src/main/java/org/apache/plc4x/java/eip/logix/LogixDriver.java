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
package org.apache.plc4x.java.eip.logix;

import org.apache.plc4x.java.eip.base.EipTcpConnection;
import org.apache.plc4x.java.eip.base.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.base.configuration.EipTcpTransportConfiguration;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.readwrite.Constants;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;

/**
 * Thin alias over {@link org.apache.plc4x.java.eip.base.EIPDriver} that
 * registers the {@code logix} protocol code and forces little-endian wire
 * encoding (the Logix family of controllers always speaks LE on CIP, whereas
 * generic EIP devices may go either way). All actual protocol logic lives in
 * {@link EipTcpConnection} and {@link EIPConfiguration}; this class just
 * provides the {@code logix://...} URL prefix that existing Logix tooling
 * documentation uses.
 */
public class LogixDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "logix";
    }

    @Override
    public String getProtocolName() {
        return "Logix CIP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return EIPConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("tcp".equals(transport.getTransportCode())) {
            return EipTcpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("tcp");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("tcp", "test");
    }

    @Override
    public Set<Integer> defaultPorts(String transportCode) {
        if ("tcp".equalsIgnoreCase(transportCode)) {
            return Set.of(Constants.EIPTCPDEFAULTPORT);
        }
        return Collections.emptySet();
    }

    @Override
    protected boolean canRead() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        // Logix CIP is always little-endian on the wire; pass false explicitly
        // so the user doesn't have to remember big-endian=false on the URL.
        return new EipTcpConnection((EIPConfiguration) configuration, transportInstance, auditLog, false);
    }

    @Override
    public EipTag prepareTag(String tagAddress) {
        return EipTag.of(tagAddress);
    }

}
