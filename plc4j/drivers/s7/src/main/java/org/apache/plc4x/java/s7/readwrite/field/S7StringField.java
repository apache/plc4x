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
package org.apache.plc4x.java.s7.readwrite.field;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import org.apache.plc4x.java.s7.readwrite.types.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.types.TransportSize;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;

public class S7StringField extends S7Field {

    private final int stringLength;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    protected S7StringField(@JsonProperty("dataType") TransportSize dataType, @JsonProperty("memoryArea") MemoryArea memoryArea,
                            @JsonProperty("blockNumber") int blockNumber, @JsonProperty("byteOffset") int byteOffset,
                            @JsonProperty("bitOffset") byte bitOffset, @JsonProperty("numElements") int numElements,
                            @JsonProperty("stringLength") int stringLength) {
        super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
        this.stringLength = stringLength;
    }

    public int getStringLength() {
        return stringLength;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext(getClass().getSimpleName());

        String memoryArea = getMemoryArea().name();
        writeBuffer.writeString("memoryArea", memoryArea.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), memoryArea);

        writeBuffer.writeUnsignedInt("blockNumber", 16, getBlockNumber());
        writeBuffer.writeUnsignedInt("byteOffset", 16, getByteOffset());
        writeBuffer.writeUnsignedInt("bitOffset", 8, getBitOffset());
        writeBuffer.writeUnsignedInt("numElements", 16, getNumberOfElements());
        writeBuffer.writeInt("stringLength", 64, getStringLength());

        String dataType = getDataType().name();
        writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
