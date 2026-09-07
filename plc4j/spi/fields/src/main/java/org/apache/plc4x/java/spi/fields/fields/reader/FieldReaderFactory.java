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

package org.apache.plc4x.java.spi.fields.fields.reader;

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.data.reader.ParseSupplier;

import java.math.BigInteger;
import java.util.List;
import java.util.function.Function;
import java.util.function.Supplier;

public class FieldReaderFactory {

    // TODO: only used as lazy workaround
    @Deprecated
    public static <T> List<T> readCountArrayField(DataReader<T> dataReader, BigInteger count, WithOption... options) throws BufferException {
        if (count.bitLength() > 64) {
            throw new IllegalStateException("can't handle more than 64 bit. Actual: " + count.bitLength());
        }
        return readCountArrayField(dataReader, count.longValue(), options);
    }

    public static <T> List<T> readCountArrayField(DataReader<T> dataReader, long count, WithOption... options) throws BufferException {
        return new FieldReaderArray<T>().readArrayFieldCount(dataReader, count, options);
    }

    public static <T> List<T> readLengthArrayField(DataReader<T> dataReader, int length, WithOption... options) throws BufferException {
        return new FieldReaderArray<T>().readArrayFieldLength(dataReader, length, options);
    }

    /**
     * In some protocols a long value is used as length, but we simply can't address that many bytes,
     * so we simply cast it down to int as on java we couldn't even read more bytes as MAX-INT.
     *
     * @param dataReader the dataReader used to retrieve this field
     * @param length     the length of the array
     * @param options    optional read args
     * @param <T>        the type of the array elements
     * @return the read length array
     * @throws BufferException if something went wrong parsing
     */
    public static <T> List<T> readLengthArrayField(DataReader<T> dataReader, long length, WithOption... options) throws BufferException {
        return new FieldReaderArray<T>().readArrayFieldLength(dataReader, (int) length, options);
    }

    public static <T> List<T> readTerminatedArrayField(DataReader<T> dataReader, Supplier<Boolean> termination, WithOption... options) throws BufferException {
        return new FieldReaderArray<T>().readArrayFieldTerminated(dataReader, termination, options);
    }

    public static <T> T readAssertField(DataReader<T> dataReader, T expectedValue, WithOption... options) throws BufferException {
        return new FieldReaderAssert<T>().readAssertField(dataReader, expectedValue, options);
    }

    public static <T> T readChecksumField(DataReader<T> dataReader, T expectedValue, WithOption... options) throws BufferException {
        return new FieldReaderChecksum<T>().readChecksumField(dataReader, expectedValue, options);
    }

    public static <T> T readConstField(DataReader<T> dataReader, T expectedValue, WithOption... options) throws BufferException {
        return new FieldReaderConst<T>().readConstField(dataReader, expectedValue, options);
    }

    public static <T> T readDiscriminatorField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderDiscriminator<T>().readDiscriminatorField(dataReader, options);
    }

    public static <T> T readDiscriminatorEnumField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderDiscriminatorEnum<T>().readDiscriminatorEnumField(dataReader, options);
    }

    public static <T> T readEnumField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderEnum<T>().readEnumField(dataReader, options);
    }

    public static <T> T readImplicitField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderImplicit<T>().readImplicitField(dataReader, options);
    }

    public static <T> T readOptionalField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderOptional<T>().readOptionalField(dataReader, true, options);
    }

    public static <T> T readOptionalField(DataReader<T> dataReader, boolean condition, WithOption... options) throws BufferException {
        return new FieldReaderOptional<T>().readOptionalField(dataReader, condition, options);
    }

    public static byte[] readManualByteArrayField(ReadBuffer readBuffer, Function<List<Byte>, Boolean> termination, ParseSupplier<Byte> parse, WithOption... options) throws BufferException {
        return new FieldReaderManualArray<Byte>().readManualByteArrayField(readBuffer, termination, parse, options);
    }

    public static <T> List<T> readManualArrayField(ReadBuffer readBuffer, Function<List<T>, Boolean> termination, ParseSupplier<T> parse, WithOption... options) throws BufferException {
        return new FieldReaderManualArray<T>().readManualArrayField(readBuffer, termination, parse, options);
    }

    public static <T> T readManualField(ReadBuffer readBuffer, ParseSupplier<T> parse, WithOption... options) throws BufferException {
        return new FieldReaderManual<T>().readManualField(readBuffer, parse, options);
    }

    public static <T> void readPaddingField(DataReader<T> dataReader, int timesPadding, WithOption... options) throws BufferException {
        new FieldReaderPadding<T>().readPaddingField(dataReader, timesPadding, options);
    }

    public static <T> T readReservedField(DataReader<T> dataReader, T expectedValue, WithOption... options) throws BufferException {
        return new FieldReaderReserved<T>().readReservedField(dataReader, expectedValue, options);
    }

    public static <T> T readSimpleField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderSimple<T>().readSimpleField(dataReader, options);
    }

    public static <T> T readUnknownField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderUnknown<T>().readUnknownField(dataReader, options);
    }

    public static <T> T readVirtualField(Class<T> type, Object valueExpression, WithOption... options) throws BufferException {
        return new FieldReaderVirtual<T>().readVirtualField(type, valueExpression, options);
    }

    public static <T> T readPeekField(DataReader<T> dataReader, WithOption... options) throws BufferException {
        return new FieldReaderPeek<T>().readPeekField(dataReader, 0, options);
    }

    public static <T> T readPeekField(DataReader<T> dataReader, int offset, WithOption... options) throws BufferException {
        return new FieldReaderPeek<T>().readPeekField(dataReader, offset, options);
    }

}
