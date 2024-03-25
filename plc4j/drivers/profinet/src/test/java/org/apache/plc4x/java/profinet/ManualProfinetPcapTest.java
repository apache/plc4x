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

package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.profinet.readwrite.PnDcp_Pdu;
import org.apache.plc4x.java.profinet.readwrite.PnDcp_Pdu_RealTimeCyclic;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.Dot1qVlanTagPacket;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.Packet;
import org.pcap4j.packet.UnknownPacket;

import java.io.EOFException;
import java.sql.Timestamp;
import java.util.Arrays;
import java.util.concurrent.TimeoutException;

public class ManualProfinetPcapTest {

    public static void main(String[] args) throws Exception {
        try (PcapHandle handle = Pcaps.openOffline("/Users/cdutz/Projects/Apache/PLC4X/profinet.pcapng", PcapHandle.TimestampPrecision.NANO);){
            int lastIncomingCycleCounter = 0;
            double averageIncomingCycleCounter = 0.0;
            int numberIncomingPackets = 0;
            Timestamp lastIncomingCycleTime = null;
            int lastOutgoingCycleCounter = 0;
            Timestamp lastOutgoingCycleTime = null;
            while (true) {
                try {
                    Packet packet = handle.getNextPacketEx();
                    Timestamp timestamp = handle.getTimestamp();

                    EthernetPacket.EthernetHeader packetHeader = (EthernetPacket.EthernetHeader) packet.getHeader();
                    boolean fromDevice = Arrays.equals(new byte[]{(byte) 0xF8, (byte) 0xE4, (byte) 0x3B, (byte) 0xB6, (byte) 0x9B, (byte) 0xBF}, packetHeader.getSrcAddr().getAddress());

                    // All packets from the device have VLan wrapped around the packet.
                    if(fromDevice) {
                        if(packet.getPayload() instanceof Dot1qVlanTagPacket) {
                            Dot1qVlanTagPacket vlanTagPacket = (Dot1qVlanTagPacket) packet.getPayload();
                            UnknownPacket payload = (UnknownPacket) vlanTagPacket.getPayload();
                            ReadBufferByteBased readBuffer = new ReadBufferByteBased(payload.getRawData());
                            PnDcp_Pdu pnDcpPdu = PnDcp_Pdu.staticParse(readBuffer);
                            if(pnDcpPdu instanceof PnDcp_Pdu_RealTimeCyclic) {
                                PnDcp_Pdu_RealTimeCyclic pnDcpPdu1 = (PnDcp_Pdu_RealTimeCyclic) pnDcpPdu;

                                if(lastIncomingCycleCounter != 0) {
                                    int lastCycles;
                                    if(pnDcpPdu1.getCycleCounter() > lastIncomingCycleCounter) {
                                        lastCycles = pnDcpPdu1.getCycleCounter() - lastIncomingCycleCounter;
                                    } else {
                                        lastCycles = (pnDcpPdu1.getCycleCounter() + 0xFFFF) - lastIncomingCycleCounter;
                                    }
                                    averageIncomingCycleCounter = (numberIncomingPackets * averageIncomingCycleCounter + lastCycles) / (numberIncomingPackets + 1);

                                    if (lastIncomingCycleTime != null) {
                                        System.out.printf("--> %3d %3d            %f\n", lastCycles, (timestamp.getNanos() - lastIncomingCycleTime.getNanos()) / 1000000, averageIncomingCycleCounter);
                                    }
                                }

                                numberIncomingPackets++;
                                lastIncomingCycleCounter = pnDcpPdu1.getCycleCounter();
                                lastIncomingCycleTime = timestamp;
                            }
                        } else {
                            System.out.println("Other packet");
                        }
                    } else if (packet.getPayload() instanceof UnknownPacket) {
                        UnknownPacket payload = (UnknownPacket) packet.getPayload();
                        ReadBufferByteBased readBuffer = new ReadBufferByteBased(payload.getRawData());
                        PnDcp_Pdu pnDcpPdu = PnDcp_Pdu.staticParse(readBuffer);
                        if(pnDcpPdu instanceof PnDcp_Pdu_RealTimeCyclic) {
                            PnDcp_Pdu_RealTimeCyclic pnDcpPdu1 = (PnDcp_Pdu_RealTimeCyclic) pnDcpPdu;
                            if(lastOutgoingCycleCounter != 0) {
                                int lastCycles;
                                if(pnDcpPdu1.getCycleCounter() > lastOutgoingCycleCounter) {
                                    lastCycles = pnDcpPdu1.getCycleCounter() - lastOutgoingCycleCounter;
                                } else {
                                    lastCycles = (pnDcpPdu1.getCycleCounter() + 0xFFFF) - lastOutgoingCycleCounter;
                                }

                                if (lastOutgoingCycleTime != null) {
                                    System.out.printf("<--          %3d %3d\n", lastCycles, (timestamp.getNanos() - lastOutgoingCycleTime.getNanos()) / 1000000);
                                }
                            }
                            lastOutgoingCycleCounter = pnDcpPdu1.getCycleCounter();
                            lastOutgoingCycleTime = timestamp;
                        }
                    } else {
                        System.out.println("Other packet");
                    }
                } catch (TimeoutException e) {
                    System.out.println("TIMEOUT");
                } catch (EOFException e) {
                    System.out.println("EOF");
                    break;
                }
            }
        }
    }

}
