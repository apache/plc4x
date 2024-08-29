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
package org.apache.plc4x.java.transport.virtualcan;

import static org.apache.plc4x.java.spi.codegen.fields.FieldWriterFactory.*;
import static org.apache.plc4x.java.spi.codegen.fields.FieldReaderFactory.*;
import static org.apache.plc4x.java.spi.codegen.io.DataWriterFactory.*;

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.MessageInput;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class VirtualCANFrame implements Message {

    public final static MessageInput<VirtualCANFrame> PARSER = new MessageInput<VirtualCANFrame>() {

        @Override
        public VirtualCANFrame parse(ReadBuffer io) throws ParseException {
            WithOption withOption = WithOption.WithByteOrder(ByteOrder.LITTLE_ENDIAN);

            short length = io.readUnsignedShort("length", 8, withOption);
            int nodeId = io.readInt("nodeId", 32, withOption);
            byte[] data = io.readByteArray("data", length, withOption);

            return new VirtualCANFrame(nodeId, data);
        }
    };

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
        WithOption withOption = WithOption.WithByteOrder(ByteOrder.LITTLE_ENDIAN);

        writeSimpleField("length", (short) data.length, writeUnsignedShort(writeBuffer, 8), withOption);
        writeSimpleField("nodeId", nodeId, writeSignedInt(writeBuffer, 32), withOption);
        writeByteArrayField("data", data, writeByteArray(writeBuffer, 8), withOption);
    }


    @Override
    public int getLengthInBytes() {
        return 1 + 4 + data.length;
    }

    @Override
    public int getLengthInBits() {
        return getLengthInBytes() * 8;
    }

    public String toString() {
        return "VirtualCANFrame " + nodeId + "[" + data.length + "]" + Hex.encodeHexString(data);
    }
}
