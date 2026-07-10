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
import org.apache.plc4x.java.spi.fields.fields.FieldCommons;
import org.apache.plc4x.java.spi.fields.utils.ThreadLocalHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;

public class FieldReaderArray<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderArray.class);

    /**
     * Upper bound for the initial capacity of an array whose element count comes from the wire.
     * Prevents a bogus, attacker-controlled count from triggering an eager multi-GB allocation
     * before any element is actually read. The list still grows past this if that many elements
     * really are present on the wire.
     */
    private static final int MAX_INITIAL_CAPACITY = 1024;

    public List<T> readArrayFieldCount(DataReader<T> dataReader, long count, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}. Count: {}", getName(options), count);
        if (count > Integer.MAX_VALUE) {
            throw new BufferException("Array count of " + count + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
        }
        if (count < 0) {
            return null;
        }
        // Ensure we have the render as list argument present
        //readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        dataReader.pushContext(options);
        int itemCount = Math.max(0, (int) count);
        // Don't eagerly pre-size the backing array to an untrusted, wire-supplied count: a bogus
        // count (e.g. 0x7fffffff) would trigger a multi-GB allocation before a single element is
        // read (OOM / DoS). Cap the initial capacity to a modest bound and let the list grow as
        // elements are actually decoded - once the buffer is exhausted, dataReader.read() throws.
        List<T> result = new ArrayList<>(Math.min(itemCount, MAX_INITIAL_CAPACITY));
        for (int curItem = 0; curItem < itemCount; curItem++) {
            // Make some variables available that would be otherwise challenging to forward.
            ThreadLocalHelper.curItemThreadLocal.set(curItem);
            ThreadLocalHelper.lastItemThreadLocal.set(curItem == itemCount - 1);
            result.add(dataReader.read(options));
        }
        dataReader.popContext(options);
        LOGGER.debug("done reading field {}", getName(options));
        return result;
    }

    public List<T> readArrayFieldLength(DataReader<T> dataReader, int lengthInBytes, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}. Length in bytes: {}", getName(options), lengthInBytes);
        // Ensure we have the render as list argument present
        //readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        dataReader.pushContext(options);
        int startPos = dataReader.getPositionInBits();
        List<T> result = new ArrayList<>();
        int numberOfElements = 0;
        int stopPosition = startPos + (lengthInBytes * 8);
        LOGGER.debug("start reading at pos {} while < {}", startPos, stopPosition);
        while (dataReader.getPositionInBits() < stopPosition) {
            numberOfElements++;
            T element = dataReader.read(options);
            LOGGER.debug("Read element[{}] {}", numberOfElements, element);
            result.add(element);
        }
        dataReader.popContext(options);
        LOGGER.debug("done reading field {}", getName(options));
        return result;
    }

    public List<T> readArrayFieldTerminated(DataReader<T> dataReader, Supplier<Boolean> termination, WithOption... options) throws BufferException {
        LOGGER.debug("reading field {}", getName(options));
        // Ensure we have the render as list argument present
        //readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        dataReader.pushContext(options);
        List<T> result = new ArrayList<>();
        while (!termination.get()) {
            result.add(dataReader.read(options));
        }
        dataReader.popContext(options);
        LOGGER.debug("done reading field {}", getName(options));
        return result;
    }

}
