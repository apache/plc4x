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
package org.apache.plc4x.java.knxnetip;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpUdpTransportConfiguration;
import org.apache.plc4x.java.knxnetip.discovery.KnxNetIpPlcDiscoverer;
import org.apache.plc4x.java.knxnetip.tag.KnxNetIpTag;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.CompletableFuture;

public class KnxNetIpDriver extends DriverBase implements PlcDiscoverer {

    public static final int KNXNET_IP_PORT = 3671;

    @Override
    public String getProtocolCode() {
        return "knxnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "KNXNet/IP";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return KnxNetIpConfiguration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("udp".equals(transport.getTransportCode())) {
            return KnxNetIpUdpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("udp");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return List.of("udp");
    }

    @Override
    public Set<Integer> defaultPorts(String transportCode) {
        if ("udp".equalsIgnoreCase(transportCode)) {
            return Set.of(KNXNET_IP_PORT);
        }
        return Collections.emptySet();
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
    protected boolean canSubscribe() {
        return true;
    }

    @Override
    protected boolean canBrowse() {
        return true;
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return new KnxNetIpPlcDiscoverer().discover(discoveryRequest);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest,
                                                                       PlcDiscoveryItemHandler handler) {
        return new KnxNetIpPlcDiscoverer().discoverWithHandler(discoveryRequest, handler);
    }

    @Override
    protected ConnectionBase<?> getConnection(Configuration configuration,
                                              TransportInstance<?> transportInstance,
                                              AuditLog auditLog) {
        return new KnxNetIpConnection((KnxNetIpConfiguration) configuration, transportInstance, auditLog);
    }

    @Override
    public KnxNetIpTag prepareTag(String tagAddress) {
        return KnxNetIpTag.of(tagAddress);
    }

}
