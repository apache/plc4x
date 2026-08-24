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
package org.apache.plc4x.java.transport.rawsocket.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.Required;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;

public class RawSocketTransportConfiguration implements TransportConfiguration {

    /**
     * Network interface name to use (e.g., "eth0", "en0", "\\Device\\NPF_{GUID}" on Windows)
     * If not specified, will use the first available interface.
     */
    @ConfigurationParameter( "interface-name")
    @Description("Network interface name to use (e.g., \"eth0\", \"en0\", \"\\\\Device\\\\NPF_{GUID}\" on Windows). If not specified, will use the first available interface.")
    public String interfaceName;

    /**
     * Local MAC address to send from. If not set, uses interface's MAC address.
     */
    @ConfigurationParameter( "local-address")
    @Description( "Local MAC address to send from. If not set, uses interface's MAC address.")
    public String localAddress;

    /**
     * Remote MAC address to send to.
     */
    @ConfigurationParameter( "remote-address")
    @Description( "Remote MAC address to send to.")
    @Required
    public String remoteAddress;

    /**
     * EtherType / Protocol number (e.g., 0x88B5 for PROFINET, 0x88CC for LLDP, custom values)
     */
    @ConfigurationParameter( "protocol-id")
    @Description( "EtherType / Protocol number (e.g., 0x88B5 for PROFINET, 0x88CC for LLDP, custom values)")
    @Required
    public int protocolId;

    /**
     * Promiscuous mode - capture all packets on the network, not just those destined for this interface.
     */
    @ConfigurationParameter( "promiscuous-mode")
    @Description( "Promiscuous mode - capture all packets on the network, not just those destined for this interface.")
    @BooleanDefaultValue(false)
    public boolean promiscuousMode;

    /**
     * Packet capture timeout in milliseconds.
     */
    @ConfigurationParameter( "capture-timeout")
    @Description( "Packet capture timeout in milliseconds.")
    @IntDefaultValue(1000)
    public int captureTimeout;

    /**
     * Snapshot length - maximum bytes to capture per packet. 0 means capture entire packet.
     */
    @ConfigurationParameter( "snapshot-length")
    @Description( "Snapshot length - maximum bytes to capture per packet. 0 means capture entire packet.")
    @IntDefaultValue(65536)
    public int snapshotLength;

    /**
     * Buffer size for packet capture in bytes.
     */
    @ConfigurationParameter( "buffer-size")
    @Description( "Buffer size for packet capture in bytes.")
    @IntDefaultValue(1048576)
    public int bufferSize;

    /**
     * How many bytes of captured frames may wait for a consumer.
     *
     * <p>A raw socket accepts whatever matches its filter, and nothing obliges anything to drain
     * what it captures. Counted in bytes rather than frames because frames are not one size and
     * memory is what runs out. Frames arriving past this are dropped whole, and the count of what
     * was dropped is logged.</p>
     */
    @ConfigurationParameter("receive-queue-size")
    @Description("Bytes of captured frames that may wait for a consumer before frames are dropped.")
    @IntDefaultValue(1048576)
    public int receiveQueueSize;

    /**
     * VLAN ID (0 means no VLAN tag). Range: 0-4095
     */
    @ConfigurationParameter( "vlan-id")
    @Description( "VLAN ID (0 means no VLAN tag). Range: 0-4095")
    @IntDefaultValue(0)
    public int vlanId;

    /**
     * VLAN priority (0-7). Only used if vlanId > 0.
     */
    @ConfigurationParameter( "vlan-priority")
    @Description( "VLAN priority (0-7). Only used if vlanId > 0.")
    @IntDefaultValue(0)
    public int vlanPriority;

    /**
     * Reuse the underlying network interface across multiple transport instances.
     * When true, instances with the same interface and protocol will share a pcap handle.
     * This is useful for protocols where multiple logical connections share one Ethernet type.
     */
    @ConfigurationParameter( "reuse-interface")
    @Description( "Reuse the underlying network interface across multiple transport instances. When true, instances with the same interface and protocol will share a pcap handle. This is useful for protocols where multiple logical connections share one Ethernet type.")
    @BooleanDefaultValue(false)
    public boolean reuseInterface;

    /**
     * BPF (Berkeley Packet Filter) expression to filter packets.
     * If not specified, will filter by EtherType and MAC addresses automatically.
     */
    @ConfigurationParameter( "bpf-filter")
    @Description( "BPF (Berkeley Packet Filter) expression to filter packets.")
    public String bpfFilter;

    /**
     * Maximum frame size (MTU) in bytes.
     */
    @ConfigurationParameter( "max-frame-size")
    @Description( "Maximum frame size (MTU) in bytes.")
    @IntDefaultValue(1500)
    public int maxFrameSize;

    /**
     * Read timeout for blocking reads in milliseconds. 0 means no timeout.
     */
    @ConfigurationParameter( "read-timeout")
    @Description( "Read timeout for blocking reads in milliseconds.")
    @IntDefaultValue(0)
    public int readTimeout;

    /**
     * When {@code true} the transport delivers full Ethernet frames (including the
     * destination/source MAC + EtherType header) to the driver, and expects the
     * driver to write full Ethernet frames on send (the transport will not wrap
     * payload in an Ethernet header). Required for L2 protocols (PROFINET, EtherCAT)
     * that build their own Ethernet frames and need visibility into source MAC for
     * routing.
     */
    @ConfigurationParameter("include-ethernet-header")
    @Description("Deliver full Ethernet frames to the driver and accept raw Ethernet frames on send. Required for L2 protocols that build their own Ethernet headers.")
    @BooleanDefaultValue(false)
    public boolean includeEthernetHeader;
}
