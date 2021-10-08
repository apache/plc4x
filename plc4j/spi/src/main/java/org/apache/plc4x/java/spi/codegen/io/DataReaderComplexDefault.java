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

import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;

import java.util.function.Supplier;

public class DataReaderComplexDefault<T> implements DataReaderComplex<T> {
    private final Supplier<T> complexSupplier;

    // TODO: maybe replace with resetable or something so its clear that that's the only purpose
    private final ReadBuffer readBuffer;

    public DataReaderComplexDefault(Supplier<T> complexSupplier, ReadBuffer readBuffer) {
        this.complexSupplier = complexSupplier;
        this.readBuffer = readBuffer;
    }

    @Override
    public int getPos() {
        return readBuffer.getPos();
    }

    @Override
    public void setPos(int position) {
        readBuffer.reset(position);
    }

    @Override
    public T read(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        return read(logicalName, complexSupplier, readerArgs);
    }

    public T read(String logicalName, Supplier<T> complexSupplier, WithReaderArgs... readerArgs) throws ParseException {
        readBuffer.pullContext(logicalName,readerArgs);
        final T t = complexSupplier.get();
        readBuffer.closeContext(logicalName,readerArgs);
        return t;
    }
}

