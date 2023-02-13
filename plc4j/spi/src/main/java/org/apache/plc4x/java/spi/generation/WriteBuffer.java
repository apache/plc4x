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

    default void writeBit(boolean value, WithWriterArgs... writerArgs) throws SerializationException {
        writeBit("", value, writerArgs);
    }

    void writeByte(String logicalName, byte value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeByte(byte value, WithWriterArgs... writerArgs) throws SerializationException {
        writeByte("", value, writerArgs);
    }

    void writeByteArray(String logicalName, byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeByteArray(byte[] bytes, WithWriterArgs... writerArgs) throws SerializationException {
        writeByteArray("", bytes, writerArgs);
    }

    void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedByte(int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        writeUnsignedByte("", bitLength, value, writerArgs);
    }

    void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedShort(int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        writeUnsignedShort("", bitLength, value, writerArgs);
    }

    void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedInt(int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        writeUnsignedInt("", bitLength, value, writerArgs);
    }

    void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedLong(int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        writeUnsignedLong("", bitLength, value, writerArgs);
    }

    void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeUnsignedBigInteger(int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        writeUnsignedBigInteger("", bitLength, value, writerArgs);
    }

    void writeSignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeSignedByte(int bitLength, byte value, WithWriterArgs... writerArgs) throws SerializationException {
        writeSignedByte("", bitLength, value, writerArgs);
    }

    void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeShort(int bitLength, short value, WithWriterArgs... writerArgs) throws SerializationException {
        writeShort("", bitLength, value, writerArgs);
    }

    void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeInt(int bitLength, int value, WithWriterArgs... writerArgs) throws SerializationException {
        writeInt("", bitLength, value, writerArgs);
    }

    void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeLong(int bitLength, long value, WithWriterArgs... writerArgs) throws SerializationException {
        writeLong("", bitLength, value, writerArgs);
    }

    void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeBigInteger(int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws SerializationException {
        writeBigInteger("", bitLength, value, writerArgs);
    }

    void writeFloat(String logicalName, int bitLength, float value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeFloat(int bitLength, float value, WithWriterArgs... writerArgs) throws SerializationException {
        writeFloat("", bitLength, value, writerArgs);
    }

    void writeDouble(String logicalName, int bitLength, double value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeDouble(int bitLength, double value, WithWriterArgs... writerArgs) throws SerializationException {
        writeDouble("", bitLength, value, writerArgs);
    }

    void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeBigDecimal(int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws SerializationException {
        writeBigDecimal("", bitLength, value, writerArgs);
    }

    void writeString(String logicalName, int bitLength, String value, WithWriterArgs... writerArgs) throws SerializationException;

    default void writeVirtual(String logicalName, Object value, WithWriterArgs... writerArgs) throws SerializationException {
        // No-Op
    }

    default void writeString(int bitLength, String value, WithWriterArgs... writerArgs) throws SerializationException {
        writeString("", bitLength, value, writerArgs);
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
