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
package org.apache.plc4x.java.spi.messages.utils;

import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;

public class ResponseItem<T> implements Serializable {

    private final PlcResponseCode code;
    private final T value;

    public ResponseItem(PlcResponseCode code, T value) {
        this.code = code;
        this.value = value;
    }

    public PlcResponseCode getCode() {
        return code;
    }

    public T getValue() {
        return value;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("ResponseItem");
        String codeName = code.name();
        writeBuffer.writeString("result", codeName.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), codeName);
        if (value != null) {
            if (!(value instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Tag value doesn't implement XmlSerializable");
            }
            ((Serializable) value).serialize(writeBuffer);
        }
        writeBuffer.popContext("ResponseItem");
    }

}
