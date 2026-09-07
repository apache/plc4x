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
package org.apache.plc4x.java.canopen;

import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;

import java.util.function.Consumer;

/**
 * Decodes raw SocketCAN-style 16-byte frames into {@link CANOpenFrame}s and back.
 *
 * <p>Wire layout (matches what the {@code can-socketcan} / {@code can-virtualcan}
 * transports deliver):
 * <pre>
 *   canId(4 BE) + dlc(1) + padding(3) + data(8, zero-padded)
 * </pre>
 *
 * <p>The 11-bit CAN identifier is split by CANopen into a 4-bit
 * {@link CANOpenService} (top nibble) and a 7-bit node id (lower bits) —
 * we discover the service by finding the {@code CANOpenService} whose
 * {@code [min,max]} range contains the id, then compute
 * {@code nodeId = canId - service.min}.</p>
 */
public class CANOpenMessageCodec extends MessageCodecBase<CANOpenFrame> {

    private static final int FRAME_SIZE = 16;
    private static final int EFF_FLAG = 0x80000000;
    private static final int STANDARD_ID_MASK = 0x7FF;
    private static final int EXTENDED_ID_MASK = 0x1FFFFFFF;

    public CANOpenMessageCodec(TransportInstance<?> transportInstance, Consumer<CANOpenFrame> messageHandler) {
        super("CANopen", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return FRAME_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        return FRAME_SIZE;
    }

    @Override
    protected CANOpenFrame parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        byte[] frame = readBuffer.getBytes();
        int rawId = ((frame[0] & 0xFF) << 24) | ((frame[1] & 0xFF) << 16)
            | ((frame[2] & 0xFF) << 8) | (frame[3] & 0xFF);
        boolean extended = (rawId & EFF_FLAG) != 0;
        int canId = rawId & (extended ? EXTENDED_ID_MASK : STANDARD_ID_MASK);
        int dlc = frame[4] & 0xFF;
        if (dlc < 0 || dlc > 8) {
            dlc = 0;
        }
        CANOpenService service = resolveService(canId);
        if (service == null) {
            return null;
        }
        int nodeId = canId - service.getMin();
        byte[] data = new byte[dlc];
        System.arraycopy(frame, 8, data, 0, dlc);
        // Parse just the payload portion. The frame's nodeId / service come
        // from the CAN identifier, not the payload bytes, so we can't use the
        // generated CANOpenFrame.staticParse directly.
        ReadBufferByteBased payloadBuffer = new ReadBufferByteBased(data, CANOpenConnection.LITTLE_ENDIAN_OPTIONS);
        CANOpenPayload payload = CANOpenPayload.staticParse(payloadBuffer, service);
        return new CANOpenFrame((short) nodeId, service, payload);
    }

    public void sendFrame(CANOpenFrame frame) throws MessageCodecException {
        byte[] bytes = new byte[FRAME_SIZE];
        int canId = frame.getService().getMin() + frame.getNodeId();
        bytes[0] = (byte) ((canId >> 24) & 0xFF);
        bytes[1] = (byte) ((canId >> 16) & 0xFF);
        bytes[2] = (byte) ((canId >> 8) & 0xFF);
        bytes[3] = (byte) (canId & 0xFF);
        try {
            CANOpenPayload payload = frame.getPayload();
            WriteBufferByteBased payloadBuffer = new WriteBufferByteBased(
                new byte[payload.getLengthInBytes()], CANOpenConnection.LITTLE_ENDIAN_OPTIONS);
            payload.serialize(payloadBuffer);
            byte[] payloadBytes = payloadBuffer.getBytes();
            int dlc = Math.min(payloadBytes.length, 8);
            bytes[4] = (byte) dlc;
            System.arraycopy(payloadBytes, 0, bytes, 8, dlc);
            transportInstance.write(bytes);
        } catch (BufferException | TransportException e) {
            throw new MessageCodecException("Failed to send CANopen frame", e);
        }
    }

    private static CANOpenService resolveService(int canId) {
        for (CANOpenService service : CANOpenService.values()) {
            if (canId >= service.getMin() && canId <= service.getMax()) {
                return service;
            }
        }
        return null;
    }

}
