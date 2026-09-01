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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.functions.PlcDiscoverer;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryItem;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcDiscoveryResponse;
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

    /** MBAP header: transaction id, protocol id and the length field itself. */
    private static final int MODBUS_TCP_HEADER_SIZE = 6;

    /**
     * The largest a Modbus TCP ADU can be - the MBAP header, a unit identifier and a PDU of at most
     * 253 bytes. A declared length past this is not a Modbus response, whatever else it may be.
     */
    private static final int MODBUS_TCP_MAX_ADU_SIZE = 260;

    /**
     * Works out how long the ADU is from the two bytes of the MBAP length field.
     *
     * <p>Unsigned, and in int arithmetic. Read as a signed short and adjusted by the header size, a
     * declared length near the top of the field came out negative - and the next thing done with it
     * was to allocate an array that long.</p>
     *
     * @return the total ADU length, or -1 if the field cannot describe a Modbus ADU
     */
    static int aduLength(byte[] packetLengthBytes) {
        int declared = ByteBuffer.wrap(packetLengthBytes).getShort() & 0xFFFF;
        int total = declared + MODBUS_TCP_HEADER_SIZE;
        if (total <= MODBUS_TCP_HEADER_SIZE || total > MODBUS_TCP_MAX_ADU_SIZE) {
            return -1;
        }
        return total;
    }

    private static void sleepBriefly() {
        try {
            Thread.sleep(1);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

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
        final CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        Thread discoveryThread = new Thread(() -> executeDiscovery(future, discoveryRequest, handler));
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
        } catch (Throwable e) {
            logger.error("Error collecting list of possible IP addresses", e);
            future.complete(new DefaultPlcDiscoveryResponse(
                discoveryRequest, PlcResponseCode.INTERNAL_ERROR, Collections.emptyList()));
            return;
        }
        try {
            possibleAddresses.add(InetAddress.getByName("localhost"));
        } catch (UnknownHostException e) {
            throw new PlcRuntimeException(e);
        }

        possibleAddresses = possibleAddresses.stream().filter(
            distinctByKey(InetAddress::getHostAddress)).collect(Collectors.toList());

        Queue<PlcDiscoveryItem> discoveryItems = new ConcurrentLinkedQueue<>();
        possibleAddresses.stream().parallel().forEach(possibleAddress -> {
            try {
                logger.info("Trying address: {}", possibleAddress);
                Socket socket = new Socket(possibleAddress.getHostAddress(), Constants.MODBUSTCPDEFAULTPORT);

                logger.info("Connected: {}", possibleAddress);

                final OutputStream outputStream = socket.getOutputStream();
                final InputStream inputStream = new BufferedInputStream(socket.getInputStream());

                int transactionIdentifier = 1;
                for (short unitIdentifier = 1; unitIdentifier <= 247; unitIdentifier++) {
                    ModbusTcpADU packet = new ModbusTcpADU(transactionIdentifier++, unitIdentifier,
                        new ModbusPDUReadCoilsRequest(1, 1));
                    byte[] deviceIdentificationBytes = null;
                    try {
                        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[packet.getLengthInBytes()]);
                        packet.serialize(writeBuffer);
                        deviceIdentificationBytes = writeBuffer.getBytes();
                    } catch (BufferException e) {
                        logger.error("Error creating the device identification request", e);
                    }
                    if (deviceIdentificationBytes == null) {
                        future.complete(new DefaultPlcDiscoveryResponse(
                            discoveryRequest, PlcResponseCode.INTERNAL_ERROR, Collections.emptyList()));
                        return;
                    }

                    outputStream.write(deviceIdentificationBytes);
                    outputStream.flush();

                    byte[] responseBytes = null;
                    final long endTime = System.currentTimeMillis() + 100;
                    while (responseBytes == null) {
                        // Checked here rather than in one branch of the wait: a host that sent six
                        // bytes and then stopped kept the other branch busy forever, since it always
                        // had enough for a header and never enough for a packet.
                        if (System.currentTimeMillis() > endTime) {
                            break;
                        }
                        if (inputStream.available() >= MODBUS_TCP_HEADER_SIZE) {
                            inputStream.mark(MODBUS_TCP_HEADER_SIZE);
                            inputStream.skip(4);
                            byte[] packetLengthBytes = new byte[2];
                            int bytesRead = inputStream.read(packetLengthBytes);
                            inputStream.reset();
                            if (bytesRead != 2) {
                                continue;
                            }
                            int packetLength = aduLength(packetLengthBytes);
                            if (packetLength < 0) {
                                // Whatever is answering, it is not answering Modbus.
                                break;
                            }
                            if (inputStream.available() >= packetLength) {
                                responseBytes = new byte[packetLength];
                                bytesRead = inputStream.read(responseBytes);
                                if (bytesRead != packetLength) {
                                    responseBytes = null;
                                    break;
                                }
                            } else {
                                sleepBriefly();
                            }
                        } else {
                            sleepBriefly();
                        }
                    }
                    if (responseBytes != null) {
                        ReadBufferByteBased readBuffer = new ReadBufferByteBased(responseBytes);
                        try {
                            ModbusTcpADU response = (ModbusTcpADU) ModbusTcpADU.staticParse(readBuffer, DriverType.MODBUS_TCP, true);
                            boolean found = false;
                            if (response.getPdu().getErrorFlag()) {
                                ModbusPDUError errorPdu = (ModbusPDUError) response.getPdu();
                                if (errorPdu.getExceptionCode() == ModbusErrorCode.ILLEGAL_DATA_ADDRESS) {
                                    found = true;
                                }
                            } else {
                                found = true;
                            }
                            if (found) {
                                PlcDiscoveryItem discoveryItem = new DefaultPlcDiscoveryItem(
                                    "modbus-tcp", "tcp", possibleAddress.getHostAddress(),
                                    Collections.singletonMap("unit-identifier", Integer.toString(unitIdentifier)),
                                    "unknown", Collections.emptyMap());
                                discoveryItems.add(discoveryItem);

                                if (handler != null) {
                                    handler.handle(discoveryItem);
                                }
                                break;
                            }
                        } catch (BufferException | RuntimeException e) {
                            // Whatever answered is not a Modbus device we recognise. That is the
                            // ordinary outcome of scanning an address, not a reason to stop.
                            logger.debug("Ignoring an unreadable response from {}", possibleAddress, e);
                        }
                    }
                }
            } catch (IOException e) {
                // Normal for non-existent devices.
            }
        });

        future.complete(new DefaultPlcDiscoveryResponse(discoveryRequest, PlcResponseCode.OK,
            Arrays.asList(discoveryItems.toArray(new PlcDiscoveryItem[0]))));
    }

}
