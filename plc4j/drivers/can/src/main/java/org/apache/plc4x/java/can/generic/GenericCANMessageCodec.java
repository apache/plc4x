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
package org.apache.plc4x.java.can.generic;

import org.apache.plc4x.java.can.generic.transport.GenericFrame;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;

import java.util.function.Consumer;

/**
 * MessageCodec for raw CAN frames as the new-SPI CAN transports
 * ({@code can-socketcan}, {@code can-virtualcan}) serialize them onto the
 * transport's byte stream.
 *
 * <p>Each frame is a fixed 16-byte block:
 * <pre>
 *   canId(4 BE) + dlc(1) + padding(3) + data(8, zero-padded)
 * </pre>
 * Standard/extended discrimination uses the top bit of canId (EFF flag, 0x80000000).</p>
 */
public class GenericCANMessageCodec extends MessageCodecBase<GenericFrame> {

    private static final int FRAME_SIZE = 16;
    private static final int EFF_FLAG = 0x80000000;
    private static final int STANDARD_ID_MASK = 0x7FF;
    private static final int EXTENDED_ID_MASK = 0x1FFFFFFF;

    public GenericCANMessageCodec(TransportInstance<?> transportInstance, Consumer<GenericFrame> messageHandler) {
        super("Generic CAN", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return FRAME_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        return FRAME_SIZE;
    }

    @Override
    protected GenericFrame parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        // The transport hands us frames already aligned to FRAME_SIZE; just
        // read the bytes back out and pick apart the id + payload by hand.
        // Using the buffer's underlying byte array avoids a second copy.
        byte[] frame = readBuffer.getBytes();
        int rawId = ((frame[0] & 0xFF) << 24) | ((frame[1] & 0xFF) << 16)
            | ((frame[2] & 0xFF) << 8) | (frame[3] & 0xFF);
        boolean extended = (rawId & EFF_FLAG) != 0;
        int canId = rawId & (extended ? EXTENDED_ID_MASK : STANDARD_ID_MASK);
        int dlc = frame[4] & 0xFF;
        if (dlc < 0 || dlc > 8) {
            dlc = 0;
        }
        byte[] data = new byte[dlc];
        System.arraycopy(frame, 8, data, 0, dlc);
        return new GenericFrame(canId, data);
    }

    /**
     * Serializes a {@link GenericFrame} into the 16-byte transport wire format
     * and writes it directly to the transport. Bypasses the
     * {@link MessageCodecBase#send(org.apache.plc4x.java.spi.buffers.api.Message)}
     * path because {@code GenericFrame} doesn't implement {@code Message}.
     */
    public void sendFrame(GenericFrame frame) throws MessageCodecException {
        byte[] bytes = new byte[FRAME_SIZE];
        int canId = frame.getNodeId();
        bytes[0] = (byte) ((canId >> 24) & 0xFF);
        bytes[1] = (byte) ((canId >> 16) & 0xFF);
        bytes[2] = (byte) ((canId >> 8) & 0xFF);
        bytes[3] = (byte) (canId & 0xFF);
        byte[] data = frame.getData();
        int dlc = Math.min(data.length, 8);
        bytes[4] = (byte) dlc;
        System.arraycopy(data, 0, bytes, 8, dlc);
        try {
            transportInstance.write(bytes);
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to send CAN frame", e);
        }
    }

}
