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
import java.util.Arrays;
import java.util.concurrent.TimeoutException;

public class ManualProfinetPcapTest {

    public static void main(String[] args) throws Exception {
        try (PcapHandle handle = Pcaps.openOffline("/Users/cdutz/Projects/Apache/PLC4X/profinet-slow.pcapng", PcapHandle.TimestampPrecision.NANO);){
            int lastIncomingCycleTime = 0;
            int lastOutgoingCycleTime = 0;
            int minDelay = 65000;
            int maxDelay = 0;
            while (true) {
                try {
                    Packet packet = handle.getNextPacketEx();
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
                                lastIncomingCycleTime = pnDcpPdu1.getCycleCounter();
                                //System.out.printf("--> %d\n", pnDcpPdu1.getCycleCounter());
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
                            int difference = (pnDcpPdu1.getCycleCounter() < lastIncomingCycleTime) ? lastIncomingCycleTime - pnDcpPdu1.getCycleCounter() : pnDcpPdu1.getCycleCounter() - lastIncomingCycleTime;
                            if(difference < 60000 && difference > 100) {
                                if (difference < minDelay) {
                                    minDelay = difference;
                                }
                                if (difference > maxDelay) {
                                    maxDelay = difference;
                                }
                                System.out.printf("<-- %10d %10d   %10d-%10d\n", difference, pnDcpPdu1.getCycleCounter() - lastOutgoingCycleTime, minDelay, maxDelay);
                            }
                            lastOutgoingCycleTime = pnDcpPdu1.getCycleCounter();
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
