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
import org.apache.plc4x.java.knxnetip.readwrite.HPAIDiscoveryEndpoint;
import org.apache.plc4x.java.knxnetip.readwrite.HostProtocolCode;
import org.apache.plc4x.java.knxnetip.readwrite.IPAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpMessage;
import org.apache.plc4x.java.knxnetip.readwrite.SearchRequest;
import org.apache.plc4x.java.knxnetip.readwrite.SearchResponse;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
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

    private final Logger logger = LoggerFactory.getLogger(KnxNetIpPlcDiscoverer.class);

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler) {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Map<String, PlcDiscoveryItem> values = new ConcurrentHashMap<>();

        // Send out a discovery request to every non-loopback device with IPv4 address.
        List<DatagramSocket> openSockets = new ArrayList<>();
        try {
            for (NetworkInterface networkInterface : Collections.list(NetworkInterface.getNetworkInterfaces())) {
                if (!networkInterface.isLoopback()) {
                    for (InterfaceAddress interfaceAddress : networkInterface.getInterfaceAddresses()) {
                        if ((interfaceAddress.getBroadcast() != null) && (interfaceAddress.getAddress() instanceof Inet4Address)) {
                            Inet4Address inet4Address = (Inet4Address) interfaceAddress.getAddress();
                            // Open a listening socket on the AMS discovery default port for taking in responses.
                            DatagramSocket discoverySocket = new DatagramSocket(KnxNetIpDriver.KNXNET_IP_PORT, inet4Address);
                            discoverySocket.setBroadcast(true);

                            openSockets.add(discoverySocket);

                            // Start listening for incoming messages.
                            Thread thread = new Thread(() -> {
                                try {
                                    while (true) {
                                        // Wait for an incoming packet.
                                        byte[] buffer = new byte[512];
                                        DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                                        discoverySocket.receive(packet);

                                        ReadBuffer readBuffer = new ReadBufferByteBased(packet.getData());
                                        try {
                                            KnxNetIpMessage knxNetIpMessage = KnxNetIpMessage.staticParse(readBuffer);
                                            if(knxNetIpMessage instanceof SearchResponse) {
                                                SearchResponse searchResponse = (SearchResponse) knxNetIpMessage;
                                                IPAddress ipAddress = searchResponse.getHpaiControlEndpoint().getIpAddress();
                                                int port = searchResponse.getHpaiControlEndpoint().getIpPort();
                                                String name = new String(searchResponse.getDibDeviceInfo().getDeviceFriendlyName()).trim();

                                                PlcDiscoveryItem plcDiscoveryItem = new DefaultPlcDiscoveryItem(
                                                    "knxnet-ip",
                                                    "udp",
                                                    InetAddress.getByAddress(ipAddress.getAddr()).toString().substring(1) + ":" + port,
                                                    Collections.emptyMap(),
                                                    name,
                                                    Collections.emptyMap());

                                                // If we've got an explicit handler, pass the new item to that.
                                                if ((handler != null) && !values.containsKey(plcDiscoveryItem.getConnectionUrl())){
                                                    handler.handle(plcDiscoveryItem);
                                                }

                                                // Simply add the item to the list.
                                                values.put(plcDiscoveryItem.getConnectionUrl(), plcDiscoveryItem);
                                            }
                                        } catch (ParseException e) {
                                            throw new RuntimeException(e);
                                        }
                                    }
                                } catch (SocketException e) {
                                    // If we're closing the socket at the end, a "Socket closed"
                                    // exception is thrown.
                                    if(!"Socket closed".equals(e.getMessage())) {
                                        logger.error("Error receiving EIP discovery response", e);
                                    }
                                } catch (IOException e) {
                                    logger.error("Error reading EIP discovery response", e);
                                }
                            });
                            thread.start();

                            // Send the discovery request.
                            try {
                                // TODO: Replace with the local ip address and the local udp port number.

                                SearchRequest searchRequest = new SearchRequest(new HPAIDiscoveryEndpoint(HostProtocolCode.IPV4_UDP, new IPAddress(discoverySocket.getLocalAddress().getAddress()), discoverySocket.getLocalPort()));

                                // Serialize the message.
                                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(searchRequest.getLengthInBytes()/*, ByteOrder.LITTLE_ENDIAN*/);
                                searchRequest.serialize(writeBuffer);

                                // Get the broadcast address for this interface.
                                InetAddress knxDiscoveryAddress = InetAddress.getByAddress(new byte[]{(byte) 224, (byte) 0, (byte) 23, (byte) 12});

                                // Create the UDP packet to the broadcast address.
                                DatagramPacket discoveryRequestPacket = new DatagramPacket(writeBuffer.getBytes(), writeBuffer.getBytes().length, knxDiscoveryAddress, KnxNetIpDriver.KNXNET_IP_PORT);
                                discoverySocket.send(discoveryRequestPacket);
                            } catch (SerializationException e) {
                                logger.error("Error serializing EIP discovery request", e);
                            } catch (IOException e) {
                                logger.error("Error sending EIP discover request", e);
                            }
                        }
                        try {
                            Thread.sleep(100);
                        } catch (InterruptedException e) {
                            throw new RuntimeException(e);
                        }
                    }
                }
            }
        } catch (SocketException e) {
            throw new RuntimeException(e);
        } finally {
            for (DatagramSocket openSocket : openSockets) {
                openSocket.close();
            }
        }

        // Create a timer that completes the future after a given time with all the responses it found till then.
        Timer timer = new Timer("Discovery Timeout");
        timer.schedule(new TimerTask() {
            public void run() {
                PlcDiscoveryResponse response =
                    new DefaultPlcDiscoveryResponse(discoveryRequest, PlcResponseCode.OK, new ArrayList<>(values.values()));
                timer.cancel();
                timer.purge();
                future.complete(response);
            }
        }, 5000L);

        return future;
    }

    public static void main(String[] args) throws Exception {
        KnxNetIpPlcDiscoverer discoverer = new KnxNetIpPlcDiscoverer();
        CompletableFuture<PlcDiscoveryResponse> discover = discoverer.discover(null);
        PlcDiscoveryResponse plcDiscoveryResponse = discover.get(6000L, TimeUnit.MILLISECONDS);
        for (PlcDiscoveryItem value : plcDiscoveryResponse.getValues()) {
            System.out.println(value.getConnectionUrl() + " (" + value.getName() + ")");
        }
    }

}
