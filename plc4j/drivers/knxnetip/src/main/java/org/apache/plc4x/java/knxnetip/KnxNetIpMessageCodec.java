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
package org.apache.plc4x.java.knxnetip;

import org.apache.plc4x.java.knxnetip.readwrite.KnxNetIpMessage;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for the KNXNet/IP wire format.
 *
 * <p>Every KNXNet/IP frame starts with a 6-byte common header:
 * <pre>
 *   +--------+--------+--------+--------+--------+--------+
 *   | 0x06   | 0x10   | service type    | total length    |
 *   +--------+--------+--------+--------+--------+--------+
 * </pre>
 * The {@code total length} field at offset 4 (big-endian, unsigned 16-bit)
 * counts the entire frame including the header — so it doubles as the
 * frame size for codec framing. UDP datagrams are already message-bounded,
 * but reading the length out of the header lets us share one codec with any
 * future stream-oriented transport without changes.</p>
 */
public class KnxNetIpMessageCodec extends MessageCodecBase<KnxNetIpMessage> {

    private static final int HEADER_SIZE = 6;
    private static final int LENGTH_OFFSET = 4;

    /**
     * Buffer options for KNXNet/IP — big-endian, standard binary encodings
     * for numeric fields, UTF-8 for embedded strings. The generated codec
     * classes require explicit encodings on the buffer, otherwise reads of
     * unsigned/signed integers throw {@link BufferException}.
     */
    public static final WithOption[] BUFFER_OPTIONS = {
        WithByteBasedOption.WithByteOrder("BIG_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    public KnxNetIpMessageCodec(TransportInstance<?> transportInstance, Consumer<KnxNetIpMessage> messageHandler) {
        super("KNXNet/IP", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        if ((header[0] & 0xFF) != 0x06 || (header[1] & 0xFF) != 0x10) {
            throw new MessageCodecException(String.format(
                "Invalid KNXNet/IP frame header: 0x%02X%02X (expected 0x0610)",
                header[0] & 0xFF, header[1] & 0xFF));
        }
        return ((header[LENGTH_OFFSET] & 0xFF) << 8) | (header[LENGTH_OFFSET + 1] & 0xFF);
    }

    @Override
    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data, BUFFER_OPTIONS);
    }

    @Override
    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size], BUFFER_OPTIONS);
    }

    @Override
    protected KnxNetIpMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return KnxNetIpMessage.staticParse(readBuffer);
    }

}
