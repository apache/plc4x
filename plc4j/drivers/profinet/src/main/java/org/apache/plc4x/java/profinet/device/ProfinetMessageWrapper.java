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

import io.netty.channel.Channel;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.profinet.readwrite.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.utils.rawsockets.netty.RawSocketChannel;


import java.io.IOException;
import java.net.DatagramPacket;
import java.util.Random;

public class ProfinetMessageWrapper {

    public static void sendUdpMessage(ProfinetCallable<DceRpc_Packet> callable, ProfinetDevice context) throws RuntimeException {
        try {
            DceRpc_Packet packet = callable.create();
            Random rand = new Random();
            // Serialize it to a byte-payload
            Ethernet_FramePayload_IPv4 udpFrame = new Ethernet_FramePayload_IPv4(
                rand.nextInt(65356),
                true,
                false,
                (short) 64,
                new IpAddress(context.getLocalIpAddress().getAddress()),
                new IpAddress(context.getIpAddress().getAddress()),
                50000,
                context.getPort(),
                packet
            );
            MacAddress srcAddress = context.getLocalMacAddress();
            MacAddress dstAddress = context.getMacAddress();
            Ethernet_Frame frame = new Ethernet_Frame(
                dstAddress,
                srcAddress,
                udpFrame);

            context.getChannel().send(frame, callable);

        } catch (IOException e) {
            throw new RuntimeException(e);
        } catch (PlcException e) {
            throw new RuntimeException(e);
        }
    }
}
