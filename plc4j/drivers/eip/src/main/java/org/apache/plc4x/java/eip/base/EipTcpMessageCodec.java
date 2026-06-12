/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.eip.readwrite.EipPacket;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * MessageCodec for EtherNet/IP (CIP-encapsulated) packets.
 * The EIP encapsulation header is fixed at 24 bytes; the next 2 bytes after the
 * 2-byte command code contain the length of the body.
 */
public class EipTcpMessageCodec extends MessageCodecBase<EipPacket> {

    private static final int EIP_HEADER_SIZE = 24;

    private final boolean bigEndian;

    public EipTcpMessageCodec(TransportInstance<?> transportInstance, Consumer<EipPacket> messageHandler, boolean bigEndian) {
        super("EthernetIP", transportInstance, messageHandler);
        this.bigEndian = bigEndian;
    }

    @Override
    protected int getMinimumHeaderSize() {
        return EIP_HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        int packetLength;
        if (bigEndian) {
            packetLength = ((header[2] & 0xFF) << 8) | (header[3] & 0xFF);
        } else {
            packetLength = ((header[3] & 0xFF) << 8) | (header[2] & 0xFF);
        }
        return EIP_HEADER_SIZE + packetLength;
    }

    @Override
    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data,
            WithByteBasedOption.WithByteOrder(bigEndian ? "BIG_ENDIAN" : "LITTLE_ENDIAN"),
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"));
    }

    @Override
    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size],
            WithByteBasedOption.WithByteOrder(bigEndian ? "BIG_ENDIAN" : "LITTLE_ENDIAN"),
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"));
    }

    @Override
    protected EipPacket parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return EipPacket.staticParse(readBuffer, true);
    }

}
