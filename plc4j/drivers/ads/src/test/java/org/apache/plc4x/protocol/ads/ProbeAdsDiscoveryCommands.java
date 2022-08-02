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
package org.apache.plc4x.protocol.ads;

import org.apache.plc4x.java.ads.discovery.readwrite.*;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.util.Arrays;

public class ProbeAdsDiscoveryCommands {

    public static void main(String[] args) {
        // Create the discovery request message for this device.
        AmsNetId amsNetId = new AmsNetId((byte) 192, (byte) 168, (byte) 23, (byte) 200, (byte) 1, (byte) 1);
        AdsDiscovery discoveryRequestMessage = new AdsDiscovery(1, Operation.ADD_OR_UPDATE_ROUTE_REQUEST, amsNetId, AdsPortNumbers.SYSTEM_SERVICE,
            //Collections.emptyList()
            Arrays.asList(
                new AdsDiscoveryBlockRouteName(new AmsString("route-name")),
                new AdsDiscoveryBlockAmsNetId(amsNetId),
                new AdsDiscoveryBlockUserName(new AmsString("username")),
                new AdsDiscoveryBlockPassword(new AmsString("password")),
                new AdsDiscoveryBlockHostName(new AmsString("host-name-or-ip"))
            ));

        try (DatagramSocket adsDiscoverySocket = new DatagramSocket(AdsDiscoveryConstants.ADSDISCOVERYUDPDEFAULTPORT)) {
            // Serialize the message.
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(discoveryRequestMessage.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            discoveryRequestMessage.serialize(writeBuffer);

            // Get the broadcast address for this interface.
            InetAddress address = InetAddress.getByAddress(new byte[]{(byte) 192, (byte) 168, (byte) 23, (byte) 20});

            // Create the UDP packet to the broadcast address.
            DatagramPacket discoveryRequestPacket = new DatagramPacket(writeBuffer.getBytes(), writeBuffer.getBytes().length, address, AdsDiscoveryConstants.ADSDISCOVERYUDPDEFAULTPORT);
            adsDiscoverySocket.send(discoveryRequestPacket);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
