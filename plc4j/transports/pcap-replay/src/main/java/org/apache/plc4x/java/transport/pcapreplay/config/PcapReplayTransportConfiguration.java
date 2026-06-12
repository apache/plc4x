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
package org.apache.plc4x.java.transport.pcapreplay.config;


import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.Required;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.pcapreplay.PcapFilePlayer;

public class PcapReplayTransportConfiguration implements TransportConfiguration {

    /**
     * Path to the PCAP file to replay.
     */
    @ConfigurationParameter("pcap-file")
    @Description("Path to the PCAP file to replay.")
    @Required
    public String pcapFile;

    /**
     * Local MAC address to simulate (acts as this device).
     * If not set, uses the first MAC found in the PCAP.
     */
    @ConfigurationParameter("local-address")
    @Description( "Local MAC address to simulate (acts as this device).")
    public String localAddress;

    /**
     * Remote MAC address to communicate with.
     * If not set, uses the first different MAC found in the PCAP.
     */
    @ConfigurationParameter("remote-address")
    @Description( "Remote MAC address to communicate with.")
    public String remoteAddress;

    /**
     * EtherType / Protocol number to filter (e.g., 0x88B5 for PROFINET).
     * If 0, accepts all EtherTypes.
     */
    @ConfigurationParameter( "protocol-id")
    @Description( "EtherType / Protocol number to filter (e.g., 0x88B5 for PROFINET).")
    @IntDefaultValue(0)
    public int protocolId;

    /**
     * Replay speed multiplier. 1.0 = real-time, 0.0 = as fast as possible, 2.0 = double speed.
     */
    @ConfigurationParameter( "speed-factor")
    @Description( "Replay speed multiplier. 1.0 = real-time, 0.0 = as fast as possible, 2.0 = double speed.")
    @IntDefaultValue(0)
    public double speedFactor;

    /**
     * Loop the PCAP file when it reaches the end.
     */
    @ConfigurationParameter( "loop")
    @Description( "Loop the PCAP file when it reaches the end.")
    @BooleanDefaultValue(false)
    public boolean loop;

    /**
     * Only replay packets from remote to local (simulates device responses).
     * When false, replays all packets.
     */
    @ConfigurationParameter( "only-incoming-packets")
    @Description("Only replay packets from remote to local (simulates device responses). When false, replays all packets.")
    @BooleanDefaultValue(true)
    public boolean onlyIncomingPackets;

    /**
     * Only replay packets from local to remote (simulates device requests).
     * When false, replays all packets.
     */
    @ConfigurationParameter( "only-outgoing-packets")
    @Description("Only replay packets from local to remote (simulates device requests).")
    @BooleanDefaultValue(false)
    public boolean onlyOutgoingPackets;

    /**
     * Maximum frame size in bytes.
     */
    @ConfigurationParameter( "max-frame-size")
    @Description("Maximum frame size in bytes.")
    @IntDefaultValue(1500)
    public int maxFrameSize;

    /**
     * Read timeout for blocking reads in milliseconds. 0 means no timeout.
     */
    @ConfigurationParameter("read-timeout")
    @Description("Read timeout for blocking reads in milliseconds.")
    @IntDefaultValue(0)
    public int readTimeout;

    /**
     * Buffer size for packet queue.
     */
    @ConfigurationParameter( "packet-queue-size")
    @Description("Buffer size for packet queue.")
    @IntDefaultValue(1000)
    public int packetQueueSize;

    /**
     * Auto-detect MAC addresses from PCAP file based on first packet.
     */
    @ConfigurationParameter( "auto-detect-mac-addresses")
    @Description("Auto-detect MAC addresses from PCAP file based on first packet.")
    @BooleanDefaultValue(true)
    public boolean autoDetectMacAddresses;

    /**
     * Start replay automatically when transport is created.
     */
    @ConfigurationParameter( "auto-start")
    @Description("Start replay automatically when transport is created.")
    @BooleanDefaultValue(true)
    public boolean autoStart;

    /**
     * Filter by VLAN ID (0 means no VLAN filtering).
     */
    @ConfigurationParameter( "vlan-id")
    @Description("Filter by VLAN ID (0 means no VLAN filtering).")
    @IntDefaultValue(0)
    public int vlanId;

    public PcapFilePlayer mockPlayer;

}
