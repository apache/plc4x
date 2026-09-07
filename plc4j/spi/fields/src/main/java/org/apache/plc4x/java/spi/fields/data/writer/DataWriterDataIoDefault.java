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
import org.apache.plc4x.java.spi.fields.data.reader.DataIoSerializerFunction;
import org.apache.plc4x.java.api.value.PlcValue;

public class DataWriterDataIoDefault implements DataWriterComplex<PlcValue> {

    protected final WriteBuffer writeBuffer;
    protected final DataIoSerializerFunction<WriteBuffer> serializer;

    public DataWriterDataIoDefault(WriteBuffer writeBuffer, DataIoSerializerFunction<WriteBuffer> serializer) {
        this.writeBuffer = writeBuffer;
        this.serializer = serializer;
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
    public void write(PlcValue value, WithOption... options) throws BufferException {
        writeBuffer.pushContext(options);
        serializer.apply(writeBuffer, value);
        writeBuffer.popContext(options);
    }

    @Override
    public WriteBuffer getWriteBuffer() {
        return writeBuffer;
    }

}
