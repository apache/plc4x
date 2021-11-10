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

import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.math.BigInteger;

public class DataWriterFactory {

    public static DataWriter<Boolean> writeBoolean(WriteBuffer writeBuffer) {
        return new DataWriterSimpleBoolean(writeBuffer, 1);
    }

    public static DataWriter<Byte> writeUnsignedByte(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleUnsignedByte(writeBuffer, bitLength);
    }

    public static DataWriter<Byte> writeByte(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleByte(writeBuffer, bitLength);
    }

    public static DataWriter<Short> writeUnsignedShort(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleUnsignedShort(writeBuffer, bitLength);
    }

    public static DataWriter<Integer> writeUnsignedInt(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleUnsignedInt(writeBuffer, bitLength);
    }

    public static DataWriter<Long> writeUnsignedLong(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleUnsignedLong(writeBuffer, bitLength);
    }

    public static DataWriter<BigInteger> writeUnsignedBigInteger(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleUnsignedBigInteger(writeBuffer, bitLength);
    }

    public static DataWriter<Byte> writeSignedByte(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleSignedByte(writeBuffer, bitLength);
    }

    public static DataWriter<Short> writeSignedShort(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleSignedShort(writeBuffer, bitLength);
    }

    public static DataWriter<Integer> writeSignedInt(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleSignedInt(writeBuffer, bitLength);
    }

    public static DataWriter<Long> writeSignedLong(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleSignedLong(writeBuffer, bitLength);
    }

    public static DataWriter<BigInteger> writeSignedBigInteger(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleSignedBigInteger(writeBuffer, bitLength);
    }

    public static DataWriter<Float> writeFloat(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleFloat(writeBuffer, bitLength);
    }

    public static DataWriter<Double> writeDouble(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleDouble(writeBuffer, bitLength);
    }

    public static DataWriter<String> writeString(WriteBuffer writeBuffer, int bitLength) {
        return new DataWriterSimpleString(writeBuffer, bitLength);
    }

    /*public static <T, I> DataWriterEnumDefault<T, I> writeEnum(Function<I, T> enumResolver, DataWriter<I> dataWriter) {
        return new DataWriterEnumDefault<>(enumResolver, dataWriter);
    }*/

}
