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
package org.apache.plc4x.java.s7.readwrite.tag;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;

public class S7StringTag extends S7Tag {

    private final int stringLength;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    protected S7StringTag(@JsonProperty("dataType") TransportSize dataType, @JsonProperty("memoryArea") MemoryArea memoryArea,
                          @JsonProperty("blockNumber") int blockNumber, @JsonProperty("byteOffset") int byteOffset,
                          @JsonProperty("bitOffset") byte bitOffset, @JsonProperty("numElements") int numElements,
                          @JsonProperty("stringLength") int stringLength,
                          @JsonProperty("stringEncoding") String stringEncoding) {
        super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements,stringEncoding);
        this.stringLength = stringLength;
    }

    public int getStringLength() {
        return stringLength;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        String memoryArea = getMemoryArea().name();
        writeBuffer.writeString("memoryArea",
            memoryArea.getBytes(StandardCharsets.UTF_8).length * 8,
            memoryArea, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.writeUnsignedInt("blockNumber", 16, getBlockNumber());
        writeBuffer.writeUnsignedInt("byteOffset", 16, getByteOffset());
        writeBuffer.writeUnsignedInt("bitOffset", 8, getBitOffset());
        writeBuffer.writeUnsignedInt("numElements", 16, getNumberOfElements());
        writeBuffer.writeInt("stringLength", 64, getStringLength());

        String dataType = getDataType().name();
        writeBuffer.writeString("dataType",
            dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
