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
package org.apache.plc4x.java.transport.virtualcan;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class VirtualCANFrame implements Message {
    private final int nodeId;
    private final byte[] data;

    public VirtualCANFrame(int nodeId, byte[] data) {
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
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeUnsignedShort("length", 8, (short) data.length);
        writeBuffer.writeUnsignedInt("nodeId", 32, nodeId);
        writeBuffer.writeByteArray("data", data);
    }

    @Override
    public int getLengthInBytes() {
        return 0;
    }

    @Override
    public int getLengthInBits() {
        return 0;
    }

    public String toString() {
        return "VirtualCANFrame " + nodeId + "[" + data.length + "]" + Hex.encodeHexString(data);
    }
}
