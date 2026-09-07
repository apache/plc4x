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
package org.apache.plc4x.java.umas;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.umas.configuration.UmasConfiguration;
import org.apache.plc4x.java.umas.configuration.UmasTcpTransportConfiguration;
import org.apache.plc4x.java.umas.readwrite.Constants;
import org.apache.plc4x.java.umas.tag.SymbolicUmasTag;
import org.apache.plc4x.java.umas.tag.UmasTag;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;

/**
 * PLC4J driver for the UMAS protocol (Schneider Electric Modicon PLCs).
 * UMAS is tunneled inside Modbus/TCP using function code 0x5A.
 * Connection URL format: {@code umas:tcp://host:port?unit-identifier=0}
 */
public class UmasDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "umas";
    }

    @Override
    public String getProtocolName() {
        return "UMAS (Schneider Electric)";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return UmasConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("tcp".equals(transport.getTransportCode())) {
            return UmasTcpTransportConfiguration.class;
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
            return Set.of(Constants.UMASTCPDEFAULTPORT);
        }
        return Collections.emptySet();
    }

    @Override
    protected boolean canPing() {
        return true;
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
    protected boolean canBrowse() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new UmasConnection((UmasConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public UmasTag prepareTag(String tagAddress) {
        return SymbolicUmasTag.of(tagAddress);
    }

}
