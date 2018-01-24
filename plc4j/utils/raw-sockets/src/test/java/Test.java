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
import org.pcap4j.core.*;
import org.pcap4j.core.PcapNetworkInterface.PromiscuousMode;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.*;
import org.pcap4j.util.ByteArrays;
import org.pcap4j.util.LinkLayerAddress;
import org.pcap4j.util.MacAddress;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.concurrent.*;

public class Test {

    private static final int SNAPLEN = 65536;
    private static final int READ_TIMEOUT = 10;

    public static void main(String[] args) throws Exception {
        // ICMP has the IP protocol number 1
        int protocolNumber = 1;
        String hostName = "10.10.56.1";
        InetAddress hostIpAddress = InetAddress.getByName(hostName);

        InetAddress localAddress = getLocalNetworkAddress(hostIpAddress);
        if(localAddress == null) {
            System.out.println("Unable to connect to the remote host " + hostName);
            return;
        }
        PcapNetworkInterface nif = Pcaps.getDevByAddress(localAddress);

        /////////////////////////////////////////////////////////////////////////
        // Setup handling incoming packets.
        /////////////////////////////////////////////////////////////////////////

        ExecutorService pool = Executors.newSingleThreadExecutor();
        // Set the filter to only accept packets of the same protocol and
        // addresses at the same MAC address.

        // Filter packets to contain only the ip protocol number of the current protocol.
        StringBuilder ab = new StringBuilder("ip protochain ").append(protocolNumber);
        MacAddress srcMacAddress = null;
        ArrayList<LinkLayerAddress> linkLayerAddresses = nif.getLinkLayerAddresses();
        // Restrict the filter to only replies from the mac address
        // we will be sending to the current devices mac address.
        if(!linkLayerAddresses.isEmpty()) {
            ab.append(" and( ");
            Iterator<LinkLayerAddress> deviceIterator = nif.getLinkLayerAddresses().iterator();
            while(deviceIterator.hasNext()) {
                LinkLayerAddress linkLayerAddress = deviceIterator.next();
                srcMacAddress = MacAddress.getByName(linkLayerAddress.toString(), ":");
                ab.append("ether dst ").append(linkLayerAddress.toString());
                if(deviceIterator.hasNext()) {
                    ab.append(" or ");
                }
            }
            ab.append(")");
        }
        // Create a receive handle for incoming packets
        PcapHandle receiveHandle = nif.openLive(SNAPLEN, PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);
        // Set the filter.
        receiveHandle.setFilter(ab.toString(), BpfProgram.BpfCompileMode.OPTIMIZE);
        // Create a listener to simply output the response on the console.
        CompletableFuture<Packet> result = new CompletableFuture<>();
        PacketListener listener
            = result::complete;
        // Create a task to listen for incoming packets.
        Task t = new Task(receiveHandle, listener);
        pool.execute(t);

        /////////////////////////////////////////////////////////////////////////
        // Send an outgoing packet.
        /////////////////////////////////////////////////////////////////////////

        MacAddress hostMacAddress = getMacAddress(nif, srcMacAddress, localAddress, hostIpAddress);
        if(hostMacAddress == null) {
            throw new RuntimeException("Unable to lookup mac address for host " + hostIpAddress.getHostAddress());
        }
        PcapHandle sendHandle = nif.openLive(SNAPLEN, PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);

        ///////////////////////////////////////////////////////////////////
        // Version using pcap4j built in datatypes
        ///////////////////////////////////////////////////////////////////
        /*byte[] echoData = new byte[10];
        for (int i = 0; i < echoData.length; i++) {
            echoData[i] = (byte)i;
        }

        IcmpV4EchoPacket.Builder echoBuilder = new IcmpV4EchoPacket.Builder();
        echoBuilder.identifier((short)1)
            .payloadBuilder(new UnknownPacket.Builder().rawData(echoData));


        IcmpV4CommonPacket.Builder icmpV4CommonBuilder = new IcmpV4CommonPacket.Builder();
        icmpV4CommonBuilder
            .type(IcmpV4Type.ECHO)
            .code(IcmpV4Code.NO_CODE)
            .payloadBuilder(echoBuilder)
            .correctChecksumAtBuild(true);*/

        ///////////////////////////////////////////////////////////////////
        // Version with hand-crafted icmp byte data
        ///////////////////////////////////////////////////////////////////
        UnknownPacket.Builder packetBuilder = new UnknownPacket.Builder();
        byte[] rawData = new byte[] {
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
        packetBuilder.rawData(rawData);

        IpV4Packet.Builder ipV4Builder = new IpV4Packet.Builder();
        try {
            ipV4Builder
                .version(IpVersion.IPV4)
                .tos(IpV4Rfc791Tos.newInstance((byte)0))
                .ttl((byte)100)
                .protocol(new IpNumber((byte) protocolNumber, "test")/*IpNumber.ICMPV4*/)
                .srcAddr((Inet4Address)InetAddress.getByName(localAddress.getHostAddress()))
                .dstAddr((Inet4Address) hostIpAddress)
                .payloadBuilder(packetBuilder)
                .correctChecksumAtBuild(true)
                .correctLengthAtBuild(true);
        } catch (UnknownHostException e1) {
            throw new IllegalArgumentException(e1);
        }
        ipV4Builder.identification((short)1);

        EthernetPacket.Builder etherBuilder = new EthernetPacket.Builder();
        etherBuilder.dstAddr(hostMacAddress)
            .srcAddr(srcMacAddress)
            .type(EtherType.IPV4)
            .paddingAtBuild(true);
        etherBuilder.payloadBuilder(
            new AbstractPacket.AbstractBuilder() {
                @Override
                public Packet build() {
                    return ipV4Builder.build();
                }
            }
        );

        Packet p = etherBuilder.build();
        System.out.println(p.toString());
        sendHandle.sendPacket(p);

        Packet packet = result.get();
        System.out.println(packet.toString());
    }

    /**
     * Iterate through all devices and find the first that would be able to connect to the given address.
     *
     * @param remoteAddress address we want to connect to.
     * @return PcapNetworkInterface interface that should be able to connect to the given address.
     * @throws Exception something went wrong.
     */
    protected static InetAddress getLocalNetworkAddress(InetAddress remoteAddress) throws Exception {
        byte[] remoteIp = remoteAddress.getAddress();
        // Iterate over all network interfaces.
        for(PcapNetworkInterface dev : Pcaps.findAllDevs()) {
            System.out.println(dev.getName());
            // Iterate over all addresses configured for this interface.
            for(PcapAddress addr : dev.getAddresses()) {
                // Only check addresses matching the IP-version of the remote address.
                if(remoteAddress.getClass().equals(addr.getAddress().getClass())) {
                    byte[] localIp = addr.getAddress().getAddress();
                    byte[] netMask = addr.getNetmask().getAddress();
                    boolean matches = true;
                    // Iterate over all bytes of the address and see if they match
                    // after applying the net mask filter.
                    for(int i = 0; i < localIp.length; i++) {
                        if((localIp[i] & netMask[i]) != (remoteIp[i] & netMask[i])) {
                            matches = false;
                            break;
                        }
                    }
                    // If the current address would be able to connect to the remote
                    // address, return this device.
                    if(matches) {
                        return addr.getAddress();
                    }
                }
            }
        }
        return null;
    }

    protected static MacAddress getMacAddress(PcapNetworkInterface nif, MacAddress localMacAddress, InetAddress localIpAddress, InetAddress remoteIpAddress) throws Exception {
        PcapHandle receiveHandle
            = nif.openLive(SNAPLEN, PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);
        PcapHandle sendHandle
            = nif.openLive(SNAPLEN, PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);

        ExecutorService pool = Executors.newSingleThreadExecutor();

        try {
            receiveHandle.setFilter(
                "arp and src host " + remoteIpAddress.getHostAddress()
                    + " and dst host " + localIpAddress.getHostAddress()
                    + " and ether dst " + Pcaps.toBpfString(localMacAddress),
                BpfProgram.BpfCompileMode.OPTIMIZE
            );

            CompletableFuture<MacAddress> resolvedAddress = new CompletableFuture<>();
            PacketListener listener
                = packet -> {
                    if (packet.contains(ArpPacket.class)) {
                        ArpPacket arp = packet.get(ArpPacket.class);
                        if (arp.getHeader().getOperation().equals(ArpOperation.REPLY)) {
                            resolvedAddress.complete(arp.getHeader().getSrcHardwareAddr());
                        }
                    }
                };

            Task t = new Task(receiveHandle, listener);
            pool.execute(t);

            ArpPacket.Builder arpBuilder = new ArpPacket.Builder();
            arpBuilder
                .hardwareType(ArpHardwareType.ETHERNET)
                .protocolType(EtherType.IPV4)
                .hardwareAddrLength((byte)MacAddress.SIZE_IN_BYTES)
                .protocolAddrLength((byte) ByteArrays.INET4_ADDRESS_SIZE_IN_BYTES)
                .operation(ArpOperation.REQUEST)
                .srcHardwareAddr(localMacAddress)
                .srcProtocolAddr(localIpAddress)
                .dstHardwareAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
                .dstProtocolAddr(remoteIpAddress);

            EthernetPacket.Builder etherBuilder = new EthernetPacket.Builder();
            etherBuilder.dstAddr(MacAddress.ETHER_BROADCAST_ADDRESS)
                .srcAddr(localMacAddress)
                .type(EtherType.ARP)
                .payloadBuilder(arpBuilder)
                .paddingAtBuild(true);

            Packet p = etherBuilder.build();
            sendHandle.sendPacket(p);
            try {
                return resolvedAddress.get(10000, TimeUnit.MILLISECONDS);
            } catch (TimeoutException e) {
                return null;
            }
        } finally {
            /*if (receiveHandle.isOpen()) {
                receiveHandle.close();
            }
            if (sendHandle.isOpen()) {
                sendHandle.close();
            }
            if (!pool.isShutdown()) {
                pool.shutdownNow();
            }*/
        }
    }

    private static class Task implements Runnable {

        private PcapHandle handle;
        private PacketListener listener;

        public Task(PcapHandle handle, PacketListener listener) {
            this.handle = handle;
            this.listener = listener;
        }

        @Override
        public void run() {
            try {
                handle.loop(-1, listener);
            } catch (PcapNativeException e) {
                e.printStackTrace();
            } catch (InterruptedException e) {
            } catch (NotOpenException e) {
                e.printStackTrace();
            }
        }
    }

}
