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

package org.apache.plc4x.java.spi.fields.fields.reader;

import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferUnderflowException;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferValueException;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.exceptions.ParseAssertException;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.apache.plc4x.java.spi.fields.fields.WithFieldOption;
import org.apache.plc4x.java.spi.utils.StaticHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.Optional;

public class FieldReaderOptional<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderOptional.class);

    public T readOptionalField(DataReader<T> dataReader, boolean condition, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}", getName(options));
        if (!condition) {
            LOGGER.debug("Condition doesn't match for field {}", getName(options));
            return null;
        }

        // Check if a nullByteHex is set.
        // If it is, peek the equivalent number of bytes and compare.
        // If they match, return null.
        Optional<String> nullByteHexOptional = WithFieldOption.extractNullBytesHex(options);
        if (nullByteHexOptional.isPresent()) {
            String nullByteHex = nullByteHexOptional.get();
            try {
                byte[] nullBytes = StaticHelper.DECODE_HEX(nullByteHex);
                ReadBuffer readBuffer = dataReader.getReadBuffer();
                int pos = readBuffer.getPositionInBits();
                byte[] curBytes = readBuffer.readBits(nullBytes.length, options);
                // Compare them, if they equal, return null, if not, reset the position and try to read it again.
                if (Arrays.equals(nullBytes, curBytes)) {
                    // Abort with null
                    return null;
                } else {
                    readBuffer.setPositionInBits(pos);
                }
            } catch (Exception e) {
                // Ignore.
            }
        }

        int curPosInBits = dataReader.getPositionInBits();
        try {
            T field = dataReader.read(options);
            LOGGER.debug("done reading field {}. Value: {}", getName(options), field);
            return field;
        } catch (BufferValueException e) {
            String nullBytesHex = WithFieldOption.extractNullBytesHex(options).orElseThrow(() -> new BufferException("Buffer value exception for field " + getName(options), e));
            byte[] nullBytes = StaticHelper.DECODE_HEX(nullBytesHex);
            Object value = e.getValue();
            if (value instanceof byte[]) {
                if (Arrays.equals(nullBytes, (byte[]) value)) {
                    return null;
                }
            }
            throw new BufferException("Buffer value exception for field " + getName(options), e);
        } catch (ParseAssertException e) {
            LOGGER.debug("Assertion doesn't match for field {}. Resetting read position to {}", getName(options), curPosInBits, e);
            dataReader.setPositionInBits(curPosInBits);
            return null;
        } catch (BufferUnderflowException e) {
            // The message ended before the field. Absent is the answer the caller can work with, and
            // it is only ever given for data that ran out - never for data that made no sense.
            LOGGER.debug("Not enough bytes for {}. Resetting read position to {}", getName(options), curPosInBits, e);
            dataReader.setPositionInBits(curPosInBits);
            return null;
        }
    }

}
