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
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Arrays;
import java.util.Optional;
import java.util.stream.Stream;

public class FieldReaderOptional<T> implements FieldCommons, WithReaderWriterArgs {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderOptional.class);

    public T readOptionalField(String logicalName, DataReader<T> dataReader, boolean condition, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}", logicalName);
        if (!condition) {
            LOGGER.debug("Condition doesn't match for field {}", logicalName);
            return null;
        }

        // Check if a nullByteHex is set.
        // If it is, peek the equivalent number of bytes and compare.
        // If they match, return null.
        Optional<String> nullByteHexOptional = extractNullBytesHex(readerArgs);
        if(nullByteHexOptional.isPresent()) {
            String nullByteHex = nullByteHexOptional.get();
            try {
                byte[] nullBytes = Hex.decodeHex(nullByteHex);
                ReadBuffer readBuffer = dataReader.getReadBuffer();
                int pos = readBuffer.getPos();
                byte[] curBytes = readBuffer.readByteArray("logicalName", nullBytes.length, readerArgs);
                // Compare them, if they equal, return null, if not reset the position and try to read it again.
                if (Arrays.equals(nullBytes, curBytes)) {
                    // Abort with null
                    return null;
                } else {
                    readBuffer.reset(pos);
                }
            } catch (DecoderException e) {
                // Ignore.
            }
        }

        int curPos = dataReader.getPos();
        try {
            T field = switchParseByteOrderIfNecessary(() -> dataReader.read(logicalName, readerArgs), dataReader, extractByteOrder(readerArgs).orElse(null));
            LOGGER.debug("done reading field {}. Value: {}", logicalName, field);
            return field;
        } catch (ParseAssertException e) {
            LOGGER.debug("Assertion doesn't match for field {}. Resetting read position to {}", logicalName, curPos, e);
            dataReader.setPos(curPos);
            return null;
        } catch (ArrayIndexOutOfBoundsException e) {
            LOGGER.debug("Not enough bytes for {}. Resetting read position to {}", logicalName, curPos, e);
            dataReader.setPos(curPos);
            return null;
        }
    }

}
