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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.slmp.readwrite.SlmpMessage;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * Frames {@link SlmpMessage} over a TCP byte stream using the 3E length field.
 * The fixed 3E response header is 9 bytes; the 2-byte responseDataLength is
 * little-endian at offset 7 (NB: modbus's MBAP length is big-endian — do not
 * copy that read here).
 */
public class SlmpMessageCodec extends MessageCodecBase<SlmpMessage> {

    private static final int SLMP_3E_HEADER_SIZE = 9; // subHeader(2) + accessRoute(5) + length(2)

    public SlmpMessageCodec(TransportInstance<?> transportInstance, Consumer<SlmpMessage> messageHandler) {
        super("SLMP", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return SLMP_3E_HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // responseDataLength: little-endian uint16 at byte offset 7..8
        int length = (header[7] & 0xFF) | ((header[8] & 0xFF) << 8);
        return SLMP_3E_HEADER_SIZE + length;
    }

    @Override
    protected SlmpMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return SlmpMessage.staticParse(readBuffer);
    }
}
