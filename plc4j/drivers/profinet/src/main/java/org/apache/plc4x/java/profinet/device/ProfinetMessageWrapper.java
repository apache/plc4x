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

import org.apache.plc4x.java.profinet.context.ProfinetDeviceContext;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.InetAddress;
import java.util.Random;

public class ProfinetMessageWrapper implements MessageWrapper {

    private static final Logger logger = LoggerFactory.getLogger(ProfinetMessageWrapper.class);

    public void sendUdpMessage(ProfinetCallable<DceRpc_Packet> callable, ProfinetDeviceContext context) throws RuntimeException {
        try {
            DceRpc_Packet packet = callable.create();
            Random rand = new Random();
            // Serialize it to a byte-payload
            Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
                rand.nextInt(65536),
                true,
                false,
                (short) 64,
                new IpAddress(context.getNetworkInterface().getIpAddressAsByteArray()),
                new IpAddress(InetAddress.getByName(context.getIpAddress()).getAddress()),
                context.getSourcePort(),
                context.getDestinationPort(),
                packet
            );
            MacAddress srcAddress = context.getLocalMacAddress();
            MacAddress dstAddress = context.getMacAddress();
            Ethernet_Frame frame = new Ethernet_Frame(
                dstAddress,
                srcAddress,
                udpFrame);

            context.getChannel().send(frame);

        } catch (IOException e) {
            logger.error("Unable to send UDP message");
            context.setState(ProfinetDeviceState.ABORT);
        }
    }

    public void sendPnioMessage(ProfinetCallable<Ethernet_Frame> callable, ProfinetDeviceContext context) throws RuntimeException {
        Ethernet_Frame packet = callable.create();
        context.getChannel().send(packet);
    }
}
