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
package org.apache.plc4x.java.transport.virtualcan.io;

import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.transport.virtualcan.VirtualCANFrame;

/**
 * A manual IO type which writes basic structures to pipeline:
 * -   8 bits: length
 * -  32 bits: nodeId
 * - variable: data
 */
public class VirtualCANFrameIO implements MessageInput<VirtualCANFrame> {

    public static final MessageInput<? extends Message> INSTANCE = new VirtualCANFrameIO();

    @Override
    public VirtualCANFrame parse(ReadBuffer io, Object... args) throws ParseException {
        short length = io.readUnsignedShort("length", 8);
        int nodeId = io.readUnsignedInt("nodeId", 32);
        byte[] data = io.readByteArray("data", length);

        return new VirtualCANFrame(nodeId, data);
    }

    /*@Override
    public void serialize(WriteBuffer io, VirtualCANFrame value, Object... args) throws ParseException {
        io.writeUnsignedShort("length", 8, (short) value.getData().length);
        io.writeUnsignedInt("nodeId", 32, value.getNodeId());
        io.writeByteArray("data", value.getData());
    }*/

}
