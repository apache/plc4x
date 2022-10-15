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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.profinet.protocol.ProfinetProtocolLogic;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.DatagramPacket;
import java.net.InetAddress;
import java.util.HashMap;

public class ProfinetUdpMessageHandler {

    private final Logger logger = LoggerFactory.getLogger(ProfinetUdpMessageHandler.class);

    private HashMap<String, ProfinetDevice> configuredDevices;

    public void handle(DatagramPacket packet) {
        InetAddress address = packet.getAddress();
        String ss = address.getHostName();
        logger.debug(address.getHostName());
    }

    public void setConfiguredDevices(HashMap<String, ProfinetDevice> configuredDevices) {
        this.configuredDevices = configuredDevices;
    }
}
