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
import org.pcap4j.core.*;
import org.pcap4j.packet.UdpPacket;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.Map;
import java.util.StringTokenizer;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class PoC {

    private static final Logger valueLogger = LoggerFactory.getLogger(PoC.class);

    private static final int SNAPLEN = 65536;
    private static final int READ_TIMEOUT = 10;

    private PcapHandle receiveHandle;

    private PoC() throws Exception {
        PcapNetworkInterface nif = null;
        for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
            if("en7".equals(dev.getName())) {
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

        // Set the filter.
        String filterString = "udp port 18507";
        receiveHandle.setFilter(filterString, BpfProgram.BpfCompileMode.OPTIMIZE);

        Map<String, Object> values = new HashMap<>();

        byte[] timeBytes = ByteBuffer.allocate(8).putLong(System.currentTimeMillis()).array();
        System.out.println("Current Time: " + Hex.encodeHexString(timeBytes));

        PacketListener packetListener = packet -> {
            try {
                UdpPacket udpPacket = (UdpPacket) packet.getPayload().getPayload();
                ByteBuf dis = Unpooled.wrappedBuffer(udpPacket.getPayload().getRawData());
                //DataInputStream dis = new DataInputStream(new ByteArrayInputStream(udpPacket.getPayload().getRawData()));
                dis.skipBytes(4);
                short messageType = dis.readShort();
                // We're only interested in type 2 messages.
                if(messageType == 0x0002) {
                    dis.skipBytes(10);
                    short payloadType = dis.readShort();
                    if(payloadType == 0x0403) {
                        System.out.println("----------------------------------------------------------------------------------------");
//                        System.out.println(Hex.encodeHexString(udpPacket.getPayload().getRawData()).replaceAll("(.{2})", "$1 ").replaceAll("(.{48})", "$1\n"));
//                        System.out.println("----------------------");
                        // Skip the rest of the header.
                        dis.skipBytes(39);
                        int endOfLastBlock = dis.readerIndex();
                        int lastBlockSize = 0;
                        short currentContext = 0;
                        for(byte code = dis.readByte(); dis.readableBytes() > 2; code = dis.readByte()) {
                            short blockId = dis.readShort();
                            byte type = dis.readByte();

                            // First check the code of the next block ...
                            switch (code) {
                                case (byte) 0x01: {
                                    switch (type) {
                                        case (byte) 0x01: {
                                            // - It seems that the ids of a variable seem to occur multiple times
                                            // - Also does it seem that this type of block sets some sort of context for following blocks
                                            // - After setting up a machine with a new OS, the type of every of these is 0x00

                                            // Found blocks:
                                            // 01 00 23 01 1a 04 32 1c fd (size 7)
                                            currentContext = blockId;
                                            dis.skipBytes(5);
                                            outputDetectedBlock("-- Switch Context --", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x00: {
                                            // Is seems this simply signals the end of a packet.
                                            currentContext = blockId;
                                            dis.skipBytes(5);
                                            outputDetectedBlock("-- Switch Context --", dis, endOfLastBlock);
                                            break;
                                        }
                                        default: {
                                            dumpAndExit(dis, endOfLastBlock, lastBlockSize, "Unexpected 0x01 type code: " + Hex.encodeHexString(new byte[]{type}));
                                        }
                                    }
                                    break;
                                }
                                case (byte) 0x02: {
                                    // Now inspect the block content ...
                                    switch (type) {
                                        case (byte) 0x01: {
                                            // Possibly boolean value?
                                            dis.skipBytes(1);
                                            outputDetectedBlock("BOOL value", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x03: {
                                            dis.skipBytes(5);
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x05: {
                                            // NOTE:
                                            // - Each packet seems to have one of these
                                            // - For each following packet the content is identical
                                            // Found Block:
                                            // 02 00 0c 05: 00 02 00 13 63 00 00 69 9c 1a
                                            // 02 00 0c 05: 00 01 00 47 00 64 04 2a 17 53
                                            dis.skipBytes(10);
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x06: {
                                            // Possibly Parse 16 bit int?
                                            String id = "(U)INT-" + currentContext + "-" + blockId;
                                            short shortValue = dis.readShort();
                                            outputDetectedBlock("(U)INT value", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x07: {
                                            // Possibly Parse 32 bit int?
                                            String id = "(U)DINT-" + currentContext + "-" + blockId;
                                            int intValue = dis.readInt();
                                            outputDetectedBlock("(U)DINT value", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x08: {
                                            // Parse float
                                            String id = "REAL-" + currentContext + "-" + blockId;
                                            float floatValue = dis.readFloat();
                                            outputDetectedBlock("REAL value", dis, endOfLastBlock);
                                            floatValue = Math.round(floatValue * 100.0f) / 100.0f;
                                            if (!values.containsKey(id)) {
                                                valueLogger.info(String.format("Variable with id: %s set to: %f", id, floatValue));
                                                values.put(id, floatValue);
                                            } else if (!values.get(id).equals(floatValue)) {
                                                float oldValue = (float) values.get(id);
                                                valueLogger.info(String.format("Variable with id: %s changed from: %f to: %f", id, oldValue, floatValue));
                                                values.put(id, floatValue);
                                            }
                                            break;
                                        }
                                        case (byte) 0x21: {
                                            // From having a look at the byte values these could be 32bit floating point values with some sort of parameters
                                            String id = "REAL(P)-" + currentContext + "-" + blockId;
                                            byte param = dis.readByte();
                                            decodeParam(param);
                                            float floatValue = dis.readFloat();
                                            outputDetectedBlock("REAL(P) value", dis, endOfLastBlock);
                                            floatValue = Math.round(floatValue * 100.0f) / 100.0f;
                                            if (!values.containsKey(id)) {
                                                valueLogger.info(String.format("Variable with id: %s set to: %f with params %s", id, floatValue, Hex.encodeHexString(new byte[]{param})));
                                                values.put(id, floatValue);
                                            } else if (!values.get(id).equals(floatValue)) {
                                                float oldValue = (float) values.get(id);
                                                valueLogger.info(String.format("Variable with id: %s changed from: %f to: %f with params %s", id, oldValue, floatValue, Hex.encodeHexString(new byte[]{param})));
                                                values.put(id, floatValue);
                                            }
                                            break;
                                        }
                                        case (byte) 0x22: {
                                            // Parse boolean (From what I learnt, this could be a flagged boolean, where the first byte is some sort of param)
                                            String id = "BOOL(P)-" + currentContext + "-" + blockId;
                                            byte param = dis.readByte();
                                            decodeParam(param);
                                            byte booleanByteValue = dis.readByte();
                                            outputDetectedBlock("BOOL(P) value", dis, endOfLastBlock);
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
                                            if (!values.containsKey(id)) {
                                                valueLogger.info(String.format("Variable with id: %s set to: %b with params %s", id, booleanValue, Hex.encodeHexString(new byte[]{param})));
                                                values.put(id, booleanValue);
                                            } else if (!values.get(id).equals(booleanValue)) {
                                                boolean oldValue = (boolean) values.get(id);
                                                valueLogger.info(String.format("Variable with id: %s changed from: %b to: %b with params %s", id, oldValue, booleanValue, Hex.encodeHexString(new byte[]{param})));
                                                values.put(id, booleanValue);
                                            }
                                            break;
                                        }
                                        case (byte) 0x24: {
                                            // No idea what this type is.
                                            // NOTE:
                                            // - It seems that the last byte seems to mirror the id of the block (Maybe the field id is just one byte and not a short)
                                            // - It seems that these blocks are contained in every packet.
                                            byte[] tmp = new byte[13]; // Has to be 13 in case of 0x0201 but some times 12
                                            dis.readBytes(tmp);
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x25: {
                                            dis.skipBytes(6);
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x47: {
                                            // No idea what this type is.
                                            // NOTE:
                                            // - Seems to be sent as soon as a user confirms an alarm.
                                            // - Seems the length is variable
                                            // - Seems content is terminated by a "0x0000" value
                                            // - All content seems to be encoded as short values with the first byte set to "0x00".
                                            //
                                            // Found Blocks:
                                            // 00 4b 00 22 00 49 00 6e 00 69 00 74 00 69 00 61
                                            // 00 6c 00 69 00 73 00 69 00 65 00 72 00 75 00 6e
                                            // 00 67 00 20 00 2e 00 2e 00 2e 00 2e 00 2e 00 20
                                            // 00 62 00 69 00 74 00 74 00 65 00 20 00 77 00 61
                                            // 00 72 00 74 00 65 00 6e 00 00 02 00 67 47 00 1d
                                            // 00 0b 00 57 00 41 00 52 00 54 00 45 00 4e 00 20
                                            // 00 2e 00 2e 00 2e 00 20 00 00
                                            short val = dis.readShort();
                                            while(val != 0x0000) {
                                                val = dis.readShort();
                                            }
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x48: {
                                            // No idea what this type is.
                                            // NOTE:
                                            // - Seems to be sent as soon as an alarm is fired, changed or removed from the controller.
                                            // - There seem to be only two types of values: 0x8000 and 0x8001
                                            byte[] tmp = new byte[2];
                                            dis.readBytes(tmp);
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x49: {
                                            // - Judging from the 0x80 first byte I would assume this is again one of these parametrized values
                                            // - Would suggest this is a 32 bit integer value.
                                            // Found blocks:
                                            // 80 00 00 06 0d
                                            String id = "(U)DINT(P)-" + currentContext + "-" + blockId;
                                            byte param = dis.readByte();
                                            decodeParam(param);
                                            int intValue = dis.readInt();
                                            if (!values.containsKey(id)) {
                                                valueLogger.info(String.format("Variable with id: %s set to: %d with params %s", id, intValue, Hex.encodeHexString(new byte[]{param})));
                                                values.put(id, intValue);
                                            } else if (!values.get(id).equals(intValue)) {
                                                int oldValue = (int) values.get(id);
                                                valueLogger.info(String.format("Variable with id: %s changed from: %d to: %d with params %s", id, oldValue, intValue, Hex.encodeHexString(new byte[]{param})));
                                                values.put(id, intValue);
                                            }
                                            outputDetectedBlock("(U)DINT(P) value", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x5B: {
                                            // No idea what this type is.
                                            dis.readShort();
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                        }
                                        case (byte) 0x63: {
                                            // No idea what this type is.
                                            // NOTE:
                                            // - It seems that this block is contained in every packet exactly once for id 5
                                            // Found blocks:
                                            // 02 00 06 63: 64 00 19 b9 88
                                            byte[] tmp = new byte[5];
                                            dis.readBytes(tmp);
//                                            System.out.println(String.format("Got 0x63 type for id %s with content: %s", blockId, Hex.encodeHexString(tmp)));
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x75: {
                                            // No idea what this type is.
                                            // NOTE:
                                            // - Exactly 3 blocks of this type with extremely similar content is being sent every 60 seconds for the ids: 17, 16 and 34
                                            //                            001600280d0100000000280015f360000000000100
                                            //                            001600280d0100000000280015f360000000000100
                                            int size = "001600280d0100000000280015f360000000000100".length() / 2; //21
                                            byte[] tmp = new byte[size];
                                            dis.readBytes(tmp);
//                                            System.out.println(String.format("Got 0x75 type for id %s with content: %s", blockId, Hex.encodeHexString(tmp)));
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0x76: {
                                            // No idea what this type is.
                                            // These strange blocks containing a repeating pattern of 0x00 and 0xFF
                                            // NOTE:
                                            // - These blocks seem to be transferred whenever a boolean value is changed.
                                            // - There seem to be two variants:
                                            //   - Variant 1 (shorter) is transferred as soon as a boolean value is set
                                            //   - Variant 2 (longer) is transferred as soon as a boolean values is unset
                                            // - Variant always looks the same no matter what combination of boolean values is set
                                            // - The blocks always refer to ids 0, 1 and 2
                                            // - The additional part of Variant 2 always starts with:
                                            //   "000700420049004e005f0041004c004d000000180018000300002ae7"
                                            //   The last 4 bytes (maybe more) seem to be an always increasing value
                                            //   (Maybe some sort of timestamp)
                                            short length = (short) (dis.readShort() - 3);
                                            byte[] tmp = new byte[length];
                                            dis.readBytes(tmp);
                                            String hexBlock = Hex.encodeHexString(tmp).replaceAll("(.{32})", "$1\n");
//                                            System.out.println(String.format("Got 0x76 type for id %s with content: \n%s", blockId, hexBlock));
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        case (byte) 0xF6: {
                                            // Only seen in 0x0102 blocks
                                            dis.skipBytes(4);
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);
                                            break;
                                        }
                                        default: {
                                            dumpAndExit(dis, endOfLastBlock, lastBlockSize, "Unexpected 0x02 type code: " + Hex.encodeHexString(new byte[]{type}));
                                            /*if(code == (byte) 0x01) {
                                                dis.skipBytes(4);
                                            } else {
                                                dumpAndExit(dis, endOfLastBlock, lastBlockSize, "Unknown variable type 0x" + Hex.encodeHexString(new byte[]{type}));
                                            }
                                            outputDetectedBlock("Unknown", dis, endOfLastBlock);*/
                                        }

                                    }
                                    break;
                                }
                                case (byte) 0x03: {
                                    // TODO: Check if these other types still exist ..
                                    // Found blocks:
                                    // 03 00 23 00 00 00 4a             (size 6)
                                    // 03 01 00 27 01 1e 04 36 1d       (size 8)
                                    // 03 01 00 24 01 1b 04 33 1c fe    (size 9)
                                    switch (type) {
                                        case (byte) 0x00: {
                                            dis.skipBytes(3);
                                            break;
                                        }
                                        default: {
                                            dumpAndExit(dis, endOfLastBlock, lastBlockSize, "Unexpected 0x03 type code: " + Hex.encodeHexString(new byte[]{type}));
                                        }
                                    }
                                    break;
                                }
                                default: {
                                    dumpAndExit(dis, endOfLastBlock, lastBlockSize, "Unexpected code: " + Hex.encodeHexString(new byte[]{code}));
                                }
                            }
                            lastBlockSize = dis.readerIndex() - endOfLastBlock;
                            endOfLastBlock = dis.readerIndex();
                        }

                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        };

        ExecutorService pool = Executors.newScheduledThreadPool(2);
        pool.execute(() -> {
            try {
                receiveHandle.loop(-1, packetListener);
            } catch (PcapNativeException | InterruptedException | NotOpenException e) {
                e.printStackTrace();
            }
        });
    }

    protected void outputDetectedBlock(String name, ByteBuf byteBuf, int endOfLastBlock) {
        int blockSize = byteBuf.readerIndex() - endOfLastBlock;
        byte[] blockContent = new byte[blockSize];
        byteBuf.getBytes(endOfLastBlock, blockContent);
        String content = "   " + Hex.encodeHexString(blockContent).replaceAll("(.{2})", "$1 ");
        System.out.println(String.format("Block: %20s %s", name, content));
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
        new PoC();
    }

}
