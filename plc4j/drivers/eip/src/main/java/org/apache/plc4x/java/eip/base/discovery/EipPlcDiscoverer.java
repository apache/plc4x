/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.eip.base.discovery;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.eip.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.*;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeUnit;

public class EipPlcDiscoverer implements PlcDiscoverer {

    private final Logger logger = LoggerFactory.getLogger(EipPlcDiscoverer.class);

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler) {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Queue<PlcDiscoveryItem> values = new ConcurrentLinkedQueue<>();

        // Send out a discovery request to every non-loopback device with IPv4 address.
        List<DatagramSocket> openSockets = new ArrayList<>();
        try {
            for (NetworkInterface networkInterface : Collections.list(NetworkInterface.getNetworkInterfaces())) {
                if (!networkInterface.isLoopback()) {
                    for (InterfaceAddress interfaceAddress : networkInterface.getInterfaceAddresses()) {
                        if ((interfaceAddress.getBroadcast() != null) && (interfaceAddress.getAddress() instanceof Inet4Address)) {
                            Inet4Address inet4Address = (Inet4Address) interfaceAddress.getAddress();
                            // Open a listening socket on the AMS discovery default port for taking in responses.
                            DatagramSocket discoverySocket = new DatagramSocket(EipConstants.EIPUDPDISCOVERYDEFAULTPORT, inet4Address);
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

                                        InetAddress plcAddress = packet.getAddress();
                                        ReadBuffer readBuffer = new ReadBufferByteBased(packet.getData(), ByteOrder.LITTLE_ENDIAN);
                                        try {
                                            EipPacket eipPacket = EipPacket.staticParse(readBuffer, true);

                                            // Check if this is actually a discovery response.
                                            if ((eipPacket.getCommand() == 0x0063) && (eipPacket.getResponse())){
                                                EipListIdentityResponse listIdentityResponse = (EipListIdentityResponse) eipPacket;
                                                for (CommandSpecificDataItem commandSpecificDataItem : listIdentityResponse.getItems()) {
                                                    if(commandSpecificDataItem instanceof CipIdentity) {
                                                        CipIdentity identityItem = (CipIdentity) commandSpecificDataItem;

                                                        // Add an entry to the results.
                                                        PlcDiscoveryItem plcDiscoveryItem = new DefaultPlcDiscoveryItem(
                                                            "eip", "tcp",
                                                            plcAddress.getHostAddress() + ":" + EipConstants.EIPTCPDEFAULTPORT,
                                                            Collections.emptyMap(), identityItem.getProductName(), Collections.emptyMap());

                                                        // If we've got an explicit handler, pass the new item to that.
                                                        if (handler != null) {
                                                            handler.handle(plcDiscoveryItem);
                                                        }

                                                        // Simply add the item to the list.
                                                        values.add(plcDiscoveryItem);
                                                    }
                                                }
                                            }
                                        } catch (ParseException e) {
                                            logger.error("Error parsing EIP discovery response", e);
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
                                // Create the discovery request message for this device.
                                EipPacket discoveryPacket = new EipListIdentityRequest(0, 0, new byte[] {0,0,0,0,0,0,0,0}, 0);

                                // Serialize the message.
                                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(discoveryPacket.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
                                discoveryPacket.serialize(writeBuffer);

                                // Get the broadcast address for this interface.
                                InetAddress broadcastAddress = interfaceAddress.getBroadcast();

                                // Create the UDP packet to the broadcast address.
                                DatagramPacket discoveryRequestPacket = new DatagramPacket(writeBuffer.getBytes(), writeBuffer.getBytes().length, broadcastAddress, EipConstants.EIPUDPDISCOVERYDEFAULTPORT);
                                discoverySocket.send(discoveryRequestPacket);
                            } catch (SerializationException e) {
                                logger.error("Error serializing EIP discovery request", e);
                            } catch (IOException e) {
                                logger.error("Error sending EIP discover request", e);
                            }

                            try {
                                Thread.sleep(3000);
                            } catch (InterruptedException e) {
                                throw new RuntimeException(e);
                            }
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
                    new DefaultPlcDiscoveryResponse(discoveryRequest, PlcResponseCode.OK, new ArrayList<>(values));
                future.complete(response);
            }
        }, 5000L);

        return future;
    }

    public static void main(String[] args) throws Exception {
        EipPlcDiscoverer discoverer = new EipPlcDiscoverer();
        CompletableFuture<PlcDiscoveryResponse> discover = discoverer.discover(null);
        PlcDiscoveryResponse plcDiscoveryResponse = discover.get(6000L, TimeUnit.MILLISECONDS);
        for (PlcDiscoveryItem value : plcDiscoveryResponse.getValues()) {
            System.out.println(value.getConnectionUrl() + " (" + value.getName() + ")");
        }
    }

}
