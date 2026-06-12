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
package org.apache.plc4x.java.transport.serial.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;

public class SerialTransportConfiguration implements TransportConfiguration {

    /**
     * Baud rate (bits per second)
     */
    @ConfigurationParameter( "baud-rate")
    @Description( "Baud rate (bits per second)")
    @IntDefaultValue(9600)
    public int baudRate;

    /**
     * Number of data bits (5, 6, 7, or 8)
     */
    @ConfigurationParameter( "data-bits")
    @Description( "Number of data bits (5, 6, 7, or 8)")
    @IntDefaultValue(8)
    public int dataBits;

    /**
     * Number of stop bits (1 or 2)
     */
    @ConfigurationParameter( "stop-bits")
    @Description( "Number of stop bits (1 or 2)")
    @IntDefaultValue(1)
    public int stopBits;

    /**
     * Parity: NONE, ODD, EVEN, MARK, SPACE
     */
    @ConfigurationParameter( "parity")
    @Description( "Parity: NONE, ODD, EVEN, MARK, SPACE")
    @StringDefaultValue("NONE")
    public String parity;

    /**
     * Flow control: NONE, RTS_CTS, XON_XOFF, RTS_CTS_XON_XOFF
     */
    @ConfigurationParameter( "flow-control")
    @Description( "Flow control: NONE, RTS_CTS, XON_XOFF, RTS_CTS_XON_XOFF")
    @StringDefaultValue("NONE")
    public String flowControl;

    /**
     * Read timeout in milliseconds. 0 means blocking read.
     */
    @ConfigurationParameter( "read-timeout")
    @Description( "Read timeout in milliseconds. 0 means blocking read.")
    @IntDefaultValue(1000)
    public int readTimeout;

    /**
     * Write timeout in milliseconds.
     */
    @ConfigurationParameter( "write-timeout")
    @Description( "Write timeout in milliseconds.")
    @IntDefaultValue(1000)
    public int writeTimeout;

    /**
     * Enable DTR (Data Terminal Ready) signal
     */
    @ConfigurationParameter("dtr")
    @Description( "Enable DTR (Data Terminal Ready) signal")
    @BooleanDefaultValue(false)
    public boolean dtr;

    /**
     * Enable RTS (Request To Send) signal
     */
    @ConfigurationParameter("rts")
    @Description( "Enable RTS (Request To Send) signal")
    @BooleanDefaultValue(false)
    public boolean rts;

    /**
     * Reuse the underlying serial port across multiple transport instances.
     * When true, instances with the same port will share a connection.
     * This is useful for protocols where multiple logical connections share one serial port.
     * WARNING: Only use this when the protocol supports multiplexing (like Modbus RTU with different slave IDs).
     */
    @ConfigurationParameter( "reuse-port")
    @Description( "Reuse the underlying serial port across multiple transport instances. When true, instances with the same port will share a connection. This is useful for protocols where multiple logical connections share one serial port.")
    @BooleanDefaultValue(false)
    public boolean reusePort;

    /**
     * Receive buffer size in bytes. 0 uses system default.
     */
    @ConfigurationParameter( "receive-buffer-size")
    @Description( "Receive buffer size in bytes. 0 uses system default.")
    @IntDefaultValue(4096)
    public int receiveBufferSize;

    /**
     * Send buffer size in bytes. 0 uses system default.
     */
    @ConfigurationParameter( "send-buffer-size")
    @Description( "Send buffer size in bytes. 0 uses system default.")
    @IntDefaultValue(4096)
    public int sendBufferSize;

    /**
     * Enable break signal
     */
    @ConfigurationParameter( "break-enabled")
    @Description( "Enable break signal")
    @BooleanDefaultValue(false)
    public boolean breakEnabled;

    /**
     * Interframe delay in milliseconds for protocols that need spacing between messages.
     * Only applies when reusePort is true.
     */
    @ConfigurationParameter( "interframe-delay")
    @Description( "Interframe delay in milliseconds for protocols that need spacing between messages. Only applies when reusePort is true.")
    @IntDefaultValue(0)
    public int interframeDelay;
}
