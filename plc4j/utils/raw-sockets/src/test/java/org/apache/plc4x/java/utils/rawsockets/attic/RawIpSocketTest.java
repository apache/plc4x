/*
 * Copyright 2014 The Netty Project
 *
 * The Netty Project licenses this file to you under the Apache License,
 * version 2.0 (the "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at:
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.utils.rawsockets.attic;

import org.apache.plc4x.test.RequirePcap;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.pcap4j.core.PcapAddress;
import org.pcap4j.core.PcapNetworkInterface;
import org.pcap4j.core.Pcaps;

import java.net.Inet4Address;
import java.util.Arrays;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.fail;

public class RawIpSocketTest {

    @Test
    @Disabled("Need to make tests run in Docker container first as this test requires libpcap or the entrie application to be run as 'root'")
    @RequirePcap
    public void testPingPacket() throws Exception {
        // Protocol number 1 = ICMP (Ping)
        RawIpSocket rawIpSocket = new RawIpSocket(1);

        CompletableFuture<Boolean> result = new CompletableFuture<>();
        // Simply print the result to the console
        rawIpSocket.addListener(rawData -> {
            System.out.println("Got response:");
            System.out.println(Arrays.toString(rawData));
            result.complete(true);
        });

        // Connect to the remote address
        // This doesn't really connect in the usual sense, it
        // does the ARP MAC address lookup and sets up the listener
        // to accept packets sent from that mac address to the
        // current machines with the given IP protocol id.
        // In this test we simply look for a real network device
        // (The loopback device doesn't have a MAC address)
        // and ping itself.
        for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
            System.out.println("Trying to read on device " + dev);
            if (!dev.getLinkLayerAddresses().isEmpty()) {
                for (PcapAddress pcapAddress : dev.getAddresses()) {
                    if (pcapAddress.getAddress() instanceof Inet4Address) {
                        System.out.println("Trying to connect on PcapAddress " + pcapAddress);
                        rawIpSocket.connect("plc4x.apache.org");
                    }
                }
            }
        }
        // On travis we won't have any interface at all so we don't need to run there
        assertNotNull(rawIpSocket);

        // Simple ICMP (Ping packet)
        byte[] rawData = new byte[]{
            // Type (ICMP Ping Request) & Code (just 0)
            (byte) 0x08, (byte) 0x00,
            // Checksum
            (byte) 0xe3, (byte) 0xe5,
            // Identifier
            (byte) 0x00, (byte) 0x01,
            // Sequence Number
            (byte) 0x00, (byte) 0x00,
            // Payload (Just random data that was used to fit to the checksum)
            (byte) 0x00, (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04, (byte) 0x05, (byte) 0x06, (byte) 0x07, (byte) 0x08, (byte) 0x09};

        // Write the raw packet to the remote host.
        rawIpSocket.write(rawData);

        try {
            result.get(3, TimeUnit.SECONDS);
        } catch (TimeoutException e) {
            fail("Request timed out.");
        }
    }

}
