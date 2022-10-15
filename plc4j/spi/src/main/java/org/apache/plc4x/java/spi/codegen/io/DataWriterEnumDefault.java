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
package org.apache.plc4x.java.spi.codegen.io;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WithReaderWriterArgs;
import org.apache.plc4x.java.spi.generation.WithWriterArgs;
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
    public ByteOrder getByteOrder() {
        return dataWriter.getByteOrder();
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        dataWriter.setByteOrder(byteOrder);
    }

    @Override
    public void write(String logicalName, T value, WithWriterArgs... writerArgs) throws SerializationException {
        write(logicalName, value, enumSerializer, enumNamer, dataWriter, writerArgs);
    }

    public void write(String logicalName, T value, Function<T, I> enumSerializer, Function<T, String> enumNamer, DataWriter<I> rawWriter, WithWriterArgs... writerArgs) throws SerializationException {
        if (value == null) {
            LOGGER.warn("Trying to serialize null value for {}", logicalName);
            return;
        }
        final I rawValue = enumSerializer.apply(value);
        rawWriter.write(logicalName, rawValue, ArrayUtils.addAll(writerArgs, WithReaderWriterArgs.WithAdditionalStringRepresentation(enumNamer.apply(value))));
    }

    @Override
    public void pushContext(String logicalName, WithWriterArgs... writerArgs) {
        dataWriter.pushContext(logicalName, writerArgs);
    }

    @Override
    public void popContext(String logicalName, WithWriterArgs... writerArgs) {
        dataWriter.popContext(logicalName, writerArgs);
    }

}
