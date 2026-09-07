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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.slmp.config.SlmpConfiguration;
import org.apache.plc4x.java.slmp.config.SlmpTcpTransportConfiguration;
import org.apache.plc4x.java.slmp.tag.SlmpTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.List;
import java.util.Optional;
import java.util.Set;

public class SlmpDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "slmp";
    }

    @Override
    public String getProtocolName() {
        return "SLMP (MELSEC) 3E";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return SlmpConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("tcp".equals(transport.getTransportCode())) {
            return SlmpTcpTransportConfiguration.class;
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
            return Set.of(5007);
        }
        return Set.of();
    }

    @Override
    public PlcTag prepareTag(String tagAddress) {
        return SlmpTag.of(tagAddress);
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new SlmpConnection((SlmpConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    protected boolean canRead() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }
}
