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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.s7.readwrite.*;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for S7 PDUs as carried inside COTP-DT payloads. The COTP transport already
 * unwraps TPKT/COTP framing, so the codec sees raw S7 message bytes starting with 0x32.
 * <p>
 * S7 header layout (10 bytes): magic (1) | type (1) | reserved (2) | tpduRef (2) | paramLen (2)
 * | dataLen (2). For response/ack-data (types 2 and 3) two extra bytes (errorClass/errorCode)
 * follow the header before parameters/data start, so the total is 12 + paramLen + dataLen.
 */
public class S7CotpMessageCodec extends MessageCodecBase<S7Message> {

    private static final int S7_MIN_HEADER_SIZE = 10;
    private static final byte S7_MAGIC = 0x32;
    private static final int S7_MSG_TYPE_ACK = 2;
    private static final int S7_MSG_TYPE_ACK_DATA = 3;

    public S7CotpMessageCodec(TransportInstance<?> transportInstance, Consumer<S7Message> messageHandler) {
        super("S7", transportInstance, messageHandler);
    }

    public void close() {
        try {
            transportInstance.close();
        } catch (Exception ignored) {
            // ignore
        }
    }

    @Override
    protected int getMinimumHeaderSize() {
        return S7_MIN_HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        // No S7 message starts with anything else, and no byte arriving later changes the one
        // already here, so say so and let the byte be dropped.
        if (header[0] != S7_MAGIC) {
            return DESYNCHRONIZED;
        }
        int msgType = header[1] & 0xFF;
        int paramLen = ((header[6] & 0xFF) << 8) | (header[7] & 0xFF);
        int dataLen = ((header[8] & 0xFF) << 8) | (header[9] & 0xFF);
        int extra = (msgType == S7_MSG_TYPE_ACK || msgType == S7_MSG_TYPE_ACK_DATA) ? 2 : 0;
        return S7_MIN_HEADER_SIZE + extra + paramLen + dataLen;
    }

    @Override
    protected S7Message parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return S7Message.staticParse(readBuffer);
    }

    /**
     * S7 PDUs are big-endian; the generated serializers/parsers ask the buffer for an
     * unsigned-binary integer encoding and a signed twos-complement encoding for fields
     * that don't pass an explicit option. The default {@link MessageCodecBase} buffers do
     * not pre-configure any of those, which trips field writers like {@code writeUnsignedShort}.
     */
    @Override
    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
    }

    @Override
    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
    }

}
