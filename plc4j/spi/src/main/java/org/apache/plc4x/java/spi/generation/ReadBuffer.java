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

import java.math.BigDecimal;
import java.math.BigInteger;

public interface ReadBuffer extends ByteOrderAware, PositionAware {
    int getPos();

    void reset(int pos);

    boolean hasMore(int numBits);

    void pullContext(String logicalName, WithReaderArgs... readerArgs);

    boolean readBit(String logicalName, WithReaderArgs... readerArgs) throws ParseException;

    default boolean readBit() throws ParseException {
        return readBit("");
    }

    byte readByte(String logicalName, WithReaderArgs... readerArgs) throws ParseException;

    default byte readByte() throws ParseException {
        return readSignedByte("", 8);
    }

    byte[] readByteArray(String logicalName, int numberOfBytes, WithReaderArgs... readerArgs) throws ParseException;

    default byte[] readByteArray(int numbersOfBytes, WithReaderArgs... readerArgs) throws ParseException {
        return readByteArray("", numbersOfBytes, readerArgs);
    }

    byte readUnsignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default byte readUnsignedByte(int bitLength) throws ParseException {
        return readUnsignedByte("", bitLength);
    }

    short readUnsignedShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default short readUnsignedShort(int bitLength) throws ParseException {
        return readUnsignedShort("", bitLength);
    }

    int readUnsignedInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default int readUnsignedInt(int bitLength) throws ParseException {
        return readUnsignedInt("", bitLength);
    }

    long readUnsignedLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default long readUnsignedLong(int bitLength) throws ParseException {
        return readUnsignedLong("", bitLength);
    }

    BigInteger readUnsignedBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default BigInteger readUnsignedBigInteger(int bitLength) throws ParseException {
        return readUnsignedBigInteger("", bitLength);
    }

    byte readSignedByte(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default byte readSignedByte(int bitLength) throws ParseException {
        return readSignedByte("", bitLength);
    }

    short readShort(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default short readShort(int bitLength) throws ParseException {
        return readShort("", bitLength);
    }

    int readInt(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default int readInt(int bitLength) throws ParseException {
        return readInt("", bitLength);
    }

    long readLong(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default long readLong(int bitLength) throws ParseException {
        return readLong("", bitLength);
    }

    BigInteger readBigInteger(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default BigInteger readBigInteger(int bitLength) throws ParseException {
        return readBigInteger("", bitLength);
    }

    float readFloat(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default float readFloat(int bitLength) throws ParseException {
        return readFloat("", bitLength);
    }

    double readDouble(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default double readDouble(int bitLength) throws ParseException {
        return readDouble("", bitLength);
    }

    BigDecimal readBigDecimal(String logicalName, int bitLength, WithReaderArgs... readerArgs) throws ParseException;

    default BigDecimal readBigDecimal(int bitLength) throws ParseException {
        return readBigDecimal("", bitLength);
    }

    String readString(String logicalName, int bitLength, String encoding, WithReaderArgs... readerArgs) throws ParseException;

    default String readString(int bitLength, String encoding) throws ParseException {
        return readString("", bitLength, encoding);
    }


    void closeContext(String logicalName, WithReaderArgs... readerArgs);
}
