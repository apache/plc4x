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
package org.apache.plc4x.java.spi.codegen.fields;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataWriter;
import org.apache.plc4x.java.spi.generation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.Optional;

public class FieldWriterOptional<T> implements FieldCommons, WithReaderWriterArgs {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldWriterOptional.class);

    public void writeOptionalField(String logicalName, T value, DataWriter<T> dataWriter, boolean condition, WithWriterArgs... writerArgs) throws SerializationException {
        LOGGER.debug("write field {}", logicalName);
        if (condition && value != null) {
            switchSerializeByteOrderIfNecessary(() -> dataWriter.write(logicalName, value, writerArgs), dataWriter, extractByteOrder(writerArgs).orElse(null));
        } else {
            WriteBuffer writeBuffer = dataWriter.getWriteBuffer();
            // This is very special to byte based buffers, it would just confuse the others.
            if(writeBuffer instanceof WriteBufferByteBased) {
                // Check if a nullByteHex is set.
                // If it is, peek the equivalent number of bytes and compare.
                // If they match, return null.
                Optional<String> nullByteHexOptional = extractNullBytesHex(writerArgs);
                if (nullByteHexOptional.isPresent()) {
                    String nullByteHex = nullByteHexOptional.get();
                    try {
                        byte[] nullBytes = Hex.decodeHex(nullByteHex);
                        writeBuffer.writeByteArray(logicalName, nullBytes, writerArgs);
                    } catch (DecoderException e) {
                        // Ignore.
                    }
                } else {
                    LOGGER.debug("field {} not written because value is null({}) or condition({}) didn't evaluate to true", logicalName, value != null, condition);
                }
            }
        }
    }

}
