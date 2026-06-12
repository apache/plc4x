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

package org.apache.plc4x.java.profinet.context;

import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.device.ProfinetDeviceMessageHandler;
import java.net.DatagramSocket;

/**
 * Per-connection state. Used to be wired up as the old SPI's
 * {@code DriverContext}; with the new SPI the connection just owns one of
 * these directly.
 */
public class ProfinetDriverContext {

    public static final int DEFAULT_UDP_PORT = 34964;
    private ProfinetDeviceMessageHandler handler;
    private ProfinetConfiguration configuration;
    private DatagramSocket socket;
    private ProfinetChannel channel;

    public ProfinetChannel getChannel() {
        return channel;
    }

    public void setChannel(ProfinetChannel channel) {
        this.channel = channel;
    }

    public ProfinetDeviceMessageHandler getHandler() {
        return handler;
    }

    public void setHandler(ProfinetDeviceMessageHandler handler) {
        this.handler = handler;
    }

    public ProfinetConfiguration getConfiguration() {
        return configuration;
    }

    public void setConfiguration(ProfinetConfiguration configuration) {
        this.configuration = configuration;
    }

    public DatagramSocket getSocket() {
        return socket;
    }

    public void setSocket(DatagramSocket socket) {
        this.socket = socket;
    }

}
