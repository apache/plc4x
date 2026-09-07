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
package org.apache.plc4x.java.profinet.packets;

import org.apache.plc4x.java.profinet.readwrite.Ethernet_Frame;
import org.apache.plc4x.java.profinet.readwrite.IpAddress;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.function.Predicate;

/**
 * Minimal request/response abstraction the {@link PnDcpPacketFactory} uses to
 * send PROFINET DCE-RPC requests and wait for a matching reply. Mirrors the
 * canopen {@code CANConversation} pattern but tracks PROFINET-specific endpoint
 * data the packet factory needs to build complete Ethernet/IPv4/UDP frames.
 */
public interface ProfinetConversation {

    MacAddress getLocalMacAddress();

    MacAddress getRemoteMacAddress();

    IpAddress getLocalIp();

    IpAddress getRemoteIp();

    /** Sends a frame on the bus without waiting for a response. */
    void sendToWire(Ethernet_Frame frame);

    /**
     * Registers a one-shot listener that completes with the next incoming frame
     * matching {@code predicate}, or fails with {@link java.util.concurrent.TimeoutException}
     * after {@code timeout}.
     */
    CompletableFuture<Ethernet_Frame> expect(Predicate<Ethernet_Frame> predicate, Duration timeout);

}
