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

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.device.ProfinetChannel;
import org.apache.plc4x.java.profinet.device.ProfinetDeviceMessageHandler;
import org.apache.plc4x.java.profinet.device.ProfinetSubscriptionHandle;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.io.IOException;
import java.net.DatagramSocket;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;

public class ProfinetDriverContext implements DriverContext {

    public static final int DEFAULT_UDP_PORT = 34964;
    private Map<Long, ProfinetSubscriptionHandle> subscriptions = new HashMap<>();
    private ProfinetDeviceMessageHandler handler = new ProfinetDeviceMessageHandler();
    private ProfinetConfiguration configuration;
    private DatagramSocket socket;
    private ProfinetChannel channel;
    private final Map<String, ProfinetISO15745Profile> gsdFiles = new HashMap<>();

    public ProfinetChannel getChannel() {
        return channel;
    }

    public void setChannel(ProfinetChannel channel) {
        this.channel = channel;
    }

    public Map<Long, ProfinetSubscriptionHandle> getSubscriptions() {
        return subscriptions;
    }

    public void setSubscriptions(Map<Long, ProfinetSubscriptionHandle> subscriptions) {
        this.subscriptions = subscriptions;
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

    public Map<String, ProfinetISO15745Profile> getGsdFiles() {
        return gsdFiles;
    }

}
