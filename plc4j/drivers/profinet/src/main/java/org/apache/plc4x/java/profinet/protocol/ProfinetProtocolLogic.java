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

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.pcap4j.core.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.*;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeoutException;
import java.util.function.Consumer;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> implements HasConfiguration<ProfinetConfiguration>, PlcSubscriber {

    private final Logger LOGGER = LoggerFactory.getLogger(ProfinetProtocolLogic.class);
    public static final Pattern SUB_MODULE_ARRAY_PATTERN = Pattern.compile("^\\[((\\[[\\w, ]*\\]){1}[ ,]{0,2})*\\]");
    public static final Pattern SUB_MODULE_SPLIT_ARRAY_PATTERN = Pattern.compile("(?:\\[(?:\\[([\\w, ]*)\\]){1}(?:[ ,]{0,2}))*\\]");
    public LinkedHashMap<String, ProfinetDevice> configuredDevices = new LinkedHashMap<>();
    private ProfinetDriverContext driverContext;

    public ProfinetProtocolLogic() {
        super();
        setDriverContext(new ProfinetDriverContext());
    }

    public void setDriverContext(ProfinetDriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.driverContext = driverContext;
    }

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {
        driverContext.setConfiguration(configuration);
    }

    @Override
    public void setContext(ConversationContext<Ethernet_Frame> context) {
        super.setContext(context);

        // Open the receiving UDP port.
        try {
            driverContext.setSocket(new DatagramSocket(ProfinetDriverContext.DEFAULT_UDP_PORT));
        } catch (SocketException e) {
            throw new RuntimeException(e);
        }

        driverContext.getHandler().setConfiguredDevices(configuredDevices);

        for (Map.Entry<String, ProfinetDevice> device : configuredDevices.entrySet()) {
            device.getValue().setContext(context, this.driverContext.getChannel());
        }
    }

    private void onDeviceDiscovery() throws InterruptedException {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer(this.driverContext.getChannel());
        this.driverContext.getChannel().setDiscoverer(discoverer);
        DefaultPlcDiscoveryRequest request = new DefaultPlcDiscoveryRequest(discoverer, new LinkedHashMap<>());

        discoverer.ongoingDiscoverWithHandler(request, driverContext.getHandler(), 5000L, 30000L);
        waitForDeviceDiscovery();
    }

    private void waitForDeviceDiscovery() throws InterruptedException {
        // Once we receive an LLDP and PN-DCP message for each device move on.
        boolean discovered = false;
        int count = 0;
        while (!discovered) {
            discovered = true;
            for (Map.Entry<String, ProfinetDevice> device : this.configuredDevices.entrySet()) {
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
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();
        Map<String, List<PlcBrowseItem>> values = new HashMap<>();
        Map<String, PlcResponseCode> codes = new HashMap<>();

        for (Map.Entry<String, ProfinetDevice> device : this.configuredDevices.entrySet()) {
            device.getValue().browseTags(values);
            codes.put(device.getKey(), PlcResponseCode.OK);
        }

        DefaultPlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, codes, values);
        future.complete(response);
        return future;
    }

    @Override
    public void onConnect(ConversationContext<Ethernet_Frame> context) {
        InetAddress localIpAddress = null;
        try {
            RawSocketChannel channel = (RawSocketChannel) context.getChannel();
            String localAddress = channel.getLocalAddress().toString().substring(1).split(":")[0];
            localIpAddress = InetAddress.getByName(localAddress);
            PcapNetworkInterface devByAddress = Pcaps.getDevByAddress(localIpAddress);
            driverContext.setChannel(new ProfinetChannel(Collections.singletonList(devByAddress)));
            driverContext.getChannel().setConfiguredDevices(this.configuredDevices);
        } catch (PcapNativeException | UnknownHostException e) {
            throw new RuntimeException(e);
        }

        try {
            onDeviceDiscovery();
        } catch (InterruptedException ignored) {}

        for (Map.Entry<String, ProfinetDevice> device : configuredDevices.entrySet()) {
            device.getValue().getDeviceContext().setChannel(driverContext.getChannel());
            device.getValue().getDeviceContext().setLocalIpAddress(localIpAddress);
        }

        try {
            for (Map.Entry<String, ProfinetDevice> device : this.configuredDevices.entrySet()) {
                device.getValue().onConnect(this);
            }
            context.fireConnected();

        } catch (ExecutionException | InterruptedException | TimeoutException e) {
            throw new RuntimeException(e);
        }
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

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        return CompletableFuture.supplyAsync(() -> {
            Map<String, ResponseItem<PlcSubscriptionHandle>> values = new HashMap<>();
            long subscriptionId = 0;
            ArrayList<String> fields = new ArrayList<>(subscriptionRequest.getTagNames());
            long cycleTime = (subscriptionRequest.getTag(fields.get(0))).getDuration().orElse(Duration.ofMillis(1000)).toMillis();

            for (String fieldName : subscriptionRequest.getTagNames()) {
                final DefaultPlcSubscriptionTag fieldDefaultPlcSubscription = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(fieldName);
                if (!(fieldDefaultPlcSubscription.getTag() instanceof ProfinetTag)) {
                    values.put(fieldName, new ResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                } else {
                    values.put(fieldName, new ResponseItem<>(PlcResponseCode.OK, driverContext.getSubscriptions().get(subscriptionId)));
                }
            }
            return new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
        });
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        List<PlcConsumerRegistration> registrations = new LinkedList<>();
        // Register the current consumer for each of the given subscription handles
        for (PlcSubscriptionHandle subscriptionHandle : handles) {
            LOGGER.debug("Registering Consumer");
            final PlcConsumerRegistration consumerRegistration = subscriptionHandle.register(consumer);
            registrations.add(consumerRegistration);
        }
        return new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {

    }

    @Override
    protected void decode(ConversationContext<Ethernet_Frame> context, Ethernet_Frame msg) throws Exception {
        super.decode(context, msg);
    }
}
