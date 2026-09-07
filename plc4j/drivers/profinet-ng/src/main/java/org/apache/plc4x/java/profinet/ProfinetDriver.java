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
package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.profinet.channel.ProfinetChannel;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.config.ProfinetRawSocketTransportConfiguration;
import org.apache.plc4x.java.profinet.discovery.ProfinetDiscoverer;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.core.Pcaps;

import java.util.List;
import java.util.Optional;
import java.util.regex.Pattern;

public class ProfinetDriver extends DriverBase {

    public static final Pattern MAC_ADDRESS = Pattern.compile(
        "^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})?");

    public static final String DRIVER_CODE = "profinet";

    @Override
    public String getProtocolCode() {
        return DRIVER_CODE;
    }

    @Override
    public String getProtocolName() {
        return "Profinet";
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        try {
            // The discoverer opens its own pcap4j handles to broadcast on all
            // interfaces — it doesn't share the per-connection transport.
            ProfinetChannel channel = new ProfinetChannel(Pcaps.findAllDevs());
            ProfinetDiscoverer discoverer = new ProfinetDiscoverer(channel);
            return new DefaultPlcDiscoveryRequest.Builder(discoverer);
        } catch (PcapNativeException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return ProfinetConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("raw-socket".equals(transport.getTransportCode())) {
            return ProfinetRawSocketTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("raw-socket");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("raw-socket");
    }

    @Override
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected boolean canBrowse() {
        return true;
    }

    @Override
    protected boolean canDiscover() {
        return true;
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new ProfinetConnection((ProfinetConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public ProfinetTag prepareTag(String query) {
        return ProfinetTag.of(query);
    }

}
