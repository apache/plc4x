package org.apache.plc4x.java.utils.rawsockets.netty;

import org.pcap4j.core.*;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.*;
import org.pcap4j.util.ByteArrays;
import org.pcap4j.util.MacAddress;

import java.net.Inet4Address;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.*;
import java.util.concurrent.*;

public class RawSocket {

    private static final int SNAPLEN = 65536;
    private static final int READ_TIMEOUT = 10;

    private static final Map<InetAddress, MacAddress>  arpCache = new HashMap<>();

    // The id of the protocol we will be communicating in.
    private final int protocolNumber;

    private PcapNetworkInterface nif;
    private InetAddress remoteIpAddress;
    private MacAddress remoteMacAddress;
    private InetAddress localIpAddress;
    private MacAddress localMacAddress;
    private ExecutorService pool = Executors.newSingleThreadExecutor();
    private PcapHandle receiveHandle;

    private final List<RawSocketListener> listeners = new LinkedList<>();

    public RawSocket(int protocolNumber) {
        this.protocolNumber = protocolNumber;
    }

    public void connect(String remoteAddress) throws RawSocketException {
        try {
            pool = Executors.newScheduledThreadPool(2);

            remoteIpAddress = InetAddress.getByName(remoteAddress);

            localIpAddress = getLocalNetworkAddress(remoteIpAddress);
            if (localIpAddress == null) {
                throw new RawSocketException("Unable to connect to the remote address " + remoteAddress);
            }
            nif = Pcaps.getDevByAddress(localIpAddress);
            localMacAddress = MacAddress.getByAddress(nif.getLinkLayerAddresses().get(0).getAddress());

            remoteMacAddress = getMacAddress(remoteIpAddress);

            // Setup receiving of packets and redirecting them to the corresponding listeners.
            // Filter packets to contain only the ip protocol number of the current protocol.
            receiveHandle = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);
            // Set the filter.
            String filterString = "ip protochain " + protocolNumber +
                " and ether dst " + localMacAddress.toString() +
                " and ether src " + remoteMacAddress.toString();
            receiveHandle.setFilter(filterString, BpfProgram.BpfCompileMode.OPTIMIZE);
            PacketListener packetListener
                = packet -> {
                for(RawSocketListener listener : listeners) {
                    listener.packetReceived(packet.getRawData());
                }
            };
            pool.execute(() -> {
                try {
                    receiveHandle.loop(-1, packetListener);
                } catch (PcapNativeException | InterruptedException | NotOpenException e) {
                    // TODO: This is not nice ... fix this ...
                    throw new RuntimeException("Error receiving ARP lookup", e);
                }
            });
        } catch (PcapNativeException | NotOpenException | UnknownHostException e) {
            throw new RawSocketException("Error setting up RawSocket", e);
        }
    }

    public void disconnect() throws RawSocketException {
        // TODO: Terminate all the listeners and the thread pool.
    }

    public void write(byte[] rawData) throws RawSocketException {
        PcapHandle sendHandle = null;
        try {
            sendHandle = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);
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
            etherBuilder.dstAddr(remoteMacAddress)
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
            System.out.println(p.toString());
            sendHandle.sendPacket(p);
        } catch (PcapNativeException | NotOpenException e) {
            throw new RawSocketException("Error sending packet.", e);
        } finally {
            if((sendHandle != null) && sendHandle.isOpen()) {
                sendHandle.close();
            }
        }
    }

    public void addListener(RawSocketListener listener) {
        listeners.add(listener);
    }

    public void removeListener(RawSocketListener listener) {
        listeners.remove(listener);
    }

    protected MacAddress getMacAddress(InetAddress address) throws RawSocketException {
        if(!arpCache.containsKey(address)) {
            MacAddress macAddress = lookupMacAddress(address);
            arpCache.put(address, macAddress);
            return macAddress;
        }
        return arpCache.get(address);
    }

    protected MacAddress lookupMacAddress(InetAddress remoteIpAddress) throws RawSocketException {
        try {
            PcapHandle receiveHandle
                = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);
            PcapHandle sendHandle
                = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);

            // Setup the filter to accept only the arp packets sent back to the current
            // host from the address of the remote host we wanted to get the mac address
            // from.
            receiveHandle.setFilter(
                "arp and src host " + remoteIpAddress.getHostAddress()
                    + " and dst host " + localIpAddress.getHostAddress()
                    + " and ether dst " + localMacAddress.toString(),
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
                } catch (PcapNativeException | InterruptedException | NotOpenException e) {
                    // TODO: This is not nice ... fix this ...
                    throw new RuntimeException("Error receiving ARP lookup", e);
                }
            });

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
                return resolutionFuture.get(10000, TimeUnit.MILLISECONDS);
            } catch (TimeoutException e) {
                return null;
            }
        } catch (PcapNativeException | InterruptedException | ExecutionException | NotOpenException e) {
            throw new RawSocketException("Error looking up mac address.", e);
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

    /**
     * Iterate through all devices and find the first that would be able to connect to the given address.
     *
     * @param remoteAddress address we want to connect to.
     * @return PcapNetworkInterface interface that should be able to connect to the given address.
     * @throws Exception something went wrong.
     */
    protected InetAddress getLocalNetworkAddress(InetAddress remoteAddress) throws RawSocketException {
        byte[] remoteIp = remoteAddress.getAddress();
        // Iterate over all network interfaces.
        try {
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                System.out.println(dev.getName());
                // Iterate over all addresses configured for this interface.
                for (PcapAddress addr : dev.getAddresses()) {
                    // Only check addresses matching the IP-version of the remote address.
                    if (remoteAddress.getClass().equals(addr.getAddress().getClass())) {
                        byte[] localIp = addr.getAddress().getAddress();
                        byte[] netMask = addr.getNetmask().getAddress();
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
                            return addr.getAddress();
                        }
                    }
                }
            }
        } catch (PcapNativeException e) {
            throw new RawSocketException("Error finding a device to communicate with remote address.", e);
        }
        return null;
    }

}
