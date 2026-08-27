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
package org.apache.plc4x.java.transport.udp.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;

public class UdpTransportConfiguration implements TransportConfiguration {

    public final static int NO_DEFAULT_PORT = -1;

    /**
     * Local address to bind to. If not set, binds to all interfaces.
     */
    @ConfigurationParameter( "local-address")
    @Description( "Local address to bind to. If not set, binds to all interfaces.")
    public String localAddress;

    /**
     * Local port to bind to. 0 uses ephemeral port.
     * When reuseSocket is true, this should be set to share the same socket.
     */
    @ConfigurationParameter( "local-port")
    @Description( "Local port to bind to. 0 uses ephemeral port.")
    @IntDefaultValue(0)
    public int localPort;

    /**
     * Socket read timeout in milliseconds. 0 means no timeout.
     */
    @ConfigurationParameter( "read-timeout-ms")
    @Description( "Socket read timeout in milliseconds. 0 means no timeout.")
    @IntDefaultValue(0)
    public int readTimeout;

    /**
     * Maximum UDP packet size in bytes.
     */
    @ConfigurationParameter( "max-packet-size")
    @Description( "Maximum UDP packet size in bytes.")
    @IntDefaultValue(65507)
    public int maxPacketSize;

    /**
     * Send buffer size in bytes. 0 uses system default.
     */
    @ConfigurationParameter( "send-buffer-size")
    @Description( "Send buffer size in bytes. 0 uses system default.")
    @IntDefaultValue(0)
    public int sendBufferSize;

    /**
     * Receive buffer size in bytes. 0 uses system default.
     */
    @ConfigurationParameter( "receive-buffer-size")
    @Description( "Receive buffer size in bytes. 0 uses system default.")
    @IntDefaultValue(0)
    public int receiveBufferSize;

    /**
     * Enable SO_BROADCAST for sending broadcast packets.
     */
    @ConfigurationParameter( "broadcast")
    @Description( "Enable SO_BROADCAST for sending broadcast packets.")
    @BooleanDefaultValue(false)
    public boolean broadcast;

    /**
     * Enable SO_REUSEADDR to allow multiple bindings to the same address/port.
     */
    @ConfigurationParameter( "reuse-address")
    @Description( "Enable SO_REUSEADDR to allow multiple bindings to the same address/port.")
    @BooleanDefaultValue(false)
    public boolean reuseAddress;

    /**
     * Share the underlying UDP socket across multiple transport instances.
     * When true, instances with the same localAddress:localPort will share a socket.
     * This is useful for protocols where multiple logical connections share one UDP port.
     */
    @ConfigurationParameter( "share-socket")
    @Description( "Share the underlying UDP socket across multiple transport instances. When true, instances with the same localAddress:localPort will share a socket. This is useful for protocols where multiple logical connections share one UDP port.")
    @BooleanDefaultValue(false)
    public boolean shareSocket;

    /**
     * Time-to-live for multicast packets (1-255).
     */
    @ConfigurationParameter( "multicast-ttl")
    @Description( "Time-to-live for multicast packets (1-255).")
    @IntDefaultValue(1)
    public int multicastTtl;

    public int getDefaultPort() {
        return NO_DEFAULT_PORT;
    }

}
