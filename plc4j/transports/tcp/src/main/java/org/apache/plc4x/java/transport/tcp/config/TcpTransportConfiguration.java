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
package org.apache.plc4x.java.transport.tcp.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;

public class TcpTransportConfiguration implements TransportConfiguration {

    public final static int NO_DEFAULT_PORT = -1;

    /**
     * Connection timeout in milliseconds.
     */
    @ConfigurationParameter("connect-timeout-ms")
    @IntDefaultValue(5000)
    public int connectTimeout;

    /**
     * Socket read timeout in milliseconds. 0 means no timeout.
     */
    @ConfigurationParameter("read-timeout-ms")
    @IntDefaultValue(0)
    public int readTimeout;

    /**
     * Socket write timeout in milliseconds. 0 means no timeout.
     */
    @ConfigurationParameter("write-timeout-ms")
    @IntDefaultValue(0)
    public int writeTimeout;

    /**
     * Enable TCP_NODELAY (disable Nagle's algorithm).
     */
    @ConfigurationParameter("no-delay")
    @BooleanDefaultValue(true)
    public boolean tcpNoDelay;

    /**
     * Enable SO_KEEPALIVE.
     */
    @ConfigurationParameter("keep-alive")
    @BooleanDefaultValue(false)
    public boolean keepAlive;

    /**
     * Send buffer size in bytes. 0 uses system default.
     */
    @ConfigurationParameter("send-buffer-size")
    @IntDefaultValue(81920)
    public int sendBufferSize;

    /**
     * Receive buffer size in bytes. 0 uses system default.
     */
    @ConfigurationParameter("receive-buffer-size")
    @IntDefaultValue(81920)
    public int receiveBufferSize;

    /**
     * Local address to bind to (optional). If not set, uses default.
     */
    @ConfigurationParameter("local-address")
    public String localAddress;

    /**
     * Local port to bind to (optional). 0 uses ephemeral port.
     */
    @ConfigurationParameter("local-port")
    @IntDefaultValue(0)
    public int localPort;

    public int getDefaultPort() {
        return NO_DEFAULT_PORT;
    }

}
