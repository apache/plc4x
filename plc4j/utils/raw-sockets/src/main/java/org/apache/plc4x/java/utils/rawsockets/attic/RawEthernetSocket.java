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

import org.apache.plc4x.java.utils.pcap.netty.exception.PcapException;
import org.pcap4j.core.*;
import org.pcap4j.packet.*;
import org.pcap4j.packet.namednumber.*;
import org.pcap4j.util.LinkLayerAddress;
import org.pcap4j.util.MacAddress;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.*;

/**
 * The raw ethernet socket relies on a layer 3 IP protocol implementation for finding the
 * Network device able to connect to a given mac address. In Contrast to the {@link RawIpSocket}
 * the Layer 2 protocol does not support gateways and routers, so we don't need the logic
 * for looking up the default gateway, which makes this implementation a lot simpler.
 */
public class RawEthernetSocket {

    private static final Logger logger = LoggerFactory.getLogger(RawEthernetSocket.class);

    private static final int SNAPLEN = 65536;
    private static final int READ_TIMEOUT = 10;

    // The EtherType of the protocol we will be communicating in.
    private final EtherType etherType;

    private PcapNetworkInterface nif;
    private MacAddress remoteMacAddress;
    private MacAddress localMacAddress;
    private ExecutorService pool = Executors.newSingleThreadExecutor();
    private PcapHandle receiveHandle;

    private final List<RawSocketListener> listeners = new LinkedList<>();

    public RawEthernetSocket(EtherType etherType) {
        this.etherType = etherType;
    }

    public void connect(String localMacAddress, String remoteMacAddress) throws PcapException {
        try {
            pool = Executors.newScheduledThreadPool(2);

            this.localMacAddress = MacAddress.getByName(localMacAddress);
            this.remoteMacAddress = MacAddress.getByName(remoteMacAddress);
            // Find out which network device is able to connect to the given mac address.
            nif = null;
            for (PcapNetworkInterface dev : Pcaps.findAllDevs()) {
                if(!dev.isLoopBack()) {
                    for (LinkLayerAddress macAddress : dev.getLinkLayerAddresses()) {
                        if(Arrays.equals(macAddress.getAddress(), this.localMacAddress.getAddress())) {
                            nif = dev;
                        }
                    }
                }
            }
            if(nif == null) {
                throw new PcapException(
                    "Unable to find local network device with mac address " + remoteMacAddress);
            }

            // Setup receiving of packets and redirecting them to the corresponding listeners.
            // Filter packets to contain only the ip protocol number of the current protocol.
            receiveHandle = nif.openLive(SNAPLEN, PcapNetworkInterface.PromiscuousMode.PROMISCUOUS, READ_TIMEOUT);

            // Set the filter.
            String filterString = "ether proto " + etherType.valueAsString() +
                " and ether src " + this.remoteMacAddress.toString() +
                " and ether dst " + this.localMacAddress.toString();

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
                        etherType.valueAsString(), remoteMacAddress, e);
                } catch (PcapNativeException | NotOpenException e) {
                    logger.error("Error receiving packet for protocol {} from MAC address {}",
                        etherType.valueAsString(), remoteMacAddress, e);
                }
            });
        } catch (PcapNativeException | NotOpenException e) {
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

            EthernetPacket.Builder etherBuilder = new EthernetPacket.Builder();
            etherBuilder.dstAddr(remoteMacAddress)
                .srcAddr(localMacAddress)
                .type(etherType)
                .paddingAtBuild(true);
            etherBuilder.payloadBuilder(
                new AbstractPacket.AbstractBuilder() {
                    @Override
                    public Packet build() {
                        return packetBuilder.build();
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

}
