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

import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataWriter;
import org.apache.plc4x.java.spi.generation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class FieldWriterArray<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldWriterArray.class);

    public void writeByteArrayField(String logicalName, byte[] values, DataWriter<byte[]> dataWriter, WithWriterArgs... writerArgs) throws SerializationException {
        LOGGER.debug("write field {}", logicalName);
        if (values != null) {
            dataWriter.write(logicalName, values, writerArgs);
        }
    }

    public void writeSimpleTypeArrayField(String logicalName, List<T> values, DataWriter<T> dataWriter, WithWriterArgs... writerArgs) throws SerializationException {
        LOGGER.debug("write field {}", logicalName);
        switchSerializeByteOrderIfNecessary(() -> {
            if (values != null) {
                dataWriter.pushContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
                for (T value : values) {
                    dataWriter.write("value", value, writerArgs);
                }
                dataWriter.popContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
            }
        }, dataWriter, extractByteOder(writerArgs).orElse(null));
    }

    public void writeComplexTypeArrayField(String logicalName, List<? extends Message> values, WriteBuffer writeBuffer, WithWriterArgs... writerArgs) throws SerializationException {
        LOGGER.debug("write field {}", logicalName);
        switchSerializeByteOrderIfNecessary(() -> {
            if (values != null) {
                writeBuffer.pushContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
                for (Message value : values) {
                    value.serialize(writeBuffer);
                }
                writeBuffer.popContext(logicalName, WithReaderWriterArgs.WithRenderAsList(true));
            }
        }, writeBuffer, extractByteOder(writerArgs).orElse(null));
    }

}
