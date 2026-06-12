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
package org.apache.plc4x.java.plc4x;

import org.apache.plc4x.java.plc4x.readwrite.Plc4xMessage;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * Frames the PLC4X proxy protocol:
 * <pre>
 *   byte 0     = version (always 0x01)
 *   byte 1..2  = total packet length (big endian, includes header)
 *   byte 3..4  = request id
 *   byte 5     = request type
 *   ...        = type-specific payload
 * </pre>
 * Header is 3 bytes (version + length); the rest of the message is consumed
 * by the generated parser.
 */
public class Plc4xMessageCodec extends MessageCodecBase<Plc4xMessage> {

    private static final int HEADER_SIZE = 3;
    private static final int VERSION = 0x01;

    public Plc4xMessageCodec(TransportInstance<?> transportInstance, Consumer<Plc4xMessage> messageHandler) {
        super("PLC4X-Proxy", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        if ((header[0] & 0xFF) != VERSION) {
            throw new MessageCodecException("Unexpected PLC4X proxy protocol version: 0x"
                + Integer.toHexString(header[0] & 0xFF));
        }
        return ((header[1] & 0xFF) << 8) | (header[2] & 0xFF);
    }

    @Override
    protected Plc4xMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return Plc4xMessage.staticParse(readBuffer);
    }

}
