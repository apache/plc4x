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
package org.apache.plc4x.java.ads;

import org.apache.plc4x.java.ads.readwrite.AmsTCPPacket;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for the ADS protocol over TCP.
 * The AMS/TCP header is 6 bytes: 2 reserved bytes + 4-byte little-endian length of the userdata that follows.
 */
public class AdsTcpMessageCodec extends MessageCodecBase<AmsTCPPacket> {

    private static final int AMS_TCP_HEADER_SIZE = 6;

    /**
     * The AMS/TCP length field is 32 bits wide, but an AMS packet carries at most a 16-bit data
     * length plus its headers, so anything approaching the width of the field is not a length this
     * protocol can produce. Kept generous rather than exact: the point is to tell a plausible
     * length from a number that could only wedge us waiting for bytes that will never come.
     */
    private static final long AMS_TCP_MAX_LENGTH = 16 * 1024 * 1024;

    public AdsTcpMessageCodec(TransportInstance<?> transportInstance, Consumer<AmsTCPPacket> messageHandler) {
        super("ADS", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return AMS_TCP_HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // Length field is bytes 2-5, little endian (32-bit unsigned).
        long length = ((long) (header[2] & 0xFF))
            | (((long) (header[3] & 0xFF)) << 8)
            | (((long) (header[4] & 0xFF)) << 16)
            | (((long) (header[5] & 0xFF)) << 24);
        // Decide in long arithmetic and only then narrow. Casting first turned a length near the
        // top of the field into a small or negative int, and a length merely implausible - tens of
        // megabytes - left us waiting for bytes the receive buffer could never hold.
        if (length > AMS_TCP_MAX_LENGTH) {
            return DESYNCHRONIZED;
        }
        return (int) (AMS_TCP_HEADER_SIZE + length);
    }

    @Override
    protected AmsTCPPacket parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return AmsTCPPacket.staticParse(readBuffer);
    }

}
