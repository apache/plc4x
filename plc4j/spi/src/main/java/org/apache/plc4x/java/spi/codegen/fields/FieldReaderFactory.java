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
package org.apache.plc4x.java.spi.codegen.fields;

import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.codegen.io.ParseSupplier;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;

import java.math.BigInteger;
import java.util.List;
import java.util.function.Function;
import java.util.function.Supplier;

public class FieldReaderFactory {

    @SuppressWarnings("unused")
    public static <T> T readAbstractField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderAbstract<T>().readAbstractField(logicalName, dataReader, readerArgs);
    }

    // TODO: only used as lazy workaround
    @Deprecated
    public static <T> List<T> readCountArrayField(String logicalName, DataReader<T> dataReader, BigInteger count, WithReaderArgs... readerArgs) throws ParseException {
        if (count.bitLength() > 64) {
            throw new IllegalStateException("can't handle more than 64 bit. Actual: " + count.bitLength());
        }
        return readCountArrayField(logicalName, dataReader, count.longValue(), readerArgs);
    }

    public static <T> List<T> readCountArrayField(String logicalName, DataReader<T> dataReader, long count, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderArray<T>().readFieldCount(logicalName, dataReader, count, readerArgs);
    }

    public static <T> List<T> readLengthArrayField(String logicalName, DataReader<T> dataReader, int length, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderArray<T>().readFieldLength(logicalName, dataReader, length, readerArgs);
    }

    /**
     * In some protocols a long is used as length, but we simply can't address that many bytes,
     * so we simply cast it down to int as on java we couldn't even read more bytes as MAX-INT.
     *
     * @param logicalName the logical name of this field
     * @param dataReader  the dataReader used to retrieve this field
     * @param length      the length of the array
     * @param readerArgs  optional read args
     * @param <T>         the type of the array elements
     * @return the read length array
     * @throws ParseException if something went wrong parsing
     */
    public static <T> List<T> readLengthArrayField(String logicalName, DataReader<T> dataReader, long length, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderArray<T>().readFieldLength(logicalName, dataReader, (int) length, readerArgs);
    }

    public static <T> List<T> readTerminatedArrayField(String logicalName, DataReader<T> dataReader, Supplier<Boolean> termination, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderArray<T>().readFieldTerminated(logicalName, dataReader, termination, readerArgs);
    }

    public static <T> T readAssertField(String logicalName, DataReader<T> dataReader, T expectedValue, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderAssert<T>().readAssertField(logicalName, dataReader, expectedValue, readerArgs);
    }

    public static <T> T readChecksumField(String logicalName, DataReader<T> dataReader, T expectedValue, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderChecksum<T>().readChecksumField(logicalName, dataReader, expectedValue, readerArgs);
    }

    public static <T> T readConstField(String logicalName, DataReader<T> dataReader, T expectedValue, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderConst<T>().readConstField(logicalName, dataReader, expectedValue, readerArgs);
    }

    public static <T> T readDiscriminatorField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderDiscriminator<T>().readDiscriminatorField(logicalName, dataReader, readerArgs);
    }

    public static <T> T readEnumField(String logicalName, String innerName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderEnum<T>().readEnumField(logicalName, innerName, dataReader, readerArgs);
    }

    public static <T> T readImplicitField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderImplicit<T>().readImplicitField(logicalName, dataReader, readerArgs);
    }

    public static <T> T readOptionalField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderOptional<T>().readOptionalField(logicalName, dataReader, true, readerArgs);
    }

    public static <T> T readOptionalField(String logicalName, DataReader<T> dataReader, boolean condition, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderOptional<T>().readOptionalField(logicalName, dataReader, condition, readerArgs);
    }

    public static byte[] readManualByteArrayField(String logicalName, ReadBuffer readBuffer, Function<byte[], Boolean> termination, ParseSupplier<Byte> parse, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderManualArray<Byte>().readManualByteArrayField(logicalName, readBuffer, termination, parse, readerArgs);
    }

    public static <T> List<T> readManualArrayField(String logicalName, ReadBuffer readBuffer, Function<List<T>, Boolean> termination, ParseSupplier<T> parse, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderManualArray<T>().readManualArrayField(logicalName, readBuffer, termination, parse, readerArgs);
    }

    public static <T> T readManualField(String logicalName, ReadBuffer readBuffer, ParseSupplier<T> parse, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderManual<T>().readManualField(logicalName, readBuffer, parse, readerArgs);
    }

    public static <T> void readPaddingField(DataReader<T> dataReader, int timesPadding, WithReaderArgs... readerArgs) throws ParseException {
        new FieldReaderPadding<T>().readPaddingField(dataReader, timesPadding, readerArgs);
    }

    public static <T> T readReservedField(String logicalName, DataReader<T> dataReader, T expectedValue, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderReserved<T>().readReservedField(logicalName, dataReader, expectedValue, readerArgs);
    }

    public static <T> T readSimpleField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderSimple<T>().readSimpleField(logicalName, dataReader, readerArgs);
    }

    public static <T> T readUnknownField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderUnknown<T>().readUnknownField(logicalName, dataReader, readerArgs);
    }

    public static <T> T readVirtualField(String logicalName, Class<T> type, Object valueExpression, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderVirtual<T>().readVirtualField(logicalName, type, valueExpression, readerArgs);
    }

    public static <T> T readPeekField(String logicalName, DataReader<T> dataReader, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderPeek<T>().readPeekField(logicalName, dataReader, 0, readerArgs);
    }

    public static <T> T readPeekField(String logicalName, DataReader<T> dataReader, int offset, WithReaderArgs... readerArgs) throws ParseException {
        return new FieldReaderPeek<T>().readPeekField(logicalName, dataReader, offset, readerArgs);
    }

}
