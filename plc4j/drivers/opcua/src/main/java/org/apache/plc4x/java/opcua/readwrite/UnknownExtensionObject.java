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
package org.apache.plc4x.java.opcua.readwrite;

import java.util.Arrays;
import java.util.Objects;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.generation.WriteBufferBoxBased;

/**
 * Represents an ExtensionObject whose type is not known to the parser.
 * Stores the raw body bytes so they can be interpreted later based on context
 * (e.g., the tag's dataType hint).
 */
public class UnknownExtensionObject extends ExtensionObjectDefinition {

    private final byte[] bodyBytes;

    public UnknownExtensionObject(byte[] bodyBytes) {
        super();
        this.bodyBytes = bodyBytes;
    }

    @Override
    public Integer getExtensionId() {
        return 0;
    }

    public byte[] getBodyBytes() {
        return bodyBytes;
    }

    @Override
    protected void serializeExtensionObjectDefinitionChild(WriteBuffer writeBuffer)
        throws SerializationException {
        writeBuffer.pushContext("UnknownExtensionObject");
        writeBuffer.writeByteArray("bodyBytes", bodyBytes);
        writeBuffer.popContext("UnknownExtensionObject");
    }

    @Override
    public int getLengthInBytes() {
        return (int) Math.ceil((float) getLengthInBits() / 8.0);
    }

    @Override
    public int getLengthInBits() {
        int lengthInBits = super.getLengthInBits();
        lengthInBits += bodyBytes.length * 8;
        return lengthInBits;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof UnknownExtensionObject)) return false;
        UnknownExtensionObject that = (UnknownExtensionObject) o;
        return super.equals(that) && Arrays.equals(bodyBytes, that.bodyBytes);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), Arrays.hashCode(bodyBytes));
    }

    @Override
    public String toString() {
        WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased(true, true);
        try {
            writeBufferBoxBased.writeSerializable(this);
        } catch (SerializationException e) {
            throw new RuntimeException(e);
        }
        return "\n" + writeBufferBoxBased.getBox().toString() + "\n";
    }
}
