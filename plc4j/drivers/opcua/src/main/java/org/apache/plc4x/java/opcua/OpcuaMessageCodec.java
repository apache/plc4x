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
package org.apache.plc4x.java.opcua;

import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for OPC UA Binary TCP framing.
 *
 * <p>Every OPC UA TCP chunk starts with a fixed 8-byte header:
 * <ul>
 *   <li>3 bytes — message type ASCII ({@code HEL}, {@code ACK}, {@code ERR},
 *       {@code MSG}, {@code OPN}, {@code CLO}, {@code RHE})</li>
 *   <li>1 byte — chunk flag ({@code F}=final, {@code C}=intermediate,
 *       {@code A}=abort)</li>
 *   <li>4 bytes — little-endian uint32 total chunk size including this header</li>
 * </ul>
 * The {@link OpcuaAPU} parser handles everything from byte 0 of the header onwards.
 */
public class OpcuaMessageCodec extends MessageCodecBase<OpcuaAPU> {

    private static final int HEADER_SIZE = 8;

    public OpcuaMessageCodec(TransportInstance<?> transportInstance, Consumer<OpcuaAPU> messageHandler) {
        super("OPC UA", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        return ((header[4] & 0xFF))
            | ((header[5] & 0xFF) << 8)
            | ((header[6] & 0xFF) << 16)
            | ((header[7] & 0xFF) << 24);
    }

    @Override
    protected OpcuaAPU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return OpcuaAPU.staticParse(readBuffer, true, true);
    }

    @Override
    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data, PayloadConverter.LITTLE_ENDIAN);
    }

    @Override
    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size], PayloadConverter.LITTLE_ENDIAN);
    }
}
