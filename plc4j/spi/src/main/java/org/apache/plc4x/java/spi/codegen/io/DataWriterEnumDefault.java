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

import org.apache.plc4x.java.spi.generation.*;

import java.util.function.Function;

public class DataWriterEnumDefault<T, I> implements DataWriterEnum<T> {

    private final Function<I, T> enumResolver;
    private final DataWriter<I> dataWriter;

    public DataWriterEnumDefault(Function<I, T> enumResolver, DataWriter<I> dataWriter) {
        this.enumResolver = enumResolver;
        this.dataWriter = dataWriter;
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
    public T write(String logicalName, T value, WithWriterArgs... writerArgs) throws SerializationException {
        throw new SerializationException("Unsupported");
    }

    public T write(String logicalName, T value, Function<I, T> enumResolver, DataWriter<T> rawWriter, WithWriterArgs... writerArgs) throws SerializationException {
        final T enumValue = value;//enumResolver.apply(value);
        rawWriter.write(logicalName, enumValue, writerArgs);
        return null;
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
