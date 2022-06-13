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
package org.apache.plc4x.java.spi.values;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcByteArray extends PlcIECValue<byte[]> {

    public static PlcByteArray of(Object value) {
        if (value instanceof byte[]) {
            return new PlcByteArray((byte[]) value);
        }
        throw new IllegalArgumentException("Only byte[] supported here");
    }

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcByteArray(@JsonProperty("value")byte[] value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return Hex.encodeHexString(value);
    }

    @JsonIgnore
    public byte[] getBytes() {
        return value;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeByteArray(getClass().getSimpleName(), value);
    }

}
