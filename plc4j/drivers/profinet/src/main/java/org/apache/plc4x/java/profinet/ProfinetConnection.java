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

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.config.ConfigurationProfinetDevice;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.device.ProfinetDeviceMessageHandler;
import org.apache.plc4x.java.profinet.device.ProfinetMessageWrapper;
import org.apache.plc4x.java.profinet.device.ProfinetNetworkInterface;
import org.apache.plc4x.java.profinet.device.ProfinetSubscriptionHandle;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.profinet.tag.ProfinetTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionResponse;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.pcap4j.core.PcapNativeException;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

/**
 * PROFINET connection — collapses the old SPI's {@code ProfinetProtocolLogic}
 * onto the new SPI {@link ConnectionBase}. The protocol logic itself is small
 * because the heavy lifting lives in the {@code device/} subsystem
 * ({@link ProfinetDevice}, {@link ProfinetChannel}, etc.) which talks to the
 * wire via its own pcap4j handles rather than going through the SPI transport.
 *
 * <p>The SPI transport is used essentially only to learn the local IPv4 address
 * (so we can pick the right pcap interface). All real PROFINET traffic flows
 * through {@link ProfinetChannel}.</p>
 */
public class ProfinetConnection extends ConnectionBase<ProfinetConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(ProfinetConnection.class);

    private final ProfinetDriverContext profinetDriverContext = new ProfinetDriverContext();
    private final Map<String, ProfinetDevice> devices = new HashMap<>();
    private volatile boolean connected = false;

    public ProfinetConnection(ProfinetConfiguration configuration,
                              TransportInstance<?> transportInstance,
                              AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
        initializeDevices();
    }

    private void initializeDevices() {
        profinetDriverContext.setConfiguration(configuration);
        Map<String, ConfigurationProfinetDevice> configuredDevices =
            configuration.getDevices().getConfiguredDevices();
        for (Map.Entry<String, ConfigurationProfinetDevice> entry : configuredDevices.entrySet()) {
            ProfinetDevice device = new ProfinetDevice(
                new ProfinetMessageWrapper(),
                entry.getValue().getDevicename(),
                entry.getValue().getDeviceaccess(),
                entry.getValue().getSubmodules(),
                entry.getValue().getGsdHandler());
            device.setIpAddress(entry.getValue().getIpaddress());
            devices.put(entry.getKey(), device);
        }
        profinetDriverContext.setHandler(new ProfinetDeviceMessageHandler(devices));
        for (ProfinetDevice device : devices.values()) {
            device.getDeviceContext().setConfiguration(configuration);
        }
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new ProfinetTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    public boolean isConnected() {
        return connected;
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        if (!(transportInstance instanceof RawSocketTransportInstance rawSocket)) {
            throw new PlcConnectionException(
                "PROFINET requires a raw-socket transport (got " + transportInstance.getClass().getName() + ")");
        }
        byte[] localIpBytes = rawSocket.getLocalIpAddress();
        if (localIpBytes == null || localIpBytes.length != 4) {
            throw new PlcConnectionException(
                "PROFINET requires an interface with a local IPv4 address");
        }
        try {
            InetAddress localIpAddress = InetAddress.getByAddress(localIpBytes);
            PcapNetworkInterface devByAddress = Pcaps.getDevByAddress(localIpAddress);
            profinetDriverContext.setChannel(new ProfinetChannel(
                Collections.singletonList(devByAddress), devices));
            profinetDriverContext.getChannel().setConfiguredDevices(devices);

            // Wire each device to a per-connection callback so it can signal
            // when its handshake is done. Replaces the old SPI's
            // {@code ConversationContext.fireConnected()}.
            CountDownLatch deviceReady = new CountDownLatch(devices.size());
            for (ProfinetDevice device : devices.values()) {
                device.getDeviceContext().setNetworkInterface(new ProfinetNetworkInterface(devByAddress));
                device.getDeviceContext().setChannel(profinetDriverContext.getChannel());
                device.setContext(deviceReady::countDown, profinetDriverContext.getChannel());
            }

            // Open the UDP port the devices send replies to.
            profinetDriverContext.setSocket(new DatagramSocket(ProfinetDriverContext.DEFAULT_UDP_PORT));
            profinetDriverContext.getHandler().setConfiguredDevices(devices);

            // Step 1: discover devices (LLDP + PN-DCP) so we learn the remote MACs.
            runDeviceDiscovery();
            waitForDeviceDiscovery();

            // Step 2: kick off the per-device handshake. The devices each call
            // their {@code OnConnectedCallback} once they reach ApplicationReady.
            for (ProfinetDevice device : devices.values()) {
                device.onConnect();
            }
            if (!deviceReady.await(60, TimeUnit.SECONDS)) {
                throw new PlcConnectionException("Timed out waiting for devices to reach ApplicationReady");
            }
            connected = true;
        } catch (PlcConnectionException e) {
            throw e;
        } catch (PcapNativeException | UnknownHostException | SocketException | InterruptedException
                 | ExecutionException | TimeoutException | PlcException e) {
            throw new PlcConnectionException("PROFINET connect failed", e);
        }
    }

    @Override
    public void close() throws Exception {
        connected = false;
        if (profinetDriverContext.getSocket() != null) {
            profinetDriverContext.getSocket().close();
        }
        super.close();
    }

    private void runDeviceDiscovery() throws PlcException, InterruptedException {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer(profinetDriverContext.getChannel());
        profinetDriverContext.getChannel().setDiscoverer(discoverer);
        DefaultPlcDiscoveryRequest request = new DefaultPlcDiscoveryRequest(discoverer, new LinkedHashMap<>());
        discoverer.ongoingDiscoverWithHandler(request, profinetDriverContext.getHandler(), 5000L, 30000L);
    }

    /**
     * Block until every configured device has been seen on the bus (DCP plus
     * optionally LLDP). Mirrors the original implementation including its 5-cycle
     * give-up threshold and per-cycle 3-second sleep.
     */
    private void waitForDeviceDiscovery() throws PlcConnectionException, InterruptedException {
        boolean discovered = false;
        int count = 0;
        while (!discovered) {
            List<ProfinetDevice> missing = new ArrayList<>();
            discovered = true;
            for (ProfinetDevice device : devices.values()) {
                if (!device.hasDcpPdu()) {
                    discovered = false;
                    missing.add(device);
                }
            }
            if (count > 5) {
                for (ProfinetDevice m : missing) {
                    if (m.hasDcpPdu()) {
                        LOGGER.info("For device {} we only managed to get a DCP discovery response — is the device on an LLDP-enabled switch?",
                            m.getDeviceId());
                    }
                }
                throw new PlcConnectionException("One device failed to respond to discovery packet");
            }
            if (!discovered) {
                Thread.sleep(3000L);
                count += 1;
            }
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Browse
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
        List<PlcBrowseItem> values = new ArrayList<>();
        for (ProfinetDevice device : devices.values()) {
            device.browseTags(values);
        }
        Map<String, PlcResponseCode> codes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> responseValues = new HashMap<>();
        for (String queryName : browseRequest.getQueryNames()) {
            responseValues.put(queryName, values);
            codes.put(queryName, PlcResponseCode.OK);
        }
        return CompletableFuture.completedFuture(new DefaultPlcBrowseResponse(browseRequest, codes, responseValues));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Subscribe
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
            for (String fieldName : subscriptionRequest.getTagNames()) {
                PlcSubscriptionTag subscriptionTag = subscriptionRequest.getTag(fieldName);
                DefaultPlcSubscriptionTag fieldDefault = (DefaultPlcSubscriptionTag) subscriptionTag;
                String deviceString = fieldDefault.getAddressString().split("\\.")[0].toUpperCase();
                ProfinetDevice device = devices.get(deviceString);
                if (device == null) {
                    values.put(fieldName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                    continue;
                }
                ProfinetSubscriptionHandle handle = new ProfinetSubscriptionHandle(this, fieldName, subscriptionTag);
                device.getDeviceContext().addSubscriptionHandle(fieldDefault.getAddressString(), handle);
                if (!(fieldDefault.getTag() instanceof ProfinetTag)) {
                    values.put(fieldName, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                } else {
                    values.put(fieldName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
                }
            }
            return new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
        });
    }

}
