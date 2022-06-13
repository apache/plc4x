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
package org.apache.plc4x.java.modbus.tcp.discovery;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryResponse;
import org.apache.plc4x.java.spi.messages.PlcDiscoverer;
import org.apache.plc4x.java.utils.rawsockets.netty.utils.ArpUtils;
import org.pcap4j.core.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetAddress;
import java.net.Socket;
import java.net.UnknownHostException;
import java.nio.ByteBuffer;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.function.Function;
import java.util.function.Predicate;
import java.util.stream.Collectors;

public class ModbusPlcDiscoverer implements PlcDiscoverer {

    private final Logger logger = LoggerFactory.getLogger(ModbusPlcDiscoverer.class);

    public static <T> Predicate<T> distinctByKey(Function<? super T, ?> keyExtractor) {
        Set<Object> seen = ConcurrentHashMap.newKeySet();
        return t -> seen.add(keyExtractor.apply(t));
    }
    @Override
    public CompletableFuture<PlcDiscoveryResponse> discover(PlcDiscoveryRequest discoveryRequest) {
        return discoverWithHandler(discoveryRequest, null);
    }

    @Override
    public CompletableFuture<PlcDiscoveryResponse> discoverWithHandler(PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler) {
        // Get a list of all reachable IP addresses from the current system.
        // TODO: add an option to fine tune the network device or ip subnet to scan and maybe some timeouts and delays to prevent flooding.
        final CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Thread discoveryThread = new Thread(() -> {
            executeDiscovery(future, discoveryRequest, handler);
        });
        discoveryThread.start();
        return future;
    }

    private void executeDiscovery(CompletableFuture<PlcDiscoveryResponse> future, PlcDiscoveryRequest discoveryRequest, PlcDiscoveryItemHandler handler) {
        List<InetAddress> possibleAddresses = new ArrayList<>();
        try {
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                logger.info("Scanning network {} for alive IP addresses", dev.getName());
                final Set<InetAddress> inetAddresses = ArpUtils.scanNetworkDevice(dev);
                logger.debug("Found {} addresses: {}", inetAddresses.size(), inetAddresses);
                possibleAddresses.addAll(inetAddresses);
            }
        } catch (PcapNativeException e) {
            logger.error("Error collecting list of possible IP addresses", e);
            future.complete(new DefaultPlcDiscoveryResponse(
                discoveryRequest, PlcResponseCode.INTERNAL_ERROR, Collections.emptyList()));
            return;
        }
        try {
            possibleAddresses.add(InetAddress.getByName("localhost"));
        } catch (UnknownHostException e) {
            e.printStackTrace();
        }

        // Filter out duplicates.
        possibleAddresses = possibleAddresses.stream().filter(
            distinctByKey(InetAddress::getHostAddress)).collect(Collectors.toList());

        Queue<PlcDiscoveryItem> discoveryItems = new ConcurrentLinkedQueue<>();
        possibleAddresses.stream().parallel().forEach(possibleAddress -> {
            try {
                logger.info("Trying address: {}", possibleAddress);
                // Try to get a connection to the given host and port.
                Socket socket = new Socket(possibleAddress.getHostAddress(), ModbusConstants.MODBUSTCPDEFAULTPORT);

                logger.info("Connected: {}", possibleAddress);

                // Send the request to the target device
                final OutputStream outputStream = socket.getOutputStream();
                final InputStream inputStream = new BufferedInputStream(socket.getInputStream());

                // As we not only need to provide the IP but also the unit-identifier, we need
                // to iterate over all possible values until we find one that works.
                // Unfortunately the way devices react to invalid requests differ:
                // - My heating system doesn't sends an error response for an invalid uint-identifier
                // - The Modbus Server on a S7 simply accepts any unit-identifier
                // - Modbus pal only responds to correct unit-identifiers and doesn't send anything
                //   for invalid ones.
                //
                // So-far the only way I have found to reliably check if a Modbus device exists,
                // was by trying to read a coil or register. Even if Modbus generally supports
                // commands for diagnosing connections, it turns out none of these were actually
                // supported by any device I came across. The spec is a bit unclear here, but
                // it seems as if these are only supported on Serial (Modbus RTU)
                // TODO: We should probably not only try to read a coil, but try any of the types and if one works, that's a match.
                // Possibly we can fine tune this to speed up things.
                int transactionIdentifier = 1;
                for(short unitIdentifier = 1; unitIdentifier <= 247; unitIdentifier++) {
                    ModbusTcpADU packet = new ModbusTcpADU(transactionIdentifier++, unitIdentifier,
                        new ModbusPDUReadCoilsRequest(1, 1), false);
                    byte[] deviceIdentificationBytes = null;
                    try {
                        // Serialize the request to its byte form.
                        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(packet.getLengthInBytes());
                        packet.serialize(writeBuffer);
                        deviceIdentificationBytes = writeBuffer.getBytes();
                    } catch (SerializationException e) {
                        logger.error("Error creating the device identification request", e);
                    }
                    if (deviceIdentificationBytes == null) {
                        future.complete(new DefaultPlcDiscoveryResponse(
                            discoveryRequest, PlcResponseCode.INTERNAL_ERROR, Collections.emptyList()));
                        return;
                    }
                    byte[] finalDeviceIdentificationBytes = deviceIdentificationBytes;

                    outputStream.write(finalDeviceIdentificationBytes);
                    outputStream.flush();

                    // Wait for a response.
                    byte[] responseBytes = null;
                    final long endTime = System.currentTimeMillis() + 100;
                    while (responseBytes == null) {
                        // If we've got enough bytes to find out the size of the packet, try to check this.
                        if (inputStream.available() >= 6) {
                            inputStream.mark(6);
                            inputStream.skip(4);
                            byte[] packetLengthBytes = new byte[2];
                            int bytesRead = inputStream.read(packetLengthBytes);
                            inputStream.reset();
                            // Only if we really read 2 bytes, does it make sense to continue using it.
                            if (bytesRead != 2) {
                                continue;
                            }
                            final short packetLength = (short) (ByteBuffer.wrap(packetLengthBytes).getShort() + 6);
                            if (inputStream.available() >= packetLength) {
                                responseBytes = new byte[packetLength];
                                bytesRead = inputStream.read(responseBytes);
                                if (bytesRead != packetLength) {
                                    responseBytes = null;
                                    break;
                                }
                            }
                        } else {
                            try {
                                Thread.sleep(1);
                            } catch (InterruptedException e) {
                                Thread.currentThread().interrupt();
                            }
                            if (System.currentTimeMillis() > endTime) {
                                break;
                            }
                        }
                    }
                    if (responseBytes != null) {
                        ReadBuffer readBuffer = new ReadBufferByteBased(responseBytes);
                        try {
                            ModbusTcpADU response = (ModbusTcpADU) ModbusTcpADU.staticParse(readBuffer, DriverType.MODBUS_TCP, true);
                            PlcDiscoveryItem discoveryItem;
                            if (!response.getPdu().getErrorFlag()) {
                                discoveryItem = new DefaultPlcDiscoveryItem(
                                    "modbus-tcp", "tcp", possibleAddress.getHostAddress(), Collections.singletonMap("unit-identifier", Integer.toString(unitIdentifier)), "unknown");
                                discoveryItems.add(discoveryItem);

                                // Give a handler the chance to react on the found device.
                                if(handler != null) {
                                    handler.handle(discoveryItem);
                                }
                                break;
                            }
                        } catch (ParseException e) {
                            // Ignore.
                        }
                    }
                }
            } catch (IOException e) {
                // Well this is actually sort of normal in case of us trying to connect to a non-existent device.
            }
        });

        future.complete(new DefaultPlcDiscoveryResponse(discoveryRequest, PlcResponseCode.OK, Arrays.asList(discoveryItems.toArray(new PlcDiscoveryItem[0]))));
    }

}
