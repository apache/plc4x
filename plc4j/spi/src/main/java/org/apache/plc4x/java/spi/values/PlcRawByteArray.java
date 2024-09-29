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

import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

public class PlcRawByteArray extends PlcIECValue<byte[]> {

    public static PlcRawByteArray of(Object value) {
        if (value instanceof PlcRawByteArray) {
            return (PlcRawByteArray) value;
        } else if (value instanceof byte[]) {
            return new PlcRawByteArray((byte[]) value);
        }
        throw new IllegalArgumentException("Only byte[] supported here");
    }

    public PlcRawByteArray(byte[] value) {
        this.value = value;
        this.isNullable = false;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.RAW_BYTE_ARRAY;
    }

    @Override
    public byte[] getRaw() {
        return value;
    }
    
    @Override
    public String toString() {
        return Hex.encodeHexString(value);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.writeByteArray(getClass().getSimpleName(), value);
    }

}
