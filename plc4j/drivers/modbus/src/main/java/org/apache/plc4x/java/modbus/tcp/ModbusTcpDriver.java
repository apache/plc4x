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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.readwrite.Constants;
import org.apache.plc4x.java.modbus.tcp.config.*;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.modbus.tcp.discovery.ModbusPlcDiscoverer;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.bouncycastle.tls.UDPTransport;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;

public class ModbusTcpDriver extends DriverBase {

    @Override
    public String getProtocolCode() {
        return "modbus-tcp";
    }

    @Override
    public String getProtocolName() {
        return "Modbus TCP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return ModbusTcpConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("tcp".equals(transport.getTransportCode())) {
            return ModbusTcpTcpTransportConfiguration.class;
        } else if ("tls".equals(transport.getTransportCode())) {
            return ModbusTcpTlsTransportConfiguration.class;
        } else if ("tls-psk".equals(transport.getTransportCode())) {
            return ModbusTcpPskTlsTransportConfiguration.class;
        } else if (transport instanceof UDPTransport) {
            return ModbusTcpUdpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("tcp");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("tcp", "tls", "tls-psk", "udp", "test");
    }

    @Override
    public Set<Integer> defaultPorts(String transportCode) {
        if ("tcp".equalsIgnoreCase(transportCode)) {
            return Set.of(Constants.MODBUSTCPDEFAULTPORT);
        } else if ("tls".equalsIgnoreCase(transportCode) || "tls-psk".equalsIgnoreCase(transportCode)) {
            return Set.of(Constants.MODBUSTCPTLSDEFAULTPORT);
        } else if ("udp".equalsIgnoreCase(transportCode)) {
            return Set.of(Constants.MODBUSUDPDEFAULTPORT);
        }
        return Collections.emptySet();
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        return new DefaultPlcDiscoveryRequest.Builder(new ModbusPlcDiscoverer());
    }

    @Override
    protected boolean canDiscover() {
        return true;
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
    protected ConnectionBase<ModbusTcpConfiguration> getConnection(Configuration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        return new ModbusTcpConnection((ModbusTcpConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public ModbusTag prepareTag(String tagAddress) {
        return ModbusTag.of(tagAddress);
    }

}
