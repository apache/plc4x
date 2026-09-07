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

package org.apache.plc4x.java.spi.fields.data.reader;

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
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
    public int getPositionInBits() {
        return dataReader.getPositionInBits();
    }

    @Override
    public void setPositionInBits(int positionInBits) {
        dataReader.setPositionInBits(positionInBits);
    }

    @Override
    public T read(WithOption... options) throws BufferException {
        return read(enumResolver, options);
    }

    public T read(Function<I, T> enumResolver, WithOption... options) throws BufferException {
        I rawValue = dataReader.read(options);
        T enumValue = enumResolver.apply(rawValue);
        if (enumValue == null) {
            LOGGER.debug("No {} enum found for value {}", WithOption.extractName(options), rawValue);
        }
        return enumValue;
    }

    @Override
    public void pushContext(WithOption... options) throws BufferException {
        dataReader.pushContext(options);
    }

    @Override
    public void popContext(WithOption... options) throws BufferException {
        dataReader.popContext(options);
    }

    @Override
    public ReadBuffer getReadBuffer() {
        return dataReader.getReadBuffer();
    }

}
