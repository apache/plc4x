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

import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.function.Function;

public class DataReaderEnumDefault<T, I> implements DataReaderEnum<T> {

    private static final Logger LOGGER = LoggerFactory.getLogger(DataReaderEnumDefault.class);

    private final Function<I, T> enumResolver;
    private final DataReader<I> dataReader;

    public DataReaderEnumDefault(Function<I, T> enumResolver, DataReader<I> dataReader) {
        this.enumResolver = enumResolver;
        this.dataReader = dataReader;
    }

    @Override
    public int getPos() {
        return dataReader.getPos();
    }

    @Override
    public void setPos(int position) {
        dataReader.setPos(position);
    }

    @Override
    public ByteOrder getByteOrder() {
        return dataReader.getByteOrder();
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        dataReader.setByteOrder(byteOrder);
    }

    @Override
    public T read(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        return read(logicalName, enumResolver, readerArgs);
    }

    public T read(String logicalName, Function<I, T> enumResolver, WithReaderArgs... readerArgs) throws ParseException {
        I rawValue = dataReader.read(logicalName, readerArgs);
        T enumValue = enumResolver.apply(rawValue);
        if (enumValue == null) {
            LOGGER.debug("No {} enum found for value {}", logicalName, rawValue);
        }
        return enumValue;
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
        dataReader.pullContext(logicalName, readerArgs);
    }

    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
        dataReader.closeContext(logicalName, readerArgs);
    }

    @Override
    public ReadBuffer getReadBuffer() {
        return dataReader.getReadBuffer();
    }

}
