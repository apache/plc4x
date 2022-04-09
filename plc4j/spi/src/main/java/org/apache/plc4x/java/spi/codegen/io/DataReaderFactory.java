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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.plc4x.java.spi.generation.ReadBuffer;

import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.function.Function;

public class DataReaderFactory {

    public static DataReader<Boolean> readBoolean(ReadBuffer readBuffer) {
        return new DataReaderSimpleBoolean(readBuffer);
    }

    public static DataReader<Byte> readUnsignedByte(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleUnsignedByte(readBuffer, bitLength);
    }

    public static DataReader<Byte> readByte(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleByte(readBuffer, bitLength);
    }

    public static DataReader<Short> readUnsignedShort(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleUnsignedShort(readBuffer, bitLength);
    }

    public static DataReader<Integer> readUnsignedInt(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleUnsignedInt(readBuffer, bitLength);
    }

    public static DataReader<Long> readUnsignedLong(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleUnsignedLong(readBuffer, bitLength);
    }

    public static DataReader<BigInteger> readUnsignedBigInteger(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleUnsignedBigInteger(readBuffer, bitLength);
    }

    public static DataReader<Byte> readSignedByte(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleSignedByte(readBuffer, bitLength);
    }

    public static DataReader<Short> readSignedShort(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleSignedShort(readBuffer, bitLength);
    }

    public static DataReader<Integer> readSignedInt(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleSignedInt(readBuffer, bitLength);
    }

    public static DataReader<Long> readSignedLong(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleSignedLong(readBuffer, bitLength);
    }

    public static DataReader<BigInteger> readSignedBigInteger(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleSignedBigInteger(readBuffer, bitLength);
    }

    public static DataReader<Float> readFloat(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleFloat(readBuffer, bitLength);
    }

    public static DataReader<Double> readDouble(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleDouble(readBuffer, bitLength);
    }

    public static DataReader<String> readString(ReadBuffer readBuffer, int bitLength) {
        return new DataReaderSimpleString(readBuffer, bitLength);
    }

    public static <T, I> DataReaderEnumDefault<T, I> readEnum(Function<I, T> enumResolver, DataReader<I> dataReader) {
        return new DataReaderEnumDefault<>(enumResolver, dataReader);
    }

    public static DataReader<LocalDate> readDate(ReadBuffer readBuffer) {
        return new DataReaderSimpleDate(readBuffer);
    }

    public static DataReader<LocalDateTime> readDateTime(ReadBuffer readBuffer) {
        return new DataReaderSimpleDateTime(readBuffer);
    }

    public static DataReader<LocalTime> readTime(ReadBuffer readBuffer) {
        return new DataReaderSimpleTime(readBuffer);
    }

}
