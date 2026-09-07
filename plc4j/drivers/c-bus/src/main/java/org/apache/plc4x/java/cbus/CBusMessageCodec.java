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
package org.apache.plc4x.java.cbus;

import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.cbus.readwrite.CBusOptions;
import org.apache.plc4x.java.cbus.readwrite.utils.StaticHelper;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;

import java.util.function.Consumer;

/**
 * MessageCodec for the Clipsal C-Bus wire format.
 *
 * <p>C-Bus over TCP is line-based: each command is terminated by {@code \r}
 * with an optional trailing {@code \n}. Total message size is determined by
 * scanning the available bytes for the terminator — same logic as the legacy
 * {@code CBusDriver.ByteLengthEstimator}.</p>
 */
public class CBusMessageCodec extends MessageCodecBase<CBusCommand> {

    private static final int MIN_HEADER = 1;

    private final CBusOptions cBusOptions;

    public CBusMessageCodec(TransportInstance<?> transportInstance,
                            Consumer<CBusCommand> messageHandler,
                            CBusOptions cBusOptions) {
        super("C-Bus", transportInstance, messageHandler);
        this.cBusOptions = cBusOptions;
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MIN_HEADER;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        // The 1-byte header alone isn't enough to find the CR/LF boundary —
        // peek the entire currently-available payload and scan for terminators.
        try {
            byte[] all = transportInstance.peekReadableBytes(availableBytes);
            for (int i = 0; i < all.length; i++) {
                boolean isCR = all[i] == '\r';
                boolean hasMore = i + 1 < all.length;
                if (isCR) {
                    if (hasMore && all[i + 1] == '\n') {
                        return i + 2;
                    }
                    if (hasMore) {
                        return i + 1;
                    }
                    // CR at end of buffer: wait for the optional LF.
                    return -1;
                }
            }
            return -1;
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to peek incoming C-Bus bytes", e);
        }
    }

    @Override
    protected CBusCommand parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return CBusCommand.staticParse(readBuffer, cBusOptions);
    }

    @Override
    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data, StaticHelper.OPTIONS);
    }

    @Override
    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size], StaticHelper.OPTIONS);
    }

}
