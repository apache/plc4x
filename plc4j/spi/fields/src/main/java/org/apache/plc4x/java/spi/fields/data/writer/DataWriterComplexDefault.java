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

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class DataWriterComplexDefault<T extends Message> implements DataWriterComplex<T> {

    private static final Logger LOGGER = LoggerFactory.getLogger(DataWriterComplexDefault.class);

    protected final WriteBuffer writeBuffer;

    public DataWriterComplexDefault(WriteBuffer writeBuffer) {
        this.writeBuffer = writeBuffer;
    }

    @Override
    public void pushContext(WithOption... options) {
        try {
            writeBuffer.pushContext(options);
        } catch (BufferException e) {
            throw new RuntimeException(e.getMessage(), e.getCause());
        }
    }

    @Override
    public void popContext(WithOption... options) {
        try {
            writeBuffer.popContext(options);
        } catch (BufferException e) {
            throw new RuntimeException(e.getMessage(), e.getCause());
        }
    }

    @Override
    public void write(T value, WithOption... options) throws BufferException {
        if (value == null) {
            LOGGER.warn("Trying to serialize null value");
        }
        writeBuffer.pushContext(options);
        writeBuffer.writeMessage(value);
        writeBuffer.popContext(options);
    }

    @Override
    public WriteBuffer getWriteBuffer() {
        return writeBuffer;
    }

}
