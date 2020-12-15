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

import org.apache.commons.lang3.SystemUtils;
import org.apache.plc4x.java.utils.pcap.netty.exception.PcapException;
import org.pcap4j.core.*;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.*;
import org.pcap4j.util.ByteArrays;
import org.pcap4j.util.LinkLayerAddress;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.*;

public class RawIpSocket {

    private static final Logger logger = LoggerFactory.getLogger(RawIpSocket.class);

    private static final int SNAPLEN = 65536;
    private static final int READ_TIMEOUT = 10000;

    private static final String GATEWAY_ONLY_NETMASK = "255.255.255.255";

    private static final Map<InetAddress, MacAddress> arpCache = new HashMap<>();

    // The id of the protocol we will be communicating in.
    private final int protocolNumber;

    private PcapNetworkInterface nif;
    private InetAddress remoteIpAddress;
    private MacAddress firstHopMacAddress;
    private InetAddress localIpAddress;
    private MacAddress localMacAddress;
    private ExecutorService pool = Executors.newSingleThreadExecutor();
    private PcapHandle receiveHandle;

    private final List<RawSocketListener> listeners = new LinkedList<>();

    public RawIpSocket(int protocolNumber) {
        this.protocolNumber = protocolNumber;
    }

    public void connect(String remoteAddress) throws PcapException {
        try {
            pool = Executors.newScheduledThreadPool(2);

            remoteIpAddress = InetAddress.getByName(remoteAddress);

            // As we have to create the Ethernet packets ourselves, and
            // in case of non local remote addresses we need to go through
            // routers and gateways, we need do differentiate between the
            // next ethernet node and the remote ip address.
            //
            // We therefore need to know the following information:
            // 1. Can we connect to the remote directly (a) or do we
            //    need a gateway (b)?
            // 2. The local IP and Mac Address of our device
            // 3a. The remote IP and Mac address of the target device
            // 3b. The remote IP of the target device & the Mac address
            //     of the gateway.

            // Check if we can connect directly to the destination.
            FirstHop firstHop = getFirstHop(remoteIpAddress);
            if (firstHop == null) {
                // If this wouldn't work, try to figure out the default
                // gateway and use that as next hop.
                InetAddress defaultGatewayAddress = getDefaultGatewayAddress();
                if (defaultGatewayAddress != null) {
                    firstHop = getFirstHop(defaultGatewayAddress);
                    if (firstHop == null) {
                        // If this didn't work, we simply can't reach the
                        // destination and give up with an exception. Not
                        // much we can do here.
                        throw new PcapException("Unable to connect to " + remoteAddress);
                    }
                } else {
                    throw new PcapException("Unable to connect to " + remoteAddress + " no default gateway");
                }
            }

            nif = firstHop.networkInterface;

            if (nif.isLoopBack()) {
                throw new PcapException("Can't use RawSocket on loopback devices");
            }

            localMacAddress = MacAddress.getByAddress(firstHop.localMacAddress.getAddress());
            localIpAddress = firstHop.localInetAddress;

            firstHopMacAddress = MacAddress.getByAddress(firstHop.remoteMacAddress.getAddress());

            // Setup receiving of packets and redirecting them to the corresponding listeners.
            // Filter packets to contain only the ip protocol number of the current protocol.
            receiveHandle = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);

            // Set the filter.
            String filterString = "ip protochain " + protocolNumber +
                " and ether dst " + localMacAddress.toString() +
                " and ip dst " + localIpAddress.getHostAddress() +
                " and ether src " + firstHopMacAddress.toString() +
                " and ip src " + remoteIpAddress.getHostAddress();

            receiveHandle.setFilter(filterString, BpfProgram.BpfCompileMode.OPTIMIZE);
            PacketListener packetListener = packet -> {
                for (RawSocketListener listener : listeners) {
                    listener.packetReceived(packet.getRawData());
                }
            };

            pool.execute(() -> {
                try {
                    receiveHandle.loop(-1, packetListener);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    logger.error("Error receiving packet for protocol {} from MAC address {}",
                        protocolNumber, firstHopMacAddress, e);
                } catch (PcapNativeException | NotOpenException e) {
                    logger.error("Error receiving packet for protocol {} from MAC address {}",
                        protocolNumber, firstHopMacAddress, e);
                }
            });
        } catch (PcapNativeException | NotOpenException | UnknownHostException e) {
            throw new PcapException("Error setting up RawSocket", e);
        }
    }

    public void disconnect() throws PcapException {
        // TODO: Terminate all the listeners and the thread pool.
    }

    public void write(byte[] rawData) throws PcapException {
        try (PcapHandle sendHandle =
                 nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT)) {
            UnknownPacket.Builder packetBuilder = new UnknownPacket.Builder();
            packetBuilder.rawData(rawData);

            IpV4Packet.Builder ipV4Builder = new IpV4Packet.Builder();
            ipV4Builder
                .version(IpVersion.IPV4)
                .tos(IpV4Rfc791Tos.newInstance((byte) 0))
                .ttl((byte) 100)
                .protocol(new IpNumber((byte) protocolNumber, "plc4x"))
                .srcAddr((Inet4Address) localIpAddress)
                .dstAddr((Inet4Address) remoteIpAddress)
                .payloadBuilder(packetBuilder)
                .correctChecksumAtBuild(true)
                .correctLengthAtBuild(true);
            ipV4Builder.identification((short) 1);

            EthernetPacket.Builder etherBuilder = new EthernetPacket.Builder();
            etherBuilder.dstAddr(firstHopMacAddress)
                .srcAddr(localMacAddress)
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
            sendHandle.sendPacket(p);
        } catch (PcapNativeException | NotOpenException e) {
            throw new PcapException("Error sending packet.", e);
        }
    }

    public void addListener(RawSocketListener listener) {
        listeners.add(listener);
    }

    public void removeListener(RawSocketListener listener) {
        listeners.remove(listener);
    }

    private MacAddress getMacAddress(PcapNetworkInterface dev, InetAddress localIpAddress, InetAddress remoteIpAddress) throws PcapException {
        if (!arpCache.containsKey(remoteIpAddress)) {
            MacAddress macAddress = lookupMacAddress(dev, localIpAddress, remoteIpAddress);
            arpCache.put(remoteIpAddress, macAddress);
            return macAddress;
        }
        return arpCache.get(remoteIpAddress);
    }

    private MacAddress lookupMacAddress(PcapNetworkInterface dev, InetAddress localIpAddress, InetAddress remoteIpAddress) throws PcapException {
        try (PcapHandle receiveHandle =
                 dev.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);
             PcapHandle sendHandle =
                 dev.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT)){

            try {
                // Setup the filter to accept only the arp packets sent back to the current
                // host from the address of the remote host we wanted to get the mac address
                // from.
                receiveHandle.setFilter(
                    "arp and src host " + remoteIpAddress.getHostAddress()
                        + " and dst host " + localIpAddress.getHostAddress(),
                    BpfProgram.BpfCompileMode.OPTIMIZE
                );

                CompletableFuture<MacAddress> resolutionFuture = new CompletableFuture<>();
                PacketListener listener
                    = packet -> {
                    if (packet.contains(ArpPacket.class)) {
                        ArpPacket arp = packet.get(ArpPacket.class);
                        if (arp.getHeader().getOperation().equals(ArpOperation.REPLY)) {
                            resolutionFuture.complete(arp.getHeader().getSrcHardwareAddr());
                        }
                    }
                };

                pool.execute(() -> {
                    try {
                        receiveHandle.loop(-1, listener);
                    } catch (PcapNativeException | NotOpenException e) {
                        logger.error("Error receiving ARP lookup", e);
                    } catch (InterruptedException e) {
                        logger.error("Interrupted! Error receiving ARP lookup", e);
                        Thread.currentThread().interrupt();
                    }
                });

                MacAddress localMacAddress = MacAddress.getByAddress(
                    dev.getLinkLayerAddresses().iterator().next().getAddress());
                ArpPacket.Builder arpBuilder = new ArpPacket.Builder();
                arpBuilder
                    .hardwareType(ArpHardwareType.ETHERNET)
                    .protocolType(EtherType.IPV4)
                    .hardwareAddrLength((byte) MacAddress.SIZE_IN_BYTES)
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
                    return resolutionFuture.get(3000, TimeUnit.MILLISECONDS);
                } catch (TimeoutException e) {
                    logger.info("Couldn't resolve MAC address for ip address {}", remoteIpAddress.getHostAddress(), e);
                    return null;
                }
            } finally {
                if (sendHandle.isOpen()) {
                    sendHandle.close();
                }
                if (receiveHandle.isOpen()) {
                    // Terminate the receive loop first.
                    receiveHandle.breakLoop();
                    receiveHandle.close();
                }
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PcapException("Error looking up MAC address for ip address " +
                remoteIpAddress.getHostAddress() + " on device " + dev.getName(), e);
        } catch (PcapNativeException | ExecutionException | NotOpenException e) {
            throw new PcapException("Error looking up MAC address for ip address " +
                remoteIpAddress.getHostAddress() + " on device " + dev.getName(), e);
        }
    }

    /**
     * Iterate through all devices and find the first that would be able to connect to the given address
     * because it's ip address and subnet mask would allow direct communication.
     *
     * @param remoteAddress address we want to connect to.
     * @return PcapNetworkInterface interface that should be able to connect to the given address.
     * @throws PcapException something went wrong.
     */
    private FirstHop getFirstHop(InetAddress remoteAddress) throws PcapException {
        byte[] remoteIp = remoteAddress.getAddress();

        // Iterate over all network interfaces.
        try {
            // First try if we can connect to the remote device directly.
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                // Iterate over all addresses configured for this interface.
                for (PcapAddress localAddress : dev.getAddresses()) {
                    if(GATEWAY_ONLY_NETMASK.equals(localAddress.getNetmask().getHostAddress())) {
                        return new FirstHop(dev, localAddress.getAddress(),
                            dev.getLinkLayerAddresses().iterator().next(),
                            getMacAddress(dev, localAddress.getAddress(), remoteAddress));
                    }
                    // Only check addresses matching the IP-version of the remote address.
                    if (remoteAddress.getClass().equals(localAddress.getAddress().getClass())) {
                        byte[] localIp = localAddress.getAddress().getAddress();
                        byte[] netMask = localAddress.getNetmask().getAddress();
                        boolean matches = true;
                        // Iterate over all bytes of the address and see if they match
                        // after applying the net mask filter.
                        for (int i = 0; i < localIp.length; i++) {
                            if ((localIp[i] & netMask[i]) != (remoteIp[i] & netMask[i])) {
                                matches = false;
                                break;
                            }
                        }
                        // If the current address would be able to connect to the remote
                        // address, return this device.
                        if (matches) {
                            if (dev.getLinkLayerAddresses().isEmpty()) {
                                continue;
                            }
                            LinkLayerAddress localMacAddress = dev.getLinkLayerAddresses().iterator().next();
                            return new FirstHop(dev, localAddress.getAddress(),
                                localMacAddress,
                                getMacAddress(dev, localAddress.getAddress(), remoteAddress));
                        }
                    }
                }
            }
            return null;
        } catch (PcapNativeException e) {
            throw new PcapException("Error finding a device to communicate with remote address.", e);
        }
    }

    /**
     * Unfortunately there is no way to get the ip address of the default gateway the
     * system uses to route traffic to remote networks. Luckily the 'netstat' command
     * is available on all systems and the only difference is the output. So in this
     * case we fallback to executing this command and parsing it's output depending on
     * the system it is run on.
     *
     * @return InetAddress representing the address of the internet gateway.
     */
    @SuppressWarnings("squid:S1313")
    private InetAddress getDefaultGatewayAddress() {
        try {
            Runtime rt = Runtime.getRuntime();
            String[] commands = {"netstat", "-rn"};
            Process proc = rt.exec(commands);

            BufferedReader stdInput = new BufferedReader(new InputStreamReader(proc.getInputStream()));

            String linePrefix;
            int gatewayColumn;
            if (SystemUtils.IS_OS_WINDOWS) {
                linePrefix = "0.0.0.0";
                gatewayColumn = 2;
            } else if (SystemUtils.IS_OS_MAC_OSX) {
                linePrefix = "default";
                gatewayColumn = 1;
            } else if (SystemUtils.IS_OS_LINUX) {
                linePrefix = "0.0.0.0";
                gatewayColumn = 1;
            } else {
                return null;
            }

            String s;
            while ((s = stdInput.readLine()) != null) {
                if (s.trim().startsWith(linePrefix)) {
                    String[] columns = s.trim().split("\\s+");
                    return InetAddress.getByName(columns[gatewayColumn]);
                }
            }
        } catch (IOException e) {
            logger.debug("error caught", e);
            return null;
        }

        // Command on all platforms: "netstat -rn"

        /* Windows 7: line.trim().startsWith("0.0.0.0")

===========================================================================
Schnittstellenliste
 14...00 1c 42 98 16 7d ......Intel(R) PRO/1000 MT-Netzwerkverbindung #2
 11...00 1c 42 2e f3 40 ......Intel(R) PRO/1000 MT-Netzwerkverbindung
  1...........................Software Loopback Interface 1
 13...00 00 00 00 00 00 00 e0 Microsoft-ISATAP-Adapter
 12...00 00 00 00 00 00 00 e0 Teredo Tunneling Pseudo-Interface
 15...00 00 00 00 00 00 00 e0 Microsoft-ISATAP-Adapter #2
===========================================================================

IPv4-Routentabelle
===========================================================================
Aktive Routen:
     Netzwerkziel    Netzwerkmaske          Gateway    Schnittstelle Metrik
          0.0.0.0          0.0.0.0      10.211.55.1      10.211.55.3     10
      10.211.55.0    255.255.255.0   Auf Verbindung       10.211.55.3    266
      10.211.55.3  255.255.255.255   Auf Verbindung       10.211.55.3    266
    10.211.55.255  255.255.255.255   Auf Verbindung       10.211.55.3    266
        127.0.0.0        255.0.0.0   Auf Verbindung         127.0.0.1    306
        127.0.0.1  255.255.255.255   Auf Verbindung         127.0.0.1    306
  127.255.255.255  255.255.255.255   Auf Verbindung         127.0.0.1    306
      192.168.0.0      255.255.0.0   Auf Verbindung      192.168.0.43    266
     192.168.0.43  255.255.255.255   Auf Verbindung      192.168.0.43    266
  192.168.255.255  255.255.255.255   Auf Verbindung      192.168.0.43    266
        224.0.0.0        240.0.0.0   Auf Verbindung         127.0.0.1    306
        224.0.0.0        240.0.0.0   Auf Verbindung       10.211.55.3    266
        224.0.0.0        240.0.0.0   Auf Verbindung      192.168.0.43    266
  255.255.255.255  255.255.255.255   Auf Verbindung         127.0.0.1    306
  255.255.255.255  255.255.255.255   Auf Verbindung       10.211.55.3    266
  255.255.255.255  255.255.255.255   Auf Verbindung      192.168.0.43    266
===========================================================================
Ständige Routen:
  Keine

IPv6-Routentabelle
===========================================================================
Aktive Routen:
 If Metrik Netzwerkziel             Gateway
 11    266 ::/0                     fe80::21c:42ff:fe00:18
  1    306 ::1/128                  Auf Verbindung
 11     18 fdb2:2c26:f4e4::/64      Auf Verbindung
 11    266 fdb2:2c26:f4e4:0:24b4:9398:1a69:664/128
                                    Auf Verbindung
 11    266 fdb2:2c26:f4e4:0:7147:34d4:e033:a879/128
                                    Auf Verbindung
 14     18 fdb2:2c26:f4e4:1::/64    Auf Verbindung
 14    266 fdb2:2c26:f4e4:1:797d:18a9:3dd6:8105/128
                                    Auf Verbindung
 14    266 fdb2:2c26:f4e4:1:bcd1:eeb5:c8c1:cf05/128
                                    Auf Verbindung
 11    266 fe80::/64                Auf Verbindung
 14    266 fe80::/64                Auf Verbindung
 11    266 fe80::7147:34d4:e033:a879/128
                                    Auf Verbindung
 14    266 fe80::797d:18a9:3dd6:8105/128
                                    Auf Verbindung
  1    306 ff00::/8                 Auf Verbindung
 11    266 ff00::/8                 Auf Verbindung
 14    266 ff00::/8                 Auf Verbindung
===========================================================================
Ständige Routen:
  Keine


        Windows 10

===========================================================================
Schnittstellenliste
 14...00 1c 42 98 16 7d ......Intel(R) PRO/1000 MT-Netzwerkverbindung #2
 11...00 1c 42 2e f3 40 ......Intel(R) PRO/1000 MT-Netzwerkverbindung
  1...........................Software Loopback Interface 1
 13...00 00 00 00 00 00 00 e0 Microsoft-ISATAP-Adapter
 12...00 00 00 00 00 00 00 e0 Teredo Tunneling Pseudo-Interface
 15...00 00 00 00 00 00 00 e0 Microsoft-ISATAP-Adapter #2
===========================================================================

IPv4-Routentabelle
===========================================================================
Aktive Routen:
     Netzwerkziel    Netzwerkmaske          Gateway    Schnittstelle Metrik
          0.0.0.0          0.0.0.0      10.211.55.1      10.211.55.3     10
      10.211.55.0    255.255.255.0   Auf Verbindung       10.211.55.3    266
      10.211.55.3  255.255.255.255   Auf Verbindung       10.211.55.3    266
    10.211.55.255  255.255.255.255   Auf Verbindung       10.211.55.3    266
        127.0.0.0        255.0.0.0   Auf Verbindung         127.0.0.1    306
        127.0.0.1  255.255.255.255   Auf Verbindung         127.0.0.1    306
  127.255.255.255  255.255.255.255   Auf Verbindung         127.0.0.1    306
      192.168.0.0      255.255.0.0   Auf Verbindung      192.168.0.43    266
     192.168.0.43  255.255.255.255   Auf Verbindung      192.168.0.43    266
  192.168.255.255  255.255.255.255   Auf Verbindung      192.168.0.43    266
        224.0.0.0        240.0.0.0   Auf Verbindung         127.0.0.1    306
        224.0.0.0        240.0.0.0   Auf Verbindung       10.211.55.3    266
        224.0.0.0        240.0.0.0   Auf Verbindung      192.168.0.43    266
  255.255.255.255  255.255.255.255   Auf Verbindung         127.0.0.1    306
  255.255.255.255  255.255.255.255   Auf Verbindung       10.211.55.3    266
  255.255.255.255  255.255.255.255   Auf Verbindung      192.168.0.43    266
===========================================================================
Ständige Routen:
  Keine

IPv6-Routentabelle
===========================================================================
Aktive Routen:
 If Metrik Netzwerkziel             Gateway
 11    266 ::/0                     fe80::21c:42ff:fe00:18
  1    306 ::1/128                  Auf Verbindung
 11     18 fdb2:2c26:f4e4::/64      Auf Verbindung
 11    266 fdb2:2c26:f4e4:0:24b4:9398:1a69:664/128
                                    Auf Verbindung
 11    266 fdb2:2c26:f4e4:0:7147:34d4:e033:a879/128
                                    Auf Verbindung
 14     18 fdb2:2c26:f4e4:1::/64    Auf Verbindung
 14    266 fdb2:2c26:f4e4:1:797d:18a9:3dd6:8105/128
                                    Auf Verbindung
 14    266 fdb2:2c26:f4e4:1:bcd1:eeb5:c8c1:cf05/128
                                    Auf Verbindung
 11    266 fe80::/64                Auf Verbindung
 14    266 fe80::/64                Auf Verbindung
 11    266 fe80::7147:34d4:e033:a879/128
                                    Auf Verbindung
 14    266 fe80::797d:18a9:3dd6:8105/128
                                    Auf Verbindung
  1    306 ff00::/8                 Auf Verbindung
 11    266 ff00::/8                 Auf Verbindung
 14    266 ff00::/8                 Auf Verbindung
===========================================================================
Ständige Routen:
  Keine


        Ubuntu: line.trim().startsWith("0.0.0.0")

Kernel IP routing table
Destination     Gateway         Genmask         Flags   MSS Window  irtt Iface
0.0.0.0         10.211.55.1     0.0.0.0         UG        0 0          0 enp0s5
10.211.55.0     0.0.0.0         255.255.255.0   U         0 0          0 enp0s5
169.254.0.0     0.0.0.0         255.255.0.0     U         0 0          0 enp0s5

        Mac: line.trim().startsWith("default")

Routing tables

Internet:
Destination        Gateway            Flags        Refs      Use   Netif Expire
default            10.10.56.1         UGSc          109     1893     en0
10.10.56/24        link#9             UCS             0        0     en0
10.10.56.1/32      link#9             UCS             3        0     en0
10.10.56.1         0:90:7f:a2:7a:a3   UHLWIir        39       16     en0   1188
10.10.56.8/32      link#9             UCS             1        0     en0
10.10.56.8         8c:85:90:18:6f:a9  UHLWI           0        5     lo0
10.37.129/24       link#22            UC              1        0   vnic1
10.211.55/24       link#21            UC              2        0   vnic0
10.211.55.3        0:1c:42:2e:f3:40   UHLWIi          1        0   vnic0   1120
10.211.55.4        0:1c:42:ec:9d:4d   UHLWI           0        0   vnic0   1075
127                127.0.0.1          UCS             0        0     lo0
127.0.0.1          127.0.0.1          UH             26   320774     lo0
169.254            link#9             UCS             1        0     en0
169.254            link#7             UCSI            0        0     en7
169.254            link#20            UCSI            1        0     en8
169.254.233.80/32  link#20            UCS             0        0     en8
169.254.255.255    link#9             UHLSW           6       30     en0
192.168.0/16       link#7             UCS             1        0     en7
192.168.0.99/32    link#7             UCS             1        0     en7
192.168.42.1       link#7             UHLWIi          1       10     en7
224.0.0/4          link#9             UmCS            3        0     en0
224.0.0/4          link#7             UmCSI           1        0     en7
224.0.0/4          link#20            UmCSI           1        0     en8
224.0.0.251        1:0:5e:0:0:fb      UHmLWI          0      106     en0
224.0.0.252        1:0:5e:0:0:fc      UHmLWI          0       32     en0
239.255.255.250    1:0:5e:7f:ff:fa    UHmLWI          0      124     en7
239.255.255.250    1:0:5e:7f:ff:fa    UHmLWI          0      682     en0
239.255.255.250    1:0:5e:7f:ff:fa    UHmLWI          0       64     en8
255.255.255.255/32 link#9             UCS             0        0     en0
255.255.255.255/32 link#7             UCSI            0        0     en7
255.255.255.255/32 link#20            UCSI            0        0     en8

Internet6:
Destination                             Gateway                         Flags         Netif Expire
default                                 fe80::%utun0                    UGcI          utun0
default                                 fe80::%utun1                    UGcI          utun1
::1                                     ::1                             UHL             lo0
fe80::%lo0/64                           fe80::1%lo0                     UcI             lo0
fe80::1%lo0                             link#1                          UHLI            lo0
fe80::%en7/64                           link#7                          UCI             en7
fe80::4a8:61b9:6131:3848%en7            48:65:ee:12:d2:c7               UHLI            lo0
fe80::%en5/64                           link#8                          UCI             en5
fe80::aede:48ff:fe00:1122%en5           ac:de:48:0:11:22                UHLI            lo0
fe80::aede:48ff:fe33:4455%en5           ac:de:48:33:44:55               UHLWIi          en5
fe80::%en0/64                           link#9                          UCI             en0
fe80::10cf:c2ea:7baa:626b%en0           8c:85:90:18:6f:a9               UHLI            lo0
fe80::%awdl0/64                         link#11                         UCI           awdl0
fe80::cbd:62ff:fe3e:406c%awdl0          e:bd:62:3e:40:6c                UHLI            lo0
fe80::%utun0/64                         fe80::3a17:f866:7728:6e8d%utun0 UcI           utun0
fe80::3a17:f866:7728:6e8d%utun0         link#17                         UHLI            lo0
fe80::%utun1/64                         fe80::eb7a:3ecf:562c:167b%utun1 UcI           utun1
fe80::eb7a:3ecf:562c:167b%utun1         link#18                         UHLI            lo0
fe80::%en8/64                           link#20                         UCI             en8
fe80::415:f342:498d:14d2%en8            42:4d:7f:8a:b3:83               UHLI            lo0
ff01::%lo0/32                           ::1                             UmCI            lo0
ff01::%en7/32                           link#7                          UmCI            en7
ff01::%en5/32                           link#8                          UmCI            en5
ff01::%en0/32                           link#9                          UmCI            en0
ff01::%awdl0/32                         link#11                         UmCI          awdl0
ff01::%utun0/32                         fe80::3a17:f866:7728:6e8d%utun0 UmCI          utun0
ff01::%utun1/32                         fe80::eb7a:3ecf:562c:167b%utun1 UmCI          utun1
ff01::%en8/32                           link#20                         UmCI            en8
ff02::%lo0/32                           ::1                             UmCI            lo0
ff02::%en7/32                           link#7                          UmCI            en7
ff02::%en5/32                           link#8                          UmCI            en5
ff02::%en0/32                           link#9                          UmCI            en0
ff02::%awdl0/32                         link#11                         UmCI          awdl0
ff02::%utun0/32                         fe80::3a17:f866:7728:6e8d%utun0 UmCI          utun0
ff02::%utun1/32                         fe80::eb7a:3ecf:562c:167b%utun1 UmCI          utun1
ff02::%en8/32                           link#20                         UmCI            en8


         */
        return null;
    }

    private static class FirstHop {
        private final PcapNetworkInterface networkInterface;
        private final InetAddress localInetAddress;
        private final LinkLayerAddress localMacAddress;
        private final LinkLayerAddress remoteMacAddress;

        private FirstHop(PcapNetworkInterface networkInterface, InetAddress localInetAddress, LinkLayerAddress localMacAddress,
                         LinkLayerAddress remoteMacAddress) {
            this.networkInterface = networkInterface;
            this.localInetAddress = localInetAddress;
            this.localMacAddress = localMacAddress;
            this.remoteMacAddress = remoteMacAddress;
        }
    }

}
