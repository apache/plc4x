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
package org.apache.plc4x.java.can.generic.transport;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

/**
 * Wrapper for one CAN frame's address + payload.
 *
 * <p>Implements {@link Message} so it fits the {@code MessageCodecBase<M extends Message>}
 * generic; serialization is the same 16-byte wire layout the SocketCAN / VirtualCAN
 * transports already speak: canId(4 BE) + dlc(1) + padding(3) + data(8).</p>
 */
public class GenericFrame implements Message {

    private static final int FRAME_SIZE = 16;

    private final int nodeId;
    private final byte[] data;

    public GenericFrame(int nodeId, byte[] data) {
        this.nodeId = nodeId;
        this.data = data;
    }

    public int getNodeId() {
        return nodeId;
    }

    public byte[] getData() {
        return data;
    }

    @Override
    public int getLengthInBytes() {
        return FRAME_SIZE;
    }

    @Override
    public int getLengthInBits() {
        return FRAME_SIZE * 8;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        // writeUnsignedInt only supports 1..31 bits; a full 32-bit CAN ID has
        // to go through writeUnsignedLong.
        writeBuffer.writeUnsignedLong(32, nodeId);
        int dlc = Math.min(data.length, 8);
        writeBuffer.writeUnsignedInt(8, dlc);
        writeBuffer.writeUnsignedInt(8, 0);
        writeBuffer.writeUnsignedInt(8, 0);
        writeBuffer.writeUnsignedInt(8, 0);
        for (int i = 0; i < 8; i++) {
            writeBuffer.writeUnsignedInt(8, i < dlc ? (data[i] & 0xFF) : 0);
        }
    }

}
