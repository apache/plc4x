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
     * Default local TSAP used when no explicit {@code local-tsap} is configured.
     */
    public static final int DEFAULT_LOCAL_TSAP = 0x0311;

    /**
     * Default remote TSAP used when no explicit {@code remote-tsap} is configured.
     */
    public static final int DEFAULT_REMOTE_TSAP = 0x0100;

    /**
     * Raw local TSAP (Transport Service Access Point) override. A value of {@code 0} means
     * "not set". Always read the effective value via {@link #getLocalTsap()} rather than this
     * field directly: subclasses may derive the TSAP from other parameters (e.g. rack/slot).
     */
    @ConfigurationParameter("local-tsap")
    @Description("Local TSAP (Transport Service Access Point) identifier.")
    @IntDefaultValue(0)
    public int localTsap = 0;

    /**
     * Raw remote TSAP (Transport Service Access Point) override. A value of {@code 0} means
     * "not set". Always read the effective value via {@link #getRemoteTsap()} rather than this
     * field directly: subclasses may derive the TSAP from other parameters (e.g. rack/slot).
     */
    @ConfigurationParameter("remote-tsap")
    @Description("Remote TSAP (Transport Service Access Point) identifier.")
    @IntDefaultValue(0)
    public int remoteTsap = 0;

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

    /**
     * @return the effective local TSAP: the explicit {@code local-tsap} override when set
     * (non-zero), otherwise {@link #DEFAULT_LOCAL_TSAP}. Subclasses may override to derive it.
     */
    public int getLocalTsap() {
        return localTsap != 0 ? localTsap : DEFAULT_LOCAL_TSAP;
    }

    /**
     * @return the effective remote TSAP: the explicit {@code remote-tsap} override when set
     * (non-zero), otherwise {@link #DEFAULT_REMOTE_TSAP}. Subclasses may override to derive it.
     */
    public int getRemoteTsap() {
        return remoteTsap != 0 ? remoteTsap : DEFAULT_REMOTE_TSAP;
    }

    @Override
    public String toString() {
        return "CotpTransportConfiguration{" +
            "localTsap=0x" + Integer.toHexString(getLocalTsap()) +
            ", remoteTsap=0x" + Integer.toHexString(getRemoteTsap()) +
            ", cotpTpduSize=" + cotpTpduSize +
            ", cotpConnectionTimeout=" + cotpConnectionTimeout +
            ", protocolClass=" + protocolClass +
            ", tcp=" + super.toString() +
            '}';
    }
}
