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

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.profinet.config.ConfigurationProfinetDevice;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.context.ProfinetDriverContext;
import org.apache.plc4x.java.profinet.device.*;
import org.apache.plc4x.java.profinet.discovery.ProfinetPlcDiscoverer;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.profinet.tag.ProfinetTag;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;
import org.pcap4j.core.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.*;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeoutException;

public class ProfinetProtocolLogic extends Plc4xProtocolBase<Ethernet_Frame> implements HasConfiguration<ProfinetConfiguration> {

    private final Logger LOGGER = LoggerFactory.getLogger(ProfinetProtocolLogic.class);

    private ProfinetDriverContext driverContext;
    private Map<String, ProfinetDevice> devices = new HashMap<>();

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

        Map<String, ConfigurationProfinetDevice> configuredDevices = configuration.getDevices().getConfiguredDevices();

        for (Map.Entry<String, ConfigurationProfinetDevice> entry : configuredDevices.entrySet()) {
            devices.put(entry.getKey(),
                new ProfinetDevice(
                    new ProfinetMessageWrapper(),
                    entry.getValue().getDevicename(),
                    entry.getValue().getDeviceaccess(),
                    entry.getValue().getSubmodules(),
                    entry.getValue().getGsdHandler()
                )
            );
            devices.get(entry.getValue().getDevicename()).setIpAddress(entry.getValue().getIpaddress());
        }

        driverContext.setHandler(new ProfinetDeviceMessageHandler(devices));
        for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
            device.getValue().getDeviceContext().setConfiguration(configuration);
        }
    }

    @Override
    public void setContext(ConversationContext<Ethernet_Frame> context) {
        super.setContext(context);

        // Open the receiving UDP port and keep it open.
        try {
            driverContext.setSocket(new DatagramSocket(ProfinetDriverContext.DEFAULT_UDP_PORT));
        } catch (SocketException e) {
            throw new RuntimeException(e);
        }

        driverContext.getHandler().setConfiguredDevices(devices);

        for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
            device.getValue().setContext(context, this.driverContext.getChannel());
        }
    }

    private void onDeviceDiscovery() throws InterruptedException, PlcConnectionException {
        ProfinetPlcDiscoverer discoverer = new ProfinetPlcDiscoverer(driverContext.getChannel());
        driverContext.getChannel().setDiscoverer(discoverer);
        DefaultPlcDiscoveryRequest request = new DefaultPlcDiscoveryRequest(discoverer, new LinkedHashMap<>());
        discoverer.ongoingDiscoverWithHandler(request, driverContext.getHandler(), 5000L, 30000L);
        waitForDeviceDiscovery();
    }

    private void waitForDeviceDiscovery() throws InterruptedException, PlcConnectionException {
        // Once we receive an LLDP and PN-DCP message for each device move on.
        boolean discovered = false;
        int count = 0;
        while (!discovered) {
            discovered = true;
            for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
                if (!device.getValue().hasLldpPdu() || !device.getValue().hasDcpPdu()) {
                    discovered = false;
                }
            }
            if (!discovered) {
                Thread.sleep(3000L);
                count += 1;
            }
            if (count > 5) {
                throw new PlcConnectionException("One device failed to respond to discovery packet");
            }
        }
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        CompletableFuture<PlcBrowseResponse> future = new CompletableFuture<>();
        List<PlcBrowseItem> values = new ArrayList<>();
        Map<String, PlcResponseCode> codes = new HashMap<>();
        Map<String, List<PlcBrowseItem>> responseValues = new HashMap<>();

        for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
            device.getValue().browseTags(values);
        }

        for (String queryname : browseRequest.getQueryNames()) {
            responseValues.put(queryname, values);
        }

        DefaultPlcBrowseResponse response = new DefaultPlcBrowseResponse(browseRequest, codes, responseValues);
        future.complete(response);
        return future;
    }

    @Override
    public void onConnect(ConversationContext<Ethernet_Frame> context) {
        InetAddress localIpAddress;
        try {
            RawSocketChannel channel = (RawSocketChannel) context.getChannel();
            String localAddress = channel.getLocalAddress().toString().substring(1).split(":")[0];
            localIpAddress = InetAddress.getByName(localAddress);
            PcapNetworkInterface devByAddress = Pcaps.getDevByAddress(localIpAddress);
            driverContext.setChannel(new ProfinetChannel(Collections.singletonList(devByAddress), devices));
            driverContext.getChannel().setConfiguredDevices(devices);
            for (Map.Entry<String, ProfinetDevice> entry : devices.entrySet()) {
                entry.getValue().setNetworkInterface(new ProfinetNetworkInterface(devByAddress));
            }
        } catch (PcapNativeException | UnknownHostException e) {
            throw new RuntimeException(e);
        }

        try {
            onDeviceDiscovery();
        } catch (PlcException | InterruptedException e) {
            throw new RuntimeException(e);
        }

        for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
            device.getValue().getDeviceContext().setChannel(driverContext.getChannel());
            device.getValue().getDeviceContext().setLocalIpAddress(localIpAddress);
        }

        try {
            for (Map.Entry<String, ProfinetDevice> device : devices.entrySet()) {
                device.getValue().onConnect();
            }
            context.fireConnected();

        } catch (ExecutionException | InterruptedException | TimeoutException e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void close(ConversationContext<Ethernet_Frame> context) {
        // TODO:- Do something here
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

            for (String fieldName : subscriptionRequest.getTagNames()) {
                PlcSubscriptionTag tag = subscriptionRequest.getTag(fieldName);
                final DefaultPlcSubscriptionTag fieldDefaultPlcSubscription = (DefaultPlcSubscriptionTag) subscriptionRequest.getTag(fieldName);
                String deviceString = fieldDefaultPlcSubscription.getAddressString().split("\\.")[0].toUpperCase();
                ProfinetDevice device = devices.get(deviceString);

                ProfinetSubscriptionHandle subscriptionHandle = new ProfinetSubscriptionHandle(device, fieldName, tag);
                device.getDeviceContext().addSubscriptionHandle(fieldDefaultPlcSubscription.getAddressString(), subscriptionHandle);

                if (!(fieldDefaultPlcSubscription.getTag() instanceof ProfinetTag)) {
                    values.put(fieldName, new ResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
                } else {
                    values.put(fieldName, new ResponseItem<>(PlcResponseCode.OK, subscriptionHandle));
                }
            }
            return new DefaultPlcSubscriptionResponse(subscriptionRequest, values);
        });
    }

    @Override
    protected void decode(ConversationContext<Ethernet_Frame> context, Ethernet_Frame msg) throws Exception {
        super.decode(context, msg);
    }
}
