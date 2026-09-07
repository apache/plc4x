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

package org.apache.plc4x.java.spi.fields.data.writer;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.function.Function;

public class DataWriterEnumDefault<T, I> implements DataWriterEnum<T> {

    private static final Logger LOGGER = LoggerFactory.getLogger(DataWriterEnumDefault.class);

    private final Function<T, I> enumSerializer;
    private final Function<T, String> enumNamer;
    private final DataWriter<I> dataWriter;

    public DataWriterEnumDefault(Function<T, I> enumSerializer, Function<T, String> enumNamer, DataWriter<I> dataWriter) {
        this.enumSerializer = enumSerializer;
        this.dataWriter = dataWriter;
        this.enumNamer = enumNamer;
    }

    @Override
    public void write(T value, WithOption... options) throws BufferException {
        write(value, enumSerializer, enumNamer, dataWriter, options);
    }

    public void write(T value, Function<T, I> enumSerializer, Function<T, String> enumNamer, DataWriter<I> rawWriter, WithOption... options) throws BufferException {
        if (value == null) {
            LOGGER.warn("Trying to serialize null value");
            return;
        }
        final I rawValue = enumSerializer.apply(value);

        WithOption[] newOptions = WithOption.AddOptions(options, WithOption.WithAdditionalStringRepresentation(enumNamer.apply(value)));
        rawWriter.write(rawValue, newOptions);
    }

    @Override
    public void pushContext(WithOption... options) throws BufferException {
        dataWriter.pushContext(options);
    }

    @Override
    public void popContext(WithOption... options) throws BufferException {
        dataWriter.popContext(options);
    }

    @Override
    public WriteBuffer getWriteBuffer() {
        return dataWriter.getWriteBuffer();
    }

}
