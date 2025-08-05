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

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.spi.codegen.FieldCommons;
import org.apache.plc4x.java.spi.codegen.ThreadLocalHelper;
import org.apache.plc4x.java.spi.codegen.io.DataReader;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WithReaderArgs;
import org.apache.plc4x.java.spi.generation.WithReaderWriterArgs;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;

public class FieldReaderArray<T> implements FieldCommons {

    private static final Logger LOGGER = LoggerFactory.getLogger(FieldReaderArray.class);

    public List<T> readFieldCount(String logicalName, DataReader<T> dataReader, long count, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}. Count: {}", logicalName, count);
        if (count > Integer.MAX_VALUE) {
            throw new ParseException("Array count of " + count + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
        }
        if (count < 0) {
            return null;
        }
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        dataReader.pullContext(logicalName, readerArgs);
        int itemCount = Math.max(0, (int) count);
        List<T> result = new ArrayList<>(itemCount);
        for (int curItem = 0; curItem < itemCount; curItem++) {
            // Make some variables available that would be otherwise challenging to forward.
            ThreadLocalHelper.curItemThreadLocal.set(curItem);
            ThreadLocalHelper.lastItemThreadLocal.set(curItem == itemCount - 1);
            result.add(dataReader.read("", readerArgs));
        }
        dataReader.closeContext(logicalName, readerArgs);
        LOGGER.debug("done reading field {}", logicalName);
        return result;
    }

    public List<T> readFieldLength(String logicalName, DataReader<T> dataReader, int length, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}. Length: {}", logicalName, length);
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        dataReader.pullContext(logicalName, readerArgs);
        int startPos = dataReader.getPos();
        List<T> result = new ArrayList<>();
        int numberOfElements = 0;
        int stopPosition = startPos + length;
        LOGGER.debug("start reading at pos {} while < {}", startPos, stopPosition);
        while (dataReader.getPos() < stopPosition) {
            numberOfElements++;
            T element = dataReader.read("", readerArgs);
            LOGGER.debug("Read element[{}] {}", numberOfElements, element);
            result.add(element);
        }
        dataReader.closeContext(logicalName, readerArgs);
        LOGGER.debug("done reading field {}", logicalName);
        return result;
    }

    public List<T> readFieldTerminated(String logicalName, DataReader<T> dataReader, Supplier<Boolean> termination, WithReaderArgs... readerArgs) throws ParseException {
        LOGGER.debug("reading field {}", logicalName);
        // Ensure we have the render as list argument present
        readerArgs = ArrayUtils.add(readerArgs, WithReaderWriterArgs.WithRenderAsList(true));
        dataReader.pullContext(logicalName, readerArgs);
        List<T> result = new ArrayList<>();
        while (!termination.get()) {
            result.add(dataReader.read("", readerArgs));
        }
        dataReader.closeContext(logicalName, readerArgs);
        LOGGER.debug("done reading field {}", logicalName);
        return result;
    }

}
