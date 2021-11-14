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
package org.apache.plc4x.java.spi.codegen.fields;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.spi.codegen.io.DataWriter;
import org.apache.plc4x.java.spi.generation.*;

import java.util.List;

public class FieldWriterArray<T> implements FieldWriter<T> {

    @Override
    public void writeField(String logicalName, T value, DataWriter<T> dataWriter, WithWriterArgs... writerArgs) throws SerializationException {
        throw new NotImplementedException();
    }

    public void writeByteArrayField(String logicalName, byte[] values, DataWriter<byte[]> dataWriter, WithWriterArgs... writerArgs) throws SerializationException {
        if(values != null) {
            dataWriter.write(logicalName, values, writerArgs);
        }
    }

    public void writeSimpleTypeArrayField(String logicalName, List<T> values, DataWriter<T> dataWriter, WithWriterArgs... writerArgs) throws SerializationException {
        dataWriter.pushContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
        if(values != null) {
            for (T value : values) {
                switchSerializeByteOrderIfNecessary(() -> dataWriter.write("value", value, writerArgs), dataWriter, extractByteOder(writerArgs).orElse(null));
            }
        }
        dataWriter.popContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
    }

    public void writeComplexTypeArrayField(String logicalName, List<? extends Message> values, WriteBuffer writeBuffer, WithWriterArgs... writerArgs) throws SerializationException {
        writeBuffer.pushContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
        if(values != null) {
            for (Message value : values) {
                switchSerializeByteOrderIfNecessary(() -> {
                    value.serialize(writeBuffer);
                    return null;
                }, writeBuffer, extractByteOder(writerArgs).orElse(null));
            }
        }
        writeBuffer.popContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
    }

}
