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

public class DataReaderFactory {

    public static DataReader<Boolean>  readBoolean(ReadBuffer readBuffer) {
        return new DataReaderSimpleBoolean(readBuffer);
    }

    public static DataReader<Byte> readUnsignedByte(ReadBuffer readBuffer) {
        return new DataReaderSimpleUnsignedByte(readBuffer);
    }

    public static DataReader<Short> readUnsignedShort(ReadBuffer readBuffer) {
        return new DataReaderSimpleUnsignedShort(readBuffer);
    }

    public static DataReader<Integer> readUnsignedInt(ReadBuffer readBuffer) {
        return new DataReaderSimpleUnsignedInt(readBuffer);
    }

    public static DataReader<Long> readUnsignedLong(ReadBuffer readBuffer) {
        return new DataReaderSimpleUnsignedLong(readBuffer);
    }

    public static DataReader<BigInteger> readUnsignedBigInteger(ReadBuffer readBuffer) {
        return new DataReaderSimpleUnsignedBigInteger(readBuffer);
    }

    public static DataReader<Byte> readSignedByte(ReadBuffer readBuffer) {
        return new DataReaderSimpleSignedByte(readBuffer);
    }

    public static DataReader<Short> readSignedShort(ReadBuffer readBuffer) {
        return new DataReaderSimpleSignedShort(readBuffer);
    }

    public static DataReader<Integer> readSignedInt(ReadBuffer readBuffer) {
        return new DataReaderSimpleSignedInt(readBuffer);
    }

    public static DataReader<Long> readSignedLong(ReadBuffer readBuffer) {
        return new DataReaderSimpleSignedLong(readBuffer);
    }

    public static DataReader<BigInteger> readSignedBigInteger(ReadBuffer readBuffer) {
        return new DataReaderSimpleSignedBigInteger(readBuffer);
    }

    public static DataReader<Float> readFloat(ReadBuffer readBuffer) {
        return new DataReaderSimpleFloat(readBuffer);
    }

    public static DataReader<Double> readDouble(ReadBuffer readBuffer) {
        return new DataReaderSimpleDouble(readBuffer);
    }

    public static DataReader<String > readString(ReadBuffer readBuffer) {
        return new DataReaderSimpleString(readBuffer);
    }

}
