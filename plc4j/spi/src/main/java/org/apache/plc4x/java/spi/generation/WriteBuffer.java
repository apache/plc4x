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
package org.apache.plc4x.java.spi.generation;

import org.apache.plc4x.java.spi.codegen.io.ByteOrderAware;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.math.BigDecimal;
import java.math.BigInteger;

public interface WriteBuffer extends ByteOrderAware, PositionAware {

    int getPos();

    void pushContext(String logicalName, WithWriterArgs... writerArgs);

    void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeBit(boolean value) throws SerializationException {
        writeBit("", value);
    }

    void writeByte(String logicalName, byte value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeByte(byte value) throws SerializationException {
        writeByte("", value);
    }

    void writeByteArray(String logicalName, byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeByteArray(byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException {
        writeByteArray("", bytes, writerArgs);
    }

    void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedByte(int bitLength, byte value) throws SerializationException {
        writeUnsignedByte("", bitLength, value);
    }

    void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedShort(int bitLength, short value) throws SerializationException {
        writeUnsignedShort("", bitLength, value);
    }

    void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedInt(int bitLength, int value) throws SerializationException {
        writeUnsignedInt("", bitLength, value);
    }

    void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedLong(int bitLength, long value) throws SerializationException {
        writeUnsignedLong("", bitLength, value);
    }

    void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedBigInteger(int bitLength, BigInteger value) throws SerializationException {
        writeUnsignedBigInteger("", bitLength, value);
    }

    void writeSignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeSignedByte(int bitLength, byte value) throws SerializationException {
        writeSignedByte("", bitLength, value);
    }

    void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeShort(int bitLength, short value) throws SerializationException {
        writeShort("", bitLength, value);
    }

    void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeInt(int bitLength, int value) throws SerializationException {
        writeInt("", bitLength, value);
    }

    void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeLong(int bitLength, long value) throws SerializationException {
        writeLong("", bitLength, value);
    }

    void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeBigInteger(int bitLength, BigInteger value) throws SerializationException {
        writeBigInteger("", bitLength, value);
    }

    void writeFloat(String logicalName, int bitLength, float value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeFloat(int bitLength, float value) throws SerializationException {
        writeFloat("", bitLength, value);
    }

    void writeDouble(String logicalName, int bitLength, double value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeDouble(int bitLength, double value) throws SerializationException {
        writeDouble("", bitLength, value);
    }

    void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeBigDecimal(int bitLength, BigDecimal value) throws SerializationException {
        writeBigDecimal("", bitLength, value);
    }

    void writeString(String logicalName, int bitLength, String encoding, String value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeVirtual(String logicalName, Object value, WithWriterArgs... writerArgs) throws SerializationException {
        // No-Op
    }

    default void writeString(int bitLength, String encoding, String value) throws SerializationException {
        writeString("", bitLength, encoding, value);
    }

    /**
     * this method can be used to influence serializing (e.g. intercept whole types and render them in a simplified form)
     *
     * @param value the value to be serialized
     * @throws SerializationException if something goes wrong
     */
    default void writeSerializable(Serializable value) throws SerializationException {
        if (value == null) {
            return;
        }
        value.serialize(this);
    }

    void popContext(String logicalName, WithWriterArgs... writerArgs);
}
