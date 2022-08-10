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
package org.apache.plc4x.java.ads.discovery;

import org.apache.plc4x.java.ads.discovery.readwrite.*;
import org.apache.plc4x.java.ads.readwrite.AdsConstants;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.*;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.TimeUnit;

public class AdsPlcDiscoverer implements PlcDiscoverer {

    private final Logger logger = LoggerFactory.getLogger(AdsPlcDiscoverer.class);

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
                        if (interfaceAddress.getAddress() instanceof Inet4Address) {
                            Inet4Address inet4Address = (Inet4Address) interfaceAddress.getAddress();
                            // Open a listening socket on the AMS discovery default port for taking in responses.
                            DatagramSocket adsDiscoverySocket = new DatagramSocket(AdsDiscoveryConstants.ADSDISCOVERYUDPDEFAULTPORT, inet4Address);
                            adsDiscoverySocket.setBroadcast(true);

                            openSockets.add(adsDiscoverySocket);

                            // Start listening for incoming messages.
                            Thread thread = new Thread(() -> {
                                try {
                                    while (true) {
                                        // Wait for an incoming packet.
                                        byte[] buffer = new byte[512];
                                        DatagramPacket packet = new DatagramPacket(buffer, buffer.length);
                                        adsDiscoverySocket.receive(packet);

                                        InetAddress plcAddress = packet.getAddress();
                                        ReadBuffer readBuffer = new ReadBufferByteBased(packet.getData(), ByteOrder.LITTLE_ENDIAN);
                                        AdsDiscovery adsDiscoveryResponse = AdsDiscovery.staticParse(readBuffer);

                                        // Check if this is actually a discovery response.
                                        if ((adsDiscoveryResponse.getRequestId() == 0) &&
                                            (adsDiscoveryResponse.getPortNumber() == AdsPortNumbers.SYSTEM_SERVICE) &&
                                            (adsDiscoveryResponse.getOperation() == Operation.DISCOVERY_RESPONSE)) {

                                            AmsNetId remoteAmsNetId = adsDiscoveryResponse.getAmsNetId();
                                            AdsDiscoveryBlockHostName hostNameBlock = null;
                                            AdsDiscoveryBlockOsData osDataBlock = null;
                                            AdsDiscoveryBlockVersion versionBlock = null;
                                            AdsDiscoveryBlockFingerprint fingerprintBlock = null;
                                            for (AdsDiscoveryBlock block : adsDiscoveryResponse.getBlocks()) {
                                                switch (block.getBlockType()) {
                                                    case HOST_NAME:
                                                        hostNameBlock = (AdsDiscoveryBlockHostName) block;
                                                        break;
                                                    case OS_DATA:
                                                        osDataBlock = (AdsDiscoveryBlockOsData) block;
                                                        break;
                                                    case VERSION:
                                                        versionBlock = (AdsDiscoveryBlockVersion) block;
                                                        break;
                                                    case FINGERPRINT:
                                                        fingerprintBlock = (AdsDiscoveryBlockFingerprint) block;
                                                        break;
                                                    default:
                                                        logger.info(String.format("Unexpected block type: %s", block.getBlockType().toString()));
                                                }
                                            }

                                            if (hostNameBlock != null) {
                                                Map<String, String> options = new HashMap<>();
                                                options.put("sourceAmsNetId", "65534");
                                                options.put("sourceAmsPort", inet4Address.getHostAddress() + ".1.1");
                                                options.put("targetAmsNetId", remoteAmsNetId.getOctet1() + "." + remoteAmsNetId.getOctet2() + "." + remoteAmsNetId.getOctet3() + "." + remoteAmsNetId.getOctet4() + "." + remoteAmsNetId.getOctet5() + "." + remoteAmsNetId.getOctet6());
                                                // TODO: Check if this is legit, or if we can get the information from somewhere.
                                                options.put("targetAmsPort", "851");

                                                Map<String, PlcValue> attributes = new HashMap<>();
                                                attributes.put("hostName", new PlcSTRING(hostNameBlock.getHostName().getText()));
                                                if (versionBlock != null) {
                                                    byte[] versionData = versionBlock.getVersionData();
                                                    int patchVersion = ((int) versionData[3] & 0xFF) << 8 | ((int) versionData[2] & 0xFF);
                                                    attributes.put("twinCatVersion", new PlcSTRING(String.format("%d.%d.%d", (short) versionData[0] & 0xFF, (short) versionData[1] & 0xFF, patchVersion)));
                                                }
                                                if (fingerprintBlock != null) {
                                                    attributes.put("fingerprint", new PlcSTRING(new String(fingerprintBlock.getData())));
                                                }
                                                // TODO: Find out how to handle the OS Data

                                                // Add an entry to the results.
                                                PlcDiscoveryItem plcDiscoveryItem = new DefaultPlcDiscoveryItem(
                                                    "ads", "tcp",
                                                    plcAddress.getHostAddress() + ":" + AdsConstants.ADSTCPDEFAULTPORT,
                                                    options, hostNameBlock.getHostName().getText(), attributes);

                                                // If we've got an explicit handler, pass the new item to that.
                                                if (handler != null) {
                                                    handler.handle(plcDiscoveryItem);
                                                }

                                                // Simply add the item to the list.
                                                values.add(plcDiscoveryItem);
                                            }
                                        }
                                    }
                                } catch (SocketException e) {
                                    // If we're closing the socket at the end, a "Socket closed"
                                    // exception is thrown.
                                    if(!"Socket closed".equals(e.getMessage())) {
                                        logger.error("Error receiving ADS discovery response", e);
                                    }
                                } catch (IOException e) {
                                    logger.error("Error reading ADS discovery response", e);
                                } catch (ParseException e) {
                                    logger.error("Error parsing ADS discovery response", e);
                                }
                            });
                            thread.start();

                            // Send the discovery request.
                            try {
                                // Create the discovery request message for this device.
                                AmsNetId amsNetId = new AmsNetId(inet4Address.getAddress()[0], inet4Address.getAddress()[1], inet4Address.getAddress()[2], inet4Address.getAddress()[3], (byte) 1, (byte) 1);
                                AdsDiscovery discoveryRequestMessage = new AdsDiscovery(0, Operation.DISCOVERY_REQUEST, amsNetId, AdsPortNumbers.SYSTEM_SERVICE, Collections.emptyList());

                                // Serialize the message.
                                WriteBufferByteBased writeBuffer = new WriteBufferByteBased(discoveryRequestMessage.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
                                discoveryRequestMessage.serialize(writeBuffer);

                                // Get the broadcast address for this interface.
                                InetAddress broadcastAddress = interfaceAddress.getBroadcast();

                                // Create the UDP packet to the broadcast address.
                                DatagramPacket discoveryRequestPacket = new DatagramPacket(writeBuffer.getBytes(), writeBuffer.getBytes().length, broadcastAddress, AdsDiscoveryConstants.ADSDISCOVERYUDPDEFAULTPORT);
                                adsDiscoverySocket.send(discoveryRequestPacket);
                            } catch (SerializationException e) {
                                logger.error("Error serializing ADS discovery request", e);
                            } catch (IOException e) {
                                logger.error("Error sending ADS discover request", e);
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
        AdsPlcDiscoverer discoverer = new AdsPlcDiscoverer();
        CompletableFuture<PlcDiscoveryResponse> discover = discoverer.discover(null);
        PlcDiscoveryResponse plcDiscoveryResponse = discover.get(6000L, TimeUnit.MILLISECONDS);
        System.out.println(plcDiscoveryResponse);
    }

}
