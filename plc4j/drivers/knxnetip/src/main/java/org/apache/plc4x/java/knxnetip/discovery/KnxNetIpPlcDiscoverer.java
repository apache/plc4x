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
package org.apache.plc4x.java.knxnetip.discovery;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.knxnetip.KnxNetIpDriver;
import org.apache.plc4x.java.knxnetip.KnxNetIpMessageCodec;
import org.apache.plc4x.java.knxnetip.readwrite.HPAIDiscoveryEndpoint;
import org.apache.plc4x.java.knxnetip.readwrite.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.IPAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpMessage;
import org.apache.plc4x.java.knxnetip.readwrite.SearchRequest;
import org.apache.plc4x.java.knxnetip.readwrite.SearchResponse;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.InterfaceAddress;
import java.net.NetworkInterface;
import java.net.SocketException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;

public class KnxNetIpPlcDiscoverer implements PlcDiscoverer {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpPlcDiscoverer.class);
    private static final long DISCOVERY_TIMEOUT_MS = 5000L;
    private static final byte[] KNX_MULTICAST = new byte[]{(byte) 224, 0, 23, 12};

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest,
                                                                       PlcDiscoveryItemHandler handler) {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Map<String, PlcDiscoveryItem> values = new ConcurrentHashMap<>();
        List<DatagramSocket> openSockets = new ArrayList<>();

        try {
            for (NetworkInterface networkInterface : Collections.list(NetworkInterface.getNetworkInterfaces())) {
                if (networkInterface.isLoopback()) {
                    continue;
                }
                for (InterfaceAddress interfaceAddress : networkInterface.getInterfaceAddresses()) {
                    if (interfaceAddress.getBroadcast() == null) {
                        continue;
                    }
                    if (!(interfaceAddress.getAddress() instanceof Inet4Address inet4Address)) {
                        continue;
                    }
                    DatagramSocket discoverySocket = new DatagramSocket(KnxNetIpDriver.KNXNET_IP_PORT, inet4Address);
                    discoverySocket.setBroadcast(true);
                    openSockets.add(discoverySocket);

                    Thread receiver = new Thread(() -> receiveLoop(discoverySocket, values, handler),
                        "KnxDiscovery-" + inet4Address.getHostAddress());
                    receiver.setDaemon(true);
                    receiver.start();

                    sendSearchRequest(discoverySocket);
                }
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    break;
                }
            }
        } catch (SocketException e) {
            future.completeExceptionally(new RuntimeException("Error preparing KNXnet/IP discovery sockets", e));
            return future;
        }

        Timer timer = new Timer("KnxDiscoveryTimeout", true);
        timer.schedule(new TimerTask() {
            @Override
            public void run() {
                openSockets.forEach(DatagramSocket::close);
                future.complete(new DefaultPlcDiscoveryResponse(
                    discoveryRequest, PlcResponseCode.OK, new ArrayList<>(values.values())));
                timer.cancel();
            }
        }, DISCOVERY_TIMEOUT_MS);

        return future;
    }

    private void receiveLoop(DatagramSocket socket,
                             Map<String, PlcDiscoveryItem> values,
                             PlcDiscoveryItemHandler handler) {
        byte[] buffer = new byte[512];
        try {
            while (!socket.isClosed()) {
                DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                socket.receive(packet);
                try {
                    ReadBufferByteBased readBuffer = new ReadBufferByteBased(packet.getData(),
                        KnxNetIpMessageCodec.BUFFER_OPTIONS);
                    KnxNetIpMessage message = KnxNetIpMessage.staticParse(readBuffer);
                    if (!(message instanceof SearchResponse searchResponse)) {
                        continue;
                    }
                    IPAddress ipAddress = searchResponse.getHpaiControlEndpoint().getIpAddress();
                    int port = searchResponse.getHpaiControlEndpoint().getIpPort();
                    String name = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName()).trim();
                    String transportUrl = InetAddress.getByAddress(ipAddress.getAddr()).getHostAddress() + ":" + port;
                    PlcDiscoveryItem item = new DefaultPlcDiscoveryItem(
                        "knxnet-ip", "udp", transportUrl,
                        Collections.emptyMap(), name, Collections.emptyMap());
                    if (handler != null && !values.containsKey(item.getConnectionUrl())) {
                        handler.handle(item);
                    }
                    values.put(item.getConnectionUrl(), item);
                } catch (BufferException e) {
                    LOGGER.warn("Failed to parse incoming KNXnet/IP discovery datagram", e);
                }
            }
        } catch (SocketException e) {
            if (!"Socket closed".equals(e.getMessage())) {
                LOGGER.error("Error receiving KNXnet/IP discovery response", e);
            }
        } catch (IOException e) {
            LOGGER.error("Error reading KNXnet/IP discovery response", e);
        }
    }

    private void sendSearchRequest(DatagramSocket discoverySocket) {
        try {
            SearchRequest searchRequest = new SearchRequest(
                new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP,
                    new IPAddress(discoverySocket.getLocalAddress().getAddress()),
                    discoverySocket.getLocalPort()));
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(
                new byte[searchRequest.getLengthInBytes()], KnxNetIpMessageCodec.BUFFER_OPTIONS);
            searchRequest.serialize(writeBuffer);
            InetAddress knxDiscoveryAddress = InetAddress.getByAddress(KNX_MULTICAST);
            DatagramPacket packet = new DatagramPacket(writeBuffer.getBytes(), writeBuffer.getBytes().length,
                knxDiscoveryAddress, KnxNetIpDriver.KNXNET_IP_PORT);
            discoverySocket.send(packet);
        } catch (BufferException e) {
            LOGGER.error("Error serializing KNXnet/IP discovery request", e);
        } catch (IOException e) {
            LOGGER.error("Error sending KNXnet/IP discovery request", e);
        }
    }

    public static void main(String[] args) throws Exception {
        KnxNetIpPlcDiscoverer discoverer = new KnxNetIpPlcDiscoverer();
        PlcDiscoveryResponse response = discoverer.discover(null).get(10, TimeUnit.SECONDS);
        for (PlcDiscoveryItem value : response.getValues()) {
            System.out.println(value.getConnectionUrl() + " (" + value.getName() + ")");
        }
    }

}
