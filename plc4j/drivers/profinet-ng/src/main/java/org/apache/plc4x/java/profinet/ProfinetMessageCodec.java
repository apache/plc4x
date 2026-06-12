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
package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.profinet.readwrite.Ethernet_Frame;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * PROFINET runs directly over Ethernet (EtherType 0x8892 for PnDcp, plus IPv4/UDP
 * for DCE-RPC). The new SPI raw-socket transport, when configured with
 * {@code include-ethernet-header=true}, delivers each captured Ethernet frame as
 * one self-contained byte block and accepts pre-built frames on write. So our
 * codec is effectively a 1:1 wrapper around {@code Ethernet_Frame.staticParse}.
 *
 * <p>Each packet arrives whole — there's no streaming framing to worry about.
 * We expose {@link #getMinimumHeaderSize()} as a small value just so the base
 * class lets us through; once we have a packet we treat it as a complete frame.</p>
 */
public class ProfinetMessageCodec extends MessageCodecBase<Ethernet_Frame> {

    public ProfinetMessageCodec(TransportInstance<?> transportInstance, Consumer<Ethernet_Frame> messageHandler) {
        super("PROFINET", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        // 14 = dst(6) + src(6) + ethertype(2); enough to decide a frame exists.
        return 14;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // Raw socket transport hands us whole packets, so consume everything
        // currently available — that is the full frame.
        return Math.max(availableBytes, header.length);
    }

    @Override
    protected Ethernet_Frame parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return Ethernet_Frame.staticParse(readBuffer);
    }

}
