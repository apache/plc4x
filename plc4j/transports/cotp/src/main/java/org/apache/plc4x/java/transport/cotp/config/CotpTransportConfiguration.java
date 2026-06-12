/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.transport.cotp.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;

/**
 * Configuration for COTP (Connection Oriented Transport Protocol) transport.
 * COTP runs on top of TCP and provides connection-oriented communication
 * with TPKT framing as defined in ISO 8073/RFC 1006.
 */
public class CotpTransportConfiguration extends TcpTransportConfiguration implements TransportConfiguration {

    /**
     * Local TSAP (Transport Service Access Point) identifier.
     * Used in COTP connection request. Default is 0x0100.
     */
    @ConfigurationParameter("local-tsap")
    @Description("Local TSAP (Transport Service Access Point) identifier.")
    @IntDefaultValue(0x0311)
    public int localTsap = 0x0311;

    /**
     * Remote TSAP (Transport Service Access Point) identifier.
     * Used in COTP connection request. Default is 0x0102.
     */
    @ConfigurationParameter("remote-tsap")
    @Description("Remote TSAP (Transport Service Access Point) identifier.")
    @IntDefaultValue(0x0100)
    public int remoteTsap = 0x0100;

    /**
     * COTP PDU size for data transmission.
     * Valid values: 128, 256, 512, 1024, 2048, 4096, 8192.
     * Default is 8192 bytes.
     */
    @ConfigurationParameter("cotp-tpdu-size")
    @Description("COTP PDU size for data transmission. Valid values: 128, 256, 512, 1024, 2048, 4096, 8192.")
    @IntDefaultValue(8192)
    public int cotpTpduSize = 8192;

    /**
     * Connection timeout for COTP handshake in milliseconds.
     * Default is 5000ms (5 seconds).
     */
    @ConfigurationParameter("cotp-connection-timeout")
    @Description("Connection timeout for COTP handshake in milliseconds.")
    @IntDefaultValue(5000)
    public int cotpConnectionTimeout = 5000;

    /**
     * COTP protocol class to use.
     * Class 0 is most commonly used (simple class, no flow control).
     * Default is 0.
     */
    @ConfigurationParameter("protocol-class")
    @Description("COTP protocol class to use. Class 0 is most commonly used (simple class, no flow control).")
    @IntDefaultValue(0)
    public int protocolClass = 0;

    public CotpTransportConfiguration() {
        super();
        // COTP typically uses default TCP-settings but you can override
    }

    @Override
    public String toString() {
        return "CotpTransportConfiguration{" +
            "localTsap=0x" + Integer.toHexString(localTsap) +
            ", remoteTsap=0x" + Integer.toHexString(remoteTsap) +
            ", cotpTpduSize=" + cotpTpduSize +
            ", cotpConnectionTimeout=" + cotpConnectionTimeout +
            ", protocolClass=" + protocolClass +
            ", tcp=" + super.toString() +
            '}';
    }
}
