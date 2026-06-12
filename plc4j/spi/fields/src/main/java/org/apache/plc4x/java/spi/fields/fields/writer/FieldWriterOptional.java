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

package org.apache.plc4x.java.spi.fields.fields.writer;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.writer.DataWriter;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.apache.plc4x.java.spi.fields.fields.WithFieldOption;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Optional;

public class FieldWriterOptional<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldWriterOptional.class);

    public void writeOptionalField(T value, DataWriter<T> dataWriter, boolean condition, WithOption... options) throws BufferException {
        LOGGER.debug("write field {}", getName(options));
        if (condition && value != null) {
            dataWriter.write(value, options);
        } else {
            WriteBuffer writeBuffer = dataWriter.getWriteBuffer();
            // This is very special to byte-based buffers, it would just confuse the others.
            if (writeBuffer.isByteBased()) {
                // Check if a nullByteHex is set.
                // If it is, peek the equivalent number of bytes and compare.
                // If they match, return null.
                Optional<String> nullByteHexOptional = WithFieldOption.extractNullBytesHex(options);
                if (nullByteHexOptional.isPresent()) {
                    String nullByteHex = nullByteHexOptional.get();
                    try {
                        byte[] nullBytes = StaticHelper.DECODE_HEX(nullByteHex);
                        dataWriter.getWriteBuffer().writeBits(nullBytes.length * 8, nullBytes);
                    } catch (Exception e) {
                        // Ignore.
                    }
                } else {
                    LOGGER.debug("field {} not written because value is null({}) or condition({}) didn't evaluate to true", getName(options), value != null, condition);
                }
            }
        }
    }

}
