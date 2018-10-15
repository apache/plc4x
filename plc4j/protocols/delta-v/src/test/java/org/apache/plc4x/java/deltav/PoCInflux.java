/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.java.deltav;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.commons.codec.binary.Hex;
import org.influxdb.BatchOptions;
import org.influxdb.InfluxDB;
import org.influxdb.InfluxDBFactory;
import org.influxdb.dto.Point;
import org.influxdb.dto.Pong;
import org.pcap4j.core.*;
import org.pcap4j.packet.UdpPacket;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class PoCInflux {

    private static final Logger logger = LoggerFactory.getLogger(PoCInflux.class);

    private static final int SNAPLEN = 65536;
    private static final int READ_TIMEOUT = 10;

    private PcapHandle receiveHandle;
    private InfluxDB connection;

    private List<String> missingNames = new LinkedList<>();
    private Map<String, String> testpointNames = new HashMap<>();
    private Map<String, Map<Short, String>> testpointFieldNames = new HashMap<>();
    private Map<String, Map<Short, Object>> testpointFieldValues = new HashMap<>();

    private PoCInflux(String inputPath) throws Exception {
        if(inputPath == null) {
            PcapNetworkInterface nif = null;
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                if ("en7".equals(dev.getName())) {
                    nif = dev;
                    break;
                }
                else if("eth0".equals(dev.getName())) {
                    nif = dev;
                    break;
                }
            }

            if(nif == null) {
                throw new RuntimeException("Couldn't find network device");
            }

            // Setup receiving of packets and redirecting them to the corresponding listeners.
            // Filter packets to contain only the ip protocol number of the current protocol.
            receiveHandle = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);

            System.out.println("Running in Network-Mode with device: " + nif.getName());
        } else {
            File input = new File(inputPath);
            if(!input.exists() || !input.isFile()) {
                throw new RuntimeException("Couldn't open the pcap file as it doesn't seem to exist.");
            }

            receiveHandle = Pcaps.openOffline(input.getAbsolutePath(), PcapHandle.TimestampPrecision.NANO);

            System.out.println("Running in Simulated-Mode using PCAPNG file: " + inputPath);
        }

        // Set the filter.
        String filterString = "udp port 18507";
        receiveHandle.setFilter(filterString, BpfProgram.BpfCompileMode.OPTIMIZE);

        byte[] timeBytes = ByteBuffer.allocate(8).putLong(System.currentTimeMillis()).array();
        System.out.println("Current Time: " + Hex.encodeHexString(timeBytes));

        PacketListener packetListener = packet -> {
            try {
                UdpPacket udpPacket = (UdpPacket) packet.getPayload().getPayload();
                ByteBuf buf = Unpooled.wrappedBuffer(udpPacket.getPayload().getRawData());
                short header = buf.readShort();
                if(header != (short) 0xFACE) {
                    return;
                }
                short packetLength = buf.readShort();
                short messageType = buf.readShort();
                short messageId = buf.readShort();
                short senderId = buf.readShort();
                buf.skipBytes(3); // Timestamp
                buf.skipBytes(3); // 0x800400 or 0x000400

                // Messages with payload 0 are usually responses.
                if(packetLength == 0) {
                    return;
                }

                // We're only interested in type 2 messages.
                if(messageType == 0x0002) {
                    short payloadType = buf.readShort();

                    System.out.println("Got DeltaV packet : " +
                        Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}));

                    switch(payloadType) {
                        case 0x0201: {
                            // Notes:
                            // - Sent from the OS to the controller
                            // Found packets:
                            // 02 01 00 42 0a 08 00 65 01 5a 00 59 00 00 00 00
                            // 00 00 00 00
                            //
                            // 02 01 00 00 00 00 00 65 01 de 00 18 00 00 00 00
                            // 00 00 00 00
//                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet from " + senderId);
//                            outputPacket(buf);
                            break;
                        }
                        case 0x0202: {
                            // Note:
                            // - Seems to occur during connection establishment phase.
                            // - Response to a 0x0201
                            // - Sent from the Controller to the OS
                            // Found packets:
                            // 02 02 00 00 00 00 00 65 01 de 00 18 00 00 00 00
                            // 00 06 00 00 00 00 00 00
                            break;
                        }
                        case 0x0301: {
                            // Note:
                            // - Sent from the OS to the controllers
                            // - It seems that there are two variants of these packets:
                            //   - Short ones containing almost no information
                            //   - Long ones containing a pattern 14 bytes each starting with 0x04 an an incrementing byte value
                            // - Short messages seem to start with a short numeric value followed by
                            //      0a 08 00 66 00 00 00 00 00 00 00 00 00 04 00 00 00 00 00 00
                            //   finished by an increasing one-byte value which is increased by 2 for every packet
                            //   - The incremented number seems to be increased by 3 whenever a overrun would occur.
                            //     This causes the counter to have just even numbers for one run and then switch to odd ones
                            //     after the next over-run and then back to even ones after the next.
                            // - Big messages seem to contain some pattern of 14 bytes
                            //   - Each pattern seems to have the following pattern:
                            //     04 {incremented byte} 17 {ib + 0x29} 00 {ib + 0x1D} 00 {ib + 0x3A} 00 00 03 e8 00 00
                            //
                            // Found packets:
                            // (small)
                            // 03 01 00 2e 0a 08 00 66 00 00 00 00 00 00 00 00
                            // 00 04 00 00 00 00 00 00 c5
                            //
                            // (medium)
                            // 03 01 00 00 00 00 00 07 02 95 00 bf 00 04 00 00
                            // 00 08 00 01 00 00 00 00 01 04 8f 17 b8 00 2f 01
                            // 36 00 00 3a 98 00 01 00 04 00 44 00 4e 00 4c 00
                            // 44 00 00 00 02 00 06 00 4f 00 49 00 4e 00 54 00
                            // 45 00 47 00 00 ff ff 88 ad
                            // Contains Text: DNLD and OINTEG
                            //
                            // (medium)
                            // 03 01 00 00 00 00 00 07 02 ec 01 16 00 04 00 00
                            // 00 08 00 01 00 00 00 00 01 04 8f 17 b8 00 2f 01
                            // 36 00 00 3a 98 00 03 00 06 00 53 00 49 00 4e 00
                            // 54 00 45 00 47 00 00 ff ff
                            // Contains Text: SINTEG
                            //
                            // (big)
                            // 03 01 00 00 00 00 00 53 00 00 00 00 00 00 00 00
                            // 00 04 00 00 00 00 00 00 23 04 1f 17 48 00 3c 00
                            // 59 00 00 03 e8 00 00 04 20 17 49 00 3d 00 5a 00
                            // 00 03 e8 00 00 04 21 17 4a 00 3e 00 5b 00 00 03
                            // e8 00 00 04 25 17 4e 00 42 00 5f 00 00 03 e8 00
                            // 00 04 26 17 4f 00 43 00 60 00 00 03 e8 00 00 04
                            // 27 17 50 00 44 00 61 00 00 03 e8 00 00 04 28 17
                            // 51 00 45 00 62 00 00 03 e8 00 00 04 23 17 4c 00
                            // 40 00 5d 00 00 03 e8 00 00 04 24 17 4d 00 41 00
                            // 5e 00 00 03 e8 00 00 04 2a 17 53 00 47 00 64 00
                            // 00 03 e8 00 00 04 2e 17 57 00 4b 00 68 00 00 03
                            // e8 00 00 04 31 17 5a 00 4e 00 6b 00 00 03 e8 00
                            // 00 04 36 17 5f 00 53 00 70 00 00 03 e8 00 00 04
                            // 37 17 60 00 54 00 71 00 00 03 e8 00 00 04 32 17
                            // 5b 00 4f 00 6c 00 00 03 e8 00 00 00 00 00 00
//                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet from " + senderId);
//                            outputPacket(buf);
                            break;
                        }
                        case 0x0302: {
                            // Notes:
                            // - Seems to be one possible response to a 0x0301 message
                            // - Sent from the Controller to the OS
                        }
                        case 0x0304: {
                            // Note:
                            // - The Operator Systems seem to be sending these messages.
                            // Found packets:
                            //   03 04 00 00 00 00 00 6b 00 00 00 00 00 00 00 00
                            //   00 04 00 00 00 00 00 00 49 00 86 00 01 1a 02 00
                            //   00 20 03 00 04 00 0b b6 9c 00 6b 00 00 2a cb e2
                            //   e6 f8 e0 00 00 2a cb e2 e6 f8 d4 00 00 00 00 00
                            //
                            //   11 00 50 00 54 00 30 00 39 00 2d 00 30 00 31 00
                            //   2f 00 4d 00 41 00 49 00 4e 00 54 00 5f 00 41 00
                            //   4c 00 4d 00 00
                            //
                            //   47 00
                            //
                            //   18 00 23 00 44 00 23 00 31
                            //   00 2c 00 34 00 2c 00 30 00 30 00 30 00 30 00 30
                            //   00 30 00 30 00 31 00 2c 00 30 00 30 00 30 00 30
                            //   00 30 00 30 00 30 00 30 00 00
                            //
                            //   08 00 00 00 00 ff
                            //   ff 00 86 00 01 1a 02 00 00 20 03 00 01 00 0b b6
                            //   9c 00 6b 00 00 2a cb e3 21 90 b3 00 00 2a cb e2
                            //   e6 f8 d4 00 00 00 00 00 11 00 50 00 54 00 30 00
                            //   39 00 2d 00 30 00 31 00 2f 00 4d 00 41 00 49 00
                            //   4e 00 54 00 5f 00 41 00 4c 00 4d 00 00 47 00 18
                            //   00 23 00 44 00 23 00 30 00 2c 00 34 00 2c 00 30
                            //   00 30 00 30 00 30 00 30 00 30 00 30 00 31 00 2c
                            //   00 30 00 30 00 30 00 30 00 30 00 30 00 30 00 30
                            //   00 00 08 00 00 00 00 ff ff 00 50 00 01 02 02 00
                            //   00 00 3f 00 36 00 18 b6 9c 00 6b 00 00 2a cb e2
                            //   e3 b6 a4 00 00 00 00 00 00 00 00 00 00 00 00 00
                            //   08 00 41 00 43 00 4e 00 20 00 43 00 4f 00 4d 00
                            //   4d 00 00 47 00 07 00 45 00 53 00 31 00 30 00 36
                            //   00 30 00 31 00 00 47 ff ff ff ff 00 88 00 01 1a
                            //   01 00 00 20 04 00 02 00 07 b6 9c 00 6b 00 00 2a
                            //   cb cf 72 da a9 00 00 2a cb ce 66 db 1c 00 00 00
                            //   00 00 12 00 50 00 54 00 30 00 39 00 2d 00 30 00
                            //   31 00 2f 00 41 00 44 00 56 00 49 00 53 00 45 00
                            //   5f 00 41 00 4c 00 4d 00 00 47 00 18 00 23 00 44
                            //   00 23 00 30 00 2c 00 32 00 2c 00 30 00 30 00 30
                            //   00 39 00 30 00 30 00 30 00 45 00 2c 00 30 00 30
                            //   00 30 00 30 00 30 00 31 00 30 00 31 00 00 08 00
                            //   00 00 00 ff ff 00 7c 00 01 00 01 26 00 00 2d 00
                            //   05 00 03 b6 9c 00 6b 00 00 2a cb ce 66 42 66 00
                            //   00 2a cb ce 66 42 66 00 00 00 00 00 24 00 3a 00
                            //   55 00 4e 00 49 00 54 00 5f 00 50 00 54 00 30 00
                            //   39 00 3a 00 50 00 54 00 30 00 39 00 41 00 4c 00
                            //   41 00 43 00 31 00 4d 00 43 00 54 00 33 00 30 00
                            //   2f 00 4d 00 41 00 58 00 5f 00 53 00 43 00 48 00
                            //   41 00 4c 00 54 00 00 08 c2 20 00 00 08 43 43 00
                            //   00 ff ff 00 7c 00 01 00 01 26 00 00 35 00 05 00
                            //   03 b6 9c 00 6b 00 00 2a cb ce 66 42 6b 00 00 2a
                            //   cb ce 66 42 6b 00 00 00 00 00 24 00 3a 00 55 00
                            //   4e 00 49 00 54 00 5f 00 50 00 54 00 30 00 39 00
                            //   3a 00 50 00 54 00 30 00 39 00 41 00 4c 00 41 00
                            //   43 00 31 00 4d 00 43 00 54 00 33 00 30 00 2f 00
                            //   4d 00 49 00 4e 00 5f 00 53 00 43 00 48 00 41 00
                            //   4c 00 54 00 00 08 c2 20 00 00 08 c2 18 00 00 ff
                            //   ff ff ff
//                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet from " + senderId);
//                            outputPacket(buf);
                            break;
                        }
                        case 0x0401: {
                            // Note:
                            // - Sent by the OS during connection phase.
                            // Found packets:
                            // Connection OS sent to Controller:
                            // 04 01 00 00 00 00 00 64 03 57 01 82 00 00 00 00
                            // 00 19 00 00 00 05 00 24 00 00 00 00 00 00 00 fa
                            // ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff
                            // ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff c0
                            // 01 2e
                            // Later on:
                            // 04 01 00 2a 0a 08 00 64 02 26 01 38 00 00 00 00
                            // 00 19 00 00 00 05 00 24 00 00 00 00 00 00 00 fa
                            // ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff
                            // ff ff ff ff ff ff ff ff ff ff ff ff ff ff ff c0
                            // 01 2e
//                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet from " + senderId);
                            //outputPacket(buf);
                            break;
                        }
                        case 0x0402: {
                            // Note:
                            // - Seems to be a response to a 0x0401 packet and it seems to replicate 5 bytes sent in the 0x0401
                            // Found packets:
                            // 04 02 00 00 00 00 00 64 02 7d 01 8f 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 06 00 00 00 00 00 00
//                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet from " + senderId);
                            //outputPacket(buf);
                            break;
                        }
                        case 0x1B02:
                            // 1b 02 00 00 00 00 00 07 03 ac 01 d6 00 04 00 00
                            // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 01 00 00 00 00 00 32 01
                            // 39 04 90 17 b9. 02 00 13 47 00 4b 00 22 00 49 00
                            // 6e 00 69 00 74 00 69 00 61 00 6c 00 69 00 73 00
                            // 69 00 65 00 72 00 75 00 6e 00 67 00 20 00 2e 00
                            // 2e 00 2e 00 2e 00 2e 00 20 00 62 00 69 00 74 00
                            // 74 00 65 00 20 00 77 00 61 00 72 00 74 00 65 00
                            // 6e 00 00. 03 00 14 00 00 00 1f. 02 00 15 03 00 00.
                            // 02 00 16 47 00 1d 00 0b 00 57 00 41 00 52 00 54
                            //---------^
                            //Last block started: 1 bytes before error and had a size of: 9
                            //Unexpected code: 47
                            //
                            //
                            // 00 45 00 4e 00 20 00 2e 00 2e 00 2e 00 20 00 00
                            // 01
                        case 0x0403: {
                            // Skip the rest of the header.
                            if(payloadType == 0x1B02) {
                                buf.skipBytes(118); // 0x76

                                // Opening some detail dialog.
                                // 1b 02 00 00 00 00 00 07 03 ac 01 d6 00 04 00 00
                                // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00
                                // - 01 00 00 00 00 00 32 01 39 04 90 17 b9
                                // 02 00 13 47 00 4b 00 22 00 49 00 6e 00 69 00 74 00 69 00 61 00 6c 00 69 00 73 00 69 00 65 00 72 00 75 00 6e 00 67 00 20 00 2e 00 2e 00 2e 00 2e 00 2e 00 20 00 62 00 69 00 74 00 74 00 65 00 20 00 77 00 61 00 72 00 74 00 65 00 6e 00 00
                                // 03 00 14 00 00 00 1f
                                // 02 00 15 03 00 00
                                // 02 00 16 47 00 1d 00 0b 00 57 00 41 00 52 00 54 00 45 00 4e 00 20 00 2e 00 2e 00 2e 00 20 00 00
                                // - 01
                                //
                                // 1b 02 00 00 00 00 00 07 03 ad 01 d7 00 04 00 00
                                // 00 06 00 01 00 06 00 94 00 06 01 37 00 06 02 00
                                // 00 06 02 f5 00 06 03 98 00 06 04 4d 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00
                                // - 07 00 00 00 00 00 3f 01 46 04 9d 17 c6
                                // 02 00 00 76 00 66 05 01 01 00 00 00 06 00 46 00 45 00 48 00 4c 00 45 00 52 00 00 00 3b 00 3b 01 0d 00 00 2a f2 70 c7 be 79 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 01 63 00 7c 0d a6 a2
                                // 02 00 02 01 00
                                // 02 00 03 05 01
                                // 02 00 04 01 01
                                // 02 00 05 01 00
                                // - 01 00 00 00 00 00 40 01 47 04 9e 17 c7
                                // 02 00 00 08 42 48 d5 53
                                // 02 00 01 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 02 01 00
                                // 02 00 03 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 04 47 00 0b 00 02 00 54 00 49 00 00
                                // 02 00 05 05 01
                                // 02 00 06 01 00
                                // 02 00 07 01 00
                                // - 01 00 00 00 00 00 41 01 48 04 9f 17 c8
                                // 02 00 00 01 00
                                // 02 00 01 01 01
                                // 02 00 02 08 00 00 00 00
                                // 02 00 03 24 00 0e 41 a0 00 00 00 00 00 00 05 45 01
                                // 02 00 04 76 00 7e 05 01 01 00 00 00 07 00 4d 00 49 00 4e 00 5f 00 41 00 4c 00 4d 00 00 00 32 00 32 01 0d 00 00 2a f2 70 c7 21 f3 00 0a 00 4d 00 49 00 4e 00 5f 00 53 00 43 00 48 00 41 00 4c 00 54 00 00 00 36 00 36 00 03 00 00 2a f2 70 c7 21 f6 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 05 01 00
                                // 02 00 06 47 00 0b 00 02 00 46 00 49 00 00
                                // 02 00 07 05 01
                                // - 01 00 00 00 00 00 42 01 49 04 a0 17 c9
                                // 02 00 00 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 01 08 42 c8 00 00
                                // 02 00 02 63 00 00 19 c9 9c
                                // 02 00 03 48 80 01
                                // 02 00 04 01 00
                                // 02 00 05 21 80 42 49 1c 72
                                // 02 00 06 01 00
                                // 02 00 07 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 08 47 00 0d 00 03 00 54 00 49 00 43 00 00
                                // 02 00 09 08 00 00 00 00
                                // 02 00 0a 01 00
                                // 02 00 0b 08 42 49 1c 72
                                // 02 00 0c 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 0d 21 c3 42 17 aa ab
                                // 02 00 0e 63 00 00 19 9c 7b
                                // 02 00 0f 01 00
                                // 02 00 10 05 01
                                // - 01 00 00 00 00 00 43 01 4a 04 a1 17 ca
                                // 02 00 00 05 01
                                // 02 00 01 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 02 01 00
                                // 02 00 03 01 00
                                // 02 00 04 08 42 49 1c 72
                                // 02 00 05 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 06 01 00
                                // 02 00 07 47 00 0b 00 02 00 54 00 49 00 00
                                // - 01 00 00 00 00 00 44 01 4b 04 a2 17 cb
                                // 02 00 00 76 00 7e 05 01 01 00 00 00 0a 00 4d 00 41 00 58 00 5f 00 53 00 43 00 48 00 41 00 4c 00 54 00 00 00 2e 00 2e 00 03 00 00 2a f2 71 9b 0f 2c 00 07 00 4d 00 41 00 58 00 5f 00 41 00 4c 00 4d 00 00 00 29 00 29 00 03 00 00 2a f2 71 9b 0f 2b ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 01 01 00
                                // 02 00 02 21 80 3c ce e0 00
                                // 02 00 03 01 00
                                // 02 00 04 05 01
                                // 02 00 05 01 01
                                // 02 00 06 01 01
                                // 02 00 07 01 00
                                // - 01 00 00 00 00 00 45 01 4c 04 a3 17 cc
                                // 02 00 00 01 00
                                // 02 00 01 01 00
                                // 02 00 02 05 01
                                // - 01
                                //
                                // 1b 02 00 00 00 00 00 07 03 ae 01 d8 00 04 00 00
                                // 00 06 00 01 00 06 00 be 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 <!-- No testpoint id?!?
                                // 02 00 00 00 00 00 4b 01 52 04 a9 17 d2 02 00 00 01 00
                                // 02 00 01 76 00 86 05 01 01 00 00 00 0a 00 56 00 45 00 52 00 52 00 49 00 45 00 47 00 45 00 4c 00 54 00 00 00 54 00 54 00 04 00 00 2a f2 71 97 40 42 00 0b 00 4d 00 49 00 4e 00 5f 00 41 00 5f 00 41 00 4c 00 4d 00 5f 00 59 00 00 00 35 00 35 00 03 00 00 2a f2 71 96 46 47 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 02 01 00
                                // 02 00 03 01 00
                                // 02 00 04 21 80 00 00 00 00
                                // 02 00 05 05 01
                                // 02 00 06 01 00
                                // 02 00 07 01 01
                                // - 01 00 00 00 00 00 4c 01 53 04 aa 17 d3
                                // 02 00 00 07 00 00 00 00
                                // 02 00 01 05 01
                                // 02 00 02 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 03 01 01
                                // 02 00 04 01 01
                                // 02 00 05 01 00
                                // 02 00 06 01 00
                                // 01 b2
                                // 29
                                //
                                // 1b 02 00 00 00 00 00 07 03 af 01 d9 00 04 00 00
                                // 00 06 00 01 00 06 00 af 00 06 01 42 00 06 02 3d
                                // 00 06 02 d0 00 06 03 63 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00
                                // - 06 00 00 00 00 00 45 01 4c 04 a3 17 cc
                                // 02 00 03 76 00 86 05 01 01 00 00 00 0a 00 56 00 45 00 52 00 52 00 49 00 45 00 47 00 45 00 4c 00 54 00 00 00 54 00 54 00 04 00 00 2a f2 70 d0 7b d3 00 0b 00 4d 00 49 00 4e 00 5f 00 41 00 5f 00 41 00 4c 00 4d 00 5f 00 59 00 00 00 35 00 35 00 03 00 00 2a f1 c7 5e 3e a7 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 04 01 00
                                // 02 00 05 01 01
                                // 02 00 06 21 80 00 00 00 00
                                // 02 00 07 01 00
                                // - 01 00 00 00 00 00 46 01 4d 04 a4 17 cd
                                // 02 00 00 63 01 7c 0d a6 d7
                                // 02 00 01 01 00
                                // 02 00 02 05 01
                                // 02 00 03 76 00 66 05 01 01 00 00 00 06 00 46 00 45 00 48 00 4c 00 45 00 52 00 00 00 52 00 52 01 0d 00 00 2a f1 c7 f0 48 88 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 04 01 01
                                // 02 00 05 01 00
                                // - 01 00 00 00 00 00 47 01 4e 04 a5 17 ce
                                // 02 00 00 01 00
                                // 02 00 01 47 00 0d 00 03 00 54 00 49 00 43 00 00
                                // 02 00 02 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 03 01 00
                                // 02 00 04 63 01 00 19 c9 9c
                                // 02 00 05 08 42 4c 00 00
                                // 02 00 06 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 07 21 c0 42 95 fe c7
                                // 02 00 08 08 42 48 00 00
                                // 02 00 09 01 00
                                // 02 00 0a 05 01
                                // 02 00 0b 48 80 01
                                // 02 00 0c 63 01 00 19 9c 7b
                                // 02 00 0d 01 00
                                // 02 00 0e 21 80 42 49 1c 72
                                // 02 00 0f 01 00
                                // 02 00 10 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 11 21 80 42 95 fe c7
                                // - 01 00 00 00 00 00 48 01 4f 04 a6 17 cf
                                // 02 00 00 01 00
                                // 02 00 01 05 01
                                // 02 00 02 76 00 66 05 01 01 00 00 00 06 00 46 00 45 00 48 00 4c 00 45 00 52 00 00 00 52 00 52 01 0d 00 00 2a f2 70 c7 30 06 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 03 01 00
                                // 02 00 04 01 00
                                // 02 00 05 63 00 7c 0d a6 d7
                                // - 01 00 00 00 00 00 49 01 50 04 a7 17 d0
                                // 02 00 00 01 00
                                // 02 00 01 76 00 66 05 01 01 00 00 00 06 00 46 00 45 00 48 00 4c 00 45 00 52 00 00 00 52 00 52 01 0d 00 00 2a f1 c8 bc f6 68 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 02 01 00
                                // 02 00 03 01 00
                                // 02 00 04 05 01
                                // 02 00 05 63 00 7c 0d a6 d7
                                // - 01 00 00 00 00 00 4a 01 51 04 a8 17 d1
                                // 02 00 00 01 00
                                // 02 00 01 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // 02 00 02 47 00 0b 00 02 00 54 00 49 00 00
                                // 02 00 03 05 01
                                // 02 00 04 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                // 02 00 05 07 00 00 00 00
                                // 02 00 06 01 00
                                // 02 00 07 21 80 42 48 00 00
                                // - 01
                                //
                                // 1b 02 00 00 00 00 00 07 01 ed 00 16 00 04 00 00
                                // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00
                                // - 01 00 00 00 00 00 40 01 47 04 9e 17 c7
                                // 03 00 08 00 00 00 1f
                                // - 01
                                //
                                // 1b 02 00 00 00 00 00 07 01 ee 00 17 00 04 00 00
                                // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00
                                // - 01 00 00 00 00 00 40 01 47 04 9e 17 c7
                                // 02 00 09 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // - 01
                                //
                                // 1b 02 00 00 00 00 00 07 01 ef 00 18 00 04 00 00
                                // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                                // 00 00 00 00 00 00 00 00
                                // - 01 00 00 00 00 00 40 01 47 04 9e 17 c7
                                // 02 00 0a 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                // - 01
                            } else {
                                buf.skipBytes(0x1A);
//                                outputPacket(buf);
                            }
                            buf.skipBytes(5);
                            byte[] testpointId = new byte[8];
                            buf.readBytes(testpointId, 4, 4);
                            buf.readBytes(testpointId, 0, 4);
                            String testpoint = Hex.encodeHexString(testpointId);

                            int endOfLastBlock = buf.readerIndex();
                            int lastBlockSize = 0;
                            for (byte code = buf.readByte(); buf.readableBytes() > 2; code = buf.readByte()) {
                                // First check the code of the next block ...
                                switch (code) {
                                    case (byte) 0x01: {
                                        // First make sure the current testpoints state is flushed to the db.
                                        flushTestpointValues(testpoint);
                                        // 0x01 blocks for 0x1B02 messages contain 4 bytes more.
                                        if(payloadType == 0x1B02) {
                                            buf.skipBytes(4);
                                        }
                                        if(buf.readableBytes() < 8) {
                                            return;
                                        }
                                        // Now switch to the next testpoint.
                                        buf.readBytes(testpointId, 4, 4);
                                        buf.readBytes(testpointId, 0, 4);
                                        testpoint = Hex.encodeHexString(testpointId);
                                        break;
                                    }
                                    case (byte) 0x02: {
                                        short fieldId = buf.readShort();
                                        byte type = buf.readByte();

                                        // Now inspect the block content ...
                                        switch (type) {
                                            case (byte) 0x01: {
                                                // Possibly boolean value?
                                                byte booleanByteValue = buf.readByte();
                                                boolean booleanValue = false;
                                                switch (booleanByteValue) {
                                                    case (byte) 0x00:
                                                        booleanValue = false;
                                                        break;
                                                    case (byte) 0x01:
                                                        booleanValue = true;
                                                        break;
                                                    default:
                                                        System.out.println("Unknown second byte for boolean value 0x" + Hex.encodeHexString(new byte[]{booleanByteValue}));
                                                }
                                                updateValue(testpoint, fieldId, booleanValue);
                                                break;
                                            }
                                            case (byte) 0x03: {
                                                // Note:
                                                // - Name suggests a time-related field "A/ALA_T9002_Y/RESTZEIT"
                                                //
                                                // Found Blocks:
                                                // 02 00 15 03 00 00
                                                buf.skipBytes(2);
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x04: {
                                                // Notes:
                                                // - Refers to variable "C/MERROR_MOD"
                                                // - seems to be a 4 byte integer value.
                                                // Found Blocks:
                                                // 02 00 19 04 00 00 01 08
                                                // 02 00 19 04 00 00 00 00
                                                buf.skipBytes(4);
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x05: {
                                                // NOTE:
                                                // - Name of the variable "A/GFE_ALARM", "BA"
                                                //
                                                // Found Block:
                                                // 02 00 0c 05: 00
                                                buf.skipBytes(1);
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown (BA))", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x06: {
                                                // Possibly Parse 16 bit int?
                                                short shortValue = buf.readShort();
                                                updateValue(testpoint, fieldId, shortValue);
                                                break;
                                            }
                                            case (byte) 0x07: {
                                                // Possibly Parse 32 bit int?
                                                int intValue = buf.readInt();
                                                updateValue(testpoint, fieldId, intValue);
                                                break;
                                            }
                                            case (byte) 0x08: {
                                                // Parse float
                                                float floatValue = buf.readFloat();
                                                //floatValue = Math.round(floatValue * 100.0f) / 100.0f;
                                                updateValue(testpoint, fieldId, floatValue);
                                                break;
                                            }
                                            case (byte) 0x21: {
                                                // From having a look at the byte values these could be 32bit floating point values with some sort of parameters
                                                byte param = buf.readByte();
                                                //decodeParam(param);
                                                float floatValue = buf.readFloat();
                                                //floatValue = Math.round(floatValue * 100.0f) / 100.0f;
                                                updateValue(testpoint, fieldId, floatValue);
                                                break;
                                            }
                                            case (byte) 0x22: {
                                                // Parse boolean (From what I learnt, this could be a flagged boolean, where the first byte is some sort of param)
                                                byte param = buf.readByte();
                                                //decodeParam(param);
                                                byte booleanByteValue = buf.readByte();
                                                boolean booleanValue = false;
                                                switch (booleanByteValue) {
                                                    case (byte) 0x00:
                                                        booleanValue = false;
                                                        break;
                                                    case (byte) 0x01:
                                                        booleanValue = true;
                                                        break;
                                                    default:
                                                        System.out.println("Unknown second byte for boolean value 0x" + Hex.encodeHexString(new byte[]{booleanByteValue}));
                                                }
                                                updateValue(testpoint, fieldId, booleanValue);
                                                break;
                                            }
                                            case (byte) 0x24: {
                                                // NOTE:
                                                // - Name of the variable  "SKAL1", "SKAL2", "SKAL3".
                                                // - Might be scaling variables providing information on how to scale a value
                                                // - Potentially needed to display a bar graph for a value.
                                                // - Quite a lot of these are used.
                                                // - I would assume it's 3 32 bit unsigned integers and a trailing flag-byte
                                                // - The flag byte seems to contain the values: 0, 1, 2 and 3 (So it's either two boolean flags or 4 types)
                                                //
                                                // Found Packets:
                                                // 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                                                //
                                                // 00 0e 42 c8 00 00 00 00 00 00 04 40 00
                                                //
                                                // 00 0e 41 20 00 00 00 00 00 00 05 45 01
                                                //
                                                // 00 0e 3f 80 00 00 bf 80 00 00 82 df 03
                                                long val1 = buf.readUnsignedInt();
                                                long val2 = buf.readUnsignedInt();
                                                long val3 = buf.readUnsignedInt();
                                                byte flag = buf.readByte();
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Trend Scaling) - " + val1 + ", " + val2 + ", " + val3 + ", " + flag, buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x25: {
                                                buf.skipBytes(6);
                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x47: {
                                                // NOTE:
                                                // - Name of the variable "A/ALA_T9002_Y/SCHRITTTEXT_ALT", "A/ALA_T9002_Y/SCHRITTTEXT"
                                                // - Seems the most important variable name however is "GF_NAME" which seems to give the testpoint it's name
                                                // - Seems to be a text-value
                                                // - The longer UTF-16 Text seems to start at byte 4 and ends at 0x0000
                                                // - The shorter Text seems to come after some bytes after the longer text and again end with 0x0000
                                                // - Seems to be sent as soon as a user confirms an alarm.
                                                // - Seems the length is variable
                                                //
                                                // Found Blocks:
                                                // 00 4b 00 22 00 49 00 6e 00 69 00 74 00 69 00 61
                                                // 00 6c 00 69 00 73 00 69 00 65 00 72 00 75 00 6e
                                                // 00 67 00 20 00 2e 00 2e 00 2e 00 2e 00 2e 00 20
                                                // 00 62 00 69 00 74 00 74 00 65 00 20 00 77 00 61
                                                // 00 72 00 74 00 65 00 6e 00 00
                                                //
                                                // 02 00 67 47 00 1d 00 0b 00 57 00 41 00 52 00 54
                                                // 00 45 00 4e 00 20 00 2e 00 2e 00 2e 00 20 00 00
                                                // Decoded:
                                                // K"Initialisierung ..... bitte warten Ȁ杇WARTEN ...
                                                //
                                                // Name: BESCH3
                                                // 02 00 00 47 00 0b 00 02 00 44 00 54 00 00
                                                //
                                                // Name: DYN_BESCH1
                                                // 02 00 02 47 00 0b 00 02 00 56 00 4c 00 00
                                                //
                                                // Name: DYN_BESCH1
                                                // 02 00 04 47 00 07 00 00 00 00
                                                //
                                                // It seems the parameter provides the complete length of the block in bytes
                                                // (including the final 0x0000 and the type and size indicator.
                                                short numBytes = (short) (buf.readShort() - 7);
                                                // This seems to be some sort of flag ...
                                                short flag = buf.readShort();
                                                byte[] bytes = new byte[numBytes];
                                                buf.readBytes(bytes);
                                                // Should be the trailing 0x0000
                                                buf.skipBytes(2);
                                                String text = new String(bytes, StandardCharsets.UTF_16);
                                                if(flag == (short) 9) {
                                                    if(!testpointNames.containsKey(testpoint)) {
                                                        testpointNames.put(testpoint, text);
                                                    }
                                                    if(!testpointNames.get(testpoint).equals(text)) {
                                                        System.out.println("Guess I was wrong");
                                                    }
                                                }
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Typed Text)", text, null, buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x48: {
                                                // NOTE:
                                                // - Name of the variable "C/GW_AW_MIN", "C/GW_AW_MAX", "C/GW_SW_MIN", "C/GW_SW_MAX", "C/GW_GW_MIN", "C/GW_GW_MAX"
                                                // - Seems to be sent as soon as an alarm is fired, changed or removed from the controller.
                                                // - There seem to be only two types of values: 0x8000 and 0x8001
                                                // - Might be something similar to "BOOL(P)"
                                                //
                                                // Found packets:
                                                // 02 00 10 48: 80 00
                                                byte[] tmp = new byte[2];
                                                buf.readBytes(tmp);
                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x49: {
                                                // NOTE:
                                                // - Judging from the 0x80 first byte I would assume this is again one of these parametrized values
                                                // - Would suggest this is a 32 bit integer value.
                                                // Found blocks:
                                                // 80 00 00 06 0d
                                                byte param = buf.readByte();
                                                //decodeParam(param);
                                                int intValue = buf.readInt();
                                                updateValue(testpoint, fieldId, intValue);
                                                break;
                                            }
                                            case (byte) 0x5B: {
                                                // No idea what this type is.
                                                buf.readShort();
                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x63: {
                                                // NOTE:
                                                // - Name of the variable "IST_FW_KURZ"
                                                // - Looks like a 4-byte value followed by a (P)aram byte
                                                // Found blocks:
                                                // 02 00 06 63: 64 00 19 b9 88
                                                byte[] tmp = new byte[5];
                                                buf.readBytes(tmp);
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x75: {
                                                // NOTE:
                                                // - Name of the variable "MIN_ALM", "MAX_ALM", "MIN_SCHALT", "MAX_SCHALT"
                                                // - Exactly 3 blocks of this type with extremely similar content is being sent every 60 seconds for the ids: 17, 16 and 34
                                                //                            001600280d0100000000280015f360000000000100
                                                //                            001600280d0100000000280015f360000000000100
                                                // Found blocks:
                                                // BIN_ALM
                                                // 02 00 3a 75 00 16 00 18 0d 01 00 00 18 00 18 00 15 f3 60 00 15 e4 10 00 01
                                                //
                                                // STOERUNG
                                                // 02 00 3b 75 00 16 00 39 0f 01 00 00 39 00 39 00 15 f3 60 00 00 00 00 01 00

                                                int size = "001600280d0100000000280015f360000000000100".length() / 2; //21
                                                byte[] tmp = new byte[size];
                                                buf.readBytes(tmp);
//                                                 outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0x76: {
                                                // NOTE:
                                                // - After knowing the name of the variables, it turns out that these are all "Alarm" variables.
                                                // - Possibly a data-structure is returned which contains all information about a modules alarms
                                                //   The last 4 bytes (maybe more) seem to be an always increasing value
                                                //   (Maybe some sort of timestamp)
                                                // Found blocks:
                                                // 02 00 04 76 00 7a 05 01 01 00 00 00 08 00 53 00
                                                // 54 00 4f 00 45 00 52 00 55 00 4e 00 47 00 00 00
                                                // 39 00 39 01 0f 00 00 2a f4 8e c7 c6 36 00 07 00
                                                // 42 00 49 00 4e 00 5f 00 41 00 4c 00 4d 00 00 00
                                                // 18 00 18 00 03 00 00 2a f4 8e c3 e5 2e ff ff 00
                                                // 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00
                                                // 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00
                                                // 00 00 00 00 ff 00 00 00 00 00 00 00 00
                                                //
                                                // 02 00 04 76 00 58 05 01 01 00 00 ff ff 00 00 00
                                                // 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00
                                                // 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00
                                                // 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00
                                                // 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00
                                                // 00 00 ff 00 00 00 00 00 00 00 00

                                                short blockLength = (short) (buf.readShort() - 3);
                                                buf.skipBytes(blockLength);
                                                /*int start = buf.readerIndex();
                                                buf.skipBytes(5);
                                                short textLength = buf.readShort();
                                                if(textLength != (short) 0xFFFF) {
                                                    byte[] tmp = new byte[textLength * 2];
                                                    buf.readBytes(tmp);
                                                    String text1 = new String(tmp, StandardCharsets.UTF_16);
                                                } else {
                                                    outputPacket(buf);
                                                }
                                                buf.skipBytes(2);
                                                buf.skipBytes(14);
                                                textLength = buf.readShort();
                                                if(textLength != (short) 0xFFFF) {
                                                    byte[] tmp = new byte[textLength * 2];
                                                    buf.readBytes(tmp);
                                                    buf.skipBytes(2);
                                                    String text2 = new String(tmp, StandardCharsets.UTF_16);
                                                } else {
                                                    outputPacket(buf);
                                                }

                                                int curPos = buf.readerIndex();
                                                buf.skipBytes(blockLength - (curPos - start));
                                                */
//                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown (Alarms))", buf, endOfLastBlock);
                                                break;
                                            }
                                            case (byte) 0xF6: {
                                                // TODO: Potentially obsolete ...
                                                // Only seen in 0x0102 blocks
                                                buf.skipBytes(4);
                                                outputDetectedBlock(getTestpointName(testpoint, fieldId) + " (Unknown)", buf, endOfLastBlock);
                                                break;
                                            }
                                            default: {
                                                dumpAndExit(buf, endOfLastBlock, lastBlockSize, "Unexpected 0x02 type code: " + Hex.encodeHexString(new byte[]{type}));
                                                /*if(code == (byte) 0x01) {
                                                    buf.skipBytes(4);
                                                } else {
                                                    dumpAndExit(buf, endOfLastBlock, lastBlockSize, "Unknown variable type 0x" + Hex.encodeHexString(new byte[]{type}));
                                                }
                                                outputDetectedBlock("Unknown", buf, endOfLastBlock);*/
                                            }

                                        }
                                        break;
                                    }
                                    case (byte) 0x03: {
                                        short fieldId = buf.readShort();
                                        byte type = buf.readByte();

                                        // Note:
                                        // - Might be error responses ...
                                        // Found blocks:
                                        // 03 00 23 00 00 00 4a             (size 6)
                                        // 03 00 00 00 00 00 1f
                                        switch (type) {
                                            case (byte) 0x00: {
                                                buf.skipBytes(3);
                                                break;
                                            }
                                            default: {
                                                dumpAndExit(buf, endOfLastBlock, lastBlockSize, "Unexpected 0x03 type code: " + Hex.encodeHexString(new byte[]{type}));
                                            }
                                        }
                                        break;
                                    }
                                    default: {
                                        dumpAndExit(buf, endOfLastBlock, lastBlockSize, "Unexpected code: " + Hex.encodeHexString(new byte[]{code}));
                                    }
                                }
                                lastBlockSize = buf.readerIndex() - endOfLastBlock;
                                endOfLastBlock = buf.readerIndex();
                            }
                            break;
                        }
                        case 0x0404: {
                            // Note:
                            // - Seems to be used during connection phase.
                            // Found Packets:
                            // 04 04 00 00 00 00 01 23 00 00 00 00 00 00 00 00
                            // 00 0a 00 00 00 00 00 00 00 00 00 00 00 07 00 45
                            // 00 53 00 31 00 30 00 36 00 30 00 31 00 00 ee a1
                            // ES10601 (The end seems to be a two byte per char encoded string with an ending 0x0000 value
                            // as well as a length information before the text (7 = number of chars in the string)
                        }
                        case 0x0501: {
                            // Seems to contain version information of the Operator System.
                            break;
                        }
                        case 0x0502: {
                            // Seems to contain version information of the Controller.
                            break;
                        }
                        case 0x0506: {
                            // Notes:
                            // - Sent from the Controller to the OS
                            // Found packets:
                            // 05 06 00 00 00 00 01 03 00 00 00 00 00 00 00 00
                            // 00 06 00 00 00 04 00 04 00 04 00 1f 00 00 00 00
                            // 00 00 01 90 00 00 0b 00 45 00 52 00 53 00 41 00
                            // 54 00 5a 00 53 00 54 00 52 00 4f 00 4d 00 00 01
                            // Contained a text at the end: ERSATZSTROM
                            //
                            // 05 06 00 00 00 01 01 03 00 00 00 00 00 00 00 00
                            // 00 06 00 00 00 04 00 04 00 04 00 29 00 00 00 00
                            // 00 00 01 90 00 00 10 00 56 00 45 00 52 00 52 00
                            // 5f 00 51 00 55 00 45 00 4c 00 4c 00 41 00 42 00
                            // 53 00 41 00 55 00 47 00 00 01
                            // Contained a text at the end: VERR_QUELLABSAUG
                            //
                            // 05 06 00 00 00 00 01 03 00 00 00 00 00 00 00 00
                            // 00 06 00 00 00 04 00 04 00 04 00 45 00 00 00 00
                            // 00 00 01 90 00 00 0e 00 50 00 53 00 30 00 32 00
                            // 41 00 57 00 56 00 42 00 31 00 4d 00 43 00 54 00
                            // 31 00 30 00 00 00 0e 00 50 00 53 00 30 00 32 00
                            // 41 00 57 00 56 00 42 00 31 00 4d 00 43 00 54 00
                            // 30 00 33 00 00 02
                            // Contained a text at the end: PS02AWVB1MCT10 and PS02AWVB1MCT03
//                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet");
//                            outputPacket(buf);
                            break;
                        }
                        // 0x08xx: Possibly explicit read request/response
                        case 0x0801: {
                            // Notes:
                            // - Sent from the OS to the controller
                            // Found packets:
                            // 08 01 00 00 00 00 00 00 02 92 00 bc 00 08 00 00
                            // 00 03 00 04 00 0a 00 08 00 06 00 18 00 0a 00 1c
                            // 00 02 00 24 00 02 00 26 00 00 00 00 04 8f 17 b8
                            // ff ff ff ff 00 06 00 4f 00 49 00 4e 00 54 00 45
                            // 00 47 00 00 00 00 00 00 00 02 00 43 00 56 00 00
                            // 00 00 00 00
                            // Contained a text OINTEG
                            break;
                        }
                        case 0x0802: {
                            // Notes:
                            // - Response for a 0x0801 message
                            // - Sent from the Controller to the OS
                            // Found packets:
                            // 08 02 00 00 00 00 00 00 02 92 00 bc 00 00 00 00
                            // 00 00 00 00 00 1a 00 04 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 06 00 00 00 00 00 00
                            // 01 01
                            break;
                        }
                        // Seem to be read subscriptions.
                        // Each 0x1B01 is confirmed by a 0x1B02 containing values, but after a 0x1B01 also we start getting 0x0403 packets with content.
                        case 0x1B01: {
                            // NOTE:
                            // - Seems to be sent as soon as the OS start up and a screen is opened.
                            // - Sent from the OS to the Controller
                            // - Structure of parts is always:
                            //
                            // Found packets:
                            // 1b 01 00 00 00 00 00 07 03 8a 01 c4 00 04 00 1e
                            // 00 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            // 04 cf 1a 0b 00 35 01 2c 00 00 03 e8
                            // 00 00 00 05 00 49 00 44 00 49 00 53 00 50 00 00
                            // ff ff 01
                            //
                            //
                            // Another short one:
                            // 1b 01 00 00 00 00 00 07 03 8c 01 c6 00 04 00 30
                            // 00 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            // 04 cf 1a 0b 00 35 01 2c 00 00 03 e8
                            // 00 01 00 0e 00 43 00 2f 00 53 00 57 00 5f 00 4d 00 41 00 58 00 5f 00 53 00 43 00 41 00 4c 00 45 00 00
                            // ff ff 01
                            //
                            //
                            // When opening a new screen:
                            // 1b 01 00 00 00 00 00 07 03 6d 01 97 00 04 05 04
                            // 00 08 00 00 00 08 01 78 00 08 02 90 00 08 03 fa
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            //
                            // 04 96 17 bf 00 38 01 3f 00 00 03 e8
                            // 00 00 00 0d 00 41 00 2f 00 49 00 53 00 54 00 5f 00 46 00 57 00 5f 00 4b 00 55 00 52 00 5a 00 00
                            // 00 01 00 05 00 53 00 4b 00 41 00 4c 00 31 00 00
                            // 00 02 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 03 00 05 00 53 00 4b 00 41 00 4c 00 32 00 00
                            // 00 04 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 05 00 0b 00 41 00 4b 00 54 00 5f 00 53 00 43 00 48 00 52 00 49 00 54 00 54 00 00
                            // 00 06 00 0e 00 41 00 2f 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 07 00 04 00 41 00 2f 00 42 00 41 00 00
                            // 00 08 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // 00 09 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 32 00 00
                            // 00 0a 00 07 00 47 00 46 00 5f 00 4e 00 41 00 4d 00 45 00 00
                            // 00 0b 00 06 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 0c 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 0d 00 06 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 0e 00 0b 00 41 00 2f 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 0f 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 31 00 00
                            // ff ff
                            // 04 97 17 c0 00 39 01 40 00 00 03 e8
                            // 00 00 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 01 00 05 00 53 00 4b 00 41 00 4c 00 31 00 00
                            // 00 02 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // 00 03 00 06 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 04 00 04 00 41 00 2f 00 41 00 31 00 00
                            // 00 05 00 0d 00 41 00 2f 00 49 00 53 00 54 00 5f 00 46 00 57 00 5f 00 4b 00 55 00 52 00 5a 00 00
                            // 00 06 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 07 00 07 00 47 00 46 00 5f 00 4e 00 41 00 4d 00 45 00 00
                            // 00 08 00 04 00 41 00 2f 00 42 00 41 00 00
                            // 00 09 00 0b 00 41 00 2f 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 0a 00 0e 00 41 00 2f 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 0b 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 31 00 00
                            // ff ff
                            // 04 98 17 c1 00 3a 01 41 00 00 03 e8
                            // 00 00 00 05 00 53 00 4b 00 41 00 4c 00 32 00 00
                            // 00 01 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 31 00 00
                            // 00 02 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // 00 03 00 05 00 53 00 4b 00 41 00 4c 00 31 00 00
                            // 00 04 00 04 00 41 00 2f 00 41 00 31 00 00
                            // 00 05 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 06 00 0b 00 41 00 2f 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 07 00 07 00 47 00 46 00 5f 00 4e 00 41 00 4d 00 45 00 00
                            // 00 08 00 04 00 41 00 2f 00 42 00 41 00 00
                            // 00 09 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 0a 00 06 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 0b 00 06 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 0c 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 32 00 00 00
                            // 0d 00 0e 00 41 00 2f 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 0e 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 0f 00 0d 00 41 00 2f 00 49 00 53 00 54 00 5f 00 46 00 57 00 5f 00 4b 00 55 00 52 00 5a 00 00
                            // ff ff
                            // 04 99 17 c2 00 3b 01 42 00 00 03 e8
                            // 00 00 00 06 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 01 00 05 00 53 00 4b 00 41 00 4c 00 31 00 00
                            // 00 02 00 0e 00 41 00 2f 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 03 00 07 00 47 00 46 00 5f 00 4e 00 41 00 4d 00 45 00 00
                            // 00 04 00 0b 00 41 00 2f 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 05 00 0d 00 41 00 2f 00 49 00 53 00 54 00 5f 00 46 00 57 00 5f 00 4b 00 55 00 52 00 5a 00 00
                            // 00 06 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 07 00 04 00 41 00 2f 00 42 00 41 00 00
                            // 00 08 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // 00 09 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 0a 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 31 00 00
                            // ff ff 04
                            //
                            //
                            // Another big packet:
                            // 1b 01 00 00 00 00 00 07 03 6e 01 98 00 04 05 10
                            // 00 08 00 00 00 08 00 5a 00 08 01 e0 00 08 04 d8
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            //
                            // 04 93 17 bc 00 35 01 3c 00 00 03 e8
                            // 00 06 00 0e 00 41 00 2f 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 07 00 0b 00 41 00 2f 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 08 00 04 00 41 00 2f 00 41 00 33 00 00
                            // ff ff
                            // 04 94 17 bd 00 36 01 3d 00 00 03 e8
                            // 00 00 00 06 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 01 00 05 00 53 00 4b 00 41 00 4c 00 31 00 00
                            // 00 02 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 03 00 0e 00 41 00 2f 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 04 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // 00 05 00 05 00 53 00 4b 00 41 00 4c 00 32 00 00
                            // 00 06 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 07 00 0b 00 41 00 2f 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 08 00 0a 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 09 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 32 00 00
                            // 00 0a 00 04 00 41 00 2f 00 41 00 31 00 00
                            // 00 0b 00 06 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 0c 00 07 00 47 00 46 00 5f 00 4e 00 41 00 4d 00 45 00 00
                            // 00 0d 00 04 00 41 00 2f 00 42 00 41 00 00
                            // 00 0e 00 0d 00 41 00 2f 00 49 00 53 00 54 00 5f 00 46 00 57 00 5f 00 4b 00 55 00 52 00 5a 00 00
                            // 00 0f 00 0b 00 41 00 4b 00 54 00 5f 00 53 00 43 00 48 00 52 00 49 00 54 00 54 00 00
                            // 00 10 00 08 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 31 00 00
                            // ff ff
                            // 04 95 17 be 00 37 01 3e 00 00 03 e8
                            // 00 00 00 17 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 5f 00 46 00 41 00 52 00 42 00 45 00 31 00 00
                            // 00 01 00 19 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 5f 00 45 00 4e 00 41 00 42 00 4c 00 45 00 44 00 32 00 00
                            // 00 02 00 09 00 47 00 46 00 45 00 5f 00 41 00 4c 00 41 00 52 00 4d 00 00
                            // 00 03 00 13 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 04 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // 00 05 00 13 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 06 00 0b 00 49 00 53 00 54 00 5f 00 46 00 57 00 5f 00 4b 00 55 00 52 00 5a 00 00
                            // 00 07 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 08 00 17 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 32 00 00
                            // 00 09 00 17 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 5f 00 46 00 41 00 52 00 42 00 45 00 32 00 00
                            // 00 0a 00 17 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 5f 00 53 00 43 00 41 00 4c 00 45 00 31 00 00
                            // 00 0b 00 0c 00 42 00 41 00 5f 00 47 00 46 00 45 00 53 00 5f 00 41 00 55 00 54 00 4f 00 00
                            // 00 0c 00 02 00 42 00 41 00 00
                            // 00 0d 00 11 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 32 00 00
                            // 00 0e 00 17 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 5f 00 53 00 43 00 41 00 4c 00 45 00 32 00 00
                            // 00 0f 00 19 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 5f 00 45 00 4e 00 41 00 42 00 4c 00 45 00 44 00 31 00 00
                            // 00 10 00 07 00 44 00 59 00 4e 00 5f 00 42 00 45 00 5a 00 00
                            // 00 11 00 17 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 44 00 59 00 4e 00 5f 00 42 00 45 00 53 00 43 00 48 00 31 00 00
                            // 00 12 00 11 00 44 00 59 00 4e 00 5f 00 49 00 53 00 54 00 57 00 45 00 52 00 54 00 45 00 2f 00 4f 00 55 00 54 00 31 00 00
                            // ff ff
                            // 04 9a 17 c3 00 3c 01 43 00 00 03 e8
                            // 00 00 00 06 00 41 00 4c 00 41 00 52 00 4d 00 53 00 00
                            // 00 01 00 09 00 47 00 46 00 5f 00 53 00 54 00 41 00 54 00 55 00 53 00 00
                            // ff ff
                            // 04

                            // Skip the header part.
                            buf.skipBytes(0x76);
                            while(buf.readableBytes() > 12) {
                                byte[] testpointId = new byte[8];
                                buf.readBytes(testpointId);
                                buf.skipBytes(4);

                                // The first block is part of the id of an address.
                                for (short fieldId = buf.readShort(); fieldId != (short) 0xFFFF; fieldId = buf.readShort()) {
                                    // This is followed by a short value providing the length of the following
                                    // string.
                                    short length = buf.readShort();
                                    // Then comes the string, each character encoded as two bytes.
                                    byte[] bytes = new byte[length * 2];
                                    buf.readBytes(bytes);
                                    buf.skipBytes(2);
                                    String address = new String(bytes, StandardCharsets.UTF_16);

                                    String tespoint = Hex.encodeHexString(testpointId);
                                    if(!testpointFieldNames.containsKey(tespoint)) {
                                        testpointFieldNames.put(tespoint, new HashMap<>());
                                    }
                                    testpointFieldNames.get(tespoint).put(fieldId, address);
                                }
                            }
                            break;
                        }
                        /*case 0x1B02: {
                            // NOTE:
                            // - Seems to be sent as soon as the OS start up and a screen is opened.
                            // - Sent from the Controller to the OS
                            // - Seems to be a response to a previous 0x1B01 message
                            // - Size of the response seems to be in direct correlation with the size of the corresponding 0x1B01 message
                            // - Also contains content for each subscription in a 0x1B01 message
                            // - The reason might be that in order to be notified about value changes, the client needs to be informed about the initial value.
                            //
                            // Found packets:
                            // 1b 02 00 00 00 00 00 07 03 8a 01 c4 00 04 00 00
                            // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            // 01 00 00 00 00 00 35 01 2c 04 cf 1a 0b
                            // 03 00 00 00 00 00 1f <!-- Error response with return code? (0x0000001F)
                            // 01
                            //
                            //
                            // Another short one:
                            // 1b 02 00 00 00 00 00 07 03 8c 01 c6 00 04 00 00
                            // 00 06 00 01 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            // 01 00 00 00 00 00 35 01 2c 04 cf 1a 0b
                            // 02 00 01 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                            // 01
                            //
                            //
                            // When opening a new screen:
                            // 1b 02 00 00 00 00 00 07 03 6d 01 97 00 04 00 00
                            // 00 06 00 01 00 06 01 51 00 06 02 3c 00 06 03 59
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            //
                            // 04 00 00 00 00 00 38 01 3f 04 96 17 bf
                            // 02 00 00 47 00 17 00 08 00 4b 00 65 00 69 00 6e 00 65 00 20 00 46 00 57 00 00
                            // 02 00 01 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                            // 02 00 02 76 00 76 05 01 01 00 00 00 0e 00 47 00 46 00 45 00 5f 00 56 00 45 00 52 00 52 00 49 00 45 00 47 00 45 00 4c 00 54 00 00 00 54 00 54 01 07 00 00 2a f1 c7 82 e4 d8 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 03 24 00 0e 40 20 00 00 00 00 00 00 80 16 03
                            // 02 00 04 47 00 11 00 05 00 4d 00 43 00 50 00 30 00 31 00 00
                            // 02 00 05 47 00 07 00 00 00 00
                            // 02 00 06 01 00
                            // 02 00 07 05 00
                            // 02 00 08 63 00 00 19 c4 2f
                            // 02 00 09 08 3f 89 68 4c
                            // 02 00 0a 47 00 19 00 09 00 50 00 54 00 30 00 39 00 41 00 4e 00 41 00 44 00 31 00 00
                            // 02 00 0b 47 00 09 00 01 00 54 00 00
                            // 02 00 0c 47 00 11 00 05 00 4d 00 43 00 54 00 30 00 31 00 00
                            // 02 00 0d 47 00 09 00 01 00 50 00 00
                            // 02 00 0e 05 01
                            // 02 00 0f 08 42 49 aa aa
                            //
                            // 01 00 00 00 00 00 39 01 40 04 97 17 c0
                            // 02 00 00 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 01 24 00 0e 41 28 00 00 00 00 00 00 05 45 01
                            // 02 00 02 63 00 00 19 c4 2f
                            // 02 00 03 47 00 09 00 01 00 46 00 00
                            // 02 00 04 05 01
                            // 02 00 05 47 00 17 00 08 00 4b 00 65 00 69 00 6e 00 65 00 20 00 46 00 57 00 00
                            // 02 00 06 47 00 07 00 00 00 00
                            // 02 00 07 47 00 19 00 09 00 50 00 54 00 30 00 39 00 41 00 44 00 42 00 51 00 31 00 00
                            // 02 00 08 05 00
                            // 02 00 09 05 01
                            // 02 00 0a 01 00
                            // 02 00 0b 21 80 00 00 00 00
                            //
                            // 01 00 00 00 00 00 3a 01 41 04 98 17 c1
                            // 02 00 00 24 00 0e 42 c8 00 00 00 00 00 00 06 34 00
                            // 02 00 01 08 3f 7f ff ff
                            // 02 00 02 63 00 00 19 c4 2f
                            // 02 00 03 24 00 0e 3f 99 99 9a 00 00 00 00 80 16 03
                            // 02 00 04 05 01
                            // 02 00 05 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 06 05 00
                            // 02 00 07 47 00 19 00 09 00 50 00 54 00 30 00 39 00 41 00 44 00 48 00 4f 00 31 00 00
                            // 02 00 08 05 00
                            // 02 00 09 47 00 07 00 00 00 00
                            // 02 00 0a 47 00 07 00 00 00 00
                            // 02 00 0b 47 00 09 00 01 00 50 00 00
                            // 02 00 0c 08 00 00 00 00
                            // 02 00 0d 01 01
                            // 02 00 0e 47 00 0d 00 03 00 56 00 41 00 4b 00 00
                            // 02 00 0f 47 00 17 00 08 00 4b 00 65 00 69 00 6e 00 65 00 20 00 46 00 57 00 00
                            //
                            // 01 00 00 00 00 00 3b 01 42 04 99 17 c2
                            // 02 00 00 47 00 09 00 01 00 46 00 00
                            // 02 00 01 24 00 0e 41 20 00 00 00 00 00 00 05 45 01
                            // 02 00 02 01 00
                            // 02 00 03 47 00 19 00 09 00 50 00 54 00 30 00 39 00 41 00 41 00 4e 00 4b 00 31 00 00
                            // 02 00 04 05 01
                            // 02 00 05 47 00 17 00 08 00 4b 00 65 00 69 00 6e 00 65 00 20 00 46 00 57 00 00
                            // 02 00 06 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 07 05 00
                            // 02 00 08 63 00 00 19 c4 2f
                            // 02 00 09 47 00 07 00 00 00 00
                            // 02 00 0a 08 00 00 00 00
                            //
                            // 01
                            //
                            //
                            // Another big one:
                            // 1b 02 00 00 00 01 00 07 03 6e 01 98 00 04 00 00
                            // 00 06 00 01 00 06 00 1d 00 06 01 52 00 06 02 83
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
                            // 00 00 00 00 00 00 00 00
                            //
                            // 04 00 00 00 00 00 35 01 3c 04 93 17 bc
                            // 02 00 06 01 00
                            // 02 00 07 05 01
                            // 02 00 08 05 01
                            //
                            // 01 00 00 00 00 00 36 01 3d 04 94 17 bd
                            // 02 00 00 47 00 09 00 01 00 50 00 00
                            // 02 00 01 24 00 0e 45 35 40 00 00 00 00 00 82 dc 00
                            // 02 00 02 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 03 01 01
                            // 02 00 04 63 00 00 19 c4 2f
                            // 02 00 05 24 00 0e 43 c8 00 00 c2 c8 00 00 82 dd 01
                            // 02 00 06 47 00 07 00 00 00 00
                            // 02 00 07 05 01
                            // 02 00 08 47 00 11 00 05 00 44 00 52 00 55 00 43 00 4b 00 00
                            // 02 00 09 21 88 c2 20 00 00
                            // 02 00 0a 05 01
                            // 02 00 0b 47 00 09 00 01 00 53 00 00
                            // 02 00 0c 47 00 19 00 09 00 50 00 54 00 30 00 39 00 41 00 44 00 47 00 4f 00 31 00 00
                            // 02 00 0d 05 00
                            // 02 00 0e 47 00 17 00 08 00 4b 00 65 00 69 00 6e 00 65 00 20 00 46 00 57 00 00
                            // 02 00 0f 47 00 0d 00 03 00 41 00 75 00 73 00 00
                            // 02 00 10 21 80 00 00 00 00
                            //
                            // 01 00 00 00 00 00 37 01 3e 04 95 17 be
                            // 02 00 00 05 00
                            // 02 00 01 01 01
                            // 02 00 02 05 00
                            // 02 00 03 47 00 09 00 01 00 54 00 00
                            // 02 00 04 63 00 00 19 c4 2f
                            // 02 00 05 47 00 09 00 01 00 54 00 00
                            // 02 00 06 63 64 00 19 b9 88
                            // 02 00 07 76 00 6a 05 01 01 00 00 00 08 00 53 00 54 00 4f 00 45 00 52 00 55 00 4e 00 47 00 00 00 39 00 39 00 0f 00 00 2a f1 c6 c4 7a 4f ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 08 47 00 0b 00 02 00 52 00 4c 00 00
                            // 02 00 09 05 00
                            // 02 00 0a 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                            // 02 00 0b 01 00
                            // 02 00 0c 05 00
                            // 02 00 0d 08 41 a0 8e 39
                            // 02 00 0e 24 00 0e 43 48 00 00 c2 20 00 00 03 e9 01
                            // 02 00 0f 01 01
                            // 02 00 10 47 00 19 00 09 00 50 00 54 00 30 00 39 00 41 00 4c 00 42 00 43 00 31 00 00
                            // 02 00 11 47 00 0b 00 02 00 56 00 4c 00 00
                            // 02 00 12 08 41 a0 8e 39
                            //
                            // 01 00 00 00 00 00 3c 01 43 04 9a 17 c3
                            // 02 00 00 76 00 58 05 01 01 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 ff 00 00 00 00 00 00 00 00
                            // 02 00 01 63 00 00 19 c4 2f
                            //
                            // 01

                            // Skip the header part.
                            buf.skipBytes(118);
                            for (byte code = buf.readByte(); buf.readableBytes() > 2; code = buf.readByte()) {
                                switch (code) {
                                    case (byte) 0x01:
                                    case (byte) 0x04: {
                                        // Skip the 0x04 block and continue reading the variables.
                                        buf.skipBytes(4);
                                        int blockId = buf.readInt();
                                        buf.skipBytes(4);



                                        return;
                                    }
                                    default: {

                                    }
                                }
                            }
                            break;
                        }*/
                        default: {
                            System.out.println("Got 0x" + Hex.encodeHexString(new byte[]{(byte)(payloadType >> 8), (byte)(payloadType & 0xFF)}) + " packet");
                            outputPacket(buf);
                        }
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }

            if((inputPath != null) && !testpointFieldNames.isEmpty()) {
                try {
                    TimeUnit.MILLISECONDS.sleep(10);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        };

        // Start an Elasticsearch node.
        connection = connect();
        Pong pong = connection.ping();
        if(pong != null) {
            System.out.println("Connected to InfluxDB");
        }
        // Make sure the indexes exist prior to writing to them.
//        prepareIndexes(connection);

        // Start the packet capturing.
        ExecutorService pool = Executors.newScheduledThreadPool(2);
        pool.execute(() -> {
            try {
                receiveHandle.loop(-1, packetListener);
            } catch (PcapNativeException | InterruptedException | NotOpenException e) {
                e.printStackTrace();
            }
        });
    }

    private InfluxDB connect() {
        InfluxDB connection = InfluxDBFactory.connect("http://10.10.64.222:8086", "admin", "password");
        connection.setLogLevel(InfluxDB.LogLevel.BASIC);
        connection.createDatabase("delta-v-demo");
        connection.setDatabase("delta-v-demo");
        connection.enableBatch(BatchOptions.DEFAULTS);
        return connection;
    }

    private void prepareIndexes(InfluxDB connection) {
    }

    private String getTestpointName(String testpointId, Short fieldId) {
        if(testpointFieldNames.containsKey(testpointId)) {
            return testpointFieldNames.get(testpointId).get(fieldId);
        }
        return "-unknown-";
    }

    private Object getTestpointValue(String testpointId, Short fieldId) {
        if(testpointFieldValues.containsKey(testpointId)) {
            return testpointFieldValues.get(testpointId).get(fieldId);
        }
        return null;
    }

    /**
     * Updates a single field value for a given testpoint.
     *
     * @param testpointId id of the testpoint the field belongs to.
     * @param fieldId id of the field.
     * @param newValue value the field is set to.
     */
    private void updateValue(String testpointId, Short fieldId, Object newValue) {
        if(!testpointFieldValues.containsKey(testpointId)) {
            testpointFieldValues.put(testpointId, new HashMap<>());
        }
        testpointFieldValues.get(testpointId).put(fieldId, newValue);
    }

    /**
     * Writes all current values for a fiven testpoint to the database.
     *
     * @param testpointId id of the testpoint.
     */
    private void flushTestpointValues(String testpointId) {
        // If we don't have any values yet, no need to flush
        if(!testpointFieldValues.containsKey(testpointId)) {
            return;
        }
        // If we don't have a name for the testpoint, we shouldn't flush
        if(!testpointNames.containsKey(testpointId)) {
            missingNames.add(testpointId);
            return;
        }

        String testpointName = testpointNames.get(testpointId);

        Point.Builder builder = Point.measurement("Messstelle")
            .time(System.currentTimeMillis(), TimeUnit.MILLISECONDS)
            .tag("MESSSTELLE", testpointName);

        Map<Short, Object> values = testpointFieldValues.get(testpointId);
        for (Map.Entry<Short, Object> fieldValues : values.entrySet()) {
            if(fieldValues.getValue() instanceof Number) {
                builder.addField(getTestpointName(testpointId, fieldValues.getKey()), (Number) fieldValues.getValue());
            } else if(fieldValues.getValue() instanceof Boolean) {
                builder.addField(getTestpointName(testpointId, fieldValues.getKey()), (Boolean) fieldValues.getValue());
            } else {
                System.out.println("unknown type");
            }
        }

        Point point = builder.build();

        // Write the data to the db.
        connection.write(point);
    }

    protected void outputDetectedBlock(String name, ByteBuf byteBuf, int endOfLastBlock) {
        outputDetectedBlock(name, null, null, byteBuf, endOfLastBlock);
    }

    protected void outputDetectedBlock(String name, Object value, Byte param, ByteBuf byteBuf, int endOfLastBlock) {
        int blockSize = byteBuf.readerIndex() - endOfLastBlock;
        byte[] blockContent = new byte[blockSize];
        byteBuf.getBytes(endOfLastBlock, blockContent);
        String valueString = (value != null) ? value.toString() : "";
        String paramString = (param != null) ? Hex.encodeHexString(new byte[] {param}) : "";
        String content = "   " + Hex.encodeHexString(blockContent).replaceAll("(.{2})", "$1 ");
        System.out.println(String.format("Block: %30s %10s %2s %s", name, valueString, paramString, content));
    }

    protected void outputPacket(ByteBuf byteBuf) {
        String packetAsHexString = Hex.encodeHexString(byteBuf.array()).replaceAll("(.{2})", "$1 ").replaceAll("(.{48})", "$1\n");
        System.out.println(packetAsHexString);
    }

    protected void dumpAndExit(ByteBuf byteBuf, int endOfLastBlock, int lastBlockSize, String message) {
        int errorPos = byteBuf.readerIndex();
        int lastBlockStart = errorPos - endOfLastBlock;
        byteBuf.resetReaderIndex();
        System.out.println("-------------------- ERROR --------------------");
        String packetAsHexString = Hex.encodeHexString(byteBuf.array()).replaceAll("(.{2})", "$1 ").replaceAll("(.{48})", "$1\n");
        StringTokenizer stringTokenizer = new StringTokenizer(packetAsHexString, "\n");
        while (stringTokenizer.hasMoreElements()) {
            String line = stringTokenizer.nextToken();
            System.out.println(line);
            if((errorPos < 16) && (errorPos >= 0)) {
                StringBuffer sb = new StringBuffer();
                for(int i = 0; i < errorPos - 1; i++) {
                    sb.append("---");
                }
                sb.append("^");
                System.out.println(sb);
                System.out.println("Last block started: " + lastBlockStart + " bytes before error and had a size of: " + lastBlockSize);
                System.out.println(message);
                System.out.println("\n");
            }
            errorPos -= 16;
        }
        throw new RuntimeException("Error");
    }

    // These seem to be the values decoded for parametrized values ...
    private void decodeParam(byte param) {
        switch (param) {
            case (byte) 0x00: // 00000000
            case (byte) 0x08: // 00001000
            case (byte) 0x88: // 10001000
            case (byte) 0x84: // 10000100
            case (byte) 0xC3: // 11000011
            case (byte) 0x0C: // 00001100
            case (byte) 0x80: // 10000000
            case (byte) 0xC0: // 11000000
                break;
            default:
                throw new RuntimeException("Unexpected param value " + param);
        }

    }

    public static void main(String[] args) throws Exception {
        if(args.length == 0) {
            new PoCInflux(null);
        } else {
            new PoCInflux(args[0]);
        }
    }

}
