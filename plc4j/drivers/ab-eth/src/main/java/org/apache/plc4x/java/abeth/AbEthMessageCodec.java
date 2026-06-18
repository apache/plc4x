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
package org.apache.plc4x.java.abeth;

import org.apache.plc4x.java.abeth.readwrite.CIPEncapsulationPacket;
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
 * MessageCodec for the Allen Bradley AB-ETH CIP encapsulation wire format.
 *
 * <p>Every packet starts with a 28-byte CIP encapsulation header. The
 * {@code length} field at offset 2 (little-endian, unsigned 16-bit)
 * counts only the payload <em>after</em> the header — so the total frame
 * length is {@code length + 28}. That matches the legacy
 * {@code AbEthDriver.ByteLengthEstimator}.</p>
 */
public class AbEthMessageCodec extends MessageCodecBase<CIPEncapsulationPacket> {

    private static final int HEADER_SIZE = 4;
    private static final int LENGTH_OFFSET = 2;
    private static final int HEADER_OVERHEAD = 28;

    /**
     * AB-ETH wire format is little-endian. Reads of unsigned integers and the
     * occasional signed/float field need explicit encodings on the buffer.
     */
    static final WithOption[] BUFFER_OPTIONS = {
        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    public AbEthMessageCodec(TransportInstance<?> transportInstance, Consumer<CIPEncapsulationPacket> messageHandler) {
        super("AB-ETH", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        // Big-endian unsigned 16-bit at offset 2
        int payloadLength = ((header[LENGTH_OFFSET] & 0xFF) << 8) | (header[LENGTH_OFFSET + 1] & 0xFF);
        return payloadLength + HEADER_OVERHEAD;
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
    protected CIPEncapsulationPacket parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return CIPEncapsulationPacket.staticParse(readBuffer);
    }

}
