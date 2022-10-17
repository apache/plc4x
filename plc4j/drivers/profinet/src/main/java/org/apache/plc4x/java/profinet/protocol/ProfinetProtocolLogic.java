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
package org.apache.plc4x.java.profinet.protocol;

import io.netty.channel.Channel;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.profinet.ProfinetDriver;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.device.ProfinetDeviceMessageHandler;
import org.apache.plc4x.java.profinet.device.ProfinetUdpMessageHandler;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.pcap4j.core.*;
import org.pcap4j.util.Inet4NetworkAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.*;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> implements HasConfiguration<ProfinetConfiguration>, PlcSubscriber {

    public static final Duration REQUEST_TIMEOUT = Duration.ofMillis(10000);
    private static final int DEFAULT_UDP_PORT = 34964;
    private final Logger logger = LoggerFactory.getLogger(ProfinetProtocolLogic.class);

    private ProfinetDriverContext profinetDriverContext;
    private boolean connected = false;


    private RawSocketChannel rawSocketChannel;
    private Channel channel;

    private String boundIpAddress;

    private ProfinetDeviceMessageHandler handler = new ProfinetDeviceMessageHandler();

    private ProfinetConfiguration configuration;
    private Map<MacAddress, PcapHandle> openHandles;


    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
        this.handler.setConfiguredDevices(configuration.configuredDevices);
    }

    @Override
    public void setContext(ConversationContext<Ethernet_Frame> context) {
        super.setContext(context);
        this.profinetDriverContext = (ProfinetDriverContext) driverContext;
        try {
            PcapNetworkInterface devByAddress = Pcaps.getDevByAddress(InetAddress.getByName(this.configuration.transportConfig.split(":")[0]));
            this.profinetDriverContext.setChannel(new ProfinetChannel(Collections.singletonList(devByAddress)));
        } catch (UnknownHostException | PcapNativeException e) {
            throw new RuntimeException(e);
        }

        for (Map.Entry<String, ProfinetDevice> device : configuration.configuredDevices.entrySet()) {
            device.getValue().setContext(context, this.profinetDriverContext.getChannel());
        }
        try {
            onDeviceDiscovery();
        } catch (InterruptedException e) {
        }
    }

    private void onDeviceDiscovery() throws InterruptedException {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer(this.profinetDriverContext.getChannel());
        this.profinetDriverContext.getChannel().setDiscoverer(discoverer);
        DefaultPlcDiscoveryRequest request = new DefaultPlcDiscoveryRequest(
            discoverer,
            new LinkedHashMap<>()
        );

        // TODO:- Add handler for un-requested messages
        discoverer.ongoingDiscoverWithHandler(
            request,
            handler,
            5000L,
            30000L
        );
        waitForDeviceDiscovery();
    }

    private void waitForDeviceDiscovery() throws InterruptedException {
        // Once we receive an LLDP and PN-DCP message for each device move on.
        boolean discovered = false;
        int count = 0;
        while (!discovered) {
            discovered = true;
            for (Map.Entry<String, ProfinetDevice> device : configuration.configuredDevices.entrySet()) {
                if (!device.getValue().hasLldpPdu() || !device.getValue().hasDcpPdu()) {
                    discovered = false;
                }
            }
            if (!discovered) {
                Thread.sleep(3000L);
                count += 1;
            }
            if (count > 5) {
                break;
            }
        }
    }

    @Override
    public void onConnect(ConversationContext<Ethernet_Frame> context) {

        ////////////////////////////////////////////////////////////////////////////////////////////////////////////////
        // Initialize some important datastructures, that will be used a lot.
        // TODO: Possibly we can remove the ARP lookup and simply use the mac address in the connection-response.
        // Local connectivity attributes

        for (Map.Entry<String, ProfinetDevice> device : configuration.configuredDevices.entrySet()) {
            device.getValue().onConnect();
        }

        context.fireConnected();
        connected = true;
    }

    @Override
    public void close(ConversationContext<Ethernet_Frame> context) {
        // Nothing to do here ...
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        CompletableFuture<PlcReadResponse> future = new CompletableFuture<>();
        future.completeExceptionally(new NotImplementedException());
        return future;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        CompletableFuture<PlcWriteResponse> future = new CompletableFuture<>();
        future.completeExceptionally(new NotImplementedException());
        return future;
    }

    private Ethernet_FramePayload_PnDcp createProfinetCyclicDataRequest() {
        return new Ethernet_FramePayload_PnDcp(
            new PnDcp_Pdu_RealTimeCyclic(
                0x8000,
                new PnIo_CyclicServiceDataUnit((short) 0, (short) 0, (short) 0),
                16696,
                false,
                false,
                false,
                false,
                false,
                false));
    }


    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        CompletableFuture<PlcSubscriptionResponse> future = new CompletableFuture<>();
        if (!connected) {
            throw new RuntimeException("Not Connected");
        }

        try {
            Thread.sleep(10000L);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }

        return future;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        return null;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {

    }

    @Override
    protected void decode(ConversationContext<Ethernet_Frame> context, Ethernet_Frame msg) throws Exception {
        super.decode(context, msg);
    }


    private Optional<PcapNetworkInterface> getNetworkInterfaceForConnection(InetAddress address) {
        try {
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                // We're only interested in real running network interfaces, skip the rest.
                if (dev.isLoopBack() || !dev.isRunning() || dev.isUp()) {
                    continue;
                }

                for (PcapAddress curAddress : dev.getAddresses()) {

                }
            }
        } catch (PcapNativeException e) {
            logger.warn(String.format("Error finding network device for connection to %s", address.toString()), e);
        }
        return Optional.empty();
    }


}
