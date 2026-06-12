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
package org.apache.plc4x.java.openprotocol;

import org.apache.plc4x.java.openprotocol.readwrite.OpenProtocolMessage;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for the Open-Protocol wire format.
 *
 * <p>Every frame starts with an ASCII-encoded 4-digit length field followed
 * by the rest of the message. The length field itself counts the bytes
 * <em>after</em> it minus one (matching the legacy {@code ByteLengthEstimator}
 * in {@code OpenProtocolDriver}), so the total frame length is
 * {@code parsedLength + 6}: 4 length digits + the trailing '\0' byte plus a
 * one-byte body offset adjustment baked into the spec.</p>
 */
public class OpenProtocolMessageCodec extends MessageCodecBase<OpenProtocolMessage> {

    /** Number of header bytes needed to determine the total frame size. */
    private static final int HEADER_SIZE = 6;
    /** Offset of the 2-byte ASCII length suffix inside the header. */
    private static final int LENGTH_OFFSET = 4;
    private static final int LENGTH_OVERHEAD = 6;

    /**
     * The wire format is ASCII-encoded numerics throughout — every uint field
     * is read as digits parsed from ASCII bytes. The generated codec's
     * pushContext for {@code OpenProtocolMessage} already sets this for the
     * top-level message, but we apply the same options at buffer construction
     * so peeks via the codec helpers see the same context.
     */
    static final WithOption[] BUFFER_OPTIONS = {
        WithOption.WithEncoding("ASCII")
    };

    /** Protocol-fixed revision used when parsing incoming messages. */
    private final int revision;

    public OpenProtocolMessageCodec(TransportInstance<?> transportInstance,
                                    Consumer<OpenProtocolMessage> messageHandler,
                                    int revision) {
        super("Open-Protocol", transportInstance, messageHandler);
        this.revision = revision;
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        // The 2 bytes at offset 4 carry the low-order digits of the ASCII-
        // encoded length field. Matches the legacy ByteLengthEstimator,
        // which read an unsigned 16-bit value at the same offset and added 6.
        int low = ((header[LENGTH_OFFSET] & 0xFF) << 8) | (header[LENGTH_OFFSET + 1] & 0xFF);
        return low + LENGTH_OVERHEAD;
    }

    @Override
    protected OpenProtocolMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return OpenProtocolMessage.staticParse(readBuffer, revision);
    }

}
