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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.s7.configuration.S7CotpTransportConfiguration;
import org.apache.plc4x.java.s7.discovery.S7PlcDiscoverer;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.s7.tag.S7PlcTagHandler;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.DriverBase;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.CompletableFuture;
import java.util.regex.Matcher;

/**
 * Apache PLC4X driver for the Siemens S7 protocol over the COTP transport. Replaces both the
 * old Netty-based {@code s7} driver and the {@code s7-light} driver.
 */
public class S7Driver extends DriverBase implements PlcDiscoverer {

    private static final Logger LOGGER = LoggerFactory.getLogger(S7Driver.class);

    public static final int ISO_ON_TCP_PORT = 102;

    @Override
    public String getProtocolCode() {
        return "s7";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Basic)";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationClass() {
        return S7Configuration.class;
    }

    @Override
    protected Class<? extends TransportConfiguration> getTransportConfigurationClass(Transport<?> transport) {
        if ("cotp".equals(transport.getTransportCode())) {
            return S7CotpTransportConfiguration.class;
        }
        return super.getTransportConfigurationClass(transport);
    }

    @Override
    public Optional<String> getDefaultTransportCode() {
        return Optional.of("cotp");
    }

    @Override
    public List<String> getSupportedTransportCodes() {
        return Collections.singletonList("cotp");
    }

    @Override
    public Set<Integer> defaultPorts(String transportCode) {
        if ("cotp".equalsIgnoreCase(transportCode)) {
            return Set.of(ISO_ON_TCP_PORT);
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
    protected ConnectionBase<S7Configuration> getConnection(Configuration configuration, TransportInstance<?> transportInstance, AuditLog auditLog) {
        return new S7CotpConnection((S7Configuration) configuration, transportInstance, auditLog);
    }

    /**
     * Override the standard URL-to-connection pipeline to recognize the dual-path
     * <em>S7H</em> syntax: {@code s7://primary-host/secondary-host?...}. When detected,
     * we delegate to {@link DriverBase#getConnection(String, PlcAuthentication)} twice —
     * once per host — to get two fully-configured {@link S7CotpConnection} instances,
     * then wrap them in a {@link S7HCotpConnection} that handles failover. Single-host
     * URLs fall through to the standard pipeline unchanged.
     */
    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication plcAuthentication) throws PlcConnectionException {
        Matcher matcher = URI_PATTERN.matcher(connectionString);
        if (matcher.matches()) {
            String transportConfig = matcher.group("transportConfig");
            if (transportConfig != null && transportConfig.contains("/")) {
                String[] parts = transportConfig.split("/", 2);
                String primaryHost = parts[0];
                String secondaryHost = parts[1];
                if (!primaryHost.isBlank() && !secondaryHost.isBlank()) {
                    final String primaryUrl = rebuildUrl(matcher, primaryHost);
                    final String secondaryUrl = rebuildUrl(matcher, secondaryHost);
                    LOGGER.info("S7H dual-path connection — primary='{}', secondary='{}'",
                        primaryHost, secondaryHost);
                    // Build the primary once eagerly so we can grab its config + transport for
                    // the wrapper's parent ConnectionBase. The wrapper's factories will rebuild
                    // either inner from scratch on heartbeat-detected failures.
                    PlcConnection primaryConn = super.getConnection(primaryUrl, plcAuthentication);
                    if (!(primaryConn instanceof S7CotpConnection p)) {
                        throw new PlcConnectionException(
                            "Internal error: S7H primary inner is not an S7CotpConnection");
                    }
                    java.util.function.Supplier<S7CotpConnection> primaryFactory = () -> {
                        try {
                            PlcConnection c = super.getConnection(primaryUrl, plcAuthentication);
                            if (!(c instanceof S7CotpConnection s7)) {
                                throw new RuntimeException("Primary factory produced non-S7CotpConnection");
                            }
                            return s7;
                        } catch (PlcConnectionException ex) {
                            throw new RuntimeException(ex);
                        }
                    };
                    java.util.function.Supplier<S7CotpConnection> secondaryFactory = () -> {
                        try {
                            PlcConnection c = super.getConnection(secondaryUrl, plcAuthentication);
                            if (!(c instanceof S7CotpConnection s7)) {
                                throw new RuntimeException("Secondary factory produced non-S7CotpConnection");
                            }
                            return s7;
                        } catch (PlcConnectionException ex) {
                            throw new RuntimeException(ex);
                        }
                    };
                    // The eagerly-built primary becomes the first inner. The wrapper holds the
                    // factory and will use it to build secondaries / replacements.
                    java.util.concurrent.atomic.AtomicReference<S7CotpConnection> firstPrimary =
                        new java.util.concurrent.atomic.AtomicReference<>(p);
                    java.util.function.Supplier<S7CotpConnection> primaryFirstThenFactory = () -> {
                        S7CotpConnection initial = firstPrimary.getAndSet(null);
                        return initial != null ? initial : primaryFactory.get();
                    };
                    return new S7HCotpConnection(
                        p.getConfiguration(), p.getTransportInstance(), getAuditLog(),
                        primaryFirstThenFactory, secondaryFactory);
                }
            }
        }
        return super.getConnection(connectionString, plcAuthentication);
    }

    /** Reconstruct a single-host URL from a matched {@link #URI_PATTERN}. */
    private static String rebuildUrl(Matcher matcher, String host) {
        StringBuilder sb = new StringBuilder();
        sb.append(matcher.group("protocolCode"));
        if (matcher.group("transportCode") != null) {
            sb.append(':').append(matcher.group("transportCode"));
        }
        sb.append("://").append(host);
        if (matcher.group("paramString") != null) {
            sb.append('?').append(matcher.group("paramString"));
        }
        return sb.toString();
    }

    @Override
    protected boolean canDiscover() {
        // PROFINET DCP discovery requires the optional pcap4j runtime (and libpcap on the
        // host). DriverBase only exposes the discoveryRequestBuilder when canDiscover() is
        // true AND we implement PlcDiscoverer — gate explicitly here so callers without
        // pcap installed get the standard "not implemented" response rather than a
        // ClassNotFoundError at request time.
        return true;
    }

    /**
     * Each call constructs a fresh {@link S7PlcDiscoverer} (and underlying pcap handles),
     * runs one discovery cycle, and discards them. This avoids holding raw sockets open
     * across calls and matches the ADS driver's pattern.
     */
    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest,
                                                                       PlcDiscoveryItemHandler handler) {
        try {
            S7PlcDiscoverer discoverer = new S7PlcDiscoverer();
            return discoverer.discoverWithHandler(discoveryRequest, handler);
        } catch (Throwable t) {
            // Most commonly: pcap4j missing at runtime, or libpcap not installed on host.
            LOGGER.info("S7 discovery unavailable: {}", t.getMessage());
            return CompletableFuture.failedFuture(t);
        }
    }

    public S7Tag prepareTag(String tagAddress) {
        return (S7Tag) new S7PlcTagHandler().parseTag(tagAddress);
    }
}
