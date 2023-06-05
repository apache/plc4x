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

package org.apache.plc4x.java.profinet.messages;

import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.pcap4j.packet.EthernetPacket;
import org.pcap4j.packet.IllegalRawDataException;
import org.pcap4j.packet.Packet;

import java.util.Collections;

public class DiscoveryMessageFactory {

    public static Ethernet_Frame createDiscoveryRequestMessage(MacAddress localMacAddress, MacAddress remoteMacAddress) {
        return new Ethernet_Frame(
            remoteMacAddress,
            localMacAddress,
            new Ethernet_FramePayload_VirtualLan(VirtualLanPriority.BEST_EFFORT, false, 0,
                new Ethernet_FramePayload_PnDcp(
                    new PnDcp_Pdu_IdentifyReq(PnDcp_FrameId.DCP_Identify_ReqPDU.getValue(),
                        1,
                        256,
                        Collections.singletonList(
                            new PnDcp_Block_ALLSelector()
                        )))));
    }

    public static Packet createDiscoveryRequestPacket(MacAddress localMacAddress, MacAddress remoteMacAddress) {
        Ethernet_Frame discoveryRequest = createDiscoveryRequestMessage(localMacAddress, remoteMacAddress);
        WriteBufferByteBased buffer = new WriteBufferByteBased(discoveryRequest.getLengthInBytes());
        try {
            discoveryRequest.serialize(buffer);
        } catch (SerializationException e) {
            throw new RuntimeException(e);
        }
        try {
            byte[] payload = buffer.getBytes();
            return EthernetPacket.newPacket(payload, 0, payload.length);
        } catch (IllegalRawDataException e) {
            throw new RuntimeException(e);
        }
    }

}
