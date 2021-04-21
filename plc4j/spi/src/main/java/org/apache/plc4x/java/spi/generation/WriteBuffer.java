/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.java.spi.generation;

import java.math.BigDecimal;
import java.math.BigInteger;

public interface WriteBuffer {
    // TODO: check if this is really needed or if this is just an artifact
    int getPos();

    void pushContext(String logicalName, WithWriterArgs... writerArgs);

    void writeBit(String logicalName, boolean value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeBit(boolean value) throws ParseException {
        writeBit("", value);
    }

    void writeUnsignedByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeUnsignedByte(int bitLength, byte value) throws ParseException {
        writeUnsignedByte("", bitLength, value);
    }

    void writeUnsignedShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeUnsignedShort(int bitLength, short value) throws ParseException {
        writeUnsignedShort("", bitLength, value);
    }

    void writeUnsignedInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeUnsignedInt(int bitLength, int value) throws ParseException {
        writeUnsignedInt("", bitLength, value);
    }

    void writeUnsignedLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeUnsignedLong(int bitLength, long value) throws ParseException {
        writeUnsignedLong("", bitLength, value);
    }

    void writeUnsignedBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeUnsignedBigInteger(int bitLength, BigInteger value) throws ParseException {
        writeUnsignedBigInteger("", bitLength, value);
    }

    void writeByte(String logicalName, int bitLength, byte value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeByte(int bitLength, byte value) throws ParseException {
        writeByte("", bitLength, value);
    }

    void writeShort(String logicalName, int bitLength, short value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeShort(int bitLength, short value) throws ParseException {
        writeShort("", bitLength, value);
    }

    void writeInt(String logicalName, int bitLength, int value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeInt(int bitLength, int value) throws ParseException {
        writeInt("", bitLength, value);
    }

    void writeLong(String logicalName, int bitLength, long value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeLong(int bitLength, long value) throws ParseException {
        writeLong("", bitLength, value);
    }

    void writeBigInteger(String logicalName, int bitLength, BigInteger value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeBigInteger(int bitLength, BigInteger value) throws ParseException {
        writeBigInteger("", bitLength, value);
    }

    void writeFloat(String logicalName, float value, int bitsExponent, int bitsMantissa, WithWriterArgs... writerArgs) throws ParseException;

    default void writeFloat(float value, int bitsExponent, int bitsMantissa) throws ParseException {
        writeFloat("", value, bitsExponent, bitsMantissa);
    }

    void writeDouble(String logicalName, double value, int bitsExponent, int bitsMantissa, WithWriterArgs... writerArgs) throws ParseException;

    default void writeDouble(double value, int bitsExponent, int bitsMantissa) throws ParseException {
        writeDouble("", value, bitsExponent, bitsMantissa);
    }

    void writeBigDecimal(String logicalName, int bitLength, BigDecimal value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeBigDecimal(int bitLength, BigDecimal value) throws ParseException {
        writeBigDecimal("", bitLength, value);
    }

    void writeString(String logicalName, int bitLength, String encoding, String value, WithWriterArgs... writerArgs) throws ParseException;

    default void writeString(int bitLength, String encoding, String value) throws ParseException {
        writeString("", bitLength, encoding, value);
    }

    void popContext(String logicalName, WithWriterArgs... writerArgs);
}
