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

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;

public class DataReaderComplexDefault<T> implements DataReaderComplex<T> {

    private final ComplexTypeSupplier<T> complexSupplier;

    // TODO: maybe replace with resetable or something so its clear that that's the only purpose
    private final ReadBuffer readBuffer;

    public DataReaderComplexDefault(ComplexTypeSupplier<T> complexSupplier, ReadBuffer readBuffer) {
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
    public ByteOrder getByteOrder() {
        return readBuffer.getByteOrder();
    }

    @Override
    public void setByteOrder(ByteOrder byteOrder) {
        readBuffer.setByteOrder(byteOrder);
    }

    @Override
    public T read(String logicalName, WithReaderArgs... readerArgs) throws ParseException {
        return read(logicalName, complexSupplier, readerArgs);
    }

    public T read(String logicalName, ComplexTypeSupplier<T> complexSupplier, WithReaderArgs... readerArgs) throws ParseException {
        // TODO: it might be even better if we default to value like in other places... on the other hand a complex type has always a proper logical name so this might be fine like that
        boolean hasLogicalName = StringUtils.isNotBlank(logicalName);
        if (hasLogicalName) {
            readBuffer.pullContext(logicalName, readerArgs);
        }
        final T t = complexSupplier.get();
        if (hasLogicalName) {
            readBuffer.closeContext(logicalName, readerArgs);
        }
        return t;
    }

    @Override
    public void pullContext(String logicalName, WithReaderArgs... readerArgs) {
        readBuffer.pullContext(logicalName, readerArgs);
    }

    @Override
    public void closeContext(String logicalName, WithReaderArgs... readerArgs) {
        readBuffer.closeContext(logicalName, readerArgs);
    }

}

