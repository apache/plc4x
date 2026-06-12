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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.eip.base.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.base.configuration.EipTcpTransportConfiguration;
import org.apache.plc4x.java.eip.base.discovery.EipPlcDiscoverer;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.readwrite.Constants;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;

public class EIPDriver extends DriverBase {

    public static final int PORT = Constants.EIPTCPDEFAULTPORT;

    @Override
    public String getProtocolCode() {
        return "eip";
    }

    @Override
    public String getProtocolName() {
        return "EthernetIP";
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
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        return new DefaultPlcDiscoveryRequest.Builder(new EipPlcDiscoverer());
    }

    @Override
    protected boolean canDiscover() {
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
    protected ConnectionBase<?> getConnection(Configuration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        return new EipTcpConnection((EIPConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public EipTag prepareTag(String tagAddress) {
        return EipTag.of(tagAddress);
    }

}
