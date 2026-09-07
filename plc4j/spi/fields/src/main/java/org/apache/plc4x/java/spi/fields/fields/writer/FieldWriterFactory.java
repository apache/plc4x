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

package org.apache.plc4x.java.spi.fields.fields.writer;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.writer.DataWriter;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;

import java.util.ArrayList;
import java.util.List;

public class FieldWriterFactory {

    public static <T> void writeSimpleTypeArrayField(List<T> value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterArray<T>().writeSimpleTypeArrayField(value, dataWriter, options);
    }

    public static void writeComplexTypeArrayField(List<? extends Message> value, WriteBuffer writeBuffer, WithOption... options) throws BufferException {
        new FieldWriterArray<Message>().writeComplexTypeArrayField(value, writeBuffer, options);
    }

    public static <T> void writeByteArrayField(byte[] value, DataWriter<byte[]> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterArray<T>().writeByteArrayField(value, dataWriter, options);
    }

    public static <T> void writeChecksumField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterChecksum<T>().writeChecksumField(value, dataWriter, options);
    }

    public static <T> void writeConstField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterConst<T>().writeConstField(value, dataWriter, options);
    }

    public static <T> void writeEnumField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterEnum<T>().writeEnumField(value, dataWriter, options);
    }

    public static <T> void writeDiscriminatorField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterDiscriminator<T>().writeDiscriminatorField(value, dataWriter, options);
    }

    public static <T> void writeDiscriminatorEnumField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterDiscriminatorEnum<T>().writeDiscriminatorEnumField(value, dataWriter, options);
    }

    public static <T> void writeImplicitField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterImplicit<T>().writeImplicitField(value, dataWriter, options);
    }

    // TODO: Check if we shouldn't pass in the value directly.
    public static <T> void writeManualField(FieldCommons.RunSerializeWrapped runnable, WriteBuffer writeBuffer, WithOption... options) throws BufferException {
        new FieldWriterManual<>().writeManualField(runnable, writeBuffer, options);
    }

    public static void writeManualArrayField(byte[] bytes, FieldCommons.ConsumeSerializeWrapped<Byte> runnable, WriteBuffer writeBuffer, WithOption... options) throws BufferException {
        List<Byte> list = new ArrayList<>(bytes.length);
        for (byte b : bytes) {
            list.add(b);
        }
        new FieldWriterManualArray<Byte>().writeManualArrayField(list, runnable, writeBuffer, options);
    }

    public static <T> void writeManualArrayField(List<T> values, FieldCommons.ConsumeSerializeWrapped<T> runnable, WriteBuffer writeBuffer, WithOption... options) throws BufferException {
        new FieldWriterManualArray<T>().writeManualArrayField(values, runnable, writeBuffer, options);
    }

    public static <T> void writeOptionalField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterOptional<T>().writeOptionalField(value, dataWriter, true, options);
    }

    public static <T> void writeOptionalField(T value, DataWriter<T> dataWriter, boolean condition, WithOption... options) throws BufferException {
        new FieldWriterOptional<T>().writeOptionalField(value, dataWriter, condition, options);
    }

    public static <T> void writeOptionalEnumField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterOptionalEnum<T>().writeOptionalEnumField(value, dataWriter, true, options);
    }

    public static <T> void writeOptionalEnumField(T value, DataWriter<T> dataWriter, boolean condition, WithOption... options) throws BufferException {
        new FieldWriterOptionalEnum<T>().writeOptionalEnumField(value, dataWriter, condition, options);
    }

    public static <T> void writePaddingField(int timesPadding, T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterPadding<T>().writePaddingField(timesPadding, value, dataWriter, options);
    }

    public static <T> void writeReservedField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterReserved<T>().writeReservedField(value, dataWriter, options);
    }

    public static <T> void writeSimpleField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterSimple<T>().writeSimpleField(value, dataWriter, options);
    }

    public static <T> void writeSimpleEnumField(T value, DataWriter<T> dataWriter, WithOption... options) throws BufferException {
        new FieldWriterSimpleEnum<T>().writeSimpleEnumField(value, dataWriter, options);
    }

}
