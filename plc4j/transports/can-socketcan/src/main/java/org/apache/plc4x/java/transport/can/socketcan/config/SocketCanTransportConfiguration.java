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
package org.apache.plc4x.java.transport.can.socketcan.config;


import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.Required;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.can.config.CanTransportConfiguration;

/**
 * Configuration for the SocketCAN transport.
 * <p>
 * Extends the base CAN transport configuration with SocketCAN-specific parameters
 * for interface selection, resource sharing, and timeout control.
 */
public class SocketCanTransportConfiguration extends CanTransportConfiguration {

    /**
     * The Linux CAN interface name to use (e.g., "can0", "vcan0").
     * This is a required parameter — the transport cannot be created without it.
     */
    @ConfigurationParameter("interface-name")
    @Description("Linux CAN interface name (e.g., \"can0\", \"vcan0\")")
    @Required
    public String interfaceName;

    /**
     * Whether to share the underlying CAN socket across multiple transport instances.
     * When true, instances using the same interface name will share a single CAN socket
     * via the {@link org.apache.plc4x.java.transport.can.socketcan.SharedCanManager}.
     */
    @ConfigurationParameter("reuse-interface")
    @Description("Share CAN socket across multiple transport instances on the same interface")
    @BooleanDefaultValue(false)
    public boolean reuseInterface;

    /**
     * Read timeout in milliseconds for blocking reads on the CAN socket.
     */
    @ConfigurationParameter("read-timeout-ms")
    @Description("Read timeout in milliseconds for blocking reads on the CAN socket")
    @IntDefaultValue(1000)
    public int readTimeout = 1000;
}
