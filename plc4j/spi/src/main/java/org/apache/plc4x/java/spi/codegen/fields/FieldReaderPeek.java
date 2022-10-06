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

import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.ParseAssertException;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class FieldReaderPeek<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderPeek.class);

    public T readPeekField(String logicalName, DataReader<T> dataReader, int offset, WithReaderArgs... readerArgs) throws ParseException {
        int curPos = dataReader.getPos();
        try {
            // TODO: implement offset. We either need to pass the readBuffer or add a instruction to the dataReader
            T field = switchParseByteOrderIfNecessary(() -> dataReader.read(logicalName, readerArgs), dataReader, extractByteOrder(readerArgs).orElse(null));
            LOGGER.debug("done reading field {}. Value: {}", logicalName, field);
            return field;
        } catch (ParseAssertException e) {
            LOGGER.debug("Peeking failed for field {}. Resetting read position to {}", logicalName, curPos, e);
            return null;
        } finally {
            dataReader.setPos(curPos);
        }
    }
}
