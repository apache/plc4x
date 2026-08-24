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

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.fields.data.reader.DataReader;
import org.apache.plc4x.java.spi.fields.exceptions.ParseAssertException;
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class FieldReaderPeek<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderPeek.class);

    public T readPeekField(DataReader<T> dataReader, int offset, WithOption... options) throws BufferException {
        int curPosInBits = dataReader.getPositionInBits();
        // A peek is undone whatever happens, so the nesting budget goes back with the position. A
        // peeked complex type that failed left its context open, since only a read that finishes
        // pops what it pushed.
        int curContextDepth = dataReader.getReadBuffer().getContextDepth();
        try {
            // TODO: implement offset. We either need to pass the readBuffer or add a instruction to the dataReader
            T field = dataReader.read(options);
            LOGGER.debug("done reading field {}. Value: {}", getName(options), field);
            return field;
        } catch (ParseAssertException e) {
            LOGGER.debug("Peeking failed for field {}. Resetting read position to {}", getName(options), curPosInBits, e);
            return null;
        } finally {
            dataReader.setPositionInBits(curPosInBits);
            dataReader.getReadBuffer().resetContextDepth(curContextDepth);
        }
    }
}
